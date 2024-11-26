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
		MarkdownDescription: "A guardrail governs how cloud resources can be created and altered, preventing infrastructure misconfiguration. Before infrastructure is provisioned, Resourcely examines the changes being made and prevents a merge if any guardrail requirements are violated. Some examples of guardrails include:\n\n- Require approval for making a public S3 bucket\n- Restrict the allowed compute instance types or images\n\nGuardrails are specified using the [Really policy language](https://docs.resourcely.io/build/setting-up-guardrails/authoring-your-own-guardrails).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of this guar.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the guardrail.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Incrementing version number for this current version of the guardrail.",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the guardrail.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the guardrail's purpose or policy.",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider that this guardrail targets.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The category of this guardrail.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The [state](https://docs.resourcely.io/build/setting-up-guardrails/releasing-guardrails#guardrail-status) of the guardrail.",
				Computed:            true,
			},
			"content": schema.StringAttribute{
				MarkdownDescription: "The guardrail policy written in the [Really policy language](https://docs.resourcely.io/build/setting-up-guardrails/authoring-your-own-guardrails).",
				Computed:            true,
			},
			"guardrail_template_series_id": schema.StringAttribute{
				MarkdownDescription: "The series id of the guardrail template used to render the policy.",
				Computed:            true,
			},
			"guardrail_template_inputs": schema.StringAttribute{
				CustomType:          jsontypes.NormalizedType{},
				MarkdownDescription: "A JSON encoding of values for the guardrail template inputs.`",
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
	var state GuardrailResourceModel
	resp.Diagnostics.Append(FlattenGuardrail(guardrail, &state)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
