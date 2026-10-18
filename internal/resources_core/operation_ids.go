package resources_core

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// BuildActionID builds a stable ID for action resources.
func BuildActionID(instanceUID string, operation string, runID string) string {
	return strings.TrimSpace(instanceUID) + ":" + strings.TrimSpace(operation) + ":" + strings.TrimSpace(runID)
}

// BuildSubresourceID builds a stable ID for subresource resources based on params.
func BuildSubresourceID(instanceUID string, subresource string, params map[string]string) string {
	payload := canonicalParamsString(params)
	h := sha256.Sum256([]byte(payload))
	return strings.TrimSpace(instanceUID) + ":" + strings.TrimSpace(subresource) + ":" + hex.EncodeToString(h[:])
}

func canonicalParamsString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	// Ключи храним как есть (camelCase, напр. "recordName") — иначе lookup по
	// lower-ключу в map вернёт пустую строку (баг: id становился константой).
	// Сортировка идёт по нормализованному (lower) ключу, но в payload ключ
	// нормализуется, а значение берётся по ИСХОДНОМУ ключу.
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(strings.TrimSpace(keys[i])) < strings.ToLower(strings.TrimSpace(keys[j]))
	})
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, strings.ToLower(strings.TrimSpace(k))+"="+strings.TrimSpace(params[k]))
	}
	return strings.Join(parts, ";")
}
