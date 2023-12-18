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
    curl \
    ripgrep \
    eza \
    node

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
