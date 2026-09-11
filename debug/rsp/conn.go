package rsp

import (
	"bufio"
	"fmt"
	"io"
)

// Conn is one framed RSP conversation over a byte transport: a full
// request/response round trip with acknowledgment, checksum validation,
// and retransmission on a negative ack or a corrupt packet. The
// conversation is synchronous - one command in flight at a time, which
// is all the debugging flow needs (the executor stops while we ask).
type Conn struct {
	r    *bufio.Reader
	w    io.Writer
	rerr error // sticky read error: a closed or broken transport
}

// NewConn is a conversation over an established transport (a TCP
// connection to the gdbstub, a pipe to a fake target in tests).
func NewConn(rw io.ReadWriter) *Conn {
	return &Conn{
		r: bufio.NewReader(rw),
		w: rw,
	}
}

// maxRetransmits bounds the resend attempts of one round trip: a target
// that keeps nacking or corrupting its checksums is broken, not slow.
const maxRetransmits = 3

// Interrupt asks a running target to stop: the raw 0x03 byte, the only
// protocol element outside the packet framing. The stop report itself
// arrives as the reply of the pending Continue/Step.
func (c *Conn) Interrupt() error {
	if _, err := c.w.Write([]byte{0x03}); err != nil {
		return fmt.Errorf("assembly/rsp: interrupt: %w", err)
	}

	return nil
}

// roundTrip sends one command and returns the payload of the reply.
// It consumes the target's ack of the request, validates the reply's
// checksum (acking it with "+" on success), and resends the request on
// a nack ("-") or a bad checksum. Interrupt bytes (0x03) echoed inside
// the stream are skipped: they are addressed to the target, not us.
func (c *Conn) roundTrip(cmd string) (string, error) {
	pkt := encodePacket(cmd)
	for attempt := 0; attempt <= maxRetransmits; attempt++ {
		if _, err := c.w.Write(pkt); err != nil {
			return "", fmt.Errorf("assembly/rsp: write: %w", err)
		}

		switch ack := c.readByte(); ack {
		case '+': // accepted - the reply follows
		case '-': // rejected - resend
			continue
		default:
			return "", c.transportError(fmt.Sprintf("unexpected ack byte %#02x", ack))
		}

		payload, ok := c.readPacket()
		if !ok {
			if c.rerr != nil {
				// the target died mid-reply (a powered-off machine): no
				// point resending
				return "", c.transportError("no reply")
			}

			// a corrupt reply leaves the stream position undefined to
			// us but not to the target: nack it and resend the request
			if _, err := c.w.Write([]byte{'-'}); err != nil {
				return "", fmt.Errorf("assembly/rsp: nack: %w", err)
			}

			continue
		}

		if _, err := c.w.Write([]byte{'+'}); err != nil {
			return "", fmt.Errorf("assembly/rsp: ack: %w", err)
		}

		return replyPayload(payload)
	}

	return "", fmt.Errorf("assembly/rsp: no valid reply after %d retransmits", maxRetransmits)
}

// readByte is the next stream byte; the 0x03 of a concurrent interrupt
// is skipped (it is our own byte echoed by a line discipline or a
// target artifact - never part of a reply). A read error sticks in
// rerr and reads as zero.
func (c *Conn) readByte() byte {
	for {
		b, err := c.r.ReadByte()
		if err != nil {
			c.rerr = err
			return 0
		}

		if b != 0x03 {
			return b
		}
	}
}

// transportError is a message with the sticky transport failure
// attached (a closed connection reads as an empty stream otherwise).
func (c *Conn) transportError(msg string) error {
	if c.rerr != nil {
		return fmt.Errorf("assembly/rsp: %s: connection closed: %w", msg, c.rerr)
	}

	return fmt.Errorf("assembly/rsp: %s", msg)
}

// readPacket is one received packet: the bytes from '$' to '#cs' with
// the checksum verified. ok is false on framing or checksum damage.
func (c *Conn) readPacket() (string, bool) {
	if b := c.readByte(); b != '$' {
		return "", false
	}

	pkt := []byte{'$'}
	for {
		b := c.readByte()
		if b == 0 || b == '$' {
			return "", false // framing lost mid-packet
		}

		pkt = append(pkt, b)
		if b == '#' {
			// the two checksum digits complete the packet
			hi := c.readByte()
			lo := c.readByte()
			if hi == 0 || lo == 0 {
				return "", false
			}

			return parsePacket(append(pkt, hi, lo))
		}
	}
}
