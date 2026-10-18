package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ===== Резолв refSvc-параметров (по ID) =====
//
// NOTE(resourceRealm): доступные resourceRealm для сервиса получают через
// GET /resourceRealms/available?svcId=<serviceId> (через index.cfm proxy).
//
// NOTE(refSvcId lookup): если параметр операции create ссылается на другой сервис,
// то в cfsParams будет проставлен refSvcId — значение должно быть UUID инстанса
// указанного сервиса.
//
// ВАЖНО (uuid-case, 2026-03): API облака возвращает UUID всегда в нижнем регистре.
// Для отправки в API нормализуем к lowercase. В state регистр пользователя сохраняется
// потому, что ref_svc-атрибуты ИСКЛЮЧЕНЫ из read-back (шаблон instance.go добавляет
// InputField только при `eq .RefSvcId 0`), а UUID-подстроки ВНУТРИ JSON нормализуются
// при отправке (`resources_core.BuildJSON` -> `jsonutil.LowercaseUUIDsInText`).
// (Блока "Restore user-provided casing" в instance.go нет — он описан только в доке
// docs/60_strategy/terraform_case_sensitivity_fix.md как исторический подход.)

func (c *UniversalClient) resolveRefSvcParamValues(ctx context.Context, opParams []universalCfsParam, params map[int]string) (map[int]string, error) {
	if len(params) == 0 || len(opParams) == 0 {
		return params, nil
	}

	refById := make(map[int]int)
	for _, p := range opParams {
		if p.RefSvcId != nil && *p.RefSvcId > 0 {
			refById[p.SvcOperationCfsParamId] = *p.RefSvcId
		}
	}
	if len(refById) == 0 {
		// Даже без refSvcId — резолвим s3Uid внутри map-fixed (соглашение об именах).
		return c.ResolveS3UidInAllMapFixed(params, opParams), nil
	}

	resolved := make(map[int]string, len(params))
	for paramId, value := range params {
		newValue := value
		if refSvcId, ok := refById[paramId]; ok {
			if strings.TrimSpace(value) == "" {
				newValue = ""
			} else if !isUUIDLike(value) {
				uid, err := c.findInstanceUidByDisplayNameRefSvc(ctx, refSvcId, value)
				if err != nil {
					return nil, err
				}
				if uid == "" {
					return nil, fmt.Errorf("не удалось разрешить параметр %d: не найден экземпляр для serviceId=%d displayName=%s", paramId, refSvcId, value)
				}
				newValue = uid
			} else {
				newValue = strings.ToLower(value)
			}
		} else if strings.TrimSpace(value) != "" {
			// map-fixed: резолвим s3Uid-подобные uuid sub-params внутри JSON.
			newValue = c.resolveS3UidInMapFixed(paramId, value, opParams)
		}
		resolved[paramId] = newValue
	}

	return resolved, nil
}

// ResolveS3UidInAllMapFixed проходит по всем параметрам и для map-fixed
// резолвит s3Uid-подобные значения в UUID. Используется когда нет ни одного
// параметра с refSvcId (ранний выход resolveRefSvcParamValues).
func (c *UniversalClient) ResolveS3UidInAllMapFixed(params map[int]string, opParams []universalCfsParam) map[int]string {
	result := make(map[int]string, len(params))
	for paramId, value := range params {
		if strings.TrimSpace(value) != "" && value != "{}" && strings.HasPrefix(strings.TrimSpace(value), "{") {
			result[paramId] = c.resolveS3UidInMapFixed(paramId, value, opParams)
		} else {
			result[paramId] = value
		}
	}
	return result
}

// resolveS3UidInMapFixed парсит JSON map-fixed-параметра, находит ключи
// по паттерну s3.*uid (case-insensitive) и резолвит значения-не-UUID
// в UUID через S3 (сервис 12). Соглашение об именах: s3Uid, s3_uid, S3Uid → S3.
func (c *UniversalClient) resolveS3UidInMapFixed(paramId int, jsonValue string, opParams []universalCfsParam) string {
	if jsonValue == "" || jsonValue == "{}" {
		return jsonValue
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(jsonValue), &obj); err != nil {
		return jsonValue
	}

	changed := false
	for key, val := range obj {
		lowKey := strings.ToLower(key)
		if !strings.Contains(lowKey, "s3") || !strings.Contains(lowKey, "uid") {
			continue
		}
		strVal, ok := val.(string)
		if !ok || strVal == "" || isUUIDLike(strVal) {
			continue
		}
		uid, err := c.findInstanceUidByDisplayNameRefSvc(context.Background(), 12, strVal)
		if err != nil || uid == "" {
			continue
		}
		obj[key] = uid
		changed = true
	}

	if !changed {
		return jsonValue
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return jsonValue
	}
	return string(b)
}
