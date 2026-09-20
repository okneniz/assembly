// Package prog - the result shape of the per-arch chain packages
// (prog/arm64, prog/loong64, prog/riscv): the assembled bytes, the
// symbol table, the line map and the data stream - over the unit
// output's vocabulary (its Pos and LineEntry; the chains deposit those
// directly).
package prog

import (
	"github.com/okneniz/assembly/unit"
)

// Result - the assembly result: the code (the text stream), the data
// stream and its memory size (>= len(Data): the bss tail), the symbol
// table, the line map (address → source position, in address order),
// and the errors.
type Result struct {
	Code    []byte
	Data    []byte
	DataMem int
	Syms    map[string]uint64
	Lines   []unit.LineEntry
	Errs    []error
}

func NewResult(
	code []byte,
	syms map[string]uint64,
	lines []unit.LineEntry,
	errs []error,
) *Result {
	return &Result{
		Code:  code,
		Syms:  syms,
		Lines: lines,
		Errs:  errs,
	}
}
