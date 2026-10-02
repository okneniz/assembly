package riscv

import "fmt"

// Builder - the instruction vocabulary of the package: per-instruction
// constructor methods named after the mnemonic they build (Addi, Lw,
// Bl...). Stateless by design — a namespace, not a factory.
type Builder struct{}

// New - the Builder for instruction constructors.
func New() Builder {
	return Builder{}
}

// Add - add rd, rs1, rs2.
func (Builder) Add(rd, rs1, rs2 Reg) Instr {
	return Add{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Addi - addi rd, rs1, imm (rs1 = zero is printed as li, imm = 0 - mv).
func (Builder) Addi(rd, rs1 Reg, imm Imm12) Instr {
	return Addi{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Addiw - addiw rd, rs1, imm (imm = 0 is printed as sext.w).
func (Builder) Addiw(rd, rs1 Reg, imm Imm12) Instr {
	return Addiw{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Addw - addw rd, rs1, rs2.
func (Builder) Addw(rd, rs1, rs2 Reg) Instr {
	return Addw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoaddD - amoadd.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] += rs2.
func (Builder) AmoaddD(rd, rs1, rs2 Reg) Instr {
	return AmoaddD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoaddW - amoadd.w rd, rs2, (rs1).
func (Builder) AmoaddW(rd, rs1, rs2 Reg) Instr {
	return AmoaddW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoandD - amoand.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] &= rs2.
func (Builder) AmoandD(rd, rs1, rs2 Reg) Instr {
	return AmoandD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoandW - amoand.w rd, rs2, (rs1).
func (Builder) AmoandW(rd, rs1, rs2 Reg) Instr {
	return AmoandW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmomaxD - amomax.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = max(MEM[rs1], rs2), signed.
func (Builder) AmomaxD(rd, rs1, rs2 Reg) Instr {
	return AmomaxD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmomaxW - amomax.w rd, rs2, (rs1).
func (Builder) AmomaxW(rd, rs1, rs2 Reg) Instr {
	return AmomaxW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmomaxuD - amomaxu.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = max(MEM[rs1], rs2), unsigned.
func (Builder) AmomaxuD(rd, rs1, rs2 Reg) Instr {
	return AmomaxuD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmomaxuW - amomaxu.w rd, rs2, (rs1).
func (Builder) AmomaxuW(rd, rs1, rs2 Reg) Instr {
	return AmomaxuW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmominD - amomin.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = min(MEM[rs1], rs2), signed.
func (Builder) AmominD(rd, rs1, rs2 Reg) Instr {
	return AmominD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmominW - amomin.w rd, rs2, (rs1).
func (Builder) AmominW(rd, rs1, rs2 Reg) Instr {
	return AmominW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmominuD - amominu.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = min(MEM[rs1], rs2), unsigned.
func (Builder) AmominuD(rd, rs1, rs2 Reg) Instr {
	return AmominuD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmominuW - amominu.w rd, rs2, (rs1).
func (Builder) AmominuW(rd, rs1, rs2 Reg) Instr {
	return AmominuW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoorD - amoor.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] |= rs2.
func (Builder) AmoorD(rd, rs1, rs2 Reg) Instr {
	return AmoorD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoorW - amoor.w rd, rs2, (rs1).
func (Builder) AmoorW(rd, rs1, rs2 Reg) Instr {
	return AmoorW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoswapD - amoswap.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] = rs2.
func (Builder) AmoswapD(rd, rs1, rs2 Reg) Instr {
	return AmoswapD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoswapW - amoswap.w rd, rs2, (rs1).
func (Builder) AmoswapW(rd, rs1, rs2 Reg) Instr {
	return AmoswapW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoxorD - amoxor.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] ^= rs2.
func (Builder) AmoxorD(rd, rs1, rs2 Reg) Instr {
	return AmoxorD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// AmoxorW - amoxor.w rd, rs2, (rs1).
func (Builder) AmoxorW(rd, rs1, rs2 Reg) Instr {
	return AmoxorW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// And - and rd, rs1, rs2.
func (Builder) And(rd, rs1, rs2 Reg) Instr {
	return And{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Andi - andi rd, rs1, imm (imm = 0xff is printed as zext.b).
func (Builder) Andi(rd, rs1 Reg, imm Imm12) Instr {
	return Andi{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Auipc - auipc rd, imm20 (imm20 = the raw U-type field: rd = pc + imm20<<12).
func (Builder) Auipc(rd Reg, imm Imm20) Instr {
	return Auipc{
		rd:  rd.name(),
		imm: immNum(imm.v),
	}
}

// Beq - beq rs1, rs2, off (the pc-relative byte offset; the absolute
// target is off + the instruction address).
func (Builder) Beq(rs1, rs2 Reg, off int64) Instr {
	return Beq{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Bge - bge rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bge(rs1, rs2 Reg, off int64) Instr {
	return Bge{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Bgeu - bgeu rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bgeu(rs1, rs2 Reg, off int64) Instr {
	return Bgeu{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Blt - blt rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Blt(rs1, rs2 Reg, off int64) Instr {
	return Blt{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Bltu - bltu rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bltu(rs1, rs2 Reg, off int64) Instr {
	return Bltu{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Bne - bne rs1, rs2, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Bne(rs1, rs2 Reg, off int64) Instr {
	return Bne{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off),
	}
}

// Csrrc - csrrc rd, csr, rs1; csr is a 12-bit CSR number 0..4095
// (the canonical name is resolved for display when known).
func (Builder) Csrrc(rd Reg, csr uint16, rs1 Reg) Instr {
	return Csrrc{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		rs1:   rs1.name(),
	}
}

// Csrrci - csrrci rd, csr, zimm; csr is a 12-bit CSR number
// 0..4095, zimm a 5-bit immediate 0..31.
func (Builder) Csrrci(rd Reg, csr uint16, zimm uint8) Instr {
	return Csrrci{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		zimm:  immNum(int64(zimm)),
	}
}

// Csrrs - csrrs rd, csr, rs1; csr is a 12-bit CSR number 0..4095
// (the canonical name is resolved for display when known).
func (Builder) Csrrs(rd Reg, csr uint16, rs1 Reg) Instr {
	return Csrrs{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		rs1:   rs1.name(),
	}
}

// Csrrsi - csrrsi rd, csr, zimm; csr is a 12-bit CSR number
// 0..4095, zimm a 5-bit immediate 0..31.
func (Builder) Csrrsi(rd Reg, csr uint16, zimm uint8) Instr {
	return Csrrsi{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		zimm:  immNum(int64(zimm)),
	}
}

// Csrrw - csrrw rd, csr, rs1; csr is a 12-bit CSR number 0..4095
// (the canonical name is resolved for display when known).
func (Builder) Csrrw(rd Reg, csr uint16, rs1 Reg) Instr {
	return Csrrw{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		rs1:   rs1.name(),
	}
}

// Csrrwi - csrrwi rd, csr, zimm; csr is a 12-bit CSR number
// 0..4095, zimm a 5-bit immediate 0..31.
func (Builder) Csrrwi(rd Reg, csr uint16, zimm uint8) Instr {
	return Csrrwi{
		csrOp: newCsrOp(rd.name(), int64(csr)),
		zimm:  immNum(int64(zimm)),
	}
}

// Div - div rd, rs1, rs2.
func (Builder) Div(rd, rs1, rs2 Reg) Instr {
	return Div{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Divu - divu rd, rs1, rs2.
func (Builder) Divu(rd, rs1, rs2 Reg) Instr {
	return Divu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Divuw - divuw rd, rs1, rs2.
func (Builder) Divuw(rd, rs1, rs2 Reg) Instr {
	return Divuw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Divw - divw rd, rs1, rs2.
func (Builder) Divw(rd, rs1, rs2 Reg) Instr {
	return Divw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// FaddD - fadd.d fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FaddD(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FaddD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FaddS - fadd.s fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FaddS(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FaddS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FdivD - fdiv.d fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FdivD(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FdivD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FdivS - fdiv.s fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FdivS(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FdivS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// Fence - fence fm; fm is a 4-bit fence modifier 0..15
// (fm 0 is printed as the bare "fence").
func (Builder) Fence(fm uint8) Instr {
	return Fence{
		fm: immNum(int64(fm)),
	}
}

// Fld - fld rd, off(rs1); rd is an FP register taken by number
// (Reg 0..31 is printed ft0/fa0/...).
func (Builder) Fld(rd, rs1 Reg, off Off) Instr {
	return Fld{
		rd:  fpName(rd),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Flw - flw rd, off(rs1); rd is an FP register taken by number
// (Reg 0..31 is printed ft0/fa0/...).
func (Builder) Flw(rd, rs1 Reg, off Off) Instr {
	return Flw{
		rd:  fpName(rd),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// FmaddD - fmadd.d fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmaddD(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FmaddD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FmaddS - fmadd.s fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmaddS(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FmaddS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FmsubD - fmsub.d fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmsubD(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FmsubD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FmsubS - fmsub.s fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmsubS(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FmsubS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FmulD - fmul.d fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmulD(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FmulD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FmulS - fmul.s fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FmulS(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FmulS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FnmaddD - fnmadd.d fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FnmaddD(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FnmaddD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FnmaddS - fnmadd.s fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FnmaddS(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FnmaddS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FnmsubD - fnmsub.d fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FnmsubD(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FnmsubD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// FnmsubS - fnmsub.s fd, fs1, fs2, fs3; the registers are FP registers taken
// by number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode
// 0..7 (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FnmsubS(rd, rs1, rs2, rs3 Reg, rm uint8) Instr {
	return FnmsubS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rs3: fpName(rs3),
		rm:  immNum(int64(rm)),
	}
}

// Fsd - fsd rs2, off(rs1); rs2 is an FP register taken by number
// (Reg 0..31 is printed ft0/fa0/...).
func (Builder) Fsd(rs2, rs1 Reg, off Off) Instr {
	return Fsd{
		rs1: rs1.name(),
		rs2: fpName(rs2),
		off: immNum(off.v),
	}
}

// FsubD - fsub.d fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FsubD(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FsubD{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// FsubS - fsub.s fd, fs1, fs2; the registers are FP registers taken by
// number (Reg 0..31 is printed ft0/fa0/...); rm is the rounding mode 0..7
// (0 RNE, 1 RTZ, 2 RDN, 3 RUP, 4 RMM, 7 DYN).
func (Builder) FsubS(rd, rs1, rs2 Reg, rm uint8) Instr {
	return FsubS{
		rd:  fpName(rd),
		rs1: fpName(rs1),
		rs2: fpName(rs2),
		rm:  immNum(int64(rm)),
	}
}

// Fsw - fsw rs2, off(rs1); rs2 is an FP register taken by number
// (Reg 0..31 is printed ft0/fa0/...).
func (Builder) Fsw(rs2, rs1 Reg, off Off) Instr {
	return Fsw{
		rs1: rs1.name(),
		rs2: fpName(rs2),
		off: immNum(off.v),
	}
}

// Imm12 - a validated value; an error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	if v < -2048 || v > 2047 {
		return Imm12{}, fmt.Errorf("riscv.New().Imm12: value %d outside -2048..2047", v)
	}

	return Imm12{v: v}, nil
}

// Imm20 - a validated value; an error when out of range.
func (Builder) Imm20(v int64) (Imm20, error) {
	if v < 0 || v > 0xfffff {
		return Imm20{}, fmt.Errorf("riscv.New().Imm20: value %d outside 0..%d", v, 0xfffff)
	}

	return Imm20{v: v}, nil
}

// Jal - jal rd, off (the pc-relative byte offset; the absolute target is off + the instruction address).
func (Builder) Jal(rd Reg, off int64) Instr {
	return Jal{
		rd:  rd.name(),
		off: immNum(off),
	}
}

// Jalr - jalr rd, off(rs1) (off = byte offset from rs1).
func (Builder) Jalr(rd, rs1 Reg, off Off) Instr {
	return Jalr{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// JalrReg - jalr rs (indirect call: the fixed 32-bit form jalr ra, 0(rs)).
func (Builder) JalrReg(rs1 Reg) Instr {
	return JalrReg{
		rs1: rs1.name(),
	}
}

// Lb - lb rd, off(rs1).
func (Builder) Lb(rd, rs1 Reg, off Off) Instr {
	return Lb{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Lbu - lbu rd, off(rs1).
func (Builder) Lbu(rd, rs1 Reg, off Off) Instr {
	return Lbu{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Ld - ld rd, off(rs1).
func (Builder) Ld(rd, rs1 Reg, off Off) Instr {
	return Ld{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Lh - lh rd, off(rs1).
func (Builder) Lh(rd, rs1 Reg, off Off) Instr {
	return Lh{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Lhu - lhu rd, off(rs1).
func (Builder) Lhu(rd, rs1 Reg, off Off) Instr {
	return Lhu{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Lui - lui rd, imm20.
func (Builder) Lui(rd Reg, imm Imm20) Instr {
	return Lui{
		rd:  rd.name(),
		imm: immNum(imm.v),
	}
}

// Lw - lw rd, off(rs1).
func (Builder) Lw(rd, rs1 Reg, off Off) Instr {
	return Lw{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Lwu - lwu rd, off(rs1).
func (Builder) Lwu(rd, rs1 Reg, off Off) Instr {
	return Lwu{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

// Mul - mul rd, rs1, rs2.
func (Builder) Mul(rd, rs1, rs2 Reg) Instr {
	return Mul{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Mulh - mulh rd, rs1, rs2.
func (Builder) Mulh(rd, rs1, rs2 Reg) Instr {
	return Mulh{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Mulhsu - mulhsu rd, rs1, rs2.
func (Builder) Mulhsu(rd, rs1, rs2 Reg) Instr {
	return Mulhsu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Mulhu - mulhu rd, rs1, rs2.
func (Builder) Mulhu(rd, rs1, rs2 Reg) Instr {
	return Mulhu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Mulw - mulw rd, rs1, rs2.
func (Builder) Mulw(rd, rs1, rs2 Reg) Instr {
	return Mulw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Mv - mv rd, rs2 (the c.mv halfword).
func (Builder) Mv(rd, rs2 Reg) Instr {
	h := uint32(0x8002) | uint32(rd.Num())<<7 | uint32(rs2.Num())<<2

	return Mv{
		base: newHalfBase(h),
		rd:   rd.name(),
		rs2:  rs2.name(),
	}
}

// Off - a validated value; an error when out of range.
func (Builder) Off(v int64) (Off, error) {
	if v < -2048 || v > 2047 {
		return Off{}, fmt.Errorf("riscv.New().Off: value %d outside -2048..2047", v)
	}

	return Off{v: v}, nil
}

// Or - or rd, rs1, rs2.
func (Builder) Or(rd, rs1, rs2 Reg) Instr {
	return Or{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Ori - ori rd, rs1, imm.
func (Builder) Ori(rd, rs1 Reg, imm Imm12) Instr {
	return Ori{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Rem - rem rd, rs1, rs2.
func (Builder) Rem(rd, rs1, rs2 Reg) Instr {
	return Rem{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Remu - remu rd, rs1, rs2.
func (Builder) Remu(rd, rs1, rs2 Reg) Instr {
	return Remu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Remuw - remuw rd, rs1, rs2.
func (Builder) Remuw(rd, rs1, rs2 Reg) Instr {
	return Remuw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Remw - remw rd, rs1, rs2.
func (Builder) Remw(rd, rs1, rs2 Reg) Instr {
	return Remw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Sb - sb rs2, off(rs1).
func (Builder) Sb(rs2, rs1 Reg, off Off) Instr {
	return Sb{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

// Sd - sd rs2, off(rs1).
func (Builder) Sd(rs2, rs1 Reg, off Off) Instr {
	return Sd{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

// Sh - sh rs2, off(rs1).
func (Builder) Sh(rs2, rs1 Reg, off Off) Instr {
	return Sh{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

// Sll - sll rd, rs1, rs2.
func (Builder) Sll(rd, rs1, rs2 Reg) Instr {
	return Sll{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Slli - slli rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Slli(rd, rs1 Reg, shamt Imm12) Instr {
	return Slli{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Slliw - slliw rd, rs1, shamt (shamt5).
func (Builder) Slliw(rd, rs1 Reg, shamt Imm12) Instr {
	return Slliw{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Sllw - sllw rd, rs1, rs2.
func (Builder) Sllw(rd, rs1, rs2 Reg) Instr {
	return Sllw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Slt - slt rd, rs1, rs2.
func (Builder) Slt(rd, rs1, rs2 Reg) Instr {
	return Slt{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Slti - slti rd, rs1, imm.
func (Builder) Slti(rd, rs1 Reg, imm Imm12) Instr {
	return Slti{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Sltiu - sltiu rd, rs1, imm (imm = 1 is printed as seqz).
func (Builder) Sltiu(rd, rs1 Reg, imm Imm12) Instr {
	return Sltiu{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}

// Sltu - sltu rd, rs1, rs2 (rs1 = zero is printed as snez).
func (Builder) Sltu(rd, rs1, rs2 Reg) Instr {
	return Sltu{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Sra - sra rd, rs1, rs2.
func (Builder) Sra(rd, rs1, rs2 Reg) Instr {
	return Sra{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Srai - srai rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Srai(rd, rs1 Reg, shamt Imm12) Instr {
	return Srai{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Sraiw - sraiw rd, rs1, shamt (shamt5).
func (Builder) Sraiw(rd, rs1 Reg, shamt Imm12) Instr {
	return Sraiw{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Sraw - sraw rd, rs1, rs2.
func (Builder) Sraw(rd, rs1, rs2 Reg) Instr {
	return Sraw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Srl - srl rd, rs1, rs2.
func (Builder) Srl(rd, rs1, rs2 Reg) Instr {
	return Srl{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Srli - srli rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Srli(rd, rs1 Reg, shamt Imm12) Instr {
	return Srli{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Srliw - srliw rd, rs1, shamt (shamt5).
func (Builder) Srliw(rd, rs1 Reg, shamt Imm12) Instr {
	return Srliw{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

// Srlw - srlw rd, rs1, rs2.
func (Builder) Srlw(rd, rs1, rs2 Reg) Instr {
	return Srlw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Sub - sub rd, rs1, rs2 (rs1 = zero is printed as neg).
func (Builder) Sub(rd, rs1, rs2 Reg) Instr {
	return Sub{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Subw - subw rd, rs1, rs2 (rs1 = zero is printed as negw).
func (Builder) Subw(rd, rs1, rs2 Reg) Instr {
	return Subw{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Sw - sw rs2, off(rs1).
func (Builder) Sw(rs2, rs1 Reg, off Off) Instr {
	return Sw{
		rs1: rs1.name(),
		rs2: rs2.name(),
		off: immNum(off.v),
	}
}

// Xor - xor rd, rs1, rs2.
func (Builder) Xor(rd, rs1, rs2 Reg) Instr {
	return Xor{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

// Xori - xori rd, rs1, imm (imm = -1 is printed as not).
func (Builder) Xori(rd, rs1 Reg, imm Imm12) Instr {
	return Xori{
		rd:  rd.name(),
		rs1: rs1.name(),
		imm: immNum(imm.v),
	}
}
