// Package resources_core — общая логика CRUD для сгенерированных ресурсов.
//
// Предоставляет функции, используемые всеми ресурсами из resources_gen:
//   - Create/Read/Update/Delete — базовые CRUD-операции
//   - State Refresh — обновление состояния из API
//   - Param Mapping — маппинг параметров Terraform ↔ API
//   - Validation — валидация обязательных полей и ref_svc
//   - Plan Modifiers — UUID, JSON-нормализация
//
// ⛔ НЕ ИЗМЕНЯТЬ НИЧЕГО В CRUD БЕЗ ПРЯМОГО РАЗРЕШЕНИЯ ОПЕРАТОРА.
package resources_core

import (
	"context"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// CreateResource uses the universal client flow for create/resume/adopt.
func CreateResource(ctx context.Context, client *core.UniversalClient, serviceID int, displayName string, resumeIfExists bool, params map[int]string) (string, error) {
	return CreateResourceWithTimeout(ctx, client, serviceID, displayName, resumeIfExists, params, "")
}

// CreateResourceWithTimeout uses optional operation timeout override for all create/resume waits.
func CreateResourceWithTimeout(ctx context.Context, client *core.UniversalClient, serviceID int, displayName string, resumeIfExists bool, params map[int]string, operationTimeout string) (string, error) {
	ctxWithTimeout, err := core.WithOperationTimeout(ctx, operationTimeout)
	if err != nil {
		return "", err
	}

	existing, err := client.FindInstanceByDisplayName(ctxWithTimeout, serviceID, displayName)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return adoptExistingInstanceOnCreate(ctxWithTimeout, client, serviceID, displayName, resumeIfExists, params, existing)
	}

	id, err := client.CreateGenericInstanceUniversalV6(ctxWithTimeout, serviceID, displayName, params)
	if err == nil {
		return id, nil
	}
	if strings.TrimSpace(id) == "" {
		return "", err
	}

	// Ошибка операции create, но запись инстанса в облаке могла остаться.
	// Если инстанс НЕ создан (полный откат: isCreated=false) — возвращаем пустой id,
	// иначе он попадёт в state и ресурс окажется в дедлоке: ни apply (запись уже есть),
	// ни destroy (провайдер отказывается удалять "not created"). Такой недосозданный
	// инстанс стирается только вручную в ЛК.
	// См. NOTES/10_plans/PLAN_CREATION_FAILED_GATE_2026-10-02.md
	state, stateErr := client.GetInstanceStateRaw(ctxWithTimeout, id)
	if !shouldKeepIDOnCreateError(state, stateErr) {
		return "", err
	}
	// Инстанс создан, либо состояние определить не удалось (legacy-эндпоинт без
	// isCreated, таймаут IAM и т.п.) — сохраняем текущее поведение: id остаётся,
	// исходная ошибка не маскируется, Read досинхронизирует состояние.
	return id, err
}

// shouldKeepIDOnCreateError решает, оставлять ли id инстанса при ошибке операции create.
//
// true  — сохранять id: инстанс создан (isCreated=true) либо состояние неизвестно
//
//	(stateErr != nil / nil state / IsCreated == nil у legacy-эндпоинта);
//
// false — не сохранять id: инстанс НЕ создан (isCreated=false), писать его в state нельзя.
func shouldKeepIDOnCreateError(state *core.InstanceStateResponse, stateErr error) bool {
	if stateErr != nil || state == nil || state.IsCreated == nil {
		return true
	}
	return *state.IsCreated
}

// UpdateResource uses universal modify flow with defaults.
func UpdateResource(ctx context.Context, client *core.UniversalClient, instanceID string, params map[int]string) error {
	return UpdateResourceWithTimeout(ctx, client, instanceID, params, "")
}

// UpdateResourceWithTimeout uses optional operation timeout override for modify waits.
func UpdateResourceWithTimeout(ctx context.Context, client *core.UniversalClient, instanceID string, params map[int]string, operationTimeout string) error {
	if strings.TrimSpace(instanceID) == "" {
		return fmt.Errorf("missing instance id for modify")
	}
	unlock := client.LockInstance(instanceID)
	defer unlock()

	ctxWithTimeout, err := core.WithOperationTimeout(ctx, operationTimeout)
	if err != nil {
		return err
	}
	return client.RunInstanceOperationUniversalWithDefaults(ctxWithTimeout, instanceID, "modify", params)
}

// DeleteResource runs destroy behavior (suspend/state-only) for instance resources.
func DeleteResource(ctx context.Context, client *core.UniversalClient, instanceID string, destroyBehavior string) error {
	return DeleteResourceWithTimeout(ctx, client, instanceID, destroyBehavior, "")
}

// DeleteResourceWithTimeout uses optional operation timeout override for suspend waits.
func DeleteResourceWithTimeout(ctx context.Context, client *core.UniversalClient, instanceID string, destroyBehavior string, operationTimeout string) error {
	mode := strings.ToLower(strings.TrimSpace(destroyBehavior))
	if mode == "" {
		mode = "state_only"
	}

	if mode != "state_only" && mode != "detach" {
		unlock := client.LockInstance(instanceID)
		defer unlock()
	}

	ctxWithTimeout, err := core.WithOperationTimeout(ctx, operationTimeout)
	if err != nil {
		return err
	}

	switch mode {
	case "state_only", "detach":
		return nil
	case "suspend":
		return client.RunInstanceOperationUniversal(ctxWithTimeout, instanceID, mode, nil)
	case "delete":
		return client.RunInstanceOperationUniversal(ctxWithTimeout, instanceID, mode, nil)
	default:
		return fmt.Errorf("invalid destroy behavior: %s", destroyBehavior)
	}
}

func isStatusSuspended(status string) bool {
	return strings.Contains(status, "suspend") || strings.Contains(status, "suspended")
}

func isStatusNonAdoptable(status string) bool {
	if status == "" {
		return true
	}
	if strings.Contains(status, "not created") {
		return true
	}
	if strings.Contains(status, "pending") || strings.Contains(status, "creating") {
		return true
	}
	if strings.Contains(status, "failed") || strings.Contains(status, "error") {
		return true
	}
	return false
}

// RunOperationByCode triggers a non-CRUD operation using params keyed by code.
func RunOperationByCode(ctx context.Context, client *core.UniversalClient, instanceID string, operation string, params map[string]string) error {
	return RunOperationByCodeWithTimeout(ctx, client, instanceID, operation, params, "")
}

// RunOperationByCodeWithTimeout triggers a non-CRUD operation using params keyed by code
// and optional per-request timeout override.
func RunOperationByCodeWithTimeout(ctx context.Context, client *core.UniversalClient, instanceID string, operation string, params map[string]string, operationTimeout string) error {
	if strings.TrimSpace(instanceID) == "" {
		return fmt.Errorf("missing instance id for operation")
	}
	if client == nil {
		return fmt.Errorf("missing client for operation")
	}
	unlock := client.LockInstance(instanceID)
	defer unlock()

	ctxWithTimeout, err := core.WithOperationTimeout(ctx, operationTimeout)
	if err != nil {
		return err
	}
	return client.RunInstanceOperationUniversalByCode(ctxWithTimeout, instanceID, operation, params)
}

// RunOperationByCodeIdempotent — то же, что RunOperationByCodeWithTimeout, но
// с idempotency pre-check (desired==current → skip run). Для модификаторов
// с idempotency: check_before_run.
func RunOperationByCodeIdempotent(ctx context.Context, client *core.UniversalClient, instanceID string, operation string, params map[string]string, operationTimeout string) error {
	if strings.TrimSpace(instanceID) == "" {
		return fmt.Errorf("missing instance id for operation")
	}
	if client == nil {
		return fmt.Errorf("missing client for operation")
	}
	unlock := client.LockInstance(instanceID)
	defer unlock()

	ctxWithTimeout, err := core.WithOperationTimeout(ctx, operationTimeout)
	if err != nil {
		return err
	}
	return client.RunInstanceOperationUniversalByIdempotent(ctxWithTimeout, instanceID, operation, params)
}

func adoptExistingInstanceOnCreate(ctx context.Context, client *core.UniversalClient, serviceID int, displayName string, resumeIfExists bool, params map[int]string, existing *core.InstanceStateResponse) (string, error) {
	// Проверка: операция уже выполняется
	if existing.OperationIsInProgress || existing.OperationIsPending {
		return "", fmt.Errorf("инстанс с resource_name %s: операция в процессе (operation_in_progress=%v, operation_pending=%v). Дождитесь завершения текущей операции и повторите apply. %s", displayName, existing.OperationIsInProgress, existing.OperationIsPending, formatInstanceDetails(existing, serviceID, displayName))
	}

	if !resumeIfExists {
		return "", fmt.Errorf("инстанс с resource_name %s уже существует; выберите другое имя, если нужен новый ресурс, или выполните импорт через усыновление: установите adopt_existing_on_create=true и повторите apply. %s", displayName, formatInstanceDetails(existing, serviceID, displayName))
	}
	state := classifyStatus(existing)
	statusText := normalizeStatus(existing.ExplainedStatus)

	switch state {
	case StateNotCreated:
		// Философия (provider_philosophy.md раздел 8): not created → hard error,
		// без auto-adopt/create. Пользователь должен проверить ресурс в ЛК.
		return "", fmt.Errorf("инстанс %s в статусе 'not created' — проверьте ресурс в личном кабинете Nubes. %s",
			existing.InstanceUid, formatInstanceDetails(existing, serviceID, displayName))

	case StateRunning, StateRunningPending:
		refIssues, err := ValidateRefParamsOnAdopt(ctx, client, serviceID, params, existing.InstanceUid)
		if err != nil {
			return "", fmt.Errorf("ошибка валидации ref-параметров при adopt: %s", err)
		}
		if HasHardErrors(refIssues) {
			return "", fmt.Errorf("%s\n%s", FormatRefParamIssues(refIssues), formatInstanceDetails(existing, serviceID, displayName))
		}
		return existing.InstanceUid, nil

	case StateSuspended:
		mismatches, err := RequiredParamsMismatch(ctx, client, serviceID, params, existing.InstanceUid)
		if err != nil {
			return "", err
		}
		if len(mismatches) > 0 {
			return "", fmt.Errorf("required params mismatch for resource_name %s: %s. %s", displayName, FormatRequiredParamMismatchMessage(mismatches), formatInstanceDetails(existing, serviceID, displayName))
		}
		unlock := client.LockInstance(existing.InstanceUid)
		defer unlock()
		if err := client.RunInstanceOperationUniversal(ctx, existing.InstanceUid, "resume", nil); err != nil {
			return "", err
		}
		resumed, err := client.GetInstanceState(ctx, existing.InstanceUid)
		if err != nil {
			return "", err
		}
		resumedState := classifyStatus(resumed)
		if !isAdoptable(resumedState) {
			return "", fmt.Errorf("resource not ready after resume: %s. %s", resumed.ExplainedStatus, formatInstanceDetails(resumed, serviceID, displayName))
		}
		refIssues, err := ValidateRefParamsOnAdopt(ctx, client, serviceID, params, existing.InstanceUid)
		if err != nil {
			return "", fmt.Errorf("ошибка валидации ref-параметров после resume: %s", err)
		}
		if HasHardErrors(refIssues) {
			return "", fmt.Errorf("%s\n%s", FormatRefParamIssues(refIssues), formatInstanceDetails(existing, serviceID, displayName))
		}
		return existing.InstanceUid, nil

	case StateCreating, StateModifying, StateSuspending, StateResuming, StateDeleting:
		return "", fmt.Errorf("инстанс с resource_name %s: операция в процессе (статус: %s). Дождитесь завершения и повторите apply. %s",
			displayName, statusText, formatInstanceDetails(existing, serviceID, displayName))

	case StateCreationFailed:
		return "", fmt.Errorf("инстанс с resource_name %s: создание не удалось (статус: %s). Удалите его в Личном кабинете или установите adopt_existing_on_create=true для авто-очистки. %s",
			displayName, statusText, formatInstanceDetails(existing, serviceID, displayName))

	case StateModificationFailed, StateSuspendFailed, StateResumeFailed, StateDeletionFailed:
		return "", fmt.Errorf("инстанс с resource_name %s: предыдущая операция завершилась ошибкой (статус: %s). Проверьте состояние в Личном кабинете. %s",
			displayName, statusText, formatInstanceDetails(existing, serviceID, displayName))

	case StateDeleted, StateOrphaned:
		return "", fmt.Errorf("инстанс с resource_name %s: удалён или осиротел (статус: %s). Выберите другое имя для нового ресурса. %s",
			displayName, statusText, formatInstanceDetails(existing, serviceID, displayName))

	default:
		return "", fmt.Errorf("инстанс с resource_name %s: неизвестный статус %q. %s",
			displayName, statusText, formatInstanceDetails(existing, serviceID, displayName))
	}
}
