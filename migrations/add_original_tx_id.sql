-- 添加 original_tx_id 字段到 webhook_data 表
-- 这个字段用于存储原始委托交易的ID

-- 检查字段是否已存在，如果不存在则添加
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'webhook_data' 
        AND column_name = 'original_tx_id'
    ) THEN
        ALTER TABLE webhook_data ADD COLUMN original_tx_id VARCHAR(255) UNIQUE;
        RAISE NOTICE '已添加 original_tx_id 字段到 webhook_data 表';
    ELSE
        RAISE NOTICE 'original_tx_id 字段已存在';
    END IF;
END $$;

-- 创建索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_webhook_data_original_tx_id ON webhook_data(original_tx_id);
CREATE INDEX IF NOT EXISTS idx_webhook_data_tx_hash_original_tx_id ON webhook_data(tx_hash, original_tx_id);

-- 验证字段添加成功
SELECT 
    column_name, 
    data_type, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'webhook_data' 
AND column_name = 'original_tx_id'; 