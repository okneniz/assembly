package arm64

// Bitfield-family assembler constructors: lsl/lsr/asr #imm (UBFM/SBFM),
// ror #imm (EXTR with rn==rm). Priorities follow armAliasMap: the
// UBFM/SBFM imm form is primary, the register form (the LSLV encoding) is
// the fallback of the newer RorReg/LslReg constructor. The aliases
// sxtb/sxth/sxtw and ubfiz/ubfx/sbfiz/sbfx (lsb+width → immr/imms) —
// asm/arm64/alias.

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// makeUbfmLSL — the UBFM encoding of lsl rd, rn, #sh: immr = -sh mod
// regsize, imms = regsize-1 (scalable width).
func makeUbfmLSL(rd, rn string, sh uint32, isf bool) (arch.Ubfm, bool) {
	regsize := uint32(64)
	if !isf {
		regsize = 32
	}

	if sh == 0 || sh >= regsize {
		return arch.Ubfm{}, false
	}

	u, err := UbfmOfC(rd, rn, (regsize-sh)%regsize, regsize-1, isf)
	if err != nil {
		return arch.Ubfm{}, false
	}

	return u, true
}

// newLslArm — lsl rd, rn, #sh → UBFM (primary); reg form → LslReg;
// anything unfit — legacy candidates.
func newLslArm(ops []vOp) (Instr, error) {
	if len(ops) == 3 && ops[2].Kind() == arch.ArmOpReg {
		return newLslReg(ops)
	}

	rd, rn, sh, err := rri(ops, "lsl")
	if err != nil {
		return nil, err
	}

	if u, ok := makeUbfmLSL(rd, rn, sh, rd[0] == 'x'); ok {
		return u, nil
	}

	return nil, errors.New("lsl: imm out of range — legacy path")
}

// newLsrArm — lsr rd, rn, #sh → UBFM immr=sh, imms=regsize-1.
func newLsrArm(ops []vOp) (Instr, error) {
	if len(ops) == 3 && ops[2].Kind() == arch.ArmOpReg {
		return newLsrReg(ops)
	}

	rd, rn, sh, err := rri(ops, "lsr")
	if err != nil {
		return nil, err
	}

	isf := rd[0] == 'x'
	regsize := uint32(64)
	if !isf {
		regsize = 32
	}

	if sh == 0 || sh >= regsize {
		return nil, errors.New("lsr: imm out of range")
	}

	return UbfmOf(rd, rn, sh, regsize-1, isf)
}

// newAsrArm — asr rd, rn, #sh → SBFM immr=sh, imms=regsize-1.
func newAsrArm(ops []vOp) (Instr, error) {
	if len(ops) == 3 && ops[2].Kind() == arch.ArmOpReg {
		return newAsrReg(ops)
	}

	rd, rn, sh, err := rri(ops, "asr")
	if err != nil {
		return nil, err
	}

	isf := rd[0] == 'x'
	regsize := uint32(64)
	if !isf {
		regsize = 32
	}

	if sh == 0 || sh >= regsize {
		return nil, errors.New("asr: imm out of range")
	}

	return SbfmOf(rd, rn, sh, regsize-1, isf)
}

// rri — rd, rn, #imm parsing.
func rri(ops []vOp, name string) (string, string, uint32, error) {
	if len(ops) != 3 || ops[2].Kind() != arch.ArmOpImm || ops[2].Sym() != "" {
		return "", "", 0, fmt.Errorf("%s: want rd, rn, #imm", name)
	}

	rd, rn, err := armReg2(ops, name)
	if err != nil {
		return "", "", 0, err
	}

	v := ops[2].Num()
	if v < 0 || v > 63 {
		return "", "", 0, fmt.Errorf("%s: bad imm", name)
	}

	return rd, rn, uint32(v), nil
}

// newRorArm2 — ror rd, rn, #imm → EXTR rn, rn.
func newRorArm2(ops []vOp) (Instr, error) {
	if len(ops) == 3 && ops[2].Kind() == arch.ArmOpReg {
		return newRorReg(ops)
	}

	rd, rn, sh, err := rri(ops, "ror")
	if err != nil {
		return nil, err
	}

	if sh > 63 {
		return nil, errors.New("ror: imm out of range")
	}

	return ExtrOf(rd, rn, rn, sh)
}
