package qemu

import "io"

// Options are the executor knobs: the console drain target (nil -
// discard), the gdbstub port (0 - a fresh ephemeral port), and the
// deterministic virtual clock (-icount shift=auto: repeatable
// instruction-count-driven time).
type Options struct {
	Serial io.Writer
	Port   int
	Icount bool
}

func NewOptions() Options {
	return Options{}
}
