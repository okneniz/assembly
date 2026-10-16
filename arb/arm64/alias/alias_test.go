package alias

import (
	"math/rand/v2"
	"testing"

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
			return Cmp(rnd).Generate().Instr()
		}),
		newInstrCase("cmn", func() arm64.Instr {
			return Cmn(rnd).Generate().Instr()
		}),
		newInstrCase("neg", func() arm64.Instr {
			return Neg(rnd).Generate().Instr()
		}),
		newInstrCase("negs", func() arm64.Instr {
			return Negs(rnd).Generate().Instr()
		}),
		newInstrCase("tst", func() arm64.Instr {
			return Tst(rnd).Generate().Instr()
		}),
		newInstrCase("mvn", func() arm64.Instr {
			return Mvn(rnd).Generate().Instr()
		}),
		newInstrCase("mov", func() arm64.Instr {
			return Mov(rnd).Generate().Instr()
		}),
		newInstrCase("mul", func() arm64.Instr {
			return Mul(rnd).Generate().Instr()
		}),
		newInstrCase("mneg", func() arm64.Instr {
			return Mneg(rnd).Generate().Instr()
		}),
		newInstrCase("cset", func() arm64.Instr {
			return Cset(rnd).Generate().Instr()
		}),
		newInstrCase("csetm", func() arm64.Instr {
			return Csetm(rnd).Generate().Instr()
		}),
		newInstrCase("cinc", func() arm64.Instr {
			return Cinc(rnd).Generate().Instr()
		}),
		newInstrCase("cinv", func() arm64.Instr {
			return Cinv(rnd).Generate().Instr()
		}),
		newInstrCase("cneg", func() arm64.Instr {
			return Cneg(rnd).Generate().Instr()
		}),
		newInstrCase("sxtb", func() arm64.Instr {
			return Sxtb(rnd).Generate().Instr()
		}),
		newInstrCase("sxth", func() arm64.Instr {
			return Sxth(rnd).Generate().Instr()
		}),
		newInstrCase("sxtw", func() arm64.Instr {
			return Sxtw(rnd).Generate().Instr()
		}),
		newInstrCase("ubfiz", func() arm64.Instr {
			return Ubfiz(rnd).Generate().Instr()
		}),
		newInstrCase("ubfx", func() arm64.Instr {
			return Ubfx(rnd).Generate().Instr()
		}),
		newInstrCase("sbfiz", func() arm64.Instr {
			return Sbfiz(rnd).Generate().Instr()
		}),
		newInstrCase("sbfx", func() arm64.Instr {
			return Sbfx(rnd).Generate().Instr()
		}),
	}
}
