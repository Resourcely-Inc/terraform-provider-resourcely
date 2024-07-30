package provider

import (
	"context"
	"fmt"

	"github.com/Resourcely-Inc/terraform-provider-resourcely-internal/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GuardrailDataSource{}

func NewGuardrailDataSource() datasource.DataSource {
	return &GuardrailDataSource{}
}

// GuardrailDataSource defines the data source implementation.
type GuardrailDataSource struct {
	service *client.GuardrailsService
}

func (d *GuardrailDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guardrail"
}

func (d *GuardrailDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the guardrail",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
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
			"category": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"rego_policy": schema.StringAttribute{
				MarkdownDescription: "String of the Rego Policy",
				Computed:            true,
			},
			"cue_policy": schema.StringAttribute{
				MarkdownDescription: "String of the Cue Policy",
				Computed:            true,
			},
			"json_validation": schema.StringAttribute{
				MarkdownDescription: "JSON string of the JSON Validation used on the frontend",
				Computed:            true,
			},
			"guardrail_template_id": schema.StringAttribute{
				MarkdownDescription: "The id of the guardrail template used to render the policies",
				Computed:            true,
			},
			"guardrail_template_inputs": schema.StringAttribute{
				MarkdownDescription: "JSON string of the guardrail template inputs used on the frontend",
				Computed:            true,
			},
		},
	}
}

func (d *GuardrailDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.service = client.Guardrails
}

func (d *GuardrailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config
	var config GuardrailResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	guardrailSeriesId := config.SeriesId.ValueString()

	guardrail, _, err := d.service.GetGuardrailBySeriesId(ctx, guardrailSeriesId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading guardrail",
			"Could not read guardrail id "+guardrailSeriesId+": "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed value
	state := FlattenGuardrail(guardrail)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
