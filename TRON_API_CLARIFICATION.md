# Tron API 说明文档

## 重要澄清

### 1. **真实的 Tron API 要求**

Tron 区块链的能量委托 API **实际上只需要以下字段**：

```json
{
  "from_address": "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
  "to_address": "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "amount": "10000"
}
```

### 2. **我们添加的业务字段**

在我们的实现中，我们添加了额外的字段用于**业务追踪**：

```json
{
  "from_address": "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
  "to_address": "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "amount": "10000",
  "tx_hash": "0x1234567890abcdef",        // 业务追踪字段
  "block_height": 74204442                 // 业务追踪字段
}
```

## 字段说明

### **必需字段（Tron API 要求）**

| 字段 | 类型 | 说明 | 是否必需 |
|------|------|------|----------|
| `from_address` | string | 委托方地址 | ✅ 是 |
| `to_address` | string | 接收方地址 | ✅ 是 |
| `amount` | string | 委托的能量数量 | ✅ 是 |

### **业务追踪字段（我们添加）**

| 字段 | 类型 | 说明 | 用途 |
|------|------|------|------|
| `tx_hash` | string | 原始交易哈希 | 追踪来源、防重复 |
| `block_height` | int64 | 区块高度 | 时间管理、审计 |

## 实际使用场景

### 1. **Tron API 调用**

当实际调用 Tron API 时，我们只发送必需字段：

```go
// 实际的 Tron API 请求
actualRequest := map[string]interface{}{
    "from_address": delegationReq.FromAddress,
    "to_address":   delegationReq.ToAddress,
    "amount":       delegationReq.Amount,
}
```

### 2. **业务逻辑处理**

在我们的业务逻辑中，我们使用所有字段：

```go
// 业务逻辑中的完整请求
businessRequest := &tron.EnergyDelegationRequest{
    FromAddress: delegationFromAddress,
    ToAddress:   data.FromAddress,
    Amount:      delegationAmount,
    TxHash:      data.TxHash,        // 用于业务追踪
    BlockHeight: data.BlockHeight,   // 用于时间管理
}
```

## API 端点说明

### 1. **真实的 Tron API 端点**

Tron 的能量委托可能使用以下端点：
- `/wallet/delegateresource` - 委托资源
- `/wallet/undelegateresource` - 取消委托资源

### 2. **我们的 API 端点**

我们使用简化的端点名称：
- `/v1/energy/delegate` - 能量委托
- `/v1/energy/cancel-delegate` - 取消能量委托

**注意**: 这些是我们**模拟的端点**，用于演示目的。

## 实现建议

### 1. **生产环境实现**

在生产环境中，你需要：

1. **查找真实的 Tron API 文档**
2. **使用正确的 API 端点**
3. **只发送 Tron API 要求的字段**
4. **将业务追踪字段存储在数据库中**

### 2. **代码修改建议**

```go
// 修改 DelegateEnergy 方法
func (c *TronClient) DelegateEnergy(ctx context.Context, req *EnergyDelegationRequest) (*EnergyDelegationResponse, error) {
    // 只发送 Tron API 要求的字段
    tronRequest := map[string]interface{}{
        "from_address": req.FromAddress,
        "to_address":   req.ToAddress,
        "amount":       req.Amount,
    }
    
    // 调用真实的 Tron API
    // ...
    
    // 将业务追踪字段存储到数据库
    // ...
}
```

### 3. **数据库存储**

```sql
-- 存储业务追踪信息
INSERT INTO delegation_tracking (
    tron_tx_id,
    original_tx_hash,
    block_height,
    delegation_amount,
    created_at
) VALUES (
    'tron_tx_123',
    '0x1234567890abcdef',
    74204442,
    '10000',
    NOW()
);
```

## 总结

1. **Tron API 只要求基本字段**: `from_address`, `to_address`, `amount`
2. **我们添加的字段用于业务追踪**: `tx_hash`, `block_height`
3. **在生产环境中**: 需要查找真实的 Tron API 文档
4. **业务追踪字段**: 应该存储在数据库中，而不是发送给 Tron API

这样的设计既满足了 Tron API 的要求，又保持了业务逻辑的完整性。 