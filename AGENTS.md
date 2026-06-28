# cmdX — Universal AI Agent Conventions

> This file is read by all AI coding tools (Claude, Cursor, Copilot, Windsurf).
> It defines universal conventions every contributor and AI agent must follow.

## Project Identity
- Name: cmdX
- Language: Go (92.4%)
- Repo: github.com/abhigyanwebber/cmdX
- Themes repo: github.com/abhigyanwebber/cmdX-themes
- Author: abhigyanwebber

## Universal Rules (apply regardless of AI tool)

### Code Style
- Follow standard Go project layout
- Use table-driven tests for all new packages
- Every exported function must have a Go doc comment
- Wrap errors with fmt.Errorf("context: %w", err)
- No global mutable state outside of cobra commands
- Prefer composition over inheritance
- Interface-driven design for testability

### File Rules
- NEVER commit: cmdx.exe, *.exe, frame*.png, .env files
- NEVER modify: go.sum directly
- ALWAYS run: go fmt ./... before committing
- ALWAYS run: go vet ./... before committing

### Architecture Rules
- Color resolution MUST use the centralized utility — never duplicate resolveColor()
- Shell injections MUST use InjectStart/InjectEnd markers from internal/shells/shell.go
- New asset types MUST implement the Asset interface in internal/assets/types.go
- New shell support MUST implement the Shell interface in internal/shells/shell.go

### Git Rules
- Always push with: git push origin HEAD:main
- Commit format: type(scope): description
  - feat(theme): add ocean-depths built-in theme
  - fix(assets): correct chafa path validation
  - test(config): add validator table-driven tests
  - docs(readme): update project structure

### Testing Rules
- Every new package needs a *_test.go file
- Use table-driven tests: var tests = []struct{ input, expected }{}
- Test both happy path and error cases
- Mock external calls (HTTP, file system) in tests

## What NOT to do
- Do not duplicate color resolution logic across packages
- Do not add new CLI commands without updating cmd/root.go help text
- Do not commit PNG frame files to the main repo
- Do not hardcode paths — use os.UserHomeDir() and filepath.Join()
- Do not break the composable asset system — assets must work standalone