package client

import (
	"context"
	"encoding/json"
	"enterprise-microservice-system/common/circuitbreaker"
	"enterprise-microservice-system/common/errors"
	"enterprise-microservice-system/services/order-service/internal/model"
	"fmt"
	"net/http"
	"time"
)

// UserClient handles communication with the user service
type UserClient struct {
	baseURL        string
	client         *http.Client
	circuitBreaker *circuitbreaker.CircuitBreaker
}

// NewUserClient creates a new user service client
func NewUserClient(baseURL string, cb *circuitbreaker.CircuitBreaker) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		circuitBreaker: cb,
	}
}

// UserResponse represents the response from user service
type UserResponse struct {
	Success bool        `json:"success"`
	Data    *model.User `json:"data"`
	Error   *ErrorInfo  `json:"error"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetUser retrieves a user by ID from user service
func (c *UserClient) GetUser(ctx context.Context, userID uint) (*model.User, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", c.baseURL, userID)

	// Execute request with circuit breaker protection
	result, err := c.circuitBreaker.ExecuteWithContext(ctx, func() (interface{}, error) {
		return c.doGetRequest(ctx, url)
	})

	if err != nil {
		return nil, err
	}

	user := result.(*model.User)
	return user, nil
}

// doGetRequest performs the actual HTTP GET request
func (c *UserClient) doGetRequest(ctx context.Context, url string) (*model.User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.NewInternal("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.NewInternal("failed to call user service", err)
	}
	defer resp.Body.Close()

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errors.NewNotFound("user")
		}
		return nil, errors.New(
			errors.ErrCodeServiceUnavail,
			fmt.Sprintf("user service returned status %d", resp.StatusCode),
			nil,
		)
	}

	// Parse response
	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, errors.NewInternal("failed to decode response", err)
	}

	if !userResp.Success {
		return nil, errors.New(
			userResp.Error.Code,
			userResp.Error.Message,
			nil,
		)
	}

	return userResp.Data, nil
}

// GetCircuitBreakerState returns the current circuit breaker state
func (c *UserClient) GetCircuitBreakerState() float64 {
	return c.circuitBreaker.StateAsFloat()
}
