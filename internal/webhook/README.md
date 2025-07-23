# Webhook 工具函数

## 概述

本模块提供了用于解析和处理webhook数据的工具函数，主要功能是将接收到的JSON数据转换为`WebhookData`结构体数组，并支持批量插入数据库。

## 主要功能

### ParseWebhookData

将webhook请求的JSON数据解析为`WebhookData`数组。

```go
func ParseWebhookData(body []byte) ([]WebhookData, error)
```

**参数:**
- `body []byte`: webhook请求的JSON数据

**返回值:**
- `[]WebhookData`: 解析后的交易数据数组
- `error`: 解析错误

### ConvertToWebhookDataModel

将单个`WebhookData`转换为`WebhookDataModel`。

```go
func ConvertToWebhookDataModel(data WebhookData) *db.WebhookDataModel
```

### ConvertToWebhookDataModelSlice

将`WebhookData`切片转换为`WebhookDataModel`切片，用于批量插入。

```go
func ConvertToWebhookDataModelSlice(dataList []WebhookData) []*db.WebhookDataModel
```

## 数据结构

### WebhookRequest
完整的webhook请求结构：
```go
type WebhookRequest struct {
    Data     []TransactionData `json:"data"`
    Metadata Metadata          `json:"metadata"`
}
```

### TransactionData
单个交易数据结构：
```go
type TransactionData struct {
    BlockHash        string `json:"blockHash"`
    BlockNumber      string `json:"blockNumber"`
    From             string `json:"from"`
    Hash             string `json:"hash"`
    Timestamp        string `json:"timestamp"`
    To               string `json:"to"`
    Value            string `json:"value"`
    // ... 其他字段
}
```

### WebhookData
转换后的数据结构：
```go
type WebhookData struct {
    BlockHeight int64  `json:"blockNumber"`
    TxHash      string `json:"hash"`
    FromAddress string `json:"from"`
    ToAddress   string `json:"to"`
    Value       string `json:"value"`
    BlockTime   int64  `json:"timestamp"`
    ExpireTime  int64  `json:"expire_time"`
    Status      int16  `json:"status"`
}
```

## 字段转换规则

1. **BlockNumber**: 十六进制字符串 → int64
   - 输入: `"0x46c451a"`
   - 输出: `74204442`

2. **Timestamp**: 十六进制字符串 → int64
   - 输入: `"0x6880ce30"`
   - 输出: `1753271856`

3. **Value**: 十六进制字符串 → 十进制字符串
   - 输入: `"0x6"`
   - 输出: `"6"`

4. **其他字段**: 直接复制
   - Hash, From, To 等字段保持不变

## 批量插入功能

### 使用示例

```go
// 读取webhook请求体
body, err := io.ReadAll(c.Request.Body)
if err != nil {
    return err
}

// 解析数据
webhookDataList, err := ParseWebhookData(body)
if err != nil {
    return err
}

// 转换为WebhookDataModel并批量插入
webhookDataModels := ConvertToWebhookDataModelSlice(webhookDataList)
err = db.BatchInsertWebhookData(ctx, pool, webhookDataModels)
if err != nil {
    return err
}

fmt.Printf("成功插入 %d 条记录\n", len(webhookDataList))
```

### 批量插入的优势

1. **性能提升**: 相比逐条插入，批量插入可以显著提高数据库写入性能
2. **事务一致性**: 所有数据在同一个事务中插入，保证数据一致性
3. **减少网络开销**: 减少与数据库的网络交互次数
4. **自动处理**: 自动处理时间戳和默认值

### 数据库函数

```go
// BatchInsertWebhookData 批量插入 webhook_data 记录
func BatchInsertWebhookData(ctx context.Context, pool *pgxpool.Pool, data []*WebhookDataModel) error
```

**参数:**
- `ctx context.Context`: 上下文
- `pool *pgxpool.Pool`: 数据库连接池
- `data []*WebhookDataModel`: 要插入的数据数组

**返回值:**
- `error`: 插入错误

## 测试

运行测试：
```bash
go test ./internal/webhook -v
```

测试覆盖：
- ✅ 数据解析功能
- ✅ 十六进制转换功能
- ✅ 数据模型转换功能
- ✅ 批量转换功能

## 注意事项

1. 函数会自动处理十六进制到十进制的转换
2. 对于空字符串或无效的十六进制值，会返回默认值
3. `ExpireTime`会自动设置为`BlockTime + 1小时`（3600秒）
4. `Status`字段会设置为默认值（0），可根据业务需求调整
5. 支持批量处理多个交易数据
6. 批量插入会自动设置`CreateTime`和`UpdateTime`字段
7. 使用事务保证数据一致性，如果任何一条记录插入失败，整个批次都会回滚 