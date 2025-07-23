package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"lending-trx/internal/cronjob"
	"lending-trx/internal/db"
	"lending-trx/internal/webhook"

	"github.com/sunjiangjun/xlog"
)

func main() {
	ctx := context.Background()
	pool, err := db.InitDB(ctx)
	if err != nil {

		panic(err)
	}
	defer pool.Close()

	LOG := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildLevel(xlog.InfoLevel).BuildFormatter(xlog.FORMAT_JSON).BuildFile("logs/lending-trx.log", 24*time.Hour)

	// 启动定时任务
	cronjob.StartCron(ctx, pool, LOG)

	// 启动 gin HTTP 服务
	r := gin.Default()
	webhook.RegisterRoutes(r, ctx, pool, LOG)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("服务已启动，监听端口:", port)
	r.Run(":" + port)
}
