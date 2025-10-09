# Installation script for my-context (Windows PowerShell)
# Installs to $env:USERPROFILE\bin without requiring admin privileges

$ErrorActionPreference = "Stop"

$INSTALL_DIR = Join-Path $env:USERPROFILE "bin"
$BINARY_NAME = "my-context.exe"
$INSTALL_PATH = Join-Path $INSTALL_DIR $BINARY_NAME

Write-Host "Installing my-context for Windows (PowerShell)..."

# Detect existing installation
if (Test-Path $INSTALL_PATH) {
    Write-Host "Existing installation found."
    
    # Get current version
    try {
        & $INSTALL_PATH --version 2>$null
    } catch {
        Write-Host "Current version: unknown"
    }
    
    # Prompt for upgrade
    $upgrade = Read-Host "Existing installation found. Upgrade? (y/N)"
    if ($upgrade -ne "y" -and $upgrade -ne "Y") {
        Write-Host "Installation cancelled."
        exit 0
    }
    
    # Backup old binary
    $backupPath = $INSTALL_PATH + ".bak"
    Copy-Item $INSTALL_PATH $backupPath -Force
    Write-Host "Backed up existing binary"
}

# Create install directory
if (-not (Test-Path $INSTALL_DIR)) {
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Copy binary (assume binary provided as $args[0] or from bin/)
if ($args.Count -gt 0) {
    Copy-Item $args[0] $INSTALL_PATH -Force
} else {
    # Auto-detect binary from bin/
    $binarySrc = "bin\my-context-windows-amd64.exe"
    if (Test-Path $binarySrc) {
        Copy-Item $binarySrc $INSTALL_PATH -Force
    } else {
        Write-Error "Binary not found at $binarySrc"
        Write-Host "Run './scripts/build-all.sh' first to build binaries"
        exit 1
    }
}

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not ($currentPath -like "*$INSTALL_DIR*")) {
    Write-Host "Adding $INSTALL_DIR to user PATH..."
    $newPath = $currentPath + ";" + $INSTALL_DIR
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added to PATH. Restart terminal for changes to take effect."
    
    # Update current session PATH
    $env:Path = $newPath
} else {
    Write-Host "$INSTALL_DIR already in PATH"
}

# Verify installation
try {
    $version = & $INSTALL_PATH --version 2>&1
    Write-Host ""
    Write-Host "✓ Installation complete!" -ForegroundColor Green
    Write-Host $version
    
    # Remove backup if successful
    $backupPath = $INSTALL_PATH + ".bak"
    if (Test-Path $backupPath) {
        Remove-Item $backupPath -Force
    }
} catch {
    Write-Host ""
    Write-Error "Installation verification failed"
    exit 1
}

Write-Host ""
Write-Host "Note: ~/.my-context/ data directory is preserved (separate from binary)"
Write-Host "Restart your terminal for PATH changes to take effect"

