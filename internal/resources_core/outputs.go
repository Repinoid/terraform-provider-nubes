package resources_core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InstanceOutputs struct {
	StateParams     types.Map
	StateOut        types.Map
	StateParamsFlat types.Map
	StateOutFlat    types.Map
	VaultSecrets    types.Map
	VaultUrl        types.String
	VaultUserPath   types.String
	VaultFields     types.List
}

// FetchInstanceOutputs loads instance state params/out and vault secrets into TF types.
func FetchInstanceOutputs(ctx context.Context, client *core.UniversalClient, instanceID string) (InstanceOutputs, diag.Diagnostics) {
	var diags diag.Diagnostics
	outputs := emptyInstanceOutputs()
	if client == nil || strings.TrimSpace(instanceID) == "" {
		return outputs, diags
	}

	details, err := client.GetInstanceStateDetails(ctx, instanceID)
	if err != nil {
		return outputs, diags
	}
	if details == nil {
		return outputs, diags
	}

	params := normalizeStringMap(details.Params)
	out := normalizeStringMap(details.Out)
	paramsFlat := flattenStringMap(details.RawParams)
	outFlat := flattenStringMap(details.RawOut)
	fields := normalizeStringSlice(details.Vault.Fields)

	outputs.StateParams, diags = mapValueFrom(ctx, params, diags)
	outputs.StateOut, diags = mapValueFrom(ctx, out, diags)
	outputs.StateParamsFlat, diags = mapValueFrom(ctx, paramsFlat, diags)
	outputs.StateOutFlat, diags = mapValueFrom(ctx, outFlat, diags)

	secrets, err := client.GetInstanceVaultSecrets(ctx, instanceID, fields)
	if err != nil {
		secrets = nil
	}
	outputs.VaultSecrets, diags = mapValueFrom(ctx, normalizeStringMap(secrets), diags)
	outputs.VaultFields, diags = listValueFrom(ctx, fields, diags)
	outputs.VaultUrl = stringValue(details.Vault.Url)
	outputs.VaultUserPath = stringValue(details.Vault.UserPath)

	return outputs, diags
}

func emptyInstanceOutputs() InstanceOutputs {
	return InstanceOutputs{
		StateParams:     types.MapNull(types.StringType),
		StateOut:        types.MapNull(types.StringType),
		StateParamsFlat: types.MapNull(types.StringType),
		StateOutFlat:    types.MapNull(types.StringType),
		VaultSecrets:    types.MapNull(types.StringType),
		VaultUrl:        types.StringNull(),
		VaultUserPath:   types.StringNull(),
		VaultFields:     types.ListNull(types.StringType),
	}
}

func flattenStringMap(value map[string]interface{}) map[string]string {
	if value == nil {
		return map[string]string{}
	}
	flat := make(map[string]string)
	for key, val := range value {
		flattenValue(flat, key, val)
	}
	return flat
}

func flattenValue(out map[string]string, prefix string, value interface{}) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, child := range v {
			childKey := prefix + "." + key
			flattenValue(out, childKey, child)
		}
	case []interface{}:
		for idx, child := range v {
			childKey := fmt.Sprintf("%s[%d]", prefix, idx)
			flattenValue(out, childKey, child)
		}
	default:
		out[prefix] = normalizeFlatValue(value)
	}
}

func normalizeFlatValue(value interface{}) string {
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
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
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

func normalizeStringMap(value map[string]string) map[string]string {
	if value == nil {
		return map[string]string{}
	}
	return value
}

func normalizeStringSlice(value []string) []string {
	if value == nil {
		return []string{}
	}
	return value
}

func stringValue(value string) types.String {
	if strings.TrimSpace(value) == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func mapValueFrom(ctx context.Context, value map[string]string, diags diag.Diagnostics) (types.Map, diag.Diagnostics) {
	result, d := types.MapValueFrom(ctx, types.StringType, value)
	diags.Append(d...)
	return result, diags
}

func listValueFrom(ctx context.Context, value []string, diags diag.Diagnostics) (types.List, diag.Diagnostics) {
	result, d := types.ListValueFrom(ctx, types.StringType, value)
	diags.Append(d...)
	return result, diags
}

// StringMapFromTypesMap extracts a map[string]string from a TF map for state_params sync.
func StringMapFromTypesMap(ctx context.Context, value types.Map) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value.IsNull() || value.IsUnknown() {
		return map[string]string{}, diags
	}

	var result map[string]string
	if d := value.ElementsAs(ctx, &result, false); d.HasError() {
		diags.Append(d...)
		return map[string]string{}, diags
	}
	if result == nil {
		return map[string]string{}, diags
	}
	return result, diags
}

func FixInternalConnectMasterSuffix(outFlat map[string]string, resourceRealm string) bool {
	if len(outFlat) == 0 || strings.TrimSpace(resourceRealm) == "" {
		return false
	}
	value := strings.TrimSpace(outFlat["internalConnect.master"])
	if value == "" {
		return false
	}
	if !strings.HasSuffix(value, ".svc.") {
		return false
	}
	fixed := value + strings.TrimPrefix(resourceRealm, ".")
	if fixed == value {
		return false
	}
	outFlat["internalConnect.master"] = fixed
	return true
}
