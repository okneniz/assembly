package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DupElem — the vector-source Advanced SIMD copy forms: DUP (element)
// dup.Arr vd, vn[idx], INS (element) ins.sz vd[idx], vn[idx], and the
// scalar DUP alias mov.sz vd, vn (llvm prints mov.d/mov.s). imm5 is the
// size one-hot plus the destination lane index (imm5 >> size+1); INS
// carries the source lane in imm4 (bits [14:11]); the scalar form's
// imm5 is one-hot (size only - the index is always 0).
type DupElem struct {
	base

	op     string // dup | ins | mov (the scalar alias)
	size   uint32 // 0=b 1=h 2=s 3=d
	idx    uint32 // the dup source / ins destination lane
	srcIdx uint32 // the ins source lane
	q      uint32
	rd, rn string
	rdN    uint32
	rnN    uint32
}

// decodeSimdDupElem — DUP (element): 0x0E000400/0xBFE0FC00.
func decodeSimdDupElem(w uint32, addr uint64) Instr {
	imm5 := w >> 16 & 0x1f
	// the size is the lowest set bit (0=b..3=d); unencodable sizes go to
	// the .word fallback (the schema mask does not express this)
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w, addr)
	}

	size := uint32(bitsCtz(imm5))
	return DupElem{
		base: newBase(addr, w),
		op:   "dup",
		size: size,
		idx:  imm5 >> (size + 1),
		q:    w >> 30 & 1,
		rd:   vReg(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		rdN:  w & 0x1f,
		rnN:  w >> 5 & 0x1f,
	}
}

// decodeSimdInsElem — INS (element): 0x6E000400/0xFFE08400 (imm4 at
// bits [14:11] is the source lane).
func decodeSimdInsElem(w uint32, addr uint64) Instr {
	d, ok := decodeSimdDupElem(w, addr).(DupElem)
	if !ok {
		return decodeUnknown(w, addr) // the .word fallback
	}

	d.op = "ins"
	d.srcIdx = w >> 11 & 0xf
	return d
}

// decodeSimdDupScalar — the scalar DUP alias (llvm mov.d/mov.s):
// 0x5E000400/0xFFE0FC00; imm5 is one-hot (size only).
func decodeSimdDupScalar(w uint32, addr uint64) Instr {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || imm5 > 8 || imm5&(imm5-1) != 0 {
		return decodeUnknown(w, addr)
	}

	return DupElem{
		base: newBase(addr, w),
		op:   "mov",
		size: uint32(bitsCtz(imm5)),
		rd:   vReg(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		rdN:  w & 0x1f,
		rnN:  w >> 5 & 0x1f,
	}
}

func (i DupElem) ObjDump(_ disasm.ViewCtx) string {
	sz := [...]string{"b", "h", "s", "d"}[i.size]
	switch i.op {
	case "dup":
		return fmt.Sprintf("dup.%s %s, %s[%d]",
			decodeArrangement(i.q, i.size), i.rd, i.rn, i.idx)
	case "ins":
		return fmt.Sprintf("ins.%s %s[%d], %s[%d]", sz, i.rd, i.idx, i.rn, i.srcIdx)
	default:
		return fmt.Sprintf("mov.%s %s, %s", sz, i.rd, i.rn)
	}
}

func (i DupElem) Encode(w io.Writer, pc uint64) (int64, error) {
	maxIdx := uint32(16 >> i.size)
	var word uint32
	switch i.op {
	case "dup":
		if i.idx >= maxIdx {
			return 0, fmt.Errorf("dup: lane index %d out of range (0..%d)",
				i.idx, maxIdx-1)
		}

		word = 0x0E000400 | i.q<<30 | (1<<i.size|i.idx<<(i.size+1))<<16 |
			i.rnN<<5 | i.rdN
	case "ins":
		if i.idx >= maxIdx || i.srcIdx >= maxIdx {
			return 0, fmt.Errorf("ins: lane index out of range (0..%d)", maxIdx-1)
		}

		word = 0x6E000400 | (1<<i.size|i.idx<<(i.size+1))<<16 |
			i.srcIdx<<11 | i.rnN<<5 | i.rdN
	default:
		word = 0x5E000400 | 1<<i.size<<16 | i.rnN<<5 | i.rdN
	}

	return writeWord(w, word)
}

func (i DupElem) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"dup",
		i.ObjDump(disasm.DefaultViewCtx()),
		"SIMD copy",
		map[string]any{"op": i.op, "size": i.size},
	)
}
