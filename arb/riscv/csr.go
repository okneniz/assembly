package riscv

// The CSR core: the 6 atomic read/write/modify forms over a 12-bit CSR
// address (the register forms draw a register, the immediate forms a
// 5-bit zimm). The generator walks the named CSR pool — the text layer
// prints and parses CSR names; the full 460-name table is
// arch-tested. The per-family generators live in their own files
// (csrrw.go, ...).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// csrPool — the named CSR addresses of the generator pool (the base
// ISA status registers).
var csrPool = []uint16{
	0x001, // fflags
	0x002, // frm
	0x003, // fcsr
	0x300, // mstatus
	0x301, // misa
	0x304, // mie
	0x305, // mtvec
	0x340, // mscratch
	0x341, // mepc
	0x342, // mcause
	0x343, // mtval
	0xc00, // cycle
	0xc01, // time
	0xc02, // instret
}

// CsrOp — the CSR mnemonic.
type CsrOp int

const (
	CsrCsrrw CsrOp = iota
	CsrCsrrs
	CsrCsrrc
	CsrCsrrwi
	CsrCsrrsi
	CsrCsrrci
	csrOpCount
)

// CsrParams — parameters of the CSR families (the register form reads
// Rs1, the immediate form Zimm).
type CsrParams struct {
	Op   CsrOp
	Rd   riscv.Reg
	Csr  uint16
	Rs1  riscv.Reg
	Zimm uint8
}

func NewCsrParams(op CsrOp, rd riscv.Reg, csr uint16, rs1 riscv.Reg, zimm uint8) CsrParams {
	return CsrParams{
		Op:   op,
		Rd:   rd,
		Csr:  csr,
		Rs1:  rs1,
		Zimm: zimm,
	}
}

func (p CsrParams) Instr() riscv.Instr {
	b := riscv.New()
	switch p.Op {
	case CsrCsrrw:
		return b.Csrrw(p.Rd, p.Csr, p.Rs1)
	case CsrCsrrs:
		return b.Csrrs(p.Rd, p.Csr, p.Rs1)
	case CsrCsrrc:
		return b.Csrrc(p.Rd, p.Csr, p.Rs1)
	case CsrCsrrwi:
		return b.Csrrwi(p.Rd, p.Csr, p.Zimm)
	case CsrCsrrsi:
		return b.Csrrsi(p.Rd, p.Csr, p.Zimm)
	default:
		return b.Csrrci(p.Rd, p.Csr, p.Zimm)
	}
}

func (p CsrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// immediate — the ci/si forms take a zimm instead of a register.
func (o CsrOp) immediate() bool {
	return o >= CsrCsrrwi
}

// csrGen — the shared generator of the CSR group.
type csrGen struct {
	rnd *rand.Rand
	op  CsrOp
}

// Csr — an arbitrary CSR instruction of one family (the per-family
// files bind their mnemonic).
func Csr(rnd *rand.Rand, op CsrOp) ohsnap.Arbitrary[CsrParams] {
	return csrGen{rnd: rnd, op: op}
}

func (g csrGen) Generate() iter.Seq[CsrParams] {
	return arb.Stream(func() CsrParams {
		return NewCsrParams(
			g.op,
			reg(g.rnd),
			csrPool[g.rnd.IntN(len(csrPool))],
			reg(g.rnd),
			uint8(g.rnd.IntN(32)),
		)
	})
}

func (g csrGen) Shrink(p CsrParams) iter.Seq[CsrParams] {
	out := make([]CsrParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewCsrParams(p.Op, r, p.Csr, p.Rs1, p.Zimm))
	}

	if p.Op.immediate() {
		for _, z := range zimmShrunk(p.Zimm) {
			out = append(out, NewCsrParams(p.Op, p.Rd, p.Csr, p.Rs1, z))
		}

		return slices.Values(out)
	}

	for _, r := range regShrunk(p.Rs1) {
		out = append(out, NewCsrParams(p.Op, p.Rd, p.Csr, r, p.Zimm))
	}

	return slices.Values(out)
}

// zimmShrunk — the zimm shrink candidates: zero and the halved value.
func zimmShrunk(v uint8) []uint8 {
	if v == 0 {
		return nil
	}

	out := []uint8{0}
	if h := v / 2; h > 0 {
		out = append(out, h)
	}

	return out
}
