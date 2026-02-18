package provider

const testAccDataSourceUserConfig = `
data "rootly_user" "test" {
	email = "bot+tftests@rootly.com"
}
`
