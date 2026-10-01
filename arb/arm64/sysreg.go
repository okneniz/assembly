package arm64

// The system register name generator (mrs/msr): uniform over the arch
// table; the names are sorted — the map order is random, and both the
// generation and the shrink target must be deterministic. The shrink
// collapses to the alphabetically first name.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
)

// sysregArb — a system register name from the arch table.
type sysregArb struct {
	rnd   *rand.Rand
	names []string
}

func newSysregArb(rnd *rand.Rand) sysregArb {
	return sysregArb{
		rnd:   rnd,
		names: sysregNamesSorted(),
	}
}

// Sysreg — an arbitrary system register name (mrs/msr operand).
func Sysreg(rnd *rand.Rand) ohsnap.Arbitrary[string] {
	return newSysregArb(rnd)
}

func (a sysregArb) Generate() iter.Seq[string] {
	return arb.Stream(func() string {
		return a.names[a.rnd.IntN(len(a.names))]
	})
}

func (a sysregArb) Shrink(name string) iter.Seq[string] {
	if name == a.names[0] {
		return ohsnap.Empty[string]()
	}

	return slices.Values([]string{a.names[0]})
}

// sysregNamesSorted — the arch table names, sorted: the deterministic
// order of the generation and the shrink target.
func sysregNamesSorted() []string {
	table := arm64.SysregNames()
	names := make([]string, 0, len(table))
	for _, name := range table {
		names = append(names, name)
	}

	slices.Sort(names)
	return names
}
