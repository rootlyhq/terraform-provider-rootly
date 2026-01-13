package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/acctest"
)

func TestAccResourceDashboardPanel(t *testing.T) {
	resName := "rootly_dashboard_panel.test"
	dashboardName := acctest.RandomWithPrefix("tf-dashboard")
	dashboardPanelName := acctest.RandomWithPrefix("tf-dashboard-panel")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName, `
					position {
						x = 3
						y = 4
						h = 5
						w = 6
					}

					params {
						display = "line_chart"
						description = "description"
						// No Legend
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
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(dashboardPanelName)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("position").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"x": knownvalue.Int64Exact(3),
						"y": knownvalue.Int64Exact(4),
						"h": knownvalue.Int64Exact(5),
						"w": knownvalue.Int64Exact(6),
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("display"), knownvalue.StringExact("line_chart")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact("description")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("legend"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("datasets").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"name":       knownvalue.StringExact(""),
						"collection": knownvalue.StringExact("incidents"),
						"filter": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"operation": knownvalue.StringExact("and"),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.MapExact(map[string]knownvalue.Check{
										"operation": knownvalue.StringExact("and"),
										"condition": knownvalue.StringExact("="),
										"key":       knownvalue.StringExact("status"),
										"value":     knownvalue.StringExact("started"),
									}),
								}),
							}),
						}),
						"group_by": knownvalue.StringExact("severity"),
						"aggregate": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"cumulative": knownvalue.Bool(false),
								"key":        knownvalue.StringExact("results"),
								"operation":  knownvalue.StringExact("count"),
							}),
						}),
					})),
				),
			},
			{
				Config: testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName+"-updated", `
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
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(dashboardPanelName+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("position").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"x": knownvalue.Int64Exact(30),
						"y": knownvalue.Int64Exact(40),
						"h": knownvalue.Int64Exact(50),
						"w": knownvalue.Int64Exact(60),
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("display"), knownvalue.StringExact("line_chart")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact("description-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("legend").AtSliceIndex(0).AtMapKey("groups"), knownvalue.StringExact("charted")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("datasets").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"name":       knownvalue.StringExact(""),
						"collection": knownvalue.StringExact("incidents"),
						"filter": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"operation": knownvalue.StringExact("or"),
								"rules": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.MapExact(map[string]knownvalue.Check{
										"operation": knownvalue.StringExact("or"),
										"condition": knownvalue.StringExact("="),
										"key":       knownvalue.StringExact("status"),
										"value":     knownvalue.StringExact("started"),
									}),
								}),
							}),
						}),
						"group_by": knownvalue.StringExact("severity"),
						"aggregate": knownvalue.ListExact([]knownvalue.Check{
							knownvalue.MapExact(map[string]knownvalue.Check{
								"cumulative": knownvalue.Bool(true),
								"key":        knownvalue.StringExact("results"),
								"operation":  knownvalue.StringExact("count"),
							}),
						}),
					})),
				),
			},
			// Empty description
			{
				Config: testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName+"-updated", `
					position {
						x = 300
						y = 400
						h = 500
						w = 600
					}

					params {
						display = "line_chart"
						legend {
							groups = "charted"
						}
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(dashboardPanelName+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("position").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"x": knownvalue.Int64Exact(300),
						"y": knownvalue.Int64Exact(400),
						"h": knownvalue.Int64Exact(500),
						"w": knownvalue.Int64Exact(600),
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("display"), knownvalue.StringExact("line_chart")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact("")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("legend").AtSliceIndex(0).AtMapKey("groups"), knownvalue.StringExact("charted")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("datasets"), knownvalue.ListSizeExact(0)),
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

func TestAccResourceDashboardPanel_Import(t *testing.T) {
	resName := "rootly_dashboard_panel.test"
	dashboardName := acctest.RandomWithPrefix("tf-dashboard")
	dashboardPanelName := acctest.RandomWithPrefix("tf-dashboard-panel")

	configStateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(resName, tfjsonpath.New("id"), knownvalue.NotNull()),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName, `
					position {
						x = 3
						y = 4
						h = 5
						w = 6
					}

					params {
						display = "line_chart"
						description = "description"
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(dashboardPanelName)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("position").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"x": knownvalue.Int64Exact(3),
						"y": knownvalue.Int64Exact(4),
						"h": knownvalue.Int64Exact(5),
						"w": knownvalue.Int64Exact(6),
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("display"), knownvalue.StringExact("line_chart")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact("description")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("legend"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("datasets"), knownvalue.ListSizeExact(0)),
				),
			},
			{
				Config: testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName+"-updated", `
					position {
						x = 30
						y = 40
						h = 50
						w = 60
					}

					params {
						display = "line_chart"
						description = "description-updated"
					}
				`),
				ConfigStateChecks: append(
					configStateChecks,
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact(dashboardPanelName+"-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("position").AtSliceIndex(0), knownvalue.MapExact(map[string]knownvalue.Check{
						"x": knownvalue.Int64Exact(30),
						"y": knownvalue.Int64Exact(40),
						"h": knownvalue.Int64Exact(50),
						"w": knownvalue.Int64Exact(60),
					})),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("display"), knownvalue.StringExact("line_chart")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("description"), knownvalue.StringExact("description-updated")),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("legend"), knownvalue.ListSizeExact(0)),
					statecheck.ExpectKnownValue(resName, tfjsonpath.New("params").AtSliceIndex(0).AtMapKey("datasets"), knownvalue.ListSizeExact(0)),
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

func testAccResourceDashboardPanelConfig(dashboardName, dashboardPanelName, extra string) string {
	return testAccResourceDashboardConfig(dashboardName) + fmt.Sprintf(`
resource "rootly_dashboard_panel" "test" {
	dashboard_id = rootly_dashboard.test.id
	name         = "%s"
	%s
}
`, dashboardPanelName, extra)
}
