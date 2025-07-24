package telegram_bot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TelegramBot Telegram Bot结构体
type TelegramBot struct {
	Token      string
	APIBaseURL string
	HTTPClient *http.Client
	Offset     int64
	mu         sync.Mutex

	// 监控相关
	monitoring     map[int64]bool // chatID -> isMonitoring
	monitoringMu   sync.RWMutex
	stopMonitoring map[int64]chan struct{} // chatID -> stop channel
}

// Config 配置结构体
type Config struct {
	TelegramToken       string
	APIBaseURL          string
	MonitorIntervalMins int
	HTTPTimeout         time.Duration
	LongPollingTimeout  int
	MaxRetries          int
	RetryDelay          time.Duration
}

// Update Telegram更新消息结构体
type Update struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

// Message Telegram消息结构体
type Message struct {
	MessageID int64  `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

// User Telegram用户结构体
type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Chat Telegram聊天结构体
type Chat struct {
	ID    int64  `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

// SendMessageRequest 发送消息请求结构体
type SendMessageRequest struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// SendMessageResponse 发送消息响应结构体
type SendMessageResponse struct {
	OK     bool    `json:"ok"`
	Result Message `json:"result"`
}

// DelegationAccountInfo 委托账户信息结构体
type DelegationAccountInfo struct {
	Status string `json:"status"`
	Data   struct {
		Address     string `json:"address"`
		Balance     string `json:"balance"`
		Energy      string `json:"energy"`
		EnergyLimit string `json:"energy_limit"`
		EnergyUsed  string `json:"energy_used"`
	} `json:"data"`
}

// TelegramResponse Telegram API响应结构体
type TelegramResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// NewTelegramBot 创建新的Telegram Bot实例
func NewTelegramBot(config *Config) *TelegramBot {
	return &TelegramBot{
		Token:      config.TelegramToken,
		APIBaseURL: "https://api.telegram.org/bot" + config.TelegramToken,
		HTTPClient: &http.Client{
			Timeout: config.HTTPTimeout,
		},
		Offset:         0,
		monitoring:     make(map[int64]bool),
		stopMonitoring: make(map[int64]chan struct{}),
	}
}

// GetUpdates 获取更新消息
func (bot *TelegramBot) GetUpdates() ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=30", bot.APIBaseURL, bot.Offset)

	resp, err := bot.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response TelegramResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("telegram API error: %s", string(body))
	}

	return response.Result, nil
}

// SendMessage 发送消息
func (bot *TelegramBot) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s/sendMessage", bot.APIBaseURL)

	request := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := bot.HTTPClient.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var response SendMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.OK {
		return fmt.Errorf("telegram API error: %s", string(body))
	}

	return nil
}

// GetDelegationAccountInfo 获取委托账户信息
func (bot *TelegramBot) GetDelegationAccountInfo(apiURL string) (*DelegationAccountInfo, error) {
	url := fmt.Sprintf("%s/api/delegation-account", apiURL)

	resp, err := bot.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var accountInfo DelegationAccountInfo
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &accountInfo, nil
}

// FormatAccountInfo 格式化账户信息
func (bot *TelegramBot) FormatAccountInfo(info *DelegationAccountInfo) string {
	if info.Status != "ok" {
		return "❌ 获取账户信息失败"
	}

	// 解析余额 (SUN -> TRX)
	balanceSUN, err := strconv.ParseInt(info.Data.Balance, 10, 64)
	if err != nil {
		return "❌ 解析余额失败"
	}
	balanceTRX := float64(balanceSUN) / 1000000.0

	// 解析能量信息
	energy, err := strconv.ParseInt(info.Data.Energy, 10, 64)
	if err != nil {
		return "❌ 解析能量信息失败"
	}

	energyLimit, err := strconv.ParseInt(info.Data.EnergyLimit, 10, 64)
	if err != nil {
		return "❌ 解析能量限制失败"
	}

	energyUsed, err := strconv.ParseInt(info.Data.EnergyUsed, 10, 64)
	if err != nil {
		return "❌ 解析已用能量失败"
	}

	// 计算使用率
	var usagePercent float64
	if energyLimit > 0 {
		usagePercent = float64(energyUsed) / float64(energyLimit) * 100
	}

	// 构建消息
	message := fmt.Sprintf(`📊 <b>委托账户状态报告</b>

🏦 <b>账户地址:</b>
<code>%s</code>

💰 <b>账户余额:</b>
%.6f TRX

⚡ <b>能量状态:</b>
• 可用能量: %d
• 能量限制: %d
• 已用能量: %d
• 使用率: %.2f%%

📈 <b>状态:</b>
• 余额状态: %s
• 能量状态: %s`,
		info.Data.Address,
		balanceTRX,
		energy,
		energyLimit,
		energyUsed,
		usagePercent,
		getBalanceStatus(balanceTRX),
		getEnergyStatus(energy),
	)

	// 添加告警信息
	alerts := getAlerts(balanceTRX, energy)
	if alerts != "" {
		message += "\n\n🚨 <b>告警信息:</b>\n" + alerts
	}

	return message
}

// getBalanceStatus 获取余额状态
func getBalanceStatus(balanceTRX float64) string {
	if balanceTRX < 10 {
		return "⚠️ 余额不足"
	}
	return "✅ 余额充足"
}

// getEnergyStatus 获取能量状态
func getEnergyStatus(energy int64) string {
	if energy < 1000 {
		return "⚠️ 能量不足"
	}
	return "✅ 能量充足"
}

// getAlerts 获取告警信息
func getAlerts(balanceTRX float64, energy int64) string {
	var alerts []string

	if balanceTRX < 10 {
		alerts = append(alerts, "• 账户余额不足 (少于10 TRX)")
	}

	if energy < 1000 {
		alerts = append(alerts, "• 可用能量不足 (少于1000)")
	}

	return strings.Join(alerts, "\n")
}

// StartMonitoring 开始监控指定聊天
func (bot *TelegramBot) StartMonitoring(chatID int64, apiURL string, intervalMins int) {
	bot.monitoringMu.Lock()
	defer bot.monitoringMu.Unlock()

	// 如果已经在监控，先停止
	if bot.monitoring[chatID] {
		close(bot.stopMonitoring[chatID])
	}

	// 创建停止通道
	stopChan := make(chan struct{})
	bot.stopMonitoring[chatID] = stopChan
	bot.monitoring[chatID] = true

	// 启动监控goroutine
	go func() {
		ticker := time.NewTicker(time.Duration(intervalMins) * time.Minute)
		defer ticker.Stop()

		log.Printf("🔄 开始监控聊天 %d (间隔: %d分钟)", chatID, intervalMins)

		// 立即执行一次检查
		bot.sendAccountStatus(chatID, apiURL)

		for {
			select {
			case <-ticker.C:
				bot.sendAccountStatus(chatID, apiURL)
			case <-stopChan:
				log.Printf("⏹️ 停止监控聊天 %d", chatID)
				return
			}
		}
	}()
}

// StopMonitoring 停止监控指定聊天
func (bot *TelegramBot) StopMonitoring(chatID int64) {
	bot.monitoringMu.Lock()
	defer bot.monitoringMu.Unlock()

	if bot.monitoring[chatID] {
		close(bot.stopMonitoring[chatID])
		delete(bot.monitoring, chatID)
		delete(bot.stopMonitoring, chatID)
	}
}

// IsMonitoring 检查是否正在监控
func (bot *TelegramBot) IsMonitoring(chatID int64) bool {
	bot.monitoringMu.RLock()
	defer bot.monitoringMu.RUnlock()
	return bot.monitoring[chatID]
}

// sendAccountStatus 发送账户状态
func (bot *TelegramBot) sendAccountStatus(chatID int64, apiURL string) {
	accountInfo, err := bot.GetDelegationAccountInfo(apiURL)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ 监控查询失败: %v", err)
		bot.SendMessage(chatID, errorMsg)
		return
	}

	formattedInfo := bot.FormatAccountInfo(accountInfo)
	bot.SendMessage(chatID, formattedInfo)
}

// StartLongPolling 开始Long Polling
func (bot *TelegramBot) StartLongPolling(config *Config, commandHandler *CommandHandler) {
	log.Println("🤖 Telegram Bot 启动中...")
	log.Printf("📡 API地址: %s", config.APIBaseURL)
	log.Printf("⏱️ 监控间隔: %d分钟", config.MonitorIntervalMins)
	log.Println("⏳ 开始Long Polling...")

	for {
		updates, err := bot.GetUpdates()
		if err != nil {
			log.Printf("❌ 获取更新失败: %v", err)
			time.Sleep(config.RetryDelay)
			continue
		}

		for _, update := range updates {
			// 更新offset
			bot.mu.Lock()
			if update.UpdateID >= bot.Offset {
				bot.Offset = update.UpdateID + 1
			}
			bot.mu.Unlock()

			// 处理消息
			if update.Message.Text != "" {
				log.Printf("📨 收到消息: %s (来自: %s)", update.Message.Text, update.Message.From.Username)
				commandHandler.HandleCommand(update.Message)
			}
		}

		// 短暂等待
		time.Sleep(1 * time.Second)
	}
}
