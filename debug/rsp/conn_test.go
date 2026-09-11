package rsp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

// fakeStep is one dialog turn of the fake target: the request payload
// it must receive, and the payload it answers with. nacks is how many
// "-" answers precede the "+" of a damaged request (checksum testing).
type fakeStep struct {
	expect string
	reply  string
	nacks  int
}

// runFakeTarget speaks the server side of the protocol over conn. Per
// step: it reads the request packet, answers each expected "-" with a
// re-read of the resent packet (the retransmit handshake), then "+" and
// its own reply packet. The goroutine ends when the dialog or the
// connection ends; the return value reports mismatches.
func runFakeTarget(conn net.Conn, steps []fakeStep) <-chan error {
	done := make(chan error, 1)
	go func() {
		defer close(done)
		r := bufio.NewReader(conn)
		w := conn
		fail := func(format string, args ...any) {
			done <- errors.Join(fmt.Errorf(format, args...), conn.Close())
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

			payload, ok := parsePacket(pkt)
			return payload, ok
		}

		for _, st := range steps {
			payload, ok := readRequest()
			if !ok {
				fail("fake: bad request packet")
				return
			}

			if payload != st.expect {
				fail("fake: request %q, want %q", payload, st.expect)
				return
			}

			for range st.nacks {
				if _, werr := w.Write([]byte{'-'}); werr != nil {
					fail("fake: nack write: %v", werr)
					return
				}

				resend, ok := readRequest()
				if !ok || resend != st.expect {
					fail("fake: bad resend %q", resend)
					return
				}
			}

			if _, werr := w.Write([]byte{'+'}); werr != nil {
				fail("fake: ack write: %v", werr)
				return
			}

			if _, werr := w.Write(encodePacket(st.reply)); werr != nil {
				fail("fake: reply write: %v", werr)
				return
			}

			// the client acks our reply before its next request - drain
			// it, or the synchronous pipe sees that '+' hit a closed end
			if b, err := r.ReadByte(); err != nil || b != '+' {
				fail("fake: reply ack %#02x (%v)", b, err)
				return
			}
		}

		done <- conn.Close()
	}()

	return done
}

// dialFake is a client conversation with a fake target bound to the
// dialog; wait collects the fake's verdict after the test's call.
func dialFake(t *testing.T, steps []fakeStep) (*Conn, func() error) {
	t.Helper()

	client, server := net.Pipe()
	t.Cleanup(func() { require.NoError(t, client.Close()) })

	done := runFakeTarget(server, steps)
	return NewConn(client), func() error {
		require.NoError(t, client.Close())
		return <-done
	}
}

func TestRoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		steps []fakeStep
		cmd   string
		want  string
	}{
		{
			name:  "halt reason",
			steps: []fakeStep{{expect: "?", reply: "T05thread:p1.1;"}},
			cmd:   "?",
			want:  "T05thread:p1.1;",
		},
		{
			name:  "retransmit after nack",
			steps: []fakeStep{{expect: "g", reply: "1122334455667788", nacks: 1}},
			cmd:   "g",
			want:  "1122334455667788",
		},
		{
			name:  "double nack",
			steps: []fakeStep{{expect: "m0,4", reply: "deadbeef", nacks: 2}},
			cmd:   "m0,4",
			want:  "deadbeef",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, wait := dialFake(t, c.steps)
			got, err := conn.roundTrip(c.cmd)
			require.NoError(t, err)
			require.Equal(t, c.want, got)
			require.NoError(t, wait())
		})
	}
}

func TestRoundTripRequestMismatch(t *testing.T) {
	// the fake expects another request: the fake errors, the client
	// sees a broken conversation
	conn, wait := dialFake(t, []fakeStep{{expect: "?", reply: "S05"}})
	_, err := conn.roundTrip("c")
	require.Error(t, err)
	require.Error(t, wait())
}

func TestRoundTripErrorReply(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "mdeadbeef,4", reply: "E14"}})
	_, err := conn.roundTrip("mdeadbeef,4")
	require.ErrorContains(t, err, "target error 14")
	require.NoError(t, wait())
}

func TestInterruptByte(t *testing.T) {
	// the interrupt is one raw byte on the wire, outside any framing;
	// the pipe is synchronous, so the reader runs while Interrupt writes
	client, server := net.Pipe()
	t.Cleanup(func() {
		require.NoError(t, client.Close())
		require.NoError(t, server.Close())
	})

	type readByte struct {
		b   byte
		err error
	}

	got := make(chan readByte, 1)
	go func() {
		b := make([]byte, 1)
		_, err := server.Read(b)
		got <- readByte{b: b[0], err: err}
	}()

	conn := NewConn(client)
	require.NoError(t, conn.Interrupt())

	r := <-got
	require.NoError(t, r.err)
	require.Equal(t, byte(0x03), r.b)
}
