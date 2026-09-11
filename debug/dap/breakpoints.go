package dap

import (
	"path/filepath"
	"sort"

	"github.com/okneniz/assembly/debug/session"
)

// breakpointSet tracks the addresses the adapter installed, per
// breakpoint kind: a fresh set* request replaces the whole kind (the
// DAP contract), so the stale addresses must be cleared first. (A
// source and a function breakpoint on one address collapse into the
// same stub slot - accepted, the editor never stacks them.)
type breakpointSet struct {
	files  map[string][]uint64
	labels []uint64
	addrs  []uint64
}

func newBreakpointSet() breakpointSet {
	return breakpointSet{files: map[string][]uint64{}}
}

// takeFile returns the installed addresses of the source file and
// empties the slot (the caller clears them before installing anew).
func (b *breakpointSet) takeFile(path string) []uint64 {
	out := b.files[path]
	delete(b.files, path)
	return out
}

// putFile remembers the installed addresses of the source file.
func (b *breakpointSet) putFile(path string, addrs []uint64) {
	b.files[path] = addrs
}

// takeLabels returns the installed function-breakpoint addresses and
// empties the slot.
func (b *breakpointSet) takeLabels() []uint64 {
	out := b.labels
	b.labels = nil
	return out
}

// putLabels remembers the installed function-breakpoint addresses.
func (b *breakpointSet) putLabels(addrs []uint64) {
	b.labels = addrs
}

// takeAddrs returns the installed instruction-breakpoint addresses and
// empties the slot.
func (b *breakpointSet) takeAddrs() []uint64 {
	out := b.addrs
	b.addrs = nil
	return out
}

// putAddrs remembers the installed instruction-breakpoint addresses.
func (b *breakpointSet) putAddrs(addrs []uint64) {
	b.addrs = addrs
}

// unverified is the verdict for every requested source breakpoint of
// an ended run (the editor may re-send breakpoints until the
// disconnect - the grey dot with a reason, no error).
func unverified(requested []sourceBreakpoint) []breakpoint {
	out := make([]breakpoint, 0, len(requested))
	for _, bp := range requested {
		out = append(
			out,
			breakpoint{Verified: false, Line: bp.Line, Message: "the program has finished"},
		)
	}

	return out
}

// resolveLine is the address of the source line: the first line-map
// entry of the line (the entries come address-ordered, so the earliest
// instruction of the line wins). The file matches by full path, with
// the base name as the fallback (the editor may spell the path
// differently than the launch configuration did).
func resolveLine(lines []session.Line, file string, line int) (uint64, bool) {
	for _, e := range lines {
		if e.Line != line {
			continue
		}

		if e.File == file || filepath.Base(e.File) == filepath.Base(file) {
			return e.Addr, true
		}
	}

	return 0, false
}

// sortByAddr is the line map in address order (the resolution and the
// editor's expectations both want it; the producers deliver sorted,
// this makes the contract local).
func sortByAddr(lines []session.Line) []session.Line {
	out := append([]session.Line{}, lines...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Addr < out[j].Addr
	})

	return out
}
