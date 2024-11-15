package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccContextQuestionResource_basic(t *testing.T) {
	rLabel := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccContextQuestionResourceConfig_basic(rLabel, "what is your prompt?"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_context_question.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_context_question.basic", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "prompt", "what is your prompt?"),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "qtype", "QTYPE_SINGLE_SELECT"),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "blueprint_categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "answer_choices.0.label", "tenant-context Option 1"),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "label", rLabel),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "regex_pattern", `regex`),
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "priority", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resourcely_context_question.basic",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: importContextQuestionBySeriesId("resourcely_context_question.basic"),
			},
			// Update and Read testing
			{
				Config: testAccContextQuestionResourceConfig_basic(rLabel, "basic_test_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resourcely_context_question.basic", "prompt", "basic_test_update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccContextQuestionResource_useDefaults(t *testing.T) {
	rLabel := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccContextQuestionResourceConfig_useDefaults(rLabel, "what is your prompt?"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("resourcely_context_question.usedefaults", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("resourcely_context_question.usedefaults", "series_id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("resourcely_context_question.usedefaults", "prompt", "what is your prompt?"),
					resource.TestCheckResourceAttr("resourcely_context_question.usedefaults", "qtype", "QTYPE_TEXT"),
					resource.TestCheckResourceAttr("resourcely_context_question.usedefaults", "blueprint_categories.0", "BLUEPRINT_BLOB_STORAGE"),
					resource.TestCheckResourceAttr("resourcely_context_question.usedefaults", "priority", "0"),
					resource.TestCheckResourceAttr("resourcely_context_question.usedefaults", "label", rLabel),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func importContextQuestionBySeriesId(contextQuestionName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		contextQuestion := s.RootModule().Resources[contextQuestionName]
		if contextQuestion == nil {
			return "", fmt.Errorf("Cannot find %s in terraform state", contextQuestionName)
		}
		seriesId, found := contextQuestion.Primary.Attributes["series_id"]
		if !found {
			return "", fmt.Errorf("ContextQuestion %s does not have series_id in Terraform state", contextQuestionName)
		}
		return seriesId, nil
	}
}

func testAccContextQuestionResourceConfig_basic(label, prompt string) string {
	return fmt.Sprintf(`
resource "resourcely_context_question" "basic" {
	prompt = "%s"
	qtype = "QTYPE_SINGLE_SELECT"
	scope = "SCOPE_TENANT"
	blueprint_categories = ["BLUEPRINT_BLOB_STORAGE"]
	answer_choices = [{label: "tenant-context Option 1"}]
	label = "%s"
	regex_pattern = "regex"
	excluded_blueprint_series = []
	priority = 2
}
`, prompt, label)
}

func testAccContextQuestionResourceConfig_useDefaults(label, prompt string) string {
	return fmt.Sprintf(`
resource "resourcely_context_question" "usedefaults" {
	prompt = "%s"
	qtype = "QTYPE_TEXT"
	scope = "SCOPE_TENANT"
	blueprint_categories = ["BLUEPRINT_BLOB_STORAGE"]
	label = "%s"
}
`, prompt, label)
}
