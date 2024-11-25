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
		MarkdownDescription: "A [context question](https://docs.resourcely.io/concepts/other-features-and-settings/global-context-and-values) is used to gather data from developers before provisioning a resource. They are designed to gather and store insightful data related to the resource.\n\nSome examples include:\n\n- What type of data will be stored in this infrastructure?\n- What application is this infrastructure associated with?\n- What is the email address the person/team responsible for this infrastructure?\n\nThree types of context questions are supported:\n\n- Text\n- Single Choice\n- Multiple Choice\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of this context question.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the context question.",
				Required:            true,
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Incrementing version number of the context question.",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "A key used to reference the context question in blueprints and guardrails. Is unique within your Resourcley tenant.",
				Computed:            true,
			},
			"prompt": schema.StringAttribute{
				MarkdownDescription: "The question that Resourcely will ask your developers.",
				Computed:            true,
			},
			"qtype": schema.StringAttribute{
				MarkdownDescription: "The type of the question. Will be one of `QYTPE_TEXT`, `QYTPE_SINGLE_SELECT`, or `QTYPE_MULTI_SELECT`.",
				Computed:            true,
			},
			"answer_format": schema.StringAttribute{
				MarkdownDescription: "A format validation for acceptable answers to the context question. Applicable only when `qtype` is `QTYPE_TEXT` . Will be one of `ANSWER_TEXT`, `ANSWER_NUMBER`, `ANSWER_EMAIL`, or `ANSWER_REGEX`. If `ANSWER_REGEX`, the `regex_pattern` property will also be set.",

				Computed: true,
			},
			"answer_choices": schema.SetNestedAttribute{
				MarkdownDescription: "The answer choices from which the developer can select. Applicable only when `qtype` is `QTYPE_SINGLE_SELECT` or `QTYPE_MULTI_SELECT`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label": schema.StringAttribute{
							MarkdownDescription: "The value for the answer choice.",
							Computed:            true,
						},
					},
				},
			},
			"blueprint_categories": schema.SetAttribute{
				MarkdownDescription: "The blueprint categories to which this context question applies. This question will be asked whenever a developer uses a blueprint in these categories.",
				ElementType:         basetypes.StringType{},
				Computed:            true,
			},
			"regex_pattern": schema.StringAttribute{
				MarkdownDescription: "A regex validation for the acceptable answers to the context question. Applicable only when both `qtype` is `QTYPE_TEXT` and `answer_format` is `ANSWER_REGEX`.",
				Computed:            true,
			},
			"excluded_blueprint_series": schema.SetAttribute{
				MarkdownDescription: "The series_ids of blueprints for which this question should not be asked, even if those blueprints belong to the context question's blueprint_categories.",
				ElementType:         basetypes.StringType{},
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of this question, relative to others. 0=high, 1=medium, 2=low",
				Computed:            true,
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
