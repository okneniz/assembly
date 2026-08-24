package arm64

// Common helpers of ASIMD representations: v register names and reglists.

import (
	"fmt"
	"strings"
)

// vReg — the v register name.
func vReg(n uint32) string {
	return fmt.Sprintf("v%d", n)
}

// regListStr — "{ v0, v1, ... }" from the first number and count (for the ld1 family).
func regListStr(first uint32, count int) string {
	s := fmt.Sprintf("v%d", first)
	var sSb13 strings.Builder
	for k := 1; k < count; k++ {
		fmt.Fprintf(&sSb13, ", v%d", first+uint32(k))
	}

	s += sSb13.String()
	return s
}
