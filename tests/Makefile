# Test targets of the project: all module go tests, rebuilding the test
# binaries (examples), the round-trip corpus, and the behavioral VM matrix.
#
# Run from the repository root:
#
#	make tests              # full run of everything with one command
#	make -C tests <target>  # a specific target
#
# Recipes execute from the repository root (cd ..), so the paths inside are
# root-level (tests/..., arch/..., cmd/...). The toolchain part of the full
# run (clang, llvm-mc/llvm-objdump/llvm-readelf) lives in the docker
# container run by `docker` (the engine is colima); on the host only
# qemu-system-* is needed, for rt-vm.

.PHONY: all test docker examples file-fixtures rt-cli rt-run rt-report rt-vm

# all - the full run of all gates: the toolchain gates in the docker
# container (examples → all unit and oracle tests → round-trip over the
# corpus) → the behavioral VM matrix on the host (qemu).
all: docker rt-vm

test:
	cd .. && go test ./...

# --- toolchain container: LLVM oracles without a local install ---

# DOCKER - the docker invocation for the toolchain container. colima is the
# engine: the context is pinned so a running Docker Desktop does not take
# the builds silently. Override: make -C tests docker DOCKER=docker
DOCKER ?= docker --context colima
IMAGE ?= assembly-tests

# DOCKER_RT_JOBS - rt-run parallelism inside the container. Each assembly-rt
# holds the full listing of a multi-megabyte .text in memory, and the colima
# VM (8 GiB by default, half of a 16 GiB host) is tighter than a bare-metal
# host: at 4 one of the processes gets oom-killed mid-run.
DOCKER_RT_JOBS ?= 2

# docker - the toolchain gates inside the Linux container: examples (clang)
# → go test ./... (the llvm-oracle tests; the QEMU tests skip inside - the
# image has no emulator) → rt-run (the corpus round-trip; its artifacts
# land in the mounted repo and feed rt-vm on the host). The repo is
# bind-mounted at /src; the Go caches live in named volumes so repeated
# runs do not re-download and re-build the world.
docker:
	@cd .. && $(DOCKER) build -t $(IMAGE) tests/docker
	@cd .. && $(DOCKER) run --rm \
		-v "$$(pwd)":/src -w /src \
		-v assembly-gocache:/root/.cache/go-build \
		-v assembly-gomod:/go/pkg/mod \
		$(IMAGE) make -C tests examples test rt-run RT_JOBS=$(DOCKER_RT_JOBS)

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
# + a fat binary + an arm64e slice of /bin/ls (lipo) - these need the Apple
# toolchain and skip gracefully elsewhere. Tools are optional — missing
# pieces are skipped gracefully.
file-fixtures:
	@cd .. && printf 'package main\nfunc main() { println("hi") }\n' > /tmp/hello.go
	@cd .. && for t in linux/arm64 linux/riscv64 linux/loong64 linux/386 linux/ppc64 linux/s390x; do \
		CGO_ENABLED=0 GOOS=$${t%/*} GOARCH=$${t#*/} go build -ldflags="-s -w" \
			-o file/elf/testdata/hello-$${t#*/} /tmp/hello.go; done
	@cd .. && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -buildmode=pie -ldflags="-s -w" \
		-o file/elf/testdata/dyn-pie-arm64 /tmp/hello.go
	@cd .. && printf 'int glob = 42;\nextern int ext(int);\nint call(int x) { return ext(x) + glob; }\nint *ptr = &glob;\n' > /tmp/reloc.c
	@cd .. && for tgt in aarch64-linux-gnu riscv64-linux-gnu loongarch64-linux-gnu i386-linux-gnu powerpc-linux-gnu; do \
		"$(CLANG)" -c -target $$tgt -x c -o file/elf/testdata/reloc-$$(echo $$tgt | cut -d- -f1).o /tmp/reloc.c; done
	@cd .. && printf 'int main(void) { return 42; }\n' > /tmp/mh.c
	@cd .. && "$(CLANG)" -arch arm64 -o file/macho/testdata/exec-arm64 /tmp/mh.c || \
		echo "skip: exec-arm64 (needs the Apple toolchain)"
	@cd .. && "$(CLANG)" -arch arm64 -c -o file/macho/testdata/obj-arm64.o /tmp/mh.c || \
		echo "skip: obj-arm64.o (needs the Apple toolchain)"
	@cd .. && "$(CLANG)" -arch arm64 -dynamiclib -o file/macho/testdata/dylib-arm64.dylib /tmp/mh.c || \
		echo "skip: dylib-arm64.dylib (needs the Apple toolchain)"
	@cd .. && "$(CLANG)" -arch arm64 -arch x86_64 -o file/macho/testdata/fat-exec /tmp/mh.c || \
		echo "skip: fat-exec (needs the Apple toolchain)"
	@cd .. && "$(CLANG)" -arch i386 -c -o file/macho/testdata/obj-i386.o /tmp/mh.c 2>/dev/null || \
		echo "skip: i386 object (no 32-bit toolchain)"
	@cd .. && lipo /bin/ls -thin arm64e -output file/macho/testdata/ls-arm64e 2>/dev/null || \
		echo "skip: ls-arm64e (macOS lipo only)"
	@echo "rebuilt file/*/testdata fixtures"

# --- round-trip final test: prebuilt binaries → disasm → asm → bytes ---

# Paths, not names: assembly-diff is a test utility living in tests/cmd.
RT_CMDS := cmd/assembly tests/cmd/assembly-diff cmd/assembly-server
RT_ARCHES := arm64 riscv64 loong64

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

# RT_JOBS - how many corpus binaries round-trip concurrently. Parallelism
# is per process only: inside one assembly-rt the arb seed and parsec are
# not thread-safe. Each process holds the full listing of a multi-megabyte
# .text in memory - lower this on RAM-tight machines.
RT_JOBS ?= 4

# rt-run - round-trip fidelity over the whole corpus: 3 iterations of
# "disasm → listing → assembly" per binary, a sha256 gate on .text, and a
# rebuild of the whole binary with our ELF writer into tests/corpus/out.
# The binaries are independent (each writes its own report file), so they
# run in batches of RT_JOBS; the merged report (jsonl) is
# tests/corpus/rt-report.jsonl.
rt-run: rt-cli
	@cd .. && mkdir -p tests/corpus/out tests/corpus/rt
	@cd .. && rm -f tests/corpus/rt-report.jsonl && rm -rf tests/corpus/rt/*
	@cd .. && go build -o tests/corpus/rt/assembly-rt ./tests/cmd/assembly-rt
	@cd .. && rc=0; set -- tests/corpus/src/*; \
		while [ $$# -gt 0 ]; do \
			pids=""; \
			i=0; \
			while [ $$i -lt $(RT_JOBS) ] && [ $$# -gt 0 ]; do \
				echo "== $$1"; \
				( tests/corpus/rt/assembly-rt -bin $$1 -iterations 3 \
					-rebuild tests/corpus/out/$$(basename $$1).rt.elf \
					-json > tests/corpus/rt/$$(basename $$1).json ) & \
				pids="$$pids $$!"; \
				shift; \
				i=$$((i + 1)); \
			done; \
			for p in $$pids; do wait $$p || rc=1; done; \
		done; \
		cat tests/corpus/rt/*.json > tests/corpus/rt-report.jsonl; exit $$rc

# rt-report prints a summary of the last rt-run.
rt-report:
	@cat corpus/rt-report.jsonl 2>/dev/null || echo "no corpus/rt-report.jsonl - run make -C tests rt-run"

# rt-vm - the QEMU gates on the host: the behavioral matrix original vs
# rebuilt in isolated VMs (TestMatrixVM; the /init template is
# tests/vm/init.sh - one machine per architecture, booted in parallel) and
# the bare-metal hello demo (TestHelloVM). The matrix consumes the rt-run
# artifacts (corpus binaries + rebuilds; produced by `docker` in the
# container, or by a local rt-run) and the downloaded corpus/vm/cache
# (kernels/minirootfs). ASSEMBLY_RTVM_STRICT turns missing tools/artifacts
# from skips into failures - the full gate must not stay silent; -count=1
# keeps the go test result cache from substituting a stale plain-mode run.
rt-vm:
	cd .. && ASSEMBLY_RTVM_STRICT=1 go test -count=1 -run 'TestMatrixVM|TestHelloVM' ./tests -v
