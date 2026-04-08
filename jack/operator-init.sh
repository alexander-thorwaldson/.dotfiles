#!/bin/bash
# Registers or logs in the operator on the local Matrix homeserver
# and writes the token to ~/.jack/operator/token for use with jack msg.
#
# Usage:
#   operator-init.sh [username]
#
# The token is also exported as JACK_MSG_TOKEN. Source this script or
# add the fish integration to have it available in every shell.
set -euo pipefail

HOMESERVER="http://localhost:6167"
REGISTRATION_TOKEN="jack-local-dev"
USERNAME="${1:-operator}"
DATA_DIR="${JACK_DATA_DIR:-$HOME/.jack}"
TOKEN_DIR="$DATA_DIR/operator"
TOKEN_FILE="$TOKEN_DIR/token"

mkdir -p "$TOKEN_DIR"

# If we already have a token, try to validate it.
if [ -f "$TOKEN_FILE" ]; then
    EXISTING_TOKEN=$(cat "$TOKEN_FILE")
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
        -H "Authorization: Bearer $EXISTING_TOKEN" \
        "$HOMESERVER/_matrix/client/v3/account/whoami")
    if [ "$STATUS" = "200" ]; then
        echo "operator already authenticated as @${USERNAME}:localhost"
        echo "$EXISTING_TOKEN"
        exit 0
    fi
    echo "existing token expired, re-authenticating..."
fi

# Try to register first (first user gets admin).
REGISTER_BODY=$(cat <<EOF
{
    "auth": {"type": "m.login.registration_token", "token": "$REGISTRATION_TOKEN"},
    "username": "$USERNAME",
    "password": "$USERNAME"
}
EOF
)

RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$REGISTER_BODY" \
    "$HOMESERVER/_matrix/client/v3/register")

TOKEN=$(echo "$RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

# If registration failed (user exists), try login.
if [ -z "$TOKEN" ]; then
    LOGIN_BODY=$(cat <<EOF
{
    "type": "m.login.password",
    "user": "$USERNAME",
    "password": "$USERNAME"
}
EOF
    )

    RESPONSE=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$LOGIN_BODY" \
        "$HOMESERVER/_matrix/client/v3/login")

    TOKEN=$(echo "$RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ]; then
    echo "failed to authenticate: $RESPONSE" >&2
    exit 1
fi

# Store the token.
echo -n "$TOKEN" > "$TOKEN_FILE"
chmod 600 "$TOKEN_FILE"

echo "authenticated as @${USERNAME}:localhost"
echo "$TOKEN"
