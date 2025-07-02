package services

import (
	"context"
	"math"
	"slices"
	"time"
)

// RetryConfig holds configuration for retry logic
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []ErrorType
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:    3,
		InitialDelay:  500 * time.Millisecond,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []ErrorType{
			ErrorTypeNetwork,
			ErrorTypeTimeout,
			ErrorTypeRateLimit,
		},
	}
}

// RetryableFunc represents a function that can be retried
type RetryableFunc func() error

// WithRetry executes a function with retry logic based on the provided configuration
func WithRetry(ctx context.Context, config RetryConfig, fn RetryableFunc) error {
	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return NewTimeoutError("retry cancelled due to context", "", ctx.Err())
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if this is the last attempt
		if attempt == config.MaxRetries {
			break
		}

		// Check if the error is retryable
		if !isRetryableError(err, config.RetryableErrors) {
			return err // Don't retry non-retryable errors
		}

		// Calculate delay for next attempt
		delay := calculateDelay(attempt, config)

		// Wait before retrying
		select {
		case <-ctx.Done():
			return NewTimeoutError("retry cancelled during backoff", "", ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}

// isRetryableError checks if an error should be retried
func isRetryableError(err error, retryableTypes []ErrorType) bool {
	if ghErr, ok := err.(*GitHubError); ok {
		if slices.Contains(retryableTypes, ghErr.Type) {
			return true
		}
		// Also check if the error itself says it's retryable
		return ghErr.IsRetryable()
	}
	return false
}

// calculateDelay calculates the delay for the next retry attempt using exponential backoff
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt))

	// Cap at maximum delay
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}

	return time.Duration(delay)
}
