package provider

import (
	"github.com/Resourcely-Inc/terraform-provider-resourcely-internal/internal/client"

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

	IsActive types.Bool `tfsdk:"is_active"`

	Content types.String `tfsdk:"content"`
}

func FlattenGuardrail(guardrail *client.Guardrail) GuardrailResourceModel {
	var data GuardrailResourceModel

	data.Id = types.StringValue(guardrail.Id)
	data.SeriesId = types.StringValue(guardrail.SeriesId)
	data.Version = types.Int64Value(guardrail.Version)
	data.Scope = types.StringValue(guardrail.Scope)

	data.Name = types.StringValue(guardrail.Name)
	data.Description = types.StringValue(guardrail.Description)
	data.Provider = types.StringValue(guardrail.Provider)
	data.Category = types.StringValue(guardrail.Category)

	data.IsActive = types.BoolValue(guardrail.IsActive)

	data.Content = types.StringValue(guardrail.Content)

	return data
}

func NormalizeGuardrail(state *GuardrailResourceModel, prior *GuardrailResourceModel) diag.Diagnostics {
	return diag.Diagnostics{}
}
