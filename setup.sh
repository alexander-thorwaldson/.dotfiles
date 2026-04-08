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
# Kuang — MCP tool gateway daemon
# ──────────────────────────────────────────────
echo "#> Building kuang..."
(cd "$DOTFILES_DIR/kuang" && go build -o kuang ./cmd)

KUANG_BIN="$DOTFILES_DIR/kuang/kuang"
PLIST_DST=~/Library/LaunchAgents/com.dotfiles.kuang.plist
sed -e "s|__KUANG_BIN__|${KUANG_BIN}|g" -e "s|__PATH__|${PATH}|g" "$DOTFILES_DIR/kuang/com.dotfiles.kuang.plist.in" > "$PLIST_DST"
launchctl bootout gui/$(id -u) "$PLIST_DST" 2>/dev/null || true
launchctl bootstrap gui/$(id -u) "$PLIST_DST"
echo "#> kuang daemon registered at http://localhost:8080"

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
