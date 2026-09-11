// Package qemu is the debugger's executor: a qemu-system process per
// debug run, booted frozen (-S) with its gdbstub (-gdb tcp::PORT) as
// the RSP endpoint. The package is arch-blind - the machine, cpu, and
// image placement come from the debug.Target of the arch.
package qemu

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/okneniz/assembly/debug"
)

// dialTimeout bounds the wait for the stub to come up: the process
// starts, maps the image, and listens well within it.
const dialTimeout = 10 * time.Second

// Start boots the image under the arch's qemu: frozen, with the gdbstub
// listening. The image bytes go to a temp file (the loader devices
// want a path); the console output is drained into opts.Serial so the
// pipe never fills up and blocks the machine. ctx bounds the executor
// lifetime (a cancel tears the process down with the connection).
func Start(ctx context.Context, tgt debug.Target, img []byte, opts Options) (*Machine, error) {
	if len(img) == 0 {
		return nil, errors.New("assembly/qemu: no image")
	}

	f, err := os.CreateTemp("", "assembly-debug-*.img")
	if err != nil {
		return nil, fmt.Errorf("assembly/qemu: image file: %w", err)
	}

	imgPath := f.Name()
	_, werr := f.Write(img)
	if cerr := f.Close(); werr != nil || cerr != nil {
		return nil, errors.Join(
			fmt.Errorf("assembly/qemu: image file: %w", errors.Join(werr, cerr)),
			os.Remove(imgPath),
		)
	}

	port := opts.Port
	if port == 0 {
		port, err = freePort(ctx)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("assembly/qemu: port: %w", err), os.Remove(imgPath))
		}
	}

	bin, err := findQemu(tgt.QemuBinary(), qemuDirs())
	if err != nil {
		return nil, errors.Join(fmt.Errorf("assembly/qemu: %w", err), os.Remove(imgPath))
	}

	cmd := exec.CommandContext(
		ctx,
		bin,
		fullArgs(tgt.QemuArgs(imgPath), port, opts.Icount)...,
	)

	// The console rides OUR pipe, never a caller fd: a tty on qemu's
	// stdio makes it switch the shared terminal to non-blocking, and
	// the tty is shared with our own stdin - an interactive REPL died
	// at its first idle read with EAGAIN. The drain goroutine carries
	// the serial bytes to the caller's sink.
	consoleR, consoleW, err := os.Pipe()
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("assembly/qemu: console pipe: %w", err),
			os.Remove(imgPath),
		)
	}

	cmd.Stdout = consoleW
	cmd.Stderr = consoleW
	if err := cmd.Start(); err != nil {
		return nil, errors.Join(
			fmt.Errorf("assembly/qemu: start: %w", err),
			consoleR.Close(),
			consoleW.Close(),
			os.Remove(imgPath),
		)
	}

	go drainConsole(consoleR, opts.Serial)

	conn, err := dialStub(ctx, port)
	if err != nil {
		return nil, errors.Join(
			err,
			cmd.Process.Kill(),
			cmd.Wait(),
			consoleR.Close(),
			consoleW.Close(),
			os.Remove(imgPath),
		)
	}

	return &Machine{
		cmd:     cmd,
		conn:    conn,
		imgPath: imgPath,
		console: consoleW,
	}, nil
}

// drainConsole copies the executor console into the sink until the
// machine side closes the pipe; the read end goes with the copy. A nil
// sink discards (the pipe stays drained either way).
func drainConsole(r *os.File, sink io.Writer) {
	defer func() {
		// a close error of our single-use read end: the machine
		// teardown already finished the drain
		_ = r.Close() //nolint:errcheck // the drain outcome is not actionable
	}()

	dst := sink
	if dst == nil {
		dst = io.Discard
	}

	// a copy error is the sink closing before the machine did: the
	// console had nowhere to go, the drain is over
	_, _ = io.Copy(dst, r) //nolint:errcheck // the drain outcome is not actionable
}

// findQemu resolves the executor binary: PATH first, then the given
// extra directories. The fallback exists for GUI-spawned consumers (the
// VSCode adapter): a dock-launched editor hands its children the minimal
// launchd PATH, where the homebrew prefixes are absent.
func findQemu(bin string, dirs []string) (string, error) {
	if path, err := exec.LookPath(bin); err == nil {
		return path, nil
	}

	for _, dir := range dirs {
		at := filepath.Join(dir, bin)
		if _, err := exec.LookPath(at); err == nil {
			return at, nil
		}
	}

	return "", fmt.Errorf("exec: %q: executable file not found in PATH or %v", bin, dirs)
}

// qemuDirs are the executor prefixes outside the minimal launchd PATH:
// homebrew on Apple Silicon and on Intel.
func qemuDirs() []string {
	return []string{"/opt/homebrew/bin", "/usr/local/bin"}
}

// fullArgs is the complete executor command line: the arch part plus
// the frozen boot with the gdbstub (and the deterministic clock). The
// serial is wired explicitly: -nographic would mux the monitor onto
// the same stream, and only serial bytes belong on the console.
func fullArgs(archArgs []string, port int, icount bool) []string {
	args := append([]string{}, archArgs...)
	args = append(
		args,
		"-display", "none",
		"-monitor", "none",
		"-serial", "stdio",
		"-S",
		"-gdb", fmt.Sprintf("tcp::%d", port),
	)
	if icount {
		args = append(args, "-icount", "shift=auto")
	}

	return args
}

// freePort is an unused TCP port (a listen-and-close probe; the
// window between the close and qemu's bind is the usual ephemeral
// race, acceptable against our own sequential starts).
func freePort(ctx context.Context) (port int, err error) {
	var lc net.ListenConfig
	l, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	addr, ok := l.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.Join(
			fmt.Errorf("assembly/qemu: port probe: %T is not tcp", l.Addr()),
			l.Close(),
		)
	}

	port = addr.Port
	defer func() { err = errors.Join(err, l.Close()) }()
	return port, nil
}

// dialStub waits for the gdbstub to accept connections: qemu listens
// once the machine is up, retrying covers the process startup race.
func dialStub(ctx context.Context, port int) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: time.Second}
	deadline := time.Now().Add(dialTimeout)
	for {
		conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return conn, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("assembly/qemu: gdbstub not reachable on port %d: %w", port, err)
		}

		time.Sleep(50 * time.Millisecond)
	}
}
