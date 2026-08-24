package riscv

// Bit layouts of the RISC-V immediate fields (I/S/B/U/J formats and shamt):
// scattered fields are gathered from the 32-bit word per format; the reverse
// packing for the assembler (encI/encS/...) is in asm_encode.go.

func signExtendN(v uint32, bits uint) int64 {
	s := int64(v)
	sign := int64(1) << (bits - 1)
	if s&sign != 0 {
		s |= ^int64((uint64(1) << bits) - 1)
	}

	return s
}

// I-type imm: bits[31:20], 12-bit signed (addi/lb/.../jalr offset).
func iImm(w uint32) int64 {
	return signExtendN((w>>20)&0xfff, 12)
}

// S-type imm: bits[31:25]<<5 | bits[11:7], 12-bit signed (sb/sh/sw/sd).
func sImm(w uint32) int64 {
	v := ((w>>25)&0x7f)<<5 | ((w >> 7) & 0x1f)
	return signExtendN(v, 12)
}

// B-type imm: bits[31|30:25|11:8|7] → imm[12|10:5|4:1|11], 13-bit signed (branch).
func bImm(w uint32) int64 {
	v := ((w>>31)&1)<<12 | ((w>>25)&0x3f)<<5 | ((w>>8)&0xf)<<1 | ((w>>7)&1)<<11
	return signExtendN(v, 13)
}

// U-type imm: bits[31:12], printed as the raw 20-bit value (lui/auipc).
func uImm(w uint32) uint32 {
	return (w >> 12) & 0xfffff
}

// J-type imm: bits[31|30:21|20|19:12] → imm[20|10:1|11|19:12], 21-bit signed (jal).
func jImm(w uint32) int64 {
	v := ((w>>31)&1)<<20 | ((w>>21)&0x3ff)<<1 | ((w>>20)&1)<<11 | ((w>>12)&0xff)<<12
	return signExtendN(v, 21)
}

// RV64 shift amount: bits[25:20] (6 bit) for slli/srli/srai.
func shamt6(w uint32) uint32 {
	return (w >> 20) & 0x3f
}

// RV32/RV64 W shift amount: bits[24:20] (5 bit) for slliw/srliw/sraw.
func shamt5(w uint32) uint32 {
	return (w >> 20) & 0x1f
}
