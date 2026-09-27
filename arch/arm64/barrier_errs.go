package arm64

import "fmt"

// unknownBarrier / unknownDomain - the Barrier constructor errors.
func unknownBarrier(name string) error {
	return fmt.Errorf("arm64: unknown barrier %q", name)
}

func unknownDomain(name string, domain BarrierDomain) error {
	return fmt.Errorf("arm64: %s: unknown domain %#x", name, uint8(domain))
}
