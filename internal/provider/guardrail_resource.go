package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &GuardrailResource{}
	_ resource.ResourceWithImportState = &GuardrailResource{}
)

func NewGuardrailResource() resource.Resource {
	return &GuardrailResource{}
}

// GuardrailResource defines the resource implementation.
type GuardrailResource struct {
	service *client.GuardrailsService
}

func (r *GuardrailResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_guardrail"
}

func (r *GuardrailResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A resourcely guardrail",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the guardrail",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the guardrail",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name to associate with the guardrail",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of what the guardrail applies to",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider that this guardrail targets. Can be one of `PROVIDER_AMAZON`, `PROVIDER_AZURE`, `PROVIDER_CONDUCTORONE`, `PROVIDER_DATADOG`, `PROVIDER_GITHUB`, `PROVIDER_GITLAB`, `PROVIDER_GOOGLE`, `PROVIDER_JUMPCLOUD`, `PROVIDER_OKTA`, `PROVIDER_RESOURCELY`, `PROVIDER_SNOWFLAKE`, `PROVIDER_SPACELIFT`, `PROVIDER_OTHER`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PROVIDER_AMAZON",
						"PROVIDER_AZURE",
						"PROVIDER_CONDUCTORONE",
						"PROVIDER_DATADOG",
						"PROVIDER_GITHUB",
						"PROVIDER_GITLAB",
						"PROVIDER_GOOGLE",
						"PROVIDER_JUMPCLOUD",
						"PROVIDER_OKTA",
						"PROVIDER_RESOURCELY",
						"PROVIDER_SNOWFLAKE",
						"PROVIDER_SPACELIFT",
						"PROVIDER_OTHER",
					),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category to assign to the guardrail. Can be one of `GUARDRAIL_ACCESS_CONTROL`, `GUARDRAIL_BEST_PRACTICES`, `GUARDRAIL_CHANGE_MANAGEMENT`, `GUARDRAIL_COST_EFFICIENCY`, `GUARDRAIL_ENCRYPTION`, `GUARDRAIL_GLOBALIZATION`, `GUARDRAIL_IAM`, `GUARDRAIL_LOGGING`, `GUARDRAIL_MODULE_INPUTS`, `GUARDRAIL_PRIVACY_COMPLIANCE`, `GUARDRAIL_RELIABILITY`, `GUARDRAIL_STORAGE_AND_SCALE`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"GUARDRAIL_ACCESS_CONTROL",
						"GUARDRAIL_BEST_PRACTICES",
						"GUARDRAIL_CHANGE_MANAGEMENT",
						"GUARDRAIL_COST_EFFICIENCY",
						"GUARDRAIL_ENCRYPTION",
						"GUARDRAIL_GLOBALIZATION",
						"GUARDRAIL_IAM",
						"GUARDRAIL_LOGGING",
						"GUARDRAIL_MODULE_INPUTS",
						"GUARDRAIL_PRIVACY_COMPLIANCE",
						"GUARDRAIL_RELIABILITY",
						"GUARDRAIL_STORAGE_AND_SCALE",
					),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state to set the guardrail to. Can be one of `GUARDRAIL_STATE_INACTIVE`, `GUARDRAIL_STATE_EVALUATE_ONLY`, `GUARDRAIL_STATE_ACTIVE`. If not provided state is set to `GUARDRAIL_STATE_ACTIVE`",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("GUARDRAIL_STATE_ACTIVE"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"GUARDRAIL_STATE_INACTIVE",
						"GUARDRAIL_STATE_EVALUATE_ONLY",
						"GUARDRAIL_STATE_ACTIVE",
					),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The content of the guardrail, written in Resourcely's Really language",
				Required:            true,
			},
		},
	}
}

func (r *GuardrailResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *client.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.service = client.Guardrails
}

func (r *GuardrailResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Get the plan
	var plan GuardrailResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	newGuardrail := &client.NewGuardrail{
		CommonGuardrailFields: client.CommonGuardrailFields{
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueString(),
			Provider:    plan.Provider.ValueString(),
			Category:    plan.Category.ValueString(),
			State:       plan.State.ValueString(),
			Content:     plan.Content.ValueString(),
		},
	}

	guardrail, _, err := r.service.CreateGuardrail(ctx, newGuardrail)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating guardrail",
			"Could not create guardrail, unexpected error: "+err.Error(),
		)
		return
	}

	// Set the resource state
	state := FlattenGuardrail(guardrail)
	resp.Diagnostics.Append(NormalizeGuardrail(&state, &plan)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GuardrailResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get the current state
	var previousState GuardrailResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &previousState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh value from the remote API
	guardrail, httpResp, err := r.service.GetGuardrailBySeriesId(
		ctx,
		previousState.SeriesId.ValueString(),
	)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Guardrail "+previousState.SeriesId.ValueString()+" was not found in Resourcely",
				"The guardrail may have been deleted outside of Terraform",
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error reading guardrail",
				"Could not read guardrail series id "+previousState.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite state with refreshed value
	state := FlattenGuardrail(guardrail)
	resp.Diagnostics.Append(NormalizeGuardrail(&state, &previousState)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GuardrailResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve the plan and state
	var plan GuardrailResourceModel
	var state GuardrailResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	updatedGuardrail := &client.UpdatedGuardrail{
		SeriesId: state.SeriesId.ValueString(),
		CommonGuardrailFields: client.CommonGuardrailFields{
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueString(),
			Provider:    plan.Provider.ValueString(),
			Category:    plan.Category.ValueString(),
			State:       plan.State.ValueString(),
			Content:     plan.Content.ValueString(),
		},
	}

	guardrail, _, err := r.service.UpdateGuardrail(ctx, updatedGuardrail)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating guardrail",
			"Could not update guardrail series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the resource state
	state = FlattenGuardrail(guardrail)
	resp.Diagnostics.Append(NormalizeGuardrail(&state, &plan)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GuardrailResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve from state
	var state GuardrailResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.service.DeleteGuardrail(ctx, state.SeriesId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting guardrail",
			"Could not delete guardrail series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *GuardrailResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("series_id"), req, resp)
}
