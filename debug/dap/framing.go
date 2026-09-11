package dap

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

// framing is the DAP wire: LSP-style headers over a byte transport
// ("Content-Length: N", a blank line, then the N bytes of the JSON
// body). Writes are serialized - responses from the request loop and
// events from the session goroutine share the channel.
type framing struct {
	r  *bufio.Reader
	w  io.Writer
	mu sync.Mutex
}

func newFraming(r io.Reader, w io.Writer) *framing {
	return &framing{r: bufio.NewReader(r), w: w}
}

// read is the next envelope body (the byte count its Content-Length
// header declares). Unknown headers are skipped; a conversation that
// ends without a length reports the transport error.
func (f *framing) read() ([]byte, error) {
	length := -1
	for {
		line, err := f.r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		header := strings.TrimRight(line, "\r\n")
		if header == "" {
			break
		}

		name, value, ok := strings.Cut(header, ":")
		if !ok {
			return nil, fmt.Errorf("assembly/dap: malformed header %q", header)
		}

		if strings.EqualFold(name, "Content-Length") {
			length, err = strconv.Atoi(strings.TrimSpace(value))
			if err != nil || length < 0 {
				return nil, fmt.Errorf("assembly/dap: bad Content-Length %q", value)
			}
		}
	}

	if length < 0 {
		return nil, errors.New("assembly/dap: no Content-Length header")
	}

	body := make([]byte, length)
	if _, err := io.ReadFull(f.r, body); err != nil {
		return nil, fmt.Errorf("assembly/dap: body: %w", err)
	}

	return body, nil
}

// write sends one envelope: the header and the body of it, serialized
// against concurrent senders.
func (f *framing) write(body []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	head := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := io.WriteString(f.w, head); err != nil {
		return fmt.Errorf("assembly/dap: header: %w", err)
	}

	if _, err := f.w.Write(body); err != nil {
		return fmt.Errorf("assembly/dap: body: %w", err)
	}

	return nil
}
