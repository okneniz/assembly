package arm64

import (
	"fmt"
	"math"
)

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

// resolvedInstr — an instruction with computed operands: the input of the
// arch's constructors and legacy handlers.
type resolvedInstr struct {
	mnem string
	ops  []vOp
}

// contextTransforms — transforms that require addr (branches, literal
// offsets): their inverses are not in the registry, format handlers put
// ready bits into the fields, and packFields passes them through the empty
// transform.
var contextTransforms = map[string]bool{
	"brOff26": true, "brOff19": true, "brOff14": true,
	"sext19": true, "sext9": true,
	"lsPairImm7_64": true, "lsPairImm7_32": true,
}

// packFields assembles the word: Schema.Value | fields (inverse transforms).
func packFields(s *Schema, fields map[string]any) (uint32, error) {
	w := s.Value
	for _, f := range s.Fields {
		v, ok := fields[f.Name]
		if !ok {
			// field not set by a handler — must be part of Value
			continue
		}

		var bits uint32
		var err error
		if contextTransforms[f.Transform] {
			bits, err = applyInverseTransform("", v) // value already = bits
			if err == nil && f.Width < 32 {
				bits &= (1 << f.Width) - 1 // signed raw bits (imm7=-2 → 0x7e)
			}
		} else {
			bits, err = applyInverseTransform(f.Transform, v)
		}

		if err != nil {
			return 0, fmt.Errorf("field %s: %w", f.Name, err)
		}

		if f.Width < 32 && bits>>f.Width != 0 {
			return 0, fmt.Errorf("field %s: %#x does not fit in %d bits", f.Name, bits, f.Width)
		}

		w |= bits << f.Offset
	}

	return w, nil
}

// opRegOf — the register operand at position i (or an error).
func opRegOf(in resolvedInstr, i int) (string, error) {
	if i >= len(in.ops) {
		return "", fmt.Errorf("operand %d: register expected", i+1)
	}

	op := in.ops[i]
	if op.kind != armOpReg {
		return "", fmt.Errorf("operand %d: register expected", i+1)
	}

	return op.reg, nil
}

// opImmOf — the number operand at position i (the value is already computed).
func opImmOf(in resolvedInstr, i int) (int64, error) {
	if i >= len(in.ops) {
		return 0, fmt.Errorf("operand %d: immediate expected", i+1)
	}

	op := in.ops[i]
	if op.kind == armOpImm {
		return op.num, nil
	}

	if op.kind == armOpFloat {
		return int64(math.Float64bits(op.fval)), nil
	}

	return 0, fmt.Errorf("operand %d: immediate expected", i+1)
}

// opMemOf — the memory operand at position i.
func opMemOf(in resolvedInstr, i int) (*vMem, error) {
	if i >= len(in.ops) || in.ops[i].kind != armOpMem {
		return nil, fmt.Errorf("operand %d: memory operand expected", i+1)
	}

	return in.ops[i].mem, nil
}

// opShiftOf — a trailing lsl/... modifier (after position from).
func opShiftOf(in resolvedInstr, from int) (name string, amt int64, ok bool) {
	for i := from; i < len(in.ops); i++ {
		if in.ops[i].kind == armOpShift {
			return in.ops[i].shift, in.ops[i].num, true
		}
	}

	return "", 0, false
}

// brBits — the target's absolute address → signed imm bits shifted right by 2.
func brBits(target, addr int64, bits int) (uint32, error) {
	off := target - addr
	if off%4 != 0 {
		return 0, fmt.Errorf("branch target %#x is not 4-byte aligned from %#x", target, addr)
	}

	off >>= 2
	lim := int64(1) << uint(bits-1)
	if off < -lim || off >= lim {
		return 0, fmt.Errorf("branch target %#x out of ±%d range", target, lim*4)
	}

	return uint32(off) & ((1 << uint(bits)) - 1), nil
}

// ctx — the instruction encoding context: address (legacy-path branches).
type ctx struct {
	Addr uint64
}
