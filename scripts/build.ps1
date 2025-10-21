# build.ps1 - PowerShell build script with proper metadata
# Usage: .\scripts\build.ps1 [-Install]

param(
    [switch]$Install
)

Write-Host "Building my-context for Windows..." -ForegroundColor Blue

# Get version from git
$VERSION = git describe --tags --exact-match 2>$null
if (-not $VERSION) {
    $VERSION = git describe --tags --always 2>$null
}
if (-not $VERSION) {
    $VERSION = "dev"
}

# Get build time (UTC)
$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

# Get git commit
$GIT_COMMIT = git rev-parse --short HEAD 2>$null
if (-not $GIT_COMMIT) {
    $GIT_COMMIT = "unknown"
}

Write-Host "   Version: $VERSION" -ForegroundColor Gray
Write-Host "   Build: $BUILD_TIME" -ForegroundColor Gray
Write-Host "   Commit: $GIT_COMMIT" -ForegroundColor Gray
Write-Host ""

# Build with ldflags
$ldflags = "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"
go build -ldflags $ldflags -o my-context.exe cmd/my-context/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Build complete: my-context.exe" -ForegroundColor Green
Write-Host ""
.\my-context.exe --version

if ($Install) {
    Write-Host ""
    Write-Host "Installing to user PATH..." -ForegroundColor Yellow

    # Install to %USERPROFILE%\bin or %LOCALAPPDATA%\Programs
    $installDir = "$env:USERPROFILE\bin"
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir | Out-Null
        Write-Host "Created directory: $installDir" -ForegroundColor Gray
    }

    Copy-Item my-context.exe $installDir\my-context.exe -Force
    Write-Host "Installed: $installDir\my-context.exe" -ForegroundColor Green

    # Check if in PATH
    if ($env:PATH -notlike "*$installDir*") {
        Write-Host ""
        Write-Host "Add to PATH:" -ForegroundColor Yellow
        Write-Host "  setx PATH `"%PATH%;$installDir`"" -ForegroundColor Gray
        Write-Host "  (Restart terminal after running)" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "To install: .\scripts\build.ps1 -Install"
Write-Host "Or manually: copy my-context.exe to your PATH"
