# Test targets of the project: all module go tests, rebuilding the test
# binaries (examples), the round-trip corpus, and the behavioral VM matrix.
#
# Run from the repository root:
#
#	make tests              # full run of everything with one command
#	make -C tests <target>  # a specific target
#
# Recipes execute from the repository root (cd ..), so the paths inside are
# root-level (tests/..., arch/..., cmd/...). A full run needs objdump and
# qemu-system-* on PATH; for examples - clang with the RISC-V target.

.PHONY: all test examples file-fixtures rt-cli rt-run rt-report rt-vm

# all - the full run of all gates: rebuild the test binaries → all unit and
# integration tests → round-trip over the corpus (it builds the corpus itself
# via rt-cli) → the behavioral VM matrix.
all: examples test rt-run rt-vm

test:
	cd .. && go test ./...

# examples rebuilds the test binaries for every supported architecture:
#   - arm64:   a Go Mach-O (the canonical example the test suite diffs against)
#   - riscv64: an RV64I ELF object (no compressed) compiled from seed.c
# CLANG must be an LLVM build with the RISC-V target (Apple's /usr/bin/clang
# lacks it). Override: make -C tests examples CLANG=/opt/homebrew/opt/llvm/bin/clang
CLANG ?= clang
examples:
	cd .. && GOOS=darwin GOARCH=arm64 go build -ldflags='-s -w' -o tests/examples/hello-world/hello-world ./tests/examples/hello-world
	cd .. && "$(CLANG)" --target=riscv64-linux-gnu -march=rv64gc -mabi=lp64d -O2 -c \
		-o tests/examples/hello-riscv/hello-riscv.o tests/examples/hello-riscv/seed.c
	@echo "built examples: tests/examples/hello-world/hello-world (arm64 Mach-O), tests/examples/hello-riscv/hello-riscv.o (riscv64 ELF)"

# file-fixtures rebuilds the parser test corpus (committed to git; regenerating
# is optional). ELF: Go cross-builds covering 32/64-bit LE/BE + clang -c .o
# with symtab/relocs + a PIE with PT_DYNAMIC. Mach-O: clang exec/dylib/.o/i386
# + a fat binary. Tools are optional — missing pieces are skipped gracefully.
file-fixtures:
	@cd .. && printf 'package main\nfunc main() { println("hi") }\n' > /tmp/hello.go
	@cd .. && for t in linux/arm64 linux/riscv64 linux/386 linux/ppc64 linux/s390x; do \
		CGO_ENABLED=0 GOOS=$${t%/*} GOARCH=$${t#*/} go build -ldflags="-s -w" \
			-o file/elf/testdata/hello-$${t#*/} /tmp/hello.go; done
	@cd .. && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -buildmode=pie -ldflags="-s -w" \
		-o file/elf/testdata/dyn-pie-arm64 /tmp/hello.go
	@cd .. && printf 'int glob = 42;\nextern int ext(int);\nint call(int x) { return ext(x) + glob; }\nint *ptr = &glob;\n' > /tmp/reloc.c
	@cd .. && for tgt in aarch64-linux-gnu riscv64-linux-gnu i386-linux-gnu powerpc-linux-gnu; do \
		"$(CLANG)" -c -target $$tgt -x c -o file/elf/testdata/reloc-$$(echo $$tgt | cut -d- -f1).o /tmp/reloc.c; done
	@cd .. && printf 'int main(void) { return 42; }\n' > /tmp/mh.c
	@cd .. && "$(CLANG)" -arch arm64 -o file/macho/testdata/exec-arm64 /tmp/mh.c
	@cd .. && "$(CLANG)" -arch arm64 -c -o file/macho/testdata/obj-arm64.o /tmp/mh.c
	@cd .. && "$(CLANG)" -arch arm64 -dynamiclib -o file/macho/testdata/dylib-arm64.dylib /tmp/mh.c
	@cd .. && "$(CLANG)" -arch arm64 -arch x86_64 -o file/macho/testdata/fat-exec /tmp/mh.c
	@cd .. && "$(CLANG)" -arch i386 -c -o file/macho/testdata/obj-i386.o /tmp/mh.c 2>/dev/null || \
		echo "skip: i386 object (no 32-bit toolchain)"
	@echo "rebuilt file/*/testdata fixtures"

# --- round-trip final test: prebuilt binaries → disasm → asm → bytes ---

# Paths, not names: assembly-diff is a test utility living in tests/cmd.
RT_CMDS := cmd/assembly-asm tests/cmd/assembly-diff cmd/assembly-server
RT_ARCHES := arm64 riscv64

# rt-cli cross-builds the final-test corpus: our own CLIs for linux/arm64 and
# linux/riscv64 (CGO off - static, no libc). Artifacts go to tests/corpus/src
# (gitignored), checksums of the originals to tests/corpus/SHA256SUMS.
rt-cli:
	@cd .. && mkdir -p tests/corpus/src
	@cd .. && for c in $(RT_CMDS); do \
		for a in $(RT_ARCHES); do \
			CGO_ENABLED=0 GOOS=linux GOARCH=$$a go build -ldflags='-s -w' \
				-o tests/corpus/src/$${c##*/}-linux-$$a ./$$c || exit 1; \
		done; \
	done
	@cd .. && cd tests/corpus/src && shasum -a 256 * > ../SHA256SUMS
	@cd .. && echo "rt corpus: $$(ls tests/corpus/src | wc -l | tr -d ' ') binaries in tests/corpus/src, sums in tests/corpus/SHA256SUMS"

# rt-run - round-trip fidelity over the whole corpus: 3 iterations of
# "disasm → listing → assembly" per binary, a sha256 gate on .text, and a
# rebuild of the whole binary with our ELF writer into tests/corpus/out.
# The report (jsonl) is tests/corpus/rt-report.jsonl.
rt-run: rt-cli
	@cd .. && mkdir -p tests/corpus/out
	@cd .. && rm -f tests/corpus/rt-report.jsonl
	@cd .. && rc=0; for bin in tests/corpus/src/*; do \
		echo "== $$bin"; \
		go run ./tests/cmd/assembly-rt -bin $$bin -iterations 3 \
			-rebuild tests/corpus/out/$$(basename $$bin).rt.elf \
			-json >> tests/corpus/rt-report.jsonl || rc=1; \
	done; exit $$rc

# rt-report prints a summary of the last rt-run.
rt-report:
	@cat corpus/rt-report.jsonl 2>/dev/null || echo "no corpus/rt-report.jsonl - run make -C tests rt-run"

# rt-vm - the behavioral matrix original vs rebuilt in isolated VMs
# (qemu-system-aarch64/-riscv64, Alpine initramfs). Script: tests/vm/rt-vm.sh.
rt-vm: rt-run
	cd .. && sh tests/vm/rt-vm.sh
