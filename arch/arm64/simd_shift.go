package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SimdShift — op.Arr vd, vn, #shift (immh:immb; ushr/sshr/sri invert it).
type SimdShift struct {
	base

	op         string
	rd, rn     string
	immh, immb uint32
	q          uint32
	enc        uint32
}

func decodeSimdShiftOf(op string, enc uint32) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return SimdShift{
			base: newBase(addr, w),
			op:   op,
			rd:   vReg(w & 0x1f),
			rn:   vReg(w >> 5 & 0x1f),
			immh: w >> 19 & 0xf,
			immb: w >> 16 & 7,
			q:    w >> 30 & 1,
			enc:  enc,
		}
	}
}

func (i SimdShift) ObjDump(_ disasm.ViewCtx) string {
	size, shift := simdShiftAmount(i.op, i.immh, i.immb)
	arr := decodeArrangement(i.q, size)
	return fmt.Sprintf("%s.%s %s, %s, #0x%x", i.op, arr, i.rd, i.rn, shift)
}

func (i SimdShift) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	if i.immh > 0xf || i.immb > 7 {
		return 0, fmt.Errorf("%s: imm out of range", i.op)
	}

	// immh carries size; the enc constants are in the Q=0 form
	return writeWord(w, i.enc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

func (i SimdShift) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.op,
		i.ObjDump(disasm.DefaultViewCtx()),
		"ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
