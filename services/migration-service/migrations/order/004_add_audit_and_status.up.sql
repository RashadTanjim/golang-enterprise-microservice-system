ALTER TABLE orders RENAME COLUMN status TO order_status;

ALTER TABLE orders ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS created_by VARCHAR(100) NOT NULL DEFAULT 'system';
ALTER TABLE orders ADD COLUMN IF NOT EXISTS updated_by VARCHAR(100) NOT NULL DEFAULT 'system';

ALTER TABLE orders DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_orders_deleted_at;
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (status);
CREATE INDEX IF NOT EXISTS idx_orders_order_status ON orders (order_status);
