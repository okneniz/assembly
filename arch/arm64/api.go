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

// DecodeArrangement — the arrangement string by Q+size.
func DecodeArrangement(q, size uint32) string { return decodeArrangement(q, size) }

// RegNums2/3 — register names to numbers.
func RegNums2(a, b string) (uint32, uint32, error)            { return regNums2(a, b) }
func RegNums3(a, b, c string) (uint32, uint32, uint32, error) { return regNums3(a, b, c) }

// AddSubRegName — the register name by number, sf and flags.
func AddSubRegName(n uint32, isf bool, is64 bool) string { return addSubRegName(n, isf, is64) }

// ArmReg3 — three registers.
func ArmReg3(ops []ArmOp, name string) (string, string, string, error) { return armReg3(ops, name) }

// NewCsel — the Csel base of the conditional-select family.
func NewCsel(rd, rn, rm, cond string) Csel { return newCsel(rd, rn, rm, cond) }

// CondNum — a condition name to its 4-bit index.
func CondNum(name string) (uint32, error) { return condNum(name) }

// RegIndex — the register number of a vector register name.
func RegIndex(reg string) uint32 { return regIndex(reg) }

// RegListStr — the rendered "{ vN, vN+1... }" list.
func RegListStr(rt0 uint32, count int) string { return regListStr(rt0, count) }

// IsSimd3Logical — whether the three-same mnemonic uses the logical
// arrangement convention (only Q selects).
func IsSimd3Logical(name string) bool { return isSimd3Logical(name) }

// InvSysReg — a system register name → 15-bit key.
func InvSysReg(v any) (uint32, error) { return invSysReg(v) }

// VfpExpandImm64/32 — the fmov immediate expansion.
func VfpExpandImm64(imm8 uint32) float64 { return vfpExpandImm64(imm8) }
func VfpExpandImm32(imm8 uint32) float32 { return vfpExpandImm32(imm8) }

// CondNames — the condition names table.
func CondNames() [16]string { return condNames }

// SysregNames — the system register names table.
func SysregNames() map[uint32]string { return sysregNames }

// BrBits — the branch target to signed offset bits.
func OffBits(off int64, bits int) (uint32, error) { return offBits(off, bits) }
