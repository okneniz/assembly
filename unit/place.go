package unit

// Place is the placement policy of the two unit streams: given the text
// size, the data file size, and the data memory size (>= the file size:
// the bss tail), it returns the base address of each stream. The policy
// owns the layout knowledge - file.MachoPlaceStreams is the Mach-O one,
// built on the same engine as the writer, so deferred records resolve
// against exactly the addresses the image will use and no caller ever
// computes one.
type Place func(textSize, dataSize, dataMem int) (textAddr, dataAddr uint64)
