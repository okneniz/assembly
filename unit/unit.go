// Package unit - the instruction output: the ordered stream of records a
// program is deposited into, in two typed phases. Unit is the unresolved
// output: producers (a compiler through arch builders, the asm text
// parser through fragments) deposit ready records and deferred (Sym)
// ones; nothing is resolved and nothing is encoded yet. Resolve is the
// phase boundary - the policy places the streams, every Sym resolves
// through a Ctx, and the result is Fixed, the resolved output of purely
// encodable records. The resolve vocabulary (Ctx, Resolved) lives here
// for every producer to share; the text assembler speaks it too.
package unit

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Unit - a program being deposited: an ordered sequence of records
// (built instructions, deferred ones, data, bss reserves) in two streams
// (text, data), with labels recorded at their stream positions. Nothing
// is encoded until Resolve. Every deposit carries its source position
// explicitly - the producer always knows where a record came from (a
// compiler its source node, a chain its caller); there is no resolver to
// guess it.
type Unit struct {
	slots  []slot
	labels []labelAt
	entry  string
	cur    int    // the deposit stream: 0 text (default), 1 data
	mem    [2]int // the running memory size per stream (the bss tail included)
	file   [2]int // the running file size per stream (bss excluded)
	errs   []error
}

// slot is one deposited record: a ready one (res), a deferred one (sym),
// or a bss reserve (nobits) - exactly one of the three - in its stream,
// at its source position, reserving size bytes.
type slot struct {
	res    Resolved
	sym    Sym
	nobits int
	stream int
	size   int
	pos    Pos
}

// labelAt is a label definition: its stream and its memory offset inside
// it (bss reserves counted).
type labelAt struct {
	name   string
	stream int
	off    int
}

// New - an empty output.
func New() *Unit {
	return &Unit{}
}

// Entry - the label the program starts at.
func (u *Unit) Entry(name string) *Unit {
	u.entry = name
	return u
}

// Text - deposit into the text stream: the instructions and any read-only
// data placed before the first Data() call live here.
func (u *Unit) Text() *Unit {
	u.cur = 0
	return u
}

// Data - deposit into the data stream: the writable statics of the
// program.
func (u *Unit) Data() *Unit {
	u.cur = 1
	return u
}

// Label - define a label at the current position of the current stream.
// A redefinition moves the label (the later position wins).
func (u *Unit) Label(name string) *Unit {
	u.labels = append(u.labels, labelAt{name: name, stream: u.cur, off: u.mem[u.cur]})
	return u
}

// Instr - deposit a built instruction (a Builder product: arch
// instructions are resolved records by shape). err is the construction
// error, when the builder refused the operands: it is recorded and
// surfaces at Resolve - nothing is deposited.
func (u *Unit) Instr(at Pos, r Resolved, err error) *Unit {
	if err != nil {
		return u.fail(err)
	}

	return u.ready(at, r)
}

// Sym - deposit a deferred record of the current stream: its bytes
// appear at Resolve, when its address and the program symbols are known.
// This is the add-instruction of the unresolved side - everything with a
// hole (a branch to a label, an address pair, an asm fragment, a
// symbolic data word) enters the stream through here.
func (u *Unit) Sym(at Pos, s Sym) *Unit {
	u.slots = append(u.slots, slot{sym: s, stream: u.cur, size: s.Size(), pos: at})
	u.mem[u.cur] += s.Size()
	u.file[u.cur] += s.Size()
	return u
}

// fail records a producer-side error (a refused instruction, a guard);
// it surfaces in the errors of Resolve, like every other one.
func (u *Unit) fail(err error) *Unit {
	u.errs = append(u.errs, err)
	return u
}

// Ascii - string data appended verbatim (no terminating zero), into the
// current stream.
func (u *Unit) Ascii(at Pos, s string) *Unit {
	return u.ready(at, blob(s))
}

// Bytes - raw data bytes, into the current stream.
func (u *Unit) Bytes(at Pos, b ...byte) *Unit {
	return u.ready(at, blob(b))
}

// Half - 16-bit little-endian values, into the current stream.
func (u *Unit) Half(at Pos, vs ...uint16) *Unit {
	b := make([]byte, 2*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint16(b[2*i:], v)
	}

	return u.ready(at, blob(b))
}

// Word - 32-bit little-endian values, into the current stream.
func (u *Unit) Word(at Pos, vs ...uint32) *Unit {
	b := make([]byte, 4*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint32(b[4*i:], v)
	}

	return u.ready(at, blob(b))
}

// Quad - 64-bit little-endian values, into the current stream.
func (u *Unit) Quad(at Pos, vs ...uint64) *Unit {
	b := make([]byte, 8*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint64(b[8*i:], v)
	}

	return u.ready(at, blob(b))
}

// QuadSym - 64-bit little-endian addresses of labels, into the current
// stream: each label's final image address as one deferred record - the
// address initializer of a static (a Quad whose value only the layout
// knows).
func (u *Unit) QuadSym(at Pos, labels ...string) *Unit {
	for _, name := range labels {
		u.Sym(at, newQuadSym(name))
	}

	return u
}

// Bss - a zero-fill reserve of the data stream: memory the kernel zeroes,
// no file bytes. Labels after it resolve past the reserved range.
func (u *Unit) Bss(at Pos, reserve int) *Unit {
	if u.cur != 1 {
		return u.fail(errors.New(
			"a bss reserve belongs to the data stream",
		))
	}

	if reserve <= 0 {
		return u.fail(fmt.Errorf("reserve %d is not positive", reserve))
	}

	u.slots = append(u.slots, slot{nobits: reserve, stream: 1, size: reserve, pos: at})
	u.mem[1] += reserve
	return u
}

// ready - deposit a resolved record: an instruction of the text stream,
// a data blob of either - anything purely encodable.
func (u *Unit) ready(at Pos, r Resolved) *Unit {
	size, err := encodedSize([]Resolved{r})
	if err != nil {
		return u.fail(err)
	}

	u.slots = append(u.slots, slot{res: r, stream: u.cur, size: size, pos: at})
	u.mem[u.cur] += size
	u.file[u.cur] += size
	return u
}
