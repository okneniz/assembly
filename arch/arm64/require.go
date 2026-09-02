package arm64

// Contextual checks for instruction constructors: constraints that bind
// several operands together (register classes and widths, offset
// alignment per access size) - everything a single role type cannot
// check on its own.

import (
	"fmt"
	"slices"
)

// requireClass - register of an allowed class, otherwise an error with a hint.
func requireClass(r Reg, instr, op, hint string, ok ...regClass) error {
	if slices.Contains(ok, r.class) {
		return nil
	}

	return fmt.Errorf("arm64.New%s: operand %s: %s does not fit - %s", instr, op, r.name(), hint)
}

// requireWidth - equal width of all registers of an instruction.
func requireWidth(instr string, regs ...Reg) error {
	for _, r := range regs[1:] {
		if r.Is64() != regs[0].Is64() {
			return fmt.Errorf("arm64.New%s: register widths must match: %s vs %s",
				instr, regs[0].name(), r.name())
		}
	}

	return nil
}

// requireHwW - the 32-bit form of movz/movk never has shift #32/#48.
func requireHwW(rd Reg, instr string, hw Hw) error {
	if !rd.Is64() && hw > Hw1 {
		return fmt.Errorf("arm64.New%s: the 32-bit form allows Hw0/Hw1, not %s", instr, hw)
	}

	return nil
}

// requireShift - constraints of the shifted form of add/sub: only
// lsl/lsr/asr (ror is a logical-family encoding, unallocated for add/sub);
// the 32-bit form limits the shift amount to 0..31 (imm6 >= 32 is
// unpredictable).
func requireShift(rd Reg, instr string, imm Imm6, sh Shift) error {
	if sh == ROR {
		return fmt.Errorf("arm64.New%s: ror is not allowed for add/sub (only lsl/lsr/asr)", instr)
	}

	if !rd.Is64() && imm.v > 31 {
		return fmt.Errorf(
			"arm64.New%s: the 32-bit form allows shift 0..31, not #%d",
			instr,
			imm.v,
		)
	}

	return nil
}

// requireOff - offset must be non-negative, aligned, and within imm12.
func requireOff(instr string, off Off, scale uint32) error {
	if off < 0 || uint64(off)&(uint64(1)<<scale-1) != 0 || uint64(off)>>scale > 0xfff {
		return fmt.Errorf(
			"arm64.New%s: offset %s is out of the form's range (0..0x%x, alignment %d)",
			instr,
			off,
			uint64(0xfff)<<scale,
			uint64(1)<<scale,
		)
	}

	return nil
}

// lsOperand - rt: x/w register (sp not allowed), rn: x/sp base.
func lsOperand(rt, rn Reg, instr string) error {
	if err := requireClass(rt, instr, "rt", "x/w register (register 31 in rt reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return err
	}

	return requireClass(rn, instr, "rn", "x register or SP (register 31 in the base reads as sp)",
		classX, classSP)
}

// requirePairOff — the pair imm7 constraint: the offset is scaled to the
// access size, aligned, and within the signed imm7 range.
func requirePairOff(instr string, off Off, scale uint32) error {
	if off&(Off(1)<<scale-1) != 0 || off>>scale < -64 || off>>scale > 63 {
		return fmt.Errorf(
			"arm64.New%s: offset %s is out of the form's range (-0x%x..0x%x, alignment %d)",
			instr,
			off,
			64<<scale,
			63<<scale,
			1<<scale,
		)
	}

	return nil
}

// requireUnscaledOff — the ldur/stur family constraint: the offset is a
// signed imm9, any alignment.
func requireUnscaledOff(instr string, off Off) error {
	if off < -256 || off > 255 {
		return fmt.Errorf(
			"arm64.New%s: offset %s is out of the unscaled form's range (-0x100..0xff)",
			instr,
			off,
		)
	}

	return nil
}
