package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceFormField(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-form-field")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFormFieldConfig(rName),
				Check:  resource.TestCheckResourceAttr("rootly_form_field.test", "resource_type", "incident"),
			},
		},
	})
}

func TestAccResourceFormFieldInvalidResourceType(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-form-field-invalid")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "rootly_form_field" "test" {
  kind          = "custom"
  name          = %q
  resource_type = "bogus"
}
`, rName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("expected resource_type to be one of"),
			},
		},
	})
}

func TestAccResourceFormFieldProblem(t *testing.T) {
	if os.Getenv("ROOTLY_TEST_PROBLEM_MANAGEMENT") != "1" {
		t.Skip("Set ROOTLY_TEST_PROBLEM_MANAGEMENT=1 for an organization with problem management enabled")
	}
	rName := acctest.RandomWithPrefix("tf-form-field-problem")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProblemFormFieldConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_form_field.test", "kind", "custom"),
					resource.TestCheckResourceAttr("rootly_form_field.test", "input_kind", "select"),
					resource.TestCheckResourceAttr("rootly_form_field.test", "name", rName),
					resource.TestCheckResourceAttr("rootly_form_field.test", "resource_type", "problem"),
					resource.TestCheckResourceAttr("rootly_form_field.test", "shown.#", "0"),
					resource.TestCheckResourceAttr("rootly_form_field.test", "required.#", "0"),
				),
			},
			{
				ResourceName:      "rootly_form_field.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceFormFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_field" "test" {
	kind = "custom"
	name = "%s"
}
`, name)
}

func testAccResourceProblemFormFieldConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_form_field" "test" {
  kind          = "custom"
  name          = %q
  input_kind    = "select"
  resource_type = "problem"
}

resource "rootly_form_field_option" "investigating" {
  form_field_id = rootly_form_field.test.id
  value         = "Investigating"
}

resource "rootly_form_field_option" "resolved" {
  form_field_id = rootly_form_field.test.id
  value         = "Resolved"
}
`, name)
}
