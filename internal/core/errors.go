package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ===== Ошибки ядра =====

// instanceDeletedError marks deleted instances to be treated as not found.
type instanceDeletedError struct {
	InstanceUID string
}

func (e *instanceDeletedError) Error() string {
	return fmt.Sprintf("экземпляр %s удален", e.InstanceUID)
}

func isInstanceDeletedError(err error) bool {
	var e *instanceDeletedError
	return errors.As(err, &e)
}

// ===== Текст ошибок API =====

// extractAPIErrorMessage вытаскивает человекочитаемый текст из тела ошибки API.
//
// Платформа отдаёт РАЗНЫЕ формы тела, и все они должны разбираться, иначе юзеру
// уйдёт сырой JSON (см. NOTES/30_analysis/ERROR_MATRIX_dns_zones_2026-10-16.md):
//
//  1. {"ERROR":"...","DETAIL":"..."}                     — legacy-прокси;
//  2. {"TITLE":"Bad Request","DETAIL":"...","TYPE":"..."} — RFC7807-подобная;
//  3. {"error":{"code":"...","message":"..."}}            — шлюз (401 и др.);
//  4. {"message":"..."} / {"detail":"..."}                — отдельные сервисы.
//
// Возвращает "", если тело пустое/не JSON/нет знакомых полей — вызывающий обязан
// в этом случае показать юзеру хоть что-то (сырое тело либо http.StatusText).
func extractAPIErrorMessage(body []byte) string {
	if strings.TrimSpace(string(body)) == "" {
		return ""
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return ""
	}

	// 3) вложенный объект error{code,message} / error{errors:[…]} / строка error.
	if nested, ok := raw["error"]; ok {
		switch v := nested.(type) {
		case string:
			if s := unescapeAPIText(v); s != "" {
				return s
			}
		case map[string]interface{}:
			code := unescapeAPIText(jsonScalarString(v["code"]))
			message := unescapeAPIText(jsonScalarString(v["message"]))
			switch {
			case code != "" && message != "":
				return code + ": " + message
			case message != "":
				return message
			case code != "":
				return code
			}
			// Google-style: {"error":{"errors":[{"message":…}, …]}}
			if s := errorArrayText(v["errors"]); s != "" {
				return s
			}
		}
	}

	// 1) legacy-прокси: ERROR (+ DETAIL). Проверяем ДО TITLE/DETAIL, иначе для
	// {"ERROR":"Internal error","DETAIL":"boom"} потерялось бы «Internal error».
	if legacy := unescapeAPIText(jsonScalarString(raw["ERROR"])); legacy != "" {
		if d := unescapeAPIText(jsonScalarString(raw["DETAIL"])); d != "" {
			return legacy + ": " + d
		}
		return legacy
	}

	// 2) RFC7807-подобная форма платформы: TITLE (+ DETAIL).
	title := unescapeAPIText(jsonScalarString(raw["TITLE"]))
	detail := unescapeAPIText(jsonScalarString(raw["DETAIL"]))
	switch {
	case title != "" && detail != "":
		return title + ": " + detail
	case title != "":
		return title
	case detail != "":
		return detail
	}

	// 4) message / detail в нижнем регистре.
	if m := unescapeAPIText(jsonScalarString(raw["message"])); m != "" {
		return m
	}
	if d := unescapeAPIText(jsonScalarString(raw["detail"])); d != "" {
		return d
	}

	// 5) массив ошибок: {"errors":["…", …]} / {"errors":[{"message":…}]}.
	if s := errorArrayText(raw["errors"]); s != "" {
		return s
	}

	// 6) Страховка: регистронезависимый поиск знакомых ключей (Message, Detail,
	// Description, …) — порядок фиксирован, чтобы результат был детерминированным.
	for _, want := range []string{"message", "detail", "error", "title", "description"} {
		for key, val := range raw {
			if strings.EqualFold(key, want) {
				if s := unescapeAPIText(jsonScalarString(val)); s != "" {
					return s
				}
			}
		}
	}

	return ""
}

// jsonScalarString аккуратно достаёт строку из произвольного JSON-значения.
// Объекты/массивы в тексты ошибок не тащим — для них сработает fallback (сырое тело).
func jsonScalarString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	default:
		return ""
	}
}

// errorArrayText достаёт текст из массива ошибок: []string или []{message:…}.
// Несколько сообщений объединяются через "; ".
func errorArrayText(v interface{}) string {
	arr, ok := v.([]interface{})
	if !ok {
		return ""
	}
	parts := make([]string, 0, len(arr))
	for _, item := range arr {
		switch t := item.(type) {
		case string:
			if s := unescapeAPIText(t); s != "" {
				parts = append(parts, s)
			}
		case map[string]interface{}:
			if s := unescapeAPIText(jsonScalarString(t["message"])); s != "" {
				parts = append(parts, s)
			} else if s := unescapeAPIText(jsonScalarString(t["error"])); s != "" {
				parts = append(parts, s)
			}
		}
	}
	return strings.Join(parts, "; ")
}

// unescapeAPIText убирает HTML-энтити из текстов платформы: они приходят прямо из
// описаний параметров спеки («integer &gt; 0», «&#39;A&#39;»). Один проход, без
// повторного разбора: &amp; стоит последним, чтобы «&amp;gt;» не стал «>».
func unescapeAPIText(s string) string {
	replacer := strings.NewReplacer(
		"&gt;", ">",
		"&lt;", "<",
		"&quot;", `"`,
		"&#39;", "'",
		"&apos;", "'",
		"&nbsp;", " ",
		"&amp;", "&",
	)
	return strings.TrimSpace(replacer.Replace(strings.TrimSpace(s)))
}
