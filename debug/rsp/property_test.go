package rsp

import (
	"testing"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
)

// payloadAlphabet is the protocol's working alphabet: commands and
// replies are built of it. The framing bytes ($, #) and the escaping
// byte (}) are excluded - the v1 codec does not escape payload
// content, so they are not legal inside a payload (wire-level
// escaping is future work, not a property of this codec).
const payloadAlphabet = "0123456789abcdefABCDEFmxcsgpHqvMX,:;?.=+-"

// TestPacketRoundTrip - the codec law: what was encoded, parses back.
func TestPacketRoundTrip(t *testing.T) {
	rnd := arb.Rnd(42)
	payloads := ohsnap.ArbitraryString(rnd, payloadAlphabet, 0, 40)
	ohsnap.Check(t, 12000, payloads, func(p string) bool {
		got, ok := parsePacket(encodePacket(p))
		return ok && got == p
	})
}

// TestPacketSingleFlipDetected - no single-bit flip of a valid packet
// (framing, payload, or checksum byte) can forge different content:
// the flip either breaks parsing or parses back the SAME payload. The
// payload case is the strong one: one flipped byte always changes the
// sum modulo 256, so the checksum cannot accidentally match again. The
// one harmless survivor is the case of a hex letter in the checksum
// digits ("9a" -> "9A"): the value is the same, the packet is
// equivalent - found by this very property in its first, too-strict
// wording.
func TestPacketSingleFlipDetected(t *testing.T) {
	rnd := arb.Rnd(7)
	payloads := ohsnap.ArbitraryString(rnd, payloadAlphabet, 0, 40)
	ohsnap.Check(t, 12000, payloads, func(p string) bool {
		pkt := encodePacket(p)
		for i := range pkt {
			for bit := range 8 {
				corrupt := append([]byte{}, pkt...)
				corrupt[i] ^= 1 << bit
				if got, ok := parsePacket(corrupt); ok && got != p {
					return false
				}
			}
		}

		return true
	})
}

// TestPacketParseTotalOnGarbage - the parser is total: arbitrary bytes
// never panic it, the answer is always (payload, ok).
func TestPacketParseTotalOnGarbage(t *testing.T) {
	rnd := arb.Rnd(9)
	raw := ohsnap.ArbitrarySlice(rnd, ohsnap.ArbitraryByte(rnd, 0, 255), 0, 32)
	ohsnap.Check(t, 12000, raw, func(b []byte) bool {
		parsePacket(b) // the call itself is the property: no panic
		return true
	})
}
