package arm64

// The by-element ctor registrations and constructor (moved from
// arch/arm64/by_elem.go: asm-layer code registering into armCtors).
// registerByElem is called from buildArmCtors: the map is assembled in
// one place, no post-init mutation. Every entry builds through the
// arch Builder - the per-instruction methods of the family.

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// byElemSpec - one mnemonic's registration shape: the key set is
// generated from it (arrangements x Q), the encoding lives in the
// Builder method the name dispatches to.
type byElemSpec struct {
	name     string
	long, fp bool // long: the result is one lane wider; fp: .2s/.4s/.2d
}

// byElemOps - the family's mnemonics (fcmla rides along: its rotation
// is a fourth operand, its arrangements are the integer ones).
var byElemOps = []byElemSpec{
	{"mla", false, false}, {"mls", false, false}, {"mul", false, false},
	{"sqdmulh", false, false}, {"sqrdmulh", false, false},
	{"sqrdmlah", false, false}, {"sqrdmlsh", false, false},
	{"smlal", true, false}, {"sqdmlal", true, false},
	{"smlsl", true, false}, {"sqdmlsl", true, false},
	{"smull", true, false}, {"sqdmull", true, false},
	{"umlal", true, false}, {"umlsl", true, false}, {"umull", true, false},
	{"fmla", false, true}, {"fmls", false, true},
	{"fmul", false, true}, {"fmulx", false, true},
	{"fcmla", false, false},
}

// registerByElem adds the by-element mnemonics (name.arr keys) to m.
func registerByElem(m map[string]func(ops []vOp) (Instr, error)) {
	for _, op := range byElemOps {
		switch {
		case op.long:
			// Q=0 and Q=1 forms print the same result arrangement; the
			// keys differ by the "2" suffix
			for size := range uint32(3) {
				arr := decodeArrangement(1, size+1)
				m[op.name+"."+arr] = newByElemArm(op.name, 0, size, true)
				m[op.name+"2."+arr] = newByElemArm(op.name, 1, size, true)
			}
		case op.fp:
			for q := range uint32(2) {
				m[op.name+"."+decodeArrangement(q, 2)] = newByElemArm(op.name, q, 0, false)
			}

			m[op.name+".2d"] = newByElemArm(op.name, 1, 3, false)
		default:
			if op.name == "fcmla" {
				for q := range uint32(2) {
					for size := uint32(1); size < 3; size++ {
						m[op.name+"."+decodeArrangement(q, size)] =
							newByElemArm(op.name, q, size, false)
					}
				}

				continue
			}

			for q := range uint32(2) {
				for size := uint32(1); size < 3; size++ {
					m[op.name+"."+decodeArrangement(q, size)] =
						newByElemArm(op.name, q, size, false)
				}
			}
		}
	}
}

// newByElemArm — name.Arr vd, vn, vm[idx]{, #rot (fcmla)}. The
// captured (q, size) is authoritative (the operand suffix is checked
// by the Builder's arrangement validation).
func newByElemArm(name string, q, size uint32, long bool) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) < 3 || ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: want vd.Arr, vn, vm[idx]", name)
		}

		rd, err := wantV(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		rm, err := wantV(ops[2], name)
		if err != nil {
			return nil, err
		}

		idx := ops[2].Num()
		if idx < 0 {
			return nil, fmt.Errorf("%s: bad index", name)
		}

		b := arch.Builder{}
		if long {
			arr := decodeArrangement(1, size+1)
			two := q == 1
			switch name {
			case "smlal":
				return b.SmlalElem(rd, rn, rm, arr, two, uint32(idx))
			case "sqdmlal":
				return b.SqdmlalElem(rd, rn, rm, arr, two, uint32(idx))
			case "smlsl":
				return b.SmlslElem(rd, rn, rm, arr, two, uint32(idx))
			case "sqdmlsl":
				return b.SqdmlslElem(rd, rn, rm, arr, two, uint32(idx))
			case "smull":
				return b.SmullElem(rd, rn, rm, arr, two, uint32(idx))
			case "sqdmull":
				return b.SqdmullElem(rd, rn, rm, arr, two, uint32(idx))
			case "umlal":
				return b.UmlalElem(rd, rn, rm, arr, two, uint32(idx))
			case "umlsl":
				return b.UmlslElem(rd, rn, rm, arr, two, uint32(idx))
			default:
				return b.UmullElem(rd, rn, rm, arr, two, uint32(idx))
			}
		}

		if name == "fcmla" {
			if len(ops) < 4 || ops[3].Kind() != arch.ArmOpImm {
				return nil, errors.New("fcmla: rotation expected (#0/#90/#180/#270)")
			}

			r := ops[3].Num()
			if r%90 != 0 || r < 0 || r > 270 {
				return nil, errors.New("fcmla: bad rotation")
			}

			arr := decodeArrangement(q, size)
			if size == 3 {
				arr = "2d"
			}

			return b.FcmlaElem(rd, rn, rm, arr, uint32(idx), uint32(r))
		}

		if size == 0 || size == 3 { // fp: .2s/.4s/.2d
			arr := decodeArrangement(q, 2)
			if size == 3 {
				arr = "2d"
			}

			switch name {
			case "fmla":
				return b.FmlaElem(rd, rn, rm, arr, uint32(idx))
			case "fmls":
				return b.FmlsElem(rd, rn, rm, arr, uint32(idx))
			case "fmulx":
				return b.FmulxElem(rd, rn, rm, arr, uint32(idx))
			default:
				return b.FmulElem(rd, rn, rm, arr, uint32(idx))
			}
		}

		arr := decodeArrangement(q, size)
		switch name {
		case "mla":
			return b.MlaElem(rd, rn, rm, arr, uint32(idx))
		case "mls":
			return b.MlsElem(rd, rn, rm, arr, uint32(idx))
		case "mul":
			return b.MulElem(rd, rn, rm, arr, uint32(idx))
		case "sqdmulh":
			return b.SqdmulhElem(rd, rn, rm, arr, uint32(idx))
		case "sqrdmulh":
			return b.SqrdmulhElem(rd, rn, rm, arr, uint32(idx))
		case "sqrdmlah":
			return b.SqrdmlahElem(rd, rn, rm, arr, uint32(idx))
		default:
			return b.SqrdmlshElem(rd, rn, rm, arr, uint32(idx))
		}
	}
}
