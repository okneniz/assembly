package unit

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
)

// Fixed - the resolved output: the two streams as purely encodable
// records in deposit order (labels resolved away into Syms), the data
// memory size (>= the encoded data: the bss tail), the line map, and
// every error collected on the way (deposit errors and resolve errors
// together). Encode materializes the bytes; nothing here needs an
// environment.
type Fixed struct {
	Text    []Resolved
	Data    []Resolved
	DataMem int
	Syms    map[string]uint64
	Lines   []LineEntry
	Errs    []error
}

// EncodeData materializes the data stream - the file bytes; the bss tail
// is the DataMem beyond them (memory the kernel zeroes, no records).
func (f *Fixed) EncodeData() ([]byte, error) {
	var buf bytes.Buffer
	var errs []error
	for _, r := range f.Data {
		if _, err := r.Encode(&buf); err != nil {
			errs = append(errs, err)
		}
	}

	return buf.Bytes(), errors.Join(errs...)
}

// EncodeText materializes the text stream. The records were verified at
// resolve (each deferred record encodes to exactly its reserved size),
// so an error here can only be a record's own late encode failure.
func (f *Fixed) EncodeText() ([]byte, error) {
	var buf bytes.Buffer
	var errs []error
	for _, r := range f.Text {
		if _, err := r.Encode(&buf); err != nil {
			errs = append(errs, err)
		}
	}

	return buf.Bytes(), errors.Join(errs...)
}

// Resolve is the phase boundary: the policy places the two streams (see
// Place), every label becomes base+offset, every deferred record resolves
// through a Ctx at its final address and lands in the output as purely
// encodable records. A record that fails (an undefined label, an
// out-of-range offset) leaves its error and its reserved size in place -
// the addresses of the records around it stay true.
func (u *Unit) Resolve(place Place) *Fixed {
	f := &Fixed{
		DataMem: u.mem[1],
		Syms:    make(map[string]uint64, len(u.labels)),
		Errs:    slices.Clone(u.errs),
	}

	textAddr, dataAddr := place(u.mem[0], u.file[1], u.mem[1])
	base := [2]uint64{textAddr, dataAddr}

	for _, l := range u.labels {
		f.Syms[l.name] = base[l.stream] + uint64(l.off)
	}

	resolve := func(name string) (uint64, bool) {
		v, ok := f.Syms[name]
		return v, ok
	}

	// pass 2: resolve per stream; the line map covers both.
	var cur [2]int
	streams := [2]*[]Resolved{&f.Text, &f.Data}

	for _, s := range u.slots {
		addr := base[s.stream] + uint64(cur[s.stream])
		f.Lines = append(f.Lines, NewLineEntry(addr, s.size, s.pos))

		switch {
		case s.nobits > 0:
			cur[1] += s.nobits
		case s.res != nil:
			*streams[s.stream] = append(*streams[s.stream], s.res)
			cur[s.stream] += s.size
		case s.sym != nil:
			rs, err := evalSym(s.sym, addr, resolve)
			if err == nil {
				*streams[s.stream] = append(*streams[s.stream], rs...)
			} else {
				f.Errs = append(f.Errs, err)
			}

			cur[s.stream] += s.size
		}
	}

	if u.entry != "" {
		if _, ok := f.Syms[u.entry]; !ok {
			f.Errs = append(f.Errs, fmt.Errorf("entry: undefined label %q", u.entry))
		}
	}

	return f
}

// evalSym - one deferred record at its final address (the record's own
// Resolve, not this unit's phase transition), with the guard that it
// encodes to exactly its reserved size - the layout counted that many
// bytes, and a mismatch would desync every label after it.
func evalSym(
	s Sym,
	addr uint64,
	resolve func(string) (uint64, bool),
) ([]Resolved, error) {
	rs, err := s.Resolve(NewCtx(addr, resolve))
	if err != nil {
		return nil, err
	}

	n, err := encodedSize(rs)
	if err != nil {
		return nil, err
	}

	if n != s.Size() {
		return nil, fmt.Errorf("resolved %d bytes, reserved %d", n, s.Size())
	}

	return rs, nil
}
