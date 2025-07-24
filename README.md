# TRX委托服务

一个用于管理TRX能量委托的Go应用程序，包含HTTP API服务、定时任务处理和Telegram Bot监控功能。

## 🚀 快速开始

### 使用命令行工具

```bash
# 构建项目
make build

# 启动完整服务 (HTTP API + 定时任务)
make server

# 启动Telegram Bot
make bot

# 仅启动定时任务
make cron

# 开发模式运行
make dev
```

### 直接使用二进制文件

```bash
# 构建
go build -o lending-trx cmd/root/*.go

# 启动完整服务
./lending-trx server

# 启动Telegram Bot
./lending-trx bot

# 启动定时任务
./lending-trx cron

# 显示帮助
./lending-trx --help
```

## 📁 项目结构

```
lending-trx/
├── cmd/
│   └── root/           # 命令行工具
│       ├── main.go     # 根命令入口
│       ├── server.go   # server子命令
│       └── bot.go      # bot子命令
├── internal/
│   ├── cronjob/        # 定时任务处理
│   ├── db/             # 数据库操作
│   ├── tron/           # TRON API客户端
│   └── webhook/        # HTTP API处理
├── pkg/
│   └── telegram_bot/   # Telegram Bot包
├── migrations/         # 数据库迁移
├── main.go            # 传统版本入口
├── Makefile           # 构建和部署脚本
└── README.md          # 项目文档
```

## 🔧 服务说明

### 1. 完整服务 (server)
启动HTTP API服务和定时任务处理：
- HTTP API端点：`/api/delegation-account`
- 定时处理webhook数据
- 自动能量委托管理

### 2. Telegram Bot (bot)
启动Telegram Bot监控服务：
- 查询委托账户状态
- 持续监控功能
- 告警通知

### 3. 定时任务 (cron)
仅启动定时任务处理：
- 处理webhook数据
- 执行能量委托逻辑
- 无HTTP服务

## ⚙️ 配置

### 环境变量

创建 `.env` 文件：

```bash
# 数据库配置
DATABASE_URL=postgresql://username:password@localhost:5432/lending_trx

# TRON API配置
TRON_API_URL=https://api.trongrid.io
DELEGATION_FROM_ADDRESS=your_delegation_address

# Telegram Bot配置
TELEGRAM_BOT_TOKEN=your_bot_token
API_BASE_URL=http://localhost:8080
MONITOR_INTERVAL_MINUTES=5

# HTTP配置
PORT=8080
```

### 配置示例

```bash
# 复制配置示例
cp config.example.env .env

# 编辑配置
vim .env
```

## 🛠️ 开发

### 安装依赖

```bash
make deps
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行Bot测试
make test-bot
```

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

## 🐳 Docker部署

### 构建镜像

```bash
make docker-build
```

### 运行容器

```bash
make docker-run
```

## 📊 API接口

### 查询委托账户信息

```bash
GET /api/delegation-account
```

响应示例：
```json
{
  "status": "ok",
  "data": {
    "address": "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
    "balance": "1000000000",
    "energy": "50000",
    "energy_limit": "100000",
    "energy_used": "30000"
  }
}
```

## 🤖 Telegram Bot

### 可用命令

- `/start` - 启动Bot
- `/help` - 显示帮助信息
- `/status` - 查询委托账户状态
- `/monitor` - 开始持续监控
- `/stop` - 停止监控

### 功能特性

- 实时账户状态查询
- 自动监控和告警
- 余额和能量阈值告警
- 支持多用户监控

## 🔄 委托逻辑

根据交易金额自动计算委托数量：

- 1 TRX → 委托 `delegationBase` (15000)
- 2 TRX → 委托 `2 * delegationBase` (30000)
- 最大委托限制：`2 * delegationBase`

## 📝 日志

日志文件位置：`logs/lending-trx.log`

日志级别可通过环境变量配置：
```bash
export LOG_LEVEL=debug
```

## 🚨 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 `DATABASE_URL` 配置
   - 确认PostgreSQL服务运行正常

2. **TRON API调用失败**
   - 检查网络连接
   - 验证API密钥配置

3. **Telegram Bot无响应**
   - 检查 `TELEGRAM_BOT_TOKEN` 配置
   - 确认Bot已启动

4. **委托失败**
   - 检查账户余额
   - 验证委托地址配置

## 📄 许可证

本项目遵循MIT许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📞 支持

如有问题，请通过以下方式联系：

- 提交GitHub Issue
- 发送邮件至项目维护者