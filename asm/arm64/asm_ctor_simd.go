package arm64

// SIMD assembler constructors: the simd3 family (and/add/cmeq/addp/
// sqrshl/eor/bic/orn/eon/bics + the mov alias), simd2 (cnt/rev32/not/abs/
// rbit), shifts (shl/sri/ushr/sshr), aese/aesmc, dup,
// tbl, uaddlv, ld1 + structural ld1-ld4/st1-st4 (reglist/element
// forms), mov.16b/mov.8b. The .Arr() suffix carries Q/size (arrQSize —
// the inverse of decodeArrangement).

import (
	"errors"
	"fmt"
	"strings"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// wantV — a vector operand (a v register with or without a suffix).
func wantV(op vOp, name string) (string, error) {
	if op.Reg() == "" || op.Reg()[0] != 'v' {
		return "", fmt.Errorf("%s: vector register expected", name)
	}

	if _, err := armRegNum(op.Reg()); err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return op.Reg(), nil
}

// newSimd3 — op.Arr vd, vn, vm.
func newSimd3(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", op)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected (.16b)", op)
		}

		q, size, err := arrQSize(ops[0].Arr())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if isSimd3Logical(op) && ops[0].Arr() != "8b" && ops[0].Arr() != "16b" {
			// in the logical group bits 23:22 are an opcode, not an arrangement
			return nil, fmt.Errorf("%s: only .8b/.16b arrangements", op)
		}

		rd, err := wantV(ops[0], op)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		rm, err := wantV(ops[2], op)
		if err != nil {
			return nil, err
		}

		return arch.NewSimd3(op, rd, rn, rm, enc, q, size), nil
	}
}

// newSimd2 — op.Arr vd, vn.
func newSimd2(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want vd, vn", op)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(ops[0].Arr())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rd, err := wantV(ops[0], op)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		return arch.NewSimd2(op, rd, rn, ops[0].Arr(), enc, q, size), nil
	}
}

// newSimdShift — op.Arr vd, vn, #shift: immh:immb from the element width
// and the amount (ushr/sshr/sri invert — like simdShiftAmount).
func newSimdShift(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 || ops[2].Kind() != arch.ArmOpImm {
			return nil, fmt.Errorf("%s: want vd, vn, #shift", op)
		}

		arr := ops[0].Arr()
		if arr == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(arr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		sh := ops[2].Num()
		if sh < 0 {
			return nil, fmt.Errorf("%s: bad shift", op)
		}

		imm := uint32(sh)
		// inverse of simdShiftAmount: ushr/sshr/sri store esize-shift
		if op == "ushr" || op == "sshr" || op == "sri" {
			imm = (uint32(8) << size) - imm
		}

		immh := imm >> 3
		immb := imm & 7
		if immh == 0 || immh > 0xf {
			return nil, fmt.Errorf("%s: shift out of range", op)
		}

		if immh>>3 == 1 || size == 3 {
			// size is confirmed by immh: a mismatch is an error
			if immh>>3 != 1 && size == 3 {
				return nil, fmt.Errorf("%s: immh/arr mismatch", op)
			}
		}

		rd, err := wantV(ops[0], op)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		return arch.NewSimdShift(op, rd, rn, immh, immb, q, enc), nil
	}
}

// newAes — aese/aesmc vd, vn (.16b implicit).
func newAes(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want vd, vn", op)
		}

		rd, err := wantV(ops[0], op)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		return arch.NewV1arr(op, rd, rn, enc), nil
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

	q, size, err := arrQSize(ops[0].Arr())
	if err != nil {
		return nil, fmt.Errorf("dup: %w", err)
	}

	rd, err := wantV(ops[0], "dup")
	if err != nil {
		return nil, err
	}

	rdN, err := armRegNum(rd)
	if err != nil {
		return nil, fmt.Errorf("dup: %w", err)
	}

	// DUP (element): the source is a lane of a vector register
	if ops[1].Kind() == arch.ArmOpReg && ops[1].Reg() != "" && ops[1].Reg()[0] == 'v' {
		idx := ops[1].Num()
		if !ops[1].LaneIdx() || idx < 0 || idx >= 16>>size {
			return nil, fmt.Errorf("dup: want vd.Arr, vn[idx] (lane 0..%d)",
				(16>>size)-1)
		}

		rn, err := wantV(ops[1], "dup")
		if err != nil {
			return nil, err
		}

		rnN, err := armRegNum(rn)
		if err != nil {
			return nil, fmt.Errorf("dup: %w", err)
		}

		return arch.NewDupElem("dup", size, uint32(idx), 0, q, rd, rn, rdN, rnN), nil
	}

	rn, err := wantAReg(ops[1], "dup")
	if err != nil {
		return nil, err
	}

	rnN, err := armRegNum(rn)
	if err != nil {
		return nil, fmt.Errorf("dup: %w", err)
	}

	return arch.NewSimdCopyGPR("dup", rd, rn, size, 0, q, rdN, rnN, false), nil
}

// newInsElemArm — INS (element): ins.sz vd[idx], vn[idx]. (The GPR-source
// form has no spelling here: llvm prints it as the mov alias, so the
// canonical input is mov.sz vd[idx], wn — newMovInsArm.)
func newInsElemArm(size uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 || !ops[0].LaneIdx() || !ops[1].LaneIdx() {
			return nil, errors.New("ins: want vd[idx], vn[idx]")
		}

		maxIdx := int64(16 >> size)
		if ops[0].Num() < 0 || ops[0].Num() >= maxIdx ||
			ops[1].Num() < 0 || ops[1].Num() >= maxIdx {
			return nil, fmt.Errorf("ins: lane index out of range (0..%d)", maxIdx-1)
		}

		rd, err := wantV(ops[0], "ins")
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], "ins")
		if err != nil {
			return nil, err
		}

		rdN, err := armRegNum(rd)
		if err != nil {
			return nil, fmt.Errorf("ins: %w", err)
		}

		rnN, err := armRegNum(rn)
		if err != nil {
			return nil, fmt.Errorf("ins: %w", err)
		}

		return arch.NewDupElem("ins", size,
			uint32(ops[0].Num()), uint32(ops[1].Num()), 0, rd, rn, rdN, rnN), nil
	}
}

// newSmovUmovArm — SMOV/UMOV: op wd, vn.sz[idx] (the element size is the
// source suffix; an x destination sets Q).
func newSmovUmovArm(op string) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 || ops[1].Kind() != arch.ArmOpReg || !ops[1].LaneIdx() {
			return nil, fmt.Errorf("%s: want wd, vn.sz[idx]", op)
		}

		var size uint32
		switch ops[1].Arr() {
		case "b":
			size = 0
		case "h":
			size = 1
		case "s":
			size = 2
		case "d":
			size = 3
		default:
			return nil, fmt.Errorf("%s: element size suffix (b/h/s/d) expected", op)
		}

		maxIdx := int64(16 >> size)
		if ops[1].Num() < 0 || ops[1].Num() >= maxIdx {
			return nil, fmt.Errorf("%s: lane index out of range (0..%d)", op, maxIdx-1)
		}

		gpr, err := wantAReg(ops[0], op)
		if err != nil {
			return nil, err
		}

		vd, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		gprN, err := armRegNum(gpr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		vdN, err := armRegNum(vd)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		var q uint32
		if gpr[0] == 'x' {
			q = 1
		}

		if op == "smov" && size == 3 {
			return nil, errors.New("smov: .d elements are not allowed")
		}

		// the element must fill (smov) or fit (umov) the destination
		// register: UMOV takes .d only into x and .b/.h/.s only into w
		// (llvm: "invalid operand" otherwise)
		if op == "umov" && (size == 3) != (gpr[0] == 'x') {
			return nil, fmt.Errorf("%s: %s destination expected for .%s elements",
				op, map[bool]string{true: "x", false: "w"}[size == 3],
				[...]string{"b", "h", "s", "d"}[size])
		}

		return arch.NewSimdCopyGPR(
			op,
			vd,
			gpr,
			size,
			uint32(ops[1].Num()),
			q,
			vdN,
			gprN,
			true,
		), nil
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

	rn := ops[1].List()[0].Reg()
	if rn == "" || rn[0] != 'v' {
		return nil, errors.New("tbl: vector in list expected")
	}

	rm, err := wantV(ops[2], "tbl")
	if err != nil {
		return nil, err
	}

	return arch.NewTbl(rd, rn, rm), nil
}

// newUaddlv — uaddlv.Arr hd/sd/dd, vn (dest scalar by size).
func newUaddlv(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].Arr() == "" {
		return nil, errors.New("uaddlv: want vd, vn with .Arr()")
	}

	q, size, err := arrQSize(ops[0].Arr())
	if err != nil {
		return nil, fmt.Errorf("uaddlv: %w", err)
	}

	rd, err := wantV(ops[0], "uaddlv")
	if err != nil {
		return nil, err
	}

	scalar := rd
	switch size {
	case 0:
		scalar = fmt.Sprintf("h%d", regIndex(rd))
	case 1:
		scalar = fmt.Sprintf("s%d", regIndex(rd))
	}

	rn, err := wantV(ops[1], "uaddlv")
	if err != nil {
		return nil, err
	}

	return arch.NewUaddlv(scalar, rn, q, size), nil
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

		list := "{ " + regListStr(regIndex(rt0), count) + " }"
		return arch.NewLd1(list, rn, dname, arr, "", postImm, hasPost, enc,
			regIndex(rt0), count, opcode, size, q, false), nil
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

		enc := uint32(0x4EA01C00)
		if arr == "8b" {
			enc &^= 1 << 30
		}

		return arch.NewMovSimd(rd, rm, arr, enc), nil
	}
}

// newSimdWidenArm — s{add,sub}w{,2}.Arr vd, vn, vm (widening three-same):
// the arrangement sets the RESULT (one size wider than the Rm source).
func newSimdWidenArm(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", op)
		}

		if ops[0].Arr() == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(ops[0].Arr())
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if size == 0 {
			return nil, fmt.Errorf("%s: arrangement too narrow", op)
		}

		size-- // the source is one size narrower than the result
		rd, err := wantV(ops[0], op)
		if err != nil {
			return nil, err
		}

		rn, err := wantV(ops[1], op)
		if err != nil {
			return nil, err
		}

		rm, err := wantV(ops[2], op)
		if err != nil {
			return nil, err
		}

		rdN, rnN, rmN, err := regNums3(rd, rn, rm)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		return arch.NewSimdWiden(op, q, size, rd, rn, rm, enc, rdN, rnN, rmN), nil
	}
}

// newMovInsArm — mov.sz vd[idx], rn (INS general: inserting a GPR into a lane)
// and the scalar DUP alias mov.sz vd, vn (llvm prints mov.d/mov.s). The size
// arrives with the registration key (mov.b/h/s/d); the index — in
// ops[0].Num() (the laneIdx flag).
func newMovInsArm(size uint32) func([]vOp) (Instr, error) {
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

			rdN, err := armRegNum(rd)
			if err != nil {
				return nil, fmt.Errorf("mov: %w", err)
			}

			rnN, err := armRegNum(rn)
			if err != nil {
				return nil, fmt.Errorf("mov: %w", err)
			}

			return arch.NewDupElem("mov", size, 0, 0, 0, rd, rn, rdN, rnN), nil
		}

		if len(ops) != 2 || ops[0].Kind() != arch.ArmOpReg {
			return nil, errors.New("mov: want vd[idx], rn")
		}

		idx := ops[0].Num()
		if idx < 0 {
			return nil, errors.New("mov: bad index")
		}

		if idx >= 16>>size {
			return nil, fmt.Errorf("mov: index %d out of range", idx)
		}

		vd, err := wantV(ops[0], "mov")
		if err != nil {
			return nil, err
		}

		rn, err := wantAReg(ops[1], "mov")
		if err != nil {
			return nil, err
		}

		vdN, err := armRegNum(vd)
		if err != nil {
			return nil, err
		}

		rnN, err := armRegNum(rn)
		if err != nil {
			return nil, err
		}

		return arch.NewSimdCopyGPR("ins", vd, rn, size, uint32(idx), 1, vdN, rnN, false), nil
	}
}
