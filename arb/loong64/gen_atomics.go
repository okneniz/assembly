package loong64

// Atomic family: the am* group and sc.q. The LoongArch assembly order is
// "rd, rk, rj" (the @orig_fmt=DKJ note: the second operand packs the K
// field [14:10], the third the J field [9:5]) — the ctor parameter names
// of arch/loong64 already carry it, so the table is a plain R3 table.

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	arch "github.com/okneniz/assembly/arch/loong64"
)

// atomics — the am*/sc.q family in DKJ assembly order (53 ctors).
var atomics = []r3Entry{
	{name: "amcas.b", ctor: arch.NewAmcasB},
	{name: "amcas.h", ctor: arch.NewAmcasH},
	{name: "amcas.w", ctor: arch.NewAmcasW},
	{name: "amcas.d", ctor: arch.NewAmcasD},
	{name: "amcas_db.b", ctor: arch.NewAmcasDbB},
	{name: "amcas_db.h", ctor: arch.NewAmcasDbH},
	{name: "amcas_db.w", ctor: arch.NewAmcasDbW},
	{name: "amcas_db.d", ctor: arch.NewAmcasDbD},
	{name: "amswap.b", ctor: arch.NewAmswapB},
	{name: "amswap.h", ctor: arch.NewAmswapH},
	{name: "amswap.w", ctor: arch.NewAmswapW},
	{name: "amswap.d", ctor: arch.NewAmswapD},
	{name: "amswap_db.b", ctor: arch.NewAmswapDbB},
	{name: "amswap_db.h", ctor: arch.NewAmswapDbH},
	{name: "amswap_db.w", ctor: arch.NewAmswapDbW},
	{name: "amswap_db.d", ctor: arch.NewAmswapDbD},
	{name: "amadd.b", ctor: arch.NewAmaddB},
	{name: "amadd.h", ctor: arch.NewAmaddH},
	{name: "amadd.w", ctor: arch.NewAmaddW},
	{name: "amadd.d", ctor: arch.NewAmaddD},
	{name: "amadd_db.b", ctor: arch.NewAmaddDbB},
	{name: "amadd_db.h", ctor: arch.NewAmaddDbH},
	{name: "amadd_db.w", ctor: arch.NewAmaddDbW},
	{name: "amadd_db.d", ctor: arch.NewAmaddDbD},
	{name: "amand.w", ctor: arch.NewAmandW},
	{name: "amand.d", ctor: arch.NewAmandD},
	{name: "amand_db.w", ctor: arch.NewAmandDbW},
	{name: "amand_db.d", ctor: arch.NewAmandDbD},
	{name: "amor.w", ctor: arch.NewAmorW},
	{name: "amor.d", ctor: arch.NewAmorD},
	{name: "amor_db.w", ctor: arch.NewAmorDbW},
	{name: "amor_db.d", ctor: arch.NewAmorDbD},
	{name: "amxor.w", ctor: arch.NewAmxorW},
	{name: "amxor.d", ctor: arch.NewAmxorD},
	{name: "amxor_db.w", ctor: arch.NewAmxorDbW},
	{name: "amxor_db.d", ctor: arch.NewAmxorDbD},
	{name: "ammax.w", ctor: arch.NewAmmaxW},
	{name: "ammax.d", ctor: arch.NewAmmaxD},
	{name: "ammax.wu", ctor: arch.NewAmmaxWu},
	{name: "ammax.du", ctor: arch.NewAmmaxDu},
	{name: "ammax_db.w", ctor: arch.NewAmmaxDbW},
	{name: "ammax_db.d", ctor: arch.NewAmmaxDbD},
	{name: "ammax_db.wu", ctor: arch.NewAmmaxDbWu},
	{name: "ammax_db.du", ctor: arch.NewAmmaxDbDu},
	{name: "ammin.w", ctor: arch.NewAmminW},
	{name: "ammin.d", ctor: arch.NewAmminD},
	{name: "ammin.wu", ctor: arch.NewAmminWu},
	{name: "ammin.du", ctor: arch.NewAmminDu},
	{name: "ammin_db.w", ctor: arch.NewAmminDbW},
	{name: "ammin_db.d", ctor: arch.NewAmminDbD},
	{name: "ammin_db.wu", ctor: arch.NewAmminDbWu},
	{name: "ammin_db.du", ctor: arch.NewAmminDbDu},
	{name: "sc.q", ctor: arch.NewScQ},
}

// Atomics — an arbitrary am*/sc.q instruction.
func Atomics(rnd *rand.Rand) ohsnap.Arbitrary[R3Params] {
	return newR3Gen(rnd, atomics)
}
