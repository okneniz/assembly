package loong64

// Shift-immediate families: the .w forms (ui5 shift amounts) and the .d
// forms (ui6). slli/srli/srai/rotri of each width over the shared
// "reg, reg, role" generator.

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	arch "github.com/okneniz/assembly/arch/loong64"
)

// shiftW — the .w shift-immediate family (4 ctors).
var shiftW = []r2RoleEntry[arch.UImm5]{
	{name: "slli.w", ctor: arch.NewSlliW},
	{name: "srli.w", ctor: arch.NewSrliW},
	{name: "srai.w", ctor: arch.NewSraiW},
	{name: "rotri.w", ctor: arch.NewRotriW},
}

// ShiftW — an arbitrary .w shift-immediate instruction.
func ShiftW(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.UImm5]] {
	return newR2RoleGen(rnd, shiftW, UImm5(rnd))
}

// shiftD — the .d shift-immediate family (4 ctors).
var shiftD = []r2RoleEntry[arch.UImm6]{
	{name: "slli.d", ctor: arch.NewSlliD},
	{name: "srli.d", ctor: arch.NewSrliD},
	{name: "srai.d", ctor: arch.NewSraiD},
	{name: "rotri.d", ctor: arch.NewRotriD},
}

// ShiftD — an arbitrary .d shift-immediate instruction.
func ShiftD(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.UImm6]] {
	return newR2RoleGen(rnd, shiftD, UImm6(rnd))
}
