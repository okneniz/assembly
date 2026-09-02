// Package pseudo provides LoongArch pseudo-instructions (GNU as + the
// ObjDump text's own aliases): a layer ABOVE asm/loong64 (the grammar)
// and arch/loong64 (the real instructions). A pseudo-instruction is
// assembler syntax: the processor executes only real encodings, and
// decoding never returns pseudo-instructions (the pseudo text on the
// decode side comes from formatters of the real structures).
//
// Two classes (as in the mnemonic dictionary):
//   - multi-word (la) - its own evaluated form with a fixed length,
//     for deterministic assembler passes;
//   - single-word (nop/move/not/ret/jr/call/tail/bltz/... and the
//     li.w/li.d ladders) - expansion into base forms and the syntax
//     layer's ResolveForm.
package pseudo

import (
	"bytes"
	"fmt"
	"io"

	"github.com/okneniz/parsec"
	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/expr"
	loong64 "github.com/okneniz/assembly/asm/loong64"
)

// NewASMBackend returns a LoongArch Syntax with pseudo-instructions: a
// decorator over the asm/loong64 syntax layer for asm.Assemble.
func NewASMBackend() asm.Syntax {
	be := loong64.New()

	return source{be: be}
}

// Assemble is the full LoongArch assembly (syntax layer + pseudo) by
// the asm core.
func Assemble(src string, base uint64) (*asm.Result, []asm.AsmError) {
	return asm.Assemble(src, base, NewASMBackend())
}

// source is an asm.Syntax decorator: it parses pseudo-mnemonics
// itself and delegates the rest to the syntax layer's inner grammar.
type source struct {
	be *loong64.Backend
}

// pInstr is an unevaluated pseudo-instruction (mnemonic + operands).
type pInstr struct {
	mnem string
	ops  []loong64.Op
	be   *loong64.Backend
}

func newPInstr(mnem string, ops []loong64.Op, be *loong64.Backend) pInstr {
	return pInstr{
		mnem: mnem,
		ops:  ops,
		be:   be,
	}
}

// ctors are the evaluators of the multi-word pseudo-instructions
// (their own fixed-length forms).
var ctors = map[string]func(ops []loong64.Op, ctx asm.Ctx) (asm.Resolved, error){
	"la": resolveLa,
}

// cPseudoMnemonic is the trie of pseudo-mnemonics (longest-match).
// Built in a var initialization (as in the syntax layer): MapStrings
// builds a trie at creation.
var cPseudoMnemonic = parsecstrings.MapStrings("mnemonic", pseudoMnemonics)

var pseudoMnemonics = buildPseudoMnemonics()

func buildPseudoMnemonics() map[string]string {
	m := map[string]string{}
	for p := range ctors {
		m[p] = p
	}

	for _, p := range expandable {
		m[p] = p
	}

	return m
}

// Instruction is the grammar "mnemonic operands": a pseudo-mnemonic
// from its own trie (with the same boundary check), otherwise the
// syntax layer's inner grammar.
func (s source) Instruction() parsec.Combinator[rune, parsecstrings.Position, asm.Unresolved] {
	return func(buf parsec.Buffer[rune, parsecstrings.Position]) (asm.Unresolved, parsec.Error[parsecstrings.Position]) {
		pos := buf.Position()
		expr.SkipSpaces(buf)
		if name, err := cPseudoMnemonic(buf); err == nil {
			// mnemonic boundary: followed by a space/comma/end of
			// line; otherwise - rewind to the inner grammar
			if r, ok := expr.PeekRune(buf); !ok || r == ' ' || r == '\t' || r == ',' || r == '\n' {
				ops, oerr := loong64.ParseOps(buf)
				if oerr != nil {
					return nil, oerr
				}

				return newPInstr(name, ops, s.be), nil
			}
		}

		if rerr := expr.Rewind(buf, pos); rerr != nil {
			return nil, rerr
		}

		return s.be.Instruction()(buf)
	}
}

// Comment parses a comment ('#' and '//' to the end of the line), as
// in the syntax layer.
func (s source) Comment() parsec.Combinator[rune, parsecstrings.Position, string] {
	return s.be.Comment()
}

// ApplyOption is delegated to the syntax layer.
func (s source) ApplyOption(name string) error {
	return s.be.ApplyOption(name)
}

// ResetOptions is delegated to the syntax layer.
func (s source) ResetOptions() {
	s.be.ResetOptions()
}

// Resolve evaluates the pseudo-instruction (with expansion): the
// multi-word ones get their own forms, the single-word ones expand
// into base forms + ResolveForm. The core calls it twice (layout with
// a placeholder environment and code generation) - the decisions (the
// li ladder) do not depend on symbol values (a symbolic li takes the
// fixed worst-case chain), so the passes converge.
func (p pInstr) Resolve(ctx asm.Ctx) (asm.Resolved, error) {
	if ctor, ok := ctors[p.mnem]; ok {
		return ctor(p.ops, ctx)
	}

	forms, err := expandPseudo(p.mnem, p.ops, ctx)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	for _, f := range forms {
		insAddr := ctx.Addr() + uint64(out.Len())
		res, rerr := p.be.ResolveForm(f.mnem, f.ops, formCtx{base: ctx, addr: insAddr})
		if rerr != nil {
			return nil, fmt.Errorf("%s: %w", f.mnem, rerr)
		}

		if _, werr := res.Encode(&out); werr != nil {
			return nil, fmt.Errorf("%s: %w", f.mnem, werr)
		}
	}

	return encoded{b: out.Bytes()}, nil
}

// formCtx is a derived per-form environment: its own address, a
// shared resolver (pseudo expansion: each form is encoded from its
// own address).
type formCtx struct {
	base asm.Ctx
	addr uint64
}

func (c formCtx) Addr() uint64 {
	return c.addr
}

func (c formCtx) Resolve(name string) (uint64, bool) {
	return c.base.Resolve(name)
}

// encoded is an evaluated pseudo-instruction as ready-made bytes
// (expansion).
type encoded struct {
	b []byte
}

func (e encoded) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(e.b)
	return int64(n), err
}
