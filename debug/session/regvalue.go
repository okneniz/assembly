package session

// RegValue is one register of a display dump.
type RegValue struct {
	Name  string
	Value uint64
}

func NewRegValue(name string, value uint64) RegValue {
	return RegValue{
		Name:  name,
		Value: value,
	}
}
