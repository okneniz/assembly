package unit

// symDiff - one symbolic data word: the difference of two labels' final
// addresses as width little-endian bytes (see QuadDiff, WordDiff) - the
// label-difference initializer of a static, the .quad a - b form of the
// sources (both labels of the unit, no relocation: two references and a
// subtract at resolve).

import (
	"encoding/binary"
	"fmt"
)

type symDiff struct {
	from  string
	to    string
	width int
}

func newSymDiff(from, to string, width int) symDiff {
	return symDiff{from: from, to: to, width: width}
}

func (d symDiff) Resolve(ctx Ctx) ([]Resolved, error) {
	from, ok := ctx.Resolve(d.from)
	if !ok {
		return nil, fmt.Errorf("diff: undefined label %q", d.from)
	}

	to, ok := ctx.Resolve(d.to)
	if !ok {
		return nil, fmt.Errorf("diff: undefined label %q", d.to)
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, from-to)
	return []Resolved{blob(b[:d.width])}, nil
}

func (d symDiff) Size() int {
	return d.width
}
