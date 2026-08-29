package loong64

import "strconv"

// imm - the value of an immediate operand: a concrete number. Symbolic
// slots are evaluated by the syntax layer (asm/loong64) BEFORE the
// structure is built - a computed instruction contains no holes. Branch
// and jump targets are absolute addresses.
type imm struct {
	val int64
}

// immNum - a concrete value (decoding and construction).
func immNum(v int64) imm {
	return imm{val: v}
}

// text - the %d output (for ObjDump of decoded instructions: LoongArch
// objdump prints immediates and branch targets in decimal).
func (m imm) text() string {
	return strconv.FormatInt(m.val, 10)
}
