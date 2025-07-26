-- 运行 tx_hash 唯一约束迁移
-- 执行命令: psql -d lending_trx -f run_migration.sql

\echo '开始执行 tx_hash 唯一约束迁移...'

\i migrations/add_tx_hash_unique_constraint.sql

\echo '迁移完成!' 