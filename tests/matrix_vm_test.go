package tests

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/asm/riscv/pseudo"
)

// The behavioral matrix of the final test: the corpus CLIs (tests/corpus)
// - the cross-compiled originals and the binaries rebuilt by our ELF
// writer - run inside isolated VMs, and their observable behavior
// (stdout/stderr/exit codes/HTTP responses, see tests/vm/init.sh) must
// match byte for byte. One qemu-system machine per architecture; the
// subtests are t.Parallel - three sequential boots used to dominate the
// matrix wall time. This kind of parallelism shares no in-process state
// (the machines are separate processes), unlike the arb/parsec state that
// keeps the rest of the suite serial.
//
// Nothing is downloaded by the test: kernels and Alpine minirootfs
// snapshots live in the gitignored corpus/vm/cache, curated by hand with
// sha256 pins (see tests/README.md). Outside make -C tests rt-vm missing
// artifacts or emulators only skip the test; ASSEMBLY_RTVM_STRICT=1
// (set by the make target) turns them into failures - a gate must not
// stay silent in the full run.

const (
	// vmMatrixDeadline bounds one whole machine: Alpine boot + matrix +
	// poweroff fits into well under a minute even under TCG; the ceiling
	// only catches hangs (a machine with no poweroff would idle forever).
	vmMatrixDeadline = 5 * time.Minute

	// vmMemory is the guest RAM of every matrix machine, in MB.
	vmMemory = "1024"
)

// matrixPassLines - what /init prints when the whole matrix passes; every
// line must be present (the old shell gate checked only that MATRIX lines
// existed at all, so a FAIL run passed the gate too).
var matrixPassLines = []string{
	"MATRIX asm: PASS",
	"MATRIX diff: PASS",
	"MATRIX server: PASS",
	"MATRIX FINAL: PASS",
}

// matrixCommands - the corpus CLIs assembled into the initramfs, both the
// originals and the ELF-writer rebuilds.
var matrixCommands = []string{
	"assembly",
	"assembly-diff",
	"assembly-server",
}

// kernelSrc - a vmlinuz extracted once from a cached tarball into the
// cache (the loongarch64 kernel lives inside an apk, which is also a
// gzip tar).
type kernelSrc struct {
	tarball string
	inner   string
	out     string
}

// matrixVM - the launch description of one matrix machine.
type matrixVM struct {
	arch     string // corpus name: arm64|riscv64|loong64
	qemu     string
	console  string
	kernel   kernelSrc
	mini     string // minirootfs tarball in the cache
	src      string // asm source: copied in as matrix/in/hello.s and assembled into blob.bin
	blobBase uint64
	accel    []string
}

// TestMatrixVM - the behavioral matrix in isolated VMs, one parallel
// subtest per machine: it builds the initramfs and boots it.
func TestMatrixVM(t *testing.T) {
	vms, ok := matrixVMs(t)
	if !ok {
		return
	}

	for i := range vms {
		vm := &vms[i]
		t.Run(vm.arch, func(t *testing.T) {
			t.Parallel()

			initrd := buildInitrd(t, vm)
			runMatrixVM(t, vm, initrd)
		})
	}
}

// matrixVMs - the per-arch machine descriptions after checking that
// everything the matrix needs is in place (emulators, cache, corpus).
func matrixVMs(t *testing.T) ([]matrixVM, bool) {
	t.Helper()

	strict := os.Getenv("ASSEMBLY_RTVM_STRICT") == "1"
	miss := func(what string, err error) bool {
		if err == nil {
			return false
		}

		if strict {
			t.Fatalf("VM matrix: %s: %v (gates must not stay silent)", what, err)
		}

		t.Skipf("%s: %v (make -C tests rt-run; see tests/README.md)", what, err)
		return true
	}

	vms := []matrixVM{
		{
			arch:    "arm64",
			qemu:    "qemu-system-aarch64",
			console: "ttyAMA0",
			kernel: kernelSrc{
				tarball: "corpus/vm/cache/alpine-netboot-3.21.7-aarch64.tar.gz",
				inner:   "boot/vmlinuz-virt",
				out:     "corpus/vm/cache/kernel-aarch64/boot/vmlinuz-virt",
			},
			mini:     "corpus/vm/cache/alpine-minirootfs-3.21.7-aarch64.tar.gz",
			src:      "examples/hello-asm/hello-linux.s",
			blobBase: 0,
			accel:    hostAccel(),
		},
		{
			arch:    "riscv64",
			qemu:    "qemu-system-riscv64",
			console: "ttyS0",
			kernel: kernelSrc{
				tarball: "corpus/vm/cache/netboot.tar.gz",
				inner:   "debian-installer/riscv64/linux",
				out:     "corpus/vm/cache/kernel-riscv64/debian-installer/riscv64/linux",
			},
			mini:     "corpus/vm/cache/alpine-minirootfs-3.21.7-riscv64.tar.gz",
			src:      "examples/hello-asm/hello-riscv.s",
			blobBase: 0x80000000,
			accel:    []string{},
		},
		{
			// The loong64 kernel panics in i8042_flush on the qemu virt
			// machine before init starts (the port probes a nonexistent
			// PS/2 controller); blacklisting the i8042 initcall skips it.
			arch:    "loong64",
			qemu:    "qemu-system-loongarch64",
			console: "ttyS0",
			kernel: kernelSrc{
				tarball: "corpus/vm/cache/linux-lts-6.18.48-r0.apk",
				inner:   "boot/vmlinuz-lts",
				out:     "corpus/vm/cache/kernel-loongarch64/vmlinuz-lts",
			},
			mini:     "corpus/vm/cache/alpine-minirootfs-20251016-loongarch64.tar.gz",
			src:      "examples/hello-asm/hello-loongarch.s",
			blobBase: 0x1c000000,
			accel:    []string{},
		},
	}

	for i := range vms {
		vm := &vms[i]
		if miss("emulator", execLookPath(vm.qemu)) {
			return nil, false
		}

		if miss("cache "+vm.kernel.tarball, fileExists(vm.kernel.tarball)) {
			return nil, false
		}

		if miss("cache "+vm.mini, fileExists(vm.mini)) {
			return nil, false
		}

		if miss("example "+vm.src, fileExists(vm.src)) {
			return nil, false
		}

		for _, cmd := range matrixCommands {
			orig := fmt.Sprintf("corpus/src/%s-linux-%s", cmd, vm.arch)
			rt := fmt.Sprintf("corpus/out/%s-linux-%s.rt.elf", cmd, vm.arch)

			if miss("corpus "+orig, fileExists(orig)) {
				return nil, false
			}

			if miss("corpus "+rt, fileExists(rt)) {
				return nil, false
			}
		}
	}

	if miss("init template", fileExists("vm/init.sh")) {
		return nil, false
	}

	return vms, true
}

// hostAccel - host acceleration for the aarch64 guest: hvf on macOS, kvm
// on Linux; both require the host CPU (a named core like cortex-a53 is
// TCG-only). Everything else stays on TCG.
func hostAccel() []string {
	if runtime.GOOS == "darwin" {
		return []string{"-accel", "hvf", "-cpu", "host"}
	}

	if st, err := os.Stat("/dev/kvm"); err == nil && st.Mode()&os.ModeCharDevice != 0 {
		return []string{"-accel", "kvm", "-cpu", "host"}
	}

	return []string{"-cpu", "cortex-a53"}
}

// buildInitrd - the initramfs of one machine: Alpine minirootfs + the
// corpus CLIs (orig and rt) + inputs + /init (tests/vm/init.sh with the
// arch substituted), packed as newc cpio, gzip -1.
func buildInitrd(t *testing.T, vm *matrixVM) string {
	t.Helper()

	cpio := &bytes.Buffer{}

	packRootfs(t, cpio, vm.mini)
	packMatrix(t, cpio, vm)
	cpioTrailer(cpio)

	out := fmt.Sprintf("corpus/vm/initrd-%s.cpio.gz", vm.arch)
	f, err := os.Create(out)
	require.NoError(t, err)

	gz, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
	require.NoError(t, err)

	_, err = cpio.WriteTo(gz)
	require.NoError(t, err)
	require.NoError(t, gz.Close())
	require.NoError(t, f.Close())

	return out
}

// packRootfs appends the minirootfs contents. Device nodes are skipped:
// /init mounts devtmpfs over /dev before anything looks at it. Hard links
// (rare in the minirootfs) are materialized as copies - the newc writer
// below supports only files, dirs and symlinks.
func packRootfs(t *testing.T, w *bytes.Buffer, tarball string) {
	t.Helper()

	f, err := os.Open(tarball)
	require.NoError(t, err)

	gz, err := gzip.NewReader(f)
	require.NoError(t, err, tarball)

	tr := tar.NewReader(gz)
	files := map[string][]byte{}
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err, tarball)

		name := tarName(hdr.Name)
		if name == "" {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			cpioEntry(w, cpioModeDir, name, nil)
		case tar.TypeReg:
			data, err := io.ReadAll(tr)
			require.NoError(t, err, name)

			files[name] = data
			cpioEntry(w, cpioRegMode(hdr.Mode), name, data)
		case tar.TypeSymlink:
			// a symlink's data is its target
			cpioEntry(w, cpioModeSymlink, name, []byte(hdr.Linkname))
		case tar.TypeLink:
			data, ok := files[hdr.Linkname]
			require.True(t, ok, "%s: hard link target %q not seen", name, hdr.Linkname)

			cpioEntry(w, cpioRegMode(hdr.Mode), name, data)
		}
	}

	require.NoError(t, gz.Close())
	require.NoError(t, f.Close())
}

// packMatrix appends the matrix payload: /init, the CLIs of both sides,
// and the inputs.
func packMatrix(t *testing.T, w *bytes.Buffer, vm *matrixVM) {
	t.Helper()

	src, err := os.ReadFile(vm.src)
	require.NoError(t, err)

	tmpl, err := os.ReadFile("vm/init.sh")
	require.NoError(t, err)

	for _, dir := range []string{"matrix", "matrix/orig", "matrix/rt", "matrix/in"} {
		cpioEntry(w, cpioModeDir, dir, nil)
	}

	init := strings.ReplaceAll(string(tmpl), "__ARCH__", vm.arch)
	cpioEntry(w, cpioModeExec, "init", []byte(init))
	cpioEntry(w, cpioModeFile, "matrix/in/hello.s", src)
	cpioEntry(w, cpioModeFile, "matrix/in/blob.bin", matrixBlob(t, vm, string(src)))

	for _, cmd := range matrixCommands {
		orig, err := os.ReadFile(fmt.Sprintf("corpus/src/%s-linux-%s", cmd, vm.arch))
		require.NoError(t, err)

		rt, err := os.ReadFile(fmt.Sprintf("corpus/out/%s-linux-%s.rt.elf", cmd, vm.arch))
		require.NoError(t, err)

		cpioEntry(w, cpioModeExec, "matrix/orig/"+cmd, orig)
		cpioEntry(w, cpioModeExec, "matrix/rt/"+cmd, rt)
	}
}

// matrixBlob - the --disasm input: the example assembled in-process by
// the facades, the sections concatenated in order (what the assembly
// CLI writes in the raw format on the host; NOBITS sections carry no
// bytes).
func matrixBlob(t *testing.T, vm *matrixVM, src string) []byte {
	t.Helper()

	var (
		res  *asm.Result
		errs []asm.AsmError
	)

	switch vm.arch {
	case "arm64":
		res, errs = alias.Assemble(src, vm.blobBase)
	case "riscv64":
		res, errs = pseudo.Assemble(src, vm.blobBase)
	case "loong64":
		res, errs = lpseudo.Assemble(src, vm.blobBase)
	default:
		t.Fatalf("unknown arch %q", vm.arch)
	}

	require.Empty(t, errs, "assemble %s", vm.src)

	blob := []byte{}
	for _, sec := range res.Sections {
		blob = append(blob, sec.Data...)
	}

	return blob
}

// runMatrixVM boots one machine and checks the serial log: every matrix
// line must be a PASS. The qemu exit status is not a signal - poweroff
// codes vary between machines - the log is; a deadline kill (a hang) is
// tolerated here precisely because it leaves no MATRIX lines and fails
// the check below.
func runMatrixVM(t *testing.T, vm *matrixVM, initrd string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), vmMatrixDeadline)
	defer cancel()

	logPath := fmt.Sprintf("corpus/vm/matrix-%s.log", vm.arch)
	log, err := os.Create(logPath)
	require.NoError(t, err)

	appendLine := "console=" + vm.console
	if vm.arch == "loong64" {
		appendLine += " initcall_blacklist=i8042_init"
	}

	args := slices.Concat(
		[]string{"-M", "virt", "-m", vmMemory, "-nographic"},
		vm.accel,
		[]string{
			"-kernel",
			vm.kernel.out,
			"-initrd",
			initrd,
			"-append",
			appendLine,
		},
	)

	// Stdin nil - the process reads the null device, like the old
	// `< /dev/null` redirect of the shell matrix
	cmd := exec.CommandContext(ctx, vm.qemu, args...)
	cmd.Stdout = log
	cmd.Stderr = log

	runErr := cmd.Run()
	require.NoError(t, log.Close())

	serial, err := os.ReadFile(logPath)
	require.NoError(t, err)

	if ctx.Err() == nil {
		require.NoError(t, runErr, "%s\nserial tail:\n%s",
			vm.qemu, tailLines(string(serial), 20))
	}

	missing := []string{}
	for _, line := range matrixPassLines {
		if !strings.Contains(string(serial), line) {
			missing = append(missing, line)
		}
	}

	require.Empty(t, missing, "matrix %s\nserial tail:\n%s",
		vm.arch, tailLines(string(serial), 20))
}

// tailLines - the last n lines of s (failure diagnostics).
func tailLines(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return strings.Join(lines, "\n")
}

// tarName normalizes an archive member name: no leading "./" or "/".
func tarName(name string) string {
	name = strings.TrimPrefix(name, "./")
	return strings.TrimPrefix(name, "/")
}

func execLookPath(qemu string) error {
	_, err := exec.LookPath(qemu)
	return err
}

func fileExists(path string) error {
	_, err := os.Stat(path)
	return err
}

// --- newc cpio writer (the format the Linux initramfs unpacker eats) ---

const (
	// cpioMagic is the newc header magic; cpioHeaderLen is the whole
	// 110-byte header (magic + 13 hex fields).
	cpioMagic     = "070701"
	cpioHeaderLen = 110
	cpioAlign     = 4

	cpioModeFile    = 0o100644
	cpioModeExec    = 0o100755
	cpioModeDir     = 0o040755
	cpioModeSymlink = 0o120777

	// cpioTrailerName is the conventional last entry of the archive.
	cpioTrailerName = "TRAILER!!!"
)

// cpioEntry appends one newc record: a 110-byte header, the name
// (NUL-terminated), both padded to 4 bytes, then the data, padded too.
// uid/gid/mtime are zero - the kernel does not care; hard links are not
// supported (materialized as copies by packRootfs).
func cpioEntry(w *bytes.Buffer, mode int64, name string, data []byte) {
	header := &bytes.Buffer{}
	header.WriteString(cpioMagic)

	fields := []int64{
		0, // inode (the kernel allocates its own)
		mode,
		0,                    // uid
		0,                    // gid
		1,                    // nlink
		0,                    // mtime
		int64(len(data)),     // filesize
		0,                    // devmajor
		0,                    // devminor
		0,                    // rdevmajor
		0,                    // rdevminor
		int64(len(name) + 1), // namesize (with the NUL)
		0,                    // check (newc ignores it)
	}
	for _, f := range fields {
		fmt.Fprintf(header, "%08X", uint64(f))
	}

	header.WriteString(name)
	header.WriteByte(0)

	w.Write(header.Bytes())
	cpioPad(w, cpioHeaderLen+len(name)+1)
	w.Write(data)
	cpioPad(w, len(data))
}

// cpioPad aligns the stream to 4 bytes after n payload bytes.
func cpioPad(w *bytes.Buffer, n int) {
	for i := (cpioAlign - n%cpioAlign) % cpioAlign; i > 0; i-- {
		w.WriteByte(0)
	}
}

// cpioTrailer marks the end of the archive.
func cpioTrailer(w *bytes.Buffer) {
	cpioEntry(w, 0, cpioTrailerName, nil)
}

// cpioRegMode maps a tar file mode onto a cpio one.
func cpioRegMode(tarMode int64) int64 {
	if tarMode&0o100 != 0 {
		return cpioModeExec
	}

	return cpioModeFile
}
