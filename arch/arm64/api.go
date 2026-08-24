package arm64

// Export surface for the alias package (asm/arm64/alias): aliases sit ABOVE
// the arch and reuse the real structs and the constructor machinery (an
// alias — a mnemonic without its own encoding: the constructor builds
// another struct). No new common structs — only aliases of existing types,
// accessor methods and wrappers; fields stay unexported.

// Public names of existing types (aliases, not copies). ArmOp — a computed
// construction operand (vOp): encodeARM resolves slots before dispatching
// to constructors.
type (
	ArmOp  = vOp
	ArmMem = vMem
)

// ArmCtor — the constructor of a per-instruction struct from computed
// operands.
type ArmCtor = func(ops []ArmOp) (Instr, error)

// --- operand access (vOp fields stay unexported) ---

func (o vOp) IsReg() bool {
	return o.kind == armOpReg
}
func (o vOp) IsImm() bool {
	return o.kind == armOpImm
}
func (o vOp) IsShift() bool {
	return o.kind == armOpShift
}
func (o vOp) IsExtend() bool {
	return o.kind == armOpExtend
}
func (o vOp) IsLit() bool {
	return o.kind == armOpLit
}
func (o vOp) Reg() string {
	return o.reg
}
func (o vOp) Num() int64 {
	return o.num
}
func (o vOp) Sym() string {
	return o.sym
}
func (o vOp) ShiftName() string {
	return o.shift
}

// --- operand validation and parsing ---

func WantAReg(op ArmOp, name string) (string, error) {
	return wantAReg(op, name)
}
func WantCond(op ArmOp, name string) (string, error) {
	return wantCond(op, name)
}

func ArmReg2(ops []ArmOp, name string) (string, string, error) {
	return armReg2(ops, name)
}

func ArmRegNum(name string) (uint32, error) {
	return armRegNum(name)
}

// RegNameXOf — the X register name by number (same numbering as the W name).
func RegNameXOf(num uint32) string {
	return regNameX(num)
}

func ZeroReg(rd string) string {
	return zeroReg(rd)
}
func IsGPR(name string) bool {
	return isGPR(name)
}
func ShiftAmt(op ArmOp) int64 {
	return shiftAmt(op)
}
func InvertCond(c string) string {
	return invertCond(c)
}

// EncodeBitMasks — logical immediate encoding (tst/mov bitmasks).
func EncodeBitMasks(is64 bool, value uint64) (n, immr, imms uint32, ok bool) {
	return encodeBitMasks(is64, value)
}

// --- family constructor machinery (reused by aliases) ---

// AddSubThird — the third operand (imm | reg | reg+modifier) selects the
// add/sub family form; cmp/cmn call it with rdN=31 (zr).
func AddSubThird(ops []ArmOp, base string, rdN, rnN uint32, idx int) (Instr, error) {
	return addSubThird(ops, base, rdN, rnN, idx)
}

// MaddCtor/Msub3 — constructors of the madd/msub family (mul/mneg: ra = zr).
func MaddCtor(ops []ArmOp, name string) (Instr, error) {
	return makeMaddCtor(ops, name)
}
func Msub3(ops []ArmOp, name string) (Instr, error) {
	return msub3(ops, name)
}

// CsOf — assemble a csel-family struct by the base encoding name.
func CsOf(base, rd, rn, rm, cond string) (Instr, error) {
	return csOf(base, rd, rn, rm, cond)
}

// --- alias struct builders (struct fields stay unexported) ---

func SubShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) Instr {
	return SubShift{
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

func SubsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) Instr {
	return SubsShift{
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

func AndsImmOf(rd, rn string, immr, imms uint32, n, is64 bool) Instr {
	return AndsImm{logImm: newLogImm(rd, rn, immr, imms, n, is64)}
}

func AndsShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) Instr {
	return AndsShift{
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

func OrnShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) Instr {
	return OrnShift{
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

func OrrShiftOf(rd, rn, rm string, imm6 uint32, shift string, isf bool) Instr {
	return OrrShift{
		rd:    rd,
		rn:    rn,
		rm:    rm,
		imm6:  imm6,
		shift: shift,
		isf:   isf,
	}
}

func OrrImmOf(rd, rn string, immr, imms uint32, n, is64 bool) Instr {
	return OrrImm{logImm: newLogImm(rd, rn, immr, imms, n, is64)}
}

func MovzOf(rd string, imm16, hw uint32) Instr {
	return Movz{
		rd:    rd,
		imm16: imm16,
		hw:    hw,
	}
}

func MovnOf(rd string, imm16, hw uint32) Instr {
	return Movn{
		rd:    rd,
		imm16: imm16,
		hw:    hw,
	}
}

func SbfmOf(rd, rn string, immr, imms uint32, isf bool) Instr {
	return Sbfm{
		rd:   rd,
		rn:   rn,
		immr: immr,
		imms: imms,
		isf:  isf,
	}
}

func UbfmOf(rd, rn string, immr, imms uint32, isf bool) Instr {
	return Ubfm{
		rd:   rd,
		rn:   rn,
		immr: immr,
		imms: imms,
		isf:  isf,
	}
}

// VerifyBitMasks — a self-test of logical immediate encoding
// (encodeBitMasks ↔ decodeBitMasks).
func VerifyBitMasks() error {
	return verifyBitMasks()
}
