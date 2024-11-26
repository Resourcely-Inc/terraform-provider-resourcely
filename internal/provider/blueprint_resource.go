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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
		MarkdownDescription: "A blueprint is a configuration template used to provision cloud infrastructure resources. Blueprints allow you to:\n\n- Define which options are available for properties of your resource(s).\n- Apply gaurdrails to your resource(s) to prevent misconfiguration.\n- Define what information to collect from your developers before provisioning the resource.\n\nOnce a blueprint is configured and published, it becomes available for use in your Resourcely service catalog.\n\nThe template is specified using Resourcely's TFT templating language. See the [Authoring Your Own Blueprints](https://docs.resourcely.io/build/setting-up-blueprints/authoring-your-own-blueprints) docs for details about TFT. The [Resourcely Foundry](https://portal.resourcely.io/foundry?mode=blueprint) provides an IDE to assist with authoring the template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of the blueprint.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Incrementing version number for the current version of the blueprint.",
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
				MarkdownDescription: "The name of the blueprint.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the blueprint's purpose or functionality.",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider that this blueprint targets. Can be one of `PROVIDER_AMAZON`, `PROVIDER_AZURE`, `PROVIDER_CONDUCTORONE`, `PROVIDER_DATABRICKS`, `PROVIDER_DATADOG`, `PROVIDER_GITHUB`, `PROVIDER_GITLAB`, `PROVIDER_GOOGLE`, `PROVIDER_HYPERV`, `PROVIDER_IBM`, `PROVIDER_JUMPCLOUD`, `PROVIDER_KUBERNETES`, `PROVIDER_OKTA`, `PROVIDER_ORACLE`, `PROVIDER_RESOURCELY`, `PROVIDER_SNOWFLAKE`, `PROVIDER_SPACELIFT`, `PROVIDER_VMWARE`, `PROVIDER_OTHER`",
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PROVIDER_AMAZON",
						"PROVIDER_AZURE",
						"PROVIDER_CONDUCTORONE",
						"PROVIDER_DATABRICKS",
						"PROVIDER_DATADOG",
						"PROVIDER_GITHUB",
						"PROVIDER_GITLAB",
						"PROVIDER_GOOGLE",
						"PROVIDER_HYPERV",
						"PROVIDER_IBM",
						"PROVIDER_JUMPCLOUD",
						"PROVIDER_KUBERNETES",
						"PROVIDER_OKTA",
						"PROVIDER_ORACLE",
						"PROVIDER_RESOURCELY",
						"PROVIDER_SNOWFLAKE",
						"PROVIDER_SPACELIFT",
						"PROVIDER_VMWARE",
						"PROVIDER_OTHER",
					),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The templated Terraform configuration specified using Resourcely's TFT format. See the [Authoring Your Own Blueprints](https://docs.resourcely.io/build/setting-up-blueprints/authoring-your-own-blueprints) docs for details. The [Resourcely Foundry](https://portal.resourcely.io/foundry?mode=blueprint) provides an IDE to assist with authoring the content.",
				Required:            true,
			},
			"guidance": schema.StringAttribute{
				MarkdownDescription: "Guidance to help your users know when and how to use this blueprint.",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"labels": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "Additional keywords to help your users discover this blueprint.",
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
			"is_published": schema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "A published blueprint is available for use by developers to create resources through the Resourcely portal. If left unset, the blueprint will start as unpublished, and you may safely change this property in the Resourcely portal.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_context_question_series": schema.SetAttribute{
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				ElementType:         basetypes.StringType{},
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "The series_ids for context questions that won't be used with this blueprint, even if this blueprint matches the context questions' blueprint_categories",
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
		IsPublished:           plan.IsPublished.ValueBool(), // defaults to false if not explicitly set
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

	// Compute which update methods needs to be called.
	//
	// Some fields are changed via update (PUT) and some are changed
	// via patch (PATCH).
	var blueprint *client.Blueprint
	var err error
	needsUpdate, needsPatch := r.computeUpdateActions(ctx, state, plan)

	// Update the resource
	if needsUpdate {
		updatedBlueprint := &client.UpdatedBlueprint{
			SeriesId:              state.SeriesId.ValueString(),
			CommonBlueprintFields: r.buildCommonFields(ctx, plan),
		}

		blueprint, _, err = r.service.UpdateBlueprint(ctx, updatedBlueprint)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating blueprint",
				"Could not put blueprint series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Patch the resource
	if needsPatch {
		patchedBlueprint := &client.PatchedBlueprint{
			SeriesId:    state.SeriesId.ValueString(),
			IsPublished: plan.IsPublished.ValueBool(),
		}
		blueprint, _, err = r.service.PatchBlueprint(ctx, patchedBlueprint)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating blueprint",
				"Could not patch blueprint series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
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

// Most blueprint fields are updated via an Update (Put) API call, but
// a IsPublished is updated via a Patch API call. This function
// examines the changed properties to determine which are needed.
func (r *BlueprintResource) computeUpdateActions(ctx context.Context, state BlueprintResourceModel, plan BlueprintResourceModel) (needsUpdate, needsPatch bool) {
	needsUpdate = false
	needsPatch = false

	// Determine if the Patch fields have changed
	isPublishedKnown := !plan.IsPublished.IsNull() && !plan.IsPublished.IsUnknown()
	isPublishedChanged := !state.IsPublished.Equal(plan.IsPublished)
	if isPublishedKnown && isPublishedChanged {
		needsPatch = true
	}

	// Determine if the Update fields have changed

	// If patch is not needed, we know an update is needed. Terraform
	// only calls us if there was some change. And it was't to a patch
	// field...
	if !needsPatch {
		needsUpdate = true
		return
	}

	// Otherwise, check all update fields
	if !plan.Categories.Equal(state.Categories) {
		needsUpdate = true
	}

	if !plan.Content.Equal(state.Content) {
		needsUpdate = true
	}

	if !plan.Description.Equal(state.Description) {
		needsUpdate = true
	}

	if !plan.ExcludedContextQuestionSeries.Equal(state.ExcludedContextQuestionSeries) {
		needsUpdate = true
	}

	if !plan.Guidance.Equal(state.Guidance) {
		needsUpdate = true
	}

	if !plan.Labels.Equal(state.Labels) {
		needsUpdate = true
	}

	if !plan.Name.Equal(state.Name) {
		needsUpdate = true
	}

	return
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
