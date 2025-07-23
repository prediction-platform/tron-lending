# Lending TRX - Tron 能量委托系统

## 概述

这是一个基于 Go 语言开发的 Tron 区块链能量委托系统，支持自动化的能量委托和取消委托功能。系统通过 webhook 接收交易数据，使用定时任务处理能量委托业务逻辑。

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Webhook       │    │   CronJob       │    │   Tron API      │
│   Handler       │    │   Processor     │    │   Client        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Database      │    │   Business      │    │   Blockchain    │
│   Operations    │    │   Logic         │    │   Integration   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 核心模块

### 1. Webhook 模块 (`internal/webhook/`)
- 接收和处理 webhook 数据
- JSON 解析和字段转换
- 批量数据插入
- HTTP 接口处理

### 2. 数据库模块 (`internal/db/`)
- 数据库连接和初始化
- 数据模型定义
- 查询和更新操作
- 批量操作优化

### 3. 定时任务模块 (`internal/cronjob/`)
- 定时处理 webhook 数据
- 能量委托业务逻辑
- 状态管理和更新
- 错误处理和重试

### 4. Tron API 模块 (`internal/tron/`)
- Tron 区块链 API 客户端
- 能量委托和取消委托
- 账户信息查询
- 交易状态监控

## 功能特性

### 🚀 核心功能
- ✅ **Webhook 数据接收**: 支持批量 JSON 数据接收
- ✅ **数据解析转换**: 自动处理十六进制和数值转换
- ✅ **批量数据库操作**: 高效的批量插入和更新
- ✅ **定时任务处理**: 自动化的数据处理流程
- ✅ **能量委托管理**: 智能的能量委托和取消
- ✅ **状态跟踪**: 完整的数据状态管理

### 🔧 技术特性
- ✅ **模块化设计**: 清晰的模块分离和职责划分
- ✅ **错误处理**: 完善的错误处理和日志记录
- ✅ **配置管理**: 环境变量配置支持
- ✅ **测试覆盖**: 完整的单元测试
- ✅ **文档完善**: 详细的 API 和使用文档

## 快速开始

### 1. 环境准备

```bash
# 克隆项目
git clone <repository-url>
cd lending-trx

# 安装依赖
go mod tidy

# 复制配置文件
cp config.example.env .env
```

### 2. 配置环境变量

编辑 `.env` 文件：

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

### 3. 启动数据库

```bash
# 使用 Docker Compose
docker-compose up -d

# 或者直接启动 PostgreSQL
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=postgres \
  -p 5432:5432 \
  postgres:13
```

### 4. 运行应用

```bash
# 启动应用
go run main.go

# 或者构建后运行
go build -o lending-trx
./lending-trx
```

## API 使用

### 1. 接收 Webhook 数据

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

### 2. 查看数据统计

```bash
curl http://localhost:8080/stats
```

## 数据流程

### 1. 数据接收流程
```
Webhook 数据 → JSON 解析 → 字段转换 → 批量插入 → 数据库
```

### 2. 处理流程
```
定时任务 → 查询待处理数据 → 执行能量委托 → 保存原始交易ID → 更新状态
```

### 3. 过期处理流程
```
定时任务 → 查询过期数据 → 获取原始交易ID → 取消能量委托 → 更新状态
```

## 状态管理

| 状态 | 说明 | 处理逻辑 |
|------|------|----------|
| 0 | 初始化 | 等待处理 |
| 1 | 执行中 | 正在处理 |
| 2 | 已授权 | 等待过期 |
| 3 | 已回收 | 最终状态 |

## 原始交易ID管理

### 字段说明
- **original_tx_id**: 存储 Tron 能量委托成功后返回的交易ID
- **用途**: 用于后续取消委托时提供必要的交易标识
- **唯一性**: 确保每个委托交易ID的唯一性

### 数据库操作
```go
// 保存原始交易ID
err := db.UpdateOriginalTxIDByID(ctx, pool, data.ID, delegationResp.TxID)

// 获取原始交易ID
originalTxID, err := db.GetOriginalTxIDByID(ctx, pool, data.ID)

// 根据交易哈希更新
err := db.UpdateOriginalTxIDByTxHash(ctx, pool, txHash, originalTxID)
```

## 能量委托策略

### 统一委托地址
- 所有能量委托都使用环境变量 `DELEGATION_FROM_ADDRESS` 指定的统一地址
- 原始交易中的 `FromAddress` 仅用于记录，实际委托使用统一地址
- 支持不同环境使用不同的委托地址

### 委托数量计算
- **小额交易** (≤ 10万): 委托 10% 可用能量
- **中额交易** (10万-100万): 委托 8% 可用能量  
- **大额交易** (> 100万): 委托 5% 可用能量
- **最小委托**: 1000 能量

### 过期时间设置
- 默认过期时间: `BlockTime + 1小时`
- 可配置的过期策略
- 自动过期处理

## 开发指南

### 1. 项目结构

```
lending-trx/
├── internal/
│   ├── webhook/     # Webhook 处理模块
│   ├── db/          # 数据库操作模块
│   ├── cronjob/     # 定时任务模块
│   └── tron/        # Tron API 客户端
├── main.go          # 应用入口
├── go.mod           # Go 模块文件
├── go.sum           # 依赖校验文件
├── docker-compose.yml # Docker 编排文件
└── README.md        # 项目文档
```

### 2. 添加新功能

#### 扩展数据库操作
```go
// 在 internal/db/db.go 中添加新函数
func NewDatabaseFunction(ctx context.Context, pool *pgxpool.Pool) error {
    // 实现新功能
    return nil
}
```

#### 扩展 Tron API
```go
// 在 internal/tron/client.go 中添加新方法
func (c *TronClient) NewAPIMethod(ctx context.Context) error {
    // 实现新的 API 调用
    return nil
}
```

#### 扩展业务逻辑
```go
// 在 internal/cronjob/cron.go 中添加新方法
func (c *CronJob) newBusinessLogic(data *db.WebhookDataModel) error {
    // 实现新的业务逻辑
    return nil
}
```

### 3. 测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定模块测试
go test ./internal/webhook -v
go test ./internal/tron -v
go test ./internal/cronjob -v
go test ./internal/db -v

# 运行测试并生成覆盖率报告
go test ./... -cover
```

## 部署

### 1. Docker 部署

```bash
# 构建镜像
docker build -t lending-trx .

# 运行容器
docker run -d \
  --name lending-trx \
  --env-file .env \
  -p 8080:8080 \
  lending-trx
```

### 2. Docker Compose 部署

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### 3. 生产环境配置

```bash
# 设置生产环境变量
export APP_ENV=production
export LOG_LEVEL=warn
export TRON_API_KEY=your-production-api-key

# 启动应用
./lending-trx
```

## 监控和日志

### 1. 日志级别
- `DEBUG`: 详细的调试信息
- `INFO`: 一般信息
- `WARN`: 警告信息
- `ERROR`: 错误信息

### 2. 关键指标
- Webhook 接收成功率
- 数据库操作性能
- Tron API 调用成功率
- 能量委托成功率

### 3. 告警配置
- API 调用失败告警
- 数据库连接异常告警
- 定时任务执行失败告警

## 故障排除

### 1. 常见问题

#### 数据库连接失败
```bash
# 检查数据库状态
docker ps | grep postgres

# 检查连接配置
echo $PG_DSN
```

#### Tron API 调用失败
```bash
# 检查 API 密钥
echo $TRON_API_KEY

# 检查网络连接
curl -I https://api.trongrid.io
```

#### 定时任务不执行
```bash
# 检查日志
docker-compose logs cronjob

# 检查环境变量
echo $CRON_SCHEDULE
```

### 2. 调试模式

```bash
# 启用调试日志
export LOG_LEVEL=debug

# 启动应用
go run main.go
```

## 贡献指南

### 1. 代码规范
- 遵循 Go 语言规范
- 使用 `gofmt` 格式化代码
- 添加适当的注释和文档

### 2. 提交规范
- 使用清晰的提交信息
- 每个提交只包含一个功能
- 添加测试用例

### 3. 测试要求
- 新功能必须包含测试
- 保持测试覆盖率 > 80%
- 集成测试覆盖主要流程

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue: [GitHub Issues](https://github.com/your-repo/issues)
- 邮箱: your-email@example.com

---

**注意**: 在生产环境中使用前，请确保：
1. 正确配置 Tron API 密钥
2. 测试所有功能
3. 配置适当的监控和告警
4. 备份重要数据