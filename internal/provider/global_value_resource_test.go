package provider

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/matoous/go-nanoid/v2"
)

func importGlobalValueBySeriesId(global_valueName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		global_value := s.RootModule().Resources[global_valueName]
		if global_value == nil {
			return "", fmt.Errorf("Cannot find %s in terraform state", global_valueName)
		}
		seriesId, found := global_value.Primary.Attributes["series_id"]
		if !found {
			return "", fmt.Errorf("GlobalValue %s does not have series_id in Terraform state", global_valueName)
		}
		return seriesId, nil
	}
}

func TestAccGlobalValueResource_basic_text(t *testing.T) {
	id := gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyz", 16)
	key := "basic_text_" + id

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalValueResourceConfig_basic_text(key, "Basic Text Test"),
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
			// ImportState testing
			{
				ResourceName:      "resourcely_global_value.basic_text",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importGlobalValueBySeriesId("resourcely_global_value.basic_text"),
			},
			// Update and Read testing
			{
				Config: testAccGlobalValueResourceConfig_basic_text(key, "Basic Text Test Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_global_value.basic_text", "name", "Basic Text Test Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGlobalValueResourceConfig_basic_text(key string, name string) string {
	return fmt.Sprintf(`
resource "resourcely_global_value" "basic_text" {
  key         = "%s"
  name        = "%s"
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
`, key, name)
}

func TestAccGlobalValueResource_basic_object(t *testing.T) {
	id := gonanoid.MustGenerate("abcdefghijklmnopqrstuvwxyz", 16)
	key := "basic_object_" + id

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalValueResourceConfig_basic_object(key, "Basic Object Test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_global_value.basic_object", "id", UUID_REGEX),
					resource.TestMatchResourceAttr("resourcely_global_value.basic_object", "series_id", UUID_REGEX),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "key", key),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "name", "Basic Object Test"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "description", "This is a basic object test"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "type", "PRESET_VALUE_OBJECT"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.0.key", "option_0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.0.label", "Option 0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.0.description", "This is option 0"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.0.value", "{\"bool\":false,\"list\":[\"a\",\"b\",\"c\"],\"number\":0}"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.1.key", "option_1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.1.label", "Option 1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.1.description", "This is option 1"),
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "options.1.value", "{\"bool\":true,\"list\":[\"d\",\"e\",\"f\"],\"number\":1}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resourcely_global_value.basic_object",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importGlobalValueBySeriesId("resourcely_global_value.basic_object"),
			},
			// Update and Read testing
			{
				Config: testAccGlobalValueResourceConfig_basic_object(key, "Basic Object Test Updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_global_value.basic_object", "name", "Basic Object Test Updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGlobalValueResourceConfig_basic_object(key string, name string) string {
	return fmt.Sprintf(`
resource "resourcely_global_value" "basic_object" {
  key         = "%s"
  name        = "%s"
  description = "This is a basic object test"

  type    = "PRESET_VALUE_OBJECT"
  options = [
    {
      key         = "option_0"
      label       = "Option 0"
      description = "This is option 0"
      value       = jsonencode({
                      bool   = false
                      number = 0
                      list   = ["a", "b", "c"]
                    })
    },
    {
      key         = "option_1"
      label       = "Option 1"
      description = "This is option 1"
      value       = jsonencode({
                      bool   = true
                      number = 1
                      list   = ["d", "e", "f"]
                    })
    }
  ]
}
`, key, name)
}

func TestAccGlobalValueResource_errorNoOptions(t *testing.T) {
	expectedErrors := []string{
		"Attribute options list must contain at least 1 elements",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrors...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalValueResourceConfig_errorNoOptions,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_global_value.error_no_options", "options"),
				),
			},
		},
	})
}

const testAccGlobalValueResourceConfig_errorNoOptions = `
resource "resourcely_global_value" "error_no_options" {
  key     = "error_no_options"
  name    = "Error No Options"
  type    = "PRESET_VALUE_TEXT"
  options = []
}
`

func TestAccGlobalValueResource_errorInvalidJson(t *testing.T) {
	expectedErrors := []string{
		"Attribute options[0].value string must be valid JSON",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrors...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGlobalValueResourceConfig_errorInvalidJson,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_global_value.error_invalid_json", "options"),
				),
			},
		},
	})
}

const testAccGlobalValueResourceConfig_errorInvalidJson = `
resource "resourcely_global_value" "error_invalid_json" {
  key     = "error_invalid_json"
  name    = "Error Invalid JSON"
  type    = "PRESET_VALUE_TEXT"
  options = [
    {
      key   = "option_0"
      label = "Option 0"
      value = "{not valid json"
    }
  ]
}
`

func ErrorCheckExpectedErrorMessagesContaining(t *testing.T, messages ...string) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return errors.New("Expecting an error, got none")
		}

		for _, message := range messages {
			errorMessage := err.Error()
			if strings.Contains(errorMessage, message) {
				t.Skipf("found expected error message: %s", errorMessage)
				return nil
			}
		}

		return err
	}
}
