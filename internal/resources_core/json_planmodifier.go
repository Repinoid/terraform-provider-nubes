package resources_core

// json_planmodifier.go — план-модификатор для нормализации JSON-атрибутов.
//
// Проблема: API всегда возвращает компактный JSON, но пользователь может написать
// json_parameters = "{ \"log_connections\": \"off\" }" с пробелами.
// Terraform сравнивает план (с пробелами) с результатом Apply (компактный) и видит расхождение →
// "Provider produced inconsistent result after apply".
//
// Решение: plan modifier JsonNormalize() сворачивает JSON до компактного вида ещё на этапе плана,
// чтобы план и результат Read всегда совпадали.

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/repinoid/terraform-provider-nubes/internal/core/jsonutil"
)

// jsonNormalizePlanModifier реализует planmodifier.String.
// Сворачивает JSON-строку до компактного представления при планировании.
type jsonNormalizePlanModifier struct{}

// JsonNormalize возвращает plan modifier, нормализующий JSON к компактному виду.
// Использовать на атрибутах с data_type: json.
func JsonNormalize() planmodifier.String {
	return jsonNormalizePlanModifier{}
}

func (m jsonNormalizePlanModifier) Description(_ context.Context) string {
	return "Normalizes JSON string to compact form to prevent inconsistent-result errors."
}

func (m jsonNormalizePlanModifier) MarkdownDescription(_ context.Context) string {
	return "Normalizes JSON string to compact form to prevent inconsistent-result errors."
}

// PlanModifyString сворачивает JSON до компактного вида.
// Если значение не является корректным JSON — оставляет как есть, не добавляет ошибку.
func (m jsonNormalizePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}
	raw := req.PlanValue.ValueString()
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(raw)); err != nil {
		// Невалидный JSON — оставляем значение без изменений, не прерываем планирование.
		return
	}
	// FIX(uuid-case): UUID-подстроки — к lowercase, чтобы план не расходился с
	// состоянием только из-за регистра UUID (облако отдаёт их в разных регистрах).
	resp.PlanValue = types.StringValue(jsonutil.LowercaseUUIDsInText(buf.String()))
}
