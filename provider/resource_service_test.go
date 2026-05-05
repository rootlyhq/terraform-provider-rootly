package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceService(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceCreateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", rName),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_aliases.0.name", "eng-terraform"),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", "terraform"),
				),
			},
			{
				Config: testAccResourceServiceUpdateConfig(rName + "-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", rName+"-updated"),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", "terraform"),
					resource.TestCheckNoResourceAttr("rootly_service.tf", "slack_aliases.0.id"),
				),
			},
		},
	})
}

func testAccResourceServiceCreateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_service" "tf" {
	name = "%s"
	slack_aliases {
      id   = "S0883KV6123"
      name = "eng-terraform"
	}
	slack_channels {
      id   = "C08836PQ123"
	  name = "terraform"
	}
}
`, name)
}

func testAccResourceServiceUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "rootly_service" "tf" {
	name = "%s"
	slack_channels {
      id   = "C08836PQ123"
	  name = "terraform"
	}
}
`, name)
}
