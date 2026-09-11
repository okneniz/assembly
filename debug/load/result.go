package load

import "github.com/okneniz/assembly/debug/session"

// Result is everything a debugging run needs: the executor image and
// the debug views over it (symbols, the line map, the source text).
type Result struct {
	Img      []byte
	Syms     map[string]uint64
	Lines    []session.Line
	SrcLines []string
	SrcName  string
}
