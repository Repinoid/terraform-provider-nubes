package resources_core

import (
	"context"
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

var _ resource.Resource = &NsxtSnatResource{}
var _ resource.ResourceWithConfigure = &NsxtSnatResource{}
var _ resource.ResourceWithImportState = &NsxtSnatResource{}

// NsxtSnatResource включает/выключает SNAT у СУЩЕСТВУЮЩЕГО сетевого шлюза периметра
// (сервис 22, vc_nsxt) через операцию modify с параметром ipSpaceName (id 372).
//
// Зачем отдельный ресурс: ipSpaceName есть ТОЛЬКО в операции modify (в create его нет),
// поэтому одним ресурсом «create + modify» в одном apply не сделать.
//
// Канонические значения (HAR/edge_.har, NOTES/30_analysis/HAR_SNAT_MODIFY_FINDINGS.md):
//   - включить SNAT: ip_space_name = "<имя ipSpace из аллокации организации>";
//   - выключить SNAT: ip_space_name = "no-needed" (легальное значение платформы).
type NsxtSnatResource struct {
	client *core.UniversalClient
}

type NsxtSnatModel struct {
	ID            types.String `tfsdk:"id"`
	NsxtUID       types.String `tfsdk:"nsxt_uid"`
	IpSpaceName   types.String `tfsdk:"ip_space_name"`
	KeepOnDestroy types.Bool   `tfsdk:"keep_on_destroy"`
}

// noNeededIpSpace — каноническое значение «SNAT не нужен».
const noNeededIpSpace = "no-needed"

func NewNsxtSnatResource() resource.Resource {
	return &NsxtSnatResource{}
}

func (r *NsxtSnatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vc_nsxt_snat"
}

func (r *NsxtSnatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "SNAT (ipSpaceName) на существующем сетевом шлюзе периметра. " +
			"Шлюз создаётся отдельным ресурсом `nubes_vc_nsxt`, здесь задаётся только SNAT. " +
			"Значение `no-needed` выключает SNAT.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nsxt_uid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "UUID существующей услуги «Сетевой шлюз периметра (Edge)».",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_space_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Имя ipSpace для внешнего IP (SNAT). Значение `no-needed` выключает SNAT. " +
					"Имя должно быть выделено на организации (см. `nubes_vc_org_ip_allocation`).",
			},
			"keep_on_destroy": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "Не выключать SNAT при `destroy` (по умолчанию `false` — отправляется " +
					"`ipSpaceName = \"no-needed\"`).",
			},
		},
	}
}

func (r *NsxtSnatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NsxtSnatModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.setSnat(ctx, plan.NsxtUID, plan.IpSpaceName); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(strings.TrimSpace(plan.NsxtUID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NsxtSnatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NsxtSnatModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.setSnat(ctx, plan.NsxtUID, plan.IpSpaceName); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(strings.TrimSpace(plan.NsxtUID.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NsxtSnatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NsxtSnatModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nsxtUID := strings.TrimSpace(state.NsxtUID.ValueString())
	if nsxtUID == "" || r.client == nil {
		return
	}

	remove, err := ShouldRemoveFromState(ctx, r.client, nsxtUID)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}
	if remove {
		resp.State.RemoveResource(ctx)
		return
	}

	live, err := r.client.GetInstanceStateParams(ctx, nsxtUID)
	if err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	// ВАЖНО: в Required-атрибут нельзя писать null — после apply state обязан совпасть с планом,
	// иначе Terraform вернёт "Provider produced inconsistent result after apply". Если ключа ещё нет
	// (SNAT ни разу не включали, HAR fresh-create) — оставляем текущее значение state.
	if raw, ok := live["ipSpaceName"]; ok && strings.TrimSpace(raw) != "" {
		state.IpSpaceName = types.StringValue(strings.TrimSpace(raw))
	}

	state.ID = types.StringValue(nsxtUID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NsxtSnatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NsxtSnatModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nsxtUID := strings.TrimSpace(state.NsxtUID.ValueString())
	if nsxtUID == "" || r.client == nil {
		return
	}

	if !state.KeepOnDestroy.IsNull() && !state.KeepOnDestroy.IsUnknown() && state.KeepOnDestroy.ValueBool() {
		resp.Diagnostics.AddWarning(
			"SNAT не выключался",
			fmt.Sprintf("keep_on_destroy = true: ipSpaceName шлюза %s оставлен без изменений.", nsxtUID),
		)
		return
	}

	remove, err := ShouldRemoveFromState(ctx, r.client, nsxtUID)
	if err != nil {
		// Реальная ошибка API (не «шлюза нет») — нельзя молча терять SNAT: ресурс уйдёт из state,
		// а SNAT останется включённым.
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}
	if remove {
		resp.Diagnostics.AddWarning(
			"SNAT не выключался",
			fmt.Sprintf("шлюз %s не найден — обратный modify пропущен.", nsxtUID),
		)
		return
	}

	unlock := r.client.LockInstance(nsxtUID)
	defer unlock()

	// Обратный modify: каноническое «SNAT выключен» = no-needed (подтверждено HAR).
	if err := r.client.RunInstanceOperationUniversalByCode(ctx, nsxtUID, "modify", map[string]string{
		"ipSpaceName": noNeededIpSpace,
	}); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	resp.Diagnostics.AddWarning(
		"SNAT выключен",
		fmt.Sprintf("по шлюзу %s отправлен modify с ipSpaceName = %q.", nsxtUID, noNeededIpSpace),
	)
}

func (r *NsxtSnatResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NsxtSnatResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	uid := strings.TrimSpace(req.ID)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), uid)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("nsxt_uid"), uid)...)

	// Required-атрибут ip_space_name нельзя оставить null после импорта: пользователь
	// обязан указать его в конфиге. Заполняем live-значением; если SNAT не выставлен —
	// каноном «выключено» (no-needed).
	var live map[string]string
	if r.client != nil {
		live, _ = r.client.GetInstanceStateParams(ctx, uid)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_space_name"), importIpSpaceName(live))...)
}

// importIpSpaceName выбирает значение ip_space_name для состояния после импорта:
// live-значение, если есть; иначе канон «SNAT выключен» (no-needed).
func importIpSpaceName(live map[string]string) string {
	if v := strings.TrimSpace(live["ipSpaceName"]); v != "" {
		return v
	}
	return noNeededIpSpace
}

// setSnat отправляет modify только с ipSpaceName. Остальные параметры операции
// (needEnableAVI, virtualServicesCount, qosProfile, routedNetConfiguration) досылаются
// клиентом из LIVE-состояния инстанса — приоритет live → paramValue формы → default
// (core/operation_run_bycode.go), поэтому частичный payload ничего не затирает.
func (r *NsxtSnatResource) setSnat(ctx context.Context, nsxtUID types.String, ipSpaceName types.String) error {
	uid := strings.TrimSpace(nsxtUID.ValueString())
	if uid == "" {
		return fmt.Errorf("nsxt_uid обязателен")
	}
	if r.client == nil {
		return fmt.Errorf("клиент не инициализирован")
	}

	// Пустую строку молча подменять нельзя (скрытое поведение + риск вечного diff).
	// Выключение SNAT — явное каноническое значение "no-needed".
	value := strings.TrimSpace(ipSpaceName.ValueString())
	if value == "" {
		return fmt.Errorf("ip_space_name не может быть пустым: укажите имя ipSpace или %q для выключения SNAT", noNeededIpSpace)
	}

	unlock := r.client.LockInstance(uid)
	defer unlock()

	// Idempotent: core.RunInstanceOperationUniversalByIdempotent выполняет pre-check
	// по LIVE-значениям инстанса (state.params) и пропускает modify, если
	// ipSpaceName уже совпадает — повторный apply не дёргает платформу зря.
	return r.client.RunInstanceOperationUniversalByIdempotent(ctx, uid, "modify", map[string]string{
		"ipSpaceName": value,
	})
}
