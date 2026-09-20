package unit

// Pos is the source position a deposit reports: the file and line the
// record came from - a compiler's source, an .s file, a Go chain call;
// the origin is always the producer's.
type Pos struct {
	File string
	Line int
}

func NewPos(file string, line int) Pos {
	return Pos{
		File: file,
		Line: line,
	}
}
