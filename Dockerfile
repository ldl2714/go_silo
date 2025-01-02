# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.17-alpine

# 设置容器内的当前工作目录
WORKDIR /app

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载所有依赖项。如果 go.mod 和 go.sum 文件没有更改，依赖项将被缓存
RUN go mod download

# 将源代码从当前目录复制到容器内的工作目录
COPY . .

# 构建 Go 应用程序
RUN go build -o main .

# 向外界暴露端口 8080
EXPOSE 8080

# 运行可执行文件的命令
CMD ["./main"]