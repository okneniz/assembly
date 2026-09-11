package dap

import (
	"io"
	"strings"
	"sync"
)

// tailLimit bounds the kept tail: the interesting failure text of a
// `go run` (the compiler's verdict) sits in the last lines.
const tailLimit = 512

// tailWriter forwards everything to the underlying writer and keeps
// the last tailLimit bytes: when the relayed child dies before
// answering, its stderr tail rides along in the launch failure (the
// editor's adapter log is easy to miss; the error toast is not).
type tailWriter struct {
	w    io.Writer
	mu   sync.Mutex
	tail []byte
}

func newTailWriter(w io.Writer) *tailWriter {
	return &tailWriter{w: w}
}

func (t *tailWriter) Write(p []byte) (int, error) {
	t.mu.Lock()
	t.tail = append(t.tail, p...)
	if len(t.tail) > tailLimit {
		t.tail = t.tail[len(t.tail)-tailLimit:]
	}

	t.mu.Unlock()

	return t.w.Write(p)
}

// String is the kept tail, trimmed for a message line.
func (t *tailWriter) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return strings.TrimSpace(string(t.tail))
}
