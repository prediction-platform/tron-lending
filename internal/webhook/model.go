package webhook

// WebhookDataModel 用于表示 webhook_data 表结构
// 注意：Value 用 string 以兼容大数，时间戳用 int64
// create_time, update_time 用 string（或 time.Time，视序列化需求）
type WebhookDataModel struct {
	ID          int64  `json:"id"`           // 主键唯一ID
	BlockHeight int64  `json:"block_height"` // 区块高度
	TxHash      string `json:"tx_hash"`      // 交易哈希
	FromAddress string `json:"from_address"` // 发送方地址
	ToAddress   string `json:"to_address"`   // 接收方地址
	Value       string `json:"value"`        // 交易金额（大整数，字符串存储）
	BlockTime   int64  `json:"block_time"`   // 区块时间（毫秒时间戳）
	CreateTime  string `json:"create_time"`  // 创建时间
	UpdateTime  string `json:"update_time"`  // 更新时间
	ExpireTime  int64  `json:"expire_time"`  // 有效期（毫秒时间戳）
	Status      int16  `json:"status"`       // 状态（0:初始化，1:执行中，2:已授权，3:已回收）
}
