# cmd-customizer Windows Installer
# Run with: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

$REPO        = "https://github.com/abhigyanwebber/cmd-customizer"
$INSTALL_DIR = "$env:USERPROFILE\.cmdx"
$BIN_DIR     = "$env:USERPROFILE\.local\bin"
$THEMES_DIR  = "$INSTALL_DIR\themes"

# ── Colors ────────────────────────────────────────────────
function Write-Banner {
    Write-Host ""
    Write-Host " ██████╗███╗   ███╗██████╗ ██╗  ██╗" -ForegroundColor Cyan
    Write-Host "██╔════╝████╗ ████║██╔══██╗╚██╗██╔╝" -ForegroundColor Cyan
    Write-Host "██║     ██╔████╔██║██║  ██║ ╚███╔╝ " -ForegroundColor Cyan
    Write-Host "██║     ██║╚██╔╝██║██║  ██║ ██╔██╗ " -ForegroundColor Cyan
    Write-Host "╚██████╗██║ ╚═╝ ██║██████╔╝██╔╝ ██╗" -ForegroundColor Cyan
    Write-Host " ╚═════╝╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "  cmd-customizer installer" -ForegroundColor Cyan
    Write-Host "  break free from the boring terminal" -ForegroundColor Yellow
    Write-Host ""
}

function Write-Step  { param($msg) Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok    { param($msg) Write-Host "  [OK] $msg" -ForegroundColor Green }
function Write-Warn  { param($msg) Write-Host "  [!] $msg" -ForegroundColor Yellow }
function Write-Fail  { param($msg) Write-Host "  [FAIL] $msg" -ForegroundColor Red; exit 1 }

# ── Checks ────────────────────────────────────────────────
Write-Banner

Write-Step "Checking dependencies"

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Fail "Go is not installed. Install it from https://go.dev/dl/ and re-run."
}
Write-Ok "Go found: $(go version)"

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Fail "Git is not installed. Install it from https://git-scm.com and re-run."
}
Write-Ok "Git found: $(git version)"

# ── Clone ─────────────────────────────────────────────────
Write-Step "Cloning repository"

if (Test-Path $INSTALL_DIR) {
    Write-Warn "Found existing installation at $INSTALL_DIR — updating"
    Set-Location $INSTALL_DIR
    git pull origin main
} else {
    git clone $REPO $INSTALL_DIR
    Set-Location $INSTALL_DIR
}
Write-Ok "Repository ready"

# ── Build ─────────────────────────────────────────────────
Write-Step "Building cmdx binary"

go mod tidy
go build -o "$INSTALL_DIR\cmdx.exe" .\cmd\
Write-Ok "Binary built"

# ── Install ───────────────────────────────────────────────
Write-Step "Installing to $BIN_DIR"

if (-not (Test-Path $BIN_DIR)) {
    New-Item -ItemType Directory -Path $BIN_DIR | Out-Null
}

Copy-Item "$INSTALL_DIR\cmdx.exe" "$BIN_DIR\cmdx.exe" -Force
Write-Ok "Binary installed to $BIN_DIR"

# ── Themes ────────────────────────────────────────────────
Write-Step "Installing themes"

if (-not (Test-Path $THEMES_DIR)) {
    New-Item -ItemType Directory -Path $THEMES_DIR | Out-Null
}

Copy-Item "$INSTALL_DIR\themes\*.json" $THEMES_DIR -Force
Write-Ok "Themes installed to $THEMES_DIR"

# ── PATH ──────────────────────────────────────────────────
Write-Step "Configuring PATH"

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

if ($currentPath -notlike "*$BIN_DIR*") {
    [Environment]::SetEnvironmentVariable(
        "PATH",
        "$currentPath;$BIN_DIR",
        "User"
    )
    Write-Ok "Added $BIN_DIR to user PATH"
} else {
    Write-Ok "PATH already configured"
}

# ── PowerShell Profile ────────────────────────────────────
Write-Step "Setting up PowerShell profile"

$profileDir = Split-Path $PROFILE
if (-not (Test-Path $profileDir)) {
    New-Item -ItemType Directory -Path $profileDir | Out-Null
}

if (-not (Test-Path $PROFILE)) {
    New-Item -ItemType File -Path $PROFILE | Out-Null
}

$profileEntry = "`n# cmd-customizer`n`$env:PATH += `";$BIN_DIR`""

if (-not (Select-String -Path $PROFILE -Pattern "cmd-customizer" -Quiet)) {
    Add-Content -Path $PROFILE -Value $profileEntry
    Write-Ok "Added to PowerShell profile: $PROFILE"
} else {
    Write-Ok "PowerShell profile already configured"
}

# ── Done ──────────────────────────────────────────────────
Write-Host ""
Write-Host "  Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "  Restart your terminal, then try:" -ForegroundColor White
Write-Host "  cmdx theme list" -ForegroundColor Cyan
Write-Host "  cmdx theme preview cyberpunk" -ForegroundColor Cyan
Write-Host "  cmdx theme apply cyberpunk" -ForegroundColor Cyan
Write-Host ""