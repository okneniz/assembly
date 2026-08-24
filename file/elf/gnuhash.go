package elf

import "errors"

// ErrNoGnuHash means the file has no DT_GNU_HASH (no .gnu.hash table).
var ErrNoGnuHash = errors.New("elf: no DT_GNU_HASH")

// ErrNoSysVHash means the file has no DT_HASH (no classic .hash table).
var ErrNoSysVHash = errors.New("elf: no DT_HASH")

// GnuHash is a parsed .gnu.hash table (DT_GNU_HASH): a hash table for fast
// symbol lookup at runtime. The first SymOffset .dynsym symbols (local ones,
// without hashes) are not included in the table.
//
// Layout:
//
//	nbuckets(4) symoffset(4) bloom_size(4) bloom_shift(4)
//	bloom[bloom_size] × uint64 (uint32 for ELF32)
//	buckets[nbuckets] × uint32
//	chains[] × uint32 - one array per hashed symbol; the end of a chain is
//	the element with the low bit set.
type GnuHash struct {
	NBuckets   uint32
	SymOffset  uint32
	BloomSize  uint32
	BloomShift uint32
	Bloom      []uint64
	Buckets    []uint32
	Chains     []uint32
}

func NewGnuHash(nBuckets uint32, symOffset uint32, bloomSize uint32, bloomShift uint32) *GnuHash {
	return &GnuHash{
		NBuckets:   nBuckets,
		SymOffset:  symOffset,
		BloomSize:  bloomSize,
		BloomShift: bloomShift,
	}
}

// SysVHash is the .hash table (DT_HASH, the classic SysV format).
//
//	nbucket(4) nchain(4)
//	buckets[nbucket] × uint32
//	chains[nchain] × uint32 - parallel to the .dynsym entries
type SysVHash struct {
	NBucket uint32
	NChain  uint32
	Buckets []uint32
	Chains  []uint32
}

func NewSysVHash(nBucket uint32, nChain uint32) *SysVHash {
	return &SysVHash{
		NBucket: nBucket,
		NChain:  nChain,
	}
}

// GnuHash parses .gnu.hash; an absent table yields ErrNoGnuHash. There is no
// section - the address is taken from DT_GNU_HASH and mapped through PT_LOAD.
func (f *File) GnuHash() (*GnuHash, error) {
	v, ok := f.dynamicValue(DT_GNU_HASH)
	if !ok {
		return nil, ErrNoGnuHash
	}

	off, err := f.vaddrToOff(v)
	if err != nil {
		return nil, err
	}

	is64 := f.class == CLASS64
	hdr, err := readBytes(f.buf, int(off), 16)
	if err != nil {
		return nil, err
	}

	h := NewGnuHash(
		u32(hdr, f.order),
		u32(hdr[4:], f.order),
		u32(hdr[8:], f.order),
		u32(hdr[12:], f.order),
	)
	if h.NBuckets == 0 {
		return nil, errf("gnu.hash: nbuckets = 0")
	}

	bloomBytes := h.BloomSize * 4
	if is64 {
		bloomBytes = h.BloomSize * 8
	}

	pos := int(off) + 16
	bloomRaw, err := readBytes(f.buf, pos, int(bloomBytes))
	if err != nil {
		return nil, err
	}

	if is64 {
		h.Bloom = make([]uint64, h.BloomSize)
		for i := range h.Bloom {
			h.Bloom[i] = u64(bloomRaw[i*8:], f.order)
		}
	} else {
		h.Bloom = make([]uint64, h.BloomSize)
		for i := range h.Bloom {
			h.Bloom[i] = uint64(u32(bloomRaw[i*4:], f.order))
		}
	}

	pos += int(bloomBytes)
	bucketsRaw, err := readBytes(f.buf, pos, int(h.NBuckets)*4)
	if err != nil {
		return nil, err
	}

	h.Buckets = make([]uint32, h.NBuckets)
	for i := range h.Buckets {
		h.Buckets[i] = u32(bucketsRaw[i*4:], f.order)
	}

	// The chains tail is read up to the end of the segment containing the
	// table (the format has no field with the size). The upper bound is the
	// segment's filesz.
	chainLen, err := f.gnuHashChainLen(v, pos+int(h.NBuckets)*4)
	if err != nil {
		return nil, err
	}

	chainsRaw, err := readBytes(f.buf, pos+int(h.NBuckets)*4, chainLen)
	if err != nil {
		return nil, err
	}

	h.Chains = make([]uint32, len(chainsRaw)/4)
	for i := range h.Chains {
		h.Chains[i] = u32(chainsRaw[i*4:], f.order)
	}

	return h, nil
}

// gnuHashChainLen computes the size of the chains array: up to the end of the
// SHT_GNU_HASH section at this address, otherwise up to the end of the
// containing PT_LOAD segment.
func (f *File) gnuHashChainLen(vaddr uint64, startOff int) (int, error) {
	for _, s := range f.sections {
		if s.Type == SHT_GNU_HASH && s.Addr == vaddr && s.Off <= uint64(startOff) && s.Size > 0 {
			end := int(s.Off + s.Size)
			if end >= startOff {
				return ((end - startOff) / 4) * 4, nil
			}
		}
	}

	for _, p := range f.progs {
		if p.Type != PT_LOAD {
			continue
		}

		if vaddr >= p.Vaddr && vaddr < p.Vaddr+p.Filesz {
			end := int(p.Off + p.Filesz)
			n := max((end-startOff)/4, 0)
			return n * 4, nil
		}
	}

	return 0, errf("gnu.hash: table is not covered by PT_LOAD")
}

// SysVHash parses .hash; an absent table yields ErrNoSysVHash.
func (f *File) SysVHash() (*SysVHash, error) {
	v, ok := f.dynamicValue(DT_HASH)
	if !ok {
		return nil, ErrNoSysVHash
	}

	off, err := f.vaddrToOff(v)
	if err != nil {
		return nil, err
	}

	hdr, err := readBytes(f.buf, int(off), 8)
	if err != nil {
		return nil, err
	}

	h := NewSysVHash(u32(hdr, f.order), u32(hdr[4:], f.order))
	if h.NBucket == 0 {
		return nil, errf(".hash: nbucket = 0")
	}

	data, err := readBytes(f.buf, int(off)+8, int(h.NBucket+h.NChain)*4)
	if err != nil {
		return nil, err
	}

	h.Buckets = make([]uint32, h.NBucket)
	for i := range h.Buckets {
		h.Buckets[i] = u32(data[i*4:], f.order)
	}

	h.Chains = make([]uint32, h.NChain)
	for i := range h.Chains {
		h.Chains[i] = u32(data[(int(h.NBucket)+i)*4:], f.order)
	}

	return h, nil
}
