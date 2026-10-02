package loong64

import "fmt"

// Builder — the construction vocabulary of the package: instruction
// constructors named after the mnemonic they build and operand role
// constructors. Stateless by design — a namespace, not a factory.
type Builder struct{}

// New — the Builder.
func New() Builder {
	return Builder{}
}

// AddD - add.d rd, rj, rk.
func (Builder) AddD(rd, rj, rk Reg) Instr {
	return AddD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// AddW - add.w rd, rj, rk.
func (Builder) AddW(rd, rj, rk Reg) Instr {
	return AddW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// AddiD - addi.d rd, rj, si12.
func (Builder) AddiD(rd, rj Reg, v Imm12) Instr {
	return AddiD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// AddiW - addi.w rd, rj, si12.
func (Builder) AddiW(rd, rj Reg, v Imm12) Instr {
	return AddiW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Addu16iD - addu16i.d rd, rj, si16.
func (Builder) Addu16iD(rd, rj Reg, v Imm16) Instr {
	return Addu16iD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// AlslD - alsl.d rd, rj, rk, shift.
func (Builder) AlslD(rd, rj, rk Reg, shift Shift3) Instr {
	return AlslD{
		rd:    rd.Num(),
		rj:    rj.Num(),
		rk:    rk.Num(),
		shift: immNum(shift.Val()),
	}
}

// AlslW - alsl.w rd, rj, rk, shift.
func (Builder) AlslW(rd, rj, rk Reg, shift Shift3) Instr {
	return AlslW{
		rd:    rd.Num(),
		rj:    rj.Num(),
		rk:    rk.Num(),
		shift: immNum(shift.Val()),
	}
}

// AlslWu - alsl.wu rd, rj, rk, shift.
func (Builder) AlslWu(rd, rj, rk Reg, shift Shift3) Instr {
	return AlslWu{
		rd:    rd.Num(),
		rj:    rj.Num(),
		rk:    rk.Num(),
		shift: immNum(shift.Val()),
	}
}

// AmaddB - amadd.b rd, rk, rj.
func (Builder) AmaddB(rd, rk, rj Reg) Instr {
	return AmaddB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddD - amadd.d rd, rk, rj.
func (Builder) AmaddD(rd, rk, rj Reg) Instr {
	return AmaddD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddDbB - amadd_db.b rd, rk, rj.
func (Builder) AmaddDbB(rd, rk, rj Reg) Instr {
	return AmaddDbB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddDbD - amadd_db.d rd, rk, rj.
func (Builder) AmaddDbD(rd, rk, rj Reg) Instr {
	return AmaddDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddDbH - amadd_db.h rd, rk, rj.
func (Builder) AmaddDbH(rd, rk, rj Reg) Instr {
	return AmaddDbH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddDbW - amadd_db.w rd, rk, rj.
func (Builder) AmaddDbW(rd, rk, rj Reg) Instr {
	return AmaddDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddH - amadd.h rd, rk, rj.
func (Builder) AmaddH(rd, rk, rj Reg) Instr {
	return AmaddH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmaddW - amadd.w rd, rk, rj.
func (Builder) AmaddW(rd, rk, rj Reg) Instr {
	return AmaddW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmandD - amand.d rd, rk, rj.
func (Builder) AmandD(rd, rk, rj Reg) Instr {
	return AmandD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmandDbD - amand_db.d rd, rk, rj.
func (Builder) AmandDbD(rd, rk, rj Reg) Instr {
	return AmandDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmandDbW - amand_db.w rd, rk, rj.
func (Builder) AmandDbW(rd, rk, rj Reg) Instr {
	return AmandDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmandW - amand.w rd, rk, rj.
func (Builder) AmandW(rd, rk, rj Reg) Instr {
	return AmandW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasB - amcas.b rd, rk, rj.
func (Builder) AmcasB(rd, rk, rj Reg) Instr {
	return AmcasB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasD - amcas.d rd, rk, rj.
func (Builder) AmcasD(rd, rk, rj Reg) Instr {
	return AmcasD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasDbB - amcas_db.b rd, rk, rj.
func (Builder) AmcasDbB(rd, rk, rj Reg) Instr {
	return AmcasDbB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasDbD - amcas_db.d rd, rk, rj.
func (Builder) AmcasDbD(rd, rk, rj Reg) Instr {
	return AmcasDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasDbH - amcas_db.h rd, rk, rj.
func (Builder) AmcasDbH(rd, rk, rj Reg) Instr {
	return AmcasDbH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasDbW - amcas_db.w rd, rk, rj.
func (Builder) AmcasDbW(rd, rk, rj Reg) Instr {
	return AmcasDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasH - amcas.h rd, rk, rj.
func (Builder) AmcasH(rd, rk, rj Reg) Instr {
	return AmcasH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmcasW - amcas.w rd, rk, rj.
func (Builder) AmcasW(rd, rk, rj Reg) Instr {
	return AmcasW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxD - ammax.d rd, rk, rj.
func (Builder) AmmaxD(rd, rk, rj Reg) Instr {
	return AmmaxD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxDbD - ammax_db.d rd, rk, rj.
func (Builder) AmmaxDbD(rd, rk, rj Reg) Instr {
	return AmmaxDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxDbDu - ammax_db.du rd, rk, rj.
func (Builder) AmmaxDbDu(rd, rk, rj Reg) Instr {
	return AmmaxDbDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxDbW - ammax_db.w rd, rk, rj.
func (Builder) AmmaxDbW(rd, rk, rj Reg) Instr {
	return AmmaxDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxDbWu - ammax_db.wu rd, rk, rj.
func (Builder) AmmaxDbWu(rd, rk, rj Reg) Instr {
	return AmmaxDbWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxDu - ammax.du rd, rk, rj.
func (Builder) AmmaxDu(rd, rk, rj Reg) Instr {
	return AmmaxDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxW - ammax.w rd, rk, rj.
func (Builder) AmmaxW(rd, rk, rj Reg) Instr {
	return AmmaxW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmmaxWu - ammax.wu rd, rk, rj.
func (Builder) AmmaxWu(rd, rk, rj Reg) Instr {
	return AmmaxWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminD - ammin.d rd, rk, rj.
func (Builder) AmminD(rd, rk, rj Reg) Instr {
	return AmminD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminDbD - ammin_db.d rd, rk, rj.
func (Builder) AmminDbD(rd, rk, rj Reg) Instr {
	return AmminDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminDbDu - ammin_db.du rd, rk, rj.
func (Builder) AmminDbDu(rd, rk, rj Reg) Instr {
	return AmminDbDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminDbW - ammin_db.w rd, rk, rj.
func (Builder) AmminDbW(rd, rk, rj Reg) Instr {
	return AmminDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminDbWu - ammin_db.wu rd, rk, rj.
func (Builder) AmminDbWu(rd, rk, rj Reg) Instr {
	return AmminDbWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminDu - ammin.du rd, rk, rj.
func (Builder) AmminDu(rd, rk, rj Reg) Instr {
	return AmminDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminW - ammin.w rd, rk, rj.
func (Builder) AmminW(rd, rk, rj Reg) Instr {
	return AmminW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmminWu - ammin.wu rd, rk, rj.
func (Builder) AmminWu(rd, rk, rj Reg) Instr {
	return AmminWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmorD - amor.d rd, rk, rj.
func (Builder) AmorD(rd, rk, rj Reg) Instr {
	return AmorD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmorDbD - amor_db.d rd, rk, rj.
func (Builder) AmorDbD(rd, rk, rj Reg) Instr {
	return AmorDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmorDbW - amor_db.w rd, rk, rj.
func (Builder) AmorDbW(rd, rk, rj Reg) Instr {
	return AmorDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmorW - amor.w rd, rk, rj.
func (Builder) AmorW(rd, rk, rj Reg) Instr {
	return AmorW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapB - amswap.b rd, rk, rj.
func (Builder) AmswapB(rd, rk, rj Reg) Instr {
	return AmswapB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapD - amswap.d rd, rk, rj.
func (Builder) AmswapD(rd, rk, rj Reg) Instr {
	return AmswapD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapDbB - amswap_db.b rd, rk, rj.
func (Builder) AmswapDbB(rd, rk, rj Reg) Instr {
	return AmswapDbB{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapDbD - amswap_db.d rd, rk, rj.
func (Builder) AmswapDbD(rd, rk, rj Reg) Instr {
	return AmswapDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapDbH - amswap_db.h rd, rk, rj.
func (Builder) AmswapDbH(rd, rk, rj Reg) Instr {
	return AmswapDbH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapDbW - amswap_db.w rd, rk, rj.
func (Builder) AmswapDbW(rd, rk, rj Reg) Instr {
	return AmswapDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapH - amswap.h rd, rk, rj.
func (Builder) AmswapH(rd, rk, rj Reg) Instr {
	return AmswapH{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmswapW - amswap.w rd, rk, rj.
func (Builder) AmswapW(rd, rk, rj Reg) Instr {
	return AmswapW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmxorD - amxor.d rd, rk, rj.
func (Builder) AmxorD(rd, rk, rj Reg) Instr {
	return AmxorD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmxorDbD - amxor_db.d rd, rk, rj.
func (Builder) AmxorDbD(rd, rk, rj Reg) Instr {
	return AmxorDbD{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmxorDbW - amxor_db.w rd, rk, rj.
func (Builder) AmxorDbW(rd, rk, rj Reg) Instr {
	return AmxorDbW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// AmxorW - amxor.w rd, rk, rj.
func (Builder) AmxorW(rd, rk, rj Reg) Instr {
	return AmxorW{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// And - and rd, rj, rk.
func (Builder) And(rd, rj, rk Reg) Instr {
	return And{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Andi - andi rd, rj, ui12.
func (Builder) Andi(rd, rj Reg, v UImm12) Instr {
	return Andi{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Andn - andn rd, rj, rk.
func (Builder) Andn(rd, rj, rk Reg) Instr {
	return Andn{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// AsrtgtD - asrtgt.d rj, rk.
func (Builder) AsrtgtD(rj, rk Reg) Instr {
	return AsrtgtD{
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// AsrtleD - asrtle.d rj, rk.
func (Builder) AsrtleD(rj, rk Reg) Instr {
	return AsrtleD{
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// B - b offs (the pc-relative byte offset).
func (Builder) B(off int64) Instr {
	return B{
		off: immNum(off),
	}
}

// Beq - beq rj, rd, offs (the pc-relative byte offset).
func (Builder) Beq(rj, rd Reg, off int64) Instr {
	return Beq{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Beqz - beqz rj, offs (the pc-relative byte offset).
func (Builder) Beqz(rj Reg, off int64) Instr {
	return Beqz{
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bge - bge rj, rd, offs (the pc-relative byte offset).
func (Builder) Bge(rj, rd Reg, off int64) Instr {
	return Bge{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bgeu - bgeu rj, rd, offs (the pc-relative byte offset).
func (Builder) Bgeu(rj, rd Reg, off int64) Instr {
	return Bgeu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bl - bl offs (the pc-relative byte offset).
func (Builder) Bl(off int64) Instr {
	return Bl{
		off: immNum(off),
	}
}

// Blt - blt rj, rd, offs (the pc-relative byte offset).
func (Builder) Blt(rj, rd Reg, off int64) Instr {
	return Blt{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bltu - bltu rj, rd, offs (the pc-relative byte offset).
func (Builder) Bltu(rj, rd Reg, off int64) Instr {
	return Bltu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bne - bne rj, rd, offs (the pc-relative byte offset).
func (Builder) Bne(rj, rd Reg, off int64) Instr {
	return Bne{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Bnez - bnez rj, offs (the pc-relative byte offset).
func (Builder) Bnez(rj Reg, off int64) Instr {
	return Bnez{
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// Break - break code (a 15-bit code).
func (Builder) Break(code Code15) Instr {
	return Break{
		code: immNum(code.Val()),
	}
}

// BstrinsD - bstrins.d rd, rj, msb, lsb.
func (Builder) BstrinsD(rd, rj Reg, msb, lsb UImm6) Instr {
	return BstrinsD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

// BstrinsW - bstrins.w rd, rj, msb, lsb.
func (Builder) BstrinsW(rd, rj Reg, msb, lsb UImm5) Instr {
	return BstrinsW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

// BstrpickD - bstrpick.d rd, rj, msb, lsb.
func (Builder) BstrpickD(rd, rj Reg, msb, lsb UImm6) Instr {
	return BstrpickD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

// BstrpickW - bstrpick.w rd, rj, msb, lsb.
func (Builder) BstrpickW(rd, rj Reg, msb, lsb UImm5) Instr {
	return BstrpickW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

// BytepickD - bytepick.d rd, rj, rk, sel.
func (Builder) BytepickD(rd, rj, rk Reg, sel UImm3) Instr {
	return BytepickD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		rk:  rk.Num(),
		sel: immNum(sel.Val()),
	}
}

// BytepickW - bytepick.w rd, rj, rk, sel.
func (Builder) BytepickW(rd, rj, rk Reg, sel UImm2) Instr {
	return BytepickW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		rk:  rk.Num(),
		sel: immNum(sel.Val()),
	}
}

// Cacop - cacop op, rj, si12 (the assembly operand order).
func (Builder) Cacop(op UImm5, rj Reg, off Imm12) Instr {
	return Cacop{
		op:  immNum(op.Val()),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// CloD - clo.d rd, rj.
func (Builder) CloD(rd, rj Reg) Instr {
	return CloD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// CloW - clo.w rd, rj.
func (Builder) CloW(rd, rj Reg) Instr {
	return CloW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// ClzD - clz.d rd, rj.
func (Builder) ClzD(rd, rj Reg) Instr {
	return ClzD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// ClzW - clz.w rd, rj.
func (Builder) ClzW(rd, rj Reg) Instr {
	return ClzW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Code15 - a validated code; an error when out of range.
func (Builder) Code15(v int64) (Code15, error) {
	x, err := newImm(v, 0, 32767, "code")
	if err != nil {
		return Code15{}, err
	}

	return Code15{v: x}, nil
}

// Cpucfg - cpucfg rd, rj.
func (Builder) Cpucfg(rd, rj Reg) Instr {
	return Cpucfg{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// CrcWBW - crc.w.b.w rd, rj, rk.
func (Builder) CrcWBW(rd, rj, rk Reg) Instr {
	return CrcWBW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrcWDW - crc.w.d.w rd, rj, rk.
func (Builder) CrcWDW(rd, rj, rk Reg) Instr {
	return CrcWDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrcWHW - crc.w.h.w rd, rj, rk.
func (Builder) CrcWHW(rd, rj, rk Reg) Instr {
	return CrcWHW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrcWWW - crc.w.w.w rd, rj, rk.
func (Builder) CrcWWW(rd, rj, rk Reg) Instr {
	return CrcWWW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrccWBW - crcc.w.b.w rd, rj, rk.
func (Builder) CrccWBW(rd, rj, rk Reg) Instr {
	return CrccWBW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrccWDW - crcc.w.d.w rd, rj, rk.
func (Builder) CrccWDW(rd, rj, rk Reg) Instr {
	return CrccWDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrccWHW - crcc.w.h.w rd, rj, rk.
func (Builder) CrccWHW(rd, rj, rk Reg) Instr {
	return CrccWHW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// CrccWWW - crcc.w.w.w rd, rj, rk.
func (Builder) CrccWWW(rd, rj, rk Reg) Instr {
	return CrccWWW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Csrrd - csrrd rd, csr.
func (Builder) Csrrd(rd Reg, csr UImm14) Instr {
	return Csrrd{
		rd:  rd.Num(),
		csr: immNum(csr.Val()),
	}
}

// Csrwr - csrwr rd, csr.
func (Builder) Csrwr(rd Reg, csr UImm14) Instr {
	return Csrwr{
		rd:  rd.Num(),
		csr: immNum(csr.Val()),
	}
}

// Csrxchg - csrxchg rd, rj, csr.
func (Builder) Csrxchg(rd, rj Reg, csr UImm14) Instr {
	return Csrxchg{
		rd:  rd.Num(),
		rj:  rj.Num(),
		csr: immNum(csr.Val()),
	}
}

// CtoD - cto.d rd, rj.
func (Builder) CtoD(rd, rj Reg) Instr {
	return CtoD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// CtoW - cto.w rd, rj.
func (Builder) CtoW(rd, rj Reg) Instr {
	return CtoW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// CtzD - ctz.d rd, rj.
func (Builder) CtzD(rd, rj Reg) Instr {
	return CtzD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// CtzW - ctz.w rd, rj.
func (Builder) CtzW(rd, rj Reg) Instr {
	return CtzW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Dbar - dbar hint (a 15-bit code).
func (Builder) Dbar(code Code15) Instr {
	return Dbar{
		code: immNum(code.Val()),
	}
}

// Dbcl - dbcl code.
func (Builder) Dbcl(code Code15) Instr {
	return Dbcl{
		code: immNum(code.Val()),
	}
}

// DivD - div.d rd, rj, rk.
func (Builder) DivD(rd, rj, rk Reg) Instr {
	return DivD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// DivDu - div.du rd, rj, rk.
func (Builder) DivDu(rd, rj, rk Reg) Instr {
	return DivDu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// DivW - div.w rd, rj, rk.
func (Builder) DivW(rd, rj, rk Reg) Instr {
	return DivW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// DivWu - div.wu rd, rj, rk.
func (Builder) DivWu(rd, rj, rk Reg) Instr {
	return DivWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Ertn - ertn (no operands).
func (Builder) Ertn() Instr {
	return Ertn{}
}

// ExtWB - ext.w.b rd, rj.
func (Builder) ExtWB(rd, rj Reg) Instr {
	return ExtWB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// ExtWH - ext.w.h rd, rj.
func (Builder) ExtWH(rd, rj Reg) Instr {
	return ExtWH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Ibar - ibar hint (a 15-bit code).
func (Builder) Ibar(code Code15) Instr {
	return Ibar{
		code: immNum(code.Val()),
	}
}

// Idle - idle code.
func (Builder) Idle(code Code15) Instr {
	return Idle{
		code: immNum(code.Val()),
	}
}

// Imm12 - a validated value; an error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	x, err := newImm(v, -2048, 2047, "si12 value")
	if err != nil {
		return Imm12{}, err
	}

	return Imm12{v: x}, nil
}

// Imm14 - a validated word-aligned value; an error otherwise.
func (Builder) Imm14(v int64) (Imm14, error) {
	if v%4 != 0 {
		return Imm14{}, fmt.Errorf("loong64: si14 byte offset %d not word-aligned", v)
	}

	x, err := newImm(v, -16380, 16380, "si14 byte offset")
	if err != nil {
		return Imm14{}, err
	}

	return Imm14{v: x}, nil
}

// Imm16 - a validated value; an error when out of range.
func (Builder) Imm16(v int64) (Imm16, error) {
	x, err := newImm(v, -32768, 32767, "si16 value")
	if err != nil {
		return Imm16{}, err
	}

	return Imm16{v: x}, nil
}

// Imm20 - a validated value; an error when out of range.
func (Builder) Imm20(v int64) (Imm20, error) {
	x, err := newImm(v, -524288, 524287, "si20 value")
	if err != nil {
		return Imm20{}, err
	}

	return Imm20{v: x}, nil
}

// Invtlb - invtlb op, rj, rk (the assembly operand order).
func (Builder) Invtlb(op UImm5, rj, rk Reg) Instr {
	return Invtlb{
		rj: rj.Num(),
		rk: rk.Num(),
		op: immNum(op.Val()),
	}
}

// IocsrrdB - iocsrrd.b rd, rj.
func (Builder) IocsrrdB(rd, rj Reg) Instr {
	return IocsrrdB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrrdD - iocsrrd.d rd, rj.
func (Builder) IocsrrdD(rd, rj Reg) Instr {
	return IocsrrdD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrrdH - iocsrrd.h rd, rj.
func (Builder) IocsrrdH(rd, rj Reg) Instr {
	return IocsrrdH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrrdW - iocsrrd.w rd, rj.
func (Builder) IocsrrdW(rd, rj Reg) Instr {
	return IocsrrdW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrwrB - iocsrwr.b rd, rj.
func (Builder) IocsrwrB(rd, rj Reg) Instr {
	return IocsrwrB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrwrD - iocsrwr.d rd, rj.
func (Builder) IocsrwrD(rd, rj Reg) Instr {
	return IocsrwrD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrwrH - iocsrwr.h rd, rj.
func (Builder) IocsrwrH(rd, rj Reg) Instr {
	return IocsrwrH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// IocsrwrW - iocsrwr.w rd, rj.
func (Builder) IocsrwrW(rd, rj Reg) Instr {
	return IocsrwrW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Jirl - jirl rd, rj, offs (the byte offset from rj).
func (Builder) Jirl(rd, rj Reg, off int64) Instr {
	return Jirl{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

// LdB - ld.b rd, rj, si12.
func (Builder) LdB(rd, rj Reg, v Imm12) Instr {
	return LdB{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdBu - ld.bu rd, rj, si12.
func (Builder) LdBu(rd, rj Reg, v Imm12) Instr {
	return LdBu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(v.Val()),
	}
}

// LdD - ld.d rd, rj, si12.
func (Builder) LdD(rd, rj Reg, v Imm12) Instr {
	return LdD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdH - ld.h rd, rj, si12.
func (Builder) LdH(rd, rj Reg, v Imm12) Instr {
	return LdH{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdHu - ld.hu rd, rj, si12.
func (Builder) LdHu(rd, rj Reg, v Imm12) Instr {
	return LdHu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(v.Val()),
	}
}

// LdW - ld.w rd, rj, si12.
func (Builder) LdW(rd, rj Reg, v Imm12) Instr {
	return LdW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdWu - ld.wu rd, rj, si12.
func (Builder) LdWu(rd, rj Reg, v Imm12) Instr {
	return LdWu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Lddir - lddir rd, rj, ui8.
func (Builder) Lddir(rd, rj Reg, v UImm8) Instr {
	return Lddir{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdgtB - ldgt.b rd, rj, rk.
func (Builder) LdgtB(rd, rj, rk Reg) Instr {
	return LdgtB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdgtD - ldgt.d rd, rj, rk.
func (Builder) LdgtD(rd, rj, rk Reg) Instr {
	return LdgtD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdgtH - ldgt.h rd, rj, rk.
func (Builder) LdgtH(rd, rj, rk Reg) Instr {
	return LdgtH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdgtW - ldgt.w rd, rj, rk.
func (Builder) LdgtW(rd, rj, rk Reg) Instr {
	return LdgtW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdleB - ldle.b rd, rj, rk.
func (Builder) LdleB(rd, rj, rk Reg) Instr {
	return LdleB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdleD - ldle.d rd, rj, rk.
func (Builder) LdleD(rd, rj, rk Reg) Instr {
	return LdleD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdleH - ldle.h rd, rj, rk.
func (Builder) LdleH(rd, rj, rk Reg) Instr {
	return LdleH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdleW - ldle.w rd, rj, rk.
func (Builder) LdleW(rd, rj, rk Reg) Instr {
	return LdleW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Ldpte - ldpte rj, ui8.
func (Builder) Ldpte(rj Reg, v UImm8) Instr {
	return Ldpte{
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// LdptrD - ldptr.d rd, rj, offs (the byte offset).
func (Builder) LdptrD(rd, rj Reg, off Imm14) Instr {
	return LdptrD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// LdptrW - ldptr.w rd, rj, offs (the byte offset).
func (Builder) LdptrW(rd, rj Reg, off Imm14) Instr {
	return LdptrW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// LdxB - ldx.b rd, rj, rk.
func (Builder) LdxB(rd, rj, rk Reg) Instr {
	return LdxB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxBu - ldx.bu rd, rj, rk.
func (Builder) LdxBu(rd, rj, rk Reg) Instr {
	return LdxBu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxD - ldx.d rd, rj, rk.
func (Builder) LdxD(rd, rj, rk Reg) Instr {
	return LdxD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxH - ldx.h rd, rj, rk.
func (Builder) LdxH(rd, rj, rk Reg) Instr {
	return LdxH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxHu - ldx.hu rd, rj, rk.
func (Builder) LdxHu(rd, rj, rk Reg) Instr {
	return LdxHu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxW - ldx.w rd, rj, rk.
func (Builder) LdxW(rd, rj, rk Reg) Instr {
	return LdxW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LdxWu - ldx.wu rd, rj, rk.
func (Builder) LdxWu(rd, rj, rk Reg) Instr {
	return LdxWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// LlD - ll.d rd, rj, offs (the byte offset from rj).
func (Builder) LlD(rd, rj Reg, off Imm14) Instr {
	return LlD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// LlW - ll.w rd, rj, offs (the byte offset from rj).
func (Builder) LlW(rd, rj Reg, off Imm14) Instr {
	return LlW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// LlacqD - llacq.d rd, rj.
func (Builder) LlacqD(rd, rj Reg) Instr {
	return LlacqD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// LlacqW - llacq.w rd, rj.
func (Builder) LlacqW(rd, rj Reg) Instr {
	return LlacqW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Lu12iW - lu12i.w rd, si20.
func (Builder) Lu12iW(rd Reg, v Imm20) Instr {
	return Lu12iW{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Lu32iD - lu32i.d rd, si20.
func (Builder) Lu32iD(rd Reg, v Imm20) Instr {
	return Lu32iD{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Lu52iD - lu52i.d rd, rj, si12.
func (Builder) Lu52iD(rd, rj Reg, v Imm12) Instr {
	return Lu52iD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Maskeqz - maskeqz rd, rj, rk.
func (Builder) Maskeqz(rd, rj, rk Reg) Instr {
	return Maskeqz{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Masknez - masknez rd, rj, rk.
func (Builder) Masknez(rd, rj, rk Reg) Instr {
	return Masknez{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// ModD - mod.d rd, rj, rk.
func (Builder) ModD(rd, rj, rk Reg) Instr {
	return ModD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// ModDu - mod.du rd, rj, rk.
func (Builder) ModDu(rd, rj, rk Reg) Instr {
	return ModDu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// ModW - mod.w rd, rj, rk.
func (Builder) ModW(rd, rj, rk Reg) Instr {
	return ModW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// ModWu - mod.wu rd, rj, rk.
func (Builder) ModWu(rd, rj, rk Reg) Instr {
	return ModWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulD - mul.d rd, rj, rk.
func (Builder) MulD(rd, rj, rk Reg) Instr {
	return MulD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulW - mul.w rd, rj, rk.
func (Builder) MulW(rd, rj, rk Reg) Instr {
	return MulW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulhD - mulh.d rd, rj, rk.
func (Builder) MulhD(rd, rj, rk Reg) Instr {
	return MulhD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulhDu - mulh.du rd, rj, rk.
func (Builder) MulhDu(rd, rj, rk Reg) Instr {
	return MulhDu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulhW - mulh.w rd, rj, rk.
func (Builder) MulhW(rd, rj, rk Reg) Instr {
	return MulhW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulhWu - mulh.wu rd, rj, rk.
func (Builder) MulhWu(rd, rj, rk Reg) Instr {
	return MulhWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulwDW - mulw.d.w rd, rj, rk.
func (Builder) MulwDW(rd, rj, rk Reg) Instr {
	return MulwDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// MulwDWu - mulw.d.wu rd, rj, rk.
func (Builder) MulwDWu(rd, rj, rk Reg) Instr {
	return MulwDWu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Nor - nor rd, rj, rk.
func (Builder) Nor(rd, rj, rk Reg) Instr {
	return Nor{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Off16 - a validated word-aligned value; an error otherwise.
func (Builder) Off16(v int64) (Off16, error) {
	if v%4 != 0 {
		return Off16{}, fmt.Errorf("loong64: si16 byte offset %d not word-aligned", v)
	}

	x, err := newImm(v, -131068, 131068, "si16 byte offset")
	if err != nil {
		return Off16{}, err
	}

	return Off16{v: x}, nil
}

// Or - or rd, rj, rk.
func (Builder) Or(rd, rj, rk Reg) Instr {
	return Or{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Ori - ori rd, rj, ui12.
func (Builder) Ori(rd, rj Reg, v UImm12) Instr {
	return Ori{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Orn - orn rd, rj, rk.
func (Builder) Orn(rd, rj, rk Reg) Instr {
	return Orn{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Pcaddi - pcaddi rd, si20.
func (Builder) Pcaddi(rd Reg, v Imm20) Instr {
	return Pcaddi{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Pcaddu12i - pcaddu12i rd, si20.
func (Builder) Pcaddu12i(rd Reg, v Imm20) Instr {
	return Pcaddu12i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Pcaddu18i - pcaddu18i rd, si20.
func (Builder) Pcaddu18i(rd Reg, v Imm20) Instr {
	return Pcaddu18i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Pcalau12i - pcalau12i rd, si20.
func (Builder) Pcalau12i(rd Reg, v Imm20) Instr {
	return Pcalau12i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

// Preld - preld hint, rj, si12.
func (Builder) Preld(hint UImm5, rj Reg, off Imm12) Instr {
	return Preld{
		rj:   rj.Num(),
		hint: immNum(hint.Val()),
		off:  immNum(off.Val()),
	}
}

// Preldx - preldx hint, rj, rk.
func (Builder) Preldx(hint UImm5, rj, rk Reg) Instr {
	return Preldx{
		rj:   rj.Num(),
		rk:   rk.Num(),
		hint: immNum(hint.Val()),
	}
}

// RdtimeD - rdtime.d rd, rj.
func (Builder) RdtimeD(rd, rj Reg) Instr {
	return RdtimeD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RdtimehW - rdtimeh.w rd, rj.
func (Builder) RdtimehW(rd, rj Reg) Instr {
	return RdtimehW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RdtimelW - rdtimel.w rd, rj.
func (Builder) RdtimelW(rd, rj Reg) Instr {
	return RdtimelW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revb2H - revb.2h rd, rj.
func (Builder) Revb2H(rd, rj Reg) Instr {
	return Revb2H{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revb2W - revb.2w rd, rj.
func (Builder) Revb2W(rd, rj Reg) Instr {
	return Revb2W{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revb4H - revb.4h rd, rj.
func (Builder) Revb4H(rd, rj Reg) Instr {
	return Revb4H{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RevbD - revb.d rd, rj.
func (Builder) RevbD(rd, rj Reg) Instr {
	return RevbD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revbit4B - bitrev.4b rd, rj.
func (Builder) Revbit4B(rd, rj Reg) Instr {
	return Revbit4B{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revbit8B - bitrev.8b rd, rj.
func (Builder) Revbit8B(rd, rj Reg) Instr {
	return Revbit8B{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RevbitD - bitrev.d rd, rj.
func (Builder) RevbitD(rd, rj Reg) Instr {
	return RevbitD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RevbitW - bitrev.w rd, rj.
func (Builder) RevbitW(rd, rj Reg) Instr {
	return RevbitW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Revh2W - revh.2w rd, rj.
func (Builder) Revh2W(rd, rj Reg) Instr {
	return Revh2W{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RevhD - revh.d rd, rj.
func (Builder) RevhD(rd, rj Reg) Instr {
	return RevhD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// RotrD - rotr.d rd, rj, rk.
func (Builder) RotrD(rd, rj, rk Reg) Instr {
	return RotrD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// RotrW - rotr.w rd, rj, rk.
func (Builder) RotrW(rd, rj, rk Reg) Instr {
	return RotrW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// RotriD - rotri.d rd, rj, ui6.
func (Builder) RotriD(rd, rj Reg, v UImm6) Instr {
	return RotriD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// RotriW - rotri.w rd, rj, ui5.
func (Builder) RotriW(rd, rj Reg, v UImm5) Instr {
	return RotriW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// ScD - sc.d rd, rj, offs (the byte offset from rj).
func (Builder) ScD(rd, rj Reg, off Imm14) Instr {
	return ScD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// ScQ - sc.q rd, rk, rj.
func (Builder) ScQ(rd, rk, rj Reg) Instr {
	return ScQ{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

// ScW - sc.w rd, rj, offs (the byte offset from rj).
func (Builder) ScW(rd, rj Reg, off Imm14) Instr {
	return ScW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// ScrelD - screl.d rd, rj.
func (Builder) ScrelD(rd, rj Reg) Instr {
	return ScrelD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// ScrelW - screl.w rd, rj.
func (Builder) ScrelW(rd, rj Reg) Instr {
	return ScrelW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

// Shift3 - a validated alsl shift amount (1..4).
func (Builder) Shift3(v int64) (Shift3, error) {
	x, err := newImm(v, 1, 4, "alsl shift amount")
	if err != nil {
		return Shift3{}, err
	}

	return Shift3{v: x}, nil
}

// SllD - sll.d rd, rj, rk.
func (Builder) SllD(rd, rj, rk Reg) Instr {
	return SllD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SllW - sll.w rd, rj, rk.
func (Builder) SllW(rd, rj, rk Reg) Instr {
	return SllW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SlliD - slli.d rd, rj, ui6.
func (Builder) SlliD(rd, rj Reg, v UImm6) Instr {
	return SlliD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// SlliW - slli.w rd, rj, ui5.
func (Builder) SlliW(rd, rj Reg, v UImm5) Instr {
	return SlliW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Slt - slt rd, rj, rk.
func (Builder) Slt(rd, rj, rk Reg) Instr {
	return Slt{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Slti - slti rd, rj, si12.
func (Builder) Slti(rd, rj Reg, v Imm12) Instr {
	return Slti{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// Sltu - sltu rd, rj, rk.
func (Builder) Sltu(rd, rj, rk Reg) Instr {
	return Sltu{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Sltui - sltui rd, rj, si12.
func (Builder) Sltui(rd, rj Reg, v Imm12) Instr {
	return Sltui{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// SraD - sra.d rd, rj, rk.
func (Builder) SraD(rd, rj, rk Reg) Instr {
	return SraD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SraW - sra.w rd, rj, rk.
func (Builder) SraW(rd, rj, rk Reg) Instr {
	return SraW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SraiD - srai.d rd, rj, ui6.
func (Builder) SraiD(rd, rj Reg, v UImm6) Instr {
	return SraiD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// SraiW - srai.w rd, rj, ui5.
func (Builder) SraiW(rd, rj Reg, v UImm5) Instr {
	return SraiW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// SrlD - srl.d rd, rj, rk.
func (Builder) SrlD(rd, rj, rk Reg) Instr {
	return SrlD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SrlW - srl.w rd, rj, rk.
func (Builder) SrlW(rd, rj, rk Reg) Instr {
	return SrlW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SrliD - srli.d rd, rj, ui6.
func (Builder) SrliD(rd, rj Reg, v UImm6) Instr {
	return SrliD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// SrliW - srli.w rd, rj, ui5.
func (Builder) SrliW(rd, rj Reg, v UImm5) Instr {
	return SrliW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// StB - st.b rd, rj, si12.
func (Builder) StB(rd, rj Reg, v Imm12) Instr {
	return StB{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// StD - st.d rd, rj, si12.
func (Builder) StD(rd, rj Reg, v Imm12) Instr {
	return StD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// StH - st.h rd, rj, si12.
func (Builder) StH(rd, rj Reg, v Imm12) Instr {
	return StH{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// StW - st.w rd, rj, si12.
func (Builder) StW(rd, rj Reg, v Imm12) Instr {
	return StW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}

// StgtB - stgt.b rd, rj, rk.
func (Builder) StgtB(rd, rj, rk Reg) Instr {
	return StgtB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StgtD - stgt.d rd, rj, rk.
func (Builder) StgtD(rd, rj, rk Reg) Instr {
	return StgtD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StgtH - stgt.h rd, rj, rk.
func (Builder) StgtH(rd, rj, rk Reg) Instr {
	return StgtH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StgtW - stgt.w rd, rj, rk.
func (Builder) StgtW(rd, rj, rk Reg) Instr {
	return StgtW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StleB - stle.b rd, rj, rk.
func (Builder) StleB(rd, rj, rk Reg) Instr {
	return StleB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StleD - stle.d rd, rj, rk.
func (Builder) StleD(rd, rj, rk Reg) Instr {
	return StleD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StleH - stle.h rd, rj, rk.
func (Builder) StleH(rd, rj, rk Reg) Instr {
	return StleH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StleW - stle.w rd, rj, rk.
func (Builder) StleW(rd, rj, rk Reg) Instr {
	return StleW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StptrD - stptr.d rd, rj, offs (the byte offset).
func (Builder) StptrD(rd, rj Reg, off Imm14) Instr {
	return StptrD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// StptrW - stptr.w rd, rj, offs (the byte offset).
func (Builder) StptrW(rd, rj Reg, off Imm14) Instr {
	return StptrW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off.Val()),
	}
}

// StxB - stx.b rd, rj, rk.
func (Builder) StxB(rd, rj, rk Reg) Instr {
	return StxB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StxD - stx.d rd, rj, rk.
func (Builder) StxD(rd, rj, rk Reg) Instr {
	return StxD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StxH - stx.h rd, rj, rk.
func (Builder) StxH(rd, rj, rk Reg) Instr {
	return StxH{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// StxW - stx.w rd, rj, rk.
func (Builder) StxW(rd, rj, rk Reg) Instr {
	return StxW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SubD - sub.d rd, rj, rk.
func (Builder) SubD(rd, rj, rk Reg) Instr {
	return SubD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// SubW - sub.w rd, rj, rk.
func (Builder) SubW(rd, rj, rk Reg) Instr {
	return SubW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Syscall - syscall code (a 15-bit code).
func (Builder) Syscall(code Code15) Instr {
	return Syscall{
		code: immNum(code.Val()),
	}
}

// Tlbclr - tlbclr (no operands).
func (Builder) Tlbclr() Instr {
	return Tlbclr{}
}

// Tlbfill - tlbfill (no operands).
func (Builder) Tlbfill() Instr {
	return Tlbfill{}
}

// Tlbflush - tlbflush (no operands).
func (Builder) Tlbflush() Instr {
	return Tlbflush{}
}

// Tlbrd - tlbrd (no operands).
func (Builder) Tlbrd() Instr {
	return Tlbrd{}
}

// Tlbsrch - tlbsrch (no operands).
func (Builder) Tlbsrch() Instr {
	return Tlbsrch{}
}

// Tlbwr - tlbwr (no operands).
func (Builder) Tlbwr() Instr {
	return Tlbwr{}
}

// UImm12 - a validated value; an error when out of range.
func (Builder) UImm12(v int64) (UImm12, error) {
	x, err := newImm(v, 0, 4095, "ui12 value")
	if err != nil {
		return UImm12{}, err
	}

	return UImm12{v: x}, nil
}

// UImm14 - a validated CSR number; an error when out of range.
func (Builder) UImm14(v int64) (UImm14, error) {
	x, err := newImm(v, 0, 16383, "csr number")
	if err != nil {
		return UImm14{}, err
	}

	return UImm14{v: x}, nil
}

// UImm2 - a validated value; an error when out of range.
func (Builder) UImm2(v int64) (UImm2, error) {
	x, err := newImm(v, 0, 3, "ui2 value")
	if err != nil {
		return UImm2{}, err
	}

	return UImm2{v: x}, nil
}

// UImm3 - a validated value; an error when out of range.
func (Builder) UImm3(v int64) (UImm3, error) {
	x, err := newImm(v, 0, 7, "ui3 value")
	if err != nil {
		return UImm3{}, err
	}

	return UImm3{v: x}, nil
}

// UImm5 - a validated value; an error when out of range.
func (Builder) UImm5(v int64) (UImm5, error) {
	x, err := newImm(v, 0, 31, "ui5 value")
	if err != nil {
		return UImm5{}, err
	}

	return UImm5{v: x}, nil
}

// UImm6 - a validated value; an error when out of range.
func (Builder) UImm6(v int64) (UImm6, error) {
	x, err := newImm(v, 0, 63, "ui6 value")
	if err != nil {
		return UImm6{}, err
	}

	return UImm6{v: x}, nil
}

// UImm8 - a validated value; an error when out of range.
func (Builder) UImm8(v int64) (UImm8, error) {
	x, err := newImm(v, 0, 255, "ui8 value")
	if err != nil {
		return UImm8{}, err
	}

	return UImm8{v: x}, nil
}

// Xor - xor rd, rj, rk.
func (Builder) Xor(rd, rj, rk Reg) Instr {
	return Xor{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

// Xori - xori rd, rj, ui12.
func (Builder) Xori(rd, rj Reg, v UImm12) Instr {
	return Xori{
		rd:  rd.Num(),
		rj:  rj.Num(),
		imm: immNum(v.Val()),
	}
}
