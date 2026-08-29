package loong64

// The assemble side: the operand shapes of every scalar integer
// instruction. The shapes validate the operand kinds and counts and
// delegate the range validation to the role-type constructors; a
// declarative table (the set is fully regular per family) instead of
// per-instruction new* functions. asm/loong64 (the syntax layer above)
// resolves expressions into computed operands and calls BuildInstr.

import (
	"errors"
	"fmt"
)

// asmCtors - mnemonic -> (computed operands -> instruction).
var asmCtors = map[string]func(ops []Op) (Instr, error){
	"add.w":       regs3(NewAddW),
	"add.d":       regs3(NewAddD),
	"sub.w":       regs3(NewSubW),
	"sub.d":       regs3(NewSubD),
	"slt":         regs3(NewSlt),
	"sltu":        regs3(NewSltu),
	"maskeqz":     regs3(NewMaskeqz),
	"masknez":     regs3(NewMasknez),
	"nor":         regs3(NewNor),
	"and":         regs3(NewAnd),
	"or":          regs3(NewOr),
	"xor":         regs3(NewXor),
	"orn":         regs3(NewOrn),
	"andn":        regs3(NewAndn),
	"sll.w":       regs3(NewSllW),
	"srl.w":       regs3(NewSrlW),
	"sra.w":       regs3(NewSraW),
	"sll.d":       regs3(NewSllD),
	"srl.d":       regs3(NewSrlD),
	"sra.d":       regs3(NewSraD),
	"rotr.w":      regs3(NewRotrW),
	"rotr.d":      regs3(NewRotrD),
	"addi.w":      regRegI12(NewAddiW),
	"addi.d":      regRegI12(NewAddiD),
	"slti":        regRegI12(NewSlti),
	"sltui":       regRegI12(NewSltui),
	"andi":        regRegU12(NewAndi),
	"ori":         regRegU12(NewOri),
	"xori":        regRegU12(NewXori),
	"addu16i.d":   regRegI16(NewAddu16iD),
	"slli.w":      regRegU5(NewSlliW),
	"srli.w":      regRegU5(NewSrliW),
	"srai.w":      regRegU5(NewSraiW),
	"rotri.w":     regRegU5(NewRotriW),
	"slli.d":      regRegU6(NewSlliD),
	"srli.d":      regRegU6(NewSrliD),
	"srai.d":      regRegU6(NewSraiD),
	"rotri.d":     regRegU6(NewRotriD),
	"lu12i.w":     regI20(NewLu12iW),
	"lu32i.d":     regI20(NewLu32iD),
	"lu52i.d":     regRegI12(NewLu52iD),
	"pcaddi":      regI20(NewPcaddi),
	"pcalau12i":   regI20(NewPcalau12i),
	"pcaddu12i":   regI20(NewPcaddu12i),
	"pcaddu18i":   regI20(NewPcaddu18i),
	"ld.b":        regRegI12(NewLdB),
	"ld.h":        regRegI12(NewLdH),
	"ld.w":        regRegI12(NewLdW),
	"ld.wu":       regRegI12(NewLdWu),
	"ld.bu":       regRegI12(NewLdBu),
	"ld.hu":       regRegI12(NewLdHu),
	"ld.d":        regRegI12(NewLdD),
	"st.b":        regRegI12(NewStB),
	"st.h":        regRegI12(NewStH),
	"st.w":        regRegI12(NewStW),
	"st.d":        regRegI12(NewStD),
	"ldptr.w":     regRegI14(NewLdptrW),
	"ldptr.d":     regRegI14(NewLdptrD),
	"stptr.w":     regRegI14(NewStptrW),
	"stptr.d":     regRegI14(NewStptrD),
	"ldx.b":       regs3(NewLdxB),
	"ldx.h":       regs3(NewLdxH),
	"ldx.w":       regs3(NewLdxW),
	"ldx.wu":      regs3(NewLdxWu),
	"ldx.bu":      regs3(NewLdxBu),
	"ldx.hu":      regs3(NewLdxHu),
	"ldx.d":       regs3(NewLdxD),
	"stx.b":       regs3(NewStxB),
	"stx.h":       regs3(NewStxH),
	"stx.w":       regs3(NewStxW),
	"stx.d":       regs3(NewStxD),
	"preld":       u5RegI12(NewPreld),
	"preldx":      u5Regs2(NewPreldx),
	"beq":         regs2Target(NewBeq),
	"bne":         regs2Target(NewBne),
	"blt":         regs2Target(NewBlt),
	"bge":         regs2Target(NewBge),
	"bltu":        regs2Target(NewBltu),
	"bgeu":        regs2Target(NewBgeu),
	"beqz":        regTarget(NewBeqz),
	"bnez":        regTarget(NewBnez),
	"b":           target(NewB),
	"bl":          target(NewBl),
	"jirl":        regs2Off(NewJirl),
	"break":       code15(NewBreak),
	"syscall":     code15(NewSyscall),
	"dbcl":        code15(NewDbcl),
	"dbar":        code15(NewDbar),
	"ibar":        code15(NewIbar),
	"mul.w":       regs3(NewMulW),
	"mulh.w":      regs3(NewMulhW),
	"mulh.wu":     regs3(NewMulhWu),
	"mul.d":       regs3(NewMulD),
	"mulh.d":      regs3(NewMulhD),
	"mulh.du":     regs3(NewMulhDu),
	"mulw.d.w":    regs3(NewMulwDW),
	"mulw.d.wu":   regs3(NewMulwDWu),
	"div.w":       regs3(NewDivW),
	"mod.w":       regs3(NewModW),
	"div.wu":      regs3(NewDivWu),
	"mod.wu":      regs3(NewModWu),
	"div.d":       regs3(NewDivD),
	"mod.d":       regs3(NewModD),
	"div.du":      regs3(NewDivDu),
	"mod.du":      regs3(NewModDu),
	"ext.w.b":     regs2(NewExtWB),
	"ext.w.h":     regs2(NewExtWH),
	"rdtimel.w":   regs2(NewRdtimelW),
	"rdtimeh.w":   regs2(NewRdtimehW),
	"rdtime.d":    regs2(NewRdtimeD),
	"cpucfg":      regs2(NewCpucfg),
	"clo.w":       regs2(NewCloW),
	"clz.w":       regs2(NewClzW),
	"cto.w":       regs2(NewCtoW),
	"ctz.w":       regs2(NewCtzW),
	"clo.d":       regs2(NewCloD),
	"clz.d":       regs2(NewClzD),
	"cto.d":       regs2(NewCtoD),
	"ctz.d":       regs2(NewCtzD),
	"revb.2h":     regs2(NewRevb2H),
	"revb.4h":     regs2(NewRevb4H),
	"revb.2w":     regs2(NewRevb2W),
	"revb.d":      regs2(NewRevbD),
	"revh.2w":     regs2(NewRevh2W),
	"revh.d":      regs2(NewRevhD),
	"bitrev.4b":   regs2(NewRevbit4B),
	"bitrev.w":    regs2(NewRevbitW),
	"bitrev.8b":   regs2(NewRevbit8B),
	"bitrev.d":    regs2(NewRevbitD),
	"bstrins.w":   regRegU5U5(NewBstrinsW),
	"bstrpick.w":  regRegU5U5(NewBstrpickW),
	"bstrins.d":   regRegU6U6(NewBstrinsD),
	"bstrpick.d":  regRegU6U6(NewBstrpickD),
	"alsl.w":      regs3Shift(NewAlslW),
	"alsl.wu":     regs3Shift(NewAlslWu),
	"alsl.d":      regs3Shift(NewAlslD),
	"bytepick.w":  regs3U2(NewBytepickW),
	"bytepick.d":  regs3U3(NewBytepickD),
	"crc.w.b.w":   regs3(NewCrcWBW),
	"crc.w.h.w":   regs3(NewCrcWHW),
	"crc.w.w.w":   regs3(NewCrcWWW),
	"crc.w.d.w":   regs3(NewCrcWDW),
	"crcc.w.b.w":  regs3(NewCrccWBW),
	"crcc.w.h.w":  regs3(NewCrccWHW),
	"crcc.w.w.w":  regs3(NewCrccWWW),
	"crcc.w.d.w":  regs3(NewCrccWDW),
	"asrtle.d":    regs2(NewAsrtleD),
	"asrtgt.d":    regs2(NewAsrtgtD),
	"ldgt.b":      regs3(NewLdgtB),
	"ldgt.h":      regs3(NewLdgtH),
	"ldgt.w":      regs3(NewLdgtW),
	"ldgt.d":      regs3(NewLdgtD),
	"ldle.b":      regs3(NewLdleB),
	"ldle.h":      regs3(NewLdleH),
	"ldle.w":      regs3(NewLdleW),
	"ldle.d":      regs3(NewLdleD),
	"stgt.b":      regs3(NewStgtB),
	"stgt.h":      regs3(NewStgtH),
	"stgt.w":      regs3(NewStgtW),
	"stgt.d":      regs3(NewStgtD),
	"stle.b":      regs3(NewStleB),
	"stle.h":      regs3(NewStleH),
	"stle.w":      regs3(NewStleW),
	"stle.d":      regs3(NewStleD),
	"ll.w":        regRegI14(NewLlW),
	"ll.d":        regRegI14(NewLlD),
	"sc.w":        regRegI14(NewScW),
	"sc.d":        regRegI14(NewScD),
	"llacq.w":     regs2(NewLlacqW),
	"llacq.d":     regs2(NewLlacqD),
	"screl.w":     regs2(NewScrelW),
	"screl.d":     regs2(NewScrelD),
	"sc.q":        regs3rev(NewScQ),
	"amcas.b":     regs3rev(NewAmcasB),
	"amcas.h":     regs3rev(NewAmcasH),
	"amcas.w":     regs3rev(NewAmcasW),
	"amcas.d":     regs3rev(NewAmcasD),
	"amcas_db.b":  regs3rev(NewAmcasDbB),
	"amcas_db.h":  regs3rev(NewAmcasDbH),
	"amcas_db.w":  regs3rev(NewAmcasDbW),
	"amcas_db.d":  regs3rev(NewAmcasDbD),
	"amswap.b":    regs3rev(NewAmswapB),
	"amswap.h":    regs3rev(NewAmswapH),
	"amswap.w":    regs3rev(NewAmswapW),
	"amswap.d":    regs3rev(NewAmswapD),
	"amswap_db.b": regs3rev(NewAmswapDbB),
	"amswap_db.h": regs3rev(NewAmswapDbH),
	"amswap_db.w": regs3rev(NewAmswapDbW),
	"amswap_db.d": regs3rev(NewAmswapDbD),
	"amadd.b":     regs3rev(NewAmaddB),
	"amadd.h":     regs3rev(NewAmaddH),
	"amadd.w":     regs3rev(NewAmaddW),
	"amadd.d":     regs3rev(NewAmaddD),
	"amadd_db.b":  regs3rev(NewAmaddDbB),
	"amadd_db.h":  regs3rev(NewAmaddDbH),
	"amadd_db.w":  regs3rev(NewAmaddDbW),
	"amadd_db.d":  regs3rev(NewAmaddDbD),
	"amand.w":     regs3rev(NewAmandW),
	"amand.d":     regs3rev(NewAmandD),
	"amand_db.w":  regs3rev(NewAmandDbW),
	"amand_db.d":  regs3rev(NewAmandDbD),
	"amor.w":      regs3rev(NewAmorW),
	"amor.d":      regs3rev(NewAmorD),
	"amor_db.w":   regs3rev(NewAmorDbW),
	"amor_db.d":   regs3rev(NewAmorDbD),
	"amxor.w":     regs3rev(NewAmxorW),
	"amxor.d":     regs3rev(NewAmxorD),
	"amxor_db.w":  regs3rev(NewAmxorDbW),
	"amxor_db.d":  regs3rev(NewAmxorDbD),
	"ammax.w":     regs3rev(NewAmmaxW),
	"ammax.d":     regs3rev(NewAmmaxD),
	"ammax.wu":    regs3rev(NewAmmaxWu),
	"ammax.du":    regs3rev(NewAmmaxDu),
	"ammax_db.w":  regs3rev(NewAmmaxDbW),
	"ammax_db.d":  regs3rev(NewAmmaxDbD),
	"ammax_db.wu": regs3rev(NewAmmaxDbWu),
	"ammax_db.du": regs3rev(NewAmmaxDbDu),
	"ammin.w":     regs3rev(NewAmminW),
	"ammin.d":     regs3rev(NewAmminD),
	"ammin.wu":    regs3rev(NewAmminWu),
	"ammin.du":    regs3rev(NewAmminDu),
	"ammin_db.w":  regs3rev(NewAmminDbW),
	"ammin_db.d":  regs3rev(NewAmminDbD),
	"ammin_db.wu": regs3rev(NewAmminDbWu),
	"ammin_db.du": regs3rev(NewAmminDbDu),
	"csrrd":       regU14(NewCsrrd),
	"csrwr":       regU14(NewCsrwr),
	"csrxchg":     regRegU14(NewCsrxchg),
	"cacop":       u5RegI12(NewCacop),
	"lddir":       regRegU8(NewLddir),
	"ldpte":       regU8(NewLdpte),
	"iocsrrd.b":   regs2(NewIocsrrdB),
	"iocsrrd.h":   regs2(NewIocsrrdH),
	"iocsrrd.w":   regs2(NewIocsrrdW),
	"iocsrrd.d":   regs2(NewIocsrrdD),
	"iocsrwr.b":   regs2(NewIocsrwrB),
	"iocsrwr.h":   regs2(NewIocsrwrH),
	"iocsrwr.w":   regs2(NewIocsrwrW),
	"iocsrwr.d":   regs2(NewIocsrwrD),
	"tlbclr":      noOperands(NewTlbclr),
	"tlbflush":    noOperands(NewTlbflush),
	"tlbsrch":     noOperands(NewTlbsrch),
	"tlbrd":       noOperands(NewTlbrd),
	"tlbwr":       noOperands(NewTlbwr),
	"tlbfill":     noOperands(NewTlbfill),
	"ertn":        noOperands(NewErtn),
	"idle":        code15(NewIdle),
	"invtlb":      u5Regs2(NewInvtlb),
}

// --- operand validation helpers (kind/count; ranges are the roles') ---

// wantCount - the operand count.
func wantCount(ops []Op, n int) error {
	if len(ops) != n {
		return fmt.Errorf("want %d operands, got %d", n, len(ops))
	}

	return nil
}

// theReg - operand i as a register.
func theReg(ops []Op, i int) (Reg, error) {
	if ops[i].kind != opRegK {
		return Reg{}, fmt.Errorf("operand %d: want a register", i+1)
	}

	return NewReg(ops[i].reg), nil
}

// theNum - operand i as a number.
func theNum(ops []Op, i int) (int64, error) {
	if ops[i].kind != opNumK {
		return 0, fmt.Errorf("operand %d: want an immediate", i+1)
	}

	return ops[i].num, nil
}

// theRole - operand i as a validated role value.
func theRole[T any](ops []Op, i int, role func(int64) (T, error)) (T, error) {
	v, err := theNum(ops, i)
	if err != nil {
		var zero T

		return zero, err
	}

	return role(v)
}

// --- the shape factories ---

// regs3 - "rd, rj, rk" (and the am* "rd, rk, rj": the ctor parameter
// names carry the order).
func regs3(ctor func(a, b, c Reg) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		c, err := theReg(ops, 2)
		if err != nil {
			return nil, err
		}

		return ctor(a, b, c), nil
	}
}

// regs3rev - the am* operand order (rd, rk, rj): a regs3 twin named for
// the K/J slot swap (the ctor takes them in assembly order).
var regs3rev = regs3

// regs2 - two register operands (rd, rj; asrtle/asrtgt: rj, rk).
func regs2(ctor func(a, b Reg) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 2); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		return ctor(a, b), nil
	}
}

// noOperands - the operandless forms (tlbsrch and friends).
func noOperands(ctor func() Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 0); err != nil {
			return nil, err
		}

		return ctor(), nil
	}
}

// regI20 - "rd, si20" (lu12i.w and the pcaddi family).
func regI20(ctor func(rd Reg, v Imm20) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 2); err != nil {
			return nil, err
		}

		rd, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		v, err := theRole(ops, 1, NewImm20)
		if err != nil {
			return nil, err
		}

		return ctor(rd, v), nil
	}
}

// code15 - a single ui15 code (break/syscall/dbcl/dbar/ibar/idle).
func code15(ctor func(code Code15) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 1); err != nil {
			return nil, err
		}

		v, err := theRole(ops, 0, NewCode15)
		if err != nil {
			return nil, err
		}

		return ctor(v), nil
	}
}

// target - a single absolute target (b, bl).
func target(ctor func(t int64) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 1); err != nil {
			return nil, err
		}

		t, err := theNum(ops, 0)
		if err != nil {
			return nil, err
		}

		return ctor(t), nil
	}
}

// regTarget - "rj, target" (beqz/bnez).
func regTarget(ctor func(rj Reg, t int64) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 2); err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		t, err := theNum(ops, 1)
		if err != nil {
			return nil, err
		}

		return ctor(rj, t), nil
	}
}

// regs2Target - "rj, rd, target" (the compare-and-branch family).
func regs2Target(ctor func(a, b Reg, t int64) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		t, err := theNum(ops, 2)
		if err != nil {
			return nil, err
		}

		return ctor(a, b, t), nil
	}
}

// regs2Off - "rd, rj, offset" (jirl; the raw int64 offset).
func regs2Off(ctor func(a, b Reg, off int64) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		off, err := theNum(ops, 2)
		if err != nil {
			return nil, err
		}

		return ctor(a, b, off), nil
	}
}

// regregRole - "reg, reg, role" (the si12/ui12/shift families).
func regregRole[T any](
	ctor func(a, b Reg, v T) Instr,
	role func(int64) (T, error),
) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		v, err := theRole(ops, 2, role)
		if err != nil {
			return nil, err
		}

		return ctor(a, b, v), nil
	}
}

// The concrete regregRole bindings.
func regRegI12(ctor func(a, b Reg, v Imm12) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewImm12)
}

func regRegU12(ctor func(a, b Reg, v UImm12) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewUImm12)
}

func regRegI16(ctor func(a, b Reg, v Imm16) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewImm16)
}

func regRegU5(ctor func(a, b Reg, v UImm5) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewUImm5)
}

func regRegU6(ctor func(a, b Reg, v UImm6) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewUImm6)
}

func regRegI14(ctor func(a, b Reg, v Imm14) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewImm14)
}

func regRegU8(ctor func(a, b Reg, v UImm8) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewUImm8)
}

func regRegU14(ctor func(a, b Reg, v UImm14) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, NewUImm14)
}

// regs3Role - three registers plus a role operand (bytepick, alsl).
func regs3Role[T any](
	ctor func(a, b, c Reg, v T) Instr,
	role func(int64) (T, error),
) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 4); err != nil {
			return nil, err
		}

		a, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		b, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		c, err := theReg(ops, 2)
		if err != nil {
			return nil, err
		}

		v, err := theRole(ops, 3, role)
		if err != nil {
			return nil, err
		}

		return ctor(a, b, c, v), nil
	}
}

func regs3U2(ctor func(a, b, c Reg, v UImm2) Instr) func([]Op) (Instr, error) {
	return regs3Role(ctor, NewUImm2)
}

func regs3U3(ctor func(a, b, c Reg, v UImm3) Instr) func([]Op) (Instr, error) {
	return regs3Role(ctor, NewUImm3)
}

func regs3Shift(ctor func(a, b, c Reg, v Shift3) Instr) func([]Op) (Instr, error) {
	return regs3Role(ctor, NewShift3)
}

// regRegU5U5 - "rd, rj, msb, lsb" (bstrins/bstrpick .w; the field bounds
// cross-check: msb >= lsb).
func regRegU5U5(ctor func(rd, rj Reg, msb, lsb UImm5) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 4); err != nil {
			return nil, err
		}

		rd, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		msb, err := theRole(ops, 2, NewUImm5)
		if err != nil {
			return nil, err
		}

		lsb, err := theRole(ops, 3, NewUImm5)
		if err != nil {
			return nil, err
		}

		if msb.Val() < lsb.Val() {
			return nil, errors.New("msb is below lsb")
		}

		return ctor(rd, rj, msb, lsb), nil
	}
}

// regRegU6U6 - the .d bstrins/bstrpick form.
func regRegU6U6(ctor func(rd, rj Reg, msb, lsb UImm6) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 4); err != nil {
			return nil, err
		}

		rd, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		msb, err := theRole(ops, 2, NewUImm6)
		if err != nil {
			return nil, err
		}

		lsb, err := theRole(ops, 3, NewUImm6)
		if err != nil {
			return nil, err
		}

		if msb.Val() < lsb.Val() {
			return nil, errors.New("msb is below lsb")
		}

		return ctor(rd, rj, msb, lsb), nil
	}
}

// u5RegI12 - "ui5, rj, si12" (preld, cacop).
func u5RegI12(ctor func(op UImm5, rj Reg, off Imm12) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		op, err := theRole(ops, 0, NewUImm5)
		if err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		off, err := theRole(ops, 2, NewImm12)
		if err != nil {
			return nil, err
		}

		return ctor(op, rj, off), nil
	}
}

// u5Regs2 - "ui5, rj, rk" (preldx, invtlb).
func u5Regs2(ctor func(op UImm5, rj, rk Reg) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 3); err != nil {
			return nil, err
		}

		op, err := theRole(ops, 0, NewUImm5)
		if err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		rk, err := theReg(ops, 2)
		if err != nil {
			return nil, err
		}

		return ctor(op, rj, rk), nil
	}
}

// regU8 - "rj, ui8" (ldpte).
func regU8(ctor func(rj Reg, v UImm8) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 2); err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		v, err := theRole(ops, 1, NewUImm8)
		if err != nil {
			return nil, err
		}

		return ctor(rj, v), nil
	}
}

// regU14 - "rd, csr" (csrrd/csrwr).
func regU14(ctor func(rd Reg, v UImm14) Instr) func([]Op) (Instr, error) {
	return func(ops []Op) (Instr, error) {
		if err := wantCount(ops, 2); err != nil {
			return nil, err
		}

		rd, err := theReg(ops, 0)
		if err != nil {
			return nil, err
		}

		v, err := theRole(ops, 1, NewUImm14)
		if err != nil {
			return nil, err
		}

		return ctor(rd, v), nil
	}
}
