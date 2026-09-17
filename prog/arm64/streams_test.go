package arm64

// The data stream: Text/Data switching, the integer emitters, bss
// reserves, La pairs resolving across streams, and the layout-mode
// assembly (AssembleLayout). The flat Assemble keeps its behavior and
// rejects programs that switched streams.

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// streamsProg - a program reading and writing a data static and a bss
// one through an La pair, exiting with the data value it re-read.
func streamsProg() *Program {
	p := New().Entry("start")

	p.Label("start")
	p.La(X0, "counter")
	p.Ldr(X1, X0, 0)
	p.AddImm(X1, X1, 1, arch.NoSh12)
	p.Str(X1, X0, 0)
	p.Ldr(X0, X0, 0)
	p.Movz(X16, 0x200, arch.Hw1)
	p.Movk(X16, 1, arch.Hw0)
	p.Svc(0x80)

	p.Data()
	p.Label("counter").Quad(7)
	p.Label("buf").Bss(16)

	return p
}

func TestStreamsLayout(t *testing.T) {
	bin, errs := streamsProg().Build()
	require.Empty(t, errs)

	// a deterministic fake policy: the addresses are arbitrary, the
	// assertions only check that everything resolved against them
	res := bin.AssembleLayout(func(text, data, dataMem int) (uint64, uint64) {
		require.Equal(t, 8+7*4, text) // the La pair plus seven instructions
		require.Equal(t, 8, data)
		require.Equal(t, 24, dataMem) // 8 of data + the 16-byte bss
		return 0x1000, 0x8000
	})
	require.Empty(t, res.Errs)

	require.Len(t, res.Code, 8+7*4)
	require.Equal(t, []byte{7, 0, 0, 0, 0, 0, 0, 0}, res.Data)
	require.Equal(t, 24, res.DataMem)

	// the symbols landed on the policy's addresses
	require.Equal(t, map[string]uint64{
		"start":   0x1000,
		"counter": 0x8000,
		"buf":     0x8008,
	}, res.Syms)

	// the La pair: adrp x0, #7 pages (0x8000 - 0x1000); add x0, x0, #0
	require.Equal(t, uint32(0xF0000020), binary.LittleEndian.Uint32(res.Code[0:4]))
	require.Equal(t, uint32(0x91000000), binary.LittleEndian.Uint32(res.Code[4:8]))

	// the line map: the text lines then the data lines, one entry each
	require.Len(t, res.Lines, 10)
	require.Equal(t, uint64(0x1000), res.Lines[0].Addr)
	require.Equal(t, uint64(0x8000), res.Lines[8].Addr)
	require.Equal(t, 8, res.Lines[8].Size)
	require.Equal(t, uint64(0x8008), res.Lines[9].Addr)
	require.Equal(t, 16, res.Lines[9].Size)
}

func TestStreamsFlatRejected(t *testing.T) {
	bin, errs := streamsProg().Build()
	require.Empty(t, errs)

	res := bin.Assemble(0)
	require.Len(t, res.Errs, 1)
	require.ErrorContains(t, res.Errs[0], "needs AssembleLayout")
}

func TestStreamsDirectivesGuard(t *testing.T) {
	// an instruction in the data stream
	_, errs := New().Data().Nop().Build()
	require.Len(t, errs, 1)
	require.ErrorContains(t, errs[0], "instruction in the data stream")

	// a bss reserve outside the data stream
	_, errs = New().Bss(8).Build()
	require.Len(t, errs, 1)
	require.ErrorContains(t, errs[0], "belongs to the data stream")

	// a non-positive reserve
	_, errs = New().Data().Bss(0).Build()
	require.Len(t, errs, 1)
	require.ErrorContains(t, errs[0], "not positive")
}

func TestStreamsIntegerEmitters(t *testing.T) {
	bin, errs := New().
		Data().
		Half(0x0102).
		Word(0x03040506).
		Quad(0x0708090a0b0c0d0e).
		Build()
	require.Empty(t, errs)

	res := bin.AssembleLayout(nopPlace)
	require.Empty(t, res.Errs)
	require.Equal(t, []byte{
		0x02, 0x01,
		0x06, 0x05, 0x04, 0x03,
		0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07,
	}, res.Data)
	require.Equal(t, 14, res.DataMem)
}

func TestStreamsLaFlat(t *testing.T) {
	// the La pair also assembles flat, against the flat base
	bin, errs := New().
		Label("start").
		La(X0, "msg").
		Label("msg").
		Ascii("hi").
		Build()
	require.Empty(t, errs)

	res := bin.Assemble(0x1000)
	require.Empty(t, res.Errs)
	require.Len(t, res.Code, 8+2)

	// adrp x0, #0 (same page); add x0, x0, #8 (the label right after)
	require.Equal(t, uint32(0x90000000), binary.LittleEndian.Uint32(res.Code[0:4]))
	require.Equal(t, uint32(0x91002000), binary.LittleEndian.Uint32(res.Code[4:8]))
}

// nopPlace - a policy that does not care.
func nopPlace(text, data, dataMem int) (uint64, uint64) {
	return 0x1000, 0x8000
}

// compile-time: the policy type is the plain function the layers share.
var _ prog.Place = nopPlace
