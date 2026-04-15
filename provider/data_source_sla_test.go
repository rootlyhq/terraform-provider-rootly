package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceSLA_byName looks up an SLA by its name.
func TestAccDataSourceSLA_byName(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-ds-sla")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSLAConfig(name),
				Check: resource.ComposeTestCheckFunc(
					// The resource was created.
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					// The data source finds it.
					resource.TestCheckResourceAttr("data.rootly_sla.by_name", "name", name),
					resource.TestCheckResourceAttrPair(
						"data.rootly_sla.by_name", "id",
						"rootly_sla.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.rootly_sla.by_name", "slug",
						"rootly_sla.test", "slug",
					),
				),
			},
		},
	})
}

// TestAccDataSourceSLA_bySlug looks up an SLA by its slug.
func TestAccDataSourceSLA_bySlug(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-test-ds-sla-slug")

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:        false,
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSLABySlugConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_sla.test", "name", name),
					resource.TestCheckResourceAttrPair(
						"data.rootly_sla.by_slug", "id",
						"rootly_sla.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.rootly_sla.by_slug", "name",
						"rootly_sla.test", "name",
					),
				),
			},
		},
	})
}

func testAccDataSourceSLAConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "test" {
	name = "tf-test-ds-sla-manager"
}

resource "rootly_sla" "test" {
	name                              = "%s"
	assignment_deadline_days          = 3
	assignment_deadline_parent_status = "started"
	completion_deadline_days          = 7
	completion_deadline_parent_status = "resolved"
	manager_role_id                   = rootly_incident_role.test.id
}

data "rootly_sla" "by_name" {
	name       = "%s"
	depends_on = [rootly_sla.test]
}
`, name, name)
}

func testAccDataSourceSLABySlugConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_incident_role" "test" {
	name = "tf-test-ds-sla-slug-manager"
}

resource "rootly_sla" "test" {
	name                              = "%s"
	assignment_deadline_days          = 3
	assignment_deadline_parent_status = "started"
	completion_deadline_days          = 7
	completion_deadline_parent_status = "resolved"
	manager_role_id                   = rootly_incident_role.test.id
}

data "rootly_sla" "by_slug" {
	slug       = rootly_sla.test.slug
	depends_on = [rootly_sla.test]
}
`, name)
}
