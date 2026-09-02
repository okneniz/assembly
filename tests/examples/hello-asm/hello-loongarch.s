# hello-loongarch.s - bare-metal LoongArch64 hello for the QEMU virt
# machine: one string through the 16550A UART, then an idle loop.
#
# Build: assembly -arch loong64 -base 0x1c000000 -format bin
# Run:   qemu-system-loongarch64 -machine virt -nographic \
#          -device loader,file=hello-loongarch.bin,addr=0x1c000000,cpu-num=0
#
# Unlike the aarch64/riscv64 virt machines, LoongArch resets the CPU into
# the flash at 0x1c000000 (not into RAM), so the image is placed there;
# and the upstream virt machine has no bare-metal poweroff register (the
# old vendor-tree 0x100100 PM is gone, shutdown went through ACPI/GED),
# so the program idles after printing - the VM gate observes the output
# and stops the machine.

	.text
start:
	# $t0 = the UART data register (16550A, byte-wide at 0x1fe001e0).
	lu12i.w	$t0, 0x1fe00
	ori	$t0, $t0, 0x1e0

	# $t1 = the message cursor.
	la	$t1, msg

	# for (; *p; p++) *uart = *p;
1:
	ld.bu	$t2, $t1, 0
	beq	$t2, $zero, 2f
	st.b	$t2, $t0, 0
	addi.w	$t1, $t1, 1
	b	1b

	# idle: the gate stops the machine once the line is out.
2:
	b	2b

	.data
msg:
	.ascii "hello world\n"
	.byte 0
