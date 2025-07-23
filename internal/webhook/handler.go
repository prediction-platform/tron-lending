package webhook

import (
	"context"
	"net/http"

	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunjiangjun/xlog"
)

type WebhookData struct {
	BlockHeight int64  `json:"block_height"`
	TxHash      string `json:"tx_hash"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Value       string `json:"value"`
	BlockTime   int64  `json:"block_time"`
	ExpireTime  int64  `json:"expire_time"`
	Status      int16  `json:"status"`
}

// RegisterRoutes 注册 webhook 路由
func RegisterRoutes(r *gin.Engine, ctx context.Context, pool *pgxpool.Pool, log *xlog.XLog) {
	l := log.WithField("module", "webhook")
	r.POST("/webhook", func(c *gin.Context) {
		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			l.Error("读取请求体失败", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
			return
		}
		// 打印请求体内容
		l.Info("webhook 请求体", string(body))

		// 打印请求头信息
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			headers[k] = strings.Join(v, ",")
		}
		l.Info("webhook 请求头", headers)

		// 重新设置请求体供 ShouldBindJSON 使用
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var data WebhookData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		var id int64
		err = pool.QueryRow(ctx, `INSERT INTO webhook_data 
			(block_height, tx_hash, from_address, to_address, value, block_time, expire_time, status, create_time, update_time)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW()) RETURNING id`,
			data.BlockHeight, data.TxHash, data.FromAddress, data.ToAddress, data.Value, data.BlockTime, data.ExpireTime, data.Status).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": id})
	})
}
