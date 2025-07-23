# 数据库模块 (DB)

## 概述

本模块提供了完整的数据库操作功能，包括webhook数据的增删改查、状态管理和统计信息。

## 核心功能

### 1. 数据模型

```go
type WebhookDataModel struct {
    ID          int64  `json:"id"`           // 主键唯一ID
    BlockHeight int64  `json:"block_height"` // 区块高度
    TxHash      string `json:"tx_hash"`      // 交易哈希
    FromAddress string `json:"from_address"` // 发送方地址
    ToAddress   string `json:"to_address"`   // 接收方地址
    Value       string `json:"value"`        // 交易金额（大整数，字符串存储）
    BlockTime   int64  `json:"block_time"`   // 区块时间（毫秒时间戳）
    CreateTime  string `json:"create_time"`  // 创建时间
    UpdateTime  string `json:"update_time"`  // 更新时间
    ExpireTime  int64  `json:"expire_time"`  // 有效期（毫秒时间戳）
    Status      int16  `json:"status"`       // 状态（0:初始化，1:执行中，2:已授权，3:已回收）
}
```

### 2. 数据库操作函数

#### 初始化函数
```go
// InitDB 连接数据库并初始化表
func InitDB(ctx context.Context) (*pgxpool.Pool, error)
```

#### 查询函数
```go
// QueryPendingWebhookData 查询待处理的数据 (status=0)
func QueryPendingWebhookData(ctx context.Context, pool *pgxpool.Pool) ([]*WebhookDataModel, error)

// QueryExpiredWebhookData 查询已过期且已授权的数据 (status=2)
func QueryExpiredWebhookData(ctx context.Context, pool *pgxpool.Pool) ([]*WebhookDataModel, error)

// GetWebhookDataStats 获取webhook数据统计信息
func GetWebhookDataStats(ctx context.Context, pool *pgxpool.Pool) (map[int16]int, error)
```

#### 更新函数
```go
// UpdateWebhookStatusByID 更新单个记录的status
func UpdateWebhookStatusByID(ctx context.Context, pool *pgxpool.Pool, id int64, status int16) error

// UpdateWebhookStatus 批量更新指定 id 的 status
func UpdateWebhookStatus(ctx context.Context, pool *pgxpool.Pool, ids []int64, status int16) error

// UpdateWebhookStatusAndExpireTime 批量更新指定 id 的 status 和 expire_time
func UpdateWebhookStatusAndExpireTime(ctx context.Context, pool *pgxpool.Pool, ids []int64, status int16, expireTime int64) error
```

#### 插入函数
```go
// BatchInsertWebhookData 批量插入 webhook_data 记录
func BatchInsertWebhookData(ctx context.Context, pool *pgxpool.Pool, data []*WebhookDataModel) error
```

## 使用示例

### 1. 初始化数据库
```go
ctx := context.Background()
pool, err := db.InitDB(ctx)
if err != nil {
    log.Fatal(err)
}
defer pool.Close()
```

### 2. 查询待处理数据
```go
pendingData, err := db.QueryPendingWebhookData(ctx, pool)
if err != nil {
    log.Printf("查询失败: %v", err)
    return
}

for _, item := range pendingData {
    fmt.Printf("ID: %d, Status: %d, TxHash: %s\n", 
        item.ID, item.Status, item.TxHash)
}
```

### 3. 查询过期数据
```go
expiredData, err := db.QueryExpiredWebhookData(ctx, pool)
if err != nil {
    log.Printf("查询失败: %v", err)
    return
}

for _, item := range expiredData {
    fmt.Printf("ID: %d, ExpireTime: %d\n", 
        item.ID, item.ExpireTime)
}
```

### 4. 更新状态
```go
// 单个更新
err := db.UpdateWebhookStatusByID(ctx, pool, 1, 1)
if err != nil {
    log.Printf("更新失败: %v", err)
}

// 批量更新
ids := []int64{1, 2, 3}
err = db.UpdateWebhookStatus(ctx, pool, ids, 2)
if err != nil {
    log.Printf("批量更新失败: %v", err)
}
```

### 5. 批量插入
```go
data := []*db.WebhookDataModel{
    {
        BlockHeight: 74204442,
        TxHash:      "0x1234567890abcdef",
        FromAddress: "0xabcdef1234567890",
        ToAddress:   "0x1234567890abcdef",
        Value:       "1000000",
        BlockTime:   time.Now().UnixMilli(),
        ExpireTime:  time.Now().Add(time.Hour).UnixMilli(),
        Status:      0,
    },
}

err := db.BatchInsertWebhookData(ctx, pool, data)
if err != nil {
    log.Printf("批量插入失败: %v", err)
}
```

### 6. 获取统计信息
```go
stats, err := db.GetWebhookDataStats(ctx, pool)
if err != nil {
    log.Printf("获取统计信息失败: %v", err)
    return
}

for status, count := range stats {
    fmt.Printf("状态 %d: %d 条记录\n", status, count)
}
```

## 状态管理

| 状态 | 说明 | 处理逻辑 |
|------|------|----------|
| 0 | 初始化 | 等待处理 |
| 1 | 执行中 | 正在处理 |
| 2 | 已授权 | 等待过期 |
| 3 | 已回收 | 最终状态 |

## 数据库表结构

### webhook_data 表
```sql
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
  status SMALLINT
);
```

### logs 表
```sql
CREATE TABLE IF NOT EXISTS logs (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## 性能优化

### 1. 连接池管理
- 使用 `pgxpool` 管理数据库连接池
- 自动处理连接的创建和回收
- 支持连接池大小配置

### 2. 批量操作
- 支持批量插入和批量更新
- 减少数据库交互次数
- 提高操作效率

### 3. 查询优化
- 使用索引优化查询性能
- 支持分页查询
- 结果集大小控制

## 错误处理

### 1. 连接错误
```go
if err != nil {
    return nil, fmt.Errorf("无法连接数据库: %w", err)
}
```

### 2. 查询错误
```go
if err != nil {
    return nil, fmt.Errorf("查询失败: %w", err)
}
```

### 3. 扫描错误
```go
if err != nil {
    return nil, fmt.Errorf("扫描数据失败: %w", err)
}
```

## 监控和日志

### 1. 关键指标
- 数据库连接数
- 查询响应时间
- 错误率统计
- 数据量统计

### 2. 日志记录
- 连接状态日志
- 查询执行日志
- 错误详情日志
- 性能监控日志

## 配置选项

### 1. 数据库连接
```go
dsn := os.Getenv("PG_DSN")
if dsn == "" {
    dsn = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
}
```

### 2. 连接池配置
- 最大连接数
- 最小连接数
- 连接超时时间
- 查询超时时间

## 注意事项

1. **连接管理**: 正确关闭数据库连接
2. **事务处理**: 使用事务保证数据一致性
3. **错误处理**: 完善的错误处理和重试机制
4. **性能监控**: 监控数据库性能指标
5. **数据备份**: 定期备份重要数据 