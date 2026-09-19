package assembly_test

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/okneniz/assembly/asm/arm64/alias"
)

// clangPath - the Apple clang oracle (empty = absent, the check skips).
func clangPath() string {
	for _, c := range []string{"/usr/bin/clang", "clang"} {
		if _, err := exec.LookPath(c); err == nil {
			return c
		}
	}

	return ""
}

// TestSimdClangDual - the decomposed SIMD families byte-compared with
// the host clang oracle over the whole corpus .s (the assembler twin
// of the prog chain pin; the llvm-mc gate covers it in docker too).
func TestSimdDualCheck(t *testing.T) {
	clang := clangPath()
	if clang == "" {
		t.Skip("no clang found on this host")
	}

	src, _ := os.ReadFile("tests/examples/simd/simd-arm64.s")
	res, errs := alias.Assemble(string(src), 0)
	if len(errs) > 0 {
		t.Fatalf("our asm errors: %v", errs)
	}

	dir := t.TempDir()
	cc := exec.Command(clang, "-arch", "arm64", "-c", "-x", "assembler", "-o", filepath.Join(dir, "o.o"), "-")
	cc.Stdin = strings.NewReader(string(src))
	if out, err := cc.CombinedOutput(); err != nil {
		t.Fatalf("clang: %v: %s", err, out)
	}

	od := exec.Command("/usr/bin/objdump", "-d", filepath.Join(dir, "o.o"))
	out, _ := od.Output()
	want := []uint32{}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		for _, f := range fields {
			if len(f) == 8 && isHex(f) {
				if v, err := parseHex(f); err == nil {
					want = append(want, v)
				}
				break
			}
		}
	}

	got := []uint32{}
	for i := 0; i+4 <= len(res.Sections[0].Data); i += 4 {
		got = append(got, binary.LittleEndian.Uint32(res.Sections[0].Data[i:]))
	}

	if len(want) != len(got) {
		t.Fatalf("count: clang %d vs ours %d", len(want), len(got))
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("word %d: clang %08x vs ours %08x", i, want[i], got[i])
		}
	}
}

func isHex(s string) bool {
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return len(s) > 0
}

func parseHex(s string) (uint32, error) {
	var v uint32
	for _, c := range s {
		v <<= 4
		switch {
		case c >= '0' && c <= '9':
			v |= uint32(c - '0')
		case c >= 'a' && c <= 'f':
			v |= uint32(c-'a') + 10
		}
	}
	return v, nil
}
