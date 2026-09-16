package riscv

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/riscv"
)

// builderInstrs - the instruction methods of the arch Builder (the ones
// returning arch.Instr; the operand roles Imm12/Imm20/Off drop out by
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

// untwinable - Builder instruction methods with no chain twin, each for
// a structural reason: both are arch-side pseudos that predate the
// no-alias rule.
// Mv - the c.mv HALFWORD: a 2-byte encoder cannot ride the fixed
// 4-byte line layout (and mv is a pseudo of addi anyway);
// JalrReg - the "jalr rs" pseudo (jalr ra, 0(rs)): the real form is
// Jalr(Ra, rs, 0).
func untwinable() map[string]bool {
	return map[string]bool{
		"Mv":      true,
		"JalrReg": true,
	}
}

func TestBuilderParity(t *testing.T) {
	// every twin-able instruction method of the Builder has a
	// same-named chain twin on *Program with the same operand count
	program := reflect.TypeFor[*Program]()
	exempt := untwinable()
	twins := 0
	for name, params := range builderInstrs() {
		if exempt[name] {
			continue
		}

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

	// 109 instruction methods, 2 untwinable: 107 twins, no more, no less
	require.Equal(t, 107, twins)
}

func TestChainMethodsAreTwinsOrDirectives(t *testing.T) {
	// the no-alias rule, mechanically: every chain method (returning
	// *Program) is a DSL directive, a sanctioned macro (La - the
	// auipc+addi address pair; Ecall - a real instruction whose arch
	// side is not a Builder method), or the twin of a Builder
	// instruction method - nothing else may chain.
	directives := map[string]bool{
		"WithPos": true,
		"Label":   true,
		"Entry":   true,
		"Ascii":   true,
		"Bytes":   true,
	}
	sanctioned := map[string]bool{
		"La":    true,
		"Ecall": true,
	}

	instrs := builderInstrs()
	program := reflect.TypeFor[*Program]()
	for m := range program.Methods() {
		if m.Type.NumOut() == 1 && m.Type.Out(0) == program {
			_, isInstr := instrs[m.Name]
			require.True(
				t, directives[m.Name] || sanctioned[m.Name] || isInstr,
				"%s: not a directive, sanctioned macro, or Builder twin (no aliases)", m.Name,
			)
		}
	}
}
