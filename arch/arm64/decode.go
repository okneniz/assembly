package arm64

import (
	"errors"
	"fmt"
)

func decodeAbs(w uint32) (Instr, error) {
	in, err := newAbs(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeAdc(w uint32) (Instr, error) {
	in, err := newAdc(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAdd(w uint32) (Instr, error) {
	in, err := newAdd(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeAddExt(w uint32) (Instr, error) {
	in, err := newAddExt(
		numReg(w&0x1f, w>>31&1 == 1),    // rd 31 reads as sp/wsp in the plain add/sub ext words
		numReg(w>>5&0x1f, w>>31&1 == 1), // rn 31 reads as sp/wsp in the ext words
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		extName(w>>13&7),
		w>>10&7,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAddImm(w uint32) (Instr, error) {
	in, err := newAddImm(
		numReg(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAddShift(w uint32) (Instr, error) {
	in, err := newAddShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAddp(w uint32) (Instr, error) {
	in, err := newAddp(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeAddsExt(w uint32) (Instr, error) {
	in, err := newAddsExt(
		gprOf(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), // rn 31 reads as sp/wsp in the ext words
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		extName(w>>13&7),
		w>>10&7,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAddsImm(w uint32) (Instr, error) {
	in, err := newAddsImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAddsShift(w uint32) (Instr, error) {
	in, err := newAddsShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAdr(w uint32) (Instr, error) {
	raw := (w>>5&0x7ffff)<<2 | w>>29&3
	in, err := newAdr(
		gprOf(w&0x1f, true),
		signExtendN(raw, 21),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAdrp(w uint32) (Instr, error) {
	raw := (w>>5&0x7ffff)<<2 | w>>29&3
	imm21 := signExtendN(raw, 21)
	in, err := newAdrp(
		gprOf(w&0x1f, true),
		imm21,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAese(w uint32) (Instr, error) {
	return newAese(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	), nil
}

func decodeAesmc(w uint32) (Instr, error) {
	return newAesmc(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	), nil
}

func decodeAndImm(w uint32) (Instr, error) {
	in, err := newAndImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		decodeBitMasks(w>>22&1 == 1, w>>16&0x3f, w>>10&0x3f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAndShift(w uint32) (Instr, error) {
	in, err := newAndShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAndsImm(w uint32) (Instr, error) {
	in, err := newAndsImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		decodeBitMasks(w>>22&1 == 1, w>>16&0x3f, w>>10&0x3f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAndsShift(w uint32) (Instr, error) {
	in, err := newAndsShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeAsrReg(w uint32) (Instr, error) {
	in, err := newAsrReg(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeB(w uint32) (Instr, error) {
	in, err := newB(immNum(signExtendN(w&0x3ffffff, 26) * 4))
	if err != nil {
		return nil, err
	}

	return in, nil
}

// decodeBarrierOf - the decoder for one barrier mnemonic: the CRm nibble
// back to the domain (an unassigned nibble is not this instruction).
func decodeBarrierOf(name string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		domain := BarrierDomain((w >> 8) & 0xF)
		if _, ok := barrierNames[domain]; !ok {
			return nil, unknownDomain(name, domain)
		}

		return newBarrier(name, domain)
	}
}

func decodeBcondOf(cond string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return newBcond(
			cond,
			immNum(signExtendN(w>>5&0x7ffff, 19)*4),
		), nil
	}
}

func decodeBfmInstr(w uint32) (Instr, error) {
	in, err := newBfm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		w>>16&0x3f,
		w>>10&0x3f,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeBicShift(w uint32) (Instr, error) {
	in, err := newBicShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeBicsShift(w uint32) (Instr, error) {
	in, err := newBicsShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeBl(w uint32) (Instr, error) {
	in, err := newBl(immNum(signExtendN(w&0x3ffffff, 26) * 4))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeBlr(w uint32) (Instr, error) {
	in, err := newBlr(gprOf(w>>5&0x1f, true))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeBr(w uint32) (Instr, error) {
	in, err := newBr(gprOf(w>>5&0x1f, true))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCbnz(w uint32) (Instr, error) {
	in, err := newCbnz(gprOf(w&0x1f, w>>31&1 == 1), signExtendN(w>>5&0x7ffff, 19)*4)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCbz(w uint32) (Instr, error) {
	in, err := newCbz(gprOf(w&0x1f, w>>31&1 == 1), signExtendN(w>>5&0x7ffff, 19)*4)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCcmp(w uint32) (Instr, error) {
	in, err := newCcmp(
		gprOf(w>>5&0x1f, true),
		gprOf(w>>16&0x1f, true),
		w&0xf,
		condName(w>>12&0xf),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCls(w uint32) (Instr, error) {
	in, err := newCls(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeClz(w uint32) (Instr, error) {
	in, err := newClz(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCmeq(w uint32) (Instr, error) {
	in, err := newCmeq(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeCmge(w uint32) (Instr, error) {
	in, err := newCmge(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeCmtst(w uint32) (Instr, error) {
	in, err := newCmtst(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeCnt(w uint32) (Instr, error) {
	in, err := newCnt(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeCsel(w uint32) (Instr, error) {
	in, err := newCsel(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		condName(w>>12&0xf),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeCsinc(w uint32) (Instr, error) {
	ci, cerr := decodeCsel(w)
	if cerr != nil {
		return nil, cerr
	}

	c, ok := ci.(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csinc{}, nil
	}

	return Csinc{Csel: c}, nil
}

func decodeCsinv(w uint32) (Instr, error) {
	ci, cerr := decodeCsel(w)
	if cerr != nil {
		return nil, cerr
	}

	c, ok := ci.(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csinv{}, nil
	}

	return Csinv{Csel: c}, nil
}

func decodeCsneg(w uint32) (Instr, error) {
	ci, cerr := decodeCsel(w)
	if cerr != nil {
		return nil, cerr
	}

	c, ok := ci.(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csneg{}, nil
	}

	return Csneg{Csel: c}, nil
}

func decodeEonShift(w uint32) (Instr, error) {
	in, err := newEonShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeEorImm(w uint32) (Instr, error) {
	in, err := newEorImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		decodeBitMasks(w>>22&1 == 1, w>>16&0x3f, w>>10&0x3f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeEorShift(w uint32) (Instr, error) {
	in, err := newEorShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeExtr(w uint32) (Instr, error) {
	if w>>22&1 != w>>31&1 { // N must equal sf — the fixed 0/1 classes
		return nil, errors.New("extr: N != sf")
	}

	in, err := newExtr(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFadd(w uint32) (Instr, error) {
	in, err := newFadd(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcmlaElem(w uint32) (Instr, error) {
	size := w >> 22 & 3
	// the lane index lives in {L (b21), H (b11)}: .h reads both bits,
	// .s only H (see fcmla_elem.go)
	idx := w >> 11 & 1
	if size == 1 {
		idx = idx<<1 | w>>21&1
	}

	in, err := newFcmlaElem(
		w>>30&1,
		size,
		idx,
		w>>12&0xf/2*90%360,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeFcmpReg(w uint32) (Instr, error) {
	in, err := newFcmp(
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
		true,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcmpZero(w uint32) (Instr, error) {
	in, err := newFcmp(
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		FReg{},
		false,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcvt(w uint32) (Instr, error) {
	in, err := newFcvt(
		newFReg(uint8(w&0x1f), w>>15&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcvtzs(w uint32) (Instr, error) {
	in, err := newFcvtzs(
		gprOf(w&0x1f, w>>31&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFcvtzu(w uint32) (Instr, error) {
	in, err := newFcvtzu(
		gprOf(w&0x1f, w>>31&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFdiv(w uint32) (Instr, error) {
	in, err := newFdiv(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmadd(w uint32) (Instr, error) {
	is64 := w>>22&1 == 1
	in, err := newFmadd(
		newFReg(uint8(w&0x1f), is64),
		newFReg(uint8(w>>5&0x1f), is64),
		newFReg(uint8(w>>16&0x1f), is64),
		newFReg(uint8(w>>10&0x1f), is64),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmax(w uint32) (Instr, error) {
	in, err := newFmax(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmin(w uint32) (Instr, error) {
	in, err := newFmin(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmlaElem(w uint32) (Instr, error) {
	in, err := newFmlaElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeFmlsElem(w uint32) (Instr, error) {
	in, err := newFmlsElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeFmov(w uint32) (Instr, error) {
	in, err := newFmov(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmovFromGpr(w uint32) (Instr, error) {
	is64 := w>>22&1 == 1
	in, err := newFmovFromGpr(
		newFReg(uint8(w&0x1f), is64),
		gprOf(w>>5&0x1f, is64),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmovImmOf(isS bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		imm8 := w >> 13 & 0xff
		rd := newFReg(uint8(w&0x1f), !isS)
		if isS {
			v := vfpExpandImm32(imm8)
			return newFmovImm(rd, float64(v), fmt.Sprintf("%.8f", v)), nil
		}

		v := vfpExpandImm64(imm8)
		return newFmovImm(rd, v, fmt.Sprintf("%.8f", v)), nil
	}
}

func decodeFmovToGpr(w uint32) (Instr, error) {
	is64 := w>>22&1 == 1
	in, err := newFmovToGpr(
		gprOf(w&0x1f, is64),
		newFReg(uint8(w>>5&0x1f), is64),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmul(w uint32) (Instr, error) {
	in, err := newFmul(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFmulElem(w uint32) (Instr, error) {
	in, err := newFmulElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeFmulxElem(w uint32) (Instr, error) {
	in, err := newFmulxElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeFneg(w uint32) (Instr, error) {
	in, err := newFneg(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFnmsub(w uint32) (Instr, error) {
	is64 := w>>22&1 == 1
	in, err := newFnmsub(
		newFReg(uint8(w&0x1f), is64),
		newFReg(uint8(w>>5&0x1f), is64),
		newFReg(uint8(w>>16&0x1f), is64),
		newFReg(uint8(w>>10&0x1f), is64),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeFsub(w uint32) (Instr, error) {
	in, err := newFsub(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeLd1Of(enc uint32) func(w uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		opcode, size, q := w>>12&0xf, w>>10&3, w>>30&1
		name, arr, count, isElem := ldStructDecode(opcode, size, q, w>>22&1)
		rt := w & 0x1f
		list := "{ " + regListStr(rt, count) + " }"
		if w>>24&1 == 1 {
			// single-structure (bits[29:24]=001101): the layout is DIFFERENT,
			// not like multi — size lives in bits[15:14], the element index =
			// Q with the truncated bits [13:10], bits[15:14]=11 → ld1r.
			name, arr, list, size = ldElemDecode(w, q, rt)
		}

		i := Ld1{
			regList: list,
			rn:      regNameXSP(w >> 5 & 0x1f),
			name:    name,
			arr:     arr,
			enc:     enc,
			rtNum:   rt,
			count:   count,
			opcode:  opcode,
			size:    size,
			q:       q,
			isElem:  isElem,
		}
		// post-index: bit23 (Rm=11111 → immediate, otherwise post-index by
		// register — objdump prints the register name); bit22 = L.
		if w>>23&1 == 1 {
			if rm := w >> 16 & 0x1f; rm != 0x1f {
				i.hasPost = true
				i.postReg = regNameX(rm)
				return i, nil
			}

			elemBytes := uint32(1) << size
			regBytes := uint32(8)
			if q == 1 {
				regBytes = 16
			}

			postImm := regBytes * uint32(count)
			if opcode == 0xe {
				postImm = elemBytes * 4
			}

			if opcode == 0xc {
				postImm = elemBytes
			}

			if w>>24&1 == 1 {
				postImm = elemBytes
				if name == "ld1r" {
					postImm = regBytes // ld1r writes the whole register
				}

				if name == "ld4r" {
					postImm = regBytes // ld4r: post = one register (objdump: #16 for .4s)
				}
			}

			i.hasPost, i.postImm = true, postImm
		}

		return i, nil
	}
}

func decodeLdarOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldar{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeLdarbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldarb{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeLdaxrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldaxr{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeLdaxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Ldaxrb{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeLdpOf(enc uint32, scale uint32, x64 bool, rtKind string) func(uint32) (Instr, error) {
	kind := rtKind
	if kind == "" {
		if x64 {
			kind = "x"
		} else {
			kind = "w"
		}
	}

	return func(w uint32) (Instr, error) {
		rt, rt2, rn, k, off, load := pairDecode(w, scale, kind)
		if !load {
			return newStpBase(
				newPairBase(rt, rt2, rn, k, off, scale, enc&^(1<<22)),
			), nil
		}

		return Ldp{
			pairBase: newPairBase(rt, rt2, rn, k, off, scale, enc|1<<22),
		}, nil
	}
}

func decodeLdpsw(w uint32) (Instr, error) {
	in, err := newLdpsw(
		gprOf(w&0x1f, true),
		gprOf(w>>10&0x1f, true),
		xspOf(w>>5&0x1f),
		Off(signExtendN(w>>15&0x7f, 7)<<2),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeLdrOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldr{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdrbOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameW(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldrb{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdrhOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameW(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldrh{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdrsb(w uint32) (Instr, error) {
	in, err := newLdrsb(
		gprOf(w&0x1f, true),
		xspOf(w>>5&0x1f),
		Off(int64(w>>10&0xfff)<<0),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeLdrsh(w uint32) (Instr, error) {
	in, err := newLdrsh(
		gprOf(w&0x1f, true),
		xspOf(w>>5&0x1f),
		Off(int64(w>>10&0xfff)<<1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeLdrswOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameX(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldrsw{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdurOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldur{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdurbOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldurb{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLdurhOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldurh{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeLslReg(w uint32) (Instr, error) {
	in, err := newLslReg(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeLsrReg(w uint32) (Instr, error) {
	in, err := newLsrReg(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMadd(w uint32) (Instr, error) {
	in, err := newMadd(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		gprOf(w>>10&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMlaElem(w uint32) (Instr, error) {
	in, err := newMlaElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeMlsElem(w uint32) (Instr, error) {
	in, err := newMlsElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeMovk(w uint32) (Instr, error) {
	in, err := newMovk(
		gprOf(w&0x1f, w>>31&1 == 1),
		imm16Of(w>>5&0xffff),
		hwOf(w>>21&0x3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMovn(w uint32) (Instr, error) {
	in, err := newMovn(
		gprOf(w&0x1f, w>>31&1 == 1),
		imm16Of(w>>5&0xffff),
		hwOf(w>>21&0x3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMovz(w uint32) (Instr, error) {
	in, err := newMovz(
		gprOf(w&0x1f, w>>31&1 == 1),
		imm16Of(w>>5&0xffff),
		hwOf(w>>21&0x3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMrs(w uint32) (Instr, error) {
	in, err := newMrs(gprOf(w&0x1f, true), sysRegName(w>>5&0x7fff))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMsr(w uint32) (Instr, error) {
	in, err := newMsr(
		sysRegName(w>>5&0x7fff),
		gprOf(w&0x1f, true),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeMsub(w uint32) (Instr, error) {
	mi, merr := decodeMadd(w)
	if merr != nil {
		return nil, merr
	}

	m, ok := mi.(Madd)
	if !ok {
		// decodeMadd always returns Madd; this branch guards against schema desynchronization
		return Msub{}, nil
	}

	return Msub{Madd: m}, nil
}

func decodeMulElem(w uint32) (Instr, error) {
	in, err := newMulElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeNop(w uint32) (Instr, error) {
	in, err := newNop()
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeNot(w uint32) (Instr, error) {
	in, err := newNot(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeOrnShift(w uint32) (Instr, error) {
	in, err := newOrnShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeOrrImm(w uint32) (Instr, error) {
	in, err := newOrrImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		decodeBitMasks(w>>22&1 == 1, w>>16&0x3f, w>>10&0x3f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeOrrShift(w uint32) (Instr, error) {
	in, err := newOrrShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodePrfm(w uint32) (Instr, error) {
	in, err := newPrfm(xspOf(w >> 5 & 0x1f))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRbit(w uint32) (Instr, error) {
	in, err := newRbit(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRbitV(w uint32) (Instr, error) {
	arr := "8b"
	if w>>30&1 == 1 {
		arr = "16b"
	}

	in, err := newRbitV(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		arr,
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeRet(w uint32) (Instr, error) {
	in, err := newRet(gprOf(w>>5&0x1f, true))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRev(w uint32) (Instr, error) {
	in, err := newRev(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRev16(w uint32) (Instr, error) {
	in, err := newRev16(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRev32(w uint32) (Instr, error) {
	in, err := newRev32(gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeRev32V(w uint32) (Instr, error) {
	in, err := newRev32V(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeRorReg(w uint32) (Instr, error) {
	in, err := newRorReg(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSaddw(w uint32) (Instr, error) {
	in, err := newSaddw(
		w>>30&1,
		w>>22&3,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSbfm(w uint32) (Instr, error) {
	in, err := newSbfm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		w>>16&0x3f,
		w>>10&0x3f,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeScvtf(w uint32) (Instr, error) {
	in, err := newScvtf(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSdiv(w uint32) (Instr, error) {
	in, err := newSdiv(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeShl(w uint32) (Instr, error) {
	in, err := newShl(
		w>>30&1,
		w>>19&0xf,
		w>>16&7,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

// decodeSimdDupElem — DUP (element): 0x0E000400/0xBFE0FC00.
func decodeSimdDupElem(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	// the size is the lowest set bit (0=b..3=d); unencodable sizes go to
	// the .word fallback (the schema mask does not express this)
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	in, err := newDupElem(
		w>>30&1,
		size,
		imm5>>(size+1),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// decodeSimdDupScalar — the scalar DUP alias: 0x5E000400/0xFFE0FC00;
// imm5 is the size one-hot plus the lane index above it (llvm prints
// mov b/h/s/d<n>, vn.sz[idx]).
func decodeSimdDupScalar(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	in, err := newDupScalar(
		size,
		imm5>>(size+1),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// decodeSimdInsElem — INS (element): 0x6E000400/0xFFE08400.
func decodeSimdInsElem(w uint32) (Instr, error) {
	imm5 := w >> 16 & 0x1f
	if imm5 == 0 || bitsCtz(imm5) > 3 {
		return decodeUnknown(w)
	}

	size := uint32(bitsCtz(imm5))
	in, err := newInsElem(
		size,
		imm5>>(size+1),
		w>>11&0xf>>size,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSmlalElem(w uint32) (Instr, error) {
	in, err := newSmlalElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSmlslElem(w uint32) (Instr, error) {
	in, err := newSmlslElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSmulh(w uint32) (Instr, error) {
	in, err := newSmulh(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSmullElem(w uint32) (Instr, error) {
	in, err := newSmullElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqdmlalElem(w uint32) (Instr, error) {
	in, err := newSqdmlalElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqdmlslElem(w uint32) (Instr, error) {
	in, err := newSqdmlslElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqdmulhElem(w uint32) (Instr, error) {
	in, err := newSqdmulhElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqdmullElem(w uint32) (Instr, error) {
	in, err := newSqdmullElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqrdmlahElem(w uint32) (Instr, error) {
	in, err := newSqrdmlahElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqrdmlshElem(w uint32) (Instr, error) {
	in, err := newSqrdmlshElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqrdmulhElem(w uint32) (Instr, error) {
	in, err := newSqrdmulhElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSqrshl(w uint32) (Instr, error) {
	in, err := newSqrshl(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSri(w uint32) (Instr, error) {
	in, err := newSri(
		w>>30&1,
		w>>19&0xf,
		w>>16&7,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSshr(w uint32) (Instr, error) {
	in, err := newSshr(
		w>>30&1,
		w>>19&0xf,
		w>>16&7,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeSsubw(w uint32) (Instr, error) {
	in, err := newSsubw(
		w>>30&1,
		w>>22&3,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeStlrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlr{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeStlrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlrb{
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}, nil
	}
}

func decodeStlxrOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlxr{
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func decodeStlxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stlxrb{
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func decodeStrOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Str{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeStrbOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameW(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Strb{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeStrhOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameW(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Strh{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeSturOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Stur{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeSturbOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Sturb{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeSturhOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Sturh{
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func decodeStxrbOf(enc uint32, x64 bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return Stxrb{
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}, nil
	}
}

func decodeSubExt(w uint32) (Instr, error) {
	in, err := newSubExt(
		numReg(w&0x1f, w>>31&1 == 1),    // rd 31 reads as sp/wsp in the plain add/sub ext words
		numReg(w>>5&0x1f, w>>31&1 == 1), // rn 31 reads as sp/wsp in the ext words
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		extName(w>>13&7),
		w>>10&7,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSubImm(w uint32) (Instr, error) {
	in, err := newSubImm(
		numReg(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSubShift(w uint32) (Instr, error) {
	in, err := newSubShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSubsExt(w uint32) (Instr, error) {
	in, err := newSubsExt(
		gprOf(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), // rn 31 reads as sp/wsp in the ext words
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		extName(w>>13&7),
		w>>10&7,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSubsImm(w uint32) (Instr, error) {
	in, err := newSubsImm(
		gprOf(w&0x1f, w>>31&1 == 1),
		numReg(w>>5&0x1f, w>>31&1 == 1), imm12Of(w>>10&0xfff), sh12Of(w>>22&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSubsShift(w uint32) (Instr, error) {
	in, err := newSubsShift(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		imm6Of(w>>10&0x3f),
		shiftOf(w>>22&3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeSysFixedOf(name, group string, enc uint32) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return sysFixed{
			name:  name,
			group: group,
			enc:   enc,
		}, nil
	}
}

// decodeSysImmOf — the imm16 sits at bits 5..20 for every member of the
// family (udf/brk/hlt/svc/hvc/smc).
func decodeSysImmOf(name string, enc uint32) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		return sysImm{
			name:  name,
			imm16: w >> 5 & 0xffff,
			enc:   enc,
			shift: 5,
		}, nil
	}
}

func decodeTbl(w uint32) (Instr, error) {
	return newTbl(
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	), nil
}

func decodeTbzOf(isTbnz bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		x64 := w>>31&1 == 1
		return Tbz{
			rt:     armRegName(w&0x1f, x64),
			bit:    w>>19&0x1f | w>>26&0x20,
			off:    immNum(signExtendN(w>>5&0x3fff, 14) * 4),
			isTbnz: isTbnz,
		}, nil
	}
}

func decodeUaddlv(w uint32) (Instr, error) {
	rd := vReg(w & 0x1f)
	size := w >> 22 & 3
	switch size {
	case 0:
		rd = fmt.Sprintf("h%d", regIndex(rd))
	case 1:
		rd = fmt.Sprintf("s%d", regIndex(rd))
	}

	return Uaddlv{
		rd:   rd,
		rn:   vReg(w >> 5 & 0x1f),
		q:    w >> 30 & 1,
		size: size,
	}, nil
}

func decodeUaddw(w uint32) (Instr, error) {
	in, err := newUaddw(
		w>>30&1,
		w>>22&3,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeUbfm(w uint32) (Instr, error) {
	in, err := newUbfm(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		w>>16&0x3f,
		w>>10&0x3f,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeUcvtf(w uint32) (Instr, error) {
	in, err := newUcvtf(
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeUdiv(w uint32) (Instr, error) {
	in, err := newUdiv(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeUmlalElem(w uint32) (Instr, error) {
	in, err := newUmlalElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeUmlslElem(w uint32) (Instr, error) {
	in, err := newUmlslElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeUmulh(w uint32) (Instr, error) {
	in, err := newUmulh(
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func decodeUmullElem(w uint32) (Instr, error) {
	in, err := newUmullElem(
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeUshr(w uint32) (Instr, error) {
	in, err := newUshr(
		w>>30&1,
		w>>19&0xf,
		w>>16&7,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}

func decodeUsubw(w uint32) (Instr, error) {
	in, err := newUsubw(
		w>>30&1,
		w>>22&3,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
