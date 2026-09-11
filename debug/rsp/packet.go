// Package rsp is a client of the GDB Remote Serial Protocol: the wire
// the qemu gdbstub (and every other gdb-compatible executor) speaks.
// The protocol is arch-neutral - register numbers and layouts come from
// the target description, not from here.
//
// Packet form: "$" payload "#" checksum, the checksum being the low
// byte of the payload's sum as two hex digits. The receiver
// acknowledges with "+" (accept) or "-" (retransmit). An interrupt is
// the single raw byte 0x03, outside any packet.
package rsp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// checksum is the packet checksum: the sum of the payload bytes modulo
// 256, on the wire as two hex digits.
func checksum(payload string) byte {
	var sum byte
	for i := range len(payload) {
		sum += payload[i]
	}

	return sum
}

// encodePacket is the wire form of a payload: $payload#cs.
func encodePacket(payload string) []byte {
	return []byte(fmt.Sprintf("$%s#%02x", payload, checksum(payload)))
}

// parsePacket is the payload of one received wire packet; ok is false
// when the framing or the checksum does not hold.
func parsePacket(pkt []byte) (string, bool) {
	if len(pkt) < 4 || pkt[0] != '$' || pkt[len(pkt)-3] != '#' {
		return "", false
	}

	payload := string(pkt[1 : len(pkt)-3])
	want, err := strconv.ParseUint(string(pkt[len(pkt)-2:]), 16, 8)
	if err != nil || checksum(payload) != byte(want) {
		return "", false
	}

	return payload, true
}

// replyPayload classifies a raw reply payload: a normal payload, an
// error, or an empty packet (command not supported by this target).
func replyPayload(payload string) (string, error) {
	if payload == "" {
		return "", errors.New("assembly/rsp: empty reply (command not supported)")
	}

	if strings.HasPrefix(payload, "E") {
		return "", newReplyError(payload[1:])
	}

	return payload, nil
}
