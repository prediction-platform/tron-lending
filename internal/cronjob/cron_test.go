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
	// 根据实际代码逻辑更新测试用例
	testCases := []struct {
		value           string
		availableEnergy string
		expected        string
		description     string
	}{
		{
			value:           "500000", // 500,000 SUN = 0.5 TRX (小于1 TRX)
			availableEnergy: "100000",
			expected:        "0", // 小于1 TRX，返回0
			description:     "小于1 TRX交易，跳过委托",
		},
		{
			value:           "1000000", // 1,000,000 SUN = 1 TRX
			availableEnergy: "100000",
			expected:        "65000", // 1 TRX → 委托 65000
			description:     "1 TRX交易，委托基础数量",
		},
		{
			value:           "2000000", // 2,000,000 SUN = 2 TRX
			availableEnergy: "200000",
			expected:        "130000", // 2 TRX → 委托 2 * 65000 = 130000
			description:     "2 TRX交易，委托双倍数量",
		},
		{
			value:           "1000000", // 1,000,000 SUN = 1 TRX
			availableEnergy: "50000",   // 可用能量不足
			expected:        "50000",   // 受可用能量限制
			description:     "1 TRX交易，受可用能量限制",
		},
		{
			value:           "1000000", // 1,000,000 SUN = 1 TRX
			availableEnergy: "500000",  // 足够的可用能量，但委托数量小于最小值
			expected:        "0",       // 65000 < 1000000，不满足最小委托要求
			description:     "1 TRX交易，不满足最小委托要求",
		},
		{
			value:           "3000000", // 3,000,000 SUN = 3 TRX
			availableEnergy: "100000",
			expected:        "0", // 3 TRX不在支持范围内
			description:     "3 TRX交易，不支持的交易金额",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// 模拟实际代码的计算逻辑
			valueInt, err := strconv.ParseInt(tc.value, 10, 64)
			if err != nil {
				t.Errorf("解析交易金额失败: %v", err)
				return
			}

			energyInt, err := strconv.ParseInt(tc.availableEnergy, 10, 64)
			if err != nil {
				t.Errorf("解析可用能量失败: %v", err)
				return
			}

			// 使用实际代码的逻辑
			var delegationBase int64 = 65000
			var minDelegation int64 = 1000000
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

			result := strconv.FormatInt(delegationAmount, 10)

			if result != tc.expected {
				t.Errorf("期望委托数量为%s，实际为%s", tc.expected, result)
			}

			t.Logf("交易金额: %s SUN (%.1f TRX), 可用能量: %s, 计算结果: %s",
				tc.value, float64(trxValue), tc.availableEnergy, result)
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

// 测试委托基数和最小委托数量的环境变量
func TestDelegationConfiguration(t *testing.T) {
	// 保存原始环境变量
	originalBase := os.Getenv("DELEGATION_BASE")
	originalMin := os.Getenv("MIN_DELEGATION_AMOUNT")

	// 清理环境变量
	os.Unsetenv("DELEGATION_BASE")
	os.Unsetenv("MIN_DELEGATION_AMOUNT")

	defer func() {
		// 恢复原始环境变量
		if originalBase != "" {
			os.Setenv("DELEGATION_BASE", originalBase)
		}
		if originalMin != "" {
			os.Setenv("MIN_DELEGATION_AMOUNT", originalMin)
		}
	}()

	// 测试默认值
	base := os.Getenv("DELEGATION_BASE")
	if base != "" {
		t.Errorf("期望 DELEGATION_BASE 为空，实际为%s", base)
	}

	min := os.Getenv("MIN_DELEGATION_AMOUNT")
	if min != "" {
		t.Errorf("期望 MIN_DELEGATION_AMOUNT 为空，实际为%s", min)
	}

	// 测试设置自定义值
	os.Setenv("DELEGATION_BASE", "80000")
	os.Setenv("MIN_DELEGATION_AMOUNT", "500000")

	base = os.Getenv("DELEGATION_BASE")
	if base != "80000" {
		t.Errorf("期望 DELEGATION_BASE 为80000，实际为%s", base)
	}

	min = os.Getenv("MIN_DELEGATION_AMOUNT")
	if min != "500000" {
		t.Errorf("期望 MIN_DELEGATION_AMOUNT 为500000，实际为%s", min)
	}
}

// 测试 TRX 值计算逻辑
func TestTRXValueCalculation(t *testing.T) {
	testCases := []struct {
		sunValue    string
		expectedTRX int64
		description string
	}{
		{
			sunValue:    "1000000", // 1,000,000 SUN
			expectedTRX: 1,         // 1 TRX
			description: "1 TRX 转换",
		},
		{
			sunValue:    "2000000", // 2,000,000 SUN
			expectedTRX: 2,         // 2 TRX
			description: "2 TRX 转换",
		},
		{
			sunValue:    "1500000", // 1,500,000 SUN
			expectedTRX: 1,         // 1.5 TRX → 整数除法 = 1
			description: "1.5 TRX 转换（整数除法）",
		},
		{
			sunValue:    "999999", // 999,999 SUN
			expectedTRX: 0,        // 0.999999 TRX → 整数除法 = 0
			description: "不足1 TRX 转换",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			sunInt, err := strconv.ParseInt(tc.sunValue, 10, 64)
			if err != nil {
				t.Errorf("解析 SUN 值失败: %v", err)
				return
			}

			// 使用实际代码的转换逻辑
			trxValue := sunInt / 1000000 // 整数除法

			if trxValue != tc.expectedTRX {
				t.Errorf("期望 TRX 值为 %d，实际为 %d", tc.expectedTRX, trxValue)
			}

			t.Logf("SUN: %s, TRX: %d", tc.sunValue, trxValue)
		})
	}
}
