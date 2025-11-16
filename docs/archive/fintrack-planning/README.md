# FinTrack Planning Archive

This directory contains the planning documentation created during the initial design phase of FinTrack, a terminal-based personal finance tracking and budgeting application.

## Repository Location

FinTrack now lives in its own repository:
- **GitHub:** https://github.com/jcaldwell1066/fintrack
- **Local:** ~/projects/fintrack

## Archived Documents

- **FINANCE_TRACKER_PLAN.md** - Complete system design (36KB)
  - Architecture overview
  - Database schema design
  - Command interface specifications
  - Data flow examples
  - Development phases

- **FINTRACK_QUICKREF.md** - Command reference guide (10KB)
  - Quick command cheat sheet
  - Common workflows
  - Examples and tips

- **FINTRACK_ROADMAP.md** - Implementation timeline (20KB)
  - 5-phase development plan
  - Technical tasks and milestones
  - Success criteria
  - Testing strategy

- **fintrack_schema.sql** - PostgreSQL database schema (20KB)
  - Complete table definitions
  - Indexes and triggers
  - Seed data
  - Materialized views

- **fintrack_config.example.yaml** - Configuration template (10KB)
  - All configuration options documented
  - Bank CSV format mappings
  - Default settings

## Timeline

- **Planning:** November 16, 2025
- **Phase 1 Foundation:** November 16, 2025
- **Moved to Separate Repository:** November 16, 2025

## Status

FinTrack Phase 1 (MVP) foundation is complete:
- ✅ Account management (CRUD operations)
- ✅ PostgreSQL backend with GORM
- ✅ Configuration management
- ✅ Output formatters (table/JSON)
- ✅ Unit test suite
- ✅ Build system

**Next:** Continue Phase 1 with transaction management, categories, CSV import, and reporting.

## Related

These planning documents were created as part of the my-context project but represent a separate application. They are archived here for reference and historical context.

For active development, see the FinTrack repository.
