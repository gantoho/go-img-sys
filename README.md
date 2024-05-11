## 说明
快速启动一个图片服务器

```shell
$ go clone https://github.com/gantoho/go-img-sys.git
$ go mod tidy
$ go run main.go
```

打包二进制文件
linux
```shell
$ go build -o main main.go
```

windows
```shell
$ go build -o main.exe main.go
```

windows build linux
```shell
$ $env:GOOS="linux"
$ go build -o main main.go
```