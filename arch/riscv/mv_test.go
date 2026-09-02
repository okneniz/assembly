package riscv

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMvCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		half  uint16
	}{
		{"mv a0, a1", New().Mv(xreg(t, 10), xreg(t, 11)), 0x852e},
		{"mv ra, a5", New().Mv(xreg(t, 1), xreg(t, 15)), 0x80be},
		{"mv sp, t6", New().Mv(xreg(t, 2), xreg(t, 31)), 0x817e},
	} {
		t.Run(c.name, func(t *testing.T) {
			b := ctorBytes(t, c.instr)
			require.Len(t, b, 2, "mv is the c.mv halfword")
			require.Equal(t, c.half, binary.LittleEndian.Uint16(b))
		})
	}
}
