package alias

import (
	"math/rand/v2"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// instrCase — one alias family: its mnemonic + generator.
type instrCase struct {
	mnem string
	gen  func() arm64.Instr
}

// TestAliasGenValid — a generated alias text always assembles into an
// instruction (the generator produces no invalid combinations).
func TestAliasGenValid(t *testing.T) {
	rnd := arb.Rnd(42)
	for _, c := range instrCases(rnd) {
		for range 300 {
			in := c.gen()
			require.NotNil(t, in, "%s: the text did not assemble", c.mnem)
		}
	}
}

func newInstrCase(mnem string, gen func() arm64.Instr) instrCase {
	return instrCase{
		mnem: mnem,
		gen:  gen,
	}
}

// instrCases — all alias families (mirroring the aliasCtors table of
// asm/arm64/alias; tst imm waits for the structural bitmask generator).
func instrCases(rnd *rand.Rand) []instrCase {
	return []instrCase{
		newInstrCase("cmp", func() arm64.Instr {
			return ohsnap.First(Cmp(rnd).Generate()).Instr()
		}),
		newInstrCase("cmn", func() arm64.Instr {
			return ohsnap.First(Cmn(rnd).Generate()).Instr()
		}),
		newInstrCase("neg", func() arm64.Instr {
			return ohsnap.First(Neg(rnd).Generate()).Instr()
		}),
		newInstrCase("negs", func() arm64.Instr {
			return ohsnap.First(Negs(rnd).Generate()).Instr()
		}),
		newInstrCase("tst", func() arm64.Instr {
			return ohsnap.First(Tst(rnd).Generate()).Instr()
		}),
		newInstrCase("mvn", func() arm64.Instr {
			return ohsnap.First(Mvn(rnd).Generate()).Instr()
		}),
		newInstrCase("mov", func() arm64.Instr {
			return ohsnap.First(Mov(rnd).Generate()).Instr()
		}),
		newInstrCase("mul", func() arm64.Instr {
			return ohsnap.First(Mul(rnd).Generate()).Instr()
		}),
		newInstrCase("mneg", func() arm64.Instr {
			return ohsnap.First(Mneg(rnd).Generate()).Instr()
		}),
		newInstrCase("cset", func() arm64.Instr {
			return ohsnap.First(Cset(rnd).Generate()).Instr()
		}),
		newInstrCase("csetm", func() arm64.Instr {
			return ohsnap.First(Csetm(rnd).Generate()).Instr()
		}),
		newInstrCase("cinc", func() arm64.Instr {
			return ohsnap.First(Cinc(rnd).Generate()).Instr()
		}),
		newInstrCase("cinv", func() arm64.Instr {
			return ohsnap.First(Cinv(rnd).Generate()).Instr()
		}),
		newInstrCase("cneg", func() arm64.Instr {
			return ohsnap.First(Cneg(rnd).Generate()).Instr()
		}),
		newInstrCase("sxtb", func() arm64.Instr {
			return ohsnap.First(Sxtb(rnd).Generate()).Instr()
		}),
		newInstrCase("sxth", func() arm64.Instr {
			return ohsnap.First(Sxth(rnd).Generate()).Instr()
		}),
		newInstrCase("sxtw", func() arm64.Instr {
			return ohsnap.First(Sxtw(rnd).Generate()).Instr()
		}),
		newInstrCase("ubfiz", func() arm64.Instr {
			return ohsnap.First(Ubfiz(rnd).Generate()).Instr()
		}),
		newInstrCase("ubfx", func() arm64.Instr {
			return ohsnap.First(Ubfx(rnd).Generate()).Instr()
		}),
		newInstrCase("sbfiz", func() arm64.Instr {
			return ohsnap.First(Sbfiz(rnd).Generate()).Instr()
		}),
		newInstrCase("sbfx", func() arm64.Instr {
			return ohsnap.First(Sbfx(rnd).Generate()).Instr()
		}),
	}
}
