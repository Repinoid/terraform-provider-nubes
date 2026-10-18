package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

// ===== Internal HTTP helpers =====

func (c *UniversalClient) doRequest(ctx context.Context, method, path string, payload interface{}) ([]byte, http.Header, error) {
	// User-Agent: браузерный, чтобы пройти DDoS-Guard (см. docs/ops/API_TOKENS.md).
	// Go-http-client по умолчанию блокируется фильтром ddos-guard.
	const maxRetries = 3
	baseDelay := c.RetryBaseDelay
	if baseDelay <= 0 {
		baseDelay = 2 * time.Second
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<(attempt-1)) // 2s, 4s, 8s
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, nil, ctx.Err()
			}
		}

		var body io.Reader
		if payload != nil {
			b, err := json.Marshal(payload)
			if err != nil {
				return nil, nil, err
			}
			body = bytes.NewBuffer(b)
		}

		reqURL := c.buildURL(path)
		req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
		if err != nil {
			return nil, nil, err
		}

		// Форсируем новое TCP-соединение для POST/PUT/PATCH: DDoS-Guard может блочить их на keep-alive.
		// GET-запросы (поиск, чтение) оставляем на keep-alive — req.Close на них ломает DDoS-Guard.
		if method != "GET" {
			req.Close = true
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "*/*")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		if c.ApiToken != "" {
			req.Header.Set("Authorization", "Bearer "+c.ApiToken)
		}

		// DEBUG.
		// ⚠️ ВАЖНО: httputil.DumpRequestOut читает Body, а req.Clone НЕ копирует
		// тело (делит reader с оригиналом). Без подмены тела дамп высасывает
		// req.Body, и реальный запрос уходит пустым:
		//   Post "...": http: ContentLength=N with Body length 0
		// Поэтому для дампа берём СВЕЖУЮ копию тела через req.GetBody
		// (устанавливается автоматически http.NewRequest* для bytes.Buffer).
		if isHTTPDebugEnabled() {
			reqForDump := req.Clone(req.Context())
			reqForDump.Header = sanitizeAuthHeader(req.Header)
			if req.GetBody != nil {
				if freshBody, bodyErr := req.GetBody(); bodyErr == nil {
					reqForDump.Body = freshBody
				}
			}
			dump, _ := httputil.DumpRequestOut(reqForDump, body != nil)
			fmt.Fprintf(os.Stderr, "\n>>> REQ %s %s\n%s\n", method, path, dump)
		}

		// DEBUG в файл
		if isHTTPDebugEnabled() {
			f, _ := os.OpenFile(debugLogPath("nubes_debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if f != nil {
				fmt.Fprintf(f, ">>> %s %s\n", method, req.URL.String())
				f.Close()
			}
		}
		resp, err := c.HttpClient.Do(req)
		if err != nil {
			lastErr = err
			// Сетевые ошибки: retry только для идемпотентного GET.
			if attempt < maxRetries && method == "GET" {
				continue
			}
			return nil, nil, err
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			if attempt < maxRetries {
				continue
			}
			return nil, nil, err
		}

		// Retry только для transient-ошибок и только для GET
		if resp.StatusCode >= 400 {
			if attempt < maxRetries && method == "GET" && isRetryable(resp.StatusCode) {
				lastErr = formatAPIError(resp.StatusCode, respBody)
				continue
			}
			return nil, nil, formatAPIError(resp.StatusCode, respBody)
		}

		return respBody, resp.Header, nil
	}

	return nil, nil, fmt.Errorf("doRequest failed after %d retries: %w", maxRetries, lastErr)
}

// isRetryable returns true for transient HTTP errors that can be retried.
// 401 включён: Gateway может транзиентно отклонять валидный JWT (см. ARCHITECTURE.md,
// «API Resilience»). Токен между попытками НЕ обновляется (обновление — вне провайдера),
// поэтому при постоянном 401 это даст 3 холостых повтора с backoff.
// Ретраится только для GET (см. условие в doRequest).
func isRetryable(statusCode int) bool {
	return statusCode == http.StatusUnauthorized || // 401
		statusCode == http.StatusTooManyRequests || // 429
		statusCode == http.StatusServiceUnavailable || // 503
		statusCode == http.StatusBadGateway || // 502
		statusCode == http.StatusGatewayTimeout // 504
}

func (c *UniversalClient) postIgnoreResponse(ctx context.Context, path string, payload interface{}, returnLocation bool) (string, error) {
	respBody, headers, err := c.doRequest(ctx, "POST", path, payload)
	if err != nil {
		return "", err
	}

	if returnLocation {
		if loc := headers.Get("Location"); loc != "" {
			return extractUIDFromLocation(loc), nil
		}
	}

	var justId string
	if err := json.Unmarshal(respBody, &justId); err == nil && justId != "" {
		return justId, nil
	}

	return "", nil
}

func extractUIDFromLocation(loc string) string {
	if loc == "" {
		return ""
	}
	return strings.TrimPrefix(loc, "./")
}

// formatAPIError превращает ответ об ошибке API в читаемый текст.
//
// Текст ошибки выводится ВСЕГДА (требование оператора):
//  1. пробуем разобрать тело (см. extractAPIErrorMessage — все формы платформы);
//  2. не вышло — отдаём сырое тело как есть;
//  3. тела нет — стандартный текст HTTP-статуса.
//
// Так юзер никогда не остаётся без объяснения причины.
func formatAPIError(statusCode int, body []byte) error {
	msg := extractAPIErrorMessage(body)
	if msg == "" {
		msg = strings.TrimSpace(string(body))
	}
	if msg == "" {
		msg = http.StatusText(statusCode)
	}
	return fmt.Errorf("ошибка API %d: %s", statusCode, msg)
}
