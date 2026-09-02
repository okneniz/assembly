package loong64

// Immediate field gather/scatter for the LoongArch formats. The bit
// positions are the format-notation segments (see gen/loong/opcodes):
// index characters d/j/k/a/m/n start at bits 0/5/10/15/16/18; long
// immediates split MSB-first (Sd5k16: the low 5 bits hold imm bits
// 20..16, bits 25..10 hold imm bits 15..0).
//
// Range validation lives in the operand role types (Imm12 and friends)
// - the scatter side is pure bit packing for values already validated.
// The one dynamic value is the pc-relative branch offset (target - pc):
// encPs2 validates its alignment and range before the scatter.

import "fmt"

// signExtendN - v (already masked to its field) as a signed n-bit value.
func signExtendN(v uint32, n int) int64 {
	v &= 1<<n - 1

	return int64(int32(v<<(32-n))) >> (32 - n)
}

// uField - the unsigned field [lo+width-1 : lo] of w.
func uField(w uint32, lo, width int) uint32 {
	return w >> lo & (1<<width - 1)
}

// sField - the sign-extended field [lo+width-1 : lo] of w.
func sField(w uint32, lo, width int) int64 {
	return signExtendN(uField(w, lo, width), width)
}

// scatterU - the unsigned scatter of v into [lo+width-1 : lo] (v must be
// in range - the role types guarantee it).
func scatterU(v int64, lo, width int) uint32 {
	mask := uint32(1<<width-1) << lo

	return uint32(v) << lo & mask
}

// scatterS - the signed scatter (two's complement truncation).
func scatterS(v int64, lo, width int) uint32 {
	return uint32(v&(1<<width-1)) << lo
}

// d5k16Imm - the split offs21 field of beqz/bnez: the low 5 bits hold
// imm bits 20..16, bits 25..10 hold imm bits 15..0 (signed).
func d5k16Imm(w uint32) int64 {
	return signExtendN(uField(w, 0, 5)<<16|uField(w, 10, 16), 21)
}

// scatterD5k16 - the inverse of d5k16Imm: imm bits 20..16 to the low
// 5 bits, imm bits 15..0 to bits 25..10 (v must be a validated 21-bit
// word offset).
func scatterD5k16(v int64) uint32 {
	return uint32(v>>16)&0x1f | uint32(v)&0xffff<<10
}

// d10k16Imm - the split offs26 field of b/bl: the low 10 bits hold imm
// bits 25..16, bits 25..10 hold imm bits 15..0 (signed).
func d10k16Imm(w uint32) int64 {
	return signExtendN(uField(w, 0, 10)<<16|uField(w, 10, 16), 26)
}

// scatterD10k16 - the inverse of d10k16Imm: imm bits 25..16 to the low
// 10 bits, imm bits 15..0 to bits 25..10 (v must be a validated 26-bit
// word offset).
func scatterD10k16(v int64) uint32 {
	return uint32(v>>16)&0x3ff | uint32(v)&0xffff<<10
}

// encPs2 - the word-scaled byte offset diff (target - pc) as an n-bit
// signed word count: validates the 4-byte alignment and the range.
func encPs2(diff int64, width int, what string) (int64, error) {
	if diff%4 != 0 {
		return 0, fmt.Errorf("%s %d is not word-aligned", what, diff)
	}

	lo := -int64(1) << (width + 1)
	hi := int64(1) << (width + 1)
	if diff < lo || diff > hi-4 {
		return 0, fmt.Errorf("%s %d does not fit in %d signed bits", what, diff, width)
	}

	return diff >> 2, nil
}
