package resources_core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/repinoid/terraform-provider-nubes/internal/core/jsonutil"
)

func FormatBool(v types.Bool) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	if v.ValueBool() {
		return "true"
	}
	return "false"
}

func FormatInt64(v types.Int64) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return fmt.Sprintf("%d", v.ValueInt64())
}

func FormatString(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return v.ValueString()
}

// ParseBool converts string params from state_params into typed values.
func ParseBool(value string) types.Bool {
	if strings.TrimSpace(value) == "" {
		return types.BoolNull()
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		return types.BoolNull()
	}
	return types.BoolValue(parsed)
}

// ParseInt64 converts string params from state_params into typed values.
func ParseInt64(value string) types.Int64 {
	if strings.TrimSpace(value) == "" {
		return types.Int64Null()
	}
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return types.Int64Null()
	}
	return types.Int64Value(parsed)
}

// ParseString converts string params from state_params into typed values.
func ParseString(value string) types.String {
	if strings.TrimSpace(value) == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// CompactParams drops empty values to avoid sending empty strings for optional params.
func CompactParams(params map[string]string) map[string]string {
	if len(params) == 0 {
		return params
	}
	out := make(map[string]string, len(params))
	for k, v := range params {
		if strings.TrimSpace(v) == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// BuildJSON строит JSON-строку из map[string]string значений.
// Ключи не кавычатся если они уже в кавычках.
//
// FIX(uuid-case): после сборки приводим UUID-подстроки к lowercase. Облако хранит
// UUID в lowercase, но ресурсы могут отдавать id в UPPERCASE (например nsxtUid у
// кластера Штурвала), а платформа сравнивает регистр при create и падает с
// «Edge не развёрнут в указанном vDC». Без этой нормализации приходилось писать
// lower(...) в HCL у каждого такого поля.
func BuildJSON(parts map[string]string) string {
	if len(parts) == 0 {
		return "{}"
	}
	pairs := make([]string, 0, len(parts))
	for k, v := range parts {
		pairs = append(pairs, fmt.Sprintf(`"%s":%s`, k, v))
	}
	return jsonutil.LowercaseUUIDsInText("{" + strings.Join(pairs, ",") + "}")
}
