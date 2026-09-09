package loong64

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr).

import (
	asm "github.com/okneniz/assembly/asm"
)

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr).
type instr struct {
	mnem string
	ops  []Op
}

// newInstr - an unevaluated instruction.
func newInstr(mnem string, ops []Op) instr {
	return instr{
		mnem: mnem,
		ops:  ops,
	}
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
	return newInstr(mnem, ops).resolve(ctx)
}
