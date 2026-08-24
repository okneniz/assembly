package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleAsm(t *testing.T) {
	s := New("127.0.0.1:0")
	body, err := json.Marshal(
		newAsmRequest("riscv64", "li a0, 0x42\nret\n", "0x1000"),
	)
	require.NoError(t, err, "marshal")
	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/asm",
		bytes.NewReader(body),
	)
	rec := httptest.NewRecorder()
	s.handleAsm(rec, req)
	require.Equal(
		t,
		http.StatusOK,
		rec.Code,
		"status %d: %s",
		rec.Code,
		rec.Body.String(),
	)
	var resp asmResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Len(t, resp.Sections, 1, "bad sections: %+v", resp.Sections)
	require.NotZero(t, resp.Sections[0].Size, "bad sections: %+v", resp.Sections)
	require.Equal(t, "0x1000", resp.Sections[0].Addr, "addr")
	require.Empty(t, resp.Errors, "errors: %+v", resp.Errors)

	// parse errors come back with positions
	body, err = json.Marshal(newAsmRequest("arm64", "frobnicate x0\n", ""))
	require.NoError(t, err, "marshal")
	req = httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/asm",
		bytes.NewReader(body),
	)
	rec = httptest.NewRecorder()
	s.handleAsm(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "status")
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.NotEmpty(t, resp.Errors, "expected line-1 error")
	require.Equal(t, 1, resp.Errors[0].Line, "expected line-1 error, got %+v", resp.Errors)
}
