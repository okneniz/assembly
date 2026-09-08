#!/bin/sh
# init - the behavioral matrix of the final test inside an isolated VM
# (qemu-system, full system emulation). Compares the original binaries
# (/matrix/orig) with the ones rebuilt by our ELF writer (/matrix/rt):
# the commands' stdout/stderr/exit codes must
# match byte for byte. __ARCH__ is substituted at image build time
# (arm64|riscv64|loong64).

ARCH=__ARCH__

/bin/busybox mkdir -p /proc /sys /dev /tmp
/bin/busybox mount -t proc none /proc
/bin/busybox mount -t sysfs none /sys
/bin/busybox mount -t devtmpfs none /dev
/bin/busybox --install -s

FAIL=0

# dump_pair prints both logs to the console on a divergence (FAIL diagnostics).
dump_pair() { # log file name prefix
	echo "---- $1: ORIG ----"
	head -30 "/tmp/$1-orig.log" 2>&1
	echo "---- $1: RT ----"
	head -30 "/tmp/$1-rt.log" 2>&1
	echo "---- $1: diff ----"
	diff "/tmp/$1-orig.log" "/tmp/$1-rt.log" 2>&1 | head -20
}

# --- assembly: assembling the example, usage exit, disasm of the blob ---
for side in orig rt; do
	BIN=/matrix/$side/assembly
	{
		echo "=== asm:hex ==="
		"$BIN" -arch "$ARCH" -base 0 --hex /matrix/in/hello.s
		echo "exit=$?"
		echo "=== asm:usage ==="
		"$BIN"
		echo "exit=$?"
		echo "=== asm:disasm ==="
		"$BIN" -arch "$ARCH" --disasm /matrix/in/blob.bin
		echo "exit=$?"
	} > "/tmp/asm-$side.log" 2>&1
done
if cmp -s /tmp/asm-orig.log /tmp/asm-rt.log; then
	echo "MATRIX asm: PASS"
else
	echo "MATRIX asm: FAIL"
	dump_pair asm
	FAIL=1
fi

# --- assembly-diff: without arguments and on a fixture (there is no objdump
# in the VM - both sides must fail identically) ---
for side in orig rt; do
	BIN=/matrix/$side/assembly-diff
	{
		echo "=== diff:noargs ==="
		"$BIN"
		echo "exit=$?"
		echo "=== diff:fixture ==="
		"$BIN" /matrix/in/blob.bin
		echo "exit=$?"
	} > "/tmp/diff-$side.log" 2>&1
done
if cmp -s /tmp/diff-orig.log /tmp/diff-rt.log; then
	echo "MATRIX diff: PASS"
else
	echo "MATRIX diff: FAIL"
	dump_pair diff
	FAIL=1
fi

if [ "$FAIL" -eq 0 ]; then
	echo "MATRIX FINAL: PASS"
	poweroff -f
else
	echo "MATRIX FINAL: FAIL"
	# do not shut the machine down right away - the serial log has already
	# gone out; but nothing depends on this
	poweroff -f
fi
