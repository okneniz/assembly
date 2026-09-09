package arm64

// The SIMD/FP part of the Builder vocabulary (scalar families are next
// to their structs; this completes the programmatic construction
// surface so the assembler constructors (asm/arm64 after the
// arch-exodus) build through exported API only).

// Fp3 — an FP three-register operation (fadd/fsub/fmul/fdiv/fmax/fmin
// ...): the mnemonic, three FP registers and the family encoding base
// (the s/d variant chosen by the caller).
func NewFp3(op string, fd, fn, fm string, enc uint32) Instr {
	return Fp3{
		op:  op,
		rd:  fd,
		rn:  fn,
		rm:  fm,
		enc: enc,
	}
}

// Fp4 — an FP four-register operation (fmadd/fmsub/...).
func NewFp4(op string, fd, fn, fm, fa string, enc uint32) Instr {
	return Fp4{
		op:  op,
		rd:  fd,
		rn:  fn,
		rm:  fm,
		ra:  fa,
		enc: enc,
	}
}

// Simd2 — a two-register SIMD operation (cnt/abs/not/...): q/size carry
// the arrangement.
func NewSimd2(op, rd, rn string, arr string, enc, q, size uint32) Instr {
	return Simd2{
		op:   op,
		rd:   rd,
		rn:   rn,
		arr:  arr,
		enc:  enc,
		q:    q,
		size: size,
	}
}

// Simd3 — a three-same SIMD operation (add/cmeq/eor/...): q/size carry
// the arrangement.
func NewSimd3(op, rd, rn, rm string, enc, q, size uint32) Instr {
	return Simd3{
		op:   op,
		rd:   rd,
		rn:   rn,
		rm:   rm,
		enc:  enc,
		q:    q,
		size: size,
	}
}

// SimdShift — a SIMD shift by immediate (shl/ssra/...): immh/immb carry
// the shift amount and lane size.
func NewSimdShift(op, rd, rn string, immh, immb, q, enc uint32) Instr {
	return SimdShift{
		op: op, rd: rd, rn: rn, immh: immh, immb: immb, q: q, enc: enc,
	}
}

// SimdWiden — a widening three-same SIMD operation (saddw/ssubw/...):
// registers as strings and numbers (the struct prints from the strings,
// encodes from the numbers).
func NewSimdWiden(
	op string,
	q, size uint32,
	rd, rn, rm string,
	enc uint32,
	rdN, rnN, rmN uint32,
) Instr {
	return SimdWiden{
		op: op, q: q, size: size, rd: rd, rn: rn, rm: rm,
		enc: enc, rdN: rdN, rnN: rnN, rmN: rmN,
	}
}

// SimdCopyGPR — the GPR-side Advanced SIMD copy (dup/ins/smov/umov):
// op selects the operation, idx the lane (dup has none), isDest the
// smov/umov register-slot swap.
func NewSimdCopyGPR(
	op, vd, gpr string,
	size, idx, q uint32,
	vdNum, gprNum uint32,
	isDest bool,
) Instr {
	return SimdCopy{
		op: op, size: size, idx: idx, q: q, vd: vd, gpr: gpr,
		enc: 0x0E000000, vdNum: vdNum, gprNum: gprNum, isDest: isDest,
	}
}

// DupElem — the vector-source Advanced SIMD copy (dup element/ins element/
// scalar mov alias): registers as strings and numbers (the struct prints
// from the strings, encodes from the numbers), idx the dup source / ins
// destination lane, srcIdx the ins source lane.
func NewDupElem(
	op string,
	size, idx, srcIdx, q uint32,
	rd, rn string,
	rdN, rnN uint32,
) Instr {
	return DupElem{
		op: op, size: size, idx: idx, srcIdx: srcIdx, q: q,
		rd: rd, rn: rn, rdN: rdN, rnN: rnN,
	}
}

// Uaddlv — uaddlv vd, vn: q/size carry the arrangement.
func NewUaddlv(rd, rn string, q, size uint32) Instr {
	return Uaddlv{
		rd:   rd,
		rn:   rn,
		q:    q,
		size: size,
	}
}

// V1arr — a one-register-and-arrangement SIMD operation (aes/aesmc...).
func NewV1arr(op, rd, rn string, enc uint32) Instr {
	return V1arr{
		op:  op,
		rd:  rd,
		rn:  rn,
		enc: enc,
	}
}

// MovSimd — the whole-register vector mov (orr alias): the arrangement
// in arr.
func NewMovSimd(rd, rm string, arr string, enc uint32) Instr {
	return MovSimd{
		rd:  rd,
		rm:  rm,
		arr: arr,
		enc: enc,
	}
}

// ByElem — the by-element SIMD family (24 mnemonics; fcmla with
// rotation). The Builder takes the fully computed layout.
func NewByElem(
	name string,
	q, size uint32,
	rd, rn, rm string,
	idx int64,
	long bool,
	rot uint32,
	rdN, rnN, rmN uint32,
) ByElem {
	return ByElem{
		name: name, q: q, size: size,
		rd: rd, rn: rn, rm: rm,
		idx: uint32(idx), long: long, rot: rot,
		rdN: rdN, rnN: rnN, rmN: rmN,
	}
}
