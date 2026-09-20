package loong64

// Binary - a built program: the unit output of the chain, ready to be
// resolved at any base. The assemble entry point is the unit's resolve
// phase adapted to the chain's result shape - the chain itself never
// resolves anything.

import (
	"github.com/okneniz/assembly/prog"
	"github.com/okneniz/assembly/unit"
)

// Binary - a built program: the deposited unit output, ready to be
// assembled at any base.
type Binary struct {
	u *unit.Unit
}

// Assemble - encode the program at base as one flat stream: labels become
// absolute addresses, label-directed lines receive their targets. Returns
// the result: the code, the symbol table, the line map, and the assembly
// errors (undefined labels and entries, encode failures).
func (b *Binary) Assemble(base uint64) *prog.Result {
	f := b.u.Resolve(func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return base, base
	})

	code, codeErr := f.EncodeText()

	errs := f.Errs
	if codeErr != nil {
		errs = append(errs, codeErr)
	}

	return prog.NewResult(code, f.Syms, f.Lines, errs)
}
