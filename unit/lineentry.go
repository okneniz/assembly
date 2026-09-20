package unit

// LineEntry - one record at its final address: its bytes start at Addr
// and occupy Size of them, deposited from Pos (the debug view of the
// program). Labels carry no bytes and emit no entries; a deferred pair
// (an la) is one entry - one deposit.
type LineEntry struct {
	Addr uint64
	Size int
	Pos  Pos
}

func NewLineEntry(addr uint64, size int, pos Pos) LineEntry {
	return LineEntry{
		Addr: addr,
		Size: size,
		Pos:  pos,
	}
}
