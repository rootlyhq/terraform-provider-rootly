package acctest

import (
	pluginacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func RandInt() int {
	return pluginacctest.RandInt()
}

func RandomWithPrefix(name string) string {
	return pluginacctest.RandomWithPrefix(name)
}
