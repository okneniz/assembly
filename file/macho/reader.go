package macho

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/okneniz/parsec"
)

// errf is a helper for creating Mach-O parse errors.
func errf(f string, args ...any) error {
	return fmt.Errorf("assembly/file/macho: "+f, args...)
}

// The read helpers are plain byte loops over the shared parsec buffer:
// a fixed-offset read has no alternatives to backtrack, so combinators
// add nothing (and every combinator call builds fresh closures - the
// thing the style rules forbid); random access stays seek + consume +
// restore on the single streaming buffer.

// readU32At reads a uint32 at absolute position off.
func readU32At(buf parsec.Buffer[byte, int], order binary.ByteOrder, off int) (uint32, error) {
	bs, err := readRawAt(buf, off, 4)
	if err != nil {
		return 0, err
	}

	return order.Uint32(bs), nil
}

// readU64At reads a uint64 at absolute position off.
func readU64At(buf parsec.Buffer[byte, int], order binary.ByteOrder, off int) (uint64, error) {
	bs, err := readRawAt(buf, off, 8)
	if err != nil {
		return 0, err
	}

	return order.Uint64(bs), nil
}

// readBytes reads size bytes at absolute position off.
func readBytes(buf parsec.Buffer[byte, int], off, size int) ([]byte, error) {
	return readRawAt(buf, off, size)
}

// checkSpan rejects negative values and off+size overflow.
func checkSpan(off, size int) error {
	if off < 0 || size < 0 || off > math.MaxInt-size {
		return errf("invalid read span: off=%d size=%d", off, size)
	}

	return nil
}

// cstr returns a string from bytes b up to the first null byte.
func cstr(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}

	return string(b)
}

// readRawAt reads n bytes at the absolute position off, returning the
// buffer to its original position. Random access on top of a single
// streaming buffer.
func readRawAt(buf parsec.Buffer[byte, int], off, n int) ([]byte, error) {
	if err := checkSpan(off, n); err != nil {
		return nil, err
	}

	prev := buf.Position()
	if err := buf.Seek(off); err != nil {
		return nil, err
	}

	// a corrupt header may name an absurd size - the capacity is
	// capped, a truncated read fails at EOF instead of at allocation
	out := make([]byte, 0, min(n, 4096))
	var readErr error
	for range n {
		b, err := buf.Read(true)
		if err != nil {
			readErr = errf("read %d bytes at %d: %w", n, off, err)
			break
		}

		out = append(out, b)
	}

	// restore position: a seek-back error matters only if the read succeeded
	if seekErr := buf.Seek(prev); seekErr != nil && readErr == nil {
		return nil, seekErr
	}

	if readErr != nil {
		return nil, readErr
	}

	return out, nil
}
