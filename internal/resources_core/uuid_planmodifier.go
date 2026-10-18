package resources_core

// uuid_planmodifier.go — план-модификатор для нормализации UUID-атрибутов к нижнему регистру.
//
// Проблема: пользователь может написать UUID в верхнем регистре ("6214BA32-..."),
// но API всегда возвращает UUID в нижнем регистре ("6214ba32-...").
// Terraform сравнивает план (uppercase) с результатом Apply (lowercase) и видит расхождение →
// "Provider produced inconsistent result after apply".
//
// Решение: plan modifier UUIDNormalize() нормализует UUID к нижнему регистру ещё на этапе плана,
// чтобы план и результат Read всегда совпадали.

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// uuidNormalizePlanModifier реализует planmodifier.String.
// Нормализует UUID-строку к нижнему регистру при планировании.
type uuidNormalizePlanModifier struct{}

// UUIDNormalize возвращает plan modifier, нормализующий значение к нижнему регистру.
// Использовать на атрибутах с ref_svc (типа UUID).
func UUIDNormalize() planmodifier.String {
	return uuidNormalizePlanModifier{}
}

func (m uuidNormalizePlanModifier) Description(_ context.Context) string {
	return "Normalizes UUID string to lowercase to prevent inconsistent-result errors."
}

func (m uuidNormalizePlanModifier) MarkdownDescription(_ context.Context) string {
	return "Normalizes UUID string to lowercase to prevent inconsistent-result errors."
}

// PlanModifyString нормализует значение к нижнему регистру.
func (m uuidNormalizePlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}
	normalized := strings.ToLower(req.PlanValue.ValueString())
	if normalized != req.PlanValue.ValueString() {
		resp.PlanValue = types.StringValue(normalized)
	}
}
