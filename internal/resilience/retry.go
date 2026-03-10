// Package resilience provides retry policy utilities.
// 重试策略模块，提供指数退避与可重试判定。
package resilience

import (
	"context"
	"errors"
	"math"
	"time"
)

type Policy struct {
	MaxAttempts int
	MinDelay    time.Duration
	MaxDelay    time.Duration
}

func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		MinDelay:    time.Second,
		MaxDelay:    20 * time.Second,
	}
}

type TryFunc func(context.Context) error

func Execute(ctx context.Context, p Policy, fn TryFunc, retryable func(error) bool) error {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.MinDelay <= 0 {
		p.MinDelay = 200 * time.Millisecond
	}
	if p.MaxDelay < p.MinDelay {
		p.MaxDelay = p.MinDelay
	}

	var lastErr error
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := fn(ctx); err != nil {
			lastErr = err
			if attempt == p.MaxAttempts-1 || (retryable != nil && !retryable(err)) {
				return err
			}
			delay := backoff(p, attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			continue
		}
		return nil
	}
	if lastErr == nil {
		return errors.New("retry exhausted")
	}
	return lastErr
}

func backoff(p Policy, attempt int) time.Duration {
	m := math.Pow(2, float64(attempt))
	d := time.Duration(float64(p.MinDelay) * m)
	if d > p.MaxDelay {
		return p.MaxDelay
	}
	return d
}
