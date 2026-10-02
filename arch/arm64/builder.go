package arm64

import "fmt"

// Builder — the construction vocabulary of the package: instruction
// constructors named after the mnemonic they build and operand role
// constructors. Stateless by design — a namespace, not a factory.
type Builder struct{}

// New — the Builder.
func New() Builder {
	return Builder{}
}

func (Builder) Abs(rd, rn VReg, arr string) (Instr, error) {
	return newAbs(base{}, rd, rn, arr)
}

func (Builder) Adc(rd, rn, rm Reg) (Instr, error) {
	return newAdc(base{}, rd, rn, rm)
}

func (Builder) Add(rd, rn, rm VReg, arr string) (Instr, error) {
	return newAdd(base{}, rd, rn, rm, arr)
}

func (Builder) AddExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newAddExt(base{}, rd, rn, rm, ext, imm3)
}

func (Builder) AddImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newAddImm(base{}, rd, rn, imm, sh)
}

func (Builder) AddShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAddShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Addp(rd, rn, rm VReg, arr string) (Instr, error) {
	return newAddp(base{}, rd, rn, rm, arr)
}

func (Builder) AddsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newAddsExt(base{}, rd, rn, rm, ext, imm3)
}

func (Builder) AddsImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newAddsImm(base{}, rd, rn, imm, sh)
}

func (Builder) AddsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAddsShift(base{}, rd, rn, rm, imm, sh)
}

// Adr — adr rd, #off: off — the signed byte offset from the address
// of the instruction (the imm21 form, -0x100000..0xfffff). rd — only x
// registers (register 31 reads as zr).
func (Builder) Adr(rd Reg, off int64) (Instr, error) {
	return newAdr(base{}, rd, off)
}

// Adrp — adrp rd, #off: off — the signed imm21 count of 4KB pages
// from the page of the instruction (-0x100000..0xfffff pages). The
// absolute-page annotation of the decoded form needs the instruction
// address and stays zero here. rd — only x registers (register 31 reads
// as zr).
func (Builder) Adrp(rd Reg, off int64) (Instr, error) {
	return newAdrp(base{}, rd, off)
}

func (Builder) Aese(rd, rn VReg) (Instr, error) {
	return newAese(base{}, rd, rn)
}

func (Builder) Aesmc(rd, rn VReg) (Instr, error) {
	return newAesmc(base{}, rd, rn)
}

func (Builder) And(rd, rn, rm VReg, arr string) (Instr, error) {
	return newAnd(base{}, rd, rn, rm, arr)
}

func (Builder) AndImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newAndImm(base{}, rd, rn, imm)
}

func (Builder) AndShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAndShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) AndsImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newAndsImm(base{}, rd, rn, imm)
}

func (Builder) AndsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newAndsShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) AsrReg(rd, rn, rm Reg) (Instr, error) {
	return newAsrReg(base{}, rd, rn, rm)
}

// B — b off (the pc-relative byte offset; the absolute target is off +
// the instruction address).
func (Builder) B(off int64) Instr {
	in, err := newB(base{}, immNum(off))
	if err != nil {
		panic(err) // a plain offset cannot fail
	}

	return in
}

// Bcond — b.cond off: off — the pc-relative byte offset of the branch
// destination (the ±1MB imm19 range is checked at encode time; the
// absolute target is off + the instruction address); cond — the
// standard condition names (eq/ne/.../nv, see condNum).
func (Builder) Bcond(cond string, off int64) (Instr, error) {
	if _, err := condNum(cond); err != nil {
		return nil, fmt.Errorf("arm64.NewBcond: operand cond: %w", err)
	}

	return Bcond{cond: cond, off: immNum(off)}, nil
}

func (Builder) Bfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newBfm(base{}, rd, rn, immr, imms)
}

func (Builder) Bic(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBic(base{}, rd, rn, rm, arr)
}

func (Builder) BicShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newBicShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) BicsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newBicsShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Bif(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBif(base{}, rd, rn, rm, arr)
}

func (Builder) Bit(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBit(base{}, rd, rn, rm, arr)
}

// Bl — bl off: off — the pc-relative byte offset of the call
// destination (the ±128MB imm26 range is checked at encode time; the
// absolute target is off + the instruction address).
func (Builder) Bl(off int64) Instr {
	in, err := newBl(base{}, immNum(off))
	if err != nil {
		panic(err) // a plain offset cannot fail
	}

	return in
}

func (Builder) Blr(rn Reg) (Instr, error) {
	return newBlr(base{}, rn)
}

func (Builder) Br(rn Reg) (Instr, error) {
	return newBr(base{}, rn)
}

// Brk — brk #imm16.
func (Builder) Brk(imm Imm16) Instr {
	return sysImm{
		name:  "brk",
		imm16: imm.v,
		enc:   0xD4200000,
		shift: 5,
	}
}

func (Builder) Bsl(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBsl(base{}, rd, rn, rm, arr)
}

func (Builder) Cbnz(rt Reg, off int64) (Instr, error) {
	return newCbnz(base{}, rt, off)
}

func (Builder) Cbz(rt Reg, off int64) (Instr, error) {
	return newCbz(base{}, rt, off)
}

func (Builder) Ccmp(rn, rm Reg, nzcv uint32, cond string) (Instr, error) {
	return newCcmp(base{}, rn, rm, nzcv, cond)
}

func (Builder) Cls(rd, rn Reg) (Instr, error) {
	return newCls(base{}, rd, rn)
}

func (Builder) Clz(rd, rn Reg) (Instr, error) {
	return newClz(base{}, rd, rn)
}

func (Builder) Cmeq(rd, rn, rm VReg, arr string) (Instr, error) {
	return newCmeq(base{}, rd, rn, rm, arr)
}

func (Builder) Cmge(rd, rn, rm VReg, arr string) (Instr, error) {
	return newCmge(base{}, rd, rn, rm, arr)
}

func (Builder) Cmtst(rd, rn, rm VReg, arr string) (Instr, error) {
	return newCmtst(base{}, rd, rn, rm, arr)
}

func (Builder) Cnt(rd, rn VReg, arr string) (Instr, error) {
	return newCnt(base{}, rd, rn, arr)
}

func (Builder) Csel(rd, rn, rm Reg, cond string) (Instr, error) {
	return newCsel(base{}, rd, rn, rm, cond)
}

// Csinc — csinc rd, rn, rm, cond (cset/cinc pseudos); the operand
// constraints are those of Csel.
func (Builder) Csinc(rd, rn, rm Reg, cond string) (Instr, error) {
	base, err := New().Csel(rd, rn, rm, cond)
	if err != nil {
		return nil, err
	}

	c, ok := base.(Csel)
	if !ok {
		return nil, fmt.Errorf("csinc: internal: want Csel, got %T", base)
	}

	return Csinc{Csel: c}, nil
}

// Csinv — csinv rd, rn, rm, cond (csetm/cinv pseudos); the operand
// constraints are those of Csel.
func (Builder) Csinv(rd, rn, rm Reg, cond string) (Instr, error) {
	base, err := New().Csel(rd, rn, rm, cond)
	if err != nil {
		return nil, err
	}

	c, ok := base.(Csel)
	if !ok {
		return nil, fmt.Errorf("csinv: internal: want Csel, got %T", base)
	}

	return Csinv{Csel: c}, nil
}

// Csneg — csneg rd, rn, rm, cond (cneg pseudo); the operand
// constraints are those of Csel.
func (Builder) Csneg(rd, rn, rm Reg, cond string) (Instr, error) {
	base, err := New().Csel(rd, rn, rm, cond)
	if err != nil {
		return nil, err
	}

	c, ok := base.(Csel)
	if !ok {
		return nil, fmt.Errorf("csneg: internal: want Csel, got %T", base)
	}

	return Csneg{Csel: c}, nil
}

// Dmb - dmb domain (the memory barrier of the domain).
func (Builder) Dmb(domain BarrierDomain) (Instr, error) {
	return newBarrier(base{}, "dmb", domain)
}

// Dsb - dsb domain (the barrier completes before anything continues).
func (Builder) Dsb(domain BarrierDomain) (Instr, error) {
	return newBarrier(base{}, "dsb", domain)
}

func (Builder) Dup(vd VReg, wn Reg, arr string) (Instr, error) {
	return newDup(base{}, vd, wn, arr)
}

func (Builder) DupElem(rd, rn VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewDupElem: %w", err)
	}

	return newDupElem(base{}, q, size, idx, rd, rn)
}

func (Builder) DupScalar(rd, rn VReg, elem string, idx uint32) (Instr, error) {
	size, err := elemSize("DupScalar", elem)
	if err != nil {
		return nil, err
	}

	return newDupScalar(base{}, size, idx, rd, rn)
}

func (Builder) EonShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newEonShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Eor(rd, rn, rm VReg, arr string) (Instr, error) {
	return newEor(base{}, rd, rn, rm, arr)
}

func (Builder) EorImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newEorImm(base{}, rd, rn, imm)
}

func (Builder) EorShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newEorShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Extr(rd, rn, rm Reg, lsb Imm6) (Instr, error) {
	return newExtr(base{}, rd, rn, rm, lsb)
}

func (Builder) Fadd(rd, rn, rm FReg) (Instr, error) {
	return newFadd(base{}, rd, rn, rm)
}

// FcmlaElem - the Builder entry: .4h/.8h/.2s/.4s/.2d lanes.
func (Builder) FcmlaElem(rd, rn, rm VReg, arr string, idx, rot uint32) (Instr, error) {
	// only the three complex arrangements exist: the .2s and .2d words
	// are unallocated (clang refuses, the word decodes as unknown)
	switch arr {
	case "4h", "8h", "4s":
	default:
		return nil, fmt.Errorf(
			"arm64.NewFcmlaElem: arrangement %q is not one of [4h 8h 4s]", arr,
		)
	}

	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewFcmlaElem: %w", err)
	}

	return newFcmlaElem(base{}, q, size, idx, rot, rd, rn, rm)
}

func (Builder) Fcmp(rn, rm FReg) (Instr, error) {
	return newFcmp(base{}, rn, rm, true)
}

// FcmpZero — fcmp fn, #0.0 (the immediate form: Rm=0).
func (Builder) FcmpZero(rn FReg) (Instr, error) {
	return newFcmp(base{}, rn, FReg{}, false)
}

func (Builder) Fcvt(rd, rn FReg) (Instr, error) {
	return newFcvt(base{}, rd, rn)
}

func (Builder) Fcvtzs(rd Reg, rn FReg) (Instr, error) {
	return newFcvtzs(base{}, rd, rn)
}

func (Builder) Fcvtzu(rd Reg, rn FReg) (Instr, error) {
	return newFcvtzu(base{}, rd, rn)
}

func (Builder) Fdiv(rd, rn, rm FReg) (Instr, error) {
	return newFdiv(base{}, rd, rn, rm)
}

func (Builder) Fmadd(rd, rn, rm, ra FReg) (Instr, error) {
	return newFmadd(base{}, rd, rn, rm, ra)
}

func (Builder) Fmax(rd, rn, rm FReg) (Instr, error) {
	return newFmax(base{}, rd, rn, rm)
}

func (Builder) Fmin(rd, rn, rm FReg) (Instr, error) {
	return newFmin(base{}, rd, rn, rm)
}

// FmlaElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmlaElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmlaElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmlaElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmlaElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmlaElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

// FmlsElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmlsElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmlsElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmlsElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmlsElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmlsElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

func (Builder) Fmov(rd, rn FReg) (Instr, error) {
	return newFmov(base{}, rd, rn)
}

func (Builder) FmovFromGpr(rd FReg, rn Reg) (Instr, error) {
	return newFmovFromGpr(base{}, rd, rn)
}

func (Builder) FmovImm(rd FReg, val float64) (Instr, error) {
	return newFmovImm(base{}, rd, val, fmt.Sprintf("%.8f", val))
}

func (Builder) FmovToGpr(rd Reg, rn FReg) (Instr, error) {
	return newFmovToGpr(base{}, rd, rn)
}

func (Builder) Fmul(rd, rn, rm FReg) (Instr, error) {
	return newFmul(base{}, rd, rn, rm)
}

// FmulElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmulElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmulElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmulElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmulElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmulElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

// FmulxElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmulxElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmulxElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmulxElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmulxElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmulxElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

func (Builder) Fneg(rd, rn FReg) (Instr, error) {
	return newFneg(base{}, rd, rn)
}

func (Builder) Fnmsub(rd, rn, rm, ra FReg) (Instr, error) {
	return newFnmsub(base{}, rd, rn, rm, ra)
}

func (Builder) Fsub(rd, rn, rm FReg) (Instr, error) {
	return newFsub(base{}, rd, rn, rm)
}

// Imm12 — validated value; error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	return newImm12(v)
}

// Imm16 — validated value; error when out of range.
func (Builder) Imm16(v int64) (Imm16, error) {
	return newImm16(v)
}

// Imm6 — validated value; error when out of range.
func (Builder) Imm6(v int64) (Imm6, error) {
	return newImm6(v)
}

func (Builder) Ins(vd VReg, idx uint32, wn Reg, elem string) (Instr, error) {
	size, err := elemSize("Ins", elem)
	if err != nil {
		return nil, err
	}

	return newIns(base{}, size, idx, vd, wn)
}

func (Builder) InsElem(rd, rn VReg, elem string, idx, srcIdx uint32) (Instr, error) {
	size, err := elemSize("InsElem", elem)
	if err != nil {
		return nil, err
	}

	return newInsElem(base{}, size, idx, srcIdx, rd, rn)
}

// Isb - isb (the instruction barrier; the full-system domain is its
// only form, so there is nothing to choose).
func (Builder) Isb() (Instr, error) {
	return newBarrier(base{}, "isb", Sy)
}

func (Builder) Ldar(rt, rn Reg) (Instr, error) {
	return newLdar(base{}, rt, rn)
}

func (Builder) Ldarb(rt, rn Reg) (Instr, error) {
	return newLdarb(base{}, rt, rn)
}

// Ldaxr — ldaxr rt, [rn]: rt — x/w register (register 31 reads as
// zr), rn — x register or SP (register 31 in the base reads as sp).
func (Builder) Ldaxr(rt, rn Reg) (Instr, error) {
	return newLdaxr(base{}, rt, rn)
}

// Ldaxrb — ldaxrb rt, [rn]: byte access, rt — w register only
// (register 31 reads as wzr), rn — x register or SP (register 31 in the
// base reads as sp).
func (Builder) Ldaxrb(rt, rn Reg) (Instr, error) {
	return newLdaxrb(base{}, rt, rn)
}

func (Builder) Ldp(rt, rt2, rn Reg, off Off) (Instr, error) {
	return newLdp(base{}, rt, rt2, rn, off)
}

func (Builder) Ldpsw(rt, rt2, rn Reg, off Off) (Instr, error) {
	return newLdpsw(base{}, rt, rt2, rn, off)
}

func (Builder) Ldr(rt, rn Reg, off Off) (Instr, error) {
	return newLdr(base{}, rt, rn, off)
}

// LdrF — ldr st|dt, [xn, #off] (the FP/SIMD register form).
func (Builder) LdrF(rt FReg, rn Reg, off Off) (Instr, error) {
	err := requireClass(
		rn,
		"LdrF",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX, classSP,
	)
	if err != nil {
		return Ldr{}, err
	}

	enc, scale := ldrFDEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldrFSEnc, 2
	}

	if err := requireOff("LdrF", off, scale); err != nil {
		return Ldr{}, err
	}

	return newLdrBase(
		base{},
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

func (Builder) Ldrb(rt, rn Reg, off Off) (Instr, error) {
	return newLdrb(base{}, rt, rn, off)
}

func (Builder) Ldrh(rt, rn Reg, off Off) (Instr, error) {
	return newLdrh(base{}, rt, rn, off)
}

func (Builder) Ldrsb(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsb(base{}, rt, rn, off)
}

func (Builder) Ldrsh(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsh(base{}, rt, rn, off)
}

func (Builder) Ldrsw(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsw(base{}, rt, rn, off)
}

func (Builder) Ldur(rt, rn Reg, off Off) (Instr, error) {
	return newLdur(base{}, rt, rn, off)
}

func (Builder) Ldurb(rt, rn Reg, off Off) (Instr, error) {
	return newLdurb(base{}, rt, rn, off)
}

func (Builder) Ldurh(rt, rn Reg, off Off) (Instr, error) {
	return newLdurh(base{}, rt, rn, off)
}

func (Builder) LslReg(rd, rn, rm Reg) (Instr, error) {
	return newLslReg(base{}, rd, rn, rm)
}

func (Builder) LsrReg(rd, rn, rm Reg) (Instr, error) {
	return newLsrReg(base{}, rd, rn, rm)
}

func (Builder) Madd(rd, rn, rm, ra Reg) (Instr, error) {
	return newMadd(base{}, rd, rn, rm, ra)
}

// MlaElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) MlaElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewMlaElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewMlaElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newMlaElem(base{}, q, size, idx, rd, rn, rm)
}

// MlsElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) MlsElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewMlsElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewMlsElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newMlsElem(base{}, q, size, idx, rd, rn, rm)
}

// MovSimdBuilder entry: mov.8b/16b vd, vm (ORR-vector, Rn=31).
func (Builder) MovSimd(rd, rm VReg, arr string) (Instr, error) {
	err := requireArr("MovSimd", arr, "8b", "16b")
	if err != nil {
		return nil, err
	}

	enc := uint32(0x4EA01C00) // mov.16b (Q=1)
	if arr == "8b" {
		enc &^= 1 << 30
	}

	return MovSimd{
		rd:  rd.name(),
		rm:  rm.name(),
		arr: arr,
		enc: enc,
	}, nil
}

func (Builder) Movk(rd Reg, imm Imm16, hw Hw) (Instr, error) {
	return newMovk(base{}, rd, imm, hw)
}

func (Builder) Movn(rd Reg, imm Imm16, hw Hw) (Instr, error) {
	return newMovn(base{}, rd, imm, hw)
}

func (Builder) Movz(rd Reg, imm Imm16, hw Hw) (Instr, error) {
	return newMovz(base{}, rd, imm, hw)
}

func (Builder) Mrs(rd Reg, sysreg string) (Instr, error) {
	return newMrs(base{}, rd, sysreg)
}

func (Builder) Msr(sysreg string, rt Reg) (Instr, error) {
	return newMsr(base{}, sysreg, rt)
}

// Msub — msub rd, rn, rm, ra (mneg when Ra = zr); the operand
// constraints are those of Madd.
func (Builder) Msub(rd, rn, rm, ra Reg) (Instr, error) {
	base, err := New().Madd(rd, rn, rm, ra)
	if err != nil {
		return nil, err
	}

	m, ok := base.(Madd)
	if !ok {
		return nil, fmt.Errorf("msub: internal: want Madd, got %T", base)
	}

	return Msub{Madd: m}, nil
}

// MulElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) MulElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewMulElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewMulElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newMulElem(base{}, q, size, idx, rd, rn, rm)
}

// Nop — nop (no operands, fixed encoding).
func (Builder) Nop() Instr {
	in, err := newNop(base{})
	if err != nil {
		panic(err) // no operands - cannot fail
	}

	return in
}

func (Builder) Not(rd, rn VReg, arr string) (Instr, error) {
	return newNot(base{}, rd, rn, arr)
}

func (Builder) Orn(rd, rn, rm VReg, arr string) (Instr, error) {
	return newOrn(base{}, rd, rn, rm, arr)
}

func (Builder) OrnShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newOrnShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Orr(rd, rn, rm VReg, arr string) (Instr, error) {
	return newOrr(base{}, rd, rn, rm, arr)
}

func (Builder) OrrImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newOrrImm(base{}, rd, rn, imm)
}

func (Builder) OrrShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newOrrShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) Prfm(rn Reg) (Instr, error) {
	return newPrfm(base{}, rn)
}

func (Builder) Rbit(rd, rn Reg) (Instr, error) {
	return newRbit(base{}, rd, rn)
}

func (Builder) RbitV(rd, rn VReg, arr string) (Instr, error) {
	return newRbitV(base{}, rd, rn, arr)
}

func (Builder) Ret(rn Reg) (Instr, error) {
	return newRet(base{}, rn)
}

func (Builder) Rev(rd, rn Reg) (Instr, error) {
	return newRev(base{}, rd, rn)
}

func (Builder) Rev16(rd, rn Reg) (Instr, error) {
	return newRev16(base{}, rd, rn)
}

func (Builder) Rev32(rd, rn Reg) (Instr, error) {
	return newRev32(base{}, rd, rn)
}

func (Builder) Rev32V(rd, rn VReg, arr string) (Instr, error) {
	return newRev32V(base{}, rd, rn, arr)
}

func (Builder) RorReg(rd, rn, rm Reg) (Instr, error) {
	return newRorReg(base{}, rd, rn, rm)
}

func (Builder) Saddw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSaddw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf(
			"arm64.NewSaddw: arrangement too narrow (the source is one lane narrower)",
		)
	}

	return newSaddw(base{}, q, size-1, rd, rn, rm)
}

func (Builder) Sbfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newSbfm(base{}, rd, rn, immr, imms)
}

func (Builder) Scvtf(rd FReg, rn Reg) (Instr, error) {
	return newScvtf(base{}, rd, rn)
}

func (Builder) Sdiv(rd, rn, rm Reg) (Instr, error) {
	return newSdiv(base{}, rd, rn, rm)
}

// Shl - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Shl(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewShl: %w", err)
	}

	if shift >= 8<<size {
		return nil, fmt.Errorf("arm64.NewShl: shift %d out of range for .%s",
			shift, arr)
	}

	// the immediate field carries the size marker above the shift
	imm := 1<<(3+size) | shift
	immh, immb := imm>>3, imm&7

	return newShl(base{}, q, immh, immb, rd, rn)
}

// Smc — smc #imm16 (the secure-monitor call; PSCI rides it at #0).
func (Builder) Smc(imm Imm16) Instr {
	return sysImm{
		name:  "smc",
		imm16: imm.v,
		enc:   0xD4000003,
		shift: 5,
	}
}

// SmlalElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SmlalElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSmlalElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSmlalElem: %w", err)
	}

	return newSmlalElem(base{}, q, size-1, idx, rd, rn, rm)
}

// SmlslElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SmlslElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSmlslElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSmlslElem: %w", err)
	}

	return newSmlslElem(base{}, q, size-1, idx, rd, rn, rm)
}

func (Builder) Smov(wd Reg, vn VReg, elem string, idx uint32) (Instr, error) {
	size, err := elemSize("Smov", elem)
	if err != nil {
		return nil, err
	}

	return newSmov(base{}, 0, size, idx, vn, wd)
}

func (Builder) Smulh(rd, rn, rm Reg) (Instr, error) {
	return newSmulh(base{}, rd, rn, rm)
}

// SmullElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SmullElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSmullElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSmullElem: %w", err)
	}

	return newSmullElem(base{}, q, size-1, idx, rd, rn, rm)
}

// SqdmlalElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SqdmlalElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSqdmlalElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqdmlalElem: %w", err)
	}

	return newSqdmlalElem(base{}, q, size-1, idx, rd, rn, rm)
}

// SqdmlslElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SqdmlslElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSqdmlslElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqdmlslElem: %w", err)
	}

	return newSqdmlslElem(base{}, q, size-1, idx, rd, rn, rm)
}

// SqdmulhElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqdmulhElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqdmulhElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewSqdmulhElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newSqdmulhElem(base{}, q, size, idx, rd, rn, rm)
}

// SqdmullElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SqdmullElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSqdmullElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqdmullElem: %w", err)
	}

	return newSqdmullElem(base{}, q, size-1, idx, rd, rn, rm)
}

// SqrdmlahElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqrdmlahElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqrdmlahElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewSqrdmlahElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newSqrdmlahElem(base{}, q, size, idx, rd, rn, rm)
}

// SqrdmlshElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqrdmlshElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqrdmlshElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewSqrdmlshElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newSqrdmlshElem(base{}, q, size, idx, rd, rn, rm)
}

// SqrdmulhElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqrdmulhElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqrdmulhElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf(
			"arm64.NewSqrdmulhElem: arrangement %q is not one of the integer lane widths",
			arr,
		)
	}

	return newSqrdmulhElem(base{}, q, size, idx, rd, rn, rm)
}

func (Builder) Sqrshl(rd, rn, rm VReg, arr string) (Instr, error) {
	return newSqrshl(base{}, rd, rn, rm, arr)
}

// Sri - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Sri(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSri: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewSri: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newSri(base{}, q, immh, immb, rd, rn)
}

// Sshr - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Sshr(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSshr: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewSshr: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newSshr(base{}, q, immh, immb, rd, rn)
}

func (Builder) Ssubw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSsubw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf(
			"arm64.NewSsubw: arrangement too narrow (the source is one lane narrower)",
		)
	}

	return newSsubw(base{}, q, size-1, rd, rn, rm)
}

func (Builder) Stlr(rt, rn Reg) (Instr, error) {
	return newStlr(base{}, rt, rn)
}

func (Builder) Stlrb(rt, rn Reg) (Instr, error) {
	return newStlrb(base{}, rt, rn)
}

func (Builder) Stlxr(rs, rt, rn Reg) (Instr, error) {
	return newStlxr(base{}, rs, rt, rn)
}

func (Builder) Stlxrb(rs, rt, rn Reg) (Instr, error) {
	return newStlxrb(base{}, rs, rt, rn)
}

func (Builder) Stp(rt, rt2, rn Reg, off Off) (Instr, error) {
	return newStp(base{}, rt, rt2, rn, off)
}

func (Builder) Str(rt, rn Reg, off Off) (Instr, error) {
	return newStr(base{}, rt, rn, off)
}

// StrF — str st|dt, [xn, #off] (the FP/SIMD register form).
func (Builder) StrF(rt FReg, rn Reg, off Off) (Instr, error) {
	err := requireClass(
		rn,
		"StrF",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX, classSP,
	)
	if err != nil {
		return Str{}, err
	}

	enc, scale := strFDEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = strFSEnc, 2
	}

	if err := requireOff("StrF", off, scale); err != nil {
		return Str{}, err
	}

	return newStrBase(
		base{},
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

// Strb — strb rt, [rn, #off]: byte access, rt — w register only
// (register 31 reads as wzr), rn — x register or SP (register 31 in the
// base reads as sp); the offset is an unscaled imm12 (0..0xfff).
func (Builder) Strb(rt, rn Reg, off Off) (Instr, error) {
	return newStrb(base{}, rt, rn, off)
}

func (Builder) Strh(rt, rn Reg, off Off) (Instr, error) {
	return newStrh(base{}, rt, rn, off)
}

func (Builder) Stur(rt, rn Reg, off Off) (Instr, error) {
	return newStur(base{}, rt, rn, off)
}

func (Builder) Sturb(rt, rn Reg, off Off) (Instr, error) {
	return newSturb(base{}, rt, rn, off)
}

func (Builder) Sturh(rt, rn Reg, off Off) (Instr, error) {
	return newSturh(base{}, rt, rn, off)
}

func (Builder) Stxrb(rs, rt, rn Reg) (Instr, error) {
	return newStxrb(base{}, rs, rt, rn)
}

func (Builder) SubExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newSubExt(base{}, rd, rn, rm, ext, imm3)
}

func (Builder) SubImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newSubImm(base{}, rd, rn, imm, sh)
}

func (Builder) SubShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newSubShift(base{}, rd, rn, rm, imm, sh)
}

func (Builder) SubsExt(rd, rn, rm Reg, ext string, imm3 uint32) (Instr, error) {
	return newSubsExt(base{}, rd, rn, rm, ext, imm3)
}

func (Builder) SubsImm(rd, rn Reg, imm Imm12, sh Sh12) (Instr, error) {
	return newSubsImm(base{}, rd, rn, imm, sh)
}

func (Builder) SubsShift(rd, rn, rm Reg, imm Imm6, sh Shift) (Instr, error) {
	return newSubsShift(base{}, rd, rn, rm, imm, sh)
}

// Svc — svc #imm16.
func (Builder) Svc(imm Imm16) Instr {
	return sysImm{
		name:  "svc",
		imm16: imm.v,
		enc:   0xD4000001,
		shift: 5,
	}
}

func (Builder) Tbl(rd, rn, rm VReg) (Instr, error) {
	return newTbl(base{}, rd, rn, rm)
}

// Tbz — tbz rt, #bit, off: off — the pc-relative byte offset of the
// branch destination (the ±32KB imm14 range is checked at encode time;
// the absolute target is off + the instruction address). The register
// width is dictated by the bit number — the sf bit
// of the encoding is the bit's b5: bits 32..63 need an x register, bits
// 0..31 — a w one (register 31 reads as zr — use XZR/WZR).
func (Builder) Tbz(rt Reg, bit uint32, off int64) (Instr, error) {
	err := requireClass(
		rt,
		"Tbz",
		"rt",
		"x/w register (register 31 reads as zr — use XZR/WZR)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return nil, err
	}

	if bit > 63 {
		return nil, fmt.Errorf("arm64.NewTbz: operand bit: %d is out of 0..63", bit)
	}

	if rt.Is64() != (bit >= 32) {
		want := "w"
		if bit >= 32 {
			want = "x"
		}

		return nil, fmt.Errorf(
			"arm64.NewTbz: operand bit: %d needs a %s register (the sf bit of the encoding is the bit's b5), got %s",
			bit,
			want,
			rt.name(),
		)
	}

	return Tbz{
		rt:     rt.name(),
		bit:    bit,
		off:    immNum(off),
		isTbnz: false,
	}, nil
}

// Uaddlv - the Builder entry: .8b/.16b/.4h/.8h lanes.
func (Builder) Uaddlv(rd, rn VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUaddlv: %w", err)
	}

	if size > 1 {
		return nil, fmt.Errorf("arm64.NewUaddlv: arrangement %q is not one of [8b 16b 4h 8h]", arr)
	}

	return newUaddlv(base{}, q, size, rd, rn)
}

func (Builder) Uaddw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUaddw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf(
			"arm64.NewUaddw: arrangement too narrow (the source is one lane narrower)",
		)
	}

	return newUaddw(base{}, q, size-1, rd, rn, rm)
}

func (Builder) Ubfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newUbfm(base{}, rd, rn, immr, imms)
}

func (Builder) Ucvtf(rd FReg, rn Reg) (Instr, error) {
	return newUcvtf(base{}, rd, rn)
}

func (Builder) Udiv(rd, rn, rm Reg) (Instr, error) {
	return newUdiv(base{}, rd, rn, rm)
}

// UmlalElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) UmlalElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewUmlalElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUmlalElem: %w", err)
	}

	return newUmlalElem(base{}, q, size-1, idx, rd, rn, rm)
}

// UmlslElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) UmlslElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewUmlslElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUmlslElem: %w", err)
	}

	return newUmlslElem(base{}, q, size-1, idx, rd, rn, rm)
}

func (Builder) Umov(wd Reg, vn VReg, elem string, idx uint32) (Instr, error) {
	size, err := elemSize("Umov", elem)
	if err != nil {
		return nil, err
	}

	return newUmov(base{}, 0, size, idx, vn, wd)
}

func (Builder) Umulh(rd, rn, rm Reg) (Instr, error) {
	return newUmulh(base{}, rd, rn, rm)
}

// UmullElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) UmullElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewUmullElem: arrangement %q is not one of [4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUmullElem: %w", err)
	}

	return newUmullElem(base{}, q, size-1, idx, rd, rn, rm)
}

// Ushr - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Ushr(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUshr: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewUshr: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newUshr(base{}, q, immh, immb, rd, rn)
}

func (Builder) Usubw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUsubw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf(
			"arm64.NewUsubw: arrangement too narrow (the source is one lane narrower)",
		)
	}

	return newUsubw(base{}, q, size-1, rd, rn, rm)
}
