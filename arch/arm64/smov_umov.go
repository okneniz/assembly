package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// elemSize - the element width number by its letter (b/h/s/d).
func elemSize(instr, elem string) (uint32, error) {
	switch elem {
	case "b":
		return 0, nil
	case "h":
		return 1, nil
	case "s":
		return 2, nil
	case "d":
		return 3, nil
	}

	return 0, fmt.Errorf("arm64.New%s: element size %q is not one of b/h/s/d",
		instr, elem,
	)
}

// elemName - the element letter by its width number.
func elemName(size uint32) string {
	return [...]string{"b", "h", "s", "d"}[size]
}

// requireElemSize - a sane element width number (decode-side values).
func requireElemSize(instr string, size uint32) error {
	if size > 3 {
		return fmt.Errorf("arm64.New%s: element size %d is out of 0..3", instr, size)
	}

	return nil
}

// Smov — smov wd|xd, vn.sz[idx] (sign-extend one lane into the
// integer register; .d elements do not exist).
type Smov struct {
	base

	size, idx uint32
	q         uint32 // an x destination
	vd, gpr   string
}

// newSmov - the Smov constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSmov(b base, q, size, idx uint32, vd VReg, gpr Reg) (Smov, error) {
	if err := requireElemSize("Smov", size); err != nil {
		return Smov{}, err
	}

	if err := requireLaneIdx("Smov", size, idx); err != nil {
		return Smov{}, err
	}

	if size == 3 {
		return Smov{}, fmt.Errorf("arm64.NewSmov: .d elements are not allowed")
	}

	if err := requireGprClass(gpr, "Smov", "wd"); err != nil {
		return Smov{}, err
	}

	qv := uint32(0)
	if gpr.Is64() {
		qv = 1
	}

	return Smov{
		base: b,
		size: size,
		idx:  idx,
		q:    qv,
		vd:   vd.name(),
		gpr:  gpr.name(),
	}, nil
}

const smovEnc uint32 = 0x0E002C00 // smov wd, vn.sz[idx] (op bits 10)

func (i Smov) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("smov %s, %s.%s[%d]", i.gpr, i.vd, elemName(i.size), i.idx)
}

func (i Smov) Encode(w io.Writer) (int64, error) {
	vd, err := armRegNum(i.vd)
	if err != nil {
		return 0, fmt.Errorf("smov: %w", err)
	}

	gpr, err := armRegNum(i.gpr)
	if err != nil {
		return 0, fmt.Errorf("smov: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, smovEnc|i.q<<30|imm5<<16|vd<<5|gpr)
}

func (Builder) Smov(wd Reg, vn VReg, elem string, idx uint32) (Instr, error) {
	size, err := elemSize("Smov", elem)
	if err != nil {
		return nil, err
	}

	return newSmov(base{}, 0, size, idx, vn, wd)
}

// Umov — umov wd|xd, vn.sz[idx] (move one lane into the integer
// register). llvm prints the form that fills the whole register
// (.s into w, .d into x) as the mov alias.
type Umov struct {
	base

	size, idx uint32
	q         uint32 // .d into x
	vd, gpr   string
}

// newUmov - the Umov constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newUmov(b base, q, size, idx uint32, vd VReg, gpr Reg) (Umov, error) {
	if err := requireElemSize("Umov", size); err != nil {
		return Umov{}, err
	}

	if err := requireLaneIdx("Umov", size, idx); err != nil {
		return Umov{}, err
	}

	if err := requireGprClass(gpr, "Umov", "wd"); err != nil {
		return Umov{}, err
	}

	// the element must fill the destination: umov takes .d only into x
	// and .b/.h/.s only into w (llvm: "invalid operand" otherwise)
	if (size == 3) != gpr.Is64() {
		return Umov{}, fmt.Errorf(
			"arm64.NewUmov: %s destination expected for .%s elements",
			map[bool]string{true: "x", false: "w"}[size == 3], elemName(size),
		)
	}

	qv := uint32(0)
	if gpr.Is64() {
		qv = 1
	}

	return Umov{
		base: b,
		size: size,
		idx:  idx,
		q:    qv,
		vd:   vd.name(),
		gpr:  gpr.name(),
	}, nil
}

const umovEnc uint32 = 0x0E003C00 // umov wd, vn.sz[idx] (op bits 11)

func (i Umov) ObjDump(_ disasm.ViewCtx) string {
	sz := elemName(i.size)
	if (i.size == 2 && i.q == 0) || (i.size == 3 && i.q == 1) {
		return fmt.Sprintf("mov %s, %s.%s[%d]", i.gpr, i.vd, sz, i.idx)
	}

	return fmt.Sprintf("umov %s, %s.%s[%d]", i.gpr, i.vd, sz, i.idx)
}

func (i Umov) Encode(w io.Writer) (int64, error) {
	vd, err := armRegNum(i.vd)
	if err != nil {
		return 0, fmt.Errorf("umov: %w", err)
	}

	gpr, err := armRegNum(i.gpr)
	if err != nil {
		return 0, fmt.Errorf("umov: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, umovEnc|i.q<<30|imm5<<16|vd<<5|gpr)
}

func (Builder) Umov(wd Reg, vn VReg, elem string, idx uint32) (Instr, error) {
	size, err := elemSize("Umov", elem)
	if err != nil {
		return nil, err
	}

	return newUmov(base{}, 0, size, idx, vn, wd)
}
