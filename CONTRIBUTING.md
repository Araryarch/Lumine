# Contributing to Lumine

Thank you for your interest in contributing to Lumine! 🌟

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/Araryarch/lumine.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit: `git commit -m "Add your feature"`
7. Push: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

### Prerequisites
- Go 1.21+
- Docker 20.10+
- Git

### Build from Source
```bash
git clone https://github.com/Araryarch/lumine.git
cd lumine
go mod download
go build -o lumine
./lumine
```

### Project Structure
```
lumine/
├── main.go                 # Entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── docker/            # Docker integration
│   ├── domain/            # Domain management
│   ├── project/           # Project creation & management
│   ├── runtime/           # Runtime version management
│   ├── ui/                # TUI components
│   ├── compose/           # Docker Compose generation
│   └── nginx/             # Nginx config generation
├── .github/
│   └── workflows/         # CI/CD pipelines
└── docs/                  # Documentation
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Run `go vet` before committing
- Add comments for exported functions
- Keep functions small and focused

## Adding New Features

### Adding a New Framework

1. Add framework to `internal/project/templates.go`:
```go
func (m *Manager) createYourFrameworkProject(ctx context.Context, name, path string) error {
    // Implementation
}
```

2. Add to switch case in `internal/project/manager.go`:
```go
case "yourframework":
    return m.createYourFrameworkProject(ctx, name, path)
```

3. Add to UI in `internal/ui/create_project.go`:
```go
{"YourFramework", "Description", "Runtime"},
```

4. Add badge style in `internal/ui/projects_panel.go`

### Adding a New Runtime

1. Add to `internal/config/config.go`:
```go
type Runtimes struct {
    // ... existing
    YourRuntime string `yaml:"yourruntime"`
}
```

2. Add to default config initialization

3. Add to UI panels

## Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./internal/package
```

## Pull Request Guidelines

- Keep PRs focused on a single feature/fix
- Update documentation if needed
- Add tests for new features
- Ensure all tests pass
- Update CHANGELOG.md
- Follow commit message conventions

### Commit Message Format
```
type(scope): subject

body (optional)

footer (optional)
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:
```
feat(ui): add Rust project support
fix(docker): resolve container startup issue
docs(readme): update installation instructions
```

## Reporting Bugs

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md) and include:
- OS and version
- Lumine version
- Docker version
- Steps to reproduce
- Expected vs actual behavior
- Logs/screenshots

## Feature Requests

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md) and describe:
- The problem you're trying to solve
- Your proposed solution
- Alternative solutions considered
- Any additional context

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Help others learn and grow

## Questions?

- Open a [Discussion](https://github.com/Araryarch/lumine/discussions)
- Join our community chat (if available)
- Check existing issues and PRs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Lumine! 🚀
