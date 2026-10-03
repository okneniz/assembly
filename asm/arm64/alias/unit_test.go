package alias

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/unit"
)

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

	fixed := u.Resolve(flatPlace(0x1000))
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

	fixed := u.Resolve(flatPlace(0x1000))
	require.Empty(t, fixed.Errs, "errs: %v", fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)

	// bl @0x1000, ldr @0x1004, b @0x1008, then the pool: the x-slot sits
	// at its natural 8-byte alignment (the 12-byte tail is 4 mod 8 - the
	// zero pad word is the udf #0 of the gas dump), the C nops
	// @0x1018/0x101C (kernel_init = 0x1018)
	require.Len(t, text, 3*4+4+8+2*4)
	words := []uint32{
		binary.LittleEndian.Uint32(text[0:]),
		binary.LittleEndian.Uint32(text[4:]),
		binary.LittleEndian.Uint32(text[8:]),
	}
	require.Equal(t, uint32(0x94000006), words[0], "bl kernel_init")
	require.Equal(t, uint32(0x58000064), words[1], "ldr x4, =kernel_stack+16")
	require.Equal(t, uint32(0x17FFFFFE), words[2], "b _start")

	require.Equal(t, []byte{0, 0, 0, 0}, text[12:16], "the alignment pad word")
	require.Equal(
		t,
		uint64(0x80000018),
		binary.LittleEndian.Uint64(text[16:24]),
		"the pool slot value",
	)
	require.Equal(
		t,
		uint32(0xD503201F),
		binary.LittleEndian.Uint32(text[24:]),
		"the C nop before kernel_init",
	)
	require.Equal(t, uint64(0x1018), fixed.Syms["kernel_init"])
	require.Equal(t, uint64(0x1000), fixed.Syms["_start"])

	// without the C side the externals stay undefined at resolve
	bare := unit.New()
	require.Empty(t, AssembleUnit(bare, "head.S", src))
	bfixed := bare.Resolve(flatPlace(0x1000))
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

	fixed := u.Resolve(flatPlace(0x1000))
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

// TestUnitPoolAlignment - the literal pool slots sit at their natural
// alignment in the unit mode too (the head.S shape of the first kernel
// image: the quad ldr= slot after a tail at 4 mod 8 was the alignment
// fault); the unit bytes stay identical to the byte mode.
func TestUnitPoolAlignment(t *testing.T) {
	src := `
  nop
  nop
  ldr x19, =0x1122334455667788
`
	res, errs := Assemble(src, 0x41000000)
	require.Empty(t, errs, "bytes mode")

	u := unit.New()
	errs = AssembleUnit(u, "head.S", src)
	require.Empty(t, errs, "unit mode")

	fixed := u.Resolve(flatPlace(0x41000000))
	require.Empty(t, fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)
	require.Equal(t, res.Sections[0].Data, text, "unit bytes == byte-mode bytes")

	require.Len(t, text, 12+4+8, "code + the pad word + the slot")
	require.Equal(t, uint32(0xD503201F), binary.LittleEndian.Uint32(text[0:]))
	require.Equal(t, uint32(0xD503201F), binary.LittleEndian.Uint32(text[4:]))
	require.Equal(
		t,
		uint32(0x58000053),
		binary.LittleEndian.Uint32(text[8:]),
		"ldr x19 -> the slot @+8",
	)
	require.Equal(t, []byte{0, 0, 0, 0}, text[12:16], "the udf #0 pad word")
	require.Equal(
		t,
		uint64(0x1122334455667788),
		binary.LittleEndian.Uint64(text[16:]),
		"the slot (LE64)",
	)
}

// TestUnitSysOps - the machine.h system-operation idiom in the unit
// mode: the same bytes as the byte mode (the fragment oracle shape,
// through the .S deposit path dcc links with).
func TestUnitSysOps(t *testing.T) {
	src := `vspace_switch:
  dsb sy
  ic iallu
  dsb sy
  isb
  tlbi vmalls12e1
  dsb sy
  isb
  dc zva, x0
  ic ivau, x1
  msr sp_el0, x2
  eret
`

	res, errs := Assemble(src, 0x1000)
	require.Empty(t, errs, "bytes mode")

	u := unit.New()
	errs = AssembleUnit(u, "t.s", src)
	require.Empty(t, errs, "unit mode")

	fixed := u.Resolve(flatPlace(0x1000))
	require.Empty(t, fixed.Errs)

	text, err := fixed.EncodeText()
	require.NoError(t, err)
	require.Equal(t, res.Sections[0].Data, text, "unit bytes == byte-mode bytes")
	require.Len(t, text, 44, "eleven words")
}

// flatPlace places the text at textBase and the data at the 2GB mark.
func flatPlace(textBase uint64) unit.Place {
	return func(textSize, dataSize, dataMem int) (uint64, uint64) {
		return textBase, 0x80000000
	}
}

// errStrings flattens the errors (resolve or assemble) for matching.
func errStrings[E interface{ Error() string }](errs []E) []string {
	out := make([]string, 0, len(errs))
	for _, e := range errs {
		out = append(out, e.Error())
	}

	return out
}
