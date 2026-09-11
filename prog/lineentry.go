package prog

// LineEntry - one line at its final address: its bytes start at Addr and
// occupy Size of them, appended from the Go position Pos. Labels carry
// no bytes and emit no entries; an la pair is one entry (one chain
// call).
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
