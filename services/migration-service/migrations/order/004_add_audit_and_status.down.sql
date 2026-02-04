ALTER TABLE orders DROP COLUMN IF EXISTS status;
ALTER TABLE orders DROP COLUMN IF EXISTS created_by;
ALTER TABLE orders DROP COLUMN IF EXISTS updated_by;

ALTER TABLE orders RENAME COLUMN order_status TO status;

ALTER TABLE orders ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ NULL;

DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_order_status;
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders (deleted_at);
