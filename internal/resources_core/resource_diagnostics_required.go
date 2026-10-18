package resources_core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// PlanExistingResourceDiagnosticsWithParams returns diagnostics when resource exists and compares required params.
// Устаревшая обёртка: считает, что сервис поддерживает suspend/resume (старое поведение).
func PlanExistingResourceDiagnosticsWithParams(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool, desiredParams map[int]string) diag.Diagnostics {
	return PlanExistingResourceDiagnosticsWithParamsAndDomainAndServices(ctx, client, serviceID, resourceName, adoptExistingOnCreate, desiredParams, "", nil, true)
}

// PlanExistingResourceDiagnosticsWithParamsAndDomain — см. PlanExistingResourceDiagnosticsWithParams.
func PlanExistingResourceDiagnosticsWithParamsAndDomain(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool, desiredParams map[int]string, desiredDomain string) diag.Diagnostics {
	return PlanExistingResourceDiagnosticsWithParamsAndDomainAndServices(ctx, client, serviceID, resourceName, adoptExistingOnCreate, desiredParams, desiredDomain, nil, true)
}

// PlanExistingResourceDiagnosticsWithParamsAndDomainAndServices — create-time проверка
// «инстанс с таким именем уже существует».
//
// supportsSuspend — поддерживает ли сервис операции suspend/resume (из YAML-спеки).
// Если НЕТ (например vc_nsxt/Edge), то:
//   - советы про resume/adopt для состояния Suspend бессмысленны — такое состояние
//     для такого сервиса вообще недостижимо, а усыновление требует resume;
//   - корректные действия пользователя: удалить инстанс в ЛК или взять другое имя.
func PlanExistingResourceDiagnosticsWithParamsAndDomainAndServices(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool, desiredParams map[int]string, desiredDomain string, allowedDomainServiceIDs []int, supportsSuspend bool) diag.Diagnostics {
	f, _ := os.OpenFile("/tmp/nubes_find_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		fmt.Fprintf(f, "[FIND-DEBUG] PlanExistingResourceDiagnostics entered: serviceId=%d name=%q adopt=%v client=%v\n", serviceID, resourceName, adoptExistingOnCreate, client != nil)
		f.Close()
	}

	var diags diag.Diagnostics
	appendLocalPlanInputDiagnostics(serviceID, desiredParams, desiredDomain, &diags)
	appendRemotePlanInputDiagnostics(ctx, client, serviceID, desiredParams, &diags)

	if client == nil {
		return diags
	}
	if strings.TrimSpace(resourceName) == "" {
		return diags
	}

	existing, err := client.FindInstanceByDisplayName(ctx, serviceID, resourceName)
	if err != nil || existing == nil {
		appendDomainCollisionDiagnostics(ctx, client, serviceID, resourceName, desiredDomain, allowedDomainServiceIDs, nil, &diags)
		return diags
	}
	appendDomainCollisionDiagnostics(ctx, client, serviceID, resourceName, desiredDomain, allowedDomainServiceIDs, existing, &diags)
	if diags.HasError() {
		return diags
	}
	if diags.HasError() {
		return diags
	}

	status := normalizeStatus(existing.ExplainedStatus)
	if isStatusNotCreated(status) {
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (NOT CREATED)",
			"Инстанс найден, но находится в состоянии Not Created.\nУдалите его в Личном кабинете (Управление облаком), затем повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if isStatusRunning(status) {
		if adoptExistingOnCreate {
			appendAdoptOrphanGuardDiagnostics(ctx, client, existing, serviceID, resourceName, &diags)
			return diags
		}
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (RUNNING)",
			"Инстанс с таким resource_name уже работает.\nАвто-усыновление отключено (adopt_existing_on_create=false).\nВыберите другое имя, если нужен новый ресурс."+runningConflictHint(supportsSuspend)+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if !isStatusSuspended(status) {
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (НЕПОДДЕРЖИВАЕМЫЙ СТАТУС)",
			"Инстанс найден, но статус не подходит для авто-усыновления.\nСтатус: "+status+". Выберите другое имя или устраните проблему в ЛК."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}

	if !adoptExistingOnCreate {
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (SUSPEND)",
			"Инстанс найден в состоянии Suspend, но adopt_existing_on_create=false."+suspendConflictHint(supportsSuspend)+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}

	mismatches, err := RequiredParamsMismatch(ctx, client, serviceID, desiredParams, existing.InstanceUid)
	if err != nil {
		diags.AddError("Не удалось сравнить параметры", err.Error())
		return diags
	}
	if len(mismatches) == 0 {
		return diags
	}

	details := FormatRequiredParamMismatchMessage(mismatches)
	diags.AddError(
		"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (SUSPEND, ПАРАМЕТРЫ ОТЛИЧАЮТСЯ)",
		"Обязательные параметры не совпадают с существующим инстансом в состоянии Suspend: "+details+".\nИсправьте конфиг и повторите plan/apply."+formatInstanceDetails(existing, serviceID, resourceName),
	)
	return diags
}

func appendRemotePlanInputDiagnostics(ctx context.Context, client *core.UniversalClient, serviceID int, desiredParams map[int]string, diags *diag.Diagnostics) {
	if diags == nil || client == nil || serviceID <= 0 || len(desiredParams) == 0 {
		return
	}

	specByID, err := getServiceParamValidationSpecMap(serviceID)
	if err != nil || len(specByID) == 0 {
		return
	}

	realmCache := []string(nil)
	realmLoaded := false
	refOptionsCache := map[int][]core.RefServiceInstanceOption{}

	for id, rawValue := range desiredParams {
		spec, ok := specByID[id]
		if !ok {
			continue
		}
		value := strings.TrimSpace(rawValue)
		if value == "" {
			continue
		}
		code := strings.TrimSpace(spec.Code)

		if len(spec.ValueList) > 0 && !containsValueCaseSensitive(spec.ValueList, value) {
			diags.AddError(
				"Некорректное значение параметра",
				"Параметр "+code+" содержит недопустимое значение: \""+value+"\".\nДоступные значения: "+joinLimited(spec.ValueList, 25),
			)
		}

		// TODO: realm validation disabled — backend bug in /resourceRealms/available (500, available_resource_realms.cfc:79)
		// Re-enable when backend fixes cyrillic 'п' in out.arпuments

		if strings.EqualFold(spec.Func, "getAvailableResourceRealms") {
			if !realmLoaded {
				loaded, loadErr := client.ListAvailableResourceRealms(ctx, serviceID)
				if loadErr != nil {
					diags.AddWarning(
						"Проверка realm пропущена",
						"API /resourceRealms/available вернул ошибку — список доступных realm получить не удалось. Проверка параметра "+code+" не выполнялась.",
					)
					realmLoaded = true
					realmCache = nil
				} else {
					realmCache = loaded
					realmLoaded = true
				}
			}
			if realmLoaded && realmCache != nil && len(realmCache) > 0 && !containsValueInsensitive(realmCache, value) {
				diags.AddError(
					"Некорректное значение параметра realm",
					"Параметр "+code+" содержит недопустимое значение: \""+value+"\".\nДоступные realm: "+joinLimited(realmCache, 25),
				)
			}
		}

		if spec.RefSvcId > 0 {
			options, cached := refOptionsCache[spec.RefSvcId]
			if !cached {
				loaded, loadErr := client.ListRefServiceInstances(ctx, spec.RefSvcId)
				if loadErr != nil {
					diags.AddWarning(
						"Проверка ref-параметра пропущена",
						"API /instances вернул ошибку для service_id="+fmt.Sprintf("%d", spec.RefSvcId)+" — список инстансов получить не удалось. Проверка значения \""+value+"\" не выполнялась.",
					)
					refOptionsCache[spec.RefSvcId] = nil
					continue
				}
				options = loaded
				refOptionsCache[spec.RefSvcId] = options
			}

			if !isValidRefValue(options, spec.RefSvcId, value) {
				diags.AddError(
					"Некорректное значение ссылочного параметра",
					"Параметр "+code+" содержит недопустимое значение: \""+value+"\".\nВыберите значение из списка доступных инстансов service_id="+fmt.Sprintf("%d", spec.RefSvcId)+": "+formatRefOptions(options, spec.RefSvcId, 20),
				)
			}
		}
	}
}

// CreateExistingResourceDiagnosticsWithDomain — устаревшая обёртка (suspend считается
// поддерживаемым). См. CreateExistingResourceDiagnosticsWithDomainAndServices.
func CreateExistingResourceDiagnosticsWithDomain(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool, desiredDomain string) diag.Diagnostics {
	return CreateExistingResourceDiagnosticsWithDomainAndServices(ctx, client, serviceID, resourceName, adoptExistingOnCreate, desiredDomain, nil, true)
}

// CreateExistingResourceDiagnosticsWithDomainAndServices — см.
// PlanExistingResourceDiagnosticsWithParamsAndDomainAndServices (supportsSuspend).
func CreateExistingResourceDiagnosticsWithDomainAndServices(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, adoptExistingOnCreate bool, desiredDomain string, allowedDomainServiceIDs []int, supportsSuspend bool) diag.Diagnostics {
	var diags diag.Diagnostics
	if client == nil {
		return diags
	}
	if strings.TrimSpace(resourceName) == "" {
		return diags
	}

	existing, err := client.FindInstanceByDisplayName(ctx, serviceID, resourceName)
	if err != nil || existing == nil {
		appendDomainCollisionDiagnostics(ctx, client, serviceID, resourceName, desiredDomain, allowedDomainServiceIDs, nil, &diags)
		return diags
	}

	appendDomainCollisionDiagnostics(ctx, client, serviceID, resourceName, desiredDomain, allowedDomainServiceIDs, existing, &diags)
	if diags.HasError() {
		return diags
	}

	status := normalizeStatus(existing.ExplainedStatus)
	if isStatusNotCreated(status) {
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (NOT CREATED)",
			"Инстанс найден, но находится в состоянии Not Created.\nУдалите его в Личном кабинете (Управление облаком), затем повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if isStatusRunning(status) {
		if adoptExistingOnCreate {
			return diags
		}
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (RUNNING)",
			"Инстанс с таким resource_name уже работает.\nАвто-усыновление отключено (adopt_existing_on_create=false).\nВыберите другое имя, если нужен новый ресурс."+runningConflictHint(supportsSuspend)+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}
	if isStatusSuspended(status) {
		if adoptExistingOnCreate {
			return diags
		}
		diags.AddError(
			"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (SUSPEND)",
			"Инстанс найден в состоянии Suspend, но adopt_existing_on_create=false."+suspendConflictHint(supportsSuspend)+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return diags
	}

	diags.AddError(
		"РЕСУРС С ТАКИМ ИМЕНЕМ УЖЕ СУЩЕСТВУЕТ (НЕПОДДЕРЖИВАЕМЫЙ СТАТУС)",
		"Инстанс найден, но статус не подходит для авто-усыновления.\nСтатус: "+status+". Выберите другое имя или устраните проблему в ЛК."+formatInstanceDetails(existing, serviceID, resourceName),
	)
	return diags
}

// runningConflictHint — что предлагать, если инстанс с таким именем уже работает.
//
// Для сервисов БЕЗ suspend/resume усыновление работающего инстанса технически возможно
// (приостановка не требуется — adopt просто берёт существующий uid), но пользователю
// часто нужен именно destroy, поэтому называем оба пути явно.
func runningConflictHint(supportsSuspend bool) string {
	if !supportsSuspend {
		return "\nЛибо усыновите существующий инстанс: adopt_existing_on_create=true (он уже работает, suspend не требуется).\nЛибо удалите инстанс в Личном кабинете Nubes и повторите apply."
	}
	return "\nИли выполните импорт существующего ресурса через усыновление: установите adopt_existing_on_create=true и повторите apply."
}

// suspendConflictHint — что предлагать, если инстанс найден в состоянии Suspend.
//
// Для сервисов БЕЗ suspend/resume состояние Suspend недостижимо, а adopt требует
// resume — значит совет «включите adopt_existing_on_create=true» заведомо невыполним.
// Корректный выход один: удалить инстанс и/или взять другое имя.
func suspendConflictHint(supportsSuspend bool) string {
	if !supportsSuspend {
		return "\nСервис НЕ поддерживает suspend/resume, поэтому усыновление (adopt) для него невозможно.\nУдалите инстанс в Личном кабинете Nubes и повторите apply, либо выберите другое resource_name."
	}
	return "\nДля восстановления/импорта этим ресурсом включите adopt_existing_on_create=true и повторите apply."
}

func appendLocalPlanInputDiagnostics(serviceID int, desiredParams map[int]string, desiredDomain string, diags *diag.Diagnostics) {
	if diags == nil {
		return
	}

	if domain := strings.TrimSpace(desiredDomain); domain != "" {
		if err := validateDomainValue(domain); err != nil {
			diags.AddError("Некорректный domain", err.Error())
		}
	}

	if len(desiredParams) == 0 || serviceID <= 0 {
		return
	}

	codeMap, err := getServiceParamCodeMap(serviceID)
	if err != nil || len(codeMap) == 0 {
		return
	}

	for id, rawValue := range desiredParams {
		value := strings.TrimSpace(rawValue)
		code := strings.TrimSpace(codeMap[id])
		if code == "" {
			continue
		}
		switch strings.ToLower(code) {
		case "jsonenv", "jsonparameters":
			if value == "" {
				continue
			}
			if !json.Valid([]byte(value)) {
				diags.AddError("Некорректный JSON", "Параметр "+code+" должен быть валидным JSON.")
			}
		case "domain":
			if value == "" {
				continue
			}
			if err := validateDomainValue(value); err != nil {
				diags.AddError("Некорректный domain", "Параметр "+code+": "+err.Error())
			}
		case "healthpath":
			if value == "" {
				continue
			}
			if !strings.HasPrefix(value, "/") {
				diags.AddError("Некорректный health_path", "Параметр "+code+" должен начинаться с '/'.")
			}
		case "gitpath":
			if value == "" {
				continue
			}
			u, parseErr := url.Parse(value)
			if parseErr != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") || strings.TrimSpace(u.Host) == "" {
				diags.AddError("Некорректный git_path", "Параметр "+code+" должен быть корректным http/https URL.")
			}
		case "resourcecpu", "resourcememory", "resourceinstances", "resourcedisk":
			if value == "" {
				continue
			}
			num, convErr := strconv.ParseInt(value, 10, 64)
			if convErr != nil {
				diags.AddError("Некорректное числовое значение", "Параметр "+code+" должен быть целым числом.")
				continue
			}
			if num <= 0 {
				diags.AddError("Некорректный диапазон", "Параметр "+code+" должен быть больше 0.")
			}
		}
	}
}

func appendAdoptOrphanGuardDiagnostics(ctx context.Context, client *core.UniversalClient, existing *core.InstanceStateResponse, serviceID int, resourceName string, diags *diag.Diagnostics) {
	if diags == nil || client == nil || existing == nil {
		return
	}
	instanceUID := strings.TrimSpace(existing.InstanceUid)
	if instanceUID == "" {
		diags.AddError(
			"Невозможно выполнить adopt_existing_on_create",
			"Найден инстанс без instance_uid. Усыновление невозможно, удалите объект в Личном кабинете и повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return
	}
	params, err := client.GetInstanceStateParams(ctx, instanceUID)
	if err != nil {
		diags.AddError(
			"Невозможно выполнить adopt_existing_on_create",
			"Найден существующий инстанс, но не удалось прочитать его state-параметры.\nСчитаем ресурс orphan/broken для adopt; удалите объект в Личном кабинете и повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
		return
	}
	if len(params) == 0 {
		diags.AddError(
			"Невозможно выполнить adopt_existing_on_create",
			"Найден существующий инстанс без state-параметров (orphan).\nУсыновление заблокировано: удалите инстанс в Личном кабинете и повторите apply."+formatInstanceDetails(existing, serviceID, resourceName),
		)
	}
}

var domainLabelRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

func validateDomainValue(domain string) error {
	trimmed := strings.TrimSpace(domain)
	if trimmed == "" {
		return nil
	}
	if trimmed != strings.ToLower(trimmed) {
		return fmt.Errorf("domain должен быть в lowercase")
	}
	if strings.ContainsAny(trimmed, " \t\n\r") {
		return fmt.Errorf("domain не должен содержать пробелы")
	}
	labels := strings.Split(trimmed, ".")
	for _, label := range labels {
		if !domainLabelRegex.MatchString(label) {
			return fmt.Errorf("domain должен быть валидным label/FQDN")
		}
	}
	return nil
}

func appendDomainCollisionDiagnostics(ctx context.Context, client *core.UniversalClient, serviceID int, resourceName string, desiredDomain string, allowedDomainServiceIDs []int, existingByName *core.InstanceStateResponse, diags *diag.Diagnostics) {
	if diags == nil {
		return
	}
	if client == nil {
		return
	}
	targetDomain := strings.TrimSpace(desiredDomain)
	if targetDomain == "" {
		return
	}

	owner, err := findDomainOwnerInstance(ctx, client, targetDomain, allowedDomainServiceIDs)
	if err != nil || owner == nil {
		return
	}

	if existingByName != nil && strings.TrimSpace(existingByName.InstanceUid) != "" && strings.EqualFold(existingByName.InstanceUid, owner.InstanceUID) {
		return
	}
	if strings.TrimSpace(resourceName) != "" && strings.EqualFold(strings.TrimSpace(resourceName), strings.TrimSpace(owner.DisplayName)) {
		return
	}

	diags.AddError(
		"ДОМЕН УЖЕ ЗАНЯТ ДРУГИМ РЕСУРСОМ",
		"Запрошенный domain уже используется другим инстансом.\nВыберите уникальный domain.\nИли, если это тот же самый объект, используйте корректный resource_name и adopt_existing_on_create=true для усыновления.\nDetails:"+
			"\n- requested_domain: "+targetDomain+
			"\n- requested_service_id: "+fmt.Sprintf("%d", serviceID)+
			"\n- requested_resource_name: "+resourceName+
			"\n- owner_instance_uid: "+owner.InstanceUID+
			"\n- owner_service_id: "+fmt.Sprintf("%d", owner.ServiceID)+
			"\n- owner_resource_name: "+owner.DisplayName,
	)
}

func containsValueCaseSensitive(values []string, target string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == target {
			return true
		}
	}
	return false
}

func containsValueInsensitive(values []string, target string) bool {
	for _, v := range values {
		if strings.EqualFold(strings.TrimSpace(v), target) {
			return true
		}
	}
	return false
}

func joinLimited(values []string, limit int) string {
	if len(values) == 0 {
		return "(пусто)"
	}
	if limit <= 0 || len(values) <= limit {
		return strings.Join(values, ", ")
	}
	return strings.Join(values[:limit], ", ") + ", ..."
}

func sortedCopy(values []string) []string {
	copyValues := make([]string, 0, len(values))
	copyValues = append(copyValues, values...)
	sort.Strings(copyValues)
	return copyValues
}

func isValidRefValue(options []core.RefServiceInstanceOption, refSvcId int, value string) bool {
	if len(options) == 0 {
		return true
	}
	for _, opt := range options {
		if strings.EqualFold(strings.TrimSpace(opt.InstanceUID), value) {
			return true
		}
		if strings.EqualFold(strings.TrimSpace(opt.DisplayName), value) {
			return true
		}
	}
	return false
}

func formatRefOptions(options []core.RefServiceInstanceOption, refSvcId int, limit int) string {
	if len(options) == 0 {
		return "список пуст"
	}
	formatted := make([]string, 0, len(options))
	for _, opt := range options {
		uid := strings.TrimSpace(opt.InstanceUID)
		name := strings.TrimSpace(opt.DisplayName)
		if name != "" && uid != "" {
			formatted = append(formatted, name+" ("+uid+")")
		} else if uid != "" {
			formatted = append(formatted, uid)
		} else if name != "" {
			formatted = append(formatted, name)
		}
	}
	if len(formatted) == 0 {
		return "список пуст"
	}
	return joinLimited(formatted, limit)
}

func shouldSurfaceAuxValidationWarning(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	if strings.Contains(msg, "ошибка api 500") || strings.Contains(msg, "ошибка api 502") || strings.Contains(msg, "ошибка api 503") || strings.Contains(msg, "ошибка api 504") {
		return false
	}
	return true
}
