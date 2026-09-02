package elf

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
)

// errf is a helper for creating ELF parsing errors.
func errf(f string, args ...any) error {
	return fmt.Errorf("assembly/file/elf: "+f, args...)
}

// readU16At reads a uint16 at the absolute position off (restoring the position).
func readU16At(buf parsec.Buffer[byte, int], order binary.ByteOrder, off int) (uint16, error) {
	return readAt(buf, bytes.ReadAs[uint16](2, "elf: u16", order), off)
}

// readU32At reads a uint32 at the absolute position off.
func readU32At(buf parsec.Buffer[byte, int], order binary.ByteOrder, off int) (uint32, error) {
	return readAt(buf, bytes.ReadAs[uint32](4, "elf: u32", order), off)
}

// readU64At reads a uint64 at the absolute position off.
func readU64At(buf parsec.Buffer[byte, int], order binary.ByteOrder, off int) (uint64, error) {
	return readAt(buf, bytes.ReadAs[uint64](8, "elf: u64", order), off)
}

// readByteAt reads a single byte at the absolute position off.
func readByteAt(buf parsec.Buffer[byte, int], off int) (byte, error) {
	return readAt(buf, bytes.Any(), off)
}

// readBytes reads size bytes at the absolute position off.
func readBytes(buf parsec.Buffer[byte, int], off, size int) ([]byte, error) {
	if err := checkSpan(off, size); err != nil {
		return nil, err
	}

	return readAt(buf, parsec.Count(size, "elf: raw bytes", bytes.Any()), off)
}

// checkSpan rejects negative values and overflow of off+size (without this,
// parsec-Count panics on garbage offsets from corrupted files).
func checkSpan(off, size int) error {
	if off < 0 || size < 0 || off > math.MaxInt-size {
		return errf("invalid read range: off=%d size=%d", off, size)
	}

	return nil
}

// readCStringAt reads a null-terminated string starting at position off
// (ManyTill combinator: any bytes up to the null one).
func readCStringAt(buf parsec.Buffer[byte, int], off int) (string, error) {
	if off < 0 {
		return "", errf("read string at negative offset %d", off)
	}

	prev := buf.Position()
	if err := buf.Seek(off); err != nil {
		return "", err
	}

	bs, err := parsec.ManyTill(
		16, "elf: cstring",
		bytes.Any(),
		parsec.Eq[byte, int]("elf: string terminator", 0),
	)(buf)
	// restore position: a seek-back error matters only if the read succeeded
	if seekErr := buf.Seek(prev); seekErr != nil && err == nil {
		return "", seekErr
	}

	if err != nil {
		return "", err
	}

	return string(bs), nil
}

// readAt runs the combinator c at the absolute position off, returning the
// buffer to its original position. Random access on top of a single streaming
// buffer.
func readAt[T any](
	buf parsec.Buffer[byte, int],
	c parsec.Combinator[byte, int, T],
	off int,
) (T, error) {
	if off < 0 {
		var zero T
		return zero, errf("read at negative offset %d", off)
	}

	prev := buf.Position()
	if err := buf.Seek(off); err != nil {
		var zero T
		return zero, err
	}

	v, err := c(buf)
	// restore position: a seek-back error matters only if the read succeeded
	if seekErr := buf.Seek(prev); seekErr != nil && err == nil {
		var zero T
		return zero, seekErr
	}

	return v, err
}
