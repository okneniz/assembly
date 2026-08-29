package loong64

// Register-only ALU families: the 3R shapes (arithmetic/logic/shift by
// register, multiply/divide, crc) and the 2R shapes (bit scans, byte
// reversals, extensions). One generator per family over a table of the
// family's arch constructors; the tables carry the mnemonic so the test
// inventory can prove full coverage of arch.Mnemonics().

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

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

func (g r3Gen) Generate() R3Params {
	e := g.ctors[g.rnd.IntN(len(g.ctors))]
	return NewR3Params(reg(g.rnd), reg(g.rnd), reg(g.rnd), e.ctor)
}

func (g r3Gen) Shrink(p R3Params) []R3Params {
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

	return out
}

// alu3R — the 3R arithmetic/logic/shift-register family (46 ctors).
var alu3R = []r3Entry{
	{name: "add.w", ctor: arch.NewAddW},
	{name: "add.d", ctor: arch.NewAddD},
	{name: "sub.w", ctor: arch.NewSubW},
	{name: "sub.d", ctor: arch.NewSubD},
	{name: "slt", ctor: arch.NewSlt},
	{name: "sltu", ctor: arch.NewSltu},
	{name: "maskeqz", ctor: arch.NewMaskeqz},
	{name: "masknez", ctor: arch.NewMasknez},
	{name: "nor", ctor: arch.NewNor},
	{name: "and", ctor: arch.NewAnd},
	{name: "or", ctor: arch.NewOr},
	{name: "xor", ctor: arch.NewXor},
	{name: "orn", ctor: arch.NewOrn},
	{name: "andn", ctor: arch.NewAndn},
	{name: "sll.w", ctor: arch.NewSllW},
	{name: "srl.w", ctor: arch.NewSrlW},
	{name: "sra.w", ctor: arch.NewSraW},
	{name: "sll.d", ctor: arch.NewSllD},
	{name: "srl.d", ctor: arch.NewSrlD},
	{name: "sra.d", ctor: arch.NewSraD},
	{name: "rotr.w", ctor: arch.NewRotrW},
	{name: "rotr.d", ctor: arch.NewRotrD},
	{name: "mul.w", ctor: arch.NewMulW},
	{name: "mulh.w", ctor: arch.NewMulhW},
	{name: "mulh.wu", ctor: arch.NewMulhWu},
	{name: "mul.d", ctor: arch.NewMulD},
	{name: "mulh.d", ctor: arch.NewMulhD},
	{name: "mulh.du", ctor: arch.NewMulhDu},
	{name: "mulw.d.w", ctor: arch.NewMulwDW},
	{name: "mulw.d.wu", ctor: arch.NewMulwDWu},
	{name: "div.w", ctor: arch.NewDivW},
	{name: "mod.w", ctor: arch.NewModW},
	{name: "div.wu", ctor: arch.NewDivWu},
	{name: "mod.wu", ctor: arch.NewModWu},
	{name: "div.d", ctor: arch.NewDivD},
	{name: "mod.d", ctor: arch.NewModD},
	{name: "div.du", ctor: arch.NewDivDu},
	{name: "mod.du", ctor: arch.NewModDu},
	{name: "crc.w.b.w", ctor: arch.NewCrcWBW},
	{name: "crc.w.h.w", ctor: arch.NewCrcWHW},
	{name: "crc.w.w.w", ctor: arch.NewCrcWWW},
	{name: "crc.w.d.w", ctor: arch.NewCrcWDW},
	{name: "crcc.w.b.w", ctor: arch.NewCrccWBW},
	{name: "crcc.w.h.w", ctor: arch.NewCrccWHW},
	{name: "crcc.w.w.w", ctor: arch.NewCrccWWW},
	{name: "crcc.w.d.w", ctor: arch.NewCrccWDW},
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

func (g r2Gen) Generate() R2Params {
	e := g.ctors[g.rnd.IntN(len(g.ctors))]
	return NewR2Params(reg(g.rnd), reg(g.rnd), e.ctor)
}

func (g r2Gen) Shrink(p R2Params) []R2Params {
	rd, rj := regShrunk(p.Rd), regShrunk(p.Rj)
	out := make([]R2Params, 0, len(rd)+len(rj))
	for _, r := range rd {
		out = append(out, NewR2Params(r, p.Rj, p.Ctor))
	}

	for _, r := range rj {
		out = append(out, NewR2Params(p.Rd, r, p.Ctor))
	}

	return out
}

// alu2R — the 2R bit-manipulation family (26 ctors).
var alu2R = []r2Entry{
	{name: "ext.w.b", ctor: arch.NewExtWB},
	{name: "ext.w.h", ctor: arch.NewExtWH},
	{name: "rdtimel.w", ctor: arch.NewRdtimelW},
	{name: "rdtimeh.w", ctor: arch.NewRdtimehW},
	{name: "rdtime.d", ctor: arch.NewRdtimeD},
	{name: "cpucfg", ctor: arch.NewCpucfg},
	{name: "clo.w", ctor: arch.NewCloW},
	{name: "clz.w", ctor: arch.NewClzW},
	{name: "cto.w", ctor: arch.NewCtoW},
	{name: "ctz.w", ctor: arch.NewCtzW},
	{name: "clo.d", ctor: arch.NewCloD},
	{name: "clz.d", ctor: arch.NewClzD},
	{name: "cto.d", ctor: arch.NewCtoD},
	{name: "ctz.d", ctor: arch.NewCtzD},
	{name: "revb.2h", ctor: arch.NewRevb2H},
	{name: "revb.4h", ctor: arch.NewRevb4H},
	{name: "revb.2w", ctor: arch.NewRevb2W},
	{name: "revb.d", ctor: arch.NewRevbD},
	{name: "revh.2w", ctor: arch.NewRevh2W},
	{name: "revh.d", ctor: arch.NewRevhD},
	{name: "bitrev.4b", ctor: arch.NewRevbit4B},
	{name: "bitrev.w", ctor: arch.NewRevbitW},
	{name: "bitrev.8b", ctor: arch.NewRevbit8B},
	{name: "bitrev.d", ctor: arch.NewRevbitD},
	{name: "asrtle.d", ctor: arch.NewAsrtleD},
	{name: "asrtgt.d", ctor: arch.NewAsrtgtD},
}

// Alu2R — an arbitrary 2R bit-manipulation instruction.
func Alu2R(rnd *rand.Rand) ohsnap.Arbitrary[R2Params] {
	return newR2Gen(rnd, alu2R)
}
