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
		MarkdownDescription: "A blueprint is a configuration template used to provision cloud infrastructure resources. Blueprints allow you to:\n\n- Define which options are available for properties of your resource(s).\n- Apply gaurdrails to your resource(s) to prevent misconfiguration.\n- Define what information to collect from your developers before provisioning the resource.\n\nOnce a blueprint is configured and published, it becomes available for use in your Resourcely service catalog.\n\nThe template is specified using Resourcely's TFT templating language. See the [Authoring Your Own Blueprints](https://docs.resourcely.io/build/setting-up-blueprints/authoring-your-own-blueprints) docs for details about TFT. The [Resourcely Foundry](https://portal.resourcely.io/foundry?mode=blueprint) provides an IDE to assist with authoring the template.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of the blueprint.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the blueprint",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Increment version number for the current version of the blueprint.",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the blueprint.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the blueprints's purpose or functionality.",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider that this blueprint targets.",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The templated Terraform configuration specified using Resourcely's TFT format.",
				Computed:            true,
			},
			"guidance": schema.StringAttribute{
				MarkdownDescription: "Guidance to help your users know when and how to use this blueprint.",
				Computed:            true,
			},
			"categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "The category to assign to this blueprint.",
			},
			"labels": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "Additional keywords to help your users discover this blueprint.",
			},
			"is_published": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "A published blueprint is available for use by developers to create resources through the Resourcely portal.",
			},
			"excluded_context_question_series": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "The series_ids for context questions that won't be used with this blueprint, even if this blueprint matches the context questions' blueprint_categories",
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
