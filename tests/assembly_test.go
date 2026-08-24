package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/okneniz/parsec/bytes"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/tests/cmd/objdump"
	"github.com/okneniz/assembly/text"
)

const (
	testBinaryPath  = "examples/hello-world/hello-world"   // arm64 Mach-O
	riscvBinaryPath = "examples/hello-riscv/hello-riscv.o" // riscv64 ELF (rv64gc, RVC)
)

func TestDisasmVsObjdump(t *testing.T) {
	diffAgainstObjdump(t, testBinaryPath, 90.0)
}
func TestDisasmVsObjdumpRISCV(t *testing.T) {
	diffAgainstObjdump(t, riscvBinaryPath, 90.0)
}

// diffAgainstObjdump loads the code section from path, disassembles it
// (the architecture is selected automatically from metadata), and compares
// against objdump by address. The test fails if the match rate is below the
// threshold percentage.
func diffAgainstObjdump(t *testing.T, path string, threshold float64) {
	t.Helper()

	ff, err := file.Detect(path)
	require.NoError(t, err)
	sec, err := ff.CodeSection()
	require.NoError(t, err)

	// The architecture is selected from the binary's metadata; the
	// arch-specific Parse is called directly (there is no common instruction/
	// decoder type). The lines are assembled by disasm (the instruction gives
	// the ObjDump text, the package adds the address and code columns); the
	// code column style depends on the format (Mach-O - bytes, ELF - hex
	// word), same as objdump.
	style := text.StyleFor(ff.Name())
	opts := disasm.NewOptions(style)
	ours := map[uint64]string{}
	archName := ""
	switch ff.ArchKind() {
	case file.ArchARM64:
		archName = arm64.Name
		insts, err := arm64.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		require.NoError(t, err)
		for _, inst := range insts {
			a := inst.Addr()
			ours[a] = objdump.Normalize(disasm.Line(a, sec.Data[a-sec.Addr:], inst, opts))
		}
	case file.ArchRISCV64:
		archName = riscv.Name
		insts, err := riscv.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		require.NoError(t, err)
		for _, inst := range insts {
			ours[inst.Addr()] = objdump.Normalize(
				disasm.Line(inst.Addr(), sec.Data[inst.Addr()-sec.Addr:], inst, opts),
			)
		}
	default:
		require.Fail(t, "unsupported architecture", "kind %d", ff.ArchKind())
	}

	// objdump output. objdump decodes executable sections itself and JUMPS
	// between addresses (skipping data bytes inside .text), so the comparison
	// is by address, not by ordinal index. The output is parsed by the objdump
	// package (parsec grammar, shared with tests/cmd/assembly-diff).
	out, err := objdump.Run(context.Background(), objdump.Args(ff.Name(), ff.ArchKind(), path))
	if err != nil {
		t.Skipf("no objdump on PATH can disassemble %s: %v", path, err)
	}

	objdumpOut := objdump.ParseByAddr(string(out))
	require.NotEmpty(t, objdumpOut)

	// Compare by objdump's addresses (it is the reference). Addresses skipped
	// by objdump as data are not penalized. Comments are cut off: ';'
	// (symbol stub, adrp hints) and ' <symbol>' (RISC-V branch/jump target
	// annotations).
	matched, mismatched, missingInOurs := 0, 0, 0
	var samples [10]string
	for addr, objLine := range objdumpOut {
		ourLine, ok := ours[addr]
		if !ok {
			missingInOurs++
			continue
		}

		objCmp := objdump.StripComments(objLine)
		ourCmp := objdump.StripComments(ourLine)
		if ourCmp == objCmp {
			matched++
		} else {
			mismatched++
			if mismatched <= 10 {
				samples[mismatched-1] = fmt.Sprintf("  0x%x\n    ours:   %s\n    objdump: %s",
					addr, ourLine, objLine)
			}
		}
	}

	total := len(objdumpOut)
	pct := 0.0
	if total > 0 {
		pct = float64(matched) * 100 / float64(total)
	}

	t.Logf("%s: objdump %d addrs | matched %d (%.2f%%), mismatched %d, not-in-ours %d",
		archName, total, matched, pct, mismatched, missingInOurs)

	require.GreaterOrEqual(t, pct, threshold)
	if mismatched > 0 {
		t.Logf("first %d mismatches:", mismatched)
		for _, s := range samples {
			if s != "" {
				t.Log(s)
			}
		}
	}
}
