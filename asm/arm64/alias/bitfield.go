package alias

// Aliases of the bitfield family: sxtb/sxth/sxtw (SBFM immr=0,
// imms=7/15/31, Rn is a W-register) and ubfiz/ubfx/sbfiz/sbfx
// (lsb+width -> immr/imms).

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// newSxt is the sxtb/sxth/sxtw alias: SBFM immr=0, imms=7/15/31; Rn
// is a W-register.
func newSxt(name string, imms uint32) arch.ArmCtor {
	return func(ops []arch.ArmOp) (arch.Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want rd, rn", name)
		}

		rd, err := arch.WantAReg(ops[0], name)
		if err != nil {
			return nil, err
		}

		rnW, err := arch.WantAReg(ops[1], name)
		if err != nil {
			return nil, err
		}

		if rnW[0] != 'w' {
			return nil, fmt.Errorf("%s: W-register expected", name)
		}

		num, err := arch.ArmRegNum(rnW)
		if err != nil {
			return nil, err
		}

		rn := arch.RegNameXOf(num) // the encoding uses the same number
		return arch.SbfmOf(rd, rn, 0, imms, true), nil
	}
}

// newBf is the ubfiz/ubfx/sbfiz/sbfx alias: lsb+width -> immr/imms.
func newBf(name string, isU, isFiz bool) arch.ArmCtor {
	return func(ops []arch.ArmOp) (arch.Instr, error) {
		if len(ops) != 4 {
			return nil, fmt.Errorf("%s: want rd, rn, #lsb, #width", name)
		}

		rd, rn, err := arch.ArmReg2(ops, name)
		if err != nil {
			return nil, err
		}

		lsbV := ops[2].Num()
		wV := ops[3].Num()

		isf := rd[0] == 'x'
		regsize := uint32(64)
		if !isf {
			regsize = 32
		}

		lsb, w := uint32(lsbV), uint32(wV)
		if lsb >= regsize || w == 0 || lsb+w > regsize {
			return nil, fmt.Errorf("%s: out of range", name)
		}

		var immr, imms uint32
		if isFiz {
			immr = (regsize - lsb) % regsize
			imms = w - 1
		} else {
			immr = lsb
			imms = lsb + w - 1
		}

		if isU {
			return arch.UbfmOf(rd, rn, immr, imms, isf), nil
		}

		return arch.SbfmOf(rd, rn, immr, imms, isf), nil
	}
}
