package file

// Tests of the universal Mach-O writer: golden pins of the legacy bytes
// (the generated tables must reproduce the ld template byte for byte),
// parser round-trips of full images (segments, sections, symbols, the
// exports trie, function starts, the entry), constructor validation, the
// exact-16K-page regression the template writer could not build, and
// native execution of images with a writable __DATA, a zero-filled __bss,
// and an entry past the old template limit.

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/macho"
)

// machoTestText - a deterministic code blob of n bytes.
func machoTestText(n int) []byte {
	text := make([]byte, n)
	for i := range text {
		text[i] = byte(i*7 + 3)
	}

	return text
}

// TestMachOGolden pins the exact bytes WriteMachO produced before the
// universal writer: same input, same image, forever. The hashes were
// captured from the template writer at 7673a0e.
func TestMachOGolden(t *testing.T) {
	rows := []struct {
		textLen int
		entry   uint64
		size    int
		sha     string
	}{
		{4, 0, 16837, "93aa323b129d9c94f3239c5469914c1031f47b9fe598929f25d6b410168040de"},
		{16, 0, 16837, "94e69d70e63ea648f2af5b3399c553f574969675da921dc0bb86d1d6d9e7e1ad"},
		{16, 4, 16837, "a5404bcf8cf3ba646ea22e91c94c37980ef6a3d7309c09d44d006999c40479bf"},
		{100, 12, 16837, "83b05726541146f4404f9624a7fd2fe39e311caca7046b00dc02e20f75d0d572"},
		{4096, 0, 16837, "7fc29c38138ff1d4b6d94e57506b8de9ea008d55fc157e7f49f61ea101259105"},
		{15687, 0, 16837, "e675bae474e5b89a838608e4aa99d33c98c6af22055afba806209bbd15cb09be"},
		{16300, 7000, 33349, "9cf1c343dd286a96564da75aeb530288b24ec49fa2e28884845c94fc229b00a5"},
		{16384, 0, 33349, "7cc2aa30f00e2c369af60fca1392bd658d66a8aa79c6ab191e60e9e9ae1099a7"},
	}

	for _, r := range rows {
		bin, err := WriteMachO(machoTestText(r.textLen), r.entry)
		require.NoError(t, err)
		require.Len(t, bin, r.size)

		sum := sha256.Sum256(bin)
		require.Equal(t, r.sha, hex.EncodeToString(sum[:]))
	}
}

// TestMachOExactPage - 696 + len landing exactly on a 16K page boundary
// panicked the template writer (its aligned branch divided the size by the
// page); the placement engine rounds up instead.
func TestMachOExactPage(t *testing.T) {
	for _, r := range []struct {
		textLen int
		vmsize  uint64
	}{
		{15688, 16384}, // 696 + 15688 == 16384 exactly
		{15689, 32768},
	} {
		bin, err := WriteMachO(machoTestText(r.textLen), 0)
		require.NoError(t, err)

		path := filepath.Join(t.TempDir(), "exact")
		require.NoError(t, os.WriteFile(path, bin, 0o755))

		f, err := macho.Open(path)
		require.NoError(t, err)

		var seg *macho.Segment
		for _, s := range f.Segments() {
			if s.SegName == "__TEXT" {
				seg = s
			}
		}
		require.NotNil(t, seg)
		require.Equal(t, r.vmsize, seg.Vmsize)
		require.Equal(t, r.vmsize, seg.Filesize)
	}
}

// TestMachOImageRoundTrip - a full image (globals and a local in __TEXT,
// data, a bss, a custom segment, more __DATA sections declared after the
// custom one, and a data span engineered so the zero-fill tail reaches
// past the file-backed pages) parsed back with the package's own parser.
func TestMachOImageRoundTrip(t *testing.T) {
	code := machoTestText(36)
	data := []byte{7, 0, 0, 0, 1, 2, 3, 4}
	big := machoTestText(16376) // 8 + 16376 == 16384: file span exactly one page

	sections := []MachOSection{
		{Segment: "__TEXT", Name: "__text", Data: code, Align: 4},
		{Segment: "__DATA", Name: "__data", Data: data, Align: 8},
		{Segment: "__DATA", Name: "__bss", Nobits: 16, Align: 8},
		{Segment: "__CUSTOM", Name: "__mine", Data: []byte{9, 9}, Align: 2},
		{Segment: "__DATA", Name: "__big", Data: big, Align: 8},
		{Segment: "__DATA", Name: "__bigbss", Nobits: 32, Align: 8},
	}

	syms := []MachOSym{
		{Name: "_start", Section: "__text", Global: true},
		{Name: "_fn2", Section: "__text", Off: 8, Global: true},
		{Name: "_localfn", Section: "__text", Off: 12},
		{Name: "_counter", Section: "__data", Off: 4, Global: true},
		{Name: "_arr", Section: "__bss", Global: true},
		{Name: "_mine", Section: "__mine", Global: true},
	}

	img, err := NewMachOImage(sections, syms, "_start")
	require.NoError(t, err)
	bin := img.Bytes()

	path := filepath.Join(t.TempDir(), "full")
	require.NoError(t, os.WriteFile(path, bin, 0o755))

	f, err := macho.Open(path)
	require.NoError(t, err)

	// segments: PAGEZERO, __TEXT, __DATA, __CUSTOM, __LINKEDIT
	segs := f.Segments()
	require.Len(t, segs, 5)
	require.Equal(t, "__PAGEZERO", segs[0].SegName)
	require.Equal(t, "__TEXT", segs[1].SegName)
	require.Equal(t, uint64(machoVMAddr), segs[1].Vmaddr)
	require.Equal(t, uint32(7), segs[1].Maxprot)
	require.Equal(t, uint32(5), segs[1].Initprot)
	require.Equal(t, "__DATA", segs[2].SegName)
	require.Equal(t, uint64(machoVMAddr)+segs[1].Vmsize, segs[2].Vmaddr)
	require.Equal(t, uint32(3), segs[2].Initprot)
	require.Equal(t, "__CUSTOM", segs[3].SegName)
	require.Equal(t, segs[2].Vmaddr+segs[2].Vmsize, segs[3].Vmaddr)
	require.Equal(t, "__LINKEDIT", segs[4].SegName)

	// every segment starts on a whole 16K page, file and memory
	for _, s := range segs[1:] {
		require.Zero(t, s.Fileoff%machoKPage)
		require.Zero(t, s.Vmaddr%machoKPage)
	}

	// the zero-fill tail: __DATA file span ends exactly at 16384, the bss
	// reserves push vmsize one page past filesize
	require.Equal(t, uint64(16384), segs[2].Filesize)
	require.Equal(t, uint64(32768), segs[2].Vmsize)

	// sections
	sect := f.Section("__data")
	require.NotNil(t, sect)
	require.Equal(t, "__DATA", sect.SegName)
	require.Equal(t, segs[2].Vmaddr, sect.Addr)
	require.Equal(t, uint64(len(data)), sect.Size)
	require.Equal(t, data, mustSectData(t, sect))

	bss := f.Section("__bss")
	require.NotNil(t, bss)
	require.Equal(t, uint32(1), bss.Flags&0xff) // S_ZEROFILL
	require.Empty(t, mustSectData(t, bss))

	text := f.Section("__text")
	require.NotNil(t, text)
	require.Equal(t, uint64(machoVMAddr)+uint64(img.place.codeOff), text.Addr)
	require.Equal(t, uint32(0x80000400), text.Flags)

	// the symbol table: __mh_execute_header plus the six input symbols
	symtab := mustSymbols(t, f)
	require.Len(t, symtab, 7)
	require.Equal(t, "__mh_execute_header", symtab[0].Name)
	require.True(t, symtab[0].IsExternal())
	require.Equal(t, uint64(machoVMAddr), symtab[0].Value)

	for i, s := range syms {
		got := symtab[i+1]
		require.Equal(t, s.Name, got.Name)
		require.Equal(t, s.Global, got.IsExternal())
		require.Equal(t, img.addrs[i], got.Value)
	}

	// the exports trie: __mh_execute_header plus the globals
	exports := mustExports(t, f)
	require.Len(t, exports, 6)
	at := map[string]uint64{}
	for _, e := range exports {
		at[e.Name] = e.Addr
	}
	require.Contains(t, at, "__mh_execute_header")
	require.Equal(t, uint64(machoVMAddr), at["__mh_execute_header"])
	require.Equal(t, img.addrs[0], at["_start"])
	require.Equal(t, img.addrs[3], at["_counter"])
	require.Equal(t, img.addrs[4], at["_arr"])
	require.Equal(t, img.addrs[5], at["_mine"])
	require.NotContains(t, at, "_localfn")

	// function starts: the __TEXT symbols only, sorted by address
	starts := mustFstarts(t, f)
	require.Equal(t, []uint64{img.addrs[0], img.addrs[1], img.addrs[2]}, starts)

	// the entry
	entry, ok := f.Entry()
	require.True(t, ok)
	require.Equal(t, img.addrs[0], entry)
}

func mustSectData(t *testing.T, s *macho.Section) []byte {
	t.Helper()

	d, err := s.Data()
	require.NoError(t, err)

	return d
}

func mustSymbols(t *testing.T, f *macho.File) []macho.Symbol {
	t.Helper()

	syms, err := f.Symbols()
	require.NoError(t, err)

	return syms
}

func mustExports(t *testing.T, f *macho.File) []macho.Export {
	t.Helper()

	e, err := f.Exports()
	require.NoError(t, err)

	return e
}

func mustFstarts(t *testing.T, f *macho.File) []uint64 {
	t.Helper()

	s, err := f.FunctionStarts()
	require.NoError(t, err)

	return s
}

// TestMachOImageErrors - the constructor rejects every malformed input.
func TestMachOImageErrors(t *testing.T) {
	text := machoTestText(16)

	rows := []struct {
		name     string
		sections []MachOSection
		syms     []MachOSym
		entry    string
		want     string
	}{
		{
			name: "no sections",
			want: "no sections",
		},
		{
			name:     "no text segment",
			sections: []MachOSection{{Segment: "__DATA", Name: "__data", Data: text}},
			entry:    "_start",
			want:     "no __TEXT",
		},
		{
			name:     "empty segment",
			sections: []MachOSection{{Segment: "", Name: "__text", Data: text}},
			want:     "empty segment",
		},
		{
			name:     "empty section name",
			sections: []MachOSection{{Segment: "__TEXT", Name: "", Data: text}},
			want:     "empty name",
		},
		{
			name:     "reserved segment",
			sections: []MachOSection{{Segment: "__LINKEDIT", Name: "__x", Data: text}},
			want:     "synthesized by the writer",
		},
		{
			name:     "data and nobits",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text, Nobits: 4}},
			want:     "both data and a nobits",
		},
		{
			name:     "alignment not a power of two",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text, Align: 3}},
			want:     "not a power of two",
		},
		{
			name: "duplicate section",
			sections: []MachOSection{
				{Segment: "__TEXT", Name: "__text", Data: text},
				{Segment: "__DATA", Name: "__text", Data: text},
			},
			want: "duplicate section",
		},
		{
			name: "text after another segment",
			sections: []MachOSection{
				{Segment: "__DATA", Name: "__d", Data: text},
				{Segment: "__TEXT", Name: "__text", Data: text},
			},
			want: "before the other segments",
		},
		{
			name:     "no entry name",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			want:     "no entry symbol",
		},
		{
			name:     "entry not defined",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			entry:    "_start",
			want:     "not defined",
		},
		{
			name:     "empty symbol name",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			syms:     []MachOSym{{Name: "", Section: "__text"}},
			entry:    "_start",
			want:     "empty name",
		},
		{
			name:     "reserved symbol",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			syms:     []MachOSym{{Name: "__mh_execute_header", Section: "__text"}},
			entry:    "_start",
			want:     "synthesized by the writer",
		},
		{
			name:     "symbol section missing",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			syms:     []MachOSym{{Name: "_start", Section: "__nope"}},
			entry:    "_start",
			want:     "no section",
		},
		{
			name:     "symbol offset outside",
			sections: []MachOSection{{Segment: "__TEXT", Name: "__text", Data: text}},
			syms:     []MachOSym{{Name: "_start", Section: "__text", Off: 17}},
			entry:    "_start",
			want:     "outside section",
		},
		{
			name: "entry outside text",
			sections: []MachOSection{
				{Segment: "__TEXT", Name: "__text", Data: text},
				{Segment: "__DATA", Name: "__data", Data: text},
			},
			syms:  []MachOSym{{Name: "_start", Section: "__data"}},
			entry: "_start",
			want:  "must sit in a __TEXT",
		},
	}

	for _, r := range rows {
		t.Run(r.name, func(t *testing.T) {
			_, err := NewMachOImage(r.sections, r.syms, r.entry)
			require.ErrorContains(t, err, r.want)
		})
	}
}

// The Darwin arm64 exit syscall: movz x16, #0x200, lsl 16; movk x16, #1
// (0x2000001 = SYSCALL_CLASS_UNIX|exit), then svc #0x80.
var machoExitSetup = []uint32{0xD2A04010, 0xF2800030}

// TestMachOImageExec - whole images running natively on an arm64 Mac:
// a writable __DATA read-modify-write, a zero-filled __bss, and an entry
// past the 8192-byte template limit of the old writer. The adrp/add
// immediates are patched from the image placement itself - the same
// numbers the assembler policies hand out.
func TestMachOImageExec(t *testing.T) {
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("native arm64 macOS only")
	}

	// rmw builds a program that loads from the named section, increments,
	// stores back, reloads, and exits with the value - proof the section
	// is mapped, addressable via adrp+add, writable, and re-readable.
	rmw := func(t *testing.T, sect MachOSection, inc uint32, want int) []byte {
		code := machoWords([]uint32{
			0, // adrp x0, <sect>     - patched
			0, // add x0, x0, #lo12   - patched
			0xB9400001,          // ldr w1, [x0]
			0x11000000 | inc<<10 | 1<<5 | 1, // add w1, w1, #inc
			0xB9000001,                      // str w1, [x0]
			0xB9400002,                      // ldr w2, [x0]
		})
		code = append(code, machoWords(machoExitSetup)...)
		code = append(code, machoWords([]uint32{
			0xAA0203E0, // mov x0, x2
			0xD4001001, // svc #0x80
		})...)

		sections := []MachOSection{
			{Segment: "__TEXT", Name: "__text", Data: code, Align: 4},
			sect,
		}
		syms := []MachOSym{{Name: "_start", Section: "__text", Global: true}}

		img, err := NewMachOImage(sections, syms, "_start")
		require.NoError(t, err)

		addr := img.place.sectOf(1).addr
		binary.LittleEndian.PutUint32(code[0:4], machoAdrp(0, addr, uint64(machoVMAddr)))
		binary.LittleEndian.PutUint32(code[4:8], machoAddLo(0, addr))

		img, err = NewMachOImage(sections, syms, "_start")
		require.NoError(t, err)

		_ = want
		return img.Bytes()
	}

	rows := []struct {
		name     string
		exitCode int
		build    func(t *testing.T) []byte
	}{
		{"data read-modify-write", 8, func(t *testing.T) []byte {
			return rmw(t, MachOSection{Segment: "__DATA", Name: "__data",
				Data: []byte{7, 0, 0, 0}, Align: 8}, 1, 8)
		}},
		{"bss zero on first touch", 5, func(t *testing.T) []byte {
			return rmw(t, MachOSection{Segment: "__DATA", Name: "__bss",
				Nobits: 16, Align: 8}, 5, 5)
		}},
		{"entry past the template limit", 42, func(t *testing.T) []byte {
			code := make([]uint32, 2100) // 8400 bytes of nop
			for i := range code {
				code[i] = 0xD503201F
			}

			code = append(code, machoExitSetup...)
			code = append(code, 0xD2800540, 0xD4001001) // mov x0, #42; svc

			sections := []MachOSection{
				{Segment: "__TEXT", Name: "__text", Data: machoWords(code), Align: 4},
			}
			syms := []MachOSym{{
				Name:    "_start",
				Section: "__text",
				Off:     4 * 2100, // 8400 > 8192: the old writer refused
				Global:  true,
			}}

			img, err := NewMachOImage(sections, syms, "_start")
			require.NoError(t, err)

			return img.Bytes()
		}},
	}

	for _, r := range rows {
		t.Run(r.name, func(t *testing.T) {
			bin := r.build(t)

			path := filepath.Join(t.TempDir(), "prog")
			require.NoError(t, os.WriteFile(path, bin, 0o755))

			cmd := exec.CommandContext(context.Background(), path)
			runErr := cmd.Run()

			var exitErr *exec.ExitError
			require.ErrorAs(t, runErr, &exitErr)
			require.Equal(t, r.exitCode, exitErr.ExitCode())
		})
	}
}

// machoWords - u32 instructions as little-endian bytes.
func machoWords(code []uint32) []byte {
	out := make([]byte, 4*len(code))
	for i, w := range code {
		binary.LittleEndian.PutUint32(out[4*i:], w)
	}

	return out
}

// machoAdrp - the adrp encoding for rd: pc to target as a page delta.
func machoAdrp(rd uint32, target, pc uint64) uint32 {
	imm := uint32(((target &^ 0xFFF) - (pc &^ 0xFFF)) >> 12)
	return 0x90000000 | (imm&3)<<29 | (imm>>2)<<5 | rd
}

// machoAddLo - add rd, rd, #target&0xfff.
func machoAddLo(rd uint32, target uint64) uint32 {
	return 0x91000000 | uint32(target&0xFFF)<<10 | rd<<5 | rd
}
