package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContextQuestionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccContextQuestionDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.resourcely_context_question.by_series_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.resourcely_context_question.by_series_id", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "prompt", "what is your prompt?"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "qtype", "QTYPE_SINGLE_SELECT"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "blueprint_categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "answer_choices.0.label", "tenant-context Option 1"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "label", "marketing"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "regex_pattern", "regex"),
					resource.TestCheckResourceAttr("data.resourcely_context_question.by_series_id", "priority", "2"),
				),
			},
		},
	})
}

const testAccContextQuestionDataSourceConfig_basic = `
resource "resourcely_context_question" "basic" {
	prompt = "what is your prompt?"
	qtype = "QTYPE_SINGLE_SELECT"
	scope = "SCOPE_TENANT"
	blueprint_categories = ["BLUEPRINT_BLOB_STORAGE"]
	answer_choices = [{label: "tenant-context Option 1"}]
	label = "marketing"
	regex_pattern = "regex"
	excluded_blueprint_series = []
	priority = 2
}

data "resourcely_context_question" "by_series_id" {
  series_id = resourcely_context_question.basic.series_id
}
`
