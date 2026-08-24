package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ubfm - ubfm rd, rn, #immr, #imms; aliases lsl/lsr/ubfiz/ubfx.
type Ubfm struct {
	base

	rd, rn     string
	immr, imms uint32
	isf        bool
}

const (
	ubfmX uint32 = 0xD3400000
	ubfmW uint32 = 0x53000000
)

func decodeUbfm(w uint32, addr uint64) Instr {
	return Ubfm{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		immr: w >> 16 & 0x3f,
		imms: w >> 10 & 0x3f,
		isf:  w>>31&1 == 1,
	}
}

func (i Ubfm) ObjDump(_ disasm.ViewCtx) string {
	regsize := bfmRegsize(i.rd, i.immr, i.imms)
	if i.imms != regsize-1 && i.imms+1 == i.immr { // LSL alias
		return fmt.Sprintf("lsl %s, %s, #%d", i.rd, i.rn, regsize-i.immr)
	}

	if i.imms == regsize-1 { // LSR alias
		return fmt.Sprintf("lsr %s, %s, #%d", i.rd, i.rn, i.immr)
	}

	if i.imms < i.immr { // UBFIZ: lsb = regsize-immr, width = imms+1
		return fmt.Sprintf("ubfiz %s, %s, #%d, #%d", i.rd, i.rn, regsize-i.immr, i.imms+1)
	}

	return fmt.Sprintf("ubfx %s, %s, #%d, #%d", i.rd, i.rn, i.immr, i.imms-i.immr+1)
}

func (i Ubfm) Encode(w io.Writer, pc uint64) (int64, error) {
	return bfmWrite(w, ubfmX, ubfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (i Ubfm) MarshalJSON() ([]byte, error) {
	return i.marshal("ubfm", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
