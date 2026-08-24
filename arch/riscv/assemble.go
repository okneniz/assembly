package riscv

// Registry of new* constructors (mnemonic -> structure) and construction
// operand validation. Grammar, expression slots, and evaluation live in
// asm/riscv (the syntax layer above the arch); already computed operands
// (op.go) arrive here, and from here consumers (asm/riscv, pseudo) call
// BuildInstr.

import "errors"

// asmCtors - constructors of per-instruction structures from computed
// operands (mnemonic -> structure). The instruction encodes itself
// in Encode.
var asmCtors = map[string]func(ops []Op) (Instr, error){
	"lui":   newLui,
	"auipc": newAuipc,
	"add":   newAdd, "sub": newSub, "sll": newSll, "slt": newSlt, "sltu": newSltu,
	"xor": newXor, "srl": newSrl, "sra": newSra, "or": newOr, "and": newAnd,
	"mul": newMul, "mulh": newMulh, "mulhsu": newMulhsu, "mulhu": newMulhu,
	"div": newDiv, "divu": newDivu, "rem": newRem, "remu": newRemu,
	"addw": newAddw, "subw": newSubw, "sllw": newSllw, "srlw": newSrlw, "sraw": newSraw,
	"mulw": newMulw, "divw": newDivw, "divuw": newDivuw, "remw": newRemw, "remuw": newRemuw,
	"lb": newLb, "lh": newLh, "lw": newLw, "ld": newLd, "lbu": newLbu,
	"lhu": newLhu, "lwu": newLwu, "flw": newFlw, "fld": newFld,
	"sb": newSb, "sh": newSh, "sw": newSw, "sd": newSd, "fsw": newFsw, "fsd": newFsd,
	"beq": newBeq, "bne": newBne, "blt": newBlt, "bge": newBge,
	"bltu": newBltu, "bgeu": newBgeu,
	"jal": newJal, "jalr": newJalr,
	"fadd.s":    newFaddS,
	"fsub.s":    newFsubS,
	"fmul.s":    newFmulS,
	"fdiv.s":    newFdivS,
	"fadd.d":    newFaddD,
	"fsub.d":    newFsubD,
	"fmul.d":    newFmulD,
	"fdiv.d":    newFdivD,
	"fmadd.s":   newFmaddS,
	"fmsub.s":   newFmsubS,
	"fnmsub.s":  newFnmsubS,
	"fnmadd.s":  newFnmaddS,
	"fmadd.d":   newFmaddD,
	"fmsub.d":   newFmsubD,
	"fnmsub.d":  newFnmsubD,
	"fnmadd.d":  newFnmaddD,
	"amoadd.w":  newAmoaddW,
	"amoswap.w": newAmoswapW,
	"amoxor.w":  newAmoxorW,
	"amoor.w":   newAmoorW,
	"amoand.w":  newAmoandW,
	"amomin.w":  newAmominW,
	"amomax.w":  newAmomaxW,
	"amominu.w": newAmominuW,
	"amomaxu.w": newAmomaxuW,
	"amoadd.d":  newAmoaddD,
	"amoswap.d": newAmoswapD,
	"amoxor.d":  newAmoxorD,
	"amoor.d":   newAmoorD,
	"amoand.d":  newAmoandD,
	"amomin.d":  newAmominD,
	"amomax.d":  newAmomaxD,
	"amominu.d": newAmominuD,
	"amomaxu.d": newAmomaxuD,
	"csrrw":     newCsrrw, "csrrs": newCsrrs, "csrrc": newCsrrc,
	"csrrwi": newCsrrwi, "csrrsi": newCsrrsi, "csrrci": newCsrrci,
	"fence": newFence, "ecall": newSystem("ecall"), "ebreak": newSystem("ebreak"),
	"mret": newSystem("mret"), "sret": newSystem("sret"), "wfi": newSystem("wfi"),
	"addi": newAddi, "slti": newSlti, "sltiu": newSltiu,
	"xori": newXori, "ori": newOri, "andi": newAndi,
	"addiw": newAddiw,
	"slli":  newSlli, "srli": newSrli, "srai": newSrai,
	"slliw": newSlliw, "srliw": newSrliw, "sraiw": newSraiw,
}

// wantReg - validates a register operand (fp - a floating register is expected).
func wantReg(op Op, fp bool) (string, error) {
	if op.kind != opRegK {
		return "", errors.New("want register operand")
	}

	if _, err := rvRegNumOf(op.reg, fp); err != nil {
		return "", err
	}

	return op.reg, nil
}

// wantExpr - validates a numeric operand (immediate).
func wantExpr(op Op) (imm, error) {
	if op.kind != opNumK {
		return imm{}, errors.New("want immediate operand")
	}

	return immNum(op.num), nil
}
