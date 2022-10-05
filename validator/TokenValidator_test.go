package validator

import (
	"fmt"
	"testing"

	scanner "gitee.com/piecat/text-scanner"
)

func TestValidator(t *testing.T) {
	v := Validator{}

	fmt.Printf("v.isToken([]rune(\"this is\"), scanner.Match{S: 2, E: 7}): %v\n", v.isToken([]rune("this is"), scanner.Match{S: 2, E: 7}))
}
