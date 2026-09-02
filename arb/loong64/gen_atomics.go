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
	{name: "amcas.b", ctor: arch.New().AmcasB},
	{name: "amcas.h", ctor: arch.New().AmcasH},
	{name: "amcas.w", ctor: arch.New().AmcasW},
	{name: "amcas.d", ctor: arch.New().AmcasD},
	{name: "amcas_db.b", ctor: arch.New().AmcasDbB},
	{name: "amcas_db.h", ctor: arch.New().AmcasDbH},
	{name: "amcas_db.w", ctor: arch.New().AmcasDbW},
	{name: "amcas_db.d", ctor: arch.New().AmcasDbD},
	{name: "amswap.b", ctor: arch.New().AmswapB},
	{name: "amswap.h", ctor: arch.New().AmswapH},
	{name: "amswap.w", ctor: arch.New().AmswapW},
	{name: "amswap.d", ctor: arch.New().AmswapD},
	{name: "amswap_db.b", ctor: arch.New().AmswapDbB},
	{name: "amswap_db.h", ctor: arch.New().AmswapDbH},
	{name: "amswap_db.w", ctor: arch.New().AmswapDbW},
	{name: "amswap_db.d", ctor: arch.New().AmswapDbD},
	{name: "amadd.b", ctor: arch.New().AmaddB},
	{name: "amadd.h", ctor: arch.New().AmaddH},
	{name: "amadd.w", ctor: arch.New().AmaddW},
	{name: "amadd.d", ctor: arch.New().AmaddD},
	{name: "amadd_db.b", ctor: arch.New().AmaddDbB},
	{name: "amadd_db.h", ctor: arch.New().AmaddDbH},
	{name: "amadd_db.w", ctor: arch.New().AmaddDbW},
	{name: "amadd_db.d", ctor: arch.New().AmaddDbD},
	{name: "amand.w", ctor: arch.New().AmandW},
	{name: "amand.d", ctor: arch.New().AmandD},
	{name: "amand_db.w", ctor: arch.New().AmandDbW},
	{name: "amand_db.d", ctor: arch.New().AmandDbD},
	{name: "amor.w", ctor: arch.New().AmorW},
	{name: "amor.d", ctor: arch.New().AmorD},
	{name: "amor_db.w", ctor: arch.New().AmorDbW},
	{name: "amor_db.d", ctor: arch.New().AmorDbD},
	{name: "amxor.w", ctor: arch.New().AmxorW},
	{name: "amxor.d", ctor: arch.New().AmxorD},
	{name: "amxor_db.w", ctor: arch.New().AmxorDbW},
	{name: "amxor_db.d", ctor: arch.New().AmxorDbD},
	{name: "ammax.w", ctor: arch.New().AmmaxW},
	{name: "ammax.d", ctor: arch.New().AmmaxD},
	{name: "ammax.wu", ctor: arch.New().AmmaxWu},
	{name: "ammax.du", ctor: arch.New().AmmaxDu},
	{name: "ammax_db.w", ctor: arch.New().AmmaxDbW},
	{name: "ammax_db.d", ctor: arch.New().AmmaxDbD},
	{name: "ammax_db.wu", ctor: arch.New().AmmaxDbWu},
	{name: "ammax_db.du", ctor: arch.New().AmmaxDbDu},
	{name: "ammin.w", ctor: arch.New().AmminW},
	{name: "ammin.d", ctor: arch.New().AmminD},
	{name: "ammin.wu", ctor: arch.New().AmminWu},
	{name: "ammin.du", ctor: arch.New().AmminDu},
	{name: "ammin_db.w", ctor: arch.New().AmminDbW},
	{name: "ammin_db.d", ctor: arch.New().AmminDbD},
	{name: "ammin_db.wu", ctor: arch.New().AmminDbWu},
	{name: "ammin_db.du", ctor: arch.New().AmminDbDu},
	{name: "sc.q", ctor: arch.New().ScQ},
}

// Atomics — an arbitrary am*/sc.q instruction.
func Atomics(rnd *rand.Rand) ohsnap.Arbitrary[R3Params] {
	return newR3Gen(rnd, atomics)
}
