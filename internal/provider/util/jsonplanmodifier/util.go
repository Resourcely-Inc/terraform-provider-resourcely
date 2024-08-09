package jsonplanmodifier

import (
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// IsKnown returns true if the String represents known, non-null
// value.
func IsKnown(s basetypes.StringValue) bool {
	return !s.IsNull() && !s.IsUnknown()
}
