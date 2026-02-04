# High-Level Design (HLD) – Enterprise Microservice System

## 1) Overview
This system is a production-ready microservice architecture built with Go. It consists of three core services (User Service, Order Service, and Audit Log Service) that share a single PostgreSQL database (tables per service), a shared Redis cache for read-heavy endpoints, shared libraries for cross-cutting concerns, and supporting infrastructure for observability, authentication, and CI.

## 2) Goals
- Provide clean, maintainable microservices with clear boundaries.
- Support full CRUD operations for users, orders, and audit logs.
- Ensure resilience with circuit breakers and rate limits.
- Provide observability (metrics, logs, health checks).
- Secure APIs with JWT authentication and role-based access control.
- Enable automated migrations and CI checks.

## 3) Non-Goals
- API gateway or service mesh integration.
- Distributed tracing (planned, not implemented).
- Multi-region, production-grade HA deployment out of the box.

## 4) Architecture

### 4.1 Services
- **User Service** (port 8081)
  - Manages user CRUD, validation, pagination, and filtering.
  - Issues JWT tokens to clients.

- **Order Service** (port 8082)
  - Manages order CRUD and validation.
  - Calls User Service to validate user existence and status.
  - Uses circuit breaker for resilience.

- **Audit Log Service** (port 8083)
  - Captures audit log entries for compliance and traceability.
  - Enforces role-based access for write operations.

### 4.2 Frontend Portal
- **Vue.js Portal** (port 8080 via gateway)
  - Enterprise UI for operations, observability, and documentation.
  - Connects through the Nginx gateway for API routing.

### 4.3 Gateway
- **Nginx Gateway**
  - Serves the Vue SPA and routes API traffic to backend services.
  - Provides consolidated paths for `/api`, `/health`, `/metrics`, and `/swagger`.

### 4.4 Shared Libraries (`common/`)
- **auth**: JWT token generation and parsing.
- **middleware**: CORS, request ID, recovery, logging, metrics, rate limiting, JWT auth.
- **logger**: Structured logging with Zap.
- **metrics**: Prometheus metrics.
- **errors/response**: Standard error types and API responses.

### 4.5 Data Stores
- **Shared DB**: PostgreSQL (appdb) with per-service tables
- **Redis Cache**: Shared cache for read-heavy endpoints

Services isolate data by table and access patterns within a shared database.

### 4.6 Audit & Soft Delete
- All tables include `created_at`, `updated_at`, `created_by`, `updated_by`, and `status`.
- `status` is used for visibility (soft delete) across services.

### 4.7 Diagram Source
- Mermaid source: `docs/hld-diagram.mmd`
- Rendered SVG: `docs/hld-diagram.svg`

## 5) Data Flow

### 5.1 Auth & Authorization
1. Client requests token from `POST /api/v1/auth/token` (User Service).
2. User Service validates client credentials and returns JWT.
3. Client calls protected endpoints with `Authorization: Bearer <token>`.
4. Middleware validates token and enforces role-based access.

### 5.3 Read Caching
1. Services check Redis for cached reads (detail or list endpoints).
2. If a cache hit exists, the response is returned without a database call.
3. On writes, services update or invalidate relevant cache entries.

### 5.2 Order Creation
1. Client calls `POST /api/v1/orders` on Order Service.
2. Order Service validates input.
3. Order Service calls User Service to validate the user (JWT service token).
4. If valid, order is persisted to Order DB.

## 6) Security Model
- **JWT authentication** for all `/api/v1` endpoints.
- **Role-based access control**:
  - `admin` can create/update/delete and list resources.
  - `user` can read and create orders.
  - `service` role for service-to-service calls.

## 7) Migrations
- **golang-migrate** used with embedded SQL migrations.
- Migrations run on service startup and can be triggered manually with:
  - `make migrate-user`
  - `make migrate-order`
- Audit Log Service uses GORM auto-migrations on startup.

## 8) Observability
- **Health checks**: `/health/user`, `/health/order`, `/health/audit-log` via gateway.
- **Metrics**: `/metrics/user`, `/metrics/order`, `/metrics/audit-log` (Prometheus format).
- **Structured logs**: Zap-based JSON logs with request IDs.

## 9) CI/CD
- **GitHub Actions** runs:
  - `make test`
  - Link checks (lychee)

## 10) Deployment
- **Docker Compose** for local development.
- Services can be built and run independently.
- Environment variables drive configuration (see `.env.example`).

## 11) API Docs
- **Swagger UI** per service:
  - Gateway: `http://localhost:8080/swagger/user/`, `http://localhost:8080/swagger/order/`, `http://localhost:8080/swagger/audit-log/`
  - Direct: `http://localhost:8081/swagger/index.html`, `http://localhost:8082/swagger/index.html`, `http://localhost:8083/swagger/index.html`

## 12) Risks and Future Enhancements
- Harden gateway with WAF and rate limiting policies.
- Add distributed tracing (Jaeger/Zipkin).
- Add caching layer (Redis) for read-heavy operations.
- Expand CI with security scanning and linting.
