package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &BlueprintResource{}
	_ resource.ResourceWithImportState = &BlueprintResource{}
)

func NewBlueprintResource() resource.Resource {
	return &BlueprintResource{}
}

// BlueprintResource defines the resource implementation.
type BlueprintResource struct {
	service *client.BlueprintsService
}

func (r *BlueprintResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (r *BlueprintResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Resourcely Blueprint",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the blueprint",
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
				MarkdownDescription: "The cloud provider that this blueprint targets. Can be one of `PROVIDER_AMAZON`, `PROVIDER_AZURE`, `PROVIDER_CONDUCTORONE`, `PROVIDER_DATADOG`, `PROVIDER_GITHUB`, `PROVIDER_GITLAB`, `PROVIDER_GOOGLE`, `PROVIDER_JUMPCLOUD`, `PROVIDER_RESOURCELY`, `PROVIDER_SNOWFLAKE`, `PROVIDER_SPACELIFT`, `PROVIDER_OTHER`",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
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
			"content": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"guidance": schema.StringAttribute{
				MarkdownDescription: "",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"labels": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				Computed:            true,
				Optional:            true,
			},
			"categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "The category to assign to this blueprint. Can be one of `BLUEPRINT_ASYNC_PROCESSING`, `BLUEPRINT_BLOB_STORAGE`, `BLUEPRINT_COMPUTE`, `BLUEPRINT_CONTAINERIZATION`, `BLUEPRINT_DATABASE`, `BLUEPRINT_GITHUB_REPO`, `BLUEPRINT_GITHUB_REPO_TEAM`, `BLUEPRINT_IAM`, `BLUEPRINT_LOGS_AND_METRICS`, `BLUEPRINT_NETWORKING`, `BLUEPRINT_SERVERLESS_COMPUTE`",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				Computed:            true,
				Optional:            true,
				Validators: []validator.Set{
					// All list items must pass the nested validators
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						"BLUEPRINT_ASYNC_PROCESSING",
						"BLUEPRINT_BLOB_STORAGE",
						"BLUEPRINT_COMPUTE",
						"BLUEPRINT_CONTAINERIZATION",
						"BLUEPRINT_DATABASE",
						"BLUEPRINT_GITHUB_REPO",
						"BLUEPRINT_GITHUB_REPO_TEAM",
						"BLUEPRINT_IAM",
						"BLUEPRINT_LOGS_AND_METRICS",
						"BLUEPRINT_NETWORKING",
						"BLUEPRINT_SERVERLESS_COMPUTE",
					),
					),
				},
			},
			"excluded_context_question_series": schema.SetAttribute{
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				ElementType:         basetypes.StringType{},
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "series_id for context questions that won't be used with this blueprint, even if this blueprint matches the context questions' blueprint_categories",
			},
		},
	}
}

func (r *BlueprintResource) Configure(
	_ context.Context,
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

	r.service = client.Blueprints
}

func (r *BlueprintResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Get the plan
	var plan BlueprintResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	newBlueprint := &client.NewBlueprint{
		CommonBlueprintFields: r.buildCommonFields(ctx, plan),
		Provider:              plan.Provider.ValueString(),
		IsTerraformManaged:    true,
	}

	blueprint, _, err := r.service.CreateBlueprint(ctx, newBlueprint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating blueprint",
			"Could not create blueprint: "+err.Error(),
		)
		return
	}

	// Set the resource state
	state := FlattenBlueprint(blueprint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get the current state
	var state BlueprintResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh value from the remote API
	blueprint, httpResp, err := r.service.GetBlueprintBySeriesId(
		ctx,
		state.SeriesId.ValueString(),
	)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Blueprint "+state.SeriesId.ValueString()+" was not found in Resourcely",
				"The blueprint may have been deleted outside of Terraform",
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error reading blueprint",
				"Could not read blueprint series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite state with refreshed value
	state = FlattenBlueprint(blueprint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve the plan and state
	var plan BlueprintResourceModel
	var state BlueprintResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	updatedBlueprint := &client.UpdatedBlueprint{
		SeriesId:              state.SeriesId.ValueString(),
		CommonBlueprintFields: r.buildCommonFields(ctx, plan),
	}

	blueprint, _, err := r.service.UpdateBlueprint(ctx, updatedBlueprint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating blueprint",
			"Could not update blueprint series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the resource state
	state = FlattenBlueprint(blueprint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve from state
	var state *BlueprintResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.service.DeleteBlueprint(ctx, state.SeriesId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting blueprint",
			"Could not delete blueprint series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *BlueprintResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("series_id"), req, resp)
}

func (r *BlueprintResource) buildCommonFields(
	ctx context.Context,
	plan BlueprintResourceModel,
) client.CommonBlueprintFields {
	commonFields := client.CommonBlueprintFields{
		Name:                          plan.Name.ValueString(),
		Description:                   plan.Description.ValueString(),
		Content:                       plan.Content.ValueString(),
		Guidance:                      plan.Guidance.ValueString(),
		Categories:                    nil,
		Labels:                        nil,
		ExcludedContextQuestionSeries: nil,
	}

	var labels []string
	plan.Labels.ElementsAs(ctx, &labels, false)
	for _, label := range labels {
		commonFields.Labels = append(commonFields.Labels, client.Label{Label: label})
	}

	plan.Categories.ElementsAs(ctx, &commonFields.Categories, false)
	plan.ExcludedContextQuestionSeries.ElementsAs(ctx, &commonFields.ExcludedContextQuestionSeries, false)

	return commonFields
}
