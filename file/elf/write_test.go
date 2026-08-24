package elf_test

// Tests for the minimal ELF writer: the NOBITS tail (p_memsz > p_filesz) and
// the "tail only on the last blob" contract. The phdr fields are read directly
// from the output bytes (e_phoff @0x20, p_filesz @+32, p_memsz @+40).

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/elf"
)

func phdrOf(t *testing.T, out []byte) (filesz, memsz uint64) {
	t.Helper()
	phOff := binary.LittleEndian.Uint64(out[0x20:])
	return binary.LittleEndian.Uint64(out[phOff+32:]), binary.LittleEndian.Uint64(out[phOff+40:])
}

func TestWriteMemszTail(t *testing.T) {
	base := uint64(0x1000)
	// file part of 4 bytes, memory 0x100: a tail of 0xFC zeros with no file space
	blob := elf.NewNobitsBlob(base, []byte{1, 2, 3, 4}, 0x100)
	out, err := elf.Write(elf.EM_RISCV, base, base, []elf.Blob{blob})
	require.NoError(t, err, "Write")

	filesz, memsz := phdrOf(t, out)
	// file: image (4) + phdr (56); memory: + a 0x100-4 tail
	require.Equal(t, uint64(4+56), filesz, "p_filesz")
	require.Equal(t, uint64(4+56+0x100-4), memsz, "p_memsz")
	require.Equal(t, int(0x1000+4+56), len(out), "the tail occupies no file space")
}

func TestWriteMemszTailNotLast(t *testing.T) {
	b1 := elf.NewNobitsBlob(0x1000, []byte{1}, 0x10)
	b2 := elf.NewBlob(0x2000, []byte{2})
	_, err := elf.Write(elf.EM_RISCV, 0x1000, 0x1000, []elf.Blob{b1, b2})
	require.ErrorContains(t, err, "last blob", "a tail not on the last blob must be an error")
}

func TestWriteMemszTooSmall(t *testing.T) {
	b1 := elf.NewNobitsBlob(0x1000, []byte{1, 2, 3}, 2) // memory smaller than the file
	_, err := elf.Write(elf.EM_RISCV, 0x1000, 0x1000, []elf.Blob{b1})
	require.ErrorContains(t, err, "MemSize", "memory smaller than the file must be an error")
}
