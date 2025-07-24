package db

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WebhookDataModel 用于表示 webhook_data 表结构
// 注意：Value 用 string 以兼容大数，时间戳用 int64
// create_time, update_time 用 string（或 time.Time，视序列化需求）
type WebhookDataModel struct {
	ID           int64  `json:"id"`             // 主键唯一ID
	BlockHeight  int64  `json:"block_height"`   // 区块高度
	TxHash       string `json:"tx_hash"`        // 交易哈希
	FromAddress  string `json:"from_address"`   // 发送方地址
	ToAddress    string `json:"to_address"`     // 接收方地址
	Value        string `json:"value"`          // 交易金额（大整数，字符串存储）
	BlockTime    int64  `json:"block_time"`     // 区块时间（毫秒时间戳）
	CreateTime   string `json:"create_time"`    // 创建时间
	UpdateTime   string `json:"update_time"`    // 更新时间
	ExpireTime   int64  `json:"expire_time"`    // 有效期（毫秒时间戳）
	Status       int16  `json:"status"`         // 状态（0:初始化，1:执行中，2:已授权，3:已回收）
	OriginalTxID string `json:"original_tx_id"` // 原始委托交易ID
}

const createWebhookTableSQL = `
CREATE TABLE IF NOT EXISTS webhook_data (
  id SERIAL PRIMARY KEY,
  block_height BIGINT,
  tx_hash VARCHAR(128),
  from_address VARCHAR(128),
  to_address VARCHAR(128),
  value NUMERIC(36,0),
  block_time BIGINT,
  create_time TIMESTAMP NOT NULL DEFAULT NOW(),
  update_time TIMESTAMP NOT NULL DEFAULT NOW(),
  expire_time BIGINT,
  status SMALLINT,
  original_tx_id VARCHAR(255) UNIQUE
);`

const createLogTableSQL = `
CREATE TABLE IF NOT EXISTS logs (
	id SERIAL PRIMARY KEY,
	message TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);`

// InitDB 连接数据库并初始化表，返回连接池
func InitDB(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:password@localhost:5432/lending_trx?sslmode=disable"
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("无法连接数据库: %w", err)
	}
	// 初始化表
	if _, err := pool.Exec(ctx, createWebhookTableSQL); err != nil {
		return nil, fmt.Errorf("创建 webhook_data 表失败: %w", err)
	}
	if _, err := pool.Exec(ctx, createLogTableSQL); err != nil {
		return nil, fmt.Errorf("创建 logs 表失败: %w", err)
	}
	return pool, nil
}

// UpdateWebhookStatus 批量更新指定 id 的 status
func UpdateWebhookStatus(ctx context.Context, pool *pgxpool.Pool, ids []int64, status int16) error {
	if len(ids) == 0 {
		return nil
	}
	var params []interface{}
	for _, id := range ids {
		params = append(params, id)
	}
	inClause := make([]string, len(ids))
	for i := range ids {
		inClause[i] = fmt.Sprintf("$%d", i+2)
	}
	query := fmt.Sprintf("UPDATE webhook_data SET status=$1, update_time=NOW() WHERE id IN (%s)", strings.Join(inClause, ","))
	params = append([]interface{}{status}, params...)
	_, err := pool.Exec(ctx, query, params...)
	return err
}

// UpdateWebhookStatusAndExpireTime 批量更新指定 id 的 status 和 expire_time
func UpdateWebhookStatusAndExpireTime(ctx context.Context, pool *pgxpool.Pool, ids []int64, status int16, expireTime int64) error {
	if len(ids) == 0 {
		return nil
	}
	var params []interface{}
	for _, id := range ids {
		params = append(params, id)
	}
	inClause := make([]string, len(ids))
	for i := range ids {
		inClause[i] = fmt.Sprintf("$%d", i+3)
	}
	query := fmt.Sprintf("UPDATE webhook_data SET status=$1, expire_time=$2, update_time=NOW() WHERE id IN (%s)", strings.Join(inClause, ","))
	params = append([]interface{}{status, expireTime}, params...)
	_, err := pool.Exec(ctx, query, params...)
	return err
}

// BatchInsertWebhookData 批量插入 webhook_data 记录
func BatchInsertWebhookData(ctx context.Context, pool *pgxpool.Pool, data []*WebhookDataModel) error {
	if len(data) == 0 {
		return nil
	}
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*9)
	for i, d := range data {
		idx := i * 9
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8, idx+9))
		valueArgs = append(valueArgs,
			d.BlockHeight, d.TxHash, d.FromAddress, d.ToAddress, d.Value, d.BlockTime, d.ExpireTime, d.Status, time.Now())
	}
	query := "INSERT INTO webhook_data (block_height, tx_hash, from_address, to_address, value, block_time, expire_time, status, create_time) VALUES " + strings.Join(valueStrings, ",")
	_, err := pool.Exec(ctx, query, valueArgs...)
	return err
}

// QueryPendingWebhookData 查询待处理的数据 (status=0)
func QueryPendingWebhookData(ctx context.Context, pool *pgxpool.Pool) ([]*WebhookDataModel, error) {
	query := `
		SELECT id, block_height, tx_hash, from_address, to_address, value, 
		       block_time, create_time, update_time, expire_time, status, original_tx_id 
		FROM webhook_data 
		WHERE status=0
		ORDER BY create_time ASC
	`

	return queryWebhookData(ctx, pool, query)
}

// QueryExpiredWebhookData 查询已过期且已授权的数据 (status=2)
func QueryExpiredWebhookData(ctx context.Context, pool *pgxpool.Pool) ([]*WebhookDataModel, error) {
	now := time.Now().UnixMilli()
	query := `
		SELECT id, block_height, tx_hash, from_address, to_address, value, 
		       block_time, create_time, update_time, expire_time, status, original_tx_id 
		FROM webhook_data 
		WHERE status=2 AND expire_time < $1
		ORDER BY expire_time ASC
	`

	return queryWebhookDataWithParams(ctx, pool, query, now)
}

// UpdateWebhookStatusByID 更新单个记录的status
func UpdateWebhookStatusByID(ctx context.Context, pool *pgxpool.Pool, id int64, status int16) error {
	query := `
		UPDATE webhook_data 
		SET status = $1, update_time = NOW() 
		WHERE id = $2
	`

	_, err := pool.Exec(ctx, query, status, id)
	return err
}

// queryWebhookData 执行查询并返回结果
func queryWebhookData(ctx context.Context, pool *pgxpool.Pool, query string) ([]*WebhookDataModel, error) {
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	return scanWebhookDataRows(rows)
}

// queryWebhookDataWithParams 执行带参数的查询并返回结果
func queryWebhookDataWithParams(ctx context.Context, pool *pgxpool.Pool, query string, args ...interface{}) ([]*WebhookDataModel, error) {
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	return scanWebhookDataRows(rows)
}

// scanWebhookDataRows 扫描数据库行并转换为模型
func scanWebhookDataRows(rows pgx.Rows) ([]*WebhookDataModel, error) {
	var result []*WebhookDataModel

	for rows.Next() {
		var data WebhookDataModel
		err := rows.Scan(
			&data.ID, &data.BlockHeight, &data.TxHash, &data.FromAddress,
			&data.ToAddress, &data.Value, &data.BlockTime, &data.CreateTime,
			&data.UpdateTime, &data.ExpireTime, &data.Status, &data.OriginalTxID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}
		result = append(result, &data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rows: %w", err)
	}

	return result, nil
}

// GetWebhookDataStats 获取webhook数据统计信息
func GetWebhookDataStats(ctx context.Context, pool *pgxpool.Pool) (map[int16]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM webhook_data 
		GROUP BY status
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query statistics: %w", err)
	}
	defer rows.Close()

	stats := make(map[int16]int)
	for rows.Next() {
		var status int16
		var count int
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan statistics: %w", err)
		}
		stats[status] = count
	}

	return stats, nil
}

// UpdateOriginalTxIDByID 根据ID更新original_tx_id字段
func UpdateOriginalTxIDByID(ctx context.Context, pool *pgxpool.Pool, id int64, originalTxID string) error {
	query := `
		UPDATE webhook_data 
		SET original_tx_id = $1, update_time = NOW() 
		WHERE id = $2
	`

	_, err := pool.Exec(ctx, query, originalTxID, id)
	return err
}

// UpdateOriginalTxIDByTxHash 根据交易哈希更新original_tx_id字段
func UpdateOriginalTxIDByTxHash(ctx context.Context, pool *pgxpool.Pool, txHash string, originalTxID string) error {
	query := `
		UPDATE webhook_data 
		SET original_tx_id = $1, update_time = NOW() 
		WHERE tx_hash = $2
	`

	_, err := pool.Exec(ctx, query, originalTxID, txHash)
	return err
}

// GetOriginalTxIDByID 根据ID获取original_tx_id
func GetOriginalTxIDByID(ctx context.Context, pool *pgxpool.Pool, id int64) (string, error) {
	query := `
		SELECT original_tx_id 
		FROM webhook_data 
		WHERE id = $1
	`

	var originalTxID string
	err := pool.QueryRow(ctx, query, id).Scan(&originalTxID)
	if err != nil {
		return "", fmt.Errorf("failed to get original_tx_id: %w", err)
	}

	return originalTxID, nil
}

// GetOriginalTxIDByTxHash 根据交易哈希获取original_tx_id
func GetOriginalTxIDByTxHash(ctx context.Context, pool *pgxpool.Pool, txHash string) (string, error) {
	query := `
		SELECT original_tx_id 
		FROM webhook_data 
		WHERE tx_hash = $1
	`

	var originalTxID string
	err := pool.QueryRow(ctx, query, txHash).Scan(&originalTxID)
	if err != nil {
		return "", fmt.Errorf("failed to get original_tx_id: %w", err)
	}

	return originalTxID, nil
}
