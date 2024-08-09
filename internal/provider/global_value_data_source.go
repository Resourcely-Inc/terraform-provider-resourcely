package provider

import (
	"context"
	"fmt"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

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
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A Resourcely global value",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the global value",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the global value",
				Computed:            true,
			},
			"is_deprecated": schema.BoolAttribute{
				MarkdownDescription: "True if the global value should not be used in new blueprints or guardrails",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "An immutable identifier used to reference this global value in blueprints or guardrails.\n\nMust start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A short display name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A longer description",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of options in the global value. Can be one of `PRESET_VALUE_TEXT`, `PRESET_VALUE_NUMBER`, `PRESET_VALUE_LIST`, `PRESET_VALUE_OBJECT`",
				Computed:            true,
			},
			"options": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "An immutable identifier for ths option.\n\nMust start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "A unique short display name",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "A longer description",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "A JSON encoding of the option's value. This value must match the declared type of the global value.\n\nExample: `value = jsonencode(\"a\")`\n\nExample: `value = jsonencode([\"a\", \"b\"])`",
							Computed:            true,
						},
					},
				},
				MarkdownDescription: "The list of value options for this global value",
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
