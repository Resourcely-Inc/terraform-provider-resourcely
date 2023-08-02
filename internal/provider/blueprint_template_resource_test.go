package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlueprintTemplateResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccBlueprintTemplateResourceConfig_basic("basic_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_blueprint_template.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_blueprint_template.basic", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "labels.0", "marketing"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "guidance", "How to use this template"),
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "content",
						`resource "aws_s3_bucket" "{{ resource_name }}" {
  bucket = "{{ bucket }}"
}
`),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resourcely_blueprint_template.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importBlueprintTemplateBySeriesId("resourcely_blueprint_template.basic"),
			},
			// Update and Read testing
			{
				Config: testAccBlueprintTemplateResourceConfig_basic("basic_test_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_blueprint_template.basic", "name", "basic_test_update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func importBlueprintTemplateBySeriesId(blueprintTemplateName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		blueprintTemplate := s.RootModule().Resources[blueprintTemplateName]
		if blueprintTemplate == nil {
			return "", fmt.Errorf("Cannot find %s in terraform state", blueprintTemplateName)
		}
		seriesId, found := blueprintTemplate.Primary.Attributes["series_id"]
		if !found {
			return "", fmt.Errorf("BlueprintTemplate %s does not have series_id in Terraform state", blueprintTemplateName)
		}
		return seriesId, nil
	}
}

func testAccBlueprintTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "resourcely_blueprint_template" "basic" {
  name = "%s"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  categories = ["BLUEPRINT_BLOB_STORAGE"]
  labels = ["marketing"]
  guidance = "How to use this template"
  content = <<-EOT
              resource "aws_s3_bucket" "{{ resource_name }}" {
                bucket = "{{ bucket }}"
              }
            EOT
}
`, name)
}
