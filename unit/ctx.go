package unit

// Ctx is the evaluation environment of a deferred record: the address
// the record lands at (its own pc) and the symbol resolver. The resolve
// phase of Unit builds one per record; the text assembler (asm) builds
// its own against its symbol table - the vocabulary is shared, the
// tables are not.
type Ctx interface {
	// Addr is the absolute address of the record itself.
	Addr() uint64

	// Resolve is the value of a symbol by name; ok=false for unknown ones.
	Resolve(name string) (uint64, bool)
}

// NewCtx is a Ctx at addr with the resolver (a nil resolver resolves
// nothing).
func NewCtx(addr uint64, resolve func(string) (uint64, bool)) Ctx {
	return addrCtx{addr: addr, resolve: resolve}
}

// addrCtx is the concrete Ctx.
type addrCtx struct {
	addr    uint64
	resolve func(string) (uint64, bool)
}

func (c addrCtx) Addr() uint64 {
	return c.addr
}

func (c addrCtx) Resolve(name string) (uint64, bool) {
	if c.resolve == nil {
		return 0, false
	}

	return c.resolve(name)
}
