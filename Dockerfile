# 使用官方 Go 镜像作为构建环境
FROM golang:1.21 AS builder

WORKDIR /app

# 拷贝 go.mod 和 go.sum 并下载依赖
COPY go.mod .
RUN go mod download

# 拷贝源代码
COPY . .

# 构建可执行文件
RUN go build -o app main.go

# 使用更小的基础镜像运行
FROM debian:bullseye-slim
WORKDIR /app

# 安装 ca-certificates 以支持 https
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# 拷贝编译好的二进制文件
COPY --from=builder /app/app .

# 设置环境变量（可选，实际运行时可覆盖）
ENV PG_DSN="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

# 启动应用
CMD ["./app"] 