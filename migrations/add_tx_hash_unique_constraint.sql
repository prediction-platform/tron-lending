-- 为 tx_hash 字段添加唯一约束
-- 这确保同一个交易哈希不会被重复处理

-- 首先检查是否已经存在唯一约束
DO $$
BEGIN
    -- 检查约束是否存在
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE table_name = 'webhook_data' 
        AND constraint_type = 'UNIQUE'
        AND constraint_name = 'webhook_data_tx_hash_key'
    ) THEN
        -- 在添加约束之前，先删除可能存在的重复数据
        -- 保留每个 tx_hash 的第一条记录（按 id 排序）
        DELETE FROM webhook_data 
        WHERE id NOT IN (
            SELECT MIN(id) 
            FROM webhook_data 
            GROUP BY tx_hash
        );
        
        -- 添加唯一约束
        ALTER TABLE webhook_data ADD CONSTRAINT webhook_data_tx_hash_key UNIQUE (tx_hash);
        RAISE NOTICE '已为 tx_hash 字段添加唯一约束';
    ELSE
        RAISE NOTICE 'tx_hash 字段的唯一约束已存在';
    END IF;
END $$;

-- 创建索引以提高查询性能（如果还没有的话）
CREATE INDEX IF NOT EXISTS idx_webhook_data_tx_hash ON webhook_data(tx_hash);

-- 验证约束添加成功
SELECT 
    constraint_name,
    constraint_type,
    table_name
FROM information_schema.table_constraints 
WHERE table_name = 'webhook_data' 
AND constraint_type = 'UNIQUE'
AND constraint_name = 'webhook_data_tx_hash_key'; 