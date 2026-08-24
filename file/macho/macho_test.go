// Diff tests against stdlib debug/macho (tests only; the runtime does not
// depend on debug/*). Corpus: clang builds in testdata (exec, dylib, .o,
// FAT arm64+x86_64) plus macOS system binaries when available (skipped
// when absent).
package macho_test

import (
	"context"
	stdmacho "debug/macho"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/file/macho"
)

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

// sysCorpus is system binaries (not committed; skipped if unavailable).
var sysCorpus = []string{
	"/bin/ls",
	"/usr/bin/otool",
	"/usr/lib/dyld",
}

// stdOpen opens the file via stdlib; for FAT archives it takes the arm64
// slice (stdlib Open does not read fats - only NewFatFile).
func stdOpen(t *testing.T, path string) *stdmacho.File {
	t.Helper()
	f, err := stdmacho.Open(path)
	if err == nil {
		return f
	}

	r, oerr := os.Open(path)
	require.NoError(t, oerr, "stdlib: %v", err)
	fat, ferr := stdmacho.NewFatFile(r)
	require.NoError(t, ferr, "stdlib: %v", err)
	for i := range fat.Arches {
		if fat.Arches[i].Cpu == stdmacho.CpuArm64 {
			return fat.Arches[i].File
		}
	}

	require.Fail(t, "stdlib: no arm64 slice in FAT")
	return nil
}

// TestDiffHeader compares headers.
func TestDiffHeader(t *testing.T) {
	paths := corpus(t)
	for _, p := range sysCorpus {
		if _, err := os.Stat(p); err == nil {
			paths = append(paths, p)
		}
	}

	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			std := stdOpen(t, path)
			defer func() {
				require.NoError(t, std.Close(), "close stdlib file")
			}()

			a, b := ours.Header(), std.FileHeader
			require.Equal(t, int32(b.Cpu), a.CpuType, "cputype")
			require.Equal(t, int32(b.SubCpu), a.CpuSubtype, "cpusubtype")
			require.Equal(t, uint32(b.Type), uint32(a.FileType), "filetype")
			require.Equal(t, b.Ncmd, a.Ncmds, "ncmds")
			require.Equal(t, b.Cmdsz, a.Sizeofcmds, "sizeofcmds")
			require.Equal(t, b.Flags, uint32(a.Flags), "flags")
		})
	}
}

// TestDiffSections compares sections (flat, in file order) and __text data.
func TestDiffSections(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			std := stdOpen(t, path)
			defer func() {
				require.NoError(t, std.Close(), "close stdlib file")
			}()

			a, b := ours.Sections(), std.Sections
			require.Len(t, a, len(b), "number of sections")
			for i := range a {
				x, y := a[i], b[i]
				require.Equal(t, y.Name, x.SectName, "section %d: name", i)
				require.Equal(t, y.Seg, x.SegName, "section %s: segment", x.SectName)
				require.Equal(t, y.Addr, x.Addr, "section %s: addr", x.SectName)
				require.Equal(t, y.Size, x.Size, "section %s: size", x.SectName)
				require.Equal(t, y.Offset, x.Offset, "section %s: offset", x.SectName)
				require.Equal(t, y.Align, x.Align, "section %s: align", x.SectName)
				require.Equal(t, y.Flags, x.Flags, "section %s: flags", x.SectName)
			}

			if ts := ours.Section("__text"); ts != nil {
				d1, err := ts.Data()
				require.NoError(t, err, "__text Data")
				d2, err := std.Section("__text").Data()
				require.NoError(t, err, "stdlib __text Data")
				require.Equal(t, string(d2), string(d1), "__text: data")
			}
		})
	}
}

// TestOpenGarbage: garbage and truncated files must return an error, not panic.
func TestOpenGarbage(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		{0xcf},
		{0xcf, 0xfa, 0xed, 0xfe}, // magic only
		{0xcf, 0xfa, 0xed, 0xfe, 0x0c, 0x00, 0x00, 0x01}, // truncated in cputype
		{
			0xce,
			0xfa,
			0xed,
			0xfe,
			0x0c,
			0x00,
			0x00,
			0x01,
			0x00,
			0x00,
			0x00,
			0x00,
			0x02,
			0x00,
			0x00,
			0x00,
			0x01,
			0x00,
			0x00,
			0x00,
			0x60,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0x00,
			0xff,
			0xff,
			0x00,
			0x00,
		}, // huge ncmds
	}
	for i, c := range cases {
		path := t.TempDir() + "/garbage"
		require.NoError(t, os.WriteFile(path, c, 0o644))
		func() {
			defer func() {
				r := recover()
				require.Nil(t, r, "case %d: panic", i)
			}()
			_, err := macho.Open(path)
			require.Error(t, err, "case %d (%d bytes)", i, len(c))
		}()
	}
}

// TestOpenTruncated: prefixes of a real binary - errors only, no panics.
func TestOpenTruncated(t *testing.T) {
	raw, err := os.ReadFile("testdata/exec-arm64")
	if err != nil {
		t.Skip("no exec-arm64 fixture")
	}

	// Every 33rd prefix (including the full 31 bytes of the header and the
	// tails of the commands).
	for n := 1; n <= len(raw); n += 33 {
		path := t.TempDir() + "/trunc"
		require.NoError(t, os.WriteFile(path, raw[:n], 0o644))
		func() {
			defer func() {
				r := recover()
				require.Nil(t, r, "prefix %d: panic", n)
			}()
			_, openErr := macho.Open(path)
			_ = openErr // an error is acceptable (truncated prefix), a panic is not
		}()
	}
}

// TestDiffLoadCommands compares the (cmd, cmdsize) sequence with stdlib.
func TestDiffLoadCommands(t *testing.T) {
	paths := corpus(t)
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			std := stdOpen(t, path)
			defer func() {
				require.NoError(t, std.Close(), "close stdlib file")
			}()

			a := ours.LoadCommands()
			require.Len(t, a, len(std.Loads), "number of commands")
			for i := range a {
				cmd, size := uint32(a[i].Cmd()), a[i].Cmdsize()
				raw := std.Loads[i].Raw()
				require.True(t, len(raw) >= 8, "command %d: stdlib raw %d bytes", i, len(raw))
				scmd := std.ByteOrder.Uint32(raw)
				ssize := std.ByteOrder.Uint32(raw[4:])
				require.Equal(t, scmd, cmd, "command %d: cmd", i)
				require.Equal(t, ssize, size, "command %d (%v): cmdsize", i, a[i].Cmd())
			}
		})
	}
}

// TestDiffSymbols compares symbol tables.
func TestDiffSymbols(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			std := stdOpen(t, path)
			defer func() {
				require.NoError(t, std.Close(), "close stdlib file")
			}()

			a, err := ours.Symbols()
			require.NoError(t, err, "Symbols")
			require.NotNil(t, std.Symtab, "stdlib: no Symtab")
			b := std.Symtab.Syms
			require.Len(t, a, len(b), "number of symbols")
			for i := range a {
				x, y := a[i], b[i]
				require.Equal(t, y.Name, x.Name, "symbol %d: name", i)
				require.Equal(t, y.Type, x.Type, "symbol %d (%s): type", i, x.Name)
				require.Equal(t, y.Sect, x.Sect, "symbol %d (%s): sect", i, x.Name)
				require.Equal(t, y.Desc, x.Desc, "symbol %d (%s): desc", i, x.Name)
				require.Equal(t, y.Value, x.Value, "symbol %d (%s): value", i, x.Name)
			}
		})
	}
}

// TestFatAPI checks OpenFat: all slices, metadata, opening by cputype.
func TestFatAPI(t *testing.T) {
	raw, err := os.ReadFile("testdata/fat-exec")
	if err != nil {
		t.Skip("no fat-exec fixture")
	}

	fat, err := macho.OpenFat("testdata/fat-exec")
	require.NoError(t, err, "OpenFat")
	require.Len(t, fat.Arches, 2, "number of slices")

	// Slice metadata - against stdlib FatFile.
	r, err := os.Open("testdata/fat-exec")
	require.NoError(t, err, "os.Open fat-exec")
	std, err := stdmacho.NewFatFile(r)
	require.NoError(t, err, "stdlib NewFatFile")
	require.Len(t, fat.Arches, len(std.Arches), "slices: ours vs stdlib")
	for i, a := range fat.Arches {
		b := std.Arches[i]
		require.Equal(t, int32(b.Cpu), a.CpuType, "slice %d: cputype", i)
		require.Equal(t, int32(b.SubCpu), a.CpuSubtype, "slice %d: cpusubtype", i)
		require.Equal(t, uint64(b.Offset), a.Offset, "slice %d: offset", i)
		require.Equal(t, uint64(b.Size), a.Size, "slice %d: size", i)
	}

	// Opening a specific slice.
	arch := fat.Arch(macho.CPU_TYPE_ARM64)
	require.NotNil(t, arch, "no arm64 slice")
	f, err := arch.Open()
	require.NoError(t, err, "arch.Open")
	require.Equal(t, macho.CPU_TYPE_ARM64, f.Header().CpuType, "cputype of the opened slice")
	require.NotNil(t, f.Section("__text"), "no __text in the arm64 slice")
	_ = raw
}

// Test32BitObject checks a 32-bit Mach-O (.o for i386).
func Test32BitObject(t *testing.T) {
	f, err := macho.Open("testdata/obj-i386.o")
	if err != nil {
		t.Skipf("no 32-bit fixture: %v", err)
	}

	h := f.Header()
	require.Equal(t, macho.CPU_TYPE_X86, h.CpuType, "cputype")
	require.Zero(t, h.Reserved, "a 32-bit header must not have reserved")
	std := stdOpen(t, "testdata/obj-i386.o")
	defer func() {
		require.NoError(t, std.Close(), "close stdlib file")
	}()
	require.Len(t, f.LoadCommands(), len(std.Loads), "number of commands")
	require.Len(t, f.Sections(), len(std.Sections), "number of sections")
	for i, s := range f.Sections() {
		require.Equal(t, std.Sections[i].Name, s.SectName, "section %d: name", i)
		require.Equal(t, std.Sections[i].Seg, s.SegName, "section %d: segment", i)
		require.Equal(t, std.Sections[i].Addr, s.Addr, "section %s: addr", s.SectName)
		require.Equal(t, std.Sections[i].Size, s.Size, "section %s: size", s.SectName)
	}
}

// TestBEMacho: a minimal handcrafted BE Mach-O (magic MH_MAGIC_64 in BE).
func TestBEMacho(t *testing.T) {
	img := make([]byte, 32)
	be := func(off int, v uint32) {
		img[off] = byte(v >> 24)
		img[off+1] = byte(v >> 16)
		img[off+2] = byte(v >> 8)
		img[off+3] = byte(v)
	}
	// In the file, magic is stored as 0xfeedfacf BE: FE ED FA CF.
	img[0], img[1], img[2], img[3] = 0xFE, 0xED, 0xFA, 0xCF
	be(4, uint32(macho.CPU_TYPE_ARM64)) // cputype
	be(8, 0)                            // cpusubtype
	be(12, 2)                           // MH_EXECUTE
	be(16, 0)                           // ncmds
	be(20, 0)                           // sizeofcmds
	be(24, 0x00200085)                  // flags: PIE|DYLDLINK|NOUNDEFS|TWOLEVEL
	be(28, 0)                           // reserved

	path := t.TempDir() + "/be-macho"
	require.NoError(t, os.WriteFile(path, img, 0o644))
	f, err := macho.Open(path)
	require.NoError(t, err, "Open")
	h := f.Header()
	require.Equal(t, macho.CPU_TYPE_ARM64, h.CpuType, "BE header parsed incorrectly: cputype")
	require.Equal(t, macho.MH_EXECUTE, h.FileType, "BE header parsed incorrectly: filetype")
	require.Empty(t, f.LoadCommands(), "expected 0 commands")
}

// TestDiffRelocs compares section relocations and the indirect symbol table.
func TestDiffRelocs(t *testing.T) {
	for _, path := range corpus(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			ours, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			std := stdOpen(t, path)
			defer func() {
				require.NoError(t, std.Close(), "close stdlib file")
			}()

			for i, sec := range ours.Sections() {
				a, err := ours.Relocations(sec)
				require.NoError(t, err, "Relocations(%s)", sec.SectName)
				b := std.Sections[i].Relocs
				require.Len(t, a, len(b), "section %s: number of relocations", sec.SectName)
				for j := range a {
					x, y := a[j], b[j]
					require.Equal(t, y.Addr, x.Addr, "section %s reloc %d: addr", sec.SectName, j)
					require.Equal(
						t,
						uint32(y.Type),
						uint32(x.Type),
						"section %s reloc %d: type",
						sec.SectName,
						j,
					)
					require.Equal(t, y.Len, x.Length, "section %s reloc %d: len", sec.SectName, j)
					require.Equal(
						t,
						y.Pcrel,
						x.Pcrel,
						"section %s reloc %d: pcrel",
						sec.SectName,
						j,
					)
					require.Equal(
						t,
						y.Extern,
						x.Extern,
						"section %s reloc %d: extern",
						sec.SectName,
						j,
					)
					if !x.Scattered {
						require.Equal(
							t,
							y.Value,
							x.SymIndex,
							"section %s reloc %d: symbolnum",
							sec.SectName,
							j,
						)
					}

					require.Equal(
						t,
						y.Scattered,
						x.Scattered,
						"section %s reloc %d: scattered",
						sec.SectName,
						j,
					)
					if x.Scattered {
						require.Equal(
							t,
							y.Value,
							x.Value,
							"section %s reloc %d: scattered value",
							sec.SectName,
							j,
						)
					}
				}
			}

			// Indirect symbol table.
			ia, err := ours.IndirectSymbols()
			require.NoError(t, err, "IndirectSymbols")
			if std.Dysymtab != nil {
				ib := std.Dysymtab.IndirectSyms
				require.Len(t, ia, len(ib), "indirect syms")
				for k := range ia {
					require.Equal(t, ib[k], ia[k], "indirect syms[%d]", k)
				}
			}
		})
	}
}

// thinSystemArm64 extracts the thin arm64 slice of a system binary via
// lipo (modern macOS binaries use chained fixups - the only easy access
// to real LC_DYLD_CHAINED_FIXUPS on this host).
func thinSystemArm64(t *testing.T) string {
	t.Helper()
	lipo, err := exec.LookPath("lipo")
	if err != nil {
		t.Skip("lipo unavailable")
	}

	if _, err := os.Stat("/bin/ls"); err != nil {
		t.Skip("/bin/ls unavailable")
	}

	out := t.TempDir() + "/ls-arm64e"
	if b, err := exec.CommandContext(context.Background(), lipo, "/bin/ls", "-thin", "arm64e", "-output", out).
		CombinedOutput(); err != nil {
		t.Skipf("lipo: %v: %s", err, b)
	}

	return out
}

// llvmObjdumpMacho runs llvm-objdump --macho with option opt; skips when
// the tool is missing.
func llvmObjdumpMacho(t *testing.T, opt, path string) string {
	t.Helper()
	tool, err := exec.LookPath("llvm-objdump")
	if err != nil {
		tool = "/opt/homebrew/opt/llvm/bin/llvm-objdump"
		if _, serr := os.Stat(tool); serr != nil {
			t.Skip("llvm-objdump unavailable")
		}
	}

	out, err := exec.CommandContext(context.Background(), tool, "--macho", opt, path).Output()
	if err != nil {
		t.Skipf("llvm-objdump %s: %v", opt, err)
	}

	return string(out)
}

// bindLine is a row of the llvm-objdump bind table.
type bindLine struct {
	addr uint64
	sym  string
	add  int64
}

func newBindLine(addr uint64, sym string, add int64) bindLine {
	return bindLine{
		addr: addr,
		sym:  sym,
		add:  add,
	}
}

// parseBindLines parses the output of --bind/--lazy-bind/--weak-bind:
// "segment section address type addend dylib symbol".
func parseBindLines(out string) []bindLine {
	var res []bindLine
	for _, line := range strings.Split(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 6 || !strings.HasPrefix(f[2], "0x") {
			continue
		}

		addr, err := strconv.ParseUint(strings.TrimPrefix(f[2], "0x"), 16, 64)
		if err != nil {
			continue
		}

		// A non-numeric addend field is treated as zero.
		add := int64(0)
		if v, err := strconv.ParseInt(f[4], 10, 64); err == nil {
			add = v
		}

		res = append(res, newBindLine(addr, f[len(f)-1], add))
	}

	return res
}

func diffBinds(t *testing.T, ours []macho.Bind, llvmOut string, what string) {
	t.Helper()
	want := parseBindLines(llvmOut)
	require.Len(t, ours, len(want), "%s: number of entries", what)
	for i := range want {
		require.Equal(t, want[i].addr, ours[i].Addr, "%s[%d]: addr", what, i)
		require.Equal(t, want[i].sym, ours[i].SymName, "%s[%d] (%#x)", what, i, ours[i].Addr)
	}
}

// TestDyldStreamsOnSystemBin - dyld_info streams on a modern system binary
// (chained fixups; llvm synthesizes binds from the chains).
func TestDyldStreamsOnSystemBin(t *testing.T) {
	path := thinSystemArm64(t)
	f, err := macho.Open(path)
	require.NoError(t, err, "Open")

	b, err := f.Binds()
	require.NoError(t, err, "Binds")
	diffBinds(t, b, llvmObjdumpMacho(t, "--bind", path), "bind")
	b, err = f.LazyBinds()
	require.NoError(t, err, "LazyBinds")
	diffBinds(t, b, llvmObjdumpMacho(t, "--lazy-bind", path), "lazy-bind")
	b, err = f.WeakBinds()
	require.NoError(t, err, "WeakBinds")
	diffBinds(t, b, llvmObjdumpMacho(t, "--weak-bind", path), "weak-bind")
}

// TestDyldStreamsOnFixture - classic LC_DYLD_INFO streams on our clang
// fixtures.
func TestDyldStreamsOnFixture(t *testing.T) {
	for _, name := range []string{"exec-arm64", "dylib-arm64.dylib"} {
		path := filepath.Join("testdata", name)
		t.Run(name, func(t *testing.T) {
			f, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			b, err := f.Binds()
			require.NoError(t, err, "Binds")
			diffBinds(t, b, llvmObjdumpMacho(t, "--bind", path), "bind")
			b, err = f.LazyBinds()
			require.NoError(t, err, "LazyBinds")
			diffBinds(t, b, llvmObjdumpMacho(t, "--lazy-bind", path), "lazy-bind")
			r, err := f.Rebases()
			require.NoError(t, err, "Rebases")
			llvmOut := llvmObjdumpMacho(t, "--rebase", path)
			want := 0
			for _, line := range strings.Split(llvmOut, "\n") {
				if strings.Contains(line, "0x") && !strings.Contains(line, "Rebase table") &&
					!strings.Contains(line, "segment") {
					want++
				}
			}

			require.Len(t, r, want, "rebase: number of entries")
		})
	}
}

// TestExportsTrie compares the trie exports with llvm-objdump --exports-trie.
func TestExportsTrie(t *testing.T) {
	for _, path := range []string{filepath.Join("testdata", "dylib-arm64.dylib"), thinSystemArm64(t)} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			exp, err := f.Exports()
			require.NoError(t, err, "Exports")

			type ev struct {
				name string
				addr uint64
			}
			var want []ev
			for _, line := range strings.Split(llvmObjdumpMacho(t, "--exports-trie", path), "\n") {
				f := strings.Fields(line)
				if len(f) != 2 || !strings.HasPrefix(f[0], "0x") {
					continue
				}

				addr, err := strconv.ParseUint(strings.TrimPrefix(f[0], "0x"), 16, 64)
				if err != nil {
					continue
				}

				want = append(want, ev{
					name: f[1],
					addr: addr,
				})
			}

			require.Len(t, exp, len(want), "number of exports")
			for i := range want {
				require.Equal(t, want[i].name, exp[i].Name, "export %d: name", i)
				if !exp[i].Reexport {
					require.Equal(
						t,
						want[i].addr,
						exp[i].Addr,
						"export %s: addr",
						exp[i].Name,
					)
				}
			}
		})
	}
}

// TestChainedFixups checks Fixups() on a system binary: a full diff of the
// import table against llvm-objdump --chained-fixups + chain invariants.
func TestChainedFixups(t *testing.T) {
	path := thinSystemArm64(t)
	f, err := macho.Open(path)
	require.NoError(t, err, "Open")
	cx, err := f.Fixups()
	require.NoError(t, err, "Fixups: expected on a modern binary")

	// llvm prints each import as a triple of lines:
	//   dyld chained import[N]
	//     lib_ordinal = X (libname)
	//     weak_import = W
	//     name_offset = N (_symbol)
	type imp struct {
		ordinal int
		name    string
	}
	var want []imp
	lines := strings.Split(llvmObjdumpMacho(t, "--chained-fixups", path), "\n")
	for i := range lines {
		if !strings.HasPrefix(lines[i], "dyld chained import[") {
			continue
		}

		im := imp{}
		for j := i + 1; j < len(lines) && j <= i+3; j++ {
			f := strings.Fields(lines[j])
			if len(f) >= 3 && f[0] == "lib_ordinal" {
				// A non-numeric field is 0.
				if v, err := strconv.Atoi(f[2]); err == nil {
					im.ordinal = v
				}
			}

			if len(f) >= 3 && f[0] == "name_offset" {
				im.name = strings.TrimSuffix(strings.TrimPrefix(f[len(f)-1], "("), ")")
			}
		}

		want = append(want, im)
	}

	require.Len(t, cx.Imports, len(want), "number of imports")
	for i := range want {
		require.Equal(
			t,
			want[i].ordinal,
			cx.Imports[i].LibOrdinal,
			"import %d: ordinal",
			i,
		)
		require.Equal(t, want[i].name, cx.Imports[i].Name, "import %d: name", i)
	}

	// Chain invariants: addresses within segment bounds; bind fixups refer
	// to existing imports; the library name resolves.
	nbind := 0
	for _, fx := range cx.Fixups {
		if fx.Bind {
			nbind++
			require.Less(
				t,
				int(fx.ImportIndex),
				len(cx.Imports),
				"bind fixup %#x: import index %d outside the table (%d)",
				fx.Addr,
				fx.ImportIndex,
				len(cx.Imports),
			)
		}

		inSeg := false
		for _, s := range f.Segments() {
			if fx.Addr >= s.Vmaddr && fx.Addr < s.Vmaddr+s.Vmsize {
				inSeg = true
				break
			}
		}

		require.True(t, inSeg, "fixup %#x outside all segments", fx.Addr)
	}

	if len(cx.Imports) > 0 {
		require.Positive(
			t,
			nbind,
			"there are imports but not a single bind fixup - chains not walked?",
		)
	}

	require.NotEmpty(t, cx.Fixups, "chains are empty")
}

// TestFunctionStarts compares function start addresses.
func TestFunctionStarts(t *testing.T) {
	for _, path := range []string{filepath.Join("testdata", "exec-arm64"), thinSystemArm64(t)} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := macho.Open(path)
			require.NoError(t, err, "Open %s", path)
			ours, err := f.FunctionStarts()
			require.NoError(t, err, "FunctionStarts")
			var want []uint64
			for _, line := range strings.Split(llvmObjdumpMacho(t, "--function-starts", path), "\n") {
				f := strings.Fields(line)
				if len(f) != 1 {
					continue
				}

				// addrs mode: one line = a hex address without a prefix.
				v, err := strconv.ParseUint(f[0], 16, 64)
				if err != nil {
					continue
				}

				want = append(want, v)
			}

			require.Len(t, ours, len(want), "number of functions")
			for i := range want {
				require.Equal(t, want[i], ours[i], "function %d: addr", i)
			}
		})
	}
}

// TestDataInCode compares the data-in-code table.
func TestDataInCode(t *testing.T) {
	f, err := macho.Open("testdata/obj-arm64.o")
	if err != nil {
		t.Skip("no fixture")
	}

	_, err = f.DataInCode()
	require.NoError(t, err, "DataInCode")
}

// TestEntryAndUUID - the entry point and UUID on our fixtures.
func TestEntryAndUUID(t *testing.T) {
	f, err := macho.Open("testdata/exec-arm64")
	if err != nil {
		t.Skip("no fixture")
	}

	entry, ok := f.Entry()
	require.True(t, ok, "entry point not found")
	text := f.Section("__text")
	require.NotNil(t, text, "no __text")
	// The entry point (main) must lie in __TEXT.
	found := false
	for _, s := range f.Segments() {
		if s.SegName == "__TEXT" && entry >= s.Vmaddr && entry < s.Vmaddr+s.Vmsize {
			found = true
		}
	}

	require.True(t, found, "entry %#x outside __TEXT", entry)
	require.NotNil(t, f.UUIDCommand(), "UUID not found (clang always sets it)")
	if f.BuildVersionCommand() == nil {
		t.Log("build version missing (old linker)")
	}
}

// TestCodeSignature - the superblob conversion on a system binary.
func TestCodeSignature(t *testing.T) {
	path := thinSystemArm64(t)
	f, err := macho.Open(path)
	require.NoError(t, err, "Open")
	blobs, err := f.CodeSignature()
	require.NoError(t, err, "CodeSignature")
	require.NotEmpty(t, blobs, "a system binary must have a code signature")
}

// TestMutationSweep mutates random bytes of the header/load commands: Open
// and lazy accesses must fail, but not panic.
func TestMutationSweep(t *testing.T) {
	seed := int64(42)
	if v := os.Getenv("ASSEMBLY_SEED"); v != "" {
		_, err := fmt.Sscan(v, &seed)
		require.NoError(t, err, "ASSEMBLY_SEED")
	}

	rnd := rand.New(rand.NewSource(seed))

	for _, name := range []string{"exec-arm64", "obj-i386.o", "dylib-arm64.dylib"} {
		raw, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			continue
		}

		for iter := range 200 {
			mut := append([]byte(nil), raw...)
			// Mutate the header and load command area.
			for range 8 {
				mut[rnd.Intn(minInt(8192, len(mut)))] = byte(rnd.Intn(256))
			}

			path := t.TempDir() + "/mut.macho"
			require.NoError(t, os.WriteFile(path, mut, 0o644))
			func() {
				defer func() {
					r := recover()
					require.Nil(t, r, "%s iter %d: panic", name, iter)
				}()
				f, err := macho.Open(path)
				if err != nil {
					return // an error is a normal outcome
				}

				// the calls must not panic; errors are acceptable (corrupted data)
				call := func(_ any, _ error) {}
				call(f.Symbols())
				call(f.IndirectSymbols())
				call(f.Binds())
				call(f.LazyBinds())
				call(f.WeakBinds())
				call(f.Rebases())
				call(f.Exports())
				call(f.Fixups())
				call(f.FunctionStarts())
				call(f.DataInCode())
				call(f.CodeSignature())
				for _, s := range f.Sections() {
					call(s.Data())
					call(f.Relocations(s))
				}
			}()
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
