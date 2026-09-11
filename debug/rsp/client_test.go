package rsp

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// hexOf is the wire form of qXfer data: the payload as hex.
func hexOf(s string) string {
	return hex.EncodeToString([]byte(s))
}

func TestSupported(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{
		expect: ourFeatures,
		reply:  "PacketSize=1000;qXfer:features:read+;multiprocess-;swbreak+",
	}})
	defer func() { require.NoError(t, wait()) }()

	features, err := NewClient(conn).Supported()
	require.NoError(t, err)
	require.Equal(t, map[string]string{
		"PacketSize":          "1000",
		"qXfer:features:read": "+",
		"multiprocess":        "-",
		"swbreak":             "+",
	}, features)
}

func TestTargetXML(t *testing.T) {
	// two chunks: 'm' continues, 'l' closes; the payload is hex
	conn, wait := dialFake(t, []fakeStep{
		{
			expect: "qXfer:features:read:target.xml:0,fff",
			reply:  "m" + hexOf("<target version=\"1.0\"><architecture>aarch64</architecture>"),
		},
		{
			expect: "qXfer:features:read:target.xml:fff,fff",
			reply:  "l" + hexOf("</target>"),
		},
	})
	defer func() { require.NoError(t, wait()) }()

	xml, err := NewClient(conn).TargetXML()
	require.NoError(t, err)
	require.Equal(t, `<target version="1.0"><architecture>aarch64</architecture></target>`, xml)
}

func TestHaltReason(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "?", reply: "T05thread:p1.1;core:0;"}})
	defer func() { require.NoError(t, wait()) }()

	stop, err := NewClient(conn).HaltReason()
	require.NoError(t, err)
	require.Equal(t, NewStopReply('T', 5, map[string]string{
		"thread": "p1.1",
		"core":   "0",
	}), stop)
}

func TestSelectThread(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "Hg1", reply: "OK"}})
	defer func() { require.NoError(t, wait()) }()

	require.NoError(t, NewClient(conn).SelectThread(1))
}

func TestReadRegs(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "g", reply: "0102030405060708"}})
	defer func() { require.NoError(t, wait()) }()

	regs, err := NewClient(conn).ReadRegs()
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8}, regs)
}

func TestReadReg(t *testing.T) {
	cases := []struct {
		name  string
		reg   int
		reply string
		want  uint64
	}{
		// little-endian on the wire: the first pair is the low byte
		{"pc", 32, "00104000", 0x401000},
		{"low byte only", 0, "7f", 0x7f},
		// a wide register truncates to its low half: the FIRST eight
		// bytes (little-endian: the low doubleword)
		{"wide register truncated", 31, "0102030405060708090a0b0c", 0x0807060504030201},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, wait := dialFake(t, []fakeStep{{
				expect: fmt.Sprintf("p%x", c.reg),
				reply:  c.reply,
			}})
			defer func() { require.NoError(t, wait()) }()

			got, err := NewClient(conn).ReadReg(c.reg)
			require.NoError(t, err)
			require.Equal(t, c.want, got)
		})
	}
}

func TestWriteReg(t *testing.T) {
	// the value goes out as eight little-endian bytes
	conn, wait := dialFake(t, []fakeStep{{
		expect: "P20=0010400000000000",
		reply:  "OK",
	}})
	defer func() { require.NoError(t, wait()) }()

	require.NoError(t, NewClient(conn).WriteReg(0x20, 0x401000))
}

func TestReadMem(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{
		expect: "m401000,4",
		reply:  "deadbeef",
	}})
	defer func() { require.NoError(t, wait()) }()

	mem, err := NewClient(conn).ReadMem(0x401000, 4)
	require.NoError(t, err)
	require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, mem)
}

func TestWriteMem(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{
		expect: "M401000,3:010203",
		reply:  "OK",
	}})
	defer func() { require.NoError(t, wait()) }()

	require.NoError(t, NewClient(conn).WriteMem(0x401000, []byte{1, 2, 3}))
}

func TestBreakpoints(t *testing.T) {
	cases := []struct {
		name   string
		set    bool
		expect string
	}{
		{name: "set", set: true, expect: "Z0,401000,4"},
		{name: "clear", set: false, expect: "z0,401000,4"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, wait := dialFake(t, []fakeStep{{expect: c.expect, reply: "OK"}})
			defer func() { require.NoError(t, wait()) }()

			cl := NewClient(conn)
			var err error
			if c.set {
				err = cl.SetBreak(0x401000, 4)
			} else {
				err = cl.ClearBreak(0x401000, 4)
			}

			require.NoError(t, err)
		})
	}
}

func TestContinueStopReply(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "c", reply: "T05thread:p1.1;"}})
	defer func() { require.NoError(t, wait()) }()

	stop, err := NewClient(conn).Continue()
	require.NoError(t, err)
	require.Equal(t, byte('T'), stop.Kind)
}

func TestContinueExit(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "c", reply: "W00"}})
	defer func() { require.NoError(t, wait()) }()

	stop, err := NewClient(conn).Continue()
	require.NoError(t, err)
	require.True(t, stop.Exited())
}

func TestStep(t *testing.T) {
	conn, wait := dialFake(t, []fakeStep{{expect: "s", reply: "S07"}})
	defer func() { require.NoError(t, wait()) }()

	stop, err := NewClient(conn).Step()
	require.NoError(t, err)
	require.Equal(t, NewStopReply('S', 7, map[string]string{}), stop)
}
