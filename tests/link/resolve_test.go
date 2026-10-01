package linktests

// The resolve-side laws of the generated specimens: every program links
// clean (all the reference directions, the shared .Lloop names, the
// data/bss shapes), its data memory size is exactly the generator's
// bookkeeping, and linking is deterministic - the same sources give the
// same records and the same symbol table twice. Runs anywhere (no
// execution, no writer).

import (
	"bytes"
	"maps"
	mrnd "math/rand/v2"
	"os"
	"strconv"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/link"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
	"github.com/okneniz/assembly/unit"
)

// deps - the arm64 assembler as the link driver's injection.
func deps() link.Deps {
	return link.Deps{Parse: alias.ParseSourceUnit}
}

// flatPlace - fixed bases: the laws here are layout facts.
func flatPlace(textSize, dataSize, dataMem int) (uint64, uint64) {
	return 0x1000, 0x8000
}

func TestLinkProgramResolves(t *testing.T) {
	ohsnap.Check(t, 300, linkarb.NewProgramArb(seedRnd(t)), func(p linkarb.Program) bool {
		sources := sources(p)

		f := link.Link(deps(), sources, "", flatPlace)
		if len(f.Errs) != 0 {
			t.Logf("resolve errors: %v", f.Errs)
			return false
		}

		if f.DataMem != p.DataMem() {
			t.Logf("datamem %d, want %d", f.DataMem, p.DataMem())
			return false
		}

		// determinism: the same sources link to the same program
		again := link.Link(deps(), sources, "", flatPlace)
		textA, errA := f.EncodeText()
		textB, errB := again.EncodeText()
		dataA, derrA := f.EncodeData()
		dataB, derrB := again.EncodeData()
		if errA != nil || errB != nil || derrA != nil || derrB != nil {
			return false
		}

		return bytes.Equal(textA, textB) &&
			bytes.Equal(dataA, dataB) &&
			maps.Equal(f.Syms, again.Syms)
	})
}

// sources adapts the specimen to the driver's input.
func sources(p linkarb.Program) []link.Source {
	ss := p.Sources()
	out := make([]link.Source, len(ss))
	for i := range ss {
		out[i] = link.Source{File: ss[i].Name, Src: ss[i].Src}
	}

	return out
}

// compile-time: the flat policy speaks the unit's shape.
var _ unit.Place = flatPlace

// seedRnd - the deterministic seed of the property suite (ASSEMBLY_SEED,
// default 42), logged so a failure can be reproduced.
func seedRnd(t *testing.T) *mrnd.Rand {
	t.Helper()

	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	t.Logf("seed: %d (ASSEMBLY_SEED)", seed)
	return mrnd.New(mrnd.NewPCG(seed, seed))
}
