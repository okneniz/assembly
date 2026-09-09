package arm64

// String-operands constructors for the scalar families the assembler
// ctors build (the *Of convention of api.go; the typed Builder methods
// serve programmatic generation). imm-typed fields take int64.

// AdcOf — the Adc family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdcOf(rd string, rn string, rm string) Adc {
	return newAdc(base{}, rd, rn, rm)
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
	return newAsrReg(base{}, rd, rn, rm)
}

// BOf — the B family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BOf(off int64) B {
	return newB(base{}, immNum(off))
}

// BcondOf — the Bcond family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BcondOf(cond string, off int64) Bcond {
	return Bcond{cond: cond, off: immNum(off)}
}

// BlOf — the Bl family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlOf(off int64) Bl {
	return newBl(base{}, immNum(off))
}

// BlrOf — the Blr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlrOf(rn string) Blr {
	return newBlr(base{}, rn)
}

// BrOf — the Br family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BrOf(rn string) Br {
	return newBr(base{}, rn)
}

// CbnzOf — the Cbnz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbnzOf(rt string, off int64) Cbnz {
	return newCbnz(base{}, rt, immNum(off))
}

// CbzOf — the Cbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbzOf(rt string, off int64) Cbz {
	return newCbz(base{}, rt, immNum(off))
}

// CcmpOf — the Ccmp family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CcmpOf(rn string, rm string, immVal uint32, cond string) Ccmp {
	return newCcmp(base{}, rn, rm, immVal, cond)
}

// ClsOf — the Cls family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClsOf(rd string, rn string) Cls {
	return newCls(base{}, rd, rn)
}

// ClzOf — the Clz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClzOf(rd string, rn string) Clz {
	return newClz(base{}, rd, rn)
}

// CselOf — the Csel family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CselOf(rd string, rn string, rm string, cond string) Csel {
	return newCsel(base{}, rd, rn, rm, cond)
}

// ExtrOf — the Extr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ExtrOf(rd string, rn string, rm string, lsb uint32) Extr {
	return newExtr(base{}, rd, rn, rm, lsb)
}

// LdrsbOf — the Ldrsb family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrsbOf(rt string, rn string, off int64) Ldrsb {
	return newLdrsb(base{}, rt, rn, off)
}

// LdrshOf — the Ldrsh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrshOf(rt string, rn string, off int64) Ldrsh {
	return newLdrsh(base{}, rt, rn, off)
}

// LslRegOf — the LslReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LslRegOf(rd string, rn string, rm string) LslReg {
	return newLslReg(base{}, rd, rn, rm)
}

// LsrRegOf — the LsrReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LsrRegOf(rd string, rn string, rm string) LsrReg {
	return newLsrReg(base{}, rd, rn, rm)
}

// MaddOf — the Madd family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MaddOf(rd string, rn string, rm string, ra string) Madd {
	return newMadd(base{}, rd, rn, ra, rm)
}

// MovkOf — the Movk family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MovkOf(rd string, imm16 uint32, hw uint32) Movk {
	return newMovk(base{}, rd, imm16, hw)
}

// MrsOf — the Mrs family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MrsOf(rd string, sysreg string) Mrs {
	return newMrs(base{}, rd, sysreg)
}

// MsrOf — the Msr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MsrOf(rt string, sysreg string) Msr {
	return newMsr(base{}, rt, sysreg)
}

// PrfmOf — the Prfm family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func PrfmOf(rn string) Prfm {
	return newPrfm(base{}, rn)
}

// RbitOf — the Rbit family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RbitOf(rd string, rn string) Rbit {
	return newRbit(base{}, rd, rn)
}

// RetOf — the Ret family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RetOf(rn string) Ret {
	return newRet(base{}, rn)
}

// RevOf — the Rev family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RevOf(rd string, rn string) Rev {
	return newRev(base{}, rd, rn)
}

// Rev16Of — the Rev16 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev16Of(rd string, rn string) Rev16 {
	return newRev16(base{}, rd, rn)
}

// Rev32Of — the Rev32 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev32Of(rd string, rn string) Rev32 {
	return newRev32(base{}, rd, rn)
}

// RorRegOf — the RorReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RorRegOf(rd string, rn string, rm string) RorReg {
	return newRorReg(base{}, rd, rn, rm)
}

// SdivOf — the Sdiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SdivOf(rd string, rn string, rm string) Sdiv {
	return newSdiv(base{}, rd, rn, rm)
}

// SmulhOf — the Smulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SmulhOf(rd string, rn string, rm string) Smulh {
	return newSmulh(base{}, rd, rn, rm)
}

// TbzOf — the Tbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func TbzOf(rt string, bit uint32, off int64, isTbnz bool) Tbz {
	return Tbz{
		rt:     rt,
		bit:    bit,
		off:    immNum(off),
		isTbnz: isTbnz,
	}
}

// UdivOf — the Udiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UdivOf(rd string, rn string, rm string) Udiv {
	return newUdiv(base{}, rd, rn, rm)
}

// UmulhOf — the Umulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UmulhOf(rd string, rn string, rm string) Umulh {
	return newUmulh(base{}, rd, rn, rm)
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
func LdrbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldrb {
	return Ldrb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrhOf — the Ldrh family by the memory-operand base.
func LdrhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldrh {
	return Ldrh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrOf — the Str family by the memory-operand base.
func StrOf(rt, rn string, kind MemKind, off int64, enc uint32, rm, option string, amt uint32) Str {
	return Str{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrbOf — the Strb family by the memory-operand base.
func StrbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Strb {
	return Strb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// StrhOf — the Strh family by the memory-operand base.
func StrhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Strh {
	return Strh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurOf — the Ldur family by the memory-operand base.
func LdurOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldur {
	return Ldur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturOf — the Stur family by the memory-operand base.
func SturOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Stur {
	return Stur{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurbOf — the Ldurb family by the memory-operand base.
func LdurbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldurb {
	return Ldurb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdurhOf — the Ldurh family by the memory-operand base.
func LdurhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldurh {
	return Ldurh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturbOf — the Sturb family by the memory-operand base.
func SturbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Sturb {
	return Sturb{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// SturhOf — the Sturh family by the memory-operand base.
func SturhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Sturh {
	return Sturh{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrswOf — the Ldrsw family by the memory-operand base.
func LdrswOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) Ldrsw {
	return Ldrsw{lsBase: makeLSBase(rt, rn, kind, off, enc, rm, option, amt)}
}

// LdrLitOf — the ldr literal form (ldr rt, =addr / label; lit — the
// pc-relative byte offset of the literal).
func LdrLitOf(rt string, lit int64, enc uint32) Ldr {
	return Ldr{lsBase: newLsBase(rt, "", memLiteral, 0, lit, enc, "", "", 0)}
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
	return newAddExt(base{}, newExtBase(rdNum, rnNum, rmNum, option, imm3, isf))
}

// AddsExtOf — the AddsExt family (register-extend operand).
func AddsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) AddsExt {
	return newAddsExt(base{}, newExtBase(rdNum, rnNum, rmNum, option, imm3, isf))
}

// SubExtOf — the SubExt family (register-extend operand).
func SubExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) SubExt {
	return newSubExt(base{}, newExtBase(rdNum, rnNum, rmNum, option, imm3, isf))
}

// SubsExtOf — the SubsExt family (register-extend operand).
func SubsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) SubsExt {
	return newSubsExt(base{}, newExtBase(rdNum, rnNum, rmNum, option, imm3, isf))
}

// AndImmOf — the AndImm family (logical immediate).
func AndImmOf(rd, rn string, immr, imms uint32, n, is64 bool) AndImm {
	return newAndImm(base{}, newLogImm(rd, rn, immr, imms, n, is64))
}

// EorImmOf — the EorImm family (logical immediate).
func EorImmOf(rd, rn string, immr, imms uint32, n, is64 bool) EorImm {
	return EorImm{logImm: newLogImm(rd, rn, immr, imms, n, is64)}
}

// AddsShiftOf — the AddsShift family (shifted register).
func AddsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AddsShift {
	return newAddsShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// AddShiftOf — the AddShift family (shifted register).
func AddShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AddShift {
	return newAddShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// AndShiftOf — the AndShift family (shifted register).
func AndShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) AndShift {
	return newAndShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// EorShiftOf — the EorShift family (shifted register).
func EorShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) EorShift {
	return newEorShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// EonShiftOf — the EonShift family (shifted register).
func EonShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) EonShift {
	return newEonShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// BicShiftOf — the BicShift family (shifted register).
func BicShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) BicShift {
	return newBicShift(base{}, rd, rn, rm, imm6, shift, isf)
}

// BicsShiftOf — the BicsShift family (shifted register).
func BicsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) BicsShift {
	return newBicsShift(base{}, rd, rn, rm, imm6, shift, isf)
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
	return Msub{Madd: newMadd(base{}, rd, rn, ra, rm)}
}

// CsincOf/CsinvOf/CsnegOf — the conditional-select family by the Csel base.
func CsincOf(c Csel) Csinc { return Csinc{Csel: c} }
func CsinvOf(c Csel) Csinv { return Csinv{Csel: c} }
func CsnegOf(c Csel) Csneg { return Csneg{Csel: c} }

// AddImmOf/AddsImmOf/SubImmOf/SubsImmOf — the immediate families
// (register numbers: 31 means sp/wsp).
func AddImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) AddImm {
	return newAddImm(base{}, rdNum, rnNum, imm12, shift, isf)
}

func AddsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) AddsImm {
	return newAddsImm(base{}, rdNum, rnNum, imm12, shift, isf)
}

func SubImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) SubImm {
	return newSubImm(base{}, rdNum, rnNum, imm12, shift, isf)
}

func SubsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) SubsImm {
	return newSubsImm(base{}, rdNum, rnNum, imm12, shift, isf)
}

// UbfmOfC — the concrete UBFM (api.go's UbfmOf returns Instr; this form
// returns the struct for the (Ubfm, bool) helpers).
func UbfmOfC(rd, rn string, immr, imms uint32, isf bool) Ubfm {
	return newUbfm(base{}, rd, rn, immr, imms, isf)
}
