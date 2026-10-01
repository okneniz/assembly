package arm64

// The barrier shareability domain generator (dsb/dmb): uniform over the
// eight arch domains; the shrink collapses to Sy — the canonical "full
// system" spelling of the minimal counterexample.

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
)

// BarrierDomain — an arbitrary barrier domain; shrinks to arm64.Sy.
func BarrierDomain(rnd *rand.Rand) ohsnap.Arbitrary[arm64.BarrierDomain] {
	return arb.Enum(
		rnd,
		arm64.Sy,
		arm64.St,
		arm64.Ishst,
		arm64.Ish,
		arm64.Nshst,
		arm64.Nsh,
		arm64.Oshst,
		arm64.Osh,
	)
}
