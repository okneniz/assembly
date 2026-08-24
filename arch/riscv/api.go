package riscv

// Export surface for the layers ABOVE the arch: asm/riscv (the syntax layer)
// and its pseudo - building computed instructions (BuildInstr + the op.go
// operands) and encoding primitives. No new shared structures - only
// aliases and one-line wrappers; fields stay unexported, instruction
// files are unchanged.

import (
	"fmt"
	"io"
)

// Base - the instruction bookkeeping record (for MarshalDTO of pseudo structures).
type Base = base

// BuildInstr builds a computed instruction from a mnemonic and operands
// (asmCtors lookup + constructor). Errors - in encoding-path form.
func BuildInstr(mnem string, ops []Op) (Instr, error) {
	ctor, ok := asmCtors[mnem]
	if !ok {
		return nil, fmt.Errorf("unknown instruction %q", mnem)
	}

	st, err := ctor(ops)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mnem, err)
	}

	return st, nil
}

// RegNumOf - register number by name (ABI or xN/fN); fp - the expected
// register class (operand validation by the syntax layer).
func RegNumOf(name string, fp bool) (uint32, error) {
	return rvRegNumOf(name, fp)
}

// EncodingWord - the base encoding word from the generated table.
func EncodingWord(name string) uint32 {
	return riscvEncodings[name][0]
}

// Encoding primitives (wrappers).
func EncU(imm int64) (uint32, error) {
	return encU(imm)
}
func EncI(imm int64) (uint32, error) {
	return encI(imm)
}
func EncJ(off int64) (uint32, error) {
	return encJ(off)
}
func Fits12(v int64) bool {
	return fits12(v)
}
func PcrelHiLo(rel int64) (hi, lo int64) {
	return pcrelHiLo(rel)
}
func RegBits(name string) uint32 {
	return regBits(name)
}

func WriteWord(w io.Writer, word uint32) (int64, error) {
	return writeWord(w, word)
}
func WriteHalf(w io.Writer, h uint16) (int64, error) {
	return writeHalf(w, h)
}

// MarshalDTO - the instruction JSON skeleton (as for real structures).
func MarshalDTO(b Base, mnem, full, group string, fields map[string]any) ([]byte, error) {
	return b.marshalDTO(mnem, full, group, fields)
}

// Mnemonics - the mnemonics of the decode table (the assembler grammar
// dictionary; the pseudo layer adds its own on top).
func Mnemonics() []string {
	names := make([]string, 0, len(decodeTable))
	for _, e := range decodeTable {
		names = append(names, e.name)
	}

	return names
}

// RegNames - canonical ABI names of integer registers (index = number).
func RegNames() [32]string {
	return rvRegNames
}

// FpRegNames - canonical names of floating registers (index = number).
func FpRegNames() [32]string {
	return rvFRegNames
}
