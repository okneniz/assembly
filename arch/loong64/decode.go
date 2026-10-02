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
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeAddW(w uint32) Instr {
	return AddW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeAddiD(w uint32) Instr {
	return AddiD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeAddiW(w uint32) Instr {
	return AddiW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeAddu16iD(w uint32) Instr {
	return Addu16iD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 16)),
	}
}

func decodeAlslD(w uint32) Instr {
	return AlslD{
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAlslW(w uint32) Instr {
	return AlslW{
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAlslWu(w uint32) Instr {
	return AlslWu{
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func decodeAmaddB(w uint32) Instr {
	return AmaddB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddD(w uint32) Instr {
	return AmaddD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbB(w uint32) Instr {
	return AmaddDbB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbD(w uint32) Instr {
	return AmaddDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbH(w uint32) Instr {
	return AmaddDbH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddDbW(w uint32) Instr {
	return AmaddDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddH(w uint32) Instr {
	return AmaddH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmaddW(w uint32) Instr {
	return AmaddW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandD(w uint32) Instr {
	return AmandD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandDbD(w uint32) Instr {
	return AmandDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandDbW(w uint32) Instr {
	return AmandDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmandW(w uint32) Instr {
	return AmandW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasB(w uint32) Instr {
	return AmcasB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasD(w uint32) Instr {
	return AmcasD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbB(w uint32) Instr {
	return AmcasDbB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbD(w uint32) Instr {
	return AmcasDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbH(w uint32) Instr {
	return AmcasDbH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasDbW(w uint32) Instr {
	return AmcasDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasH(w uint32) Instr {
	return AmcasH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmcasW(w uint32) Instr {
	return AmcasW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxD(w uint32) Instr {
	return AmmaxD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbD(w uint32) Instr {
	return AmmaxDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbDu(w uint32) Instr {
	return AmmaxDbDu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbW(w uint32) Instr {
	return AmmaxDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDbWu(w uint32) Instr {
	return AmmaxDbWu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxDu(w uint32) Instr {
	return AmmaxDu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxW(w uint32) Instr {
	return AmmaxW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmmaxWu(w uint32) Instr {
	return AmmaxWu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminD(w uint32) Instr {
	return AmminD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbD(w uint32) Instr {
	return AmminDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbDu(w uint32) Instr {
	return AmminDbDu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbW(w uint32) Instr {
	return AmminDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDbWu(w uint32) Instr {
	return AmminDbWu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminDu(w uint32) Instr {
	return AmminDu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminW(w uint32) Instr {
	return AmminW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmminWu(w uint32) Instr {
	return AmminWu{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorD(w uint32) Instr {
	return AmorD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorDbD(w uint32) Instr {
	return AmorDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorDbW(w uint32) Instr {
	return AmorDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmorW(w uint32) Instr {
	return AmorW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapB(w uint32) Instr {
	return AmswapB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapD(w uint32) Instr {
	return AmswapD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbB(w uint32) Instr {
	return AmswapDbB{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbD(w uint32) Instr {
	return AmswapDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbH(w uint32) Instr {
	return AmswapDbH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapDbW(w uint32) Instr {
	return AmswapDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapH(w uint32) Instr {
	return AmswapH{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmswapW(w uint32) Instr {
	return AmswapW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorD(w uint32) Instr {
	return AmxorD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorDbD(w uint32) Instr {
	return AmxorDbD{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorDbW(w uint32) Instr {
	return AmxorDbW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAmxorW(w uint32) Instr {
	return AmxorW{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeAnd(w uint32) Instr {
	return And{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeAndi(w uint32) Instr {
	return Andi{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 12))),
	}
}

func decodeAndn(w uint32) Instr {
	return Andn{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeAsrtgtD(w uint32) Instr {
	return AsrtgtD{
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeAsrtleD(w uint32) Instr {
	return AsrtleD{
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeB(w uint32) Instr {
	return B{
		off: immNum(d10k16Imm(w) << 2),
	}
}

func decodeBeq(w uint32) Instr {
	return Beq{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBeqz(w uint32) Instr {
	return Beqz{
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(d5k16Imm(w) << 2),
	}
}

func decodeBge(w uint32) Instr {
	return Bge{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBgeu(w uint32) Instr {
	return Bgeu{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBl(w uint32) Instr {
	return Bl{
		off: immNum(d10k16Imm(w) << 2),
	}
}

func decodeBlt(w uint32) Instr {
	return Blt{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBltu(w uint32) Instr {
	return Bltu{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBne(w uint32) Instr {
	return Bne{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeBnez(w uint32) Instr {
	return Bnez{
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(d5k16Imm(w) << 2),
	}
}

func decodeBreak(w uint32) Instr {
	return Break{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeBstrinsD(w uint32) Instr {
	return BstrinsD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		msb: immNum(int64(uField(w, 16, 6))),
		lsb: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeBstrinsW(w uint32) Instr {
	return BstrinsW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		msb: immNum(int64(uField(w, 16, 5))),
		lsb: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeBstrpickD(w uint32) Instr {
	return BstrpickD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		msb: immNum(int64(uField(w, 16, 6))),
		lsb: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeBstrpickW(w uint32) Instr {
	return BstrpickW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		msb: immNum(int64(uField(w, 16, 5))),
		lsb: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeBytepickD(w uint32) Instr {
	return BytepickD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		rk:  uint8(w >> 10 & 0x1f),
		sel: immNum(int64(uField(w, 15, 3))),
	}
}

func decodeBytepickW(w uint32) Instr {
	return BytepickW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		rk:  uint8(w >> 10 & 0x1f),
		sel: immNum(int64(uField(w, 15, 2))),
	}
}

func decodeCacop(w uint32) Instr {
	return Cacop{
		op:  immNum(int64(uField(w, 0, 5))),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 12)),
	}
}

func decodeCloD(w uint32) Instr {
	return CloD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCloW(w uint32) Instr {
	return CloW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeClzD(w uint32) Instr {
	return ClzD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeClzW(w uint32) Instr {
	return ClzW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCpucfg(w uint32) Instr {
	return Cpucfg{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCrcWBW(w uint32) Instr {
	return CrcWBW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWDW(w uint32) Instr {
	return CrcWDW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWHW(w uint32) Instr {
	return CrcWHW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrcWWW(w uint32) Instr {
	return CrcWWW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWBW(w uint32) Instr {
	return CrccWBW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWDW(w uint32) Instr {
	return CrccWDW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWHW(w uint32) Instr {
	return CrccWHW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCrccWWW(w uint32) Instr {
	return CrccWWW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeCsrrd(w uint32) Instr {
	return Csrrd{
		rd:  uint8(w & 0x1f),
		csr: immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCsrwr(w uint32) Instr {
	return Csrwr{
		rd:  uint8(w & 0x1f),
		csr: immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCsrxchg(w uint32) Instr {
	return Csrxchg{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		csr: immNum(int64(uField(w, 10, 14))),
	}
}

func decodeCtoD(w uint32) Instr {
	return CtoD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCtoW(w uint32) Instr {
	return CtoW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCtzD(w uint32) Instr {
	return CtzD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeCtzW(w uint32) Instr {
	return CtzW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeDbar(w uint32) Instr {
	return Dbar{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeDbcl(w uint32) Instr {
	return Dbcl{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeDivD(w uint32) Instr {
	return DivD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeDivDu(w uint32) Instr {
	return DivDu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeDivW(w uint32) Instr {
	return DivW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeDivWu(w uint32) Instr {
	return DivWu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeErtn(w uint32) Instr {
	return Ertn{}
}

func decodeExtWB(w uint32) Instr {
	return ExtWB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeExtWH(w uint32) Instr {
	return ExtWH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIbar(w uint32) Instr {
	return Ibar{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeIdle(w uint32) Instr {
	return Idle{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeInvtlb(w uint32) Instr {
	return Invtlb{
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		op: immNum(int64(uField(w, 0, 5))),
	}
}

func decodeIocsrrdB(w uint32) Instr {
	return IocsrrdB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdD(w uint32) Instr {
	return IocsrrdD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdH(w uint32) Instr {
	return IocsrrdH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrrdW(w uint32) Instr {
	return IocsrrdW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrB(w uint32) Instr {
	return IocsrwrB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrD(w uint32) Instr {
	return IocsrwrD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrH(w uint32) Instr {
	return IocsrwrH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeIocsrwrW(w uint32) Instr {
	return IocsrwrW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeJirl(w uint32) Instr {
	return Jirl{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 16) << 2),
	}
}

func decodeLdB(w uint32) Instr {
	return LdB{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeLdBu(w uint32) Instr {
	return LdBu{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 12)),
	}
}

func decodeLdD(w uint32) Instr {
	return LdD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeLdH(w uint32) Instr {
	return LdH{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeLdHu(w uint32) Instr {
	return LdHu{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 12)),
	}
}

func decodeLdW(w uint32) Instr {
	return LdW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeLdWu(w uint32) Instr {
	return LdWu{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeLddir(w uint32) Instr {
	return Lddir{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 8))),
	}
}

func decodeLdgtB(w uint32) Instr {
	return LdgtB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtD(w uint32) Instr {
	return LdgtD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtH(w uint32) Instr {
	return LdgtH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdgtW(w uint32) Instr {
	return LdgtW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleB(w uint32) Instr {
	return LdleB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleD(w uint32) Instr {
	return LdleD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleH(w uint32) Instr {
	return LdleH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdleW(w uint32) Instr {
	return LdleW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdpte(w uint32) Instr {
	return Ldpte{
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 8))),
	}
}

func decodeLdptrD(w uint32) Instr {
	return LdptrD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLdptrW(w uint32) Instr {
	return LdptrW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLdxB(w uint32) Instr {
	return LdxB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxBu(w uint32) Instr {
	return LdxBu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxD(w uint32) Instr {
	return LdxD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxH(w uint32) Instr {
	return LdxH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxHu(w uint32) Instr {
	return LdxHu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxW(w uint32) Instr {
	return LdxW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLdxWu(w uint32) Instr {
	return LdxWu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeLlD(w uint32) Instr {
	return LlD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLlW(w uint32) Instr {
	return LlW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeLlacqD(w uint32) Instr {
	return LlacqD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeLlacqW(w uint32) Instr {
	return LlacqW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeLu12iW(w uint32) Instr {
	return Lu12iW{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodeLu32iD(w uint32) Instr {
	return Lu32iD{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodeLu52iD(w uint32) Instr {
	return Lu52iD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeMaskeqz(w uint32) Instr {
	return Maskeqz{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMasknez(w uint32) Instr {
	return Masknez{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeModD(w uint32) Instr {
	return ModD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeModDu(w uint32) Instr {
	return ModDu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeModW(w uint32) Instr {
	return ModW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeModWu(w uint32) Instr {
	return ModWu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulD(w uint32) Instr {
	return MulD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulW(w uint32) Instr {
	return MulW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhD(w uint32) Instr {
	return MulhD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhDu(w uint32) Instr {
	return MulhDu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhW(w uint32) Instr {
	return MulhW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulhWu(w uint32) Instr {
	return MulhWu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulwDW(w uint32) Instr {
	return MulwDW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeMulwDWu(w uint32) Instr {
	return MulwDWu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeNor(w uint32) Instr {
	return Nor{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeOr(w uint32) Instr {
	return Or{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeOri(w uint32) Instr {
	return Ori{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 12))),
	}
}

func decodeOrn(w uint32) Instr {
	return Orn{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodePcaddi(w uint32) Instr {
	return Pcaddi{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodePcaddu12i(w uint32) Instr {
	return Pcaddu12i{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodePcaddu18i(w uint32) Instr {
	return Pcaddu18i{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodePcalau12i(w uint32) Instr {
	return Pcalau12i{
		rd:  uint8(w & 0x1f),
		imm: immNum(sField(w, 5, 20)),
	}
}

func decodePreld(w uint32) Instr {
	return Preld{
		rj:   uint8(w >> 5 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
		off:  immNum(sField(w, 10, 12)),
	}
}

func decodePreldx(w uint32) Instr {
	return Preldx{
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		hint: immNum(int64(uField(w, 0, 5))),
	}
}

func decodeRdtimeD(w uint32) Instr {
	return RdtimeD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRdtimehW(w uint32) Instr {
	return RdtimehW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRdtimelW(w uint32) Instr {
	return RdtimelW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb2H(w uint32) Instr {
	return Revb2H{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb2W(w uint32) Instr {
	return Revb2W{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevb4H(w uint32) Instr {
	return Revb4H{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbD(w uint32) Instr {
	return RevbD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbit4B(w uint32) Instr {
	return Revbit4B{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbit8B(w uint32) Instr {
	return Revbit8B{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbitD(w uint32) Instr {
	return RevbitD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevbitW(w uint32) Instr {
	return RevbitW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevh2W(w uint32) Instr {
	return Revh2W{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRevhD(w uint32) Instr {
	return RevhD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeRotrD(w uint32) Instr {
	return RotrD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeRotrW(w uint32) Instr {
	return RotrW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeRotriD(w uint32) Instr {
	return RotriD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeRotriW(w uint32) Instr {
	return RotriW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeScD(w uint32) Instr {
	return ScD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeScQ(w uint32) Instr {
	return ScQ{
		rd: uint8(w & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeScW(w uint32) Instr {
	return ScW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeScrelD(w uint32) Instr {
	return ScrelD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeScrelW(w uint32) Instr {
	return ScrelW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
	}
}

func decodeSllD(w uint32) Instr {
	return SllD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSllW(w uint32) Instr {
	return SllW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSlliD(w uint32) Instr {
	return SlliD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSlliW(w uint32) Instr {
	return SlliW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeSlt(w uint32) Instr {
	return Slt{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSlti(w uint32) Instr {
	return Slti{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeSltu(w uint32) Instr {
	return Sltu{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSltui(w uint32) Instr {
	return Sltui{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeSraD(w uint32) Instr {
	return SraD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSraW(w uint32) Instr {
	return SraW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSraiD(w uint32) Instr {
	return SraiD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSraiW(w uint32) Instr {
	return SraiW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeSrlD(w uint32) Instr {
	return SrlD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSrlW(w uint32) Instr {
	return SrlW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSrliD(w uint32) Instr {
	return SrliD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 6))),
	}
}

func decodeSrliW(w uint32) Instr {
	return SrliW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 5))),
	}
}

func decodeStB(w uint32) Instr {
	return StB{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeStD(w uint32) Instr {
	return StD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeStH(w uint32) Instr {
	return StH{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeStW(w uint32) Instr {
	return StW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(sField(w, 10, 12)),
	}
}

func decodeStgtB(w uint32) Instr {
	return StgtB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtD(w uint32) Instr {
	return StgtD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtH(w uint32) Instr {
	return StgtH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStgtW(w uint32) Instr {
	return StgtW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStleB(w uint32) Instr {
	return StleB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStleD(w uint32) Instr {
	return StleD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStleH(w uint32) Instr {
	return StleH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStleW(w uint32) Instr {
	return StleW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStptrD(w uint32) Instr {
	return StptrD{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeStptrW(w uint32) Instr {
	return StptrW{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		off: immNum(sField(w, 10, 14) << 2),
	}
}

func decodeStxB(w uint32) Instr {
	return StxB{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStxD(w uint32) Instr {
	return StxD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStxH(w uint32) Instr {
	return StxH{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeStxW(w uint32) Instr {
	return StxW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSubD(w uint32) Instr {
	return SubD{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSubW(w uint32) Instr {
	return SubW{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeSyscall(w uint32) Instr {
	return Syscall{
		code: immNum(int64(uField(w, 0, 15))),
	}
}

func decodeTlbclr(w uint32) Instr {
	return Tlbclr{}
}

func decodeTlbfill(w uint32) Instr {
	return Tlbfill{}
}

func decodeTlbflush(w uint32) Instr {
	return Tlbflush{}
}

func decodeTlbrd(w uint32) Instr {
	return Tlbrd{}
}

func decodeTlbsrch(w uint32) Instr {
	return Tlbsrch{}
}

func decodeTlbwr(w uint32) Instr {
	return Tlbwr{}
}

func decodeXor(w uint32) Instr {
	return Xor{
		rd: uint8(w & 0x1f),
		rj: uint8(w >> 5 & 0x1f),
		rk: uint8(w >> 10 & 0x1f),
	}
}

func decodeXori(w uint32) Instr {
	return Xori{
		rd:  uint8(w & 0x1f),
		rj:  uint8(w >> 5 & 0x1f),
		imm: immNum(int64(uField(w, 10, 12))),
	}
}
