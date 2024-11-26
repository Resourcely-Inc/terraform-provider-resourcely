package provider

import (
	"context"
	"fmt"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GlobalValueDataSource{}

func NewGlobalValueDataSource() datasource.DataSource {
	return &GlobalValueDataSource{}
}

// GlobalValueDataSource defines the data source implementation.
type GlobalValueDataSource struct {
	service *client.GlobalValuesService
}

func (d *GlobalValueDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_value"
}

func (d *GlobalValueDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A [global value](https://docs.resourcely.io/concepts/other-features-and-settings/global-values) allows admins to define custom drop-downs for customizing Terraform infrastructure resource properties before they are provisioned.  They are useful for providing access to lists of relatively static values like VPC IDs, allowed regions, department or team names, etc.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of the global value.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the global value.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Incrementing version number for the current version of the global value.",
				Computed:            true,
			},
			"is_deprecated": schema.BoolAttribute{
				MarkdownDescription: "Set to true if the global value should not be used in new blueprints or guardrails",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "An immutable identifier used to reference this global value in blueprints or guardrails.\n\nMust start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the global value.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the purpose of the global value.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of options in the global value. Will be one of `PRESET_VALUE_TEXT`, `PRESET_VALUE_NUMBER`, `PRESET_VALUE_LIST`, `PRESET_VALUE_OBJECT`",
				Computed:            true,
			},
			"options": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "An immutable identifier for ths option. Must start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "A unique display name",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A description of this option's meaning.",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							CustomType:          jsontypes.NormalizedType{},
							MarkdownDescription: "A JSON encoding of the option's value.`",
							Computed:            true,
						},
					},
				},
				MarkdownDescription: "The list of value options for this global value.",
				Computed:            true,
			},
		},
	}
}

func (d *GlobalValueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.service = client.GlobalValues
}

func (d *GlobalValueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config
	var config GlobalValueResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	globalValueSeriesId := config.SeriesId.ValueString()

	globalValue, _, err := d.service.GetGlobalValueBySeriesId(ctx, globalValueSeriesId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading global value",
			"Could not read global value series id "+globalValueSeriesId+": "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed value
	var state GlobalValueResourceModel
	resp.Diagnostics.Append(FlattenGlobalValue(globalValue, &state)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
