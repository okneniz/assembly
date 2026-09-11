// Command assembly-debug-dap is the editor adapter of the debugger:
// one Debug Adapter Protocol conversation over stdio (the editor's
// debug type "assembly-debug" spawns it; vscode/assembly-debug is the
// extension that declares it). The launch configuration picks the
// session kind: .s sources (assembled in-process) and prebuilt images
// boot right here, while a `command` attribute relays the conversation
// to a program serving DAP itself (a prog example in its -debug mode).
package main

import (
	"fmt"
	"os"

	"github.com/okneniz/assembly/debug/dap"
)

func main() {
	srv, err := dap.NewServer(os.Stdin, os.Stdout, dap.Boot)
	if err != nil {
		fmt.Fprintln(os.Stderr, "assembly-debug-dap:", err)
		os.Exit(1)
	}

	if err := srv.Serve(); err != nil {
		fmt.Fprintln(os.Stderr, "assembly-debug-dap:", err)
		os.Exit(1)
	}
}
