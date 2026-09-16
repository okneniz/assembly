package loong64

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/loong64"
)

// builderInstrs - the instruction methods of the arch Builder (the
// ones returning arch.Instr; the operand roles Imm12/UImm12/... drop
// out by signature).
func builderInstrs() map[string]bool {
	instr := reflect.TypeFor[arch.Instr]()
	out := map[string]bool{}

	builder := reflect.TypeFor[arch.Builder]()
	for m := range builder.Methods() {
		if m.Type.NumOut() >= 1 && m.Type.Out(0) == instr {
			out[m.Name] = true
		}
	}

	return out
}

func TestChainMethodsAreTwinsOrDirectives(t *testing.T) {
	// the no-alias rule, mechanically: every chain method (returning
	// *Program) is a DSL directive, a sanctioned macro (La - the
	// pcalau12i+addi.d address pair), or the twin of a Builder
	// instruction method - nothing else may chain. Full Builder parity
	// (a twin for every instruction method) is a separate pass.
	directives := map[string]bool{
		"WithPos": true,
		"Label":   true,
		"Entry":   true,
		"Ascii":   true,
		"Bytes":   true,
	}
	sanctioned := map[string]bool{
		"La": true,
	}

	instrs := builderInstrs()
	program := reflect.TypeFor[*Program]()
	for m := range program.Methods() {
		if m.Type.NumOut() == 1 && m.Type.Out(0) == program {
			require.True(
				t, directives[m.Name] || sanctioned[m.Name] || instrs[m.Name],
				"%s: not a directive, sanctioned macro, or Builder twin (no aliases)", m.Name,
			)
		}
	}
}
