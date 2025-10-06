# Troubleshooting Guide

This guide covers common installation and usage issues for my-context across different platforms.

---

## Installation Issues

### WSL (Windows Subsystem for Linux)

#### Problem: "my-context.exe: cannot execute binary file"

**Cause**: You're trying to run the Windows executable in a Linux environment.

**Solution**:
```bash
# Download or build the Linux binary
wget https://github.com/user/my-context-copilot/releases/latest/download/my-context-linux-amd64
chmod +x my-context-linux-amd64
sudo mv my-context-linux-amd64 ~/.local/bin/my-context

# Or build from source
cd ~/projects/my-context-copilot
go build -o my-context ./cmd/my-context/
```

#### Problem: "Go: command not found" when building

**Cause**: Go is not installed in your WSL environment.

**Solution**:
```bash
# Install Go in WSL
sudo apt update
sudo apt install golang-go

# Verify installation
go version

# Then rebuild
cd ~/projects/my-context-copilot
go build -o my-context ./cmd/my-context/
```

---

### Windows (Native cmd.exe / PowerShell)

#### Problem: "'my-context' is not recognized as an internal or external command"

**Cause**: Binary not in PATH or PATH not refreshed after installation.

**Solution 1** (Restart shell):
```cmd
# Close and reopen your terminal window
# PATH changes require new session to take effect
```

**Solution 2** (Manual PATH check):
```cmd
# Check if binary exists
dir %USERPROFILE%\bin\my-context.exe

# Check current PATH
echo %PATH%

# If directory missing from PATH, add manually:
setx PATH "%PATH%;%USERPROFILE%\bin"
# Then restart terminal
```

**Solution 3** (PowerShell):
```powershell
# Check PATH
$env:Path

# Add to PATH if needed
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$env:USERPROFILE\bin", "User")
# Then restart terminal
```

#### Problem: "install.bat" fails with "Access Denied"

**Cause**: Trying to install to system directory without admin privileges.

**Solution**:
```cmd
# Use user-specific installation (default)
install.bat

# Binary installs to: %USERPROFILE%\bin\
# No admin required
```

#### Problem: PowerShell says "running scripts is disabled"

**Cause**: PowerShell execution policy blocks .ps1 scripts.

**Solution**:
```powershell
# Check current policy
Get-ExecutionPolicy

# Allow scripts for current user (no admin needed)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Then run installer
.\install.ps1
```

---

### macOS

#### Problem: "my-context cannot be opened because the developer cannot be verified"

**Cause**: macOS Gatekeeper blocks unsigned binaries.

**Solution**:
```bash
# Option 1: Right-click > Open (then click "Open" in dialog)
# This creates a one-time exception

# Option 2: Remove quarantine attribute
xattr -d com.apple.quarantine ~/Downloads/my-context-darwin-*
chmod +x ~/Downloads/my-context-darwin-*

# Option 3: Allow via System Preferences
# System Preferences > Security & Privacy > General > "Allow Anyway"
```

#### Problem: ARM Mac downloads Intel binary by mistake

**Cause**: Downloaded wrong binary for architecture.

**Solution**:
```bash
# Check your Mac architecture
uname -m
# Output: arm64 (Apple Silicon) or x86_64 (Intel)

# Download correct binary:
# arm64 → my-context-darwin-arm64
# x86_64 → my-context-darwin-amd64

# Intel Macs CAN run ARM binaries via Rosetta 2, but native is faster
```

---

### Linux

#### Problem: "Permission denied" when running binary

**Cause**: Binary not marked as executable.

**Solution**:
```bash
chmod +x my-context-linux-amd64
./my-context-linux-amd64 --version
```

#### Problem: "~/.local/bin not in PATH"

**Cause**: Installation directory not added to PATH.

**Solution**:
```bash
# Add to your shell rc file
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# For zsh users:
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify
echo $PATH | grep ".local/bin"
```

---

## Upgrade Issues

### Problem: "Old contexts not visible after upgrade"

**Cause**: This should NOT happen - backward compatibility is guaranteed.

**Solution**:
```bash
# Check if contexts still exist
ls -la ~/.my-context/

# Verify state file
cat ~/.my-context/state.json

# Try listing with --all flag
my-context list --all

# If contexts still missing, report bug with:
my-context --version
ls -la ~/.my-context/
```

### Problem: "New commands not recognized after upgrade"

**Cause**: Old binary still in use (cached or wrong PATH priority).

**Solution**:
```bash
# Check which binary is being used
which my-context

# Check version
my-context --version
# Should show v1.1.0 or later for Sprint 2 features

# If old version, reinstall:
rm $(which my-context)
./scripts/install.sh
```

---

## Runtime Issues

### Problem: "Context not found" but you just created it

**Cause**: Context name normalization (spaces → underscores).

**Solution**:
```bash
# Names with spaces are stored with underscores
my-context start "Bug Fix"
# Creates: Bug_Fix

# Use exact name from list command:
my-context list
my-context show Bug_Fix

# Or use original name - tool normalizes automatically
my-context show "Bug Fix"
```

### Problem: Special characters (like $) missing from notes

**Cause**: Bug in Sprint 1 (fixed in Sprint 2).

**Solution**:
```bash
# Upgrade to Sprint 2 (v1.1.0+)
my-context --version

# If old version, upgrade:
curl -sSL https://raw.githubusercontent.com/.../install.sh | bash
```

### Problem: "Cannot archive active context"

**Cause**: You must stop the context before archiving it.

**Solution**:
```bash
# Check active context
my-context show

# Stop it first
my-context stop

# Then archive
my-context archive "Context Name"

# Or switch to new context (auto-stops previous)
my-context start "New Context"
my-context archive "Old Context"
```

---

## Data Issues

### Problem: "Lost my context data!"

**Cause**: Data is stored in `~/.my-context/` which might be in different locations depending on environment.

**Solution**:
```bash
# Check where data is stored
echo $HOME/.my-context

# Windows users in WSL might have data in Windows home:
ls /mnt/c/Users/$USER/.my-context/

# Set custom location (optional)
export MY_CONTEXT_HOME=/custom/path/.my-context
my-context list
```

### Problem: Corrupted context files

**Cause**: Rare - might occur from interrupted writes or filesystem issues.

**Solution**:
```bash
# Check meta.json for valid JSON
cat ~/.my-context/Context_Name/meta.json | jq .

# If corrupted, manually fix JSON or restore from backup
# Context directory is self-contained - just copy from backup

# Last resort: delete corrupted context
rm -rf ~/.my-context/Corrupted_Context_Name/
```

---

## Build from Source Issues

### Problem: "go.mod not found" or module errors

**Cause**: Not in correct directory or Go modules not initialized.

**Solution**:
```bash
# Ensure you're in repository root
cd ~/projects/my-context-copilot

# Verify go.mod exists
ls go.mod

# Install dependencies
go mod download

# Build
go build -o my-context ./cmd/my-context/
```

### Problem: Build fails with dependency errors

**Cause**: Outdated dependencies or Go version mismatch.

**Solution**:
```bash
# Check Go version (need 1.21+)
go version

# Update dependencies
go mod tidy

# Clean build cache
go clean -cache -modcache

# Rebuild
go build ./cmd/my-context/
```

---

## Cross-Platform Issues

### Problem: Paths don't work correctly on Windows

**Cause**: Windows uses backslashes, internal storage uses forward slashes.

**Solution**: This should be handled automatically by the tool. If not:
```bash
# Use forward slashes in commands (tool converts)
my-context file ~/documents/file.txt

# Or use backslashes (tool normalizes)
my-context file C:\Users\Name\documents\file.txt

# Both work identically
```

### Problem: Context data not shared between Git Bash and WSL

**Cause**: Different home directories (`~` resolves differently).

**Solution**:
```bash
# Check where each environment stores data
# In Git Bash:
echo $HOME/.my-context
# Likely: /c/Users/YourName/.my-context

# In WSL:
echo $HOME/.my-context  
# Likely: /home/yourname/.my-context

# Option 1: Use MY_CONTEXT_HOME to point both to same location
export MY_CONTEXT_HOME=/mnt/c/Users/YourName/.my-context

# Option 2: Symlink
ln -s /mnt/c/Users/YourName/.my-context ~/.my-context
```

---

## Getting Help

### Check Version and Installation

```bash
# Verify version
my-context --version

# Check installation location
which my-context

# Verify data directory
ls -la ~/.my-context/
```

### Enable Debug Output (Future Feature)

Currently not available. For debugging:
```bash
# Check if context exists
cat ~/.my-context/state.json

# View context metadata
cat ~/.my-context/Context_Name/meta.json

# View logs directly
cat ~/.my-context/Context_Name/notes.log
cat ~/.my-context/transitions.log
```

### Report a Bug

1. Check you're on the latest version: `my-context --version`
2. Try to reproduce with minimal steps
3. Collect diagnostic info:
   ```bash
   my-context --version
   uname -a  # or systeminfo on Windows
   echo $SHELL
   ls -la ~/.my-context/
   ```
4. Open issue at: [GitHub repository URL]

---

## FAQ

**Q: Can I use my-context for non-development work?**  
A: Absolutely! It works for any context-switching workflow (event planning, research, project management, etc.).

**Q: Where is my data stored?**  
A: `~/.my-context/` directory. Plain text files you can view with any editor.

**Q: Can I use my-context on multiple computers?**  
A: Yes, but contexts are local to each machine. You can export contexts to share via `my-context export`.

**Q: How do I back up my contexts?**  
A: Copy `~/.my-context/` directory: `cp -r ~/.my-context ~/.my-context-backup` or use `my-context export --all`.

**Q: Can I edit context data manually?**  
A: Yes! All files are plain text JSON and logs. Just maintain the format.

**Q: Does this track my code or files automatically?**  
A: No, you manually associate files with `my-context file <path>`. Automatic tracking might come in future version.

---

*For additional help, see [README.md](../README.md) and [GitHub Issues](link)*
