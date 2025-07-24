package telegram_bot

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ 未找到.env文件，使用系统环境变量")
	}

	config := &Config{
		TelegramToken:       getEnv("TELEGRAM_BOT_TOKEN", ""),
		APIBaseURL:          getEnv("API_BASE_URL", "http://localhost:8080"),
		MonitorIntervalMins: getEnvAsInt("MONITOR_INTERVAL_MINUTES", 5),
		HTTPTimeout:         getEnvAsDuration("HTTP_TIMEOUT", 30*time.Second),
		LongPollingTimeout:  getEnvAsInt("LONG_POLLING_TIMEOUT", 30),
		MaxRetries:          getEnvAsInt("MAX_RETRIES", 3),
		RetryDelay:          getEnvAsDuration("RETRY_DELAY", 5*time.Second),
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("⚠️ 环境变量 %s 不是有效的整数，使用默认值 %d", key, defaultValue)
	}
	return defaultValue
}

// getEnvAsDuration 获取环境变量并转换为时间间隔
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("⚠️ 环境变量 %s 不是有效的时间间隔，使用默认值 %v", key, defaultValue)
	}
	return defaultValue
}
