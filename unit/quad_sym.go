package unit

// quadSym - one symbolic data word: the final address of a label as 8
// little-endian bytes (see QuadSym) - the address initializer of a
// static (a Quad whose value only the layout knows).

import (
	"encoding/binary"
	"fmt"
)

type quadSym struct {
	name string
}

func newQuadSym(name string) quadSym {
	return quadSym{name: name}
}

func (q quadSym) Resolve(ctx Ctx) ([]Resolved, error) {
	v, ok := ctx.Resolve(q.name)
	if !ok {
		return nil, fmt.Errorf("quad: undefined label %q", q.name)
	}

	return []Resolved{blob(binary.LittleEndian.AppendUint64(nil, v))}, nil
}

func (q quadSym) Size() int {
	return 8
}
