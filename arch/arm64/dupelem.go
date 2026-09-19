package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// DupElem — dup.Arr vd, vn[idx] (DUP element: every lane = the source
// lane).
type DupElem struct {
	base

	q, size, idx uint32
	rd, rn       string
}

// newDupElem - the DupElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newDupElem(b base, q, size, idx uint32, rd, rn VReg) (DupElem, error) {
	if err := requireElemSize("DupElem", size); err != nil {
		return DupElem{}, err
	}

	if err := requireLaneIdx("DupElem", size, idx); err != nil {
		return DupElem{}, err
	}

	return DupElem{
		base: b,
		q:    q,
		size: size,
		idx:  idx,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const dupElemEnc uint32 = 0x0E000400 // dup vd, vn[idx] (Q=0 form)

func (i DupElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("dup.%s %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.idx)
}

func (i DupElem) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("dup: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, dupElemEnc|i.q<<30|imm5<<16|rn<<5|rd)
}

func (Builder) DupElem(rd, rn VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewDupElem: %w", err)
	}

	return newDupElem(base{}, q, size, idx, rd, rn)
}

// decodeSimdDupElem — DUP (element): 0x0E000400/0xBFE0FC00.
func decodeSimdDupElem(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	// the size is the lowest set bit (0=b..3=d); unencodable sizes go to
	// the .word fallback (the schema mask does not express this)
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	in, err := newDupElem(
		newBase(w),
		w>>30&1,
		size,
		imm5>>(size+1),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// InsElem — ins.sz vd[idx], vn[idx] (INS element: copy one lane into
// another; the source lane rides imm4, bits [14:11]).
type InsElem struct {
	base

	size, idx, srcIdx uint32
	rd, rn            string
}

// newInsElem - the InsElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newInsElem(b base, size, idx, srcIdx uint32, rd, rn VReg) (InsElem, error) {
	if err := requireElemSize("InsElem", size); err != nil {
		return InsElem{}, err
	}

	if err := requireLaneIdx("InsElem", size, idx); err != nil {
		return InsElem{}, err
	}

	if err := requireLaneIdx("InsElem", size, srcIdx); err != nil {
		return InsElem{}, err
	}

	return InsElem{
		base:   b,
		size:   size,
		idx:    idx,
		srcIdx: srcIdx,
		rd:     rd.name(),
		rn:     rn.name(),
	}, nil
}

const insElemEnc uint32 = 0x6E000400 // ins vd[idx], vn[idx]

func (i InsElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ins.%s %s[%d], %s[%d]",
		elemName(i.size), i.rd, i.idx, i.rn, i.srcIdx)
}

func (i InsElem) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("ins: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, insElemEnc|imm5<<16|i.srcIdx<<i.size<<11|rn<<5|rd)
}

func (Builder) InsElem(rd, rn VReg, elem string, idx, srcIdx uint32) (Instr, error) {
	size, err := elemSize("InsElem", elem)
	if err != nil {
		return nil, err
	}

	return newInsElem(base{}, size, idx, srcIdx, rd, rn)
}

// decodeSimdInsElem — INS (element): 0x6E000400/0xFFE08400.
func decodeSimdInsElem(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	in, err := newInsElem(
		newBase(w),
		size,
		imm5>>(size+1),
		w>>11&0xf>>size,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// DupScalar — the scalar DUP alias (llvm prints mov.d/mov.s vd, vn):
// the bottom fp lane of vn replicated into vd's scalar; imm5 is
// one-hot (size only - the index is always 0).
type DupScalar struct {
	base

	size   uint32 // 2=s 3=d
	rd, rn string
}

// newDupScalar - the DupScalar constructor: validates the operands
// and assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newDupScalar(b base, size uint32, rd, rn VReg) (DupScalar, error) {
	if size < 2 || size > 3 {
		return DupScalar{}, fmt.Errorf(
			"arm64.NewDupScalar: only the .s and .d scalar forms exist",
		)
	}

	return DupScalar{
		base: b,
		size: size,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const dupScalarEnc uint32 = 0x5E000400 // mov vd, vn (the scalar DUP alias)

func (i DupScalar) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mov.%s %s, %s", elemName(i.size), i.rd, i.rn)
}

func (i DupScalar) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("mov: %w", err)
	}

	return writeWord(w, dupScalarEnc|1<<i.size<<16|rn<<5|rd)
}

func (Builder) DupScalar(rd, rn VReg, elem string) (Instr, error) {
	size, err := elemSize("DupScalar", elem)
	if err != nil {
		return nil, err
	}

	return newDupScalar(base{}, size, rd, rn)
}

// decodeSimdDupScalar — the scalar DUP alias: 0x5E000400/0xFFE0FC00;
// imm5 is one-hot (size only).
func decodeSimdDupScalar(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || imm5 > 8 || imm5&(imm5-1) != 0 {
		return decodeUnknown(w)
	}

	in, err := newDupScalar(
		newBase(w),
		uint32(bitsCtz(imm5)),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
