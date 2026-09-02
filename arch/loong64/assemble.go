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
	"add.w":       regs3(New().AddW),
	"add.d":       regs3(New().AddD),
	"sub.w":       regs3(New().SubW),
	"sub.d":       regs3(New().SubD),
	"slt":         regs3(New().Slt),
	"sltu":        regs3(New().Sltu),
	"maskeqz":     regs3(New().Maskeqz),
	"masknez":     regs3(New().Masknez),
	"nor":         regs3(New().Nor),
	"and":         regs3(New().And),
	"or":          regs3(New().Or),
	"xor":         regs3(New().Xor),
	"orn":         regs3(New().Orn),
	"andn":        regs3(New().Andn),
	"sll.w":       regs3(New().SllW),
	"srl.w":       regs3(New().SrlW),
	"sra.w":       regs3(New().SraW),
	"sll.d":       regs3(New().SllD),
	"srl.d":       regs3(New().SrlD),
	"sra.d":       regs3(New().SraD),
	"rotr.w":      regs3(New().RotrW),
	"rotr.d":      regs3(New().RotrD),
	"addi.w":      regRegI12(New().AddiW),
	"addi.d":      regRegI12(New().AddiD),
	"slti":        regRegI12(New().Slti),
	"sltui":       regRegI12(New().Sltui),
	"andi":        regRegU12(New().Andi),
	"ori":         regRegU12(New().Ori),
	"xori":        regRegU12(New().Xori),
	"addu16i.d":   regRegI16(New().Addu16iD),
	"slli.w":      regRegU5(New().SlliW),
	"srli.w":      regRegU5(New().SrliW),
	"srai.w":      regRegU5(New().SraiW),
	"rotri.w":     regRegU5(New().RotriW),
	"slli.d":      regRegU6(New().SlliD),
	"srli.d":      regRegU6(New().SrliD),
	"srai.d":      regRegU6(New().SraiD),
	"rotri.d":     regRegU6(New().RotriD),
	"lu12i.w":     regI20(New().Lu12iW),
	"lu32i.d":     regI20(New().Lu32iD),
	"lu52i.d":     regRegI12(New().Lu52iD),
	"pcaddi":      regI20(New().Pcaddi),
	"pcalau12i":   regI20(New().Pcalau12i),
	"pcaddu12i":   regI20(New().Pcaddu12i),
	"pcaddu18i":   regI20(New().Pcaddu18i),
	"ld.b":        regRegI12(New().LdB),
	"ld.h":        regRegI12(New().LdH),
	"ld.w":        regRegI12(New().LdW),
	"ld.wu":       regRegI12(New().LdWu),
	"ld.bu":       regRegI12(New().LdBu),
	"ld.hu":       regRegI12(New().LdHu),
	"ld.d":        regRegI12(New().LdD),
	"st.b":        regRegI12(New().StB),
	"st.h":        regRegI12(New().StH),
	"st.w":        regRegI12(New().StW),
	"st.d":        regRegI12(New().StD),
	"ldptr.w":     regRegI14(New().LdptrW),
	"ldptr.d":     regRegI14(New().LdptrD),
	"stptr.w":     regRegI14(New().StptrW),
	"stptr.d":     regRegI14(New().StptrD),
	"ldx.b":       regs3(New().LdxB),
	"ldx.h":       regs3(New().LdxH),
	"ldx.w":       regs3(New().LdxW),
	"ldx.wu":      regs3(New().LdxWu),
	"ldx.bu":      regs3(New().LdxBu),
	"ldx.hu":      regs3(New().LdxHu),
	"ldx.d":       regs3(New().LdxD),
	"stx.b":       regs3(New().StxB),
	"stx.h":       regs3(New().StxH),
	"stx.w":       regs3(New().StxW),
	"stx.d":       regs3(New().StxD),
	"preld":       u5RegI12(New().Preld),
	"preldx":      u5Regs2(New().Preldx),
	"beq":         regs2Target(New().Beq),
	"bne":         regs2Target(New().Bne),
	"blt":         regs2Target(New().Blt),
	"bge":         regs2Target(New().Bge),
	"bltu":        regs2Target(New().Bltu),
	"bgeu":        regs2Target(New().Bgeu),
	"beqz":        regTarget(New().Beqz),
	"bnez":        regTarget(New().Bnez),
	"b":           target(New().B),
	"bl":          target(New().Bl),
	"jirl":        regs2Off(New().Jirl),
	"break":       code15(New().Break),
	"syscall":     code15(New().Syscall),
	"dbcl":        code15(New().Dbcl),
	"dbar":        code15(New().Dbar),
	"ibar":        code15(New().Ibar),
	"mul.w":       regs3(New().MulW),
	"mulh.w":      regs3(New().MulhW),
	"mulh.wu":     regs3(New().MulhWu),
	"mul.d":       regs3(New().MulD),
	"mulh.d":      regs3(New().MulhD),
	"mulh.du":     regs3(New().MulhDu),
	"mulw.d.w":    regs3(New().MulwDW),
	"mulw.d.wu":   regs3(New().MulwDWu),
	"div.w":       regs3(New().DivW),
	"mod.w":       regs3(New().ModW),
	"div.wu":      regs3(New().DivWu),
	"mod.wu":      regs3(New().ModWu),
	"div.d":       regs3(New().DivD),
	"mod.d":       regs3(New().ModD),
	"div.du":      regs3(New().DivDu),
	"mod.du":      regs3(New().ModDu),
	"ext.w.b":     regs2(New().ExtWB),
	"ext.w.h":     regs2(New().ExtWH),
	"rdtimel.w":   regs2(New().RdtimelW),
	"rdtimeh.w":   regs2(New().RdtimehW),
	"rdtime.d":    regs2(New().RdtimeD),
	"cpucfg":      regs2(New().Cpucfg),
	"clo.w":       regs2(New().CloW),
	"clz.w":       regs2(New().ClzW),
	"cto.w":       regs2(New().CtoW),
	"ctz.w":       regs2(New().CtzW),
	"clo.d":       regs2(New().CloD),
	"clz.d":       regs2(New().ClzD),
	"cto.d":       regs2(New().CtoD),
	"ctz.d":       regs2(New().CtzD),
	"revb.2h":     regs2(New().Revb2H),
	"revb.4h":     regs2(New().Revb4H),
	"revb.2w":     regs2(New().Revb2W),
	"revb.d":      regs2(New().RevbD),
	"revh.2w":     regs2(New().Revh2W),
	"revh.d":      regs2(New().RevhD),
	"bitrev.4b":   regs2(New().Revbit4B),
	"bitrev.w":    regs2(New().RevbitW),
	"bitrev.8b":   regs2(New().Revbit8B),
	"bitrev.d":    regs2(New().RevbitD),
	"bstrins.w":   regRegU5U5(New().BstrinsW),
	"bstrpick.w":  regRegU5U5(New().BstrpickW),
	"bstrins.d":   regRegU6U6(New().BstrinsD),
	"bstrpick.d":  regRegU6U6(New().BstrpickD),
	"alsl.w":      regs3Shift(New().AlslW),
	"alsl.wu":     regs3Shift(New().AlslWu),
	"alsl.d":      regs3Shift(New().AlslD),
	"bytepick.w":  regs3U2(New().BytepickW),
	"bytepick.d":  regs3U3(New().BytepickD),
	"crc.w.b.w":   regs3(New().CrcWBW),
	"crc.w.h.w":   regs3(New().CrcWHW),
	"crc.w.w.w":   regs3(New().CrcWWW),
	"crc.w.d.w":   regs3(New().CrcWDW),
	"crcc.w.b.w":  regs3(New().CrccWBW),
	"crcc.w.h.w":  regs3(New().CrccWHW),
	"crcc.w.w.w":  regs3(New().CrccWWW),
	"crcc.w.d.w":  regs3(New().CrccWDW),
	"asrtle.d":    regs2(New().AsrtleD),
	"asrtgt.d":    regs2(New().AsrtgtD),
	"ldgt.b":      regs3(New().LdgtB),
	"ldgt.h":      regs3(New().LdgtH),
	"ldgt.w":      regs3(New().LdgtW),
	"ldgt.d":      regs3(New().LdgtD),
	"ldle.b":      regs3(New().LdleB),
	"ldle.h":      regs3(New().LdleH),
	"ldle.w":      regs3(New().LdleW),
	"ldle.d":      regs3(New().LdleD),
	"stgt.b":      regs3(New().StgtB),
	"stgt.h":      regs3(New().StgtH),
	"stgt.w":      regs3(New().StgtW),
	"stgt.d":      regs3(New().StgtD),
	"stle.b":      regs3(New().StleB),
	"stle.h":      regs3(New().StleH),
	"stle.w":      regs3(New().StleW),
	"stle.d":      regs3(New().StleD),
	"ll.w":        regRegI14(New().LlW),
	"ll.d":        regRegI14(New().LlD),
	"sc.w":        regRegI14(New().ScW),
	"sc.d":        regRegI14(New().ScD),
	"llacq.w":     regs2(New().LlacqW),
	"llacq.d":     regs2(New().LlacqD),
	"screl.w":     regs2(New().ScrelW),
	"screl.d":     regs2(New().ScrelD),
	"sc.q":        regs3rev(New().ScQ),
	"amcas.b":     regs3rev(New().AmcasB),
	"amcas.h":     regs3rev(New().AmcasH),
	"amcas.w":     regs3rev(New().AmcasW),
	"amcas.d":     regs3rev(New().AmcasD),
	"amcas_db.b":  regs3rev(New().AmcasDbB),
	"amcas_db.h":  regs3rev(New().AmcasDbH),
	"amcas_db.w":  regs3rev(New().AmcasDbW),
	"amcas_db.d":  regs3rev(New().AmcasDbD),
	"amswap.b":    regs3rev(New().AmswapB),
	"amswap.h":    regs3rev(New().AmswapH),
	"amswap.w":    regs3rev(New().AmswapW),
	"amswap.d":    regs3rev(New().AmswapD),
	"amswap_db.b": regs3rev(New().AmswapDbB),
	"amswap_db.h": regs3rev(New().AmswapDbH),
	"amswap_db.w": regs3rev(New().AmswapDbW),
	"amswap_db.d": regs3rev(New().AmswapDbD),
	"amadd.b":     regs3rev(New().AmaddB),
	"amadd.h":     regs3rev(New().AmaddH),
	"amadd.w":     regs3rev(New().AmaddW),
	"amadd.d":     regs3rev(New().AmaddD),
	"amadd_db.b":  regs3rev(New().AmaddDbB),
	"amadd_db.h":  regs3rev(New().AmaddDbH),
	"amadd_db.w":  regs3rev(New().AmaddDbW),
	"amadd_db.d":  regs3rev(New().AmaddDbD),
	"amand.w":     regs3rev(New().AmandW),
	"amand.d":     regs3rev(New().AmandD),
	"amand_db.w":  regs3rev(New().AmandDbW),
	"amand_db.d":  regs3rev(New().AmandDbD),
	"amor.w":      regs3rev(New().AmorW),
	"amor.d":      regs3rev(New().AmorD),
	"amor_db.w":   regs3rev(New().AmorDbW),
	"amor_db.d":   regs3rev(New().AmorDbD),
	"amxor.w":     regs3rev(New().AmxorW),
	"amxor.d":     regs3rev(New().AmxorD),
	"amxor_db.w":  regs3rev(New().AmxorDbW),
	"amxor_db.d":  regs3rev(New().AmxorDbD),
	"ammax.w":     regs3rev(New().AmmaxW),
	"ammax.d":     regs3rev(New().AmmaxD),
	"ammax.wu":    regs3rev(New().AmmaxWu),
	"ammax.du":    regs3rev(New().AmmaxDu),
	"ammax_db.w":  regs3rev(New().AmmaxDbW),
	"ammax_db.d":  regs3rev(New().AmmaxDbD),
	"ammax_db.wu": regs3rev(New().AmmaxDbWu),
	"ammax_db.du": regs3rev(New().AmmaxDbDu),
	"ammin.w":     regs3rev(New().AmminW),
	"ammin.d":     regs3rev(New().AmminD),
	"ammin.wu":    regs3rev(New().AmminWu),
	"ammin.du":    regs3rev(New().AmminDu),
	"ammin_db.w":  regs3rev(New().AmminDbW),
	"ammin_db.d":  regs3rev(New().AmminDbD),
	"ammin_db.wu": regs3rev(New().AmminDbWu),
	"ammin_db.du": regs3rev(New().AmminDbDu),
	"csrrd":       regU14(New().Csrrd),
	"csrwr":       regU14(New().Csrwr),
	"csrxchg":     regRegU14(New().Csrxchg),
	"cacop":       u5RegI12(New().Cacop),
	"lddir":       regRegU8(New().Lddir),
	"ldpte":       regU8(New().Ldpte),
	"iocsrrd.b":   regs2(New().IocsrrdB),
	"iocsrrd.h":   regs2(New().IocsrrdH),
	"iocsrrd.w":   regs2(New().IocsrrdW),
	"iocsrrd.d":   regs2(New().IocsrrdD),
	"iocsrwr.b":   regs2(New().IocsrwrB),
	"iocsrwr.h":   regs2(New().IocsrwrH),
	"iocsrwr.w":   regs2(New().IocsrwrW),
	"iocsrwr.d":   regs2(New().IocsrwrD),
	"tlbclr":      noOperands(New().Tlbclr),
	"tlbflush":    noOperands(New().Tlbflush),
	"tlbsrch":     noOperands(New().Tlbsrch),
	"tlbrd":       noOperands(New().Tlbrd),
	"tlbwr":       noOperands(New().Tlbwr),
	"tlbfill":     noOperands(New().Tlbfill),
	"ertn":        noOperands(New().Ertn),
	"idle":        code15(New().Idle),
	"invtlb":      u5Regs2(New().Invtlb),
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

	return newReg(ops[i].reg), nil
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

		v, err := theRole(ops, 1, New().Imm20)
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

		v, err := theRole(ops, 0, New().Code15)
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
	return regregRole(ctor, New().Imm12)
}

func regRegU12(ctor func(a, b Reg, v UImm12) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().UImm12)
}

func regRegI16(ctor func(a, b Reg, v Imm16) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().Imm16)
}

func regRegU5(ctor func(a, b Reg, v UImm5) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().UImm5)
}

func regRegU6(ctor func(a, b Reg, v UImm6) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().UImm6)
}

func regRegI14(ctor func(a, b Reg, v Imm14) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().Imm14)
}

func regRegU8(ctor func(a, b Reg, v UImm8) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().UImm8)
}

func regRegU14(ctor func(a, b Reg, v UImm14) Instr) func([]Op) (Instr, error) {
	return regregRole(ctor, New().UImm14)
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
	return regs3Role(ctor, New().UImm2)
}

func regs3U3(ctor func(a, b, c Reg, v UImm3) Instr) func([]Op) (Instr, error) {
	return regs3Role(ctor, New().UImm3)
}

func regs3Shift(ctor func(a, b, c Reg, v Shift3) Instr) func([]Op) (Instr, error) {
	return regs3Role(ctor, New().Shift3)
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

		msb, err := theRole(ops, 2, New().UImm5)
		if err != nil {
			return nil, err
		}

		lsb, err := theRole(ops, 3, New().UImm5)
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

		msb, err := theRole(ops, 2, New().UImm6)
		if err != nil {
			return nil, err
		}

		lsb, err := theRole(ops, 3, New().UImm6)
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

		op, err := theRole(ops, 0, New().UImm5)
		if err != nil {
			return nil, err
		}

		rj, err := theReg(ops, 1)
		if err != nil {
			return nil, err
		}

		off, err := theRole(ops, 2, New().Imm12)
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

		op, err := theRole(ops, 0, New().UImm5)
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

		v, err := theRole(ops, 1, New().UImm8)
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

		v, err := theRole(ops, 1, New().UImm14)
		if err != nil {
			return nil, err
		}

		return ctor(rd, v), nil
	}
}
