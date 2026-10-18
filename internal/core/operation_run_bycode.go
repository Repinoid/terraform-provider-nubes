package core

import (
	"context"
	"fmt"
	"strings"
)

// RunInstanceOperationUniversalByCode runs an operation using params keyed by code.
// It resolves param codes to IDs via operation manifest, applies defaults, validates, and runs.
func (c *UniversalClient) RunInstanceOperationUniversalByCode(ctx context.Context, instanceUid string, action string, params map[string]string) error {
	return c.runInstanceOperationByCode(ctx, instanceUid, action, params, false)
}

// RunInstanceOperationUniversalByIdempotent — то же, но с idempotency pre-check:
// перед run сверяет desired==current и, если ВСЕ поля совпали, пропускает run.
// Применяется для модификаторов с idempotency: check_before_run.
func (c *UniversalClient) RunInstanceOperationUniversalByIdempotent(ctx context.Context, instanceUid string, action string, params map[string]string) error {
	return c.runInstanceOperationByCode(ctx, instanceUid, action, params, true)
}

func (c *UniversalClient) runInstanceOperationByCode(ctx context.Context, instanceUid string, action string, params map[string]string, idempotent bool) error {
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

	// Idempotency pre-check: если все desired уже равны LIVE-значениям инстанса —
	// выходим, НЕ создавая операцию вовсе. Схема операции здесь принципиально НЕ
	// запрашивается: pre-check идёт по live-кодам (см. modifierDesiredEqualsLive и
	// ARCHITECTURE.md «Modifier idempotency»).
	if idempotent {
		equal, liveErr := c.modifierDesiredEqualsLive(ctx, instanceUid, params)
		if liveErr == nil && equal {
			return nil
		}
		// Ошибка/несовпадение live — не пропускаем: выполняем modify.
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

	codeToParam := modifierCodeMap(cfsParams)

	paramsByID := map[int]string{}
	for code, value := range params {
		key := strings.ToLower(strings.TrimSpace(code))
		p, ok := codeToParam[key]
		if !ok {
			return fmt.Errorf("код параметра %s не найден для операции %s", code, action)
		}
		paramsByID[p.SvcOperationCfsParamId] = value
	}

	paramsByID, err = c.resolveRefSvcParamValues(ctx, cfsParams, paramsByID)
	if err != nil {
		return err
	}

	sent := make(map[int]bool)
	for paramId, value := range paramsByID {
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
		// Дозаполняем ВСЕ незаданные параметры, чтобы бэкенд modify не трактовал
		// пропущенный/null как reset-to-default (иначе частичный payload затирает
		// create-поля, см. prompt_for_opus_modifier_null_bug.md).
		//
		// ПРИОРИТЕТ ИСТОЧНИКА: live state.params инстанса → paramValue операции → defaultValue.
		// КРИТИЧНО: paramValue из ?fields=cfsParams — дефолт ФОРМЫ операции, не состояние
		// инстанса (HAR/edge_.har: needEnableAVI paramValue="false" при live=true → стирало ALB).
		// Если ни live, ни paramValue, ни defaultValue НЕТ — пропускаем (не шлём синтетический
		// "0"/"false"/"[]", который может нарушить constraint "integer > 0").
		val, hasLive := lookupLiveParam(live, param)
		if !hasLive {
			if (param.ParamValue == nil || strings.TrimSpace(*param.ParamValue) == "") &&
				(param.DefaultValue == nil || strings.TrimSpace(*param.DefaultValue) == "") {
				continue
			}
			if param.ParamValue != nil {
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
