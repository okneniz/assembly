package dap

import (
	"io"

	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/debug/session"
)

// Runtime is one launched debug run: the session bound to the halted
// target, its line map, and the machine teardown. The server owns it
// from launch to disconnect.
type Runtime struct {
	Tgt    debug.Target
	Sess   *session.Session
	Lines  []session.Line
	Closer io.Closer
}

// Launcher boots a run from the launch arguments: load the input,
// start the executor (its console drained into the writer), bind the
// session. The server calls it exactly once per launch; tests
// substitute their own.
type Launcher func(args launchArgs, console io.Writer) (*Runtime, error)
