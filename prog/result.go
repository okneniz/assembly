package prog

// Result - the assembly result: the code, the symbol table, the line map
// (address → Go source position, in address order), and the errors.
type Result struct {
	Code  []byte
	Syms  map[string]uint64
	Lines []LineEntry
	Errs  []error
}

func NewResult(
	code []byte,
	syms map[string]uint64,
	lines []LineEntry,
	errs []error,
) *Result {
	return &Result{
		Code:  code,
		Syms:  syms,
		Lines: lines,
		Errs:  errs,
	}
}
