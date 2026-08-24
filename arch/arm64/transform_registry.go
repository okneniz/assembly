package arm64

// Registry of inverse bit-field transforms (final representation -> bits):
// used by the assembler's legacy path (BuildLegacy) when assembling a word
// from the decode table. Implementations are in asm_inverse.go.

import "fmt"

// inverseTransformFunc - final representation (register name, condition
// name, literal) -> field bits. Used by the assembler's legacy path
// (BuildLegacy); contextual transforms (branches with addr) are inverted
// by the format handlers, not by the registry.
type inverseTransformFunc func(any) (uint32, error)

var inverseTransforms = map[string]inverseTransformFunc{
	"regX":        invRegX,
	"regW":        invRegW,
	"regXSP":      invRegXSP,
	"regWSP":      invRegWSP,
	"regV":        invRegV,
	"fpRegD":      invFpReg('d'),
	"fpRegS":      invFpReg('s'),
	"intX":        invRegX,
	"intW":        invRegW,
	"regIndex":    invInt,
	"cond":        invCond,
	"sysreg":      invSysReg,
	"shiftName":   invShiftName,
	"immShiftLSL": invImmShiftLSL,
	"extOpt":      invExtOpt,
	"lsOpt":       invLsOpt,
}

// applyInverseTransform applies the inverse transform name to the value v.
// An empty name means "no transform": v must be a number.
// An unknown name is an error (a silent fallback is unsafe for the
// assembler).
func applyInverseTransform(name string, v any) (uint32, error) {
	if name == "" {
		switch n := v.(type) {
		case int:
			return uint32(n), nil
		case int64:
			return uint32(n), nil
		case uint32:
			return n, nil
		case uint64:
			return uint32(n), nil
		}

		return 0, fmt.Errorf("numeric field expected, got %T", v)
	}

	f, ok := inverseTransforms[name]
	if !ok {
		return 0, fmt.Errorf("no inverse transform %q", name)
	}

	return f(v)
}
