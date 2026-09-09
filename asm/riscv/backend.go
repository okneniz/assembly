// Package riscv is the RISC-V (RV64GC) assembler syntax layer: the
// instruction grammar, .option modes, and unevaluated instructions
// (operands with expression slots). It sits ABOVE arch/riscv
// (concrete instructions: decode/encode) and the asm core (sections,
// symbols, passes).
// Pseudo-instructions
// (nop/li/mv/ret/call/la/...) are a layer above this package (pseudo).
// A mirror of MakeDecoder: the operand grammar accepts both GNU-as syntax
// and ObjDump()'s own output (round-trip).
package riscv

import (
	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	asm "github.com/okneniz/assembly/asm"
)

// New returns a RISC-V Syntax for asm.Assemble. An instance is
// intended for a single assembly: ApplyOption mutates the (.option)
// modes; create separate instances for parallel assemblies.
func New() *Backend {
	b := &Backend{g: makeGrammar()}
	b.parseInstruction = b.makeInstructionParser()
	b.parseComment = makeCommentParser()
	return b
}

// Backend implements asm.Syntax; it owns the (.option) encoding modes:
// noRVC disables auto-compression, optStack holds push/pop snapshots.
// The grammar (mnemonic/operand combinators) is built once in New and
// captured here - Instruction/Comment return the stored values.
type Backend struct {
	noRVC            bool
	optStack         []bool
	g                *grammar
	parseInstruction parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved]
	parseComment     parsec.Combinator[rune, parsecstrings.Position, string]
}

// ApplyOption handles .option values: norvc/rvc toggle auto-compression,
// push/pop are a stack of mode snapshots; the rest (pic/nopic/relax/
// norelax/csr-check/arch) are outside the model (an absolute assembler
// without relaxation, a fixed RV64GC profile) and do not affect the
// bytes - they are ignored.
func (b *Backend) ApplyOption(name string) error {
	switch name {
	case "norvc":
		b.noRVC = true
	case "rvc":
		b.noRVC = false
	case "push":
		b.optStack = append(b.optStack, b.noRVC)
	case "pop":
		if n := len(b.optStack); n > 0 {
			b.noRVC = b.optStack[n-1]
			b.optStack = b.optStack[:n-1]
		}
	}

	return nil
}

// ResetOptions resets the modes to the defaults (the core calls it
// before every pass: .option is applied symmetrically from one
// initial state).
func (b *Backend) ResetOptions() {
	b.noRVC = false
	b.optStack = nil
}
