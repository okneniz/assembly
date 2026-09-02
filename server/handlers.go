package server

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/okneniz/parsec/bytes"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
	bf "github.com/okneniz/assembly/file"
)

// errorResponse is the JSON error envelope.
type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(msg string) errorResponse {
	return errorResponse{Error: msg}
}

// sectionInfo is metadata of the disassembled section.
type sectionInfo struct {
	Name string `json:"name"`
	Addr string `json:"addr"` // hex
	Size uint64 `json:"size"`
}

func newSectionInfo(name string, addr string, size uint64) sectionInfo {
	return sectionInfo{
		Name: name,
		Addr: addr,
		Size: size,
	}
}

// disasmResponse is the full POST /api/v1/disasm response. Instructions
// arrive as json.RawMessage: both architectures serialize themselves
// (MarshalJSON, see buildResponseInstr).
type disasmResponse struct {
	Arch         string            `json:"arch"`
	Section      sectionInfo       `json:"section"`
	Count        int               `json:"count"`
	Instructions []json.RawMessage `json:"instructions"`
}

func newDisasmResponse(
	arch string,
	section sectionInfo,
	count int,
	instructions []json.RawMessage,
) disasmResponse {
	return disasmResponse{
		Arch:         arch,
		Section:      section,
		Count:        count,
		Instructions: instructions,
	}
}

// handleDisasm handles POST /api/v1/disasm.
func (s *Server) handleDisasm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	// Cap the body size. MaxBytesReader aborts the read immediately on
	// overflow and returns an error on any further read attempt.
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBytes)

	// Parse multipart. The "file" field with the binary is required.
	const maxMem = 32 << 20 // 32 MiB in memory, the rest spills to disk
	if err := r.ParseMultipartForm(maxMem); err != nil {
		// Exceeding the MaxBytesReader body limit is a 413,
		// not a generic bad request.
		if mbErr, ok := errors.AsType[*http.MaxBytesError](err); ok {
			writeError(w, http.StatusRequestEntityTooLarge,
				"upload too large (limit "+humanBytes(mbErr.Limit)+")")
			return
		}

		writeError(w, http.StatusBadRequest, "invalid multipart request: "+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, `missing "file" field`)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("assembly/server: close upload: %v", err)
		}
	}()

	// Save to a temp file: bf.Detect and CodeSection work with a path.
	tmp, err := os.CreateTemp("", "assembly-*.bin")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create temp file: "+err.Error())
		return
	}

	tmpPath := tmp.Name()
	defer func() {
		// cleanup: a missing file is the norm; other failures go to the log
		if err := os.Remove(tmpPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
			log.Printf("assembly/server: remove temp file: %v", err)
		}
	}()

	if _, err := io.Copy(tmp, file); err != nil {
		// secondary error: the copy already failed, closing is best effort
		if cerr := tmp.Close(); cerr != nil {
			log.Printf("assembly/server: close temp file after failed copy: %v", cerr)
		}

		if mbErr, ok := errors.AsType[*http.MaxBytesError](err); ok {
			writeError(w, http.StatusRequestEntityTooLarge,
				"upload too large (limit "+humanBytes(mbErr.Limit)+")")
			return
		}

		writeError(w, http.StatusBadRequest, "read file: "+err.Error())
		return
	}

	if err := tmp.Close(); err != nil {
		writeError(w, http.StatusInternalServerError, "close temp file: "+err.Error())
		return
	}

	// Detect the format by magic and load the code section.
	ff, err := bf.Detect(tmpPath)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("detect format: %v", err))
		return
	}

	sec, err := ff.CodeSection()
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("load code section: %v", err))
		return
	}

	// Pick the decoder from the binary's architecture (cputype/e_machine)
	// and call the arch-specific Parse directly — there is no shared
	// instruction/decoder type.
	var resp disasmResponse
	switch ff.ArchKind() {
	case bf.ArchARM64:
		insts, err := arm64.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode: "+err.Error())
			return
		}

		resp = buildResponseInstr(arm64.Name, "text", sec.Addr, uint64(len(sec.Data)), insts)
	case bf.ArchRISCV64:
		insts, err := riscv.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode: "+err.Error())
			return
		}

		resp = buildResponseInstr(riscv.Name, "text", sec.Addr, uint64(len(sec.Data)), insts)
	case bf.ArchLOONGARCH64:
		insts, err := loong64.Parse(sec.Addr)(bytes.Buffer(sec.Data))
		if err != nil {
			writeError(w, http.StatusBadRequest, "decode: "+err.Error())
			return
		}

		resp = buildResponseInstr(loong64.Name, "text", sec.Addr, uint64(len(sec.Data)), insts)
	default:
		writeError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("unsupported architecture (kind %d)", ff.ArchKind()),
		)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	if err := enc.Encode(resp); err != nil {
		// the response has already started — a new status can't be sent, log only
		log.Printf("assembly/server: write response: %v", err)
	}

	_ = header // file name unused; reserved for future needs
}

// buildResponseInstr assembles the JSON response from self-serializing
// instructions (arm64.Instr / riscv.Instr implement MarshalJSON).
func buildResponseInstr[T json.Marshaler](
	archName, secName string,
	secAddr, secSize uint64,
	instrs []T,
) disasmResponse {
	raws := make([]json.RawMessage, 0, len(instrs))
	for _, in := range instrs {
		b, err := json.Marshal(in)
		if err != nil {
			continue
		}

		raws = append(raws, b)
	}

	return newDisasmResponse(
		archName,
		newSectionInfo(secName, fmt.Sprintf("0x%x", secAddr), secSize),
		len(raws),
		raws,
	)
}

// writeError sends a JSON error with the given status code.
func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(newErrorResponse(msg)); err != nil {
		// the response has already started — a new status can't be sent, log only
		log.Printf("assembly/server: write error response: %v", err)
	}
}

// humanBytes returns a human-readable size (e.g. "256 MiB").
func humanBytes(n int64) string {
	const unit = 1 << 20
	if n >= unit {
		return fmt.Sprintf("%d MiB", n/unit)
	}

	if n >= 1<<10 {
		return fmt.Sprintf("%d KiB", n/(1<<10))
	}

	return fmt.Sprintf("%d B", n)
}

// asmRequest is the POST /api/v1/asm body.
type asmRequest struct {
	Arch     string `json:"arch"`               // "arm64" | "riscv64"
	Source   string `json:"source"`             // assembler source
	BaseAddr string `json:"baseAddr,omitempty"` // hex/dec, default 0
}

func newAsmRequest(arch string, source string, baseAddr string) asmRequest {
	return asmRequest{
		Arch:     arch,
		Source:   source,
		BaseAddr: baseAddr,
	}
}

// asmSectionDTO is an assembled section in the response.
type asmSectionDTO struct {
	Name string `json:"name"`
	Addr string `json:"addr"` // hex
	Data string `json:"data"` // hex string of bytes
	Size int    `json:"size"`
}

func newAsmSectionDTO(name string, addr string, data string, size int) asmSectionDTO {
	return asmSectionDTO{
		Name: name,
		Addr: addr,
		Data: data,
		Size: size,
	}
}

// asmErrorDTO is an error with a position in the source.
type asmErrorDTO struct {
	Line int    `json:"line"`
	Col  int    `json:"col"`
	Msg  string `json:"msg"`
}

func newAsmErrorDTO(line int, col int, msg string) asmErrorDTO {
	return asmErrorDTO{
		Line: line,
		Col:  col,
		Msg:  msg,
	}
}

// asmResponse is the POST /api/v1/asm response.
type asmResponse struct {
	Sections []asmSectionDTO   `json:"sections"`
	Symbols  map[string]string `json:"symbols"` // name → hex addr
	Errors   []asmErrorDTO     `json:"errors"`
}

func newAsmResponse(symbols map[string]string, errors []asmErrorDTO) asmResponse {
	return asmResponse{
		Symbols: symbols,
		Errors:  errors,
	}
}

// handleAsm handles POST /api/v1/asm: assembles the source for the selected
// architecture and returns sections/symbols/errors.
func (s *Server) handleAsm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed, use POST")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, s.maxBytes)

	var req asmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body: "+err.Error())
		return
	}

	var assemble func(string, uint64) (*asm.Result, []asm.AsmError)
	switch strings.ToLower(req.Arch) {
	case "arm64", "aarch64":
		assemble = alias.Assemble
	case "riscv64", "riscv":
		assemble = pseudo.Assemble
	case "loong64", "loongarch64":
		assemble = lpseudo.Assemble
	default:
		writeError(w, http.StatusBadRequest, "unknown arch (want arm64, riscv64 or loong64)")
		return
	}

	base := uint64(0)
	if req.BaseAddr != "" {
		v, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(req.BaseAddr), "0x"), 16, 64)
		if err == nil {
			base = v
		} else if v, err := strconv.ParseUint(req.BaseAddr, 10, 64); err == nil {
			base = v
		} else {
			writeError(w, http.StatusBadRequest, "bad baseAddr")
			return
		}
	}

	res, errs := assemble(req.Source, base)
	resp := newAsmResponse(map[string]string{}, make([]asmErrorDTO, 0, len(errs)))
	for _, sec := range res.Sections {
		resp.Sections = append(
			resp.Sections,
			newAsmSectionDTO(
				sec.Name,
				fmt.Sprintf("%#x", sec.Addr),
				hex.EncodeToString(sec.Data),
				len(sec.Data),
			),
		)
	}

	for name, addr := range res.Symbols {
		resp.Symbols[name] = fmt.Sprintf("%#x", addr)
	}

	for _, e := range errs {
		resp.Errors = append(
			resp.Errors,
			newAsmErrorDTO(int(e.Line), int(e.Col), e.Msg),
		)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// the response has already started — a new status can't be sent, log only
		log.Printf("assembly/server: write response: %v", err)
	}
}
