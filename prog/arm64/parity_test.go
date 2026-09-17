package arm64

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// builderInstrs - the instruction methods of arch.Builder (the ones
// returning arch.Instr; the operand roles Imm12/Imm16/Imm6 drop out by
// signature) mapped to their operand count - the set the chain twins
// mirror.
func builderInstrs() map[string]int {
	instr := reflect.TypeFor[arch.Instr]()
	out := map[string]int{}

	builder := reflect.TypeFor[arch.Builder]()
	for m := range builder.Methods() {
		if m.Type.NumOut() >= 1 && m.Type.Out(0) == instr {
			out[m.Name] = m.Type.NumIn() - 1
		}
	}

	return out
}

func TestBuilderParity(t *testing.T) {
	// every instruction method of the Builder has a same-named chain
	// twin on *Program with the same operand count
	program := reflect.TypeFor[*Program]()
	twins := 0
	for name, params := range builderInstrs() {
		twin, ok := program.MethodByName(name)
		require.True(t, ok, "no prog twin for Builder.%s", name)

		require.Equal(t, 1, twin.Type.NumOut(), "%s: the twin must chain", name)
		require.Equal(t, program, twin.Type.Out(0), "%s: the twin must return *Program", name)
		require.Equal(
			t, params, twin.Type.NumIn()-1,
			"%s: the twin must take the Builder's operand count", name,
		)
		twins++
	}

	// a new Builder instruction method without a twin fails the loop
	// above; a removed one fails this pin
	require.Equal(t, 97, twins)
}

func TestChainMethodsAreTwinsOrDirectives(t *testing.T) {
	// the no-alias rule, mechanically: every chain method (returning
	// *Program) is either a DSL directive or the twin of a Builder
	// instruction method - nothing else may chain
	directives := map[string]bool{
		"WithPos": true,
		"Label":   true,
		"Entry":   true,
		"Ascii":   true,
		"Bytes":   true,
		"Text":    true,
		"Data":    true,
		"Half":    true,
		"Word":    true,
		"Quad":    true,
		"Bss":     true,
		"La":      true, // the adrp+add pair, the riscv/loong64 twin
	}

	instrs := builderInstrs()
	program := reflect.TypeFor[*Program]()
	for m := range program.Methods() {
		if m.Type.NumOut() == 1 && m.Type.Out(0) == program {
			_, isInstr := instrs[m.Name]
			require.True(
				t, directives[m.Name] || isInstr,
				"%s: neither a directive nor a Builder twin (no aliases in the DSL)", m.Name,
			)
		}
	}
}
