package arm64

// The SIMD copy family: the lane-addressed moves over the v file
// (ins/ins element/smov/umov/dup element/dup scalar). The lane index
// rides imm5 above the size one-hot; its bound is the lane count of
// the full vector register by the element width.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// InsParams — parameters of ins (general).
type InsParams struct {
	Vd   arm64.VReg
	Idx  uint32
	Wn   arm64.Reg
	Elem string
}

func NewInsParams(vd arm64.VReg, idx uint32, wn arm64.Reg, elem string) InsParams {
	return InsParams{
		Vd:   vd,
		Idx:  idx,
		Wn:   wn,
		Elem: elem,
	}
}

func (p InsParams) Instr() arm64.Instr {
	in, err := arm64.New().Ins(p.Vd, p.Idx, p.Wn, p.Elem)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p InsParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// insArb — the ins core: the gpr width follows the element (.d takes
// an x register).
type insArb struct {
	rnd *rand.Rand
}

// Ins — an arbitrary ins (general).
func Ins(rnd *rand.Rand) ohsnap.Arbitrary[InsParams] {
	return insArb{rnd: rnd}
}

func (a insArb) Generate() iter.Seq[InsParams] {
	return arbStream(func() InsParams {
		elem := genElem(a.rnd, []string{"b", "h", "s", "d"})
		return NewInsParams(
			genV(a.rnd),
			uint32(a.rnd.IntN(int(elemLanes(elemWidth(elem))))),
			genWidthReg(a.rnd, elem == "d"),
			elem,
		)
	})
}

func (a insArb) Shrink(p InsParams) iter.Seq[InsParams] {
	out := make([]InsParams, 0, 8)
	for _, r := range vShrunk(p.Vd) {
		out = append(out, NewInsParams(r, p.Idx, p.Wn, p.Elem))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewInsParams(p.Vd, i, p.Wn, p.Elem))
	}

	for _, r := range regShrunk(p.Wn) {
		out = append(out, NewInsParams(p.Vd, p.Idx, r, p.Elem))
	}

	for _, e := range elemShrunk(p.Elem) {
		if (e == "d") != (p.Elem == "d") {
			continue // the gpr width stays tied to the element
		}

		out = append(out, NewInsParams(p.Vd, p.Idx, p.Wn, e))
	}

	return slices.Values(out)
}

// --- ins (element): ins.sz vd[idx], vn[src] ---------------------------

// InsElemParams — parameters of ins (element).
type InsElemParams struct {
	Rd     arm64.VReg
	Rn     arm64.VReg
	Elem   string
	Idx    uint32
	SrcIdx uint32
}

func NewInsElemParams(
	rd arm64.VReg,
	rn arm64.VReg,
	elem string,
	idx uint32,
	srcIdx uint32,
) InsElemParams {
	return InsElemParams{
		Rd:     rd,
		Rn:     rn,
		Elem:   elem,
		Idx:    idx,
		SrcIdx: srcIdx,
	}
}

func (p InsElemParams) Instr() arm64.Instr {
	in, err := arm64.New().InsElem(p.Rd, p.Rn, p.Elem, p.Idx, p.SrcIdx)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p InsElemParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// InsElem — an arbitrary ins (element).
func InsElem(rnd *rand.Rand) ohsnap.Arbitrary[InsElemParams] {
	return insElemArb{rnd: rnd}
}

type insElemArb struct {
	rnd *rand.Rand
}

func (a insElemArb) Generate() iter.Seq[InsElemParams] {
	return arbStream(func() InsElemParams {
		elem := genElem(a.rnd, []string{"b", "h", "s", "d"})
		lanes := int(elemLanes(elemWidth(elem)))
		return NewInsElemParams(
			genV(a.rnd),
			genV(a.rnd),
			elem,
			uint32(a.rnd.IntN(lanes)),
			uint32(a.rnd.IntN(lanes)),
		)
	})
}

func (a insElemArb) Shrink(p InsElemParams) iter.Seq[InsElemParams] {
	out := make([]InsElemParams, 0, 8)
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewInsElemParams(r, p.Rn, p.Elem, p.Idx, p.SrcIdx))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewInsElemParams(p.Rd, r, p.Elem, p.Idx, p.SrcIdx))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewInsElemParams(p.Rd, p.Rn, p.Elem, i, p.SrcIdx))
	}

	for _, i := range u32Halved(p.SrcIdx) {
		out = append(out, NewInsElemParams(p.Rd, p.Rn, p.Elem, p.Idx, i))
	}

	for _, e := range elemShrunk(p.Elem) {
		out = append(out, NewInsElemParams(p.Rd, p.Rn, e, p.Idx, p.SrcIdx))
	}

	return slices.Values(out)
}

// --- smov: smov wd, vn.b|h|s[idx] --------------------------------------

// SmovParams — parameters of smov.
type SmovParams struct {
	Wd   arm64.Reg
	Vn   arm64.VReg
	Elem string
	Idx  uint32
}

func NewSmovParams(wd arm64.Reg, vn arm64.VReg, elem string, idx uint32) SmovParams {
	return SmovParams{
		Wd:   wd,
		Vn:   vn,
		Elem: elem,
		Idx:  idx,
	}
}

func (p SmovParams) Instr() arm64.Instr {
	in, err := arm64.New().Smov(p.Wd, p.Vn, p.Elem, p.Idx)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p SmovParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Smov — an arbitrary smov (the .s form takes an x destination).
func Smov(rnd *rand.Rand) ohsnap.Arbitrary[SmovParams] {
	return smovArb{rnd: rnd}
}

type smovArb struct {
	rnd *rand.Rand
}

func (a smovArb) Generate() iter.Seq[SmovParams] {
	return arbStream(func() SmovParams {
		elem := genElem(a.rnd, []string{"b", "h", "s"})
		is64 := elem == "s" || a.rnd.IntN(2) == 1
		return NewSmovParams(
			genWidthReg(a.rnd, is64),
			genV(a.rnd),
			elem,
			uint32(a.rnd.IntN(int(elemLanes(elemWidth(elem))))),
		)
	})
}

func (a smovArb) Shrink(p SmovParams) iter.Seq[SmovParams] {
	out := make([]SmovParams, 0, 8)
	for _, r := range regShrunk(p.Wd) {
		out = append(out, NewSmovParams(r, p.Vn, p.Elem, p.Idx))
	}

	for _, r := range vShrunk(p.Vn) {
		out = append(out, NewSmovParams(p.Wd, r, p.Elem, p.Idx))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewSmovParams(p.Wd, p.Vn, p.Elem, i))
	}

	for _, e := range elemShrunk(p.Elem) {
		out = append(out, NewSmovParams(p.Wd, p.Vn, e, p.Idx))
	}

	return slices.Values(out)
}

// --- umov: umov|mov wd, vn.sz[idx] -------------------------------------

// UmovParams — parameters of umov.
type UmovParams struct {
	Wd   arm64.Reg
	Vn   arm64.VReg
	Elem string
	Idx  uint32
}

func NewUmovParams(wd arm64.Reg, vn arm64.VReg, elem string, idx uint32) UmovParams {
	return UmovParams{
		Wd:   wd,
		Vn:   vn,
		Elem: elem,
		Idx:  idx,
	}
}

func (p UmovParams) Instr() arm64.Instr {
	in, err := arm64.New().Umov(p.Wd, p.Vn, p.Elem, p.Idx)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p UmovParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Umov — an arbitrary umov (the element fills the destination: .b/.h/.s
// into w, .d into x; the filling forms print as mov).
func Umov(rnd *rand.Rand) ohsnap.Arbitrary[UmovParams] {
	return umovArb{rnd: rnd}
}

type umovArb struct {
	rnd *rand.Rand
}

func (a umovArb) Generate() iter.Seq[UmovParams] {
	return arbStream(func() UmovParams {
		elem := genElem(a.rnd, []string{"b", "h", "s", "d"})
		return NewUmovParams(
			genWidthReg(a.rnd, elem == "d"),
			genV(a.rnd),
			elem,
			uint32(a.rnd.IntN(int(elemLanes(elemWidth(elem))))),
		)
	})
}

func (a umovArb) Shrink(p UmovParams) iter.Seq[UmovParams] {
	out := make([]UmovParams, 0, 8)
	for _, r := range regShrunk(p.Wd) {
		out = append(out, NewUmovParams(r, p.Vn, p.Elem, p.Idx))
	}

	for _, r := range vShrunk(p.Vn) {
		out = append(out, NewUmovParams(p.Wd, r, p.Elem, p.Idx))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewUmovParams(p.Wd, p.Vn, p.Elem, i))
	}

	for _, e := range elemShrunk(p.Elem) {
		if (e == "d") != (p.Elem == "d") {
			continue // the gpr width stays tied to the element
		}

		out = append(out, NewUmovParams(p.Wd, p.Vn, e, p.Idx))
	}

	return slices.Values(out)
}

// --- dup (element): dup.Arr vd, vn[idx] --------------------------------

// DupElemParams — parameters of dup (element).
type DupElemParams struct {
	Rd  arm64.VReg
	Rn  arm64.VReg
	Arr string
	Idx uint32
}

func NewDupElemParams(rd arm64.VReg, rn arm64.VReg, arr string, idx uint32) DupElemParams {
	return DupElemParams{
		Rd:  rd,
		Rn:  rn,
		Arr: arr,
		Idx: idx,
	}
}

func (p DupElemParams) Instr() arm64.Instr {
	in, err := arm64.New().DupElem(p.Rd, p.Rn, p.Arr, p.Idx)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p DupElemParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// DupElem — an arbitrary dup (element).
func DupElem(rnd *rand.Rand) ohsnap.Arbitrary[DupElemParams] {
	return dupElemArb{rnd: rnd}
}

type dupElemArb struct {
	rnd *rand.Rand
}

func (a dupElemArb) Generate() iter.Seq[DupElemParams] {
	return arbStream(func() DupElemParams {
		arr := genElem(a.rnd, arrDupElem())
		return NewDupElemParams(
			genV(a.rnd),
			genV(a.rnd),
			arr,
			uint32(a.rnd.IntN(int(elemLanes(elemWidth(arr[len(arr)-1:]))))),
		)
	})
}

func (a dupElemArb) Shrink(p DupElemParams) iter.Seq[DupElemParams] {
	out := make([]DupElemParams, 0, 8)
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewDupElemParams(r, p.Rn, p.Arr, p.Idx))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewDupElemParams(p.Rd, r, p.Arr, p.Idx))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewDupElemParams(p.Rd, p.Rn, p.Arr, i))
	}

	for _, a := range arrShrunk(p.Arr, arrDupElem()) {
		out = append(out, NewDupElemParams(p.Rd, p.Rn, a, p.Idx))
	}

	return slices.Values(out)
}

// --- dup (scalar): mov <b|h|s|d>n, vn.sz[idx] --------------------------

// DupScalarParams — parameters of dup (scalar).
type DupScalarParams struct {
	Rd   arm64.VReg
	Rn   arm64.VReg
	Elem string
	Idx  uint32
}

func NewDupScalarParams(rd arm64.VReg, rn arm64.VReg, elem string, idx uint32) DupScalarParams {
	return DupScalarParams{
		Rd:   rd,
		Rn:   rn,
		Elem: elem,
		Idx:  idx,
	}
}

func (p DupScalarParams) Instr() arm64.Instr {
	in, err := arm64.New().DupScalar(p.Rd, p.Rn, p.Elem, p.Idx)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p DupScalarParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// DupScalar — an arbitrary dup (scalar).
func DupScalar(rnd *rand.Rand) ohsnap.Arbitrary[DupScalarParams] {
	return dupScalarArb{rnd: rnd}
}

// elemLanes — the lane count of the full vector register by the element
// width number (b:16 h:8 s:4 d:2): the bound of the lane indexes.
func elemLanes(size uint32) uint32 {
	return 16 >> size
}

// elemWidth — the element width number by its letter (b/h/s/d).
func elemWidth(elem string) uint32 {
	switch elem {
	case "h":
		return 1
	case "s":
		return 2
	case "d":
		return 3
	default:
		return 0
	}
}

// genElem — an element letter from the set.
func genElem(rnd *rand.Rand, set []string) string {
	return set[rnd.IntN(len(set))]
}

// genWidthReg — a gpr of the given width class.
func genWidthReg(rnd *rand.Rand, is64 bool) arm64.Reg {
	return genReg(rnd, is64, false, true)
}

// elemShrunk — the element letter shrink: the narrowest letter whose
// width keeps the register rule valid (the caller filters by its own
// width law).
func elemShrunk(elem string) []string {
	if elem == "b" {
		return nil
	}

	return []string{"b"}
}

// --- ins (general): mov.sz vd[idx], wn --------------------------------

// arrDupElem — the dup element arrangement set (the .1d form does not
// exist: clang refuses it).
func arrDupElem() []string {
	return []string{"8b", "16b", "4h", "8h", "2s", "4s", "2d"}
}

type dupScalarArb struct {
	rnd *rand.Rand
}

func (a dupScalarArb) Generate() iter.Seq[DupScalarParams] {
	return arbStream(func() DupScalarParams {
		elem := genElem(a.rnd, []string{"b", "h", "s", "d"})
		return NewDupScalarParams(
			genV(a.rnd),
			genV(a.rnd),
			elem,
			uint32(a.rnd.IntN(int(elemLanes(elemWidth(elem))))),
		)
	})
}

func (a dupScalarArb) Shrink(p DupScalarParams) iter.Seq[DupScalarParams] {
	out := make([]DupScalarParams, 0, 8)
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewDupScalarParams(r, p.Rn, p.Elem, p.Idx))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewDupScalarParams(p.Rd, r, p.Elem, p.Idx))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewDupScalarParams(p.Rd, p.Rn, p.Elem, i))
	}

	for _, e := range elemShrunk(p.Elem) {
		out = append(out, NewDupScalarParams(p.Rd, p.Rn, e, p.Idx))
	}

	return slices.Values(out)
}
