package dap

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandDirs(t *testing.T) {
	require.Equal(
		t,
		[]string{"/usr/local/go/bin", "/opt/homebrew/bin", "/usr/local/bin"},
		commandDirs(),
	)
}

// TestLookCommand - the resolution order of a command word: PATH
// first, then the extra directories, nowhere after that.
func TestLookCommand(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "fake-go")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	// PATH branch
	t.Setenv("PATH", dir)
	path, err := lookCommand("fake-go", nil)
	require.NoError(t, err)
	require.Equal(t, bin, path)

	// fallback branch: PATH without it, an extra dir carries the binary
	t.Setenv("PATH", t.TempDir())
	path, err = lookCommand("fake-go", []string{dir})
	require.NoError(t, err)
	require.Equal(t, bin, path)

	// nowhere
	t.Setenv("PATH", t.TempDir())
	_, err = lookCommand("fake-go", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fake-go")
}

// TestRelayForwardsConversation - cat as the child: the launch prefix
// and everything after it must arrive on the child's stdout verbatim.
func TestRelayForwardsConversation(t *testing.T) {
	launch := []byte(`{"seq":2,"type":"request","command":"launch","arguments":{}}`)
	rest := strings.NewReader("Content-Length: 2\r\n\r\n{}")

	var out bytes.Buffer
	require.NoError(t, Relay(t.Context(), "cat", launch, rest, &out, &out))

	require.Equal(t, string(envelope(launch))+"Content-Length: 2\r\n\r\n{}", out.String())
}

// TestRelayEmptyCommand - nothing to spawn is an error, not a hang.
func TestRelayEmptyCommand(t *testing.T) {
	var out bytes.Buffer
	err := Relay(t.Context(), "  ", nil, strings.NewReader(""), &out, &out)
	require.ErrorContains(t, err, "empty command")
}

// TestRelayMissingCommand - an unknown word reports the resolution
// failure instead of spawning.
func TestRelayMissingCommand(t *testing.T) {
	var out bytes.Buffer
	err := Relay(t.Context(), "no-such-command-anywhere", nil, strings.NewReader(""), &out, &out)
	require.ErrorContains(t, err, "no-such-command-anywhere")
}

// TestGoRunDir - the `go run DIR` working-directory rule: the package's
// own directory for an absolute argument, nothing for any other shape.
func TestGoRunDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "main.go")
	require.NoError(t, os.WriteFile(file, []byte("package main\n"), 0o644))

	cases := []struct {
		name   string
		fields []string
		want   string
	}{
		{name: "absolute dir", fields: []string{"go", "run", dir}, want: dir},
		{name: "absolute file", fields: []string{"/usr/local/go/bin/go", "run", file}, want: dir},
		{name: "flags skipped", fields: []string{"go", "run", "-count=1", dir}, want: dir},
		{name: "relative pkg", fields: []string{"go", "run", "./riscv"}, want: ""},
		{name: "no run", fields: []string{"/usr/local/go/bin/go", "build", dir}, want: ""},
		{name: "not go", fields: []string{"cat", dir}, want: ""},
		{name: "missing dir", fields: []string{"go", "run", "/no/such/dir"}, want: ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, goRunDir(c.fields))
		})
	}
}

// TestTailWriter - everything forwards, only the tail stays.
func TestTailWriter(t *testing.T) {
	var sink bytes.Buffer
	tw := newTailWriter(&sink)

	payload := strings.Repeat("x", tailLimit+100) + "THE-END"
	n, err := tw.Write([]byte(payload))
	require.NoError(t, err)
	require.Equal(t, len(payload), n)
	require.Equal(t, len(payload), sink.Len(), "everything forwarded")

	require.Equal(t, strings.Repeat("x", tailLimit-len("THE-END"))+"THE-END", tw.String())
}
