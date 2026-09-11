package dap

import "io"

// firstWriteWatcher notes the first write through it: whether the
// relayed child ever spoke on the protocol channel (a launch response
// counts; a child that dies silently never writes).
type firstWriteWatcher struct {
	w    io.Writer
	seen *bool
}

func (f *firstWriteWatcher) Write(p []byte) (int, error) {
	*f.seen = true
	return f.w.Write(p)
}
