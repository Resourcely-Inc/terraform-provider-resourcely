package provider

import (
	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BlueprintResourceModel describes the resource data model.
type BlueprintResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SeriesId types.String `tfsdk:"series_id"`
	Version  types.Int64  `tfsdk:"version"`
	Scope    types.String `tfsdk:"scope"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Provider    types.String `tfsdk:"cloud_provider"`

	Content                       types.String `tfsdk:"content"`
	Categories                    types.Set    `tfsdk:"categories"`
	Guidance                      types.String `tfsdk:"guidance"`
	Labels                        types.Set    `tfsdk:"labels"`
	ExcludedContextQuestionSeries types.Set    `tfsdk:"excluded_context_question_series"`
}

func FlattenBlueprint(blueprint *client.Blueprint) BlueprintResourceModel {
	var data BlueprintResourceModel

	data.Id = types.StringValue(blueprint.Id)
	data.SeriesId = types.StringValue(blueprint.SeriesId)
	data.Version = types.Int64Value(blueprint.Version)
	data.Scope = types.StringValue(blueprint.Scope)

	data.Name = types.StringValue(blueprint.Name)
	data.Description = types.StringValue(blueprint.Description)
	data.Provider = types.StringValue(blueprint.Provider)

	data.Content = types.StringValue(blueprint.Content)
	data.Guidance = types.StringValue(blueprint.Guidance)

	var labels []attr.Value
	for _, label := range blueprint.Labels {
		labels = append(labels, basetypes.NewStringValue(label.Label))
	}
	data.Labels = types.SetValueMust(basetypes.StringType{}, labels)

	var categories []attr.Value
	for _, category := range blueprint.Categories {
		categories = append(categories, basetypes.NewStringValue(category))
	}
	data.Categories = types.SetValueMust(basetypes.StringType{}, categories)

	var excludedContextQuestionSeries []attr.Value
	for _, excludedCQS := range blueprint.ExcludedContextQuestionSeries {
		excludedContextQuestionSeries = append(excludedContextQuestionSeries, basetypes.NewStringValue(excludedCQS))
	}
	data.ExcludedContextQuestionSeries = types.SetValueMust(basetypes.StringType{}, excludedContextQuestionSeries)

	return data
}
