package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// instanceLiveParams возвращает ТЕКУЩИЕ params инстанса (GET /instances/{uid} → state.params),
// keyed by lower(code). Это единственный достоверный источник live-значений.
//
// ВАЖНО: paramValue из ответа GET /instanceOperations/{opUid}?fields=cfsParams — это
// дефолт ФОРМЫ операции, а НЕ состояние инстанса. Пример из HAR/edge_.har:
// для needEnableAVI в cfsParams paramValue="false", тогда как live-значение инстанса
// needEnableAVI=true. Досылка paramValue как live стирала ALB (reset-to-default).
//
// При ошибке чтения live возвращает err: тихий возврат пустой map опасен — тогда
// fallback уйдёт на paramValue (дефолт формы) и reset-баг вернётся без признаков
// (см. R3 в ревью). Вызывающий код обязан не молчать.
func (c *UniversalClient) instanceLiveParams(ctx context.Context, instanceUid string) (map[string]string, error) {
	live := map[string]string{}
	params, err := c.GetInstanceStateParams(ctx, instanceUid)
	if err != nil {
		return live, err
	}
	for k, v := range params {
		key := strings.ToLower(strings.TrimSpace(k))
		if key == "" {
			continue
		}
		live[key] = v
	}
	return live, nil
}

// lookupLiveParam ищет live-значение параметра инстанса по возможным ключам cfsParam.
func lookupLiveParam(live map[string]string, p universalCfsParam) (string, bool) {
	if len(live) == 0 {
		return "", false
	}
	for _, cand := range []string{p.Code, p.SvcOperationCfsParam, p.Name, p.Label} {
		key := strings.ToLower(strings.TrimSpace(cand))
		if key == "" {
			continue
		}
		if v, ok := live[key]; ok && strings.TrimSpace(v) != "" {
			return v, true
		}
	}
	return "", false
}

// fetchOperationCfsParams возвращает cfsParams операции. При неудаче живого
// GET /instanceOperations/{opUid}?fields=cfsParams (backend 500 в
// getResourceRealmConfig на проблемных инстансах) берёт схему из
// /instanceOperations/default/{opId}.
func (c *UniversalClient) fetchOperationCfsParams(ctx context.Context, opUid string, opId int) ([]universalCfsParam, error) {
	opDetailsResp, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/instanceOperations/%s?fields=cfsParams", opUid), nil)
	if err == nil {
		var opDetails universalOpResponse
		if uerr := json.Unmarshal(opDetailsResp, &opDetails); uerr == nil {
			return opDetails.InstanceOperation.CfsParams, nil
		}
	}

	defResp, _, defErr := c.doRequest(ctx, "GET", fmt.Sprintf("/instanceOperations/default/%d", opId), nil)
	if defErr != nil {
		return nil, fmt.Errorf("не удалось получить детали операции: %w", defErr)
	}
	var def universalOpDefaultResponse
	if uerr := json.Unmarshal(defResp, &def); uerr != nil {
		return nil, fmt.Errorf("не удалось разобрать детали операции: %w", uerr)
	}
	return def.SvcOperation.CfsParams, nil
}
