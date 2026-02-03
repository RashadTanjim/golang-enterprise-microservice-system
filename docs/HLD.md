# High-Level Design (HLD) – Enterprise Microservice System

## 1) Overview
This system is a production-ready microservice architecture built with Go. It consists of three core services (User Service, Order Service, and Repository Service), each with its own database, shared libraries for cross-cutting concerns, and supporting infrastructure for observability, authentication, and CI.

## 2) Goals
- Provide clean, maintainable microservices with clear boundaries.
- Support full CRUD operations for users, orders, and repositories.
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

- **Repository Service** (port 8083)
  - Manages repository CRUD and validation.
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
- **User DB**: PostgreSQL (userdb)
- **Order DB**: PostgreSQL (orderdb)
- **Repository DB**: PostgreSQL (repositorydb)

Each service owns its database (database-per-service).

## 5) Data Flow

### 5.1 Auth & Authorization
1. Client requests token from `POST /api/v1/auth/token` (User Service).
2. User Service validates client credentials and returns JWT.
3. Client calls protected endpoints with `Authorization: Bearer <token>`.
4. Middleware validates token and enforces role-based access.

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
- Repository Service uses GORM auto-migrations on startup.

## 8) Observability
- **Health checks**: `/health/user`, `/health/order`, `/health/repository` via gateway.
- **Metrics**: `/metrics/user`, `/metrics/order`, `/metrics/repository` (Prometheus format).
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
  - Gateway: `http://localhost:8080/swagger/user/`, `http://localhost:8080/swagger/order/`, `http://localhost:8080/swagger/repository/`
  - Direct: `http://localhost:8081/swagger/index.html`, `http://localhost:8082/swagger/index.html`, `http://localhost:8083/swagger/index.html`

## 12) Risks and Future Enhancements
- Harden gateway with WAF and rate limiting policies.
- Add distributed tracing (Jaeger/Zipkin).
- Add caching layer (Redis) for read-heavy operations.
- Expand CI with security scanning and linting.
