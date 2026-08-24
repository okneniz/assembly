package riscv

// Encoding (the reverse direction of the decoder): register name -> number
// (the inverse of rvRegNames/rvFRegNames + accepting x0-x31/f0-f31), the
// reverse scatter of split immediates (I/S/B/U/J), assembling the word from
// {match} + fields, and the RVC compressor (32-bit form -> 16-bit when the
// constraints allow - GNU as semantics).
// All compression decisions depend only on the operands and the instruction
// address, so Size and Encode take the same inputs and never diverge.

import "fmt"

// asmReg - a register from source: ABI name or xN/fN.
type asmReg struct {
	num uint32
	fp  bool
}

func newAsmReg(num uint32, fp bool) asmReg {
	return asmReg{
		num: num,
		fp:  fp,
	}
}

// asmRegNum - the name -> register table (int + FP, both notations); consumers -
// rvc_bits (r3ok/cr3/r5), the compress branches of Encode, and rvRegNumOf.
var asmRegNum = buildAsmRegNum()

func buildAsmRegNum() map[string]asmReg {
	m := map[string]asmReg{}
	for i, n := range rvRegNames {
		m[n] = newAsmReg(uint32(i), false)
		m[fmt.Sprintf("x%d", i)] = newAsmReg(uint32(i), false)
	}

	for i, n := range rvFRegNames {
		m[n] = newAsmReg(uint32(i), true)
		m[fmt.Sprintf("f%d", i)] = newAsmReg(uint32(i), true)
	}

	return m
}

// rvRegNumOf returns the register number; fp=false for an integer operand.
func rvRegNumOf(name string, wantFP bool) (uint32, error) {
	r, ok := asmRegNum[name]
	if !ok {
		return 0, fmt.Errorf("unknown register %q", name)
	}

	if r.fp != wantFP {
		kind := "integer"
		if wantFP {
			kind = "floating-point"
		}

		return 0, fmt.Errorf("register %q is not %s", name, kind)
	}

	return r.num, nil
}

// --- reverse immediate scatter (the mirror of iImm/sImm/bImm/uImm/jImm) ---

func encI(imm int64) (uint32, error) {
	if imm < -2048 || imm > 2047 {
		return 0, fmt.Errorf("immediate %d does not fit in 12 signed bits", imm)
	}

	return uint32(imm&0xfff) << 20, nil
}

func encS(imm int64) (uint32, error) {
	if imm < -2048 || imm > 2047 {
		return 0, fmt.Errorf("store offset %d does not fit in 12 signed bits", imm)
	}

	v := uint32(imm & 0xfff)
	return (v>>5)<<25 | (v&0x1f)<<7, nil
}

func encB(off int64) (uint32, error) {
	if off%2 != 0 {
		return 0, fmt.Errorf("branch offset %d is not even", off)
	}

	if off < -4096 || off > 4094 {
		return 0, fmt.Errorf("branch offset %d out of ±4KB range", off)
	}

	v := uint32(off & 0x1fff)
	return ((v>>12)&1)<<31 |
		((v>>5)&0x3f)<<25 |
		((v>>1)&0xf)<<8 |
		((v>>11)&1)<<7, nil
}

func encU(imm int64) (uint32, error) {
	// Accept either a 20-bit imm20 value (as our decoder prints it),
	// or a full constant divisible by 0x1000 (as written in GAS sources:
	// lui a0, 0x12345000).
	if imm >= 0 && imm <= 0xfffff {
		return uint32(imm) << 12, nil
	}

	if imm%0x1000 == 0 && imm >= -(1<<31) && imm < (1<<32) {
		return uint32(imm) & 0xfffff000, nil
	}

	return 0, fmt.Errorf("u-type immediate %#x out of range", imm)
}

func encJ(off int64) (uint32, error) {
	if off%2 != 0 {
		return 0, fmt.Errorf("jump offset %d is not even", off)
	}

	if off < -(1<<20) || off >= (1<<20) {
		return 0, fmt.Errorf("jump offset %d out of ±1MB range", off)
	}

	v := uint32(off & 0x1fffff)
	return ((v>>20)&1)<<31 |
		((v>>1)&0x3ff)<<21 |
		((v>>11)&1)<<20 |
		((v>>12)&0xff)<<12, nil
}

// pcrelHiLo splits a PC-relative offset into an auipc+addi pair:
// hi = (rel+0x800)>>12 (signed), lo = rel - (hi<<12) - always fits int12.
func pcrelHiLo(rel int64) (hi, lo int64) {
	hi = (rel + 0x800) >> 12
	lo = rel - (hi << 12)
	return
}

// fits12 is true for the 12-bit signed range.
func fits12(v int64) bool {
	return v >= -2048 && v <= 2047
}

// fits6 - the 6-bit signed range of CI immediates (c.addi/c.li/...).
func fits6(v int64) bool {
	return v >= -32 && v <= 31
}

// rvCSRNum - CSR name -> address (the inverse of rvCsrNames).
var rvCSRNum = buildRVCSRNum()

func buildRVCSRNum() map[string]uint32 {
	m := map[string]uint32{}
	for addr, name := range rvCsrNames {
		if _, exists := m[name]; !exists {
			m[name] = addr
		}
	}

	return m
}

// CSRNumOf - the CSR number by name (the inverse of rvCsrNames; the syntax
// layer converts CSR names to numbers before building the instruction).
func CSRNumOf(name string) (uint32, bool) {
	v, ok := rvCSRNum[name]
	return v, ok
}
