package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceService(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", "Terraform"),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_aliases.0.name", "eng-terraform"),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", "terraform"),
				),
			},
			{
				Config: testAccResourceServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", "Terraform (updated)"),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", "terraform (updated)"),
					resource.TestCheckNoResourceAttr("rootly_service.tf", "slack_aliases.0.id"),
				),
			},
		},
	})
}

const testAccResourceServiceCreate = `
resource "rootly_service" tf {
	name = "Terraform"
	slack_aliases {
      id   = "S0883KV6123"
      name = "eng-terraform"
	}
	slack_channels {
      id   = "C08836PQ123"
	  name = "terraform"
	}
}
`

const testAccResourceServiceUpdate = `
resource "rootly_service" tf {
	name = "Terraform (updated)"
	slack_channels {
      id   = "C08836PQ123"
	  name = "terraform (updated)"
	}
}
`

func TestAccResourceServiceWithKubernetesDeploymentName(t *testing.T) {
	service_name := acctest.RandomWithPrefix("tf-service")
	kubernetes_deployment_name := "namespace/" + acctest.RandomWithPrefix("deployment-name")
	tf_name := "tf-" + acctest.RandomWithPrefix("service")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceWithKubernetesDeploymentNameConfig(tf_name, service_name, kubernetes_deployment_name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", service_name),
					resource.TestCheckResourceAttr("rootly_service.tf", "kubernetes_deployment_name", kubernetes_deployment_name),
				),
			},
		},
	})
}

func testAccResourceServiceWithKubernetesDeploymentNameConfig(tf_name, service_name, kubernetes_deployment_name string) string {
	return fmt.Sprintf(`
	resource "rootly_service" "%s" {
		name = "%s"
		kubernetes_deployment_name = "%s"
	}
	`, tf_name, service_name, kubernetes_deployment_name)
}
