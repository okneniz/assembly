package riscv

// RVC (compressed, 16-bit) decoding: every c.* encoding builds the same
// per-instruction structure as its 32-bit equivalent, but with length=2
// and raw=halfword - objdump (llvm and binutils) prints c.* in expanded
// form (li / slli / add / lw / beqz / j / ...), so the alias logic
// of the structures yields the correct text automatically.
//
// Registers of the restricted quadrants are 3-bit indices x8-x15.
// Immediate layouts follow the C-extension specification.

// r3 maps a 3-bit compressed register field to its ABI name (x8-x15).
func r3(v uint32) string {
	return rvRegNames[8+(v&0x7)]
}

// half - the base of a compressed instruction (raw = zero-extended halfword).
func newHalfBase(h uint32, addr uint64) base {
	return base{
		addr:   addr,
		raw:    h,
		length: 2,
	}
}

// compressedInstruction decodes a 16-bit instruction (zero-extended
// halfword).
func compressedInstruction(halfword uint32, addr uint64) Instr {
	h := halfword & 0xffff
	switch h & 0x3 {
	case 0:
		return decodeQ0(h, addr)
	case 1:
		return decodeQ1(h, addr)
	case 2:
		return decodeQ2(h, addr)
	}

	return newUnknown(newHalfBase(h, addr))
}

// CI imm {bit12, bits[6:2]} → 6-bit signed (c.addi/c.li/c.addiw).
func ciImm(h uint32) int64 {
	return signExtendN(((h>>12)&1)<<5|((h>>2)&0x1f), 6)
}

// CI shamt {bit12, bits[6:2]} → 6-bit unsigned (c.slli).
func ciShamt(h uint32) uint32 {
	return ((h>>12)&1)<<5 | ((h >> 2) & 0x1f)
}

// --- quadrant 0: c.addi4spn / c.lw / c.ld / c.sw / c.sd ---

func decodeQ0(h uint32, addr uint64) Instr {
	rs13 := r3(h >> 7)
	switch (h >> 13) & 0x7 {
	case 0: // c.addi4spn: addi rd', sp, nzuimm
		imm := addi4spnImm(h)
		if imm == 0 {
			return newUnknown(newHalfBase(h, addr))
		}

		return cAddi(h, addr, r3(h>>2), "sp", imm)
	case 2: // c.lw
		return cLw(h, addr, r3(h>>2), rs13, clLwImm(h))
	case 3: // c.ld (RV64)
		return cLd(h, addr, r3(h>>2), rs13, clLdImm(h))
	case 6: // c.sw
		return cSw(h, addr, rs13, r3(h>>2), clLwImm(h))
	case 7: // c.sd (RV64)
		return cSd(h, addr, rs13, r3(h>>2), clLdImm(h))
	}

	return newUnknown(newHalfBase(h, addr))
}

// c.addi4spn nzuimm: [9:6]=inst[10:7], [5:4]=inst[12:11], [3]=inst[5], [2]=inst[6].
func addi4spnImm(h uint32) int64 {
	return int64(((h>>7)&0xf)<<6 | ((h>>11)&0x3)<<4 | ((h>>5)&1)<<3 | ((h>>6)&1)<<2)
}

// c.lw/c.sw offset: [5:3]=inst[12:10], [2]=inst[6], [6]=inst[5].
func clLwImm(h uint32) int64 {
	return int64(((h>>10)&0x7)<<3 | ((h>>6)&1)<<2 | ((h>>5)&1)<<6)
}

// c.ld/c.sd offset (RV64): [5:3]=inst[12:10], [7:6]=inst[6:5].
func clLdImm(h uint32) int64 {
	return int64(((h>>10)&0x7)<<3 | ((h>>5)&0x3)<<6)
}

// --- quadrant 1: c.nop/c.addi / c.addiw / c.li / c.addi16sp / c.lui /
// c.srli/srai/andi/sub/xor/or/and / c.beqz / c.bnez / c.j ---

func decodeQ1(h uint32, addr uint64) Instr {
	rd := rvRegNames[(h>>7)&0x1f]
	switch (h >> 13) & 0x7 {
	case 0: // c.addi / c.nop
		imm := ciImm(h)
		if (h>>7)&0x1f == 0 {
			return cAddi(h, addr, "zero", "zero", 0)
		}

		return cAddi(h, addr, rd, rd, imm)
	case 1: // c.addiw (RV64; c.jal on RV32)
		if (h>>7)&0x1f == 0 {
			return newUnknown(newHalfBase(h, addr))
		}

		return cAddiw(h, addr, rd, rd, ciImm(h))
	case 2: // c.li
		if (h>>7)&0x1f == 0 {
			return cAddi(h, addr, "zero", "zero", 0)
		}

		return cAddi(h, addr, rd, "zero", ciImm(h))
	case 3: // c.addi16sp (rd==2) / c.lui
		switch (h >> 7) & 0x1f {
		case 2:
			return cAddi(h, addr, "sp", "sp", addi16spImm(h))
		case 0:
			return newUnknown(newHalfBase(h, addr))
		default:
			return cLui(h, addr, rd, luiImm(h))
		}

	case 4: // c.srli / c.srai / c.andi / c.sub / c.xor / c.or / c.and
		return decodeCA(h, addr)
	case 6: // c.beqz
		return cBeq(h, addr, r3(h>>7), "zero", int64(addr)+cbImm(h))
	case 7: // c.bnez
		return cBne(h, addr, r3(h>>7), "zero", int64(addr)+cbImm(h))
	case 5: // c.j
		return cJal(h, addr, "zero", int64(addr)+cjImm(h))
	}

	return newUnknown(newHalfBase(h, addr))
}

func decodeCA(h uint32, addr uint64) Instr {
	rs := r3(h >> 7)
	switch (h >> 10) & 0x3 {
	case 0: // c.srli (bit12=0) / c.srai (bit12=1)
		sh := ((h >> 2) & 0x1f) | ((h>>12)&1)<<5
		if (h>>12)&1 == 0 {
			return cSrli(h, addr, rs, rs, int64(sh))
		}

		return cSrai(h, addr, rs, rs, int64(sh))
	case 1: // c.andi
		return cAndi(h, addr, rs, rs, signExtendN(((h>>2)&0x1f)|((h>>12)&1)<<5, 6))
	}

	// bits[11:10]=10: R-type; op in bits[6:5]; bit12 selects the RV64 /w variant.
	rs2 := r3(h >> 2)
	switch ((h >> 5) & 0x3) | ((h>>12)&1)<<2 {
	case 0:
		return cSub(h, addr, rs, rs, rs2)
	case 1:
		return cXor(h, addr, rs, rs, rs2)
	case 2:
		return cOr(h, addr, rs, rs, rs2)
	case 3:
		return cAnd(h, addr, rs, rs, rs2)
	case 4:
		return cSubw(h, addr, rs, rs, rs2)
	case 5:
		return cAddw(h, addr, rs, rs, rs2)
	}

	return newUnknown(newHalfBase(h, addr))
}

// c.addi16sp nzimm: [9]=inst[12], [8:7]=inst[4:3], [6]=inst[5], [5]=inst[2], [4]=inst[6].
func addi16spImm(h uint32) int64 {
	return signExtendN(((h>>12)&1)<<9|((h>>3)&0x3)<<7|((h>>5)&1)<<6|((h>>2)&1)<<5|((h>>6)&1)<<4, 10)
}

// c.lui nzimm: [17]=inst[12], [16:12]=inst[6:2].
func luiImm(h uint32) int64 {
	return signExtendN(((h>>12)&1)<<17|((h>>2)&0x1f)<<12, 18)
}

// c.beqz/c.bnez offset: [8]=inst[12], [7:6]=inst[6:5], [5]=inst[2],
// [4:3]=inst[11:10], [2:1]=inst[4:3].
func cbImm(h uint32) int64 {
	return signExtendN(
		((h>>12)&1)<<8|((h>>5)&0x3)<<6|((h>>2)&1)<<5|((h>>10)&0x3)<<3|((h>>3)&0x3)<<1,
		9,
	)
}

// c.j offset: [11]=inst[12], [10]=inst[8], [9:8]=inst[10:9], [7]=inst[6],
// [6]=inst[7], [5]=inst[2], [4]=inst[11], [3:1]=inst[5:3].
func cjImm(h uint32) int64 {
	return signExtendN(
		((h>>12)&1)<<11|((h>>8)&1)<<10|((h>>9)&0x3)<<8|((h>>6)&1)<<7|((h>>7)&1)<<6|((h>>2)&1)<<5|((h>>11)&1)<<4|((h>>3)&0x7)<<1,
		12,
	)
}

// --- quadrant 2: c.slli / c.lwsp / c.ldsp / c.mv/c.add/c.jr/c.jalr / c.swsp / c.sdsp ---

func decodeQ2(h uint32, addr uint64) Instr {
	rd := rvRegNames[(h>>7)&0x1f]
	switch (h >> 13) & 0x7 {
	case 0: // c.slli
		if (h>>7)&0x1f == 0 {
			return newUnknown(newHalfBase(h, addr))
		}

		return cSlli(h, addr, rd, rd, int64(ciShamt(h)))
	case 2: // c.lwsp
		if (h>>7)&0x1f == 0 {
			return newUnknown(newHalfBase(h, addr))
		}

		return cLw(h, addr, rd, "sp", lwspImm(h))
	case 3: // c.ldsp (RV64)
		if (h>>7)&0x1f == 0 {
			return newUnknown(newHalfBase(h, addr))
		}

		return cLd(h, addr, rd, "sp", ldspImm(h))
	case 4: // c.jr / c.mv (bit12=0); c.jalr / c.add / c.ebreak (bit12=1)
		rs1 := rvRegNames[(h>>7)&0x1f]
		rs2 := rvRegNames[(h>>2)&0x1f]
		if (h>>12)&1 == 0 {
			if (h>>2)&0x1f == 0 { // c.jr
				if rs1 == "zero" {
					return newUnknown(newHalfBase(h, addr))
				}

				return cJalr(h, addr, "zero", rs1, 0)
			}

			if (h>>7)&0x1f == 0 {
				return newUnknown(newHalfBase(h, addr))
			}

			return cMv(h, addr, rd, rs2) // c.mv
		}

		if (h>>2)&0x1f == 0 && (h>>7)&0x1f == 0 {
			return cSystem(h, addr, "ebreak", "RV32I")
		}

		if (h>>2)&0x1f == 0 { // c.jalr
			return cJalr(h, addr, "ra", rs1, 0)
		}

		return cAdd(h, addr, rd, rd, rs2) // c.add
	case 6: // c.swsp
		return cSw(h, addr, "sp", rvRegNames[(h>>2)&0x1f], swspImm(h))
	case 7: // c.sdsp (RV64)
		return cSd(h, addr, "sp", rvRegNames[(h>>2)&0x1f], sdspImm(h))
	}

	return newUnknown(newHalfBase(h, addr))
}

// c.lwsp offset: [5]=inst[12], [4:2]=inst[6:4], [7:6]=inst[3:2].
func lwspImm(h uint32) int64 {
	return int64(((h>>12)&1)<<5 | ((h>>4)&0x7)<<2 | ((h>>2)&0x3)<<6)
}

// c.ldsp offset (RV64): [5]=inst[12], [4:3]=inst[6:5].
func ldspImm(h uint32) int64 {
	return int64(((h>>12)&1)<<5 | ((h>>5)&0x3)<<3)
}

// c.swsp offset: [5:2]=inst[12:9], [7:6]=inst[8:7].
func swspImm(h uint32) int64 {
	return int64(((h>>9)&0xf)<<2 | ((h>>7)&0x3)<<6)
}

// c.sdsp offset (RV64): [5:3]=inst[12:10].
func sdspImm(h uint32) int64 {
	return int64(((h >> 10) & 0x7) << 3)
}

// cR3 - the shared skeleton of the compressed R forms c.sub/c.xor/c.or/c.and (basis 0x8C01)
// and c.subw/c.addw (basis 0x9C01): rd==rs1, both in x8-x15, rs2 != rd.
func cR3(rd, rs1, rs2 string, base, funct2 uint16) (uint16, bool) {
	if rd != rs1 || !r3ok(rd) || !r3ok(rs2) || rs2 == rd {
		return 0, false
	}

	return base | cr3(rd)<<7 | funct2<<5 | cr3(rs2)<<2, true
}

// cbeqz - c.beqz/c.bnez: rs1' in x8-x15, rs2 == zero, offset +/-256.
func cbeqz(rs1, rs2, name string, off int64) (uint16, bool) {
	if rs2 != "zero" || !r3ok(rs1) {
		return 0, false
	}

	if off%2 != 0 || off < -256 || off > 254 {
		return 0, false
	}

	u := uint16(off & 0x1ff) // 9-bit
	h := uint16(0xC001) | cr3(rs1)<<7
	if name == "bne" {
		h |= 1 << 13
	}

	// u[8]=h12, u[7:6]=h[6:5], u[5]=h2, u[4:3]=h[11:10], u[2:1]=h[4:3]
	h |= (u>>8&1)<<12 | (u>>6&3)<<5 | (u>>5&1)<<2 | (u>>3&3)<<10 | (u>>1&3)<<3
	return h, true
}

// cjal - c.j: jal zero, offset +/-2KB (CJ format).
func cjal(off int64) (uint16, bool) {
	if off%2 != 0 || off < -2048 || off > 2046 {
		return 0, false
	}

	u := uint16(off & 0xfff)
	bit := func(n uint) uint16 {
		return (u >> n) & 1
	}
	h := uint16(0xA001) // Q1 funct3=101
	h |= bit(11) << 12
	h |= bit(4) << 11
	h |= bit(9) << 10
	h |= bit(8) << 9
	h |= bit(10) << 8
	h |= bit(6) << 7
	h |= bit(7) << 6
	h |= bit(3) << 5
	h |= bit(2) << 4
	h |= bit(1) << 3
	h |= bit(5) << 2
	return h, true
}
