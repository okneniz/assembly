package file

// Section describes one section of a binary file.
type Section struct {
	Name    string
	Segment string
	Addr    uint64
	Offset  uint64
	Size    uint64
	Data    []byte
}

func NewSection(
	name string,
	segment string,
	addr uint64,
	offset uint64,
	size uint64,
	data []byte,
) *Section {
	return &Section{
		Name:    name,
		Segment: segment,
		Addr:    addr,
		Offset:  offset,
		Size:    size,
		Data:    data,
	}
}
