package riscv

// Decode table: the order defines the match priority
// (first-match-wins) - for example, OP-IMM shifts before OP-IMM arithmetic
// (a single opcode 0x13). match/mask are authoritative - from the generated
// riscvEncodings (Spike encoding.h); each entry builds its own
// per-instruction structure.

import "strings"

type decodeEntry struct {
	name string
	ctor func(word uint32) Instr
}

func newDecodeEntry(name string, ctor func(word uint32) Instr) decodeEntry {
	return decodeEntry{
		name: name,
		ctor: ctor,
	}
}

var decodeTable = []decodeEntry{
	newDecodeEntry("lui", decodeLui),
	newDecodeEntry("auipc", decodeAuipc),
	newDecodeEntry("jal", decodeJal),
	newDecodeEntry("jalr", decodeJalr),
	newDecodeEntry("beq", decodeBeq),
	newDecodeEntry("bne", decodeBne),
	newDecodeEntry("blt", decodeBlt),
	newDecodeEntry("bge", decodeBge),
	newDecodeEntry("bltu", decodeBltu),
	newDecodeEntry("bgeu", decodeBgeu),
	newDecodeEntry("lb", decodeLb),
	newDecodeEntry("lh", decodeLh),
	newDecodeEntry("lw", decodeLw),
	newDecodeEntry("ld", decodeLd),
	newDecodeEntry("lbu", decodeLbu),
	newDecodeEntry("lhu", decodeLhu),
	newDecodeEntry("lwu", decodeLwu),
	newDecodeEntry("sb", decodeSb),
	newDecodeEntry("sh", decodeSh),
	newDecodeEntry("sw", decodeSw),
	newDecodeEntry("sd", decodeSd),
	newDecodeEntry("slli", decodeSlli),
	newDecodeEntry("srli", decodeSrli),
	newDecodeEntry("srai", decodeSrai),
	newDecodeEntry("addi", decodeAddi),
	newDecodeEntry("slti", decodeSlti),
	newDecodeEntry("sltiu", decodeSltiu),
	newDecodeEntry("xori", decodeXori),
	newDecodeEntry("ori", decodeOri),
	newDecodeEntry("andi", decodeAndi),
	newDecodeEntry("slliw", decodeSlliw),
	newDecodeEntry("srliw", decodeSrliw),
	newDecodeEntry("sraiw", decodeSraiw),
	newDecodeEntry("addiw", decodeAddiw),
	newDecodeEntry("add", decodeAdd),
	newDecodeEntry("sub", decodeSub),
	newDecodeEntry("sll", decodeSll),
	newDecodeEntry("slt", decodeSlt),
	newDecodeEntry("sltu", decodeSltu),
	newDecodeEntry("xor", decodeXor),
	newDecodeEntry("srl", decodeSrl),
	newDecodeEntry("sra", decodeSra),
	newDecodeEntry("or", decodeOr),
	newDecodeEntry("and", decodeAnd),
	newDecodeEntry("addw", decodeAddw),
	newDecodeEntry("subw", decodeSubw),
	newDecodeEntry("sllw", decodeSllw),
	newDecodeEntry("srlw", decodeSrlw),
	newDecodeEntry("sraw", decodeSraw),
	newDecodeEntry("mul", decodeMul),
	newDecodeEntry("mulh", decodeMulh),
	newDecodeEntry("mulhsu", decodeMulhsu),
	newDecodeEntry("mulhu", decodeMulhu),
	newDecodeEntry("div", decodeDiv),
	newDecodeEntry("divu", decodeDivu),
	newDecodeEntry("rem", decodeRem),
	newDecodeEntry("remu", decodeRemu),
	newDecodeEntry("mulw", decodeMulw),
	newDecodeEntry("divw", decodeDivw),
	newDecodeEntry("divuw", decodeDivuw),
	newDecodeEntry("remw", decodeRemw),
	newDecodeEntry("remuw", decodeRemuw),
	newDecodeEntry("flw", decodeFlw),
	newDecodeEntry("fld", decodeFld),
	newDecodeEntry("fsw", decodeFsw),
	newDecodeEntry("fsd", decodeFsd),
	newDecodeEntry("fadd.s", decodeFaddS),
	newDecodeEntry("fsub.s", decodeFsubS),
	newDecodeEntry("fmul.s", decodeFmulS),
	newDecodeEntry("fdiv.s", decodeFdivS),
	newDecodeEntry("fadd.d", decodeFaddD),
	newDecodeEntry("fsub.d", decodeFsubD),
	newDecodeEntry("fmul.d", decodeFmulD),
	newDecodeEntry("fdiv.d", decodeFdivD),
	newDecodeEntry("fmadd.s", decodeFmaddS),
	newDecodeEntry("fmadd.d", decodeFmaddD),
	newDecodeEntry("fmsub.s", decodeFmsubS),
	newDecodeEntry("fmsub.d", decodeFmsubD),
	newDecodeEntry("fnmsub.s", decodeFnmsubS),
	newDecodeEntry("fnmsub.d", decodeFnmsubD),
	newDecodeEntry("fnmadd.s", decodeFnmaddS),
	newDecodeEntry("fnmadd.d", decodeFnmaddD),
	newDecodeEntry("amoadd.w", decodeAmoaddW),
	newDecodeEntry("amoadd.d", decodeAmoaddD),
	newDecodeEntry("amoswap.w", decodeAmoswapW),
	newDecodeEntry("amoswap.d", decodeAmoswapD),
	newDecodeEntry("amoxor.w", decodeAmoxorW),
	newDecodeEntry("amoxor.d", decodeAmoxorD),
	newDecodeEntry("amoor.w", decodeAmoorW),
	newDecodeEntry("amoor.d", decodeAmoorD),
	newDecodeEntry("amoand.w", decodeAmoandW),
	newDecodeEntry("amoand.d", decodeAmoandD),
	newDecodeEntry("amomin.w", decodeAmominW),
	newDecodeEntry("amomin.d", decodeAmominD),
	newDecodeEntry("amomax.w", decodeAmomaxW),
	newDecodeEntry("amomax.d", decodeAmomaxD),
	newDecodeEntry("amominu.w", decodeAmominuW),
	newDecodeEntry("amominu.d", decodeAmominuD),
	newDecodeEntry("amomaxu.w", decodeAmomaxuW),
	newDecodeEntry("amomaxu.d", decodeAmomaxuD),
	newDecodeEntry("fence", decodeFence),
	newDecodeEntry("csrrw", decodeCsrrw),
	newDecodeEntry("csrrs", decodeCsrrs),
	newDecodeEntry("csrrc", decodeCsrrc),
	newDecodeEntry("csrrwi", decodeCsrrwi),
	newDecodeEntry("csrrsi", decodeCsrrsi),
	newDecodeEntry("csrrci", decodeCsrrci),
	newDecodeEntry("ecall", decodeSystem("ecall", "RV32I")),
	newDecodeEntry("ebreak", decodeSystem("ebreak", "RV32I")),
	newDecodeEntry("mret", decodeSystem("mret", "Privileged")),
	newDecodeEntry("sret", decodeSystem("sret", "Privileged")),
	newDecodeEntry("wfi", decodeSystem("wfi", "Privileged")),
}

// encName converts a mnemonic to an riscvEncodings key ("fadd.s" -> "fadd_s").
func encName(mnem string) string {
	return strings.ReplaceAll(mnem, ".", "_")
}

func decodeAdd(w uint32) Instr {
	return Add{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAddi(w uint32) Instr {
	return Addi{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeAddiw(w uint32) Instr {
	return Addiw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeAddw(w uint32) Instr {
	return Addw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoaddD(w uint32) Instr {
	return AmoaddD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoaddW(w uint32) Instr {
	return AmoaddW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoandD(w uint32) Instr {
	return AmoandD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoandW(w uint32) Instr {
	return AmoandW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmomaxD(w uint32) Instr {
	return AmomaxD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmomaxW(w uint32) Instr {
	return AmomaxW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmomaxuD(w uint32) Instr {
	return AmomaxuD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmomaxuW(w uint32) Instr {
	return AmomaxuW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmominD(w uint32) Instr {
	return AmominD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmominW(w uint32) Instr {
	return AmominW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmominuD(w uint32) Instr {
	return AmominuD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmominuW(w uint32) Instr {
	return AmominuW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoorD(w uint32) Instr {
	return AmoorD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoorW(w uint32) Instr {
	return AmoorW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoswapD(w uint32) Instr {
	return AmoswapD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoswapW(w uint32) Instr {
	return AmoswapW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoxorD(w uint32) Instr {
	return AmoxorD{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAmoxorW(w uint32) Instr {
	return AmoxorW{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAnd(w uint32) Instr {
	return And{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeAndi(w uint32) Instr {
	return Andi{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeAuipc(w uint32) Instr {
	return Auipc{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		imm:  immNum(int64(uImm(w))),
	}
}

func decodeBeq(w uint32) Instr {
	return Beq{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeBge(w uint32) Instr {
	return Bge{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeBgeu(w uint32) Instr {
	return Bgeu{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeBlt(w uint32) Instr {
	return Blt{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeBltu(w uint32) Instr {
	return Bltu{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeBne(w uint32) Instr {
	return Bne{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(bImm(w)),
	}
}

func decodeCsrrc(w uint32) Instr {
	return Csrrc{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
}

func decodeCsrrci(w uint32) Instr {
	return Csrrci{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
}

func decodeCsrrs(w uint32) Instr {
	return Csrrs{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
}

func decodeCsrrsi(w uint32) Instr {
	return Csrrsi{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
}

func decodeCsrrw(w uint32) Instr {
	return Csrrw{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		rs1:   rvRegNames[w>>15&0x1f],
	}
}

func decodeCsrrwi(w uint32) Instr {
	return Csrrwi{
		base:  newBase(w),
		csrOp: newCsrOp(rvRegNames[w>>7&0x1f], int64(w>>20&0xfff)),
		zimm:  immNum(int64(w >> 15 & 0x1f)),
	}
}

func decodeDiv(w uint32) Instr {
	return Div{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeDivu(w uint32) Instr {
	return Divu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeDivuw(w uint32) Instr {
	return Divuw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeDivw(w uint32) Instr {
	return Divw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeFaddD(w uint32) Instr {
	return FaddD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFaddS(w uint32) Instr {
	return FaddS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFdivD(w uint32) Instr {
	return FdivD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFdivS(w uint32) Instr {
	return FdivS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFence(w uint32) Instr {
	return Fence{
		base: newBase(w),
		fm:   immNum(int64(w >> 20 & 0xf)),
	}
}

func decodeFld(w uint32) Instr {
	return Fld{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeFlw(w uint32) Instr {
	return Flw{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeFmaddD(w uint32) Instr {
	return FmaddD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFmaddS(w uint32) Instr {
	return FmaddS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFmsubD(w uint32) Instr {
	return FmsubD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFmsubS(w uint32) Instr {
	return FmsubS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFmulD(w uint32) Instr {
	return FmulD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFmulS(w uint32) Instr {
	return FmulS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFnmaddD(w uint32) Instr {
	return FnmaddD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFnmaddS(w uint32) Instr {
	return FnmaddS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFnmsubD(w uint32) Instr {
	return FnmsubD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFnmsubS(w uint32) Instr {
	return FnmsubS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rs3:  rvFRegNames[w>>27&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFsd(w uint32) Instr {
	return Fsd{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeFsubD(w uint32) Instr {
	return FsubD{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFsubS(w uint32) Instr {
	return FsubS{
		base: newBase(w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvFRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		rm:   immNum(int64(w >> 12 & 0x7)),
	}
}

func decodeFsw(w uint32) Instr {
	return Fsw{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvFRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeJal(w uint32) Instr {
	return Jal{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		off:  immNum(jImm(w)),
	}
}

func decodeJalr(w uint32) Instr {
	return Jalr{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLb(w uint32) Instr {
	return Lb{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLbu(w uint32) Instr {
	return Lbu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLd(w uint32) Instr {
	return Ld{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLh(w uint32) Instr {
	return Lh{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLhu(w uint32) Instr {
	return Lhu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLui(w uint32) Instr {
	return Lui{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		imm:  immNum(int64(uImm(w))),
	}
}

func decodeLw(w uint32) Instr {
	return Lw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeLwu(w uint32) Instr {
	return Lwu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func decodeMul(w uint32) Instr {
	return Mul{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeMulh(w uint32) Instr {
	return Mulh{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeMulhsu(w uint32) Instr {
	return Mulhsu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeMulhu(w uint32) Instr {
	return Mulhu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeMulw(w uint32) Instr {
	return Mulw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeOr(w uint32) Instr {
	return Or{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeOri(w uint32) Instr {
	return Ori{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeRem(w uint32) Instr {
	return Rem{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeRemu(w uint32) Instr {
	return Remu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeRemuw(w uint32) Instr {
	return Remuw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeRemw(w uint32) Instr {
	return Remw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSb(w uint32) Instr {
	return Sb{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeSd(w uint32) Instr {
	return Sd{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeSh(w uint32) Instr {
	return Sh{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeSll(w uint32) Instr {
	return Sll{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSlli(w uint32) Instr {
	return Slli{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

func decodeSlliw(w uint32) Instr {
	return Slliw{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt5(w))),
	}
}

func decodeSllw(w uint32) Instr {
	return Sllw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSlt(w uint32) Instr {
	return Slt{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSlti(w uint32) Instr {
	return Slti{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeSltiu(w uint32) Instr {
	return Sltiu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}

func decodeSltu(w uint32) Instr {
	return Sltu{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSra(w uint32) Instr {
	return Sra{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSrai(w uint32) Instr {
	return Srai{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

func decodeSraiw(w uint32) Instr {
	return Sraiw{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt5(w))),
	}
}

func decodeSraw(w uint32) Instr {
	return Sraw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSrl(w uint32) Instr {
	return Srl{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSrli(w uint32) Instr {
	return Srli{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

func decodeSrliw(w uint32) Instr {
	return Srliw{
		base:  newBase(w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt5(w))),
	}
}

func decodeSrlw(w uint32) Instr {
	return Srlw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSub(w uint32) Instr {
	return Sub{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSubw(w uint32) Instr {
	return Subw{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeSw(w uint32) Instr {
	return Sw{
		base: newBase(w),
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
		off:  immNum(sImm(w)),
	}
}

func decodeSystem(name, group string) func(uint32) Instr {
	return func(w uint32) Instr {
		return systemInstr{
			base:  newBase(w),
			name:  name,
			group: group,
		}
	}
}

func decodeXor(w uint32) Instr {
	return Xor{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func decodeXori(w uint32) Instr {
	return Xori{
		base: newBase(w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		imm:  immNum(iImm(w)),
	}
}
