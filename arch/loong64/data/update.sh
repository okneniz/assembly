#!/bin/sh
# Refresh the vendored loongarch-opcodes tables (the scalar integer subsets
# of the LoongArch ISA) used to generate arch/loong64/instr_generated.go.
#
# After running this, regenerate the Go table:
#     make gen-loongarch-instr
#
# Behind a TLS-intercepting proxy, pass curl options, e.g.:
#     CURL_OPTS=-k ./update.sh
set -eu

DIR="$(cd "$(dirname "$0")" && pwd)"
OpcodesBaseURL="${OpcodesBaseURL:-https://raw.githubusercontent.com/loongson-community/loongarch-opcodes/develop}"
CURL_OPTS="${CURL_OPTS:-}"

# The scalar integer scope of arch/loong64: base, multiply, atomics, bitops,
# bound checks, privileged. FP/LSX/LASX/LBT/LVZ subsets are out of scope.
SUBSETS="
	la-base-32
	la-base-64
	la-mul-32
	la-mul-64
	la-atomics-32
	la-atomics-64
	la-bitops-32
	la-bitops-64
	la-bound
	la-bound-64
	la-privileged-32
	la-privileged-64
"

count=0
for subset in $SUBSETS; do
	echo "==> downloading $subset.txt"
	echo "    $OpcodesBaseURL/$subset.txt"
	curl $CURL_OPTS -fL "$OpcodesBaseURL/$subset.txt" -o "$DIR/$subset.txt"
	count=$((count + $(grep -c . "$DIR/$subset.txt" | tr -d ' ')))
done

echo
echo "    $count instruction lines in $DIR"
echo
echo "Done. Now regenerate the Go table:"
echo "    make gen-loongarch-instr"
