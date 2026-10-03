package riscv

// The shared core of the resolve-level pseudo families (la/call/tail).
// The parameters hold the TARGET OFFSET from the assembly base: the
// laws render the absolute text at the property base and check that
// the decoded pair lands exactly on the target (String prints the
// offset form - the diagnostics face; the absolute text needs the
// base, which lives in the test). The per-family generators live in
// their own files (la.go, call.go, tail.go).

import (
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
)

// LaParams — parameters of la rd, target.
type LaParams struct {
	Rd  riscv.Reg
	Off int64 // the target offset from the assembly base
}

func NewLaParams(rd riscv.Reg, off int64) LaParams {
	return LaParams{
		Rd:  rd,
		Off: off,
	}
}

func (p LaParams) String() string {
	return fmt.Sprintf("la %s, +%#x", p.Rd, p.Off)
}

// CallParams — parameters of call target.
type CallParams struct {
	Off int64 // the target offset from the assembly base
}

func NewCallParams(off int64) CallParams {
	return CallParams{
		Off: off,
	}
}

func (p CallParams) String() string {
	return fmt.Sprintf("call +%#x", p.Off)
}

// TailParams — parameters of tail target.
type TailParams struct {
	Off int64 // the target offset from the assembly base
}

func NewTailParams(off int64) TailParams {
	return TailParams{
		Off: off,
	}
}

func (p TailParams) String() string {
	return fmt.Sprintf("tail +%#x", p.Off)
}

// The pcrel window of an auipc+addi pair: hi is the SIGNED 20-bit
// field, lo the signed 12-bit one — the reachable offsets.
const (
	pcrelMin = -(1 << 31) - 2048 // hi = -0x80000, lo = -2048
	pcrelMax = 1<<31 - 2049      // hi = 0x7ffff, lo = 2047
)

// genOff — an offset of the pcrel window (the ±2GB an auipc+addi pair
// reaches; tail's ±1MB law is the family's own).
func genOff(rnd *rand.Rand) int64 {
	switch rnd.IntN(3) {
	case 0: // small deltas
		return int64(rnd.IntN(1<<20)) - 1<<19
	case 1: // the whole pcrel window
		return rnd.Int64N(pcrelMax-pcrelMin+1) + pcrelMin
	default: // the window edges
		return []int64{
			pcrelMax, pcrelMin, 1<<31 - 4096, -(1 << 31),
			4095, -4096, 2047, -2048, 0, 1,
		}[rnd.IntN(10)]
	}
}

// genTailOff — an offset of the jal window (±1MB, even; the bounds
// are the jalOff constants of the jal family).
func genTailOff(rnd *rand.Rand) int64 {
	switch rnd.IntN(3) {
	case 0: // small deltas
		v := int64(rnd.IntN(1<<12)) - 1<<11
		if v%2 != 0 {
			v++
		}

		return v
	case 1: // the whole jal window
		return jalOffMin + 2*rnd.Int64N((jalOffMax-jalOffMin)/2+1)
	default: // the edges
		return []int64{
			jalOffMax, jalOffMin, 1<<20 - 4, -(1 << 20) + 4,
			4094, -4096, 2046, -2048, 0,
		}[rnd.IntN(9)]
	}
}

// offShrunk — the shrink axis of an offset: the boundaries of the
// reachable window first, then halving toward zero.
func offShrunk(v, from, to int64) []int64 {
	var out []int64
	for d := range immShrink(from, to)(v) {
		out = append(out, d)
	}

	return out
}

// laGen — generator for la: any register, the pcrel window.
type laGen struct {
	rnd *rand.Rand
}

func newLaGen(rnd *rand.Rand) laGen {
	return laGen{rnd: rnd}
}

func (g laGen) Generate() iter.Seq[LaParams] {
	return arb.Stream(func() LaParams {
		return NewLaParams(reg(g.rnd), genOff(g.rnd))
	})
}

func (g laGen) Shrink(p LaParams) iter.Seq[LaParams] {
	out := make([]LaParams, 0, 12)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewLaParams(r, p.Off))
	}

	for _, v := range offShrunk(p.Off, pcrelMin, pcrelMax) {
		out = append(out, NewLaParams(p.Rd, v))
	}

	return slices.Values(out)
}

// callGen — generator for call: the pcrel window (rd is always ra).
type callGen struct {
	rnd *rand.Rand
}

func newCallGen(rnd *rand.Rand) callGen {
	return callGen{rnd: rnd}
}

func (g callGen) Generate() iter.Seq[CallParams] {
	return arb.Stream(func() CallParams {
		return NewCallParams(genOff(g.rnd))
	})
}

func (g callGen) Shrink(p CallParams) iter.Seq[CallParams] {
	out := make([]CallParams, 0, 8)
	for _, v := range offShrunk(p.Off, pcrelMin, pcrelMax) {
		out = append(out, NewCallParams(v))
	}

	return slices.Values(out)
}

// tailGen — generator for tail: the jal window.
type tailGen struct {
	rnd *rand.Rand
}

func newTailGen(rnd *rand.Rand) tailGen {
	return tailGen{rnd: rnd}
}

func (g tailGen) Generate() iter.Seq[TailParams] {
	return arb.Stream(func() TailParams {
		return NewTailParams(genTailOff(g.rnd))
	})
}

func (g tailGen) Shrink(p TailParams) iter.Seq[TailParams] {
	out := make([]TailParams, 0, 8)
	for _, v := range pcOffShrunk(p.Off, jalOffMin, jalOffMax) {
		out = append(out, NewTailParams(v))
	}

	return slices.Values(out)
}
