package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Simd3 — op.Arr vd, vn, vm; the mov alias: and/orr with Rn == Rm -> mov.16b/8b.
type Simd3 struct {
	base

	op         string
	rd, rn, rm string
	enc        uint32
	q, size    uint32
}

func decodeSimd3Of(op string, enc uint32) func(w uint32, addr uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Simd3{
			base: newBase(addr, w),
			op:   op,
			rd:   vReg(w & 0x1f),
			rn:   vReg(w >> 5 & 0x1f),
			rm:   vReg(w >> 16 & 0x1f),
			enc:  enc,
			q:    w >> 30 & 1,
			size: w >> 22 & 3,
		}
	}
}

// simd3Logical - the three-same logical group: with fixed bits 0x1C00
// (and bit 21) the mnemonic is selected by the pair of U (bit 29) + size
// (bits 23:22): one table row = one U-group of 4 instructions.
var simd3Logical = [2][4]string{
	{"and", "bic", "orr", "orn"}, // U = 0
	{"eor", "bsl", "bit", "bif"}, // U = 1
}

func isSimd3Logical(op string) bool {
	for _, row := range simd3Logical {
		for _, o := range row {
			if o == op {
				return true
			}
		}
	}

	return false
}

// decodeSimd3Logical - the opcode and encoding are read FROM THE WORD: the
// group's family schemas (and/eor) do not pin the size bits and each covers
// 4 instructions, so the mnemonic cannot be baked in as an argument. The
// arrangement of logical operations depends only on Q (8b/16b), size is
// part of the opcode.
func decodeSimd3Logical(w uint32, addr uint64) Instr {
	return Simd3{
		base: newBase(addr, w),
		op:   simd3Logical[w>>29&1][w>>22&3],
		rd:   vReg(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		rm:   vReg(w >> 16 & 0x1f),
		enc:  w &^ 0x001F1F1F, // all encoding bits except Rd/Rn/Rm (Q included)
		q:    w >> 30 & 1,
		size: w >> 22 & 3,
	}
}

func (i Simd3) ObjDump(_ disasm.ViewCtx) string {
	if (i.op == "and" || i.op == "orr") && i.rn == i.rm {
		arr := "16b"
		if i.q == 0 {
			arr = "8b"
		}

		return fmt.Sprintf("mov.%s %s, %s", arr, i.rd, i.rm)
	}

	arr := decodeArrangement(i.q, i.size3())
	if isSimd3Logical(i.op) {
		// the logical group's size bits are part of the opcode, not the arrangement
		if i.q == 1 {
			arr = "16b"
		} else {
			arr = "8b"
		}
	}

	return fmt.Sprintf("%s.%s %s, %s, %s", i.op, arr, i.rd, i.rn, i.rm)
}

func (i Simd3) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", i.op, err)
	}

	// q/size are mandatory: the enc constants are registered in the .8b form
	// (Q=0), and in the logical group bits 23:22 are the opcode, not the
	// arrangement (the OR is idempotent)
	return writeWord(w, i.enc|i.q<<30|i.size<<22|rd|rn<<5|rm<<16)
}

func (i Simd3) MarshalJSON() ([]byte, error) {
	return i.marshal(i.op, i.ObjDump(disasm.DefaultViewCtx()), "ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm})
}

func (i Simd3) size3() uint32 {
	return i.raw >> 22 & 3
}
