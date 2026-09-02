package loong64

// Register-only ALU families: the 3R shapes (arithmetic/logic/shift by
// register, multiply/divide, crc) and the 2R shapes (bit scans, byte
// reversals, extensions). One generator per family over a table of the
// family's arch constructors; the tables carry the mnemonic so the test
// inventory can prove full coverage of arch.Mnemonics().

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// R3Ctor — a 3R constructor: rd, rj, rk (the am* order rd, rk, rj is a
// separate family — the parameter names carry the assembly order).
type R3Ctor func(rd, rj, rk arch.Reg) arch.Instr

// r3Entry — one constructor of a 3R family table.
type r3Entry struct {
	name string
	ctor R3Ctor
}

// R3Params — parameters of any 3R instruction (plus the chosen ctor).
type R3Params struct {
	Rd, Rj, Rk arch.Reg
	Ctor       R3Ctor
}

func NewR3Params(rd, rj, rk arch.Reg, ctor R3Ctor) R3Params {
	return R3Params{
		Rd:   rd,
		Rj:   rj,
		Rk:   rk,
		Ctor: ctor,
	}
}

func (p R3Params) Instr() arch.Instr {
	return p.Ctor(p.Rd, p.Rj, p.Rk)
}

func (p R3Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// r3Gen — generator over a 3R family: a random ctor of the table plus
// three random registers.
type r3Gen struct {
	rnd   *rand.Rand
	ctors []r3Entry
}

func newR3Gen(rnd *rand.Rand, ctors []r3Entry) r3Gen {
	return r3Gen{
		rnd:   rnd,
		ctors: ctors,
	}
}

func (g r3Gen) Generate() iter.Seq[R3Params] {
	return arb.Stream(func() R3Params {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewR3Params(reg(g.rnd), reg(g.rnd), reg(g.rnd), e.ctor)
	})
}

func (g r3Gen) Shrink(p R3Params) iter.Seq[R3Params] {
	rd, rj, rk := regShrunk(p.Rd), regShrunk(p.Rj), regShrunk(p.Rk)
	out := make([]R3Params, 0, len(rd)+len(rj)+len(rk))
	for _, r := range rd {
		out = append(out, NewR3Params(r, p.Rj, p.Rk, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewR3Params(p.Rd, r, p.Rk, p.Ctor))
	}

	for _, r := range rk {
		out = append(out, NewR3Params(p.Rd, p.Rj, r, p.Ctor))
	}

	return slices.Values(out)
}

// alu3R — the 3R arithmetic/logic/shift-register family (46 ctors).
var alu3R = []r3Entry{
	{name: "add.w", ctor: arch.New().AddW},
	{name: "add.d", ctor: arch.New().AddD},
	{name: "sub.w", ctor: arch.New().SubW},
	{name: "sub.d", ctor: arch.New().SubD},
	{name: "slt", ctor: arch.New().Slt},
	{name: "sltu", ctor: arch.New().Sltu},
	{name: "maskeqz", ctor: arch.New().Maskeqz},
	{name: "masknez", ctor: arch.New().Masknez},
	{name: "nor", ctor: arch.New().Nor},
	{name: "and", ctor: arch.New().And},
	{name: "or", ctor: arch.New().Or},
	{name: "xor", ctor: arch.New().Xor},
	{name: "orn", ctor: arch.New().Orn},
	{name: "andn", ctor: arch.New().Andn},
	{name: "sll.w", ctor: arch.New().SllW},
	{name: "srl.w", ctor: arch.New().SrlW},
	{name: "sra.w", ctor: arch.New().SraW},
	{name: "sll.d", ctor: arch.New().SllD},
	{name: "srl.d", ctor: arch.New().SrlD},
	{name: "sra.d", ctor: arch.New().SraD},
	{name: "rotr.w", ctor: arch.New().RotrW},
	{name: "rotr.d", ctor: arch.New().RotrD},
	{name: "mul.w", ctor: arch.New().MulW},
	{name: "mulh.w", ctor: arch.New().MulhW},
	{name: "mulh.wu", ctor: arch.New().MulhWu},
	{name: "mul.d", ctor: arch.New().MulD},
	{name: "mulh.d", ctor: arch.New().MulhD},
	{name: "mulh.du", ctor: arch.New().MulhDu},
	{name: "mulw.d.w", ctor: arch.New().MulwDW},
	{name: "mulw.d.wu", ctor: arch.New().MulwDWu},
	{name: "div.w", ctor: arch.New().DivW},
	{name: "mod.w", ctor: arch.New().ModW},
	{name: "div.wu", ctor: arch.New().DivWu},
	{name: "mod.wu", ctor: arch.New().ModWu},
	{name: "div.d", ctor: arch.New().DivD},
	{name: "mod.d", ctor: arch.New().ModD},
	{name: "div.du", ctor: arch.New().DivDu},
	{name: "mod.du", ctor: arch.New().ModDu},
	{name: "crc.w.b.w", ctor: arch.New().CrcWBW},
	{name: "crc.w.h.w", ctor: arch.New().CrcWHW},
	{name: "crc.w.w.w", ctor: arch.New().CrcWWW},
	{name: "crc.w.d.w", ctor: arch.New().CrcWDW},
	{name: "crcc.w.b.w", ctor: arch.New().CrccWBW},
	{name: "crcc.w.h.w", ctor: arch.New().CrccWHW},
	{name: "crcc.w.w.w", ctor: arch.New().CrccWWW},
	{name: "crcc.w.d.w", ctor: arch.New().CrccWDW},
}

// Alu3R — an arbitrary 3R arithmetic/logic/mul/div/crc instruction.
func Alu3R(rnd *rand.Rand) ohsnap.Arbitrary[R3Params] {
	return newR3Gen(rnd, alu3R)
}

// R2Ctor — a 2R constructor: rd, rj (asrtle/asrtgt: rj, rk).
type R2Ctor func(rd, rj arch.Reg) arch.Instr

// r2Entry — one constructor of a 2R family table.
type r2Entry struct {
	name string
	ctor R2Ctor
}

// R2Params — parameters of any 2R instruction (plus the chosen ctor).
type R2Params struct {
	Rd, Rj arch.Reg
	Ctor   R2Ctor
}

func NewR2Params(rd, rj arch.Reg, ctor R2Ctor) R2Params {
	return R2Params{
		Rd:   rd,
		Rj:   rj,
		Ctor: ctor,
	}
}

func (p R2Params) Instr() arch.Instr {
	return p.Ctor(p.Rd, p.Rj)
}

func (p R2Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// r2Gen — generator over a 2R family: a random ctor of the table plus
// two random registers.
type r2Gen struct {
	rnd   *rand.Rand
	ctors []r2Entry
}

func newR2Gen(rnd *rand.Rand, ctors []r2Entry) r2Gen {
	return r2Gen{
		rnd:   rnd,
		ctors: ctors,
	}
}

func (g r2Gen) Generate() iter.Seq[R2Params] {
	return arb.Stream(func() R2Params {
		e := g.ctors[g.rnd.IntN(len(g.ctors))]
		return NewR2Params(reg(g.rnd), reg(g.rnd), e.ctor)
	})
}

func (g r2Gen) Shrink(p R2Params) iter.Seq[R2Params] {
	rd, rj := regShrunk(p.Rd), regShrunk(p.Rj)
	out := make([]R2Params, 0, len(rd)+len(rj))
	for _, r := range rd {
		out = append(out, NewR2Params(r, p.Rj, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewR2Params(p.Rd, r, p.Ctor))
	}

	return slices.Values(out)
}

// alu2R — the 2R bit-manipulation family (26 ctors).
var alu2R = []r2Entry{
	{name: "ext.w.b", ctor: arch.New().ExtWB},
	{name: "ext.w.h", ctor: arch.New().ExtWH},
	{name: "rdtimel.w", ctor: arch.New().RdtimelW},
	{name: "rdtimeh.w", ctor: arch.New().RdtimehW},
	{name: "rdtime.d", ctor: arch.New().RdtimeD},
	{name: "cpucfg", ctor: arch.New().Cpucfg},
	{name: "clo.w", ctor: arch.New().CloW},
	{name: "clz.w", ctor: arch.New().ClzW},
	{name: "cto.w", ctor: arch.New().CtoW},
	{name: "ctz.w", ctor: arch.New().CtzW},
	{name: "clo.d", ctor: arch.New().CloD},
	{name: "clz.d", ctor: arch.New().ClzD},
	{name: "cto.d", ctor: arch.New().CtoD},
	{name: "ctz.d", ctor: arch.New().CtzD},
	{name: "revb.2h", ctor: arch.New().Revb2H},
	{name: "revb.4h", ctor: arch.New().Revb4H},
	{name: "revb.2w", ctor: arch.New().Revb2W},
	{name: "revb.d", ctor: arch.New().RevbD},
	{name: "revh.2w", ctor: arch.New().Revh2W},
	{name: "revh.d", ctor: arch.New().RevhD},
	{name: "bitrev.4b", ctor: arch.New().Revbit4B},
	{name: "bitrev.w", ctor: arch.New().RevbitW},
	{name: "bitrev.8b", ctor: arch.New().Revbit8B},
	{name: "bitrev.d", ctor: arch.New().RevbitD},
	{name: "asrtle.d", ctor: arch.New().AsrtleD},
	{name: "asrtgt.d", ctor: arch.New().AsrtgtD},
}

// Alu2R — an arbitrary 2R bit-manipulation instruction.
func Alu2R(rnd *rand.Rand) ohsnap.Arbitrary[R2Params] {
	return newR2Gen(rnd, alu2R)
}
