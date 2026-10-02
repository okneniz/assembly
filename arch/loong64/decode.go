package loong64

// Decode table: the order defines the match priority (first-match-wins) -
// the only overlap in the scalar integer space is the csr triple:
// csrrd/csrwr (the rj field fixed) before the general csrxchg encoding
// they specialize. match/mask are authoritative - from the generated
// loongEncodings (loongarch-opcodes tables); each entry builds its own
// per-instruction structure.

// decodeEntry - a decodeTable row: the mnemonic (the loongEncodings key)
// and the decode constructor of the per-instruction structure.
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

// decodeTable - every scalar integer instruction of the LA64 scope.
var decodeTable = []decodeEntry{
	newDecodeEntry("add.w", decodeAddW),
	newDecodeEntry("add.d", decodeAddD),
	newDecodeEntry("sub.w", decodeSubW),
	newDecodeEntry("sub.d", decodeSubD),
	newDecodeEntry("slt", decodeSlt),
	newDecodeEntry("sltu", decodeSltu),
	newDecodeEntry("maskeqz", decodeMaskeqz),
	newDecodeEntry("masknez", decodeMasknez),
	newDecodeEntry("nor", decodeNor),
	newDecodeEntry("and", decodeAnd),
	newDecodeEntry("or", decodeOr),
	newDecodeEntry("xor", decodeXor),
	newDecodeEntry("orn", decodeOrn),
	newDecodeEntry("andn", decodeAndn),
	newDecodeEntry("sll.w", decodeSllW),
	newDecodeEntry("srl.w", decodeSrlW),
	newDecodeEntry("sra.w", decodeSraW),
	newDecodeEntry("sll.d", decodeSllD),
	newDecodeEntry("srl.d", decodeSrlD),
	newDecodeEntry("sra.d", decodeSraD),
	newDecodeEntry("rotr.w", decodeRotrW),
	newDecodeEntry("rotr.d", decodeRotrD),
	newDecodeEntry("addi.w", decodeAddiW),
	newDecodeEntry("addi.d", decodeAddiD),
	newDecodeEntry("slti", decodeSlti),
	newDecodeEntry("sltui", decodeSltui),
	newDecodeEntry("andi", decodeAndi),
	newDecodeEntry("ori", decodeOri),
	newDecodeEntry("xori", decodeXori),
	newDecodeEntry("addu16i.d", decodeAddu16iD),
	newDecodeEntry("slli.w", decodeSlliW),
	newDecodeEntry("srli.w", decodeSrliW),
	newDecodeEntry("srai.w", decodeSraiW),
	newDecodeEntry("rotri.w", decodeRotriW),
	newDecodeEntry("slli.d", decodeSlliD),
	newDecodeEntry("srli.d", decodeSrliD),
	newDecodeEntry("srai.d", decodeSraiD),
	newDecodeEntry("rotri.d", decodeRotriD),
	newDecodeEntry("lu12i.w", decodeLu12iW),
	newDecodeEntry("lu32i.d", decodeLu32iD),
	newDecodeEntry("lu52i.d", decodeLu52iD),
	newDecodeEntry("pcaddi", decodePcaddi),
	newDecodeEntry("pcalau12i", decodePcalau12i),
	newDecodeEntry("pcaddu12i", decodePcaddu12i),
	newDecodeEntry("pcaddu18i", decodePcaddu18i),
	newDecodeEntry("ld.b", decodeLdB),
	newDecodeEntry("ld.h", decodeLdH),
	newDecodeEntry("ld.w", decodeLdW),
	newDecodeEntry("ld.wu", decodeLdWu),
	newDecodeEntry("ld.bu", decodeLdBu),
	newDecodeEntry("ld.hu", decodeLdHu),
	newDecodeEntry("ld.d", decodeLdD),
	newDecodeEntry("st.b", decodeStB),
	newDecodeEntry("st.h", decodeStH),
	newDecodeEntry("st.w", decodeStW),
	newDecodeEntry("st.d", decodeStD),
	newDecodeEntry("ldptr.w", decodeLdptrW),
	newDecodeEntry("ldptr.d", decodeLdptrD),
	newDecodeEntry("stptr.w", decodeStptrW),
	newDecodeEntry("stptr.d", decodeStptrD),
	newDecodeEntry("ldx.b", decodeLdxB),
	newDecodeEntry("ldx.h", decodeLdxH),
	newDecodeEntry("ldx.w", decodeLdxW),
	newDecodeEntry("ldx.wu", decodeLdxWu),
	newDecodeEntry("ldx.bu", decodeLdxBu),
	newDecodeEntry("ldx.hu", decodeLdxHu),
	newDecodeEntry("ldx.d", decodeLdxD),
	newDecodeEntry("stx.b", decodeStxB),
	newDecodeEntry("stx.h", decodeStxH),
	newDecodeEntry("stx.w", decodeStxW),
	newDecodeEntry("stx.d", decodeStxD),
	newDecodeEntry("preld", decodePreld),
	newDecodeEntry("preldx", decodePreldx),
	newDecodeEntry("beq", decodeBeq),
	newDecodeEntry("bne", decodeBne),
	newDecodeEntry("blt", decodeBlt),
	newDecodeEntry("bge", decodeBge),
	newDecodeEntry("bltu", decodeBltu),
	newDecodeEntry("bgeu", decodeBgeu),
	newDecodeEntry("beqz", decodeBeqz),
	newDecodeEntry("bnez", decodeBnez),
	newDecodeEntry("b", decodeB),
	newDecodeEntry("bl", decodeBl),
	newDecodeEntry("jirl", decodeJirl),
	newDecodeEntry("break", decodeBreak),
	newDecodeEntry("syscall", decodeSyscall),
	newDecodeEntry("dbcl", decodeDbcl),
	newDecodeEntry("dbar", decodeDbar),
	newDecodeEntry("ibar", decodeIbar),
	newDecodeEntry("mul.w", decodeMulW),
	newDecodeEntry("mulh.w", decodeMulhW),
	newDecodeEntry("mulh.wu", decodeMulhWu),
	newDecodeEntry("mul.d", decodeMulD),
	newDecodeEntry("mulh.d", decodeMulhD),
	newDecodeEntry("mulh.du", decodeMulhDu),
	newDecodeEntry("mulw.d.w", decodeMulwDW),
	newDecodeEntry("mulw.d.wu", decodeMulwDWu),
	newDecodeEntry("div.w", decodeDivW),
	newDecodeEntry("mod.w", decodeModW),
	newDecodeEntry("div.wu", decodeDivWu),
	newDecodeEntry("mod.wu", decodeModWu),
	newDecodeEntry("div.d", decodeDivD),
	newDecodeEntry("mod.d", decodeModD),
	newDecodeEntry("div.du", decodeDivDu),
	newDecodeEntry("mod.du", decodeModDu),
	newDecodeEntry("ext.w.b", decodeExtWB),
	newDecodeEntry("ext.w.h", decodeExtWH),
	newDecodeEntry("rdtimel.w", decodeRdtimelW),
	newDecodeEntry("rdtimeh.w", decodeRdtimehW),
	newDecodeEntry("rdtime.d", decodeRdtimeD),
	newDecodeEntry("cpucfg", decodeCpucfg),
	newDecodeEntry("clo.w", decodeCloW),
	newDecodeEntry("clz.w", decodeClzW),
	newDecodeEntry("cto.w", decodeCtoW),
	newDecodeEntry("ctz.w", decodeCtzW),
	newDecodeEntry("clo.d", decodeCloD),
	newDecodeEntry("clz.d", decodeClzD),
	newDecodeEntry("cto.d", decodeCtoD),
	newDecodeEntry("ctz.d", decodeCtzD),
	newDecodeEntry("revb.2h", decodeRevb2H),
	newDecodeEntry("revb.4h", decodeRevb4H),
	newDecodeEntry("revb.2w", decodeRevb2W),
	newDecodeEntry("revb.d", decodeRevbD),
	newDecodeEntry("revh.2w", decodeRevh2W),
	newDecodeEntry("revh.d", decodeRevhD),
	newDecodeEntry("bitrev.4b", decodeRevbit4B),
	newDecodeEntry("bitrev.w", decodeRevbitW),
	newDecodeEntry("bitrev.8b", decodeRevbit8B),
	newDecodeEntry("bitrev.d", decodeRevbitD),
	newDecodeEntry("bstrins.w", decodeBstrinsW),
	newDecodeEntry("bstrpick.w", decodeBstrpickW),
	newDecodeEntry("bstrins.d", decodeBstrinsD),
	newDecodeEntry("bstrpick.d", decodeBstrpickD),
	newDecodeEntry("alsl.w", decodeAlslW),
	newDecodeEntry("alsl.wu", decodeAlslWu),
	newDecodeEntry("alsl.d", decodeAlslD),
	newDecodeEntry("bytepick.w", decodeBytepickW),
	newDecodeEntry("bytepick.d", decodeBytepickD),
	newDecodeEntry("crc.w.b.w", decodeCrcWBW),
	newDecodeEntry("crc.w.h.w", decodeCrcWHW),
	newDecodeEntry("crc.w.w.w", decodeCrcWWW),
	newDecodeEntry("crc.w.d.w", decodeCrcWDW),
	newDecodeEntry("crcc.w.b.w", decodeCrccWBW),
	newDecodeEntry("crcc.w.h.w", decodeCrccWHW),
	newDecodeEntry("crcc.w.w.w", decodeCrccWWW),
	newDecodeEntry("crcc.w.d.w", decodeCrccWDW),
	newDecodeEntry("asrtle.d", decodeAsrtleD),
	newDecodeEntry("asrtgt.d", decodeAsrtgtD),
	newDecodeEntry("ldgt.b", decodeLdgtB),
	newDecodeEntry("ldgt.h", decodeLdgtH),
	newDecodeEntry("ldgt.w", decodeLdgtW),
	newDecodeEntry("ldgt.d", decodeLdgtD),
	newDecodeEntry("ldle.b", decodeLdleB),
	newDecodeEntry("ldle.h", decodeLdleH),
	newDecodeEntry("ldle.w", decodeLdleW),
	newDecodeEntry("ldle.d", decodeLdleD),
	newDecodeEntry("stgt.b", decodeStgtB),
	newDecodeEntry("stgt.h", decodeStgtH),
	newDecodeEntry("stgt.w", decodeStgtW),
	newDecodeEntry("stgt.d", decodeStgtD),
	newDecodeEntry("stle.b", decodeStleB),
	newDecodeEntry("stle.h", decodeStleH),
	newDecodeEntry("stle.w", decodeStleW),
	newDecodeEntry("stle.d", decodeStleD),
	newDecodeEntry("ll.w", decodeLlW),
	newDecodeEntry("ll.d", decodeLlD),
	newDecodeEntry("sc.w", decodeScW),
	newDecodeEntry("sc.d", decodeScD),
	newDecodeEntry("llacq.w", decodeLlacqW),
	newDecodeEntry("llacq.d", decodeLlacqD),
	newDecodeEntry("screl.w", decodeScrelW),
	newDecodeEntry("screl.d", decodeScrelD),
	newDecodeEntry("sc.q", decodeScQ),
	newDecodeEntry("amcas.b", decodeAmcasB),
	newDecodeEntry("amcas.h", decodeAmcasH),
	newDecodeEntry("amcas.w", decodeAmcasW),
	newDecodeEntry("amcas.d", decodeAmcasD),
	newDecodeEntry("amcas_db.b", decodeAmcasDbB),
	newDecodeEntry("amcas_db.h", decodeAmcasDbH),
	newDecodeEntry("amcas_db.w", decodeAmcasDbW),
	newDecodeEntry("amcas_db.d", decodeAmcasDbD),
	newDecodeEntry("amswap.b", decodeAmswapB),
	newDecodeEntry("amswap.h", decodeAmswapH),
	newDecodeEntry("amswap.w", decodeAmswapW),
	newDecodeEntry("amswap.d", decodeAmswapD),
	newDecodeEntry("amswap_db.b", decodeAmswapDbB),
	newDecodeEntry("amswap_db.h", decodeAmswapDbH),
	newDecodeEntry("amswap_db.w", decodeAmswapDbW),
	newDecodeEntry("amswap_db.d", decodeAmswapDbD),
	newDecodeEntry("amadd.b", decodeAmaddB),
	newDecodeEntry("amadd.h", decodeAmaddH),
	newDecodeEntry("amadd.w", decodeAmaddW),
	newDecodeEntry("amadd.d", decodeAmaddD),
	newDecodeEntry("amadd_db.b", decodeAmaddDbB),
	newDecodeEntry("amadd_db.h", decodeAmaddDbH),
	newDecodeEntry("amadd_db.w", decodeAmaddDbW),
	newDecodeEntry("amadd_db.d", decodeAmaddDbD),
	newDecodeEntry("amand.w", decodeAmandW),
	newDecodeEntry("amand.d", decodeAmandD),
	newDecodeEntry("amand_db.w", decodeAmandDbW),
	newDecodeEntry("amand_db.d", decodeAmandDbD),
	newDecodeEntry("amor.w", decodeAmorW),
	newDecodeEntry("amor.d", decodeAmorD),
	newDecodeEntry("amor_db.w", decodeAmorDbW),
	newDecodeEntry("amor_db.d", decodeAmorDbD),
	newDecodeEntry("amxor.w", decodeAmxorW),
	newDecodeEntry("amxor.d", decodeAmxorD),
	newDecodeEntry("amxor_db.w", decodeAmxorDbW),
	newDecodeEntry("amxor_db.d", decodeAmxorDbD),
	newDecodeEntry("ammax.w", decodeAmmaxW),
	newDecodeEntry("ammax.d", decodeAmmaxD),
	newDecodeEntry("ammax.wu", decodeAmmaxWu),
	newDecodeEntry("ammax.du", decodeAmmaxDu),
	newDecodeEntry("ammax_db.w", decodeAmmaxDbW),
	newDecodeEntry("ammax_db.d", decodeAmmaxDbD),
	newDecodeEntry("ammax_db.wu", decodeAmmaxDbWu),
	newDecodeEntry("ammax_db.du", decodeAmmaxDbDu),
	newDecodeEntry("ammin.w", decodeAmminW),
	newDecodeEntry("ammin.d", decodeAmminD),
	newDecodeEntry("ammin.wu", decodeAmminWu),
	newDecodeEntry("ammin.du", decodeAmminDu),
	newDecodeEntry("ammin_db.w", decodeAmminDbW),
	newDecodeEntry("ammin_db.d", decodeAmminDbD),
	newDecodeEntry("ammin_db.wu", decodeAmminDbWu),
	newDecodeEntry("ammin_db.du", decodeAmminDbDu),
	newDecodeEntry("csrrd", decodeCsrrd),
	newDecodeEntry("csrwr", decodeCsrwr),
	newDecodeEntry("csrxchg", decodeCsrxchg),
	newDecodeEntry("cacop", decodeCacop),
	newDecodeEntry("lddir", decodeLddir),
	newDecodeEntry("ldpte", decodeLdpte),
	newDecodeEntry("iocsrrd.b", decodeIocsrrdB),
	newDecodeEntry("iocsrrd.h", decodeIocsrrdH),
	newDecodeEntry("iocsrrd.w", decodeIocsrrdW),
	newDecodeEntry("iocsrrd.d", decodeIocsrrdD),
	newDecodeEntry("iocsrwr.b", decodeIocsrwrB),
	newDecodeEntry("iocsrwr.h", decodeIocsrwrH),
	newDecodeEntry("iocsrwr.w", decodeIocsrwrW),
	newDecodeEntry("iocsrwr.d", decodeIocsrwrD),
	newDecodeEntry("tlbclr", decodeTlbclr),
	newDecodeEntry("tlbflush", decodeTlbflush),
	newDecodeEntry("tlbsrch", decodeTlbsrch),
	newDecodeEntry("tlbrd", decodeTlbrd),
	newDecodeEntry("tlbwr", decodeTlbwr),
	newDecodeEntry("tlbfill", decodeTlbfill),
	newDecodeEntry("ertn", decodeErtn),
	newDecodeEntry("idle", decodeIdle),
	newDecodeEntry("invtlb", decodeInvtlb),
}

func decodeAddD(w uint32) Instr {
	return AddD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeAddW(w uint32) Instr {
	return AddW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeAddiD(w uint32) Instr {
	return AddiD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeAddiW(w uint32) Instr {
	return AddiW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeAddu16iD(w uint32) Instr {
	return Addu16iD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 16)),
	}
}

func decodeAlslD(w uint32) Instr {
	return AlslD{
		base:  newBase(w),
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAlslW(w uint32) Instr {
	return AlslW{
		base:  newBase(w),
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAlslWu(w uint32) Instr {
	return AlslWu{
		base:  newBase(w),
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAmaddB(w uint32) Instr {
	return AmaddB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddD(w uint32) Instr {
	return AmaddD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbB(w uint32) Instr {
	return AmaddDbB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbD(w uint32) Instr {
	return AmaddDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbH(w uint32) Instr {
	return AmaddDbH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbW(w uint32) Instr {
	return AmaddDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddH(w uint32) Instr {
	return AmaddH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddW(w uint32) Instr {
	return AmaddW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandD(w uint32) Instr {
	return AmandD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandDbD(w uint32) Instr {
	return AmandDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandDbW(w uint32) Instr {
	return AmandDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandW(w uint32) Instr {
	return AmandW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasB(w uint32) Instr {
	return AmcasB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasD(w uint32) Instr {
	return AmcasD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbB(w uint32) Instr {
	return AmcasDbB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbD(w uint32) Instr {
	return AmcasDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbH(w uint32) Instr {
	return AmcasDbH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbW(w uint32) Instr {
	return AmcasDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasH(w uint32) Instr {
	return AmcasH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasW(w uint32) Instr {
	return AmcasW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxD(w uint32) Instr {
	return AmmaxD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbD(w uint32) Instr {
	return AmmaxDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbDu(w uint32) Instr {
	return AmmaxDbDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbW(w uint32) Instr {
	return AmmaxDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbWu(w uint32) Instr {
	return AmmaxDbWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDu(w uint32) Instr {
	return AmmaxDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxW(w uint32) Instr {
	return AmmaxW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxWu(w uint32) Instr {
	return AmmaxWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminD(w uint32) Instr {
	return AmminD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbD(w uint32) Instr {
	return AmminDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbDu(w uint32) Instr {
	return AmminDbDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbW(w uint32) Instr {
	return AmminDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbWu(w uint32) Instr {
	return AmminDbWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDu(w uint32) Instr {
	return AmminDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminW(w uint32) Instr {
	return AmminW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminWu(w uint32) Instr {
	return AmminWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorD(w uint32) Instr {
	return AmorD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorDbD(w uint32) Instr {
	return AmorDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorDbW(w uint32) Instr {
	return AmorDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorW(w uint32) Instr {
	return AmorW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapB(w uint32) Instr {
	return AmswapB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapD(w uint32) Instr {
	return AmswapD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbB(w uint32) Instr {
	return AmswapDbB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbD(w uint32) Instr {
	return AmswapDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbH(w uint32) Instr {
	return AmswapDbH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbW(w uint32) Instr {
	return AmswapDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapH(w uint32) Instr {
	return AmswapH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapW(w uint32) Instr {
	return AmswapW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorD(w uint32) Instr {
	return AmxorD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorDbD(w uint32) Instr {
	return AmxorDbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorDbW(w uint32) Instr {
	return AmxorDbW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorW(w uint32) Instr {
	return AmxorW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeAnd(w uint32) Instr {
	return And{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeAndi(w uint32) Instr {
	return Andi{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 12))),
	}
}

func decodeAndn(w uint32) Instr {
	return Andn{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeAsrtgtD(w uint32) Instr {
	return AsrtgtD{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeAsrtleD(w uint32) Instr {
	return AsrtleD{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeB(w uint32) Instr {
	return B{
		base: newBase(w),
		off:  immNum(d10k16Imm(w) << 2),
	}
}

func decodeBeq(w uint32) Instr {
	return Beq{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBeqz(w uint32) Instr {
	return Beqz{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(d5k16Imm(w) << 2),
	}
}

func decodeBge(w uint32) Instr {
	return Bge{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBgeu(w uint32) Instr {
	return Bgeu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBl(w uint32) Instr {
	return Bl{
		base: newBase(w),
		off:  immNum(d10k16Imm(w) << 2),
	}
}

func decodeBlt(w uint32) Instr {
	return Blt{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBltu(w uint32) Instr {
	return Bltu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBne(w uint32) Instr {
	return Bne{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBnez(w uint32) Instr {
	return Bnez{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(d5k16Imm(w) << 2),
	}
}

func decodeBreak(w uint32) Instr {
	return Break{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeBstrinsD(w uint32) Instr {
	return BstrinsD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 6))),
		lsb:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeBstrinsW(w uint32) Instr {
	return BstrinsW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 5))),
		lsb:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeBstrpickD(w uint32) Instr {
	return BstrpickD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 6))),
		lsb:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeBstrpickW(w uint32) Instr {
	return BstrpickW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 5))),
		lsb:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeBytepickD(w uint32) Instr {
	return BytepickD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		sel:  immNum(int64(uField(w, 15, 3))),
	}
}

func decodeBytepickW(w uint32) Instr {
	return BytepickW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		sel:  immNum(int64(uField(w, 15, 2))),
	}
}

func decodeCacop(w uint32) Instr {
	return Cacop{
		base: newBase(w),
		op:   immNum(int64(uField(w, 0, 5))),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func decodeCloD(w uint32) Instr {
	return CloD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCloW(w uint32) Instr {
	return CloW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeClzD(w uint32) Instr {
	return ClzD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeClzW(w uint32) Instr {
	return ClzW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCpucfg(w uint32) Instr {
	return Cpucfg{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCrcWBW(w uint32) Instr {
	return CrcWBW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWDW(w uint32) Instr {
	return CrcWDW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWHW(w uint32) Instr {
	return CrcWHW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWWW(w uint32) Instr {
	return CrcWWW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWBW(w uint32) Instr {
	return CrccWBW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWDW(w uint32) Instr {
	return CrccWDW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWHW(w uint32) Instr {
	return CrccWHW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWWW(w uint32) Instr {
	return CrccWWW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeCsrrd(w uint32) Instr {
	return Csrrd{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCsrwr(w uint32) Instr {
	return Csrwr{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCsrxchg(w uint32) Instr {
	return Csrxchg{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		csr:  immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCtoD(w uint32) Instr {
	return CtoD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCtoW(w uint32) Instr {
	return CtoW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCtzD(w uint32) Instr {
	return CtzD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeCtzW(w uint32) Instr {
	return CtzW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeDbar(w uint32) Instr {
	return Dbar{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeDbcl(w uint32) Instr {
	return Dbcl{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeDivD(w uint32) Instr {
	return DivD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeDivDu(w uint32) Instr {
	return DivDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeDivW(w uint32) Instr {
	return DivW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeDivWu(w uint32) Instr {
	return DivWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeErtn(w uint32) Instr {
	return Ertn{
		base: newBase(w),
	}
}

func decodeExtWB(w uint32) Instr {
	return ExtWB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeExtWH(w uint32) Instr {
	return ExtWH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIbar(w uint32) Instr {
	return Ibar{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeIdle(w uint32) Instr {
	return Idle{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeInvtlb(w uint32) Instr {
	return Invtlb{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		op:   immNum(int64(uField(w, 0, 5))),
	}
}

func decodeIocsrrdB(w uint32) Instr {
	return IocsrrdB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdD(w uint32) Instr {
	return IocsrrdD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdH(w uint32) Instr {
	return IocsrrdH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdW(w uint32) Instr {
	return IocsrrdW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrB(w uint32) Instr {
	return IocsrwrB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrD(w uint32) Instr {
	return IocsrwrD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrH(w uint32) Instr {
	return IocsrwrH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrW(w uint32) Instr {
	return IocsrwrW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeJirl(w uint32) Instr {
	return Jirl{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func decodeLdB(w uint32) Instr {
	return LdB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdBu(w uint32) Instr {
	return LdBu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdD(w uint32) Instr {
	return LdD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdH(w uint32) Instr {
	return LdH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdHu(w uint32) Instr {
	return LdHu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdW(w uint32) Instr {
	return LdW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeLdWu(w uint32) Instr {
	return LdWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeLddir(w uint32) Instr {
	return Lddir{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 8))),
	}
}

func decodeLdgtB(w uint32) Instr {
	return LdgtB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtD(w uint32) Instr {
	return LdgtD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtH(w uint32) Instr {
	return LdgtH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtW(w uint32) Instr {
	return LdgtW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleB(w uint32) Instr {
	return LdleB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleD(w uint32) Instr {
	return LdleD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleH(w uint32) Instr {
	return LdleH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleW(w uint32) Instr {
	return LdleW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdpte(w uint32) Instr {
	return Ldpte{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 8))),
	}
}

func decodeLdptrD(w uint32) Instr {
	return LdptrD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLdptrW(w uint32) Instr {
	return LdptrW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLdxB(w uint32) Instr {
	return LdxB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxBu(w uint32) Instr {
	return LdxBu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxD(w uint32) Instr {
	return LdxD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxH(w uint32) Instr {
	return LdxH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxHu(w uint32) Instr {
	return LdxHu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxW(w uint32) Instr {
	return LdxW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxWu(w uint32) Instr {
	return LdxWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeLlD(w uint32) Instr {
	return LlD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLlW(w uint32) Instr {
	return LlW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLlacqD(w uint32) Instr {
	return LlacqD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeLlacqW(w uint32) Instr {
	return LlacqW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeLu12iW(w uint32) Instr {
	return Lu12iW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodeLu32iD(w uint32) Instr {
	return Lu32iD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodeLu52iD(w uint32) Instr {
	return Lu52iD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeMaskeqz(w uint32) Instr {
	return Maskeqz{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMasknez(w uint32) Instr {
	return Masknez{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeModD(w uint32) Instr {
	return ModD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeModDu(w uint32) Instr {
	return ModDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeModW(w uint32) Instr {
	return ModW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeModWu(w uint32) Instr {
	return ModWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulD(w uint32) Instr {
	return MulD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulW(w uint32) Instr {
	return MulW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhD(w uint32) Instr {
	return MulhD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhDu(w uint32) Instr {
	return MulhDu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhW(w uint32) Instr {
	return MulhW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhWu(w uint32) Instr {
	return MulhWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulwDW(w uint32) Instr {
	return MulwDW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeMulwDWu(w uint32) Instr {
	return MulwDWu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeNor(w uint32) Instr {
	return Nor{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeOr(w uint32) Instr {
	return Or{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeOri(w uint32) Instr {
	return Ori{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 12))),
	}
}

func decodeOrn(w uint32) Instr {
	return Orn{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodePcaddi(w uint32) Instr {
	return Pcaddi{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodePcaddu12i(w uint32) Instr {
	return Pcaddu12i{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodePcaddu18i(w uint32) Instr {
	return Pcaddu18i{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodePcalau12i(w uint32) Instr {
	return Pcalau12i{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func decodePreld(w uint32) Instr {
	return Preld{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
		off:  immNum(sField(w, 10, 12)),
	}
}

func decodePreldx(w uint32) Instr {
	return Preldx{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
	}
}

func decodeRdtimeD(w uint32) Instr {
	return RdtimeD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRdtimehW(w uint32) Instr {
	return RdtimehW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRdtimelW(w uint32) Instr {
	return RdtimelW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb2H(w uint32) Instr {
	return Revb2H{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb2W(w uint32) Instr {
	return Revb2W{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb4H(w uint32) Instr {
	return Revb4H{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbD(w uint32) Instr {
	return RevbD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbit4B(w uint32) Instr {
	return Revbit4B{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbit8B(w uint32) Instr {
	return Revbit8B{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbitD(w uint32) Instr {
	return RevbitD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbitW(w uint32) Instr {
	return RevbitW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevh2W(w uint32) Instr {
	return Revh2W{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRevhD(w uint32) Instr {
	return RevhD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeRotrD(w uint32) Instr {
	return RotrD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeRotrW(w uint32) Instr {
	return RotrW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeRotriD(w uint32) Instr {
	return RotriD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeRotriW(w uint32) Instr {
	return RotriW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeScD(w uint32) Instr {
	return ScD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeScQ(w uint32) Instr {
	return ScQ{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeScW(w uint32) Instr {
	return ScW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeScrelD(w uint32) Instr {
	return ScrelD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeScrelW(w uint32) Instr {
	return ScrelW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func decodeSllD(w uint32) Instr {
	return SllD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSllW(w uint32) Instr {
	return SllW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSlliD(w uint32) Instr {
	return SlliD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSlliW(w uint32) Instr {
	return SlliW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeSlt(w uint32) Instr {
	return Slt{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSlti(w uint32) Instr {
	return Slti{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeSltu(w uint32) Instr {
	return Sltu{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSltui(w uint32) Instr {
	return Sltui{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeSraD(w uint32) Instr {
	return SraD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSraW(w uint32) Instr {
	return SraW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSraiD(w uint32) Instr {
	return SraiD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSraiW(w uint32) Instr {
	return SraiW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeSrlD(w uint32) Instr {
	return SrlD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSrlW(w uint32) Instr {
	return SrlW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSrliD(w uint32) Instr {
	return SrliD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSrliW(w uint32) Instr {
	return SrliW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 5))),
	}
}

func decodeStB(w uint32) Instr {
	return StB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeStD(w uint32) Instr {
	return StD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeStH(w uint32) Instr {
	return StH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeStW(w uint32) Instr {
	return StW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(sField(w, 10, 12)),
	}
}

func decodeStgtB(w uint32) Instr {
	return StgtB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtD(w uint32) Instr {
	return StgtD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtH(w uint32) Instr {
	return StgtH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtW(w uint32) Instr {
	return StgtW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStleB(w uint32) Instr {
	return StleB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStleD(w uint32) Instr {
	return StleD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStleH(w uint32) Instr {
	return StleH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStleW(w uint32) Instr {
	return StleW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStptrD(w uint32) Instr {
	return StptrD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeStptrW(w uint32) Instr {
	return StptrW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 14) << 2),
	}
}

func decodeStxB(w uint32) Instr {
	return StxB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStxD(w uint32) Instr {
	return StxD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStxH(w uint32) Instr {
	return StxH{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeStxW(w uint32) Instr {
	return StxW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSubD(w uint32) Instr {
	return SubD{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSubW(w uint32) Instr {
	return SubW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeSyscall(w uint32) Instr {
	return Syscall{
		base: newBase(w),
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeTlbclr(w uint32) Instr {
	return Tlbclr{
		base: newBase(w),
	}
}

func decodeTlbfill(w uint32) Instr {
	return Tlbfill{
		base: newBase(w),
	}
}

func decodeTlbflush(w uint32) Instr {
	return Tlbflush{
		base: newBase(w),
	}
}

func decodeTlbrd(w uint32) Instr {
	return Tlbrd{
		base: newBase(w),
	}
}

func decodeTlbsrch(w uint32) Instr {
	return Tlbsrch{
		base: newBase(w),
	}
}

func decodeTlbwr(w uint32) Instr {
	return Tlbwr{
		base: newBase(w),
	}
}

func decodeXor(w uint32) Instr {
	return Xor{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func decodeXori(w uint32) Instr {
	return Xori{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		imm:  immNum(int64(uField(w, 10, 12))),
	}
}
