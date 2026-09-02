// Diff tests against the stdlib debug/elf (tests only; the runtime does not
// depend on debug/*). Corpus: Go cross-builds (32/64-bit, LE/BE) + clang .o
// files with symbol and relocation tables, compiled into testdata.
package elf_test

import (
	"context"
	stdelf "debug/elf"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/elf"
)

// corpus returns the testdata files (skipping missing ones - e.g., with a
// partial checkout).
func corpus(t *testing.T) []string {
	t.Helper()
	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Skipf("no testdata: %v", err)
	}

	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, filepath.Join("testdata", e.Name()))
		}
	}

	if len(out) == 0 {
		t.Skip("testdata is empty")
	}

	return out
}

// TestDiffHeader compares headers.
func TestDiffHeader(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			std, err := stdelf.Open(path)
			require.NoError(t, err, "stdlib")
			defer func() {
				require.NoError(t, std.Close(), "stdlib Close")
			}()

			h, sh := ours.Header(), std.FileHeader
			require.Equal(t, uint16(sh.Type), uint16(h.Type), "e_type")
			require.Equal(t, uint16(sh.Machine), uint16(h.Machine), "e_machine")
			require.Equal(t, sh.Entry, h.Entry, "e_entry")
			require.Equal(t, byte(sh.Class), byte(h.Class), "EI_CLASS")
			require.Equal(t, sh.ByteOrder, h.Order, "EI_DATA")
			require.Equal(t, byte(sh.OSABI), byte(h.OSABI), "EI_OSABI")
			require.Equal(t, sh.ABIVersion, h.ABIVersion, "EI_ABIVERSION")
		})
	}
}

// TestDiffSections compares the section tables field by field.
func TestDiffSections(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			std, err := stdelf.Open(path)
			require.NoError(t, err, "stdlib")
			defer func() {
				require.NoError(t, std.Close(), "stdlib Close")
			}()

			a, b := ours.Sections(), std.Sections
			require.Len(t, a, len(b), "section count")
			for i := range a {
				x, y := a[i], b[i]
				require.Equal(t, y.Name, x.Name, "section %d: name", i)
				require.Equal(t, uint32(y.Type), uint32(x.Type), "section %s: type", x.Name)
				require.Equal(t, uint64(y.Flags), uint64(x.Flags), "section %s: flags", x.Name)
				require.Equal(t, y.Addr, x.Addr, "section %s: addr", x.Name)
				require.Equal(t, y.Offset, x.Off, "section %s: offset", x.Name)
				require.Equal(t, y.Size, x.Size, "section %s: size", x.Name)
				require.Equal(t, y.Link, x.Link, "section %s: link", x.Name)
				require.Equal(t, y.Info, x.Info, "section %s: info", x.Name)
				require.Equal(t, y.Addralign, x.Addralign, "section %s: align", x.Name)
				require.Equal(t, y.Entsize, x.Entsize, "section %s: entsize", x.Name)
			}

			// The .text data must match byte for byte.
			if ts := ours.Section(".text"); ts != nil {
				d1, err := ts.Data()
				require.NoError(t, err, ".text Data")
				d2, err := std.Section(".text").Data()
				require.NoError(t, err, "stdlib .text Data")
				require.Equal(t, d2, d1, ".text: data")
			}
		})
	}
}

// TestDiffProgs compares program headers.
func TestDiffProgs(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			std, err := stdelf.Open(path)
			require.NoError(t, err, "stdlib")
			defer func() {
				require.NoError(t, std.Close(), "stdlib Close")
			}()

			a, b := ours.ProgramHeaders(), std.Progs
			require.Len(t, a, len(b), "phdr count")
			for i := range a {
				x, y := a[i], b[i]
				require.Equal(t, uint32(y.Type), uint32(x.Type), "phdr %d: type", i)
				require.Equal(t, uint32(y.Flags), uint32(x.Flags), "phdr %d (%v): flags", i, x.Type)
				require.Equal(t, y.Off, x.Off, "phdr %d (%v): off", i, x.Type)
				require.Equal(t, y.Vaddr, x.Vaddr, "phdr %d (%v): vaddr", i, x.Type)
				require.Equal(t, y.Paddr, x.Paddr, "phdr %d (%v): paddr", i, x.Type)
				require.Equal(t, y.Filesz, x.Filesz, "phdr %d (%v): filesz", i, x.Type)
				require.Equal(t, y.Memsz, x.Memsz, "phdr %d (%v): memsz", i, x.Type)
				require.Equal(t, y.Align, x.Align, "phdr %d (%v): align", i, x.Type)
			}
		})
	}
}

// TestDiffSymbols compares the symbol tables (.symtab and .dynsym).
func TestDiffSymbols(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			std, err := stdelf.Open(path)
			require.NoError(t, err, "stdlib")
			defer func() {
				require.NoError(t, std.Close(), "stdlib Close")
			}()

			check := func(what string, a []elf.Symbol, b []stdelf.Symbol) {
				t.Helper()
				require.Len(t, a, len(b), "%s: symbol count", what)
				for i := range a {
					x, y := a[i], b[i]
					require.Equal(t, y.Name, x.Name, "%s[%d]: name", what, i)
					require.Equal(t, y.Value, x.Value, "%s[%d] (%s): value", what, i, x.Name)
					require.Equal(t, y.Size, x.Size, "%s[%d] (%s): size", what, i, x.Name)
					require.Equal(t, y.Info, x.Info, "%s[%d] (%s): info", what, i, x.Name)
					require.Equal(t, y.Other, x.Other, "%s[%d] (%s): other", what, i, x.Name)
					require.Equal(
						t,
						uint32(y.Section),
						x.Shndx,
						"%s[%d] (%s): shndx",
						what,
						i,
						x.Name,
					)
				}
			}

			a1, err := ours.Symbols()
			b1, serr := std.Symbols()
			if serr != nil {
				require.Error(t, err, "Symbols: stdlib failed (%v), ours did not", serr)
			} else {
				require.NoError(t, err, "Symbols")
			}

			check("symtab", a1, b1)

			a2, err := ours.DynamicSymbols()
			b2, derr := std.DynamicSymbols()
			if derr != nil {
				require.Error(t, err, "DynamicSymbols: stdlib failed (%v), ours did not", derr)
			} else {
				require.NoError(t, err, "DynamicSymbols")
			}

			check("dynsym", a2, b2)
		})
	}
}

// llvmReadelf - the path of llvm-readelf (the oracle for the diff tests
// below); the test is skipped when the tool is unavailable.
func llvmReadelf(t *testing.T) string {
	t.Helper()
	for _, c := range []string{
		"/opt/homebrew/opt/llvm/bin/llvm-readelf",
		"/opt/homebrew/bin/llvm-readelf",
		"/usr/local/opt/llvm/bin/llvm-readelf",
		"/usr/local/bin/llvm-readelf",
		"llvm-readelf",
	} {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}

	t.Skip("llvm-readelf is unavailable")
	return ""
}

// TestDiffRelocs compares the relocations with the output of llvm-readelf
// -rW: offset and the full r_info layout (symbol index + type).
func TestDiffRelocs(t *testing.T) {
	readelf := llvmReadelf(t)

	for _, path := range corpus(t) {
		out, err := exec.CommandContext(context.Background(), readelf, "-rW", path).Output()
		if err != nil {
			continue // no relocations/unreadable - skip
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			is64 := ours.Header().Class == elf.CLASS64

			// readelf: a "Offset Info Type ..." line (fixed-width hex fields).
			want := map[uint64][2]uint64{}
			for line := range strings.SplitSeq(string(out), "\n") {
				f := strings.Fields(line)
				if len(f) < 2 || !isHex(f[0]) || !isHex(f[1]) {
					continue
				}

				// a non-numeric field - 0
				off := uint64(0)
				if v, err := strconv.ParseUint(f[0], 16, 64); err == nil {
					off = v
				}

				info := uint64(0)
				if v, err := strconv.ParseUint(f[1], 16, 64); err == nil {
					info = v
				}

				if is64 {
					want[off] = [2]uint64{info >> 32, info & 0xffffffff}
				} else {
					want[off] = [2]uint64{info >> 8, info & 0xff}
				}
			}

			if len(want) == 0 {
				t.Skip("no relocations")
			}

			got := map[uint64][2]uint64{}
			for _, sec := range ours.Sections() {
				rels, err := ours.Relocations(sec)
				require.NoError(t, err, "Relocations(%s)", sec.Name)
				for _, r := range rels {
					got[r.Off] = [2]uint64{uint64(r.SymIndex), uint64(r.Type)}
				}
			}

			require.Len(t, got, len(want), "relocation count")
			for off, wv := range want {
				require.Contains(t, got, off, "offset %#x: readelf has it, we do not", off)
				require.Equal(t, wv, got[off], "offset %#x: (sym,type)", off)
			}
		})
	}
}

func isHex(s string) bool {
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return false
		}
	}

	return len(s) > 0
}

// TestRelocNames checks the type names against the expected ones (the numbers
// verified with llvm-readelf).
func TestRelocNames(t *testing.T) {
	f, err := elf.Open("testdata/relocs-aarch64.o")
	if err != nil {
		t.Skip("no fixture")
	}

	rels, err := f.Relocations(f.Section(".text"))
	require.NoError(t, err)
	want := []string{
		"R_AARCH64_PLT32", "R_AARCH64_MOVW_PREL_G0", "R_AARCH64_MOVW_PREL_G1",
		"R_AARCH64_MOVW_GOTOFF_G1", "R_AARCH64_GOTREL64", "R_AARCH64_GOTREL32",
		"R_AARCH64_GOT_LD_PREL19", "R_AARCH64_ADR_GOT_PAGE",
		"R_AARCH64_LD64_GOTOFF_LO15", "R_AARCH64_LD64_GOT_LO12_NC",
		"R_AARCH64_LD64_GOTPAGE_LO15", "R_AARCH64_TLSGD_ADR_PREL21",
		"R_AARCH64_TLSLE_LDST128_TPREL_LO12_NC",
	}
	require.Len(t, rels, len(want), "relocation count")
	for i, w := range want {
		require.Equal(t, w, f.RelocName(rels[i].Type), "relocation %d", i)
	}
}

// TestOpenGarbage: garbage and truncated files must produce an error, not a panic.
func TestOpenGarbage(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		{0x7f},
		{0x7f, 'E', 'L'},
		{0x7f, 'E', 'L', 'F'},             // truncated right after the magic
		{0x7f, 'E', 'L', 'F', 2, 1, 1, 0}, // truncated in e_ident
		{0x7f, 'E', 'L', 'F', 3, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 183, 0},
		{
			0x7f,
			'E',
			'L',
			'F',
			2,
			3,
			1,
			0,
		}, // unknown EI_DATA
		{
			0x7f,
			'E',
			'L',
			'F',
			2,
			1,
			1,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			2,
			0,
			183,
			0,
			1,
			0,
			0,
			0,
		}, // shoff out of bounds
	}
	for i, c := range cases {
		tmp := t.TempDir()
		path := tmp + "/garbage"
		err := os.WriteFile(path, c, 0o644)
		require.NoError(t, err)
		func() {
			defer func() {
				r := recover()
				require.Nil(t, r, "case %d: panic: %v", i, r)
			}()
			_, err := elf.Open(path)
			require.Error(t, err, "case %d (%d bytes): an error is expected", i, len(c))
		}()
	}
}

// TestOpenNotElf: a file with a different magic must be an error.
func TestOpenNotElf(t *testing.T) {
	path := t.TempDir() + "/notelf"
	err := os.WriteFile(path, []byte("MZ\x90\x00 garbage"), 0o644)
	require.NoError(t, err)
	_, err = elf.Open(path)
	require.Error(t, err, "not an ELF file")
}

// TestDiffDynamic compares .dynamic with the output of llvm-readelf -dW
// (tag + hex value).
func TestDiffDynamic(t *testing.T) {
	readelf := llvmReadelf(t)

	for _, path := range corpus(t) {
		out, err := exec.CommandContext(context.Background(), readelf, "-dW", path).Output()
		if err != nil {
			continue // no dynamics - skip
		}

		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := elf.Open(path)
			require.NoError(t, err, "Open")
			entries, err := ours.Dynamic()
			require.NoError(t, err, "Dynamic")

			// readelf: " 0x0000000000000004 (HASH) 0xb1d90" or
			// " 0x000000000000000b (SYMENT) 24 (bytes)" - the address is hex or decimal.
			type tv struct {
				tag, val uint64
			}
			want := []tv{}
			for line := range strings.SplitSeq(string(out), "\n") {
				f := strings.Fields(line)
				if len(f) < 3 || !strings.HasPrefix(f[0], "0x") || !strings.HasPrefix(f[1], "(") {
					continue
				}

				// a non-numeric field - 0
				tag := uint64(0)
				if v, err := strconv.ParseUint(strings.TrimPrefix(f[0], "0x"), 16, 64); err == nil {
					tag = v
				}

				valStr := f[2]
				base := 16
				if !strings.HasPrefix(valStr, "0x") {
					base = 10
				} else {
					valStr = strings.TrimPrefix(valStr, "0x")
				}

				val, err := strconv.ParseUint(valStr, base, 64)
				if err != nil {
					continue
				}

				want = append(want, tv{
					tag,
					val,
				})
			}

			if len(want) == 0 {
				t.Skip("no dynamics")
			}

			require.GreaterOrEqual(t, len(entries), len(want), "DT entry count")
			// Compared as a subsequence: our parser does not readelf-parse the
			// flag lines (FLAGS_1 "Flags: PIE") - they are skipped, but the
			// order of the rest matches the file order.
			i := 0
			for _, e := range entries {
				if i >= len(want) {
					break
				}

				if uint64(e.Tag) != want[i].tag {
					continue
				}

				require.Equal(t, want[i].val, e.Val, "tag %v", e.Tag)
				i++
			}

			require.Equal(t, len(want), i, "matched DT entries")
		})
	}
}

// TestNotesBuildID checks notes and the build id on Go binaries.
func TestNotesBuildID(t *testing.T) {
	for _, name := range []string{"hello-riscv64", "dyn-pie-arm64"} {
		t.Run(name, func(t *testing.T) {
			f, err := elf.Open(filepath.Join("testdata", name))
			if err != nil {
				t.Skipf("no fixture: %v", err)
			}

			notes, err := f.Notes()
			require.NoError(t, err, "Notes")
			found := false
			for _, n := range notes {
				if n.Name == "Go" && n.Type == elf.NT_GO_BUILD_ID && len(n.Desc) > 0 {
					found = true
				}
			}

			require.True(t, found, "no Go build-id note (%d total)", len(notes))
			id, err := f.BuildID()
			require.NoError(t, err, "BuildID")
			require.NotEmpty(t, id, "build id")
			t.Logf("build id: %.32s...", id)
		})
	}
}

// TestSysVHash parses .hash (dyn-pie-arm64) and checks the sizes against dynsym.
func TestSysVHash(t *testing.T) {
	f, err := elf.Open("testdata/dyn-pie-arm64")
	if err != nil {
		t.Skip("no fixture")
	}

	h, err := f.SysVHash()
	require.NoError(t, err, "SysVHash (want a table in the dyn-pie fixture)")
	syms, err := f.DynamicSymbols()
	require.NoError(t, err, "DynamicSymbols")
	// nchain == the number of .dynsym entries, including the zeroth (we skip it).
	require.Equal(t, len(syms)+1, int(h.NChain), "nchain (+1 for the zeroth)")
	for _, b := range h.Buckets {
		require.LessOrEqual(t, int(b), int(h.NChain), "bucket %#x is outside chains", b)
	}
}

// TestGnuHashSynthetic builds a minimal ET_DYN with .gnu.hash/versions and
// checks their parsing (there are no real glibc binaries on the host).
func TestGnuHashSynthetic(t *testing.T) {
	img, dynOff := buildSyntheticDynamic(t)
	path := t.TempDir() + "/synth.elf"
	err := os.WriteFile(path, img, 0o644)
	require.NoError(t, err)
	f, err := elf.Open(path)
	require.NoError(t, err, "Open")
	_, err = f.Symbols()
	require.ErrorIs(t, err, elf.ErrNoSymbols, "no symtab")

	entries, err := f.Dynamic()
	require.NoError(t, err, "Dynamic")
	require.Len(t, entries, 7, "DT entry count") // 6 tags + DT_NULL

	soname, err := f.Soname()
	require.NoError(t, err, "Soname")
	require.Equal(t, "libsynth.so", soname)

	h, err := f.GnuHash()
	require.NoError(t, err, "GnuHash (want a table in the synthetic dynamics)")
	require.Equal(t, uint32(1), h.NBuckets, "gnu.hash nbuckets")
	require.Equal(t, uint32(2), h.SymOffset, "gnu.hash symoffset")
	require.Len(t, h.Bloom, 1, "gnu.hash bloom")
	require.Len(t, h.Buckets, 1, "gnu.hash buckets")
	require.Len(t, h.Chains, 1, "gnu.hash chains")

	vn, err := f.VersionNeeds()
	require.NoError(t, err, "VersionNeeds")
	require.Len(t, vn, 1, "verneed")
	require.Equal(t, "libext.so", vn[0].File, "verneed file")
	require.Len(t, vn[0].Entries, 1, "verneed entries")
	require.Equal(t, "VER_X", vn[0].Entries[0].Name, "vernaux name")
	require.Equal(t, uint16(2), vn[0].Entries[0].Idx, "vernaux idx")

	vd, err := f.VersionDefs()
	require.NoError(t, err, "VersionDefs")
	require.Len(t, vd, 1, "verdef")
	require.Equal(t, "VER_Y", vd[0].Name, "verdef name")
	require.Equal(t, uint16(3), vd[0].Idx, "verdef idx")

	vs, err := f.SymbolVersions()
	require.NoError(t, err, "SymbolVersions")
	require.Len(t, vs, 4, "versym")
	require.Equal(t, uint16(2), vs[2], "versym[2]")
	require.Equal(t, uint16(3), vs[3], "versym[3]")
	_ = dynOff
}

// buildSyntheticDynamic builds a minimal ELF64 LE ET_DYN: a PT_LOAD with an
// identity mapping vaddr == offset, PT_DYNAMIC, .dynstr, .gnu.hash,
// .gnu.version, .gnu.version_r, .gnu.version_d, and sections for them.
func buildSyntheticDynamic(t *testing.T) ([]byte, int) {
	t.Helper()
	const (
		ehdrSize = 64
		phdrSize = 56
	)

	// The dynamic string table.
	strtab := append([]byte{0}, []byte("libsynth.so\x00VER_Y\x00libext.so\x00VER_X\x00")...)

	// .gnu.hash: 1 bucket, symoffset=2, bloom=1 word.
	gnuHash := make([]byte, 0, 32)
	gnuHash = append(gnuHash, le32(1)...)          // nbuckets
	gnuHash = append(gnuHash, le32(2)...)          // symoffset
	gnuHash = append(gnuHash, le32(1)...)          // bloom_size
	gnuHash = append(gnuHash, le32(5)...)          // bloom_shift
	gnuHash = append(gnuHash, le64(0xdeadbeef)...) // bloom[0]
	gnuHash = append(gnuHash, le32(2)...)          // buckets[0]
	gnuHash = append(gnuHash, le32(0x11|1)...)     // chains[0] (terminal)

	// .gnu.version: 4 entries.
	versym := make([]byte, 0, 8)
	versym = append(versym, le16(0)...)
	versym = append(versym, le16(1)...)
	versym = append(versym, le16(2)...)
	versym = append(versym, le16(3)...)

	// .gnu.version_r: vn(version=1, cnt=1, file=off("libext.so"), aux=16, next=0)
	// + vernaux(hash, flags=0, other=2, name=off("VER_X"), next=0).
	offVerY := 1 + len("libsynth.so") + 1
	offLibext := offVerY + len("VER_Y") + 1
	offVerX := offLibext + len("libext.so") + 1
	verneed := make([]byte, 0, 32)
	verneed = append(verneed, le16(1)...) // vn_version
	verneed = append(verneed, le16(1)...) // vn_cnt
	verneed = append(verneed, le32(uint32(offLibext))...)
	verneed = append(verneed, le32(16)...)     // vn_aux
	verneed = append(verneed, le32(0)...)      // vn_next
	verneed = append(verneed, le32(0x1234)...) // vna_hash
	verneed = append(verneed, le16(0)...)      // vna_flags
	verneed = append(verneed, le16(2)...)      // vna_other
	verneed = append(verneed, le32(uint32(offVerX))...)
	verneed = append(verneed, le32(0)...) // vna_next

	// .gnu.version_d: vd(version=1, flags=0, ndx=3, cnt=1, hash, aux=20, next=0)
	// + verdaux(name=off("VER_Y"), next=0).
	verdef := make([]byte, 0, 28)
	verdef = append(verdef, le16(1)...)      // vd_version
	verdef = append(verdef, le16(0)...)      // vd_flags
	verdef = append(verdef, le16(3)...)      // vd_ndx
	verdef = append(verdef, le16(1)...)      // vd_cnt
	verdef = append(verdef, le32(0x5678)...) // vd_hash
	verdef = append(verdef, le32(20)...)     // vd_aux
	verdef = append(verdef, le32(0)...)      // vd_next
	verdef = append(verdef, le32(uint32(offVerY))...)
	verdef = append(verdef, le32(0)...) // vda_next

	// The dynamic table (the addresses are filled in after the layout).
	blocks := [][]byte{strtab, gnuHash, versym, verneed, verdef}
	offs := make([]int, len(blocks))

	phoff := ehdrSize
	dynOff := phoff + phdrSize*2
	pos := dynOff
	// Space for the dynamics: 7 entries x 16.
	pos += 7 * 16
	for i, b := range blocks {
		pos = (pos + 7) &^ 7 // 8-byte alignment
		offs[i] = pos
		pos += len(b)
	}

	strOff, hashOff, vsymOff, vneedOff, vdefOff := offs[0], offs[1], offs[2], offs[3], offs[4]

	dyn := make([]byte, 0, 7*16)
	appendDyn := func(tag, val uint64) {
		dyn = append(dyn, le64(tag)...)
		dyn = append(dyn, le64(val)...)
	}
	appendDyn(1, 1)                         // DT_NEEDED -> strtab[1] = "libsynth.so"
	appendDyn(14, 1)                        // DT_SONAME -> strtab[1]
	appendDyn(5, uint64(strOff))            // DT_STRTAB
	appendDyn(10, uint64(len(strtab)))      // DT_STRSZ
	appendDyn(0x6ffffef5, uint64(hashOff))  // DT_GNU_HASH
	appendDyn(0x6ffffffe, uint64(vneedOff)) // DT_VERNEED
	appendDyn(0, 0)                         // DT_NULL

	// Assemble the image: ehdr + 2 phdr + dyn + blocks + sections.
	shstrtab := []byte(
		"\x00.dynstr\x00.gnu.hash\x00.gnu.version\x00.gnu.version_r\x00.gnu.version_d\x00.shstrtab\x00",
	)

	secs := []struct {
		name string
		typ  uint32
		off  int
		size int
		info uint32
		link uint32
	}{
		{
			".dynstr",
			3,
			strOff,
			len(strtab),
			0,
			0,
		},
		{
			".gnu.hash",
			0x6ffffff6,
			hashOff,
			len(gnuHash),
			0,
			1,
		},
		{
			".gnu.version",
			0x6fffffff,
			vsymOff,
			len(versym),
			0,
			1,
		},
		{
			".gnu.version_r",
			0x6ffffffe,
			vneedOff,
			len(verneed),
			0,
			1,
		},
		{
			".gnu.version_d",
			0x6ffffffd,
			vdefOff,
			len(verdef),
			0,
			1,
		},
		{
			".shstrtab",
			3,
			0,
			len(shstrtab),
			0,
			0,
		},
	}
	shstrOff := pos
	pos += len(shstrtab)
	pos = (pos + 7) &^ 7
	shoff := pos

	total := shoff + 64*(len(secs)+1) // + the zero section

	img := make([]byte, total)
	copy(img[dynOff:], dyn)
	for i, b := range blocks {
		copy(img[offs[i]:], b)
	}

	copy(img[shstrOff:], shstrtab)

	// --- Ehdr ---
	e := img
	e[0] = 0x7f
	copy(e[1:4], "ELF")
	e[4] = 2            // CLASS64
	e[5] = 1            // DATA2LSB
	e[6] = 1            // version
	put16(img, 16, 3)   // ET_DYN
	put16(img, 18, 183) // EM_AARCH64
	put32(img, 20, 1)
	put64(img, 32, uint64(phoff))
	put64(img, 40, uint64(shoff))
	put16(img, 52, 64)                    // ehsize
	put16(img, 54, 56)                    // phentsize
	put16(img, 56, 2)                     // phnum
	put16(img, 58, 64)                    // shentsize
	put16(img, 60, uint16(len(secs)+1))   // shnum
	put16(img, 62, uint16(len(secs)+1-1)) // shstrndx = the last section

	// --- Phdrs: PT_LOAD (the whole file) + PT_DYNAMIC ---
	p := phoff
	put32(img, p, 1)   // PT_LOAD
	put32(img, p+4, 4) // PF_R
	put64(img, p+8, 0)
	put64(img, p+16, 0)
	put64(img, p+24, 0)
	put64(img, p+32, uint64(total))
	put64(img, p+40, uint64(total))
	put64(img, p+48, 0x1000)
	p += phdrSize
	put32(img, p, 2) // PT_DYNAMIC
	put32(img, p+4, 4)
	put64(img, p+8, uint64(dynOff))
	put64(img, p+16, uint64(dynOff))
	put64(img, p+24, uint64(dynOff))
	put64(img, p+32, uint64(len(dyn)))
	put64(img, p+40, uint64(len(dyn)))
	put64(img, p+48, 8)

	// --- Shdrs ---
	writeShdr := func(idx int, nameOff, typ uint32, off, size int, link uint32) {
		b := shoff + idx*64
		put32(img, b, nameOff)
		put32(img, b+4, typ)
		put64(img, b+16, uint64(off)) // sh_addr: vaddr == offset (identity)
		put64(img, b+24, uint64(off))
		put64(img, b+32, uint64(size))
		put32(img, b+40, link)
	}
	writeShdr(0, 0, 0, 0, 0, 0)
	for i, s := range secs {
		nameOff := shstrtabOff(shstrtab, s.name)
		writeShdr(i+1, nameOff, s.typ, s.off, s.size, s.link)
	}

	// The shstrtab section refers to itself: off is filled in separately.
	b := shoff + len(secs)*64
	put64(img, b+24, uint64(shstrOff))
	return img, dynOff
}

func shstrtabOff(tab []byte, name string) uint32 {
	i := indexOf(tab, []byte(name))
	if i < 0 {
		return 0
	}

	return uint32(i)
}

func indexOf(hay, needle []byte) int {
	for i := 0; i+len(needle) <= len(hay); i++ {
		match := true
		for j := range needle {
			if hay[i+j] != needle[j] {
				match = false
				break
			}
		}

		if match {
			return i
		}
	}

	return -1
}

func le16(v uint16) []byte {
	return []byte{byte(v), byte(v >> 8)}
}
func le32(v uint32) []byte {
	return []byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}
}
func le64(v uint64) []byte {
	return append(le32(uint32(v)), le32(uint32(v>>32))...)
}

func put16(b []byte, off int, v uint16) {
	b[off] = byte(v)
	b[off+1] = byte(v >> 8)
}
func put32(b []byte, off int, v uint32) {
	copy(b[off:], le32(v))
}
func put64(b []byte, off int, v uint64) {
	copy(b[off:], le64(v))
}
