// ref_validation.go — Валидация ref-параметров при adopt/resume.
// Проверяет что все ссылочные параметры (vappUid и т.п.) указывают на живые инстансы
// и совпадают с плановыми значениями.
package resources_core

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// RefParamIssue описывает проблему с ref-параметром при adopt.
type RefParamIssue struct {
	ParamCode    string // код параметра, напр. "vappUid"
	PlannedValue string // UUID из Terraform плана
	ActualValue  string // UUID из state.params adopted-инстанса
	RefSvcId     int    // service_id referenced ресурса
	RefStatus    string // статус referenced инстанса ("deleted", "suspended", "running", "not_found")
	IsHardError  bool   // true = блокирует adopt, false = warning
	HumanMessage string // человекочитаемое сообщение
}

// ValidateRefParamsOnAdopt проверяет ref-параметры adopted/resumed инстанса.
// Сравнивает плановые значения (params) с фактическими из API (state.params),
// и для каждого ref-параметра проверяет статус referenced инстанса.
//
// Возвращает список проблем. Если есть хотя бы один IsHardError — adopt невозможен.
func ValidateRefParamsOnAdopt(ctx context.Context, client *core.UniversalClient, serviceID int, plannedParams map[int]string, adoptedInstanceUID string) ([]RefParamIssue, error) {
	if client == nil || serviceID <= 0 || strings.TrimSpace(adoptedInstanceUID) == "" {
		return nil, nil
	}

	// Получить карту ref-параметров для этого сервиса (paramID → refSvcId)
	refMap, err := getServiceParamRefMap(serviceID)
	if err != nil {
		return nil, err
	}
	if len(refMap) == 0 {
		return nil, nil // нет ref-параметров у этого сервиса
	}

	// Получить карту paramID → code
	codeMap, err := getServiceParamCodeMap(serviceID)
	if err != nil {
		return nil, err
	}

	// Получить фактические параметры из API
	actualParams, err := client.GetInstanceStateParams(ctx, adoptedInstanceUID)
	if err != nil {
		return nil, err
	}

	var issues []RefParamIssue

	for paramID, refSvcId := range refMap {
		if refSvcId <= 0 {
			continue
		}
		code := codeMap[paramID]
		if strings.TrimSpace(code) == "" {
			continue
		}

		// Плановое значение (resolve через ResolveRefSvcParamValue)
		plannedRaw, hasPlan := plannedParams[paramID]
		if !hasPlan || strings.TrimSpace(plannedRaw) == "" {
			continue
		}
		plannedResolved, err := client.ResolveRefSvcParamValue(ctx, refSvcId, plannedRaw)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(plannedResolved) == "" {
			plannedResolved = plannedRaw
		}

		// Фактическое значение из adopted инстанса
		actualRaw := strings.TrimSpace(actualParams[code])
		if actualRaw == "" {
			continue // параметр не задан в API, пропускаем
		}
		actualResolved, err := client.ResolveRefSvcParamValue(ctx, refSvcId, actualRaw)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(actualResolved) == "" {
			actualResolved = actualRaw
		}

		// Сравнение: план vs факт
		if strings.EqualFold(plannedResolved, actualResolved) {
			// UUID совпадают — проверить только статус referenced инстанса
			issue := checkRefInstanceStatus(ctx, client, code, plannedResolved, actualResolved, refSvcId)
			if issue != nil {
				issues = append(issues, *issue)
			}
			continue
		}

		// UUID НЕ совпадают — проверить статус обоих
		actualStatus := getRefInstanceStatus(ctx, client, actualResolved)
		plannedStatus := getRefInstanceStatus(ctx, client, plannedResolved)

		issues = append(issues, RefParamIssue{
			ParamCode:    code,
			PlannedValue: plannedResolved,
			ActualValue:  actualResolved,
			RefSvcId:     refSvcId,
			RefStatus:    actualStatus,
			IsHardError:  true,
			HumanMessage: fmt.Sprintf(
				"НЕСОВПАДЕНИЕ REF-ПАРАМЕТРА '%s' ПРИ ADOPT\n"+
					"\n"+
					"Инстанс хранит ссылку %s = %s (статус: %s),\n"+
					"но конфигурация Terraform указывает на %s (статус: %s).\n"+
					"\n"+
					"Adopt невозможен — ресурс привязан к другому зависимому инстансу.\n"+
					"\n"+
					"Рекомендации:\n"+
					"  1. Удалите suspended инстанс из облака (ЛК → Управление облаком) и повторите apply.\n"+
					"  2. Или восстановите старый зависимый инстанс (%s) через ЛК, если он нужен.\n"+
					"  3. Или выполните terraform destroy и затем terraform apply для полного пересоздания.",
				code,
				code, actualResolved, actualStatus,
				plannedResolved, plannedStatus,
				actualResolved,
			),
		})
	}

	return issues, nil
}

// checkRefInstanceStatus проверяет статус referenced инстанса при совпадающих UUID.
func checkRefInstanceStatus(ctx context.Context, client *core.UniversalClient, paramCode string, plannedUID string, actualUID string, refSvcId int) *RefParamIssue {
	status := getRefInstanceStatus(ctx, client, actualUID)

	switch status {
	case "deleted":
		return &RefParamIssue{
			ParamCode:    paramCode,
			PlannedValue: plannedUID,
			ActualValue:  actualUID,
			RefSvcId:     refSvcId,
			RefStatus:    status,
			IsHardError:  true,
			HumanMessage: fmt.Sprintf(
				"REF-ПАРАМЕТР '%s' УКАЗЫВАЕТ НА УДАЛЁННЫЙ ИНСТАНС\n"+
					"\n"+
					"Параметр %s = %s, но этот инстанс в статусе 'deleted'.\n"+
					"Adopt невозможен — зависимый ресурс удалён из облака.\n"+
					"\n"+
					"Рекомендации:\n"+
					"  1. Удалите suspended инстанс из облака и пересоздайте всё через terraform apply.\n"+
					"  2. Или восстановите зависимый инстанс через ЛК.",
				paramCode, paramCode, actualUID,
			),
		}
	case "suspended":
		return &RefParamIssue{
			ParamCode:    paramCode,
			PlannedValue: plannedUID,
			ActualValue:  actualUID,
			RefSvcId:     refSvcId,
			RefStatus:    status,
			IsHardError:  false,
			HumanMessage: fmt.Sprintf(
				"ВНИМАНИЕ: ref-параметр '%s' указывает на приостановленный инстанс\n"+
					"\n"+
					"Параметр %s = %s, статус: 'suspended'.\n"+
					"Убедитесь, что зависимый ресурс будет активирован (resume) перед использованием.",
				paramCode, paramCode, actualUID,
			),
		}
	case "not_found":
		return &RefParamIssue{
			ParamCode:    paramCode,
			PlannedValue: plannedUID,
			ActualValue:  actualUID,
			RefSvcId:     refSvcId,
			RefStatus:    status,
			IsHardError:  true,
			HumanMessage: fmt.Sprintf(
				"REF-ПАРАМЕТР '%s' УКАЗЫВАЕТ НА НЕСУЩЕСТВУЮЩИЙ ИНСТАНС\n"+
					"\n"+
					"Параметр %s = %s, но инстанс не найден в API (404).\n"+
					"Проверьте облако — инстанс мог быть полностью удалён.\n"+
					"\n"+
					"Рекомендации:\n"+
					"  1. Удалите suspended инстанс из облака и пересоздайте всё.\n"+
					"  2. Проверьте зависимости в конфигурации Terraform.",
				paramCode, paramCode, actualUID,
			),
		}
	}
	// running или другой нормальный статус — всё ОК
	return nil
}

// getRefInstanceStatus возвращает упрощённый статус инстанса: "running", "suspended", "deleted", "not_found", или raw status.
func getRefInstanceStatus(ctx context.Context, client *core.UniversalClient, instanceUID string) string {
	if strings.TrimSpace(instanceUID) == "" {
		return "not_found"
	}
	state, err := client.GetInstanceStateRaw(ctx, instanceUID)
	if err != nil {
		return "not_found"
	}
	if state == nil {
		return "not_found"
	}
	if state.IsDeleted {
		return "deleted"
	}
	status := strings.ToLower(strings.TrimSpace(state.ExplainedStatus))
	if status == "deleted" {
		return "deleted"
	}
	if strings.Contains(status, "suspend") {
		return "suspended"
	}
	if strings.Contains(status, "running") || strings.Contains(status, "active") {
		return "running"
	}
	return status
}

// FormatRefParamIssues форматирует список проблем для diagnostic-сообщения.
func FormatRefParamIssues(issues []RefParamIssue) string {
	if len(issues) == 0 {
		return ""
	}
	var parts []string
	for _, issue := range issues {
		parts = append(parts, issue.HumanMessage)
	}
	return strings.Join(parts, "\n\n---\n\n")
}

// HasHardErrors возвращает true если есть хотя бы одна блокирующая проблема.
func HasHardErrors(issues []RefParamIssue) bool {
	for _, issue := range issues {
		if issue.IsHardError {
			return true
		}
	}
	return false
}

// FormatDuplicateInstancesError формирует сообщение о множественных инстансах с одним именем.
func FormatDuplicateInstancesError(displayName string, serviceID int, instances []core.InstanceStateResponse) string {
	if len(instances) <= 1 {
		return ""
	}

	var parts []string
	statusCounts := map[string]int{}
	for _, inst := range instances {
		status := strings.ToLower(strings.TrimSpace(inst.ExplainedStatus))
		statusCounts[status]++
		parts = append(parts, fmt.Sprintf("  - %s (статус: %s)", inst.InstanceUid, status))
	}

	sort.Strings(parts)

	return fmt.Sprintf(
		"ОБНАРУЖЕНО %d ИНСТАНСОВ С ИМЕНЕМ '%s' (serviceId=%d)\n"+
			"\n"+
			"Список:\n"+
			"%s\n"+
			"\n"+
			"Невозможно определить какой инстанс adopt-ить.\n"+
			"\n"+
			"Рекомендации:\n"+
			"  1. Удалите лишние инстансы через ЛК (Управление облаком).\n"+
			"  2. Или укажите конкретный UUID через 'terraform import'.",
		len(instances), displayName, serviceID, strings.Join(parts, "\n"),
	)
}
