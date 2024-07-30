package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Resourcely-Inc/terraform-provider-resourcely-internal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
				MarkdownDescription: "",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PROVIDER_AMAZON",
						"PROVIDER_GOOGLE",
						"PROVIDER_GITHUB",
						"PROVIDER_AZURE",
						"PROVIDER_JUMPCLOUD",
						"PROVIDER_GITLAB",
						"PROVIDER_OTHER",
						"PROVIDER_SNOWFLAKE",
						"PROVIDER_DATADOG",
						"PROVIDER_SPACELIFT",
						"PROVIDER_RESOURCELY",
						"PROVIDER_CONDUCTORONE",
					),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"GUARDRAIL_ACCESS_CONTROL",
						"GUARDRAIL_BEST_PRACTICES",
						"GUARDRAIL_COST_EFFICIENCY",
						"GUARDRAIL_GLOBALIZATION",
						"GUARDRAIL_PRIVACY_COMPLIANCE",
						"GUARDRAIL_STORAGE_AND_SCALE",
						"GUARDRAIL_MODULE_INPUTS",
						"GUARDRAIL_IAM",
						"GUARDRAIL_ENCRYPTION",
						"GUARDRAIL_LOGGING",
						"GUARDRAIL_CHANGE_MANAGEMENT",
						"GUARDRAIL_RELIABILITY"),
				},
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"rego_policy": schema.StringAttribute{
				MarkdownDescription: "String of the Rego Policy",
				Optional:            true,
				Computed:            true,
			},
			"cue_policy": schema.StringAttribute{
				MarkdownDescription: "String of the Cue Policy",
				Optional:            true,
				Computed:            true,
			},
			"json_validation": schema.StringAttribute{
				MarkdownDescription: "JSON string of the JSON Validation used on the frontend",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					StringIsJSON(),
				},
				PlanModifiers: []planmodifier.String{
					SuppressEquivalentJsonDiffs(),
				},
			},
			"guardrail_template_id": schema.StringAttribute{
				MarkdownDescription: "The id of the guardrail template used to render the policies",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guardrail_template_inputs": schema.StringAttribute{
				MarkdownDescription: "JSON string of the guardrail template inputs used on the frontend",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					StringIsJSON(),
				},
				PlanModifiers: []planmodifier.String{
					SuppressEquivalentJsonDiffs(),
				},
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

func (r *GuardrailResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("content"),
			path.MatchRoot("guardrail_template_id"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("guardrail_template_id"),
			path.MatchRoot("guardrail_template_inputs"),
		),
		resourcevalidator.RequiredTogether(
			path.MatchRoot("content"),
			path.MatchRoot("rego_policy"),
			path.MatchRoot("cue_policy"),
			path.MatchRoot("json_validation"),
		),
	}
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
			Name:                plan.Name.ValueString(),
			Description:         plan.Description.ValueString(),
			Provider:            plan.Provider.ValueString(),
			Category:            plan.Category.ValueString(),
			IsActive:            plan.IsActive.ValueBool(),
			Content:             plan.Content.ValueString(),
			RegoPolicy:          plan.RegoPolicy.ValueString(),
			CuePolicy:           plan.CuePolicy.ValueString(),
			GuardrailTemplateId: plan.GuardrailTemplateId.ValueString(),
		},
	}

	if plan.JsonValidation.ValueStringPointer() != nil && plan.JsonValidation.ValueString() != "" {
		newGuardrail.JsonValidation = json.RawMessage(plan.JsonValidation.ValueString())
	}

	if plan.GuardrailTemplateInputs.ValueStringPointer() != nil &&
		plan.GuardrailTemplateInputs.ValueString() != "" {
		newGuardrail.GuardrailTemplateInputs = json.RawMessage(
			plan.GuardrailTemplateInputs.ValueString(),
		)
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
			Name:                plan.Name.ValueString(),
			Description:         plan.Description.ValueString(),
			Provider:            plan.Provider.ValueString(),
			Category:            plan.Category.ValueString(),
			IsActive:            plan.IsActive.ValueBool(),
			Content:             plan.Content.ValueString(),
			RegoPolicy:          plan.RegoPolicy.ValueString(),
			CuePolicy:           plan.CuePolicy.ValueString(),
			GuardrailTemplateId: plan.GuardrailTemplateId.ValueString(),
		},
	}

	if plan.JsonValidation.ValueStringPointer() != nil && plan.JsonValidation.ValueString() != "" {
		updatedGuardrail.JsonValidation = json.RawMessage(plan.JsonValidation.ValueString())
	}

	if plan.GuardrailTemplateInputs.ValueStringPointer() != nil &&
		plan.GuardrailTemplateInputs.ValueString() != "" {
		updatedGuardrail.GuardrailTemplateInputs = json.RawMessage(
			plan.GuardrailTemplateInputs.ValueString(),
		)
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
