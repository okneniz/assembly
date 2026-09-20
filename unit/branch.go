package unit

// The canonical deferred records: the one place a hole in the stream
// turns into resolved bytes. Producers deposit these with the arch
// formula as an injected function (the builder call over the final
// target and the record's own address) - they never carry resolve logic
// of their own, and the unit never learns an architecture: the ctor
// travels as data, the resolve driver lives here alone.

import (
	"fmt"
)

// NewBranch is one deferred instruction with a label target: a fixed
// size, the builder ctor running at resolve time (b label, j label,
// adr label - one instruction whose operand is the target). src is the
// producer's name of the record - it prefixes the deferred errors.
func NewBranch(
	src, target string,
	size int,
	build func(target, pc uint64) (Resolved, error),
) Sym {
	return branchSym{src: src, target: target, size: size, build: build}
}

// branchSym - the NewBranch record.
type branchSym struct {
	src    string
	target string
	size   int
	build  func(target, pc uint64) (Resolved, error)
}

func (s branchSym) Size() int {
	return s.size
}

func (s branchSym) Resolve(ctx Ctx) ([]Resolved, error) {
	t, ok := ctx.Resolve(s.target)
	if !ok {
		return nil, fmt.Errorf("%s: undefined label %q", s.src, s.target)
	}

	r, err := s.build(t, ctx.Addr())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", s.src, err)
	}

	return []Resolved{r}, nil
}
