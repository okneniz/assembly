package riscv

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJalrRegCtor(t *testing.T) {
	// the c.jalr halfwords (llvm-mc parity: jalr rs1 compresses)
	for _, c := range []struct {
		name string
		rs1  int
		half uint16
	}{
		{"jalr t0", 5, 0x9282},
		{"jalr a0", 10, 0x9502},
		{"jalr ra", 1, 0x9082},
	} {
		t.Run(c.name, func(t *testing.T) {
			b := ctorBytes(t, New().JalrReg(xreg(t, c.rs1)))
			require.Len(t, b, 2, "c.jalr")
			require.Equal(t, c.half, binary.LittleEndian.Uint16(b))
		})
	}
}
