package file

// The exported placement half of the universal writer - the "no
// contracts" rule. The assembler layers resolve label references through
// these functions; they run the same placement engine the writer itself
// uses, so the addresses baked into the code at assembly time are exactly
// where the bytes land in the image. This package depends on nothing
// above it: the policies speak plain numbers and file types, and the
// assembler layers adopt them structurally (a prog.Place, an asm layout
// callback).

import (
	"fmt"
	"strings"
)

// MachoPlaceStreams - where the two streams of a program image land:
// textSize is the text stream size, dataSize the data file bytes, dataMem
// the data memory size (>= dataSize: the bss tail). The result fits the
// Place signature of the prog layer as-is.
func MachoPlaceStreams(textSize, dataSize, dataMem int) (uint64, uint64) {
	specs := []machoSpec{
		{segment: "__TEXT", name: "__text", size: uint64(textSize), align: 4},
	}

	switch {
	case dataSize > 0:
		specs = append(specs, machoSpec{
			segment: "__DATA", name: "__data", size: uint64(dataSize), align: 8,
		})
		if tail := dataMem - dataSize; tail > 0 {
			specs = append(specs, machoSpec{
				segment: "__DATA", name: "__bss", size: uint64(tail), nobits: true, align: 8,
			})
		}
	case dataMem > 0:
		specs = append(specs, machoSpec{
			segment: "__DATA", name: "__bss", size: uint64(dataMem), nobits: true, align: 8,
		})
	}

	place, err := placeMachOSpecs(specs)
	if err != nil {
		// size-only specs of this shape are always placeable
		panic(fmt.Sprintf("macho: program streams: %v", err))
	}

	if len(place.sect) > 1 {
		return place.sect[0].addr, place.sect[1].addr
	}

	// no data at all: the empty data base is the memory continuation
	return place.sect[0].addr, machoVMAddr + place.vmTotal
}

// MachoPlaceSections - where the given image sections land: the engine
// NewMachOImage runs, exported for the code that resolved label
// references through a policy and now needs to translate absolute symbol
// addresses back into section offsets.
func MachoPlaceSections(sections []MachOSection) ([]uint64, error) {
	place, err := placeMachO(sections)
	if err != nil {
		return nil, err
	}

	out := make([]uint64, len(sections))
	for i := range out {
		out[i] = place.sect[i].addr
	}

	return out, nil
}

// MachoGasSpec - one gas-style source section of a layout query: its
// name (.text, .data, .bss, a custom .foo), its full memory size, and
// whether it is a NOBITS reserve. The vocabulary the asm core's layout
// callback converts into (the core knows no file formats, so the types
// meet here as plain strings and numbers).
type MachoGasSpec struct {
	Name   string
	Size   int
	Nobits bool
}

// NewMachoGasSpec - the spec of a gas-style section of size (or reserve).
func NewMachoGasSpec(name string, size int, nobits bool) MachoGasSpec {
	return MachoGasSpec{Name: name, Size: size, Nobits: nobits}
}

// MachoPlaceGas - the layout policy for gas-style sources: where the
// sections of an .s program land (see MachoGasSection for the name
// mapping). The asm layer's AssembleLayout callback delegates here.
func MachoPlaceGas(specs []MachoGasSpec) ([]uint64, error) {
	placed := make([]machoSpec, len(specs))
	for i := range specs {
		segment, name, align := machoGasMap(specs[i].Name)
		placed[i] = machoSpec{
			segment: segment,
			name:    name,
			size:    uint64(specs[i].Size),
			nobits:  specs[i].Nobits,
			align:   align,
		}
	}

	place, err := placeMachOSpecs(placed)
	if err != nil {
		return nil, err
	}

	out := make([]uint64, len(specs))
	for i := range out {
		out[i] = place.sect[i].addr
	}

	return out, nil
}

// MachoGasSection - the image section of a gas-style source section:
// .text lands in __TEXT (aligned 4, executable), everything else in
// __DATA (aligned 8, the leading dot stripped: .data becomes __data,
// .bss becomes __bss, a custom .foo becomes __foo). A positive reserve
// makes it a NOBITS section (memory without file bytes).
func MachoGasSection(name string, data []byte, reserve int) MachOSection {
	segment, sect, align := machoGasMap(name)

	return MachOSection{
		Segment: segment,
		Name:    sect,
		Data:    data,
		Nobits:  reserve,
		Align:   align,
	}
}

// machoGasMap - the gas-style section name to segment, section, and
// alignment.
func machoGasMap(gas string) (segment, name string, align int) {
	if gas == ".text" {
		return "__TEXT", "__text", 4
	}

	return "__DATA", "__" + strings.TrimPrefix(gas, "."), 8
}
