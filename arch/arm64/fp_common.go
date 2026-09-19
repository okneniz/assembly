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

// fpMatch — the encoding of an FP form by the kind of the stored
// register name (d/s), the FP twin of sfMatch. The s form is 0 only in
// families without one (there are none today, but the shape stays
// symmetric with sfMatch).
func fpMatch(rd string, matchD, matchS uint32) (uint32, error) {
	if rd[0] == 'd' {
		return matchD, nil
	}

	return matchS, nil
}

// fpGprMatch — the encoding of an int↔FP conversion by the kinds of
// the stored register names (d/x, d/w, s/x, s/w): the cross-file twin
// of sfMatch.
func fpGprMatch(fp, gpr string, dx, dw, sx, sw uint32) uint32 {
	if gpr[0] == 'w' {
		if fp[0] == 's' {
			return sw
		}

		return dw
	}

	if fp[0] == 's' {
		return sx
	}

	return dx
}
