package riscv

// Op - a computed operand for instruction construction: a register, a
// number, or memory off(reg). The syntax layer (asm/riscv) evaluates the
// expression slots (expressions, rounding-mode keywords, CSR names) and
// arrives here with numbers; BuildInstr dispatches by mnemonic to the
// new* constructors.
// The decode path does not use these types - it builds structures directly
// from the word.

// opKind - the kind of a construction operand.
type opKind int

const (
	opRegK opKind = iota
	opNumK
	opMemK
)

// Op - a construction operand: exactly one of the kinds.
type Op struct {
	kind opKind
	reg  string
	num  int64
	mem  Mem
}

// Mem - a memory operand off(reg).
type Mem struct {
	base string
	off  int64
}

// OpReg - a register by canonical name (ABI or xN/fN).
func OpReg(name string) Op {
	return Op{kind: opRegK, reg: name}
}

// OpNum - a numeric operand.
func OpNum(v int64) Op {
	return Op{kind: opNumK, num: v}
}

// OpMem - memory off(reg).
func OpMem(base string, off int64) Op {
	return Op{kind: opMemK, mem: Mem{base: base, off: off}}
}
