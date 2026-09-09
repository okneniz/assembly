package riscv

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr; compression
// is value-driven - see resolve.go).

import (
	asm "github.com/okneniz/assembly/asm"
)

// instr is an unevaluated instruction (mnemonic + operand slots); it
// implements Unresolved: Resolve evaluates the expressions and builds
// the evaluated form (expr -> numbers -> arch.BuildInstr; compression
// is value-driven - see resolve.go).
type instr struct {
	mnem string
	ops  []Op
	be   *Backend
}

// newInstr - an unevaluated instruction bound to its backend.
func newInstr(mnem string, ops []Op, be *Backend) instr {
	return instr{
		mnem: mnem,
		ops:  ops,
		be:   be,
	}
}

// Resolve evaluates the expressions and builds the evaluated
// instruction; the current .option norvc mode -> NoRVC (compression
// itself is value-driven, the core relaxes the layout).
func (in instr) Resolve(ctx asm.Ctx) (asm.Resolved, error) {
	return in.resolve(ctx, in.be != nil && in.be.noRVC)
}

// ResolveForm evaluates the "mnemonic + operands" form in the ctx
// environment with the current .option modes (pseudo expansion: each
// form starts from its own address, pass a derived ctx).
func (b *Backend) ResolveForm(mnem string, ops []Op, ctx asm.Ctx) (asm.Resolved, error) {
	return newInstr(mnem, ops, b).resolve(ctx, b.noRVC)
}
