package dap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// commandDirs are the executable prefixes outside the minimal launchd
// PATH a GUI-launched editor hands its children: the Go toolchain and
// homebrew.
func commandDirs() []string {
	return []string{"/usr/local/go/bin", "/opt/homebrew/bin", "/usr/local/bin"}
}

// LookCommand resolves a command word the way the executor resolves
// qemu: PATH first, then the well-known prefixes (a dock-launched
// editor never sees /usr/local/go/bin, and `go run ...` commands need
// it).
func LookCommand(bin string) (string, error) {
	return lookCommand(bin, commandDirs())
}

func lookCommand(bin string, dirs []string) (string, error) {
	if path, err := exec.LookPath(bin); err == nil {
		return path, nil
	}

	for _, dir := range dirs {
		at := dir + "/" + bin
		if _, err := exec.LookPath(at); err == nil {
			return at, nil
		}
	}

	return "", fmt.Errorf("assembly/dap: %q: executable not found in PATH or %v", bin, dirs)
}

// envelope is the wire form of one body: the deterministic framing
// around it.
func envelope(body []byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Content-Length: %d\r\n\r\n", len(body))
	buf.Write(body)
	return buf.Bytes()
}

// goRunDir is the working directory for a `go run DIR` command: the
// package's own directory when the argument is absolute. The editor
// spawns the adapter wherever it pleases (not the repository), and go
// resolves the module from the working directory - outside it even an
// absolute package path dies with "go.mod file not found". Commands of
// any other shape inherit the adapter's directory.
func goRunDir(fields []string) string {
	if len(fields) < 3 || filepath.Base(fields[0]) != "go" || fields[1] != "run" {
		return ""
	}

	for _, arg := range fields[2:] {
		if strings.HasPrefix(arg, "-") {
			continue
		}

		if !filepath.IsAbs(arg) {
			return ""
		}

		info, err := os.Stat(arg)
		if err != nil {
			return ""
		}

		if info.IsDir() {
			return arg
		}

		return filepath.Dir(arg)
	}

	return ""
}

// Relay hands the conversation to a program serving DAP itself (a prog
// example in its -debug mode): the child receives the launch request
// that named it - and everything the editor sends after - and answers
// on this stdout. The adapter has already answered initialize by the
// time a launch can arrive, so the child never sees it and no request
// is answered twice. rest is the server's buffered reader (bytes
// beyond the launch envelope may already sit in it).
func Relay(
	ctx context.Context,
	command string,
	launch []byte,
	rest io.Reader,
	out, errW io.Writer,
) error {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return errors.New("assembly/dap: relay: empty command")
	}

	bin, err := LookCommand(fields[0])
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, bin, fields[1:]...)
	cmd.Stdout = out
	cmd.Stderr = errW
	cmd.Dir = goRunDir(fields)

	// The child's stdin is OUR pipe with OUR copier: exec's managed
	// Stdin would hold Wait hostage on the editor's (still open) stdin
	// long after a child that died early - a failed `go run` would hang
	// the launch forever. With a manual pipe, Wait returns at the
	// process exit and the close below unblocks the copier.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("assembly/dap: relay: stdin: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("assembly/dap: relay: %w", err)
	}

	// The copier is fire-and-forget on purpose: it ends when the editor
	// closes the stream (or its next write fails into the dead child's
	// pipe) - joining it would block on the editor's still-open stdin,
	// exactly the hang a failed child must not cause. The adapter
	// process ends with the conversation; the goroutine dies with it.
	go func() {
		_, err := io.Copy(stdin, io.MultiReader(bytes.NewReader(envelope(launch)), rest))
		if cerr := stdin.Close(); err == nil {
			err = cerr
		}

		// the outcome belongs to the stream: a dead child or a closed
		// editor side ended the copy, and nobody is left to hear it
		_ = err
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("assembly/dap: relay: %w", err)
	}

	return nil
}
