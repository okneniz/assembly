// Package prog - the shared vocabulary of the per-arch chain packages
// (prog/arm64, prog/loong64, prog/riscv): the source position a chain
// call reports, the line map entry, and the assembly result.
package prog

// Pos - the source position of a chain call: the Go file and line the
// line came from (the file may differ per line: helpers in another file
// append to the same program).
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
