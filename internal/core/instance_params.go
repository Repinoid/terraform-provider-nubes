package core

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// GetInstanceStateParams returns instance.state.params as a map of string values.
func (c *UniversalClient) GetInstanceStateParams(ctx context.Context, instanceUid string) (map[string]string, error) {
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
		return nil, fmt.Errorf("instance field missing in response")
	}

	state, _ := inst["state"].(map[string]interface{})
	if state == nil {
		return map[string]string{}, nil
	}

	paramsRaw, _ := state["params"].(map[string]interface{})
	if paramsRaw == nil {
		return map[string]string{}, nil
	}

	params := make(map[string]string, len(paramsRaw))
	for key, val := range paramsRaw {
		params[key] = normalizeInstanceParamValue(val)
	}

	return params, nil
}

func normalizeInstanceParamValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		if math.Mod(v, 1) == 0 {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case json.Number:
		return v.String()
	default:
		b, err := json.Marshal(v)
		if err == nil {
			return strings.TrimSpace(string(b))
		}
		return fmt.Sprintf("%v", v)
	}
}
