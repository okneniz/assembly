package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fp2 — op rd, rn (fneg/fcvt/fcvtzs/fcvtzu/scvtf/ucvtf/fmov).
type Fp2 struct {
	base

	op       string
	rd, rn   string
	enc      uint32
	rdK, rnK fpKind
}

// fpTypeBits — the fp type from bits [22:21]/[15:14]: 01 → s, 11 → d.
func fpTypeBits(v uint32) fpKind {
	if v&2 != 0 {
		return kD
	}

	return kS
}

func decodeFp2Of(op string, enc uint32, rdK, rnK fpKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		// the sf convention for widenable families: operand types are read
		// FROM THE WORD, not hard-coded by the ctor (fcvt/fcvtzs/fcvtzu —
		// their XML entries free the type bits; see isa_map widenAcceptable)
		switch op {
		case "fcvt":
			rnK = fpTypeBits(w >> 21 & 3)
			rdK = fpTypeBits(w >> 14 & 3)
		case "fcvtzs", "fcvtzu":
			rdK = kW
			if w>>31&1 == 1 {
				rdK = kX
			}

			rnK = fpTypeBits(w >> 21 & 3)
		case "scvtf", "ucvtf":
			rnK = kW
			if w>>31&1 == 1 {
				rnK = kX
			}

			rdK = fpTypeBits(w >> 21 & 3)
		}

		return Fp2{
			base: newBase(addr, w),
			op:   op,
			rd:   fpReg(w&0x1f, rdK),
			rn:   fpReg(w>>5&0x1f, rnK),
			enc:  enc,
			rdK:  rdK,
			rnK:  rnK,
		}
	}
}

func (i Fp2) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("%s %s, %s", i.op, i.rd, i.rn)
}

func (i Fp2) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	return writeWord(w, i.enc|rd|rn<<5)
}

func (i Fp2) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.op,
		i.ObjDump(disasm.DefaultViewCtx()),
		"FP",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
