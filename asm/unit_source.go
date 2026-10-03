package asm

// The unit-mode assembly: a whole .S source becomes records of a
// *unit.Unit. The named labels of the source join the unit namespace;
// every section deposits as deferred runs (unit.Sym) that resolve at the
// unit's resolve phase - a name the source does not define (a C function
// called by bl) simply waits for the rest of the program, like every
// other deferred record of the unit. The internal layout (record sizes,
// the literal pools, the numeric locals) stays inside the source, as
// with fragments.
//
// v1: the sections map to the two unit streams by name - anything ending
// in ".text" is text, ".bss" is a zero-fill reserve, everything else is
// data; subsection numbers merge into their section. The .L-prefixed
// labels are local to the source: they deposit and resolve under the
// "file:name" key, so sources linked into one unit never collide on
// them; .global (and .globl) promotes its names into the shared
// namespace. .incbin and .ltorg are not supported.

import (
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"strings"

	parsecstrings "github.com/okneniz/parsec/strings"

	"github.com/okneniz/assembly/asm/expr"
	"github.com/okneniz/assembly/unit"
)

// srcRec is one record of a section: an instruction (in, with its pool
// slot), a data element list (exprs of width), ready bytes (fixed), or a
// zero-fill reserve of a NOBITS section. Exactly one shape is set; size
// is the placeholder-pass byte count (the sizes do not depend on symbol
// values, the fragment contract).
type srcRec struct {
	in      Unresolved
	poolIdx int // the section pool slot in reads (-1: none)
	exprs   []*expr.Expr
	width   int
	fixed   []byte
	reserve int
	size    int
	off     int
	line    int
}

// srcLabel is a named label at its insertion point (the index of the
// next record; past the end after the last one).
type srcLabel struct {
	name string
	rec  int
}

// srcDef is one numeric local definition: the record coordinate of the
// reference point (Nb/Nf pick the nearest by record order).
type srcDef struct {
	rec int
	off int
}

// srcPool is one literal pool slot of the section tail.
type srcPool struct {
	expr *expr.Expr
	size int
	line int
	off  int
}

// srcSection is one output section: its unit stream, its records, its
// named labels, its numeric locals, its literal pool.
type srcSection struct {
	name   string
	stream int
	nobits bool
	recs   []srcRec
	labels []srcLabel
	defs   map[string][]srcDef
	pool   []srcPool
	size   int
}

// srcBytes is ready data of the unit mode (the local mirror of the
// unit's own data blob).
type srcBytes []byte

func (b srcBytes) Encode(w io.Writer) (int64, error) {
	n, err := w.Write(b)
	return int64(n), err
}

// srcRun is a run of section records between two named-label
// boundaries: one deferred unit record. Its instructions and data
// resolve at their own addresses against the unit - the fragment
// discipline, grown to labels, sections and directives.
type srcRun struct {
	sec   *srcSection
	scope nameScope
	recs  []srcRec
	from  int // the section record index of the first record (the locals coordinate)
	off   int // the run's section offset
	size  int
}

// Resolve encodes the run at its final address: each record at its own
// address, numeric locals and pool slots against the section, .set
// names against the source, every other name against the unit (the
// labels of this source included - they were deposited as unit labels).
func (r *srcRun) Resolve(ctx unit.Ctx) ([]unit.Resolved, error) {
	base := ctx.Addr() - uint64(r.off)
	out := make([]unit.Resolved, 0, len(r.recs))

	for i := range r.recs {
		rec := &r.recs[i]
		addr := base + uint64(rec.off)
		resolve := r.sec.resolveNames(ctx, base, r.from+i, addr, rec.poolIdx, r.scope)

		switch {
		case rec.in != nil:
			res, err := rec.in.Resolve(unit.NewCtx(addr, resolve))
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", rec.line, err)
			}

			out = append(out, res)
		case rec.exprs != nil:
			b := make([]byte, 0, rec.size)
			for _, e := range rec.exprs {
				v, err := e.Eval(resolve)
				if err != nil {
					return nil, fmt.Errorf("line %d: %w", rec.line, err)
				}

				b = appendIntLE(b, rec.width, v)
			}

			out = append(out, srcBytes(b))
		default:
			out = append(out, srcBytes(rec.fixed))
		}
	}

	return out, nil
}

// Size is the run's byte count (the placeholder sizing of the build).
func (r *srcRun) Size() int {
	return r.size
}

// nameScope is the file-scoped name environment of a source: the .set
// expressions, the .global-promoted names, and the file key that
// isolates the .L locals of different sources linked into one unit.
type nameScope struct {
	file    string
	sets    map[string]*expr.Expr
	globals map[string]bool
}

// linkName is the unit name of a source symbol: a .L local sits behind
// its file key (invisible to the other sources of the link), everything
// else - and a .global-promoted name whatever its spelling - joins the
// shared namespace.
func (n nameScope) linkName(name string) string {
	if strings.HasPrefix(name, ".L") && !n.globals[name] {
		return n.file + ":" + name
	}

	return name
}

// unitSource is the interpreted source: the sections in first-appearance
// order, the shared .set expressions, the .global-promoted names, and
// the file key of the .L locals' namespace.
type unitSource struct {
	secs    []*srcSection
	sets    map[string]*expr.Expr
	file    string
	globals map[string]bool
}

// scope is the name environment the source's runs resolve in.
func (s *unitSource) scope() nameScope {
	return nameScope{
		file:    s.file,
		sets:    s.sets,
		globals: s.globals,
	}
}

// sectionFor is the section of a switch: an existing one by name, a new
// one placed by the stream rule.
func (s *unitSource) sectionFor(name string) *srcSection {
	for _, sec := range s.secs {
		if sec.name == name {
			return sec
		}
	}

	stream, nobits := unitStreamOf(name)
	sec := newSrcSection(name, stream, nobits)
	s.secs = append(s.secs, sec)
	return sec
}

// appendRec appends a laid-out record at the section end.
func (s *srcSection) appendRec(rec srcRec) {
	rec.off = s.size
	s.size += rec.size
	s.recs = append(s.recs, rec)
}

// appendReserve appends n zero bytes - file zeros of a regular section,
// a NOBITS reserve otherwise.
func (s *srcSection) appendReserve(n int, pos parsecstrings.Position) {
	if n == 0 {
		return
	}

	if s.nobits {
		s.appendRec(srcRec{
			reserve: n,
			size:    n,
			off:     s.size,
			line:    int(pos.Line()) + 1,
		})
		return
	}

	s.appendRec(srcRec{
		fixed: make([]byte, n),
		size:  n,
		off:   s.size,
		line:  int(pos.Line()) + 1,
	})
}

// AssembleUnit assembles a whole .S source into the unit: the named
// labels join the unit's namespace (the .L locals behind the file key),
// the sections deposit as deferred runs - a name the source does not
// define (a C function called by bl) waits for the program's resolve
// phase, like every other deferred record. file names the origin in the
// unit's line map. It is the single-source shorthand of ParseSourceUnit
// + Deposit + DepositBss (see source_unit.go).
func AssembleUnit(u *unit.Unit, file, src string, be Syntax) []AsmError {
	su, errs := ParseSourceUnit(file, src, be)
	if len(errs) > 0 {
		return errs
	}

	su.Deposit(u)
	su.DepositBss(u)
	return nil
}

func newSrcRec() srcRec {
	return srcRec{poolIdx: -1}
}

func newSrcSection(name string, stream int, nobits bool) *srcSection {
	return &srcSection{
		name:   name,
		stream: stream,
		nobits: nobits,
		defs:   map[string][]srcDef{},
	}
}

// unitStreamOf is the v1 section rule: a name ending in ".text" is the
// text stream, ".bss" a zero-fill reserve, everything else data.
func unitStreamOf(name string) (stream int, nobits bool) {
	switch {
	case name == ".bss" || strings.HasSuffix(name, ".bss"):
		return 1, true
	case name == ".text" || strings.HasSuffix(name, ".text"):
		return 0, false
	default:
		return 1, false
	}
}

// buildUnitSource interprets the statements into the source: records
// sized under the placeholder environment, labels, pools, sets. The
// directive semantics mirror the byte mode one to one; the implicit
// .text materializes lazily - a source opening with another section
// deposits that one first.
func buildUnitSource(stmts []statement, be Syntax) (*unitSource, []AsmError) {
	src := &unitSource{
		sets:    map[string]*expr.Expr{},
		globals: map[string]bool{},
	}
	var errs []AsmError

	var sec *srcSection
	fail := func(pos parsecstrings.Position, format string, args ...any) {
		errs = append(errs, NewAsmError(pos.Line()+1, 0, fmt.Sprintf(format, args...)))
	}

	for i := range stmts {
		st := &stmts[i]
		if st.err != nil {
			errs = append(errs, *st.err)
			continue
		}

		if st.directive != nil {
			if len(st.labels) > 0 {
				if sec == nil {
					sec = src.sectionFor(".text")
				}

				defineLabels(sec, st.labels)
			}

			sec = doUnitDirective(src, sec, st.directive, be, st.pos, fail)
			continue
		}

		if len(st.labels) == 0 && !st.hasInstr {
			continue
		}

		if sec == nil {
			sec = src.sectionFor(".text")
		}

		defineLabels(sec, st.labels)

		if !st.hasInstr {
			continue
		}

		if sec.nobits {
			fail(st.pos, "section %s is NOBITS: instructions are not permitted", sec.name)
			continue
		}

		rec := newSrcRec()
		rec.in = st.instr
		rec.line = int(st.pos.Line()) + 1
		if pu, ok := st.instr.(PoolUser); ok {
			if e, slot, want := pu.PoolReq(); want {
				rec.poolIdx = sec.poolAdd(e, slot, rec.line)
			}
		}

		size, err := sizeOf(st.instr, unit.NewCtx(0, placeholderResolve(0)))
		if err != nil {
			fail(st.pos, "%v", err)
			continue
		}

		rec.size = size
		rec.off = sec.size
		sec.size += size
		sec.recs = append(sec.recs, rec)
	}

	// the literal pools sit at the tail of their section, in
	// first-appearance order, every slot at its natural alignment (the
	// byte mode layout): the zero gap before a wider slot is a fixed
	// record — the `udf #0` fill words of the gas dump
	for _, s := range src.secs {
		for i := range s.pool {
			p := &s.pool[i]
			pad := alignTo(s.size, p.size) - s.size
			if pad > 0 {
				s.recs = append(s.recs, srcRec{
					fixed: make([]byte, pad),
					size:  pad,
					off:   s.size,
					line:  p.line,
				})
				s.size += pad
			}

			p.off = s.size
			s.recs = append(s.recs, srcRec{
				exprs: []*expr.Expr{p.expr},
				width: p.size,
				size:  p.size,
				off:   s.size,
				line:  p.line,
			})
			s.size += p.size
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return src, nil
}

// defineLabels records the statement's labels at the section's current
// end: numeric locals into the redefinable table, named ones as run
// boundaries (a label before a section-switching directive stays in the
// section the walk is leaving).
func defineLabels(sec *srcSection, labels []string) {
	point := len(sec.recs)
	for _, lbl := range labels {
		if isNumericLabel(lbl) {
			sec.defs[lbl] = append(sec.defs[lbl], srcDef{rec: point, off: sec.size})
			continue
		}

		sec.labels = append(sec.labels, srcLabel{name: lbl, rec: point})
	}
}

// doUnitDirective interprets one directive of the unit mode (the byte
// mode semantics; the differences are the file-wide v1 notes at the top)
// and returns the section the walk continues in.
func doUnitDirective(
	src *unitSource,
	sec *srcSection,
	d *directive,
	be Syntax,
	pos parsecstrings.Position,
	fail func(parsecstrings.Position, string, ...any),
) *srcSection {
	switch d.name {
	case ".text", ".data", ".bss":
		return src.sectionFor(d.name) // subsection numbers merge (v1)
	case ".section":
		if len(d.args) > 0 {
			return src.sectionFor(d.args[0].str)
		}

		return sec
	case ".global", ".globl":
		for _, arg := range d.args {
			src.globals[arg.str] = true
		}

		return sec
	}

	if sec == nil {
		sec = src.sectionFor(".text")
	}

	switch d.name {
	case ".ltorg":
		fail(pos, ".ltorg is not supported: literal pool is emitted at the end of the section")
	case ".set", ".equ":
		if len(d.args) >= 2 {
			src.sets[d.args[0].str] = d.args[1].expr
		}
	case ".option":
		if len(d.args) > 0 {
			if err := be.ApplyOption(d.args[0].str); err != nil {
				fail(pos, ".option: %v", err)
			}
		}
	case ".word", ".half", ".short", ".byte", ".quad", ".dword":
		if sec.nobits {
			fail(pos, "section %s is NOBITS: only .zero/.skip/.align permitted", sec.name)
			return sec
		}

		width := map[string]int{
			".word": 4, ".half": 2, ".short": 2,
			".byte": 1, ".quad": 8, ".dword": 8,
		}[d.name]
		sec.appendRec(srcRec{
			exprs: dirExprs(d),
			width: width,
			size:  width * len(d.args),
			off:   sec.size,
			line:  int(pos.Line()) + 1,
		})
	case ".zero", ".space", ".skip":
		n, err := d.args[0].expr.Eval(nil) // the size must be a constant
		if err != nil || n < 0 {
			fail(pos, "%s: constant non-negative size required", d.name)
			return sec
		}

		sec.appendReserve(int(n), pos)
	case ".align", ".p2align", ".balign":
		n, err := d.args[0].expr.Eval(nil)
		if err != nil || n < 0 {
			fail(pos, "%s: constant alignment required", d.name)
			return sec
		}

		align := uint64(1) << uint(n)
		if d.name == ".balign" {
			align = uint64(n)
		}

		if align == 0 {
			if d.name == ".balign" {
				fail(pos, ".balign: zero alignment")
			} else {
				fail(pos, "%s: alignment 2^%d too large", d.name, n)
			}

			return sec
		}

		off := uint64(sec.size)
		pad := int((align - off%align) % align)
		sec.appendReserve(pad, pos)
	case ".incbin":
		fail(pos, ".incbin is not supported in the unit mode")
	case ".string", ".asciz", ".ascii":
		if sec.nobits {
			fail(pos, "section %s is NOBITS: only .zero/.skip/.align permitted", sec.name)
			return sec
		}

		var b []byte
		for _, arg := range d.args {
			b = append(b, arg.str...)
			if d.name != ".ascii" {
				b = append(b, 0)
			}
		}

		sec.appendRec(srcRec{
			fixed: b,
			size:  len(b),
			off:   sec.size,
			line:  int(pos.Line()) + 1,
		})
	}

	// .type/.size/... - recognized, carry no semantics (v1)
	return sec
}

// dirExprs is the expression list of an argsExprs directive.
func dirExprs(d *directive) []*expr.Expr {
	out := make([]*expr.Expr, 0, len(d.args))
	for _, arg := range d.args {
		out = append(out, arg.expr)
	}

	return out
}

// appendIntLE appends v as width little-endian bytes.
func appendIntLE(b []byte, width int, v int64) []byte {
	switch width {
	case 1:
		return append(b, byte(v))
	case 2:
		var t [2]byte
		binary.LittleEndian.PutUint16(t[:], uint16(v))
		return append(b, t[:]...)
	case 4:
		var t [4]byte
		binary.LittleEndian.PutUint32(t[:], uint32(v))
		return append(b, t[:]...)
	default:
		var t [8]byte
		binary.LittleEndian.PutUint64(t[:], uint64(v))
		return append(b, t[:]...)
	}
}

// depositUnit lays the source into the unit, one section class at a
// time: bss=false deposits the text and data sections (in
// first-appearance order), bss=true the NOBITS ones. A linker deposits
// the former of every source first and the latter of all of them after -
// the unit's data stream carries one bss tail, so the zero-fill reserves
// of the whole program aggregate at its end. Per section the stream is
// switched, the named labels are defined at their boundaries through the
// source's name scope (the .L locals behind their file key), the record
// runs between them deposit as one deferred record each (a NOBITS
// section reserves instead).
func depositUnit(u *unit.Unit, src *unitSource, bss bool) {
	scope := src.scope()

	for _, sec := range src.secs {
		if sec.nobits != bss {
			continue
		}

		if sec.stream == 0 {
			u.Text()
		} else {
			u.Data()
		}

		pos := func(line int) unit.Pos {
			return unit.NewPos(src.file, line)
		}

		bound := 0
		for _, lbl := range sec.labels {
			sec.depositSpan(u, scope, bound, lbl.rec, pos)
			u.Label(scope.linkName(lbl.name))
			bound = lbl.rec
		}

		sec.depositSpan(u, scope, bound, len(sec.recs), pos)
	}
}

// depositSpan deposits the records [from, to) of the section: one
// deferred run, or the NOBITS reserves.
func (s *srcSection) depositSpan(
	u *unit.Unit,
	scope nameScope,
	from, to int,
	pos func(int) unit.Pos,
) {
	if from >= to {
		return
	}

	if s.nobits {
		for _, rec := range s.recs[from:to] {
			u.Bss(pos(rec.line), rec.reserve)
		}

		return
	}

	recs := s.recs[from:to]
	u.Sym(pos(recs[0].line), &srcRun{
		sec:   s,
		scope: scope,
		recs:  recs,
		from:  from,
		off:   recs[0].off,
		size:  recs[len(recs)-1].off + recs[len(recs)-1].size - recs[0].off,
	})
}

// localAddr is "Nb"/"Nf": the nearest definition by record order - b at
// the record itself or earlier (a label on the same line precedes the
// instruction), f strictly later.
func (s *srcSection) localAddr(name string, rec int, base uint64) (uint64, bool) {
	defs := s.defs[name[:len(name)-1]]
	if name[len(name)-1] == 'b' {
		for _, d := range slices.Backward(defs) {
			if d.rec <= rec {
				return base + uint64(d.off), true
			}
		}

		return 0, false
	}

	for _, d := range defs {
		if d.rec > rec {
			return base + uint64(d.off), true
		}
	}

	return 0, false
}

// poolAdd registers a literal slot of the section (dedup by the pool
// name, as the byte mode) and returns its index.
func (s *srcSection) poolAdd(e *expr.Expr, slot int, line int) int {
	name := poolName(slot, expr.ExprKey(e))
	for i := range s.pool {
		if poolName(s.pool[i].size, expr.ExprKey(s.pool[i].expr)) == name {
			return i
		}
	}

	s.pool = append(s.pool, srcPool{expr: e, size: slot, line: line})
	return len(s.pool) - 1
}

// resolveNames is the name resolver of one record: numeric locals
// against the section, "." the record itself, PoolSelf the record's pool
// slot, .set names the source's expressions (cycle-guarded), the rest
// the unit (through the source's name scope - the .L locals behind
// their file key).
func (s *srcSection) resolveNames(
	ctx unit.Ctx,
	base uint64,
	rec int,
	addr uint64,
	poolIdx int,
	scope nameScope,
) func(string) (uint64, bool) {
	visiting := map[string]bool{}

	var resolve func(string) (uint64, bool)
	resolve = func(name string) (uint64, bool) {
		switch {
		case isLocalRef(name):
			return s.localAddr(name, rec, base)
		case name == ".":
			return addr, true
		case name == PoolSelf && poolIdx >= 0:
			return base + uint64(s.pool[poolIdx].off), true
		}

		if e, ok := scope.sets[name]; ok {
			if visiting[name] {
				return 0, false
			}

			visiting[name] = true
			defer delete(visiting, name)

			v, err := e.Eval(resolve)
			if err != nil {
				return 0, false
			}

			return uint64(v), true
		}

		return ctx.Resolve(scope.linkName(name))
	}

	return resolve
}

// compile-time: a section run is a deferred record of the unit output.
var _ unit.Sym = (*srcRun)(nil)
