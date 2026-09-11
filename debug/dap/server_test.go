package dap

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/debug/session"
)

// fakeTarget is the two-register target of the engine tests: pc is
// register 1 (register 0 doubles as the only GPR shown in dumps).
type fakeTarget struct{}

func (fakeTarget) Arch() string             { return "fake" }
func (fakeTarget) QemuBinary() string       { return "qemu-system-fake" }
func (fakeTarget) PCNum() int               { return 1 }
func (fakeTarget) SPNum() int               { return 0 }
func (fakeTarget) InstrLen([]byte) int      { return 4 }
func (fakeTarget) QemuArgs(string) []string { return nil }
func (fakeTarget) Disasm(code []byte, addr uint64) []string {
	return []string{fmt.Sprintf("%x: % x", addr, code)}
}

func (fakeTarget) Registers() []debug.Reg {
	return []debug.Reg{
		debug.NewReg("r0", 0, 64),
		debug.NewReg("pc", 1, 64),
	}
}

// dialogStep is one expected RSP request with its reply (the rsp fake
// without retransmission - the framing is covered there). A step with
// a gate holds its reply until the gate closes: the deterministic
// form of a slow target for the interrupt test. A step with cut
// closes the connection instead of replying: the riscv sifive_test
// poweroff leaves without a farewell stop reply.
type dialogStep struct {
	expect string
	reply  string
	gate   chan struct{}
	cut    bool
}

// The wire codec in miniature for the dialog stub (the rsp package
// owns the real one; its exported surface stays clean).
func testChecksum(payload string) byte {
	var sum byte
	for i := range len(payload) {
		sum += payload[i]
	}

	return sum
}

func testEncode(payload string) []byte {
	return []byte(fmt.Sprintf("$%s#%02x", payload, testChecksum(payload)))
}

func testParse(pkt []byte) (string, bool) {
	if len(pkt) < 4 || pkt[0] != '$' || pkt[len(pkt)-3] != '#' {
		return "", false
	}

	payload := string(pkt[1 : len(pkt)-3])
	var want byte
	if _, err := fmt.Sscanf(string(pkt[len(pkt)-2:]), "%02x", &want); err != nil {
		return "", false
	}

	return payload, testChecksum(payload) == want
}

// serveStub runs the fake target side of the RSP conversation over a
// TCP loopback pair (a buffered transport: a concurrent interrupt
// write must not block the way it would on a raw pipe); the returned
// wait collects its verdict (all steps served in order).
func serveStub(t *testing.T, steps []dialogStep) (net.Conn, func() error) {
	t.Helper()

	var lc net.ListenConfig
	lst, err := lc.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)

	type acceptedConn struct {
		conn net.Conn
		err  error
	}

	accepted := make(chan acceptedConn, 1)
	go func() {
		conn, aerr := lst.Accept()
		accepted <- acceptedConn{conn, errors.Join(aerr, lst.Close())}
	}()

	var dialer net.Dialer
	client, err := dialer.DialContext(t.Context(), "tcp", lst.Addr().String())
	require.NoError(t, err)

	res := <-accepted
	require.NoError(t, res.err, "stub: accept")
	server := res.conn

	done := make(chan error, 1)
	go func() {
		defer close(done)
		r := bufio.NewReader(server)
		w := server
		fail := func(format string, args ...any) {
			done <- errors.Join(fmt.Errorf(format, args...), server.Close())
		}

		readRequest := func() (string, bool) {
			for {
				b, rerr := r.ReadByte()
				if rerr != nil {
					return "", false
				}

				if b == '$' {
					pkt := []byte{'$'}
					for {
						b, rerr = r.ReadByte()
						if rerr != nil {
							return "", false
						}

						pkt = append(pkt, b)
						if b == '#' {
							hi, herr := r.ReadByte()
							lo, lerr := r.ReadByte()
							if herr != nil || lerr != nil {
								return "", false
							}

							pkt = append(pkt, hi, lo)
							break
						}
					}

					return testParse(pkt)
				}

				// everything before '$' is a concurrent interrupt byte:
				// addressed to the target, not part of any request
			}
		}

		for _, st := range steps {
			payload, ok := readRequest()
			if !ok || payload != st.expect {
				fail("stub: request %q (ok=%v), want %q", payload, ok, st.expect)
				return
			}

			if _, werr := w.Write([]byte{'+'}); werr != nil {
				fail("stub: ack write: %v", werr)
				return
			}

			if st.gate != nil {
				<-st.gate
			}

			if st.cut {
				done <- server.Close()
				return
			}

			if _, werr := w.Write(testEncode(st.reply)); werr != nil {
				fail("stub: reply write: %v", werr)
				return
			}

			for { // the reply ack; a concurrent interrupt byte is skipped
				b, rerr := r.ReadByte()
				if rerr != nil {
					fail("stub: reply ack read: %v", rerr)
					return
				}

				if b == 0x03 {
					continue
				}

				if b != '+' {
					fail("stub: reply ack %#02x", b)
					return
				}

				break
			}
		}

		done <- server.Close()
	}()

	return client, func() error {
		require.NoError(t, client.Close())
		return <-done
	}
}

// closeSpy is the Runtime teardown in the tests: did disconnect close
// the machine (atomic: Serve's shutdown writes it from its goroutine).
type closeSpy struct {
	closed atomic.Bool
}

func (c *closeSpy) Close() error {
	c.closed.Store(true)
	return nil
}

// testLauncher is the Launcher over the scripted stub dialog: the
// launch boots the session against the fake target (the New
// handshake heads the dialog), the verdict of the stub joins the test
// cleanup. The spy reports the machine teardown.
func testLauncher(
	t *testing.T,
	steps []dialogStep,
	syms map[string]uint64,
	lines []session.Line,
) (Launcher, *closeSpy) {
	t.Helper()

	handshake := []dialogStep{
		{
			expect: "qSupported:multiprocess-;swbreak+;hwbreak+;xmlRegisters=aarch64,riscv:rv64,loongarch64,i386",
			reply:  "qXfer:features:read+",
		},
		{expect: "?", reply: "T05thread:p1.1;"},
	}

	spy := &closeSpy{}
	return func(args launchArgs, console io.Writer) (*Runtime, error) {
		conn, wait := serveStub(t, append(handshake, steps...))
		t.Cleanup(func() { require.NoError(t, wait()) })

		s, err := session.New(conn, fakeTarget{}, syms, lines)
		if err != nil {
			return nil, err
		}

		return &Runtime{Tgt: fakeTarget{}, Sess: s, Lines: lines, Closer: spy}, nil
	}, spy
}

// wireMsg is the editor-side view of one DAP message.
type wireMsg struct {
	Type    string          `json:"type"`
	Event   string          `json:"event"`
	ReqSeq  int             `json:"request_seq"`
	Success bool            `json:"success"`
	Command string          `json:"command"`
	Message string          `json:"message"`
	Body    json.RawMessage `json:"body"`
}

// testClient is the editor side of the test conversation: sequential
// requests, events collected in arrival order.
type testClient struct {
	t      *testing.T
	fr     *framing
	conn   net.Conn
	seq    int
	events []wireMsg
}

func newTestClient(t *testing.T, conn net.Conn) *testClient {
	t.Helper()

	return &testClient{t: t, fr: newFraming(conn, conn), conn: conn}
}

func (c *testClient) request(command string, args any) wireMsg {
	c.t.Helper()
	c.seq++
	raw, err := json.Marshal(map[string]any{
		"seq":       c.seq,
		"type":      "request",
		"command":   command,
		"arguments": args,
	})
	require.NoError(c.t, err)
	require.NoError(c.t, c.fr.write(raw))
	return c.waitResponse()
}

func (c *testClient) waitResponse() wireMsg {
	c.t.Helper()
	for {
		msg := c.read()
		if msg.Type == "response" {
			require.Equal(c.t, c.seq, msg.ReqSeq, "reply to another request")
			return msg
		}

		c.events = append(c.events, msg)
	}
}

// event returns the next event of that name (buffered ones first),
// failing the test on any other event in between.
func (c *testClient) event(name string) wireMsg {
	c.t.Helper()
	for i, ev := range c.events {
		if ev.Event == name {
			c.events = append(c.events[:i], c.events[i+1:]...)
			return ev
		}
	}

	for {
		msg := c.read()
		if msg.Type != "event" {
			c.t.Fatalf("want event %q, got %s message", name, msg.Type)
		}

		if msg.Event == name {
			return msg
		}

		c.t.Fatalf("want event %q, got event %q", name, msg.Event)
	}
}

func (c *testClient) read() wireMsg {
	c.t.Helper()
	require.NoError(c.t, c.conn.SetReadDeadline(time.Now().Add(10*time.Second)))
	body, err := c.fr.read()
	require.NoError(c.t, err)

	var msg wireMsg
	require.NoError(c.t, json.Unmarshal(body, &msg))
	return msg
}

// decodeBody unmarshals the body of one message into out.
func (c *testClient) decodeBody(msg wireMsg, out any) {
	c.t.Helper()
	require.NoError(c.t, json.Unmarshal(msg.Body, out))
}

// newTestServer starts the server over a pipe pair with the scripted
// launcher (nil steps: no launch may happen); the spy returns for the
// teardown assertions. The Serve verdict joins the cleanup (a clean
// conversation ends in a nil).
func newTestServer(
	t *testing.T,
	steps []dialogStep,
	syms map[string]uint64,
	lines []session.Line,
) (*testClient, *closeSpy) {
	t.Helper()

	editor, adapter := net.Pipe()

	launcher, spy := testLauncher(t, steps, syms, lines)
	srv, err := NewServer(adapter, adapter, launcher)
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve() }()
	t.Cleanup(func() { require.NoError(t, <-serveErr) })
	t.Cleanup(func() { require.NoError(t, editor.Close()) })

	return newTestClient(t, editor), spy
}

// TestServerFlow - the whole editor conversation over one scripted
// session: launch, line and label breakpoints, the entry stop, a run
// to a label, the frame, the registers, the watch expression, the
// listing and memory windows, and the teardown.
func TestServerFlow(t *testing.T) {
	lines := []session.Line{
		session.NewLine("hello.s", 2, 0x401000, 4),
		session.NewLine("hello.s", 3, 0x401004, 4),
	}
	syms := map[string]uint64{"loop": 0x401004}

	steps := []dialogStep{
		{expect: "Z0,401000,4", reply: "OK"}, // line 2 breakpoint
		{expect: "Z0,401004,4", reply: "OK"}, // the loop label
		{expect: "c", reply: "T05swbreak:;"},
		{expect: "p1", reply: "0010400000000000"},
		{expect: "m401000,4", reply: "200080d2"},
		{expect: "g", reply: "4100000000000000" + "0010400000000000"},
		{expect: "m401004,8", reply: "0800000000000000"},
		{expect: "m401000,8", reply: "200080d2200080d2"},
		{expect: "m401000,4", reply: "200080d2"},
	}
	c, spy := newTestServer(t, steps, syms, lines)

	caps := capabilities{}
	c.decodeBody(c.request("initialize", nil), &caps)
	require.True(t, caps.SupportsConfigurationDoneRequest)
	require.True(t, caps.SupportsFunctionBreakpoints)
	require.True(t, caps.SupportsDisassembleRequest)

	c.request("launch", launchArgs{Arch: "fake"})
	c.event("initialized")

	// line breakpoints: the real line verified, a stray one grey
	res := setBreakpointsResult{}
	c.decodeBody(c.request("setBreakpoints", setBreakpointsArgs{
		Source:      sourceRef{Path: "hello.s"},
		Breakpoints: []sourceBreakpoint{{Line: 2}},
	}), &res)
	require.Len(t, res.Breakpoints, 1)
	require.True(t, res.Breakpoints[0].Verified)
	require.Equal(t, 2, res.Breakpoints[0].Line)

	res = setBreakpointsResult{}
	c.decodeBody(c.request("setBreakpoints", setBreakpointsArgs{
		Source:      sourceRef{Path: "stray.s"},
		Breakpoints: []sourceBreakpoint{{Line: 99}},
	}), &res)
	require.Len(t, res.Breakpoints, 1)
	require.False(t, res.Breakpoints[0].Verified)
	require.NotEmpty(t, res.Breakpoints[0].Message)

	// label breakpoint
	res = setBreakpointsResult{}
	c.decodeBody(c.request("setFunctionBreakpoints", setFunctionBreakpointsArgs{
		Breakpoints: []functionBreakpoint{{Name: "loop"}},
	}), &res)
	require.Len(t, res.Breakpoints, 1)
	require.True(t, res.Breakpoints[0].Verified)

	// the entry stop lands after the configuration
	c.request("configurationDone", nil)
	entry := stoppedBody{}
	c.decodeBody(c.event("stopped"), &entry)
	require.Equal(t, "entry", entry.Reason)

	// the run to the label, reported as a breakpoint hit
	cont := continueResult{}
	c.decodeBody(c.request("continue", continueArgs{ThreadId: oneThread}), &cont)
	require.True(t, cont.AllThreadsContinued)

	hit := stoppedBody{}
	c.decodeBody(c.event("stopped"), &hit)
	require.Equal(t, "breakpoint", hit.Reason)

	// the frame: pc, the source line, the disassembly label
	frames := stackTraceResult{}
	c.decodeBody(c.request("stackTrace", map[string]any{"threadId": oneThread}), &frames)
	require.Len(t, frames.StackFrames, 1)
	frame := frames.StackFrames[0]
	require.Equal(t, "401000: 20 00 80 d2", frame.Name)
	require.Equal(t, 2, frame.Line)
	require.Equal(t, "hello.s", frame.Source.Path)
	require.Equal(t, "0x401000", frame.InstructionPointerReference)

	// the one scope with the two registers in it
	scopes := scopesResult{}
	c.decodeBody(c.request("scopes", map[string]any{"frameId": oneFrame}), &scopes)
	require.Len(t, scopes.Scopes, 1)
	require.Equal(t, "Registers", scopes.Scopes[0].Name)
	require.Equal(t, registersRef, scopes.Scopes[0].VariablesReference)

	vars := variablesResult{}
	c.decodeBody(c.request("variables", variablesArgs{VariablesReference: registersRef}), &vars)
	require.Len(t, vars.Variables, 2)
	require.Equal(
		t,
		variable{Name: "r0", Value: "0x41", MemoryReference: "0x41"},
		vars.Variables[0],
	)
	require.Equal(
		t,
		variable{Name: "pc", Value: "0x401000", MemoryReference: "0x401000"},
		vars.Variables[1],
	)

	// a label in the watch window: the word at it
	eval := evaluateResult{}
	c.decodeBody(c.request("evaluate", evaluateArgs{Expression: "loop"}), &eval)
	require.Equal(t, "0x8", eval.Result)

	// the listing window: two instructions at the pc
	dis := disassembleResult{}
	c.decodeBody(c.request("disassemble", disassembleArgs{
		MemoryReference:  "0x401000",
		InstructionCount: 2,
	}), &dis)
	require.Len(t, dis.Instructions, 2)
	require.Equal(t, "0x401000", dis.Instructions[0].Address)
	require.Equal(t, "401000: 20 00 80 d2", dis.Instructions[0].Instruction)
	require.Equal(t, "0x401004", dis.Instructions[1].Address)

	// the memory window
	mem := readMemoryResult{}
	c.decodeBody(
		c.request("readMemory", readMemoryArgs{MemoryReference: "0x401000", Count: 4}),
		&mem,
	)
	require.Equal(t, "0x401000", mem.Address)
	require.Equal(t, []byte{0x20, 0x00, 0x80, 0xd2}, mem.Data)

	// disconnect tears the machine down (the teardown runs in the Serve
	// goroutine after its reply - hence the wait)
	c.request("disconnect", nil)
	require.Eventually(t, spy.closed.Load, time.Second, 5*time.Millisecond, "machine closed")
}

// TestServerStepsWhileRunning - pause and the stepping vocabulary:
// next reports a step, pause interrupts a continue the stub holds on
// the gate (the deterministic form of a slow target).
func TestServerStepsWhileRunning(t *testing.T) {
	lines := []session.Line{session.NewLine("hello.s", 2, 0x401000, 4)}

	gate := make(chan struct{})
	steps := []dialogStep{
		{expect: "s", reply: "T05thread:p1.1;"},
		{expect: "c", reply: "T05thread:p1.1;", gate: gate},
	}
	c, _ := newTestServer(t, steps, nil, lines)

	c.request("launch", launchArgs{Arch: "fake"})
	c.event("initialized")
	c.request("configurationDone", nil)
	c.event("stopped")

	c.request("next", nil)
	stopped := stoppedBody{}
	c.decodeBody(c.event("stopped"), &stopped)
	require.Equal(t, "step", stopped.Reason)

	c.request("continue", continueArgs{ThreadId: oneThread})

	pause := c.request("pause", nil)
	require.True(t, pause.Success)
	close(gate)

	paused := stoppedBody{}
	c.decodeBody(c.event("stopped"), &paused)
	require.Equal(t, "pause", paused.Reason)
}

// TestServerNeedsLaunch - the requests before the launch, and the
// ones that never touch the target, answer without a session.
func TestServerNeedsLaunch(t *testing.T) {
	editor, adapter := net.Pipe()

	srv, err := NewServer(
		adapter,
		adapter,
		func(args launchArgs, console io.Writer) (*Runtime, error) {
			return nil, errors.New("must not be called")
		},
	)
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve() }()
	t.Cleanup(func() { require.NoError(t, <-serveErr) })
	t.Cleanup(func() { require.NoError(t, editor.Close()) })

	c := newTestClient(t, editor)

	c.request("initialize", nil)

	failed := []string{"setBreakpoints", "setFunctionBreakpoints", "setInstructionBreakpoints",
		"configurationDone", "scopes", "variables", "evaluate",
		"disassemble", "readMemory"}
	for _, command := range failed {
		msg := c.request(command, nil)
		require.False(t, msg.Success, command)
		require.Contains(t, msg.Message, "no session", command)
	}

	// stackTrace never fails: the frames just stay empty (an error
	// response paints the editor's Frames panel red)
	frames := stackTraceResult{}
	c.decodeBody(c.request("stackTrace", map[string]any{"threadId": oneThread}), &frames)
	require.Equal(t, 0, frames.TotalFrames)

	// the target-free answers still work
	threads := threadsResult{}
	c.decodeBody(c.request("threads", nil), &threads)
	require.Equal(t, []thread{{Id: 1, Name: "cpu0"}}, threads.Threads)

	msg := c.request("attach", nil)
	require.False(t, msg.Success)

	msg = c.request("pause", nil)
	require.False(t, msg.Success)
	require.Contains(t, msg.Message, "not running")

	msg = c.request("frobnicate", nil)
	require.False(t, msg.Success)
	require.Contains(t, msg.Message, "unknown command")

	// the cleanup closes the editor side: Serve reads EOF and ends
	// cleanly (its verdict is checked in the cleanup)
}

// TestServerLaunchFailure - a launcher error is a failed launch
// response, the server stays alive for a retry.
func TestServerLaunchFailure(t *testing.T) {
	editor, adapter := net.Pipe()

	srv, err := NewServer(
		adapter,
		adapter,
		func(args launchArgs, console io.Writer) (*Runtime, error) {
			return nil, errors.New("no such arch")
		},
	)
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve() }()
	t.Cleanup(func() { require.NoError(t, <-serveErr) })
	t.Cleanup(func() { require.NoError(t, editor.Close()) })

	c := newTestClient(t, editor)

	c.request("initialize", nil)
	msg := c.request("launch", launchArgs{Arch: "fake"})
	require.False(t, msg.Success)
	require.Contains(t, msg.Message, "no such arch")

	threads := threadsResult{}
	c.decodeBody(c.request("threads", nil), &threads)
	require.Equal(t, "cpu0", threads.Threads[0].Name)
}

// TestServerExited - the W stop of the continue goroutine becomes the
// exited+terminated event pair.
func TestServerExited(t *testing.T) {
	steps := []dialogStep{
		{expect: "c", reply: "W00"},
	}
	c, _ := newTestServer(t, steps, nil, nil)

	c.request("launch", launchArgs{Arch: "fake"})
	c.event("initialized")
	c.request("configurationDone", nil)
	c.event("stopped")

	c.request("continue", continueArgs{ThreadId: oneThread})
	c.event("exited")
	c.event("terminated")
}

// TestServerPoweroffWithoutReply - the riscv-style poweroff: the stub
// cuts the connection instead of a farewell stop reply. The pending
// resume reports the end as a console line (no red ink), then the
// conversation terminates.
func TestServerPoweroffWithoutReply(t *testing.T) {
	steps := []dialogStep{
		{expect: "c", cut: true},
	}
	c, _ := newTestServer(t, steps, nil, nil)

	c.request("launch", launchArgs{Arch: "fake"})
	c.event("initialized")
	c.request("configurationDone", nil)
	c.event("stopped")

	c.request("continue", continueArgs{ThreadId: oneThread})

	// the death window: the editor may poll the frames right after the
	// continue response, while the resume goroutine is still dying with
	// the machine - the answer is a graceful empty, never an error
	frames := stackTraceResult{}
	c.decodeBody(c.request("stackTrace", map[string]any{"threadId": oneThread}), &frames)
	require.Equal(t, 0, frames.TotalFrames)

	finished := outputBody{}
	c.decodeBody(c.event("output"), &finished)
	require.Equal(t, "console", finished.Category)
	require.Contains(t, finished.Output, "the program finished")

	c.event("terminated")
}

// TestServerAfterFinished - the editor keeps asking after the
// terminated event (Zed polls stackTrace): the finished run answers
// gracefully - empty frames, unverified breakpoints, no errors, no red
// banner.
func TestServerAfterFinished(t *testing.T) {
	steps := []dialogStep{
		{expect: "c", cut: true}, // the machine vanishes mid-run
	}
	c, _ := newTestServer(t, steps, map[string]uint64{"loop": 0x401004}, nil)

	c.request("launch", launchArgs{Arch: "fake"})
	c.event("initialized")
	c.request("configurationDone", nil)
	c.event("stopped")

	c.request("continue", continueArgs{ThreadId: oneThread})
	c.event("output") // the console farewell
	c.event("terminated")

	frames := stackTraceResult{}
	c.decodeBody(c.request("stackTrace", map[string]any{"threadId": oneThread}), &frames)
	require.Equal(t, 0, frames.TotalFrames)
	require.Empty(t, frames.StackFrames)

	scopes := scopesResult{}
	c.decodeBody(c.request("scopes", map[string]any{"frameId": oneFrame}), &scopes)
	require.Empty(t, scopes.Scopes)

	res := setBreakpointsResult{}
	c.decodeBody(c.request("setBreakpoints", setBreakpointsArgs{
		Source:      sourceRef{Path: "hello.s"},
		Breakpoints: []sourceBreakpoint{{Line: 2}},
	}), &res)
	require.Len(t, res.Breakpoints, 1)
	require.False(t, res.Breakpoints[0].Verified)
	require.Contains(t, res.Breakpoints[0].Message, "finished")

	eval := evaluateResult{}
	c.decodeBody(c.request("evaluate", evaluateArgs{Expression: "loop"}), &eval)
	require.Contains(t, eval.Result, "finished")

	mem := readMemoryResult{}
	c.decodeBody(
		c.request("readMemory", readMemoryArgs{MemoryReference: "0x401000", Count: 4}),
		&mem,
	)
	require.Equal(t, 4, mem.UnreadableBytes)

	// the disconnect still ends cleanly
	c.request("disconnect", nil)
}
