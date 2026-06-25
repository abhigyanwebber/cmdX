#!/bin/bash

set -e

REPO="https://github.com/abhigyanwebber/cmd-customizer"
INSTALL_DIR="$HOME/.cmdx"
BIN_DIR="$HOME/.local/bin"
THEMES_DIR="$INSTALL_DIR/themes"

# ── Colors ────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
RESET='\033[0m'

print_banner() {
echo -e "${CYAN}${BOLD}"
cat << 'EOF'
 ██████╗███╗   ███╗██████╗ ██╗  ██╗
██╔════╝████╗ ████║██╔══██╗╚██╗██╔╝
██║     ██╔████╔██║██║  ██║ ╚███╔╝
██║     ██║╚██╔╝██║██║  ██║ ██╔██╗
╚██████╗██║ ╚═╝ ██║██████╔╝██╔╝ ██╗
 ╚═════╝╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝
EOF
echo -e "${RESET}"
echo -e "${CYAN}  cmd-customizer installer${RESET}"
echo -e "${YELLOW}  break free from the boring terminal${RESET}"
echo ""
}

step() { echo -e "${CYAN}${BOLD}==> ${RESET}${BOLD}$1${RESET}"; }
ok()   { echo -e "${GREEN}  ✓ $1${RESET}"; }
fail() { echo -e "${RED}  ✗ $1${RESET}"; exit 1; }
warn() { echo -e "${YELLOW}  ! $1${RESET}"; }

# ── Checks ────────────────────────────────────────────────
print_banner

step "Checking dependencies"

if ! command -v go &> /dev/null; then
    fail "Go is not installed. Install it from https://go.dev/dl/ and re-run this script."
fi
ok "Go found: $(go version)"

if ! command -v git &> /dev/null; then
    fail "Git is not installed. Install it from https://git-scm.com and re-run this script."
fi
ok "Git found: $(git version)"

# ── Clone ─────────────────────────────────────────────────
step "Cloning repository"

if [ -d "$INSTALL_DIR" ]; then
    warn "Found existing installation at $INSTALL_DIR — updating"
    cd "$INSTALL_DIR"
    git pull origin main
else
    git clone "$REPO" "$INSTALL_DIR"
    cd "$INSTALL_DIR"
fi
ok "Repository ready"

# ── Build ─────────────────────────────────────────────────
step "Building cmdx binary"

go mod tidy
go build -o "$INSTALL_DIR/cmdx" ./cmd/
ok "Binary built"

# ── Install ───────────────────────────────────────────────
step "Installing to $BIN_DIR"

mkdir -p "$BIN_DIR"
cp "$INSTALL_DIR/cmdx" "$BIN_DIR/cmdx"
chmod +x "$BIN_DIR/cmdx"
ok "Binary installed"

# ── Themes ────────────────────────────────────────────────
step "Installing themes"

mkdir -p "$THEMES_DIR"
cp "$INSTALL_DIR/themes/"*.json "$THEMES_DIR/"
ok "Themes installed to $THEMES_DIR"

# ── PATH ──────────────────────────────────────────────────
step "Configuring PATH"

SHELL_CONFIG=""
if [[ "$SHELL" == *"zsh"* ]]; then
    SHELL_CONFIG="$HOME/.zshrc"
elif [[ "$SHELL" == *"bash"* ]]; then
    SHELL_CONFIG="$HOME/.bashrc"
fi

if [ -n "$SHELL_CONFIG" ]; then
    if ! grep -q "$BIN_DIR" "$SHELL_CONFIG"; then
        echo "" >> "$SHELL_CONFIG"
        echo "# cmd-customizer" >> "$SHELL_CONFIG"
        echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$SHELL_CONFIG"
        ok "Added $BIN_DIR to PATH in $SHELL_CONFIG"
    else
        ok "PATH already configured"
    fi
else
    warn "Could not detect shell config. Manually add $BIN_DIR to your PATH."
fi

# ── Done ──────────────────────────────────────────────────
echo ""
echo -e "${GREEN}${BOLD}  Installation complete!${RESET}"
echo ""
echo -e "  Restart your terminal or run:"
echo -e "  ${CYAN}source $SHELL_CONFIG${RESET}"
echo ""
echo -e "  Then try:"
echo -e "  ${CYAN}cmdx theme list${RESET}"
echo -e "  ${CYAN}cmdx theme preview cyberpunk${RESET}"
echo -e "  ${CYAN}cmdx theme apply cyberpunk${RESET}"
echo ""