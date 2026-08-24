package arm64

// SIMD assembler constructors: the simd3 family (and/add/cmeq/addp/
// sqrshl/eor/bic/orn/eon/bics + the mov alias), simd2 (cnt/rev32/not/abs/
// rbit), shifts (shl/sri/ushr/sshr), aese/aesmc, dup,
// tbl, uaddlv, ld1 + structural ld1-ld4/st1-st4 (reglist/element
// forms), mov.16b/mov.8b. The .arr suffix carries Q/size (arrQSize —
// the inverse of decodeArrangement).

import (
	"errors"
	"fmt"
	"strings"
)

// arrQSize — (Q, size) from the .8b/.16b/.4h/.8h/.2s/.4s/.2d suffix.
func arrQSize(arr string) (q, size uint32, err error) {
	switch arr {
	case "8b":
		return 0, 0, nil
	case "16b":
		return 1, 0, nil
	case "4h":
		return 0, 1, nil
	case "8h":
		return 1, 1, nil
	case "2s":
		return 0, 2, nil
	case "4s":
		return 1, 2, nil
	case "2d":
		return 1, 3, nil
	}

	return 0, 0, fmt.Errorf("unknown arrangement %q", arr)
}

// wantV — a vector operand (a v register with or without a suffix).
func wantV(op vOp, name string) (string, error) {
	if op.reg == "" || op.reg[0] != 'v' {
		return "", fmt.Errorf("%s: vector register expected", name)
	}

	if _, err := armRegNum(op.reg); err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}

	return op.reg, nil
}

// newSimd3 — op.Arr vd, vn, vm.
func newSimd3(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", op)
		}

		if ops[0].arr == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected (.16b)", op)
		}

		q, size, err := arrQSize(ops[0].arr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if isSimd3Logical(op) && ops[0].arr != "8b" && ops[0].arr != "16b" {
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

		return Simd3{
			op:   op,
			rd:   rd,
			rn:   rn,
			rm:   rm,
			enc:  enc,
			q:    q,
			size: size,
		}, nil
	}
}

// newSimd2 — op.Arr vd, vn.
func newSimd2(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 {
			return nil, fmt.Errorf("%s: want vd, vn", op)
		}

		if ops[0].arr == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(ops[0].arr)
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

		return Simd2{
			op:   op,
			rd:   rd,
			rn:   rn,
			arr:  ops[0].arr,
			enc:  enc,
			q:    q,
			size: size,
		}, nil
	}
}

// newSimdShift — op.Arr vd, vn, #shift: immh:immb from the element width
// and the amount (ushr/sshr/sri invert — like simdShiftAmount).
func newSimdShift(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 || ops[2].kind != armOpImm {
			return nil, fmt.Errorf("%s: want vd, vn, #shift", op)
		}

		arr := ops[0].arr
		if arr == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(arr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		sh := ops[2].num
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

		return SimdShift{
			op:   op,
			rd:   rd,
			rn:   rn,
			immh: immh,
			immb: immb,
			q:    q,
			enc:  enc,
		}, nil
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

		return V1arr{
			op:  op,
			rd:  rd,
			rn:  rn,
			enc: enc,
		}, nil
	}
}

// newDupArm — dup.Arr vd, wn.
func newDupArm(ops []vOp) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("dup: want vd, wn")
	}

	if ops[0].arr == "" {
		return nil, errors.New("dup: arrangement suffix expected")
	}

	q, size, err := arrQSize(ops[0].arr)
	if err != nil {
		return nil, fmt.Errorf("dup: %w", err)
	}

	rd, err := wantV(ops[0], "dup")
	if err != nil {
		return nil, err
	}

	rn, err := wantAReg(ops[1], "dup")
	if err != nil {
		return nil, err
	}

	rdN, rnN, err := regNums2(rd, rn)
	if err != nil {
		return nil, fmt.Errorf("dup: %w", err)
	}

	return SimdCopy{
		op:     "dup",
		size:   size,
		q:      q,
		vd:     rd,
		gpr:    rn,
		enc:    0x0E000000,
		vdNum:  rdN,
		gprNum: rnN,
	}, nil
}

// newTblArm — tbl.16b vd, { vn }, vm.
func newTblArm(ops []vOp) (Instr, error) {
	if len(ops) != 3 || ops[1].kind != armOpList || len(ops[1].list) != 1 {
		return nil, errors.New("tbl: want vd, { vn }, vm")
	}

	rd, err := wantV(ops[0], "tbl")
	if err != nil {
		return nil, err
	}

	rn := ops[1].list[0].reg
	if rn == "" || rn[0] != 'v' {
		return nil, errors.New("tbl: vector in list expected")
	}

	rm, err := wantV(ops[2], "tbl")
	if err != nil {
		return nil, err
	}

	return Tbl{
		rd: rd,
		rn: rn,
		rm: rm,
	}, nil
}

// newUaddlv — uaddlv.Arr hd/sd/dd, vn (dest scalar by size).
func newUaddlv(ops []vOp) (Instr, error) {
	if len(ops) != 2 || ops[0].arr == "" {
		return nil, errors.New("uaddlv: want vd, vn with .arr")
	}

	q, size, err := arrQSize(ops[0].arr)
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

	return Uaddlv{
		rd:   scalar,
		rn:   rn,
		q:    q,
		size: size,
	}, nil
}

// newLdStruct — all structural load/store: ld1-ld4/st1-st4 (+ r forms
// ld1r/ld4r/...), reglist { vt... }{[idx]}, [rn]{, #post}. The reglist sets
// count → opcode (1→0x7, 2→0xa, 3→0x4, 4→0x0; r forms: 1r→0xc, 4r→0xe),
// .arr → Q/size, post-imm = regBytes*count with a check.
func newLdStruct(mnem string) func([]vOp) (Instr, error) {
	isLoad := strings.HasPrefix(mnem, "ld")
	return func(ops []vOp) (Instr, error) {
		if len(ops) < 2 || ops[0].kind != armOpList {
			return nil, fmt.Errorf("%s: want { vt... }, [rn]", mnem)
		}

		arr := ops[0].arr
		if arr == "" && len(ops[0].list) > 0 {
			arr = ops[0].list[0].arr
		}

		q, size, err := arrQSize(arr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", mnem, err)
		}

		count := len(ops[0].list)
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
		rt0, err := wantV(vOp{kind: armOpReg, reg: ops[0].list[0].reg}, mnem)
		if err != nil {
			return nil, err
		}

		m := ops[1].mem
		if m == nil {
			return nil, fmt.Errorf("%s: memory operand expected", mnem)
		}

		rn := m.base
		hasPost, postImm := false, uint32(0)
		if m.post != 0 {
			v := m.post
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
		return Ld1{
			regList: list,
			rn:      rn,
			name:    dname,
			arr:     arr,
			hasPost: hasPost,
			postImm: postImm,
			enc:     enc,
			rtNum:   regIndex(rt0),
			count:   count,
			opcode:  opcode,
			size:    size,
			q:       q,
		}, nil
	}
}

// newMovSimd — mov.16b/mov.8b vd, vm (ORR-vector, Rn=31). The .arr suffix
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

		return MovSimd{
			rd:  rd,
			rm:  rm,
			arr: arr,
			enc: enc,
		}, nil
	}
}

// newSimdWidenArm — s{add,sub}w{,2}.Arr vd, vn, vm (widening three-same):
// the arrangement sets the RESULT (one size wider than the Rm source).
func newSimdWidenArm(op string, enc uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 3 {
			return nil, fmt.Errorf("%s: want vd, vn, vm", op)
		}

		if ops[0].arr == "" {
			return nil, fmt.Errorf("%s: arrangement suffix expected", op)
		}

		q, size, err := arrQSize(ops[0].arr)
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

		return SimdWiden{
			op:   op,
			q:    q,
			size: size,
			rd:   rd,
			rn:   rn,
			rm:   rm,
			enc:  enc,
			rdN:  rdN,
			rnN:  rnN,
			rmN:  rmN,
		}, nil
	}
}

// newMovInsArm — mov.sz vd[idx], rn (INS general: inserting a GPR into a lane).
// The size arrives with the registration key (mov.b/h/s/d); the index — in
// ops[0].num (the laneIdx flag).
func newMovInsArm(size uint32) func([]vOp) (Instr, error) {
	return func(ops []vOp) (Instr, error) {
		if len(ops) != 2 || ops[0].kind != armOpReg {
			return nil, errors.New("mov: want vd[idx], rn")
		}

		idx := ops[0].num
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

		return SimdCopy{
			op:     "ins",
			size:   size,
			idx:    uint32(idx),
			q:      1,
			vd:     vd,
			gpr:    rn,
			enc:    0x0E000000,
			vdNum:  vdN,
			gprNum: rnN,
		}, nil
	}
}
