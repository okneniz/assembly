package asm

// The inline-asm fragment: the parsed body of an asm("...") template.
// The contract is instructions and numeric local labels only - the GAS
// 1:/1b/1f discipline (named labels would collide between inlined
// copies; layout directives and literal pools have no meaning inside a
// host program). The compiler-output directives the file path silently
// ignores (.arch_extension fp of the kernel headers, .cfi_*, ...) are
// tolerated the same way - they produce nothing to layout. A Fragment
// is a unit.Sym: the host program resolves it - each instruction at its
// own address, numeric locals against the fragment's own table, every
// other name against the host's symbols (asm("bl foo") calls the host's
// foo).

import (
	"fmt"
	"slices"

	"github.com/okneniz/assembly/unit"
)

// fragStmt is one instruction of the fragment: its parser product, its
// byte offset inside the fragment, its statement index (the numeric
// local coordinate, see fragDef), and its source line (errors).
type fragStmt struct {
	in   Unresolved
	off  int
	stmt int
	line uint
}

// fragDef is one numeric local label definition: its statement index and
// its byte offset (Nb/Nf pick the nearest by statement order - a label
// on the same line as the instruction counts as behind it, see
// resolveLocal).
type fragDef struct {
	stmt int
	off  int
}

// Fragment is the parsed inline-asm body (see ParseFragment): a fixed
// sequence of instructions with fragment-local numeric labels. It
// satisfies unit.Sym.
type Fragment struct {
	stmts []fragStmt
	defs  map[string][]fragDef
	size  int
}

// Resolve encodes the fragment at its final address: each instruction
// resolves at its own pc, numeric locals "Nb"/"Nf" against the
// fragment's table (the nearest by source order), "." is the instruction
// itself, and every other name goes to the host program's resolver.
func (f *Fragment) Resolve(ctx unit.Ctx) ([]unit.Resolved, error) {
	base := ctx.Addr()
	out := make([]unit.Resolved, 0, len(f.stmts))

	for i := range f.stmts {
		st := &f.stmts[i]
		addr := base + uint64(st.off)
		res, err := st.in.Resolve(unit.NewCtx(addr, f.local(ctx, st.stmt, addr)))
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", st.line, err)
		}

		out = append(out, res)
	}

	return out, nil
}

// Size is the byte size of the fragment: the sum of the instruction
// sizes, computed under the placeholder environment at parse time.
func (f *Fragment) Size() int {
	return f.size
}

// local is the resolver of one statement: numeric locals against the
// fragment's definitions, "." the statement's own address, the rest
// delegated to the host program.
func (f *Fragment) local(ctx unit.Ctx, stmt int, addr uint64) func(string) (uint64, bool) {
	return func(name string) (uint64, bool) {
		if isLocalRef(name) {
			return f.localAddr(name, stmt, ctx.Addr())
		}

		if name == "." {
			return addr, true
		}

		return ctx.Resolve(name)
	}
}

// localAddr is "Nb"/"Nf": the nearest definition by statement order - b
// at the statement itself or earlier (a label on the same line precedes
// the instruction), f strictly later.
func (f *Fragment) localAddr(name string, stmt int, base uint64) (uint64, bool) {
	defs := f.defs[name[:len(name)-1]]
	if name[len(name)-1] == 'b' {
		for _, d := range slices.Backward(defs) {
			if d.stmt <= stmt {
				return base + uint64(d.off), true
			}
		}

		return 0, false
	}

	for _, d := range defs {
		if d.stmt > stmt {
			return base + uint64(d.off), true
		}
	}

	return 0, false
}

// ParseFragment parses an inline-asm template body into a Fragment: a
// sequence of instructions with numeric local labels (1:/1b/1f). Named
// labels, layout directives, and literal pool requests (=expr) are
// rejected - a fragment is instructions, the host program owns the rest;
// the compiler-output directives the file path ignores are tolerated. The
// sizes are computed under the placeholder environment, so the layout
// knows the byte count before any symbol exists.
func ParseFragment(src string, be Syntax) (*Fragment, []AsmError) {
	be.ResetOptions()
	stmts := parseSource([]rune(src), be)

	f := &Fragment{defs: map[string][]fragDef{}}
	var errs []AsmError

	for i := range stmts {
		st := &stmts[i]
		if st.err != nil {
			errs = append(errs, *st.err)
			continue
		}

		line := st.pos.Line() + 1

		for _, lbl := range st.labels {
			if isNumericLabel(lbl) {
				f.defs[lbl] = append(f.defs[lbl], newFragDef(i, f.size))
				continue
			}

			errs = append(errs, NewAsmError(line, 0, fmt.Sprintf(
				"named label %q is not allowed in a fragment (numeric locals only)", lbl)))
		}

		if st.directive != nil {
			// the ignored class produces nothing in the file path either;
			// every other directive lays out or names - the host owns that
			if kind, ok := directiveKind(st.directive.name); !ok || kind != argsRestIgnore {
				errs = append(errs, NewAsmError(line, 0, fmt.Sprintf(
					"directive %q is not allowed in a fragment", st.directive.name)))
			}
		}

		if !st.hasInstr {
			continue
		}

		if pu, ok := st.instr.(PoolUser); ok {
			if _, _, want := pu.PoolReq(); want {
				errs = append(errs, NewAsmError(line, 0,
					"literal pools (=expr) are not allowed in a fragment"))
				continue
			}
		}

		size, err := sizeOf(st.instr, unit.NewCtx(0, placeholderResolve(0)))
		if err != nil {
			errs = append(errs, NewAsmError(line, 0, err.Error()))
			continue
		}

		f.stmts = append(f.stmts, newFragStmt(st.instr, f.size, i, line))
		f.size += size
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return f, nil
}

func newFragStmt(in Unresolved, off, stmt int, line uint) fragStmt {
	return fragStmt{in: in, off: off, stmt: stmt, line: line}
}

func newFragDef(stmt, off int) fragDef {
	return fragDef{stmt: stmt, off: off}
}

// compile-time: a Fragment is a deferred record of the unit output.
var _ unit.Sym = (*Fragment)(nil)
