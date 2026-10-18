package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOperationTimeoutValue = 30 * time.Minute
	defaultIdleTimeoutValue      = 30 * time.Minute
)

type OperationTimeouts struct {
	DefaultTimeout     time.Duration
	IdleTimeout        time.Duration
	ByOperationName    map[string]time.Duration
	ByServiceOperation map[int]map[string]time.Duration
}

type operationTimeoutsConfig struct {
	DefaultTimeout string            `json:"default_timeout"`
	IdleTimeout    string            `json:"idle_timeout"`
	Operations     map[string]string `json:"operations"`
	Overrides      map[string]string `json:"overrides"`
}

func DefaultOperationTimeouts() *OperationTimeouts {
	return &OperationTimeouts{
		DefaultTimeout: defaultOperationTimeoutValue,
		IdleTimeout:    defaultIdleTimeoutValue,
		ByOperationName: map[string]time.Duration{
			"create":          30 * time.Minute,
			"modify":          30 * time.Minute,
			"resume":          30 * time.Minute,
			"suspend":         30 * time.Minute,
			"delete_user":     10 * time.Minute,
			"create_user":     10 * time.Minute,
			"delete_database": 10 * time.Minute,
			"create_database": 10 * time.Minute,
		},
		ByServiceOperation: map[int]map[string]time.Duration{},
	}
}

func LoadOperationTimeoutsFromBytes(raw []byte) (*OperationTimeouts, error) {
	var cfg operationTimeoutsConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse operation timeouts config: %w", err)
	}

	result := DefaultOperationTimeouts()

	if strings.TrimSpace(cfg.DefaultTimeout) != "" {
		d, err := time.ParseDuration(strings.TrimSpace(cfg.DefaultTimeout))
		if err != nil {
			return nil, fmt.Errorf("invalid default_timeout value %q: %w", cfg.DefaultTimeout, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("invalid default_timeout value %q: must be > 0", cfg.DefaultTimeout)
		}
		result.DefaultTimeout = d
	}

	if strings.TrimSpace(cfg.IdleTimeout) != "" {
		d, err := time.ParseDuration(strings.TrimSpace(cfg.IdleTimeout))
		if err != nil {
			return nil, fmt.Errorf("invalid idle_timeout value %q: %w", cfg.IdleTimeout, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("invalid idle_timeout value %q: must be > 0", cfg.IdleTimeout)
		}
		result.IdleTimeout = d
	}

	for operation, timeoutRaw := range cfg.Operations {
		operationKey := strings.ToLower(strings.TrimSpace(operation))
		if operationKey == "" {
			continue
		}
		d, err := time.ParseDuration(strings.TrimSpace(timeoutRaw))
		if err != nil {
			return nil, fmt.Errorf("invalid timeout for operation %q: %w", operation, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("invalid timeout for operation %q: must be > 0", operation)
		}
		result.ByOperationName[operationKey] = d
	}

	for overrideKey, timeoutRaw := range cfg.Overrides {
		serviceID, operationKey, err := parseServiceOperationOverrideKey(overrideKey)
		if err != nil {
			return nil, err
		}
		d, err := time.ParseDuration(strings.TrimSpace(timeoutRaw))
		if err != nil {
			return nil, fmt.Errorf("invalid timeout for override %q: %w", overrideKey, err)
		}
		if d <= 0 {
			return nil, fmt.Errorf("invalid timeout for override %q: must be > 0", overrideKey)
		}
		if _, ok := result.ByServiceOperation[serviceID]; !ok {
			result.ByServiceOperation[serviceID] = map[string]time.Duration{}
		}
		result.ByServiceOperation[serviceID][operationKey] = d
	}

	return result, nil
}

func parseServiceOperationOverrideKey(raw string) (int, string, error) {
	key := strings.ToLower(strings.TrimSpace(raw))
	parts := strings.Split(key, ".")
	if len(parts) != 3 || parts[0] != "services" || strings.TrimSpace(parts[1]) == "" || strings.TrimSpace(parts[2]) == "" {
		return 0, "", fmt.Errorf("invalid override key %q: expected services.<serviceId>.<operation>", raw)
	}
	serviceID, err := strconv.Atoi(parts[1])
	if err != nil || serviceID <= 0 {
		return 0, "", fmt.Errorf("invalid override key %q: serviceId must be positive integer", raw)
	}
	return serviceID, parts[2], nil
}

func (c *UniversalClient) operationTimeoutFor(serviceID int, action string) time.Duration {
	if c == nil || c.OperationTimeouts == nil {
		return defaultOperationTimeoutValue
	}
	key := strings.ToLower(strings.TrimSpace(action))
	if serviceID > 0 && key != "" {
		if serviceOverrides, ok := c.OperationTimeouts.ByServiceOperation[serviceID]; ok {
			if d, ok := serviceOverrides[key]; ok && d > 0 {
				return d
			}
		}
	}
	if key != "" {
		if d, ok := c.OperationTimeouts.ByOperationName[key]; ok && d > 0 {
			return d
		}
	}
	if c.OperationTimeouts.DefaultTimeout > 0 {
		return c.OperationTimeouts.DefaultTimeout
	}
	return defaultOperationTimeoutValue
}

func (c *UniversalClient) idleTimeoutFor(serviceID int) time.Duration {
	if c == nil || c.OperationTimeouts == nil {
		return defaultIdleTimeoutValue
	}
	if serviceID > 0 {
		if serviceOverrides, ok := c.OperationTimeouts.ByServiceOperation[serviceID]; ok {
			if d, ok := serviceOverrides["idle"]; ok && d > 0 {
				return d
			}
		}
	}
	if c.OperationTimeouts.IdleTimeout > 0 {
		return c.OperationTimeouts.IdleTimeout
	}
	if c.OperationTimeouts.DefaultTimeout > 0 {
		return c.OperationTimeouts.DefaultTimeout
	}
	return defaultIdleTimeoutValue
}
