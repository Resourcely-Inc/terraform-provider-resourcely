package provider

import (
	"context"
	"fmt"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &BlueprintDataSource{}

func NewBlueprintDataSource() datasource.DataSource {
	return &BlueprintDataSource{}
}

// BlueprintDataSource defines the data source implementation.
type BlueprintDataSource struct {
	service *client.BlueprintsService
}

func (d *BlueprintDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint"
}

func (d *BlueprintDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A resourcely blueprint",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the blueprint",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"guidance": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "",
			},
			"labels": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "",
			},
			"excluded_context_question_series": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "series_id for context questions that won't be used with this blueprint, even if this blueprint matches the context questions' blueprint_categories",
			},
		},
	}
}

func (d *BlueprintDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.service = client.Blueprints
}

func (d *BlueprintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config
	var config BlueprintResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintSeriesId := config.SeriesId.ValueString()

	blueprint, _, err := d.service.GetBlueprintBySeriesId(ctx, blueprintSeriesId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading blueprint",
			"Could not read blueprint series id "+blueprintSeriesId+": "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed value
	state := FlattenBlueprint(blueprint)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
