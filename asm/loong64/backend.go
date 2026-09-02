// Package loong64 is the LoongArch (LA64) assembler syntax layer: the
// instruction grammar and unevaluated instructions (operands with
// expression slots). It sits ABOVE arch/loong64 (concrete instructions:
// decode/encode) and the asm core (sections, symbols, passes).
// Pseudo-instructions (nop/move/ret/li/la/call/...) are a layer above
// this package (pseudo). A mirror of Parse: the operand grammar accepts
// both GNU-as syntax and ObjDump()'s own output (round-trip).
package loong64

import (
	"fmt"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	asm "github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
)

// New returns a LoongArch Syntax for asm.Assemble. An instance is
// intended for a single assembly.
func New() *Backend {
	return &Backend{}
}

// Backend implements asm.Syntax. LoongArch has no .option modes in this
// model (fixed 32-bit words, no compression), so ApplyOption and
// ResetOptions are no-ops.
type Backend struct{}

// ApplyOption accepts and ignores .option values (outside the model:
// the encoding of an LA64 scalar integer instruction depends on
// nothing).
func (b *Backend) ApplyOption(string) error {
	return nil
}

// ResetOptions resets the modes to the defaults (the core calls it
// before every pass).
func (b *Backend) ResetOptions() {}

// Instruction is the grammar "mnemonic operands" (comma-separated
// operands).
func (b *Backend) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		skipSpaces(buf)
		name, err := cMnemonic(buf)
		if err != nil {
			return nil, err
		}

		// mnemonic boundary: followed by a space/comma/end of line
		if r, ok := peekRune(buf); ok && r != ' ' && r != '\t' && r != ',' && r != '\n' {
			return nil, parsec.NewParseError(pos, fmt.Sprintf("unknown mnemonic %q", name))
		}

		ops, err := ParseOps(buf)
		if err != nil {
			return nil, err
		}

		return instr{mnem: name, ops: ops}, nil
	}
}

// Comment parses '#' and '//' to the end of the line.
func (b *Backend) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	body := parsecstrings.Many(4, expr.CNotNL)
	hash := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.Eq("comment", '#')), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)
	slash := parsecstrings.Cast(
		parsecstrings.Skip(parsecstrings.Try(parsecstrings.String("comment", "//")), body),
		func(rs []rune) (string, error) {
			return string(rs), nil
		},
	)

	return parsecstrings.Choice("comment", parsecstrings.Try(slash), parsecstrings.Try(hash))
}

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr).
type instr struct {
	mnem string
	ops  []Op
}

// Resolve evaluates the expressions and builds the evaluated
// instruction.
func (in instr) Resolve(ctx asm.Ctx) (asm.Resolved, error) {
	return in.resolve(ctx)
}

// ResolveForm evaluates the "mnemonic + operands" form in the ctx
// environment (pseudo expansion: each form starts from its own address,
// pass a derived ctx).
func (b *Backend) ResolveForm(mnem string, ops []Op, ctx asm.Ctx) (asm.Resolved, error) {
	return instr{mnem: mnem, ops: ops}.resolve(ctx)
}
