# jcaldwell-labs Organization Guidelines

> Standards for polishing and maintaining projects across the jcaldwell-labs organization.

**Version**: 1.0
**Updated**: 2025-12-23

## Repository Structure

### Root Directory

Keep the root directory clean. Only these files should be at the root:

**Required:**
- `README.md` — Project overview and quick start
- `LICENSE` — MIT License
- `.gitignore` — Git ignore patterns
- `llms.txt` — AI discoverability file

**Language-specific (as needed):**
- `go.mod`, `go.sum` — Go projects
- `package.json` — Node.js projects
- `requirements.txt` — Python projects
- `Cargo.toml` — Rust projects
- `Makefile` — Build automation

**Move everything else** to `docs/` or `.github/planning/`.

### Documentation Hierarchy

```
docs/
├── README.md           # Documentation hub (index)
├── guides/             # User-focused how-to guides
├── tutorials/          # Step-by-step walkthroughs
├── examples/           # Sample code and configurations
└── blog/               # Articles and thought pieces

.github/
├── planning/           # Internal planning docs
│   ├── ROADMAP.md     # Feature roadmap
│   ├── backlog.md     # Feature backlog
│   └── sprints/       # Sprint notes
├── templates/          # Reusable templates
├── ISSUE_TEMPLATE/     # Issue templates
├── workflows/          # GitHub Actions
└── PULL_REQUEST_TEMPLATE.md
```

## README Best Practices

### Quality Bar

> Can a stranger understand what this does in 30 seconds?

### Required Sections

1. **Title + Badges** — Project name, license, language, PRs welcome
2. **Tagline** — One sentence explaining what it does
3. **Why [Project]?** — Value proposition and key benefits
4. **Quick Start** — Installation and first steps (< 5 minutes)
5. **Features** — Core capabilities with examples
6. **Documentation** — Links to guides and tutorials
7. **Contributing** — How to contribute
8. **License** — License type and link

### Optional Sections

- **Demo** — GIF, screenshot, or asciinema recording
- **Comparison** — How this differs from alternatives
- **Roadmap** — Link to planning docs
- **Community** — Links to discussions, Discord, etc.

## AI Discoverability

### llms.txt File

Every project should have an `llms.txt` file for AI assistant discovery:

```
# project-name

> Tagline describing the project.

Brief description of what the project does and its purpose.

## What it does
[Core capabilities]

## Quick start
[Installation and first commands]

## Common commands/patterns
[Most-used operations]

## Use cases
[When to use this tool]

## Repository
[GitHub URL]
```

### GitHub Topics

Add 5-10 relevant topics to improve discoverability:
- Language: `go`, `golang`, `cli`
- Domain: `developer-tools`, `productivity`, `context-management`
- Type: `command-line`, `terminal`, `unix`

## GitHub Configuration

### Required Settings

- [ ] Repository description (70-120 characters, keyword-rich)
- [ ] Topics/tags (5-10 relevant tags)
- [ ] Issue templates (bug report, feature request)
- [ ] Pull request template
- [ ] Branch protection on main

### Recommended

- [ ] Discussions enabled (for community projects)
- [ ] Social preview image (1200×630 px)
- [ ] Homepage URL (if docs site exists)

## Project Maturity Levels

### L1 — Experimental
- Basic README exists
- Code runs but may have rough edges
- Minimal documentation

### L2 — Functional
- Clear README with installation
- Basic documentation
- Consistent code style
- Some tests

### L3 — Production
- Polished README with badges
- Complete documentation structure
- Comprehensive tests
- CI/CD pipeline
- Issue/PR templates

### L4 — Showcase
- All L3 requirements
- Blog post or demo video
- Social proof (stars, usage)
- Active maintenance

## Launch Checklist

Before announcing a project:

- [ ] README passes 30-second test
- [ ] Quick start works in < 5 minutes
- [ ] llms.txt created
- [ ] GitHub topics added
- [ ] Issue templates configured
- [ ] Root directory clean (< 5 markdown files)
- [ ] At least one visual (screenshot/GIF/diagram)

## Templates

Reusable templates available in `.github/templates/`:
- `README-TEMPLATE.md` — Standard README structure
- `llms-template.txt` — AI discoverability file template

## Reference Projects

- **my-grid** — Terminal canvas editor (L4 showcase)
- **my-context** — Context tracking CLI (L3 production)
