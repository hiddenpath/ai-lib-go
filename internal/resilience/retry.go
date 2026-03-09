// Package resilience provides retry and rate limiting.
// 重试策略和限流，遵循 AI-Protocol retry_policy 规范。
package resilience

import (
	"context"
	"math/rand"
	"time"
)

// RetryPolicy defines retry behavior.
// 重试策略配置。
type RetryPolicy struct {
	Strategy      Strategy     // exponential_backoff, fixed, none
	MinDelay      time.Duration
	MaxDelay      time.Duration
	MaxAttempts   int
	Jitter        JitterType   // none, full, decorrelated
	RetryOnStatus []int
	RetryOnErrors []string
}

// Strategy defines the retry strategy.
type Strategy string

const (
	StrategyExponential Strategy = "exponential_backoff"
	StrategyFixed       Strategy = "fixed"
	StrategyNone        Strategy = "none"
)

// JitterType defines the jitter type.
type JitterType string

const (
	JitterNone        JitterType = "none"
	JitterFull        JitterType = "full"
	JitterDecorrelated JitterType = "decorrelated"
)

// DefaultRetryPolicy returns the default retry policy.
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		Strategy:      StrategyExponential,
		MinDelay:      1 * time.Second,
		MaxDelay:      60 * time.Second,
		MaxAttempts:   3,
		Jitter:        JitterFull,
		RetryOnStatus: []int{429, 500, 502, 503, 504},
	}
}

// ShouldRetry determines if a request should be retried.
func (p *RetryPolicy) ShouldRetry(attempt int, statusCode int, errCode string) bool {
	if attempt >= p.MaxAttempts {
		return false
	}

	// Check status codes
	for _, s := range p.RetryOnStatus {
		if s == statusCode {
			return true
		}
	}

	// Check error codes
	for _, e := range p.RetryOnErrors {
		if e == errCode {
			return true
		}
	}

	return false
}

// CalculateDelay calculates the delay for the next retry.
func (p *RetryPolicy) CalculateDelay(attempt int) time.Duration {
	if p.Strategy == StrategyNone {
		return 0
	}

	var delay time.Duration

	switch p.Strategy {
	case StrategyExponential:
		// Exponential backoff: min_delay * 2^attempt
		delay = p.MinDelay * time.Duration(1<<uint(attempt))
		if delay > p.MaxDelay {
			delay = p.MaxDelay
		}
	case StrategyFixed:
		delay = p.MinDelay
	}

	// Apply jitter
	switch p.Jitter {
	case JitterFull:
		delay = time.Duration(rand.Float64() * float64(delay))
	case JitterDecorrelated:
		// Decorrelated jitter: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
		delay = time.Duration(rand.Float64() * float64(delay*3-p.MinDelay) + float64(p.MinDelay))
	}

	return delay
}

// RetryableFunc is a function that can be retried.
type RetryableFunc func(ctx context.Context) error

// Retry executes a function with retry logic.
func Retry(ctx context.Context, policy *RetryPolicy, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		// In real implementation, extract status code from error
		if !policy.ShouldRetry(attempt, 0, "") {
			return err
		}

		// Calculate delay
		delay := policy.CalculateDelay(attempt)
		if delay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return lastErr
}

// RateLimiter implements rate limiting.
type RateLimiter struct {
	requestsPerMinute int
	tokensPerMinute   int
	tokensUsed        int
	requestsUsed      int
	windowStart       time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(requestsPerMinute, tokensPerMinute int) *RateLimiter {
	return &RateLimiter{
		requestsPerMinute: requestsPerMinute,
		tokensPerMinute:   tokensPerMinute,
		windowStart:       time.Now(),
	}
}

// Allow checks if a request is allowed.
func (r *RateLimiter) Allow(tokens int) bool {
	now := time.Now()
	windowEnd := r.windowStart.Add(time.Minute)

	// Reset window if expired
	if now.After(windowEnd) {
		r.windowStart = now
		r.requestsUsed = 0
		r.tokensUsed = 0
		return true
	}

	// Check limits
	if r.requestsPerMinute > 0 && r.requestsUsed >= r.requestsPerMinute {
		return false
	}
	if r.tokensPerMinute > 0 && r.tokensUsed+tokens > r.tokensPerMinute {
		return false
	}

	// Increment counters
	r.requestsUsed++
	r.tokensUsed += tokens
	return true
}

// WaitUntilAllowed blocks until a request is allowed.
func (r *RateLimiter) WaitUntilAllowed(ctx context.Context, tokens int) error {
	for {
		if r.Allow(tokens) {
			return nil
		}

		// Calculate wait time
		waitTime := time.Until(r.windowStart.Add(time.Minute))
		if waitTime <= 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}
	}
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	failures      int
	state         CircuitState
	lastFailTime  time.Time
}

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        CircuitClosed,
	}
}

// Allow checks if requests are allowed.
func (c *CircuitBreaker) Allow() bool {
	switch c.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		// Check if reset timeout has passed
		if time.Since(c.lastFailTime) > c.resetTimeout {
			c.state = CircuitHalfOpen
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful request.
func (c *CircuitBreaker) RecordSuccess() {
	if c.state == CircuitHalfOpen {
		c.state = CircuitClosed
		c.failures = 0
	}
}

// RecordFailure records a failed request.
func (c *CircuitBreaker) RecordFailure() {
	c.failures++
	c.lastFailTime = time.Now()

	if c.failures >= c.maxFailures {
		c.state = CircuitOpen
	}
}

// State returns the current state.
func (c *CircuitBreaker) State() CircuitState {
	return c.state
}
