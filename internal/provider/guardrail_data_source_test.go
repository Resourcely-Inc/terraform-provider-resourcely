package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGuardrailDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccGuardrailDataSourceConfig_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.resourcely_guardrail.by_series_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.resourcely_guardrail.by_series_id", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "name", "basic_test"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "category", "GUARDRAIL_BEST_PRACTICES"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "is_active", "false"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "content",
						`guardrail "basic test"
when aws_s3_bucket
require bucket = "acme-{team}-{project}"
end
`),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "rego_policy", `rego policy`),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "cue_policy", `cue policy`),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "json_validation", `{}`),
				),
			},
			{
				Config: testAccGuardrailDataSourceConfig_guardrailTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.resourcely_guardrail.by_series_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.resourcely_guardrail.by_series_id", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "name", "basic_test"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "category", "GUARDRAIL_BEST_PRACTICES"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "is_active", "false"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "content", "really render test"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "rego_policy", `rego render test`),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "cue_policy", `cue render test`),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "json_validation", "{\"json\":\"render test\"}"),
					resource.TestMatchResourceAttr("data.resourcely_guardrail.by_series_id", "guardrail_template_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "guardrail_template_inputs", "{\"test\":\"render test\"}"),
				),
			},
		},
	})
}

const testAccGuardrailDataSourceConfig_basic = `
resource "resourcely_guardrail" "basic_data_source" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  content = <<-EOT
              guardrail "basic test"
              when aws_s3_bucket
              require bucket = "acme-{team}-{project}"
              end
            EOT
  rego_policy = "rego policy"
  cue_policy = "cue policy"
  json_validation = "{}"
}

data "resourcely_guardrail" "by_series_id" {
  series_id = resourcely_guardrail.basic_data_source.series_id
}
`

const testAccGuardrailDataSourceConfig_guardrailTemplate = `
resource "resourcely_guardrail" "guardrail_template_data_source" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  guardrail_template_id = resourcely_guardrail_template.guardrail_template_data_source.id
  guardrail_template_inputs = "{\"test\":\"render test\"}"
}

resource "resourcely_guardrail_template" "guardrail_template_data_source" {
	name= "basic guardrail"
	description= "This is a nice guardrail"
	cloud_provider= "PROVIDER_GOOGLE"
	category= "GUARDRAIL_ACCESS_CONTROL"
	really_template= "really {{ test }}"
	rego_template= "rego {{ test }}"
	json_template= "{\"json\": \"{{ test }}\"}"
	cue_template= "cue {{ test }}"
	template_schema= "test:"
}

data "resourcely_guardrail" "by_series_id" {
  series_id = resourcely_guardrail.guardrail_template_data_source.series_id
}
`
