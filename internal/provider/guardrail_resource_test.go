package provider

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "is_active", "false"),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "content",
						`guardrail "basic test"
  when aws_s3_bucket
  require bucket = "acme-{team}-{project}"
end
`),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "rego_policy", `rego policy`),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "cue_policy", `cue policy`),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "guardrail_template_id", ""),
					resource.TestCheckResourceAttrWith("resourcely_guardrail.basic", "guardrail_template_inputs", testCheckResourceAttrNull),
					resource.TestCheckResourceAttr("resourcely_guardrail.basic", "json_validation", `{}`),
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
			{
				Config: testAccGuardrailResourceConfig_guardrailTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_guardrail.guardrail_template", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_guardrail.guardrail_template", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "name", "basic_test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "description", "this is a basic test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "cloud_provider", "PROVIDER_AMAZON"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "category", "GUARDRAIL_BEST_PRACTICES"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "is_active", "false"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "content", "really render test"),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "rego_policy", `rego render test`),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "cue_policy", `cue render test`),
					resource.TestCheckResourceAttr("resourcely_guardrail.guardrail_template", "json_validation", "{\"json\":\"render test\"}"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// testAccCheckExampleWidgetAttributes verifies attributes are set correctly by
// Terraform
func testCheckResourceAttrNull(value string) error {
	if value != "null" {
		return fmt.Errorf("resource attribute should be null")
	}

	return nil
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

func TestAccGuardrailResource_errorsInvalidJSON(t *testing.T) {
	expectedErrorsInvalidJSON := []string{
		"Attribute json_validation string must be valid JSON",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsInvalidJSON...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_invalidJSON("{}]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.invalidJSON", "json_validation"),
				),
			},
		},
	})
}

func TestAccGuardrailResource_errorsInvalidInputs(t *testing.T) {
	expectedErrorsInvalidInputs := []string{
		"Attribute guardrail_template_inputs string must be valid JSON",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsInvalidInputs...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_invalidInput("{}]"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.invalidInput", "guardrail_template_inputs"),
				),
			},
		},
	})
}

func TestAccGuardrailResource_errorsMissingPolicies(t *testing.T) {
	expectedErrorsValidatorErrorMissingPolicies := []string{
		"These attributes must be configured together",
		"rego_policy",
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ErrorCheck:               ErrorCheckExpectedErrorMessagesContaining(t, expectedErrorsValidatorErrorMissingPolicies...),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGuardrailResourceConfig_configValidatorErrorMissingPolicies,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("resourcely_guardrail.guardrail_template", "guardrail_template_inputs"),
				),
			},
		},
	})
}

func TestAccGuardrailResource_errorsMissingInputs(t *testing.T) {
	expectedErrorsValidatorErrorMissingInputs := []string{
		"These attributes must be configured together",
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

func TestAccGuardrailResource_errorsConflicting(t *testing.T) {
	expectedErrorsValidatorErrorConflicting := []string{
		"Exactly one of these attributes must be configured",
		"guardrail_template_id",
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

func testAccGuardrailResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "resourcely_guardrail" "basic" {
  name = "%s"
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
`, name)
}

func testAccGuardrailResourceConfig_invalidJSON(json string) string {
	return fmt.Sprintf(`
resource "resourcely_guardrail" "invalidJSON" {
  name = "test"
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
  json_validation = "%s"
}
`, json)
}

func testAccGuardrailResourceConfig_invalidInput(json string) string {
	return fmt.Sprintf(`
resource "resourcely_guardrail" "invalidInput" {
  name = "test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  guardrail_template_inputs = "%s"
}
`, json)
}

const testAccGuardrailResourceConfig_guardrailTemplate = `
resource "resourcely_guardrail" "guardrail_template" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  guardrail_template_id = resourcely_guardrail_template.guardrail_template.id
  guardrail_template_inputs = "{\"test\": \"render test\"}"
}

resource "resourcely_guardrail_template" "guardrail_template" {
	name= "basic guardrail"
	description= "This is a nice guardrail"
	cloud_provider= "PROVIDER_GOOGLE"
	category= "GUARDRAIL_ACCESS_CONTROL"
	really_template= "really {{ test }}"
	rego_template= "rego {{ test }}"
	json_template= "{\"json\":\"{{ test }}\"}"
	cue_template= "cue {{ test }}"
	template_schema= "test:"
}
`

const testAccGuardrailResourceConfig_configValidatorErrorMissingPolicies = `
resource "resourcely_guardrail" "guardrail_template" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  content = ""
}
`

const testAccGuardrailResourceConfig_configValidatorErrorMissingInputs = `
resource "resourcely_guardrail" "guardrail_template" {
  name = "basic_test"
  description = "this is a basic test"
  cloud_provider = "PROVIDER_AMAZON"
  category = "GUARDRAIL_BEST_PRACTICES"
  guardrail_template_id = resourcely_guardrail_template.guardrail_template.id
}

resource "resourcely_guardrail_template" "guardrail_template" {
	name= "basic guardrail"
	description= "This is a nice guardrail"
	cloud_provider= "PROVIDER_GOOGLE"
	category= "GUARDRAIL_ACCESS_CONTROL"
	really_template= "really {{ test }}"
	rego_template= "rego {{ test }}"
	json_template= "{\"json\":\"{{ test }}\"}"
	cue_template= "cue {{ test }}"
	template_schema= "test:"
}
`

const testAccGuardrailResourceConfig_configValidatorErrorConflicting = `
resource "resourcely_guardrail" "guardrail_template" {
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
  guardrail_template_id = resourcely_guardrail_template.guardrail_template.id
  guardrail_template_inputs = "{\"test\": \"render test\"}"
}

resource "resourcely_guardrail_template" "guardrail_template" {
	name= "basic guardrail"
	description= "This is a nice guardrail"
	cloud_provider= "PROVIDER_GOOGLE"
	category= "GUARDRAIL_ACCESS_CONTROL"
	really_template= "really {{ test }}"
	rego_template= "rego {{ test }}"
	json_template= "{\"json\":\"{{ test }}\"}"
	cue_template= "cue {{ test }}"
	template_schema= "test:"
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
