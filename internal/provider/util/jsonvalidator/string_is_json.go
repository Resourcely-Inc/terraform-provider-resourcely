package jsonvalidator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = stringIsJson{}

// stringIsJson validates that a string is a valid JSON
type stringIsJson struct{}

func StringIsJSON() validator.String {
	return stringIsJson{}
}

// Description describes the validation in plain text formatting.
func (validator stringIsJson) Description(_ context.Context) string {
	return "string must be valid JSON"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator stringIsJson) MarkdownDescription(ctx context.Context) string {
	return "string must be valid JSON"
}

// Validate performs the validation.
func (v stringIsJson) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	var js json.RawMessage
	isJSON := json.Unmarshal([]byte(value), &js) == nil

	if !isJSON {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))

		return
	}
}
