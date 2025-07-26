package cronjob

import (
	"context"
	"fmt"
	"lending-trx/internal/db"
	"lending-trx/internal/tron"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
	"github.com/sunjiangjun/xlog"
)

// CronJob 定时任务结构体
type CronJob struct {
	ctx        context.Context
	pool       *pgxpool.Pool
	log        *xlog.XLog
	tronClient *tron.TronClient
}

// NewCronJob 创建新的定时任务实例
func NewCronJob(ctx context.Context, pool *pgxpool.Pool, log *xlog.XLog) *CronJob {
	// 从环境变量获取 Tron API 配置
	baseURL := os.Getenv("TRON_API_URL")
	if baseURL == "" {
		baseURL = "https://api.trongrid.io"
	}

	apiKey := os.Getenv("TRON_API_KEY")

	tronClient := tron.NewTronClient(baseURL, apiKey)

	return &CronJob{
		ctx:        ctx,
		pool:       pool,
		log:        log,
		tronClient: tronClient,
	}
}

// StartCron 启动定时任务
func StartCron(ctx context.Context, pool *pgxpool.Pool, log *xlog.XLog) {
	job := NewCronJob(ctx, pool, log)
	job.start()
}

// start 启动定时任务
func (c *CronJob) start() {
	cronScheduler := cron.New()

	// 从环境变量获取定时任务配置
	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule == "" {
		cronSchedule = "@every 30s"
	}

	_, err := cronScheduler.AddFunc(cronSchedule, c.processWebhookData)
	if err != nil {
		c.log.Error("Failed to add cron job", err)
		return
	}

	c.log.Info("Cron job started", "schedule", cronSchedule)
	go cronScheduler.Run()
}

// processWebhookData 处理webhook数据的主函数
func (c *CronJob) processWebhookData() {
	c.log.Info("Starting to process webhook data")

	// 获取统计信息
	stats, err := db.GetWebhookDataStats(c.ctx, c.pool)
	if err != nil {
		c.log.Error("Failed to get statistics", err)
	} else {
		c.log.Info("Data statistics", "stats", stats)
	}

	// 处理待处理的数据 (status=0)
	pendingData, err := db.QueryPendingWebhookData(c.ctx, c.pool)
	if err != nil {
		c.log.Error("Failed to query pending data", err)
	} else if len(pendingData) > 0 {
		c.processPendingData(pendingData)
	}

	// 处理已过期且已授权的数据 (status=2)
	expiredData, err := db.QueryExpiredWebhookData(c.ctx, c.pool)
	if err != nil {
		c.log.Error("Failed to query expired data", err)
	} else if len(expiredData) > 0 {
		c.processExpiredData(expiredData)
	}

	c.log.Info("Webhook data processing completed")
}

// processPendingData 处理待处理的数据
func (c *CronJob) processPendingData(data []*db.WebhookDataModel) {
	c.log.Info("Processing pending data", "count", len(data))

	for _, item := range data {
		c.log.Info("Processing pending item",
			"id", item.ID,
			"tx_hash", item.TxHash,
			"from", item.FromAddress,
			"to", item.ToAddress,
			"value", item.Value,
		)

		// TODO: 执行能量委托
		err := c.executeEnergyDelegation(item)
		if err != nil {
			c.log.Error("Failed to execute energy delegation", err, "id", item.ID)
			continue
		}

		// 更新状态为处理中 (status=2)
		if err := db.UpdateWebhookStatusByID(c.ctx, c.pool, item.ID, 2); err != nil {
			c.log.Error("Failed to update status", err, "id", item.ID)
		} else {
			c.log.Info("Status updated successfully", "id", item.ID, "status", 2)
		}
	}
}

// processExpiredData 处理已过期的数据
func (c *CronJob) processExpiredData(data []*db.WebhookDataModel) {
	c.log.Info("Processing expired data", "count", len(data))

	for _, item := range data {
		c.log.Info("Processing expired item",
			"id", item.ID,
			"tx_hash", item.TxHash,
			"expire_time", item.ExpireTime,
		)

		// TODO: 取消能量委托
		err := c.cancelEnergyDelegation(item)
		if err != nil {
			c.log.Error("Failed to cancel energy delegation", err, "id", item.ID)
			continue
		}

		// 更新状态为已回收 (status=3)
		if err := db.UpdateWebhookStatusByID(c.ctx, c.pool, item.ID, 3); err != nil {
			c.log.Error("Failed to update status", err, "id", item.ID)
		} else {
			c.log.Info("Status updated successfully", "id", item.ID, "status", 3)
		}
	}
}

// TODO: 实现具体的业务逻辑函数
// executeEnergyDelegation 执行能量委托
func (c *CronJob) executeEnergyDelegation(data *db.WebhookDataModel) error {
	c.log.Info("Starting energy delegation",
		"id", data.ID,
		"from", data.FromAddress,
		"to", data.ToAddress,
		"value", data.Value,
		"tx_hash", data.TxHash,
	)

	// 金额<1trx

	valueInt, err := strconv.ParseInt(data.Value, 10, 64)
	if err != nil {
		c.log.Error("Failed to parse transaction amount", err, "value", data.Value)
		return fmt.Errorf("failed to parse transaction amount: %w", err)
	}

	if valueInt < 1000000 { // 小于 1 TRX (1,000,000 SUN)
		c.log.Info("Transaction amount is less than 1 TRX, skipping energy delegation", "id", data.ID, "value", data.Value)
		return nil
	}

	// 从环境变量获取统一的委托方地址
	delegationFromAddress := os.Getenv("DELEGATION_FROM_ADDRESS")
	if delegationFromAddress == "" {
		return fmt.Errorf("environment variable DELEGATION_FROM_ADDRESS not set")
	}

	c.log.Info("Using unified delegation address",
		"original_from", data.FromAddress,
		"delegation_from", delegationFromAddress,
	)

	// 1. 验证委托方账户信息
	accountInfo, err := c.tronClient.GetAccountInfo(c.ctx, delegationFromAddress)
	if err != nil {
		return fmt.Errorf("failed to get delegation account info: %w", err)
	}

	c.log.Info("Delegation account info",
		"address", accountInfo.Address,
		"balance", accountInfo.Balance,
		"energy", accountInfo.Energy,
		"energy_limit", accountInfo.EnergyLimit,
		"energy_used", accountInfo.EnergyUsed,
	)

	// 2. 计算可委托的能量数量
	// 这里可以根据业务逻辑计算委托数量
	// 例如：根据交易金额的百分比，或者固定数量
	delegationAmount := c.calculateDelegationAmount(data.Value, accountInfo.Energy)

	if delegationAmount == "0" {
		return fmt.Errorf("no energy available for delegation")
	}

	c.log.Info("Calculated delegation amount",
		"original_value", data.Value,
		"available_energy", accountInfo.Energy,
		"delegation_amount", delegationAmount,
	)

	// 3. 构建委托请求
	delegationReq := &tron.EnergyDelegationRequest{
		FromAddress: delegationFromAddress, // 使用统一的委托方地址
		ToAddress:   data.FromAddress,      // 委托给交易发起方
		Amount:      delegationAmount,
		TxHash:      data.TxHash,
		BlockHeight: data.BlockHeight,
	}

	// 4. 执行能量委托
	delegationResp, err := c.tronClient.DelegateEnergy(c.ctx, delegationReq)
	if err != nil {
		return fmt.Errorf("energy delegation API call failed: %w", err)
	}

	c.log.Info("Energy delegation successful",
		"tx_id", delegationResp.TxID,
		"message", delegationResp.Message,
		"from", delegationFromAddress,
		"to", data.ToAddress,
		"amount", delegationAmount,
	)

	// 5. 保存原始委托交易ID到数据库
	if delegationResp.TxID != "" {
		err = db.UpdateOriginalTxIDByID(c.ctx, c.pool, data.ID, delegationResp.TxID)
		if err != nil {
			c.log.Error("Failed to save original delegation transaction ID", err, "id", data.ID, "tx_id", delegationResp.TxID)
			// 不返回错误，因为委托已经成功
		} else {
			c.log.Info("Original delegation transaction ID saved", "id", data.ID, "original_tx_id", delegationResp.TxID)
		}
	}

	return nil
}

// cancelEnergyDelegation 取消能量委托
func (c *CronJob) cancelEnergyDelegation(data *db.WebhookDataModel) error {
	c.log.Info("Starting energy delegation cancellation",
		"id", data.ID,
		"from", data.FromAddress,
		"to", data.ToAddress,
		"tx_hash", data.TxHash,
	)

	// 从环境变量获取统一的委托方地址
	delegationFromAddress := os.Getenv("DELEGATION_FROM_ADDRESS")
	if delegationFromAddress == "" {
		return fmt.Errorf("environment variable DELEGATION_FROM_ADDRESS not set")
	}

	c.log.Info("Using unified delegation address for cancellation",
		"original_from", data.FromAddress,
		"delegation_from", delegationFromAddress,
	)

	// 1. 获取原始委托交易ID
	originalTxID, err := db.GetOriginalTxIDByID(c.ctx, c.pool, data.ID)
	if err != nil {
		return fmt.Errorf("failed to get original delegation transaction ID: %w", err)
	}

	if originalTxID == "" {
		return fmt.Errorf("original delegation transaction ID is empty, cannot cancel delegation")
	}

	c.log.Info("Retrieved original delegation transaction ID", "id", data.ID, "original_tx_id", originalTxID)

	// 2. 构建取消委托请求
	cancelReq := &tron.CancelDelegationRequest{
		FromAddress:  delegationFromAddress, // 使用统一的委托方地址
		ToAddress:    data.FromAddress,      // 取消委托给交易发起方
		OriginalTxID: originalTxID,
		TxHash:       data.TxHash,
	}

	// 3. 执行取消委托
	cancelResp, err := c.tronClient.CancelEnergyDelegation(c.ctx, cancelReq)
	if err != nil {
		return fmt.Errorf("cancel energy delegation API call failed: %w", err)
	}

	c.log.Info("Energy delegation cancellation successful",
		"tx_id", cancelResp.TxID,
		"message", cancelResp.Message,
		"from", delegationFromAddress,
		"to", data.ToAddress,
	)

	return nil
}

// calculateDelegationAmount 计算委托数量
func (c *CronJob) calculateDelegationAmount(value string, availableEnergy string) string {
	// 将字符串转换为数值进行计算
	valueInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		c.log.Error("Failed to parse transaction amount", err, "value", value)
		return "0"
	}

	energyInt, err := strconv.ParseInt(availableEnergy, 10, 64)
	if err != nil {
		c.log.Error("Failed to parse available energy", err, "energy", availableEnergy)
		return "0"
	}

	// 从环境变量获取委托配置
	delegationBaseStr := os.Getenv("DELEGATION_BASE")
	if delegationBaseStr == "" {
		delegationBaseStr = "65000"
	}
	delegationBase, err := strconv.ParseInt(delegationBaseStr, 10, 64)
	if err != nil {
		c.log.Error("Failed to parse DELEGATION_BASE", err, "value", delegationBaseStr)
		delegationBase = 65000
	}

	minDelegationStr := os.Getenv("MIN_DELEGATION_AMOUNT")
	if minDelegationStr == "" {
		minDelegationStr = "1000000"
	}
	minDelegation, err := strconv.ParseInt(minDelegationStr, 10, 64)
	if err != nil {
		c.log.Error("Failed to parse MIN_DELEGATION_AMOUNT", err, "value", minDelegationStr)
		minDelegation = 1000000
	}

	// 业务逻辑：基于交易金额计算委托数量
	// 1 TRX = 1,000,000 SUN
	var delegationAmount int64

	// 将 SUN 转换为 TRX 进行计算
	trxValue := valueInt / 1000000 // 1 TRX = 1,000,000 SUN

	if trxValue == 1 {
		// 1 TRX → 委托 delegationBase
		delegationAmount = delegationBase
	} else if trxValue == 2 {
		// 2 TRX → 委托 2 * delegationBase
		delegationAmount = 2 * delegationBase
	} else {
		// 其他情况，不进行委托
		delegationAmount = 0
	}

	// 确保不超过可用能量
	if delegationAmount > energyInt {
		delegationAmount = energyInt
	}

	// 确保最小委托数量
	if delegationAmount < minDelegation {
		delegationAmount = 0
	}

	return strconv.FormatInt(delegationAmount, 10)
}
