package tron

import (
	"testing"
)

func TestNewTronClient(t *testing.T) {
	baseURL := "https://api.trongrid.io"
	apiKey := "test-api-key"

	client := NewTronClient(baseURL, apiKey)

	if client.baseURL != baseURL {
		t.Errorf("期望baseURL为%s，实际为%s", baseURL, client.baseURL)
	}

	if client.apiKey != apiKey {
		t.Errorf("期望apiKey为%s，实际为%s", apiKey, client.apiKey)
	}

	if client.httpClient == nil {
		t.Error("httpClient不能为nil")
	}
}

func TestEnergyDelegationRequest(t *testing.T) {
	req := &EnergyDelegationRequest{
		FromAddress: "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
		ToAddress:   "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		Amount:      "10000",
		TxHash:      "0x1234567890abcdef",
		BlockHeight: 74204442,
	}

	if req.FromAddress != "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs" {
		t.Errorf("期望FromAddress为TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs，实际为%s", req.FromAddress)
	}

	if req.Amount != "10000" {
		t.Errorf("期望Amount为10000，实际为%s", req.Amount)
	}

	if req.BlockHeight != 74204442 {
		t.Errorf("期望BlockHeight为74204442，实际为%d", req.BlockHeight)
	}
}

func TestCancelDelegationRequest(t *testing.T) {
	req := &CancelDelegationRequest{
		FromAddress:  "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
		ToAddress:    "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		OriginalTxID: "original_tx_id_123",
		TxHash:       "0x1234567890abcdef",
	}

	if req.FromAddress != "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs" {
		t.Errorf("期望FromAddress为TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs，实际为%s", req.FromAddress)
	}

	if req.OriginalTxID != "original_tx_id_123" {
		t.Errorf("期望OriginalTxID为original_tx_id_123，实际为%s", req.OriginalTxID)
	}
}

func TestAccountInfo(t *testing.T) {
	account := &AccountInfo{
		Address:     "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
		Balance:     "1000000000",
		Energy:      "50000",
		Frozen:      "0",
		NetUsed:     "1000",
		NetLimit:    "10000",
		EnergyUsed:  "5000",
		EnergyLimit: "100000",
	}

	if account.Address != "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs" {
		t.Errorf("期望Address为TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs，实际为%s", account.Address)
	}

	if account.Balance != "1000000000" {
		t.Errorf("期望Balance为1000000000，实际为%s", account.Balance)
	}

	if account.Energy != "50000" {
		t.Errorf("期望Energy为50000，实际为%s", account.Energy)
	}
}

func TestEnergyDelegationResponse(t *testing.T) {
	successResp := &EnergyDelegationResponse{
		Success: true,
		TxID:    "tx_id_123456",
		Message: "能量委托成功",
	}

	if !successResp.Success {
		t.Error("期望Success为true")
	}

	if successResp.TxID != "tx_id_123456" {
		t.Errorf("期望TxID为tx_id_123456，实际为%s", successResp.TxID)
	}

	errorResp := &EnergyDelegationResponse{
		Success: false,
		Error:   "能量不足",
	}

	if errorResp.Success {
		t.Error("期望Success为false")
	}

	if errorResp.Error != "能量不足" {
		t.Errorf("期望Error为能量不足，实际为%s", errorResp.Error)
	}
}

func TestCancelDelegationResponse(t *testing.T) {
	successResp := &CancelDelegationResponse{
		Success: true,
		TxID:    "cancel_tx_id_123456",
		Message: "取消委托成功",
	}

	if !successResp.Success {
		t.Error("期望Success为true")
	}

	if successResp.TxID != "cancel_tx_id_123456" {
		t.Errorf("期望TxID为cancel_tx_id_123456，实际为%s", successResp.TxID)
	}

	errorResp := &CancelDelegationResponse{
		Success: false,
		Error:   "委托不存在",
	}

	if errorResp.Success {
		t.Error("期望Success为false")
	}

	if errorResp.Error != "委托不存在" {
		t.Errorf("期望Error为委托不存在，实际为%s", errorResp.Error)
	}
}

// 集成测试（需要真实的API环境）
func TestTronClientIntegration(t *testing.T) {
	// 这个测试需要真实的Tron API环境
	// 在实际环境中运行
	/*
		ctx := context.Background()
		client := NewTronClient("https://api.trongrid.io", "your-api-key")

		// 测试获取账户信息
		accountInfo, err := client.GetAccountInfo(ctx, "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs")
		if err != nil {
			t.Errorf("获取账户信息失败: %v", err)
			return
		}

		t.Logf("账户信息: %+v", accountInfo)

		// 测试能量委托
		delegationReq := &EnergyDelegationRequest{
			FromAddress: "TQn9Y2khDD95J42FQtQTdwVVRyc2jBEsVs",
			ToAddress:   "TRX7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			Amount:      "1000",
			TxHash:      "0x1234567890abcdef",
			BlockHeight: 74204442,
		}

		delegationResp, err := client.DelegateEnergy(ctx, delegationReq)
		if err != nil {
			t.Errorf("能量委托失败: %v", err)
			return
		}

		t.Logf("委托响应: %+v", delegationResp)
	*/
}
