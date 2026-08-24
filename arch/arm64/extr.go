package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Extr — extr rd, rn, rm, #lsb; pseudo: ror rd, rn, #imm (rn == rm).
type Extr struct {
	base

	rd, rn, rm string
	lsb        uint32
}

const extrX uint32 = 0x93000000

func decodeExtr(w uint32, addr uint64) Instr {
	return Extr{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
		lsb:  w >> 10 & 0x3f,
	}
}

func (i Extr) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm {
		return fmt.Sprintf("ror %s, %s, #0x%x", i.rd, i.rn, i.lsb)
	}

	return fmt.Sprintf("extr %s, %s, %s, #0x%x", i.rd, i.rn, i.rm, i.lsb)
}

func (i Extr) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("extr: %w", err)
	}

	if i.lsb > 63 {
		return 0, fmt.Errorf("extr: lsb %#x out of range", i.lsb)
	}

	return writeWord(w, extrX|rd|rn<<5|i.lsb<<10|rm<<16)
}

func (i Extr) MarshalJSON() ([]byte, error) {
	return i.marshal("extr", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "imms": i.lsb})
}
