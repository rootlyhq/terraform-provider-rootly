package provider

var testAccUserDataSourceConfig = `
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}
`
