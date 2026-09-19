package arm64

// SIMD assembler constructors: every entry builds through the arch
// Builder (the one construction API) - the logical and arithmetic
// three-same families (and..bif/add/cmeq/addp/sqrshl), simd2
// (cnt/rev32/not/abs/rbit), the shifts (shl/sri/ushr/sshr),
// aese/aesmc, the copy family (dup/ins/smov/umov + the element/scalar
// dup forms), tbl, uaddlv, the structural ld1-ld4/st1-st4 (reglist/
// element forms), mov.16b/mov.8b. The .Arr() suffix carries Q/size
// (arrQSize - the inverse of decodeArrangement).

import (
	"errors"
	"fmt"
	"strings"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// wantV — a vector register operand as a VReg.
func wantV(op vOp, name string) (arch.VReg, error) {
	r, err := arch.VRegOf(op.Reg())
	if err != nil {
		return arch.VReg{}, fmt.Errorf("%s: vector register expected", name)
	}

	return r, nil
}

// newV3Arm — a three-same instruction (vd, vn, vm): the Builder
// method comes in as a method expression.
func newV3Arm(
	name string,
	method func(arch.Builder, arch.VReg, arch.VReg, arch.VReg, string) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", name)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected (.16b)", name)
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

		return method(arch.Builder{}, rd, rn, rm, ops[0].Arr())
	}
}

// newV2Arm — a two-register instruction (vd, vn).
func newV2Arm(
	name string,
	method func(arch.Builder, arch.VReg, arch.VReg, string) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want vd, vn", name)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", name)
		}

		rd, err := wantV(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn, ops[0].Arr())
	}
}

// newShiftArm — a vector shift (vd, vn, #shift): the amount as
// written on the arrangement's lanes.
func newShiftArm(
	name string,
	method func(arch.Builder, arch.VReg, arch.VReg, string, uint32) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 || ops[2].Kind() != arch.ArmOpImm {
			return nil, fmt.Errorf("%s: want vd, vn, #shift", name)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", name)
		}

		sh := ops[2].Num()
		if sh < 0 {
			return nil, fmt.Errorf("%s: bad shift", name)
		}

		rd, err := wantV(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn, ops[0].Arr(), uint32(sh))
	}
}

// newAes — aese/aesmc vd, vn (.16b implicit).
func newAes(
	name string,
	method func(arch.Builder, arch.VReg, arch.VReg) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want vd, vn", name)
		}

		rd, err := wantV(ops[0], name)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		return method(arch.Builder{}, rd, rn)
	}
}

// newWidenArm — the widening three-same family (vd, vn, vm): the
// arrangement is the RESULT's, one lane wider than the source.
func newWidenArm(
	name string,
	method func(arch.Builder, arch.VReg, arch.VReg, arch.VReg, string) (arch.Instr, error),
) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", name)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", name)
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

		return method(arch.Builder{}, rd, rn, rm, ops[0].Arr())
	}
}

// newDupArm — dup.Arr vd, wn | dup.Arr vd, vn[idx].
func newDupArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("dup: want vd, wn")
	}

	if ops[0].Arr() == "" {
		return nil, errors.New("dup: arrangement suffix expected")
	}

	rd, err := wantV(ops[0], "dup")
	if err != nil {
		return nil, err
	}

	// DUP (element): the source is a lane of a vector register
	if ops[1].Kind() == arch.ArmOpReg && ops[1].Reg() != "" && ops[1].Reg()[0] == 'v' {
		idx := ops[1].Num()
		if !ops[1].LaneIdx() || idx < 0 {
			return nil, errors.New("dup: want vd.Arr, vn[idx]")
		}

		rn, err := wantV(ops[1], "dup")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).DupElem(rd, rn, ops[0].Arr(), uint32(idx))
	}

	rn, err := gpr(ops[1], "dup")
	if err != nil {
		return nil, err
	}

	return (arch.Builder{}).Dup(rd, rn, ops[0].Arr())
}

// newInsElemArm — INS (element): ins.sz vd[idx], vn[idx]. (The GPR-source
// form has no spelling here: llvm prints it as the mov alias, so the
// canonical input is mov.sz vd[idx], wn — newMovInsArm.)
func newInsElemArm(elem string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 || !ops[0].LaneIdx() || !ops[1].LaneIdx() {
			return nil, errors.New("ins: want vd[idx], vn[idx]")
		}

		rd, err := wantV(ops[0], "ins")
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], "ins")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).InsElem(
			rd, rn, elem, uint32(ops[0].Num()), uint32(ops[1].Num()),
		)
	}
}

// newSmovUmovArm — SMOV/UMOV: op wd, vn.sz[idx] (the element size is
// the source suffix; an x destination sets Q).
func newSmovUmovArm(name string, isSmov bool) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 || ops[1].Kind() != arch.ArmOpReg || !ops[1].LaneIdx() {
			return nil, fmt.Errorf("%s: want wd, vn.sz[idx]", name)
		}

		elem := ops[1].Arr()
		if elem == "" {
			return nil, fmt.Errorf("%s: element size suffix (b/h/s/d) expected", name)
		}

		gprOp, err := gpr(ops[0], name)
		if err != nil {
			return nil, err
		}

		vd, err := wantV(ops[1], name)
		if err != nil {
			return nil, err
		}

		if isSmov {
			return (arch.Builder{}).Smov(gprOp, vd, elem, uint32(ops[1].Num()))
		}

		return (arch.Builder{}).Umov(gprOp, vd, elem, uint32(ops[1].Num()))
	}
}

// newTblArm — tbl.16b vd, { vn }, vm.
func newTblArm(ops []vOp) (Instr, error) {
	if len(ops) != 3 || ops[1].Kind() != arch.ArmOpList || len(ops[1].List()) != 1 {
		return nil, errors.New("tbl: want vd, { vn }, vm")
	}

	rd, err := wantV(ops[0], "tbl")
	if err != nil {
		return nil, err
	}

	rn, err := wantV(arch.VOpReg(ops[1].List()[0].Reg(), "", false, 0), "tbl")
	if err != nil {
		return nil, err
	}

	rm, err := wantV(ops[2], "tbl")
	if err != nil {
		return nil, err
	}

	return (arch.Builder{}).Tbl(rd, rn, rm)
}

// newUaddlv — uaddlv.Arr hd/sd/dd, vn (dest scalar by size).
func newUaddlv(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].Arr() == "" {
		return nil, errors.New("uaddlv: want vd, vn with .Arr()")
	}

	rd, err := wantV(ops[0], "uaddlv")
	if err != nil {
		return nil, err
	}

	rn, err := wantV(ops[1], "uaddlv")
	if err != nil {
		return nil, err
	}

	return (arch.Builder{}).Uaddlv(rd, rn, ops[0].Arr())
}

// newLdStruct — all structural load/store: ld1-ld4/st1-st4 (+ r forms
// ld1r/ld4r/...), reglist { vt... }{[idx]}, [rn]{, #post}. The reglist sets
// count → opcode (1→0x7, 2→0xa, 3→0x4, 4→0x0; r forms: 1r→0xc, 4r→0xe),
// .Arr() → Q/size, post-imm = regBytes*count with a check.
func newLdStruct(mnem string) func([]vOp) (Instr, error) {
	isLoad := strings.HasPrefix(mnem, "ld")
	return func(ops []vOp) (Instr, error) {
		if len(ops) < 2 || ops[0].Kind() != arch.ArmOpList {
			return nil, fmt.Errorf("%s: want { vt... }, [rn]", mnem)
		}

		arr := ops[0].Arr()
		if arr == "" && len(ops[0].List()) > 0 {
			arr = ops[0].List()[0].Arr()
		}

		q, size, err := arrQSize(arr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", mnem, err)
		}

		count := len(ops[0].List())
		if count < 1 || count > 4 {
			return nil, fmt.Errorf("%s: 1-4 registers in list", mnem)
		}

		opcode := map[int]uint32{4: 0x0, 3: 0x4, 2: 0xa, 1: 0x7}[count]
		switch {
		case strings.HasSuffix(mnem, "4r"):
			opcode = 0xe
		case strings.HasSuffix(mnem, "1r"):
			opcode = 0xc
		}

		var l uint32
		if isLoad {
			l = 1
		}

		// display name: as in ldStructDecode (ld4r/st1r when L=0)
		dname, darr, _, _ := ldStructDecode(opcode, size, q, l)
		_ = darr
		enc := uint32(0x0C000000) | q<<30 | size<<10 | opcode<<12 | l<<22
		rt0, err := wantV(arch.VOpReg(ops[0].List()[0].Reg(), "", false, 0), mnem)
		if err != nil {
			return nil, err
		}

		if !ops[1].IsMem() {
			return nil, fmt.Errorf("%s: memory operand expected", mnem)
		}

		m := ops[1].Mem()

		rn := m.Base()
		hasPost, postImm := false, uint32(0)
		if m.Post() != 0 {
			v := m.Post()
			hasPost = true
			regBytes := uint32(8)
			if q == 1 {
				regBytes = 16
			}

			postImm = regBytes * uint32(count)
			if opcode == 0xe {
				postImm = uint32(1) << size * 4
			}

			if opcode == 0xc {
				postImm = uint32(1) << size
			}

			if v != int64(postImm) {
				return nil, fmt.Errorf("%s: post-imm mismatch (want %d)", mnem, postImm)
			}

			enc = enc&^0x03800000 | 0x00800000
		}

		list := "{ " + regListStr(uint32(rt0.Num()), count) + " }"
		return arch.NewLd1(list, rn, dname, arr, "", postImm, hasPost, enc,
			uint32(rt0.Num()), count, opcode, size, q, false), nil
	}
}

// newMovSimd — mov.16b/mov.8b vd, vm (ORR-vector, Rn=31). The .Arr() suffix
// sits in the mnemonic itself (the grammar yields it whole).
func newMovSimd(arr string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("mov.%s: want vd, vm", arr)
		}

		rd, err := wantV(ops[0], "mov")
		if err != nil {
			return nil, err
		}

		rm, err := wantV(ops[1], "mov")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).MovSimd(rd, rm, arr)
	}
}

// newMovInsArm — mov.sz vd[idx], rn (INS general: inserting a GPR into a lane)
// and the scalar DUP alias mov.sz vd, vn (llvm prints mov.d/mov.s). The size
// arrives with the registration key (mov.b/h/s/d); the index — in
// ops[0].Num() (the laneIdx flag).
func newMovInsArm(elem string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		// the scalar DUP alias: mov.sz vd, vn (no lane index anywhere;
		// llvm prints these as mov.d/mov.s)
		if len(ops) == 2 && ops[0].Kind() == arch.ArmOpReg && !ops[0].LaneIdx() &&
			ops[1].Kind() == arch.ArmOpReg && ops[1].Reg() != "" &&
			ops[1].Reg()[0] == 'v' && !ops[1].LaneIdx() {
			rd, err := wantV(ops[0], "mov")
			if err != nil {
				return nil, err
			}

			rn, err := wantV(ops[1], "mov")
			if err != nil {
				return nil, err
			}

			return (arch.Builder{}).DupScalar(rd, rn, elem)
		}

		if len(ops) != 2 || ops[0].Kind() != arch.ArmOpReg {
			return nil, errors.New("mov: want vd[idx], rn")
		}

		idx := ops[0].Num()
		if idx < 0 {
			return nil, errors.New("mov: bad index")
		}

		vd, err := wantV(ops[0], "mov")
		if err != nil {
			return nil, err
		}

		rn, err := gpr(ops[1], "mov")
		if err != nil {
			return nil, err
		}

		return (arch.Builder{}).Ins(vd, uint32(idx), rn, elem)
	}
}
