package rsp

import (
	"fmt"
	"strconv"
	"strings"
)

// StopReply is the target's report of why execution stopped. Kind 'T'
// (stopped with signal) and 'S' (same, no fields) carry the signal in
// Signal; 'W' and 'X' report process exit/termination with the code in
// Signal. Fields holds the annex pairs ("thread" -> "p1.1", "core" ->
// "0", "watch" -> "8021a4f0" for watchpoint reports).
type StopReply struct {
	Kind   byte
	Signal int
	Fields map[string]string
}

func NewStopReply(kind byte, signal int, fields map[string]string) StopReply {
	return StopReply{
		Kind:   kind,
		Signal: signal,
		Fields: fields,
	}
}

// Exited reports whether the reply is a process exit or termination
// (no further execution is possible on this target).
func (s StopReply) Exited() bool {
	return s.Kind == 'W' || s.Kind == 'X'
}

// String renders the reply for logs and the stop line of the UI.
func (s StopReply) String() string {
	switch s.Kind {
	case 'W':
		return fmt.Sprintf("exited with %d", s.Signal)
	case 'X':
		return fmt.Sprintf("terminated by %d", s.Signal)
	case 'S':
		return fmt.Sprintf("stopped by signal %d", s.Signal)
	}

	pairs := make([]string, 0, len(s.Fields))
	for k, v := range s.Fields {
		pairs = append(pairs, k+":"+v)
	}

	return fmt.Sprintf("stopped, signal %d (%s)", s.Signal, strings.Join(pairs, ","))
}

// parseStopReply decodes a stop reply payload ("T05thread:p1.1;...",
// "S05", "W05", "X05"); field order in the annex is arbitrary.
func parseStopReply(payload string) (StopReply, error) {
	if len(payload) < 3 {
		return StopReply{}, fmt.Errorf("assembly/rsp: malformed stop reply %q", payload)
	}

	kind := payload[0]
	signal, err := strconv.ParseUint(payload[1:3], 16, 8)
	if err != nil {
		return StopReply{}, fmt.Errorf("assembly/rsp: malformed stop reply %q", payload)
	}

	fields := map[string]string{}
	for pair := range strings.SplitSeq(payload[3:], ";") {
		if k, v, ok := strings.Cut(pair, ":"); ok {
			fields[k] = v
		}
	}

	switch kind {
	case 'T', 'S', 'W', 'X':
		return NewStopReply(kind, int(signal), fields), nil
	}

	return StopReply{}, fmt.Errorf("assembly/rsp: unknown stop reply kind %q", string(kind))
}
