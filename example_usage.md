# 统一委托地址使用示例

## 概述

本系统现在支持使用环境变量 `DELEGATION_FROM_ADDRESS` 来设置统一的委托方地址，所有能量委托操作都将使用这个统一地址，而不是原始交易中的 `FromAddress`。

## 配置设置

### 1. 环境变量配置

```bash
# 设置统一委托地址
export DELEGATION_FROM_ADDRESS="TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"

# 其他必需的环境变量
export TRON_API_BASE_URL="https://api.trongrid.io"
export TRON_API_KEY="your-tron-api-key-here"
export PG_DSN="postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
```

### 2. 配置文件示例

创建 `.env` 文件：

```bash
# 数据库配置
PG_DSN=postgres://postgres:password@localhost:5432/postgres?sslmode=disable

# Tron API 配置
TRON_API_BASE_URL=https://api.trongrid.io
TRON_API_KEY=your-tron-api-key-here

# 能量委托配置
DELEGATION_FROM_ADDRESS=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs

# 应用配置
APP_PORT=8080
APP_ENV=development
```

## 使用示例

### 1. 启动应用

```bash
# 设置环境变量
export DELEGATION_FROM_ADDRESS="TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"
export TRON_API_KEY="your-api-key"

# 启动应用
go run main.go
```

### 2. 发送 Webhook 数据

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "data": [
      {
        "block_height": "0x46b8a8a",
        "tx_hash": "0x1234567890abcdef",
        "from_address": "0xabcdef1234567890",
        "to_address": "0x1234567890abcdef",
        "value": "0x3b9aca00",
        "block_time": "0x64f8b8a8"
      }
    ]
  }'
```

### 3. 查看日志输出

应用会输出类似以下的日志：

```
INFO 开始执行能量委托 id=1 from=0xabcdef1234567890 to=0x1234567890abcdef value=1000000000 tx_hash=0x1234567890abcdef
INFO 使用统一委托地址 original_from=0xabcdef1234567890 delegation_from=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs
INFO 委托方账户信息 address=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs balance=1000000000 energy=50000 energy_limit=100000 energy_used=5000
INFO 计算委托数量 original_value=1000000000 available_energy=50000 delegation_amount=5000
INFO 能量委托成功 tx_id=tx_id_123456 message=委托成功 from=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs to=0x1234567890abcdef amount=5000
```

## 功能说明

### 1. 统一委托地址的优势

- **集中管理**: 所有委托操作使用同一个地址，便于管理
- **安全性**: 避免使用用户原始地址进行委托操作
- **灵活性**: 可以轻松切换不同的委托地址
- **监控**: 便于监控和统计委托操作

### 2. 地址映射关系

| 原始交易 | 实际委托 | 说明 |
|----------|----------|------|
| `from_address` | `DELEGATION_FROM_ADDRESS` | 原始交易发送方 |
| `to_address` | `to_address` | 委托接收方（保持不变） |

### 3. 环境变量验证

系统会在启动时验证环境变量：

```go
// 检查环境变量是否设置
delegationFromAddress := os.Getenv("DELEGATION_FROM_ADDRESS")
if delegationFromAddress == "" {
    return fmt.Errorf("环境变量 DELEGATION_FROM_ADDRESS 未设置")
}
```

### 4. 日志记录

系统会记录详细的地址映射信息：

```go
c.log.Info("使用统一委托地址", 
    "original_from", data.FromAddress,
    "delegation_from", delegationFromAddress,
)
```

## 测试验证

### 1. 环境变量测试

```bash
# 测试环境变量设置
export DELEGATION_FROM_ADDRESS="TTestAddressForTesting123456789"
echo $DELEGATION_FROM_ADDRESS
# 输出: TTestAddressForTesting123456789
```

### 2. 运行单元测试

```bash
# 运行 cronjob 测试
go test ./internal/cronjob -v

# 运行所有测试
go test ./... -v
```

### 3. 集成测试

```bash
# 启动应用
go run main.go

# 在另一个终端发送测试数据
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "data": [
      {
        "block_height": "0x46b8a8a",
        "tx_hash": "0x1234567890abcdef",
        "from_address": "0xabcdef1234567890",
        "to_address": "0x1234567890abcdef",
        "value": "0x3b9aca00",
        "block_time": "0x64f8b8a8"
      }
    ]
  }'
```

## 生产环境配置

### 1. 生产环境变量

```bash
# 生产环境配置
export DELEGATION_FROM_ADDRESS="TProductionAddressForEnergyDelegation"
export TRON_API_KEY="production-api-key"
export APP_ENV="production"
export LOG_LEVEL="warn"
```

### 2. Docker 部署

```bash
# 使用 Docker Compose
docker-compose up -d

# 或者直接运行容器
docker run -d \
  --name lending-trx \
  --env-file .env \
  -p 8080:8080 \
  lending-trx
```

### 3. 监控和告警

```bash
# 检查环境变量
docker exec lending-trx env | grep DELEGATION_FROM_ADDRESS

# 查看日志
docker logs lending-trx -f
```

## 故障排除

### 1. 环境变量未设置

**错误信息**: `环境变量 DELEGATION_FROM_ADDRESS 未设置`

**解决方案**:
```bash
# 设置环境变量
export DELEGATION_FROM_ADDRESS="TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"

# 或者添加到 .env 文件
echo "DELEGATION_FROM_ADDRESS=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs" >> .env
```

### 2. 地址格式错误

**错误信息**: `获取委托方账户信息失败`

**解决方案**:
```bash
# 确保地址格式正确（Tron 地址以 T 开头）
export DELEGATION_FROM_ADDRESS="TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs"
```

### 3. API 密钥问题

**错误信息**: `能量委托API调用失败`

**解决方案**:
```bash
# 检查 API 密钥
echo $TRON_API_KEY

# 重新设置 API 密钥
export TRON_API_KEY="your-valid-api-key"
```

## 最佳实践

### 1. 地址管理

- 使用专门的委托地址，不要与用户交易地址混用
- 定期轮换委托地址以提高安全性
- 监控委托地址的余额和能量状态

### 2. 环境分离

```bash
# 开发环境
export DELEGATION_FROM_ADDRESS="TDevAddressForTesting"

# 测试环境
export DELEGATION_FROM_ADDRESS="TTestAddressForTesting"

# 生产环境
export DELEGATION_FROM_ADDRESS="TProductionAddressForEnergy"
```

### 3. 监控和日志

- 记录所有委托操作的详细信息
- 监控委托地址的余额变化
- 设置委托失败的告警机制

### 4. 安全考虑

- 不要在代码中硬编码委托地址
- 使用环境变量管理敏感信息
- 定期审计委托操作日志 