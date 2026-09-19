package arm64

// The logical three-same group dispatch: the family's schemas (and,
// eor) free the opcode bits 23:22, so one schema covers four
// instructions — the mnemonic is selected by U (bit 29) + the opcode
// pair, each landing on its own instruction type.

// decodeSimd3Logical - the opcode is read FROM THE WORD (the group's
// family schemas free the opcode bits and each covers four
// instructions). The arrangement of the logical group depends only on
// Q (8b/16b), the size bits are part of the opcode.
func decodeSimd3Logical(w uint32) (Instr, error) {
	arr := "8b"
	if w>>30&1 == 1 {
		arr = "16b"
	}

	b := newBase(w)
	rd := newVReg(uint8(w & 0x1f))
	rn := newVReg(uint8(w>>5) & 0x1f)
	rm := newVReg(uint8(w>>16) & 0x1f)

	switch w>>29&1<<2 | w>>22&3 {
	case 0b0000:
		return newAnd(b, rd, rn, rm, arr)
	case 0b0001:
		return newBic(b, rd, rn, rm, arr)
	case 0b0010:
		return newOrr(b, rd, rn, rm, arr)
	case 0b0011:
		return newOrn(b, rd, rn, rm, arr)
	case 0b0100:
		return newEor(b, rd, rn, rm, arr)
	case 0b0101:
		return newBsl(b, rd, rn, rm, arr)
	case 0b0110:
		return newBit(b, rd, rn, rm, arr)
	default:
		return newBif(b, rd, rn, rm, arr)
	}
}

// simd3LogicalNames - the group's mnemonics (the exported membership
// check lives on them).
var simd3LogicalNames = map[string]bool{
	"and": true, "bic": true, "orr": true, "orn": true,
	"eor": true, "bsl": true, "bit": true, "bif": true,
}

// IsSimd3Logical - whether the mnemonic belongs to the logical
// three-same group (bits 23:22 are its opcode, not an arrangement).
func IsSimd3Logical(name string) bool {
	return simd3LogicalNames[name]
}
