package arm64

// The common encoding word of the bitfield family (ubfm/sbfm/bfm).

import (
	"errors"
	"io"
)

// bfmWrite — the common encoding word of the bitfield family.
func bfmWrite(
	w io.Writer,
	matchX, matchW uint32,
	isf bool,
	rd, rn string,
	immr, imms uint32,
) (int64, error) {
	match := matchX
	if !isf {
		match = matchW
	}

	r, n, err := regNums2(rd, rn)
	if err != nil {
		return 0, err
	}

	if immr > 63 || imms > 63 {
		return 0, errors.New("immr/imms out of range")
	}

	return writeWord(w, match|r|n<<5|imms<<10|immr<<16)
}
