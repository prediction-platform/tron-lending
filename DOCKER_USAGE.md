# Docker 部署指南

## 🚀 快速启动

### 1. 准备环境变量

复制环境变量示例文件：
```bash
cp docker.env.example .env
```

编辑 `.env` 文件，设置必要的环境变量：
```bash
# TRON API配置
TRON_API_KEY=your-tron-api-key-here

# 委托配置
DELEGATION_FROM_ADDRESS=TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs

# Telegram Bot配置
TELEGRAM_BOT_TOKEN=your-telegram-bot-token-here
```

### 2. 启动所有服务

```bash
# 启动所有服务 (数据库 + Server + Bot)
make docker-compose-up

# 或者直接使用docker-compose
docker-compose up -d
```

### 3. 查看服务状态

```bash
# 查看所有容器状态
docker-compose ps

# 查看服务日志
make docker-compose-logs

# 或者查看特定服务日志
docker-compose logs -f server
docker-compose logs -f bot
docker-compose logs -f postgres
```

## 📋 服务说明

### 服务架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │   Server        │    │   Bot           │
│   Database      │    │   (HTTP API)    │    │   (Telegram)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      Docker Network       │
                    └───────────────────────────┘
```

### 服务详情

1. **postgres** - PostgreSQL数据库
   - 端口: 5432
   - 数据库: lending_trx
   - 持久化: pgdata卷

2. **server** - HTTP API服务 + 定时任务
   - 端口: 8080
   - 功能: Webhook处理、委托逻辑、API接口
   - 依赖: postgres

3. **bot** - Telegram Bot服务
   - 功能: 监控、告警、状态查询
   - 依赖: server (通过内部网络访问API)

## 🛠️ 管理命令

### 启动服务
```bash
# 启动所有服务
make docker-compose-up

# 启动特定服务
docker-compose up -d postgres
docker-compose up -d server
docker-compose up -d bot
```

### 停止服务
```bash
# 停止所有服务
make docker-compose-down

# 停止特定服务
docker-compose stop server
docker-compose stop bot
```

### 重启服务
```bash
# 重启所有服务
make docker-compose-restart

# 重启特定服务
docker-compose restart server
docker-compose restart bot
```

### 查看日志
```bash
# 查看所有服务日志
make docker-compose-logs

# 查看特定服务日志
docker-compose logs -f server
docker-compose logs -f bot
docker-compose logs -f postgres
```

### 进入容器
```bash
# 进入server容器
docker-compose exec server sh

# 进入bot容器
docker-compose exec bot sh

# 进入数据库容器
docker-compose exec postgres psql -U postgres -d lending_trx
```

## 🔧 配置说明

### 环境变量

#### Server服务环境变量
- `DATABASE_URL` - 数据库连接字符串
- `TRON_API_URL` - TRON API地址
- `TRON_API_KEY` - TRON API密钥
- `DELEGATION_FROM_ADDRESS` - 委托方地址
- `PORT` - HTTP服务端口
- `LOG_LEVEL` - 日志级别
- `CRON_SCHEDULE` - 定时任务间隔
- `DELEGATION_BASE` - 委托基础数量
- `MIN_DELEGATION_AMOUNT` - 最小委托数量

#### Bot服务环境变量
- `TELEGRAM_BOT_TOKEN` - Telegram Bot令牌
- `API_BASE_URL` - API服务器地址
- `MONITOR_INTERVAL_MINUTES` - 监控间隔
- `HTTP_TIMEOUT` - HTTP超时时间
- `LONG_POLLING_TIMEOUT` - 长轮询超时
- `MAX_RETRIES` - 最大重试次数
- `RETRY_DELAY` - 重试延迟

### 网络配置

- **内部网络**: 服务间通过Docker网络通信
- **外部端口**: 只有PostgreSQL和Server暴露端口
- **Bot服务**: 仅通过内部网络访问Server API

## 📊 监控和调试

### 健康检查
```bash
# 检查数据库连接
docker-compose exec server ./lending-trx server --help

# 检查Bot状态
docker-compose exec bot ./lending-trx bot --help

# 检查API接口
curl http://localhost:8080/api/delegation-account
```

### 数据持久化
```bash
# 查看数据库数据
docker-compose exec postgres psql -U postgres -d lending_trx -c "SELECT * FROM webhook_data LIMIT 5;"

# 备份数据库
docker-compose exec postgres pg_dump -U postgres lending_trx > backup.sql

# 恢复数据库
docker-compose exec -T postgres psql -U postgres lending_trx < backup.sql
```

## 🚨 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose ps postgres
   
   # 查看数据库日志
   docker-compose logs postgres
   ```

2. **Server启动失败**
   ```bash
   # 检查环境变量
   docker-compose exec server env | grep DATABASE_URL
   
   # 查看Server日志
   docker-compose logs server
   ```

3. **Bot连接失败**
   ```bash
   # 检查Bot配置
   docker-compose exec bot env | grep TELEGRAM_BOT_TOKEN
   
   # 查看Bot日志
   docker-compose logs bot
   ```

4. **API接口无响应**
   ```bash
   # 检查Server是否运行
   curl http://localhost:8080/api/delegation-account
   
   # 检查网络连接
   docker-compose exec bot curl http://server:8080/api/delegation-account
   ```

### 日志分析
```bash
# 查看错误日志
docker-compose logs | grep ERROR

# 查看特定时间段的日志
docker-compose logs --since="2024-01-01T00:00:00" server

# 实时监控日志
docker-compose logs -f --tail=100
```

## 🔄 更新部署

### 更新代码
```bash
# 停止服务
make docker-compose-down

# 重新构建镜像
make docker-build

# 启动服务
make docker-compose-up
```

### 更新配置
```bash
# 修改环境变量后重启服务
docker-compose restart server
docker-compose restart bot
```

## 📝 注意事项

1. **数据持久化**: 数据库数据存储在Docker卷中，容器重启不会丢失数据
2. **网络隔离**: Bot服务通过内部网络访问Server，确保安全性
3. **资源限制**: 生产环境建议设置容器资源限制
4. **日志管理**: 日志文件存储在宿主机的logs目录中
5. **环境变量**: 敏感信息通过.env文件管理，不要提交到版本控制

## 🎯 生产环境建议

1. **使用外部数据库**: 生产环境建议使用独立的PostgreSQL服务
2. **配置反向代理**: 使用Nginx等反向代理管理HTTP流量
3. **监控告警**: 集成Prometheus、Grafana等监控系统
4. **日志聚合**: 使用ELK Stack等日志聚合系统
5. **备份策略**: 定期备份数据库和配置文件 