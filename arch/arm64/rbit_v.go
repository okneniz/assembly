package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RbitV — rbit.8b/16b vd, vn (bit reversal per byte lane; the family's
// size bits are part of the opcode - the arrangement rides Q only).
type RbitV struct {
	base

	rd, rn string
	arr    string // 8b/16b
}

// newRbitV - the RbitV constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newRbitV(b base, rd, rn VReg, arr string) (RbitV, error) {
	err := requireArr("RbitV", arr, "8b", "16b")
	if err != nil {
		return RbitV{}, err
	}

	return RbitV{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		arr:  arr,
	}, nil
}

const rbitVEnc uint32 = 0x6E605800 // rbit vd, vn (Q=0 form)

func (i RbitV) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rbit.%s %s, %s", i.arr, i.rd, i.rn)
}

func (i RbitV) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	return writeWord(w, rbitVEnc|q<<30|size<<22|rd|rn<<5)
}

func (Builder) RbitV(rd, rn VReg, arr string) (Instr, error) {
	return newRbitV(base{}, rd, rn, arr)
}

func decodeRbitV(w uint32) (Instr, error) {
	arr := "8b"
	if w>>30&1 == 1 {
		arr = "16b"
	}

	in, err := newRbitV(
		newBase(w),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		arr,
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
