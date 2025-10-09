---
name: Bug Report
about: Report a bug or unexpected behavior
title: '[BUG] '
labels: bug
assignees: ''
---

## Bug Description

A clear and concise description of what the bug is.

## Steps to Reproduce

1. Run command: `my-context ...`
2. Observe: ...
3. Expected: ...
4. Actual: ...

## Environment

**Platform** (check all that apply):
- [ ] Linux
- [ ] macOS (Intel)
- [ ] macOS (Apple Silicon / ARM)
- [ ] Windows (native cmd.exe)
- [ ] Windows (PowerShell)
- [ ] Windows (Git Bash)
- [ ] WSL (Windows Subsystem for Linux)

**Version**:
```
my-context --version
```

**OS Details**:
```bash
# Linux/macOS/WSL
uname -a

# Windows
systeminfo | findstr /B /C:"OS Name" /C:"OS Version"
```

## Context Data

**Number of contexts**:
```bash
my-context list --all | wc -l
```

**Active context**:
```bash
my-context show
```

## Additional Context

- Screenshots (if applicable)
- Error messages (full output)
- Workarounds you've tried

## Expected Behavior

What you expected to happen.

## Actual Behavior

What actually happened.

