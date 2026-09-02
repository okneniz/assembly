#!/bin/sh
# init - the behavioral matrix of the final test inside an isolated VM
# (qemu-system, full system emulation). Compares the original binaries
# (/matrix/orig) with the ones rebuilt by our ELF writer (/matrix/rt):
# the commands' stdout/stderr/exit codes and the server's HTTP responses must
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

# --- assembly-server: static files and the JSON API (assembling a source) ---
if [ "$ARCH" = arm64 ]; then
	ASM_SRC='mov x0, #42\n'
elif [ "$ARCH" = loong64 ]; then
	ASM_SRC='li.w $a0, 42\n'
else
	ASM_SRC='li a0, 42\n'
fi
for side in orig rt; do
	BIN=/matrix/$side/assembly-server
	"$BIN" -addr 127.0.0.1:18765 > "/tmp/srvrun-$side.log" 2>&1 &
	SPID=$!
	sleep 2
	{
		echo "=== server:static ==="
		wget -q -O - http://127.0.0.1:18765/
		echo "exit=$?"
		echo "=== server:asm-api ==="
		wget -q -O - --post-data='{"arch":"'"$ARCH"'","source":"'"$ASM_SRC"'","baseAddr":0}' \
			http://127.0.0.1:18765/api/v1/asm
		echo "exit=$?"
	} > "/tmp/srv-$side.log" 2>&1
	kill "$SPID" 2>/dev/null
	wait "$SPID" 2>/dev/null
done
if cmp -s /tmp/srv-orig.log /tmp/srv-rt.log; then
	echo "MATRIX server: PASS"
else
	echo "MATRIX server: FAIL"
	dump_pair srv
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
