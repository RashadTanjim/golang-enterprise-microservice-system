package circuitbreaker

import (
	"context"
	"enterprise-microservice-system/common/errors"
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker wraps gobreaker with custom configuration
type CircuitBreaker struct {
	breaker     *gobreaker.CircuitBreaker
	serviceName string
}

// Config holds circuit breaker configuration
type Config struct {
	MaxRequests uint32        // Max requests allowed in half-open state
	Interval    time.Duration // Time period for counting failures
	Timeout     time.Duration // Time to wait before transitioning from open to half-open
}

// DefaultConfig returns default circuit breaker configuration
func DefaultConfig() Config {
	return Config{
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
	}
}

// New creates a new circuit breaker
func New(serviceName string, config Config) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open circuit if failure rate >= 50% and at least 3 requests
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("[CircuitBreaker] %s: State changed from %s to %s\n", name, from.String(), to.String())
		},
	}

	return &CircuitBreaker{
		breaker:     gobreaker.NewCircuitBreaker(settings),
		serviceName: serviceName,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() (interface{}, error)) (interface{}, error) {
	result, err := cb.breaker.Execute(fn)
	if err != nil {
		// Check if error is due to open circuit
		if err == gobreaker.ErrOpenState {
			return nil, errors.NewCircuitOpen(cb.serviceName)
		}
		return nil, err
	}
	return result, nil
}

// ExecuteWithContext runs the given function with circuit breaker protection and context
func (cb *CircuitBreaker) ExecuteWithContext(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	// Check context before executing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return cb.Execute(fn)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() gobreaker.State {
	return cb.breaker.State()
}

// StateAsFloat returns the state as a float for metrics
// 0 = Closed, 1 = Half-Open, 2 = Open
func (cb *CircuitBreaker) StateAsFloat() float64 {
	switch cb.breaker.State() {
	case gobreaker.StateClosed:
		return 0
	case gobreaker.StateHalfOpen:
		return 1
	case gobreaker.StateOpen:
		return 2
	default:
		return -1
	}
}

// Counts returns the current counts
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.breaker.Counts()
}
