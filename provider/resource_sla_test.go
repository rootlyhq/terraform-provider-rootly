package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccResourceSLA_basic exercises the minimum required fields.
func TestAccResourceSLA_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-sla")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSLABasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					resource.TestCheckResourceAttr("rootly_sla.test", "assignment_deadline_days", "3"),
					resource.TestCheckResourceAttr("rootly_sla.test", "assignment_deadline_parent_status", "started"),
					resource.TestCheckResourceAttr("rootly_sla.test", "completion_deadline_days", "7"),
					resource.TestCheckResourceAttr("rootly_sla.test", "completion_deadline_parent_status", "resolved"),
					resource.TestCheckResourceAttrSet("rootly_sla.test", "manager_role_id"),
					resource.TestCheckResourceAttrSet("rootly_sla.test", "slug"),
					resource.TestCheckResourceAttrSet("rootly_sla.test", "id"),
				),
			},
			// ImportState step verifies that the resource can be imported by ID.
			{
				ResourceName:      "rootly_sla.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccResourceSLA_update verifies updates propagate correctly.
func TestAccResourceSLA_update(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-sla")
	updatedName := name + "-updated"

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSLABasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					resource.TestCheckResourceAttr("rootly_sla.test", "assignment_deadline_days", "3"),
				),
			},
			{
				// Re-use the same role suffix as step 1 so the incident_role
				// isn't replaced between steps (only the SLA is updated).
				Config: testAccResourceSLAUpdated(updatedName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", updatedName),
					resource.TestCheckResourceAttr("rootly_sla.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("rootly_sla.test", "assignment_deadline_days", "5"),
					resource.TestCheckResourceAttr("rootly_sla.test", "assignment_skip_weekends", "true"),
					resource.TestCheckResourceAttr("rootly_sla.test", "completion_deadline_days", "14"),
					resource.TestCheckResourceAttr("rootly_sla.test", "completion_skip_weekends", "true"),
					resource.TestCheckResourceAttr("rootly_sla.test", "condition_match_type", "ANY"),
				),
			},
		},
	})
}

// TestAccResourceSLA_withConditions exercises the nested conditions block.
func TestAccResourceSLA_withConditions(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-sla-cond")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSLAWithConditions(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					resource.TestCheckResourceAttr("rootly_sla.test", "conditions.#", "1"),
					resource.TestCheckResourceAttr("rootly_sla.test", "conditions.0.conditionable_type", "SLAs::BuiltInFieldCondition"),
					resource.TestCheckResourceAttr("rootly_sla.test", "conditions.0.property", "severity"),
					resource.TestCheckResourceAttr("rootly_sla.test", "conditions.0.operator", "is_set"),
				),
			},
		},
	})
}

// TestAccResourceSLA_withNotifications exercises the nested
// notification_configurations block across the three offset_type variants.
func TestAccResourceSLA_withNotifications(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-sla-notif")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSLAWithNotifications(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					resource.TestCheckResourceAttr("rootly_sla.test", "notification_configurations.#", "3"),
				),
			},
		},
	})
}

// managerRoleConfig returns a unique incident_role block so tests that share
// the same test binary (or previous runs that leaked roles) don't collide on
// the "name must be unique" constraint.
func managerRoleConfig(suffix string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "test" {
	name = "tf-test-sla-manager-%s"
}
`, suffix)
}

func testAccResourceSLABasic(name string) string {
	return managerRoleConfig(name) + fmt.Sprintf(`
resource "rootly_sla" "test" {
	name                              = "%s"
	assignment_deadline_days          = 3
	assignment_deadline_parent_status = "started"
	completion_deadline_days          = 7
	completion_deadline_parent_status = "resolved"
	manager_role_id                   = rootly_incident_role.test.id
}
`, name)
}

func testAccResourceSLAUpdated(name, roleSuffix string) string {
	return managerRoleConfig(roleSuffix) + fmt.Sprintf(`
resource "rootly_sla" "test" {
	name                              = "%s"
	description                       = "Updated description"
	condition_match_type              = "ANY"
	assignment_deadline_days          = 5
	assignment_deadline_parent_status = "started"
	assignment_skip_weekends          = true
	completion_deadline_days          = 14
	completion_deadline_parent_status = "resolved"
	completion_skip_weekends          = true
	manager_role_id                   = rootly_incident_role.test.id
}
`, name)
}

func testAccResourceSLAWithConditions(name string) string {
	return managerRoleConfig(name) + fmt.Sprintf(`
resource "rootly_sla" "test" {
	name                              = "%s"
	assignment_deadline_days          = 2
	assignment_deadline_parent_status = "started"
	completion_deadline_days          = 7
	completion_deadline_parent_status = "resolved"
	manager_role_id                   = rootly_incident_role.test.id

	conditions {
		conditionable_type = "SLAs::BuiltInFieldCondition"
		property           = "severity"
		operator           = "is_set"
	}
}
`, name)
}

func testAccResourceSLAWithNotifications(name string) string {
	return managerRoleConfig(name) + fmt.Sprintf(`
resource "rootly_sla" "test" {
	name                              = "%s"
	assignment_deadline_days          = 3
	assignment_deadline_parent_status = "started"
	completion_deadline_days          = 7
	completion_deadline_parent_status = "resolved"
	manager_role_id                   = rootly_incident_role.test.id

	notification_configurations {
		offset_type = "before_due"
		offset_days = 1
	}

	notification_configurations {
		offset_type = "when_due"
		offset_days = 0
	}

	notification_configurations {
		offset_type = "after_due"
		offset_days = 1
	}
}
`, name)
}
