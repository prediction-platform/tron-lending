package telegram_bot

import (
	"fmt"
	"log"
	"strings"
)

// CommandHandler 命令处理器
type CommandHandler struct {
	bot    *TelegramBot
	config *Config
}

// NewCommandHandler 创建新的命令处理器
func NewCommandHandler(bot *TelegramBot, config *Config) *CommandHandler {
	return &CommandHandler{
		bot:    bot,
		config: config,
	}
}

// HandleCommand 处理命令
func (h *CommandHandler) HandleCommand(message Message) {
	command := strings.ToLower(strings.TrimSpace(message.Text))

	switch {
	case strings.HasPrefix(command, "/start"):
		h.handleStart(message)
	case strings.HasPrefix(command, "/help"):
		h.handleHelp(message)
	case strings.HasPrefix(command, "/status"):
		h.handleStatus(message)
	case strings.HasPrefix(command, "/monitor"):
		h.handleMonitor(message)
	case strings.HasPrefix(command, "/stop"):
		h.handleStop(message)
	default:
		h.handleUnknownCommand(message)
	}
}

// handleStart 处理启动命令
func (h *CommandHandler) handleStart(message Message) {
	response := `🤖 <b>委托账户监控Bot</b>

欢迎使用委托账户监控服务！

<b>可用命令:</b>
• /status - 查询委托账户状态
• /help - 显示帮助信息
• /monitor - 开始持续监控
• /stop - 停止监控

使用 /help 查看更多信息。`

	if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
		log.Printf("❌ 发送启动消息失败: %v", err)
	}
}

// handleHelp 处理帮助命令
func (h *CommandHandler) handleHelp(message Message) {
	response := fmt.Sprintf(`📖 <b>帮助信息</b>

<b>命令说明:</b>
• /start - 启动Bot
• /status - 查询当前委托账户状态
• /help - 显示此帮助信息
• /monitor - 开始持续监控（每%d分钟）
• /stop - 停止持续监控

<b>告警阈值:</b>
• 余额少于10 TRX时告警
• 可用能量少于1000时告警

<b>API地址:</b>
%s/api/delegation-account`, h.config.MonitorIntervalMins, h.config.APIBaseURL)

	if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
		log.Printf("❌ 发送帮助消息失败: %v", err)
	}
}

// handleStatus 处理状态查询命令
func (h *CommandHandler) handleStatus(message Message) {
	// 发送查询中消息
	if err := h.bot.SendMessage(message.Chat.ID, "🔍 正在查询委托账户状态..."); err != nil {
		log.Printf("❌ 发送查询中消息失败: %v", err)
		return
	}

	// 获取账户信息
	accountInfo, err := h.bot.GetDelegationAccountInfo(h.config.APIBaseURL)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ 查询失败: %v", err)
		h.bot.SendMessage(message.Chat.ID, errorMsg)
		return
	}

	// 格式化并发送结果
	formattedInfo := h.bot.FormatAccountInfo(accountInfo)
	if err := h.bot.SendMessage(message.Chat.ID, formattedInfo); err != nil {
		log.Printf("❌ 发送状态信息失败: %v", err)
	}
}

// handleMonitor 处理监控命令
func (h *CommandHandler) handleMonitor(message Message) {
	if h.bot.IsMonitoring(message.Chat.ID) {
		response := fmt.Sprintf(`🔄 已经在监控中

Bot正在每%d分钟自动查询委托账户状态。
使用 /stop 命令停止监控。`, h.config.MonitorIntervalMins)

		if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
			log.Printf("❌ 发送监控状态消息失败: %v", err)
		}
	} else {
		response := fmt.Sprintf(`🔄 开始持续监控模式

Bot将每%d分钟自动查询一次委托账户状态。
使用 /stop 命令停止监控。`, h.config.MonitorIntervalMins)

		if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
			log.Printf("❌ 发送开始监控消息失败: %v", err)
			return
		}

		// 启动监控
		h.bot.StartMonitoring(message.Chat.ID, h.config.APIBaseURL, h.config.MonitorIntervalMins)
	}
}

// handleStop 处理停止监控命令
func (h *CommandHandler) handleStop(message Message) {
	if h.bot.IsMonitoring(message.Chat.ID) {
		h.bot.StopMonitoring(message.Chat.ID)
		response := `⏹️ 停止监控模式

已停止持续监控。使用 /monitor 重新开始监控。`

		if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
			log.Printf("❌ 发送停止监控消息失败: %v", err)
		}
	} else {
		response := `⏹️ 未在监控中

当前没有进行持续监控。使用 /monitor 开始监控。`

		if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
			log.Printf("❌ 发送未监控消息失败: %v", err)
		}
	}
}

// handleUnknownCommand 处理未知命令
func (h *CommandHandler) handleUnknownCommand(message Message) {
	response := `❓ 未知命令

使用 /help 查看可用命令。`

	if err := h.bot.SendMessage(message.Chat.ID, response); err != nil {
		log.Printf("❌ 发送未知命令消息失败: %v", err)
	}
}
