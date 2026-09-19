package arm64

// The by-element decode dispatch: one schema family covers the whole
// space, the instruction is selected by U (bit 29) and opc (bits
// 15:12) - the fp and integer entries are disjoint; the lane size
// (bits 23:22: 0=b 1=h 2=s/int-or-fp32 3=d) is validated by each
// type's constructor. Each landing is its own instruction type (see
// by_elem_common.go for the shared field layout, the per-type
// decoders live in their files).

// byElemCtors - the decode table [U][opc]; nil = an unallocated
// combination (the word falls to the .word fallback).
var byElemCtors = [2][16]func(uint32) (Instr, error){
	{ // U = 0
		0x1: decodeFmlaElem, 0x2: decodeSmlalElem, 0x3: decodeSqdmlalElem,
		0x5: decodeFmlsElem, 0x6: decodeSmlslElem, 0x7: decodeSqdmlslElem,
		0x8: decodeMulElem, 0x9: decodeFmulElem,
		0xa: decodeSmullElem, 0xb: decodeSqdmullElem,
		0xc: decodeSqdmulhElem, 0xd: decodeSqrdmulhElem,
	},
	{ // U = 1
		0x0: decodeMlaElem, 0x1: decodeFcmlaElem, 0x2: decodeUmlalElem,
		0x3: decodeFcmlaElem, 0x4: decodeMlsElem, 0x5: decodeFcmlaElem,
		0x6: decodeUmlslElem, 0x7: decodeFcmlaElem,
		0x9: decodeFmulxElem, 0xa: decodeUmullElem,
		0xd: decodeSqrdmlahElem, 0xf: decodeSqrdmlshElem,
	},
}

// decodeByElem - the word dispatch into the family's instruction
// types; unallocated combinations go to the .word fallback.
func decodeByElem(w uint32) (Instr, error) {
	ctor := byElemCtors[w>>29&1][w>>12&0xf]
	if ctor == nil {
		return decodeUnknown(w)
	}

	return ctor(w)
}
