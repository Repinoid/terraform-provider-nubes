package resources_core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &OrgIpAllocationResource{}
var _ resource.ResourceWithConfigure = &OrgIpAllocationResource{}
var _ resource.ResourceWithImportState = &OrgIpAllocationResource{}

// ===== Именованные ссылки на сервисы платформы (реестр данных, без магических чисел) =====
// Проверяется стражем TOOLS/scripts/check_hardcoded_service_ids.sh.
const (
	// svcIDVcOrg — сервис 19 «vc_org»: организация Cloud Director.
	// Организация создаётся вручную в ЛК, ресурс адресует её по имени/UUID.
	svcIDVcOrg = 19
)

// OrgIpAllocationResource управляет аллокацией внешних IP на СУЩЕСТВУЮЩЕЙ организации
// (сервис 19, vc_org) через операцию modify с параметром vIPConfigure (id 662).
//
// Организация НЕ управляется Terraform: она создаётся один раз вручную в ЛК
// и адресуется здесь по uid.
//
// Семантика операции — replace всего массива: переданное значение полностью заменяет
// текущую аллокацию (проверено тестом NOTES/30_analysis/ORG_IP_MODIFIER_TEST_2026-09-22.md).
// Поэтому ресурс владеет массивом ЦЕЛИКОМ, а не отдельным элементом.
type OrgIpAllocationResource struct {
	client *core.UniversalClient
}

type OrgIpAllocationModel struct {
	ID            types.String `tfsdk:"id"`
	Organization  types.String `tfsdk:"organization"`
	VIPConfigure  types.String `tfsdk:"vip_configure"`
	KeepOnDestroy types.Bool   `tfsdk:"keep_on_destroy"`
}

// vipAllocation — элемент массива vIPConfigure. count ВСЕГДА строка:
// ЛК присылает его строкой (HAR/globak.har), API принимает строкой.
type vipAllocation struct {
	Name  string
	Count string
}

func NewOrgIpAllocationResource() resource.Resource {
	return &OrgIpAllocationResource{}
}

func (r *OrgIpAllocationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vc_org_ip_allocation"
}

func (r *OrgIpAllocationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Аллокация внешних IP (vIPConfigure) на существующей организации Cloud Director. " +
			"Организация создаётся вручную в ЛК, в конфиге указывается её имя или UUID. " +
			"Операция имеет replace-семантику: массив перезаписывается целиком.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Организация, на которой выделяются внешние IP: имя из ЛК (например `organ`) " +
					"или её UUID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vip_configure": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "JSON-массив аллокаций: `[{\"name\":\"internet-ipv4-v1\",\"count\":\"3\"}]`. " +
					"Значение перезаписывает текущую аллокацию целиком. `count` — строка. " +
					"Порядок ключей и форматирование не важны (сравнение смысловое). " +
					"Снять аллокацию (`[]`) через этот атрибут **нельзя** — только удалением ресурса (`destroy`).",
			},
			"keep_on_destroy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Не снимать аллокацию IP при `destroy` (по умолчанию `false` — квота обнуляется, " +
					"`count=0` по каждому элементу).",
			},
		},
	}
}

func (r *OrgIpAllocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrgIpAllocationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUID, err := r.resolveOrganizationUID(ctx, plan.Organization)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	if err := r.applyAllocation(ctx, orgUID, plan.VIPConfigure); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(orgUID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgIpAllocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrgIpAllocationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgUID, err := r.resolveOrganizationUID(ctx, plan.Organization)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	if err := r.applyAllocation(ctx, orgUID, plan.VIPConfigure); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(orgUID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrgIpAllocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrgIpAllocationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.TrimSpace(state.Organization.ValueString()) == "" || r.client == nil {
		return
	}

	orgUID, err := r.resolveOrganizationUID(ctx, state.Organization)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	remove, err := ShouldRemoveFromState(ctx, r.client, orgUID)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}
	if remove {
		// Организации больше нет — ресурс тоже не нужен.
		resp.State.RemoveResource(ctx)
		return
	}

	live, err := r.client.GetInstanceStateParams(ctx, orgUID)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	// Атрибут принадлежит пользователю: НЕ переписываем его, если смысл совпал — иначе Terraform
	// увидит расхождение config vs state и покажет ложный дрейф (jsonencode отдаёт ключи по алфавиту).
	// Писать null в Required-атрибут тоже нельзя (это даёт "Provider produced inconsistent result").
	raw, ok := live["vIPConfigure"]
	if ok {
		liveItems, parseErr := parseVipConfigure(raw)
		if parseErr != nil {
			resp.Diagnostics.AddError("Ошибка чтения состояния", parseErr.Error())
			return
		}
		stateItems, _ := parseVipConfigure(state.VIPConfigure.ValueString())
		if !vipAllocationsEqual(liveItems, stateItems) {
			state.VIPConfigure = types.StringValue(formatVipConfigure(liveItems))
		}
	}

	state.ID = types.StringValue(orgUID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrgIpAllocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrgIpAllocationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if strings.TrimSpace(state.Organization.ValueString()) == "" || r.client == nil {
		return
	}

	orgUID, err := r.resolveOrganizationUID(ctx, state.Organization)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	if !state.KeepOnDestroy.IsNull() && !state.KeepOnDestroy.IsUnknown() && state.KeepOnDestroy.ValueBool() {
		resp.Diagnostics.AddWarning(
			"Аллокация IP не снималась",
			fmt.Sprintf("keep_on_destroy = true: квота внешних IP организации %s оставлена без изменений.", orgUID),
		)
		return
	}

	remove, err := ShouldRemoveFromState(ctx, r.client, orgUID)
	if err != nil {
		// Реальная ошибка API (не «инстанса нет») — нельзя молча терять квоту: ресурс уйдёт из state,
		// а выделенные IP останутся висеть.
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}
	if remove {
		resp.Diagnostics.AddWarning(
			"Аллокация IP не снималась",
			fmt.Sprintf("организация %s не найдена — обратный modify пропущен.", orgUID),
		)
		return
	}

	unlock := r.client.LockInstance(orgUID)
	defer unlock()

	// Имена берём из LIVE-состояния (что реально выделено), при неудаче — из конфигурации.
	items := []vipAllocation{}
	if live, liveErr := r.client.GetInstanceStateParams(ctx, orgUID); liveErr == nil {
		if parsed, parseErr := parseVipConfigure(live["vIPConfigure"]); parseErr == nil {
			items = parsed
		}
	}
	if len(items) == 0 {
		if parsed, parseErr := parseVipConfigure(state.VIPConfigure.ValueString()); parseErr == nil {
			items = parsed
		}
	}
	if len(items) == 0 {
		resp.Diagnostics.AddWarning(
			"Аллокация IP не снималась",
			"не удалось определить выделенные ipSpace — обратный modify пропущен.",
		)
		return
	}

	// Обратный modify: тот же массив, но count=0 (форма проверена тестом 09-22).
	// Пустой массив `[]` НЕ отправляем — его семантика на платформе не проверена.
	zero := make([]vipAllocation, 0, len(items))
	for _, item := range items {
		zero = append(zero, vipAllocation{Name: item.Name, Count: "0"})
	}

	if err := r.client.RunInstanceOperationUniversalByCode(ctx, orgUID, "modify", map[string]string{
		"vIPConfigure": formatVipConfigure(zero),
	}); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	resp.Diagnostics.AddWarning(
		"Квота IP обнулена",
		fmt.Sprintf("по организации %s отправлен modify с count=0: %s", orgUID, formatVipConfigure(zero)),
	)
}

func (r *OrgIpAllocationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*core.UniversalClient)
	if !ok {
		resp.Diagnostics.AddError("Ошибка", "Неверный тип клиента, ожидается *core.UniversalClient")
		return
	}
	r.client = client
}

func (r *OrgIpAllocationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	uid := strings.TrimSpace(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), uid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), uid)...)

	// Required-атрибут vip_configure нельзя оставить null после импорта.
	// Заполняем live-значением; если аллокации нет — пустым массивом.
	var live map[string]string
	if r.client != nil {
		live, _ = r.client.GetInstanceStateParams(ctx, uid)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vip_configure"), importVipConfigure(live))...)
}

// importVipConfigure выбирает значение vip_configure для состояния после импорта:
// канонический JSON из live (если есть непустые элементы), иначе пустой массив.
func importVipConfigure(live map[string]string) string {
	if raw, ok := live["vIPConfigure"]; ok {
		if items, err := parseVipConfigure(raw); err == nil && len(items) > 0 {
			return formatVipConfigure(items)
		}
	}
	return "[]"
}

// resolveOrganizationUID принимает имя организации из ЛК или её UUID и возвращает UUID.
// Резолв делает клиент — тем же путём, что сгенерированный nubes_vc_vdc
// (core.ResolveRefSvcParamValue, сравн. 21_vc_vdc_resource.go).
func (r *OrgIpAllocationResource) resolveOrganizationUID(ctx context.Context, organization types.String) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("клиент не инициализирован")
	}
	raw := strings.TrimSpace(organization.ValueString())
	if raw == "" {
		return "", fmt.Errorf("organization обязателен")
	}

	resolved, err := r.client.ResolveRefSvcParamValue(ctx, svcIDVcOrg, raw)
	if err != nil {
		return "", fmt.Errorf("не удалось определить организацию %q: %w", raw, err)
	}
	resolved = strings.TrimSpace(resolved)
	if resolved == "" {
		return "", fmt.Errorf("организация %q не найдена", raw)
	}
	return resolved, nil
}

// applyAllocation отправляет modify с массивом vIPConfigure целиком.
func (r *OrgIpAllocationResource) applyAllocation(ctx context.Context, orgUID string, vipConfigure types.String) error {
	uid := strings.TrimSpace(orgUID)
	if uid == "" {
		return fmt.Errorf("organization обязателен")
	}
	if r.client == nil {
		return fmt.Errorf("клиент не инициализирован")
	}

	items, err := parseVipConfigure(vipConfigure.ValueString())
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("vip_configure не содержит ни одной аллокации (name+count)")
	}

	unlock := r.client.LockInstance(uid)
	defer unlock()

	// Именно ByCode (без idempotency-pre-check): pre-check сравнивает с paramValue ФОРМЫ
	// операции, а это не live-состояние инстанса (см. core/modifier_compare.go и
	// комментарий в core/operation_cfs.go) — можно было бы ложно пропустить modify.
	return r.client.RunInstanceOperationUniversalByCode(ctx, uid, "modify", map[string]string{
		"vIPConfigure": formatVipConfigure(items),
	})
}

// parseVipConfigure разбирает значение параметра vIPConfigure.
// Пустые элементы (`{}`) — легальное состояние «не выделено» у свежей орги
// (NOTES/30_analysis/HAR_FRESH_CREATE_2026-09-24.md) и отбрасываются.
func parseVipConfigure(raw string) ([]vipAllocation, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &items); err != nil {
		return nil, fmt.Errorf("не удалось разобрать vIPConfigure %q: %w", trimmed, err)
	}

	out := make([]vipAllocation, 0, len(items))
	for _, item := range items {
		name := ""
		if v, ok := item["name"]; ok && v != nil {
			name = strings.TrimSpace(fmt.Sprint(v))
		}
		if name == "" {
			continue
		}
		count := "0"
		if v, ok := item["count"]; ok && v != nil {
			if parsed := strings.TrimSpace(fmt.Sprint(v)); parsed != "" {
				count = parsed
			}
		}
		out = append(out, vipAllocation{Name: name, Count: count})
	}
	return out, nil
}

// formatVipConfigure собирает канонический payload: [{"name":"…","count":"…"}]
// (порядок ключей name,count; count — строка). Канон ЕДИНЫЙ для отправки и для Read,
// иначе план и state расходятся по строке — см. vipConfigureCanonical.
func formatVipConfigure(items []vipAllocation) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, fmt.Sprintf(`{"name":%q,"count":%q}`, item.Name, item.Count))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// vipAllocationsEqual сравнивает аллокации по СМЫСЛУ: порядок элементов и формат не важны.
// Имена ipSpace в рамках организации уникальны, поэтому сравнение идёт по имени.
func vipAllocationsEqual(a, b []vipAllocation) bool {
	if len(a) != len(b) {
		return false
	}
	byName := make(map[string]string, len(b))
	for _, item := range b {
		byName[item.Name] = item.Count
	}
	for _, item := range a {
		count, ok := byName[item.Name]
		if !ok || count != item.Count {
			return false
		}
	}
	return true
}
