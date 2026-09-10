package arm64

// String-operands constructors for the scalar families the assembler
// ctors build (the *Of convention of api.go; the typed Builder methods
// serve programmatic generation). imm-typed fields take int64.

// AdcOf — the Adc family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdcOf(rd string, rn string, rm string) (Adc, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Adc{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Adc{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Adc{}, err
	}

	in, err := newAdc(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return Adc{}, err
	}

	return in, nil
}

// AdrOf — the Adr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdrOf(rd string, off int64) (Adr, error) {
	r, err := RegOf(rd)
	if err != nil {
		return Adr{}, err
	}

	return newAdr(base{}, r, off)
}

// AdrpOf — the Adrp family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AdrpOf(rd string, off int64) (Adrp, error) {
	r, err := RegOf(rd)
	if err != nil {
		return Adrp{}, err
	}

	return newAdrp(base{}, r, off)
}

// AsrRegOf — the AsrReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func AsrRegOf(rd string, rn string, rm string) (AsrReg, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return AsrReg{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return AsrReg{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return AsrReg{}, err
	}

	in, err := newAsrReg(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return AsrReg{}, err
	}

	return in, nil
}

// BOf — the B family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BOf(off int64) (B, error) {
	in, err := newB(base{}, immNum(off))
	if err != nil {
		return B{}, err
	}

	return in, nil
}

// BcondOf — the Bcond family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BcondOf(cond string, off int64) (Bcond, error) {
	return newBcond(base{}, cond, immNum(off)), nil
}

// BlOf — the Bl family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlOf(off int64) (Bl, error) {
	in, err := newBl(base{}, immNum(off))
	if err != nil {
		return Bl{}, err
	}

	return in, nil
}

// BlrOf — the Blr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BlrOf(rn string) (Blr, error) {
	r_rn, err := RegOf(rn)
	if err != nil {
		return Blr{}, err
	}

	in, err := newBlr(base{}, r_rn)
	if err != nil {
		return Blr{}, err
	}

	return in, nil
}

// BrOf — the Br family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func BrOf(rn string) (Br, error) {
	r_rn, err := RegOf(rn)
	if err != nil {
		return Br{}, err
	}

	in, err := newBr(base{}, r_rn)
	if err != nil {
		return Br{}, err
	}

	return in, nil
}

// CbnzOf — the Cbnz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbnzOf(rt string, off int64) (Cbnz, error) {
	r_rt, err := RegOf(rt)
	if err != nil {
		return Cbnz{}, err
	}

	in, err := newCbnz(base{}, r_rt, off)
	if err != nil {
		return Cbnz{}, err
	}

	return in, nil
}

// CbzOf — the Cbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CbzOf(rt string, off int64) (Cbz, error) {
	r_rt, err := RegOf(rt)
	if err != nil {
		return Cbz{}, err
	}

	in, err := newCbz(base{}, r_rt, off)
	if err != nil {
		return Cbz{}, err
	}

	return in, nil
}

// CcmpOf — the Ccmp family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CcmpOf(rn string, rm string, immVal uint32, cond string) (Ccmp, error) {
	r_rn, err := RegOf(rn)
	if err != nil {
		return Ccmp{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Ccmp{}, err
	}

	in, err := newCcmp(base{}, r_rn, r_rm, immVal, cond)
	if err != nil {
		return Ccmp{}, err
	}

	return in, nil
}

// ClsOf — the Cls family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClsOf(rd string, rn string) (Cls, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Cls{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Cls{}, err
	}

	in, err := newCls(base{}, r_rd, r_rn)
	if err != nil {
		return Cls{}, err
	}

	return in, nil
}

// ClzOf — the Clz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ClzOf(rd string, rn string) (Clz, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Clz{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Clz{}, err
	}

	in, err := newClz(base{}, r_rd, r_rn)
	if err != nil {
		return Clz{}, err
	}

	return in, nil
}

// CselOf — the Csel family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func CselOf(rd string, rn string, rm string, cond string) (Csel, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Csel{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Csel{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Csel{}, err
	}

	in, err := newCsel(base{}, r_rd, r_rn, r_rm, cond)
	if err != nil {
		return Csel{}, err
	}

	return in, nil
}

// ExtrOf — the Extr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func ExtrOf(rd string, rn string, rm string, lsb uint32) (Extr, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Extr{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Extr{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Extr{}, err
	}

	in, err := newExtr(base{}, r_rd, r_rn, r_rm, imm6Of(lsb))
	if err != nil {
		return Extr{}, err
	}

	return in, nil
}

// LdrsbOf — the Ldrsb family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrsbOf(rt string, rn string, off int64) (Ldrsb, error) {
	r_rt, err := RegOf(rt)
	if err != nil {
		return Ldrsb{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Ldrsb{}, err
	}

	in, err := newLdrsb(base{}, r_rt, r_rn, Off(off))
	if err != nil {
		return Ldrsb{}, err
	}

	return in, nil
}

// LdrshOf — the Ldrsh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LdrshOf(rt string, rn string, off int64) (Ldrsh, error) {
	r_rt, err := RegOf(rt)
	if err != nil {
		return Ldrsh{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Ldrsh{}, err
	}

	in, err := newLdrsh(base{}, r_rt, r_rn, Off(off))
	if err != nil {
		return Ldrsh{}, err
	}

	return in, nil
}

// LslRegOf — the LslReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LslRegOf(rd string, rn string, rm string) (LslReg, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return LslReg{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return LslReg{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return LslReg{}, err
	}

	in, err := newLslReg(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return LslReg{}, err
	}

	return in, nil
}

// LsrRegOf — the LsrReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func LsrRegOf(rd string, rn string, rm string) (LsrReg, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return LsrReg{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return LsrReg{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return LsrReg{}, err
	}

	in, err := newLsrReg(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return LsrReg{}, err
	}

	return in, nil
}

// MaddOf — the Madd family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MaddOf(rd string, rn string, rm string, ra string) (Madd, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Madd{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Madd{}, err
	}

	r_ra, err := RegOf(ra)
	if err != nil {
		return Madd{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Madd{}, err
	}

	in, err := newMadd(base{}, r_rd, r_rn, r_rm, r_ra)
	if err != nil {
		return Madd{}, err
	}

	return in, nil
}

// MovkOf — the Movk family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MovkOf(rd string, imm16 uint32, hw uint32) (Movk, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Movk{}, err
	}

	in, err := newMovk(base{}, r_rd, imm16Of(imm16), Hw(hw))
	if err != nil {
		return Movk{}, err
	}

	return in, nil
}

// MrsOf — the Mrs family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MrsOf(rd string, sysreg string) (Mrs, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Mrs{}, err
	}

	in, err := newMrs(base{}, r_rd, sysreg)
	if err != nil {
		return Mrs{}, err
	}

	return in, nil
}

// MsrOf — the Msr family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func MsrOf(rt string, sysreg string) (Msr, error) {
	r_sysreg, err := RegOf(sysreg)
	if err != nil {
		return Msr{}, err
	}

	in, err := newMsr(base{}, rt, r_sysreg)
	if err != nil {
		return Msr{}, err
	}

	return in, nil
}

// PrfmOf — the Prfm family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func PrfmOf(rn string) (Prfm, error) {
	r_rn, err := RegOf(rn)
	if err != nil {
		return Prfm{}, err
	}

	in, err := newPrfm(base{}, r_rn)
	if err != nil {
		return Prfm{}, err
	}

	return in, nil
}

// RbitOf — the Rbit family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RbitOf(rd string, rn string) (Rbit, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Rbit{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Rbit{}, err
	}

	in, err := newRbit(base{}, r_rd, r_rn)
	if err != nil {
		return Rbit{}, err
	}

	return in, nil
}

// RetOf — the Ret family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RetOf(rn string) (Ret, error) {
	r_rn, err := RegOf(rn)
	if err != nil {
		return Ret{}, err
	}

	in, err := newRet(base{}, r_rn)
	if err != nil {
		return Ret{}, err
	}

	return in, nil
}

// RevOf — the Rev family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RevOf(rd string, rn string) (Rev, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Rev{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Rev{}, err
	}

	in, err := newRev(base{}, r_rd, r_rn)
	if err != nil {
		return Rev{}, err
	}

	return in, nil
}

// Rev16Of — the Rev16 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev16Of(rd string, rn string) (Rev16, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Rev16{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Rev16{}, err
	}

	in, err := newRev16(base{}, r_rd, r_rn)
	if err != nil {
		return Rev16{}, err
	}

	return in, nil
}

// Rev32Of — the Rev32 family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func Rev32Of(rd string, rn string) (Rev32, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Rev32{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Rev32{}, err
	}

	in, err := newRev32(base{}, r_rd, r_rn)
	if err != nil {
		return Rev32{}, err
	}

	return in, nil
}

// RorRegOf — the RorReg family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func RorRegOf(rd string, rn string, rm string) (RorReg, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return RorReg{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return RorReg{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return RorReg{}, err
	}

	in, err := newRorReg(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return RorReg{}, err
	}

	return in, nil
}

// SdivOf — the Sdiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SdivOf(rd string, rn string, rm string) (Sdiv, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Sdiv{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Sdiv{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Sdiv{}, err
	}

	in, err := newSdiv(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return Sdiv{}, err
	}

	return in, nil
}

// SmulhOf — the Smulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func SmulhOf(rd string, rn string, rm string) (Smulh, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Smulh{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Smulh{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Smulh{}, err
	}

	in, err := newSmulh(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return Smulh{}, err
	}

	return in, nil
}

// TbzOf — the Tbz family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func TbzOf(rt string, bit uint32, off int64, isTbnz bool) (Tbz, error) {
	return Tbz{
		rt:     rt,
		bit:    bit,
		off:    immNum(off),
		isTbnz: isTbnz,
	}, nil
}

// UdivOf — the Udiv family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UdivOf(rd string, rn string, rm string) (Udiv, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Udiv{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Udiv{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Udiv{}, err
	}

	in, err := newUdiv(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return Udiv{}, err
	}

	return in, nil
}

// UmulhOf — the Umulh family constructor (string-operands form for the
// assembler ctors, which build outside arch).
func UmulhOf(rd string, rn string, rm string) (Umulh, error) {
	r_rd, err := RegOf(rd)
	if err != nil {
		return Umulh{}, err
	}

	r_rn, err := RegOf(rn)
	if err != nil {
		return Umulh{}, err
	}

	r_rm, err := RegOf(rm)
	if err != nil {
		return Umulh{}, err
	}

	in, err := newUmulh(base{}, r_rd, r_rn, r_rm)
	if err != nil {
		return Umulh{}, err
	}

	return in, nil
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
func LdrOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldr, error) {
	return newLdrBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdrbOf — the Ldrb family by the memory-operand base.
func LdrbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldrb, error) {
	return newLdrbBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdrhOf — the Ldrh family by the memory-operand base.
func LdrhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldrh, error) {
	return newLdrhBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// StrOf — the Str family by the memory-operand base.
func StrOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Str, error) {
	return newStrBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// StrbOf — the Strb family by the memory-operand base.
func StrbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Strb, error) {
	return newStrbBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// StrhOf — the Strh family by the memory-operand base.
func StrhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Strh, error) {
	return newStrhBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdurOf — the Ldur family by the memory-operand base.
func LdurOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldur, error) {
	return newLdurBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// SturOf — the Stur family by the memory-operand base.
func SturOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Stur, error) {
	return newSturBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdurbOf — the Ldurb family by the memory-operand base.
func LdurbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldurb, error) {
	return newLdurbBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdurhOf — the Ldurh family by the memory-operand base.
func LdurhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldurh, error) {
	return newLdurhBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// SturbOf — the Sturb family by the memory-operand base.
func SturbOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Sturb, error) {
	return newSturbBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// SturhOf — the Sturh family by the memory-operand base.
func SturhOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Sturh, error) {
	return newSturhBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdrswOf — the Ldrsw family by the memory-operand base.
func LdrswOf(
	rt, rn string,
	kind MemKind,
	off int64,
	enc uint32,
	rm, option string,
	amt uint32,
) (Ldrsw, error) {
	return newLdrswBase(base{}, makeLSBase(rt, rn, kind, off, enc, rm, option, amt)), nil
}

// LdrLitOf — the ldr literal form (ldr rt, =addr / label; lit — the
// pc-relative byte offset of the literal).
func LdrLitOf(rt string, lit int64, enc uint32) (Ldr, error) {
	return newLdrBase(base{}, newLsBase(rt, "", memLiteral, 0, lit, enc, "", "", 0)), nil
}

// LdpOf/StpOf — the load/store pair family.
func LdpOf(rt, rt2, rn string, kind MemKind, off int64, scale, enc uint32) (Ldp, error) {
	return newLdpBase(base{}, newPairBase(rt, rt2, rn, kind, off, scale, enc)), nil
}

func StpOf(rt, rt2, rn string, kind MemKind, off int64, scale, enc uint32) (Stp, error) {
	return newStpBase(base{}, newPairBase(rt, rt2, rn, kind, off, scale, enc)), nil
}

// AddExtOf — the AddExt family (register-extend operand).
func AddExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) (AddExt, error) {
	return newAddExtBase(
		base{},
		newExtBase(rdNum, rnNum, rmNum, option, imm3, isf),
	), nil
}

// AddsExtOf — the AddsExt family (register-extend operand).
func AddsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) (AddsExt, error) {
	return newAddsExtBase(
		base{},
		newExtBase(rdNum, rnNum, rmNum, option, imm3, isf),
	), nil
}

// SubExtOf — the SubExt family (register-extend operand).
func SubExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) (SubExt, error) {
	return newSubExtBase(
		base{},
		newExtBase(rdNum, rnNum, rmNum, option, imm3, isf),
	), nil
}

// SubsExtOf — the SubsExt family (register-extend operand).
func SubsExtOf(rdNum, rnNum, rmNum uint32, option string, imm3 uint32, isf bool) (SubsExt, error) {
	return newSubsExtBase(
		base{},
		newExtBase(rdNum, rnNum, rmNum, option, imm3, isf),
	), nil
}

// AndImmOf — the AndImm family (logical immediate).
func AndImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (AndImm, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return AndImm{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return AndImm{}, err
	}

	in, err := newAndImm(base{}, r1, r2, decodeBitMasks(n, immr, imms, is64))
	if err != nil {
		return AndImm{}, err
	}

	return in, nil
}

// EorImmOf — the EorImm family (logical immediate).
func EorImmOf(rd, rn string, immr, imms uint32, n, is64 bool) (EorImm, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return EorImm{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return EorImm{}, err
	}

	in, err := newEorImm(base{}, r1, r2, decodeBitMasks(n, immr, imms, is64))
	if err != nil {
		return EorImm{}, err
	}

	return in, nil
}

// AddsShiftOf — the AddsShift family (shifted register).
func AddsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (AddsShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return AddsShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return AddsShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return AddsShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return AddsShift{}, err
	}

	in, err := newAddsShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return AddsShift{}, err
	}

	return in, nil
}

// AddShiftOf — the AddShift family (shifted register).
func AddShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (AddShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return AddShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return AddShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return AddShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return AddShift{}, err
	}

	in, err := newAddShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return AddShift{}, err
	}

	return in, nil
}

// AndShiftOf — the AndShift family (shifted register).
func AndShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (AndShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return AndShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return AndShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return AndShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return AndShift{}, err
	}

	in, err := newAndShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return AndShift{}, err
	}

	return in, nil
}

// EorShiftOf — the EorShift family (shifted register).
func EorShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (EorShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return EorShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return EorShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return EorShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return EorShift{}, err
	}

	in, err := newEorShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return EorShift{}, err
	}

	return in, nil
}

// EonShiftOf — the EonShift family (shifted register).
func EonShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (EonShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return EonShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return EonShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return EonShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return EonShift{}, err
	}

	in, err := newEonShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return EonShift{}, err
	}

	return in, nil
}

// BicShiftOf — the BicShift family (shifted register).
func BicShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (BicShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return BicShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return BicShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return BicShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return BicShift{}, err
	}

	in, err := newBicShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return BicShift{}, err
	}

	return in, nil
}

// BicsShiftOf — the BicsShift family (shifted register).
func BicsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) (BicsShift, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return BicsShift{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return BicsShift{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return BicsShift{}, err
	}

	sh, err := shiftOfName(shift)
	if err != nil {
		return BicsShift{}, err
	}

	in, err := newBicsShift(base{}, r1, r2, r3, imm6Of(imm6), sh)
	if err != nil {
		return BicsShift{}, err
	}

	return in, nil
}

// LdarOf — the Ldar family (atomic load/store).
func LdarOf(rt, rn string, enc uint32) (Ldar, error) {
	return Ldar{atomic: newAtomic(rt, rn), enc: enc}, nil
}

// StlrOf — the Stlr family (atomic load/store).
func StlrOf(rt, rn string, enc uint32) (Stlr, error) {
	return Stlr{atomic: newAtomic(rt, rn), enc: enc}, nil
}

// StxrbOf — the Stxrb family (exclusive store).
func StxrbOf(rs, rt, rn string, enc uint32) (Stxrb, error) {
	return Stxrb{excl: newExcl(rs, rt, rn), enc: enc}, nil
}

// StlxrOf — the Stlxr family (exclusive store).
func StlxrOf(rs, rt, rn string, enc uint32) (Stlxr, error) {
	return Stlxr{excl: newExcl(rs, rt, rn), enc: enc}, nil
}

// MsubOf — msub/msub-3-form: the Madd base with the operands.
func MsubOf(rd, rn, rm, ra string) (Msub, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return Msub{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return Msub{}, err
	}

	r3, err := RegOf(rm)
	if err != nil {
		return Msub{}, err
	}

	r4, err := RegOf(ra)
	if err != nil {
		return Msub{}, err
	}

	m, err := newMadd(base{}, r1, r2, r4, r3)
	if err != nil {
		return Msub{}, err
	}

	return Msub{Madd: m}, nil
}

// CsincOf/CsinvOf/CsnegOf — the conditional-select family by the Csel base.
func CsincOf(c Csel) (Csinc, error) { return Csinc{Csel: c}, nil }
func CsinvOf(c Csel) (Csinv, error) { return Csinv{Csel: c}, nil }
func CsnegOf(c Csel) (Csneg, error) { return Csneg{Csel: c}, nil }

// AddImmOf/AddsImmOf/SubImmOf/SubsImmOf — the immediate families
// (register numbers: 31 means sp/wsp).
func AddImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (AddImm, error) {
	return newAddImm(base{}, numReg(rdNum, isf), numReg(rnNum, isf), imm12Of(imm12), sh12Of(shift))
}

func AddsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (AddsImm, error) {
	return newAddsImm(base{}, numReg(rdNum, isf), numReg(rnNum, isf), imm12Of(imm12), sh12Of(shift))
}

func SubImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (SubImm, error) {
	return newSubImm(base{}, numReg(rdNum, isf), numReg(rnNum, isf), imm12Of(imm12), sh12Of(shift))
}

func SubsImmOf(rdNum, rnNum, imm12 uint32, shift bool, isf bool) (SubsImm, error) {
	return newSubsImm(base{}, numReg(rdNum, isf), numReg(rnNum, isf), imm12Of(imm12), sh12Of(shift))
}

// UbfmOfC — the concrete UBFM (api.go's UbfmOf returns Instr; this form
// returns the struct for the (Ubfm, bool) helpers).
func UbfmOfC(rd, rn string, immr, imms uint32, isf bool) (Ubfm, error) {
	r1, err := RegOf(rd)
	if err != nil {
		return Ubfm{}, err
	}

	r2, err := RegOf(rn)
	if err != nil {
		return Ubfm{}, err
	}

	return newUbfm(base{}, r1, r2, immr, imms)
}
