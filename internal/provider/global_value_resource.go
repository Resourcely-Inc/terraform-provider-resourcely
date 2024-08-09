package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &GlobalValueResource{}
	_ resource.ResourceWithImportState = &GlobalValueResource{}
)

func NewGlobalValueResource() resource.Resource {
	return &GlobalValueResource{}
}

// GlobalValueResource defines the resource implementation.
type GlobalValueResource struct {
	service *client.GlobalValuesService
}

func (r *GlobalValueResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_global_value"
}

func (r *GlobalValueResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Resourcely GlobalValue",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version of the global value",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the global value",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Version of the global value",
				Computed:            true,
			},
			"is_deprecated": schema.BoolAttribute{
				MarkdownDescription: "True if the global value should not be used in new blueprints or guardrails",
				Default:             booldefault.StaticBool(false),
				Computed:            true,
				Optional:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "An immutable identifier used to reference this global value in blueprints or guardrails.\n\nMust start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9_]*$"), "Key must start with a lowercase letter in`a-z` and include only characters in `a-z0-9_`."),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A short display name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A longer description",
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of options in the global value. Can be one of `PRESET_VALUE_TEXT`, `PRESET_VALUE_NUMBER`, `PRESET_VALUE_LIST`, `PRESET_VALUE_OBJECT`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"PRESET_VALUE_TEXT",
						"PRESET_VALUE_NUMBER",
						"PRESET_VALUE_LIST",
						"PRESET_VALUE_OBJECT",
					),
				},
			},
			"options": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "An immutable identifier for ths option.\n\nMust start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(regexp.MustCompile("^[a-z][a-z0-9_]*$"), "Key must start with a lowercase letter in`a-z` and include only characters in `a-z0-9_`."),
							}},
						"label": schema.StringAttribute{
							MarkdownDescription: "A unique short display name",
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A longer description",
							Default:             stringdefault.StaticString(""),
							Computed:            true,
							Optional:            true,
						},
						"value": schema.StringAttribute{
							CustomType:          jsontypes.NormalizedType{},
							MarkdownDescription: "A JSON encoding of the option's value. This value must match the declared type of the global value.\n\nExample: `value = jsonencode(\"a\")`\n\nExample: `value = jsonencode([\"a\", \"b\"])`",
							Required:            true,
						},
					},
				},
				MarkdownDescription: "The list of value options for this global value",
				Required:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *GlobalValueResource) Configure(
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

	r.service = client.GlobalValues
}

func (r *GlobalValueResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Get the plan
	var plan GlobalValueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	var newGlobalValue client.NewGlobalValue
	resp.Diagnostics.Append(r.buildCommonFields(ctx, plan, &newGlobalValue.CommonGlobalValueFields)...)
	newGlobalValue.Key = plan.Key.ValueString()
	newGlobalValue.Type = plan.Type.ValueString()

	globalValue, _, err := r.service.CreateGlobalValue(ctx, &newGlobalValue)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating global value",
			"Could not create global value: "+err.Error(),
		)
		return
	}

	// Set the resource state
	var state GlobalValueResourceModel
	resp.Diagnostics.Append(FlattenGlobalValue(globalValue, &state)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GlobalValueResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get the current state
	var state GlobalValueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh value from the remote API
	globalValue, httpResp, err := r.service.GetGlobalValueBySeriesId(
		ctx,
		state.SeriesId.ValueString(),
	)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Global value "+state.SeriesId.ValueString()+" was not found in Resourcely",
				"The global value may have been deleted outside of Terraform",
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error reading global value",
				"Could not read global value series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite state with refreshed value
	resp.Diagnostics.Append(FlattenGlobalValue(globalValue, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GlobalValueResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Retrieve the plan and state
	var plan GlobalValueResourceModel
	var state GlobalValueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	var updatedGlobalValue client.UpdatedGlobalValue
	updatedGlobalValue.SeriesId = state.SeriesId.ValueString()
	resp.Diagnostics.Append(r.buildCommonFields(ctx, plan, &updatedGlobalValue.CommonGlobalValueFields)...)

	globalValue, _, err := r.service.UpdateGlobalValue(ctx, &updatedGlobalValue)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating global value",
			"Could not update global value series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the resource state
	resp.Diagnostics.Append(FlattenGlobalValue(globalValue, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Note: The presets API does not yet support deletion.
func (r *GlobalValueResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve from state
	var state *GlobalValueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Not implemented; the presets API does not support DELETE
	resp.Diagnostics.AddWarning(
		"Dropping global value from state without deleting",
		"The global values API does not support deletion. Dropped global value series id "+state.SeriesId.ValueString()+" from the Terraform state anyway.",
	)
}

func (r *GlobalValueResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("series_id"), req, resp)
}

func (r *GlobalValueResource) buildCommonFields(
	ctx context.Context,
	plan GlobalValueResourceModel,
	fields *client.CommonGlobalValueFields,
) diag.Diagnostics {
	var diags diag.Diagnostics

	fields.Name = plan.Name.ValueString()
	fields.Description = plan.Description.ValueString()

	fields.Options = make([]client.GlobalValueOption, len(plan.Options))
	for i, option := range plan.Options {
		diags.Append(r.buildOption(ctx, option, &fields.Options[i])...)
	}

	return diags
}

func (r *GlobalValueResource) buildOption(ctx context.Context, plan GlobalValueOptionModel, option *client.GlobalValueOption) diag.Diagnostics {
	var diags diag.Diagnostics

	option.Key = plan.Key.ValueString()
	option.Label = plan.Label.ValueString()
	option.Description = plan.Description.ValueString()

	diags.Append(plan.Value.Unmarshal(&option.Value)...)

	return diags
}
