package arm64

// The local aliases of the package-wide helpers (arb.Stream, ohsnap.Empty):
// one definition point for the generators of this package.

import (
	"iter"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
)

// arbStream — a lazy infinite sequence from f.
func arbStream[T any](f func() T) iter.Seq[T] {
	return arb.Stream(f)
}

// ohsnapEmpty — the empty candidate sequence.
func ohsnapEmpty[T any]() iter.Seq[T] {
	return ohsnap.Empty[T]()
}
