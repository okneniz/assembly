package arm64

// The SIMD copy family dispatch: op bits [13:12] (0=DUP general,
// 1=INS general, 2=SMOV, 3=UMOV); imm5 = the size one-hot plus the
// lane index above it (DUP general has no index, so its imm5 must be
// one-hot). Unencodable sizes go to the .word fallback (the schema
// mask does not express any of this).
func decodeSimdCopy(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || bitsCtz(imm5) > 3 || (w>>12&3 == 0 && imm5&(imm5-1) != 0) {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	idx := imm5 >> (size + 1)
	q := w >> 30 & 1

	switch w >> 12 & 3 {
	case 0: // DUP general: vd in Rd, wn in Rn
		return newDup(
			newBase(w),
			newVReg(uint8(w&0x1f)),
			gprOf(w>>5&0x1f, size == 3),
			decodeArrangement(q, size),
		)
	case 1: // INS general: vd in Rd, wn in Rn
		return newIns(
			newBase(w),
			size,
			idx,
			newVReg(uint8(w&0x1f)),
			gprOf(w>>5&0x1f, size == 3),
		)
	case 2: // SMOV: the GPR destination in Rd, the vector source in Rn
		return newSmov(
			newBase(w),
			q,
			size,
			idx,
			newVReg(uint8(w>>5&0x1f)),
			gprOf(w&0x1f, q == 1),
		)
	default: // UMOV: the GPR destination in Rd, the vector source in Rn
		return newUmov(
			newBase(w),
			q,
			size,
			idx,
			newVReg(uint8(w>>5&0x1f)),
			gprOf(w&0x1f, size == 3),
		)
	}
}
