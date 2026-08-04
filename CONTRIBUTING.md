# Contributing to CasaOS Fork

Thank you for your interest in contributing!

## Development Setup

### Backend (Go)
```bash
# Requires Go 1.25+
export PATH=$HOME/go-sdk/go/bin:$HOME/go/bin:$PATH
export GOPATH=$HOME/go

cd CasaOS
go build ./...
go test ./...
```

### Frontend (Vue)
```bash
cd CasaOS-UI
yarn install
yarn dev
```

## Code Style

### Go
- Follow standard Go formatting (`gofmt`)
- Use `go vet` to check for issues
- No `fmt.Errorf(variable)` - use `fmt.Errorf("%s", variable)` instead
- Sanitize all file paths with `file.SanitizePath()` before use
- Return errors instead of calling `os.Exit()` or `panic()`

### Vue/JavaScript
- Use `$t()` for all user-facing strings (i18n)
- Use `$store.state` for global state
- Handle errors gracefully with user-facing messages
- Clean up event listeners in `beforeDestroy()`

## Pull Request Process

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes
4. Run tests: `go test ./...` (backend) or `yarn test` (frontend)
5. Commit with descriptive message
6. Push and create a PR

## Commit Messages

Use conventional commits:
- `feat:` new feature
- `fix:` bug fix
- `security:` security improvement
- `refactor:` code restructuring
- `docs:` documentation changes

## Security

See [SECURITY.md](SECURITY.md) for security-related contributions.
