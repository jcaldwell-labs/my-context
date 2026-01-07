# Self-Hosted GitHub Actions Runner Setup

This guide explains how to set up a self-hosted GitHub Actions runner for automatic deployment when merging to main.

## Overview

When code is merged to `main`, the CI workflow:

1. Runs lint, test, and build jobs on GitHub-hosted runners
2. If all pass, triggers `deploy-local` job on the self-hosted runner
3. The runner executes `scripts/deploy-local.sh` which builds and installs for all users

## Prerequisites

- Linux machine (WSL2 works)
- Go 1.23+ installed
- Write access to all user `~/.local/bin` directories
- Network access to GitHub

## One-Time User Setup

Each user needs `~/.local/bin` directory created. Run as each user or with sudo:

```bash
# For each user (godev, be-dev-agent, cdev)
mkdir -p ~/.local/bin
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc
```

Or as root:

```bash
for user in godev be-dev-agent cdev; do
    sudo -u $user mkdir -p /home/$user/.local/bin
done
```

## Setup Instructions

### 1. Create Runner User (Recommended)

Create a dedicated user with permissions to write to all user install directories:

```bash
# Option A: Add godev to a shared group
sudo groupadd localbin
sudo usermod -aG localbin godev
sudo usermod -aG localbin be-dev-agent
sudo usermod -aG localbin cdev

# Set group ownership on install directories
sudo chgrp localbin /home/godev/.local/bin
sudo chgrp localbin /home/be-dev-agent/.local/bin
sudo chgrp localbin /home/cdev/.local/bin
sudo chmod g+w /home/*/.local/bin
```

### 2. Download and Configure Runner

Navigate to repository Settings > Actions > Runners > New self-hosted runner.

```bash
# Create directory for runner
mkdir -p ~/actions-runner && cd ~/actions-runner

# Download runner (check GitHub for latest version)
curl -o actions-runner-linux-x64-2.321.0.tar.gz -L \
  https://github.com/actions/runner/releases/download/v2.321.0/actions-runner-linux-x64-2.321.0.tar.gz

# Extract
tar xzf ./actions-runner-linux-x64-2.321.0.tar.gz

# Configure (get token from GitHub UI)
./config.sh --url https://github.com/jcaldwell-labs/my-context --token YOUR_TOKEN

# Labels (optional but recommended)
# When prompted for labels, add: wsl,linux,my-context
```

### 3. Install as Service

```bash
# Install service
sudo ./svc.sh install

# Start service
sudo ./svc.sh start

# Check status
sudo ./svc.sh status
```

### 4. Verify Runner

After setup, the runner should appear in:

- Repository > Settings > Actions > Runners

Status should show "Idle" (ready to receive jobs).

## Testing the Pipeline

1. Create a test branch and make a small change
2. Open PR, ensure CI passes
3. Merge to main
4. Watch Actions tab for `deploy-local` job
5. Verify installation: `my-context --version` (for each user)

## Troubleshooting

### Runner not picking up jobs

```bash
# Check service status
sudo ./svc.sh status

# View logs
journalctl -u actions.runner.jcaldwell-labs-my-context.$(hostname).service -f
```

### Permission denied during deploy

Ensure the runner user has write access:

```bash
# Test write access
touch /home/godev/.local/bin/test && rm /home/godev/.local/bin/test
touch /home/be-dev-agent/.local/bin/test && rm /home/be-dev-agent/.local/bin/test
touch /home/cdev/.local/bin/test && rm /home/cdev/.local/bin/test
```

### Go not found

Ensure Go is in the runner's PATH:

```bash
# Add to runner's environment
echo 'PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/actions-runner/.env
```

## Security Considerations

- Self-hosted runners execute code from the repository
- Only use with trusted repositories
- The runner has write access to user directories
- Consider network isolation if needed

## Workflow Configuration

The deploy job in `.github/workflows/ci.yml`:

```yaml
deploy-local:
  name: Deploy to Local Environment
  runs-on: self-hosted
  needs: [lint, test, build]
  if: github.event_name == 'push' && github.ref == 'refs/heads/main'
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: "1.23"
    - run: ./scripts/deploy-local.sh
```

## Manual Deployment

If the runner is unavailable, deploy manually:

```bash
cd ~/projects/go/my-context
git pull origin main
./scripts/deploy-local.sh
```
