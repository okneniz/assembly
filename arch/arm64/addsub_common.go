package arm64

import (
	"fmt"
	"io"
)

// extBase — the common fields of extended add/sub.
type extBase struct {
	rdNum, rnNum, rmNum uint32
	option              string // uxtb..sxtx
	imm3                uint32
	isf                 bool
}

func newExtBase(
	rdNum uint32,
	rnNum uint32,
	rmNum uint32,
	option string,
	imm3 uint32,
	isf bool,
) extBase {
	return extBase{
		rdNum:  rdNum,
		rnNum:  rnNum,
		rmNum:  rmNum,
		option: option,
		imm3:   imm3,
		isf:    isf,
	}
}

func decodeExtBase(w uint32) extBase {
	return extBase{
		rdNum:  w & 0x1f,
		rnNum:  w >> 5 & 0x1f,
		rmNum:  w >> 16 & 0x1f,
		option: extName(w >> 13 & 7),
		imm3:   w >> 10 & 7,
		isf:    w>>31&1 == 1,
	}
}

// extMod — the ", ext #imm3" modifier (empty when omitted).
func (b extBase) extMod(omitLsl bool) string {
	if b.imm3 == 0 {
		if b.option == "uxtx" || b.option == "uxtw" {
			return ""
		}

		if omitLsl && b.option == "lsl" {
			return ""
		}
	}

	return fmt.Sprintf(", %s #%d", b.option, b.imm3)
}

// extWrite — the encoding word.
func (b extBase) extWrite(w io.Writer, matchX, matchW uint32, mnem string) (int64, error) {
	match := matchX
	if !b.isf {
		match = matchW
	}

	opt, err := extNum(b.option)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", mnem, err)
	}

	if b.imm3 > 7 {
		return 0, fmt.Errorf("%s: imm3 out of range", mnem)
	}

	return writeWord(w, match|b.rdNum|b.rnNum<<5|b.imm3<<10|opt<<13|b.rmNum<<16)
}

// extName — the extension option name (like the extOpt transform).
func extName(v uint32) string {
	names := []string{"uxtb", "uxth", "uxtw", "uxtx", "sxtb", "sxth", "sxtw", "sxtx"}
	if int(v) < len(names) {
		return names[v]
	}

	return "?"
}

// extNum — the extension option number (inverse of extName).
func extNum(name string) (uint32, error) {
	for i, n := range []string{"uxtb", "uxth", "uxtw", "uxtx", "sxtb", "sxth", "sxtw", "sxtx"} {
		if n == name {
			return uint32(i), nil
		}
	}

	return 0, fmt.Errorf("unknown extension %q", name)
}
