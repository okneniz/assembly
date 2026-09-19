package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ins — INS general: insert a GPR value into one lane. llvm prints it
// as the mov alias - mov.sz vd[idx], wn - and that spelling is the
// only accepted input form (see the asm ctors).
type Ins struct {
	base

	size, idx uint32 // 0=b 1=h 2=s 3=d
	vd, gpr   string
}

// newIns - the Ins constructor: validates the operands and assembles
// the struct (the Builder method delegates here; the decoder calls it
// with values read from the word).
func newIns(b base, size, idx uint32, vd VReg, gpr Reg) (Ins, error) {
	if err := requireElemSize("Ins", size); err != nil {
		return Ins{}, err
	}

	if err := requireLaneIdx("Ins", size, idx); err != nil {
		return Ins{}, err
	}

	if err := requireGprClass(gpr, "Ins", "wn"); err != nil {
		return Ins{}, err
	}

	if (size == 3) != gpr.Is64() {
		return Ins{}, fmt.Errorf(
			"arm64.NewIns: register widths must match: .%s lane vs %s",
			elemName(size), gpr.name(),
		)
	}

	return Ins{
		base: b,
		size: size,
		idx:  idx,
		vd:   vd.name(),
		gpr:  gpr.name(),
	}, nil
}

const insEnc uint32 = 0x4E001C00 // ins vd[idx], wn (op bits 01, Q=1)

func (i Ins) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mov.%s %s[%d], %s", elemName(i.size), i.vd, i.idx, i.gpr)
}

func (i Ins) Encode(w io.Writer) (int64, error) {
	vd, err := armRegNum(i.vd)
	if err != nil {
		return 0, fmt.Errorf("mov: %w", err)
	}

	gpr, err := armRegNum(i.gpr)
	if err != nil {
		return 0, fmt.Errorf("mov: %w", err)
	}

	imm5 := 1<<i.size | i.idx<<(i.size+1)
	return writeWord(w, insEnc|imm5<<16|gpr<<5|vd)
}

func (Builder) Ins(vd VReg, idx uint32, wn Reg, elem string) (Instr, error) {
	size, err := elemSize("Ins", elem)
	if err != nil {
		return nil, err
	}

	return newIns(base{}, size, idx, vd, wn)
}
