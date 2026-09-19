package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fcmp — fcmp fn, fm | fcmp fn, #0.0 (Rm=0 encodes #0.0; sets the FP
// condition flags, no destination register).
type Fcmp struct {
	base

	rn, rm string
	withRM bool
}

// newFcmp - the Fcmp constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFcmp(b base, rn, rm FReg, withRM bool) (Fcmp, error) {
	if withRM {
		err := requireFpKind("Fcmp", rn, rm)
		if err != nil {
			return Fcmp{}, err
		}
	}

	return Fcmp{
		base:   b,
		rn:     rn.name(),
		rm:     rm.name(),
		withRM: withRM,
	}, nil
}

const (
	fcmpRD uint32 = 0x1E602000 // fcmp dn, dm
	fcmpRS uint32 = 0x1E202000 // fcmp sn, sm
	fcmp0D uint32 = 0x1E602008 // fcmp dn, #0.0
	fcmp0S uint32 = 0x1E202008 // fcmp sn, #0.0
)

func (i Fcmp) ObjDump(_ disasm.ViewCtx) string {
	if !i.withRM {
		return fmt.Sprintf("fcmp %s, #0.0", i.rn)
	}

	return fmt.Sprintf("fcmp %s, %s", i.rn, i.rm)
}

func (i Fcmp) Encode(w io.Writer) (int64, error) {
	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("fcmp: %w", err)
	}

	match, err := fpMatch(i.rn, fcmpRD, fcmpRS)
	if err != nil {
		return 0, fmt.Errorf("fcmp: %w", err)
	}

	if !i.withRM {
		match, err = fpMatch(i.rn, fcmp0D, fcmp0S)
		if err != nil {
			return 0, fmt.Errorf("fcmp: %w", err)
		}

		return writeWord(w, match|rn<<5)
	}

	rm, err := armRegNum(i.rm)
	if err != nil {
		return 0, fmt.Errorf("fcmp: %w", err)
	}

	return writeWord(w, match|rn<<5|rm<<16)
}

func (Builder) Fcmp(rn, rm FReg) (Instr, error) {
	return newFcmp(base{}, rn, rm, true)
}

// FcmpZero — fcmp fn, #0.0 (the immediate form: Rm=0).
func (Builder) FcmpZero(rn FReg) (Instr, error) {
	return newFcmp(base{}, rn, FReg{}, false)
}

func decodeFcmpReg(w uint32) (Instr, error) {
	in, err := newFcmp(
		newBase(w),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
		true,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcmpZero(w uint32) (Instr, error) {
	in, err := newFcmp(
		newBase(w),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		FReg{},
		false,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
