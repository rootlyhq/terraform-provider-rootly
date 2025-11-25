package acctest

import (
	sdkv2_acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func RandInt() int {
	return sdkv2_acctest.RandInt()
}

func RandomWithPrefix(name string) string {
	return sdkv2_acctest.RandomWithPrefix(name)
}
