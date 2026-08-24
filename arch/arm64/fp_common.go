package arm64

// fpKind — the register kind of an FP operand.
type fpKind uint8

const (
	kS fpKind = iota // s0
	kD               // d0
	kW               // w0
	kX               // x0
)

// fpReg — the register name by fpKind.
func fpReg(n uint32, k fpKind) string {
	switch k {
	case kS:
		return fpRegNameS(n)
	case kD:
		return fpRegNameD(n)
	case kW:
		return regNameW(n)
	default:
		return regNameX(n)
	}
}
