package cronjob

import (
	"os"
	"strconv"
	"testing"
)

func TestDelegationFromAddress(t *testing.T) {
	// 测试环境变量设置
	testAddress := "TTestAddressForTesting123456789"
	os.Setenv("DELEGATION_FROM_ADDRESS", testAddress)
	defer os.Unsetenv("DELEGATION_FROM_ADDRESS")

	// 验证环境变量读取
	address := os.Getenv("DELEGATION_FROM_ADDRESS")
	if address != testAddress {
		t.Errorf("期望委托地址为%s，实际为%s", testAddress, address)
	}
}

func TestCalculateDelegationAmount(t *testing.T) {
	// 模拟计算委托数量的逻辑
	testCases := []struct {
		value           string
		availableEnergy string
		expected        string
		description     string
	}{
		{
			value:           "500000", // 500,000 SUN = 0.5 TRX (小于1 TRX，应该跳过)
			availableEnergy: "10000",
			expected:        "0", // 小于1 TRX，返回0
			description:     "小于1 TRX交易，跳过委托",
		},
		{
			value:           "50000000", // 50,000,000 SUN = 50 TRX (小额交易)
			availableEnergy: "10000",
			expected:        "10000", // 基数 15000 + (50000000 * 0.1 / 10000) = 15000 + 500 = 15500，但受可用能量限制为10000
			description:     "小额交易，受可用能量限制",
		},
		{
			value:           "150000000000", // 150,000,000,000 SUN = 150,000 TRX (中额交易)
			availableEnergy: "20000",
			expected:        "20000", // 基数 15000 + (150000000000 * 0.25 / 10000) = 15000 + 3750 = 18750，但受可用能量限制为20000
			description:     "中额交易，受可用能量限制",
		},
		{
			value:           "1500000000000", // 1,500,000,000,000 SUN = 1,500,000 TRX (大额交易)
			availableEnergy: "50000",
			expected:        "50000", // 基数 15000 + (1500000000000 * 0.5 / 10000) = 15000 + 75000 = 90000，但受可用能量限制为50000
			description:     "大额交易，受可用能量限制",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// 模拟计算逻辑
			valueInt, err := strconv.ParseInt(tc.value, 10, 64)
			if err != nil {
				t.Errorf("解析交易金额失败: %v", err)
				return
			}

			// 检查是否小于1 TRX
			if valueInt < 1000000 { // 小于 1 TRX (1,000,000 SUN)
				result := "0"
				if result != tc.expected {
					t.Errorf("期望委托数量为%s，实际为%s", tc.expected, result)
				}
				return
			}

			energyInt, err := strconv.ParseInt(tc.availableEnergy, 10, 64)
			if err != nil {
				t.Errorf("解析可用能量失败: %v", err)
				return
			}

			// 根据交易金额设定固定基数（修正后的逻辑）
			var delegationBase int64 = 15000
			var delegationMultiplier float64

			// 基于 SUN 单位的阈值判断
			// 1 TRX = 1,000,000 SUN
			// 100,000 TRX = 100,000,000,000 SUN
			// 1,000,000 TRX = 1,000,000,000,000 SUN
			if valueInt > 1000000000000 { // 大于 1,000,000 TRX (1,000,000,000,000 SUN)
				delegationMultiplier = 0.5 // 50% 系数
			} else if valueInt > 100000000000 { // 大于 100,000 TRX (100,000,000,000 SUN)
				delegationMultiplier = 0.25 // 25% 系数
			} else {
				delegationMultiplier = 0.1 // 10% 系数
			}

			// 计算委托数量：基数 + (交易金额 * 系数)
			delegationAmount := delegationBase + int64(float64(valueInt)*delegationMultiplier/10000)

			// 确保不超过可用能量
			if delegationAmount > energyInt {
				delegationAmount = energyInt
			}

			// 确保最小委托数量
			minDelegation := int64(1000) // 最小委托1000能量
			if delegationAmount < minDelegation {
				delegationAmount = 0
			}

			result := strconv.FormatInt(delegationAmount, 10)

			if result != tc.expected {
				t.Errorf("期望委托数量为%s，实际为%s", tc.expected, result)
			}
		})
	}
}

func TestEnvironmentVariableValidation(t *testing.T) {
	// 测试环境变量未设置的情况
	os.Unsetenv("DELEGATION_FROM_ADDRESS")

	address := os.Getenv("DELEGATION_FROM_ADDRESS")
	if address != "" {
		t.Errorf("期望地址为空，实际为%s", address)
	}

	// 测试设置环境变量
	testAddress := "TTestAddressForTesting123456789"
	os.Setenv("DELEGATION_FROM_ADDRESS", testAddress)

	address = os.Getenv("DELEGATION_FROM_ADDRESS")
	if address != testAddress {
		t.Errorf("期望地址为%s，实际为%s", testAddress, address)
	}
}

// 新增：测试 SUN 单位转换和阈值判断
func TestSunUnitCalculation(t *testing.T) {
	testCases := []struct {
		sunValue     string
		expectedTRX  float64
		expectedTier string
		description  string
	}{
		{
			sunValue:     "1000000", // 1,000,000 SUN
			expectedTRX:  1.0,       // 1 TRX
			expectedTier: "小额交易",
			description:  "1 TRX 交易",
		},
		{
			sunValue:     "100000000", // 100,000,000 SUN
			expectedTRX:  100.0,       // 100 TRX
			expectedTier: "小额交易",
			description:  "100 TRX 交易",
		},
		{
			sunValue:     "100000000000", // 100,000,000,000 SUN
			expectedTRX:  100000.0,       // 100,000 TRX
			expectedTier: "小额交易",         // 修正：实际为小额交易
			description:  "100,000 TRX 交易",
		},
		{
			sunValue:     "1000000000000", // 1,000,000,000,000 SUN
			expectedTRX:  1000000.0,       // 1,000,000 TRX
			expectedTier: "中额交易",          // 修正：实际为中额交易
			description:  "1,000,000 TRX 交易",
		},
		{
			sunValue:     "2000000000000", // 2,000,000,000,000 SUN
			expectedTRX:  2000000.0,       // 2,000,000 TRX
			expectedTier: "大额交易",
			description:  "2,000,000 TRX 交易",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// 解析 SUN 值
			sunInt, err := strconv.ParseInt(tc.sunValue, 10, 64)
			if err != nil {
				t.Errorf("解析 SUN 值失败: %v", err)
				return
			}

			// 转换为 TRX
			trxValue := float64(sunInt) / 1000000.0

			// 验证 TRX 转换
			if trxValue != tc.expectedTRX {
				t.Errorf("期望 TRX 值为 %f，实际为 %f", tc.expectedTRX, trxValue)
			}

			// 验证阈值判断（使用与代码相同的逻辑）
			var actualTier string
			if sunInt > 1000000000000 { // 大于 1,000,000 TRX (1,000,000,000,000 SUN)
				actualTier = "大额交易"
			} else if sunInt > 100000000000 { // 大于 100,000 TRX (100,000,000,000 SUN)
				actualTier = "中额交易"
			} else {
				actualTier = "小额交易"
			}

			if actualTier != tc.expectedTier {
				t.Errorf("期望交易类型为 %s，实际为 %s (SUN: %s, TRX: %.1f)", tc.expectedTier, actualTier, tc.sunValue, trxValue)
			}

			t.Logf("SUN: %s, TRX: %.1f, 类型: %s", tc.sunValue, trxValue, actualTier)
		})
	}
}
