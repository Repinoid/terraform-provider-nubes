package core

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ===== Запуск операций инстанса (params по ID) =====
//
//   - RunInstanceOperationUniversal — операция с params по ID (без досылки дефолтов)
//   - RunInstanceOperationUniversalWithDefaults — то же + досылка required-дефолтов
//   - RunRedeployOperation — redeploy
//
// Операции с params по code — см. operation_run_bycode.go.
// Схема cfsParams (с fallback) — см. operation_cfs.go.

// RunInstanceOperationUniversal runs an available operation (modify/suspend/delete/resume) if possible.
func (c *UniversalClient) RunInstanceOperationUniversal(ctx context.Context, instanceUid string, action string, params map[int]string) error {
	state, err := c.GetInstanceState(ctx, instanceUid)
	if err != nil {
		return err
	}
	if state.OperationIsPending || state.OperationIsInProgress {
		if err := c.waitForInstanceIdle(ctx, instanceUid, c.idleTimeoutFor(state.ServiceId)); err != nil {
			return err
		}
	}

	var opId int
	for _, op := range state.AvailableOperations {
		if strings.EqualFold(op.Operation, action) {
			opId = op.SvcOperationId
			break
		}
	}
	if opId == 0 {
		return fmt.Errorf("операция %s недоступна для экземпляра %s", action, instanceUid)
	}

	payload := map[string]interface{}{
		"instanceUid":    instanceUid,
		"svcOperationId": opId,
		"operation":      action,
	}

	opUid, err := c.postIgnoreResponse(ctx, "/instanceOperations", payload, true)
	if err != nil {
		return fmt.Errorf("не удалось создать операцию %s: %w", action, err)
	}
	if opUid == "" {
		return fmt.Errorf("не удалось получить UID операции для %s", action)
	}

	for paramId, value := range params {
		pPayload := genericParamReq{
			InstanceOperationUid:   opUid,
			SvcOperationCfsParamId: paramId,
			ParamValue:             value,
		}
		_, _, err := c.doRequest(ctx, "POST", "/instanceOperationCfsParams", pPayload)
		if err != nil {
			return fmt.Errorf("не удалось установить параметр %d: %w", paramId, err)
		}
	}

	// НЕ МЕНЯТЬ: завершение операции определяется по dtFinish (runOperation
	// запускает POST /run и ждёт dtFinish; при 500 «RESPONSEHEADER» досогласовывает).
	return c.runOperation(ctx, opUid, c.operationTimeoutForContext(ctx, state.ServiceId, action))
}

// RunInstanceOperationUniversalWithDefaults runs an operation and submits required params (including defaults).
func (c *UniversalClient) RunInstanceOperationUniversalWithDefaults(ctx context.Context, instanceUid string, action string, params map[int]string) error {
	state, err := c.GetInstanceState(ctx, instanceUid)
	if err != nil {
		return err
	}
	if state.OperationIsPending || state.OperationIsInProgress {
		if err := c.waitForInstanceIdle(ctx, instanceUid, c.idleTimeoutFor(state.ServiceId)); err != nil {
			return err
		}
	}

	var opId int
	for _, op := range state.AvailableOperations {
		if strings.EqualFold(op.Operation, action) {
			opId = op.SvcOperationId
			break
		}
	}
	if opId == 0 {
		return fmt.Errorf("операция %s недоступна для экземпляра %s", action, instanceUid)
	}

	payload := map[string]interface{}{
		"instanceUid":    instanceUid,
		"svcOperationId": opId,
		"operation":      action,
	}

	opUid, err := c.postIgnoreResponse(ctx, "/instanceOperations", payload, true)
	if err != nil {
		return fmt.Errorf("не удалось создать операцию %s: %w", action, err)
	}
	if opUid == "" {
		return fmt.Errorf("не удалось получить UID операции для %s", action)
	}

	cfsParams, err := c.fetchOperationCfsParams(ctx, opUid, opId)
	if err != nil {
		return err
	}

	params, err = c.resolveRefSvcParamValues(ctx, cfsParams, params)
	if err != nil {
		return err
	}

	sent := make(map[int]bool)
	for paramId, value := range params {
		pPayload := genericParamReq{
			InstanceOperationUid:   opUid,
			SvcOperationCfsParamId: paramId,
			ParamValue:             value,
		}
		_, _, err := c.doRequest(ctx, "POST", "/instanceOperationCfsParams", pPayload)
		if err != nil {
			return fmt.Errorf("не удалось установить параметр %s: %w", paramLabelByID(cfsParams, paramId), err)
		}
		sent[paramId] = true
	}

	live, liveErr := c.instanceLiveParams(ctx, instanceUid)
	if liveErr != nil {
		return fmt.Errorf("не удалось прочитать live-значения инстанса для досылки modify: %w", liveErr)
	}
	for _, param := range cfsParams {
		if sent[param.SvcOperationCfsParamId] {
			continue
		}

		// Приоритет: live state.params инстанса → paramValue операции → defaultValue.
		// paramValue из cfsParams — дефолт формы, не состояние инстанса (см. operation_cfs.go).
		// Симметрично runInstanceOperationByCode: если ни одного источника нет — пропускаем
		// (иначе уйдёт синтетический "0"/"false"/"[]" и нарушит constraint).
		val, hasLive := lookupLiveParam(live, param)
		if !hasLive {
			if (param.ParamValue == nil || strings.TrimSpace(*param.ParamValue) == "") &&
				(param.DefaultValue == nil || strings.TrimSpace(*param.DefaultValue) == "") {
				continue
			}
			if param.ParamValue != nil && strings.TrimSpace(*param.ParamValue) != "" {
				val = *param.ParamValue
			} else if param.DefaultValue != nil {
				val = *param.DefaultValue
			}
		}
		val = normalizeUniversalValueV6(val, param)

		pPayload := genericParamReq{
			InstanceOperationUid:   opUid,
			SvcOperationCfsParamId: param.SvcOperationCfsParamId,
			ParamValue:             val,
		}
		_, _, err := c.doRequest(ctx, "POST", "/instanceOperationCfsParams", pPayload)
		if err != nil {
			return fmt.Errorf("не удалось отправить параметр по умолчанию %s: %w", paramLabel(param), err)
		}
	}

	_, _, err = c.doRequest(ctx, "GET", fmt.Sprintf("/instanceOperations/%s/validate-cfs", opUid), nil)
	if err != nil {
		return fmt.Errorf("валидация не пройдена: %w", err)
	}

	// НЕ МЕНЯТЬ: завершение операции определяется по dtFinish (runOperation
	// запускает POST /run и ждёт dtFinish; при 500 «RESPONSEHEADER» досогласовывает).
	return c.runOperation(ctx, opUid, c.operationTimeoutForContext(ctx, state.ServiceId, action))
}

// RunRedeployOperation запускает redeploy для сервисов, поддерживающих пересборку из git.
// Если params не nil — отправляет CFS-параметры перед запуском.
func (c *UniversalClient) RunRedeployOperation(ctx context.Context, instanceUid string, timeoutOverride string, params map[int]string) error {
	state, err := c.GetInstanceState(ctx, instanceUid)
	if err != nil {
		return err
	}
	if state.OperationIsPending || state.OperationIsInProgress {
		if err := c.waitForInstanceIdle(ctx, instanceUid, c.idleTimeoutFor(state.ServiceId)); err != nil {
			return err
		}
	}

	opId := 0
	for _, op := range state.AvailableOperations {
		if strings.EqualFold(op.Operation, "redeploy") {
			opId = op.SvcOperationId
			break
		}
	}
	if opId == 0 {
		return fmt.Errorf("операция redeploy недоступна для экземпляра %s", instanceUid)
	}

	payload := map[string]interface{}{
		"instanceUid":    instanceUid,
		"svcOperationId": opId,
		"operation":      "redeploy",
	}
	opUid, err := c.postIgnoreResponse(ctx, "/instanceOperations", payload, true)
	if err != nil {
		return fmt.Errorf("не удалось создать операцию redeploy: %w", err)
	}
	if opUid == "" {
		return fmt.Errorf("не удалось получить UID операции redeploy")
	}

	for paramId, value := range params {
		pPayload := genericParamReq{
			InstanceOperationUid:   opUid,
			SvcOperationCfsParamId: paramId,
			ParamValue:             value,
		}
		if _, _, err := c.doRequest(ctx, "POST", "/instanceOperationCfsParams", pPayload); err != nil {
			return fmt.Errorf("не удалось установить параметр %d для redeploy: %w", paramId, err)
		}
	}

	timeout := c.operationTimeoutForContext(ctx, state.ServiceId, "redeploy")
	if timeoutOverride != "" {
		if d, parseErr := time.ParseDuration(timeoutOverride); parseErr == nil {
			timeout = d
		}
	}
	return c.runOperation(ctx, opUid, timeout)
}
