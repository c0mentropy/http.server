# FILE CLI

`filecli `是一个基于 Go 语言开发的简易文件传输工具，既可作为文件服务器使用，也可作为客户端进行文件的上传和下载操作。它提供了命令行界面和网页界面两种交互方式，方便用户在不同场景下使用。

---

## 功能特点

- **跨平台**：`python -m http.server`已经很好用了，但很遗憾并不是所有机器上都有python，所以写个简单的文件服务，传个文件用
- **多模式运行**：支持服务器模式、客户端上传模式和客户端下载模式
- **网页界面**：提供简洁的网页上传表单，方便通过浏览器上传文件
- **命令行操作**：支持通过命令行进行文件的上传和下载
- **灵活配置**：可指定绑定地址、端口、服务目录和 HTTP 协议版本
- **日志记录**：记录访问日志，便于追踪文件传输情况

---

## 安装

### 从源码编译

1. 确保已安装 Go 1.16 或更高版本

2. 克隆代码仓库：

    ```bash
    git clone https://github.com/c0mentropy/filecli
    cd filecli
    ```

3. 编译：

    ```bash
    go build -o filecli
    ```

    也可查看`make help`编译

4. 将生成的 `filecli` 可执行文件添加到系统 PATH 中（可选）

---

## 使用方法

### 服务器模式（默认）

启动文件服务器：`filecli [options] [port]`

#### 服务器选项

| 短选项 | 长选项      | 说明                 | 默认值         |
| ------ | ----------- | -------------------- | -------------- |
| -b     | --bind      | 绑定的地址           | 所有网络接口   |
| -d     | --directory | 提供文件服务的目录   | 当前目录 (".") |
| -p     | --protocol  | 使用的 HTTP 协议版本 | HTTP/1.0       |

#### 示例

1. 使用默认配置启动服务器（端口 8000，当前目录）：

    ```bash
    filecli
    ```

2. 在端口 8080 上启动服务器，服务 `/data` 目录：

    ```bash
    filecli -d /data 8080
    ```

3. 绑定到本地回环地址，使用 HTTP/1.1 协议：

    ```bash
    filecli -b 127.0.0.1 -p HTTP/1.1 9000
    ```

启动后，可以通过以下方式访问：

- 文件浏览：`http://服务器地址:端口`
- 上传表单：`http://服务器地址:端口/upload-form`

### 客户端上传模式

通过命令行上传文件到服务器：`filecli upload [options]`

#### 上传选项

| 短选项 | 长选项   | 说明                     | 是否必须 |
| ------ | -------- | ------------------------ | -------- |
| -s     | --server | 服务器地址               | 是       |
| -p     | --port   | 服务器端口               | 是       |
| -l     | --lfile  | 本地文件路径             | 是       |
| -r     | --rfile  | 远程保存的文件名（可选） | 否       |

#### 示例

1. 上传本地文件 `test.txt` 到服务器 `192.168.1.100` 的 8000 端口：

    ```bash
    filecli upload -s 192.168.1.100 -p 8000 -l test.txt
    ```

2. 上传本地文件 `docs.pdf` 到服务器，并指定远程文件名为 `document.pdf`：

    ```bash
    filecli upload -s example.com -p 8080 -l docs.pdf -r document.pdf
    ```

### 客户端下载模式

通过命令行从服务器下载文件：`filecli download [options]`

#### 下载选项

| 短选项 | 长选项   | 说明                 | 是否必须 |
| ------ | -------- | -------------------- | -------- |
| -s     | --server | 服务器地址           | 是       |
| -p     | --port   | 服务器端口           | 是       |
| -r     | --rfile  | 远程文件路径         | 是       |
| -l     | --lfile  | 本地保存路径（可选） | 否       |

#### 示例

1. 从服务器 `192.168.1.100` 下载 `data.zip` 到本地：

    ```bash
    filecli download -s 192.168.1.100 -p 8000 -r data.zip
    ```

2. 从服务器下载 `report.csv` 并保存到本地 `./downloads/report.csv`：

    ```bash
    filecli download -s example.com -p 8080 -r report.csv -l ./downloads/report.csv
    ```

---

## 版本信息

`filecli --version`

```bash
filecli version v0.1.1
```

---

## 帮助信息

-   查看整体帮助：
    -   `filecli --help`

-   查看特定子命令的帮助：

    -   `filecli upload --help`

    -   `filecli download --help`

---

## 联系方式

- 作者：Ckyan Comentropy
- 邮箱：comentropy@foxmail.com
- GitHub：https://github.com/c0mentropy/filecli

---

## 贡献

欢迎贡献代码和报告 Bug！
欢迎在 [GitHub](https://github.com/c0mentropy/filecli) 上提交 Issue 或 Pull Request

---

## 许可证

该项目采用 MIT 许可证 - 有关详细信息，请参阅 [LICENSE](LICENSE) 文件