# Enterprise Microservice System

A production-ready, enterprise-grade microservice system built with Go, featuring clean architecture, comprehensive observability, and industry-standard patterns.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Features](#features)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Development](#development)
- [Deployment](#deployment)
- [Monitoring](#monitoring)

## Architecture Overview

This system follows a **microservice architecture** with clean separation of concerns:

```
┌─────────────────┐         ┌─────────────────┐
│   User Service  │         │  Order Service  │
│   Port: 8081    │◄────────┤   Port: 8082    │
└────────┬────────┘         └────────┬────────┘
         │                           │
    ┌────▼────┐                 ┌────▼────┐
    │ User DB │                 │ Order DB│
    └─────────┘                 └─────────┘
         │                           │
    ┌────▼───────────────────────────▼────┐
    │         Prometheus Metrics          │
    └─────────────────────────────────────┘
```

### Key Architectural Decisions

- **Clean Architecture**: Each service follows the clean architecture pattern with clear separation between transport, service, domain, and repository layers
- **Microservice Independence**: Each service has its own database (database-per-service pattern)
- **Circuit Breaker**: Order service uses circuit breaker pattern for resilient communication with user service
- **Graceful Degradation**: Services continue operating even when dependencies are unavailable
- **Observability First**: Comprehensive logging, metrics, and health checks built-in

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Logging | Zap (structured logging) |
| Metrics | Prometheus |
| Circuit Breaker | gobreaker |
| Rate Limiting | golang.org/x/time/rate (token bucket) |
| Containerization | Docker, Docker Compose |
| Hot Reload | Air |

## Project Structure

```
enterprise-microservice-system/
├── common/                          # Shared libraries
│   ├── circuitbreaker/             # Circuit breaker implementation
│   ├── errors/                     # Custom error types
│   ├── logger/                     # Structured logging
│   ├── metrics/                    # Prometheus metrics
│   ├── middleware/                 # HTTP middleware
│   │   ├── cors.go                # CORS handling
│   │   ├── logger.go              # Request logging
│   │   ├── rate_limiter.go        # Token bucket rate limiter
│   │   ├── recovery.go            # Panic recovery
│   │   └── request_id.go          # Request ID tracking
│   └── response/                   # Standard API responses
├── services/
│   ├── user-service/              # User management service
│   │   ├── cmd/
│   │   │   └── main.go           # Entry point
│   │   ├── internal/
│   │   │   ├── api/              # Route definitions
│   │   │   ├── config/           # Configuration
│   │   │   ├── handler/          # HTTP handlers
│   │   │   ├── model/            # Domain models
│   │   │   ├── repository/       # Data access layer
│   │   │   └── service/          # Business logic
│   │   ├── tests/                # Unit tests
│   │   ├── Dockerfile            # Container definition
│   │   └── .air.toml             # Hot reload config
│   └── order-service/            # Order management service
│       ├── cmd/
│       │   └── main.go           # Entry point
│       ├── internal/
│       │   ├── api/              # Route definitions
│       │   ├── client/           # Inter-service communication
│       │   ├── config/           # Configuration
│       │   ├── handler/          # HTTP handlers
│       │   ├── model/            # Domain models
│       │   ├── repository/       # Data access layer
│       │   └── service/          # Business logic
│       ├── Dockerfile            # Container definition
│       └── .air.toml             # Hot reload config
├── monitoring/
│   └── prometheus.yml            # Prometheus configuration
├── docker-compose.yml            # Container orchestration
├── Makefile                      # Build and run commands
├── go.mod                        # Go dependencies
├── .env.example                  # Environment variables template
└── README.md                     # This file
```

### Layer Responsibilities

1. **Handler Layer**: Receives HTTP requests, validates input, calls service layer, returns HTTP responses
2. **Service Layer**: Contains business logic, orchestrates operations, handles transactions
3. **Repository Layer**: Handles database operations, implements data access patterns
4. **Model Layer**: Defines domain entities and data transfer objects

## Features

### 1. CRUD Operations
- Full RESTful APIs for users and orders
- Request validation using Gin's validator
- Pagination and filtering support
- Soft delete functionality

### 2. Database Management
- PostgreSQL with GORM ORM
- Automatic migrations
- Connection pooling
- Transaction support
- Separate databases per service

### 3. Rate Limiting
- Token bucket algorithm implementation
- Per-IP rate limiting
- Configurable limits
- Thread-safe using sync.Map

### 4. Circuit Breaker
- Protects inter-service calls
- States: Closed, Half-Open, Open
- Configurable failure threshold
- Automatic recovery attempts
- Graceful fallback responses

### 5. Concurrency Features
- Goroutines for background tasks
- Worker pool pattern
- Mutex for thread-safe operations
- Channel-based communication
- Context-aware request handling

### 6. Metrics & Observability
- Prometheus-compatible metrics endpoint
- Request count, latency, error rate tracking
- Circuit breaker state monitoring
- Health check endpoints
- Structured JSON logging

### 7. Middleware Stack
- CORS handling
- Request ID generation
- Panic recovery
- Request logging
- Rate limiting
- Metrics collection

### 8. Developer Experience
- Hot reload with Air
- Comprehensive Makefile
- Docker Compose for local development
- Environment variable configuration
- Clear error messages

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- Make (optional but recommended)
- PostgreSQL (if running locally without Docker)

### Quick Start with Docker

1. Clone the repository:
```bash
cd enterprise-microservice-system
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Start all services with Docker Compose:
```bash
make docker-up
```

This will start:
- User Service (http://localhost:8081)
- Order Service (http://localhost:8082)
- PostgreSQL databases
- Prometheus (http://localhost:9090)

4. Verify services are running:
```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

### Local Development Setup

1. Install dependencies:
```bash
make deps
```

2. Install development tools:
```bash
make install-tools
```

3. Start databases:
```bash
docker-compose up -d user-db order-db
```

4. Run services locally:

In terminal 1:
```bash
make run-user
```

In terminal 2:
```bash
make run-order
```

Or run both concurrently:
```bash
make run
```

### Hot Reload Development

```bash
# Terminal 1
make dev-user

# Terminal 2
make dev-order
```

## Configuration

### Environment Variables

#### User Service
| Variable | Description | Default |
|----------|-------------|---------|
| USER_SERVICE_PORT | HTTP port | 8081 |
| USER_SERVICE_DB_HOST | Database host | localhost |
| USER_SERVICE_DB_PORT | Database port | 5432 |
| USER_SERVICE_DB_USER | Database user | postgres |
| USER_SERVICE_DB_PASSWORD | Database password | postgres |
| USER_SERVICE_DB_NAME | Database name | userdb |
| USER_SERVICE_LOG_LEVEL | Log level (debug/info/warn/error) | info |
| USER_SERVICE_RATE_LIMIT | Requests per second | 100 |

#### Order Service
| Variable | Description | Default |
|----------|-------------|---------|
| ORDER_SERVICE_PORT | HTTP port | 8082 |
| ORDER_SERVICE_DB_HOST | Database host | localhost |
| ORDER_SERVICE_DB_PORT | Database port | 5432 |
| ORDER_SERVICE_DB_USER | Database user | postgres |
| ORDER_SERVICE_DB_PASSWORD | Database password | postgres |
| ORDER_SERVICE_DB_NAME | Database name | orderdb |
| ORDER_SERVICE_LOG_LEVEL | Log level | info |
| ORDER_SERVICE_RATE_LIMIT | Requests per second | 100 |
| ORDER_SERVICE_USER_SERVICE_URL | User service URL | http://localhost:8081 |

#### Circuit Breaker
| Variable | Description | Default |
|----------|-------------|---------|
| CIRCUIT_BREAKER_MAX_REQUESTS | Max requests in half-open state | 3 |
| CIRCUIT_BREAKER_INTERVAL | Failure counting interval (seconds) | 60 |
| CIRCUIT_BREAKER_TIMEOUT | Open to half-open timeout (seconds) | 30 |

## API Documentation

### User Service (Port 8081)

#### Create User
```bash
POST /api/v1/users
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe",
  "age": 30
}
```

#### Get User
```bash
GET /api/v1/users/{id}
```

#### List Users
```bash
GET /api/v1/users?page=1&page_size=10&search=john&active=true
```

#### Update User
```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "Jane Doe",
  "age": 31,
  "active": true
}
```

#### Delete User
```bash
DELETE /api/v1/users/{id}
```

### Order Service (Port 8082)

#### Create Order
```bash
POST /api/v1/orders
Content-Type: application/json

{
  "user_id": 1,
  "product_id": "PROD-123",
  "quantity": 2,
  "total_price": 99.99
}
```

#### Get Order
```bash
GET /api/v1/orders/{id}
```

#### List Orders
```bash
GET /api/v1/orders?page=1&page_size=10&user_id=1&status=pending
```

#### Update Order
```bash
PUT /api/v1/orders/{id}
Content-Type: application/json

{
  "status": "confirmed"
}
```

#### Delete Order
```bash
DELETE /api/v1/orders/{id}
```

### Order Status Values
- `pending` - Order created
- `confirmed` - Order confirmed
- `shipped` - Order shipped
- `delivered` - Order delivered
- `cancelled` - Order cancelled

### Health & Metrics

```bash
# Health checks
GET /health

# Prometheus metrics
GET /metrics
```

## Testing

### Run All Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test
```

### Run Service-Specific Tests
```bash
make test-user
make test-order
```

### Run Linter
```bash
make lint
```

## Development

### Makefile Commands

```bash
make help              # Show all available commands
make build             # Build all services
make run              # Run all services
make run-user         # Run user service
make run-order        # Run order service
make test             # Run all tests
make lint             # Run linter
make docker-up        # Start Docker containers
make docker-down      # Stop Docker containers
make docker-clean     # Remove all Docker resources
make clean            # Clean build artifacts
make deps             # Download dependencies
make format           # Format code
make vet              # Run go vet
```

### Code Style Guidelines

1. Follow standard Go conventions and idioms
2. Use meaningful variable and function names
3. Keep functions small and focused
4. Write tests for business logic
5. Use context for cancellation and timeouts
6. Handle errors explicitly
7. Use structured logging with appropriate fields
8. Document public APIs and complex logic

### Adding a New Service

1. Create service directory under `services/`
2. Follow the same structure as existing services
3. Update `docker-compose.yml` to include new service
4. Update `Makefile` with new service commands
5. Add service-specific configuration in `.env.example`
6. Update this README with service documentation

## Deployment

### Building Production Images

```bash
# Build user service
docker build -t user-service:latest -f services/user-service/Dockerfile .

# Build order service
docker build -t order-service:latest -f services/order-service/Dockerfile .
```

### Production Considerations

1. **Security**
   - Use secrets management (HashiCorp Vault, AWS Secrets Manager)
   - Enable TLS/HTTPS
   - Implement authentication and authorization
   - Regular security audits

2. **Database**
   - Use managed PostgreSQL (RDS, Cloud SQL)
   - Enable automatic backups
   - Configure read replicas for scaling
   - Implement connection pooling

3. **Monitoring**
   - Deploy Prometheus and Grafana
   - Set up alerting rules
   - Implement distributed tracing (Jaeger, Zipkin)
   - Log aggregation (ELK Stack, Loki)

4. **Scalability**
   - Use Kubernetes for orchestration
   - Implement horizontal pod autoscaling
   - Use load balancers
   - Configure resource limits

5. **High Availability**
   - Multi-region deployment
   - Health check configuration
   - Graceful shutdown handling
   - Circuit breaker patterns

## Monitoring

### Prometheus Metrics

Access Prometheus at `http://localhost:9090`

Key metrics to monitor:

1. **Request Metrics**
   - `user_service_requests_total` - Total HTTP requests
   - `order_service_requests_total` - Total HTTP requests
   - `*_request_duration_seconds` - Request latency

2. **Error Metrics**
   - `*_errors_total{type="server_error"}` - 5xx errors
   - `*_errors_total{type="client_error"}` - 4xx errors

3. **Circuit Breaker**
   - `order_service_circuit_breaker_state` - Circuit state (0=closed, 1=half-open, 2=open)

### Example Prometheus Queries

```promql
# Request rate per second
rate(user_service_requests_total[5m])

# 95th percentile latency
histogram_quantile(0.95, rate(user_service_request_duration_seconds_bucket[5m]))

# Error rate
sum(rate(user_service_errors_total[5m])) / sum(rate(user_service_requests_total[5m]))

# Circuit breaker open events
changes(order_service_circuit_breaker_state{service="user-service"}[5m]) > 0
```

### Health Checks

```bash
# Check service health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## Troubleshooting

### Common Issues

1. **Port Already in Use**
```bash
# Find and kill process using port 8081
lsof -ti:8081 | xargs kill -9
```

2. **Database Connection Failed**
```bash
# Check if PostgreSQL is running
docker-compose ps

# View database logs
docker-compose logs user-db
```

3. **Circuit Breaker Open**
   - Check if user service is healthy
   - View circuit breaker metrics in Prometheus
   - Circuit will auto-recover after timeout period

4. **Rate Limit Exceeded**
   - Adjust `RATE_LIMIT` environment variable
   - Implement user-based rate limiting instead of IP-based

## Contributing

1. Follow the existing code structure
2. Write tests for new features
3. Update documentation
4. Follow Go best practices
5. Run linter before committing

## License

This project is licensed under the MIT License.

## Contact

For questions and support, please open an issue on GitHub.

---

Built with Go and industry best practices for production-ready microservices.
