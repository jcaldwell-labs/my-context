# GitHub Account Migration Record

**Date:** 2026-02-02
**Decision Type:** Infrastructure / Identity Management

## Summary

Migrated primary development identity from `jcaldwell1066` (accesso email) to `jcaldwell` (jeffery.caldwell@gmail.com) as the main GitHub account for all future development work.

## Accounts

| Account | Email | Role | Status |
|---------|-------|------|--------|
| `jcaldwell` | jeffery.caldwell@gmail.com | **Primary** | Active |
| `jcaldwell1066` | (accesso email) | Legacy | Retained for transition |

## Configuration Changes

### SSH Keys

- **New key generated:** `~/.ssh/id_ed25519_jcaldwell`
- **Registered to:** `jcaldwell` GitHub account
- **Fingerprint:** `SHA256:gCAjoJx/gr2mM+yID+f3rTSOtAwjIJxjXA7qmr7IjaQ`

### SSH Config (`~/.ssh/config`)

```
# GitHub - jcaldwell (primary account)
Host github.com
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_jcaldwell

# GitHub - jcaldwell1066 (legacy account)
Host github.com-jcaldwell1066
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519
```

### Git Global Config

```
user.name = Jeff Caldwell
user.email = jeffery.caldwell@gmail.com
```

### GitHub CLI (`gh`)

- Both accounts configured
- `jcaldwell` set as active account
- Switch available via: `gh auth switch`

## Organization Access

| Organization | jcaldwell | jcaldwell1066 |
|--------------|-----------|---------------|
| `jcaldwell-labs` | Owner | Owner |
| `jaxclipse-solutions` | Member | - |

## Token Scopes

**jcaldwell:**
- `gist`
- `read:org`
- `repo`

**jcaldwell1066:**
- `admin:public_key`
- `admin:org`
- `gist`
- `read:org`
- `repo`

## Migration Notes

1. The `jcaldwell` account predates `jcaldwell1066` (created March 2008 vs Feb 2017)
2. `jcaldwell1066` was created for work-related purposes (accesso)
3. `jcaldwell-labs` org was originally created under `jcaldwell1066`
4. `jcaldwell` was invited and promoted to owner of `jcaldwell-labs`

## Legacy Access Pattern

To clone/push to repos that require `jcaldwell1066` credentials:

```bash
git clone git@github.com-jcaldwell1066:owner/repo.git
```

## Future Cleanup

Once migration is complete:
- [ ] Remove `jcaldwell1066` from `jcaldwell-labs` org (optional)
- [ ] Transfer any remaining repos to `jcaldwell`
- [ ] Archive or delete `jcaldwell1066` account
