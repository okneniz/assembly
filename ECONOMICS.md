# Economics

A log of build statistics. An entry is appended after every big
feature; the numbers come from git (lines, files) and the tooling's
own usage database (sessions, requests, tokens).

## 08-09 … 08-24 — the base build

The project core: two architectures (ARM64, RISC-V), the assembler
core, ELF/Mach-O parsers, the server, the full test infrastructure.
Built as a human–AI collaboration: the entire codebase was written by
an AI coding agent, with the human owner setting direction, reviewing
architecture and making the final calls.

## 08-28 … 08-29 — the LoongArch backend

The third architecture: the full scalar integer set (248
instructions), the assembly layer with pseudo-instructions, ELF/CLI/
server wiring, the property suite, and the rt-vm matrix (three VMs).
Orchestrated differently than the base: nine subagent runs in
parallel (seven instruction slices, one rate-limit retry, one
property suite), each obliged to verify every encoding against
`llvm-mc` before writing tests.

## 09-03 … 09-08 — prog DSL, Mach-O writer, llvm-mc parity, arch exodus

Three waves. First the prog DSL (a Go-level programming interface for
machine code: every chain method is a source line, labels resolve at
assembly time) plus the Mach-O executable writer (arm64 macOS, ad-hoc
signed, no linker). Then the encoding parity push: every hello example
must assemble into exactly the bytes llvm-mc chooses — this forced the
RISC-V layout relaxation (symbolic jumps compress to RVC when the
distance fits), the loong la.pcrel/la.abs spellings, and the la.abs
64-bit ladder. Finally the arch exodus: all assembler-only code moved
out of arch/arm64 into asm/arm64 (~5,400 lines — the legacy fallback,
the 14 ctor files, the armCtors dispatch, the by-element registrations;
the Builder vocabulary was completed first, then the ctor files were
rewritten to use it, then moved; Ctor tests renamed to Build).

## 09-06 … 09-08 — the arm64 ISA-audit closeout

The arm64 decode-integrity tail of the exodus program. The five
hand-written-vs-XML audit conflicts resolved per case (llvm verdicts:
the hand-written constants were transcriptions of other instructions —
cmtst/cmge/fcsel under wrong names; a latent Fp3 type-bit bug surfaced
with them). The Advanced SIMD copy family completed: dup/ins (element),
the scalar mov alias, smov/umov, and the MOV (from vector) input alias
to full llvm print/input parity; the stale MovElem schema deleted and
the legacy mov fallback hardened against vector operands. Audit at
conflict=0; docker and rt-vm gates green throughout.

## 09-09 — the style unification

The whole codebase brought to one style under five owner rules: structs
are created by constructors; literals name their fields; complex
literals are multi-line; no vars outside functions; no heavy
computation at call sites (parsec combinators nested inside calls,
rebuilt per invocation). The rule-4/rule-5 conflict resolves with the
owner's patterns: non-recursive combinators are built as LOCAL vars
captured by the returned closure; mutually recursive grammars live in
structs whose fields are assigned in the constructor (asm/expr's
precedence ladder, the arm64 operand grammar, every backend). What
remains at package level is data Go cannot make const: the generated
and literal tables (built once by named constructors — armCtors now
includes the by-element registrations without init() mutation), the Err
sentinels, the register singletons of the public DSL, and the decode
trees (the DecodeWord export and the assembler's self-verify need them
globally). The per-call construction hot spots are gone: arch word
readers are plain byte loops (bytes.ReadAs built a fresh Count per
word), the elf/macho field readers likewise, and the rt harness holds
one shared backend instead of assembling every line with a fresh
grammar. Cost of the only API break (pre-1.0): the four exported expr
combinator vars became functions. Net effect: arm64 Parse +67%
instr/s at −77% allocations, riscv Parse ×5.1, the rt corpus run is
faster than on main. The second pass set the vocabulary and file order:
bytes combinators speak decoder (arch Parse→MakeDecoder, parser.go→
decoder.go), rune combinators speak parser (every grammar), constructors
carry the make prefix and combinator values are imperative (parseWord,
decodeHalfLE); Backend methods live in backend.go, per-instruction types
in instr.go, and the lint gate — funcorder, decorder, gochecknoinits
newly enabled — runs at zero issues. The review pass finished the constructor
program: every instruction struct is assembled ONLY in its newX constructor
(the Builder method validates and delegates, the decoder extracts fields and
delegates; the pass-through Builder{}.X calls are gone, replaced by exported
NewX constructors for the assembler layer), and the per-instruction file order
is type, constructor, methods, Builder method, decoder.

## Totals

| Period | What | Days | Go lines | Go files | Sessions | Requests | Tokens in | Tokens out |
|---|---|---|---|---|---|---|---|---|
| 08-09 … 08-24 | base build: ARM64 + RISC-V, asm core, ELF/Mach-O, server, tests | 16 | 79,562 | 440 | 117 | 9,830 | ~2.32B | ~7.06M |
| 08-28 … 08-29 | LoongArch: 248 instructions, pseudo layer, property suite, 3 VMs | 2 | 25,644 | 558 | 21 | 1,215 | ~321M | ~1.02M |
| 09-03 … 09-08 | prog DSL, Mach-O writer, llvm-mc parity, arch exodus | 5 | 9,650 | 202 | 3 | ~90 | — | — |
| 09-06 … 09-08 | arm64 ISA-audit closeout: 5 conflicts, SIMD copy family, MOV alias | 3 | 534 | 16 | 3 | — | — | — |
| 09-09 | style unification: constructors, grammar structs, decoder/parser vocabulary, lint at zero | 1 | +1,639 | 200 | 1 | — | — | — |
| **total** | | **27** | **117,029** | **1,349** | **145** | **~11,135** | **~2.64B** | **~8.08M** |

### Cost

In reality the whole run cost a flat monthly coding subscription —
no per-token billing. For reference only: the same volume at
pay-per-token API rates would have been roughly ~$200 (budget models)
to ~$5,200 (frontier models) for the base, ~$30–$740 for LoongArch,
almost all of it cache reads at their discounted rate.

Fun ratio for the base: the full pipeline burned ~29K tokens per line
of code — that is what it costs to weigh, test and review every line;
the pure output cost is only ~87 tokens per line. LoongArch halved it
to ~12K per line (highly parallel, template clones).
