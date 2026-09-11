package dap

import (
	"bytes"
	"io"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
)

func TestFramingWriteFormat(t *testing.T) {
	cases := []struct {
		name string
		body string
		wire string
	}{
		{"empty body", "", "Content-Length: 0\r\n\r\n"},
		{"json body", `{"a":1}`, "Content-Length: 7\r\n\r\n" + `{"a":1}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			f := newFraming(&buf, &buf)
			require.NoError(t, f.write([]byte(c.body)))
			require.Equal(t, c.wire, buf.String())

			got, err := f.read()
			require.NoError(t, err)
			require.Equal(t, c.body, string(got))
		})
	}
}

func TestFramingReadAcrossEnvelopes(t *testing.T) {
	var buf bytes.Buffer
	f := newFraming(&buf, &buf)
	for _, body := range []string{`{"one":1}`, `{"two":22}`, "x"} {
		require.NoError(t, f.write([]byte(body)))
	}

	for _, want := range []string{`{"one":1}`, `{"two":22}`, "x"} {
		got, err := f.read()
		require.NoError(t, err)
		require.Equal(t, want, string(got))
	}
}

func TestFramingReadErrors(t *testing.T) {
	cases := []struct {
		name string
		wire string
		want string
	}{
		{"no length header", "Some-Header: x\r\n\r\n", "no Content-Length"},
		{"bad length", "Content-Length: x\r\n\r\n", "bad Content-Length"},
		{"negative length", "Content-Length: -1\r\n\r\n", "bad Content-Length"},
		{"header without colon", "garbage\r\n\r\n", "malformed header"},
		{"truncated body", "Content-Length: 10\r\n\r\nshort", "body"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newFraming(strings.NewReader(c.wire), io.Discard)
			_, err := f.read()
			require.ErrorContains(t, err, c.want)
		})
	}
}

// TestFramingRoundTrip - the framing law: every byte sequence written
// as one envelope comes back as exactly that envelope.
func TestFramingRoundTrip(t *testing.T) {
	rnd := arb.Rnd(42)
	bodies := ohsnap.ArbitrarySlice(rnd, ohsnap.ArbitraryByte(rnd, 0, 255), 0, 64)
	ohsnap.Check(t, 12000, bodies, func(b []byte) bool {
		var buf bytes.Buffer
		f := newFraming(&buf, &buf)
		if err := f.write(b); err != nil {
			return false
		}

		got, err := f.read()
		return err == nil && bytes.Equal(b, got)
	})
}
