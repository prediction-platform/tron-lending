package cronjob

import (
	"context"
	"fmt"
	"time"

	"lending-trx/internal/webhook"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"github.com/sunjiangjun/xlog"
)

// StartCron 启动定时任务，每分钟查询并打印相关数据，并返回结果数组
func StartCron(ctx context.Context, pool *pgxpool.Pool, log *xlog.XLog) {
	l := log.WithField("module", "cronjob")
	c := cron.New()
	_, _ = c.AddFunc("@every 30s", func() {
		// 查询 status=0
		var result0 []webhook.WebhookDataModel
		rows0, err := pool.Query(ctx, "SELECT id, block_height, tx_hash, from_address, to_address, value, block_time, create_time, update_time, expire_time, status FROM webhook_data WHERE status=0")
		if err != nil {
			fmt.Printf("查询 status=0 失败: %v\n", err)
		} else {

			for rows0.Next() {
				var d webhook.WebhookDataModel
				err := rows0.Scan(&d.ID, &d.BlockHeight, &d.TxHash, &d.FromAddress, &d.ToAddress, &d.Value, &d.BlockTime, &d.CreateTime, &d.UpdateTime, &d.ExpireTime, &d.Status)
				if err != nil {
					fmt.Printf("status=0 scan error: %v\n", err)
					return
				}
				result0 = append(result0, d)
			}
			fmt.Printf("[定时任务] status=0 数量: %d\n", len(result0))
		}
		if rows0 != nil {
			rows0.Close()
		}

		//todo
		l.Info("InitList", len(result0))
		if len(result0) > 0 {
			// 执行能量委托
		}

		// 查询 status=2 且已失效
		now := time.Now().UnixMilli()
		var result2 []webhook.WebhookDataModel
		rows2, err := pool.Query(ctx, "SELECT id, block_height, tx_hash, from_address, to_address, value, block_time, create_time, update_time, expire_time, status FROM webhook_data WHERE status=2 AND expire_time < $1", now)
		if err != nil {
			fmt.Printf("查询 status=2 且已失效失败: %v\n", err)
		} else {
			for rows2.Next() {
				var d webhook.WebhookDataModel
				err := rows2.Scan(&d.ID, &d.BlockHeight, &d.TxHash, &d.FromAddress, &d.ToAddress, &d.Value, &d.BlockTime, &d.CreateTime, &d.UpdateTime, &d.ExpireTime, &d.Status)
				if err != nil {
					fmt.Printf("status=2 scan error: %v\n", err)
					return
				}
				result2 = append(result2, d)
			}
			fmt.Printf("[定时任务] status=2 且已失效数量: %d\n", len(result2))
		}
		if rows2 != nil {
			rows2.Close()
		}

		l.Info("DoneList", len(result2))
		if len(result2) > 0 {
			//todo 取消能量委托

			// update status
		}
	})
	go c.Run()
}
