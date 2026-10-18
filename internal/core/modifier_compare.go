package core

import (
	"context"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core/jsonutil"
)

// modifierRawValuesEqual сравнивает два значения БЕЗ схемы операции (dataType может быть
// недоступен до создания операции, или отличаться между default- и live-схемой):
//   - оба выглядят как JSON (`{`/`[`) — смысловое JSON-сравнение (порядок ключей не значим);
//   - иначе — строковое сравнение с нормализацией null/""/true/false (без учёта регистра).
func modifierRawValuesEqual(wanted, current string) bool {
	w := strings.TrimSpace(wanted)
	c := strings.TrimSpace(current)
	if looksLikeJSON(w) && looksLikeJSON(c) {
		return jsonutil.JSONStringsEquivalent(w, c)
	}
	return normalizeRawScalar(w) == normalizeRawScalar(c)
}

func looksLikeJSON(v string) bool {
	return strings.HasPrefix(v, "{") || strings.HasPrefix(v, "[")
}

// normalizeRawScalar приводит скаляр к каноническому виду без схемы:
// null/"" -> "", true/false — к lower-case, остальное — trim.
func normalizeRawScalar(v string) string {
	switch {
	case strings.EqualFold(v, "null"), v == `""`:
		return ""
	case strings.EqualFold(v, "true"):
		return "true"
	case strings.EqualFold(v, "false"):
		return "false"
	}
	return v
}

// modifierDesiredEqualsLive — idempotency-pre-check по ЖИВЫМ параметрам инстанса
// (state.params). Источник — ЕДИНСТВЕННЫЙ: live. Схема операции НЕ запрашивается (см.
// ARCHITECTURE.md, «Modifier idempotency»): до создания операции доступна только
// default-схема, а она может расходиться с живой — сравнение по ней рискует ложным
// пропуском. Поиск значения — по самому коду (live keyed by lower(code)).
//
// Fail-safe: live пуст, код не найден или ошибка чтения live → (false, …): modify
// лучше выполнить лишний раз, чем пропустить нужный.
func (c *UniversalClient) modifierDesiredEqualsLive(ctx context.Context, instanceUid string, desired map[string]string) (bool, error) {
	if len(desired) == 0 {
		return false, nil
	}
	live, err := c.instanceLiveParams(ctx, instanceUid)
	if err != nil {
		return false, err
	}
	if len(live) == 0 {
		return false, nil
	}
	for code, wanted := range desired {
		current, ok := live[strings.ToLower(strings.TrimSpace(code))]
		if !ok {
			return false, nil
		}
		if !modifierRawValuesEqual(wanted, current) {
			return false, nil
		}
	}
	return true, nil
}

// modifierCodeMap индексирует cfsParams по нижнему регистру Code и SvcOperationCfsParam.
func modifierCodeMap(cfsParams []universalCfsParam) map[string]universalCfsParam {
	m := make(map[string]universalCfsParam, len(cfsParams)*2)
	for _, p := range cfsParams {
		if key := strings.ToLower(strings.TrimSpace(p.Code)); key != "" {
			m[key] = p
		}
		if key := strings.ToLower(strings.TrimSpace(p.SvcOperationCfsParam)); key != "" {
			m[key] = p
		}
	}
	return m
}
