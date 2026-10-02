package arm64

// The SIMD vector-family cores: the arrangement-suffixed forms over the
// v0..v31 file. The per-family arrangement set is the family's own —
// the plain three-same set, the logical b-only set, the widen set, the
// acc set — the thin per-instruction files bind it. The shift forms
// (shl/ushr/...) carry the element-bound shift amount.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
)

// mustV — the checked VReg constructor: the input is bounded by the
// caller (0..31), the error is unreachable.
func mustV(n int) arm64.VReg {
	r, err := arm64.V(n)
	if err != nil {
		return arm64.VReg{} // unreachable: n is always in 0..31
	}

	return r
}

// genV — a vector register v0..v31.
func genV(rnd *rand.Rand) arm64.VReg {
	return mustV(rnd.IntN(32))
}

// vShrunk — the shrink candidates of a vector register (v0, then the
// halved number).
func vShrunk(r arm64.VReg) []arm64.VReg {
	var out []arm64.VReg
	if r.Num() != 0 {
		out = append(out, mustV(0))
	}

	if n := int(r.Num()); n > 1 {
		out = append(out, mustV(n/2))
	}

	return out
}

// The per-family arrangement sets.
func arrFull() []string {
	return []string{"8b", "16b", "4h", "8h", "2s", "4s", "2d"}
}

func arrLogical() []string {
	return []string{"8b", "16b"}
}

func arrWiden() []string {
	return []string{"4h", "8h", "2s", "4s", "2d"}
}

func arrHalf() []string {
	return []string{"8b", "16b", "4h", "8h"}
}

// arrShrunk — the arrangement shrink candidates: the set's first entry
// (the canonical 8b-class lane of the family).
func arrShrunk(arr string, set []string) []string {
	if arr == set[0] {
		return nil
	}

	return []string{set[0]}
}

// elemBits — the lane width of an arrangement (the shift bound).
func elemBits(arr string) uint32 {
	switch arr {
	case "4h", "8h":
		return 16
	case "2s", "4s":
		return 32
	case "2d":
		return 64
	default:
		return 8
	}
}

// V2Params — rd, rn, arr of the two-register vector forms.
type V2Params struct {
	Rd, Rn arm64.VReg
	Arr    string
}

func NewV2Params(rd arm64.VReg, rn arm64.VReg, arr string) V2Params {
	return V2Params{
		Rd:  rd,
		Rn:  rn,
		Arr: arr,
	}
}

// v2Gen — the two-register core over the family's arrangement set.
type v2Gen struct {
	rnd *rand.Rand
	set []string
}

func newV2Gen(rnd *rand.Rand, set []string) v2Gen {
	return v2Gen{
		rnd: rnd,
		set: set,
	}
}

func (g v2Gen) Generate() iter.Seq[V2Params] {
	return arbStream(func() V2Params {
		return NewV2Params(
			genV(g.rnd),
			genV(g.rnd),
			g.set[g.rnd.IntN(len(g.set))],
		)
	})
}

func (g v2Gen) Shrink(p V2Params) iter.Seq[V2Params] {
	var out []V2Params
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewV2Params(r, p.Rn, p.Arr))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewV2Params(p.Rd, r, p.Arr))
	}

	for _, a := range arrShrunk(p.Arr, g.set) {
		out = append(out, NewV2Params(p.Rd, p.Rn, a))
	}

	return slices.Values(out)
}

// V3Params — rd, rn, rm, arr of the three-register vector forms.
type V3Params struct {
	Rd, Rn, Rm arm64.VReg
	Arr        string
}

func NewV3Params(rd arm64.VReg, rn arm64.VReg, rm arm64.VReg, arr string) V3Params {
	return V3Params{
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Arr: arr,
	}
}

// v3Gen — the three-register core over the family's arrangement set.
type v3Gen struct {
	rnd *rand.Rand
	set []string
}

func newV3Gen(rnd *rand.Rand, set []string) v3Gen {
	return v3Gen{
		rnd: rnd,
		set: set,
	}
}

func (g v3Gen) Generate() iter.Seq[V3Params] {
	return arbStream(func() V3Params {
		return NewV3Params(
			genV(g.rnd),
			genV(g.rnd),
			genV(g.rnd),
			g.set[g.rnd.IntN(len(g.set))],
		)
	})
}

func (g v3Gen) Shrink(p V3Params) iter.Seq[V3Params] {
	var out []V3Params
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewV3Params(r, p.Rn, p.Rm, p.Arr))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewV3Params(p.Rd, r, p.Rm, p.Arr))
	}

	for _, r := range vShrunk(p.Rm) {
		out = append(out, NewV3Params(p.Rd, p.Rn, r, p.Arr))
	}

	for _, a := range arrShrunk(p.Arr, g.set) {
		out = append(out, NewV3Params(p.Rd, p.Rn, p.Rm, a))
	}

	return slices.Values(out)
}

// V2PlainParams — rd, rn of the arrangement-less vector forms (aese,
// tbl takes its own shape).
type V2PlainParams struct {
	Rd, Rn arm64.VReg
}

func NewV2PlainParams(rd arm64.VReg, rn arm64.VReg) V2PlainParams {
	return V2PlainParams{
		Rd: rd,
		Rn: rn,
	}
}

// v2PlainGen — the arrangement-less two-register core.
type v2PlainGen struct {
	rnd *rand.Rand
}

func newV2PlainGen(rnd *rand.Rand) v2PlainGen {
	return v2PlainGen{rnd: rnd}
}

func v2Plain(rnd *rand.Rand) v2PlainGen {
	return newV2PlainGen(rnd)
}

func (g v2PlainGen) Generate() iter.Seq[V2PlainParams] {
	return arbStream(func() V2PlainParams {
		return NewV2PlainParams(genV(g.rnd), genV(g.rnd))
	})
}

func (g v2PlainGen) Shrink(p V2PlainParams) iter.Seq[V2PlainParams] {
	var out []V2PlainParams
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewV2PlainParams(r, p.Rn))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewV2PlainParams(p.Rd, r))
	}

	return slices.Values(out)
}

// V3PlainParams — rd, rn, rm of the arrangement-less tbl.
type V3PlainParams struct {
	Rd, Rn, Rm arm64.VReg
}

func NewV3PlainParams(rd arm64.VReg, rn arm64.VReg, rm arm64.VReg) V3PlainParams {
	return V3PlainParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

// v3PlainGen — the arrangement-less three-register core.
type v3PlainGen struct {
	rnd *rand.Rand
}

func newV3PlainGen(rnd *rand.Rand) v3PlainGen {
	return v3PlainGen{rnd: rnd}
}

func v3Plain(rnd *rand.Rand) v3PlainGen {
	return newV3PlainGen(rnd)
}

func (g v3PlainGen) Generate() iter.Seq[V3PlainParams] {
	return arbStream(func() V3PlainParams {
		return NewV3PlainParams(genV(g.rnd), genV(g.rnd), genV(g.rnd))
	})
}

func (g v3PlainGen) Shrink(p V3PlainParams) iter.Seq[V3PlainParams] {
	var out []V3PlainParams
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewV3PlainParams(r, p.Rn, p.Rm))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewV3PlainParams(p.Rd, r, p.Rm))
	}

	for _, r := range vShrunk(p.Rm) {
		out = append(out, NewV3PlainParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}

// VShiftParams — rd, rn, arr, shift of the element-bound shift forms.
type VShiftParams struct {
	Rd, Rn arm64.VReg
	Arr    string
	Shift  uint32
}

func NewVShiftParams(rd arm64.VReg, rn arm64.VReg, arr string, shift uint32) VShiftParams {
	return VShiftParams{
		Rd:    rd,
		Rn:    rn,
		Arr:   arr,
		Shift: shift,
	}
}

// vShiftGen — the shift core: the full arrangement set, the shift
// uniform in the lane width (0..elemBits-1; the right-shift families
// are 1-based — min1).
type vShiftGen struct {
	rnd  *rand.Rand
	set  []string
	min1 bool
}

func newVShiftGen(rnd *rand.Rand, set []string) vShiftGen {
	return vShiftGen{
		rnd: rnd,
		set: set,
	}
}

// newVShiftGen1 — the 1-based variant (sshr/ushr/sri: 1..elemBits).
func newVShiftGen1(rnd *rand.Rand, set []string) vShiftGen {
	g := newVShiftGen(rnd, set)
	g.min1 = true
	return g
}

func (g vShiftGen) Generate() iter.Seq[VShiftParams] {
	return arbStream(func() VShiftParams {
		arr := g.set[g.rnd.IntN(len(g.set))]
		bits := int(elemBits(arr))
		shift := g.rnd.IntN(bits)
		if g.min1 {
			shift++ // the right-shift families: 1..elemBits
		}

		return NewVShiftParams(
			genV(g.rnd),
			genV(g.rnd),
			arr,
			uint32(shift),
		)
	})
}

func (g vShiftGen) Shrink(p VShiftParams) iter.Seq[VShiftParams] {
	var out []VShiftParams
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewVShiftParams(r, p.Rn, p.Arr, p.Shift))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewVShiftParams(p.Rd, r, p.Arr, p.Shift))
	}

	for _, v := range u32Halved(p.Shift) {
		if g.min1 && v == 0 {
			continue // zero is outside the 1-based range
		}

		out = append(out, NewVShiftParams(p.Rd, p.Rn, p.Arr, v))
	}

	for _, a := range arrShrunk(p.Arr, g.set) {
		sh := p.Shift
		if sh >= elemBits(a) {
			sh = 0
		}

		out = append(out, NewVShiftParams(p.Rd, p.Rn, a, sh))
	}

	return slices.Values(out)
}

// DupParams — vd, wn (a gpr source), arr of the dup-to-vector forms.
type DupParams struct {
	Rd  arm64.VReg
	Wn  arm64.Reg // x/w (the 31st reads as zr)
	Arr string
}

func NewDupParams(rd arm64.VReg, wn arm64.Reg, arr string) DupParams {
	return DupParams{
		Rd:  rd,
		Wn:  wn,
		Arr: arr,
	}
}

// dupGen — the dup core: the gpr width follows the lane size (the b/h/s
// lanes take a w register, d takes an x one).
type dupGen struct {
	rnd *rand.Rand
}

func newDupGen(rnd *rand.Rand) dupGen {
	return dupGen{rnd: rnd}
}

func (g dupGen) Generate() iter.Seq[DupParams] {
	return arbStream(func() DupParams {
		arr := ohsnap.First(arbEnumArr(g.rnd).Generate())
		is64 := arr == "2d"
		return NewDupParams(
			genV(g.rnd),
			genReg(g.rnd, is64, false, true),
			arr,
		)
	})
}

func (g dupGen) Shrink(p DupParams) iter.Seq[DupParams] {
	var out []DupParams
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewDupParams(r, p.Wn, p.Arr))
	}

	for _, r := range regShrunk(p.Wn) {
		out = append(out, NewDupParams(p.Rd, r, p.Arr))
	}

	for _, a := range arrShrunk(p.Arr, arrFull()) {
		if a == "2d" != p.Wn.Is64() {
			continue // the lane width stays tied to the gpr width
		}

		out = append(out, NewDupParams(p.Rd, p.Wn, a))
	}

	return slices.Values(out)
}

// arbEnumArr — an arrangement from the full set (the Enum adapter).
func arbEnumArr(rnd *rand.Rand) ohsnap.Arbitrary[string] {
	set := arrFull()
	vals := make([]string, len(set))
	copy(vals, set)
	return arbEnum(rnd, vals)
}
