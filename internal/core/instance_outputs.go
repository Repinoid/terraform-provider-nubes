package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type InstanceStateDetails struct {
	Params    map[string]string
	Out       map[string]string
	RawParams map[string]interface{}
	RawOut    map[string]interface{}
	Vault     InstanceVaultMeta
}

type InstanceVaultMeta struct {
	Url      string
	UserPath string
	Fields   []string
}

// GetInstanceStateDetails returns instance.state params/out and vault metadata.
func (c *UniversalClient) GetInstanceStateDetails(ctx context.Context, instanceUid string) (*InstanceStateDetails, error) {
	if strings.TrimSpace(instanceUid) == "" {
		return nil, fmt.Errorf("missing instance uid")
	}

	respBody, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/instances/%s", instanceUid), nil)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}

	inst, ok := res["instance"].(map[string]interface{})
	if !ok {
		return &InstanceStateDetails{Params: map[string]string{}, Out: map[string]string{}}, nil
	}

	state, _ := inst["state"].(map[string]interface{})
	if state == nil {
		return &InstanceStateDetails{Params: map[string]string{}, Out: map[string]string{}}, nil
	}

	paramsRaw := extractStateRawMap(state["params"])
	outRaw := extractStateRawMap(state["out"])
	params := extractStateStringMap(paramsRaw)
	out := extractStateStringMap(outRaw)
	vault := extractVaultMeta(state["vault"])

	return &InstanceStateDetails{
		Params:    params,
		Out:       out,
		RawParams: paramsRaw,
		RawOut:    outRaw,
		Vault:     vault,
	}, nil
}

// GetInstanceVaultSecrets fetches all named secrets from the vault endpoint.
func (c *UniversalClient) GetInstanceVaultSecrets(ctx context.Context, instanceUid string, fields []string) (map[string]string, error) {
	secrets := make(map[string]string)
	if strings.TrimSpace(instanceUid) == "" || len(fields) == 0 {
		return secrets, nil
	}

	var errs []string
	for _, field := range fields {
		name := strings.TrimSpace(field)
		if name == "" {
			continue
		}
		value, err := c.GetInstanceVaultSecret(ctx, instanceUid, name)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", name, err))
			continue
		}
		secrets[name] = value
	}

	if len(errs) > 0 {
		return secrets, fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return secrets, nil
}

// GetInstanceVaultSecret fetches a single secret value by name.
func (c *UniversalClient) GetInstanceVaultSecret(ctx context.Context, instanceUid string, secretName string) (string, error) {
	if strings.TrimSpace(instanceUid) == "" {
		return "", fmt.Errorf("missing instance uid")
	}
	if strings.TrimSpace(secretName) == "" {
		return "", fmt.Errorf("missing secret name")
	}

	path := fmt.Sprintf("/instances/%s/vault/%s", instanceUid, url.PathEscape(secretName))
	respBody, _, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return "", err
	}

	return extractVaultSecretValue(respBody)
}

func extractStateRawMap(value interface{}) map[string]interface{} {
	raw, _ := value.(map[string]interface{})
	if raw == nil {
		return map[string]interface{}{}
	}
	return raw
}

func extractStateStringMap(raw map[string]interface{}) map[string]string {
	if raw == nil {
		return map[string]string{}
	}

	out := make(map[string]string, len(raw))
	for key, val := range raw {
		out[key] = normalizeInstanceParamValue(val)
	}
	return out
}

func extractVaultMeta(value interface{}) InstanceVaultMeta {
	raw, _ := value.(map[string]interface{})
	if raw == nil {
		return InstanceVaultMeta{}
	}

	return InstanceVaultMeta{
		Url:      normalizeInstanceParamValue(raw["url"]),
		UserPath: normalizeInstanceParamValue(raw["userPath"]),
		Fields:   extractStringSlice(raw["fields"]),
	}
}

func extractStringSlice(value interface{}) []string {
	list, _ := value.([]interface{})
	if list == nil {
		return []string{}
	}

	out := make([]string, 0, len(list))
	for _, item := range list {
		val := strings.TrimSpace(normalizeInstanceParamValue(item))
		if val == "" {
			continue
		}
		out = append(out, val)
	}
	return out
}

func extractVaultSecretValue(body []byte) (string, error) {
	var decoded interface{}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return strings.TrimSpace(string(body)), nil
	}

	switch v := decoded.(type) {
	case string:
		return strings.TrimSpace(v), nil
	case map[string]interface{}:
		if value, ok := extractSecretFromMap(v); ok {
			return value, nil
		}
		b, _ := json.Marshal(v)
		return strings.TrimSpace(string(b)), nil
	default:
		b, _ := json.Marshal(v)
		return strings.TrimSpace(string(b)), nil
	}
}

func extractSecretFromMap(raw map[string]interface{}) (string, bool) {
	for _, key := range []string{"value", "secret", "password", "token"} {
		if value, ok := raw[key]; ok {
			out := strings.TrimSpace(normalizeInstanceParamValue(value))
			if out != "" {
				return out, true
			}
		}
	}

	if data, ok := raw["data"].(map[string]interface{}); ok {
		if value, ok := extractSecretFromMap(data); ok {
			return value, true
		}
	}
	if record, ok := raw["record"].(map[string]interface{}); ok {
		if value, ok := extractSecretFromMap(record); ok {
			return value, true
		}
	}

	if len(raw) == 1 {
		for _, value := range raw {
			out := strings.TrimSpace(normalizeInstanceParamValue(value))
			if out != "" {
				return out, true
			}
		}
	}

	return "", false
}
