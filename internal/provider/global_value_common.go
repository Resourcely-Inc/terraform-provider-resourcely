package provider

import (
	"encoding/json"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GlobalValueResourceModel describes the resource data model.
type GlobalValueResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SeriesId types.String `tfsdk:"series_id"`
	Version  types.Int64  `tfsdk:"version"`

	IsDeprecated types.Bool `tfsdk:"is_deprecated"`

	Key         types.String             `tfsdk:"key"`
	Name        types.String             `tfsdk:"name"`
	Description types.String             `tfsdk:"description"`
	Type        types.String             `tfsdk:"type"`
	Options     []GlobalValueOptionModel `tfsdk:"options"`
}

type GlobalValueOptionModel struct {
	Key         types.String         `tfsdk:"key"`
	Label       types.String         `tfsdk:"label"`
	Description types.String         `tfsdk:"description"`
	Value       jsontypes.Normalized `tfsdk:"value"`
}

func FlattenGlobalValue(global_value *client.GlobalValue, data *GlobalValueResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.Id = types.StringValue(global_value.Id)
	data.SeriesId = types.StringValue(global_value.SeriesId)
	data.Version = types.Int64Value(global_value.Version)

	data.IsDeprecated = types.BoolValue(global_value.IsDeprecated)

	data.Key = types.StringValue(global_value.Key)
	data.Name = types.StringValue(global_value.Name)
	data.Description = types.StringValue(global_value.Description)
	data.Type = types.StringValue(global_value.Type)

	data.Options = make([]GlobalValueOptionModel, len(global_value.Options))
	for i, option := range global_value.Options {
		diags.Append(FlattenGlobalValueOption(option, &data.Options[i])...)
	}

	return diags
}

func FlattenGlobalValueOption(option client.GlobalValueOption, data *GlobalValueOptionModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.Key = types.StringValue(option.Key)
	data.Label = types.StringValue(option.Label)
	data.Description = types.StringValue(option.Description)

	// JSON encode the value
	value, err := json.Marshal(option.Value)
	if err != nil {
		diags.AddError(
			"Failed to JSON encode global value option",
			"Could not JSON encode the value for global value option "+option.Key+": "+err.Error(),
		)
	}
	data.Value = jsontypes.NewNormalizedValue(string(value))

	return diags
}
