# Tron API 客户端

## 概述

本模块提供了与 Tron 区块链 API 交互的客户端，支持能量委托、取消委托、账户信息查询等功能。

## 核心功能

### 1. TronClient 客户端

```go
type TronClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}
```

### 2. 主要功能

#### 账户信息查询
```go
func (c *TronClient) GetAccountInfo(ctx context.Context, address string) (*AccountInfo, error)
```

#### 能量委托
```go
func (c *TronClient) DelegateEnergy(ctx context.Context, req *EnergyDelegationRequest) (*EnergyDelegationResponse, error)
```

#### 取消能量委托
```go
func (c *TronClient) CancelEnergyDelegation(ctx context.Context, req *CancelDelegationRequest) (*CancelDelegationResponse, error)
```

#### 交易信息查询
```go
func (c *TronClient) GetTransactionInfo(ctx context.Context, txID string) (map[string]interface{}, error)
```

## 数据结构

### EnergyDelegationRequest 能量委托请求
```go
type EnergyDelegationRequest struct {
    FromAddress string `json:"from_address"` // 委托方地址
    ToAddress   string `json:"to_address"`   // 接收方地址
    Amount      string `json:"amount"`       // 委托的能量数量
    // 以下字段用于业务追踪，不是 Tron API 必需字段
    TxHash      string `json:"tx_hash,omitempty"`      // 原始交易哈希（业务追踪）
    BlockHeight int64  `json:"block_height,omitempty"` // 区块高度（业务追踪）
}
```

**注意**: `TxHash` 和 `BlockHeight` 字段是**业务追踪字段**，不是 Tron API 的必需字段。这些字段用于：
- 追踪委托操作的来源交易
- 审计和日志记录
- 防重复处理
- 时间管理（计算过期时间）

### EnergyDelegationResponse 能量委托响应
```go
type EnergyDelegationResponse struct {
    Success bool   `json:"success"` // 是否成功
    TxID    string `json:"tx_id,omitempty"`    // 交易ID
    Message string `json:"message,omitempty"`   // 消息
    Error   string `json:"error,omitempty"`     // 错误信息
}
```

### CancelDelegationRequest 取消委托请求
```go
type CancelDelegationRequest struct {
    FromAddress  string `json:"from_address"`   // 委托方地址
    ToAddress    string `json:"to_address"`     // 接收方地址
    OriginalTxID string `json:"original_tx_id"` // 原始委托交易ID
    // 以下字段用于业务追踪，不是 Tron API 必需字段
    TxHash string `json:"tx_hash,omitempty"` // 原始交易哈希（业务追踪）
}
```

**注意**: `TxHash` 字段是**业务追踪字段**，不是 Tron API 的必需字段。这个字段用于：
- 追踪取消委托操作的来源
- 审计和日志记录
- 关联原始交易信息

### AccountInfo 账户信息
```go
type AccountInfo struct {
    Address     string `json:"address"`      // 账户地址
    Balance     string `json:"balance"`      // 余额
    Energy      string `json:"energy"`       // 能量
    Frozen      string `json:"frozen"`       // 冻结金额
    NetUsed     string `json:"net_used"`     // 已使用网络资源
    NetLimit    string `json:"net_limit"`    // 网络资源限制
    EnergyUsed  string `json:"energy_used"`  // 已使用能量
    EnergyLimit string `json:"energy_limit"` // 能量限制
}
```

## 使用示例

### 1. 创建客户端
```go
import "lending-trx/internal/tron"

// 使用默认配置
client := tron.NewTronClient("https://api.trongrid.io", "your-api-key")

// 使用环境变量
baseURL := os.Getenv("TRON_API_BASE_URL")
apiKey := os.Getenv("TRON_API_KEY")
client := tron.NewTronClient(baseURL, apiKey)
```

### 2. 查询账户信息
```go
ctx := context.Background()
address := "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"

accountInfo, err := client.GetAccountInfo(ctx, address)
if err != nil {
    log.Printf("获取账户信息失败: %v", err)
    return
}

fmt.Printf("账户地址: %s\n", accountInfo.Address)
fmt.Printf("余额: %s\n", accountInfo.Balance)
fmt.Printf("可用能量: %s\n", accountInfo.Energy)
fmt.Printf("能量限制: %s\n", accountInfo.EnergyLimit)
```

### 3. 执行能量委托
```go
delegationReq := &tron.EnergyDelegationRequest{
    FromAddress: "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
    ToAddress:   "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
    Amount:      "10000",
    TxHash:      "0x1234567890abcdef",
    BlockHeight: 74204442,
}

delegationResp, err := client.DelegateEnergy(ctx, delegationReq)
if err != nil {
    log.Printf("能量委托失败: %v", err)
    return
}

if delegationResp.Success {
    fmt.Printf("委托成功，交易ID: %s\n", delegationResp.TxID)
} else {
    fmt.Printf("委托失败: %s\n", delegationResp.Error)
}
```

### 4. 取消能量委托
```go
cancelReq := &tron.CancelDelegationRequest{
    FromAddress:  "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
    ToAddress:    "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
    OriginalTxID: "original_delegation_tx_id",
    TxHash:       "0x1234567890abcdef",
}

cancelResp, err := client.CancelEnergyDelegation(ctx, cancelReq)
if err != nil {
    log.Printf("取消委托失败: %v", err)
    return
}

if cancelResp.Success {
    fmt.Printf("取消委托成功，交易ID: %s\n", cancelResp.TxID)
} else {
    fmt.Printf("取消委托失败: %s\n", cancelResp.Error)
}
```

### 5. 查询交易信息
```go
txID := "tx_id_123456"
txInfo, err := client.GetTransactionInfo(ctx, txID)
if err != nil {
    log.Printf("查询交易信息失败: %v", err)
    return
}

fmt.Printf("交易信息: %+v\n", txInfo)
```

## 环境变量配置

### 必需的环境变量
```bash
# Tron API 基础URL
export TRON_API_BASE_URL="https://api.trongrid.io"

# Tron API 密钥
export TRON_API_KEY="your-api-key-here"

# 统一委托方地址
export DELEGATION_FROM_ADDRESS="TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"
```

### 可选的环境变量
```bash
# 自定义API端点
export TRON_API_BASE_URL="https://your-custom-tron-api.com"

# 测试网络
export TRON_API_BASE_URL="https://api.shasta.trongrid.io"

# 不同的委托方地址（测试环境）
export DELEGATION_FROM_ADDRESS="TTestAddressForTesting123456789"
```

## API 端点

### 1. 账户信息
```
GET /v1/accounts/{address}
```

### 2. 能量委托
```
POST /v1/energy/delegate
Content-Type: application/json
TRON-PRO-API-KEY: your-api-key

{
  "from_address": "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
  "to_address": "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "amount": "10000"
}
```

**注意**: 上面的示例是**简化的 Tron API 调用**。实际的 Tron 能量委托可能使用不同的 API 端点和参数格式。`tx_hash` 和 `block_height` 字段是我们添加的业务追踪字段。

### 3. 取消能量委托
```
POST /v1/energy/cancel-delegate
Content-Type: application/json
TRON-PRO-API-KEY: your-api-key

{
  "from_address": "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
  "to_address": "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "original_tx_id": "original_delegation_tx_id"
}
```

**注意**: 上面的示例是**简化的 Tron API 调用**。实际的 Tron 取消能量委托可能使用不同的 API 端点和参数格式。`tx_hash` 字段是我们添加的业务追踪字段。

### 4. 交易信息
```
GET /v1/transactions/{tx_id}
```

## 错误处理

### 1. 网络错误
```go
if err != nil {
    return fmt.Errorf("请求失败: %w", err)
}
```

### 2. API 错误
```go
if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
}
```

### 3. 业务逻辑错误
```go
if !response.Success {
    return fmt.Errorf("能量委托失败: %s", response.Error)
}
```

## 性能优化

### 1. 连接池管理
- 使用 `http.Client` 管理连接
- 设置合理的超时时间
- 复用 HTTP 连接

### 2. 请求优化
- 使用 `context` 控制超时
- 设置合适的请求头
- 处理重试逻辑

### 3. 响应处理
- 及时关闭响应体
- 使用 `io.ReadAll` 读取完整响应
- 错误时记录详细日志

## 安全考虑

### 1. API 密钥管理
- 使用环境变量存储 API 密钥
- 不要在代码中硬编码密钥
- 定期轮换 API 密钥

### 2. 请求验证
- 验证输入参数
- 检查地址格式
- 验证数值范围

### 3. 错误信息
- 不要暴露敏感信息
- 记录详细的错误日志
- 实现适当的重试机制

## 测试

### 1. 单元测试
```bash
go test ./internal/tron -v
```

### 2. 集成测试
```bash
# 设置测试环境变量
export TRON_API_KEY="test-api-key"
go test ./internal/tron -v -tags=integration
```

### 3. 模拟测试
```go
// 使用 mock 进行测试
mockClient := &MockTronClient{}
// 设置期望行为
// 执行测试
```

## 注意事项

1. **API 限制**: 注意 Tron API 的调用频率限制
2. **网络延迟**: 考虑网络延迟对交易确认的影响
3. **错误重试**: 实现适当的重试机制
4. **日志记录**: 记录详细的 API 调用日志
5. **监控告警**: 监控 API 调用成功率和响应时间 