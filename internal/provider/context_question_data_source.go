package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &ContextQuestionDataSource{}

func NewContextQuestionDataSource() datasource.DataSource {
	return &ContextQuestionDataSource{}
}

// ContextQuestionDataSource defines the data source implementation.
type ContextQuestionDataSource struct {
	service *client.ContextQuestionsService
}

func (d *ContextQuestionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context_question"
}

func (d *ContextQuestionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "A resourcely ContextQuestion",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for this version.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the global context",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Specific version of the global context",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"prompt": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"qtype": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"answer_format": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"answer_choices": schema.SetNestedAttribute{
				MarkdownDescription: "",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label": schema.StringAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
					},
				},
			},
			"blueprint_categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "Resource categories the context question applies to",
			},
			"regex_pattern": schema.StringAttribute{
				MarkdownDescription: "Regex validation for the acceptable answers to the context question",
				Computed:            true,
			},
			"excluded_blueprint_series": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				Computed:            true,
				MarkdownDescription: "series_id for Blueprints exempt from this context question even though those blueprints belong to the context question's blueprint_categories",
			},
		},
	}
}

func (d *ContextQuestionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.service = client.ContextQuestions
}

func (d *ContextQuestionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Read the config
	var config ContextQuestionResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	contextQuestionSeriesId := config.SeriesId.ValueString()

	contextQuestionRead, _, err := d.service.GetContextQuestionBySeriesId(ctx, contextQuestionSeriesId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading ContextQuestion",
			"Could not read ContextQuestion series id "+contextQuestionSeriesId+": "+err.Error(),
		)
		return
	}

	// Overwrite state with refreshed value
	state := FlattenContextQuestion(contextQuestionRead)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
