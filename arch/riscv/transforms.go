// Package riscv - per-instruction RISC-V structures: decoders, formatters
// (objdump notation), and instruction constructors — the Builder methods
// (Addi, Lw, ...), the exact inverse of decode.
//
// RISC-V register names in ABI notation (as objdump prints them): integer
// and floating; used by the decode constructors of per-instruction
// structures and by the assembler register table.
package riscv

// rvRegNames - register names in RISC-V ABI notation (as objdump prints them):
// zero, ra, sp, gp, tp, t0-t2, s0/fp, s1, a0-a7, s2-s11, t3-t6.
var rvRegNames = [32]string{
	"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
	"s0", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
	"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
	"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
}

// rvFRegNames - ABI names of FP registers (as objdump prints them):
// ft0-ft7, fs0/fs1, fa0-fa7, fs2-fs11, ft8-ft11.
var rvFRegNames = [32]string{
	"ft0", "ft1", "ft2", "ft3", "ft4", "ft5", "ft6", "ft7",
	"fs0", "fs1", "fa0", "fa1", "fa2", "fa3", "fa4", "fa5",
	"fa6", "fa7", "fs2", "fs3", "fs4", "fs5", "fs6", "fs7",
	"fs8", "fs9", "fs10", "fs11", "ft8", "ft9", "ft10", "ft11",
}

// rvCsrNames - CSR names, generated from the canonical Spike encoding.h
// (see arch/riscv/csr_generated.go and gen/cmd/gen-riscv-csr).
// var rvCsrNames defined in csr_generated.go
