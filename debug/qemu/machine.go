package qemu

import (
	"errors"
	"net"
	"os"
	"os/exec"
)

// Machine is one booted executor: the process, its gdbstub connection,
// the temp image file, and the console pipe behind both. Close tears
// the whole thing down; the machine never outlives the debugging
// session.
type Machine struct {
	cmd     *exec.Cmd
	conn    net.Conn
	imgPath string
	console *os.File // the write end of the console pipe (qemu's side)
}

// Conn is the gdbstub connection: the transport of the rsp session.
func (m *Machine) Conn() net.Conn {
	return m.conn
}

// Close tears the machine down and removes its temp image. A second
// Close is an error - the session owns the lifecycle. A machine that
// powered itself off closed the connection from the far side and left
// the process done: both of those teardown steps read as no-ops here,
// not as failures.
func (m *Machine) Close() error {
	errs := []error{}

	if err := m.conn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		errs = append(errs, err)
	}

	if err := m.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		errs = append(errs, err)
	}

	// a killed qemu exits non-zero - only the unexpected failures of
	// the teardown itself are errors
	if err := m.cmd.Wait(); err != nil && !isKillExit(err) {
		errs = append(errs, err)
	}

	// the process is gone: the console write end closes, the drain
	// goroutine ends with it
	if err := m.console.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		errs = append(errs, err)
	}

	if err := os.Remove(m.imgPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// isKillExit reports whether the Wait error is the exit of a process
// killed by us (a signal exit or the already-finished marker).
func isKillExit(err error) bool {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return exitErr.ExitCode() == -1 // terminated by a signal
	}

	return errors.Is(err, os.ErrProcessDone)
}
