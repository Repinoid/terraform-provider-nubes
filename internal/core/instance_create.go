package core

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreateGenericInstanceUniversalV6 implements the universal flow:
// instances -> instanceOperations -> get cfsParams -> submit params -> validate -> run
func (c *UniversalClient) CreateGenericInstanceUniversalV6(ctx context.Context, serviceId int, displayName string, params map[int]string) (string, error) {
	instPayload := genericInstanceReq{
		ServiceId:   serviceId,
		DisplayName: displayName,
		Descr:       c.instanceDescr(),
	}

	instResp, instHeaders, err := c.doRequest(ctx, "POST", "/instances", instPayload)
	if err != nil {
		return "", err
	}

	instanceUid := extractUIDFromLocation(instHeaders.Get("Location"))
	if instanceUid == "" {
		var instResult struct {
			InstanceUid string `json:"instanceUid"`
		}
		if err := json.Unmarshal(instResp, &instResult); err == nil && instResult.InstanceUid != "" {
			instanceUid = instResult.InstanceUid
		} else {
			var justId string
			if err2 := json.Unmarshal(instResp, &justId); err2 == nil && justId != "" {
				instanceUid = justId
			}
		}
	}
	if instanceUid == "" {
		return "", fmt.Errorf("не удалось извлечь instanceUid из ответа (Header: %s)", instHeaders.Get("Location"))
	}

	opPayload := genericOpReq{
		InstanceUid: instanceUid,
		Operation:   "create",
	}

	opResp, opHeaders, err := c.doRequest(ctx, "POST", "/instanceOperations", opPayload)
	if err != nil {
		// Инстанс УЖЕ создан: возвращаем uid вместе с ошибкой, чтобы вызывающий
		// мог зафиксировать его в state (иначе облачный инстанс осиротеет — Q5).
		return instanceUid, err
	}

	opUid := extractUIDFromLocation(opHeaders.Get("Location"))
	if opUid == "" {
		var opResult struct {
			InstanceOperationUid string `json:"instanceOperationUid"`
		}
		if err := json.Unmarshal(opResp, &opResult); err == nil && opResult.InstanceOperationUid != "" {
			opUid = opResult.InstanceOperationUid
		} else {
			var justId string
			if err2 := json.Unmarshal(opResp, &justId); err2 == nil && justId != "" {
				opUid = justId
			}
		}
	}
	if opUid == "" {
		return instanceUid, fmt.Errorf("не удалось извлечь instanceOperationUid (Header: %s)", opHeaders.Get("Location"))
	}

	opDetailsResp, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/instanceOperations/%s?fields=cfsParams", opUid), nil)
	var opDetails universalOpResponse
	hasOpDetails := false
	if err == nil {
		if jsonErr := json.Unmarshal(opDetailsResp, &opDetails); jsonErr != nil {
			return instanceUid, fmt.Errorf("не удалось разобрать детали операции: %w", jsonErr)
		}
		hasOpDetails = true
		params, err = c.resolveRefSvcParamValues(ctx, opDetails.InstanceOperation.CfsParams, params)
		if err != nil {
			return instanceUid, err
		}
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
			return instanceUid, fmt.Errorf("не удалось установить параметр %s: %w", paramLabelByID(opDetails.InstanceOperation.CfsParams, paramId), err)
		}
		sent[paramId] = true
	}

	if hasOpDetails {
		for _, param := range opDetails.InstanceOperation.CfsParams {
			if sent[param.SvcOperationCfsParamId] {
				continue
			}
			if !param.IsRequired {
				continue
			}

			val := ""
			if param.ParamValue != nil {
				val = *param.ParamValue
			} else if param.DefaultValue != nil {
				val = *param.DefaultValue
			}
			val = normalizeUniversalValueV6(val, param)

			pPayload := genericParamReq{
				InstanceOperationUid:   opUid,
				SvcOperationCfsParamId: param.SvcOperationCfsParamId,
				ParamValue:             val,
			}
			_, _, err := c.doRequest(ctx, "POST", "/instanceOperationCfsParams", pPayload)
			if err != nil {
				return instanceUid, fmt.Errorf("не удалось отправить параметр по умолчанию %s: %w", paramLabel(param), err)
			}
		}
	}

	_, _, err = c.doRequest(ctx, "GET", fmt.Sprintf("/instanceOperations/%s/validate-cfs", opUid), nil)
	if err != nil {
		return instanceUid, fmt.Errorf("валидация не пройдена: %w", err)
	}

	// НЕ МЕНЯТЬ: завершение операции определяется по dtFinish (runOperation
	// запускает POST /run и ждёт dtFinish; при 500 «RESPONSEHEADER» досогласовывает).
	if err := c.runOperation(ctx, opUid, c.operationTimeoutForContext(ctx, serviceId, "create")); err != nil {
		return instanceUid, fmt.Errorf("выполнение не удалось: %w", err)
	}
	if err := c.ensureInstanceCreated(ctx, instanceUid); err != nil {
		return instanceUid, err
	}

	return instanceUid, nil
}
