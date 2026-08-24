package assembly_test

// Property tests (oh-snap) - common part of the suite: a deterministic seed
// (ASSEMBLY_SEED, default 42) and a minimal ELF writer that both arch
// differentials use to feed arbitrary words to the real objdump.
// Arch-specific parts are property_a64_test.go and property_rv_test.go; the
// a64/rv suffixes (not arm64/riscv!) are intentional: the _<GOARCH> suffix of
// the file name is a build constraint.

import (
	"encoding/binary"
	"fmt"
	mrnd "math/rand/v2"
	"os"
	"strconv"
	"testing"

	"github.com/okneniz/assembly/arb"
)

// --- contextual Encoder/Decoder and the round-trip law ---------------------
//
// These types live only in tests: the instructions themselves do not need to
// know about them. Key invariant: both the encoder and the decoder take a
// context, and the law is proved in one and the same fixed context - as long
// as the context does not change, the round trip changes nothing.

// Encoder - conversion of instruction T to representation S in context C
// (bytes/text/...); ok=false means outside the property domain (conversion error).
type Encoder[C, T, S any] func(C, T) (S, bool)

// Decoder - the reverse conversion of representation S back to instruction T
// in the same context C.
type Decoder[C, S, T any] func(C, S) (T, bool)

// RoundTrip - the law enc∘dec∘enc == enc in the fixed context ctx:
// encoding is stable after a round trip. Returns a predicate for ohsnap.Check.
func RoundTrip[C, T, S any](
	ctx C,
	enc Encoder[C, T, S],
	dec Decoder[C, S, T],
	eq func(S, S) bool,
) func(T) bool {
	return func(x T) bool {
		s, ok := enc(ctx, x)
		if !ok {
			return false
		}

		back, ok := dec(ctx, s)
		if !ok {
			return false
		}

		again, ok := enc(ctx, back)
		if !ok {
			return false
		}

		return eq(again, s)
	}
}

const propAddr = 0x1000

// seedRnd - deterministic generator: seed from ASSEMBLY_SEED (default 42),
// logged so a failure can be reproduced.
func seedRnd(t *testing.T) *mrnd.Rand {
	t.Helper()
	seed := uint64(42)
	if s := os.Getenv("ASSEMBLY_SEED"); s != "" {
		if v, err := strconv.ParseUint(s, 0, 64); err == nil {
			seed = v
		}
	}

	t.Logf("seed: %d (ASSEMBLY_SEED)", seed)
	return arb.Rnd(seed)
}

// writeObjELF - minimal ELF64 LE (e_machine/e_flags parameterized):
// header, .text @ 0x1000 with the code, .shstrtab, section table. Needed to
// feed arbitrary words to the real objdump.
func writeObjELF(t *testing.T, code []byte, machine uint16, flags uint32) (string, error) {
	t.Helper()

	const (
		ehdrSize = 64
		shdrSize = 64
		textAddr = propAddr
	)
	shstr := []byte("\x00.text\x00.shstrtab\x00")
	textName, strName := uint32(1), uint32(7)

	// Layout: ehdr | code | shstrtab (aligned to 8) | 3 shdr.
	shstrOff := (ehdrSize + len(code) + 7) &^ 7
	shoff := shstrOff + len(shstr)

	ehdr := make([]byte, ehdrSize)
	ehdr[0], ehdr[1], ehdr[2], ehdr[3] = 0x7f, 'E', 'L', 'F'
	ehdr[4] = 2              // ELFCLASS64
	ehdr[5] = 1              // ELFDATA2LSB
	ehdr[6] = 1              // EV_CURRENT
	put16(ehdr, 16, 2)       // e_type = ET_EXEC
	put16(ehdr, 18, machine) // e_machine
	put32(ehdr, 20, 1)       // e_version
	put64(ehdr, 24, textAddr)
	put64(ehdr, 32, 0) // e_phoff
	put64(ehdr, 40, uint64(shoff))
	put32(ehdr, 48, flags) // e_flags
	put16(ehdr, 52, ehdrSize)
	put16(ehdr, 54, 56) // e_phentsize
	put16(ehdr, 56, 0)  // e_phnum
	put16(ehdr, 58, shdrSize)
	put16(ehdr, 60, 3) // e_shnum
	put16(ehdr, 62, 2) // e_shstrndx

	shdrs := make([]byte, 3*shdrSize)
	// [1] .text (addralign=4, entsize=0)
	sh := shdrs[shdrSize:]
	put32(sh, 0, textName)
	put32(sh, 4, 1) // SHT_PROGBITS
	put64(sh, 8, 0x6)
	put64(sh, 16, textAddr)
	put64(sh, 24, ehdrSize)
	put64(sh, 32, uint64(len(code)))
	put64(sh, 48, 4)
	// [2] .shstrtab (addralign=1)
	sh = shdrs[2*shdrSize:]
	put32(sh, 0, strName)
	put32(sh, 4, 3) // SHT_STRTAB
	put64(sh, 24, uint64(shstrOff))
	put64(sh, 32, uint64(len(shstr)))
	put64(sh, 48, 1)

	// Final size: ehdr+code+padding fill exactly up to shstrOff,
	// followed by the shstrtab and the section table.
	out := make([]byte, 0, shoff+len(shdrs))
	out = append(out, ehdr...)
	out = append(out, code...)
	out = append(out, make([]byte, shstrOff-(ehdrSize+len(code)))...)
	out = append(out, shstr...)
	out = append(out, shdrs...)

	path := t.TempDir() + "/prop.elf"
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return "", fmt.Errorf("write ELF: %w", err)
	}

	return path, nil
}

func put16(b []byte, off int, v uint16) {
	binary.LittleEndian.PutUint16(b[off:], v)
}

func put32(b []byte, off int, v uint32) {
	binary.LittleEndian.PutUint32(b[off:], v)
}

func put64(b []byte, off int, v uint64) {
	binary.LittleEndian.PutUint64(b[off:], v)
}
