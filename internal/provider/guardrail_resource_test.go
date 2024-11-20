package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGuardrailResource_basic_withContent(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_basic_withContent("basic_test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_guardrail.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_guardrail.basic", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "category", "GUARDRAIL_BEST_PRACTICES"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "state", "GUARDRAIL_STATE_EVALUATE_ONLY"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "content",
						`GUARDRAIL "basic test"
  WHEN aws_s3_bucket
    REQUIRE bucket = "acme-{team}-{project}"
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
				Config: testAccGuardrailResourceConfig_basic_withContent("basic_test_update"),
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

func testAccGuardrailResourceConfig_basic_withContent(name string) string {
	return fmt.Sprintf(`
resource "resourcely_guardrail" "basic" {
  name = "%s"
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
`, name)
}

// Where is
// testAccGuardrailResourceConfig_basic_withGuardrailTemplate?
//
// That test would require an actual guardrail template to exist,
// which we can't safely do.
//
// The test can't create one, because templates are managed by the
// Resourcely platform. This provider does not have a
// "guardrail_template" resource.
//
// And the test cannot hardcode a template provided by the Resourcely
// platform becuase the template series ids differ in each
// environment. These tests should pass in both our dev and prod
// environments.
//
// If this becomes a problem, we could introduce a "guardrail
// template" data source into this provider.  The data source could
// lookup a specific guardrail template by an attribute that is
// consistent across all environments (e.g., name).

func TestAccGuardrailResource_errorsMissingInputs(t *testing.T) {
	expectedErrorsValidatorErrorMissingInputs := []string{
		"These attributes must be configured together",
		"guardrail_template_series_id",
		"guardrail_template_inputs",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsValidatorErrorMissingInputs...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_configValidatorErrorMissingInputs,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.guardrail_template", "guardrail_template_inputs"),
				),
			},
		},
	})
}

const testAccGuardrailResourceConfig_configValidatorErrorMissingInputs = `
resource "resourcely_guardrail" "guardrail_template" {
  name           = "basic_test"
  description    = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category       = "GUARDRAIL_BEST_PRACTICES"

  guardrail_template_series_id = "00000000-00000000-00000000-00000000"
}
`

func TestAccGuardrailResource_errorsConflicting(t *testing.T) {
	expectedErrorsValidatorErrorConflicting := []string{
		"Exactly one of these attributes must be configured",
		"content",
		"guardrail_template_series_id",
		"guardrail_template_inputs",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsValidatorErrorConflicting...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_configValidatorErrorConflicting,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.guardrail_template", "guardrail_template_inputs"),
				),
			},
		},
	})
}

const testAccGuardrailResourceConfig_configValidatorErrorConflicting = `
resource "resourcely_guardrail" "guardrail_template" {
  name           = "basic_test"
  description    = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category       = "GUARDRAIL_BEST_PRACTICES"

  content = <<-EOT
  guardrail "basic test"
	when aws_s3_bucket
	require bucket = "acme-{team}-{project}"
  end
EOT

  guardrail_template_series_id = "00000000-00000000-00000000-00000000"
  guardrail_template_inputs    = "{\"test\": \"render test\"}"
}
`

func TestAccGuardrailResource_errorsInvalidJSON(t *testing.T) {
	expectedErrorsValidatorErrorConflicting := []string{
		"Exactly one of these attributes must be configured",
		"content",
		"guardrail_template_series_id",
		"guardrail_template_inputs",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsValidatorErrorConflicting...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_configValidatorErrorInvalidJSON,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.guardrail_template", "guardrail_template_inputs"),
				),
			},
		},
	})
}

const testAccGuardrailResourceConfig_configValidatorErrorInvalidJSON = `
resource "resourcely_guardrail" "guardrail_template" {
  name           = "basic_test"
  description    = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category       = "GUARDRAIL_BEST_PRACTICES"

  guardrail_template_series_id = "00000000-00000000-00000000-00000000"
  guardrail_template_inputs    = ""
}
`
