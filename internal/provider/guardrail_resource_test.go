package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGuardrailResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_basic("basic_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_guardrail.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_guardrail.basic", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "category", "GUARDRAIL_BEST_PRACTICES"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "state", "GUARDRAIL_STATE_EVALUATE_ONLY"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "content",
						`guardrail "basic test"
  when aws_s3_bucket
  require bucket = "acme-{team}-{project}"
end
`),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resourcely_guardrail.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importGuardrailBySeriesId("resourcely_guardrail.basic"),
			},
			// Update and Read testing
			{
				Config: testAccGuardrailResourceConfig_basic("basic_test_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "name", "basic_test_update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func importGuardrailBySeriesId(guardrailName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		guardrail := s.RootModule().Resources[guardrailName]
		if guardrail == nil {
			return "", fmt.Errorf("Cannot find %s in terraform state", guardrailName)
		}
		seriesId, found := guardrail.Primary.Attributes["series_id"]
		if !found {
			return "", fmt.Errorf("Guardrail %s does not have series_id in Terraform state", guardrailName)
		}
		return seriesId, nil
	}
}

func testAccGuardrailResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "resourcely_guardrail" "basic" {
  name = "%s"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  state = "GUARDRAIL_STATE_EVALUATE_ONLY"
  content = <<-EOT
              guardrail "basic test"
                when aws_s3_bucket
                require bucket = "acme-{team}-{project}"
              end
            EOT
}
`, name)
}
