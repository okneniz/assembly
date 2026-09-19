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

// arrBits — the Q/size bits of an arrangement string (the inverse of
// decodeArrangement; the assembler-side twin of asm arrQSize).
func arrBits(arr string) (q, size uint32, err error) {
	switch arr {
	case "8b":
		return 0, 0, nil
	case "16b":
		return 1, 0, nil
	case "4h":
		return 0, 1, nil
	case "8h":
		return 1, 1, nil
	case "2s":
		return 0, 2, nil
	case "4s":
		return 1, 2, nil
	case "2d":
		return 1, 3, nil
	}

	return 0, 0, fmt.Errorf("unknown arrangement %q", arr)
}

// bitsCtz - the number of the lowest set bit.
func bitsCtz(v uint32) int {
	n := 0
	for v&1 == 0 && v != 0 {
		n++
		v >>= 1
	}

	return n
}
