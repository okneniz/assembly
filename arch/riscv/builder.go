package riscv

// Builder - the instruction vocabulary of the package: per-instruction
// constructor methods named after the mnemonic they build (Addi, Lw,
// Bl...). Stateless by design — a namespace, not a factory.
type Builder struct{}

// New - the Builder for instruction constructors.
func New() Builder {
	return Builder{}
}
