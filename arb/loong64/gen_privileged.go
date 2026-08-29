package loong64

// Privileged families: the CSR triple (csrrd/csrwr "rd, ui14", csrxchg
// "rd, rj, ui14"), the iocsr 2R group, the page-table walk forms
// (lddir "rd, rj, ui8", ldpte "rj, ui8") and invtlb. The operandless
// privileged forms (tlbsrch and friends) have NO generator: they are a
// single fixed word each, exercised through the decodeTable-driven base
// words of the differential corpus (EmptyForms below).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	arch "github.com/okneniz/assembly/arch/loong64"
)

// csrRW — the "reg, ui14" CSR family: csrrd/csrwr (2 ctors).
var csrRW = []r1RoleEntry[arch.UImm14]{
	{name: "csrrd", ctor: arch.NewCsrrd},
	{name: "csrwr", ctor: arch.NewCsrwr},
}

// CsrRW — an arbitrary csrrd/csrwr instruction.
func CsrRW(rnd *rand.Rand) ohsnap.Arbitrary[R1RoleParams[arch.UImm14]] {
	return newR1RoleGen(rnd, csrRW, UImm14(rnd))
}

// csrXchg — the "reg, reg, ui14" CSR family: csrxchg (1 ctor).
var csrXchg = []r2RoleEntry[arch.UImm14]{
	{name: "csrxchg", ctor: arch.NewCsrxchg},
}

// CsrXchg — an arbitrary csrxchg instruction.
func CsrXchg(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.UImm14]] {
	return newR2RoleGen(rnd, csrXchg, UImm14(rnd))
}

// ioCsr — the iocsr 2R family (8 ctors).
var ioCsr = []r2Entry{
	{name: "iocsrrd.b", ctor: arch.NewIocsrrdB},
	{name: "iocsrrd.h", ctor: arch.NewIocsrrdH},
	{name: "iocsrrd.w", ctor: arch.NewIocsrrdW},
	{name: "iocsrrd.d", ctor: arch.NewIocsrrdD},
	{name: "iocsrwr.b", ctor: arch.NewIocsrwrB},
	{name: "iocsrwr.h", ctor: arch.NewIocsrwrH},
	{name: "iocsrwr.w", ctor: arch.NewIocsrwrW},
	{name: "iocsrwr.d", ctor: arch.NewIocsrwrD},
}

// IoCsr — an arbitrary iocsrrd/iocsrwr instruction.
func IoCsr(rnd *rand.Rand) ohsnap.Arbitrary[R2Params] {
	return newR2Gen(rnd, ioCsr)
}

// lddir — the page-directory walk family: lddir "rd, rj, ui8" (1 ctor).
var lddir = []r2RoleEntry[arch.UImm8]{
	{name: "lddir", ctor: arch.NewLddir},
}

// Lddir — an arbitrary lddir instruction.
func Lddir(rnd *rand.Rand) ohsnap.Arbitrary[R2RoleParams[arch.UImm8]] {
	return newR2RoleGen(rnd, lddir, UImm8(rnd))
}

// ldpte — the page-table entry family: ldpte "rj, ui8" (1 ctor).
var ldpte = []r1RoleEntry[arch.UImm8]{
	{name: "ldpte", ctor: arch.NewLdpte},
}

// Ldpte — an arbitrary ldpte instruction.
func Ldpte(rnd *rand.Rand) ohsnap.Arbitrary[R1RoleParams[arch.UImm8]] {
	return newR1RoleGen(rnd, ldpte, UImm8(rnd))
}

// invtlb — the TLB invalidate family over the "ui5, reg, reg" shape
// (1 ctor).
var invtlb = []u5rrentry{
	{name: "invtlb", ctor: arch.NewInvtlb},
}

// Invtlb — an arbitrary invtlb instruction.
func Invtlb(rnd *rand.Rand) ohsnap.Arbitrary[U5RRParams] {
	return newU5rrGen(rnd, invtlb)
}

// emptyForms — the operandless privileged forms: no generator family (a
// single fixed word each, nothing to generate); the differential corpus
// covers them through arch.EncodingWord(mnem) base words.
var emptyForms = []string{
	"tlbclr",
	"tlbflush",
	"tlbsrch",
	"tlbrd",
	"tlbwr",
	"tlbfill",
	"ertn",
}

// EmptyForms — the mnemonics of the operandless forms (a copy of the
// internal table; the inventory test checks the family tables against
// arch.Mnemonics() minus exactly these).
func EmptyForms() []string {
	out := make([]string, len(emptyForms))
	copy(out, emptyForms)
	return out
}
