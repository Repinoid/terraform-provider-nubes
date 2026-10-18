package resources_core

import (
	"context"
	"regexp"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// uuidLikeCompareRegex matches UUID-like strings (any case) so they can be
// normalised to lowercase before comparison (see normalizeCompareValue).
//
// Зачем: API облака принимает UUID в любом регистре, но сам всегда возвращает
// lowercase. Пользователь может написать "6214BA32-...", API вернёт "6214ba32-...".
// Чтобы ParamsMatchForResume (сравнение desired vs actual) работал корректно при
// suspend/resume, нормализуем UUID к lowercase с обеих сторон перед сравнением.
var uuidLikeCompareRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ParamsMatchForResume compares desired create params to instance.state.params.
// It returns true only when all non-empty desired params match actual values.
func ParamsMatchForResume(ctx context.Context, client *core.UniversalClient, serviceID int, desiredParamsByID map[int]string, instanceUID string) (bool, error) {
	if client == nil {
		return false, nil
	}
	codeMap, err := getServiceParamCodeMap(serviceID)
	if err != nil {
		return false, err
	}
	if len(codeMap) == 0 {
		return false, nil
	}
	refMap, err := getServiceParamRefMap(serviceID)
	if err != nil {
		return false, err
	}

	desiredByCode := make(map[string]string, len(desiredParamsByID))
	for id, val := range desiredParamsByID {
		code, ok := codeMap[id]
		if !ok {
			return false, nil
		}
		if refSvcId, ok := refMap[id]; ok && refSvcId > 0 {
			resolved, err := client.ResolveRefSvcParamValue(ctx, refSvcId, val)
			if err != nil {
				return false, err
			}
			if strings.TrimSpace(resolved) != "" {
				val = resolved
			}
		}
		norm := normalizeCompareValue(val)
		if norm == "" {
			continue
		}
		desiredByCode[code] = norm
	}

	if len(desiredByCode) == 0 {
		return true, nil
	}

	actualParams, err := client.GetInstanceStateParams(ctx, instanceUID)
	if err != nil {
		return false, err
	}

	for code, desired := range desiredByCode {
		actual := normalizeCompareValue(actualParams[code])
		if desired != actual {
			return false, nil
		}
	}

	return true, nil
}

func normalizeCompareValue(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.EqualFold(trimmed, "null") {
		return ""
	}
	if trimmed == "\"\"" {
		return ""
	}
	if strings.EqualFold(trimmed, "true") {
		return "true"
	}
	if strings.EqualFold(trimmed, "false") {
		return "false"
	}
	// FIX(uuid-case): нормализуем UUID к lowercase для сравнения c actual-значением
	// из API. Без этого "6214BA32-..." (из config) != "6214ba32-..." (из API state),
	// и resume считал бы параметры изменёнными, создавая новый инстанс вместо resume.
	if uuidLikeCompareRegex.MatchString(trimmed) {
		return strings.ToLower(trimmed)
	}
	return trimmed
}
