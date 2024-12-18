package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlueprintDataSource_basic(t *testing.T) {
	contextQuestionLabel := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccBlueprintDataSourceConfig_basic(contextQuestionLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.resourcely_blueprint.by_series_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.resourcely_blueprint.by_series_id", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "name", "basic_test"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "labels.0", "marketing"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "guidance", "How to use this "),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "excluded_context_question_series.#", "1"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "is_published", "true"),
					resource.TestCheckResourceAttr("data.resourcely_blueprint.by_series_id", "content",
						`resource "aws_s3_bucket" "{{ resource_name }}" {
  bucket = "{{ bucket }}"
}
`),
				),
			},
		},
	})
}

func testAccBlueprintDataSourceConfig_basic(contextQuestionLabel string) string {
	return fmt.Sprintf(`
resource "resourcely_blueprint" "basic" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  categories = ["BLUEPRINT_BLOB_STORAGE"]
  labels = ["marketing"]
  guidance = "How to use this "
  content = <<-EOT
              resource "aws_s3_bucket" "{{ resource_name }}" {
                bucket = "{{ bucket }}"
              }
            EOT

  excluded_context_question_series = [resourcely_context_question.basic.series_id]

  is_published = true
}

resource "resourcely_context_question" "basic" {
	prompt = "test"
	qtype = "QTYPE_TEXT"
	scope = "SCOPE_TENANT"
	blueprint_categories = ["BLUEPRINT_BLOB_STORAGE"]
	label = "%s"
}

data "resourcely_blueprint" "by_series_id" {
  series_id = resourcely_blueprint.basic.series_id
}
`, contextQuestionLabel)
}
