package loong64

// Builder — the construction vocabulary of the package: instruction
// constructors named after the mnemonic they build and operand role
// constructors. Stateless by design — a namespace, not a factory.
type Builder struct{}

// New — the Builder.
func New() Builder {
	return Builder{}
}
