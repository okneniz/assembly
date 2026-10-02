package arm64

// The by-element group: a vector times one lane of another. One
// parameter type over all 21 mnemonics — the Op field names the family,
// the spec table binds its arrangement set, its lane laws and its
// Builder call. The lane index and the Vm width follow the lane width,
// with the two layout exceptions the specs express: the long families
// name the RESULT arrangement (the index follows the source width) and
// fcmla carries a 5-bit Vm with a 1-2 bit index.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ByElemOp — the by-element mnemonic (the family identity of the
// parameters).
type ByElemOp int

const (
	ByElemMla ByElemOp = iota
	ByElemMls
	ByElemMul
	ByElemSqdmulh
	ByElemSqrdmulh
	ByElemSqrdmlah
	ByElemSqrdmlsh
	ByElemFmla
	ByElemFmls
	ByElemFmul
	ByElemFmulx
	ByElemSmlal
	ByElemSmlsl
	ByElemSmull
	ByElemUmlal
	ByElemUmlsl
	ByElemUmull
	ByElemSqdmlal
	ByElemSqdmlsl
	ByElemSqdmull
	ByElemFcmla
	byElemOpCount
)

// ByElemParams — parameters of the by-element families.
type ByElemParams struct {
	Op         ByElemOp
	Rd, Rn, Rm arm64.VReg
	Arr        string
	Idx        uint32
	Two        bool   // the long families' upper-half form
	Rot        uint32 // the fcmla rotation (0/90/180/270)
}

func NewByElemParams(
	op ByElemOp,
	rd arm64.VReg,
	rn arm64.VReg,
	rm arm64.VReg,
	arr string,
	idx uint32,
	two bool,
	rot uint32,
) ByElemParams {
	return ByElemParams{
		Op:  op,
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Arr: arr,
		Idx: idx,
		Two: two,
		Rot: rot,
	}
}

func (p ByElemParams) Instr() arm64.Instr {
	in, err := byElemSpecs[p.Op].mk(arm64.New(), p)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p ByElemParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// the per-kind arrangement sets.
func byElemIntArrs() []string {
	return []string{"4h", "8h", "2s", "4s"}
}

func byElemFpArrs() []string {
	return []string{"2s", "4s", "2d"}
}

func byElemLongArrs() []string {
	return []string{"4s", "2d"}
}

func fcmlaArrs() []string {
	return []string{"4h", "8h", "4s"}
}

// arrLaneWidth — the lane width number of an arrangement (h:1 s:2 d:3).
func arrLaneWidth(arr string) uint32 {
	switch arr[len(arr)-1] {
	case 'h':
		return 1
	case 's':
		return 2
	case 'd':
		return 3
	default:
		return 0
	}
}

// stdIdxMax — the lane-index bound of the standard by-element layout by
// the lane width (.h carries a 3-bit index with a 4-bit Vm, .s/.d
// shrink by one bit each).
func stdIdxMax(width uint32) uint32 {
	return 8 >> (width - 1)
}

// longSrcWidth — the source lane width of a long family (the result
// arrangement .4s reads .h sources, .2d reads .s).
func longSrcWidth(arr string) uint32 {
	if arr == "4s" {
		return 1
	}

	return 2
}

// byElemSpec — one mnemonic's laws and Builder call.
type byElemSpec struct {
	set    []string
	sizeOf func(arr string) uint32 // the lane width the index and Vm follow
	idxOf  func(arr string) uint32 // the lane-index bound
	rmAny  bool                    // the Vm field is 5-bit regardless of width
	two    bool                    // generate the upper-half flag
	rot    bool                    // generate the fcmla rotation
	mk     func(arm64.Builder, ByElemParams) (arm64.Instr, error)
}

// stdSpec — the standard layout over an arrangement set.
func stdSpec(set []string, mk func(arm64.Builder, ByElemParams) (arm64.Instr, error)) byElemSpec {
	return byElemSpec{
		set:    set,
		sizeOf: arrLaneWidth,
		idxOf:  func(arr string) uint32 { return stdIdxMax(arrLaneWidth(arr)) },
		mk:     mk,
	}
}

// longMk — the Builder call adapter of the long families (the two flag
// sits before the index in their signatures).
func longMk(
	call func(arm64.Builder, arm64.VReg, arm64.VReg, arm64.VReg, string, bool, uint32) (arm64.Instr, error),
) func(arm64.Builder, ByElemParams) (arm64.Instr, error) {
	return func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return call(b, p.Rd, p.Rn, p.Rm, p.Arr, p.Two, p.Idx)
	}
}

// longSpec — the long families: the arrangement names the result, the
// index and the Vm width follow the source width.
func longSpec(
	call func(arm64.Builder, arm64.VReg, arm64.VReg, arm64.VReg, string, bool, uint32) (arm64.Instr, error),
) byElemSpec {
	return byElemSpec{
		set:    byElemLongArrs(),
		sizeOf: longSrcWidth,
		idxOf:  func(arr string) uint32 { return stdIdxMax(longSrcWidth(arr)) },
		two:    true,
		mk:     longMk(call),
	}
}

var byElemSpecs = [byElemOpCount]byElemSpec{
	ByElemMla: stdSpec(byElemIntArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.MlaElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemMls: stdSpec(byElemIntArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.MlsElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemMul: stdSpec(byElemIntArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.MulElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemSqdmulh: stdSpec(
		byElemIntArrs(),
		func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.SqdmulhElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
		},
	),
	ByElemSqrdmulh: stdSpec(
		byElemIntArrs(),
		func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.SqrdmulhElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
		},
	),
	ByElemSqrdmlah: stdSpec(
		byElemIntArrs(),
		func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.SqrdmlahElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
		},
	),
	ByElemSqrdmlsh: stdSpec(
		byElemIntArrs(),
		func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.SqrdmlshElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
		},
	),
	ByElemFmla: stdSpec(byElemFpArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.FmlaElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemFmls: stdSpec(byElemFpArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.FmlsElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemFmul: stdSpec(byElemFpArrs(), func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
		return b.FmulElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
	}),
	ByElemFmulx: stdSpec(
		byElemFpArrs(),
		func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.FmulxElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx)
		},
	),
	ByElemSmlal:   longSpec(arm64.Builder.SmlalElem),
	ByElemSmlsl:   longSpec(arm64.Builder.SmlslElem),
	ByElemSmull:   longSpec(arm64.Builder.SmullElem),
	ByElemUmlal:   longSpec(arm64.Builder.UmlalElem),
	ByElemUmlsl:   longSpec(arm64.Builder.UmlslElem),
	ByElemUmull:   longSpec(arm64.Builder.UmullElem),
	ByElemSqdmlal: longSpec(arm64.Builder.SqdmlalElem),
	ByElemSqdmlsl: longSpec(arm64.Builder.SqdmlslElem),
	ByElemSqdmull: longSpec(arm64.Builder.SqdmullElem),
	ByElemFcmla: {
		set:    fcmlaArrs(),
		sizeOf: arrLaneWidth,
		idxOf: func(arr string) uint32 {
			if arr == "8h" {
				return 4
			}

			return 2
		},
		rmAny: true,
		rot:   true,
		mk: func(b arm64.Builder, p ByElemParams) (arm64.Instr, error) {
			return b.FcmlaElem(p.Rd, p.Rn, p.Rm, p.Arr, p.Idx, p.Rot)
		},
	},
}

// byElemValid — the params fit the family's lane laws (a .h source
// stays inside the 4-bit Vm field unless the family carries a 5-bit Vm).
func byElemValid(spec byElemSpec, p ByElemParams) bool {
	if !slices.Contains(spec.set, p.Arr) {
		return false
	}

	rmMax := 32
	if !spec.rmAny && spec.sizeOf(p.Arr) == 1 {
		rmMax = 16
	}

	return p.Idx < spec.idxOf(p.Arr) && int(p.Rm.Num()) < rmMax
}

// byElemArb — the shared arbitrary of the group: every generated or
// shrunk sample fits its family's laws. One instance per family, the op
// field binds the mnemonic.
type byElemArb struct {
	rnd *rand.Rand
	op  ByElemOp
}

// the 21 family generators (the exported faces of the group).
func MlaElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemMla}
}

func MlsElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemMls}
}

func MulElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemMul}
}

func SqdmulhElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqdmulh}
}

func SqrdmulhElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqrdmulh}
}

func SqrdmlahElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqrdmlah}
}

func SqrdmlshElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqrdmlsh}
}

func FmlaElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemFmla}
}

func FmlsElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemFmls}
}

func FmulElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemFmul}
}

func FmulxElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemFmulx}
}

func SmlalElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSmlal}
}

func SmlslElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSmlsl}
}

func SmullElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSmull}
}

func UmlalElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemUmlal}
}

func UmlslElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemUmlsl}
}

func UmullElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemUmull}
}

func SqdmlalElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqdmlal}
}

func SqdmlslElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqdmlsl}
}

func SqdmullElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemSqdmull}
}

func FcmlaElem(rnd *rand.Rand) ohsnap.Arbitrary[ByElemParams] {
	return byElemArb{rnd: rnd, op: ByElemFcmla}
}

func (a byElemArb) Generate() iter.Seq[ByElemParams] {
	return arbStream(a.gen)
}

func (a byElemArb) Shrink(p ByElemParams) iter.Seq[ByElemParams] {
	spec := byElemSpecs[p.Op]
	var out []ByElemParams
	for _, r := range vShrunk(p.Rd) {
		out = append(out, NewByElemParams(p.Op, r, p.Rn, p.Rm, p.Arr, p.Idx, p.Two, p.Rot))
	}

	for _, r := range vShrunk(p.Rn) {
		out = append(out, NewByElemParams(p.Op, p.Rd, r, p.Rm, p.Arr, p.Idx, p.Two, p.Rot))
	}

	for _, r := range vShrunk(p.Rm) {
		out = append(out, NewByElemParams(p.Op, p.Rd, p.Rn, r, p.Arr, p.Idx, p.Two, p.Rot))
	}

	for _, i := range u32Halved(p.Idx) {
		out = append(out, NewByElemParams(p.Op, p.Rd, p.Rn, p.Rm, p.Arr, i, p.Two, p.Rot))
	}

	if p.Two {
		out = append(out, NewByElemParams(p.Op, p.Rd, p.Rn, p.Rm, p.Arr, p.Idx, false, p.Rot))
	}

	if p.Rot != 0 {
		out = append(out, NewByElemParams(p.Op, p.Rd, p.Rn, p.Rm, p.Arr, p.Idx, p.Two, 0))
	}

	for _, arr := range arrShrunk(p.Arr, spec.set) {
		cand := NewByElemParams(p.Op, p.Rd, p.Rn, p.Rm, arr, p.Idx, p.Two, p.Rot)
		if byElemValid(spec, cand) {
			out = append(out, cand)
		}
	}

	return slices.Values(out)
}

// gen — one valid sample of the family.
func (a byElemArb) gen() ByElemParams {
	spec := byElemSpecs[a.op]
	arr := genElem(a.rnd, spec.set)

	rmMax := 32
	if !spec.rmAny && spec.sizeOf(arr) == 1 {
		rmMax = 16
	}

	rot := uint32(0)
	if spec.rot {
		rot = uint32(a.rnd.IntN(4)) * 90
	}

	return NewByElemParams(
		a.op,
		genV(a.rnd),
		genV(a.rnd),
		genV(a.rnd),
		arr,
		uint32(a.rnd.IntN(int(spec.idxOf(arr)))),
		spec.two && a.rnd.IntN(2) == 1,
		rot,
	).withRm(a, rmMax)
}

// withRm — regenerates the Rm until it fits the width (a .h source
// stays inside the 4-bit Vm field; at most a few rolls, the bound is
// half the file).
func (p ByElemParams) withRm(a byElemArb, rmMax int) ByElemParams {
	for int(p.Rm.Num()) >= rmMax {
		p = NewByElemParams(p.Op, p.Rd, p.Rn, genV(a.rnd), p.Arr, p.Idx, p.Two, p.Rot)
	}

	return p
}
