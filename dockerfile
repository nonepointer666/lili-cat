# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.24 AS builder

# 设置工作目录
WORKDIR /app

# 将本地代码复制到工作目录
COPY . .

# 构建 Go 应用
RUN go build -o lilicat

# 使用轻量级的 Alpine 镜像作为运行时环境
FROM alpine:latest

# 将构建好的可执行文件从构建阶段复制到运行时镜像
COPY --from=builder /app/lilicat /usr/local/bin/lilicat

# 暴露应用运行的端口
EXPOSE 8080

# 定义容器启动时运行的命令
CMD ["lilicat"]