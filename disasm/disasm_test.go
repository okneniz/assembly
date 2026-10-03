package disasm

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/text"
)

// fakeInstr is a stand-in instruction for checking the harness without a
// binding to an architecture: the package needs only the ObjDump text.
type fakeInstr struct {
	s string
}

func (f fakeInstr) ObjDump(_ ViewCtx) string {
	return f.s
}

func TestLine(t *testing.T) {
	cases := []struct {
		addr uint64
		code []byte
		in   fakeInstr
		opts Options
		want string
	}{
		// Mach-O/ARM64: byte column, 4 bytes.
		{
			0x1000,
			[]byte{0xfd, 0x23, 0x00, 0xd1},
			newFakeInstr("add x27, x28"),
			NewOptions(text.CodeBytes),
			"1000:\tfd 23 00 d1\tadd x27, x28",
		},
		// ELF/RISC-V: hex word, 2 bytes (compressed).
		{
			0x1000,
			[]byte{0x02, 0x88},
			newFakeInstr("li s0, 0x2"),
			NewOptions(text.CodeWord),
			"1000:\t8802\tli s0, 0x2",
		},
		// An instruction without operands.
		{
			0x1004,
			[]byte{0xc0, 0x03, 0x5f, 0xd6},
			newFakeInstr("ret"),
			NewOptions(text.CodeBytes),
			"1004:\tc0 03 5f d6\tret",
		},
	}
	for _, c := range cases {
		require.Equal(t, c.want, Line(c.addr, c.code, c.in, c.opts), "Line(%#x)", c.addr)
	}
}

func TestWrite(t *testing.T) {
	// A 2-byte compressed + a 4-byte full instruction: the addresses go
	// 0x80000000, 0x80000002; the code as a hex word (ELF format). size is
	// the riscv stream rule (the low two bits of the first byte - the
	// sample bytes follow it).
	code := []byte{0x02, 0x88, 0xb3, 0x03, 0xa0, 0x00}
	instrs := []fakeInstr{newFakeInstr("li s0, 0x2"), newFakeInstr("add x27, x28")}
	size := func(b []byte) int {
		if b[0]&0x3 != 0x3 {
			return 2
		}

		return 4
	}

	var buf bytes.Buffer
	err := Write(&buf, 0x80000000, code, instrs, NewOptions(text.CodeWord), size)
	require.NoError(t, err)
	want := "80000000:\t8802\tli s0, 0x2\n80000002:\t00a003b3\tadd x27, x28\n"
	require.Equal(t, want, buf.String(), "Write")
}

func newFakeInstr(s string) fakeInstr {
	return fakeInstr{
		s: s,
	}
}
