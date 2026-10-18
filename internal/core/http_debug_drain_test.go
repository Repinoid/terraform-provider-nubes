package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDebugDumpDoesNotDrainRequestBody — регресс-тест на баг дебаг-дампа.
//
// БАГ (был): при NUBES_DEBUG_HTTP=1 дамп запроса делался так:
//
//	reqForDump := req.Clone(req.Context())               // Body — ТОТ ЖЕ reader
//	dump, _ := httputil.DumpRequestOut(reqForDump, true) // вычитывает Body в конец
//
// httputil.DumpRequestOut читает тело, а http.Request.Clone НЕ копирует Body
// (делит reader с оригиналом) → req.Body оказывается пустым → реальный POST
// уходит без тела и падает в транспорте:
//
//	Post "...": http: ContentLength=N with Body length 0
//
// Проверяем РЕАЛЬНЫЙ doRequest (не копию кода): сервер обязан получить тело целиком.
func TestDebugDumpDoesNotDrainRequestBody(t *testing.T) {
	// Включаем дебаг-дамп (именно он ломал тело). Каталог — временный, чтобы
	// файловый лог не писался в /tmp.
	t.Setenv("NUBES_DEBUG_HTTP", "1")
	t.Setenv("NUBES_DEBUG_DIR", t.TempDir())

	var gotBody []byte
	var gotReadErr error
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, gotReadErr = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &UniversalClient{
		HttpClient:     &http.Client{},
		ApiEndpoint:    srv.URL,
		ApiToken:       "test-token",
		RetryBaseDelay: time.Millisecond,
	}

	payload := map[string]interface{}{
		"serviceId":   109,
		"displayName": "dns-zone-test",
	}

	if _, _, err := c.doRequest(context.Background(), "POST", "/instances", payload); err != nil {
		t.Fatalf("doRequest(POST) вернул ошибку (тело вы(с)sосано дампом?): %v", err)
	}
	if gotReadErr != nil {
		t.Fatalf("чтение тела на сервере: %v", gotReadErr)
	}
	if len(gotBody) == 0 {
		t.Fatalf("сервер получил ПУСТОЕ тело — дебаг-дамп вы(с)sосал req.Body")
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(gotBody, &parsed); err != nil {
		t.Fatalf("тело не является JSON: %v (получено %q)", err, gotBody)
	}
	if parsed["serviceId"] != float64(109) || parsed["displayName"] != "dns-zone-test" {
		t.Fatalf("тело искажено: %q", gotBody)
	}
}

// TestDebugDumpKeepsGetWithoutBody — при включённом дебаге GET без тела не ломается
// (страховка: фикс не должен влиять на запросы без payload).
func TestDebugDumpKeepsGetWithoutBody(t *testing.T) {
	t.Setenv("NUBES_DEBUG_HTTP", "1")
	t.Setenv("NUBES_DEBUG_DIR", t.TempDir())

	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := &UniversalClient{
		HttpClient:     &http.Client{},
		ApiEndpoint:    srv.URL,
		ApiToken:       "test-token",
		RetryBaseDelay: time.Millisecond,
	}

	body, _, err := c.doRequest(context.Background(), "GET", "/instances/X", nil)
	if err != nil {
		t.Fatalf("doRequest(GET) вернул ошибку: %v", err)
	}
	if hits != 1 {
		t.Fatalf("ожидался 1 запрос к серверу, получено %d", hits)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("неожиданный ответ: %q", string(body))
	}
}
