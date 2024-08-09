package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/matoous/go-nanoid/v2"
)

func TestAccGlobalValueDataSource_basic_text(t *testing.T) {
	id := gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyz", 16)
	key := "basic_text_" + id

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalValueDataSourceConfig_basic_text(key),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_global_value.basic_text", "id", UUID_REGEX),
					resource.TestMatchResourceAttr("resourcely_global_value.basic_text", "series_id", UUID_REGEX),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "key", key),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "name", "Basic Text Test"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "description", "This is a basic text test"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "type", "PRESET_VALUE_TEXT"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.0.key", "option_0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.0.label", "Option 0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.0.description", "This is option 0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.0.value", "\"option_0_value\""),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.1.key", "option_1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.1.label", "Option 1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.1.description", "This is option 1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "options.1.value", "\"option_1_value\""),
				),
			},
		},
	})
}

func testAccGlobalValueDataSourceConfig_basic_text(key string) string {
	return fmt.Sprintf(`
resource "resourcely_global_value" "basic_text" {
  key         = "%s"
  name        = "Basic Text Test"
  description = "This is a basic text test"

  type    = "PRESET_VALUE_TEXT"
  options = [
    {
      key         = "option_0"
      label       = "Option 0"
      description = "This is option 0"
      value       = "\"option_0_value\""
    },
    {
      key         = "option_1"
      label       = "Option 1"
      description = "This is option 1"
      value       = "\"option_1_value\""
    }
  ]
}
`, key)
}
