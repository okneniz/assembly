package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Udiv — udiv rd, rn, rm.
type Udiv struct {
	base

	rd, rn, rm string
}

const UdivX uint32 = 0x9AC00800

// Udiv — udiv rd, rn, rm. Register 31 reads as zr (SP/WSP are not
// allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) Udiv(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "Udiv", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Udiv", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "Udiv", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Udiv", rd, rn, rm); err != nil {
		return nil, err
	}

	return Udiv{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

func decodeUdiv(w uint32, addr uint64) Instr {
	return Udiv{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Udiv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("udiv %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Udiv) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, UdivX, 0x1AC00800)
	if err != nil {
		return 0, fmt.Errorf("udiv: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("udiv: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Udiv) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"udiv",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
