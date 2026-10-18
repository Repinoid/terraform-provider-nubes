package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ===== Запуск операции с reconcile при баге платформы на POST /run =====
//
// КОНТЕКСТ (задокументировано по факту, стенд TEST, провайдер 3.4.0):
// платформа ДЕТЕРМИНИРОВАННО (не временно) иногда отвечает на
// POST /instanceOperations/{opUid}/run кодом 500 с телом
//   {"ERROR":"key [RESPONSEHEADER] doesn't exist"}
// ПРИ ЭТОМ операцию запускает и доводит до успеха. Если провайдер трактует такой
// 500 как провал — он возвращает ошибку ресурсу, Terraform НЕ обновляет state,
// хотя в облаке операция выполнена (ложный дрейф), а повторный apply создаёт
// дубль (см. NOTES/30_analysis/DEBUG_REPORT_ZONES_V2_RECORD_DUPLICATE_AND_ORPHAN.md,
// раздел 11, и REVIEW_SONNET_run_500_verify_reconcile.md).

// responseHeaderBugMarker — подпись тела ошибки платформы на POST /run.
// Сужаем условие ровно до неё (по ревью Sonnet): НЕ прощаем произвольный 5xx,
// чтобы реальный отказ /run не маскировался.
const responseHeaderBugMarker = "RESPONSEHEADER"

// runOperation запускает операцию (POST .../run) и ждёт её завершения по dtFinish.
// Это ЕДИНАЯ точка запуска для всех run-путей (create/modify/delete/redeploy/…).
// Если POST /run вернул ошибку с подписью responseHeaderBugMarker — проверяем,
// стартовала ли операция на самом деле (operationStarted), и при старте дожидаемся
// завершения. Иначе возвращаем исходную ошибку.
func (c *UniversalClient) runOperation(ctx context.Context, opUid string, timeout time.Duration) error {
	_, _, err := c.doRequest(ctx, "POST", fmt.Sprintf("/instanceOperations/%s/run", opUid), map[string]interface{}{})
	if err == nil {
		// НЕ МЕНЯТЬ: завершение операции определяется по dtFinish.
		return c.waitForOperationFinish(ctx, opUid, timeout)
	}
	if !isResponseHeaderBug(err) {
		return err
	}
	if c.operationStarted(ctx, opUid) {
		// Баг платформы подтверждён: операция реально стартовала, несмотря на 500
		// в ответе POST /run. Дожидаемся завершения и возвращаем реальный результат
		// (успех или errorLog), а не ложный отказ.
		return c.waitForOperationFinish(ctx, opUid, timeout)
	}
	// Признаков старта нет — это настоящий отказ POST /run, возвращаем исходную 500.
	return err
}

// isResponseHeaderBug определяет, что ошибка — тот самый детерминированный 500
// платформы («key [RESPONSEHEADER] doesn't exist»), а не произвольный сбой.
func isResponseHeaderBug(err error) bool {
	return err != nil && strings.Contains(strings.ToUpper(err.Error()), responseHeaderBugMarker)
}

// operationStarted коротким окном опроса проверяет, стартовала ли операция
// (pending / in-progress / уже завершена по dtFinish). Нужно, чтобы отличить
// «500, но операция пошла» от «500, операция не стартовала». Ошибки опроса и
// неоднозначность НЕ фатальны: считаем, что операция не стартовала, и вызывающий
// вернёт исходную 500 (реальный отказ не маскируем).
func (c *UniversalClient) operationStarted(ctx context.Context, opUid string) bool {
	const attempts = 3
	const interval = 2 * time.Second
	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-time.After(interval):
			case <-ctx.Done():
				return false
			}
		}
		respBody, _, err := c.doRequest(ctx, "GET",
			fmt.Sprintf("/instanceOperations/%s?fields=dtFinish,isSuccessful,errorLog,isInProgress,isPending,duration,stages", opUid), nil)
		if err != nil {
			continue
		}
		var status operationStatusResponse
		if err := json.Unmarshal(respBody, &status); err != nil {
			continue
		}
		op := status.InstanceOperation
		if op.IsPending || op.IsInProgress {
			return true
		}
		if op.DtFinish != nil && strings.TrimSpace(*op.DtFinish) != "" {
			return true
		}
	}
	return false
}
