package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDashboardPanel_UpgradeFromVersion(t *testing.T) {
	resName := "rootly_dashboard_panel.foo"
	dashboardName := acctest.RandomWithPrefix("tf-dashboard")

	config := fmt.Sprintf(`
		resource "rootly_dashboard" "foo" {
			name = "%[1]s"
		}

		resource "rootly_dashboard_panel" "foo" {
			dashboard_id = rootly_dashboard.foo.id
			name = "test"
			position {
				x = 3
				y = 4
				h = 5
				w = 6
			}
			params {
				display = "line_chart"
				legend {
					groups = "charted"
				}
				datasets {
					collection = "incidents"
					filter {
						operation = "and"
						rules {
							operation = "and"
							condition = "="
							key = "status"
							value = "started"
						}
					}
					group_by = "severity"
					aggregate {
						cumulative = false
						key = "results"
						operation = "count"
					}
				}
			}
		}
	`, dashboardName)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resName, "name", "test"),
		resource.TestCheckResourceAttr(resName, "position.#", "1"),
		resource.TestCheckResourceAttr(resName, "position.0.x", "3"),
		resource.TestCheckResourceAttr(resName, "position.0.y", "4"),
		resource.TestCheckResourceAttr(resName, "position.0.h", "5"),
		resource.TestCheckResourceAttr(resName, "position.0.w", "6"),
		resource.TestCheckResourceAttr(resName, "params.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.display", "line_chart"),
		resource.TestCheckResourceAttr(resName, "params.0.legend.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.legend.0.groups", "charted"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.collection", "incidents"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.operation", "and"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.operation", "and"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.condition", "="),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.key", "status"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.value", "started"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.group_by", "severity"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.cumulative", "false"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.key", "results"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.operation", "count"),
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rootly": {
						Source:            "rootlyhq/rootly",
						VersionConstraint: "4.3.4",
					},
				},
				Config: config,
				Check:  check,
			},
			{
				ProviderFactories: providerFactories,
				Config:            config,
				Check:             check,
			},
		},
	})
}

func TestAccResourceDashboardPanel_Import(t *testing.T) {
	resName := "rootly_dashboard_panel.foo"
	dashboardName := acctest.RandomWithPrefix("tf-dashboard")

	config := fmt.Sprintf(`
		resource "rootly_dashboard" "foo" {
			name = "%[1]s"
		}

		resource "rootly_dashboard_panel" "foo" {
			dashboard_id = rootly_dashboard.foo.id
			name = "test"
			position {
				x = 3
				y = 4
				h = 5
				w = 6
			}
			params {
				display = "line_chart"
			}
		}
	`, dashboardName)

	check := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resName, "name", "test"),
		resource.TestCheckResourceAttr(resName, "position.#", "1"),
		resource.TestCheckResourceAttr(resName, "position.0.x", "3"),
		resource.TestCheckResourceAttr(resName, "position.0.y", "4"),
		resource.TestCheckResourceAttr(resName, "position.0.h", "5"),
		resource.TestCheckResourceAttr(resName, "position.0.w", "6"),
		resource.TestCheckResourceAttr(resName, "params.#", "1"),
		resource.TestCheckResourceAttr(resName, "params.0.display", "line_chart"),
		resource.TestCheckResourceAttr(resName, "params.0.legend.#", "0"),
		resource.TestCheckResourceAttr(resName, "params.0.datasets.#", "0"),
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDashboardPanel(t *testing.T) {
	resName := "rootly_dashboard_panel.foo"

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDashboardPanel,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "test"),
					resource.TestCheckResourceAttr(resName, "position.#", "1"),
					resource.TestCheckResourceAttr(resName, "position.0.x", "3"),
					resource.TestCheckResourceAttr(resName, "position.0.y", "4"),
					resource.TestCheckResourceAttr(resName, "position.0.h", "5"),
					resource.TestCheckResourceAttr(resName, "position.0.w", "6"),
					resource.TestCheckResourceAttr(resName, "params.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.display", "line_chart"),
					resource.TestCheckResourceAttr(resName, "params.0.description", "description"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.0.groups", "charted"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.collection", "incidents"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.operation", "and"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.operation", "and"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.condition", "="),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.key", "status"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.value", "started"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.group_by", "severity"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.cumulative", "false"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.key", "results"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.operation", "count"),
				),
			},
			{
				Config: testAccResourceDashboardPanel,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "test"),
					resource.TestCheckResourceAttr(resName, "position.#", "1"),
					resource.TestCheckResourceAttr(resName, "position.0.x", "3"),
					resource.TestCheckResourceAttr(resName, "position.0.y", "4"),
					resource.TestCheckResourceAttr(resName, "position.0.h", "5"),
					resource.TestCheckResourceAttr(resName, "position.0.w", "6"),
					resource.TestCheckResourceAttr(resName, "params.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.display", "line_chart"),
					resource.TestCheckResourceAttr(resName, "params.0.description", "description"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.0.groups", "charted"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.collection", "incidents"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.operation", "and"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.operation", "and"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.condition", "="),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.key", "status"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.value", "started"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.group_by", "severity"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.cumulative", "false"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.key", "results"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.operation", "count"),
				),
			},
			{
				Config: testAccResourceDashboardPanelUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "test-updated"),
					resource.TestCheckResourceAttr(resName, "position.#", "1"),
					resource.TestCheckResourceAttr(resName, "position.0.x", "30"),
					resource.TestCheckResourceAttr(resName, "position.0.y", "40"),
					resource.TestCheckResourceAttr(resName, "position.0.h", "50"),
					resource.TestCheckResourceAttr(resName, "position.0.w", "60"),
					resource.TestCheckResourceAttr(resName, "params.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.display", "line_chart"),
					resource.TestCheckResourceAttr(resName, "params.0.description", "description-updated"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.legend.0.groups", "charted"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.collection", "incidents"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.operation", "or"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.operation", "or"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.condition", "="),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.key", "status"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.filter.0.rules.0.value", "started"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.group_by", "severity"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.#", "1"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.cumulative", "true"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.key", "results"),
					resource.TestCheckResourceAttr(resName, "params.0.datasets.0.aggregate.0.operation", "count"),
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

const testAccResourceDashboardPanel = `
resource "rootly_dashboard" "foo" {
	name = "my-dashboard-with-panel"
}

resource "rootly_dashboard_panel" "foo" {
	dashboard_id = rootly_dashboard.foo.id
	name = "test"
	position {
		x = 3
		y = 4
		h = 5
		w = 6
	}
	params {
		display = "line_chart"
		description = "description"
		legend {
			groups = "charted"
		}
		datasets {
			collection = "incidents"
			filter {
				operation = "and"
				rules {
					operation = "and"
					condition = "="
					key = "status"
					value = "started"
				}
			}
			group_by = "severity"
			aggregate {
				cumulative = false
				key = "results"
				operation = "count"
			}
		}
	}
}
`

const testAccResourceDashboardPanelUpdated = `
resource "rootly_dashboard" "foo" {
	name = "my-dashboard-with-panel"
}

resource "rootly_dashboard_panel" "foo" {
	dashboard_id = rootly_dashboard.foo.id
	name = "test-updated"
	position {
		x = 30
		y = 40
		h = 50
		w = 60
	}
	params {
		display = "line_chart"
		description = "description-updated"
		legend {
			groups = "charted"
		}
		datasets {
			collection = "incidents"
			filter {
				operation = "or"
				rules {
					operation = "or"
					condition = "="
					key = "status"
					value = "started"
				}
			}
			group_by = "severity"
			aggregate {
				cumulative = true
				key = "results"
				operation = "count"
			}
		}
	}
}
`
