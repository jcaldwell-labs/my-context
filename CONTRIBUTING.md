# Contributing to My-Context

Thank you for considering contributing to My-Context! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue on GitHub with:
- A clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, Go version, etc.)

### Suggesting Features

Feature suggestions are welcome! Please open an issue with:
- A clear description of the feature
- Use cases and benefits
- Any implementation ideas (optional)

### Pull Requests

1. **Fork the repository** and create a feature branch
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards:
   - Write clean, readable Go code
   - Follow existing code style and patterns
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   go test ./...
   go build ./cmd/my-context/
   ```

4. **Commit your changes** with clear commit messages:
   ```bash
   git commit -m "feat: add new feature description"
   ```

   Use conventional commit prefixes:
   - `feat:` - New features
   - `fix:` - Bug fixes
   - `docs:` - Documentation changes
   - `test:` - Test additions or changes
   - `refactor:` - Code refactoring
   - `chore:` - Build/tooling changes

5. **Push to your fork** and create a pull request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Building

```bash
# Clone your fork
git clone https://github.com/YOUR-USERNAME/my-context.git
cd my-context

# Build
go build -o my-context ./cmd/my-context/

# Run tests
go test ./...
```

### Testing

- Write unit tests for new functions
- Write integration tests for new commands
- Ensure all tests pass before submitting PR
- Test on multiple platforms when possible (Windows, Linux, macOS)

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Keep functions focused and single-purpose
- Add comments for exported functions
- Use descriptive variable names

## Project Structure

```
cmd/my-context/          # CLI entry point
internal/
  commands/              # Command implementations
  core/                  # Business logic
  models/                # Data structures
  output/                # Output formatters
tests/
  integration/           # Integration tests
  unit/                  # Unit tests
```

## Questions?

Feel free to open an issue for questions or discussions about contributing.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
