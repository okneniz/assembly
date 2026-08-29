#!/bin/sh
# rt-vm - the behavioral matrix of the final test in isolated VMs
# (qemu-system-aarch64 / qemu-system-riscv64, full system machines).
# For each architecture it builds an initramfs: Alpine minirootfs + the
# original CLIs (tests/corpus/src) + the ones rebuilt by our ELF writer
# (tests/corpus/out) + inputs; /init (tests/vm/init.sh) runs the matrix and
# powers the machine off. The log is the serial console.
#
# The download cache (kernels/minirootfs, with verified sha256) is
# tests/corpus/vm/cache. HTTP artifacts are downloaded with -k (the sandbox's
# intercepted TLS), so every file MUST be verified by sha256 against the
# release metadata (see the command history; the distribution's sidecar
# files + release yamls).
set -eu

cd "$(dirname "$0")/../.."

CACHE=tests/corpus/vm/cache
ALPINE_AARCH64_NETBOOT=alpine-netboot-3.21.7-aarch64.tar.gz
ALPINE_AARCH64_MINI=alpine-minirootfs-3.21.7-aarch64.tar.gz
ALPINE_RISCV64_MINI=alpine-minirootfs-3.21.7-riscv64.tar.gz
DEBIAN_RISCV64_NETBOOT=netboot.tar.gz
# loongarch64 lives only in edge, as dated snapshots; the kernel comes
# from the linux-lts apk (no netboot tarball for the port yet).
# sha256 of the apk: db3de7c1d95f1effe448dbda838c176c30b45834613fd8bda00d94e5a2738366
# (the apk carries an apk-tools signature; the minirootfs has a sidecar
# .sha256 in the release directory)
ALPINE_LOONG64_MINI=alpine-minirootfs-20251016-loongarch64.tar.gz
LOONG64_KERNEL_APK=linux-lts-6.18.48-r0.apk

# blob.bin - input for --disasm/diff: a small binary built by OUR CLI on the
# host (a darwin/arm64 build of the same source).
make_blob() {
	go build -o /tmp/rt-assembly-asm ./cmd/assembly-asm
	if [ "$1" = arm64 ]; then
		/tmp/rt-assembly-asm -arch arm64 -base 0 \
			-o "$CACHE/blob-arm64.bin" tests/examples/hello-asm/hello-linux.s
	elif [ "$1" = loong64 ]; then
		/tmp/rt-assembly-asm -arch loong64 -base 0 \
			-o "$CACHE/blob-loong64.bin" tests/examples/hello-asm/hello-loongarch.s
	else
		/tmp/rt-assembly-asm -arch riscv64 -base 0x80000000 \
			-o "$CACHE/blob-riscv64.bin" tests/examples/hello-asm/hello-riscv.s
	fi
}

# build_initrd ARCH QEMU_KERNEL CONSOLE MINIROOTFS EXAMPLE_SRC BLOB
build_initrd() {
	arch=$1
	kernel=$2
	console=$3
	mini=$4
	example=$5
	blob=$6

	root="tests/corpus/vm/root-$arch"
	rm -rf "$root"
	mkdir -p "$root/matrix/orig" "$root/matrix/rt" "$root/matrix/in"
	tar xzf "$CACHE/$mini" -C "$root"
	for cmd in assembly-asm assembly-diff assembly-server; do
		cp "tests/corpus/src/$cmd-linux-$arch" "$root/matrix/orig/$cmd"
		cp "tests/corpus/out/$cmd-linux-$arch.rt.elf" "$root/matrix/rt/$cmd"
		chmod 0755 "$root/matrix/orig/$cmd" "$root/matrix/rt/$cmd"
	done
	cp "$example" "$root/matrix/in/hello.s"
	cp "$blob" "$root/matrix/in/blob.bin"
	sed "s/__ARCH__/$arch/" tests/vm/init.sh > "$root/init"
	chmod 0755 "$root/init"
	(
		cd "$root"
		find . -print0 | cpio -0 -o -H newc 2>/dev/null | gzip -1
	) > "tests/corpus/vm/initrd-$arch.cpio.gz"
	echo "$kernel"
}

run_vm() {
	arch=$1
	qemu=$2
	cpu=$3
	kernel=$4
	console=$5

	echo "== VM $arch"
	cpuspec=""
	[ -n "$cpu" ] && cpuspec="-cpu $cpu"
	append="console=$console"
	[ "$arch" = loong64 ] && append="$append initcall_blacklist=i8042_init"
	qemu-system-$qemu -M virt -m 1024 -nographic $cpuspec \
		-kernel "$kernel" \
		-initrd "tests/corpus/vm/initrd-$arch.cpio.gz" \
		-append "$append" \
		< /dev/null > "tests/corpus/vm/matrix-$arch.log" 2>&1 || true
	grep -E '^MATRIX ' "tests/corpus/vm/matrix-$arch.log" || {
		echo "no MATRIX lines in the $arch log:"; tail -20 "tests/corpus/vm/matrix-$arch.log"; exit 1; }
}

command -v qemu-system-aarch64 > /dev/null || { echo "no qemu-system-aarch64"; exit 1; }
command -v qemu-system-riscv64 > /dev/null || { echo "no qemu-system-riscv64"; exit 1; }
command -v qemu-system-loongarch64 > /dev/null || { echo "no qemu-system-loongarch64"; exit 1; }

for f in "$ALPINE_AARCH64_NETBOOT" "$ALPINE_AARCH64_MINI" "$ALPINE_RISCV64_MINI" "$DEBIAN_RISCV64_NETBOOT" \
	"$ALPINE_LOONG64_MINI" "$LOONG64_KERNEL_APK"; do
	[ -f "$CACHE/$f" ] || { echo "missing $CACHE/$f - see the download history and sha256"; exit 1; }
done

mkdir -p tests/corpus/vm
[ -f "$CACHE/blob-arm64.bin" ] || make_blob arm64
[ -f "$CACHE/blob-riscv64.bin" ] || make_blob riscv64
[ -f "$CACHE/blob-loong64.bin" ] || make_blob loong64

[ -d "$CACHE/kernel-aarch64" ] || {
	mkdir -p "$CACHE/kernel-aarch64"
	tar xzf "$CACHE/$ALPINE_AARCH64_NETBOOT" -C "$CACHE/kernel-aarch64" boot/vmlinuz-virt
}
[ -d "$CACHE/kernel-riscv64" ] || {
	mkdir -p "$CACHE/kernel-riscv64"
	tar xzf "$CACHE/$DEBIAN_RISCV64_NETBOOT" -C "$CACHE/kernel-riscv64" \
		debian-installer/riscv64/linux
}
[ -f "$CACHE/kernel-loongarch64/vmlinuz-lts" ] || {
	mkdir -p "$CACHE/kernel-loongarch64"
	tar xzf "$CACHE/$LOONG64_KERNEL_APK" -C "$CACHE/kernel-loongarch64" boot/vmlinuz-lts
}

build_initrd arm64 "$CACHE/kernel-aarch64/boot/vmlinuz-virt" ttyAMA0 \
	"$ALPINE_AARCH64_MINI" tests/examples/hello-asm/hello-linux.s "$CACHE/blob-arm64.bin"
run_vm arm64 aarch64 cortex-a53 "$CACHE/kernel-aarch64/boot/vmlinuz-virt" ttyAMA0

build_initrd riscv64 "$CACHE/kernel-riscv64/debian-installer/riscv64/linux" ttyS0 \
	"$ALPINE_RISCV64_MINI" tests/examples/hello-asm/hello-riscv.s "$CACHE/blob-riscv64.bin"
run_vm riscv64 riscv64 "" "$CACHE/kernel-riscv64/debian-installer/riscv64/linux" ttyS0

# The loong64 kernel panics in i8042_flush on the qemu virt machine before
# init starts (the port probes a nonexistent PS/2 controller); blacklisting
# the i8042 initcall skips it. init.sh adds the loong64 asm source branch.
build_initrd loong64 "$CACHE/kernel-loongarch64/vmlinuz-lts" ttyS0 \
	"$ALPINE_LOONG64_MINI" tests/examples/hello-asm/hello-loongarch.s "$CACHE/blob-loong64.bin"
run_vm loong64 loongarch64 "" "$CACHE/kernel-loongarch64/vmlinuz-lts" ttyS0

echo "== VM matrix done (logs: tests/corpus/vm/matrix-*.log)"
