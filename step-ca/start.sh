#!/bin/bash
# Reads the CA password from macOS Keychain and starts step-ca.
# Uses a file descriptor to avoid the password appearing in process listings or logs.
set -euo pipefail
exec step-ca "$STEPPATH/config/ca.json" --password-file <(security find-generic-password -a "step-ca" -s "dotfiles-ca" -w)
