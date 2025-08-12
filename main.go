package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
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

func main() {
    // 定义参数，短参数和长参数都定义
    bindShort := flag.String("b", "", "bind to this address (default all interfaces)")
    bindLong := flag.String("bind", "", "bind to this address (default all interfaces)")

    dirShort := flag.String("d", ".", "serve this directory (default current directory)")
    dirLong := flag.String("directory", ".", "serve this directory (default current directory)")

    protoShort := flag.String("p", "HTTP/1.0", "HTTP protocol version to use (HTTP/1.0 or HTTP/1.1)")
    protoLong := flag.String("protocol", "HTTP/1.0", "HTTP protocol version to use (HTTP/1.0 or HTTP/1.1)")

    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] [port]\n\n", os.Args[0])
        fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
        flag.PrintDefaults()
        fmt.Fprintf(flag.CommandLine.Output(), "\nPositional arguments:\n  port\t\tport number to listen on (default 8000)\n")
    }

    flag.Parse()

    bind := resolveParam("b", "bind", *bindShort, *bindLong, "")
    dir := resolveParam("d", "directory", *dirShort, *dirLong, ".")
    proto := resolveParam("p", "protocol", *protoShort, *protoLong, "HTTP/1.0")

    // 端口解析（位置参数）
    port := 8000
    args := flag.Args()
    if len(args) > 0 {
        _, err := fmt.Sscanf(args[0], "%d", &port)
        if err != nil {
            log.Fatalf("Invalid port number: %v", err)
        }
    }

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

    fs := http.FileServer(http.Dir(dir))

    handler := loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 强制修改请求协议版本，Go默认HTTP/1.1
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
