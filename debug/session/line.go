package session

// Line is one source line at its final address: the normalized line
// map entry of both producers (asm: the one source file of the run;
// prog: the Go file of each chain call).
type Line struct {
	File string
	Line int
	Addr uint64
	Size int
}

func NewLine(file string, line int, addr uint64, size int) Line {
	return Line{
		File: file,
		Line: line,
		Addr: addr,
		Size: size,
	}
}
