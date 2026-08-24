#!/bin/sh
# Refresh the vendored ARM A64 instruction XML used to generate the ARM decode
# table (arch/arm64/isa_generated.go).
#
# After running this, regenerate the Go table:
#     make gen-arm-instr
#
# The XML is the official A64 ISA release (Exploration-Tools A64 ISA) from
# developer.arm.com. Override the pinned source by setting ISA_A64_URL, e.g.
# to bump to a newer release:
#     ISA_A64_URL='https://developer.arm.com/-/cdn-downloads/permalink/Exploration-Tools-A64-ISA/ISA_A64/ISA_A64_xml_A_profile-YYYY-MM.tar.gz' ./update-instr.sh
#
# Behind a TLS-intercepting proxy, pass curl options, e.g.:
#     CURL_OPTS=-k ./update-instr.sh
set -eu

DIR="$(cd "$(dirname "$0")" && pwd)"
INSTR_DIR="$DIR/instr"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

ISA_A64_URL="${ISA_A64_URL:-https://developer.arm.com/-/cdn-downloads/permalink/Exploration-Tools-A64-ISA/ISA_A64/ISA_A64_xml_A_profile-2026-03.tar.gz}"
CURL_OPTS="${CURL_OPTS:-}"

echo "==> downloading ARM A64 ISA XML"
echo "    $ISA_A64_URL"
curl $CURL_OPTS -fL "$ISA_A64_URL" -o "$TMP/isa.tar.gz"

echo "==> extracting instruction XML (latest ISA_A64_xml_A_profile-* release)"
mkdir -p "$INSTR_DIR"
rm -f "$INSTR_DIR"/*.xml
tar xzf "$TMP/isa.tar.gz" -C "$TMP"
# The tarball holds one or more versioned ISA_A64_xml_A_profile-* dirs (+ PDFs);
# copy the *.xml from the highest-named (latest) one.
latest="$(find "$TMP" -maxdepth 1 -type d -name 'ISA_A64_xml_A_profile-*' | sort | tail -1)"
cp "$latest"/*.xml "$INSTR_DIR"/
count=$(ls "$INSTR_DIR" | wc -l | tr -d ' ')
echo "    $count instruction XML files in $INSTR_DIR"

echo
echo "Done. Now regenerate the Go table:"
echo "    make gen-arm-instr"
