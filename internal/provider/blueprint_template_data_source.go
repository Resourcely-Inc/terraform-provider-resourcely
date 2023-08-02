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
var _ datasource.DataSource = &BlueprintTemplateDataSource{}

func NewBlueprintTemplateDataSource() datasource.DataSource {
	return &BlueprintTemplateDataSource{}
}

// BlueprintTemplateDataSource defines the data source implementation.
type BlueprintTemplateDataSource struct {
	service *client.BlueprintTemplatesService
}

func (d *BlueprintTemplateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blueprint_template"
}

func (d *BlueprintTemplateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A resourcely blueprintTemplate",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint template",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the blueprint template",
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
		},
	}
}

func (d *BlueprintTemplateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.service = client.BlueprintTemplates
}

func (d *BlueprintTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config
	var config BlueprintTemplateResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	blueprintTemplateSeriesId := config.SeriesId.ValueString()

	blueprintTemplate, _, err := d.service.GetBlueprintTemplateBySeriesId(ctx, blueprintTemplateSeriesId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading blueprintTemplate",
			"Could not read blueprintTemplate series id "+blueprintTemplateSeriesId+": "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed value
	state := FlattenBlueprintTemplate(blueprintTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
