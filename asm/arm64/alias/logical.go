package alias

// Aliases of the logical family: tst (ands with Rd = zr), mvn (orn
// with Rn = zr), mov rd, rm (orr Rn=zr, byte-parity priority as in
// GAS/clang) | mov rd, #imm (movz/movn per the legacy hw layout; when
// not MOVZ/MOVN - an ORR bitmask of the xzr form, as objdump prints).

import (
	"errors"
	"fmt"
	"strings"

	arch "github.com/okneniz/assembly/arch/arm64"
	arm64 "github.com/okneniz/assembly/asm/arm64"
)

// isDigits — s is a non-empty digit run (the register number of a
// scalar view name like s9/d31).
func isDigits(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// newTst is the tst rn, #bitmask | rn, rm{, shift} alias: ands with
// Rd = zr.
func newTst(ops []arch.ArmOp) (arch.Instr, error) {
	if len(ops) != 2 && len(ops) != 3 {
		return nil, errors.New("tst: want rn, op")
	}

	rn, err := arch.WantAReg(ops[0], "tst")
	if err != nil {
		return nil, err
	}

	if !arm64.IsGPR(rn) {
		return nil, errors.New("tst: integer register expected")
	}

	zr := arch.ZeroReg(rn)
	if ops[1].IsImm() {
		v := ops[1].Num()
		is64 := rn[0] == 'x'
		n, immr, imms, ok := arch.EncodeBitMasks(is64, uint64(v))
		if !ok {
			return nil, errors.New("tst: not encodable as bitmask")
		}

		return arch.AndsImmOf(zr, rn, immr, imms, n == 1, is64)
	}

	if ops[1].IsReg() && arm64.IsGPR(ops[1].Reg()) {
		shift, amt := "lsl", uint32(0)
		if len(ops) == 3 {
			if !ops[2].IsShift() {
				return nil, errors.New("tst: shift modifier expected")
			}

			a := arch.ShiftAmt(ops[2])
			if a < 0 || a > 63 {
				return nil, errors.New("tst: bad shift")
			}

			shift, amt = ops[2].ShiftName(), uint32(a)
		}

		return arch.AndsShiftOf(zr, rn, ops[1].Reg(), amt, shift, rn[0] == 'x')
	}

	return nil, errors.New("tst: bad operand")
}

// newMvn is the mvn rd, rm{, shift} alias: orn with Rn = zr.
func newMvn(ops []arch.ArmOp) (arch.Instr, error) {
	if len(ops) != 2 && len(ops) != 3 {
		return nil, errors.New("mvn: want rd, rm")
	}

	rd, rm, err := arch.ArmReg2(ops, "mvn")
	if err != nil {
		return nil, err
	}

	if !arm64.IsGPR(rd) || !arm64.IsGPR(rm) {
		return nil, errors.New("mvn: integer registers expected")
	}

	shift, amt := "lsl", uint32(0)
	if len(ops) == 3 {
		if !ops[2].IsShift() {
			return nil, errors.New("mvn: shift modifier expected")
		}

		a := arch.ShiftAmt(ops[2])
		if a < 0 || a > 63 {
			return nil, errors.New("mvn: bad shift")
		}

		shift, amt = ops[2].ShiftName(), uint32(a)
	}

	return arch.OrnShiftOf(rd, arch.ZeroReg(rd), rm, amt, shift, rd[0] == 'x')
}

// newMov is the mov rd, rm (orr Rn=zr, byte-parity priority) |
// mov rd, #imm (movz/movn per the legacy hw layout) alias.
func newMov(ops []arch.ArmOp) (arch.Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("mov: want rd, op")
	}

	rd, err := arch.WantAReg(ops[0], "mov")
	if err != nil {
		return nil, err
	}

	// mov <b|h|s|d><n>, vn.<e>[idx]: the scalar DUP alias (the decoder
	// prints this form). The destination view class must match the
	// source element letter.
	if len(rd) >= 2 && strings.ContainsRune("bhsd", rune(rd[0])) &&
		isDigits(rd[1:]) && ops[1].IsReg() && ops[1].LaneIdx() &&
		ops[1].Arr() == string(rd[0]) {
		idx := ops[1].Num()
		if idx < 0 {
			return nil, errors.New("mov: bad lane index")
		}

		rdV, err := arch.VRegOf("v" + rd[1:])
		if err != nil {
			return nil, fmt.Errorf("mov: %w", err)
		}

		rnV, err := arch.VRegOf(ops[1].Reg())
		if err != nil {
			return nil, fmt.Errorf("mov: %w", err)
		}

		in, ierr := (arch.Builder{}).DupScalar(rdV, rnV, string(rd[0]), uint32(idx))
		if ierr != nil {
			return nil, fmt.Errorf("mov: %w", ierr)
		}

		return in, nil
	}

	if !arm64.IsGPR(rd) {
		return nil, errors.New("mov: integer register expected")
	}

	// mov wd, vn.s[n] | mov xd, vn.d[n]: the LLVM input alias of UMOV
	// (the element fills the register - the decoder prints this form).
	// Only the filling sizes have the alias (mov wd, vn.b[n] is not a
	// spelling; the canonical input there is umov).
	if ops[1].IsReg() && ops[1].LaneIdx() &&
		(ops[1].Arr() == "s" && rd[0] != 'x' || ops[1].Arr() == "d" && rd[0] == 'x') {
		vd := ops[1].Reg()
		if vd == "" || vd[0] != 'v' {
			return nil, errors.New("mov: vector register expected")
		}

		elem := "s"
		if rd[0] == 'x' {
			elem = "d"
		}

		idx := ops[1].Num()
		if idx < 0 {
			return nil, errors.New("mov: bad lane index")
		}

		gpr, err := arch.RegOf(rd)
		if err != nil {
			return nil, fmt.Errorf("mov: %w", err)
		}

		vreg, err := arch.VRegOf(vd)
		if err != nil {
			return nil, fmt.Errorf("mov: %w", err)
		}

		in, ierr := (arch.Builder{}).Umov(gpr, vreg, elem, uint32(idx))
		if ierr != nil {
			return nil, fmt.Errorf("mov: %w", ierr)
		}

		return in, nil
	}

	if ops[1].IsReg() && arm64.IsGPR(ops[1].Reg()) {
		return arch.OrrShiftOf(rd, arch.ZeroReg(rd), ops[1].Reg(), 0, "", rd[0] == 'x')
	}

	if !ops[1].IsImm() || ops[1].Sym() != "" {
		return nil, errors.New("mov: register or immediate expected")
	}

	v := ops[1].Num()
	is64 := rd[0] == 'x'
	if v >= 0 {
		if v == 0 {
			return arch.MovzOf(rd, 0, 0)
		}

		u := uint64(v)
		for hw := range uint32(4) {
			if !is64 && hw > 1 {
				break
			}

			lane := u >> (16 * hw)
			if lane <= 0xffff && u == lane<<(16*hw) {
				return arch.MovzOf(rd, uint32(lane), hw)
			}
		}

		// MOVN-encodable positive: v == ~(imm16 << 16hw) (high bits
		// set) - objdump prints mov, not orr
		for hw := range uint32(4) {
			if !is64 && hw > 1 {
				break
			}

			lane := ^u >> (16 * hw)
			if lane <= 0xffff && ^u == lane<<(16*hw) && lane != 0 {
				return arch.MovnOf(rd, uint32(lane), hw)
			}
		}

		// not MOVZ/MOVN - an ORR bitmask (the xzr form of mov), as
		// objdump does
		if n, immr, imms, ok := arch.EncodeBitMasks(is64, u); ok {
			return arch.OrrImmOf(rd, arch.ZeroReg(rd), immr, imms, n == 1, is64)
		}

		return nil, fmt.Errorf("mov: %#x not encodable", v)
	}

	// negatives: MOVN = ~(imm16 << 16hw), so x = -v is imm16<<16hw + 1
	// (x ≡ 1 mod the lane); hw 2/3 exist only in the 64-bit form.
	x := uint64(-v)
	switch {
	case x <= 0x10000:
		return arch.MovnOf(rd, uint32(x-1), 0)
	case x <= 0xffff0001 && x%0x10000 == 1:
		return arch.MovnOf(rd, uint32((x-1)/0x10000), 1)
	case x <= 0xffff00000001 && x%0x100000000 == 1:
		return arch.MovnOf(rd, uint32((x-1)/0x100000000), 2)
	case x <= 0xffff000000000001 && x%0x1000000000000 == 1:
		return arch.MovnOf(rd, uint32((x-1)/0x1000000000000), 3)
	}

	// a negative not MOVN-encodable - an ORR bitmask (the xzr form of mov);
	// v as unsigned is ^(x-1).
	if n, immr, imms, ok := arch.EncodeBitMasks(is64, ^(x - 1)); ok {
		return arch.OrrImmOf(rd, arch.ZeroReg(rd), immr, imms, n == 1, is64)
	}

	return nil, errors.New("mov: negative imm not encodable")
}
