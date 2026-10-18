// Package jsonutil — нейтральный пакет каноничного JSON-сравнения.
//
// Вынесен из resources_core, чтобы его мог импортировать core (иначе был бы
// цикл импорта: core → resources_core → core). Сравнивает JSON-строки,
// игнорируя пробелы и порядок ключей объектов (порядок элементов массивов сохраняется).
package jsonutil

import (
	"bytes"
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// uuidAnywhereRegex — UUID-подобная подстрока внутри произвольного текста
// (например, внутри JSON-строки). Регистр не учитывается.
var uuidAnywhereRegex = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

// LowercaseUUIDsInText приводит все UUID-подобные подстроки к lowercase.
//
// Зачем: облако отдаёт один и тот же UUID то в lowercase (vdcUid), то в UPPERCASE
// (nsxtUid у кластера Штурвала), а пользователь может написать в config любой
// регистр. Сравнение plan vs actual обязано игнорировать регистр UUID.
func LowercaseUUIDsInText(value string) string {
	if !uuidAnywhereRegex.MatchString(value) {
		return value
	}
	return uuidAnywhereRegex.ReplaceAllStringFunc(value, strings.ToLower)
}

// JSONStringsEquivalent compares JSON strings ignoring whitespace and map key order.
func JSONStringsEquivalent(a, b string) bool {
	na, oka := normalizeJSONIfPossible(a)
	if !oka {
		return false
	}
	nb, okb := normalizeJSONIfPossible(b)
	if !okb {
		return false
	}
	return na == nb
}

func normalizeJSONIfPossible(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}
	var v interface{}
	if err := json.Unmarshal([]byte(trimmed), &v); err != nil {
		return "", false
	}
	return encodeCanonicalJSON(normalizeJSONScalarsToStrings(v)), true
}

func encodeCanonicalJSON(v interface{}) string {
	var buf bytes.Buffer
	writeCanonicalJSON(&buf, v)
	return buf.String()
}

func writeCanonicalJSON(buf *bytes.Buffer, v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(strconv.Quote(k))
			buf.WriteByte(':')
			writeCanonicalJSON(buf, val[k])
		}
		buf.WriteByte('}')
	case []interface{}:
		buf.WriteByte('[')
		for i, item := range val {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeCanonicalJSON(buf, item)
		}
		buf.WriteByte(']')
	case string:
		buf.WriteString(strconv.Quote(val))
	case float64:
		buf.WriteString(strconv.FormatFloat(val, 'g', -1, 64))
	case bool:
		if val {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case nil:
		buf.WriteString("null")
	default:
		if encoded, err := json.Marshal(val); err == nil {
			buf.Write(encoded)
		} else {
			buf.WriteString("null")
		}
	}
}

func normalizeJSONScalarsToStrings(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, child := range val {
			out[k] = normalizeJSONScalarsToStrings(child)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, child := range val {
			out[i] = normalizeJSONScalarsToStrings(child)
		}
		return out
	case string:
		return LowercaseUUIDsInText(val)
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		if encoded, err := json.Marshal(val); err == nil {
			return strings.TrimSpace(string(encoded))
		}
		return "null"
	}
}
