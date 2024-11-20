package provider

import (
	"encoding/json"

	"github.com/Resourcely-Inc/terraform-provider-resourcely/internal/client"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GuardrailModel describes the resource data model.
type GuardrailResourceModel struct {
	Id       types.String `tfsdk:"id"`
	SeriesId types.String `tfsdk:"series_id"`
	Version  types.Int64  `tfsdk:"version"`
	Scope    types.String `tfsdk:"scope"`

	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Provider    types.String `tfsdk:"cloud_provider"`
	Category    types.String `tfsdk:"category"`

	State types.String `tfsdk:"state"`

	Content types.String `tfsdk:"content"`

	GuardrailTemplateSeriesId types.String         `tfsdk:"guardrail_template_series_id"`
	GuardrailTemplateInputs   jsontypes.Normalized `tfsdk:"guardrail_template_inputs"`
}

func FlattenGuardrail(guardrail *client.Guardrail, data *GuardrailResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	data.Id = types.StringValue(guardrail.Id)
	data.SeriesId = types.StringValue(guardrail.SeriesId)
	data.Version = types.Int64Value(guardrail.Version)
	data.Scope = types.StringValue(guardrail.Scope)

	data.Name = types.StringValue(guardrail.Name)
	data.Description = types.StringValue(guardrail.Description)
	data.Provider = types.StringValue(guardrail.Provider)
	data.Category = types.StringValue(guardrail.Category)

	data.State = types.StringValue(guardrail.State)

	data.Content = types.StringValue(guardrail.Content)

	data.GuardrailTemplateSeriesId = types.StringValue(guardrail.GuardrailTemplate.SeriesId)
	if guardrail.GuardrailTemplateInputs != nil {
		guardrailTemplateInputs, err := json.Marshal(guardrail.GuardrailTemplateInputs)
		if err != nil {
			diags.AddError(
				"Failed to JSON encode the guardrail template inputs",
				"Could not JSON encode the guardrail template inputs for guardrail "+guardrail.Id+": "+err.Error(),
			)
		}
		data.GuardrailTemplateInputs = jsontypes.NewNormalizedValue(string(guardrailTemplateInputs))
	}

	return diags
}
