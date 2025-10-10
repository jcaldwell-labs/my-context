## Description

<!-- Provide a brief description of the changes in this PR -->

## Type of Change

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Refactoring (no functional changes)

## Constitutional Compliance Checklist

### Core Principles (6)

- [ ] **I. Unix Philosophy**: Composability, text I/O, single-purpose
  - Each command does one thing well
  - Text input/output (stdin/stdout/stderr)
  - Chainable with shell tools (grep, awk, jq)
- [ ] **II. Cross-Platform Compatibility**: Windows, Linux, macOS, WSL2
  - Path handling works across platforms
  - No platform-specific code without fallbacks
  - Tested on at least 2 platforms
- [ ] **III. Stateful Context Management**: Single active context
  - Only one context active at a time
  - Automatic transitions (no orphaned states)
  - Operations default to active context
- [ ] **IV. Minimal Surface Area**: ≤12 commands total
  - Command count still within limit (currently 11/12)
  - New flags have sensible defaults
  - Single-letter aliases provided
- [ ] **V. Data Portability**: Plain text only
  - No binary databases or proprietary formats
  - JSON, logs, or markdown only
  - Manual editing with standard tools supported
- [ ] **VI. User-Driven Design**: From observed behavior
  - Feature validates real user need (not speculation)
  - Based on retrospective or observed patterns
  - Automates existing manual workflow

### Test Requirements

- [ ] **TDD Followed**: Tests written before implementation
- [ ] **Integration Tests**: Command contracts tested
- [ ] **Unit Tests**: Core logic covered
- [ ] **Backward Compatibility**: Old data works with new code
- [ ] **Performance**: Targets met (<1s for list/export/search)

### Documentation

- [ ] **README Updated**: New commands/flags documented
- [ ] **TROUBLESHOOTING**: Common issues addressed
- [ ] **Code Comments**: Public APIs and complex logic documented
- [ ] **Changelog**: User-facing changes noted

### Quality Gates

- [ ] **No Known Data Loss**: Destructive operations validated
- [ ] **Error Handling**: Edge cases handled gracefully
- [ ] **Linter Passing**: No new linter errors introduced
- [ ] **Build Successful**: All platforms compile
- [ ] **Manual Testing**: Changes tested locally

## Related Issues

<!-- Link to related issues using #issue_number -->

Closes #

## Testing Performed

<!-- Describe the testing you've done -->

### Platforms Tested
- [ ] Linux
- [ ] macOS (Intel)
- [ ] macOS (ARM)
- [ ] Windows (cmd.exe)
- [ ] Windows (PowerShell)
- [ ] Windows (git-bash)
- [ ] WSL2

### Test Results
<!-- Paste relevant test output or describe manual testing -->

```
go test ./...
```

## Performance Impact

<!-- If applicable, include benchmark results -->

```
go test -bench=. ./tests/benchmarks/
```

## Breaking Changes

<!-- List any breaking changes and migration path -->

- [ ] N/A - No breaking changes
- [ ] Breaking changes documented in CHANGELOG.md
- [ ] Migration guide provided

## Checklist Before Merge

- [ ] All constitutional compliance checks passed
- [ ] Tests passing on CI/CD
- [ ] Code reviewed and approved
- [ ] Documentation updated
- [ ] No merge conflicts with master

---

**Constitution Version**: 1.0.0 (Ratified 2025-10-09)  
**See**: `.specify/memory/constitution.md` for details

