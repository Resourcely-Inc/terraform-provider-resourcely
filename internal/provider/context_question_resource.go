package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &ContextQuestionResource{}
	_ resource.ResourceWithImportState = &ContextQuestionResource{}
)

func NewContextQuestionResource() resource.Resource {
	return &ContextQuestionResource{}
}

// ContextQuestionResource defines the resource implementation.
type ContextQuestionResource struct {
	service *client.ContextQuestionsService
}

func (r *ContextQuestionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context_question"
}

func (r *ContextQuestionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A [context question](https://docs.resourcely.io/concepts/other-features-and-settings/global-context-and-values) is used to gather data from developers before provisioning a resource. They are designed to gather and store insightful data related to the resource.\n\nSome examples include:\n\n- What type of data will be stored in this infrastructure?\n- What application is this infrastructure associated with?\n- What is the email address the person/team responsible for this infrastructure?\n\nThree types of context questions are supported:\n\n- Text\n- Single Choice\n- Multiple Choice\n",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID for the current version of the context question.",
				Computed:            true,
			},
			"series_id": schema.StringAttribute{
				MarkdownDescription: "UUID for the context question.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				MarkdownDescription: "Incrementing version number for the current version of this context question.",
				Computed:            true,
			},
			"scope": schema.StringAttribute{
				MarkdownDescription: "",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "A key used to reference the context question in blueprints and guardrails. Must be unique within your Resourcley tenant.",
				Required:            true,
			},
			"prompt": schema.StringAttribute{
				MarkdownDescription: "The question that Resourcely with ask your developers.",
				Required:            true,
			},
			"qtype": schema.StringAttribute{
				MarkdownDescription: "The type of the question. Must be one of `QTYPE_TEXT`, `QTYPE_SINGLE_SELECT`, or `QTYPE_MULTI_SELECT`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"QTYPE_TEXT",
						"QTYPE_SINGLE_SELECT",
						"QTYPE_MULTI_SELECT"),
				},
			},
			"answer_format": schema.StringAttribute{
				MarkdownDescription: "A format validation for acceptable answers to the context question. Applicable only when `qtype` is `QTYPE_TEXT` . Must be one of `ANSWER_TEXT`, `ANSWER_NUMBER`, `ANSWER_EMAIL`, or `ANSWER_REGEX`. If `ANSWER_REGEX`, must also specify the `regex_pattern` property.",
				Default:             stringdefault.StaticString("ANSWER_TEXT"),
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"ANSWER_TEXT",
						"ANSWER_NUMBER",
						"ANSWER_EMAIL",
						"ANSWER_REGEX",
					),
				},
			},
			"answer_choices": schema.SetNestedAttribute{
				MarkdownDescription: "The answer choices from which the developer can select. Applicable only when `qtype` is `QTYPE_SINGLE_SELECT` or `QTYPE_MULTI_SELECT`.",
				Default:             setdefault.StaticValue(types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"label": types.StringType}}, nil)),
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"label": schema.StringAttribute{
							MarkdownDescription: "The value for the answer choice.",
							Required:            true,
						},
					},
				},
			},
			"blueprint_categories": schema.SetAttribute{
				ElementType:         basetypes.StringType{},
				MarkdownDescription: "The blueprint categories to which this context question applies. This question will be asked whenever a developer uses a blueprint in these categories.",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				Computed:            true,
				Optional:            true,
				Validators: []validator.Set{
					// All list items must pass the nested validators
					setvalidator.ValueStringsAre(stringvalidator.OneOf(
						"BLUEPRINT_ASYNC_PROCESSING",
						"BLUEPRINT_BLOB_STORAGE",
						"BLUEPRINT_COMPUTE",
						"BLUEPRINT_CONTAINERIZATION",
						"BLUEPRINT_DATABASE",
						"BLUEPRINT_GITHUB_REPO",
						"BLUEPRINT_GITHUB_REPO_TEAM",
						"BLUEPRINT_IAM",
						"BLUEPRINT_LOGS_AND_METRICS",
						"BLUEPRINT_NETWORKING",
						"BLUEPRINT_SERVERLESS_COMPUTE",
					),
					),
				},
			},
			"regex_pattern": schema.StringAttribute{
				MarkdownDescription: "A regex validation for the acceptable answers to the context question. Applicable only when both `qtype` is `QTYPE_TEXT` and `answer_format` is `ANSWER_REGEX`.",
				Default:             stringdefault.StaticString(""),
				Optional:            true,
				Computed:            true,
			},
			"excluded_blueprint_series": schema.SetAttribute{
				MarkdownDescription: "The series_ids of blueprints for which this question should not be asked, even if those blueprints belong to the context question's blueprint_categories.",
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, nil)),
				ElementType:         basetypes.StringType{},
				Optional:            true,
				Computed:            true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of this question, relative to others. 0=high, 1=medium, 2=low",
				Default:             int64default.StaticInt64(0),
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2),
				},
			},
		},
	}
}

func (r *ContextQuestionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.service = client.ContextQuestions
}

func (r *ContextQuestionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the plan
	var plan ContextQuestionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the resource
	newContextQuestion := &client.NewContextQuestion{
		CommonContextQuestionFields: r.buildCommonFields(ctx, plan),
		IsTerraformManaged:          true,
	}

	ContextQuestion, _, err := r.service.CreateContextQuestion(ctx, newContextQuestion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Global Context",
			"Could not create Global Context: "+err.Error(),
		)
		return
	}

	// Set the resource state
	state := FlattenContextQuestion(ContextQuestion)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContextQuestionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get the current state
	var state ContextQuestionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Refresh value from the remote API
	contextQuestionResponse, httpResp, err := r.service.GetContextQuestionBySeriesId(ctx, state.SeriesId.ValueString())
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"Global Context "+state.SeriesId.ValueString()+" was not found in Resourcely",
				"The Global Context may have been deleted outside of Terraform",
			)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error reading Global Context",
				"Could not read Global Context series id "+state.SeriesId.ValueString()+": "+err.Error(),
			)
			return
		}
	}

	// Overwrite state with refreshed value
	state = FlattenContextQuestion(contextQuestionResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContextQuestionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the plan and state
	var plan ContextQuestionResourceModel
	var state ContextQuestionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the resource
	updatedContextQuestion := &client.UpdatedContextQuestion{
		SeriesId:                    state.SeriesId.ValueString(),
		CommonContextQuestionFields: r.buildCommonFields(ctx, plan),
	}

	contextQuestion, _, err := r.service.UpdateContextQuestion(ctx, updatedContextQuestion)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Global Context",
			"Could not update Global Context series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set the resource state
	state = FlattenContextQuestion(contextQuestion)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContextQuestionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve from state
	var state *ContextQuestionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.service.DeleteContextQuestion(ctx, state.SeriesId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting context question",
			"Could not delete context question series id "+state.SeriesId.ValueString()+": "+err.Error(),
		)
		return
	}
}

func (r *ContextQuestionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("series_id"), req, resp)
}

func (r *ContextQuestionResource) buildCommonFields(ctx context.Context, plan ContextQuestionResourceModel) client.CommonContextQuestionFields {
	commonFields := client.CommonContextQuestionFields{
		Label:                   plan.Label.ValueString(),
		Prompt:                  plan.Prompt.ValueString(),
		Qtype:                   plan.Qtype.ValueString(),
		AnswerFormat:            plan.AnswerFormat.ValueString(),
		Scope:                   plan.Scope.ValueString(),
		AnswerChoices:           nil,
		BlueprintCategories:     nil,
		RegexPattern:            plan.RegexPattern.ValueString(),
		ExcludedBlueprintSeries: nil,
		Priority:                plan.Priority.ValueInt64(),
	}

	var answerChoices []client.AnswerChoice
	for _, choice := range plan.AnswerChoices {
		answerChoices = append(answerChoices, client.AnswerChoice{Label: choice.Label.ValueString()})
	}
	commonFields.AnswerChoices = answerChoices

	var blueprintCategories []types.String
	_ = plan.BlueprintCategories.ElementsAs(ctx, &blueprintCategories, false)
	for _, category := range blueprintCategories {
		commonFields.BlueprintCategories = append(commonFields.BlueprintCategories, category.ValueString())
	}

	var excludedBlueprintSeries []types.String
	_ = plan.ExcludedBlueprintSeries.ElementsAs(ctx, &excludedBlueprintSeries, false)
	for _, seriesID := range excludedBlueprintSeries {
		commonFields.ExcludedBlueprintSeries = append(commonFields.ExcludedBlueprintSeries, seriesID.ValueString())
	}

	return commonFields
}
