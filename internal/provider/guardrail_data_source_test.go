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
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "state", "GUARDRAIL_STATE_EVALUATE_ONLY"),
					resource.TestCheckResourceAttr("data.resourcely_guardrail.by_series_id", "content",
						`GUARDRAIL "basic test"
  WHEN aws_s3_bucket
    REQUIRE bucket = "acme-{team}-{project}"
`),
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
  state = "GUARDRAIL_STATE_EVALUATE_ONLY"
  content = <<-EOT
              GUARDRAIL "basic test"
                WHEN aws_s3_bucket
                  REQUIRE bucket = "acme-{team}-{project}"
            EOT
}

data "resourcely_guardrail" "by_series_id" {
  series_id = resourcely_guardrail.basic_data_source.series_id
}
`
