package loong64

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// TestDiffAgainstLLVM - every decodeTable entry decodes its own base word
// (all operand fields zero) to the same text llvm-mc prints for that word
// alone at pc 0. Temporary integration gate; the property suite replaces
// it.
func TestDiffAgainstLLVM(t *testing.T) {
	mcPath, err := exec.LookPath("/opt/homebrew/opt/llvm/bin/llvm-mc")
	if err != nil {
		t.Skip("llvm-mc not available")
	}

	for _, e := range decodeTable {
		w := loongEncodings[e.name][0]

		var in bytes.Buffer
		for _, b := range binary.LittleEndian.AppendUint32(nil, w) {
			fmt.Fprintf(&in, "0x%x\n", b)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, mcPath, "--disassemble", "--triple=loongarch64")
		cmd.Stdin = &in
		var out, errb bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errb
		require.NoError(t, cmd.Run(), "%s: %s", e.name, errb.String())

		want := strings.Join(strings.Fields(strings.TrimSpace(out.String())), " ")
		got := strings.Join(
			strings.Fields(decodeOne(w, 0).ObjDump(disasm.DefaultViewCtx())),
			" ",
		)
		require.Equal(t, want, got, "%s (%#x)", e.name, w)
	}
}
