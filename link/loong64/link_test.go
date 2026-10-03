package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/link"
)

func TestELF(t *testing.T) {
	blob, err := ELF([]link.Source{
		{File: "a.s", Src: ".text\nstart:\n  li.d $a0, 1\n  jirl $zero, $ra, 0\n"},
		{File: "b.s", Src: ".data\nd:\n  .word 7\n"},
	}, "", 0x10000)

	require.NoError(t, err)
	require.NotEmpty(t, blob)
	require.Equal(t, []byte{0x7f, 'E', 'L', 'F'}, blob[:4])
	require.Equal(t, uint16(0x102), le16(blob[0x12:0x14]))   // e_machine = EM_LOONGARCH
	require.Equal(t, uint32(0x43), le32(blob[0x30:0x34]))    // e_flags: double-float base ABI
	require.Equal(t, uint64(0x10000), le64(blob[0x18:0x20])) // e_entry: start at base
}

func TestELFErrors(t *testing.T) {
	_, err := ELF([]link.Source{
		{File: "a.s", Src: ".text\nfn:\n  jirl $zero, $ra, 0\n"},
	}, "", 0x10000)
	require.ErrorContains(t, err, "no entry symbol")
}

func le16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}

func le32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func le64(b []byte) uint64 {
	var v uint64
	for i, x := range b {
		v |= uint64(x) << (8 * i)
	}

	return v
}
