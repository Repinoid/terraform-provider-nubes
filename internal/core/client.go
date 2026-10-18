// Package core — ядро провайдера: HTTP-клиент Nubes API.
//
// Содержит:
//   - UniversalClient — клиент с retry, таймаутами, Bearer-токеном
//   - HTTP-транспорт (http.go), DTO (types.go), логирование (logging.go)
//   - Методы состояния инстанса (instance_state.go), поиска (instance_find.go)
//   - Запуск и ожидание операций (operation_run.go, operation_wait.go)
//   - Резолв refSvc-параметров (refsvc.go), нормализация значений (params.go)
//
// ⛔ НЕ ИЗМЕНЯТЬ НИЧЕГО В ЯДРЕ БЕЗ ПРЯМОГО РАЗРЕШЕНИЯ ОПЕРАТОРА.
package core

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// UniversalClient handles Nubes API logic.
type UniversalClient struct {
	HttpClient        *http.Client
	ApiEndpoint       string
	ApiToken          string
	ProviderVersion   string
	OperationTimeouts *OperationTimeouts
	// LogLevel задаёт уровень вывода этапов операции: ""/"none" | "info" | "debug".
	// Может быть переопределён на уровне ресурса через context (ключ ctxKeyLogLevel).
	LogLevel string
	// PollInterval — интервал поллинга. 0 = default 5s. Для тестов.
	PollInterval time.Duration
	// StagesWriter — writer для вывода этапов. nil = ttyOut(). Для тестов.
	StagesWriter io.Writer
	// RetryBaseDelay — базовая задержка retry в doRequest. 0 = default 2s. Для тестов.
	RetryBaseDelay time.Duration
}

// instanceMutexes — глобальная карта мьютексов для сериализации операций на одном инстансе.
// Terraform параллелит ресурсы (до 10), но API Nubes не поддерживает параллельные операции
// на одном инстансе (orchestrator error, state_out inconsistency).
var instanceMutexes sync.Map

// LockInstance блокирует мьютекс для указанного instanceUid.
// Возвращает функцию unlock, которую нужно вызывать через defer.
func (c *UniversalClient) LockInstance(instanceUid string) func() {
	mu, _ := instanceMutexes.LoadOrStore(instanceUid, &sync.Mutex{})
	mu.(*sync.Mutex).Lock()
	return func() { mu.(*sync.Mutex).Unlock() }
}

// isProxyAPI returns true if ApiEndpoint uses legacy ?endpoint= proxy pattern (contains "index.cfm").
func (c *UniversalClient) isProxyAPI() bool {
	return strings.Contains(c.ApiEndpoint, "index.cfm")
}

// buildURL constructs a full URL from the base endpoint + path.
// Legacy proxy: index.cfm?endpoint=/instances&page=1
// New REST gateway: /api/v1/svc/instances?page=1
func (c *UniversalClient) buildURL(path string) string {
	if c.isProxyAPI() {
		endpointPath := path
		extraQuery := ""
		if idx := strings.Index(path, "?"); idx >= 0 {
			endpointPath = path[:idx]
			extraQuery = path[idx+1:]
		}
		rawQuery := "endpoint=" + endpointPath
		if extraQuery != "" {
			rawQuery += "&" + extraQuery
		}
		return c.ApiEndpoint + "?" + rawQuery
	}
	return c.ApiEndpoint + path
}

// userAgent — браузерный, чтобы пройти DDoS-Guard.
// Go-http-client по умолчанию блокируется фильтром ddos-guard.
const userAgent = "Mozilla/5.0"

// instanceDescr returns a stable description with provider version when available.
func (c *UniversalClient) instanceDescr() string {
	version := strings.TrimSpace(c.ProviderVersion)
	if version == "" {
		return "Создано через Nubes Terraform Universal Provider"
	}
	return fmt.Sprintf("Создано через Nubes Terraform Universal Provider %s", version)
}
