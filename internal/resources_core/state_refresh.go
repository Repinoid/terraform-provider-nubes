package resources_core

import (
	"context"
	"reflect"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StateField struct {
	Code string
}

type InputField struct {
	Code  string
	Field string
	Type  string
}

var (
	typeString = reflect.TypeOf(types.String{})
	typeBool   = reflect.TypeOf(types.Bool{})
	typeInt64  = reflect.TypeOf(types.Int64{})
	typeMap    = reflect.TypeOf(types.Map{})
	typeList   = reflect.TypeOf(types.List{})
)

func RefreshResourceState[T any](ctx context.Context, client *core.UniversalClient, instanceID string, serviceID int, state T, outputs []StateField, inputs []InputField) (T, diag.Diagnostics) {
	var diags diag.Diagnostics
	if client == nil || strings.TrimSpace(instanceID) == "" {
		return state, diags
	}

	out, outDiags := FetchInstanceOutputs(ctx, client, instanceID)
	diags.Append(outDiags...)

	v := reflect.ValueOf(&state).Elem()
	if v.Kind() != reflect.Struct {
		return state, diags
	}

	stateOutMap := map[string]string{}
	if !out.StateOut.IsNull() && !out.StateOut.IsUnknown() {
		mapped, mapDiags := StringMapFromTypesMap(ctx, out.StateOut)
		diags.Append(mapDiags...)
		if len(mapped) > 0 {
			stateOutMap = mapped
		}
	}

	for _, field := range outputs {
		code := strings.TrimSpace(field.Code)
		if code == "" {
			continue
		}
		fieldName := toCamel(code)
		fv := v.FieldByName(fieldName)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}
		switch code {
		case "state_params":
			setFieldValue(fv, out.StateParams)
		case "state_out":
			setFieldValue(fv, out.StateOut)
		case "state_params_flat":
			setFieldValue(fv, out.StateParamsFlat)
		case "state_out_flat":
			setFieldValue(fv, out.StateOutFlat)
		case "vault_secrets":
			setFieldValue(fv, out.VaultSecrets)
		case "vault_url":
			setFieldValue(fv, out.VaultUrl)
		case "vault_user_path":
			setFieldValue(fv, out.VaultUserPath)
		case "vault_fields":
			setFieldValue(fv, out.VaultFields)
		default:
			if val, ok := stateOutMap[code]; ok {
				setFieldValue(fv, ParseString(val))
			}
		}
	}

	paramsMap, paramsDiags := StringMapFromTypesMap(ctx, out.StateParams)
	diags.Append(paramsDiags...)
	if serviceID > 0 && len(paramsMap) > 0 {
		paramsMap, paramsDiags = ResolveRefSvcParamDisplayNames(ctx, client, serviceID, paramsMap)
		diags.Append(paramsDiags...)
	}

	resourceRealm := strings.TrimSpace(paramsMap["resourceRealm"])
	if resourceRealm != "" {
		outFlatMap, flatDiags := StringMapFromTypesMap(ctx, out.StateOutFlat)
		diags.Append(flatDiags...)
		if len(outFlatMap) > 0 {
			fixed := FixInternalConnectMasterSuffix(outFlatMap, resourceRealm)
			if fixed {
				out.StateOutFlat, diags = mapValueFrom(ctx, outFlatMap, diags)
				fv := v.FieldByName(toCamel("state_out_flat"))
				if fv.IsValid() && fv.CanSet() {
					setFieldValue(fv, out.StateOutFlat)
				}
			}
		}
	}

	for _, input := range inputs {
		code := strings.TrimSpace(input.Code)
		if code == "" {
			continue
		}
		// Поле модели достаём ДО проверки наличия значения в API: оно нужно
		// и в ветке "API не вернул код" (см. схлопывание unknown → null ниже).
		fieldName := strings.TrimSpace(input.Field)
		if fieldName == "" {
			fieldName = toCamel(code)
		}
		fv := v.FieldByName(fieldName)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		value, ok := paramsMap[code]
		if !ok {
			// ИНВАРИАНТ: Computed-атрибут обязан быть KNOWN после apply/read.
			//
			// Параметры, которые провайдер читает обратно, объявлены в схеме как
			// Optional+Computed (см. helpers.ShouldBeOptionalComputed). Если
			// пользователь такой параметр не задал, в плане он = unknown, и именно
			// провайдер обязан проставить конкретное значение. Когда платформа
			// не вернула код в state_params, единственное корректное конкретное
			// значение — null.
			//
			// Если оставить unknown, Terraform упадёт с
			// "Provider produced invalid result object after apply: ... was unknown".
			if inputFieldUnknown(fv) {
				setInputFieldNull(fv)
			}
			continue
		}

		if strings.EqualFold(code, "jsonEnv") && fv.Type() == typeString {
			// Preserve planned json_env when API returns equivalent JSON with different ordering.
			if planned, ok := fv.Interface().(types.String); ok && !planned.IsNull() && !planned.IsUnknown() {
				if JSONStringsEquivalent(planned.ValueString(), value) {
					fv.Set(reflect.ValueOf(planned))
					continue
				}
			}
		}

		switch strings.ToLower(strings.TrimSpace(input.Type)) {
		case "bool":
			if fv.Type() == typeBool {
				fv.Set(reflect.ValueOf(ParseBool(value)))
			}
		case "int", "int64", "number":
			if fv.Type() == typeInt64 {
				fv.Set(reflect.ValueOf(ParseInt64(value)))
			}
		default:
			if fv.Type() == typeString {
				fv.Set(reflect.ValueOf(ParseString(value)))
			}
		}
	}

	return state, diags
}

// inputFieldUnknown сообщает, находится ли поле модели в состоянии unknown.
//
// Зачем: соблюдение инварианта «Computed-атрибут обязан быть known после
// apply/read». Если платформа не вернула значение в state_params, unknown
// оставлять нельзя — его надо схлопнуть в null (см. setInputFieldNull).
func inputFieldUnknown(fv reflect.Value) bool {
	switch fv.Type() {
	case typeString:
		v, ok := fv.Interface().(types.String)
		return ok && v.IsUnknown()
	case typeBool:
		v, ok := fv.Interface().(types.Bool)
		return ok && v.IsUnknown()
	case typeInt64:
		v, ok := fv.Interface().(types.Int64)
		return ok && v.IsUnknown()
	}
	return false
}

// setInputFieldNull записывает в поле модели типизированный null.
//
// Зачем: null — это конкретное (known) значение, в отличие от unknown. Именно
// null приводит состояние Terraform в консистентный вид, когда платформа не
// сообщила значение для read-back параметра.
func setInputFieldNull(fv reflect.Value) {
	switch fv.Type() {
	case typeString:
		fv.Set(reflect.ValueOf(types.StringNull()))
	case typeBool:
		fv.Set(reflect.ValueOf(types.BoolNull()))
	case typeInt64:
		fv.Set(reflect.ValueOf(types.Int64Null()))
	}
}

func setFieldValue(field reflect.Value, value interface{}) {
	switch field.Type() {
	case typeString:
		switch v := value.(type) {
		case types.String:
			field.Set(reflect.ValueOf(v))
		case string:
			field.Set(reflect.ValueOf(ParseString(v)))
		}
	case typeBool:
		if v, ok := value.(types.Bool); ok {
			field.Set(reflect.ValueOf(v))
		}
	case typeInt64:
		if v, ok := value.(types.Int64); ok {
			field.Set(reflect.ValueOf(v))
		}
	case typeMap:
		if v, ok := value.(types.Map); ok {
			field.Set(reflect.ValueOf(v))
		}
	case typeList:
		if v, ok := value.(types.List); ok {
			field.Set(reflect.ValueOf(v))
		}
	}
}

func toCamel(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' })
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	out := strings.Join(parts, "")
	if out == "" {
		return "R"
	}
	first := rune(out[0])
	if (first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z') || first == '_' {
		return out
	}
	return "R" + out
}
