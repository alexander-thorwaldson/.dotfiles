#!/bin/bash
# Issues a certificate from the local CA using the keychain-stored password.
#
# Usage:
#   issue-cert.sh <common-name> <cert-out> <key-out> [--san value]... [--ou value]...
#
# Examples:
#   issue-cert.sh kuang ./certs/kuang.crt ./certs/kuang.key --san localhost --san 127.0.0.1
#   issue-cert.sh agent-alice ./alice.crt ./alice.key --ou gh_pr_list --ou gh_pr_view --ou pnpm_run
set -euo pipefail

if [ $# -lt 3 ]; then
    echo "Usage: $0 <common-name> <cert-out> <key-out> [--san value]... [--ou value]..."
    exit 1
fi

CN="$1"
CERT_OUT="$2"
KEY_OUT="$3"
shift 3

# Collect extra flags (--san, --ou, etc.)
EXTRA_FLAGS=()
while [ $# -gt 0 ]; do
    case "$1" in
        --san|--ou)
            EXTRA_FLAGS+=("$1" "$2")
            shift 2
            ;;
        *)
            echo "Unknown flag: $1"
            exit 1
            ;;
    esac
done

exec step ca certificate "$CN" "$CERT_OUT" "$KEY_OUT" \
    --provisioner admin \
    --provisioner-password-file <(security find-generic-password -a "step-ca" -s "dotfiles-ca" -w) \
    "${EXTRA_FLAGS[@]}"
