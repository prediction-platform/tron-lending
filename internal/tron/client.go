package tron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TronClient Tron API 客户端
type TronClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewTronClient 创建新的 Tron 客户端
func NewTronClient(baseURL, apiKey string) *TronClient {
	return &TronClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// EnergyDelegationRequest 能量委托请求
type EnergyDelegationRequest struct {
	FromAddress string `json:"from_address"` // 委托方地址
	ToAddress   string `json:"to_address"`   // 接收方地址
	Amount      string `json:"amount"`       // 委托的能量数量
	// 以下字段用于业务追踪，不是 Tron API 必需字段
	TxHash      string `json:"tx_hash,omitempty"`      // 原始交易哈希（业务追踪）
	BlockHeight int64  `json:"block_height,omitempty"` // 区块高度（业务追踪）
}

// EnergyDelegationResponse 能量委托响应
type EnergyDelegationResponse struct {
	Success bool   `json:"success"`
	TxID    string `json:"tx_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// CancelDelegationRequest 取消委托请求
type CancelDelegationRequest struct {
	FromAddress  string `json:"from_address"`   // 委托方地址
	ToAddress    string `json:"to_address"`     // 接收方地址
	OriginalTxID string `json:"original_tx_id"` // 原始委托交易ID
	// 以下字段用于业务追踪，不是 Tron API 必需字段
	TxHash string `json:"tx_hash,omitempty"` // 原始交易哈希（业务追踪）
}

// CancelDelegationResponse 取消委托响应
type CancelDelegationResponse struct {
	Success bool   `json:"success"`
	TxID    string `json:"tx_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// AccountInfo 账户信息
type AccountInfo struct {
	Address     string `json:"address"`
	Balance     string `json:"balance"`
	Energy      string `json:"energy"`
	Frozen      string `json:"frozen"`
	NetUsed     string `json:"net_used"`
	NetLimit    string `json:"net_limit"`
	EnergyUsed  string `json:"energy_used"`
	EnergyLimit string `json:"energy_limit"`
}

// GetAccountInfo 获取账户信息
func (c *TronClient) GetAccountInfo(ctx context.Context, address string) (*AccountInfo, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", c.baseURL, address)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed, status code: %d, response: %s", resp.StatusCode, string(body))
	}

	var accountInfo AccountInfo
	if err := json.Unmarshal(body, &accountInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &accountInfo, nil
}

// DelegateEnergy 执行能量委托
func (c *TronClient) DelegateEnergy(ctx context.Context, req *EnergyDelegationRequest) (*EnergyDelegationResponse, error) {
	url := fmt.Sprintf("%s/v1/energy/delegate", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response EnergyDelegationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return &response, fmt.Errorf("energy delegation failed: %s", response.Error)
	}

	return &response, nil
}

// CancelEnergyDelegation 取消能量委托
func (c *TronClient) CancelEnergyDelegation(ctx context.Context, req *CancelDelegationRequest) (*CancelDelegationResponse, error) {
	url := fmt.Sprintf("%s/v1/energy/cancel-delegate", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response CancelDelegationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !response.Success {
		return &response, fmt.Errorf("cancel energy delegation failed: %s", response.Error)
	}

	return &response, nil
}

// GetTransactionInfo 获取交易信息
func (c *TronClient) GetTransactionInfo(ctx context.Context, txID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/transactions/%s", c.baseURL, txID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var txInfo map[string]interface{}
	if err := json.Unmarshal(body, &txInfo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return txInfo, nil
}
