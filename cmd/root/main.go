package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lending-trx",
	Short: "TRX委托服务 - 管理TRX能量委托的Go应用程序",
	Long: `TRX委托服务是一个用于管理TRX能量委托的Go应用程序，
包含HTTP API服务、定时任务处理和Telegram Bot监控功能。

支持以下功能：
- HTTP API服务：提供委托账户查询接口
- 定时任务：自动处理webhook数据和能量委托
- Telegram Bot：实时监控和告警通知`,
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(botCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 执行命令失败: %v\n", err)
		os.Exit(1)
	}
}
