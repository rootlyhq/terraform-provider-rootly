package fwtypes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ExampleIsKnown() {
	// example for unknown string
	unknownString := types.StringUnknown()
	fmt.Println("IsKnown(unknownString):", IsKnown(unknownString))

	// example for known string
	knownString := types.StringValue("known")
	fmt.Println("IsKnown(knownString):", IsKnown(knownString))

	// example for null string
	nullString := types.StringNull()
	fmt.Println("IsKnown(nullString):", IsKnown(nullString))

	// Output:
	// IsKnown(unknownString): false
	// IsKnown(knownString): true
	// IsKnown(nullString): false
}
