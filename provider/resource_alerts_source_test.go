package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAlertsSource_Basic(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigBasic(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "email"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "email"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_GenericWebhook(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigGenericWebhook(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "generic_webhook"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "secret"),
					resource.TestCheckResourceAttrSet(resName, "webhook_endpoint"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_Update(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertsSourceNameUpdated := alertsSourceName + "-updated"

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigBasic(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "email"),
				),
			},
			{
				Config: testAccResourceAlertsSourceConfigBasic(alertsSourceNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceNameUpdated),
					resource.TestCheckResourceAttr(resName, "source_type", "email"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_WithDeduplication(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithDeduplication(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "deduplicate_alerts_by_key", "true"),
					resource.TestCheckResourceAttr(resName, "deduplication_key_kind", "payload"),
					resource.TestCheckResourceAttr(resName, "deduplication_key_path", "$.alert_id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_WithUrgencyRules(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	alertUrgencyName := acctest.RandomWithPrefix("tf-alert-urgency")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithUrgencyRules(alertsSourceName, alertUrgencyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "alert_source_urgency_rules_attributes.#", "1"),
					resource.TestCheckResourceAttr(resName, "alert_source_urgency_rules_attributes.0.json_path", "severity"),
					resource.TestCheckResourceAttr(resName, "alert_source_urgency_rules_attributes.0.operator", "is"),
					resource.TestCheckResourceAttr(resName, "alert_source_urgency_rules_attributes.0.value", "critical"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_WithOwnerGroups(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")
	team1Name := acctest.RandomWithPrefix("tf-team-1")
	team2Name := acctest.RandomWithPrefix("tf-team-2")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithOwnerGroups(alertsSourceName, team1Name, team2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "owner_group_ids.#", "2"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_WithSourceableAttributes(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithSourceableAttributes(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "sourceable_attributes.#", "1"),
					resource.TestCheckResourceAttr(resName, "sourceable_attributes.0.auto_resolve", "true"),
					resource.TestCheckResourceAttr(resName, "sourceable_attributes.0.resolve_state", "resolved"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_WithResolutionRule(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithResolutionRule(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "resolution_rule_attributes.#", "1"),
					resource.TestCheckResourceAttr(resName, "resolution_rule_attributes.0.enabled", "true"),
					resource.TestCheckResourceAttr(resName, "resolution_rule_attributes.0.condition_type", "all"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_AlertSourceFieldsAttributesOmitted(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigAlertSourceFieldsAttributesOmitted(alertsSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "generic_webhook"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceAlertsSource_MultipleSourceTypes(t *testing.T) {
	resName := "rootly_alerts_source.test"
	alertsSourceName := acctest.RandomWithPrefix("tf-alerts-source")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAlertsSourceConfigWithSourceType(alertsSourceName, "datadog"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "datadog"),
				),
			},
			{
				Config: testAccResourceAlertsSourceConfigWithSourceType(alertsSourceName, "new_relic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", alertsSourceName),
					resource.TestCheckResourceAttr(resName, "source_type", "new_relic"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Helper functions for generating test configurations

func testAccResourceAlertsSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "email"
}
`, name)
}

func testAccResourceAlertsSourceConfigGenericWebhook(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "generic_webhook"
}
`, name)
}

func testAccResourceAlertsSourceConfigWithDeduplication(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name                       = "%s"
  source_type                = "generic_webhook"
  deduplicate_alerts_by_key  = true
  deduplication_key_kind     = "payload"
  deduplication_key_path     = "$.alert_id"
}
`, name)
}

func testAccResourceAlertsSourceConfigWithUrgencyRules(name, urgencyName string) string {
	return fmt.Sprintf(`
resource "rootly_alert_urgency" "test" {
  name        = "%s"
  description = "%s-description"
  position    = 1
}

resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "generic_webhook"

  alert_source_urgency_rules_attributes {
    alert_urgency_id = rootly_alert_urgency.test.id
    json_path        = "severity"
    operator         = "is"
    value            = "critical"
  }
}
`, urgencyName, urgencyName, name)
}

func testAccResourceAlertsSourceConfigWithOwnerGroups(name, team1Name, team2Name string) string {
	return fmt.Sprintf(`
resource "rootly_team" "test_1" {
  name        = "%s"
  description = "%s-description"
}

resource "rootly_team" "test_2" {
  name        = "%s"
  description = "%s-description"
}

resource "rootly_alerts_source" "test" {
  name            = "%s"
  source_type     = "generic_webhook"
  owner_group_ids = [rootly_team.test_1.id, rootly_team.test_2.id]
}
`, team1Name, team1Name, team2Name, team2Name, name)
}

func testAccResourceAlertsSourceConfigWithSourceableAttributes(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "generic_webhook"

  deduplicate_alerts_by_key = true
  deduplication_key_kind = "payload"
  deduplication_key_path = "$.id"

  sourceable_attributes {
    auto_resolve  = true
    resolve_state = "resolved"

    field_mappings_attributes {
      field     = "external_id"
      json_path = "$.id"
    }

    field_mappings_attributes {
      field     = "state"
      json_path = "$.status"
    }
  }

  resolution_rule_attributes {
    condition_type = "all"
    identifier_json_path = "$.id"
    identifier_reference_kind = "payload"

    conditions_attributes {
      field = "$.status"
      operator = "is"
      value = "resolved"
      kind = "payload"
    }
  }
}
`, name)
}

func testAccResourceAlertsSourceConfigWithResolutionRule(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "email"

  resolution_rule_attributes {
    enabled                  = true
    condition_type           = "all"
    identifier_json_path     = "$.alert_id"

    conditions_attributes {
      field    = "$.status"
      operator = "is"
      value    = "resolved"
      kind     = "payload"
    }
  }
}
`, name)
}

func testAccResourceAlertsSourceConfigAlertSourceFieldsAttributesOmitted(name string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "generic_webhook"
}
`, name)
}

func testAccResourceAlertsSourceConfigWithSourceType(name, sourceType string) string {
	return fmt.Sprintf(`
resource "rootly_alerts_source" "test" {
  name        = "%s"
  source_type = "%s"
}
`, name, sourceType)
}
