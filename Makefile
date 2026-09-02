.PHONY: fmt fmt-check vet lint tests build tidy clean gen-sysregs generate update-sysreg-data gen-riscv-csr update-riscv-csr-data gen-riscv-instr gen-arm-instr update-arm-instr-data gen-loongarch-instr update-loong-data

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

build:
	go build ./...

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

