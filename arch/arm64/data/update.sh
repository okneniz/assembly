#!/bin/sh
# Refresh the vendored ARM system register XML and m1n1 apple_regs.json.
#
# After running this, regenerate the Go name table:
#     make gen-sysregs
#
# Override the pinned sources by setting the *_URL environment variables, e.g.
# to bump to a newer ARM release:
#     ARM_SYSREG_URL='https://.../SysReg_xml_v9x-YYYY-MM.tar.gz' ./update.sh
#
# Behind a TLS-intercepting proxy, pass curl options, e.g.:
#     CURL_OPTS=-k ./update.sh
set -eu

DIR="$(cd "$(dirname "$0")" && pwd)"
SYSREG_DIR="$DIR/sysreg"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

ARM_URL="${ARM_SYSREG_URL:-https://developer.arm.com/-/media/developer/products/architecture/armv8-a-architecture/2020-06/SysReg_xml_v86A-2020-06.tar.gz}"
APPLE_URL="${APPLE_REGS_URL:-https://raw.githubusercontent.com/AsahiLinux/m1n1/master/proxyclient/m1n1/apple_regs.json}"
CURL_OPTS="${CURL_OPTS:-}"

echo "==> downloading ARM System Register XML"
echo "    $ARM_URL"
curl $CURL_OPTS -fL "$ARM_URL" -o "$TMP/sysreg.tar.gz"

echo "==> extracting AArch64-*.xml"
mkdir -p "$SYSREG_DIR"
rm -f "$SYSREG_DIR"/AArch64-*.xml
tar xzf "$TMP/sysreg.tar.gz" -C "$TMP"
# tar layouts differ across releases (top-level version dir, xhtml/ subdir);
# copy only the AArch64 register XML, dropping AArch32 and HTML renderings.
find "$TMP" -name 'AArch64-*.xml' ! -path '*/xhtml/*' -exec cp {} "$SYSREG_DIR"/ \;
count=$(ls "$SYSREG_DIR" | wc -l | tr -d ' ')
echo "    $count AArch64 XML files in $SYSREG_DIR"

echo "==> downloading m1n1 apple_regs.json"
echo "    $APPLE_URL"
curl $CURL_OPTS -fL "$APPLE_URL" -o "$DIR/apple_regs.json"

echo
echo "Done. Now regenerate the Go table:"
echo "    make gen-sysregs"
