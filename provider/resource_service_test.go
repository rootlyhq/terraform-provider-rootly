package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceService(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() {
			testAccPreCheck(t)
			time.Sleep(1 * time.Second)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceService,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.foo", "name", "myservice"),
					resource.TestCheckResourceAttr("rootly_service.foo", "description", ""),
					resource.TestCheckResourceAttrSet("rootly_service.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_service.foo", "color", "#047BF8"),
				),
			},
			{
				Config: testAccResourceServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.foo", "name", "myservice2"),
					resource.TestCheckResourceAttr("rootly_service.foo", "description", "test description"),
					resource.TestCheckResourceAttr("rootly_service.foo", "public_description", "my public description"),
					resource.TestCheckResourceAttrSet("rootly_service.foo", "slug"),
					resource.TestCheckResourceAttr("rootly_service.foo", "color", "#203"),
				),
			},
		},
	})
}

const testAccResourceService = `
resource "rootly_service" "foo" {
  name = "myservice"
}
`

const testAccResourceServiceUpdate = `
resource "rootly_service" "foo" {
  name               = "myservice2"
  description        = "test description"
  color              = "#203"
  public_description = "my public description"
}
`
