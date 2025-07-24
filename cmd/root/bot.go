package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"lending-trx/pkg/telegram_bot"
)

var (
	botToken string
	apiURL   string
	interval int
	botCmd   = &cobra.Command{
		Use:   "bot",
		Short: "启动Telegram Bot服务",
		Long: `启动Telegram Bot监控服务，包括：
- 实时查询委托账户状态
- 持续监控功能
- 告警通知
- 支持多用户监控`,
		Run: runBot,
	}
)

func init() {
	botCmd.Flags().StringVarP(&botToken, "token", "t", "", "Telegram Bot Token")
	botCmd.Flags().StringVarP(&apiURL, "api", "a", "http://localhost:8080", "API服务器地址")
	botCmd.Flags().IntVarP(&interval, "interval", "i", 5, "监控间隔(分钟)")
}

func runBot(cmd *cobra.Command, args []string) {
	fmt.Println("🤖 启动Telegram Bot...")

	// 加载配置
	config := telegram_bot.LoadConfig()

	// 如果命令行参数提供了值，则覆盖配置
	if botToken != "" {
		config.TelegramToken = botToken
	}
	if apiURL != "" {
		config.APIBaseURL = apiURL
	}
	if interval > 0 {
		config.MonitorIntervalMins = interval
	}

	// 验证必需配置
	if config.TelegramToken == "" {
		log.Fatal("❌ TELEGRAM_BOT_TOKEN 环境变量未设置或未通过 --token 参数提供")
	}

	// 创建Bot实例
	bot := telegram_bot.NewTelegramBot(config)
	commandHandler := telegram_bot.NewCommandHandler(bot, config)

	fmt.Printf("✅ Telegram Bot已启动\n")
	fmt.Printf("📡 API地址: %s\n", config.APIBaseURL)
	fmt.Printf("⏱️ 监控间隔: %d分钟\n", config.MonitorIntervalMins)
	fmt.Printf("🤖 Bot Token: %s...\n", config.TelegramToken[:10])

	// 启动Long Polling
	bot.StartLongPolling(config, commandHandler)
}
