package dap

import (
	"bytes"
	"sync"
	"time"
)

// The console batching: a newline, a full buffer, or a short quiet
// period releases one output event. Without it the executor's raw
// stream turns the editor console into a per-byte letter column (the
// UART emits one byte per write).
const (
	consoleFlushBytes  = 512
	consoleFlushPeriod = 50 * time.Millisecond
)

// consoleWriter is the executor console drain: every write lands in
// the batch, the batch leaves as one output event (the program's own
// printing, category stdout).
type consoleWriter struct {
	srv *Server

	mu    sync.Mutex
	buf   []byte
	timer *time.Timer
}

func (c *consoleWriter) Write(p []byte) (int, error) {
	c.mu.Lock()
	c.buf = append(c.buf, p...)
	release := len(c.buf) >= consoleFlushBytes || bytes.ContainsRune(p, '\n')
	if !release && c.timer == nil {
		c.timer = time.AfterFunc(consoleFlushPeriod, c.flush)
	}

	c.mu.Unlock()

	if release {
		c.flush()
	}

	return len(p), nil
}

// flush releases the buffered console bytes as one output event.
func (c *consoleWriter) flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}

	if len(c.buf) == 0 {
		return
	}

	c.srv.notify("output", outputBody{Category: "stdout", Output: string(c.buf)})
	c.buf = nil
}
