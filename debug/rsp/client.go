package rsp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

// ourFeatures is what the client advertises at negotiation: software
// and hardware breakpoints understood; the target description XML
// requested for every arch we debug (a stub serves the register layout
// only when the client names the arch it can parse).
const ourFeatures = "qSupported:multiprocess-;swbreak+;hwbreak+;xmlRegisters=aarch64,riscv:rv64,loongarch64,i386"

// xmlChunk is the qXfer read granularity: one target description is a
// few KB, so a couple of chunks covers it.
const xmlChunk = 0xfff

// maxXMLChunks bounds the target description read loop.
const maxXMLChunks = 64

// Client is the typed command set of the protocol over one Conn: the
// debugging vocabulary the session speaks. Every command is a full
// round trip; Continue and Step block until the target reports a stop
// (an interrupt from another goroutine - Conn.Interrupt - unblocks it).
type Client struct {
	conn *Conn
}

// NewClient is a client over an established conversation.
func NewClient(conn *Conn) *Client {
	return &Client{conn: conn}
}

// Interrupt forwards the interrupt byte to the running target.
func (c *Client) Interrupt() error {
	return c.conn.Interrupt()
}

// Supported negotiates the protocol features: the target's answers
// ("qXfer:features:read+" -> "+", "PacketSize" -> "1000", ...). Call it
// once after connecting, before TargetXML.
func (c *Client) Supported() (map[string]string, error) {
	payload, err := c.conn.roundTrip(ourFeatures)
	if err != nil {
		return nil, err
	}

	features := map[string]string{}
	for item := range strings.SplitSeq(payload, ";") {
		if item == "" {
			continue
		}

		switch {
		case strings.HasSuffix(item, "+"):
			features[strings.TrimSuffix(item, "+")] = "+"
		case strings.HasSuffix(item, "-"):
			features[strings.TrimSuffix(item, "-")] = "-"
		default:
			k, v, _ := strings.Cut(item, "=")
			features[k] = v
		}
	}

	return features, nil
}

// TargetXML reads the target description (the register names, numbers,
// and widths) via qXfer. The annex is target.xml; the outer include
// files of a multi-part description are not followed (the register
// layout of every qemu core target sits in the main file).
func (c *Client) TargetXML() (string, error) {
	var sb strings.Builder
	for chunk := range maxXMLChunks {
		cmd := fmt.Sprintf("qXfer:features:read:target.xml:%x,%x", chunk*xmlChunk, xmlChunk)
		payload, err := c.conn.roundTrip(cmd)
		if err != nil {
			return "", err
		}

		if payload == "" || (payload[0] != 'm' && payload[0] != 'l') {
			return "", fmt.Errorf("assembly/rsp: malformed qXfer reply %q", payload)
		}

		data, err := hex.DecodeString(payload[1:])
		if err != nil {
			return "", fmt.Errorf("assembly/rsp: malformed qXfer data: %w", err)
		}

		sb.Write(data)
		if payload[0] == 'l' { // the last chunk
			return sb.String(), nil
		}
	}

	return "", fmt.Errorf("assembly/rsp: target description exceeds %d chunks", maxXMLChunks)
}

// HaltReason is the target's report of why it is stopped ("?"); the
// first thing to ask after connecting to a halted (-S) machine.
func (c *Client) HaltReason() (StopReply, error) {
	payload, err := c.conn.roundTrip("?")
	if err != nil {
		return StopReply{}, err
	}

	return parseStopReply(payload)
}

// SelectThread points the subsequent register and memory commands at
// the thread (a qemu SMP core or a bare-metal hart).
func (c *Client) SelectThread(id int) error {
	payload, err := c.conn.roundTrip(fmt.Sprintf("Hg%x", id))
	if err != nil {
		return err
	}

	if payload != "OK" {
		return fmt.Errorf("assembly/rsp: Hg reply %q", payload)
	}

	return nil
}

// ReadRegs is the whole register block of the current thread ("g"):
// raw bytes in the target description's order, per-register widths.
func (c *Client) ReadRegs() ([]byte, error) {
	payload, err := c.conn.roundTrip("g")
	if err != nil {
		return nil, err
	}

	regs, err := hex.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("assembly/rsp: malformed register block: %w", err)
	}

	return regs, nil
}

// ReadReg is one register ("p n") as a little-endian integer: the low
// bytes of the register in target order; registers wider than 64 bits
// (vector ones) are truncated to their low half.
func (c *Client) ReadReg(n int) (uint64, error) {
	payload, err := c.conn.roundTrip(fmt.Sprintf("p%x", n))
	if err != nil {
		return 0, err
	}

	raw, err := hex.DecodeString(payload)
	if err != nil || len(raw) == 0 {
		return 0, fmt.Errorf("assembly/rsp: malformed register %d reply %q", n, payload)
	}

	if len(raw) > 8 {
		raw = raw[:8]
	}

	var v uint64
	for i, b := range raw { // little-endian: first byte is the lowest
		v |= uint64(b) << (8 * i)
	}

	return v, nil
}

// WriteReg sets one register ("P n=r"): the value goes out as eight
// little-endian bytes.
func (c *Client) WriteReg(n int, v uint64) error {
	var raw [8]byte
	binary.LittleEndian.PutUint64(raw[:], v)
	payload, err := c.conn.roundTrip(fmt.Sprintf("P%x=%s", n, hex.EncodeToString(raw[:])))
	if err != nil {
		return err
	}

	if payload != "OK" {
		return fmt.Errorf("assembly/rsp: P reply %q", payload)
	}

	return nil
}

// ReadMem is n bytes at addr ("m"): the bytes, not the hex form.
func (c *Client) ReadMem(addr uint64, n int) ([]byte, error) {
	payload, err := c.conn.roundTrip(fmt.Sprintf("m%x,%x", addr, n))
	if err != nil {
		return nil, err
	}

	mem, err := hex.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("assembly/rsp: malformed memory reply: %w", err)
	}

	if len(mem) != n {
		return nil, fmt.Errorf("assembly/rsp: memory reply %d bytes, want %d", len(mem), n)
	}

	return mem, nil
}

// WriteMem writes bytes at addr ("M").
func (c *Client) WriteMem(addr uint64, b []byte) error {
	cmd := fmt.Sprintf("M%x,%x:%s", addr, len(b), hex.EncodeToString(b))
	payload, err := c.conn.roundTrip(cmd)
	if err != nil {
		return err
	}

	if payload != "OK" {
		return fmt.Errorf("assembly/rsp: M reply %q", payload)
	}

	return nil
}

// SetBreak inserts a software breakpoint ("Z0"): kind is the length in
// bytes the stub patches (the instruction width).
func (c *Client) SetBreak(addr uint64, kind int) error {
	return c.breakpoint('Z', addr, kind)
}

// ClearBreak removes a software breakpoint ("z0").
func (c *Client) ClearBreak(addr uint64, kind int) error {
	return c.breakpoint('z', addr, kind)
}

// Continue resumes the current thread ("c") and blocks until the stop
// report: a breakpoint hit, a step done, an interrupt, or an exit.
func (c *Client) Continue() (StopReply, error) {
	payload, err := c.conn.roundTrip("c")
	if err != nil {
		return StopReply{}, err
	}

	return parseStopReply(payload)
}

// Step executes one instruction of the current thread ("s") and returns
// the stop report.
func (c *Client) Step() (StopReply, error) {
	payload, err := c.conn.roundTrip("s")
	if err != nil {
		return StopReply{}, err
	}

	return parseStopReply(payload)
}

func (c *Client) breakpoint(op byte, addr uint64, kind int) error {
	cmd := fmt.Sprintf("%c0,%x,%x", op, addr, kind)
	payload, err := c.conn.roundTrip(cmd)
	if err != nil {
		return err
	}

	if payload != "OK" {
		return fmt.Errorf("assembly/rsp: %c0 reply %q", op, payload)
	}

	return nil
}
