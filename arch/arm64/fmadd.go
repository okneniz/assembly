package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmadd — fmadd fd, fn, fm, fa (fd = fa + fn*fm; double/single by the
// operand kind, the type rides bit 22).
type Fmadd struct {
	rd, rn, rm, ra string
}

// newFmadd - the Fmadd constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmadd(rd, rn, rm, ra FReg) (Fmadd, error) {
	err := requireFpKind("Fmadd", rd, rn, rm, ra)
	if err != nil {
		return Fmadd{}, err
	}

	return Fmadd{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
		ra: ra.name(),
	}, nil
}

const (
	fmaddD uint32 = 0x1F400000 // fmadd dd, dn, dm, da
	fmaddS uint32 = 0x1F000000 // fmadd sd, sn, sm, sa
)

func (i Fmadd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmadd %s, %s, %s, %s", i.rd, i.rn, i.rm, i.ra)
}

func (i Fmadd) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fmaddD, fmaddS)
	if err != nil {
		return 0, fmt.Errorf("fmadd: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmadd: %w", err)
	}

	ra, err := armRegNum(i.ra)
	if err != nil {
		return 0, fmt.Errorf("fmadd: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|ra<<10|rm<<16)
}

// Fnmsub — fnmsub fd, fn, fm, fa (fd = -(fn*fm - fa); double/single
// by the operand kind, the type rides bit 22).
type Fnmsub struct {
	rd, rn, rm, ra string
}

// newFnmsub - the Fnmsub constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFnmsub(rd, rn, rm, ra FReg) (Fnmsub, error) {
	err := requireFpKind("Fnmsub", rd, rn, rm, ra)
	if err != nil {
		return Fnmsub{}, err
	}

	return Fnmsub{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
		ra: ra.name(),
	}, nil
}

const (
	fnmsubD uint32 = 0x1F608000 // fnmsub dd, dn, dm, da
	fnmsubS uint32 = 0x1F208000 // fnmsub sd, sn, sm, sa
)

func (i Fnmsub) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fnmsub %s, %s, %s, %s", i.rd, i.rn, i.rm, i.ra)
}

func (i Fnmsub) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fnmsubD, fnmsubS)
	if err != nil {
		return 0, fmt.Errorf("fnmsub: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fnmsub: %w", err)
	}

	ra, err := armRegNum(i.ra)
	if err != nil {
		return 0, fmt.Errorf("fnmsub: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|ra<<10|rm<<16)
}
