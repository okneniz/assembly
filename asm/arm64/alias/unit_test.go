package alias

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/unit"
)

// flatPlace places the text at textBase and the data at dataBase.
func flatPlace(textBase, dataBase uint64) unit.Place {
	return func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return textBase, dataBase
	}
}

// TestUnitMatchesBytes is the unit-mode oracle: a source without
// external names assembles to the same bytes as the byte mode (labels,
// numeric locals, the literal pool, data words, alignment - the whole
// internal layout).
func TestUnitMatchesBytes(t *testing.T) {
	src := `
start:
  mov x0, #1
  ldr x1, =0x11223344
  ldr x2, =start
  cbz x0, start
  b done ; ret
  nop
1:
  b 1b
  .word 42
done:
  .align 8
  ret
`
	res, errs := Assemble(src, 0x1000)
	require.Empty(t, errs, "bytes mode")

	u := unit.New()
	errs = AssembleUnit(u, "t.s", src)
	require.Empty(t, errs, "unit mode")

	fixed := u.Resolve(flatPlace(0x1000, 0x80000000))
	require.Empty(t, fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)
	require.Equal(t, res.Sections[0].Data, text, "unit bytes == byte-mode bytes")

	// mov, ldr, ldr, cbz, b, ret, nop, b, .word = 36 bytes, then the
	// two 8-byte pool slots
	require.Equal(t, uint64(0x1000), fixed.Syms["start"])
	require.Equal(t, uint64(0x1024), fixed.Syms["done"])
}

// TestUnitDeferredExternals is the head.S shape: the source calls into
// the C side (bl kernel_init, ldr x4, =kernel_stack + 16) - the names
// are not in the source, they resolve against the unit when the C side
// has deposited its labels; an undefined one stays an error.
func TestUnitDeferredExternals(t *testing.T) {
	src := `
_start:
  bl kernel_init
  ldr x4, =kernel_stack + 16
  b _start
`
	u := unit.New()
	errs := AssembleUnit(u, "head.S", src)
	require.Empty(t, errs)

	// the C side deposits after the assembly: its labels join the same
	// namespace (the entry label discipline of the prog chains)
	var b arch.Builder
	u.Text()
	u.Label("kernel_init")
	u.Instr(unit.NewPos("c.c", 1), b.Nop(), nil)
	u.Instr(unit.NewPos("c.c", 2), b.Nop(), nil)
	u.Data()
	u.Bytes(unit.NewPos("c.c", 3), 0, 0, 0, 0, 0, 0, 0, 0)
	u.Label("kernel_stack")

	fixed := u.Resolve(flatPlace(0x1000, 0x80000000))
	require.Empty(t, fixed.Errs, "errs: %v", fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)

	// bl @0x1000, ldr @0x1004, b @0x1008, pool slot @0x100c (8 bytes),
	// the C nops @0x1014/0x1018 (kernel_init = 0x1014)
	require.Len(t, text, 3*4+8+2*4)
	words := []uint32{
		binary.LittleEndian.Uint32(text[0:]),
		binary.LittleEndian.Uint32(text[4:]),
		binary.LittleEndian.Uint32(text[8:]),
	}
	require.Equal(t, uint32(0x94000005), words[0], "bl kernel_init")
	require.Equal(t, uint32(0x58000044), words[1], "ldr x4, =kernel_stack+16")
	require.Equal(t, uint32(0x17FFFFFE), words[2], "b _start")

	require.Equal(t, uint64(0x80000018), binary.LittleEndian.Uint64(text[12:20]), "the pool slot value")
	require.Equal(t, uint32(0xD503201F), binary.LittleEndian.Uint32(text[20:]), "the C nop before kernel_init")
	require.Equal(t, uint64(0x1014), fixed.Syms["kernel_init"])
	require.Equal(t, uint64(0x1000), fixed.Syms["_start"])

	// without the C side the externals stay undefined at resolve
	bare := unit.New()
	require.Empty(t, AssembleUnit(bare, "head.S", src))
	bfixed := bare.Resolve(flatPlace(0x1000, 0x80000000))
	require.NotEmpty(t, bfixed.Errs)
	require.Contains(t, strings.Join(errStrings(bfixed.Errs), "; "), "kernel_init")
}

// TestUnitSections maps the sections onto the unit streams: .boot.text
// joins the text stream, .data is the data stream, .bss a zero-fill
// reserve with its labels past the file data.
func TestUnitSections(t *testing.T) {
	src := `
.section .boot.text, "ax"
boot:
  b main
.section .data
val: .word 5
.section .bss
buf: .zero 16
.text
main:
  ret
`
	u := unit.New()
	require.Empty(t, AssembleUnit(u, "t.s", src))

	fixed := u.Resolve(flatPlace(0x1000, 0x80000000))
	require.Empty(t, fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)
	require.Len(t, text, 8, "two words: b main, ret")

	// b main @0x1000 -> main @0x1004
	require.Equal(t, uint32(0x14000001), binary.LittleEndian.Uint32(text[0:]))
	require.Equal(t, uint32(0xD65F03C0), binary.LittleEndian.Uint32(text[4:]))

	data, err := fixed.EncodeData()
	require.NoError(t, err)
	require.Equal(t, []byte{5, 0, 0, 0}, data)
	require.Equal(t, 20, fixed.DataMem, "4 file bytes + the 16-byte reserve")

	require.Equal(t, uint64(0x1000), fixed.Syms["boot"])
	require.Equal(t, uint64(0x1004), fixed.Syms["main"])
	require.Equal(t, uint64(0x80000000), fixed.Syms["val"])
	require.Equal(t, uint64(0x80000004), fixed.Syms["buf"])
}

// errStrings flattens the resolve errors for matching.
func errStrings(errs []error) []string {
	out := make([]string, 0, len(errs))
	for _, e := range errs {
		out = append(out, e.Error())
	}

	return out
}
