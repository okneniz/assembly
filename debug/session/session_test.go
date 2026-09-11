package session

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/debug"
)

// fakeTarget is a two-register target for the engine tests: pc is
// register 1 (register 0 doubles as the only GPR shown in dumps).
type fakeTarget struct{}

func newFakeTarget() fakeTarget {
	return fakeTarget{}
}

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

// dialogStep is one expected request with its reply (the rsp fake
// without retransmission - the framing is covered there).
type dialogStep struct {
	expect string
	reply  string
}

// The wire codec in miniature for the dialog server (the rsp package
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

// serveDialog runs the server side of the conversation; the returned
// wait collects its verdict. New itself speaks two packets
// (qSupported, ?) - dialogs account for them.
func serveDialog(t *testing.T, steps []dialogStep) (net.Conn, func() error) {
	t.Helper()

	client, server := net.Pipe()
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	done := make(chan error, 1)
	go func() {
		defer close(done)
		r := bufio.NewReader(server)
		w := server
		fail := func(format string, args ...any) {
			done <- errors.Join(fmt.Errorf(format, args...), server.Close())
		}

		readRequest := func() (string, bool) {
			b, err := r.ReadByte()
			if err != nil || b != '$' {
				return "", false
			}

			pkt := []byte{'$'}
			for {
				b, err := r.ReadByte()
				if err != nil {
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

		for _, st := range steps {
			payload, ok := readRequest()
			if !ok || payload != st.expect {
				fail("dialog: request %q (ok=%v), want %q", payload, ok, st.expect)
				return
			}

			if _, werr := w.Write([]byte{'+'}); werr != nil {
				fail("dialog: ack write: %v", werr)
				return
			}

			if _, werr := w.Write(testEncode(st.reply)); werr != nil {
				fail("dialog: reply write: %v", werr)
				return
			}

			if b, err := r.ReadByte(); err != nil || b != '+' {
				fail("dialog: reply ack %#02x (%v)", b, err)
				return
			}
		}

		done <- server.Close()
	}()

	return client, func() error {
		require.NoError(t, client.Close())
		return <-done
	}
}

// newTestSession binds a session to the dialog's connection with the
// two-packet handshake of New at the head of the dialog.
func newTestSession(
	t *testing.T,
	steps []dialogStep,
	syms map[string]uint64,
	lines []Line,
) (*Session, func() error) {
	t.Helper()

	handshake := []dialogStep{
		{
			expect: "qSupported:multiprocess-;swbreak+;hwbreak+;xmlRegisters=aarch64,riscv:rv64,loongarch64,i386",
			reply:  "qXfer:features:read+",
		},
		{expect: "?", reply: "T05thread:p1.1;"},
	}

	conn, wait := serveDialog(t, append(handshake, steps...))
	s, err := New(conn, newFakeTarget(), syms, lines)
	require.NoError(t, err)

	return s, wait
}

func TestPC(t *testing.T) {
	s, wait := newTestSession(t, []dialogStep{
		{expect: "p1", reply: "00104000"},
	}, nil, nil)
	defer func() { require.NoError(t, wait()) }()

	pc, err := s.PC()
	require.NoError(t, err)
	require.Equal(t, uint64(0x401000), pc)
}

func TestBreakAtSymbolAndAddress(t *testing.T) {
	syms := map[string]uint64{"start": 0x401000}

	s, wait := newTestSession(t, []dialogStep{
		{expect: "Z0,401000,4", reply: "OK"},
		{expect: "Z0,401008,4", reply: "OK"},
	}, syms, nil)
	defer func() { require.NoError(t, wait()) }()

	addr, err := s.BreakAt("start")
	require.NoError(t, err)
	require.Equal(t, uint64(0x401000), addr)

	addr, err = s.BreakAt("0x401008")
	require.NoError(t, err)
	require.Equal(t, uint64(0x401008), addr)

	_, err = s.BreakAt("nowhere")
	require.Error(t, err)
}

func TestRegs(t *testing.T) {
	// the 'g' block of two 8-byte registers: r0 = 0x41, pc = 0x401000
	s, wait := newTestSession(t, []dialogStep{
		{expect: "g", reply: "4100000000000000" + "0010400000000000"},
	}, nil, nil)
	defer func() { require.NoError(t, wait()) }()

	regs, err := s.Regs()
	require.NoError(t, err)
	require.Equal(t, []RegValue{
		NewRegValue("r0", 0x41),
		NewRegValue("pc", 0x401000),
	}, regs)
}

func TestReadAndDisasm(t *testing.T) {
	s, wait := newTestSession(t, []dialogStep{
		{expect: "m401000,4", reply: "200080d2"},
	}, nil, nil)
	defer func() { require.NoError(t, wait()) }()

	lines, err := s.Disasm(0x401000, 1)
	require.NoError(t, err)
	require.Equal(t, []string{"401000: 20 00 80 d2"}, lines)
}

func TestLineAt(t *testing.T) {
	lines := []Line{
		NewLine("a.s", 2, 0x1000, 4),
		NewLine("a.s", 4, 0x1004, 8),
		NewLine("a.s", 6, 0x1010, 4),
	}
	s, wait := newTestSession(t, nil, nil, lines)
	defer func() { require.NoError(t, wait()) }()

	cases := []struct {
		name   string
		addr   uint64
		want   Line
		wantOK bool
	}{
		{"exact", 0x1004, lines[1], true},
		{"mid-line (8-byte line)", 0x1008, lines[1], true},
		{"after the last", 0x1014, lines[2], true},
		{"below the first", 0xffc, Line{}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := s.LineAt(c.addr)
			require.Equal(t, c.wantOK, ok)
			require.Equal(t, c.want, got)
		})
	}
}
