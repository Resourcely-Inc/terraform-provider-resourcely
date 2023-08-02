package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ContextQuestionOption struct {
	Label types.String `tfsdk:"label"`
}

type ContextQuestionResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SeriesId types.String `tfsdk:"series_id"`
	Version  types.Int64  `tfsdk:"version"`

	Label                   types.String            `tfsdk:"label"`
	Prompt                  types.String            `tfsdk:"prompt"`
	Qtype                   types.String            `tfsdk:"qtype"`
	AnswerFormat            types.String            `tfsdk:"answer_format"`
	Scope                   types.String            `tfsdk:"scope"`
	ContextQuestionOptions  []ContextQuestionOption `tfsdk:"context_question_options"`
	BlueprintCategories     types.Set               `tfsdk:"blueprint_categories"`
	RegexPattern            types.String            `tfsdk:"regex_pattern"`
	ExcludedBlueprintSeries types.Set               `tfsdk:"excluded_blueprint_series"`
}

func FlattenContextQuestion(contextQuestion *client.ContextQuestion) ContextQuestionResourceModel {
	var data ContextQuestionResourceModel
	data.Id = types.StringValue(contextQuestion.Id)
	data.SeriesId = types.StringValue(contextQuestion.SeriesId)
	data.Version = types.Int64Value(contextQuestion.Version)

	data.Label = types.StringValue(contextQuestion.Label)
	data.Prompt = types.StringValue(contextQuestion.Prompt)
	data.Qtype = types.StringValue(contextQuestion.Qtype)
	data.AnswerFormat = types.StringValue(contextQuestion.AnswerFormat)
	data.Scope = types.StringValue(contextQuestion.Scope)

	var contextQuestionOptions = make([]ContextQuestionOption, 0)
	for _, contextQuestionOption := range contextQuestion.ContextQuestionOptions {
		contextQuestionOptions = append(contextQuestionOptions, ContextQuestionOption{Label: basetypes.NewStringValue(contextQuestionOption.Label)})
	}
	data.ContextQuestionOptions = contextQuestionOptions

	var blueprintCategories []attr.Value
	for _, blueprintCategory := range contextQuestion.BlueprintCategories {
		blueprintCategories = append(blueprintCategories, basetypes.NewStringValue(blueprintCategory))
	}
	data.BlueprintCategories = types.SetValueMust(basetypes.StringType{}, blueprintCategories)

	data.RegexPattern = types.StringValue(contextQuestion.RegexPattern)

	var excludedBlueprintSeries []attr.Value
	for _, seriesID := range contextQuestion.ExcludedBlueprintSeries {
		excludedBlueprintSeries = append(excludedBlueprintSeries, basetypes.NewStringValue(seriesID))
	}
	data.ExcludedBlueprintSeries = types.SetValueMust(basetypes.StringType{}, excludedBlueprintSeries)

	return data
}
