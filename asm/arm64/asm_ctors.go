package arm64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Local aliases to the arch operand model: these files moved from
// arch/arm64 (the assembler constructors are asm-layer code); the
// historical names keep working.

type (
	vOp    = arch.VOp
	vMem   = arch.VMem
	Schema = arch.Schema
	Field  = arch.Field
	fpKind = arch.FpKind
	Csel   = arch.Csel
)

// fp/mem helpers of the moved bodies.
func decodeArrangement(q, size uint32) string       { return arch.DecodeArrangement(q, size) }
func getSchemas() []Schema                          { return arch.Schemas() }
func shiftAmt(op vOp) int64                         { return arch.ShiftAmt(op) }
func wantTarget(op vOp, name string) (int64, error) { return arch.WantTarget(op, name) }

const (
	kS = arch.FS
	kD = arch.FD
	kW = arch.FW
	kX = arch.FX
)

// Instr — the arch instruction interface.
type Instr = arch.Instr

const (
	ArmOpReg    = arch.ArmOpReg
	ArmOpImm    = arch.ArmOpImm
	ArmOpLit    = arch.ArmOpLit
	ArmOpMem    = arch.ArmOpMem
	ArmOpList   = arch.ArmOpList
	ArmOpShift  = arch.ArmOpShift
	ArmOpExtend = arch.ArmOpExtend
	ArmOpFloat  = arch.ArmOpFloat
)

type Builder = arch.Builder

// Wrappers to the shared operand helpers (arch api.go exports them).
func wantAReg(op vOp, name string) (string, error)           { return arch.WantAReg(op, name) }
func wantCond(op vOp, name string) (string, error)           { return arch.WantCond(op, name) }
func armReg2(ops []vOp, name string) (string, string, error) { return arch.ArmReg2(ops, name) }
func regNums2(a, b string) (uint32, uint32, error)           { return arch.RegNums2(a, b) }
func zeroReg(rd string) string                               { return arch.ZeroReg(rd) }
func arrQSize(arr string) (uint32, uint32, error)            { return arch.ArrQSize(arr) }
func armRegNum(name string) (uint32, error)                  { return arch.ArmRegNum(name) }

func addSubRegName(
	n uint32,
	isf, flags bool,
) string {
	return arch.AddSubRegName(n, isf, flags)
}

func newCsel(
	rd, rn, rm, cond string,
) (arch.Csel, error) {
	return arch.NewCsel(rd, rn, rm, cond)
}
func regIndex(reg string) uint32              { return arch.RegIndex(reg) }
func regListStr(rt0 uint32, count int) string { return arch.RegListStr(rt0, count) }
func vfpExpandImm64(imm8 uint32) float64      { return arch.VfpExpandImm64(imm8) }
func vfpExpandImm32(imm8 uint32) float32      { return arch.VfpExpandImm32(imm8) }
func isSimd3Logical(name string) bool         { return arch.IsSimd3Logical(name) }

// armCtors — constructors of per-instruction structs from parsed text:
// mnemonic (with .arr/.cond suffixes, where the grammar yields them
// whole) → an operand constructor. Real instructions only; aliases
// (cmp/mov/cset/sxt*/...) — asm/arm64/alias on top of NewWithCtors. The
// hook runs before the armFieldsFor legacy path; the encoding choice is
// confirmed by encodeARM's self-verify.
//
// Built once by buildArmCtors (literal entries + the by-element
// registrations); a data table, immutable after construction.
var armCtors = buildArmCtors()

func buildArmCtors() map[string]func(ops []vOp) (Instr, error) {
	m := map[string]func(ops []vOp) (Instr, error){
		"b":          newB,
		"bl":         newBl,
		"cbz":        newCbz,
		"cbnz":       newCbnz,
		"nop":        newNop,
		"ret":        newRet,
		"br":         newBr,
		"blr":        newBlr,
		"rbit":       newRbit,
		"rev16":      newRev16,
		"rev32":      newRev32,
		"rev":        newRev,
		"clz":        newClz,
		"cls":        newCls,
		"udiv":       newUdiv,
		"sdiv":       newSdiv,
		"smulh":      newSmulh,
		"umulh":      newUmulh,
		"adc":        newAdc,
		"ccmp":       newCcmp,
		"extr":       newExtr,
		"madd":       newMadd,
		"msub":       newMsub,
		"csel":       newCselArm,
		"csinc":      newCsinc,
		"csinv":      newCsinv,
		"csneg":      newCsneg,
		"movz":       newMovz,
		"movn":       newMovn,
		"movk":       newMovk,
		"add":        newAddArm,
		"adds":       newAddsArm,
		"sub":        newSubArm,
		"subs":       newSubsArm,
		"and":        newAndArm,
		"ands":       newAndsArm,
		"orr":        newOrrArm,
		"eor":        newEorArm,
		"bic":        newBicArm,
		"bics":       newBicsArm,
		"orn":        newOrnArm,
		"eon":        newEonArm,
		"ldr":        newLdrArm,
		"ldrb":       newLdrbArm,
		"ldrh":       newLdrhArm,
		"str":        newStrArm,
		"strb":       newStrbArm,
		"strh":       newStrhArm,
		"ldur":       newLdursArm,
		"stur":       newStursArm,
		"ldurb":      newLdurbArm,
		"ldurh":      newLdurhArm,
		"sturb":      newSturbArm,
		"sturh":      newSturhArm,
		"ldrsw":      newLdrswArm,
		"ldrsb":      newLdrsbArm,
		"ldrsh":      newLdrshArm,
		"ldp":        newLdpArm,
		"stp":        newStpArm,
		"ldar":       newLdarArm,
		"stlr":       newStlrArm,
		"stlxr":      newStlxrArm,
		"stxrb":      newStxrbArm,
		"lsl":        newLslArm,
		"lsr":        newLsrArm,
		"asr":        newAsrArm,
		"ror":        newRorArm2,
		"adr":        newAdr,
		"adrp":       newAdrp,
		"tbz":        newTbzArm(false),
		"tbnz":       newTbzArm(true),
		"svc":        newSvc,
		"brk":        newBrkArm,
		"udf":        newUdfArm,
		"hlt":        newHlt,
		"hvc":        newHvc,
		"mrs":        newMrsArm,
		"msr":        newMsrArm,
		"fadd":       newFp3Arm("fadd", 0x1E602800, 0x1E202800),
		"fsub":       newFp3Arm("fsub", 0x1E603800, 0x1E203800),
		"fmul":       newFp3Arm("fmul", 0x1E600800, 0x1E200800),
		"fdiv":       newFp3Arm("fdiv", 0x1E601800, 0x1E201800),
		"fmax":       newFp3Arm("fmax", 0x1E604800, 0x1E204800),
		"fmin":       newFp3Arm("fmin", 0x1E605800, 0x1E205800),
		"fneg":       newFp2Arm("fneg", 0x1E614000),
		"fmov":       newFmov,
		"fcmp":       newFcmpArm,
		"fmadd":      newFmadd("fmadd", 0x1F000000),
		"fnmsub":     newFmadd("fnmsub", 0x1F200000),
		"and.16b":    newSimd3("and", 0x0E201C00),
		"and.8b":     newSimd3("and", 0x0E201C00),
		"bic.16b":    newSimd3("bic", 0x0E601C00),
		"bic.8b":     newSimd3("bic", 0x0E601C00),
		"orr.16b":    newSimd3("orr", 0x0EA01C00),
		"orr.8b":     newSimd3("orr", 0x0EA01C00),
		"orn.16b":    newSimd3("orn", 0x0EE01C00),
		"orn.8b":     newSimd3("orn", 0x0EE01C00),
		"eor.16b":    newSimd3("eor", 0x2E201C00),
		"eor.8b":     newSimd3("eor", 0x2E201C00),
		"bsl.16b":    newSimd3("bsl", 0x2E601C00),
		"bsl.8b":     newSimd3("bsl", 0x2E601C00),
		"bit.16b":    newSimd3("bit", 0x2EA01C00),
		"bit.8b":     newSimd3("bit", 0x2EA01C00),
		"bif.16b":    newSimd3("bif", 0x2EE01C00),
		"bif.8b":     newSimd3("bif", 0x2EE01C00),
		"add.16b":    newSimd3("add", 0x0E208400),
		"add.8b":     newSimd3("add", 0x0E208400),
		"cmeq.16b":   newSimd3("cmeq", 0x2E208C00),
		"addp.16b":   newSimd3("addp", 0x0E20BC00),
		"sqrshl.16b": newSimd3("sqrshl", 0x0E205C00),
		"cnt.8b":     newSimd2("cnt", 0x0E205800),
		"cnt.16b":    newSimd2("cnt", 0x0E205800),
		"rev32.8b":   newSimd2("rev32", 0x2E200800),
		"rev32.16b":  newSimd2("rev32", 0x2E200800),
		"not.16b":    newSimd2("not", 0x2E205800),
		"not.8b":     newSimd2("not", 0x2E205800),
		"abs.16b":    newSimd2("abs", 0x0E20B800),
		"abs.8b":     newSimd2("abs", 0x0E20B800),
		"rbit.16b":   newSimd2("rbit", 0x2E205800),
		"rbit.8b":    newSimd2("rbit", 0x2E205800),
		"shl.16b":    newSimdShift("shl", 0x0F004000),
		"shl.8b":     newSimdShift("shl", 0x0F004000),
		"sri.16b":    newSimdShift("sri", 0x2F004000),
		"ushr.16b":   newSimdShift("ushr", 0x2F000000),
		"sshr.16b":   newSimdShift("sshr", 0x0F000000),
		"aese":       newAes("aese", 0x4E284800),
		"aesmc":      newAes("aesmc", 0x4E286800),
		"dup.16b":    newDupArm,
		"dup.8b":     newDupArm,
		"dup.4s":     newDupArm,
		"dup.2s":     newDupArm,
		"dup.2d":     newDupArm,
		"dup.4h":     newDupArm,
		"dup.8h":     newDupArm,
		"mov.b":      newMovInsArm(0),
		"mov.h":      newMovInsArm(1), // INS: mov.sz vd[idx], wn
		"mov.s":      newMovInsArm(2),
		"mov.d":      newMovInsArm(3),
		"ins.b":      newInsElemArm(0),
		"ins.h":      newInsElemArm(1), // INS (element): ins.sz vd[idx], vn[idx]
		"ins.s":      newInsElemArm(2),
		"ins.d":      newInsElemArm(3),
		"smov":       newSmovUmovArm("smov"),
		"umov":       newSmovUmovArm("umov"),
		"saddw.4h": newSimdWidenArm(
			"saddw",
			0x0E201000,
		),
		"saddw.2s": newSimdWidenArm("saddw", 0x0E201000),
		"saddw.1d": newSimdWidenArm("saddw", 0x0E201000),
		"saddw2.8h": newSimdWidenArm(
			"saddw2",
			0x0E201000,
		),
		"saddw2.4s": newSimdWidenArm("saddw2", 0x0E201000),
		"saddw2.2d": newSimdWidenArm("saddw2", 0x0E201000),
		"ssubw.4h": newSimdWidenArm(
			"ssubw",
			0x0E203000,
		),
		"ssubw.2s": newSimdWidenArm("ssubw", 0x0E203000),
		"ssubw.1d": newSimdWidenArm("ssubw", 0x0E203000),
		"ssubw2.8h": newSimdWidenArm(
			"ssubw2",
			0x0E203000,
		),
		"ssubw2.4s": newSimdWidenArm("ssubw2", 0x0E203000),
		"ssubw2.2d": newSimdWidenArm("ssubw2", 0x0E203000),
		"uaddw.4h": newSimdWidenArm(
			"uaddw",
			0x2E201000,
		),
		"uaddw.2s": newSimdWidenArm("uaddw", 0x2E201000),
		"uaddw.1d": newSimdWidenArm("uaddw", 0x2E201000),
		"uaddw2.8h": newSimdWidenArm(
			"uaddw2",
			0x2E201000,
		),
		"uaddw2.4s": newSimdWidenArm("uaddw2", 0x2E201000),
		"uaddw2.2d": newSimdWidenArm("uaddw2", 0x2E201000),
		"usubw.4h": newSimdWidenArm(
			"usubw",
			0x2E203000,
		),
		"usubw.2s": newSimdWidenArm("usubw", 0x2E203000),
		"usubw.1d": newSimdWidenArm("usubw", 0x2E203000),
		"usubw2.8h": newSimdWidenArm(
			"usubw2",
			0x2E203000,
		),
		"usubw2.4s":  newSimdWidenArm("usubw2", 0x2E203000),
		"usubw2.2d":  newSimdWidenArm("usubw2", 0x2E203000),
		"tbl.16b":    newTblArm,
		"uaddlv.16b": newUaddlv,
		"uaddlv.8b":  newUaddlv,
		"uaddlv.4h":  newUaddlv,
		"uaddlv.8h":  newUaddlv,
		"dmb":        newDmb,
		"yield":      newYield,
		"dc":         newDc,
		"prfm":       newPrfmArm,
		"ld1.16b":    newLdStruct("ld1"),
		"ld1.8b":     newLdStruct("ld1"),
		"ld1.4s":     newLdStruct("ld1"),
		"ld1.2s":     newLdStruct("ld1"),
		"ld1.2d":     newLdStruct("ld1"),
		"ld1.4h":     newLdStruct("ld1"),
		"ld1.8h":     newLdStruct("ld1"),
		"ld1r.16b":   newLdStruct("ld1r"),
		"ld1r.8b":    newLdStruct("ld1r"),
		"ld1r.4s":    newLdStruct("ld1r"),
		"ld1r.2s":    newLdStruct("ld1r"),
		"ld1r.2d":    newLdStruct("ld1r"),
		"ld1r.4h":    newLdStruct("ld1r"),
		"ld1r.8h":    newLdStruct("ld1r"),
		"ld2.16b":    newLdStruct("ld2"),
		"ld2.8b":     newLdStruct("ld2"),
		"ld3.16b":    newLdStruct("ld3"),
		"ld3.8b":     newLdStruct("ld3"),
		"ld4.16b":    newLdStruct("ld4"),
		"ld4.8b":     newLdStruct("ld4"),
		"st1.16b":    newLdStruct("st1"),
		"st1.8b":     newLdStruct("st1"),
		"st2.16b":    newLdStruct("st2"),
		"st3.16b":    newLdStruct("st3"),
		"st4.16b":    newLdStruct("st4"),
		"mov.16b":    newMovSimd("16b"),
		"mov.8b":     newMovSimd("8b"),
		"b.eq":       newBcondOf("eq"),
		"b.ne":       newBcondOf("ne"),
		"b.hs":       newBcondOf("hs"),
		"b.lo":       newBcondOf("lo"),
		"b.mi":       newBcondOf("mi"),
		"b.pl":       newBcondOf("pl"),
		"b.vs":       newBcondOf("vs"),
		"b.vc":       newBcondOf("vc"),
		"b.hi":       newBcondOf("hi"),
		"b.ls":       newBcondOf("ls"),
		"b.ge":       newBcondOf("ge"),
		"b.lt":       newBcondOf("lt"),
		"b.gt":       newBcondOf("gt"),
		"b.le":       newBcondOf("le"),
		"b.al":       newBcondOf("al"),
		"b.nv":       newBcondOf("nv"),
	}

	registerByElem(m)

	return m
}

// Exported for the alias layer (was arch api.go wrappers): the family
// selection and csel assembly by base name.
func IsGPR(name string) bool { return isGPR(name) }
func AddSubThird(ops []vOp, base string, rdN, rnN uint32, idx int) (Instr, error) {
	return addSubThird(ops, base, rdN, rnN, idx)
}
func MaddCtor(ops []vOp, name string) (Instr, error)    { return makeMaddCtor(ops, name) }
func Msub3(ops []vOp, name string) (Instr, error)       { return msub3(ops, name) }
func CsOf(base, rd, rn, rm, cond string) (Instr, error) { return csOf(base, rd, rn, rm, cond) }

// buildInstr — the armCtors dispatch (was arch.BuildInstr).
func buildInstr(mnem string, ops []vOp) (Instr, error) {
	ctor, ok := armCtors[mnem]
	if !ok {
		return nil, fmt.Errorf("unknown instruction %q", mnem)
	}

	st, err := ctor(ops)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", mnem, err)
	}

	return st, nil
}

func ldStructDecode(opcode, size, q, l uint32) (string, string, int, bool) {
	return arch.LdStructDecode(opcode, size, q, l)
}

// Shims to the arch *Of family constructors.
func AdcOf(rd string, rn string, rm string) (arch.Adc, error) {
	return arch.AdcOf(rd, rn, rm)
}

func AddExtOf(
	rdNum, rnNum, rmNum uint32,
	option string,
	imm3 uint32,
	isf bool,
) (arch.AddExt, error) {
	return arch.AddExtOf(rdNum, rnNum, rmNum, option, imm3, isf)
}
func AddImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (arch.AddImm, error) {
	return arch.AddImmOf(rdNum, rnNum, imm12, shift, isf)
}
func AddShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.AddShift, error) {
	return arch.AddShiftOf(rd, rn, rm, imm6, shift, isf)
}

func AddsExtOf(
	rdNum, rnNum, rmNum uint32,
	option string,
	imm3 uint32,
	isf bool,
) (arch.AddsExt, error) {
	return arch.AddsExtOf(rdNum, rnNum, rmNum, option, imm3, isf)
}
func AddsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (arch.AddsImm, error) {
	return arch.AddsImmOf(rdNum, rnNum, imm12, shift, isf)
}
func AddsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.AddsShift, error) {
	return arch.AddsShiftOf(rd, rn, rm, imm6, shift, isf)
}
func AdrOf(rd string, off int64) (arch.Adr, error) {
	return arch.AdrOf(rd, off)
}
func AdrpOf(rd string, off int64) (arch.Adrp, error) {
	return arch.AdrpOf(rd, off)
}
func AndImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (arch.AndImm, error) {
	return arch.AndImmOf(rd, rn, immr, imms, n, is64)
}
func AndShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.AndShift, error) {
	return arch.AndShiftOf(rd, rn, rm, imm6, shift, isf)
}
func AsrRegOf(rd string, rn string, rm string) (arch.AsrReg, error) {
	return arch.AsrRegOf(rd, rn, rm)
}
func BOf(target int64) (arch.B, error) {
	return arch.BOf(target)
}
func BcondOf(cond string, target int64) (arch.Bcond, error) {
	return arch.BcondOf(cond, target)
}
func BicShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.BicShift, error) {
	return arch.BicShiftOf(rd, rn, rm, imm6, shift, isf)
}
func BicsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.BicsShift, error) {
	return arch.BicsShiftOf(rd, rn, rm, imm6, shift, isf)
}
func BlOf(target int64) (arch.Bl, error) {
	return arch.BlOf(target)
}
func BlrOf(rn string) (arch.Blr, error) {
	return arch.BlrOf(rn)
}
func BrOf(rn string) (arch.Br, error) {
	return arch.BrOf(rn)
}
func CbnzOf(rt string, target int64) (arch.Cbnz, error) {
	return arch.CbnzOf(rt, target)
}
func CbzOf(rt string, target int64) (arch.Cbz, error) {
	return arch.CbzOf(rt, target)
}
func CcmpOf(rn string, rm string, immVal uint32, cond string) (arch.Ccmp, error) {
	return arch.CcmpOf(rn, rm, immVal, cond)
}
func ClsOf(rd string, rn string) (arch.Cls, error) {
	return arch.ClsOf(rd, rn)
}
func ClzOf(rd string, rn string) (arch.Clz, error) {
	return arch.ClzOf(rd, rn)
}
func CselOf(rd string, rn string, rm string, cond string) (arch.Csel, error) {
	return arch.CselOf(rd, rn, rm, cond)
}
func CsincOf(c Csel) (arch.Csinc, error) {
	return arch.CsincOf(c)
}
func CsinvOf(c Csel) (arch.Csinv, error) {
	return arch.CsinvOf(c)
}
func CsnegOf(c Csel) (arch.Csneg, error) {
	return arch.CsnegOf(c)
}
func EonShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.EonShift, error) {
	return arch.EonShiftOf(rd, rn, rm, imm6, shift, isf)
}
func EorImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (arch.EorImm, error) {
	return arch.EorImmOf(rd, rn, immr, imms, n, is64)
}
func EorShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (arch.EorShift, error) {
	return arch.EorShiftOf(rd, rn, rm, imm6, shift, isf)
}
func ExtrOf(rd string, rn string, rm string, lsb uint32) (arch.Extr, error) {
	return arch.ExtrOf(rd, rn, rm, lsb)
}
func LdarOf(rt, rn string, enc uint32) (arch.Ldar, error) {
	return arch.LdarOf(rt, rn, enc)
}
func LdpOf(rt, rt2, rn string, kind arch.MemKind, off int64, scale, enc uint32) (arch.Ldp, error) {
	return arch.LdpOf(rt, rt2, rn, kind, off, scale, enc)
}
func LdrLitOf(rt string, lit int64, enc uint32) (arch.Ldr, error) {
	return arch.LdrLitOf(rt, lit, enc)
}

func LdrOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldr, error) {
	return arch.LdrOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LdrbOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldrb, error) {
	return arch.LdrbOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LdrhOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldrh, error) {
	return arch.LdrhOf(rt, rn, kind, off, enc, rm, option, amt)
}
func LdrsbOf(rt string, rn string, off int64) (arch.Ldrsb, error) {
	return arch.LdrsbOf(rt, rn, off)
}
func LdrshOf(rt string, rn string, off int64) (arch.Ldrsh, error) {
	return arch.LdrshOf(rt, rn, off)
}

func LdrswOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldrsw, error) {
	return arch.LdrswOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LdurOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldur, error) {
	return arch.LdurOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LdurbOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldurb, error) {
	return arch.LdurbOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LdurhOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Ldurh, error) {
	return arch.LdurhOf(rt, rn, kind, off, enc, rm, option, amt)
}

func LslRegOf(
	rd string,
	rn string,
	rm string,
) (arch.LslReg, error) {
	return arch.LslRegOf(rd, rn, rm)
}

func LsrRegOf(
	rd string,
	rn string,
	rm string,
) (arch.LsrReg, error) {
	return arch.LsrRegOf(rd, rn, rm)
}

func MaddOf(
	rd string,
	rn string,
	rm string,
	ra string,
) (arch.Madd, error) {
	return arch.MaddOf(rd, rn, rm, ra)
}

func MovkOf(
	rd string,
	imm16 uint32,
	hw uint32,
) (arch.Movk, error) {
	return arch.MovkOf(rd, imm16, hw)
}
func MrsOf(rd string, sysreg string) (arch.Mrs, error) {
	return arch.MrsOf(rd, sysreg)
}
func MsrOf(rt string, sysreg string) (arch.Msr, error) {
	return arch.MsrOf(rt, sysreg)
}

func MsubOf(
	rd, rn, rm, ra string,
) (arch.Msub, error) {
	return arch.MsubOf(rd, rn, rm, ra)
}
func PrfmOf(rn string) (arch.Prfm, error) {
	return arch.PrfmOf(rn)
}
func RbitOf(rd string, rn string) (arch.Rbit, error) {
	return arch.RbitOf(rd, rn)
}
func RetOf(rn string) (arch.Ret, error) {
	return arch.RetOf(rn)
}
func Rev16Of(rd string, rn string) (arch.Rev16, error) {
	return arch.Rev16Of(rd, rn)
}
func Rev32Of(rd string, rn string) (arch.Rev32, error) {
	return arch.Rev32Of(rd, rn)
}
func RevOf(rd string, rn string) (arch.Rev, error) {
	return arch.RevOf(rd, rn)
}

func RorRegOf(
	rd string,
	rn string,
	rm string,
) (arch.RorReg, error) {
	return arch.RorRegOf(rd, rn, rm)
}
func SdivOf(rd string, rn string, rm string) (arch.Sdiv, error) {
	return arch.SdivOf(rd, rn, rm)
}

func SmulhOf(
	rd string,
	rn string,
	rm string,
) (arch.Smulh, error) {
	return arch.SmulhOf(rd, rn, rm)
}

func StlrOf(
	rt, rn string,
	enc uint32,
) (arch.Stlr, error) {
	return arch.StlrOf(rt, rn, enc)
}

func StlxrOf(
	rs, rt, rn string,
	enc uint32,
) (arch.Stlxr, error) {
	return arch.StlxrOf(rs, rt, rn, enc)
}
func StpOf(rt, rt2, rn string, kind arch.MemKind, off int64, scale, enc uint32) (arch.Stp, error) {
	return arch.StpOf(rt, rt2, rn, kind, off, scale, enc)
}

func StrOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Str, error) {
	return arch.StrOf(rt, rn, kind, off, enc, rm, option, amt)
}

func StrbOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Strb, error) {
	return arch.StrbOf(rt, rn, kind, off, enc, rm, option, amt)
}

func StrhOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Strh, error) {
	return arch.StrhOf(rt, rn, kind, off, enc, rm, option, amt)
}

func SturOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Stur, error) {
	return arch.SturOf(rt, rn, kind, off, enc, rm, option, amt)
}

func SturbOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Sturb, error) {
	return arch.SturbOf(rt, rn, kind, off, enc, rm, option, amt)
}

func SturhOf(
	rt, rn string,
	kind arch.MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (arch.Sturh, error) {
	return arch.SturhOf(rt, rn, kind, off, enc, rm, option, amt)
}
func StxrbOf(rs, rt, rn string, enc uint32) (arch.Stxrb, error) {
	return arch.StxrbOf(rs, rt, rn, enc)
}

func SubExtOf(
	rdNum, rnNum, rmNum uint32,
	option string,
	imm3 uint32,
	isf bool,
) (arch.SubExt, error) {
	return arch.SubExtOf(rdNum, rnNum, rmNum, option, imm3, isf)
}
func SubImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (arch.SubImm, error) {
	return arch.SubImmOf(rdNum, rnNum, imm12, shift, isf)
}

func SubsExtOf(
	rdNum, rnNum, rmNum uint32,
	option string,
	imm3 uint32,
	isf bool,
) (arch.SubsExt, error) {
	return arch.SubsExtOf(rdNum, rnNum, rmNum, option, imm3, isf)
}
func SubsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (arch.SubsImm, error) {
	return arch.SubsImmOf(rdNum, rnNum, imm12, shift, isf)
}
func TbzOf(rt string, bit uint32, target int64, isTbnz bool) (arch.Tbz, error) {
	return arch.TbzOf(rt, bit, target, isTbnz)
}
func UdivOf(rd string, rn string, rm string) (arch.Udiv, error) {
	return arch.UdivOf(rd, rn, rm)
}
func UmulhOf(rd string, rn string, rm string) (arch.Umulh, error) {
	return arch.UmulhOf(rd, rn, rm)
}

func SubShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (Instr, error) {
	return arch.SubShiftOf(rd, rn, rm, imm6, shift, isf)
}
func SubsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (Instr, error) {
	return arch.SubsShiftOf(rd, rn, rm, imm6, shift, isf)
}
func UbfmOf(rd, rn string, immr, imms uint32, isf bool) (Instr, error) {
	return arch.UbfmOf(rd, rn, immr, imms, isf)
}
func SbfmOf(rd, rn string, immr, imms uint32, isf bool) (Instr, error) {
	return arch.SbfmOf(rd, rn, immr, imms, isf)
}
func regNums3(a, b, c string) (uint32, uint32, uint32, error) { return arch.RegNums3(a, b, c) }
func UbfmOfC(rd, rn string, immr, imms uint32, isf bool) (arch.Ubfm, error) {
	return arch.UbfmOfC(rd, rn, immr, imms, isf)
}

func encodeBitMasks(is64 bool, value uint64) (uint32, uint32, uint32, bool) {
	return arch.EncodeBitMasks(is64, value)
}
func OrrImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (Instr, error) {
	return arch.OrrImmOf(rd, rn, immr, imms, n, is64)
}
func AndsImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (Instr, error) {
	return arch.AndsImmOf(rd, rn, immr, imms, n, is64)
}

func AndsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (Instr, error) {
	return arch.AndsShiftOf(rd, rn, rm, imm6, shift, isf)
}
func OrrShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (Instr, error) {
	return arch.OrrShiftOf(rd, rn, rm, imm6, shift, isf)
}
func OrnShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (Instr, error) {
	return arch.OrnShiftOf(rd, rn, rm, imm6, shift, isf)
}

func MovzOf(rd string, imm16, hw uint32) (Instr, error) { return arch.MovzOf(rd, imm16, hw) }
func MovnOf(rd string, imm16, hw uint32) (Instr, error) { return arch.MovnOf(rd, imm16, hw) }
func armReg3(ops []vOp, name string) (string, string, string, error) {
	return arch.ArmReg3(ops, name)
}

func offBits(off int64, bits int) (uint32, error) { return arch.OffBits(off, bits) }

func invertCond(c string) string { return arch.InvertCond(c) }
