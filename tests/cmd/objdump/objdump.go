// Package objdump - running llvm-objdump and parsing its output (a
// parsec grammar for the instruction line). Shared by the differential test
// (assembly_test.go) and tests/cmd/assembly-diff; the logic used to be
// duplicated.
package objdump

import (
	"context"
	"errors"
	"os/exec"

	"github.com/okneniz/assembly/file"
)

// Run launches llvm-objdump with the given arguments, trying candidates until
// one is found that supports the target (Apple's /usr/bin/objdump has no
// RISC-V target, so homebrew llvm-objdump is preferable - it also handles
// --macho for ARM64). ctx bounds the running time of the external tool.
func Run(ctx context.Context, args []string) ([]byte, error) {
	candidates := []string{
		"/opt/homebrew/opt/llvm/bin/llvm-objdump",
		"/opt/homebrew/bin/llvm-objdump",
		"/usr/local/opt/llvm/bin/llvm-objdump",
		"/usr/local/bin/llvm-objdump",
		"llvm-objdump",
	}
	var lastErr error
	for _, c := range candidates {
		if _, err := exec.LookPath(c); err != nil {
			continue
		}

		out, err := exec.CommandContext(ctx, c, args...).Output()
		if err == nil {
			return out, nil
		}

		lastErr = err
	}

	if lastErr == nil {
		lastErr = errors.New("no llvm-objdump found on PATH")
	}

	return nil, lastErr
}

// Args returns the objdump arguments (including the binary path) for the
// format and architecture. Mach-O needs --macho (+ --arch arm64); for ELF,
// objdump itself determines the format and architecture from the header, so
// -d is enough.
func Args(format string, kind file.ArchKind, path string) []string {
	if format == "Mach-O" {
		return []string{"-d", "--macho", "--arch", "arm64", path}
	}

	return []string{"-d", path}
}
