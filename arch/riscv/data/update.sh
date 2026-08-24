#!/bin/sh
# Refresh the vendored Spike encoding.h used to generate the RISC-V CSR name
# table (arch/riscv/csr_generated.go).
#
# After running this, regenerate the Go table:
#     make gen-riscv-csr
#
# Behind a TLS-intercepting proxy, pass curl options, e.g.:
#     CURL_OPTS=-k ./update.sh
set -eu

DIR="$(cd "$(dirname "$0")" && pwd)"
ENCODING_URL="${ENCODING_URL:-https://raw.githubusercontent.com/riscv-software-src/riscv-isa-sim/master/riscv/encoding.h}"
CURL_OPTS="${CURL_OPTS:-}"

echo "==> downloading Spike encoding.h"
echo "    $ENCODING_URL"
curl $CURL_OPTS -fL "$ENCODING_URL" -o "$DIR/encoding.h"
count=$(grep -cE '#define CSR_' "$DIR/encoding.h" | tr -d ' ')
echo "    $count CSR defines in $DIR/encoding.h"

echo
echo "Done. Now regenerate the Go table:"
echo "    make gen-riscv-csr"
