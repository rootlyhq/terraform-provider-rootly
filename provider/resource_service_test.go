package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceService(t *testing.T) {
	serviceName := acctest.RandomWithPrefix("tf-service")
	serviceNameUpdated := acctest.RandomWithPrefix("tf-service-updated")
	slackAliasName := acctest.RandomWithPrefix("tf-slack-alias")
	slackChannelName := acctest.RandomWithPrefix("tf-slack-channel")
	slackChannelNameUpdated := acctest.RandomWithPrefix("tf-slack-channel-updated")
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceServiceConfig(serviceName, slackAliasName, slackChannelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", serviceName),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_aliases.0.name", slackAliasName),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", slackChannelName),
				),
			},
			{
				Config: testAccResourceServiceConfigUpdated(serviceNameUpdated, slackChannelNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_service.tf", "name", serviceNameUpdated),
					resource.TestCheckResourceAttr("rootly_service.tf", "slack_channels.0.name", slackChannelNameUpdated),
					resource.TestCheckNoResourceAttr("rootly_service.tf", "slack_aliases.0.id"),
				),
			},
		},
	})
}

func testAccResourceServiceConfig(serviceName, slackAliasName, slackChannelName string) string {
	return fmt.Sprintf(`
resource "rootly_service" "tf" {
	name = "%s"
	slack_aliases {
      id   = "S0883KV6123"
      name = "%s"
	}
	slack_channels {
      id   = "C08836PQ123"
	  name = "%s"
	}
}
`, serviceName, slackAliasName, slackChannelName)
}

func testAccResourceServiceConfigUpdated(serviceName, slackChannelName string) string {
	return fmt.Sprintf(`
resource "rootly_service" "tf" {
	name = "%s"
	slack_channels {
      id   = "C08836PQ123"
	  name = "%s"
	}
}
`, serviceName, slackChannelName)
}

// Disabling this test until we get kubernetes deployment integration enabled where this
// API token lives.

// func TestAccResourceServiceWithKubernetesDeploymentName(t *testing.T) {
// 	service_name := acctest.RandomWithPrefix("tf-service")
// 	kubernetes_deployment_name := "namespace/" + acctest.RandomWithPrefix("deployment-name")
// 	tf_name := "tf-" + acctest.RandomWithPrefix("service")

// 	resource.UnitTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			testAccPreCheck(t)
// 		},
// 		ProviderFactories: providerFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccResourceServiceWithKubernetesDeploymentNameConfig(tf_name, service_name, kubernetes_deployment_name),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("rootly_service.tf", "name", service_name),
// 					resource.TestCheckResourceAttr("rootly_service.tf", "kubernetes_deployment_name", kubernetes_deployment_name),
// 				),
// 			},
// 		},
// 	})
// }

// func testAccResourceServiceWithKubernetesDeploymentNameConfig(tf_name, service_name, kubernetes_deployment_name string) string {
// 	return fmt.Sprintf(`
// 	resource "rootly_service" "%s" {
// 		name = "%s"
// 		kubernetes_deployment_name = "%s"
// 	}
// 	`, tf_name, service_name, kubernetes_deployment_name)
// }
