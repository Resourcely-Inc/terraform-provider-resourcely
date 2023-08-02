package provider

import (
	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// BlueprintTemplateResourceModel describes the resource data model.
type BlueprintTemplateResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SeriesId types.String `tfsdk:"series_id"`
	Version  types.Int64  `tfsdk:"version"`
	Scope    types.String `tfsdk:"scope"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Provider    types.String `tfsdk:"cloud_provider"`

	Content    types.String `tfsdk:"content"`
	Categories types.Set    `tfsdk:"categories"`
	Guidance   types.String `tfsdk:"guidance"`
	Labels     types.Set    `tfsdk:"labels"`
}

func FlattenBlueprintTemplate(blueprintTemplate *client.BlueprintTemplate) BlueprintTemplateResourceModel {
	var data BlueprintTemplateResourceModel

	data.Id = types.StringValue(blueprintTemplate.Id)
	data.SeriesId = types.StringValue(blueprintTemplate.SeriesId)
	data.Version = types.Int64Value(blueprintTemplate.Version)
	data.Scope = types.StringValue(blueprintTemplate.Scope)

	data.Name = types.StringValue(blueprintTemplate.Name)
	data.Description = types.StringValue(blueprintTemplate.Description)
	data.Provider = types.StringValue(blueprintTemplate.Provider)

	data.Content = types.StringValue(blueprintTemplate.Content)
	data.Guidance = types.StringValue(blueprintTemplate.Guidance)

	var labels []attr.Value
	for _, label := range blueprintTemplate.Labels {
		labels = append(labels, basetypes.NewStringValue(label.Label))
	}
	data.Labels = types.SetValueMust(basetypes.StringType{}, labels)

	var categories []attr.Value
	for _, category := range blueprintTemplate.Categories {
		categories = append(categories, basetypes.NewStringValue(category))
	}
	data.Categories = types.SetValueMust(basetypes.StringType{}, categories)

	return data
}
