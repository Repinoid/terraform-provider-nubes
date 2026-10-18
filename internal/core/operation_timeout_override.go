package core

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type operationTimeoutOverrideKey struct{}

// WithOperationTimeout attaches optional per-request operation timeout override to context.
// Empty value means no override and falls back to configured defaults.
func WithOperationTimeout(ctx context.Context, raw string) (context.Context, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ctx, nil
	}
	d, err := time.ParseDuration(trimmed)
	if err != nil {
		return nil, fmt.Errorf("некорректный operation_timeout %q: %w", raw, err)
	}
	if d <= 0 {
		return nil, fmt.Errorf("некорректный operation_timeout %q: значение должно быть больше 0", raw)
	}
	return context.WithValue(ctx, operationTimeoutOverrideKey{}, d), nil
}

func operationTimeoutOverrideFromContext(ctx context.Context) (time.Duration, bool) {
	if ctx == nil {
		return 0, false
	}
	v := ctx.Value(operationTimeoutOverrideKey{})
	d, ok := v.(time.Duration)
	if !ok || d <= 0 {
		return 0, false
	}
	return d, true
}

func (c *UniversalClient) operationTimeoutForContext(ctx context.Context, serviceID int, action string) time.Duration {
	if d, ok := operationTimeoutOverrideFromContext(ctx); ok {
		return d
	}
	return c.operationTimeoutFor(serviceID, action)
}
