// Package server implements a minimal HTTP server for interactive viewing
// of disassembled code.
//
// Routes:
//   - POST /api/v1/disasm — accepts a binary (multipart/form-data, field
//     "file"), disassembles its code section, and returns JSON.
//   - POST /api/v1/asm    — accepts JSON {arch, source, baseAddr}, assembles
//     it via the selected architecture's facade (asm/arm64/alias.Assemble
//     or asm/riscv/pseudo.Assemble) and returns sections/symbols/errors.
//   - GET  /             — serves index.html from the embedded static folder.
//
// The server is stateless: an uploaded file lives only for the duration
// of the request.
package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed static
var staticFS embed.FS

// Server is the HTTP server configuration.
type Server struct {
	addr     string // address like "127.0.0.1:8080"
	maxBytes int64  // request body size limit, bytes
	static   http.Handler
}

// Option is a functional configuration option for Server.
type Option func(*Server)

// WithMaxBytes sets the request body size limit (default 256 MiB).
func WithMaxBytes(n int64) Option {
	return func(s *Server) {
		if n > 0 {
			s.maxBytes = n
		}
	}
}

// New creates a server listening on addr. The architecture is selected
// automatically from the uploaded binary (cputype/e_machine); handlers
// switch on ArchKind and call the arch-specific Parse directly.
func New(addr string, opts ...Option) *Server {
	s := &Server{
		addr:     addr,
		maxBytes: 256 << 20, // 256 MiB
	}
	for _, opt := range opts {
		opt(s)
	}

	// Static assets: serve the contents of the static folder at the route root.
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		// impossible: embed guarantees "static" exists
		log.Fatalf("assembly/server: %v", err)
	}

	s.static = http.FileServer(http.FS(sub))

	return s
}

// ListenAndServe starts the HTTP server. Blocks until an error.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/disasm", s.handleDisasm)
	mux.HandleFunc("/api/v1/asm", s.handleAsm)
	// The root and everything else is frontend static.
	mux.Handle("/", s.static)

	srv := &http.Server{
		Addr:              s.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // slowloris (gosec G112)
	}
	log.Printf("assembly/server: listening on http://%s", s.addr)
	return srv.ListenAndServe()
}
