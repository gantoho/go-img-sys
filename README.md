## 说明
快速启动一个图片服务器

```shell
$ git clone https://github.com/gantoho/go-img-sys.git
$ cd go-img-sys
```

Go启动
```shell
$ go mod tidy
$ go run main.go
```

运行在本机3128端口
```
http://localhost:3128/v1
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
