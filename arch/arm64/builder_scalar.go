package arm64

// String-operands constructors for the scalar families the assembler
// ctors build (the *Of convention of api.go; the typed Builder methods
// serve programmatic generation). imm-typed fields take int64.

// AdcOf — the Adc family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdcOf(rd string, rn string, rm string) Adc {
	return Adc{rd: rd, rn: rn, rm: rm}
}

// AdrOf — the Adr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdrOf(rd string, off int64) Adr {
	return Adr{rd: rd, off: off}
}

// AdrpOf — the Adrp family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdrpOf(rd string, off int64) Adrp {
	return Adrp{rd: rd, off: off}
}

// AsrRegOf — the AsrReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AsrRegOf(rd string, rn string, rm string) AsrReg {
	return AsrReg{rd: rd, rn: rn, rm: rm}
}

// BOf — the B family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BOf(target int64) B {
	return B{target: immNum(target)}
}

// BcondOf — the Bcond family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BcondOf(cond string, target int64) Bcond {
	return Bcond{cond: cond, target: immNum(target)}
}

// BlOf — the Bl family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlOf(target int64) Bl {
	return Bl{target: immNum(target)}
}

// BlrOf — the Blr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlrOf(rn string) Blr {
	return Blr{rn: rn}
}

// BrOf — the Br family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BrOf(rn string) Br {
	return Br{rn: rn}
}

// CbnzOf — the Cbnz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbnzOf(rt string, target int64) Cbnz {
	return Cbnz{rt: rt, target: immNum(target)}
}

// CbzOf — the Cbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbzOf(rt string, target int64) Cbz {
	return Cbz{rt: rt, target: immNum(target)}
}

// CcmpOf — the Ccmp family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CcmpOf(rn string, rm string, immVal uint32, cond string) Ccmp {
	return Ccmp{rn: rn, rm: rm, immVal: immVal, cond: cond}
}

// ClsOf — the Cls family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClsOf(rd string, rn string) Cls {
	return Cls{rd: rd, rn: rn}
}

// ClzOf — the Clz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClzOf(rd string, rn string) Clz {
	return Clz{rd: rd, rn: rn}
}

// CselOf — the Csel family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CselOf(rd string, rn string, rm string, cond string) Csel {
	return Csel{rd: rd, rn: rn, rm: rm, cond: cond}
}

// ExtrOf — the Extr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ExtrOf(rd string, rn string, rm string, lsb uint32) Extr {
	return Extr{rd: rd, rn: rn, rm: rm, lsb: lsb}
}

// LdrsbOf — the Ldrsb family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrsbOf(rt string, rn string, off int64) Ldrsb {
	return Ldrsb{rt: rt, rn: rn, off: off}
}

// LdrshOf — the Ldrsh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrshOf(rt string, rn string, off int64) Ldrsh {
	return Ldrsh{rt: rt, rn: rn, off: off}
}

// LslRegOf — the LslReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LslRegOf(rd string, rn string, rm string) LslReg {
	return LslReg{rd: rd, rn: rn, rm: rm}
}

// LsrRegOf — the LsrReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LsrRegOf(rd string, rn string, rm string) LsrReg {
	return LsrReg{rd: rd, rn: rn, rm: rm}
}

// MaddOf — the Madd family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MaddOf(rd string, rn string, rm string, ra string) Madd {
	return Madd{rd: rd, rn: rn, rm: rm, ra: ra}
}

// MovkOf — the Movk family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MovkOf(rd string, imm16 uint32, hw uint32) Movk {
	return Movk{rd: rd, imm16: imm16, hw: hw}
}

// MrsOf — the Mrs family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MrsOf(rd string, sysreg string) Mrs {
	return Mrs{rd: rd, sysreg: sysreg}
}

// MsrOf — the Msr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MsrOf(rt string, sysreg string) Msr {
	return Msr{rt: rt, sysreg: sysreg}
}

// PrfmOf — the Prfm family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func PrfmOf(rn string) Prfm {
	return Prfm{rn: rn}
}

// RbitOf — the Rbit family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RbitOf(rd string, rn string) Rbit {
	return Rbit{rd: rd, rn: rn}
}

// RetOf — the Ret family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RetOf(rn string) Ret {
	return Ret{rn: rn}
}

// RevOf — the Rev family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RevOf(rd string, rn string) Rev {
	return Rev{rd: rd, rn: rn}
}

// Rev16Of — the Rev16 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev16Of(rd string, rn string) Rev16 {
	return Rev16{rd: rd, rn: rn}
}

// Rev32Of — the Rev32 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev32Of(rd string, rn string) Rev32 {
	return Rev32{rd: rd, rn: rn}
}

// RorRegOf — the RorReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RorRegOf(rd string, rn string, rm string) RorReg {
	return RorReg{rd: rd, rn: rn, rm: rm}
}

// SdivOf — the Sdiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SdivOf(rd string, rn string, rm string) Sdiv {
	return Sdiv{rd: rd, rn: rn, rm: rm}
}

// SmulhOf — the Smulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SmulhOf(rd string, rn string, rm string) Smulh {
	return Smulh{rd: rd, rn: rn, rm: rm}
}

// TbzOf — the Tbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func TbzOf(rt string, bit uint32, target int64, isTbnz bool) Tbz {
	return Tbz{rt: rt, bit: bit, target: immNum(target), isTbnz: isTbnz}
}

// UdivOf — the Udiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UdivOf(rd string, rn string, rm string) Udiv {
	return Udiv{rd: rd, rn: rn, rm: rm}
}

// UmulhOf — the Umulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UmulhOf(rd string, rn string, rm string) Umulh {
	return Umulh{rd: rd, rn: rn, rm: rm}
}

// MemKind — the memory operand form (exported alias; the constants
// mirror the decoder's memKind).
type MemKind = memKind

// The memory operand forms.
const (
	MemImm      = memImm
	MemUnscaled = memUnscaled
	MemPost     = memPost
	MemPre      = memPre
	MemLiteral  = memLiteral
	MemRegOff   = memRegOff
)

// LdrOf — the Ldr family by the memory-operand base.
func LdrOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldr {
	return Ldr{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrbOf — the Ldrb family by the memory-operand base.
func LdrbOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldrb {
	return Ldrb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrhOf — the Ldrh family by the memory-operand base.
func LdrhOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldrh {
	return Ldrh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrOf — the Str family by the memory-operand base.
func StrOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Str {
	return Str{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrbOf — the Strb family by the memory-operand base.
func StrbOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Strb {
	return Strb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrhOf — the Strh family by the memory-operand base.
func StrhOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Strh {
	return Strh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurOf — the Ldur family by the memory-operand base.
func LdurOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldur {
	return Ldur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturOf — the Stur family by the memory-operand base.
func SturOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Stur {
	return Stur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurbOf — the Ldurb family by the memory-operand base.
func LdurbOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldurb {
	return Ldurb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurhOf — the Ldurh family by the memory-operand base.
func LdurhOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldurh {
	return Ldurh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturbOf — the Sturb family by the memory-operand base.
func SturbOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Sturb {
	return Sturb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturhOf — the Sturh family by the memory-operand base.
func SturhOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Sturh {
	return Sturh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrswOf — the Ldrsw family by the memory-operand base.
func LdrswOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Ldrsw {
	return Ldrsw{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrLitOf — the ldr literal form (ldr rt, =addr / label).
func LdrLitOf(rt string, tgt uint64, enc uint32) Ldr {
	return Ldr{lsBase: newLsBase(rt, "", memLiteral, 0, tgt, enc, "", "", 0)}
}

// LdpOf/StpOf — the load/store pair family.
func LdpOf(rt, rt2, rn string, kind MemKind, off int64, scale, enc uint32) Ldp {
	return Ldp{pairBase: newPairBase(rt, rt2, rn, kind, off, scale, enc)}
}

func StpOf(rt, rt2, rn string, kind MemKind, off int64, scale, enc uint32) Stp {
	return Stp{pairBase: newPairBase(rt, rt2, rn, kind, off, scale, enc)}
}

// AddExtOf — the AddExt family (register-extend operand).
func AddExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) AddExt {
	return AddExt{extBase: newExtBase(rdNum, rnNum, rmNum, option, imm3, isf)}
}

// AddsExtOf — the AddsExt family (register-extend operand).
func AddsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) AddsExt {
	return AddsExt{extBase: newExtBase(rdNum, rnNum, rmNum, option, imm3, isf)}
}

// SubExtOf — the SubExt family (register-extend operand).
func SubExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) SubExt {
	return SubExt{extBase: newExtBase(rdNum, rnNum, rmNum, option, imm3, isf)}
}

// SubsExtOf — the SubsExt family (register-extend operand).
func SubsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) SubsExt {
	return SubsExt{extBase: newExtBase(rdNum, rnNum, rmNum, option, imm3, isf)}
}

// AndImmOf — the AndImm family (logical immediate).
func AndImmOf(rd, rn string, immr, imms uint32, n, is64 bool) AndImm {
	return AndImm{logImm: newLogImm(rd, rn, immr, imms, n, is64)}
}

// EorImmOf — the EorImm family (logical immediate).
func EorImmOf(rd, rn string, immr, imms uint32, n, is64 bool) EorImm {
	return EorImm{logImm: newLogImm(rd, rn, immr, imms, n, is64)}
}

// AddsShiftOf — the AddsShift family (shifted register).
func AddsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AddsShift {
	return AddsShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// AddShiftOf — the AddShift family (shifted register).
func AddShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AddShift {
	return AddShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// AndShiftOf — the AndShift family (shifted register).
func AndShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AndShift {
	return AndShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// EorShiftOf — the EorShift family (shifted register).
func EorShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) EorShift {
	return EorShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// EonShiftOf — the EonShift family (shifted register).
func EonShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) EonShift {
	return EonShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// BicShiftOf — the BicShift family (shifted register).
func BicShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) BicShift {
	return BicShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// BicsShiftOf — the BicsShift family (shifted register).
func BicsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) BicsShift {
	return BicsShift{rd: rd, rn: rn, rm: rm, imm6: imm6, shift: shift, isf: isf}
}

// LdarOf — the Ldar family (atomic load/store).
func LdarOf(rt, rn string, enc uint32) Ldar {
	return Ldar{atomic: newAtomic(rt, rn), enc: enc}
}

// StlrOf — the Stlr family (atomic load/store).
func StlrOf(rt, rn string, enc uint32) Stlr {
	return Stlr{atomic: newAtomic(rt, rn), enc: enc}
}

// StxrbOf — the Stxrb family (exclusive store).
func StxrbOf(rs, rt, rn string, enc uint32) Stxrb {
	return Stxrb{excl: newExcl(rs, rt, rn), enc: enc}
}

// StlxrOf — the Stlxr family (exclusive store).
func StlxrOf(rs, rt, rn string, enc uint32) Stlxr {
	return Stlxr{excl: newExcl(rs, rt, rn), enc: enc}
}

// MsubOf — msub/msub-3-form: the Madd base with the operands.
func MsubOf(rd, rn, rm, ra string) Msub {
	return Msub{Madd: Madd{rd: rd, rn: rn, rm: rm, ra: ra}}
}

// CsincOf/CsinvOf/CsnegOf — the conditional-select family by the Csel base.
func CsincOf(c Csel) Csinc { return Csinc{Csel: c} }
func CsinvOf(c Csel) Csinv { return Csinv{Csel: c} }
func CsnegOf(c Csel) Csneg { return Csneg{Csel: c} }

// AddImmOf/AddsImmOf/SubImmOf/SubsImmOf — the immediate families
// (register numbers: 31 means sp/wsp).
func AddImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) AddImm {
	return AddImm{rdNum: rdNum, rnNum: rnNum, imm12: imm12, shift: shift, isf: isf}
}

func AddsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) AddsImm {
	return AddsImm{rdNum: rdNum, rnNum: rnNum, imm12: imm12, shift: shift, isf: isf}
}

func SubImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) SubImm {
	return SubImm{rdNum: rdNum, rnNum: rnNum, imm12: imm12, shift: shift, isf: isf}
}

func SubsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) SubsImm {
	return SubsImm{rdNum: rdNum, rnNum: rnNum, imm12: imm12, shift: shift, isf: isf}
}

// UbfmOfC — the concrete UBFM (api.go's UbfmOf returns Instr; this form
// returns the struct for the (Ubfm, bool) helpers).
func UbfmOfC(rd, rn string, immr, imms uint32, isf bool) Ubfm {
	return Ubfm{rd: rd, rn: rn, immr: immr, imms: imms, isf: isf}
}
