package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	darm64 "github.com/okneniz/assembly/debug/arm64"
	"github.com/okneniz/assembly/debug/dap"
	"github.com/okneniz/assembly/debug/session"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/prog"
	aprog "github.com/okneniz/assembly/prog/arm64"
)

// dapMsg is the editor-side view of one DAP message of the e2e gate.
type dapMsg struct {
	Type    string          `json:"type"`
	Event   string          `json:"event"`
	ReqSeq  int             `json:"request_seq"`
	Success bool            `json:"success"`
	Command string          `json:"command"`
	Message string          `json:"message"`
	Body    json.RawMessage `json:"body"`
}

// dapEditor is the editor half of the gate conversation: sequential
// requests over the pipe, the events collected on the way.
type dapEditor struct {
	t   *testing.T
	r   *bufio.Reader
	c   net.Conn
	seq int

	events []dapMsg
}

func newDapEditor(t *testing.T, conn net.Conn) *dapEditor {
	t.Helper()

	return &dapEditor{t: t, r: bufio.NewReader(conn), c: conn}
}

func (e *dapEditor) request(command string, args any) dapMsg {
	e.t.Helper()
	e.seq++
	body, err := json.Marshal(map[string]any{
		"seq":       e.seq,
		"type":      "request",
		"command":   command,
		"arguments": args,
	})
	require.NoError(e.t, err)

	head := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	_, err = e.c.Write([]byte(head + string(body)))
	require.NoError(e.t, err)

	for {
		msg := e.read()
		if msg.Type == "response" {
			require.Equal(e.t, e.seq, msg.ReqSeq, "reply to another request")
			return msg
		}

		e.events = append(e.events, msg)
	}
}

// event returns the next event of that name (the buffered ones first);
// the executor's console output ("output") is skipped on the way - the
// machine may mutter as it boots.
func (e *dapEditor) event(name string) dapMsg {
	e.t.Helper()

	for {
		for i, ev := range e.events {
			if ev.Event == name {
				e.events = append(e.events[:i], e.events[i+1:]...)
				return ev
			}
		}

		msg := e.read()
		require.Equal(e.t, "event", msg.Type, "want event %q", name)
		if msg.Event == "output" {
			continue
		}

		require.Equal(e.t, name, msg.Event)
		return msg
	}
}

func (e *dapEditor) read() dapMsg {
	e.t.Helper()

	require.NoError(e.t, e.c.SetReadDeadline(time.Now().Add(20*time.Second)))
	length := -1
	for {
		line, err := e.r.ReadString('\n')
		require.NoError(e.t, err)

		if line == "\r\n" || line == "\n" {
			break
		}

		if rest, ok := cutPrefixStr(line, "Content-Length: "); ok {
			length, err = strconv.Atoi(trimEOL(rest))
			require.NoError(e.t, err)
		}
	}

	require.NotEqual(e.t, -1, length, "no Content-Length header")

	body := make([]byte, length)
	_, err := e.r.Read(body)
	require.NoError(e.t, err)

	var msg dapMsg
	require.NoError(e.t, json.Unmarshal(body, &msg))
	return msg
}

func (e *dapEditor) body(msg dapMsg, out any) {
	e.t.Helper()
	require.NoError(e.t, json.Unmarshal(msg.Body, out))
}

func cutPrefixStr(s, prefix string) (string, bool) {
	if len(s) < len(prefix) || s[:len(prefix)] != prefix {
		return "", false
	}

	return s[len(prefix):], true
}

func trimEOL(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}

	if len(s) > 0 && s[len(s)-1] == '\r' {
		s = s[:len(s)-1]
	}

	return s
}

// The typed views of the reply bodies (the untyped json map would need
// ten type assertions; the gate reads typed or not at all).
type dapBreakpointList struct {
	Breakpoints []struct {
		Verified bool `json:"verified"`
	} `json:"breakpoints"`
}

type dapStopped struct {
	Reason string `json:"reason"`
}

type dapStack struct {
	Frames []struct {
		Column int    `json:"column"`
		Line   int    `json:"line"`
		Name   string `json:"name"`
		Source *struct {
			Path string `json:"path"`
		} `json:"source"`
		IP string `json:"instructionPointerReference"`
	} `json:"stackFrames"`
}

type dapScopes struct {
	Scopes []struct {
		Name string `json:"name"`
	} `json:"scopes"`
}

type dapVariables struct {
	Variables []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"variables"`
}

type dapDisassembly struct {
	Instructions []struct {
		Address     string `json:"address"`
		Instruction string `json:"instruction"`
	} `json:"instructions"`
}

// TestDapEndToEnd - the editor gate: the full DAP conversation over
// the real qemu executor (arm64, the vm hello): launch, a label
// breakpoint, the entry stop, the run to the breakpoint, the frame
// with the source line, the register scope, the disassembly window -
// the same assertions as the session gates, one protocol up.
func TestDapEndToEnd(t *testing.T) {
	tgt := darm64.NewTarget()
	if _, err := exec.LookPath(tgt.QemuBinary()); err != nil {
		t.Skipf("%s not on PATH (brew install qemu)", tgt.QemuBinary())
	}

	src := "examples/hello-asm/hello-arm-vm.s"
	readFile(t, src) // skip when the fixture is absent

	editor, adapter := net.Pipe()

	srv, err := dap.NewServer(adapter, adapter, dap.Boot)
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve() }()
	t.Cleanup(func() { require.NoError(t, <-serveErr) })
	t.Cleanup(func() { require.NoError(t, editor.Close()) })

	e := newDapEditor(t, editor)

	caps := map[string]any{}
	e.body(e.request("initialize", nil), &caps)
	require.Equal(t, true, caps["supportsDisassembleRequest"])

	msg := e.request("launch", map[string]any{"arch": "arm64", "source": src})
	require.True(t, msg.Success, msg.Message)
	e.event("initialized")

	bps := dapBreakpointList{}
	e.body(e.request("setFunctionBreakpoints", map[string]any{
		"breakpoints": []map[string]any{{"name": "done"}},
	}), &bps)
	require.Len(t, bps.Breakpoints, 1)
	require.True(t, bps.Breakpoints[0].Verified)

	msg = e.request("configurationDone", nil)
	require.True(t, msg.Success, msg.Message)

	stopped := dapStopped{}
	e.body(e.event("stopped"), &stopped)
	require.Equal(t, "entry", stopped.Reason)

	msg = e.request("continue", map[string]any{"threadId": 1})
	require.True(t, msg.Success, msg.Message)

	e.body(e.event("stopped"), &stopped)
	require.Equal(t, "breakpoint", stopped.Reason)

	stack := dapStack{}
	e.body(e.request("stackTrace", map[string]any{"threadId": 1}), &stack)
	require.Len(t, stack.Frames, 1)
	frame := stack.Frames[0]
	require.Equal(t, 1, frame.Column)
	require.NotEmpty(t, frame.Name)
	require.NotNil(t, frame.Source)
	require.Equal(t, src, frame.Source.Path)

	scopes := dapScopes{}
	e.body(e.request("scopes", map[string]any{"frameId": 1}), &scopes)
	require.Len(t, scopes.Scopes, 1)

	vars := dapVariables{}
	e.body(e.request("variables", map[string]any{"variablesReference": 1}), &vars)
	require.NotEmpty(t, vars.Variables)

	pcSeen := false
	for _, r := range vars.Variables {
		if r.Name == "pc" {
			pcSeen = true
			require.NotEmpty(t, r.Value)
		}
	}

	require.True(t, pcSeen, "pc in the register scope")

	dis := dapDisassembly{}
	e.body(e.request("disassemble", map[string]any{
		"memoryReference":  frame.IP,
		"instructionCount": 2,
	}), &dis)
	require.Len(t, dis.Instructions, 2)

	msg = e.request("disconnect", nil)
	require.True(t, msg.Success, msg.Message)
}

// TestDapProg - the prog path over the editor protocol: a Go-written
// program debugged through dap.NewImageLauncher (the embedded-server
// form a prog example uses in its -debug mode). The frame carries the
// GO source line - a constant resolver reports a synthetic position
// for the first instruction, the injected caller resolver reports this
// test file for the rest; a breakpoint on the Go line resolves through
// the same line map.
func TestDapProg(t *testing.T) {
	tgt := darm64.NewTarget()
	if _, err := exec.LookPath(tgt.QemuBinary()); err != nil {
		t.Skipf("%s not on PATH (brew install qemu)", tgt.QemuBinary())
	}

	p := aprog.New().
		WithPos(func() prog.Pos { return prog.NewPos("synthetic.go", 42) }).
		Label("start").
		Mov(aprog.X0, 0x41).
		WithPos(callerPos).
		Label("loop").
		B("loop").
		Entry("start")

	bin, buildErrs := p.Build()
	require.Empty(t, buildErrs)

	res := bin.Assemble(0x40100000)
	require.Empty(t, res.Errs)

	lines := make([]session.Line, 0, len(res.Lines))
	for _, e := range res.Lines {
		lines = append(lines, session.NewLine(e.Pos.File, e.Pos.Line, e.Addr, e.Size))
	}

	img, err := file.WriteELF(
		file.EM_AARCH64,
		0,
		0x40100000,
		res.Syms["start"],
		[]file.Section{
			*file.NewSection(".text", "", 0x40100000, 0, uint64(len(res.Code)), res.Code),
		},
	)
	require.NoError(t, err)

	editor, adapter := net.Pipe()

	launcher := dap.NewImageLauncher(tgt, img, res.Syms, lines)
	srv, err := dap.NewServer(adapter, adapter, launcher)
	require.NoError(t, err)

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve() }()
	t.Cleanup(func() { require.NoError(t, <-serveErr) })
	t.Cleanup(func() { require.NoError(t, editor.Close()) })

	e := newDapEditor(t, editor)

	caps := map[string]any{}
	e.body(e.request("initialize", nil), &caps)
	require.Equal(t, true, caps["supportsConfigurationDoneRequest"])

	msg := e.request("launch", map[string]any{"arch": "arm64"})
	require.True(t, msg.Success, msg.Message)
	e.event("initialized")

	// a breakpoint on the synthetic Go line resolves through the line map
	bps := dapBreakpointList{}
	e.body(e.request("setBreakpoints", map[string]any{
		"source":      map[string]any{"path": "synthetic.go"},
		"breakpoints": []map[string]any{{"line": 42}},
	}), &bps)
	require.Len(t, bps.Breakpoints, 1)
	require.True(t, bps.Breakpoints[0].Verified)

	msg = e.request("configurationDone", nil)
	require.True(t, msg.Success, msg.Message)

	stopped := dapStopped{}
	e.body(e.event("stopped"), &stopped)
	require.Equal(t, "entry", stopped.Reason)

	// the frame sits on the synthetic Go line
	stack := dapStack{}
	e.body(e.request("stackTrace", map[string]any{"threadId": 1}), &stack)
	require.Len(t, stack.Frames, 1)
	frame := stack.Frames[0]
	require.Equal(t, 42, frame.Line)
	require.NotNil(t, frame.Source)
	require.Equal(t, "synthetic.go", frame.Source.Path)

	// one step lands on the caller-resolved line: this test file
	msg = e.request("next", nil)
	require.True(t, msg.Success, msg.Message)
	e.body(e.event("stopped"), &stopped)
	require.Equal(t, "step", stopped.Reason)

	stack = dapStack{}
	e.body(e.request("stackTrace", map[string]any{"threadId": 1}), &stack)
	require.Len(t, stack.Frames, 1)
	frame = stack.Frames[0]
	require.Positive(t, frame.Line)
	require.NotNil(t, frame.Source)
	require.Contains(t, frame.Source.Path, "tests/dap_test.go")

	// the register scope serves the stepped state
	vars := dapVariables{}
	e.body(e.request("variables", map[string]any{"variablesReference": 1}), &vars)
	x0 := uint64(0)
	for _, r := range vars.Variables {
		if r.Name == "x0" {
			x0 = 1
		}
	}

	require.Equal(t, uint64(1), x0, "x0 in the register scope")

	msg = e.request("disconnect", nil)
	require.True(t, msg.Success, msg.Message)
}
