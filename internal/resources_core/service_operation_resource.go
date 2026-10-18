package resources_core

import (
	"context"
	"fmt"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ServiceOperationResource{}
var _ resource.ResourceWithConfigure = &ServiceOperationResource{}

// ServiceOperationResource triggers a one-time operation for an existing instance.
type ServiceOperationResource struct {
	client *core.UniversalClient
}

type ServiceOperationModel struct {
	ID          types.String `tfsdk:"id"`
	InstanceUid types.String `tfsdk:"instance_uid"`
	Operation   types.String `tfsdk:"operation"`
	RunID       types.String `tfsdk:"run_id"`
	Params      types.Map    `tfsdk:"params"`
}

func NewServiceOperationResource() resource.Resource {
	return &ServiceOperationResource{}
}

func (r *ServiceOperationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_operation"
}

func (r *ServiceOperationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"instance_uid": schema.StringAttribute{
				Required: true,
			},
			"operation": schema.StringAttribute{
				Required: true,
			},
			"run_id": schema.StringAttribute{
				Required: true,
			},
			"params": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *ServiceOperationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceOperationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceUID := strings.TrimSpace(plan.InstanceUid.ValueString())
	operation := strings.TrimSpace(plan.Operation.ValueString())
	runID := strings.TrimSpace(plan.RunID.ValueString())
	if instanceUID == "" || operation == "" || runID == "" {
		resp.Diagnostics.AddError("Ошибка клиента", "instance_uid, operation, and run_id are required")
		return
	}

	params, diags := StringMapFromTypesMap(ctx, plan.Params)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Ошибка клиента", "клиент не инициализирован")
		return
	}

	unlock := r.client.LockInstance(instanceUID)
	defer unlock()

	if err := r.client.RunInstanceOperationUniversalByCode(ctx, instanceUID, operation, params); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(buildServiceOperationID(instanceUID, operation, runID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceOperationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceOperationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceOperationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceOperationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceUID := strings.TrimSpace(plan.InstanceUid.ValueString())
	operation := strings.TrimSpace(plan.Operation.ValueString())
	runID := strings.TrimSpace(plan.RunID.ValueString())
	if instanceUID == "" || operation == "" || runID == "" {
		resp.Diagnostics.AddError("Ошибка клиента", "instance_uid, operation, and run_id are required")
		return
	}

	params, diags := StringMapFromTypesMap(ctx, plan.Params)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if r.client == nil {
		resp.Diagnostics.AddError("Ошибка клиента", "клиент не инициализирован")
		return
	}

	unlock := r.client.LockInstance(instanceUID)
	defer unlock()

	if err := r.client.RunInstanceOperationUniversalByCode(ctx, instanceUID, operation, params); err != nil {
		resp.Diagnostics.AddError("Ошибка клиента", err.Error())
		return
	}

	plan.ID = types.StringValue(buildServiceOperationID(instanceUID, operation, runID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceOperationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: removing the resource does not trigger a remote action.
}

func (r *ServiceOperationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func buildServiceOperationID(instanceUID string, operation string, runID string) string {
	return fmt.Sprintf("%s:%s:%s", strings.TrimSpace(instanceUID), strings.TrimSpace(operation), strings.TrimSpace(runID))
}
