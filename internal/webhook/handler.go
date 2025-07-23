package webhook

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"io"

	"lending-trx/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunjiangjun/xlog"
)

// WebhookRequest 表示完整的webhook请求结构
type WebhookRequest struct {
	Data     []TransactionData `json:"data"`
	Metadata Metadata          `json:"metadata"`
}

// TransactionData 表示单个交易数据
type TransactionData struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	Timestamp        string `json:"timestamp"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
	V                string `json:"v"`
	Value            string `json:"value"`
}

// Metadata 表示元数据
type Metadata struct {
	BatchEndRange       int64  `json:"batch_end_range"`
	BatchStartRange     int64  `json:"batch_start_range"`
	DataSizeBytes       int64  `json:"data_size_bytes"`
	Dataset             string `json:"dataset"`
	EndRange            int64  `json:"end_range"`
	KeepDistanceFromTip int64  `json:"keep_distance_from_tip"`
	Network             string `json:"network"`
	StartRange          int64  `json:"start_range"`
	StreamID            string `json:"stream_id"`
	StreamName          string `json:"stream_name"`
	StreamRegion        string `json:"stream_region"`
}

type WebhookData struct {
	BlockHeight int64  `json:"blockNumber"`
	TxHash      string `json:"hash"`
	FromAddress string `json:"from"`
	ToAddress   string `json:"to"`
	Value       string `json:"value"`
	BlockTime   int64  `json:"timestamp"`
	ExpireTime  int64  `json:"expire_time"`
	Status      int16  `json:"status"`
}

// ParseWebhookData 解析webhook请求数据，返回WebhookData数组
func ParseWebhookData(body []byte) ([]WebhookData, error) {
	var request WebhookRequest
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, err
	}

	var result []WebhookData
	for _, tx := range request.Data {
		webhookData, err := convertTransactionToWebhookData(tx)
		if err != nil {
			return nil, err
		}
		result = append(result, webhookData)
	}

	return result, nil
}

// convertTransactionToWebhookData 将TransactionData转换为WebhookData
func convertTransactionToWebhookData(tx TransactionData) (WebhookData, error) {
	// 转换blockNumber从hex字符串到int64
	blockHeight, err := hexToInt64(tx.BlockNumber)
	if err != nil {
		return WebhookData{}, err
	}

	// 转换timestamp从hex字符串到int64
	blockTime, err := hexToInt64(tx.Timestamp)
	if err != nil {
		return WebhookData{}, err
	}

	// 转换value从hex字符串到十进制字符串
	value, err := hexToString(tx.Value)
	if err != nil {
		return WebhookData{}, err
	}

	return WebhookData{
		BlockHeight: blockHeight,
		TxHash:      tx.Hash,
		FromAddress: tx.From,
		ToAddress:   tx.To,
		Value:       value,
		BlockTime:   blockTime,
		ExpireTime:  blockTime + 3600, // BlockTime + 1小时 (3600秒)
		Status:      0,                // 默认状态
	}, nil
}

// hexToInt64 将十六进制字符串转换为int64
func hexToInt64(hexStr string) (int64, error) {
	// 移除0x前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		return 0, nil
	}

	// 转换为十进制
	value, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// hexToString 将十六进制字符串转换为十进制字符串
func hexToString(hexStr string) (string, error) {
	// 移除0x前缀
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		return "0", nil
	}

	// 转换为十进制
	value, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return "0", err
	}

	return strconv.FormatInt(value, 10), nil
}

// ConvertToWebhookDataModel 将WebhookData转换为WebhookDataModel
func ConvertToWebhookDataModel(data WebhookData) *db.WebhookDataModel {
	return &db.WebhookDataModel{
		BlockHeight: data.BlockHeight,
		TxHash:      data.TxHash,
		FromAddress: data.FromAddress,
		ToAddress:   data.ToAddress,
		Value:       data.Value,
		BlockTime:   data.BlockTime,
		ExpireTime:  data.ExpireTime,
		Status:      data.Status,
		CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime:  time.Now().Format("2006-01-02 15:04:05"),
	}
}

// ConvertToWebhookDataModelSlice 将WebhookData切片转换为WebhookDataModel切片
func ConvertToWebhookDataModelSlice(dataList []WebhookData) []*db.WebhookDataModel {
	result := make([]*db.WebhookDataModel, len(dataList))
	for i, data := range dataList {
		result[i] = ConvertToWebhookDataModel(data)
	}
	return result
}

// RegisterRoutes 注册 webhook 路由
func RegisterRoutes(r *gin.Engine, ctx context.Context, pool *pgxpool.Pool, log *xlog.XLog) {
	l := log.WithField("module", "webhook")
	r.POST("/webhook", AuthMiddleware(), func(c *gin.Context) {

		/*

						{
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
			}

		*/

		// 读取请求体
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			l.Error("读取请求体失败", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
			return
		}
		// 打印请求体内容
		l.Info("webhook 请求体", string(body))

		// 使用工具函数解析webhook数据
		webhookDataList, err := ParseWebhookData(body)
		if err != nil {
			l.Error("解析webhook数据失败", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json format"})
			return
		}

		// 转换为WebhookDataModel并批量插入
		webhookDataModels := ConvertToWebhookDataModelSlice(webhookDataList)
		err = db.BatchInsertWebhookData(ctx, pool, webhookDataModels)
		if err != nil {
			l.Error("批量插入数据库失败", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "inserted_count": len(webhookDataList)})
	})
}
