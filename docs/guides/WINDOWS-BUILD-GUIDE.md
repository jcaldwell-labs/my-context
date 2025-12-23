# Building my-context on Windows

**For**: Windows users building from source
**Version**: Works with v1.0.0, v2.0.0, v2.2.0-beta

---

## Quick Start

### Option 1: PowerShell (Recommended)

```powershell
# From project directory
.\scripts\build.ps1

# Or build and install in one step
.\scripts\build.ps1 -Install
```

### Option 2: Command Prompt (cmd.exe)

```cmd
REM From project directory
scripts\build.bat
```

---

## The Quoting Issue (Why Your Build Failed)

### ❌ What Doesn't Work in cmd.exe

```cmd
REM This fails with "named files must be .go files"
go build -ldflags -X main.Version=%VERSION% ... -o my-context.exe cmd/my-context/main.go
```

**Problem**: Missing quotes around ldflags value

---

### ✅ What Works in cmd.exe

```cmd
REM Note the = sign and quotes!
go build -ldflags="-X main.Version=%VERSION% -X main.BuildTime=%BUILD_TIME% -X main.GitCommit=%GIT_COMMIT%" -o my-context.exe cmd/my-context/main.go
```

**Key differences from Linux/bash**:
1. Use `-ldflags="..."` with `=` sign
2. Quotes around the entire `-X ...` string
3. Use `%VAR%` not `$VAR`
4. Output `my-context.exe` (include .exe extension)

---

## Step-by-Step Manual Build (cmd.exe)

**If scripts don't work, here's the manual process**:

```cmd
REM 1. Navigate to project
cd C:\Users\YourName\projects\my-context

REM 2. Get version from git
git describe --tags --exact-match
REM Output: v1.0.0 (or v2.0.0, v2.2.0-beta)

REM 3. Build with metadata (ONE LINE, copy carefully!)
go build -ldflags="-X main.Version=v1.0.0 -X main.BuildTime=2025-10-11T18:26:48Z -X main.GitCommit=bace3a8" -o my-context.exe cmd/my-context/main.go

REM 4. Test it
my-context.exe --version

REM Should show:
REM my-context version v1.0.0 (build: 2025-10-11T18:26:48Z, commit: bace3a8)

REM 5. Install to PATH
copy my-context.exe %USERPROFILE%\bin\my-context.exe
```

---

## Step-by-Step Manual Build (PowerShell)

**Alternative with better variable handling**:

```powershell
# 1. Navigate to project
cd C:\Users\YourName\projects\my-context

# 2. Get version
$VERSION = git describe --tags --exact-match
# Or: $VERSION = "v1.0.0"

# 3. Get build time (UTC)
$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

# 4. Get commit
$GIT_COMMIT = git rev-parse --short HEAD

# 5. Build (PowerShell handles quotes better!)
go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o my-context.exe cmd/my-context/main.go

# 6. Test
.\my-context.exe --version

# 7. Install
Copy-Item my-context.exe $env:USERPROFILE\bin\my-context.exe
```

---

## Installation Locations (Windows)

**Option 1**: User bin directory
```
%USERPROFILE%\bin\my-context.exe
C:\Users\YourName\bin\my-context.exe
```

**Option 2**: Local programs
```
%LOCALAPPDATA%\Programs\my-context\my-context.exe
C:\Users\YourName\AppData\Local\Programs\my-context\my-context.exe
```

**Option 3**: Current directory
```
C:\tools\my-context.exe
C:\Users\YourName\tools\my-context.exe
```

**Add to PATH**:
```cmd
REM Temporary (current session only)
set PATH=%PATH%;C:\Users\YourName\bin

REM Permanent (requires restart)
setx PATH "%PATH%;C:\Users\YourName\bin"
```

---

## Troubleshooting

### Error: "named files must be .go files"

**Cause**: ldflags quoting issue
**Fix**: Use `-ldflags="..."` with `=` sign and quotes

**cmd.exe syntax**:
```cmd
go build -ldflags="-X main.Version=v1.0.0 ..." -o my-context.exe cmd/my-context/main.go
```

**PowerShell syntax**:
```powershell
go build -ldflags "-X main.Version=v1.0.0 ..." -o my-context.exe cmd/my-context/main.go
```

---

### Version shows "unknown"

**Cause**: Built without ldflags
**Fix**: Rebuild using scripts or manual command above

---

### Go not found

**Cause**: Go not in PATH
**Fix**:
```cmd
REM Add Go to PATH
set PATH=%PATH%;C:\Program Files\Go\bin
```

Or reinstall Go: https://go.dev/dl/

---

## Quick Commands for Windows Users

### Using cmd.exe

```cmd
REM Clone from GitHub
git clone https://github.com/jcaldwell1066/my-context.git
cd my-context
git checkout main

REM Build v1.0.0
scripts\build.bat

REM Test
my-context.exe --version

REM Use
my-context.exe start "my-work"
my-context.exe note "task complete"
my-context.exe show
```

### Using PowerShell

```powershell
# Clone from GitHub
git clone https://github.com/jcaldwell1066/my-context.git
cd my-context
git checkout main

# Build and install
.\scripts\build.ps1 -Install

# Test
my-context --version

# Use
my-context start "my-work"
my-context note "task complete"
my-context show
```

---

## For v2.2.0-beta (Latest Features)

```cmd
REM Get latest development build
git clone https://github.com/jcaldwell1066/my-context.git
cd my-context
git checkout dev

REM Build
scripts\build.bat

REM Has signal and watch commands!
my-context.exe signal --help
my-context.exe watch --help
```

**Note**: v2.2.0-beta is development version, use v1.0.0 or v2.0.0 for stability

---

## Summary

**For v1.0.0** (stable):
- Branch: `main`
- Commit: bace3a8
- Script: `scripts\build.bat` or `scripts\build.ps1`

**For v2.0.0** (latest stable):
- Branch: `main` (after we merge)
- Currently on: `dev` branch
- Has: Lifecycle improvements, timestamps

**For v2.2.0-beta** (bleeding edge):
- Branch: `dev` or `005-context-signaling-protocol`
- Has: Signals, watch, metadata
- Status: Beta testing

**Recommendation for Windows users**: Use v1.0.0 (main branch) until v2.0.0 officially on main

---

**Create these Windows build scripts in project?**