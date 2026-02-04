# Database Schema Overview

This project uses a shared PostgreSQL database (`appdb`) with tables owned by each service. The migration service manages SQL migrations for the user and order services, while the audit log service uses GORM auto-migrations.

## Conventions

All tables include the following audit and soft-delete fields:

- `created_at` (TIMESTAMPTZ)
- `updated_at` (TIMESTAMPTZ)
- `created_by` (VARCHAR(100))
- `updated_by` (VARCHAR(100))
- `status` (VARCHAR(20))

`status` is used for soft delete and visibility. Each service enforces its own status values.

## Tables

### `users`

Owned by: User Service

Columns:
- `id` BIGSERIAL PRIMARY KEY
- `email` TEXT NOT NULL UNIQUE
- `name` TEXT NOT NULL
- `age` INTEGER NOT NULL
- `status` VARCHAR(20) NOT NULL DEFAULT 'active'
- `created_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `updated_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

Indexes:
- `idx_users_status` on (`status`)
- Unique index on `email`

Status values:
- `active`
- `inactive`
- `deleted` (soft delete)

### `orders`

Owned by: Order Service

Columns:
- `id` BIGSERIAL PRIMARY KEY
- `user_id` BIGINT NOT NULL
- `product_id` TEXT NOT NULL
- `quantity` INTEGER NOT NULL
- `total_price` NUMERIC(12,2) NOT NULL
- `order_status` VARCHAR(20) NOT NULL DEFAULT 'pending'
- `status` VARCHAR(20) NOT NULL DEFAULT 'active'
- `created_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `updated_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

Indexes:
- `idx_orders_user_id` on (`user_id`)
- `idx_orders_status` on (`status`)
- `idx_orders_order_status` on (`order_status`)

Status values:
- `status`: `active`, `deleted`
- `order_status`: `pending`, `confirmed`, `shipped`, `delivered`, `cancelled`

Relationships:
- `orders.user_id` references `users.id` at the application layer (no foreign key constraint in migrations).

### `audit_logs`

Owned by: Audit Log Service

Columns:
- `id` BIGSERIAL PRIMARY KEY
- `actor` VARCHAR(120) NOT NULL
- `action` VARCHAR(100) NOT NULL
- `resource_type` VARCHAR(100) NOT NULL
- `resource_id` VARCHAR(100) NOT NULL
- `description` TEXT
- `metadata` TEXT
- `status` VARCHAR(20) NOT NULL DEFAULT 'active'
- `created_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `updated_by` VARCHAR(100) NOT NULL DEFAULT 'system'
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

Indexes (via GORM):
- `actor`, `action`, `resource_type`, `resource_id`, `status`

Status values:
- `active`
- `deleted`

## Migration Sources

- User Service: `services/migration-service/migrations/user/`
- Order Service: `services/migration-service/migrations/order/`
- Audit Log Service: GORM auto-migrate on startup (see `services/audit-log-service/internal/model/audit_log.go`)
