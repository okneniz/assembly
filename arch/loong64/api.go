package loong64

// Export surface for the layers ABOVE the arch: asm/loong64 (the syntax
// layer) and its pseudo - building computed instructions (BuildInstr +
// the op.go operands). No new shared structures - only aliases and
// one-line wrappers; fields stay unexported, instruction files are
// unchanged.

import (
	"fmt"
	"io"
)

// Base - the instruction bookkeeping record (for MarshalDTO of pseudo structures).
type Base = base

// BuildInstr builds a computed instruction from a mnemonic and operands
// (the asmCtors shape lookup + the per-instruction constructor).
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

// asmRegNum - the register names the source accepts: the canonical $-names
// (as the decoder prints them) and the $rN numbering.
var asmRegNum = buildAsmRegNum()

func buildAsmRegNum() map[string]uint8 {
	m := make(map[string]uint8, 64)
	for i, n := range laRegNames {
		m["$"+n] = uint8(i)
	}

	return m
}

// RegNumOf - the register number by name ($zero, $t0, $r21, ...).
func RegNumOf(name string) (uint32, error) {
	n, ok := asmRegNum[name]
	if !ok {
		if len(name) > 2 && name[0] == '$' && name[1] == 'r' {
			v := 0
			for _, c := range name[2:] {
				if c < '0' || c > '9' {
					return 0, fmt.Errorf("unknown register %q", name)
				}

				v = v*10 + int(c-'0')
			}

			if v <= 31 {
				return uint32(v), nil
			}
		}

		return 0, fmt.Errorf("unknown register %q", name)
	}

	return uint32(n), nil
}

// EncodingWord - the base encoding word from the generated table.
func EncodingWord(name string) uint32 {
	return loongEncodings[name][0]
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

// RegNames - canonical names of the integer registers, with the $ prefix
// (index = number).
func RegNames() [32]string {
	names := [32]string{}
	for i, n := range laRegNames {
		names[i] = "$" + n
	}

	return names
}

// MarshalDTO - the instruction JSON skeleton (as for real structures).
func MarshalDTO(b Base, mnem, full, group string, fields map[string]any) ([]byte, error) {
	return b.marshalDTO(mnem, full, group, fields)
}

// WriteWord - writing an encoded instruction (pseudo pairs).
func WriteWord(w io.Writer, word uint32) (int64, error) {
	return writeWord(w, word)
}
