package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sunjiangjun/xlog"

	"lending-trx/internal/cronjob"
	"lending-trx/internal/db"
	"lending-trx/internal/webhook"
)

func main() {
	// 加载环境变量文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未找到.env文件，使用系统环境变量")
	}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("✅ TRX委托服务已启动，监听端口: %s\n", port)
	fmt.Printf("📡 API地址: http://localhost:%s\n", port)
	fmt.Printf("📊 委托账户查询: http://localhost:%s/api/delegation-account\n", port)
	fmt.Printf("📝 日志文件: logs/lending-trx.log\n")

	r.Run(":" + port)
}
