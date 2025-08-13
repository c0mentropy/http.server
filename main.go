package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mime/multipart"
)

func resolveParam(shortName, longName string, shortVal, longVal string, defaultVal string) string {
	shortSet := shortVal != defaultVal
	longSet := longVal != defaultVal

	if shortSet && longSet {
		if shortVal != longVal {
			log.Fatalf("Conflicting values for -%s and --%s: %q vs %q", shortName, longName, shortVal, longVal)
		}
		return shortVal
	}
	if shortSet {
		return shortVal
	}
	if longSet {
		return longVal
	}
	return defaultVal
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s - - [%s] \"%s %s %s\" %v",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			time.Since(start),
		)
	})
}

// 处理文件上传
func uploadHandler(root string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析multipart表单
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
			return
		}

		// 获取文件句柄
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 获取目标文件名
		destName := r.FormValue("dest")
		if destName == "" {
			destName = handler.Filename
		}

		// 构建目标路径
		destPath := filepath.Join(root, destName)
		
		// 创建目标文件
		dst, err := os.Create(destPath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// 复制文件内容
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully: %s", destName)
	}
}

// 提供文件上传表单
func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File Upload</title>
        <style>
            body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
            .container { background-color: #f5f5f5; padding: 20px; border-radius: 8px; }
            h1 { color: #333; }
            .form-group { margin-bottom: 15px; }
            label { display: block; margin-bottom: 5px; font-weight: bold; }
            input[type="file"] { padding: 8px; }
            input[type="text"] { width: 100%; padding: 8px; box-sizing: border-box; }
            button { background-color: #4CAF50; color: white; padding: 10px 15px; border: none; border-radius: 4px; cursor: pointer; }
            button:hover { background-color: #45a049; }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>File Upload</h1>
            <form action="/upload" method="post" enctype="multipart/form-data">
                <div class="form-group">
                    <label for="file">Select File:</label>
                    <input type="file" id="file" name="file" required>
                </div>
                <div class="form-group">
                    <label for="dest">Save As (optional):</label>
                    <input type="text" id="dest" name="dest" placeholder="Enter filename to save on server">
                </div>
                <button type="submit">Upload File</button>
            </form>
        </div>
    </body>
    </html>
    `

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, html)
}

func startServer(bind string, port int, dir string, proto string) {
	// 目录校验
	fi, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("Error accessing directory %q: %v", dir, err)
	}
	if !fi.IsDir() {
		log.Fatalf("%q is not a directory", dir)
	}

	// 协议版本校验
	if proto != "HTTP/1.0" && proto != "HTTP/1.1" {
		log.Fatalf("Unsupported protocol version: %s", proto)
	}

	// 创建文件服务器
	fs := http.FileServer(http.Dir(dir))

	// 自定义处理器，支持文件上传和协议版本控制
	handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理上传表单页面请求
		if r.URL.Path == "/upload-form" {
			uploadFormHandler(w, r)
			return
		}
		// 处理上传请求
		if r.URL.Path == "/upload" {
			uploadHandler(dir)(w, r)
			return
		}

		// 强制修改请求协议版本
		if proto == "HTTP/1.0" {
			r.Proto = "HTTP/1.0"
			r.ProtoMajor = 1
			r.ProtoMinor = 0
		} else {
			r.Proto = "HTTP/1.1"
			r.ProtoMajor = 1
			r.ProtoMinor = 1
		}
		fs.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf(":%d", port)
	if bind != "" {
		addr = fmt.Sprintf("%s:%d", bind, port)
	}

	log.Printf("Serving directory %q on %s with protocol %s\n", dir, addr, proto)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// 客户端下载文件
func clientDownload(serverAddr, sourcePath, destPath string) error {
	if !strings.HasPrefix(sourcePath, "/") {
		sourcePath = "/" + sourcePath
	}
	
	url := fmt.Sprintf("http://%s%s", serverAddr, sourcePath)
	
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Unable to connect to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returns error: %s", resp.Status)
	}

	// 如果未指定目标路径，使用源文件名并保存到当前目录
	if destPath == "" {
		parts := strings.Split(sourcePath, "/")
		destPath = parts[len(parts)-1]
	}

	// 创建目标文件
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %v", err)
	}
	defer out.Close()

	// 复制内容
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Download file failed: %v", err)
	}

	fmt.Printf("The file has been successfully downloaded to: %s\n", destPath)
	return nil
}

// 客户端上传文件
func clientUpload(serverAddr, sourcePath, destPath string) error {
	// 打开源文件
	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Unable to open file: %v", err)
	}
	defer file.Close()

	// 创建multipart表单
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	
	// 创建表单文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(sourcePath))
	if err != nil {
		return fmt.Errorf("Failed to create form: %v", err)
	}
	
	// 将文件内容复制到表单
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("Failed to read file contents: %v", err)
	}
	
	// 添加目标路径字段
	if destPath != "" {
		if err := writer.WriteField("dest", destPath); err != nil {
			return fmt.Errorf("Failed to set target path: %v", err)
		}
	}
	
	// 关闭multipart writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("Failed to close form: %v", err)
	}
	
	// 发送POST请求
	url := fmt.Sprintf("http://%s/upload", serverAddr)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("Create request failed: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Sending request failed: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Server returns error: %s - %s", resp.Status, string(body))
	}
	
	fmt.Println("File uploaded successfully")
	return nil
}

var filename = filepath.Base(os.Args[0])

func printUsage() {
	fmt.Printf("Usage: %s [options] [port]\n", filename)
	fmt.Printf("       %s [subcommand] [options]\n", filename)
	fmt.Println("\nDefault: Start the file server")
	fmt.Println("\nDownload route: /\nUpload   route: /upload-form\n")
	printServerUsage()
	fmt.Println("\nsubcommand:")
	fmt.Println("  upload    Upload files to the server")
	fmt.Println("  download  Download files from server")
	fmt.Println("\nRun 'filecli upload -h' or 'filecli download -h' to get help on the subcommands.")
}

func printServerUsage() {
    // 使用标准输出打印帮助信息
    w := os.Stdout
    fmt.Fprintf(w, "Usage: %s [options] [port]\n\n", filename)
    fmt.Fprintf(w, "Options:\n")
    
    fmt.Fprintf(w, "  -b string\n")
    fmt.Fprintf(w, "        bind to this address (default all interfaces)\n")
    fmt.Fprintf(w, "  -bind string\n")
    fmt.Fprintf(w, "        bind to this address (default all interfaces)\n")
    fmt.Fprintf(w, "  -d string\n")
    fmt.Fprintf(w, "        serve this directory (default current directory) (default \".\")\n")
    fmt.Fprintf(w, "  -directory string\n")
    fmt.Fprintf(w, "        serve this directory (default current directory) (default \".\")\n")
    fmt.Fprintf(w, "  -p string\n")
    fmt.Fprintf(w, "        HTTP protocol version to use (HTTP/1.0 or HTTP/1.1) (default \"HTTP/1.0\")\n")
    fmt.Fprintf(w, "  -protocol string\n")
    fmt.Fprintf(w, "        HTTP protocol version to use (HTTP/1.0 or HTTP/1.1) (default \"HTTP/1.0\")\n")
    
    fmt.Fprintf(w, "\nPositional arguments:\n  port\t\tport number to listen on (default 8000)\n")
}

func printUploadUsage() {
	fmt.Printf("Usage: %s upload [options]\n\n", filename)
	fmt.Println("Options:")
	fmt.Println("  -s, --server string   Server Address (required)")
	fmt.Println("  -p, --port int        Server Port (required)")
	fmt.Println("  -l, --lfile string    Local file path (required)")
	fmt.Println("  -r, --rfile string    Remote save file name (optional, defaults to the same as the local file name)")
}

func printDownloadUsage() {
	fmt.Printf("Usage: %s download [options]\n\n", filename)
	fmt.Println("Options:")
	fmt.Println("  -s, --server string   Server Address (required)")
	fmt.Println("  -p, --port int        Server Port (required)")
	fmt.Println("  -r, --rfile string    Remote file path (required)")
	fmt.Println("  -l, --lfile string    Local save path (optional, defaults to the same as the remote file name)")
}

const version = "v0.1.1"
const author = "Ckyan Comentropy"
const email = "comentropy@foxmail.com"
const github = "https://github.com/c0mentropy/filecli"

func printVersion() {
	fmt.Printf("filecli version %s\n", version)
	fmt.Printf("Author: %s\n", author)
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("GitHub: %s\n", github)
}

func main() {
	// 检查是否有子命令
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "upload":
			handleUploadCommand(os.Args[2:])
			return
		case "download":
			handleDownloadCommand(os.Args[2:])
			return
		case "-h", "--help":
			printUsage()
			return
		case "-V", "--version":
			printVersion()
			return
		}
	}

	// 默认启动服务器模式
	handleServerCommand(os.Args[1:])
}

func handleServerCommand(args []string) {

	// 创建新的flag集，处理服务器参数
	serverFlagSet := flag.NewFlagSet("server", flag.ExitOnError)
	
	// 定义服务器参数
	bindShort := serverFlagSet.String("b", "", "bind to this address (default all interfaces)")
	bindLong := serverFlagSet.String("bind", "", "bind to this address (default all interfaces)")

	dirShort := serverFlagSet.String("d", ".", "serve this directory (default current directory)")
	dirLong := serverFlagSet.String("directory", ".", "serve this directory (default current directory)")

	protoShort := serverFlagSet.String("p", "HTTP/1.0", "HTTP protocol version to use (HTTP/1.0 or HTTP/1.1)")
	protoLong := serverFlagSet.String("protocol", "HTTP/1.0", "HTTP protocol version to use (HTTP/1.0 or HTTP/1.1)")

	serverFlagSet.Usage = printServerUsage

	// 解析参数
	serverFlagSet.Parse(args)

	bind := resolveParam("b", "bind", *bindShort, *bindLong, "")
	dir := resolveParam("d", "directory", *dirShort, *dirLong, ".")
	proto := resolveParam("p", "protocol", *protoShort, *protoLong, "HTTP/1.0")

	// 端口解析（位置参数）
	port := 8000
	argsAfterFlags := serverFlagSet.Args()
	if len(argsAfterFlags) > 0 {
		_, err := fmt.Sscanf(argsAfterFlags[0], "%d", &port)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
		}
	}

	startServer(bind, port, dir, proto)
}

func handleUploadCommand(args []string) {
	// 创建上传命令的flag集
	uploadFlagSet := flag.NewFlagSet("upload", flag.ExitOnError)
	
	// 定义上传参数
	serverShort := uploadFlagSet.String("s", "", "")
	serverLong := uploadFlagSet.String("server", "", "")
	
	portShort := uploadFlagSet.Int("p", 0, "")
	portLong := uploadFlagSet.Int("port", 0, "")
	
	lfileShort := uploadFlagSet.String("l", "", "")
	lfileLong := uploadFlagSet.String("lfile", "", "")
	
	rfileShort := uploadFlagSet.String("r", "", "")
	rfileLong := uploadFlagSet.String("rfile", "", "")
	
	uploadFlagSet.Usage = printUploadUsage
	
	// 解析参数
	if err := uploadFlagSet.Parse(args); err != nil {
		log.Fatalf("Parameter parsing error: %v", err)
	}
	
	// 解析参数值
	server := resolveParam("s", "server", *serverShort, *serverLong, "")
	port := *portShort
	if port == 0 {
		port = *portLong
	}
	lfile := resolveParam("l", "lfile", *lfileShort, *lfileLong, "")
	rfile := resolveParam("r", "rfile", *rfileShort, *rfileLong, "")
	
	// 验证必填参数
	if server == "" {
		log.Fatal("Please specify the server address (-s or --server)")
	}
	if port == 0 {
		log.Fatal("Please specify the server port (-p or --port)")
	}
	if lfile == "" {
		log.Fatal("Please specify a local file path (-l or --lfile)")
	}
	
	// 构建服务器地址
	serverAddr := fmt.Sprintf("%s:%d", server, port)
	
	// 执行上传
	if err := clientUpload(serverAddr, lfile, rfile); err != nil {
		log.Fatalf("Upload failed: %v", err)
	}
}

func handleDownloadCommand(args []string) {
	// 创建下载命令的flag集
	downloadFlagSet := flag.NewFlagSet("download", flag.ExitOnError)
	
	// 定义下载参数
	serverShort := downloadFlagSet.String("s", "", "")
	serverLong := downloadFlagSet.String("server", "", "")
	
	portShort := downloadFlagSet.Int("p", 0, "")
	portLong := downloadFlagSet.Int("port", 0, "")
	
	rfileShort := downloadFlagSet.String("r", "", "")
	rfileLong := downloadFlagSet.String("rfile", "", "")
	
	lfileShort := downloadFlagSet.String("l", "", "")
	lfileLong := downloadFlagSet.String("lfile", "", "")
	
	downloadFlagSet.Usage = printDownloadUsage
	
	// 解析参数
	if err := downloadFlagSet.Parse(args); err != nil {
		log.Fatalf("Parameter parsing error: %v", err)
	}
	
	// 解析参数值
	server := resolveParam("s", "server", *serverShort, *serverLong, "")
	port := *portShort
	if port == 0 {
		port = *portLong
	}
	rfile := resolveParam("r", "rfile", *rfileShort, *rfileLong, "")
	lfile := resolveParam("l", "lfile", *lfileShort, *lfileLong, "")
	
	// 验证必填参数
	if server == "" {
		log.Fatal("Please specify the server address (-s or --server)")
	}
	if port == 0 {
		log.Fatal("Please specify the server port (-p or --port)")
	}
	if rfile == "" {
		log.Fatal("Please specify remote file path (-r or --rfile)")
	}
	
	// 构建服务器地址
	serverAddr := fmt.Sprintf("%s:%d", server, port)
	
	// 执行下载
	if err := clientDownload(serverAddr, rfile, lfile); err != nil {
		log.Fatalf("Download failed: %v", err)
	}
}
