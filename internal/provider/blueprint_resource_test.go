package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlueprintResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBlueprintResourceConfig_basic("basic_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_blueprint.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_blueprint.basic", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "labels.0", "marketing"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "guidance", "How to use this blueprint"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "excluded_context_question_series.#", "0"),
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "content",
						`resource "aws_s3_bucket" "{{ resource_name }}" {
  bucket = "{{ bucket }}"
}
`),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resourcely_blueprint.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importBlueprintBySeriesId("resourcely_blueprint.basic"),
			},
			// Update and Read testing
			{
				Config: testAccBlueprintResourceConfig_basic("basic_test_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_blueprint.basic", "name", "basic_test_update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func importBlueprintBySeriesId(blueprintName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		blueprint := s.RootModule().Resources[blueprintName]
		if blueprint == nil {
			return "", fmt.Errorf("Cannot find %s in terraform state", blueprintName)
		}
		seriesId, found := blueprint.Primary.Attributes["series_id"]
		if !found {
			return "", fmt.Errorf("Blueprint %s does not have series_id in Terraform state", blueprintName)
		}
		return seriesId, nil
	}
}

func testAccBlueprintResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "resourcely_blueprint" "basic" {
  name = "%s"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  categories = ["BLUEPRINT_BLOB_STORAGE"]
  labels = ["marketing"]
  guidance = "How to use this blueprint"
  content = <<-EOT
              resource "aws_s3_bucket" "{{ resource_name }}" {
                bucket = "{{ bucket }}"
              }
            EOT
}
`, name)
}
