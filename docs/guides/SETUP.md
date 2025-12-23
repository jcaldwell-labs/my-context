# Installation and Setup Guide

## Prerequisites - Install Go

You need to install Go 1.21 or later to build this project.

### Install Go on Windows

**Option 1: Using Chocolatey (recommended)**
```bash
choco install golang
```

**Option 2: Using winget**
```bash
winget install GoLang.Go
```

**Option 3: Manual Download**
1. Download from https://go.dev/dl/
2. Run the installer (go1.21.x.windows-amd64.msi)
3. Restart your terminal/git-bash

### Verify Installation

After installing, restart your terminal and verify:
```bash
go version
# Should show: go version go1.21.x windows/amd64
```

---

## Project Setup (Run after Go is installed)

### Step 1: Initialize Go Module
```bash
cd /c/Users/JefferyCaldwell/projects/my-context-copilot
go mod init github.com/yourusername/my-context-copilot
```

### Step 2: Install Dependencies
```bash
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/stretchr/testify@latest
go mod tidy
```

### Step 3: Verify Structure
```bash
ls -la
# Should see: cmd/, internal/, tests/, specs/, go.mod, go.sum
```

---

## Quick Start After Setup

### Build the Project
```bash
go build -o my-context.exe ./cmd/my-context/
```

### Run Tests
```bash
go test ./...
```

### Run the CLI
```bash
./my-context.exe --help
```

---

## Development Workflow

1. **T001-T003**: ✅ Directory structure created
2. **T004-T013**: Write tests first (TDD)
3. **T014-T019**: Implement models
4. **T020-T022**: Implement core logic
5. **T023-T033**: Implement commands
6. **T034-T042**: Integration, build, polish

See `specs/001-cli-context-management/tasks.md` for detailed task list.

---

## Troubleshooting

### Go command not found
- Ensure Go is installed (see above)
- Restart terminal after installation
- Check PATH: `echo $PATH | grep -i go`

### Permission issues on Windows
- Run git-bash as Administrator if needed
- Or use PowerShell/cmd.exe

### Module download issues
```bash
# Set Go proxy if behind firewall
go env -w GOPROXY=https://proxy.golang.org,direct
```

---

## Next Steps

Once Go is installed and `go version` works:

1. Run the setup commands above (Step 1-3)
2. Follow the implementation guide in `IMPLEMENTATION.md`
3. Use `HERE.md` as your scratchpad for development

**Status**: Directory structure ready, waiting for Go installation to proceed with implementation.
