package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var UUID_REGEX = regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")

// IsKnown returns true if the String represents known, non-null
// value.
func IsKnown(s basetypes.StringValue) bool {
	return !s.IsNull() && !s.IsUnknown()
}
