package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCommunicationsTemplate(t *testing.T) {
	templateName := acctest.RandomWithPrefix("tf-comm-template")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsTemplate(templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_template.test", "name", templateName),
				),
			},
		},
	})
}

func TestAccResourceCommunicationsTemplateWithStages(t *testing.T) {
	templateName := acctest.RandomWithPrefix("tf-comm-template-stages")
	stageName := acctest.RandomWithPrefix("tf-stage")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsTemplateWithStages(templateName, stageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_template.test_stages", "name", templateName),
					resource.TestCheckResourceAttr("rootly_communications_template.test_stages", "communication_template_stages.#", "1"),
				),
			},
		},
	})
}

func testAccResourceCommunicationsTemplate(templateName string) string {
	return fmt.Sprintf(`
resource "rootly_communications_type" "test" {
  name = "tf-test-type-%s"
  color = "#FF5733"
}

resource "rootly_communications_template" "test" {
  name = "%s"
  description = "Test communications template"
  communication_type_id = rootly_communications_type.test.id
}
`, templateName, templateName)
}

func testAccResourceCommunicationsTemplateWithStages(templateName, stageName string) string {
	return fmt.Sprintf(`
resource "rootly_communications_type" "test" {
  name = "tf-test-type-%s"
  color = "#FF5733"
}

resource "rootly_communications_stage" "test" {
  name = "%s"
  description = "Test stage"
}

resource "rootly_communications_template" "test_stages" {
  name = "%s"
  description = "Test communications template with stages"
  communication_type_id = rootly_communications_type.test.id

  communication_template_stages {
    data {
      attributes {
        email_subject = "Incident Update"
        email_body = "Incident details here"
        slack_content = "Slack notification content"
        sms_content = "SMS alert message"

        communication_stage {
          id = rootly_communications_stage.test.id
        }
      }
    }
  }
}
`, templateName, stageName, templateName)
}
