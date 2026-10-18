package resources_core

import (
	"context"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// FindSubresourceInStateOut checks presence of a subresource in instance state_out.
// Returns (found, known, error). If known=false, the list key was not found.
func FindSubresourceInStateOut(ctx context.Context, client *core.UniversalClient, instanceID string, listKey string, idKey string, target string) (bool, bool, error) {
	if client == nil {
		return false, false, nil
	}
	if strings.TrimSpace(instanceID) == "" {
		return false, false, nil
	}
	if strings.TrimSpace(listKey) == "" || strings.TrimSpace(idKey) == "" || strings.TrimSpace(target) == "" {
		return false, false, nil
	}

	details, err := client.GetInstanceStateDetails(ctx, instanceID)
	if err != nil {
		return false, false, err
	}
	if details == nil {
		return false, false, nil
	}

	// Список подресурсов лежит НЕ всегда в state.out. Примеры:
	//   - postgres: пользователи/БД → state.out.users / state.out.databases;
	//   - zones_v2: записи зоны      → state.params.records (в state.out только
	//     dnsServers/zoneSerial/zoneSoaServer).
	// Поэтому ищем в обеих половинах состояния. Порядок: out → params
	// (историческое поведение сохраняется для сервисов, где ключ есть в out).
	sources := []map[string]interface{}{details.RawOut, details.RawParams}
	known := false
	for _, src := range sources {
		if src == nil {
			continue
		}
		listRaw, ok := src[listKey]
		if !ok {
			continue
		}
		found, recognized := subresourceListContains(listRaw, idKey, target)
		if !recognized {
			continue
		}
		known = true
		if found {
			return true, true, nil
		}
	}

	return false, known, nil
}

// subresourceListContains ищет target в списке подресурсов.
// recognized=false означает, что значение не является списком/картой объектов —
// искать по нему бессмысленно, вызывающий должен продолжить поиск в другой
// половине состояния. Пустой список считается recognized=true (ключ есть, но
// объектов нет) — это сохраняет исходную семантику known/found.
func subresourceListContains(listRaw interface{}, idKey string, target string) (bool, bool) {
	switch list := listRaw.(type) {
	case map[string]interface{}:
		for key, item := range list {
			if strings.EqualFold(strings.TrimSpace(key), strings.TrimSpace(target)) {
				return true, true
			}
			if m, ok := item.(map[string]interface{}); ok {
				if matchSubresourceID(m, idKey, target) {
					return true, true
				}
			}
		}
		return false, true
	case []interface{}:
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				if matchSubresourceID(m, idKey, target) {
					return true, true
				}
			}
		}
		return false, true
	default:
		return false, false
	}
}

func matchSubresourceID(raw map[string]interface{}, idKey string, target string) bool {
	if raw == nil {
		return false
	}
	for key, val := range raw {
		if !strings.EqualFold(strings.TrimSpace(key), strings.TrimSpace(idKey)) {
			continue
		}
		value := strings.TrimSpace(normalizeRawValue(val))
		if value == "" {
			return false
		}
		return strings.EqualFold(value, strings.TrimSpace(target))
	}
	return false
}

func normalizeRawValue(value interface{}) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func IsSubresourceMissingError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	markers := []string{"не найден", "not found", "does not exist", "отсутствует"}
	for _, marker := range markers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}

func IsSubresourceAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	markers := []string{"уже существует", "already exists", "already exist", "duplicate", "conflict", " exist"}
	for _, marker := range markers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}

func IsSubresourceTransientDependencyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if msg == "" {
		return false
	}
	markers := []string{"ошибка api 424", "missing dependency", "orchestrator communication error:404 not found"}
	for _, marker := range markers {
		if strings.Contains(msg, marker) {
			return true
		}
	}
	return false
}

// ResolveSubresourceStateKeys returns state_out list key and identity key for a subresource.
// The second argument preserves generator-level identity priority and helps avoid hardcoding
// per-service logic in generated resources.
func ResolveSubresourceStateKeys(subresource string, identityCodes []string) (string, string) {
	listKey := SubresourceListKey(subresource)
	idKey := SubresourceIdentityKey(subresource, identityCodes)
	return listKey, idKey
}

// SubresourceListKey builds a best-effort list key in state_out.
func SubresourceListKey(subresource string) string {
	name := strings.TrimSpace(subresource)
	if name == "" {
		return ""
	}
	if strings.HasSuffix(strings.ToLower(name), "s") {
		return name
	}
	return name + "s"
}

// SubresourceIdentityKey picks identity field key for lookup in state_out.
func SubresourceIdentityKey(subresource string, identityCodes []string) string {
	name := strings.ToLower(strings.TrimSpace(subresource))
	if name == "" {
		return ""
	}

	special := map[string]string{
		"user":     "username",
		"database": "dbName",
	}
	if preferred, ok := special[name]; ok {
		if matched := findMatchingCode(identityCodes, preferred); matched != "" {
			return matched
		}
		return preferred
	}

	priority := []string{"name", name, name + "name", "id", name + "id"}
	for _, candidate := range priority {
		if matched := findMatchingCode(identityCodes, candidate); matched != "" {
			return matched
		}
	}

	if len(identityCodes) > 0 {
		return strings.TrimSpace(identityCodes[0])
	}

	return ""
}

func findMatchingCode(identityCodes []string, candidate string) string {
	target := strings.ToLower(strings.TrimSpace(candidate))
	if target == "" {
		return ""
	}
	for _, code := range identityCodes {
		trimmed := strings.TrimSpace(code)
		if strings.EqualFold(trimmed, target) {
			return trimmed
		}
	}
	return ""
}
