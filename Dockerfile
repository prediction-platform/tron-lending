# 使用官方 Go 镜像作为构建环境
FROM golang:1.23 AS builder

WORKDIR /app

# 拷贝 go.mod 和 go.sum 并下载依赖
COPY go.mod .
RUN go mod download

# 拷贝源代码
COPY . .

# 静态编译，避免 GLIBC 依赖问题
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

# 使用 alpine 镜像，更轻量
FROM alpine:latest
WORKDIR /app

# 安装 ca-certificates 以支持 https
RUN apk --no-cache add ca-certificates

# 拷贝编译好的二进制文件
COPY --from=builder /app/app .

# 设置环境变量（可选，实际运行时可覆盖）
ENV PG_DSN="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

# 启动应用
CMD ["./app"]