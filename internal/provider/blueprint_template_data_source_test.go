package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlueprintTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBlueprintTemplateDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.resourcely_blueprint_template.by_series_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.resourcely_blueprint_template.by_series_id", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "name", "basic_test"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "labels.0", "marketing"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "guidance", "How to use this template"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint_template.by_series_id", "content",
						`resource "aws_s3_bucket" "{{ resource_name }}" {
  bucket = "{{ bucket }}"
}
`),
				),
			},
		},
	})
}

const testAccBlueprintTemplateDataSourceConfig_basic = `
resource "resourcely_blueprint_template" "basic" {
  name = "basic_test"
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

data "resourcely_blueprint_template" "by_series_id" {
  series_id = resourcely_blueprint_template.basic.series_id
}
`
