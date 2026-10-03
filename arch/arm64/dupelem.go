package arm64

import (
	"fmt"
	"io"
	"strconv"

	"github.com/okneniz/assembly/disasm"
)

// DupElem — dup.Arr vd, vn[idx] (DUP element: every lane = the source
// lane).
type DupElem struct {
	q, size, idx uint32
	rd, rn       string
}

// newDupElem - the DupElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newDupElem(q, size, idx uint32, rd, rn VReg) (DupElem, error) {
	if err := requireElemSize("DupElem", size); err != nil {
		return DupElem{}, err
	}

	if err := requireLaneIdx("DupElem", size, idx); err != nil {
		return DupElem{}, err
	}

	return DupElem{
		q:    q,
		size: size,
		idx:  idx,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const dupElemEnc uint32 = 0x0E000400 // dup vd, vn[idx] (Q=0 form)

func (i DupElem) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("dup: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, dupElemEnc|i.q<<30|imm5<<16|rn<<5|rd)
}

func (i DupElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("dup.%s %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.idx)
}

// InsElem — ins.sz vd[idx], vn[idx] (INS element: copy one lane into
// another; the source lane rides imm4, bits [14:11]).
type InsElem struct {
	size, idx, srcIdx uint32
	rd, rn            string
}

// newInsElem - the InsElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newInsElem(size, idx, srcIdx uint32, rd, rn VReg) (InsElem, error) {
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
		size:   size,
		idx:    idx,
		srcIdx: srcIdx,
		rd:     rd.name(),
		rn:     rn.name(),
	}, nil
}

const insElemEnc uint32 = 0x6E000400 // ins vd[idx], vn[idx]

func (i InsElem) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("ins: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, insElemEnc|imm5<<16|i.srcIdx<<i.size<<11|rn<<5|rd)
}

func (i InsElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ins.%s %s[%d], %s[%d]",
		elemName(i.size), i.rd, i.idx, i.rn, i.srcIdx)
}

// DupScalar — the scalar DUP alias (llvm prints mov b/h/s/d<n>,
// vn.sz[idx]): the bottom fp register view of one vn lane. Q is fixed 1
// (the 0x5E000400 class); the lane index rides imm5 above the size
// one-hot, the same layout as the element DUP.
type DupScalar struct {
	size, idx uint32 // 0=b 1=h 2=s 3=d
	rd, rn    string
}

// newDupScalar - the DupScalar constructor: validates the operands
// and assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newDupScalar(size, idx uint32, rd, rn VReg) (DupScalar, error) {
	if err := requireElemSize("DupScalar", size); err != nil {
		return DupScalar{}, err
	}

	if idx >= 16>>size {
		return DupScalar{}, fmt.Errorf(
			"arm64.NewDupScalar: lane index %d out of range (0..%d)",
			idx, 16>>size-1,
		)
	}

	return DupScalar{
		size: size,
		idx:  idx,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const dupScalarEnc uint32 = 0x5E000400 // mov <b|h|s|d>n, vn.sz[idx]

func (i DupScalar) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("mov: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, dupScalarEnc|imm5<<16|rn<<5|rd)
}

func (i DupScalar) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mov %s, %s.%s[%d]",
		scalarRegName(elemName(i.size), i.rd), i.rn, elemName(i.size), i.idx)
}

// scalarRegName — the scalar view name of a vector register ("v9" and
// ".s" → "s9"): the scalar views share the vector file.
func scalarRegName(letter, name string) string {
	num, err := armRegNum(name)
	if err != nil {
		return name
	}

	return letter + strconv.Itoa(int(num))
}
