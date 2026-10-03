.PHONY: fmt autofix fmt-check vet lint tests build cli vscode tidy clean gen-sysregs generate update-sysreg-data gen-riscv-csr update-riscv-csr-data gen-riscv-instr gen-arm-instr update-arm-instr-data gen-loongarch-instr update-loong-data prop-a64 prop-a64-1 prop-a64-2 prop-a64-rest

# GOLANGCI_LINT_VERSION pins the project-local linter (see bin/golangci-lint).
GOLANGCI_LINT_VERSION ?= v2.13.2

# fmt formats all Go sources via the project-local linter — gofmt (-s),
# goimports, gci (std/external/own-module groups), golines: the formatters
# configured in .golangci.yml.
#
# House layout convention (gofmt/golines do NOT check it and do NOT fix it —
# it is upheld by how the code is written/generated):
#   - composite literals and type declarations are multiline, one field
#     per line (see arch/arm64/schemas.go); empty `T{}`/`struct{}` may
#     stay compact;
#   - no one-line functions/methods/closures: the body always goes on
#     separate lines, `{ stmt }` on a single line is never written (an
#     empty body `{}` is fine);
#   - structs are created ONLY via constructors; the composite literal
#     lives exclusively in the body of its own constructor: newT/NewT
#     (by fields), decodeT (decode), make* dispatchers (by mnemonic),
#     api wrappers (arm64 …Of, riscv Op*/, immNum/ImmOf, arb.Enum).
#     Producing methods (Generate etc.) call the constructor rather than
#     write a literal. T{} — the zero value, allowed everywhere; generic
#     instantiations (immArb[T]{…}) stay as is.
#   - after a block statement's closing `}`, before the next statement —
#     a blank line (the only layout rule checked by the linter:
#     wsl_v5/after-block in .golangci.yml, autofix —
#     `golangci-lint run --fix`).
fmt: bin/golangci-lint
	./bin/golangci-lint fmt

# fmt-check fails if anything isn't formatted (for CI).
fmt-check:
	@out=$$( gofmt -s -l . ); if [ -n "$$out" ]; then echo "$$out"; exit 1; fi

vet:
	go vet ./...

# tests — the full run of all test gates in one command: the toolchain
# gates in the docker container (colima: examples → go tests → round-trip
# corpus) → the VM matrix on the host (qemu).
# Targets — tests/Makefile, layout and external dependencies — tests/README.md;
# individually: make -C tests <target>.
tests:
	$(MAKE) -C tests all

# prop-a64 — the arm64 property suite of the root package. The
# single-instruction family table no longer fits one go-test timeout,
# so the table is split in two halves by family name; the halves and the
# remaining properties (the list round trips, the alias families, the
# decode robustness) run as separate processes. Progress lines land in
# the log via oh-snap (ProgressEvery) — visible with -v (the default
# here); the exit code is the go test one.
prop-a64: prop-a64-1 prop-a64-2 prop-a64-rest

# prop-a64-1 — the single-instruction families A..L.
prop-a64-1:
	go test -run 'TestPropertySingleInstrRoundTrip/(A|B|C|D|E|F|G|H|I|J|K|L)[a-zA-Z0-9]*$$' -count=1 -timeout 30m -v .

# prop-a64-2 — the single-instruction families M..Z.
prop-a64-2:
	go test -run 'TestPropertySingleInstrRoundTrip/(M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z)[a-zA-Z0-9]*$$' -count=1 -timeout 30m -v .

# prop-a64-rest — the list round trips, the alias families and the decode
# robustness (the sources of the lists are the same families).
prop-a64-rest:
	go test -run 'TestPropertyBytesRoundTripList|TestPropertyTextRoundTripList|TestPropertyDecodeRobustness' -count=1 -timeout 30m -v .
	go test -run 'TestPropertyAliasRoundTrip' -count=1 -timeout 30m -v .


build:
	go build ./...

# cli builds the current sources of every command-line utility (assembly,
# assembly-debug, assembly-debug-dap) into ./bin — the project-local tools
# directory (already gitignored; clean removes it together with the pinned
# linter). CGO_ENABLED=0 keeps them static, per the zero-dependency thesis.
cli:
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/ ./cmd/...

# vscode installs the assembly-debug extension into ~/.vscode/extensions
# (overwriting a previous copy): the freshly built adapter is bundled INTO
# the extension and its absolute path is written into the manifest. A
# dock-launched VSCode hands its children the minimal launchd PATH — neither
# ~/go/bin nor /opt/homebrew/bin is there, so the debug type must not depend
# on PATH at all (the qemu lookup inside the adapter carries its own
# homebrew fallback). Restart VSCode afterwards: extensions are scanned at
# startup.
VSCODE_EXT_DIR ?= $(HOME)/.vscode/extensions/okneniz.assembly-debug-0.1.0
vscode: cli
	mkdir -p $(VSCODE_EXT_DIR)
	cp -R vscode/assembly-debug/. $(VSCODE_EXT_DIR)/
	cp bin/assembly-debug-dap $(VSCODE_EXT_DIR)/
	python3 -c 'import json,os,sys; p=sys.argv[1]; m=json.load(open(p)); [d.update(program=os.path.join(os.path.dirname(p), d["program"])) for d in m["contributes"]["debuggers"]]; json.dump(m, open(p, "w"), indent=2); open(p, "a").write("\n")' $(VSCODE_EXT_DIR)/package.json

# gen-sysregs regenerates arch/arm64/sysregs_generated.go from the vendored
# ARM System Register XML and m1n1's apple_regs.json. Re-run after updating the
# data under arch/arm64/data/. It then formats the whole tree and runs vet so
# the generated code lands clean.
gen-sysregs:
	go run ./gen/cmd/gen-sysregs -i arch/arm64/data/sysreg -apple arch/arm64/data/apple_regs.json -o arch/arm64/sysregs_generated.go
	gofmt -s -w .
	go vet ./...

# generate runs all code generators.
generate: gen-sysregs gen-riscv-csr gen-riscv-instr gen-arm-instr gen-loongarch-instr

# bin/golangci-lint installs the pinned golangci-lint into ./bin (GOBIN), so
# the project doesn't depend on whatever golangci-lint is on PATH. The -w -s ldflags keep the binary small and dodge a darwin/arm64 internal
# linker failure ("no room to add dwarf info") on this large binary.
# CGO_ENABLED=0 — the linter is pure Go, and prebuilt go1.26+ toolchains
# request a MacOSX26.sdk sysroot older CommandLineTools installs don't have.
bin/golangci-lint:
	GOBIN=$(CURDIR)/bin CGO_ENABLED=0 go install -ldflags='-w -s' \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# lint runs the non-modifying checks used in CI (format + vet + errcheck).
# The linters enabled are controlled by .golangci.yml.
lint: fmt-check vet bin/golangci-lint
	./bin/golangci-lint run

# autofix applies every automatic fix available: the formatters (gofmt -s,
# goimports, gci, golines) and the auto-fix linters (modernize, intrange,
# copyloopvar, misspell, godot, whitespace, wsl_v5/after-block, ...).
# It does NOT clear issues from fixer-less linters (funcorder, prealloc,
# errcheck, unused, unparam, ...) — those are manual work; `make lint`
# lists what remains after an autofix.
autofix: bin/golangci-lint
	./bin/golangci-lint run --fix

# update-sysreg-data re-downloads the vendored ARM XML and m1n1 apple_regs.json.
# Run gen-sysregs afterwards to regenerate the Go name table. Override pinned
# sources via ARM_SYSREG_URL / APPLE_REGS_URL; use CURL_OPTS=-k behind a proxy.
update-sysreg-data:
	arch/arm64/data/update.sh

# gen-riscv-csr regenerates arch/riscv/csr_generated.go from the vendored Spike
# encoding.h. Re-run after updating arch/riscv/data/encoding.h.
gen-riscv-csr:
	go run ./gen/cmd/gen-riscv-csr -i arch/riscv/data/encoding.h -o arch/riscv/csr_generated.go
	gofmt -s -w .
	go vet ./...

# update-riscv-csr-data re-downloads the vendored Spike encoding.h.
# Run gen-riscv-csr (and gen-riscv-instr) afterwards to regenerate the Go tables.
update-riscv-csr-data:
	arch/riscv/data/update.sh

# gen-riscv-instr regenerates arch/riscv/instr_generated.go (the {match,mask}
# encoding table) from the vendored Spike encoding.h. Shares the data source
# with gen-riscv-csr; re-run after updating arch/riscv/data/encoding.h.
gen-riscv-instr:
	go run ./gen/cmd/gen-riscv-instr -i arch/riscv/data/encoding.h -o arch/riscv/instr_generated.go
	gofmt -s -w .
	go vet ./...

# gen-arm-instr regenerates arch/arm64/isa_generated.go (the A64 instruction
# encoding table) from the vendored official A64 ISA XML. Re-run after updating
# arch/arm64/data/instr via update-arm-instr-data.
gen-arm-instr:
	go run ./gen/cmd/gen-arm-instr -i arch/arm64/data/instr -o arch/arm64/isa_generated.go
	gofmt -s -w .
	go vet ./...

# update-arm-instr-data re-downloads the vendored ARM A64 ISA XML tarball and
# extracts instruction XMLs into arch/arm64/data/instr. Run gen-arm-instr afterwards.
update-arm-instr-data:
	arch/arm64/data/update-instr.sh

# gen-loongarch-instr regenerates arch/loong64/instr_generated.go (the
# {match,mask} encoding table) from the vendored loongarch-opcodes tables.
# Re-run after updating arch/loong64/data via update-loong-data.
gen-loongarch-instr:
	go run ./gen/cmd/gen-loongarch-instr -i arch/loong64/data -o arch/loong64/instr_generated.go
	gofmt -s -w .
	go vet ./...

# update-loong-data re-downloads the vendored loongarch-opcodes tables (the
# scalar integer subsets of the LoongArch ISA). Run gen-loongarch-instr afterwards.
update-loong-data:
	arch/loong64/data/update.sh

tidy:
	go mod tidy

clean:
	rm -rf bin

