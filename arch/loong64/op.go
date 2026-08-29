package loong64

// Computed operands (from the syntax layer): a register number or a
// number. Value slots are evaluated ABOVE (asm/loong64); a computed
// instruction contains no holes. Registers travel as numbers (a byte),
// not names.

type opKind int

const (
	opRegK opKind = iota
	opNumK
)

// Op - an instruction operand: a register (by number) or a number.
type Op struct {
	kind opKind
	reg  uint8
	num  int64
}

// OpReg - a register operand by number 0..31.
func OpReg(num uint8) Op {
	return Op{
		kind: opRegK,
		reg:  num,
	}
}

// OpNum - a numeric operand (an immediate or an absolute target).
func OpNum(v int64) Op {
	return Op{
		kind: opNumK,
		num:  v,
	}
}
