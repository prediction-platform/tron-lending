package db

import (
	"testing"
	"time"
)

func TestUpdateOriginalTxIDByID(t *testing.T) {
	// 这个测试需要真实的数据库连接
	// 在实际环境中运行
	/*
		ctx := context.Background()
		pool, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("初始化数据库失败: %v", err)
		}
		defer pool.Close()

		// 测试数据
		testID := int64(1)
		testOriginalTxID := "tx_test_123456"

		// 更新 original_tx_id
		err = UpdateOriginalTxIDByID(ctx, pool, testID, testOriginalTxID)
		if err != nil {
			t.Errorf("更新 original_tx_id 失败: %v", err)
		}

		// 验证更新结果
		originalTxID, err := GetOriginalTxIDByID(ctx, pool, testID)
		if err != nil {
			t.Errorf("获取 original_tx_id 失败: %v", err)
		}

		if originalTxID != testOriginalTxID {
			t.Errorf("期望 original_tx_id 为 %s，实际为 %s", testOriginalTxID, originalTxID)
		}
	*/
}

func TestUpdateOriginalTxIDByTxHash(t *testing.T) {
	// 这个测试需要真实的数据库连接
	// 在实际环境中运行
	/*
		ctx := context.Background()
		pool, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("初始化数据库失败: %v", err)
		}
		defer pool.Close()

		// 测试数据
		testTxHash := "0x1234567890abcdef"
		testOriginalTxID := "tx_test_789012"

		// 更新 original_tx_id
		err = UpdateOriginalTxIDByTxHash(ctx, pool, testTxHash, testOriginalTxID)
		if err != nil {
			t.Errorf("更新 original_tx_id 失败: %v", err)
		}

		// 验证更新结果
		originalTxID, err := GetOriginalTxIDByTxHash(ctx, pool, testTxHash)
		if err != nil {
			t.Errorf("获取 original_tx_id 失败: %v", err)
		}

		if originalTxID != testOriginalTxID {
			t.Errorf("期望 original_tx_id 为 %s，实际为 %s", testOriginalTxID, originalTxID)
		}
	*/
}

func TestGetOriginalTxIDByID(t *testing.T) {
	// 这个测试需要真实的数据库连接
	// 在实际环境中运行
	/*
		ctx := context.Background()
		pool, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("初始化数据库失败: %v", err)
		}
		defer pool.Close()

		// 测试数据
		testID := int64(1)

		// 获取 original_tx_id
		originalTxID, err := GetOriginalTxIDByID(ctx, pool, testID)
		if err != nil {
			t.Errorf("获取 original_tx_id 失败: %v", err)
		}

		t.Logf("获取到的 original_tx_id: %s", originalTxID)
	*/
}

func TestGetOriginalTxIDByTxHash(t *testing.T) {
	// 这个测试需要真实的数据库连接
	// 在实际环境中运行
	/*
		ctx := context.Background()
		pool, err := InitDB(ctx)
		if err != nil {
			t.Fatalf("初始化数据库失败: %v", err)
		}
		defer pool.Close()

		// 测试数据
		testTxHash := "0x1234567890abcdef"

		// 获取 original_tx_id
		originalTxID, err := GetOriginalTxIDByTxHash(ctx, pool, testTxHash)
		if err != nil {
			t.Errorf("获取 original_tx_id 失败: %v", err)
		}

		t.Logf("获取到的 original_tx_id: %s", originalTxID)
	*/
}

func TestWebhookDataModelWithOriginalTxID(t *testing.T) {
	// 测试 WebhookDataModel 结构体
	data := &WebhookDataModel{
		ID:           1,
		BlockHeight:  74204442,
		TxHash:       "0x1234567890abcdef",
		FromAddress:  "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
		ToAddress:    "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		Value:        "1000000000",
		BlockTime:    time.Now().UnixMilli(),
		CreateTime:   "2024-01-01 12:00:00",
		UpdateTime:   "2024-01-01 12:00:00",
		ExpireTime:   time.Now().Add(time.Hour).UnixMilli(),
		Status:       0,
		OriginalTxID: "tx_original_123456",
	}

	if data.OriginalTxID != "tx_original_123456" {
		t.Errorf("期望 OriginalTxID 为 tx_original_123456，实际为 %s", data.OriginalTxID)
	}

	t.Logf("WebhookDataModel 结构体测试通过: %+v", data)
}
