package arm64

// The remaining Builder families whose structs carry internal types
// (the literal source text) — completes the construction vocabulary
// to all families the assembler ctors build. The scalar FP families
// are per-instruction types with their own Builder methods (fadd.go,
// fcmp.go, ...) — this file keeps only the ones without them.

// FpKind — the register kind of an FP-family operand (exported alias).
type FpKind = fpKind

// The FP operand kinds.
const (
	FS = kS // s registers
	FD = kD // d registers
	FW = kW // w registers (int side)
	FX = kX // x registers (int side)
)

// Ld1 — the ldN/stN multi- and element forms: regList is the rendered
// "{ v0, v1 }" / "{ v0 }[2]" operand; postReg/postImm/hasPost the
// post-index; rtNum/count/opcode/size/q the layout fields; isElem the
// single-structure element form.
func NewLd1(
	regList, rn, name, arr, postReg string,
	postImm uint32, hasPost bool, enc uint32,
	rtNum uint32, count int, opcode, size, q uint32, isElem bool,
) Instr {
	return Ld1{
		regList: regList, rn: rn, name: name, arr: arr,
		postReg: postReg, postImm: postImm, hasPost: hasPost,
		enc: enc, rtNum: rtNum, count: count,
		opcode: opcode, size: size, q: q, isElem: isElem,
	}
}
