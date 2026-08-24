// Command assembly-server starts an HTTP server for interactive viewing
// of disassembled code.
//
// Usage:
//
//	assembly-server [-addr 127.0.0.1:8080] [-maxbytes 268435456]
//
// Open http://127.0.0.1:8080 in a browser and drag a binary (ELF or
// Mach-O) onto the page — the server will detect the format and
// architecture (cputype / e_machine), extract the section with the code,
// and disassemble it.
//
// Architectures are chosen by the server from the binary's cputype/e_machine
// and invoked directly (arm64, riscv) — there is no common architecture
// registry/interface.
package main

import (
	"flag"
	"log"

	"github.com/okneniz/assembly/server"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	maxBytes := flag.Int64("maxbytes", 256<<20, "max upload size in bytes")
	flag.Parse()

	srv := server.New(*addr, server.WithMaxBytes(*maxBytes))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("assembly-server: %v", err)
	}
}
