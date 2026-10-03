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

// NewPair is one deferred fixed pair with a label target: both
// instructions resolve together against the pair's own address (the la
// address pairs - adrp+add, auipc+addi, pcalau12i+addi.d).
func NewPair(
	src, target string,
	size int,
	build func(target, pc uint64) ([]Resolved, error),
) Sym {
	return pairSym{src: src, target: target, size: size, build: build}
}

// pairSym - the NewPair record.
type pairSym struct {
	src    string
	target string
	size   int
	build  func(target, pc uint64) ([]Resolved, error)
}

func (s pairSym) Resolve(ctx Ctx) ([]Resolved, error) {
	t, ok := ctx.Resolve(s.target)
	if !ok {
		return nil, fmt.Errorf("%s: undefined label %q", s.src, s.target)
	}

	rs, err := s.build(t, ctx.Addr())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", s.src, err)
	}

	return rs, nil
}

func (s pairSym) Size() int {
	return s.size
}
