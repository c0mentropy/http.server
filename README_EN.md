```markdown
# FILE CLI

`filecli` is a simple file transfer tool developed in Go. It can function both as a file server and as a client for uploading and downloading files. It provides both command-line and web interfaces for convenient use in different scenarios.

---

## Features

- **Cross-platform**: While `python -m http.server` is quite useful, it's unfortunately not available on all machines. This simple file service is designed for easy file transfers.
- **Multi-mode operation**: Supports server mode, client upload mode, and client download mode
- **Web interface**: Provides a clean web upload form for easy file uploads via browser
- **Command-line operation**: Supports file uploads and downloads through command line
- **Flexible configuration**: Allows specifying bind address, port, service directory, and HTTP protocol version
- **Logging**: Records access logs for easy tracking of file transfers

---

## Installation

### Compile from source

1. Ensure Go 1.16 or higher is installed

2. Clone the repository:

    ```bash
    git clone https://github.com/c0mentropy/filecli
    cd filecli
    ```

3. Compile:

    ```bash
    go build -o filecli
    ```

    You can also check `make help` for compilation options

4. Add the generated `filecli` executable to your system PATH (optional)

---

## Usage

### Server Mode (default)

Start the file server: `filecli [options] [port]`

#### Server Options

| Short | Long        | Description                  | Default value   |
|-------|-------------|------------------------------|-----------------|
| -b    | --bind      | Address to bind to           | All interfaces  |
| -d    | --directory | Directory to serve files from | Current directory (".") |
| -p    | --protocol  | HTTP protocol version to use | HTTP/1.0       |

#### Examples

1. Start server with default configuration (port 8000, current directory):

    ```bash
    filecli
    ```

2. Start server on port 8080, serving the `/data` directory:

    ```bash
    filecli -d /data 8080
    ```

3. Bind to localhost, using HTTP/1.1 protocol:

    ```bash
    filecli -b 127.0.0.1 -p HTTP/1.1 9000
    ```

After starting, you can access:

- File browsing: `http://server-address:port`
- Upload form: `http://server-address:port/upload-form`

### Client Upload Mode

Upload files to server via command line: `filecli upload [options]`

#### Upload Options

| Short | Long      | Description                               | Required |
|-------|-----------|-------------------------------------------|----------|
| -s    | --server  | Server address                            | Yes      |
| -p    | --port    | Server port                               | Yes      |
| -l    | --lfile   | Local file path                           | Yes      |
| -r    | --rfile   | Remote filename to save as (optional)     | No       |

#### Examples

1. Upload local file `test.txt` to server `192.168.1.100` on port 8000:

    ```bash
    filecli upload -s 192.168.1.100 -p 8000 -l test.txt
    ```

2. Upload local file `docs.pdf` to server, specifying remote filename as `document.pdf`:

    ```bash
    filecli upload -s example.com -p 8080 -l docs.pdf -r document.pdf
    ```

### Client Download Mode

Download files from server via command line: `filecli download [options]`

#### Download Options

| Short | Long      | Description                               | Required |
|-------|-----------|-------------------------------------------|----------|
| -s    | --server  | Server address                            | Yes      |
| -p    | --port    | Server port                               | Yes      |
| -r    | --rfile   | Remote file path                          | Yes      |
| -l    | --lfile   | Local save path (optional)                | No       |

#### Examples

1. Download `data.zip` from server `192.168.1.100` to local machine:

    ```bash
    filecli download -s 192.168.1.100 -p 8000 -r data.zip
    ```

2. Download `report.csv` from server and save to `./downloads/report.csv` locally:

    ```bash
    filecli download -s example.com -p 8080 -r report.csv -l ./downloads/report.csv
    ```

---

## Version Information

`filecli --version`
filecli version v0.1.1
---

## Help Information

- View general help:
  - `filecli --help`

- View help for specific subcommands:
  - `filecli upload --help`
  - `filecli download --help`

---

## Contact

- Author: Ckyan Comentropy
- Email: comentropy@foxmail.com
- GitHub: https://github.com/c0mentropy/filecli

---

## Contributing

Contributions and bug reports are welcome!
Please submit Issues or Pull Requests on [GitHub](https://github.com/c0mentropy/filecli)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details