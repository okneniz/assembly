package arm64

// The remaining Builder families whose structs carry internal types
// (fpKind, the literal source text) — completes the construction
// vocabulary to all families the assembler ctors build.

// FpKind — the register kind of an FP-family operand (exported alias).
type FpKind = fpKind

// The FP operand kinds.
const (
	FS = kS // s registers
	FD = kD // d registers
	FW = kW // w registers (int side)
	FX = kX // x registers (int side)
)

// Fp2 — an FP two-register operation (fneg/fmov/fcvt...): the operand
// kinds of both sides (FS/FD/FW/FX).
func NewFp2(op string, rd, rn string, enc uint32, rdK, rnK FpKind) Instr {
	return Fp2{
		op:  op,
		rd:  rd,
		rn:  rn,
		enc: enc,
		rdK: rdK,
		rnK: rnK,
	}
}

// Fcmp — fcmp/fcmpe: withRM selects the register form (vs #0.0).
func NewFcmp(fn, fm string, withRM bool, enc0, encR uint32, k FpKind) Instr {
	return Fcmp{
		rn:     fn,
		rm:     fm,
		withRM: withRM,
		enc0:   enc0,
		encR:   encR,
		k:      k,
	}
}

// FmovImm — fmov rd, #literal: text is the literal's source spelling
// (the struct prints it back verbatim).
func NewFmovImm(
	rd string,
	val float64,
	text string,
	isS bool,
	enc uint32,
	rdK FpKind,
) Instr {
	return FmovImm{
		rd:   rd,
		val:  val,
		text: text,
		isS:  isS,
		enc:  enc,
		rdK:  rdK,
	}
}

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
