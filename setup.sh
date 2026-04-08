#!/bin/bash
set -e

echo "#> Dotfiles setup — macOS provisioning"
echo ""

# ──────────────────────────────────────────────
# Homebrew
# ──────────────────────────────────────────────
if ! command -v brew &>/dev/null; then
    echo "#> Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    eval "$(/opt/homebrew/bin/brew shellenv)"
fi

# ──────────────────────────────────────────────
# Packages
# ──────────────────────────────────────────────
echo "#> Installing packages..."
brew install \
    fish \
    alacritty \
    neovim \
    tmux \
    starship \
    git \
    gh \
    curl \
    ripgrep \
    fd \
    fzf \
    eza \
    bat \
    git-delta \
    tree \
    node \
    pnpm \
    go \
    python@3.13 \
    jq \
    yq \
    httpie \
    lazygit \
    lazydocker \
    colima \
    go-jira \
    step \
    docker \
    docker-compose

# ──────────────────────────────────────────────
# Symlink dotfiles
# ──────────────────────────────────────────────
echo "#> Linking dotfiles..."
bash "$(dirname "$0")/install.sh"

# ──────────────────────────────────────────────
# Claude Code CLI
# ──────────────────────────────────────────────
if ! command -v claude &>/dev/null; then
    echo "#> Installing Claude Code CLI..."
    npm install -g @anthropic-ai/claude-code
fi

# ──────────────────────────────────────────────
# Fish shell
# ──────────────────────────────────────────────
echo "#> Setting up fish..."

FISH_PATH="$(command -v fish)"

# Add fish to allowed shells if needed
if ! grep -q "$FISH_PATH" /etc/shells 2>/dev/null; then
    echo "$FISH_PATH" | sudo tee -a /etc/shells
fi

# Set fish as default shell
if [ "$SHELL" != "$FISH_PATH" ]; then
    chsh -s "$FISH_PATH"
    echo "#> Default shell set to fish"
fi

# Install fisher
fish -c '
    curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher
    fisher update
'
echo "#> Fisher plugins installed"

# ──────────────────────────────────────────────
# ICE — prompt injection classifier daemon
# ──────────────────────────────────────────────
echo "#> Installing ice..."
DOTFILES_DIR="$(cd "$(dirname "$0")" && pwd)"
pip3 install --break-system-packages -q "$DOTFILES_DIR/ice"

# Download + cache the ONNX model on first setup
echo "#> Warming ice model cache (first run only)..."
python3 -c "from ice.classifier import Classifier; Classifier()" 2>/dev/null || true

# Install launchd plist — template in the real python3 path
mkdir -p ~/Library/LaunchAgents
PYTHON3_PATH="$(command -v python3)"
PLIST_DST=~/Library/LaunchAgents/com.dotfiles.ice.plist
sed "s|__PYTHON3__|${PYTHON3_PATH}|g" "$DOTFILES_DIR/ice/com.dotfiles.ice.plist.in" > "$PLIST_DST"
launchctl bootout gui/$(id -u) "$PLIST_DST" 2>/dev/null || true
launchctl bootstrap gui/$(id -u) "$PLIST_DST"
echo "#> ice daemon registered"

# ──────────────────────────────────────────────
# Docker infrastructure (via Colima)
# ──────────────────────────────────────────────
echo "#> Starting colima..."
brew services start colima 2>/dev/null || true
if ! colima status &>/dev/null; then
    colima start
fi

echo "#> Starting docker services..."
docker compose -f "$DOTFILES_DIR/docker-compose.yml" up -d
echo "#> tuwunel running at http://localhost:6167"

# ──────────────────────────────────────────────
# Step CA — local certificate authority
# ──────────────────────────────────────────────
STEP_HOME="$HOME/.step"
if [ ! -f "$STEP_HOME/config/ca.json" ]; then
    echo "#> Initializing step-ca..."
    CA_PASSWORD=$(openssl rand -base64 32)
    step ca init --name "dotfiles-ca" --provisioner "admin" --dns "localhost" --address ":9443" --password-file <(echo "$CA_PASSWORD")
    security add-generic-password -a "step-ca" -s "dotfiles-ca" -w "$CA_PASSWORD" -U
    # Remove any cleartext password files left by step ca init
    rm -f "$STEP_HOME/secrets/password"
fi

PLIST_DST=~/Library/LaunchAgents/com.dotfiles.step-ca.plist
sed -e "s|__DOTFILES_DIR__|${DOTFILES_DIR}|g" -e "s|__STEP_HOME__|${STEP_HOME}|g" -e "s|__PATH__|${PATH}|g" "$DOTFILES_DIR/step-ca/com.dotfiles.step-ca.plist.in" > "$PLIST_DST"
launchctl bootout gui/$(id -u) "$PLIST_DST" 2>/dev/null || true
launchctl bootstrap gui/$(id -u) "$PLIST_DST"
echo "#> step-ca daemon registered at https://localhost:9443"

# Wait for step-ca to be ready
sleep 2

# ──────────────────────────────────────────────
# Kuang — MCP tool gateway daemon
# ──────────────────────────────────────────────
echo "#> Building kuang..."
(cd "$DOTFILES_DIR/kuang" && go build -o kuang ./cmd)

KUANG_CERTS="$DOTFILES_DIR/kuang/certs"
mkdir -p "$KUANG_CERTS"

# Bundle root + intermediate CA certs for client verification
cat "$STEP_HOME/certs/root_ca.crt" "$STEP_HOME/certs/intermediate_ca.crt" > "$KUANG_CERTS/root_ca.crt"

# Issue kuang server cert if missing
if [ ! -f "$KUANG_CERTS/kuang.crt" ]; then
    echo "#> Issuing kuang server certificate..."
    "$DOTFILES_DIR/step-ca/issue-cert.sh" kuang "$KUANG_CERTS/kuang.crt" "$KUANG_CERTS/kuang.key" \
        --san localhost --san 127.0.0.1
fi

KUANG_BIN="$DOTFILES_DIR/kuang/kuang"
PLIST_DST=~/Library/LaunchAgents/com.dotfiles.kuang.plist
sed -e "s|__KUANG_BIN__|${KUANG_BIN}|g" -e "s|__CERTS_DIR__|${KUANG_CERTS}|g" -e "s|__PATH__|${PATH}|g" "$DOTFILES_DIR/kuang/com.dotfiles.kuang.plist.in" > "$PLIST_DST"
launchctl bootout gui/$(id -u) "$PLIST_DST" 2>/dev/null || true
launchctl bootstrap gui/$(id -u) "$PLIST_DST"
echo "#> kuang daemon registered at https://localhost:7117"

# Start cert auto-renewal daemon
STEP_BIN="$(command -v step)"
UID_VAL="$(id -u)"
PLIST_DST=~/Library/LaunchAgents/com.dotfiles.kuang-renew.plist
sed -e "s|__STEP_BIN__|${STEP_BIN}|g" -e "s|__STEP_HOME__|${STEP_HOME}|g" -e "s|__CERTS_DIR__|${KUANG_CERTS}|g" -e "s|__UID__|${UID_VAL}|g" "$KUANG_CERTS/com.dotfiles.kuang-renew.plist.in" > "$PLIST_DST"
launchctl bootout gui/$(id -u) "$PLIST_DST" 2>/dev/null || true
launchctl bootstrap gui/$(id -u) "$PLIST_DST"
echo "#> kuang cert auto-renewal daemon registered"

# ──────────────────────────────────────────────
# Neovim
# ──────────────────────────────────────────────
echo "#> Bootstrapping neovim plugins..."
nvim --headless "+Lazy! sync" +qa 2>/dev/null || true
echo "#> Neovim plugins synced"

# ──────────────────────────────────────────────
# Done
# ──────────────────────────────────────────────
echo ""
echo "#> Setup complete"
echo "#>   Log out and back in (or run 'fish') to use the new shell"
