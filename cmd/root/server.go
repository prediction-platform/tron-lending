package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/sunjiangjun/xlog"

	"lending-trx/internal/cronjob"
	"lending-trx/internal/db"
	"lending-trx/internal/webhook"
)

var (
	serverPort string
	serverCmd  = &cobra.Command{
		Use:   "server",
		Short: "启动完整的TRX委托服务 (HTTP API + 定时任务)",
		Long: `启动完整的TRX委托服务，包括：
- HTTP API服务：提供委托账户查询接口
- 定时任务：自动处理webhook数据和能量委托
- 数据库连接：PostgreSQL数据库操作`,
		Run: runServer,
	}
)

func init() {
	serverCmd.Flags().StringVarP(&serverPort, "port", "p", "8080", "HTTP服务端口")
}

func runServer(cmd *cobra.Command, args []string) {
	fmt.Println("🚀 启动TRX委托服务...")

	ctx := context.Background()
	pool, err := db.InitDB(ctx)
	if err != nil {
		log.Fatal("❌ 数据库初始化失败:", err)
	}
	defer pool.Close()

	LOG := xlog.NewXLogger().
		BuildOutType(xlog.FILE).
		BuildLevel(xlog.InfoLevel).
		BuildFormatter(xlog.FORMAT_JSON).
		BuildFile("logs/lending-trx.log", 24*time.Hour)

	// 启动定时任务
	cronjob.StartCron(ctx, pool, LOG)

	// 启动 gin HTTP 服务
	r := gin.Default()
	webhook.RegisterRoutes(r, ctx, pool, LOG)

	// 使用命令行参数或环境变量
	port := serverPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("✅ TRX委托服务已启动，监听端口: %s\n", port)
	fmt.Printf("📡 API地址: http://localhost:%s\n", port)
	fmt.Printf("📊 委托账户查询: http://localhost:%s/api/delegation-account\n", port)
	fmt.Printf("📝 日志文件: logs/lending-trx.log\n")

	r.Run(":" + port)
}
