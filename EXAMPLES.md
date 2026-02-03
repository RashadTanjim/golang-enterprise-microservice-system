# API Usage Examples

This document provides practical examples of how to use the microservice APIs.

## Quick Start Example

Start the services:
```bash
make docker-up
# Wait for services to be ready (about 10 seconds)
```

Get an access token and export it for the examples:
```bash
curl -X POST http://localhost:8081/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "admin",
    "client_secret": "admin123"
  }'

TOKEN="<paste access token here>"
```

All requests below require:
```
Authorization: Bearer $TOKEN
```

## User Service Examples

### 1. Create Users

```bash
# Create first user
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "alice@example.com",
    "name": "Alice Johnson",
    "age": 28
  }'

# Create second user
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "bob@example.com",
    "name": "Bob Smith",
    "age": 35
  }'

# Create third user
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "email": "carol@example.com",
    "name": "Carol Williams",
    "age": 42
  }'
```

### 2. List All Users

```bash
curl http://localhost:8081/api/v1/users \
  -H "Authorization: Bearer $TOKEN"

# With pagination
curl "http://localhost:8081/api/v1/users?page=1&page_size=2" \
  -H "Authorization: Bearer $TOKEN"

# Search by name
curl "http://localhost:8081/api/v1/users?search=alice" \
  -H "Authorization: Bearer $TOKEN"

# Filter by status
curl "http://localhost:8081/api/v1/users?status=active" \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Get Specific User

```bash
# Get user with ID 1
curl http://localhost:8081/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Update User

```bash
# Update user name and age
curl -X PUT http://localhost:8081/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Alice Cooper",
    "age": 29
  }'

# Deactivate user
curl -X PUT http://localhost:8081/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "inactive"
  }'
```

### 5. Delete User

```bash
curl -X DELETE http://localhost:8081/api/v1/users/3 \
  -H "Authorization: Bearer $TOKEN"
```

## Order Service Examples

### 1. Create Orders

```bash
# Create order for user 1
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": 1,
    "product_id": "LAPTOP-001",
    "quantity": 1,
    "total_price": 1299.99
  }'

# Create another order
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": 1,
    "product_id": "MOUSE-003",
    "quantity": 2,
    "total_price": 49.98
  }'

# Create order for user 2
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": 2,
    "product_id": "KEYBOARD-005",
    "quantity": 1,
    "total_price": 159.99
  }'
```

**Note**: The order service will call the user service to validate the user exists and is active. This demonstrates inter-service communication with circuit breaker protection.

### 2. List Orders

```bash
# List all orders
curl http://localhost:8082/api/v1/orders \
  -H "Authorization: Bearer $TOKEN"

# Filter by user
curl "http://localhost:8082/api/v1/orders?user_id=1" \
  -H "Authorization: Bearer $TOKEN"

# Filter by order status
curl "http://localhost:8082/api/v1/orders?order_status=pending" \
  -H "Authorization: Bearer $TOKEN"

# Filter by product
curl "http://localhost:8082/api/v1/orders?product_id=LAPTOP-001" \
  -H "Authorization: Bearer $TOKEN"

# Pagination
curl "http://localhost:8082/api/v1/orders?page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Get Specific Order

```bash
# Get order with ID 1
curl http://localhost:8082/api/v1/orders/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Note**: The response includes user data fetched from the user service.

### 4. Update Order Status

```bash
# Confirm order
curl -X PUT http://localhost:8082/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_status": "confirmed"
  }'

# Ship order
curl -X PUT http://localhost:8082/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_status": "shipped"
  }'

# Deliver order
curl -X PUT http://localhost:8082/api/v1/orders/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_status": "delivered"
  }'

# Cancel order
curl -X PUT http://localhost:8082/api/v1/orders/2 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_status": "cancelled"
  }'
```

### 5. Delete Order

```bash
curl -X DELETE http://localhost:8082/api/v1/orders/3 \
  -H "Authorization: Bearer $TOKEN"
```

## Health & Monitoring

### Health Checks

```bash
# User service health
curl http://localhost:8081/health

# Order service health
curl http://localhost:8082/health
```

Expected response:
```json
{
  "service": "user-service",
  "status": "healthy"
}
```

### Metrics

```bash
# User service metrics
curl http://localhost:8081/metrics

# Order service metrics
curl http://localhost:8082/metrics
```

### Prometheus Dashboard

Access Prometheus at: http://localhost:9090

Example queries:
```promql
# Request rate
rate(user_service_requests_total[1m])

# Error rate
sum(rate(order_service_errors_total[5m])) by (type)

# Circuit breaker state
order_service_circuit_breaker_state{service="user-service"}
```

## Testing Circuit Breaker

### 1. Normal Operation

```bash
# Create an order (should succeed)
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "product_id": "TEST-001",
    "quantity": 1,
    "total_price": 99.99
  }'
```

### 2. Simulate User Service Failure

```bash
# Stop user service
docker-compose stop user-service

# Try to create an order (will fail after timeout, then circuit opens)
curl -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "product_id": "TEST-002",
    "quantity": 1,
    "total_price": 99.99
  }'
```

After several failures, you'll see:
```json
{
  "success": false,
  "error": {
    "code": "CIRCUIT_OPEN",
    "message": "circuit breaker open for service: user-service"
  }
}
```

### 3. Check Circuit Breaker Metrics

```bash
curl http://localhost:8082/metrics | grep circuit_breaker_state
```

### 4. Restore User Service

```bash
# Start user service again
docker-compose start user-service

# Wait for circuit breaker to transition to half-open (30 seconds)
# Then try creating an order again
```

## Testing Rate Limiting

```bash
# Install apache bench (if not installed)
# macOS: brew install httpd
# Ubuntu: sudo apt-get install apache2-utils

# Send 200 requests with 10 concurrent connections
ab -n 200 -c 10 -H "Content-Type: application/json" \
  http://localhost:8081/api/v1/users

# Check for 429 (Too Many Requests) responses
```

Or with a simple bash loop:
```bash
# Send 150 requests rapidly
for i in {1..150}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/api/v1/users
done | sort | uniq -c
```

You should see a mix of 200 (OK) and 429 (Too Many Requests) responses.

## Error Handling Examples

### 1. Validation Error

```bash
# Missing required field
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

Response:
```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
  }
}
```

### 2. Invalid Email

```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "name": "Test User",
    "age": 25
  }'
```

### 3. Duplicate Email

```bash
# Create user
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "duplicate@example.com",
    "name": "First User",
    "age": 25
  }'

# Try to create another user with same email
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "duplicate@example.com",
    "name": "Second User",
    "age": 30
  }'
```

Response:
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "email already exists"
  }
}
```

### 4. Not Found Error

```bash
curl http://localhost:8081/api/v1/users/9999
```

Response:
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "user not found"
  }
}
```

## Complete Workflow Example

Here's a complete workflow demonstrating the system:

```bash
# 1. Create a user
USER_RESPONSE=$(curl -s -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "age": 30
  }')

echo "User created: $USER_RESPONSE"
USER_ID=$(echo $USER_RESPONSE | jq -r '.data.id')

# 2. Create an order for the user
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8082/api/v1/orders \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER_ID,
    \"product_id\": \"PROD-123\",
    \"quantity\": 2,
    \"total_price\": 199.98
  }")

echo "Order created: $ORDER_RESPONSE"
ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.data.id')

# 3. Get the order (includes user data)
curl -s http://localhost:8082/api/v1/orders/$ORDER_ID | jq '.'

# 4. Update order status
curl -s -X PUT http://localhost:8082/api/v1/orders/$ORDER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "order_status": "confirmed"
  }' | jq '.'

# 5. List all orders for the user
curl -s "http://localhost:8082/api/v1/orders?user_id=$USER_ID" | jq '.'
```

## Load Testing Example

```bash
# Install vegeta (load testing tool)
# macOS: brew install vegeta
# Ubuntu: sudo apt-get install vegeta

# Create a targets file
cat > targets.txt << EOF
POST http://localhost:8081/api/v1/users
Content-Type: application/json

{
  "email": "loadtest@example.com",
  "name": "Load Test",
  "age": 25
}
EOF

# Run load test: 10 requests/second for 30 seconds
vegeta attack -targets=targets.txt -rate=10 -duration=30s | vegeta report

# Generate plot
vegeta attack -targets=targets.txt -rate=10 -duration=30s | \
  vegeta plot > results.html
```

## Cleanup

```bash
# Stop all services
make docker-down

# Remove all data
make docker-clean
```
