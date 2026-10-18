package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type RefServiceInstanceOption struct {
	InstanceUID string
	DisplayName string
}

// ResolveRefSvcParamValue normalizes a refSvcId parameter to a UUID if possible.
// If value already looks like a UUID, it is returned as-is.
func (c *UniversalClient) ResolveRefSvcParamValue(ctx context.Context, refSvcId int, value string) (string, error) {
	if c == nil {
		return "", nil
	}
	if refSvcId <= 0 {
		return value, nil
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	// FIX(uuid-case): UUID-like values are normalised to lowercase so that
	// plan (user input may be uppercase) and API response (always lowercase)
	// produce identical strings and do not cause Terraform"s plan/state
	// inconsistency error.
	if isUUIDLike(trimmed) {
		return strings.ToLower(trimmed), nil
	}
	uid, err := c.findInstanceUidByDisplayNameRefSvc(ctx, refSvcId, trimmed)
	if err != nil {
		return "", err
	}
	return uid, nil
}

// ResolveRefSvcParamDisplayName maps UUID values back to display_name when possible.
func (c *UniversalClient) ResolveRefSvcParamDisplayName(ctx context.Context, refSvcId int, value string) (string, error) {
	if c == nil {
		return "", nil
	}
	if refSvcId <= 0 {
		return value, nil
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	if !isUUIDLike(trimmed) {
		return trimmed, nil
	}
	// Prefer direct lookup by instance UID to avoid paging the full instances list.
	displayName, err := c.getInstanceDisplayNameByUidRefSvc(ctx, refSvcId, trimmed)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(displayName) == "" {
		// Fallback to list scan if the direct lookup did not match by serviceId.
		displayName, err = c.findInstanceDisplayNameByUidRefSvc(ctx, refSvcId, trimmed)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(displayName) == "" {
			return trimmed, nil
		}
	}
	return displayName, nil
}

// ListRefServiceInstances returns active instance options for a referenced service.
func (c *UniversalClient) ListRefServiceInstances(ctx context.Context, serviceId int) ([]RefServiceInstanceOption, error) {
	if c == nil || serviceId <= 0 {
		return nil, nil
	}

	var out []RefServiceInstanceOption
	page := 1
	for {
		// isDeleted=false — API-уровень фильтрации, не тратим трафик на удалённые инстансы
		path := fmt.Sprintf("/instances?page=%d&size=100&isDeleted=false", page)
		respBody, _, err := c.doRequest(ctx, "GET", path, nil)
		if err != nil {
			return nil, err
		}

		var res struct {
			Results []struct {
				InstanceUid     string `json:"instanceUid"`
				DisplayName     string `json:"displayName"`
				ServiceId       int    `json:"serviceId"`
				IsDeleted       bool   `json:"isDeleted"`
				ExplainedStatus string `json:"explainedStatus"`
			} `json:"results"`
		}
		if err := json.Unmarshal(respBody, &res); err != nil {
			return nil, err
		}

		if len(res.Results) == 0 {
			break
		}

		for _, item := range res.Results {
			if item.ServiceId != serviceId {
				continue
			}
			// показываем только running инстансы — deleted/not_created не пригодны к использованию
			if item.IsDeleted || !strings.EqualFold(strings.TrimSpace(item.ExplainedStatus), "running") {
				continue
			}
			uid := strings.TrimSpace(item.InstanceUid)
			name := strings.TrimSpace(item.DisplayName)
			if uid == "" {
				continue
			}
			out = append(out, RefServiceInstanceOption{InstanceUID: uid, DisplayName: name})
		}

		page++
		if page > 100 {
			break
		}
	}

	return out, nil
}

// ListAvailableResourceRealms returns realm names available for the target service.
func (c *UniversalClient) ListAvailableResourceRealms(ctx context.Context, serviceId int) ([]string, error) {
	if c == nil || serviceId <= 0 {
		return nil, nil
	}

	respBody, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/resourceRealms/available?svcId=%d", serviceId), nil)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, err
	}

	parseArray := func(raw json.RawMessage) []string {
		var values []string
		_ = json.Unmarshal(raw, &values)
		return values
	}

	parseString := func(raw json.RawMessage) []string {
		var single string
		if err := json.Unmarshal(raw, &single); err != nil {
			return nil
		}
		single = strings.TrimSpace(single)
		if single == "" {
			return nil
		}
		parts := strings.Split(single, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) == 0 {
			result = append(result, single)
		}
		return result
	}

	values := parseArray(payload.Results)
	if len(values) == 0 {
		values = parseString(payload.Results)
	}

	out := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out, nil
}
