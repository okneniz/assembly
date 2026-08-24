# dtree/

A decision tree over mask rules: given a word, it selects the first matching
rule (first-match-wins) without scanning the list linearly. The production
`decodeOne` of both arches runs on it: on the ARM64 decoder registry (~3.7k
rules: schemes + generated isaTail), pure decode sped up from 120k to 7.5M
instructions/sec (x62), full Parse from 115k to 3.4M (x30); for RISC-V, from
664k to 3.4M (x5). Benchmarks: `arch/arm64/decode_bench_test.go` and
`arch/*/parser_bench_test.go` (full branch coverage).

## Mask and Match

`Mask` is which bits of the word are significant, `Match` is which values are
expected on them: a rule covers all words whose significant bits equal `Match`,
the remaining bits are "don't care". The rule
`{Mask: 0b1100_0000, Match: 0b1000_0000}` is the entire class `10xxxxxx`
(64 values). This is not an "instruction" but a "word -> class by bit mask"
classifier, where classes may intersect and the order of rules defines
priority.

A word is any 32-bit field, not a machine instruction:

- **encodings/lexers**: UTF-8 byte classification - `{0xC0, 0x80}` "continuation
  byte", `{0xE0, 0xC0}` "two-byte leading" - a transition table in four
  rules;
- **network classifiers**: an IPv4 header packed into a uint32 - "protocol is
  TCP and fragmentation forbidden": a mask over the protocol bits and the DF
  flag;
- **VM opcode dispatchers**: the high bits of the opcode select the handler,
  the low ones are operand fields;
- **protocol status bytes** (MIDI, CAN): `{0xF0, 0x90}` is "note-on on any
  channel"; device status registers where bit fields encode the error class.

Analogy: a software TCAM (associative memory with "don't care" bits, on which
network ACLs are built) or a glob for bits - `10??????`.

## Semantics

Construction with `New` preserves the semantics of the linear list unchanged:

- nodes split the group by bits covered by the masks of all rules and
  distinguishing them;
- if there are none, by a single bit covered by some of the rules (rules
  without the bit in their mask go into both branches, the original order is
  preserved in a single pass);
- rules not distinguishable by any bit (an intersection) remain a leaf list
  with priority by order.

`Lookup` returns the payload of the first (in `New` order) rule the word
belongs to; masking entries have a zero payload. Correctness is upheld by the
oracle test (tree against list on all 16-bit words) and the `decodeOne`
equivalence check on full coverage in the benchmark.

## Benchmarks (baseline, 2026-08-23)

The material is `arb.LookupCase.Generate()` (1000 cases): a "registry" of ~12k
live rules normalized to the common defined fields sf/op/opc (without
normalization, random broad masks cover each bit about half the time -
duplicates are always above the 1/4 threshold and the tree degenerates into a
leaf list); words are half miss/hit. The registry tree: 33 nodes, 32 leaves,
depth 2. The seed is `ASSEMBLY_SEED` (42).

go1.23.2, darwin/arm64. Reproduction:

	go test -bench=. -benchmem -run '^$' ./dtree/

| benchmark                         | ns/op     | throughput      | B/op / allocs |
|-----------------------------------|-----------|-----------------|---------------|
| New - random set of <=40 rules    | 2 202     | 454k tree/s     | 78 / 2        |
| New - registry of ~12k rules      | 1 597 352 | 7.5M rule/s     | 530 236 / 360 |
| Lookup - registry                 | 56        | 17.8M lookup/s  | 0 / 0         |
| Nodes / Leaves / MaxDepth - walk  | 467-489   | -               | 0 / 0         |

## Status

Used by the production decoder of both arches (arm64 - the scheme trees and
the isaTail tail with nil-ctor barriers, riscv - the decodeTable tree). A
candidate for moving to `parsec` - "mask-based parsing" is broader than a
single ISA.
