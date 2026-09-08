package arm64

import "fmt"

// Computed construction operands (vOp/vMem) and the constructors' input —
// resolvedInstr; packFields assembles the word by schema (inverse
// transforms), op*Of — operand access by position. The grammar, the
// symbolic armOp model, resolveOps and encodeARM with self-verify live in
// asm/arm64 (the syntax layer above the arch); numbers arrive here already.

type armOpKind int

const (
	armOpReg armOpKind = iota
	armOpImm           // expression (#... or bare)
	armOpLit           // =expression (literal pool, GAS: ldr xN, =literal)
	armOpMem
	armOpList
	armOpShift  // lsl/lsr/asr/ror #imm — a modifier of the previous operand
	armOpExtend // uxtw/sxtw...[#imm]
	armOpFloat  // #1.50000000 (fmov imm)
)

// armListReg — an element of a register list ({ v0.16b, x0, x1 }).
type armListReg struct {
	reg, arr string
}

// vOp — a computed construction operand (a mirror of the syntax layer's
// uncomputed armOp: expression slots turned into numbers, name operand — sym).
type vOp struct {
	kind    armOpKind
	reg     string
	arr     string
	num     int64
	sym     string
	fval    float64
	mem     *vMem
	list    []armListReg
	shift   string
	laneIdx bool // reg has a suffix index v30[1]
	hasAmt  bool // shift/extend has an amount
}

// vMem — a computed memory operand.
type vMem struct {
	base    string
	off     int64
	offReg  string
	opt     string
	optAmt  int64
	pre     bool
	post    int64
	hasOff  bool
	hasPost bool
	hasOpt  bool
}

// isKeyword — a name operand from the keyword tables (conditions, system
// registers, prfm hints): it stays a name instead of being computed by the
// resolver.
func isKeyword(sym string) bool {
	if _, err := condNum(sym); err == nil {
		return true
	}

	if invSysRegChecked(sym) != 0 {
		return true
	}

	return sym == "pldl1keep"
}

// offBits — a pc-relative byte offset → signed imm bits shifted right by 2.
func offBits(off int64, bits int) (uint32, error) {
	if off%4 != 0 {
		return 0, fmt.Errorf("branch offset %#x is not 4-byte aligned", off)
	}

	off >>= 2
	lim := int64(1) << uint(bits-1)
	if off < -lim || off >= lim {
		return 0, fmt.Errorf("branch offset %#x out of ±%d range", off, lim*4)
	}

	return uint32(off) & ((1 << uint(bits)) - 1), nil
}
