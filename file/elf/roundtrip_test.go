// Round-trip tests: serializing the parsed fixed-layout structures back into
// bytes must match the original file byte for byte (Ehdr, Shdr, Sym - at
// their offsets). This strictly verifies that the parser has read every field.
package elf_test

import (
	stdelf "debug/elf"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/elf"
)

// TestRoundtripHeader serializes Header() into an Ehdr and compares it with
// the file (the counters are the resolved ones; for files without extended
// numbering they equal the raw ones, which is what the corpus covers).
func TestRoundtripHeader(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			f, err := elf.Open(path)
			require.NoError(t, err, "Open")
			h := f.Header()
			le := h.Order
			is64 := h.Class == elf.CLASS64
			size := 52
			if is64 {
				size = 64
			}

			b := make([]byte, size)
			copy(b, h.Ident[:])
			if is64 {
				le.PutUint16(b[16:], uint16(h.Type))
				le.PutUint16(b[18:], uint16(h.Machine))
				le.PutUint32(b[20:], h.Version)
				le.PutUint64(b[24:], h.Entry)
				le.PutUint64(b[32:], h.Phoff)
				le.PutUint64(b[40:], h.Shoff)
				le.PutUint32(b[48:], h.Flags)
				le.PutUint16(b[52:], h.Ehsize)
				le.PutUint16(b[54:], h.Phentsize)
				le.PutUint16(b[56:], uint16(h.Phnum))
				le.PutUint16(b[58:], h.Shentsize)
				le.PutUint16(b[60:], uint16(h.Shnum))
				le.PutUint16(b[62:], uint16(h.Shstrndx))
			} else {
				le.PutUint16(b[16:], uint16(h.Type))
				le.PutUint16(b[18:], uint16(h.Machine))
				le.PutUint32(b[20:], h.Version)
				le.PutUint32(b[24:], uint32(h.Entry))
				le.PutUint32(b[28:], uint32(h.Phoff))
				le.PutUint32(b[32:], uint32(h.Shoff))
				le.PutUint32(b[36:], h.Flags)
				le.PutUint16(b[40:], h.Ehsize)
				le.PutUint16(b[42:], h.Phentsize)
				le.PutUint16(b[44:], uint16(h.Phnum))
				le.PutUint16(b[46:], h.Shentsize)
				le.PutUint16(b[48:], uint16(h.Shnum))
				le.PutUint16(b[50:], uint16(h.Shstrndx))
			}

			require.Equal(t, raw[:size], b, "Ehdr")
		})
	}
}

// TestRoundtripShdrs serializes the section headers and compares them with the file.
func TestRoundtripShdrs(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			f, err := elf.Open(path)
			require.NoError(t, err, "Open")
			std, err := stdelf.Open(path)
			require.NoError(t, err, "stdlib")
			defer func() {
				require.NoError(t, std.Close(), "stdlib Close")
			}()
			_ = std

			h := f.Header()
			le := h.Order
			is64 := h.Class == elf.CLASS64
			shEnt := 40
			if is64 {
				shEnt = 64
			}

			require.True(
				t,
				int(h.Shoff)+h.Shnum*shEnt <= len(raw),
				"the section table is outside the file - the parser must have failed",
			)
			for i, s := range f.Sections() {
				off := int(h.Shoff) + i*shEnt
				var b []byte
				if is64 {
					b = make([]byte, 64)
					// The name field is taken from the file as is: linkers
					// reuse suffixes (".text" inside ".rel.text"), so reverse
					// offset lookup is ambiguous. The name as such is already
					// checked in TestDiffSections.
					copy(b[0:4], raw[off:off+4])
					le.PutUint32(b[4:], uint32(s.Type))
					le.PutUint64(b[8:], uint64(s.Flags))
					le.PutUint64(b[16:], s.Addr)
					le.PutUint64(b[24:], s.Off)
					le.PutUint64(b[32:], s.Size)
					le.PutUint32(b[40:], s.Link)
					le.PutUint32(b[44:], s.Info)
					le.PutUint64(b[48:], s.Addralign)
					le.PutUint64(b[56:], s.Entsize)
				} else {
					b = make([]byte, 40)
					copy(b[0:4], raw[off:off+4])
					le.PutUint32(b[4:], uint32(s.Type))
					le.PutUint32(b[8:], uint32(s.Flags))
					le.PutUint32(b[12:], uint32(s.Addr))
					le.PutUint32(b[16:], uint32(s.Off))
					le.PutUint32(b[20:], uint32(s.Size))
					le.PutUint32(b[24:], s.Link)
					le.PutUint32(b[28:], s.Info)
					le.PutUint32(b[32:], uint32(s.Addralign))
					le.PutUint32(b[36:], uint32(s.Entsize))
				}

				require.Equal(t, raw[off:off+shEnt], b, "shdr[%d] %s", i, s.Name)
			}
		})
	}
}

// TestRoundtripSymbols serializes .symtab and compares it with the file.
func TestRoundtripSymbols(t *testing.T) {
	for _, path := range corpus(t) {
		f, err := elf.Open(path)
		require.NoError(t, err, "Open")
		syms, err := f.Symbols()
		if err != nil {
			continue // no table
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)
			var symsec *elf.Section
			for _, s := range f.Sections() {
				if s.Type == elf.SHT_SYMTAB {
					symsec = s
					break
				}
			}

			require.NotNil(t, symsec, "no .symtab")
			if symsec == nil {
				return // require has already failed the test; the exit is needed for the analysis
			}

			h := f.Header()
			le := h.Order
			is64 := h.Class == elf.CLASS64
			symEnt := 16
			if is64 {
				symEnt = 24
			}

			// Our table lacks the zero symbol: compare starting at offset+entsize.
			base := int(symsec.Off) + int(symsec.Entsize)
			for i, s := range syms {
				off := base + i*int(symsec.Entsize)
				var b []byte
				if is64 {
					b = make([]byte, 24)
					copy(b[0:4], raw[off:off+4]) // name - strtab suffix sharing
					b[4] = s.Info
					b[5] = s.Other
					le.PutUint16(b[6:], uint16(s.Shndx))
					le.PutUint64(b[8:], s.Value)
					le.PutUint64(b[16:], s.Size)
				} else {
					b = make([]byte, 16)
					copy(b[0:4], raw[off:off+4])
					le.PutUint32(b[4:], uint32(s.Value))
					le.PutUint32(b[8:], uint32(s.Size))
					b[12] = s.Info
					b[13] = s.Other
					le.PutUint16(b[14:], uint16(s.Shndx))
				}

				require.Equal(t, raw[off:off+symEnt], b, "sym[%d] %s", i, s.Name)
			}
		})
	}
}

// TestMutationSweep mutates random bytes of the header/tables: Open and the
// accessors must return errors, but never panic.
func TestMutationSweep(t *testing.T) {
	seed := int64(42)
	if v := os.Getenv("ASSEMBLY_SEED"); v != "" {
		_, err := fmt.Sscan(v, &seed)
		require.NoError(t, err, "ASSEMBLY_SEED")
	}

	rnd := rand.New(rand.NewSource(seed))

	for _, name := range []string{"hello-riscv64", "hello-386", "reloc-aarch64.o"} {
		raw, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			continue
		}

		for iter := range 200 {
			mut := append([]byte(nil), raw...)
			// Mutate the "structural" zone: the first 4K of the file (headers+tables).
			for range 8 {
				mut[rnd.Intn(min(4096, len(mut)))] = byte(rnd.Intn(256))
			}

			path := t.TempDir() + "/mut.elf"
			err := os.WriteFile(path, mut, 0o644)
			require.NoError(t, err)
			func() {
				defer func() {
					r := recover()
					require.Nil(t, r, "%s iter %d: panic: %v", name, iter, r)
				}()
				f, err := elf.Open(path)
				if err != nil {
					return // an error is a normal outcome
				}

				// the calls must not panic; errors are allowed (corrupted data)
				call := func(_ any, _ error) {}
				call(f.Symbols())
				call(f.DynamicSymbols())
				call(f.Dynamic())
				call(f.Notes())
				for _, s := range f.Sections() {
					call(s.Data())
					call(f.Relocations(s))
				}
			}()
		}
	}
}

// TestExtendedNumbering: e_shnum=0 with a non-empty table means the real
// section count lives in shdr[0].sh_size (extended numbering from gABI).
func TestExtendedNumbering(t *testing.T) {
	le := binary.LittleEndian
	img := make([]byte, 512)

	// Ehdr: 3 sections (NULL + .text + .shstrtab), e_shnum = 0.
	img[0], img[1], img[2], img[3] = 0x7f, 'E', 'L', 'F'
	img[4] = 2                  // CLASS64
	img[5] = 1                  // LE
	le.PutUint16(img[16:], 2)   // ET_EXEC
	le.PutUint16(img[18:], 183) // EM_AARCH64
	le.PutUint32(img[20:], 1)
	le.PutUint64(img[40:], 256) // e_shoff
	le.PutUint16(img[58:], 64)  // shentsize
	le.PutUint16(img[60:], 0)   // e_shnum = 0 -> extended
	le.PutUint16(img[62:], 2)   // shstrndx

	// shstrtab at offset 128: "\0.text\0.shstrtab\0".
	copy(img[128:], "\x00.text\x00.shstrtab\x00")

	// Shdrs from 256: [0] with sh_size=3 (extended count), [1] .text, [2] shstrtab.
	le.PutUint64(img[256+32:], 3) // shdr[0].sh_size = the real shnum
	le.PutUint32(img[320+0:], 1)  // shdr[1].sh_name -> ".text"
	le.PutUint32(img[320+4:], 1)  // PROGBITS
	le.PutUint64(img[320+24:], 64)
	le.PutUint64(img[320+32:], 4)
	le.PutUint32(img[384+0:], 7) // shdr[2].sh_name -> ".shstrtab"
	le.PutUint32(img[384+4:], 3) // STRTAB
	le.PutUint64(img[384+24:], 128)
	le.PutUint64(img[384+32:], 16)

	path := t.TempDir() + "/extnum.elf"
	err := os.WriteFile(path, img, 0o644)
	require.NoError(t, err)
	f, err := elf.Open(path)
	require.NoError(t, err, "Open")
	h := f.Header()
	require.Equal(t, 3, h.Shnum, "extended numbering was not resolved")
	secs := f.Sections()
	require.Len(t, secs, 3)
	require.Equal(t, ".text", secs[1].Name)
	require.Equal(t, ".shstrtab", secs[2].Name)
	require.NotNil(t, f.Section(".text"), "Section(%q) not found", ".text")
}
