package rsp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodePacket(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		want    string
	}{
		{"empty payload", "", "$#00"},
		{"halt reason", "?", "$?#3f"},
		{"continue", "c", "$c#63"},
		{"read memory", "m1000,4", "$m1000,4#8e"},
		{"register block", "g", "$g#67"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, string(encodePacket(c.payload)))
		})
	}
}

func TestParsePacket(t *testing.T) {
	cases := []struct {
		name   string
		wire   string
		want   string
		wantOK bool
	}{
		{name: "valid", wire: "$c#63", want: "c", wantOK: true},
		{name: "valid empty", wire: "$#00", want: "", wantOK: true},
		{name: "with payload", wire: "$m1000,4#8e", want: "m1000,4", wantOK: true},
		{name: "bad checksum", wire: "$c#64"},
		{name: "no dollar", wire: "c#63"},
		{name: "no hash", wire: "$c63"},
		{name: "too short", wire: "$#"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			payload, ok := parsePacket([]byte(c.wire))
			require.Equal(t, c.wantOK, ok)
			require.Equal(t, c.want, payload)
		})
	}
}

func TestChecksumWraps(t *testing.T) {
	// the sum of "mffff,fff..." overflows the byte and wraps: the
	// checksum is the low byte
	require.Equal(t, byte(0x00), checksum("\xff\x01\x00"))
	require.Equal(t, byte(0x00), checksum(""))
}
