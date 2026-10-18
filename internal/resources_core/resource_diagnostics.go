package resources_core

import (
	"context"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// PlanExistingResourceDiagnostics returns diagnostics for existing resources during планирования.
func PlanExistingResourceDiagnostics(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool) diag.Diagnostics {
	var diags diag.Diagnostics
	if client == nil {
		return diags
	}
	if strings.TrimSpace(resourceName) == "" {
		return diags
	}

	existing, err := client.FindInstanceByDisplayName(ctx, serviceID, resourceName)
	if err != nil || existing == nil {
		return diags
	}

	status := normalizeStatus(existing.ExplainedStatus)
	if isStatusNotCreated(status) {
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (NOT CREATED)",
			"Инстанс найден, но находится в состоянии Not Created.\nУдалите его в Личном кабинете (Управление облаком), затем повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if isStatusRunning(status) {
		if adoptExistingOnCreate {
			return diags
		}
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (RUNNING)",
			"Инстанс с таким resource_name уже работает.\nАвто-усыновление отключено (adopt_existing_on_create=false).\nВыберите другое имя, если нужен новый ресурс.\nИли выполните импорт существующего ресурса через усыновление: установите adopt_existing_on_create=true и повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if isStatusSuspended(status) {
		if adoptExistingOnCreate {
			return diags
		}
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (SUSPEND)",
			"Инстанс найден в состоянии Suspend, но adopt_existing_on_create=false.\nДля восстановления/импорта этим ресурсом включите adopt_existing_on_create=true и повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}

	diags.AddError(
		"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (НЕПОДДЕРЖИВАЕМЫЙ СТАТУС)",
		"Инстанс найден, но статус не подходит для авто-усыновления.\nСтатус: "+status+". Выберите другое имя или устраните проблему в ЛК."+formatInstanceDetails(existing, serviceID, resourceName),
	)
	return diags
}

// CreateExistingResourceDiagnostics returns diagnostics (errors for conflicts, warnings for adopt) during create.
func CreateExistingResourceDiagnostics(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool) diag.Diagnostics {
	return CreateExistingResourceDiagnosticsWithDomain(ctx, client, serviceID, resourceName, adoptExistingOnCreate, "")
}

// ShouldRemoveFromState checks instance status and indicates if TF state should be removed.
//
// Возвращает (true, nil) ТОЛЬКО когда инстанс реально отсутствует на платформе:
//   - состояние ответило deleted (IsDeleted=true или ExplainedStatus=="deleted"), или
//   - API ответил 404 / "not found" (инстанс полностью стёрт).
//
// В остальных случаях (живой инстанс, сетевая ошибка, 5xx, 403 — нет прав)
// возвращает (false, nil) или (false, err). Это принципиально: при временной
// недоступности API нельзя молча вычищать ресурсы из state.
func ShouldRemoveFromState(ctx context.Context, client *core.UniversalClient, instanceID string) (bool, error) {
	if client == nil {
		return false, nil
	}
	if strings.TrimSpace(instanceID) == "" {
		return false, nil
	}

	// Читаем БЕЗ валидации статуса: нужно отличить deleted/404 (убрать из state)
	// от прочих ошибок (пробросить как ошибку).
	state, err := client.GetInstanceStateRaw(ctx, instanceID)
	if err != nil {
		if isInstanceNotFoundError(err) {
			return true, nil
		}
		return false, err
	}
	if state == nil {
		return false, nil
	}

	status := strings.ToLower(strings.TrimSpace(state.ExplainedStatus))
	if state.IsDeleted || status == "deleted" {
		return true, nil
	}

	return false, nil
}

// isInstanceNotFoundError определяет, является ли ошибка чтения инстанса
// признаком того, что инстанс полностью удалён (HTTP 404 / not found),
// а не временным сбоем или отсутствием прав.
func isInstanceNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	// Точный статус 404 — основной признак «инстанс стёрт».
	if strings.Contains(msg, "ошибка api 404") || strings.Contains(msg, "error:404") || strings.Contains(msg, "status 404") {
		return true
	}
	// Дополнительные маркеры «not found» (как в IsSubresourceMissingError).
	if strings.Contains(msg, "not found") && !strings.Contains(msg, "permission") {
		return true
	}
	return false
}

func normalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isStatusRunning(status string) bool {
	if status == "" {
		return false
	}
	return strings.Contains(status, "running") || strings.Contains(status, "active") || strings.Contains(status, "работ")
}

func isStatusNotCreated(status string) bool {
	if status == "" {
		return false
	}
	return strings.Contains(status, "not created") || strings.Contains(status, "not_created") || strings.Contains(status, "notcreated")
}

func formatInstanceDetails(existing *core.InstanceStateResponse, serviceID int, resourceName string) string {
	if existing == nil {
		return ""
	}
	statusRaw := strings.TrimSpace(existing.ExplainedStatus)
	statusNorm := normalizeStatus(statusRaw)
	return "\nDetails:" +
		"\n- resource_name: " + resourceName +
		"\n- service_id: " + fmt.Sprintf("%d", serviceID) +
		"\n- instance_uid: " + existing.InstanceUid +
		"\n- status: \"" + statusNorm + "\"" +
		"\n- status_raw: \"" + statusRaw + "\"" +
		"\n- operation_pending: " + formatBool(existing.OperationIsPending) +
		"\n- operation_in_progress: " + formatBool(existing.OperationIsInProgress)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
