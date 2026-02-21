package provider

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/rootlyhq/terraform-provider-rootly/v2/internal/must"
)

func TestAccResourceEscalationPolicy(t *testing.T) {
	resName := "rootly_escalation_policy.test"
	escalationPolicyName := acctest.RandomWithPrefix("tf-ep")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEscalationPolicy(testAccResourceEscalationPolicyConfigData{
					Name:        escalationPolicyName,
					Description: "test description",
					Extras: `
						business_hours {
							time_zone  = "Europe/London"
							days       = ["M"]
							start_time = "09:00"
							end_time   = "17:00"
						}
					`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", escalationPolicyName),
					resource.TestCheckResourceAttr(resName, "description", "test description"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.time_zone", "Europe/London"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.days.#", "1"),
					resource.TestCheckTypeSetElemAttr(resName, "business_hours.0.days.*", "M"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.start_time", "09:00"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.end_time", "17:00"),
				),
			},
			{
				Config: testAccResourceEscalationPolicy(testAccResourceEscalationPolicyConfigData{
					Name:        escalationPolicyName + "-updated",
					Description: "test updated description",
					Extras: `
						business_hours {
							time_zone  = "Pacific Time (US & Canada)"
							days       = ["T"]
							start_time = "10:00"
							end_time   = "18:00"
						}
					`,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", escalationPolicyName+"-updated"),
					resource.TestCheckResourceAttr(resName, "description", "test updated description"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.time_zone", "Pacific Time (US & Canada)"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.days.#", "1"),
					resource.TestCheckTypeSetElemAttr(resName, "business_hours.0.days.*", "T"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.start_time", "10:00"),
					resource.TestCheckResourceAttr(resName, "business_hours.0.end_time", "18:00"),
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

var testAccResourceEscalationPolicyConfigTemplate = template.Must(template.New("config").Parse(`
resource "rootly_escalation_policy" "test" {
	name = "{{ .Name }}"
	description = "{{ .Description }}"
	{{ .Extras }}
}
`))

type testAccResourceEscalationPolicyConfigData struct {
	Name        string
	Description string
	Extras      string
}

func testAccResourceEscalationPolicy(data testAccResourceEscalationPolicyConfigData) string {
	var buf bytes.Buffer
	must.Do(testAccResourceEscalationPolicyConfigTemplate.Execute(&buf, data))
	return buf.String()
}
