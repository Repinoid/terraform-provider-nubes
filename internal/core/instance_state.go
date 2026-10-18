package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ===== Чтение и валидация состояния инстанса =====

func (c *UniversalClient) GetInstanceState(ctx context.Context, instanceUid string) (*InstanceStateResponse, error) {
	state, err := c.getInstanceStateWithRetry(ctx, instanceUid)
	if err != nil {
		return nil, err
	}
	if err := validateInstanceStatus(state); err != nil {
		return nil, err
	}
	return state, nil
}

func isInstanceDeleted(state *InstanceStateResponse) bool {
	if state == nil {
		return false
	}
	if state.IsDeleted {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(state.ExplainedStatus), "deleted")
}

// GetInstanceStateRaw получает состояние инстанса БЕЗ валидации статуса.
// Используется для проверки ref-параметров: нужно читать даже deleted/suspended инстансы.
func (c *UniversalClient) GetInstanceStateRaw(ctx context.Context, instanceUid string) (*InstanceStateResponse, error) {
	return c.getInstanceStateWithRetry(ctx, instanceUid)
}

func (c *UniversalClient) getInstanceStateWithRetry(ctx context.Context, instanceUid string) (*InstanceStateResponse, error) {
	respBody, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/instances/%s", instanceUid), nil)
	if err != nil {
		return nil, err
	}

	var res struct {
		Instance InstanceStateResponse `json:"instance"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}

	return &res.Instance, nil
}

func (c *UniversalClient) ensureInstanceCreated(ctx context.Context, instanceUid string) error {
	state, err := c.GetInstanceState(ctx, instanceUid)
	if err != nil {
		return fmt.Errorf("не удалось получить состояние экземпляра после создания: %w", err)
	}
	if state == nil {
		return fmt.Errorf("отсутствует состояние экземпляра после создания для %s", instanceUid)
	}
	if isInstanceDeleted(state) {
		return fmt.Errorf("экземпляр %s удален после создания", instanceUid)
	}
	status := strings.ToLower(strings.TrimSpace(state.ExplainedStatus))
	if status == "not created" {
		return fmt.Errorf("экземпляр %s не создан после завершения операции. Видимо, такой инстанс уже существует в статусе Not Created. %s", instanceUid, formatInstanceStateDetails(state))
	}
	return nil
}

func validateInstanceStatus(state *InstanceStateResponse) error {
	if state == nil {
		return fmt.Errorf("отсутствует состояние экземпляра")
	}
	if isInstanceDeleted(state) {
		return &instanceDeletedError{InstanceUID: state.InstanceUid}
	}
	if state.OperationIsPending || state.OperationIsInProgress {
		return fmt.Errorf("экземпляр %s не готов: операция в ожидании", state.InstanceUid)
	}
	status := strings.ToLower(strings.TrimSpace(state.ExplainedStatus))
	if strings.Contains(status, "not created") {
		return fmt.Errorf("экземпляр %s не создан. Видимо, такой инстанс уже существует в статусе Not Created. %s", state.InstanceUid, formatInstanceStateDetails(state))
	}
	if strings.Contains(status, "pending") {
		return fmt.Errorf("экземпляр %s в ожидании. %s", state.InstanceUid, formatInstanceStateDetails(state))
	}
	if strings.Contains(status, "failed") || strings.Contains(status, "error") {
		return fmt.Errorf("экземпляр %s завершился с ошибкой: %s. %s", state.InstanceUid, state.ExplainedStatus, formatInstanceStateDetails(state))
	}
	return nil
}

func formatInstanceStateDetails(state *InstanceStateResponse) string {
	if state == nil {
		return "details: state=nil"
	}
	statusRaw := strings.TrimSpace(state.ExplainedStatus)
	statusNorm := strings.ToLower(statusRaw)
	return "details: instance_uid=" + state.InstanceUid + ", status=\"" + statusNorm + "\", status_raw=\"" + statusRaw + "\", operation_pending=" + formatBool(state.OperationIsPending) + ", operation_in_progress=" + formatBool(state.OperationIsInProgress)
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
