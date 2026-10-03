package unit

// The core output tests: streams and labels against a fake policy, the
// integer emitters, symbolic data words, the deposit guards, and the
// phase boundary (the size lie, the resolve error that keeps the layout,
// the entry check). wordRes/stubSym stand in for the producer's records
// (a built instruction, a branch to a label).

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// wordRes is a built instruction stand-in: one little-endian word.
type wordRes struct {
	w uint32
}

func (r wordRes) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(binary.LittleEndian.AppendUint32(nil, r.w))
	return int64(n), err
}

// stubSym is a branch stand-in: a fixed size resolving to size bytes of
// one marker (err: a resolve that fails).
type stubSym struct {
	size int
	mark byte
	err  error
}

func (s stubSym) Resolve(Ctx) ([]Resolved, error) {
	if s.err != nil {
		return nil, s.err
	}

	return []Resolved{blob(bytes.Repeat([]byte{s.mark}, s.size))}, nil
}

func (s stubSym) Size() int {
	return s.size
}

// errStub - the stand-in resolve failure.
var errStub = errStubError{}

type errStubError struct{}

func (errStubError) Error() string {
	return "stub"
}

func TestUnitStreamsLayout(t *testing.T) {
	u := streamsUnit()

	// a deterministic fake policy: the addresses are arbitrary, the
	// assertions only check that everything resolved against them
	f := u.Resolve(func(text, data, dataMem int) (uint64, uint64) {
		require.Equal(t, 8+7*4, text) // the deferred pair plus seven instructions
		require.Equal(t, 8, data)
		require.Equal(t, 24, dataMem) // 8 of data + the 16-byte bss
		return 0x1000, 0x8000
	})
	require.Empty(t, f.Errs)

	code, codeErr := f.EncodeText()
	data, dataErr := f.EncodeData()
	require.NoError(t, codeErr)
	require.NoError(t, dataErr)
	require.Len(t, code, 8+7*4)
	require.Equal(t, []byte{7, 0, 0, 0, 0, 0, 0, 0}, data)
	require.Equal(t, 24, f.DataMem)

	// the deferred record leads the text stream, the instructions follow
	require.Equal(t, bytes.Repeat([]byte{0xEE}, 8), code[:8])
	require.Equal(t, uint32(1), binary.LittleEndian.Uint32(code[8:12]))

	// the symbols landed on the policy's addresses
	require.Equal(t, map[string]uint64{
		"start":   0x1000,
		"counter": 0x8000,
		"buf":     0x8008,
	}, f.Syms)

	// the line map: the text records then the data ones, one entry each
	// (labels emit nothing)
	require.Len(t, f.Lines, 10)
	require.Equal(t, uint64(0x1000), f.Lines[0].Addr)
	require.Equal(t, 8, f.Lines[0].Size)
	require.Equal(t, uint64(0x1000+8+6*4), f.Lines[7].Addr)
	require.Equal(t, uint64(0x8000), f.Lines[8].Addr)
	require.Equal(t, 8, f.Lines[8].Size)
	require.Equal(t, uint64(0x8008), f.Lines[9].Addr)
	require.Equal(t, 16, f.Lines[9].Size)
}

func TestUnitIntegerEmitters(t *testing.T) {
	f := New().
		Data().
		Half(Pos{}, 0x0102).
		Word(Pos{}, 0x03040506).
		Quad(Pos{}, 0x0708090a0b0c0d0e).
		Resolve(nopPlace)

	require.Empty(t, f.Errs)

	data, dataErr := f.EncodeData()
	require.NoError(t, dataErr)
	require.Equal(t, []byte{
		0x02, 0x01,
		0x06, 0x05, 0x04, 0x03,
		0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08, 0x07,
	}, data)
	require.Equal(t, 14, f.DataMem)
}

func TestUnitQuadSym(t *testing.T) {
	// the address initializer of a static: the label lives in the text
	// stream, the pointer word in the data one
	u := New()
	u.Label("msg")
	u.Instr(Pos{}, newWordRes(0xAAAA), nil)
	u.Data()
	u.Label("p").QuadSym(Pos{}, "msg")

	f := u.Resolve(func(text, data, dataMem int) (uint64, uint64) {
		return 0x1000, 0x8000
	})
	require.Empty(t, f.Errs)

	data, dataErr := f.EncodeData()
	require.NoError(t, dataErr)
	require.Equal(t, uint64(0x1000), binary.LittleEndian.Uint64(data))
}

func TestUnitQuadSymUndefined(t *testing.T) {
	f := New().Data().QuadSym(Pos{}, "nowhere").Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], `undefined label "nowhere"`)
}

func TestUnitGuards(t *testing.T) {
	// a bss reserve outside the data stream
	f := New().Bss(Pos{}, 8).Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], "belongs to the data stream")

	// a non-positive reserve
	f = New().Data().Bss(Pos{}, 0).Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], "not positive")

	// a refused instruction (the builder error surfaces at Resolve,
	// nothing is deposited)
	f = New().Instr(Pos{}, nil, errStub).Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.Empty(t, f.Text)

	// an undefined entry
	f = New().Entry("gone").Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], `entry: undefined label "gone"`)
}

func TestUnitResolveErrorKeepsLayout(t *testing.T) {
	// a deferred record that fails: its error is reported, its size and
	// the addresses of the records around it stay true
	f := New().
		Sym(Pos{}, newFailingSym(4)).
		Label("after").
		Instr(Pos{}, newWordRes(1), nil).
		Resolve(nopPlace)

	require.Len(t, f.Errs, 1)
	require.Equal(t, uint64(0x1004), f.Syms["after"])
	require.Len(t, f.Text, 1) // the failed record leaves no resolved bytes
	require.Equal(t, uint64(0x1004), f.Lines[1].Addr)
}

func TestUnitSizeLie(t *testing.T) {
	// a deferred record that encodes more than it reserved would desync
	// every label after it - the resolve phase rejects it
	f := New().
		Sym(Pos{}, lyingSym{declared: 4, encoded: 8}).
		Label("after").
		Resolve(nopPlace)

	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], "resolved 8 bytes, reserved 4")
}

// lyingSym declares one size and encodes another (the guard fixture).
type lyingSym struct {
	declared int
	encoded  int
}

func (s lyingSym) Resolve(Ctx) ([]Resolved, error) {
	return []Resolved{blob(make([]byte, s.encoded))}, nil
}

func (s lyingSym) Size() int {
	return s.declared
}

func TestUnitPositions(t *testing.T) {
	u := New()
	u.Instr(NewPos("t.c", 1), newWordRes(1), nil)
	u.Instr(NewPos("t.c", 2), newWordRes(2), nil)

	f := u.Resolve(nopPlace)
	require.Empty(t, f.Errs)
	require.Len(t, f.Lines, 2)
	require.Equal(t, NewPos("t.c", 1), f.Lines[0].Pos)
	require.Equal(t, NewPos("t.c", 2), f.Lines[1].Pos)
}

func TestUnitRecords(t *testing.T) {
	// the canonical deferred records: the arch formula travels as an
	// injected function, the resolve driver is the unit's alone
	branch := NewBranch("b", "target", 4, func(target, pc uint64) (Resolved, error) {
		require.Equal(t, uint64(0x100c), target) // the label sits after both records
		require.Equal(t, uint64(0x1000), pc)
		return newWordRes(uint32(target - pc)), nil
	})

	pair := NewPair("la", "target", 8, func(t, pc uint64) ([]Resolved, error) {
		return []Resolved{newWordRes(1), newWordRes(2)}, nil
	})

	u := New()
	u.Sym(NewPos("t.c", 1), branch)
	u.Sym(NewPos("t.c", 2), pair)
	u.Label("target")

	f := u.Resolve(nopPlace)
	require.Empty(t, f.Errs)

	code, codeErr := f.EncodeText()
	require.NoError(t, codeErr)
	require.Len(t, code, 12)
	require.Equal(t, uint32(12), binary.LittleEndian.Uint32(code[0:4]))

	// the undefined label and the ctor failure carry the producer's src
	f = New().Sym(Pos{}, NewBranch("b", "gone", 4, nil)).Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], `b: undefined label "gone"`)

	f = New().
		Sym(Pos{}, NewBranch("tbz", "t", 4, func(t, pc uint64) (Resolved, error) {
			return nil, errStub
		})).
		Label("t").
		Resolve(nopPlace)
	require.Len(t, f.Errs, 1)
	require.ErrorContains(t, f.Errs[0], "tbz: stub")
}

func newWordRes(w uint32) wordRes {
	return wordRes{w: w}
}

func newStubSym(size int, mark byte) stubSym {
	return stubSym{size: size, mark: mark}
}

func newFailingSym(size int) stubSym {
	return stubSym{size: size, err: errStub}
}

// streamsUnit - a program reading a data static through a branch-like
// deferred record: 8 deferred bytes plus seven instructions of text, 8
// bytes of data, a 16-byte bss tail.
func streamsUnit() *Unit {
	u := New().Entry("start")

	u.Label("start")
	u.Sym(Pos{}, newStubSym(8, 0xEE))
	for _, w := range []uint32{1, 2, 3, 4, 5, 6, 7} {
		u.Instr(Pos{}, newWordRes(w), nil)
	}

	u.Data()
	u.Label("counter").Quad(Pos{}, 7)
	u.Label("buf").Bss(Pos{}, 16)

	return u
}

// nopPlace is a policy that does not care.
func nopPlace(text, data, dataMem int) (uint64, uint64) {
	return 0x1000, 0x8000
}
