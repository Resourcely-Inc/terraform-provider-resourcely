package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &BlueprintTemplateResource{}
	_ resource.ResourceWithImportState = &BlueprintTemplateResource{}
)

func NewBlueprintTemplateResource() resource.Resource {
	return &BlueprintTemplateResource{}
}

// BlueprintTemplateResource defines the resource implementation.
type BlueprintTemplateResource struct {
	service *client.BlueprintTemplatesService
}

func (r *BlueprintTemplateResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_template"
}

func (r *BlueprintTemplateResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Resourcely Blueprint Template",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint template",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the blueprint template",
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
					),
				},
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
			},
			"guidance": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
			},
			"labels": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "",
				Optional:            true,
			},
			"categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "",
				Optional:            true,
				Validators: []validator.Set{
					// All list items must pass the nested validators
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						"BLUEPRINT_BLOB_STORAGE",
						"BLUEPRINT_NETWORKING",
						"BLUEPRINT_DATABASE",
						"BLUEPRINT_COMPUTE",
						"BLUEPRINT_SERVERLESS_COMPUTE",
						"BLUEPRINT_ASYNC_PROCESSING",
						"BLUEPRINT_CONTAINERIZATION",
						"BLUEPRINT_LOGS_AND_METRICS",
					),
					),
				},
			},
		},
	}
}

func (r *BlueprintTemplateResource) Configure(
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

	r.service = client.BlueprintTemplates
}

func (r *BlueprintTemplateResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Get the plan
	var plan BlueprintTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	newBlueprintTemplate := &client.NewBlueprintTemplate{
		CommonBlueprintTemplateFields: r.buildCommonFields(ctx, plan),
		Provider:                      plan.Provider.ValueString(),
	}

	blueprintTemplate, _, err := r.service.CreateBlueprintTemplate(ctx, newBlueprintTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating blueprint template",
			"Could not create blueprint template: "+err.Error(),
		)
		return
	}

	// Set the resource state
	state := FlattenBlueprintTemplate(blueprintTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintTemplateResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get the current state
	var state BlueprintTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh value from the remote API
	blueprintTemplate, httpResp, err := r.service.GetBlueprintTemplateBySeriesId(
		ctx,
		state.SeriesId.ValueString(),
	)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Blueprint template "+state.SeriesId.ValueString()+" was not found in Resourcely",
				"The blueprint template may have been deleted outside of Terraform",
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error reading blueprint template",
				"Could not read blueprint template series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite state with refreshed value
	state = FlattenBlueprintTemplate(blueprintTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintTemplateResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve the plan and state
	var plan BlueprintTemplateResourceModel
	var state BlueprintTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	updatedBlueprintTemplate := &client.UpdatedBlueprintTemplate{
		SeriesId:                      state.SeriesId.ValueString(),
		CommonBlueprintTemplateFields: r.buildCommonFields(ctx, plan),
	}

	blueprintTemplate, _, err := r.service.UpdateBlueprintTemplate(ctx, updatedBlueprintTemplate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating blueprint template",
			"Could not update blueprint template series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the resource state
	state = FlattenBlueprintTemplate(blueprintTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlueprintTemplateResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve from state
	var state *BlueprintTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.service.DeleteBlueprintTemplate(ctx, state.SeriesId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting blueprinteTemplate",
			"Could not delete blueprinteTemplate series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *BlueprintTemplateResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("series_id"), req, resp)
}

func (r *BlueprintTemplateResource) buildCommonFields(
	ctx context.Context,
	plan BlueprintTemplateResourceModel,
) client.CommonBlueprintTemplateFields {
	commonFields := client.CommonBlueprintTemplateFields{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Content:     plan.Content.ValueString(),
		Guidance:    plan.Guidance.ValueString(),
		Categories:  nil,
		Labels:      nil,
	}

	var labels []string
	plan.Labels.ElementsAs(ctx, &labels, false)
	for _, label := range labels {
		commonFields.Labels = append(commonFields.Labels, client.Label{Label: label})
	}

	plan.Categories.ElementsAs(ctx, &commonFields.Categories, false)

	return commonFields
}
