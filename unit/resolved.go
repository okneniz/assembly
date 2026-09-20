package unit

import "io"

// Resolved is a purely encodable record: an instruction without holes
// (arch instructions satisfy it structurally - their Encode is the same
// shape) or a byte blob. Encoding needs no environment; everything the
// record depends on is already inside it.
type Resolved interface {
	// Encode writes the encoding bytes.
	Encode(w io.Writer) (int64, error)
}

// blob is ready bytes as a resolved record: data (a string, an integer
// word) writes itself verbatim.
type blob []byte

func (b blob) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(b)
	return int64(n), err
}

// countingWriter counts written bytes (sizing without materializing).
type countingWriter struct {
	n int64
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.n += int64(len(p))
	return len(p), nil
}

// encodedSize is the byte count of the records under their pure encode.
func encodedSize(rs []Resolved) (int, error) {
	var c countingWriter
	for _, r := range rs {
		if _, err := r.Encode(&c); err != nil {
			return 0, err
		}
	}

	return int(c.n), nil
}
