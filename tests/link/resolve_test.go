package linktests

// The resolve-side laws of the generated specimens, one subtest per
// architecture: every program links clean (all the reference
// directions, the shared .Lloop names, the data/bss shapes), its data
// memory size is exactly the generator's bookkeeping, and linking is
// deterministic - the same sources give the same records and the same
// symbol table twice. Runs anywhere (no execution, no writer).

import (
	"bytes"
	"maps"
	mrnd "math/rand/v2"
	"os"
	"strconv"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	rpseudo "github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/link"
	linkarb "github.com/okneniz/assembly/tests/link/arb"
	"github.com/okneniz/assembly/unit"
)

// archCase - one architecture of the laws: its dialect, the assembler
// entry injected into the driver, and the single-source shorthand of
// the monolith law.
type archCase struct {
	name         string
	arch         linkarb.Arch
	deps         func() link.Deps
	assembleUnit func(u *unit.Unit, file, src string) []asm.AsmError
}

func TestLinkProgramResolves(t *testing.T) {
	for _, ac := range arches() {
		t.Run(ac.name, func(t *testing.T) {
			check := ohsnap.Check[linkarb.ArchProgram]
			check(
				t,
				5000,
				linkarb.NewArchProgramArb(seedRnd(t), ac.arch),
				func(a linkarb.ArchProgram) bool {
					srcs := sources(a)

					f := link.Link(ac.deps(), srcs, "", flatPlace)
					if len(f.Errs) != 0 {
						t.Logf("resolve errors: %v", f.Errs)
						return false
					}

					if f.DataMem != a.Prog.DataMem() {
						t.Logf("datamem %d, want %d", f.DataMem, a.Prog.DataMem())
						return false
					}

					// determinism: the same sources link to the same program
					again := link.Link(ac.deps(), srcs, "", flatPlace)
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
				},
			)
		})
	}
}

// arches - the full table: every dialect the generator speaks, every
// wiring the library ships.
func arches() []archCase {
	return []archCase{
		{
			name:         "arm64",
			arch:         linkarb.Arm64,
			deps:         func() link.Deps { return link.Deps{Parse: alias.ParseSourceUnit} },
			assembleUnit: alias.AssembleUnit,
		},
		{
			name:         "riscv64",
			arch:         linkarb.Riscv,
			deps:         func() link.Deps { return link.Deps{Parse: rpseudo.ParseSourceUnit} },
			assembleUnit: rpseudo.AssembleUnit,
		},
		{
			name:         "loong64",
			arch:         linkarb.Loong64,
			deps:         func() link.Deps { return link.Deps{Parse: lpseudo.ParseSourceUnit} },
			assembleUnit: lpseudo.AssembleUnit,
		},
	}
}

// flatPlace - fixed bases: the laws here are layout facts.
func flatPlace(textSize, dataSize, dataMem int) (uint64, uint64) {
	return 0x1000, 0x8000
}

// sources adapts the bound specimen to the driver's input.
func sources(a linkarb.ArchProgram) []link.Source {
	ss := a.Prog.Sources(a.Arch)
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

// arm64Deps - the arm64 wiring for the macho-only laws (the image law,
// the clang differential).
func arm64Deps() link.Deps {
	return link.Deps{Parse: alias.ParseSourceUnit}
}
