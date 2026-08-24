package asm

// Two-pass assembly: pass 1 only computes the layout (subsection counters,
// label offsets, .set, globals), then finalization - the subsections of each
// section are concatenated in ascending number order (GAS), the sections get
// consecutive addresses from base, and pass 2 encodes instructions at the
// final addresses with full symbol resolution. The walks must produce an
// identical layout: all sizes are deterministic without symbol values (see
// finalizeLayout for the single deliberate exception - RVC boundaries when
// returning between subsections).

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"slices"

	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
)

// Assemble assembles the source into sections with base address base.
// Returns the result and all accumulated errors (assembly continues past
// line errors - best-effort, like objdump).
func Assemble(src string, base uint64, be Syntax) (*Result, []AsmError) {
	a := newAssembler(
		be,
		map[string]labelRef{},
		map[string]*expr.Expr{},
		map[string][]numLabelDef{},
		base,
	)
	a.switchSection(".text", 0)
	stmts := parseSource([]rune(src), be)

	// line parse errors go into the shared list
	for i := range stmts {
		if stmts[i].err != nil {
			a.errs = append(a.errs, *stmts[i].err)
		}
	}

	// pass 1: layout (subsection counters; addresses are "live", see pass1Addr)
	a.walk(stmts, false)

	// finalization: the subsections of each section get base offsets in
	// ascending number order (GAS concatenation), the sections get
	// consecutive addresses from base
	a.finalizeLayout()

	// pass 2: encoding at the final addresses
	a.walk(stmts, true)

	// literal pools: slot values go to the tails of subsections
	a.emitPoolRecords()

	res := NewResult(map[string]uint64{}, a.globals)
	for i, s := range a.secs {
		data := s.concat()
		if len(data) == 0 && (!s.nobits || s.secSize() == 0) {
			continue
		}

		if s.nobits {
			res.Sections = append(res.Sections, NewNobitsSection(s.name, a.secAddr[i], s.secSize()))
		} else {
			res.Sections = append(res.Sections, NewSection(s.name, a.secAddr[i], data))
		}
	}

	resolve := a.resolver(-1, base) // no position - numeric references do not resolve
	for name, lr := range a.labels {
		res.Symbols[name] = a.labelAddr(lr)
	}

	for name := range a.sets {
		if v, ok := resolve(name); ok {
			res.Symbols[name] = v
		}
	}

	return res, a.errs
}

// subBuf is a subsection: an independent write stream (GAS: .text N). size
// is the layout pass counter, base is the offset within the section (assigned
// by finalizeLayout: subsections are concatenated in ascending number
// order), data is the pass 2 encoding.
type subBuf struct {
	num  int
	size int
	base int
	data []byte
}

func newSubBuf(num int) *subBuf {
	return &subBuf{num: num}
}

// secBuf is a section: name, NOBITS mode (.bss), and subsections by number.
// A section without explicit subsections is the single subsection 0.
type secBuf struct {
	name   string
	nobits bool // .bss: has a size, no file data (see doDirective)
	subs   map[int]*subBuf
}

func newSecBuf(name string, nobits bool) *secBuf {
	return &secBuf{name: name, nobits: nobits, subs: map[int]*subBuf{}}
}

// secSize is the full section size (the sum of subsections; in the layout
// pass it is "live", per the current counters).
func (s *secBuf) secSize() int {
	total := 0
	for _, sub := range s.subs {
		total += sub.size
	}

	return total
}

// sortedSubs is the subsections in ascending number order (GAS
// concatenation).
func (s *secBuf) sortedSubs() []*subBuf {
	nums := make([]int, 0, len(s.subs))
	for n := range s.subs {
		nums = append(nums, n)
	}

	slices.Sort(nums)
	out := make([]*subBuf, len(nums))
	for i, n := range nums {
		out[i] = s.subs[n]
	}

	return out
}

// concat is the section data: subsections in ascending number order (after
// the encoding pass; empty for NOBITS).
func (s *secBuf) concat() []byte {
	var out []byte
	for _, sub := range s.sortedSubs() {
		out = append(out, sub.data...)
	}

	return out
}

type labelRef struct {
	sec, sub, off int
}

func newLabelRef(sec int, sub int, off int) labelRef {
	return labelRef{
		sec: sec,
		sub: sub,
		off: off,
	}
}

// numLabelDef is one definition of a numeric local label: the statement
// index (the nearest Nb/Nf is selected by source order, as in GAS) + the
// address.
type numLabelDef struct {
	stmtIdx int
	ref     labelRef
}

func newNumLabelDef(stmtIdx int, ref labelRef) numLabelDef {
	return numLabelDef{
		stmtIdx: stmtIdx,
		ref:     ref,
	}
}

// poolEntry is a literal pool slot of a subsection (see PoolUser): the value
// is an expression evaluated when the pool is written; name is the auto-name
// of the slot.
type poolEntry struct {
	name string
	expr *expr.Expr
	slot int
}

type assembler struct {
	be        Syntax
	secs      []*secBuf
	curIdx    int
	curSub    *subBuf
	secAddr   []uint64
	labels    map[string]labelRef
	numLabels map[string][]numLabelDef // numeric locals, redefinable
	sets      map[string]*expr.Expr
	globals   []string
	errs      []AsmError
	base      uint64
	incbins   map[string][]byte       // .incbin cache: both passes read the same bytes
	pools     map[*subBuf][]poolEntry // literal pools of subsections (order of appearance)
	poolAddr  map[string]uint64       // slot addresses (after finalizeLayout)
}

func newAssembler(
	be Syntax,
	labels map[string]labelRef,
	sets map[string]*expr.Expr,
	numLabels map[string][]numLabelDef,
	base uint64,
) *assembler {
	return &assembler{
		be:        be,
		labels:    labels,
		sets:      sets,
		numLabels: numLabels,
		base:      base,
		incbins:   map[string][]byte{},
		pools:     map[*subBuf][]poolEntry{},
		poolAddr:  map[string]uint64{},
	}
}

func (a *assembler) cur() *secBuf {
	return a.secs[a.curIdx]
}

func (a *assembler) errf(pos parsecstrings.Position, format string, args ...any) {
	a.errs = append(a.errs, posErr(pos, fmt.Sprintf(format, args...)))
}

// switchSection selects the section by name and its subsection by number,
// creating them on first mention (NOBITS only for .bss: sections created by
// .section are regular PROGBITS; .section always gives subsection 0). Called
// in both passes in the same order - the sections and subsections are the
// same.
func (a *assembler) switchSection(name string, sub int) {
	for i, s := range a.secs {
		if s.name == name {
			a.curIdx = i
			a.curSub = s.subOf(sub)
			return
		}
	}

	s := newSecBuf(name, name == ".bss")
	a.secs = append(a.secs, s)
	a.curIdx = len(a.secs) - 1
	a.curSub = s.subOf(sub)
}

// subOf is the subsection by number (created on first access).
func (s *secBuf) subOf(num int) *subBuf {
	if sub, ok := s.subs[num]; ok {
		return sub
	}

	sub := newSubBuf(num)
	s.subs[num] = sub
	return sub
}

// walk is the shared walk for both passes; pass2=false computes the layout,
// pass2=true encodes. Before each pass the state is returned to the start:
// the Syntax modes are reset, the walk again starts from the first
// section/subsection - otherwise .option/sections would apply
// asymmetrically (early pass 2 instructions would see late modes/sections).
func (a *assembler) walk(stmts []statement, pass2 bool) {
	a.be.ResetOptions()
	a.curIdx = 0
	a.curSub = a.secs[0].subs[0] // starting .text: subsection 0
	for i := range stmts {
		st := &stmts[i]
		if st.err != nil {
			continue
		}

		if !pass2 {
			for _, lbl := range st.labels {
				if isNumericLabel(lbl) {
					a.numLabels[lbl] = append(a.numLabels[lbl],
						newNumLabelDef(i, newLabelRef(a.curIdx, a.curSub.num, a.curSub.size)))
					continue // redefining a numeric label is legal
				}

				if _, dup := a.labels[lbl]; dup {
					a.errf(st.pos, "label %q redefined", lbl)
					continue
				}

				a.labels[lbl] = newLabelRef(a.curIdx, a.curSub.num, a.curSub.size)
			}
		}

		switch {
		case st.directive != nil:
			a.doDirective(st, i, pass2)
		case st.hasInstr:
			a.doInstr(st, i, pass2)
		}
	}
}

// finalizeLayout assigns the layout after pass 1: the subsections of each
// section get base offsets in ascending number order (GAS concatenation:
// all subsections of a section are concatenated by number, the write order
// within a subsection is preserved), the sections get consecutive addresses
// from base without gaps. Pass 2 encodes at these FINAL addresses; if, when
// returning to an early subsection, a literal PC-relative target lands on an
// RVC range boundary, the size from the final address may diverge from the
// computed one - this is an explicit encoding error (got N bytes, layout
// pass reported M), not a silent bug.
func (a *assembler) finalizeLayout() {
	a.secAddr = make([]uint64, len(a.secs))
	next := a.base
	for i, s := range a.secs {
		a.secAddr[i] = next
		base := 0
		for _, sub := range s.sortedSubs() {
			sub.base = base
			base += sub.size

			// the subsection pool is its tail (GAS: a separate pool per
			// subsection); the slot addresses are for the pass 2 resolver
			off := base
			for _, e := range a.pools[sub] {
				a.poolAddr[e.name] = a.secAddr[i] + uint64(off)
				off += e.slot
			}

			base = off
		}

		if i+1 < len(a.secs) {
			next += uint64(base)
		}
	}
}

// poolAdd registers a literal pool slot of the current subsection (dedup by
// auto-name: PoolName(slot, ExprKey)); the order is first appearance.
func (a *assembler) poolAdd(val *expr.Expr, slot int) {
	name := poolName(slot, expr.ExprKey(val))
	for _, e := range a.pools[a.curSub] {
		if e.name == name {
			return
		}
	}

	a.pools[a.curSub] = append(a.pools[a.curSub], poolEntry{name: name, expr: val, slot: slot})
}

// emitPoolRecords appends the literal pools to the subsection data (after
// encoding: the expression values are evaluated with the full resolver; an
// expression error makes the element zeros - the AsmError with the
// instruction position was already recorded at encoding, zeros here
// silently).
func (a *assembler) emitPoolRecords() {
	for _, s := range a.secs {
		for _, sub := range s.sortedSubs() {
			for _, e := range a.pools[sub] {
				addr := a.poolAddr[e.name]
				v, err := e.expr.Eval(a.resolver(-1, addr))
				if err != nil {
					v = 0
				}

				var buf [8]byte
				binary.LittleEndian.PutUint64(buf[:], uint64(v))
				sub.data = append(sub.data, buf[:e.slot]...)
			}
		}
	}
}

func (a *assembler) doInstr(st *statement, idx int, pass2 bool) {
	if a.cur().nobits {
		if !pass2 {
			a.errf(st.pos, "section %s is NOBITS: instructions are not permitted", a.cur().name)
		}

		return
	}

	if !pass2 {
		// literal pool: the slot goes at the end of the current subsection
		// (dedup by name - identical literals share a slot; the instruction
		// size does not depend on the pool, the pool does not break the
		// layout)
		if pu, ok := st.instr.(PoolUser); ok {
			if e, slot, ok2 := pu.PoolReq(); ok2 {
				a.poolAdd(e, slot)
			}
		}
	}

	var addr uint64
	if pass2 {
		addr = a.secAddr[a.curIdx] + uint64(a.curSub.base+len(a.curSub.data))
	} else {
		addr = a.pass1Addr()
	}

	// Pool: the address of its OWN slot - via the reserved name PoolSelf
	// (the slot naming scheme does not leave the core)
	resolve := a.resolver(idx, addr)
	if pu, ok := st.instr.(PoolUser); ok {
		if e, slot, ok2 := pu.PoolReq(); ok2 {
			if pa, found := a.poolAddr[poolName(slot, expr.ExprKey(e))]; found {
				inner := resolve
				resolve = func(name string) (uint64, bool) {
					if name == PoolSelf {
						return pa, true
					}

					return inner(name)
				}
			}
		}
	}

	// Size = a trial Resolve with a placeholder environment, written to the
	// counter; deterministic across passes (see sizeOf), so in pass 2 we
	// repeat it without reporting the error - it was already recorded in
	// pass 1.
	size, err := sizeOf(st.instr, newCtx(addr, placeholderResolve(addr)))
	if err != nil {
		if !pass2 {
			a.errf(st.pos, "size: %v", err)
		}

		return
	}

	if !pass2 {
		a.curSub.size += size
		return
	}

	var buf bytes.Buffer
	res, rerr := st.instr.Resolve(newCtx(addr, resolve))
	if rerr == nil {
		_, rerr = res.Encode(&buf)
	}

	if rerr != nil {
		a.errf(st.pos, "encode: %v", rerr)
		a.curSub.data = append(a.curSub.data, make([]byte, size)...) // pass 1 layout
		return
	}

	if buf.Len() != size {
		a.errf(st.pos, "encode: got %d bytes, layout pass reported %d", buf.Len(), size)
	}

	a.curSub.data = append(a.curSub.data, buf.Bytes()...)
}

// pass1Addr is the "live" address of the current point of the first pass:
// base + the current sizes of the preceding sections + the current
// subsection base (the sum of the smaller subsections at the moment) + the
// subsection counter. It matches the final pass 2 address when subsections
// are filled in non-descending number order; when returning to earlier
// subsections, the pass 1 addresses may differ from the final ones (see
// finalizeLayout) - the divergence is caught by the encoding length check.
func (a *assembler) pass1Addr() uint64 {
	addr := a.base
	for i := range a.curIdx {
		addr += uint64(a.secs[i].secSize())
	}

	for _, sub := range a.cur().sortedSubs() {
		if sub.num >= a.curSub.num {
			break
		}

		addr += uint64(sub.size)
	}

	return addr + uint64(a.curSub.size)
}

func (a *assembler) doDirective(st *statement, idx int, pass2 bool) {
	d := st.directive
	switch d.name {
	case ".text", ".data", ".bss":
		num, ok := a.subsecNum(st, d, pass2)
		if !ok {
			return
		}

		a.switchSection(d.name, num)
	case ".section":
		if len(d.args) > 0 {
			a.switchSection(d.args[0].str, 0) // .section is always subsection 0
		}
	case ".ltorg":
		// GAS: a forced flush of the literal pool at the current point; here
		// the pool is always emitted at the tail of the subsection (see
		// emitPoolRecords) - silently ignoring it would change the layout,
		// hence the explicit error
		if !pass2 {
			a.errf(
				st.pos,
				".ltorg is not supported: literal pool is emitted at the end of the subsection",
			)
		}
	case ".set", ".equ":
		if len(d.args) >= 2 {
			a.sets[d.args[0].str] = d.args[1].expr
		}
	case ".option":
		// applied symmetrically in both passes (ResetOptions in walk); the
		// error is reported in pass 1 only, to avoid duplication
		if len(d.args) > 0 {
			if oerr := a.be.ApplyOption(d.args[0].str); oerr != nil && !pass2 {
				a.errf(st.pos, ".option: %v", oerr)
			}
		}
	case ".global", ".globl", ".local":
		if len(d.args) > 0 {
			a.globals = append(a.globals, d.args[0].str)
		}
	case ".word", ".half", ".short", ".byte", ".quad", ".dword":
		if a.cur().nobits {
			if !pass2 {
				a.errf(
					st.pos,
					"section %s is NOBITS: only .zero/.skip/.align permitted",
					a.cur().name,
				)
			}

			return
		}

		width := map[string]int{".word": 4, ".half": 2, ".short": 2, ".byte": 1, ".quad": 8, ".dword": 8}[d.name]
		if !pass2 {
			a.curSub.size += width * len(d.args)
			return
		}

		addr := a.secAddr[a.curIdx] + uint64(a.curSub.base+len(a.curSub.data))
		resolve := a.resolver(idx, addr)
		for _, arg := range d.args {
			v, err := arg.expr.Eval(resolve)
			if err != nil {
				a.errf(st.pos, "%s: %v", d.name, err)
				a.curSub.data = append(a.curSub.data, make([]byte, width)...)
				continue
			}

			a.appendInt(width, v, st.pos)
		}
	case ".zero", ".space", ".skip":
		n, err := d.args[0].expr.Eval(nil) // the size must be a constant
		if err != nil || n < 0 {
			if !pass2 {
				a.errf(st.pos, "%s: constant non-negative size required", d.name)
			}

			return
		}

		if !pass2 {
			a.curSub.size += int(n)
			return
		}

		if !a.cur().nobits { // NOBITS: the size is already counted, no file zeros
			a.curSub.data = append(a.curSub.data, make([]byte, n)...)
		}
	case ".align", ".p2align", ".balign":
		n, err := d.args[0].expr.Eval(nil)
		if err != nil || n < 0 {
			if !pass2 {
				a.errf(st.pos, "%s: constant alignment required", d.name)
			}

			return
		}

		align := uint64(1) << uint(n)
		if d.name == ".balign" {
			align = uint64(n)
		}

		// 1<<64 in Go = 0 (shift overflow), .balign 0 is zero: both would
		// give division by zero in the padding formula; GAS rejects such
		// values
		if align == 0 {
			if !pass2 {
				if d.name == ".balign" {
					a.errf(st.pos, ".balign: zero alignment")
				} else {
					a.errf(st.pos, "%s: alignment 2^%d too large", d.name, n)
				}
			}

			return
		}

		if !pass2 {
			off := uint64(a.curSub.size)
			a.curSub.size += int((align - off%align) % align)
			return
		}

		off := uint64(len(a.curSub.data))
		if pad := (align - off%align) % align; pad != 0 {
			a.curSub.data = append(a.curSub.data, make([]byte, pad)...)
		}
	case ".incbin":
		if a.cur().nobits {
			if !pass2 {
				a.errf(st.pos, "section %s is NOBITS: .incbin is not permitted", a.cur().name)
			}

			return
		}

		data, ierr := a.incbinData(d)
		if ierr != nil {
			if !pass2 {
				a.errf(st.pos, ".incbin: %v", ierr)
			}

			return
		}

		if !pass2 {
			a.curSub.size += len(data)
			return
		}

		a.curSub.data = append(a.curSub.data, data...)
	case ".string", ".asciz", ".ascii":
		if a.cur().nobits {
			if !pass2 {
				a.errf(
					st.pos,
					"section %s is NOBITS: only .zero/.skip/.align permitted",
					a.cur().name,
				)
			}

			return
		}

		total := 0
		for _, arg := range d.args {
			total += len(arg.str)
			if d.name != ".ascii" {
				total++ // NUL terminator
			}
		}

		if !pass2 {
			a.curSub.size += total
			return
		}

		for _, arg := range d.args {
			a.curSub.data = append(a.curSub.data, arg.str...)
			if d.name != ".ascii" {
				a.curSub.data = append(a.curSub.data, 0)
			}
		}
	}

	// .type/.size/.file/.loc/.cfi_*/... - recognized, carry no semantics
}

// incbinData is the .incbin data: the file at the path from the first
// argument (resolved from the process cwd - like GAS without an
// include-path), the optional skip/count are the offset and length of the
// insertion. Reading is cached for the assembly: both passes must see the
// same bytes.
func (a *assembler) incbinData(d *directive) ([]byte, error) {
	if len(d.args) == 0 || !d.args[0].isStr {
		return nil, errors.New("file path string expected")
	}

	path := d.args[0].str
	full, ok := a.incbins[path]
	if !ok {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		full = b
		a.incbins[path] = full
	}

	skip, count := int64(0), int64(len(full))
	if len(d.args) > 1 {
		v, err := d.args[1].expr.Eval(nil) // the offset must be a constant
		if err != nil || v < 0 {
			return nil, errors.New("skip: constant non-negative value required")
		}

		skip = v
	}

	if len(d.args) > 2 {
		v, err := d.args[2].expr.Eval(nil) // the length is a constant too
		if err != nil || v < 0 {
			return nil, errors.New("count: constant non-negative value required")
		}

		count = v
	}

	if skip > int64(len(full)) {
		skip = int64(len(full))
	}

	if end := skip + count; end > int64(len(full)) {
		count = int64(len(full)) - skip
	}

	if count < 0 {
		count = 0
	}

	return full[skip : skip+count], nil
}

// appendInt appends a little-endian value of width bytes with a range check
// (signed or unsigned).
func (a *assembler) appendInt(width int, v int64, pos parsecstrings.Position) {
	lo, hi := int64(-1)<<(width*8-1), int64(1)<<(width*8-1)-1
	umax := uint64(1)<<(width*8) - 1
	if (v < lo || v > hi) && uint64(v) > umax {
		a.errf(pos, "value %d does not fit in %d bytes", v, width)
		v = 0
	}

	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(v))
	a.curSub.data = append(a.curSub.data, buf[:width]...)
}

// subsecNum is the subsection number for .text/.data/.bss: an optional
// constant 0..8192 (as in GAS); without an argument, 0. The error is
// reported in pass 1.
func (a *assembler) subsecNum(st *statement, d *directive, pass2 bool) (int, bool) {
	if len(d.args) == 0 {
		return 0, true
	}

	n, err := d.args[0].expr.Eval(nil) // the number must be a constant
	if err != nil || n < 0 || n > 8192 {
		if !pass2 {
			a.errf(st.pos, "%s: subsection number must be a constant 0..8192", d.name)
		}

		return 0, false
	}

	return int(n), true
}

// resolver is the symbol resolution function for statement idx at address
// addr: named labels, "." → addr, .set lazily, with memoization and cycle
// protection; numeric locals "Nb"/"Nf" - the nearest definition relative to
// idx (see resolveLocal). refIdx < 0 is a context without a position
// (filling Symbols): numeric references do not resolve.
func (a *assembler) resolver(refIdx int, addr uint64) func(string) (uint64, bool) {
	memo := map[string]uint64{}
	visiting := map[string]bool{}
	var res func(string) (uint64, bool)
	res = func(name string) (uint64, bool) {
		if v, ok := memo[name]; ok {
			return v, true
		}

		if name == "." {
			return addr, true
		}

		if isPoolName(name) {
			if v, ok := a.poolAddr[name]; ok {
				return v, true
			}

			return 0, false
		}

		if isLocalRef(name) && refIdx >= 0 {
			v, ok := a.resolveLocal(name, refIdx)
			if !ok {
				return 0, false
			}

			memo[name] = v
			return v, true
		}

		if lr, ok := a.labels[name]; ok {
			v := a.labelAddr(lr)
			memo[name] = v
			return v, true
		}

		if e, ok := a.sets[name]; ok {
			if visiting[name] {
				return 0, false // cyclic definition
			}

			visiting[name] = true
			v, err := e.Eval(res)
			delete(visiting, name)
			if err != nil {
				return 0, false
			}

			memo[name] = uint64(v)
			return uint64(v), true
		}

		return 0, false
	}
	return res
}

// resolveLocal is the numeric local reference "Nb"/"Nf": the nearest
// definition in source order relative to statement refIdx - b at it or
// earlier (the labels of a line precede the instruction), f strictly later.
// The numLabels definitions are ordered by stmtIdx (registered in walk
// order).
func (a *assembler) resolveLocal(name string, refIdx int) (uint64, bool) {
	defs := a.numLabels[name[:len(name)-1]]
	if name[len(name)-1] == 'b' {
		for _, def := range slices.Backward(defs) {
			if def.stmtIdx <= refIdx {
				return a.labelAddr(def.ref), true
			}
		}

		return 0, false
	}

	for i := range defs {
		if defs[i].stmtIdx > refIdx {
			return a.labelAddr(defs[i].ref), true
		}
	}

	return 0, false
}

func (a *assembler) labelAddr(lr labelRef) uint64 {
	return a.secAddr[lr.sec] + uint64(a.secs[lr.sec].subs[lr.sub].base+lr.off)
}

// isNumericLabel is the name of a numeric local label: a non-empty sequence
// of digits (the definition "0:").
func isNumericLabel(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// isLocalRef is a reference name for a numeric local label: digits +
// 'b'/'f'. Named symbols cannot look like this (isIdentStart has no digits).
func isLocalRef(s string) bool {
	if len(s) < 2 || (s[len(s)-1] != 'b' && s[len(s)-1] != 'f') {
		return false
	}

	return isNumericLabel(s[:len(s)-1])
}
