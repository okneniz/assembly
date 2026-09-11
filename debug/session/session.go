// Package session is the debugger engine: one debugging run over an
// executor's RSP connection. This is the library API - a debugging
// script is an ordinary Go program (or test) driving a Session; the
// CLI REPL is just one consumer of it.
//
// The session is arch-blind: register numbers and the disassembly come
// from the debug.Target; the symbols and the line map are supplied by
// the caller (the asm core's Result.Lines for .s sources, the prog
// Result.Lines for Go-written programs - normalized into Line).
package session

import (
	"fmt"
	"net"
	"sort"
	"strconv"

	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/debug/rsp"
)

// Session is one debugging run. The lifecycle: New connects to a
// booted-frozen executor, the methods drive it, Close tears it down.
type Session struct {
	cl    *rsp.Client
	tgt   debug.Target
	syms  map[string]uint64
	lines []Line
}

// New is a session over an executor connection: it negotiates the
// protocol and synchronizes with the halted target.
func New(conn net.Conn, tgt debug.Target, syms map[string]uint64, lines []Line) (*Session, error) {
	s := &Session{
		cl:    rsp.NewClient(rsp.NewConn(conn)),
		tgt:   tgt,
		syms:  syms,
		lines: sortedLines(lines),
	}

	if _, err := s.cl.Supported(); err != nil {
		return nil, fmt.Errorf("assembly/session: negotiate: %w", err)
	}

	if _, err := s.cl.HaltReason(); err != nil {
		return nil, fmt.Errorf("assembly/session: halt reason: %w", err)
	}

	return s, nil
}

// PC is the program counter.
func (s *Session) PC() (uint64, error) {
	return s.cl.ReadReg(s.tgt.PCNum())
}

// SetPC writes the program counter.
func (s *Session) SetPC(addr uint64) error {
	return s.cl.WriteReg(s.tgt.PCNum(), addr)
}

// BreakAt sets a breakpoint at a symbol name or a hex address
// ("0x401000"): resolves through the symbol table first.
func (s *Session) BreakAt(where string) (uint64, error) {
	addr, err := s.resolve(where)
	if err != nil {
		return 0, err
	}

	if err := s.cl.SetBreak(addr, s.tgt.InstrLen(nil)); err != nil {
		return 0, fmt.Errorf("assembly/session: break %s: %w", where, err)
	}

	return addr, nil
}

// ClearAt removes the breakpoint at a symbol name or a hex address.
func (s *Session) ClearAt(where string) error {
	addr, err := s.resolve(where)
	if err != nil {
		return err
	}

	return s.cl.ClearBreak(addr, s.tgt.InstrLen(nil))
}

// Continue resumes and blocks until the stop report.
func (s *Session) Continue() (rsp.StopReply, error) {
	return s.cl.Continue()
}

// Step executes one instruction and returns the stop report.
func (s *Session) Step() (rsp.StopReply, error) {
	return s.cl.Step()
}

// Interrupt asks a running target to stop (from another goroutine
// while Continue blocks).
func (s *Session) Interrupt() error {
	return s.cl.Interrupt()
}

// Regs is the core register dump: the 'g' block sliced by the target
// layout. Registers wider than 64 bits are absent from the dump.
func (s *Session) Regs() ([]RegValue, error) {
	block, err := s.cl.ReadRegs()
	if err != nil {
		return nil, fmt.Errorf("assembly/session: registers: %w", err)
	}

	out := make([]RegValue, 0, 16)
	off := 0
	for _, r := range s.tgt.Registers() {
		if r.Bits > 64 {
			continue // vector registers: outside the display
		}

		if off+r.Width() > len(block) {
			break // the target described more than it served
		}

		var v uint64
		for i := range r.Width() {
			v |= uint64(block[off+i]) << (8 * i)
		}

		out = append(out, NewRegValue(r.Name, v))
		off += r.Width()
	}

	return out, nil
}

// Read is n bytes of target memory at addr.
func (s *Session) Read(addr uint64, n int) ([]byte, error) {
	return s.cl.ReadMem(addr, n)
}

// Write patches target memory at addr.
func (s *Session) Write(addr uint64, b []byte) error {
	return s.cl.WriteMem(addr, b)
}

// Disasm reads n instructions of memory at addr and renders them
// through the target's decoders.
func (s *Session) Disasm(addr uint64, n int) ([]string, error) {
	chunk := n * s.tgt.InstrLen(nil)
	code, err := s.Read(addr, chunk)
	if err != nil {
		return nil, fmt.Errorf("assembly/session: disasm read: %w", err)
	}

	return s.tgt.Disasm(code, addr), nil
}

// Symbols is the sorted symbol table of the debugged program.
func (s *Session) Symbols() []string {
	out := make([]string, 0, len(s.syms))
	for name := range s.syms {
		out = append(out, name)
	}

	sort.Strings(out)
	return out
}

// Symbol resolves one name to its address.
func (s *Session) Symbol(name string) (uint64, bool) {
	addr, ok := s.syms[name]
	return addr, ok
}

// LineAt is the line map entry of addr: the last entry starting at or
// before it (the instruction addr may sit mid-line for multi-instruction
// source lines). ok is false below the first entry.
func (s *Session) LineAt(addr uint64) (Line, bool) {
	i := sort.Search(len(s.lines), func(i int) bool {
		return s.lines[i].Addr > addr
	}) - 1
	if i < 0 {
		return Line{}, false
	}

	return s.lines[i], true
}

// resolve maps a breakpoint specification to an address.
func (s *Session) resolve(where string) (uint64, error) {
	if addr, ok := s.syms[where]; ok {
		return addr, nil
	}

	addr, err := strconv.ParseUint(where, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("assembly/session: %q is neither a symbol nor an address", where)
	}

	return addr, nil
}

// sortedLines orders the line map by address (Search needs it; both
// producers already deliver sorted, this makes the contract local).
func sortedLines(lines []Line) []Line {
	out := append([]Line{}, lines...)
	sort.Slice(out, func(i, j int) bool {
		return out[i].Addr < out[j].Addr
	})

	return out
}
