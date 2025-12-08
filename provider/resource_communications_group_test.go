package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCommunicationsGroup(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		IsUnitTest: false,
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCommunicationsGroupCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_group.test", "name", "TF test group"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "condition_type", "any"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.condition", "is"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.property_type", "service"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.properties.#", "3"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.0.properties.0.id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.0.properties.1.id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.0.properties.2.id"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.1.condition", "is"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.1.property_type", "group"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.1.properties.#", "1"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.1.properties.0.id"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.2.property_type", "severity"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.2.properties.#", "1"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.2.properties.0.id"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.3.property_type", "functionality"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.3.properties.#", "1"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.3.properties.0.id"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.4.property_type", "incident_type"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.4.properties.#", "1"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_conditions.4.properties.0.id"),

					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.#", "3"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.email", "test-bot+1@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.name", "test-bot+1"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.email", "test-bot+2@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.name", "test-bot+2"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.2.email", "test-bot+3@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.2.name", "test-bot+3"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_members.#", "4"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.0.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.1.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.2.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.3.user_id"),
				),
			},
			{
				Config: testAccResourceCommunicationsGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_group.test", "name", "TF test group (updated)"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "condition_type", "all"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.condition", "is"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.property_type", "service"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.properties.#", "2"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.properties.0.name", "TF test service 2"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.0.properties.1.name", "TF test service 3"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.1.condition", "is"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.1.property_type", "group"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.2.property_type", "severity"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.3.property_type", "functionality"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.4.property_type", "incident_type"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.#", "2"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.email", "test-bot+3@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.name", "test-bot+3"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.email", "test-bot+5@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.name", "test-bot+5"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_members.#", "3"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.0.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.1.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.2.user_id"),
				),
			},
			{
				Config: testAccResourceCommunicationsGroupUpdateRemoveConditions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rootly_communications_group.test", "name", "TF test group (updated removed conditions)"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "condition_type", "all"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_conditions.#", "0"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.#", "2"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.email", "test-bot+3@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.0.name", "test-bot+3"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.email", "test-bot+5@rootly.com"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_external_group_members.1.name", "test-bot+5"),
					resource.TestCheckResourceAttr("rootly_communications_group.test", "communication_group_members.#", "3"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.0.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.1.user_id"),
					resource.TestCheckResourceAttrSet("rootly_communications_group.test", "communication_group_members.2.user_id"),
				),
			},
		},
	})
}

const testAccResourceCommunicationsGroupCreate = `
	data rootly_user test1 {
		email = "bot-tftests+1@rootly.com"
	}

	data rootly_user test2 {
		email = "bot-tftests+2@rootly.com"
	}

	data rootly_user test3 {
		email = "bot-tftests+3@rootly.com"
	}

	data rootly_user test4 {
		email = "bot-tftests+4@rootly.com"
	}

	resource rootly_service test1 {
		name = "TF test service 1"
	}

	resource rootly_service test2 {
		name = "TF test service 2"
	}

	resource rootly_service test3 {
		name = "TF test service 3"
	}

	resource rootly_severity test1 {
		name = "TF test severity 1"
	}

	resource rootly_functionality test1 {
		name = "TF test functionality 1"
	}

	resource rootly_incident_type test1 {
		name = "TF test incident type 1"
	}

	resource rootly_team test {
		name = "TF test group"
	}

	resource rootly_communications_type test {
		name = "TF test type"
		color = "#FFFFFF"
	}

	resource rootly_communications_group test {
		name = "TF test group"
		communication_type_id = rootly_communications_type.test.id
		email_channel = true
		condition_type = "any"

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_service.test1.id
				name = rootly_service.test1.name
			}
			properties {
				id = rootly_service.test2.id
				name = rootly_service.test2.name
			}
			properties {
				id = rootly_service.test3.id
				name = rootly_service.test3.name
			}
			property_type = "service"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_team.test.id
				name = rootly_team.test.name
			}
			property_type = "group"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_severity.test1.id
				name = rootly_severity.test1.name
			}
			property_type = "severity"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_functionality.test1.id
				name = rootly_functionality.test1.name
			}
			property_type = "functionality"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_incident_type.test1.id
				name = rootly_incident_type.test1.name
			}
			property_type = "incident_type"
		}

		communication_external_group_members {
			email = "test-bot+1@rootly.com"
			name = "test-bot+1"
		}

		communication_external_group_members {
			email = "test-bot+2@rootly.com"
			name = "test-bot+2"
		}

		communication_external_group_members {
			email = "test-bot+3@rootly.com"
			name = "test-bot+3"
		}

		communication_group_members {
			user_id = data.rootly_user.test1.id
		}

		communication_group_members {
			user_id = data.rootly_user.test2.id
		}

		communication_group_members {
			user_id = data.rootly_user.test3.id
		}

		communication_group_members {
			user_id = data.rootly_user.test4.id
		}
	}
`

const testAccResourceCommunicationsGroupUpdate = `
	data rootly_user test1 {
		email = "bot-tftests+1@rootly.com"
	}

	data rootly_user test2 {
		email = "bot-tftests+2@rootly.com"
	}

	data rootly_user test3 {
		email = "bot-tftests+3@rootly.com"
	}

	data rootly_user test4 {
		email = "bot-tftests+4@rootly.com"
	}

	resource rootly_service test1 {
		name = "TF test service 1"
	}

	resource rootly_service test2 {
		name = "TF test service 2"
	}

	resource rootly_service test3 {
		name = "TF test service 3"
	}

	resource rootly_severity test1 {
		name = "TF test severity 1"
	}

	resource rootly_functionality test1 {
		name = "TF test functionality 1"
	}

	resource rootly_incident_type test1 {
		name = "TF test incident type 1"
	}

	resource rootly_team test {
		name = "TF test group"
	}

	resource rootly_communications_type test {
		name = "TF test type"
		color = "#FFFFFF"
	}

	resource rootly_communications_group test {
		name = "TF test group (updated)"
		communication_type_id = rootly_communications_type.test.id
		email_channel = true
		condition_type = "all"

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_service.test2.id
				name = rootly_service.test2.name
			}
			properties {
				id = rootly_service.test3.id
				name = rootly_service.test3.name
			}
			property_type = "service"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_team.test.id
				name = rootly_team.test.name
			}
			property_type = "group"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_severity.test1.id
				name = rootly_severity.test1.name
			}
			property_type = "severity"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_functionality.test1.id
				name = rootly_functionality.test1.name
			}
			property_type = "functionality"
		}

		communication_group_conditions {
			condition = "is"
			properties {
				id = rootly_incident_type.test1.id
				name = rootly_incident_type.test1.name
			}
			property_type = "incident_type"
		}

		communication_external_group_members {
			email = "test-bot+3@rootly.com"
			name = "test-bot+3"
		}

		communication_external_group_members {
			email = "test-bot+5@rootly.com"
			name = "test-bot+5"
		}

		communication_group_members {
			user_id = data.rootly_user.test1.id
		}

		communication_group_members {
			user_id = data.rootly_user.test3.id
		}

		communication_group_members {
			user_id = data.rootly_user.test4.id
		}
	}
`

const testAccResourceCommunicationsGroupUpdateRemoveConditions = `
	data rootly_user test1 {
		email = "bot-tftests+1@rootly.com"
	}

	data rootly_user test2 {
		email = "bot-tftests+2@rootly.com"
	}

	data rootly_user test3 {
		email = "bot-tftests+3@rootly.com"
	}

	data rootly_user test4 {
		email = "bot-tftests+4@rootly.com"
	}

	resource rootly_service test1 {
		name = "TF test service 1"
	}

	resource rootly_service test2 {
		name = "TF test service 2"
	}

	resource rootly_service test3 {
		name = "TF test service 3"
	}

	resource rootly_severity test1 {
		name = "TF test severity 1"
	}

	resource rootly_functionality test1 {
		name = "TF test functionality 1"
	}

	resource rootly_incident_type test1 {
		name = "TF test incident type 1"
	}

	resource rootly_team test {
		name = "TF test group"
	}

	resource rootly_communications_type test {
		name = "TF test type"
		color = "#FFFFFF"
	}

	resource rootly_communications_group test {
		name = "TF test group (updated removed conditions)"
		communication_type_id = rootly_communications_type.test.id
		email_channel = true
		condition_type = "all"

		communication_external_group_members {
			email = "test-bot+3@rootly.com"
			name = "test-bot+3"
		}

		communication_external_group_members {
			email = "test-bot+5@rootly.com"
			name = "test-bot+5"
		}

		communication_group_members {
			user_id = data.rootly_user.test1.id
		}

		communication_group_members {
			user_id = data.rootly_user.test3.id
		}

		communication_group_members {
			user_id = data.rootly_user.test4.id
		}
	}
`
