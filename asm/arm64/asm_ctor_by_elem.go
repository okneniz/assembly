package arm64

// The by-element ctor registrations and constructor (moved from
// arch/arm64/by_elem.go: asm-layer code registering into armCtors).

import (
	"errors"
	"fmt"
	arch "github.com/okneniz/assembly/arch/arm64"
)

func init() {
	// int: long × {"", "2"} × 3 arrangements; the rest × 2q × size
	for u := range 2 {
		for _, name := range arch.ByElemIntNames[u] {
			if name == "" || name == "fcmla" {
				continue
			}

			if arch.ByElemLong[name] {
				// Q=0 and Q=1 forms print the same result
				// arrangement; the keys differ by the "2" suffix
				for size := range uint32(3) {
					arr := decodeArrangement(1, size+1)
					armCtors[name+"."+arr] = newByElemArm(name, 0, size, true)
					armCtors[name+"2."+arr] = newByElemArm(name, 1, size, true)
				}
			} else {
				for q := range uint32(2) {
					for size := uint32(1); size < 3; size++ {
						k := name + "." + decodeArrangement(q, size)
						armCtors[k] = newByElemArm(name, q, size, false)
					}
				}
			}
		}
	}

	// fp: .2s/.4s (size 0), .2d (size 3)
	for u := range 2 {
		for _, name := range arch.ByElemFPNames[u] {
			for q := range uint32(2) {
				k := name + "." + decodeArrangement(q, 2)
				armCtors[k] = newByElemArm(name, q, 0, false)
			}

			armCtors[name+".2d"] = newByElemArm(name, 1, 3, false)
		}
	}

	// fcmla: 2 arrangement classes (rotation — the #rot operand, opc computed
	// at encoding time)
	for q := range uint32(2) {
		for size := uint32(1); size < 3; size++ {
			armCtors["fcmla."+decodeArrangement(q, size)] = newByElemArm("fcmla", q, size, false)
		}
	}
}

// newByElemArm — name.Arr vd, vn, vm[idx]{, #rot (fcmla)}.
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

		rdN, rnN, rmN, err := regNums3(rd, rn, rm)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}

		name2 := name
		if long && q == 1 {
			name2 += "2"
		}

		var rot uint32
		if name == "fcmla" {
			if len(ops) < 4 || ops[3].Kind() != arch.ArmOpImm {
				return nil, errors.New("fcmla: rotation expected (#0/#90/#180/#270)")
			}

			r := ops[3].Num()
			if r%90 != 0 || r < 0 || r > 270 {
				return nil, errors.New("fcmla: bad rotation")
			}

			rot = uint32(r)
		}

		return arch.Builder{}.ByElem(name2, q, size, rd, rn, rm, idx, long, rot, rdN, rnN, rmN), nil
	}
}
