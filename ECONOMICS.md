# Economics

A log of build statistics. An entry is appended after every big
feature; the numbers come from git (lines, files) and the tooling's
own usage database (sessions, requests, tokens).

## 2026-08-09 … 08-24 — the base build

The project core: two architectures (ARM64, RISC-V), the assembler
core, ELF/Mach-O parsers, the server, the full test infrastructure.
Built as a human–AI collaboration: the entire codebase was written by
an AI coding agent, with the human owner setting direction, reviewing
architecture and making the final calls. **16 days.**

| Metric | Value |
|---|---|
| Hand-written Go code | 79,562 lines (70,110 production + 9,452 tests) in 440 files |
| Other hand-written code (C / asm / shell / Makefile) | ~1,100 lines |
| Generated tables & vendored data | ~36,000 lines |
| AI sessions | 117 |
| Model requests | 9,830 |

### Token usage (actual, measured)

| Category | Tokens |
|---|---|
| Input (total) | ~2.32B |
| — of which cache reads | ~2.30B (99%) |
| Output | ~7.06M |
| Total processed | ~2.33B |

Fun ratio: the full pipeline burned ~29K tokens per line of code —
that is what it costs to weigh, test and review every line; the pure
output cost is only ~87 tokens per line.

### Cost

In reality the whole run cost a flat monthly coding subscription —
no per-token billing. For reference only: the same volume at
pay-per-token API rates would have been roughly ~$200 (budget models)
to ~$5,200 (frontier models), almost all of it cache reads at their
discounted rate.

## 2026-08-28 … 08-29 — the LoongArch backend

The third architecture: the full scalar integer set (248
instructions), the assembly layer with pseudo-instructions, ELF/CLI/
server wiring, the property suite, and the rt-vm matrix (three VMs).
**2 days** — one of discussion and planning, one of building.

Orchestrated differently than the base: nine subagent runs in
parallel (seven instruction slices, one rate-limit retry, one
property suite), each obliged to verify every encoding against
`llvm-mc` before writing tests.

| Metric | Value |
|---|---|
| Hand-written Go code | 25,644 lines (15,840 production + 9,804 tests) in 558 files |
| Other hand-written code (asm examples / shell / Makefile / docs) | ~5,990 lines |
| Generated tables | 261 lines |
| AI sessions | 21 (4 top-level + 17 subagent) |
| Model requests | 1,215 |

### Token usage (actual, measured)

| Category | Tokens |
|---|---|
| Input (total) | ~321M |
| — of which cache reads | ~318M (99.2%) |
| Output | ~1.02M |
| Tool calls | 1,745 |

Fun ratio: ~12K tokens per line this time — less than half the base
rate (the instruction slices were highly parallel and template
clones); the pure output cost is ~40 tokens per line.

### Cost

Same flat subscription. The reference pay-per-token equivalent:
~$30 (budget) to ~$740 (frontier), again almost all cache reads.

## Running total

| Metric | Base | + LoongArch | Total |
|---|---|---|---|
| Days | 16 | 2 | 18 |
| Hand-written Go | 79,562 | 25,644 | 105,206 |
| Go files | 440 | 558 | 998 |
| AI sessions | 117 | 21 | 138 |
| Model requests | 9,830 | 1,215 | 11,045 |
| Input tokens | ~2.32B | ~321M | ~2.64B |
| Output tokens | ~7.06M | ~1.02M | ~8.08M |
