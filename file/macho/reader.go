package macho

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

// errf is a helper for creating Mach-O parse errors.
func errf(f string, args ...any) error {
	return fmt.Errorf("assembly/file/macho: "+f, args...)
}

// u32c is a combinator: 4 bytes as uint32 in the given byte order.
func u32c(order binary.ByteOrder) common.Combinator[byte, int, uint32] {
	return bytes.ReadAs[uint32](4, "macho: expected 4 bytes for u32", order)
}

// u64c is a combinator: 8 bytes as uint64 in the given byte order.
func u64c(order binary.ByteOrder) common.Combinator[byte, int, uint64] {
	return bytes.ReadAs[uint64](8, "macho: expected 8 bytes for u64", order)
}

// readU32At reads a uint32 at absolute position off.
func readU32At(buf common.Buffer[byte, int], order binary.ByteOrder, off int) (uint32, error) {
	return readAt(buf, u32c(order), off)
}

// readU64At reads a uint64 at absolute position off.
func readU64At(buf common.Buffer[byte, int], order binary.ByteOrder, off int) (uint64, error) {
	return readAt(buf, u64c(order), off)
}

// readBytes reads size bytes at absolute position off.
func readBytes(buf common.Buffer[byte, int], off, size int) ([]byte, error) {
	if err := checkSpan(off, size); err != nil {
		return nil, err
	}

	return readAt(buf, common.Count(size, "macho: raw bytes", bytes.Any()), off)
}

// checkSpan rejects negative values and off+size overflow (without this,
// parsec-Count panics on garbage offsets from corrupted files).
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

// readAt runs combinator c at absolute position off, restoring the buffer
// to its original position. Random access on top of a single streaming buffer.
func readAt[T any](
	buf common.Buffer[byte, int],
	c common.Combinator[byte, int, T],
	off int,
) (T, error) {
	if off < 0 {
		var zero T
		return zero, errf("cannot read from negative offset %d", off)
	}

	prev := buf.Position()
	if err := buf.Seek(off); err != nil {
		var zero T
		return zero, err
	}

	v, err := c(buf)
	// position restore: the seek-back error matters only if the read succeeded
	if seekErr := buf.Seek(prev); seekErr != nil && err == nil {
		var zero T
		return zero, seekErr
	}

	return v, err
}
