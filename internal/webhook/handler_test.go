package webhook

import (
	"testing"
	"time"
)

func TestParseWebhookData(t *testing.T) {
	// 测试数据
	testJSON := `{
		"data": [
			{
				"blockHash": "0x00000000046c451a6f749bf87f4be6c3bc49bcdfc85e309f0529907a4de697f9",
				"blockNumber": "0x46c451a",
				"from": "0xb8a57ef5343f88712a4eee91e34290584c2d5998",
				"gas": "0x0",
				"gasPrice": "0xd2",
				"hash": "0x07e1f7519110b58ed7cdfbfccbe5b6d35ca00d7c59b21bb72ba96a77ce25675e",
				"input": "0x",
				"nonce": "0x0000000000000000",
				"r": "0xa5897110eaed6d05e5a300c797bd0fd700f50d3c3c71c5539ade0c2098a8e48d",
				"s": "0x346635157360f9e5dfce4bee3f7d0bd8c7aed8349dbe037b48170f0e3e52cac5",
				"timestamp": "0x6880ce30",
				"to": "0x678637325f9be6b2264db347021432a6a7b84c10",
				"transactionIndex": "0x0",
				"type": "0x0",
				"v": "0x1c",
				"value": "0x6"
			}
		],
		"metadata": {
			"batch_end_range": 74204442,
			"batch_start_range": 74204442,
			"data_size_bytes": 239050,
			"dataset": "block",
			"end_range": 74204442,
			"keep_distance_from_tip": 0,
			"network": "tron-mainnet",
			"start_range": 74204442,
			"stream_id": "7a67d416-30d2-4095-a48d-18d02f19ee37",
			"stream_name": "test-stream",
			"stream_region": "usa_east"
		}
	}`

	// 解析数据
	result, err := ParseWebhookData([]byte(testJSON))
	if err != nil {
		t.Fatalf("ParseWebhookData failed: %v", err)
	}

	// 验证结果
	if len(result) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(result))
	}

	data := result[0]

	// 验证字段转换
	expectedBlockHeight := int64(74204442) // 0x46c451a 的十进制值
	if data.BlockHeight != expectedBlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", expectedBlockHeight, data.BlockHeight)
	}

	expectedBlockTime := int64(1753271856) // 0x6880ce30 的十进制值
	if data.BlockTime != expectedBlockTime {
		t.Errorf("Expected BlockTime %d, got %d", expectedBlockTime, data.BlockTime)
	}

	expectedValue := "6"
	if data.Value != expectedValue {
		t.Errorf("Expected Value %s, got %s", expectedValue, data.Value)
	}

	expectedTxHash := "0x07e1f7519110b58ed7cdfbfccbe5b6d35ca00d7c59b21bb72ba96a77ce25675e"
	if data.TxHash != expectedTxHash {
		t.Errorf("Expected TxHash %s, got %s", expectedTxHash, data.TxHash)
	}

	expectedFrom := "0xb8a57ef5343f88712a4eee91e34290584c2d5998"
	if data.FromAddress != expectedFrom {
		t.Errorf("Expected FromAddress %s, got %s", expectedFrom, data.FromAddress)
	}

	expectedTo := "0x678637325f9be6b2264db347021432a6a7b84c10"
	if data.ToAddress != expectedTo {
		t.Errorf("Expected ToAddress %s, got %s", expectedTo, data.ToAddress)
	}

	// 验证ExpireTime = BlockTime + 1小时
	expectedExpireTime := expectedBlockTime + 3600
	if data.ExpireTime != expectedExpireTime {
		t.Errorf("Expected ExpireTime %d, got %d", expectedExpireTime, data.ExpireTime)
	}

	if data.Status != 0 {
		t.Errorf("Expected Status 0, got %d", data.Status)
	}
}

func TestHexToInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0x46c451a", 74204442},
		{"0x6880ce30", 1753271856},
		{"0x0", 0},
		{"0x", 0},
		{"", 0},
	}

	for _, test := range tests {
		result, err := hexToInt64(test.input)
		if err != nil {
			t.Errorf("hexToInt64(%s) failed: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("hexToInt64(%s) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestHexToString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0x6", "6"},
		{"0x0", "0"},
		{"0x", "0"},
		{"", "0"},
		{"0x1a", "26"},
	}

	for _, test := range tests {
		result, err := hexToString(test.input)
		if err != nil {
			t.Errorf("hexToString(%s) failed: %v", test.input, err)
			continue
		}
		if result != test.expected {
			t.Errorf("hexToString(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestConvertToWebhookDataModel(t *testing.T) {
	// 创建测试数据
	webhookData := WebhookData{
		BlockHeight: 74204442,
		TxHash:      "0x07e1f7519110b58ed7cdfbfccbe5b6d35ca00d7c59b21bb72ba96a77ce25675e",
		FromAddress: "0xb8a57ef5343f88712a4eee91e34290584c2d5998",
		ToAddress:   "0x678637325f9be6b2264db347021432a6a7b84c10",
		Value:       "6",
		BlockTime:   1753271856,
		ExpireTime:  1753271856 + 3600, // BlockTime + 1小时
		Status:      0,
	}

	// 转换为WebhookDataModel
	result := ConvertToWebhookDataModel(webhookData)

	// 验证转换结果
	if result.BlockHeight != webhookData.BlockHeight {
		t.Errorf("Expected BlockHeight %d, got %d", webhookData.BlockHeight, result.BlockHeight)
	}

	if result.TxHash != webhookData.TxHash {
		t.Errorf("Expected TxHash %s, got %s", webhookData.TxHash, result.TxHash)
	}

	if result.FromAddress != webhookData.FromAddress {
		t.Errorf("Expected FromAddress %s, got %s", webhookData.FromAddress, result.FromAddress)
	}

	if result.ToAddress != webhookData.ToAddress {
		t.Errorf("Expected ToAddress %s, got %s", webhookData.ToAddress, result.ToAddress)
	}

	if result.Value != webhookData.Value {
		t.Errorf("Expected Value %s, got %s", webhookData.Value, result.Value)
	}

	if result.BlockTime != webhookData.BlockTime {
		t.Errorf("Expected BlockTime %d, got %d", webhookData.BlockTime, result.BlockTime)
	}

	if result.ExpireTime != webhookData.ExpireTime {
		t.Errorf("Expected ExpireTime %d, got %d", webhookData.ExpireTime, result.ExpireTime)
	}

	if result.Status != webhookData.Status {
		t.Errorf("Expected Status %d, got %d", webhookData.Status, result.Status)
	}

	// 验证时间字段格式
	_, err := time.Parse("2006-01-02 15:04:05", result.CreateTime)
	if err != nil {
		t.Errorf("Invalid CreateTime format: %s", result.CreateTime)
	}

	_, err = time.Parse("2006-01-02 15:04:05", result.UpdateTime)
	if err != nil {
		t.Errorf("Invalid UpdateTime format: %s", result.UpdateTime)
	}
}

func TestConvertToWebhookDataModelSlice(t *testing.T) {
	// 创建测试数据
	webhookDataList := []WebhookData{
		{
			BlockHeight: 74204442,
			TxHash:      "0x07e1f7519110b58ed7cdfbfccbe5b6d35ca00d7c59b21bb72ba96a77ce25675e",
			FromAddress: "0xb8a57ef5343f88712a4eee91e34290584c2d5998",
			ToAddress:   "0x678637325f9be6b2264db347021432a6a7b84c10",
			Value:       "6",
			BlockTime:   1753271856,
			ExpireTime:  1753271856 + 3600, // BlockTime + 1小时
			Status:      0,
		},
		{
			BlockHeight: 74204443,
			TxHash:      "0x17e1f7519110b58ed7cdfbfccbe5b6d35ca00d7c59b21bb72ba96a77ce25675f",
			FromAddress: "0xc8a57ef5343f88712a4eee91e34290584c2d5999",
			ToAddress:   "0x778637325f9be6b2264db347021432a6a7b84c11",
			Value:       "16",
			BlockTime:   1753271857,
			ExpireTime:  1753271857 + 3600, // BlockTime + 1小时
			Status:      0,
		},
	}

	// 转换为WebhookDataModel切片
	result := ConvertToWebhookDataModelSlice(webhookDataList)

	// 验证结果
	if len(result) != len(webhookDataList) {
		t.Errorf("Expected %d items, got %d", len(webhookDataList), len(result))
	}

	for i, original := range webhookDataList {
		converted := result[i]
		if converted.BlockHeight != original.BlockHeight {
			t.Errorf("Item %d: Expected BlockHeight %d, got %d", i, original.BlockHeight, converted.BlockHeight)
		}
		if converted.TxHash != original.TxHash {
			t.Errorf("Item %d: Expected TxHash %s, got %s", i, original.TxHash, converted.TxHash)
		}
	}
}
