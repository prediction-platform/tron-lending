# 定时任务模块 (CronJob)

## 概述

本模块提供了优化的定时任务功能，用于处理webhook数据的状态管理和业务逻辑执行。

## 架构设计

### 核心组件

1. **CronJob 结构体**: 封装定时任务的核心逻辑
2. **模块化函数**: 将复杂逻辑分解为小的、可测试的函数
3. **错误处理**: 完善的错误处理和日志记录

### 代码结构

```
CronJob
├── NewCronJob()           # 创建定时任务实例
├── start()                # 启动定时任务
├── processWebhookData()   # 主处理函数
├── queryPendingData()     # 查询待处理数据
├── queryExpiredData()     # 查询过期数据
├── processPendingData()   # 处理待处理数据
├── processExpiredData()   # 处理过期数据
├── updateStatus()         # 更新数据状态
└── scanRows()            # 扫描数据库行
```

## 主要功能

### 1. 数据状态管理

| 状态 | 说明 | 处理逻辑 |
|------|------|----------|
| 0 | 初始化 | 执行能量委托，更新为状态1 |
| 1 | 执行中 | 等待执行完成 |
| 2 | 已授权 | 检查是否过期 |
| 3 | 已回收 | 最终状态 |

### 2. 定时处理流程

```go
// 每30秒执行一次
func (c *CronJob) processWebhookData() {
    // 1. 查询并处理待处理数据 (status=0)
    pendingData := c.queryPendingData()
    if len(pendingData) > 0 {
        c.processPendingData(pendingData)
    }
    
    // 2. 查询并处理过期数据 (status=2)
    expiredData := c.queryExpiredData()
    if len(expiredData) > 0 {
        c.processExpiredData(expiredData)
    }
}
```

### 3. 数据库操作

#### 查询待处理数据
```sql
SELECT id, block_height, tx_hash, from_address, to_address, value, 
       block_time, create_time, update_time, expire_time, status 
FROM webhook_data 
WHERE status=0
```

#### 查询过期数据
```sql
SELECT id, block_height, tx_hash, from_address, to_address, value, 
       block_time, create_time, update_time, expire_time, status 
FROM webhook_data 
WHERE status=2 AND expire_time < $1
```

#### 更新状态
```sql
UPDATE webhook_data 
SET status = $1, update_time = NOW() 
WHERE id = $2
```

## 优化亮点

### 1. 模块化设计
- **单一职责**: 每个函数只负责一个特定功能
- **可测试性**: 函数独立，易于单元测试
- **可维护性**: 代码结构清晰，易于理解和修改

### 2. 错误处理
- **优雅降级**: 单个项目失败不影响其他项目处理
- **详细日志**: 记录每个步骤的执行情况
- **状态追踪**: 记录处理的项目数量和状态变化

### 3. 性能优化
- **批量查询**: 一次性查询所有相关数据
- **连接复用**: 使用连接池管理数据库连接
- **资源管理**: 正确关闭数据库连接和事务

### 4. 可扩展性
- **插件化设计**: 业务逻辑函数可以独立实现
- **配置化**: 可以轻松调整查询条件和处理逻辑
- **监控友好**: 提供详细的处理统计信息

## 使用示例

```go
// 启动定时任务
func main() {
    ctx := context.Background()
    pool, err := db.InitDB(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    log := xlog.NewXLog()
    cronjob.StartCron(ctx, pool, log)
    
    // 保持程序运行
    select {}
}
```

## 业务逻辑扩展

### 实现能量委托
```go
func (c *CronJob) executeEnergyDelegation(data db.WebhookDataModel) error {
    // 1. 验证交易数据
    // 2. 调用区块链API
    // 3. 处理响应
    // 4. 更新状态
    return nil
}
```

### 实现取消委托
```go
func (c *CronJob) cancelEnergyDelegation(data db.WebhookDataModel) error {
    // 1. 验证当前状态
    // 2. 调用取消API
    // 3. 处理响应
    // 4. 更新状态
    return nil
}
```

## 监控和日志

### 日志级别
- **Info**: 正常处理流程
- **Error**: 错误和异常情况
- **Debug**: 详细的调试信息

### 关键指标
- 待处理数据数量
- 过期数据数量
- 处理成功率
- 平均处理时间

## 配置选项

### 定时器配置
```go
// 当前配置：每30秒执行一次
"@every 30s"

// 可选配置：
"@every 1m"    // 每分钟
"@every 5m"    // 每5分钟
"0 */5 * * *"  // 每5分钟（cron格式）
```

### 数据库配置
- 连接池大小
- 查询超时时间
- 事务隔离级别

## 注意事项

1. **并发安全**: 确保数据库操作的线程安全
2. **资源清理**: 正确关闭数据库连接和事务
3. **错误恢复**: 实现错误重试和恢复机制
4. **性能监控**: 监控处理性能和资源使用情况
5. **数据一致性**: 确保状态更新的原子性 