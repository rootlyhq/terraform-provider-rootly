package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/rootlyhq/terraform-provider-rootly/v5/internal/acctest"
)

func TestAccResourceEscalationPath_UpgradeFromVersion(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	config := testAccResourceEscalationPathConfig(name, `
		initial_delay = 5
		time_restriction_time_zone = "America/New_York"
		time_restrictions {
			start_day = "monday"
			start_time = "17:00"
			end_day = "tuesday"
			end_time = "07:00"
		}
		time_restrictions {
			start_day = "tuesday"
			start_time = "17:00"
			end_day = "wednesday"
			end_time = "07:00"
		}
	`)
	checks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("initial_delay"), knownvalue.Int64Exact(5)),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restriction_time_zone"), knownvalue.StringExact("America/New_York")),
		statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restrictions"), knownvalue.ListExact([]knownvalue.Check{
			knownvalue.ObjectExact(map[string]knownvalue.Check{
				"start_day":  knownvalue.StringExact("monday"),
				"start_time": knownvalue.StringExact("17:00"),
				"end_day":    knownvalue.StringExact("tuesday"),
				"end_time":   knownvalue.StringExact("07:00"),
			}),
			knownvalue.ObjectExact(map[string]knownvalue.Check{
				"start_day":  knownvalue.StringExact("tuesday"),
				"start_time": knownvalue.StringExact("17:00"),
				"end_day":    knownvalue.StringExact("wednesday"),
				"end_time":   knownvalue.StringExact("07:00"),
			}),
		})),
	}

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"rootly": {
						Source:            "rootlyhq/rootly",
						VersionConstraint: "5.21.0",
					},
				},
				Config:            config,
				ConfigStateChecks: checks,
			},
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Config:                   config,
				ConfigStateChecks:        checks,
			},
		},
	})
}

func TestAccResourceEscalationPath_Basic(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPathConfig(name, `
					initial_delay = 5
					time_restriction_time_zone = "America/New_York"
					time_restrictions {
						start_day = "monday"
						start_time = "17:00"
						end_day = "tuesday"
						end_time = "07:00"
					}
					time_restrictions {
						start_day = "tuesday"
						start_time = "17:00"
						end_day = "wednesday"
						end_time = "07:00"
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("initial_delay"), knownvalue.Int64Exact(5)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restriction_time_zone"), knownvalue.StringExact("America/New_York")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restrictions"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_day":  knownvalue.StringExact("monday"),
							"start_time": knownvalue.StringExact("17:00"),
							"end_day":    knownvalue.StringExact("tuesday"),
							"end_time":   knownvalue.StringExact("07:00"),
						}),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_day":  knownvalue.StringExact("tuesday"),
							"start_time": knownvalue.StringExact("17:00"),
							"end_day":    knownvalue.StringExact("wednesday"),
							"end_time":   knownvalue.StringExact("07:00"),
						}),
					})),
				},
			},
			{
				Config: testAccResourceEscalationPathConfig(name+"-updated", `
					initial_delay = 0
					time_restriction_time_zone = "Europe/London"
					time_restrictions {
						start_day = "thursday"
						start_time = "09:00"
						end_day = "friday"
						end_time = "17:00"
					}

				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("initial_delay"), knownvalue.Int64Exact(0)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restriction_time_zone"), knownvalue.StringExact("Europe/London")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("time_restrictions"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"start_day":  knownvalue.StringExact("thursday"),
							"start_time": knownvalue.StringExact("09:00"),
							"end_day":    knownvalue.StringExact("friday"),
							"end_time":   knownvalue.StringExact("17:00"),
						}),
					})),
				},
			},
		},
	})
}

func TestAccResourceEscalationPath_Deferral(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPathConfig(name, `
					path_type               = "deferral"
					after_deferral_behavior = "re_evaluate"

					rules {
						rule_type = "deferral_window"
						time_zone = "America/Los_Angeles"
						time_blocks {
							monday     = true
							tuesday    = true
							start_time = "17:00"
							end_time   = "07:00"
						}
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name)),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("path_type"), knownvalue.StringExact("deferral")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("after_deferral_behavior"), knownvalue.StringExact("re_evaluate")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("rules"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type": knownvalue.StringExact("deferral_window"),
							"time_zone": knownvalue.StringExact("America/Los_Angeles"),
							"time_blocks": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"monday":     knownvalue.Bool(true),
									"tuesday":    knownvalue.Bool(true),
									"start_time": knownvalue.StringExact("17:00"),
									"end_time":   knownvalue.StringExact("07:00"),
								}),
							}),
						}),
					})),
				},
			},
			{
				Config: testAccResourceEscalationPathConfig(name+"-updated", `
					path_type               = "deferral"
					after_deferral_behavior = "re_evaluate"

					rules {
						rule_type = "deferral_window"
						time_zone = "America/New_York"
						time_blocks {
							monday     = true
							wednesday  = true
							friday     = true
							start_time = "18:00"
							end_time   = "08:00"
						}
						time_blocks {
							saturday   = true
							sunday     = true
							all_day    = true
						}
					}
				`),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-updated")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("path_type"), knownvalue.StringExact("deferral")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("after_deferral_behavior"), knownvalue.StringExact("re_evaluate")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("rules"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type": knownvalue.StringExact("deferral_window"),
							"time_zone": knownvalue.StringExact("America/New_York"),
							"time_blocks": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"monday":     knownvalue.Bool(true),
									"wednesday":  knownvalue.Bool(true),
									"friday":     knownvalue.Bool(true),
									"start_time": knownvalue.StringExact("18:00"),
									"end_time":   knownvalue.StringExact("08:00"),
								}),
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"saturday": knownvalue.Bool(true),
									"sunday":   knownvalue.Bool(true),
									"all_day":  knownvalue.Bool(true),
								}),
							}),
						}),
					})),
				},
			},
		},
	})
}

func TestAccResourceEscalationPath_DeferralExecutePath(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "rootly_escalation_policy" "exec_test" {
						name = "%[1]s-ep"
					}

					resource "rootly_escalation_path" "target" {
						name                 = "%[1]s-target-path"
						default              = false
						escalation_policy_id = rootly_escalation_policy.exec_test.id
						path_type            = "escalation"

						rules {
							rule_type           = "working_hour"
							within_working_hour = true
						}
					}

					resource "rootly_escalation_path" "test" {
						name                    = "%[1]s-deferral-exec"
						default                 = false
						escalation_policy_id    = rootly_escalation_policy.exec_test.id
						path_type               = "deferral"
						after_deferral_behavior = "execute_path"
						after_deferral_path_id  = rootly_escalation_path.target.id

						rules {
							rule_type = "deferral_window"
							time_zone = "Etc/UTC"
							time_blocks {
								monday     = true
								tuesday    = true
								wednesday  = true
								thursday   = true
								friday     = true
								start_time = "22:00"
								end_time   = "06:00"
							}
						}
					}
				`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-deferral-exec")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("path_type"), knownvalue.StringExact("deferral")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("after_deferral_behavior"), knownvalue.StringExact("execute_path")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("after_deferral_path_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("rules"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type": knownvalue.StringExact("deferral_window"),
							"time_zone": knownvalue.StringExact("Etc/UTC"),
							"time_blocks": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"monday":     knownvalue.Bool(true),
									"tuesday":    knownvalue.Bool(true),
									"wednesday":  knownvalue.Bool(true),
									"thursday":   knownvalue.Bool(true),
									"friday":     knownvalue.Bool(true),
									"start_time": knownvalue.StringExact("22:00"),
									"end_time":   knownvalue.StringExact("06:00"),
								}),
							}),
						}),
					})),
				},
			},
		},
	})
}

func TestAccResourceEscalationPath_ServiceRule(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "rootly_escalation_policy" "exec_test" {
						name = "%[1]s-ep"
					}

					resource "rootly_service" "test" {
						name = "%[1]s-svc"
					}

					resource "rootly_escalation_path" "test" {
						name                 = "%[1]s-target-path"
						default              = false
						escalation_policy_id = rootly_escalation_policy.exec_test.id
						match_mode           = "match-any-rule"

						rules {
							rule_type   = "service"
							service_ids = [rootly_service.test.id]
						}
					}
				`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-target-path")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("match_mode"), knownvalue.StringExact("match-any-rule")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("rules"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type":   knownvalue.StringExact("service"),
							"service_ids": knownvalue.SetSizeExact(1),
						}),
					})),
				},
			},
		},
	})
}

func TestAccResourceEscalationPath_AllRuleTypes(t *testing.T) {
	addr := "rootly_escalation_path.test"
	name := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "rootly_escalation_policy" "test_rules" {
						depends_on = [rootly_alert_urgency.test, rootly_alert_field.test]

						name = "%[1]s-ep"
					}

					resource "rootly_alert_urgency" "test" {
						name        = "%[1]s-urgency"
						description = "Test urgency for escalation path rules"
					}

					resource "rootly_alert_field" "test" {
						name = "%[1]s-alert-field"
					}

					resource "rootly_escalation_path" "test" {
						name                 = "%[1]s-path"
						default              = false
						escalation_policy_id = rootly_escalation_policy.test_rules.id
						match_mode           = "match-any-rule"

						rules {
							rule_type   = "alert_urgency"
							urgency_ids = [rootly_alert_urgency.test.id]
						}

						rules {
							rule_type           = "working_hour"
							within_working_hour = true
						}

						rules {
							rule_type = "json_path"
							json_path = "$.severity"
							operator  = "is"
							value     = "critical"
						}

						rules {
							rule_type      = "field"
							fieldable_type = "AlertField"
							fieldable_id   = rootly_alert_field.test.id
							operator       = "is_one_of"
							values         = ["value1", "value2"]
						}
					}
				`, name),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("name"), knownvalue.StringExact(name+"-path")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("match_mode"), knownvalue.StringExact("match-any-rule")),
					statecheck.ExpectKnownValue(addr, tfjsonpath.New("rules"), knownvalue.ListExact([]knownvalue.Check{
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type":   knownvalue.StringExact("alert_urgency"),
							"urgency_ids": knownvalue.SetSizeExact(1),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type":           knownvalue.StringExact("working_hour"),
							"within_working_hour": knownvalue.Bool(true),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type": knownvalue.StringExact("json_path"),
							"json_path": knownvalue.StringExact("$.severity"),
							"operator":  knownvalue.StringExact("is"),
							"value":     knownvalue.StringExact("critical"),
						}),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"rule_type":      knownvalue.StringExact("field"),
							"fieldable_type": knownvalue.StringExact("AlertField"),
							"fieldable_id":   knownvalue.NotNull(),
							"operator":       knownvalue.StringExact("is_one_of"),
							"values": knownvalue.SetExact([]knownvalue.Check{
								knownvalue.StringExact("value1"),
								knownvalue.StringExact("value2"),
							}),
						}),
					})),
				},
			},
		},
	})
}

func testAccResourceEscalationPathConfig(name, extra string) string {
	return fmt.Sprintf(`
resource "rootly_escalation_policy" "test" {
	name = "%[1]s-ep"
}

resource "rootly_escalation_path" "test" {
	escalation_policy_id = rootly_escalation_policy.test.id

	name = "%[1]s"
	%[2]s
}
`, name, extra)
}
