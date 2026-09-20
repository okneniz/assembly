package arm64

// Binary - a built program: the unit output of the chain and its entry,
// ready to be resolved (or transformed) at any base. The two assemble
// entry points are the unit's resolve phase adapted to the chain's
// result shape - the chain itself never resolves anything.

import (
	"github.com/okneniz/assembly/prog"
	"github.com/okneniz/assembly/unit"
)

// Binary - a built program: the deposited unit output, ready to be
// assembled at any base.
type Binary struct {
	u       *unit.Unit
	streams bool // the program ever switched to the data stream
}

// Assemble - encode the program at base as one flat stream: labels become
// absolute addresses, label-directed lines receive their targets. Returns
// the result: the code, the symbol table, the line map, and the assembly
// errors (undefined labels and entries, encode failures). A program that
// switched to the data stream needs AssembleLayout instead.
func (b *Binary) Assemble(base uint64) *prog.Result {
	if b.streams {
		return prog.NewResult(nil, nil, nil, []error{errNeedsLayout})
	}

	return b.AssembleLayout(flatPlace(base))
}

// AssembleLayout - encode the program with its two streams placed by the
// policy: place sees the stream sizes (the data memory size includes the
// bss tail) and returns the base address of each stream. Labels resolve
// against those bases, so text may reference data statics across the page
// gap through the label-directed lines (La, Adrp, ...) without knowing
// the layout. file.MachoPlaceStreams is the Mach-O policy - the same
// numbers the writer places the bytes at.
func (b *Binary) AssembleLayout(place unit.Place) *prog.Result {
	return result(b.u.Resolve(place))
}

// flatPlace places a single text stream at base (the data stream is
// unused: the flat Assemble rejects programs that have one).
func flatPlace(base uint64) unit.Place {
	return func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return base, base
	}
}

// result adapts the unit output to the chain's result shape: the streams
// materialize to bytes, the line map keeps the chain's position type.
func result(f *unit.Fixed) *prog.Result {
	code, codeErr := f.EncodeText()
	data, dataErr := f.EncodeData()

	errs := f.Errs
	for _, err := range []error{codeErr, dataErr} {
		if err != nil {
			errs = append(errs, err)
		}
	}

	res := prog.NewResult(code, f.Syms, f.Lines, errs)
	res.Data = data
	res.DataMem = f.DataMem
	return res
}

// errNeedsLayout - the flat Assemble met a program with two streams.
var errNeedsLayout = errNeedsLayoutType{}

type errNeedsLayoutType struct{}

func (errNeedsLayoutType) Error() string {
	return "the data stream needs AssembleLayout"
}
