// Package dap is the editor frontend of the debugger: one Debug
// Adapter Protocol conversation over a byte transport (stdio for the
// VSCode adapter binary, pipes for tests). The server translates the
// protocol onto the debug/session engine and the qemu executor;
// sessions are booted through an injected Launcher (Boot is the
// production one).
package dap

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/okneniz/assembly/debug/rsp"
	"github.com/okneniz/assembly/debug/session"
)

// errDisconnect ends the conversation cleanly (the machine teardown
// belongs to Serve). errRelay hands the conversation to the spawned
// program instead (the session was never ours).
var (
	errDisconnect = errors.New("assembly/dap: disconnect")
	errRelay      = errors.New("assembly/dap: relay")
)

// The conversation's fixed coordinates: one bare-metal thread, one
// frame at the pc, the one Variables scope (the registers).
const (
	oneThread    = 1
	oneFrame     = 1
	registersRef = 1
)

// Server is one DAP conversation. Serve reads requests and answers
// them; continue and step resume in a goroutine and report back as
// stopped/terminated events. The lifecycle: initialize, launch (the
// Runtime appears, halted at the reset), the breakpoint requests,
// configurationDone (the entry stop is reported), drive, disconnect.
type Server struct {
	fr       *framing
	launcher Launcher

	mu       sync.Mutex
	seq      int
	rt       *Runtime
	lines    []session.Line
	bps      breakpointSet
	running  bool
	stepping bool
	pausing  bool
	finished bool // the run ended (exit or poweroff): the machine is gone
}

// NewServer is a conversation over the transport; launcher boots the
// run at the launch request.
func NewServer(in io.Reader, out io.Writer, launcher Launcher) (*Server, error) {
	if launcher == nil {
		return nil, errors.New("assembly/dap: no launcher")
	}

	return &Server{
		fr:       newFraming(in, out),
		launcher: launcher,
		bps:      newBreakpointSet(),
	}, nil
}

// Serve runs the conversation until the transport closes (EOF) or a
// disconnect ends the session; both tear the machine down.
func (s *Server) Serve() error {
	for {
		body, err := s.fr.read()
		if err != nil {
			s.shutdown()
			if errors.Is(err, io.EOF) {
				return nil
			}

			return fmt.Errorf("assembly/dap: read: %w", err)
		}

		if err := s.dispatch(body); err != nil {
			if errors.Is(err, errRelay) {
				return nil // the child owned everything from here on
			}

			s.shutdown()
			return nil // the one dispatch error is the clean disconnect
		}
	}
}

// dispatch decodes and runs one request; errDisconnect is the clean
// end, everything else answers the client in place.
func (s *Server) dispatch(body []byte) error {
	var req incoming
	if err := json.Unmarshal(body, &req); err != nil {
		s.warn(fmt.Errorf("assembly/dap: bad request: %w", err))
		return nil
	}

	switch req.Command {
	case "initialize":
		s.respond(req, true, "", newCapabilities())
	case "launch":
		return s.launch(req, body)
	case "attach":
		s.respond(req, false, "assembly/dap: attach is not supported (launch only)", nil)
	case "setBreakpoints":
		s.setBreakpoints(req)
	case "setFunctionBreakpoints":
		s.setFunctionBreakpoints(req)
	case "setInstructionBreakpoints":
		s.setInstructionBreakpoints(req)
	case "configurationDone":
		s.configurationDone(req)
	case "threads":
		s.threads(req)
	case "stackTrace":
		s.stackTrace(req)
	case "scopes":
		s.scopes(req)
	case "variables":
		s.variables(req)
	case "evaluate":
		s.evaluate(req)
	case "continue":
		s.resume(req, false)
	case "next", "stepIn", "stepOut":
		s.resume(req, true)
	case "pause":
		s.pause(req)
	case "disassemble":
		s.disassemble(req)
	case "readMemory":
		s.readMemory(req)
	case "disconnect", "terminate":
		s.respond(req, true, "", nil)
		return errDisconnect
	default:
		s.respond(req, false, fmt.Sprintf("assembly/dap: unknown command %q", req.Command), nil)
	}

	return nil
}

// ifFinished answers a target-touching request of an ended run: the
// editor keeps asking until the disconnect, and an error response
// paints the session red. ok reports whether the run is over and the
// graceful body went out.
func (s *Server) ifFinished(req incoming, body any) bool {
	s.mu.Lock()
	finished := s.finished
	s.mu.Unlock()

	if !finished {
		return false
	}

	s.respond(req, true, "", body)
	return true
}

// launch boots the run: the Runtime through the launcher, then the
// initialized event (the editor configures breakpoints next; the
// machine sits halted at its reset). A launch with a command delegates
// instead: the spawned program owns the whole session - it receives
// this very launch request (and everything after) and answers it, so
// no request is ever answered twice (initialize was ours, the launch
// onward is the child's).
func (s *Server) launch(req incoming, body []byte) error {
	var args launchArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: launch arguments: %v", err), nil)
		return nil
	}

	if args.Command != "" {
		// the child must not relay again: strip the command attribute
		// from the request it receives (its launcher owns the session)
		var relayed struct {
			Seq       int        `json:"seq"`
			Type      string     `json:"type"`
			Command   string     `json:"command"`
			Arguments launchArgs `json:"arguments"`
		}
		if err := json.Unmarshal(body, &relayed); err != nil {
			s.respond(req, false, fmt.Sprintf("assembly/dap: relay: %v", err), nil)
			return nil
		}

		relayed.Arguments.Command = ""
		clean, err := json.Marshal(relayed)
		if err != nil {
			s.respond(req, false, err.Error(), nil)
			return nil
		}

		// the diagnostic line lands on the editor's adapter log (the
		// launch command is exactly what breaks when an editor leaves a
		// path variable unexpanded)
		fmt.Fprintf(os.Stderr, "assembly-debug-dap: relaying to %q\n", args.Command)

		answered := false
		out := &firstWriteWatcher{w: s.fr.w, seen: &answered}
		childStderr := newTailWriter(os.Stderr)
		relayErr := Relay(context.Background(), args.Command, clean, s.fr.r, out, childStderr)
		if !answered && relayErr != nil {
			// the child died before answering the launch (a build
			// failure of `go run`, a bad path): the editor must hear it
			// from us, not spin on a dead conversation - with the
			// command and the child's own last words
			s.respond(
				req,
				false,
				fmt.Sprintf(
					"%v (command: %q; child stderr: %s)",
					relayErr,
					args.Command,
					childStderr,
				),
				nil,
			)
		}

		return errors.Join(errRelay, relayErr)
	}

	rt, err := s.launcher(args, &consoleWriter{srv: s})
	if err != nil {
		s.respond(req, false, err.Error(), nil)
		return nil
	}

	s.mu.Lock()
	s.rt = rt
	s.lines = sortByAddr(rt.Lines)
	s.mu.Unlock()

	s.respond(req, true, "", nil)
	s.notify("initialized", nil)
	return nil
}

// setBreakpoints replaces the source breakpoints of one file: every
// requested line resolves through the line map, the stale set is
// cleared before the new one goes in.
func (s *Server) setBreakpoints(req incoming) {
	var args setBreakpointsArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	if s.ifFinished(req, setBreakpointsResult{Breakpoints: unverified(args.Breakpoints)}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	for _, old := range s.bps.takeFile(args.Source.Path) {
		if err := rt.Sess.ClearAt(fmt.Sprintf("%#x", old)); err != nil {
			s.warn(err)
		}
	}

	installed := []uint64{}
	out := make([]breakpoint, 0, len(args.Breakpoints))
	for _, bp := range args.Breakpoints {
		addr, found := resolveLine(s.lines, args.Source.Path, bp.Line)
		if !found {
			out = append(
				out,
				breakpoint{Verified: false, Line: bp.Line, Message: "no code at this line"},
			)
			continue
		}

		if _, err := rt.Sess.BreakAt(fmt.Sprintf("%#x", addr)); err != nil {
			out = append(out, breakpoint{Verified: false, Line: bp.Line, Message: err.Error()})
			continue
		}

		installed = append(installed, addr)
		out = append(out, breakpoint{Verified: true, Line: bp.Line})
	}

	s.bps.putFile(args.Source.Path, installed)

	s.respond(req, true, "", setBreakpointsResult{Breakpoints: out})
}

// setFunctionBreakpoints replaces the named breakpoints: the names
// resolve through the symbol table (the labels of the program).
func (s *Server) setFunctionBreakpoints(req incoming) {
	var args setFunctionBreakpointsArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	dead := make([]breakpoint, 0, len(args.Breakpoints))
	for range args.Breakpoints {
		dead = append(dead, breakpoint{Verified: false, Message: "the program has finished"})
	}

	if s.ifFinished(req, setBreakpointsResult{Breakpoints: dead}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	for _, old := range s.bps.takeLabels() {
		if err := rt.Sess.ClearAt(fmt.Sprintf("%#x", old)); err != nil {
			s.warn(err)
		}
	}

	installed := []uint64{}
	out := make([]breakpoint, 0, len(args.Breakpoints))
	for _, bp := range args.Breakpoints {
		addr, err := rt.Sess.BreakAt(bp.Name)
		if err != nil {
			out = append(out, breakpoint{Verified: false, Message: err.Error()})
			continue
		}

		installed = append(installed, addr)
		out = append(out, breakpoint{Verified: true})
	}

	s.bps.putLabels(installed)

	s.respond(req, true, "", setBreakpointsResult{Breakpoints: out})
}

// setInstructionBreakpoints replaces the address breakpoints (the
// editor's disassembly-view breakpoints).
func (s *Server) setInstructionBreakpoints(req incoming) {
	var args setInstructionBreakpointsArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	dead := make([]breakpoint, 0, len(args.Breakpoints))
	for range args.Breakpoints {
		dead = append(dead, breakpoint{Verified: false, Message: "the program has finished"})
	}

	if s.ifFinished(req, setBreakpointsResult{Breakpoints: dead}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	for _, old := range s.bps.takeAddrs() {
		if err := rt.Sess.ClearAt(fmt.Sprintf("%#x", old)); err != nil {
			s.warn(err)
		}
	}

	installed := []uint64{}
	out := make([]breakpoint, 0, len(args.Breakpoints))
	for _, bp := range args.Breakpoints {
		base, err := strconv.ParseUint(bp.InstructionReference, 0, 64)
		if err != nil {
			out = append(
				out,
				breakpoint{
					Verified: false,
					Message:  fmt.Sprintf("assembly/dap: bad address %q", bp.InstructionReference),
				},
			)
			continue
		}

		at := int64(base) + int64(bp.Offset)
		if _, err := rt.Sess.BreakAt(fmt.Sprintf("%#x", at)); err != nil {
			out = append(out, breakpoint{Verified: false, Message: err.Error()})
			continue
		}

		installed = append(installed, uint64(at))
		out = append(out, breakpoint{Verified: true})
	}

	s.bps.putAddrs(installed)

	s.respond(req, true, "", setBreakpointsResult{Breakpoints: out})
}

// configurationDone ends the edit phase: the machine is already
// halted at its reset, report it as the entry stop.
func (s *Server) configurationDone(req incoming) {
	if _, ok := s.halted(req); !ok {
		return
	}

	s.respond(req, true, "", nil)
	s.notify("stopped", stoppedBody{Reason: "entry", ThreadId: oneThread, AllThreadsStopped: true})
}

// threads answers the single bare-metal thread without touching the
// target (safe while it runs).
func (s *Server) threads(req incoming) {
	s.mu.Lock()
	name := "cpu0"
	if s.rt != nil {
		name = s.rt.Tgt.Arch() + " cpu0"
	}

	s.mu.Unlock()

	s.respond(req, true, "", threadsResult{Threads: []thread{{Id: oneThread, Name: name}}})
}

// stackTrace is the single frame at the pc: the source line when the
// line map covers it, the instruction as the frame label, the address
// for the disassembly view.
func (s *Server) stackTrace(req incoming) {
	if s.ifFinished(req, stackTraceResult{TotalFrames: 0}) {
		return
	}

	// while a resume is in flight the pc is unaskable - the editor
	// refreshes the frames right after its next/continue response, the
	// stop event will trigger the real poll
	s.mu.Lock()
	running := s.running
	rt := s.rt
	s.mu.Unlock()

	if running || rt == nil {
		s.respond(req, true, "", stackTraceResult{TotalFrames: 0})
		return
	}

	pc, err := rt.Sess.PC()
	if err != nil {
		s.respond(req, false, err.Error(), nil)
		return
	}

	frame := stackFrame{
		Id:                          oneFrame,
		Name:                        fmt.Sprintf("%#x", pc),
		Column:                      1,
		InstructionPointerReference: fmt.Sprintf("%#x", pc),
	}

	if dis, derr := rt.Sess.Disasm(pc, 1); derr == nil && len(dis) > 0 {
		frame.Name = dis[0]
	}

	if line, found := rt.Sess.LineAt(pc); found {
		frame.Line = line.Line
		frame.Source = &sourceRef{Name: filepath.Base(line.File), Path: line.File}
	}

	s.respond(
		req,
		true,
		"",
		stackTraceResult{StackFrames: []stackFrame{frame}, TotalFrames: oneFrame},
	)
}

// scopes is the one container of the frame: the core registers.
func (s *Server) scopes(req incoming) {
	if s.ifFinished(req, scopesResult{}) {
		return
	}

	if _, ok := s.halted(req); !ok {
		return
	}

	s.respond(req, true, "", scopesResult{Scopes: []scope{{
		Name:               "Registers",
		PresentationHint:   "registers",
		VariablesReference: registersRef,
	}}})
}

// variables lists the registers, each with the memory reference at it
// (the editor's memory viewer follows sp and friends).
func (s *Server) variables(req incoming) {
	var args variablesArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	if s.ifFinished(req, variablesResult{}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	if args.VariablesReference != registersRef {
		s.respond(
			req,
			false,
			fmt.Sprintf("assembly/dap: unknown variables reference %d", args.VariablesReference),
			nil,
		)
		return
	}

	regs, err := rt.Sess.Regs()
	if err != nil {
		s.respond(req, false, err.Error(), nil)
		return
	}

	out := make([]variable, 0, len(regs))
	for _, r := range regs {
		out = append(out, variable{
			Name:            r.Name,
			Value:           fmt.Sprintf("%#x", r.Value),
			MemoryReference: fmt.Sprintf("%#x", r.Value),
		})
	}

	s.respond(req, true, "", variablesResult{Variables: out})
}

// evaluate renders a watch expression: a label or an address (hex or
// dec), shown as the 8-byte word at it.
func (s *Server) evaluate(req incoming) {
	var args evaluateArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	if s.ifFinished(req, evaluateResult{Result: "the program has finished"}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	addr, found := rt.Sess.Symbol(args.Expression)
	if !found {
		parsed, perr := strconv.ParseUint(args.Expression, 0, 64)
		if perr != nil {
			s.respond(
				req,
				false,
				fmt.Sprintf("assembly/dap: %q is neither a symbol nor an address", args.Expression),
				nil,
			)
			return
		}

		addr = parsed
	}

	data, err := rt.Sess.Read(addr, 8)
	if err != nil {
		s.respond(req, false, err.Error(), nil)
		return
	}

	s.respond(req, true, "", evaluateResult{
		Result: fmt.Sprintf("%#x", binary.LittleEndian.Uint64(data)),
		Type:   "uint64",
	})
}

// resume answers continue/step at once and finishes in a goroutine:
// the stop arrives as an event (the editor stays responsive for pause
// and disconnect meanwhile).
func (s *Server) resume(req incoming, step bool) {
	rt, ok := s.halted(req)
	if !ok {
		return
	}

	s.mu.Lock()
	s.running = true
	s.stepping = step
	s.pausing = false
	s.mu.Unlock()

	if step {
		go func() {
			stop, err := rt.Sess.Step()
			s.finishResume(stop, err)
		}()

		s.respond(req, true, "", nil)
		return
	}

	s.respond(req, true, "", continueResult{AllThreadsContinued: true})
	go func() {
		stop, err := rt.Sess.Continue()
		s.finishResume(stop, err)
	}()
}

// finishResume reports the stop of the resumed target: stopped with
// the reason (breakpoint, step, pause), or the end of the run.
func (s *Server) finishResume(stop rsp.StopReply, err error) {
	s.mu.Lock()
	s.running = false
	pausing := s.pausing
	stepping := s.stepping
	s.mu.Unlock()

	if err != nil || stop.Exited() {
		// the run is over: the finished flag arms the graceful answers
		// the same moment - no editor request can slip into the dead
		// machine anymore. Some poweroffs say farewell with a W reply
		// (arm64 PSCI), some just cut the wire (the riscv sifive_test)
		// - the broken pipe of the pending resume is the latter, not
		// an error worth red ink.
		s.mu.Lock()
		s.finished = true
		rt := s.rt
		s.rt = nil
		s.mu.Unlock()

		if rt != nil {
			// the teardown may block on the exiting executor - it runs
			// aside, the conversation is already unburdened
			go func() {
				if err := rt.Closer.Close(); err != nil {
					s.warn(err)
				}
			}()
		}

		if err == nil && stop.Kind == 'W' {
			s.notify("exited", exitedBody{ExitCode: stop.Signal})
		}

		if err != nil {
			s.notify("output", outputBody{
				Category: "console",
				Output:   "the program finished: the machine powered off\n",
			})
		}

		s.notify("terminated", nil)
		return
	}

	_, sw := stop.Fields["swbreak"]
	_, hw := stop.Fields["hwbreak"]
	reason := "breakpoint"
	if !sw && !hw {
		switch {
		case pausing:
			reason = "pause"
		case stepping:
			reason = "step"
		}
	}

	s.notify("stopped", stoppedBody{Reason: reason, ThreadId: oneThread, AllThreadsStopped: true})
}

// pause interrupts the running target: the pending resume reports the
// stop (reason pause).
func (s *Server) pause(req incoming) {
	s.mu.Lock()
	running := s.running
	rt := s.rt
	if running {
		s.pausing = true
	}

	s.mu.Unlock()

	if !running || rt == nil {
		s.respond(req, false, "assembly/dap: target is not running", nil)
		return
	}

	if err := rt.Sess.Interrupt(); err != nil {
		s.respond(req, false, err.Error(), nil)
		return
	}

	s.respond(req, true, "", nil)
}

// disassemble renders the listing window: the walk is per-instruction
// (the target's InstrLen reads the length at the head - the compressed
// riscv instructions vary), the reference plus both offsets opens it.
func (s *Server) disassemble(req incoming) {
	var args disassembleArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	if s.ifFinished(req, disassembleResult{}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	if args.InstructionCount <= 0 {
		s.respond(req, false, "assembly/dap: instructionCount must be positive", nil)
		return
	}

	base, err := strconv.ParseUint(args.MemoryReference, 0, 64)
	if err != nil {
		s.respond(
			req,
			false,
			fmt.Sprintf("assembly/dap: bad memory reference %q", args.MemoryReference),
			nil,
		)
		return
	}

	at := max(
		int64(base)+int64(args.Offset)+int64(args.InstructionOffset)*int64(rt.Tgt.InstrLen(nil)),
		0,
	)

	code, err := rt.Sess.Read(uint64(at), args.InstructionCount*rt.Tgt.InstrLen(nil))
	if err != nil {
		s.respond(req, false, err.Error(), nil)
		return
	}

	out := make([]disassembledInstruction, 0, args.InstructionCount)
	off := 0
	for len(out) < args.InstructionCount && off < len(code) {
		instr := code[off:]
		size := rt.Tgt.InstrLen(instr)
		if off+size > len(code) {
			break
		}

		text := ""
		if lines := rt.Tgt.Disasm(instr[:size], uint64(at)+uint64(off)); len(lines) > 0 {
			text = lines[0]
		}

		out = append(out, disassembledInstruction{
			Address:     fmt.Sprintf("%#x", uint64(at)+uint64(off)),
			Instruction: text,
		})

		off += size
	}

	s.respond(req, true, "", disassembleResult{Instructions: out})
}

// readMemory serves the editor's memory viewer: an unreadable window
// is a success with the unreadable tail counted (the DAP way).
func (s *Server) readMemory(req incoming) {
	var args readMemoryArgs
	if err := decodeArgs(req, &args); err != nil {
		s.respond(req, false, fmt.Sprintf("assembly/dap: arguments: %v", err), nil)
		return
	}

	if s.ifFinished(req, readMemoryResult{UnreadableBytes: args.Count}) {
		return
	}

	rt, ok := s.halted(req)
	if !ok {
		return
	}

	base, err := strconv.ParseUint(args.MemoryReference, 0, 64)
	if err != nil {
		s.respond(
			req,
			false,
			fmt.Sprintf("assembly/dap: bad memory reference %q", args.MemoryReference),
			nil,
		)
		return
	}

	at := max(int64(base)+int64(args.Offset), 0)

	if args.Count <= 0 {
		s.respond(req, false, "assembly/dap: count must be positive", nil)
		return
	}

	data, err := rt.Sess.Read(uint64(at), args.Count)
	if err != nil {
		s.respond(
			req,
			true,
			"",
			readMemoryResult{Address: fmt.Sprintf("%#x", at), UnreadableBytes: args.Count},
		)
		return
	}

	s.respond(req, true, "", readMemoryResult{Address: fmt.Sprintf("%#x", at), Data: data})
}

// halted reports whether the target is stopped and answers requests:
// the RSP conversation is one-in-flight, so nothing touches it while a
// resume is pending. ok carries the runtime for the handler.
func (s *Server) halted(req incoming) (*Runtime, bool) {
	s.mu.Lock()
	rt := s.rt
	running := s.running
	s.mu.Unlock()

	if rt == nil {
		s.respond(req, false, "assembly/dap: no session (launch first)", nil)
		return nil, false
	}

	if running {
		s.respond(req, false, "assembly/dap: target is running (pause it first)", nil)
		return nil, false
	}

	return rt, true
}

// respond answers one request: success with the body, or the message
// as the failure reason.
func (s *Server) respond(req incoming, success bool, message string, body any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	s.emit(responseMsg{
		Seq:        s.seq,
		Type:       "response",
		RequestSeq: req.Seq,
		Success:    success,
		Command:    req.Command,
		Message:    message,
		Body:       body,
	})
}

// notify fires one event to the client.
func (s *Server) notify(event string, body any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	s.emit(eventMsg{Seq: s.seq, Type: "event", Event: event, Body: body})
}

// emit marshals and writes; the caller holds mu (framing serializes
// the bytes, the seq counter needs mu). Write errors are dropped by
// design: a dying transport surfaces on the next read of the loop,
// and emit may not propagate (it would strand the conversation in a
// half-answered request).
func (s *Server) emit(msg any) {
	raw, err := json.Marshal(msg)
	if err != nil {
		// an unmarshalable DTO is a programming bug: report it without
		// the notify bookkeeping (mu is already held on this path)
		fallback, ferr := json.Marshal(eventMsg{
			Type:  "event",
			Event: "output",
			Body:  outputBody{Category: "stderr", Output: err.Error()},
		})
		if ferr != nil {
			return
		}

		raw = fallback
	}

	if werr := s.fr.write(raw); werr != nil {
		return // nothing else to do: the read loop owns the transport verdict
	}
}

// warn surfaces a non-fatal failure on the debug console.
func (s *Server) warn(err error) {
	s.notify("output", outputBody{Category: "stderr", Output: err.Error() + "\n"})
}

// shutdown tears the run down once: the machine dies with the
// conversation.
func (s *Server) shutdown() {
	s.mu.Lock()
	rt := s.rt
	s.rt = nil
	s.mu.Unlock()

	if rt == nil {
		return
	}

	if err := rt.Closer.Close(); err != nil {
		s.warn(err)
	}
}

// decodeArgs unmarshals the arguments of one request into out (an
// argument-less request leaves the zero value).
func decodeArgs(req incoming, out any) error {
	if len(req.Arguments) == 0 {
		return nil
	}

	if err := json.Unmarshal(req.Arguments, out); err != nil {
		return fmt.Errorf("assembly/dap: %w", err)
	}

	return nil
}
