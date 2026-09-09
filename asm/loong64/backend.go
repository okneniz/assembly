// Package loong64 is the LoongArch (LA64) assembler syntax layer: the
// instruction grammar and unevaluated instructions (operands with
// expression slots). It sits ABOVE arch/loong64 (concrete instructions:
// decode/encode) and the asm core (sections, symbols, passes).
// Pseudo-instructions (nop/move/ret/li/la/call/...) are a layer above
// this package (pseudo). A mirror of MakeDecoder: the operand grammar accepts
// both GNU-as syntax and ObjDump()'s own output (round-trip).
package loong64

import (
	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	asm "github.com/okneniz/assembly/asm"
)

// New returns a LoongArch Syntax for asm.Assemble. An instance is
// intended for a single assembly.
func New() *Backend {
	b := &Backend{g: makeGrammar()}
	b.parseInstruction = b.makeInstructionParser()
	b.parseComment = makeCommentParser()
	return b
}

// Backend implements asm.Syntax. LoongArch has no .option modes in this
// model (fixed 32-bit words, no compression), so ApplyOption and
// ResetOptions are no-ops. The grammar (mnemonic/operand combinators)
// is built once in New and captured here - Instruction/Comment return
// the stored values.
type Backend struct {
	g                *grammar
	parseInstruction parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved]
	parseComment     parsec.Combinator[rune, parsecstrings.Position, string]
}

// ApplyOption accepts and ignores .option values (outside the model:
// the encoding of an LA64 scalar integer instruction depends on
// nothing).
func (b *Backend) ApplyOption(string) error {
	return nil
}

// ResetOptions resets the modes to the defaults (the core calls it
// before every pass).
func (b *Backend) ResetOptions() {}
