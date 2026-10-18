package resources_core

import "strings"

// ErrorSummary возвращает человекочитаемое резюме ошибки — оно идёт в ЗАГОЛОВОК
// диагностики Terraform («Error: …»), чтобы юзер сразу понимал КАТЕГОРИЮ ошибки,
// не читая детали.
//
// Классификация — по тексту ошибки (деталь уже несёт «ошибка API 400: …» и т.п.).
// Порядок важен: конкретные коды (401/403/404) проверяются раньше общих диапазонов.
func ErrorSummary(err error) string {
	if err == nil {
		return "Ошибка"
	}
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "401") || strings.Contains(msg, "authentication") || strings.Contains(msg, "требуется авториз"):
		return "Ошибка авторизации"
	case strings.Contains(msg, "403") || strings.Contains(msg, "forbidden") || strings.Contains(msg, "доступ запрещ"):
		return "Доступ запрещён"
	case strings.Contains(msg, "404") || strings.Contains(msg, "not found") || strings.Contains(msg, "не найден"):
		return "Объект не найден"
	case strings.Contains(msg, "dial tcp") || strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "connection refused") || strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "tls") || strings.Contains(msg, "eof"):
		return "Нет связи с API"
	case strings.Contains(msg, "ошибка api 5"):
		return "Ошибка платформы"
	case strings.Contains(msg, "ошибка api 4"):
		return "Неверные данные"
	case strings.Contains(msg, "завершилась с ошибкой") || strings.Contains(msg, "таймаут"):
		return "Операция завершилась с ошибкой"
	default:
		return "Ошибка клиента"
	}
}
