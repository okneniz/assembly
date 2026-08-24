package riscv

// RVC bit utilities (C extension): packing the fields of compressed halfwords.
// Used by the compress branches of per-instruction structures.

// r3ok - the register lies in the compressed quadrant x8-x15 (3-bit addressing).
func r3ok(name string) bool {
	r, ok := asmRegNum[name]
	return ok && !r.fp && r.num >= 8 && r.num <= 15
}

// cr3 - 3-bit register index x8-x15.
func cr3(name string) uint16 {
	return uint16(asmRegNum[name].num - 8)
}

// r5 - 5-bit register index.
func r5(name string) uint16 {
	return uint16(asmRegNum[name].num)
}

// ciBits - CI immediate {bit12, bits[6:2]} from a 6-bit signed value.
func ciBits(v int64) uint16 {
	u := uint16(v & 0x3f)
	return (u>>5)&1<<12 | (u&0x1f)<<2
}

// immFits - the value is within [lo, hi] and divisible by div.
func immFits(v, lo, hi int64, div int64) bool {
	return v >= lo && v <= hi && v%div == 0
}
