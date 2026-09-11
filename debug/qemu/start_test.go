package qemu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQemuDirs(t *testing.T) {
	require.Equal(
		t,
		[]string{"/opt/homebrew/bin", "/usr/local/bin"},
		qemuDirs(),
	)
}

// TestFindQemu - the resolution order: PATH first, then the extra
// directories, nowhere after that.
func TestFindQemu(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "qemu-system-fake")
	require.NoError(t, os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755))

	// PATH branch: the fake dir on PATH resolves the binary
	t.Setenv("PATH", dir)
	path, err := findQemu("qemu-system-fake", nil)
	require.NoError(t, err)
	require.Equal(t, bin, path)

	// fallback branch: PATH without it, the extra dir carries the binary
	t.Setenv("PATH", t.TempDir())
	path, err = findQemu("qemu-system-fake", []string{dir})
	require.NoError(t, err)
	require.Equal(t, bin, path)

	// nowhere: both branches empty
	_, err = findQemu("qemu-system-fake", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "qemu-system-fake")
}
