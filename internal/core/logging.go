package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ===== Логирование и debug-утилиты =====
//
//   - ctxKeyLogLevel / CtxWithLogLevel — переопределение LogLevel на уровне ресурса
//   - debugLogPath / isHTTPDebugEnabled — файловый debug HTTP (NUBES_DEBUG_HTTP=1)
//   - sanitizeAuthHeader — маскирование Authorization в дампах
//   - formatStageMsg — форматирование stageMsg операции

// ctxKeyLogLevel — ключ для переопределения LogLevel на уровне ресурса через context.WithValue.
type ctxKeyLogLevelType struct{}

var ctxKeyLogLevel = ctxKeyLogLevelType{}

func debugLogPath(fileName string) string {
	if dir := strings.TrimSpace(os.Getenv("NUBES_DEBUG_DIR")); dir != "" {
		return filepath.Join(dir, fileName)
	}
	return filepath.Join(os.TempDir(), fileName)
}

func isHTTPDebugEnabled() bool {
	return os.Getenv("NUBES_DEBUG_HTTP") == "1"
}

func sanitizeAuthHeader(h http.Header) http.Header {
	cloned := h.Clone()
	if cloned.Get("Authorization") != "" {
		cloned.Set("Authorization", "Bearer [REDACTED]")
	}
	return cloned
}

// CtxWithLogLevel возвращает ctx с переопределённым уровнем логирования.
func CtxWithLogLevel(ctx context.Context, level string) context.Context {
	return context.WithValue(ctx, ctxKeyLogLevel, level)
}

// reTimestamp удаляет временны́е метки вида [2026-03-24T05:07:47.716Z +0.860s] из строк.
var reTimestamp = regexp.MustCompile(`\[[0-9]{4}-[0-9]{2}-[0-9]{2}T[^\]]+\]\s*`)

// formatStageMsg парсит stageMsg (JSON-массив пар [имя, текст]) и возвращает
// отформатированные строки для вывода. Временны́е метки убираются.
func formatStageMsg(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	// Пытаемся распарсить как [["name","text"],...]
	var pairs [][]string
	if err := json.Unmarshal([]byte(raw), &pairs); err != nil || len(pairs) == 0 {
		return nil
	}
	// Вычисляем максимальную ширину имени для выравнивания
	maxLen := 0
	for _, p := range pairs {
		if len(p) >= 1 && len(p[0]) > maxLen {
			maxLen = len(p[0])
		}
	}
	var result []string
	for _, p := range pairs {
		if len(p) < 2 {
			continue
		}
		name := p[0]
		text := reTimestamp.ReplaceAllString(p[1], "")
		lines := strings.Split(strings.TrimSpace(text), "\n")
		padding := strings.Repeat(" ", maxLen-len(name))
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			if i == 0 {
				result = append(result, fmt.Sprintf("        %s:%s  %s", name, padding, line))
			} else {
				result = append(result, fmt.Sprintf("        %s   %s", strings.Repeat(" ", maxLen), line))
			}
		}
	}
	return result
}
