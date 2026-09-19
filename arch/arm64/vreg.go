package arm64

import (
	"fmt"
	"strconv"
)

// VReg — the vector register operand v0..v31. The lane width is not a
// property of the register: it rides the instruction's arrangement
// (.8b/.16b/.4s/...), an operand of its own.
type VReg struct {
	num uint8 // 0..31
}

// newVReg - the register constructor: the struct is assembled only here.
func newVReg(num uint8) VReg {
	return VReg{
		num: num,
	}
}

// V — vector register v0..v31.
func V(n int) (VReg, error) {
	if n < 0 || n > 31 {
		return VReg{}, fmt.Errorf(
			"arm64.V: register number %d is out of 0..31", n,
		)
	}

	return newVReg(uint8(n)), nil
}

// VRegOf — a vector register by its source name (v0..v31): the inverse
// of (VReg).name.
func VRegOf(name string) (VReg, error) {
	if len(name) < 2 || name[0] != 'v' {
		return VReg{}, fmt.Errorf("arm64.VRegOf: unknown register %q", name)
	}

	n, err := strconv.Atoi(name[1:])
	if err != nil || n < 0 || n > 31 {
		return VReg{}, fmt.Errorf("arm64.VRegOf: unknown register %q", name)
	}

	return V(n)
}

// Num — the register number (0..31).
func (r VReg) Num() uint8 {
	return r.num
}

func (r VReg) String() string {
	return r.name()
}

// name - canonical name ("v5").
func (r VReg) name() string {
	return "v" + strconv.Itoa(int(r.num))
}
