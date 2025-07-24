# 使用官方 Go 镜像作为构建环境
FROM golang:1.23 AS builder

WORKDIR /app

# 拷贝 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码
COPY . .

# 静态编译，构建新的命令行工具
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lending-trx cmd/root/*.go

# 使用 alpine 镜像，更轻量
FROM alpine:latest
WORKDIR /app

# 安装 ca-certificates 以支持 https
RUN apk --no-cache add ca-certificates

# 拷贝编译好的二进制文件
COPY --from=builder /app/lending-trx .

# 创建日志目录
RUN mkdir -p /app/logs

# 设置默认环境变量
ENV DATABASE_URL="postgresql://postgres:password@localhost:5432/lending_trx?sslmode=disable"
ENV PORT="8080"

# 默认启动server命令
CMD ["./lending-trx", "server"]