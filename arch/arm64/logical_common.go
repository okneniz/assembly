package arm64

// logImm — the common fields of logical immediate encodings.
type logImm struct {
	rd, rn string
	immr   uint32
	imms   uint32
	n      bool
	is64   bool
}

func newLogImm(rd string, rn string, immr uint32, imms uint32, n bool, is64 bool) logImm {
	return logImm{
		rd:   rd,
		rn:   rn,
		immr: immr,
		imms: imms,
		n:    n,
		is64: is64,
	}
}

func decodeLogImm(w uint32) logImm {
	is64 := w>>31&1 == 1
	return logImm{
		rd:   armRegName(w&0x1f, is64),
		rn:   armRegName(w>>5&0x1f, is64),
		imms: w >> 10 & 0x3f,
		immr: w >> 16 & 0x3f,
		n:    w>>22&1 == 1,
		is64: is64,
	}
}

// mask — the reconstructed mask value (decodeBitMasks — ARM ARM).
func (l logImm) mask() uint64 {
	return decodeBitMasks(l.n, l.immr, l.imms, l.is64)
}

// bits — the rd/rn register numbers for word assembly (immr/imms/N are
// already in the entry's match constant).
func (l logImm) bits() (uint32, uint32, error) {
	rd, rn, err := regNums2(l.rd, l.rn)
	if err != nil {
		return 0, 0, err
	}

	return rd, rn, nil
}
