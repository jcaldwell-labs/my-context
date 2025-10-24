# My-Context Quick Start Card
## Keep this handy for the first week

---

## Your First Context (Right Now!)

```bash
# Start
my-context start "Working on [task name]"

# Work for a bit...

# Add what you learned
my-context note "Key decision: [what you decided and why]"

# Track files you changed
my-context file path/to/file.go

# Check progress anytime
my-context show

# Stop when done
my-context stop
```

---

## Most Common Commands

```bash
# Start new work
my-context start "Your task" --project service-name

# Record decisions
my-context note "Decision: Why you chose this approach"

# Track files touched
my-context file src/file.go tests/file_test.go

# View current work
my-context show

# List recent work
my-context list

# See specific project
my-context list --project service-name

# Search for something
my-context list --search "keyword"

# Stop work
my-context stop

# Clean up old work
my-context archive "Old context name"

# Share with team
my-context export "Context name" --to file.md
```

---

## Tips

| Situation | Command |
|-----------|---------|
| I'm starting a task | `my-context start "Task"` |
| I made a decision | `my-context note "Decision: ..."` |
| I edited a file | `my-context file path/to/file` |
| I'm done working | `my-context stop` |
| Show my work | `my-context show` |
| Find my work | `my-context list --search "keyword"` |
| Organize by project | `--project service-name` |
| Share with team | `my-context export "name" --to file.md` |
| Clean old contexts | `my-context archive "name"` |

---

## Notes Format (Pro Tip)

Start notes with what they are:

```bash
my-context note "DECISION: Using Redis for caching"
my-context note "BUG: Found race condition in..."
my-context note "BLOCKED: Waiting for API docs"
my-context note "QUESTION: Is this approach correct?"
my-context note "REFACTOR: Extracted auth logic"
```

Makes them searchable later!

---

## First Week Workflow

```bash
# Monday morning
my-context start "Sprint 5 work" --project myservice
my-context note "Working on payment integration"

# As you work
my-context note "Using Stripe API v3"
my-context file internal/payments/stripe.go

# Mid-day check
my-context show

# When task done
my-context stop

# End of day
my-context list  # See what you did
```

---

## Common Questions

**Q: Can I have multiple contexts at once?**
A: No, only one active. Starting a new one stops the old one.

**Q: Do my notes get saved?**
A: Yes! Everywhere. Check `~/.my-context/`

**Q: How do I share my work?**
A: `my-context export "name" --to file.md` creates a markdown file.

**Q: What if I made a mistake?**
A: Delete it: `my-context delete "wrong-context"` (safe - asks first)

**Q: Where does my data live?**
A: `~/.my-context/` - all plain text files, human-readable.

---

## Learning Resources

| Time | Read This |
|------|-----------|
| 5 min | `docs/GETTING-STARTED.md` |
| 15 min | Above + your role guide in `docs/ROLE-SPECIFIC-GUIDES.md` |
| 60 min | `docs/PROGRESSIVE-GUIDE.md` - all 6 stages |
| Deep dive | `docs/TRIGGERS-TUTORIAL.md` - automation |
| Reference | `README.md` - all commands |

---

## Your Next Step

1. Copy the "Your First Context" commands above
2. Run them now
3. Read `docs/GETTING-STARTED.md`
4. Find your role in `docs/ROLE-SPECIFIC-GUIDES.md`
5. Use daily for a week

That's it! You're now a my-context user! 🚀

---

**Keep this card handy** - bookmark it or print it out for your first week.
