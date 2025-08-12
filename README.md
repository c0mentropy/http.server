# http.server



如题，虽然`python -m http.server`已经非常好用了，但很无奈，并不是所有机器上都有python，所以写个简单server传个文件。



```bash
./bin/linux_amd64/server -h
```

```bash
Usage: ./bin/linux_amd64/server [options] [port]

Options:
  -b string
        bind to this address (default all interfaces)
  -bind string
        bind to this address (default all interfaces)
  -d string
        serve this directory (default current directory) (default ".")
  -directory string
        serve this directory (default current directory) (default ".")
  -p string
        HTTP protocol version to use (HTTP/1.0 or HTTP/1.1) (default "HTTP/1.0")
  -protocol string
        HTTP protocol version to use (HTTP/1.0 or HTTP/1.1) (default "HTTP/1.0")

Positional arguments:
  port          port number to listen on (default 8000)
```

