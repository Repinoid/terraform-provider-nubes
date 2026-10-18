package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ===== Ожидание операций (APPEND-ONLY) =====
//
// КРИТИЧЕСКИ ВАЖНО:
//  1. Критерий завершения операции — наличие dtFinish.
//  2. Эта логика едина для всех сервисов/операций и является контрактом поведения.
//  3. ЗАПРЕЩЕНО ПРАВИТЬ ЭТОТ КОД БЕЗ ЯВНОГО СОГЛАСОВАНИЯ С ОПЕРАТОРОМ.
//     Любые изменения (таймауты, критерии завершения, частота опроса, обработка ошибок)
//     должны быть согласованы заранее.
const defaultOperationTimeout = 30 * time.Minute

// ttyOut открывает /dev/tty для прямого вывода в терминал, минуя перехват terraform.
func ttyOut() *os.File {
	if f, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0); err == nil {
		return f
	}
	return os.Stderr
}

func (c *UniversalClient) waitForOperationFinish(ctx context.Context, opUid string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	interval := c.PollInterval
	if interval <= 0 {
		interval = 5 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	printedStages := make(map[string]bool)
	var lastPendingUID string

	var tty io.Writer
	if c.StagesWriter != nil {
		tty = c.StagesWriter
	} else {
		f := ttyOut()
		defer f.Close()
		tty = f
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("операция %s отменена", opUid)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("операция %s не завершилась за установленный таймаут в %.0f секунд", opUid, timeout.Seconds())
			}

			respBody, _, err := c.doRequest(ctx, "GET",
				fmt.Sprintf("/instanceOperations/%s?fields=dtFinish,isSuccessful,errorLog,isInProgress,isPending,duration,stages", opUid), nil)
			if err != nil {
				return fmt.Errorf("не удалось проверить статус операции %s: %w", opUid, err)
			}

			var status operationStatusResponse
			if err := json.Unmarshal(respBody, &status); err != nil {
				return fmt.Errorf("не удалось разобрать статус операции %s: %w", opUid, err)
			}

			// Печатаем этапы по мере выполнения — ВСЕГДА, без гейтинга.
			// Завершённые (DtFinish != nil): [OK ]/[FAIL] с duration.
			// Текущий (DtFinish == nil): [..] один раз при появлении.
			for _, stage := range status.InstanceOperation.Stages {
				if printedStages[stage.InstanceOperationStageUid] {
					continue
				}
				if stage.DtFinish != nil && *stage.DtFinish != "" {
					// Завершённый этап
					printedStages[stage.InstanceOperationStageUid] = true
					if stage.InstanceOperationStageUid == lastPendingUID {
						lastPendingUID = ""
					}
					status2 := "OK "
					if !stage.IsSuccessful {
						status2 = "FAIL"
					}
					fmt.Fprintf(tty, "  [%s] %s — %.1f sec\n", status2, stage.Stage, stage.Duration)
				} else if stage.InstanceOperationStageUid != lastPendingUID {
					// Текущий этап — показываем один раз
					lastPendingUID = stage.InstanceOperationStageUid
					fmt.Fprintf(tty, "  [..] %s\n", stage.Stage)
				}
			}

			// Критерий завершения — dtFinish (НЕ МЕНЯТЬ)
			if status.InstanceOperation.DtFinish != nil && strings.TrimSpace(*status.InstanceOperation.DtFinish) != "" {
				if status.InstanceOperation.IsSuccessful != nil && !*status.InstanceOperation.IsSuccessful {
					if status.InstanceOperation.ErrorLog != nil && strings.TrimSpace(*status.InstanceOperation.ErrorLog) != "" {
						return fmt.Errorf("операция %s завершилась с ошибкой: %s", opUid, unescapeAPIText(*status.InstanceOperation.ErrorLog))
					}
					return fmt.Errorf("операция %s завершилась с ошибкой", opUid)
				}
				if status.InstanceOperation.Duration != nil {
					fmt.Fprintf(tty, "  [DONE] %.1f sec\n", *status.InstanceOperation.Duration)
				}
				return nil
			}
		}
	}
}

// waitForInstanceIdle waits until no operation is pending/in-progress for the instance.
func (c *UniversalClient) waitForInstanceIdle(ctx context.Context, instanceUid string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("ожидание операции отменено для экземпляра %s", instanceUid)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("экземпляр %s не перешёл в состояние ожидания за установленный таймаут в %.0f секунд", instanceUid, timeout.Seconds())
			}

			state, err := c.GetInstanceState(ctx, instanceUid)
			if err != nil {
				return fmt.Errorf("не удалось проверить состояние экземпляра %s: %w", instanceUid, err)
			}
			if !state.OperationIsPending && !state.OperationIsInProgress {
				return nil
			}
		}
	}
}
