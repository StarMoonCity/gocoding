# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gocoding is a TUI project management tool built with Go and Bubble Tea. It manages projects, opens IDEs, and configures AI model providers.

**Data Storage**:
- Projects: `~/.config/gocoding/projects.json`
- Model Providers: `~/.config/gocoding/providers.json`

## Commands

```bash
# Build the project
go build ./cmd/gocoding/

# Run the application
go run ./cmd/gocoding/

# Run with provider config (skip to model provider list)
go run ./cmd/gocoding/ -p

# Install to ~/go/bin (requires ~/go/bin in PATH)
./install.sh

# Add dependencies
go get github.com/charmbracelet/bubbletea@v1.1.0 github.com/charmbracelet/lipgloss@v0.13.0 github.com/charmbracelet/bubbles@v0.20.0 github.com/spf13/viper@v1.18.2

# Tidy dependencies
go mod tidy

# Run tests
go test ./...

# Run single test
go test -v ./internal/models/...

# Lint (requires golangci-lint)
golangci-lint run ./...
```

## Architecture

### Dependency Flow
```
cmd/gocoding/main.go
    └── internal/ui/model.go (Bubble Tea Model)
        ├── internal/models/project.go (Project/ProjectStore)
        ├── internal/models/model_provider.go (ModelProvider/ModelProviderStore)
        ├── internal/commands/executor.go (IDE execution)
        ├── internal/config/config.go (Viper config)
        └── internal/store/json_store.go (JSON persistence)
```

### UI State Machine

The `Model` in `internal/ui/model.go` uses a state machine pattern:

**Project Management:**
- `StateList` - Main project list view
- `StateAddProject` - Add new project (path auto-fills name)
- `StateRenameProject` - Rename existing project
- `StateDeleteConfirm` - Delete confirmation dialog
- `StateIDEMenu` - IDE selection (Claude/VSCode/OpenCode)
- `StateViewDetail` - View project details
- `StateEditDescription` - Edit project description
- `StateSearch` - Search/filter projects

**Model Provider Management:**
- `StateProviderList` - Provider configuration list
- `StateProviderAdd` - Add new provider
- `StateProviderEdit` - Edit provider
- `StateProviderDelete` - Delete confirmation

### Key Types

- `Project`: ID, Path, Alias, Description, CreatedAt, LastOpened, OpenCount
- `ProjectStore`: CRUD operations, Search, SortByLastOpened
- `ModelProvider`: ID, Name, BaseURL, APIKey, Model, Active
- `ModelProviderStore`: Provider management with active provider tracking
- `IDEExecutor` in `internal/commands/executor.go`: Opens projects in Claude/VSCode/OpenCode

### TUI Framework

- `bubbletea` - Main application loop
- `lipgloss` - Terminal styling
- `bubbles/list` - Project and provider lists
- `bubbles/textinput` - Form inputs
- `bubbles/textarea` - Multi-line text (descriptions)
- `bubbles/viewport` - Scrollable detail views

## Code Style

### Import Grouping

Group imports in this order:
1. Standard library
2. External dependencies
3. Internal packages (relative imports)

### Error Handling

- Always check errors explicitly: `if err != nil { ... }`
- Return meaningful error messages with context (Chinese for user-facing errors)
- Handle `os.IsNotExist(err)` specifically for file-not-found cases
- Use `fmt.Errorf("context: %w", err)` for wrapping errors

### JSON Fields

Use JSON tags for all exported struct fields with snake_case keys:
```go
type Project struct {
    ID          string    `json:"id"`
    LastOpened  time.Time `json:"last_opened"`
}
```

### Comments

- Comment exported types, functions, and constants
- Use Chinese comments for internal implementation details
- Use English comments for exported/public APIs

### Dependency Injection

The codebase uses constructor injection:
```go
func NewJSONStore(store *models.ProjectStore) *JSONStore { ... }
func NewModel(store *models.ProjectStore) *Model { ... }
```

## Self-Diagnosing Workflow

### Rules
1. After editing code: run `go build ./cmd/gocoding/` to verify compilation
2. If compilation fails: analyze error, auto-fix obvious issues, retry
3. Only report completion after compilation passes
4. For repeated errors: explain tradeoffs between solutions

### Diagnostic Commands
```bash
go vet ./cmd/gocoding/ ./internal/...  # Fast syntax check
go build ./cmd/gocoding/               # Full compilation
go test ./...                         # Run tests
```

## Development Workflow

1. **Syntax checking**: Use `gopls` for Go syntax verification
2. **Build before commit**: Always run `go build ./cmd/gocoding/` before committing
3. **Incremental changes**: Make small, focused changes

## Git Workflow

- After completing git commits, run a quick verification to confirm the commit was successful (e.g., `git log -1 --oneline`)

## Refactoring

- Before refactoring, create a checklist of affected files and verify changes compile before committing

## UI Development (Debugging)

- When fixing UI-related issues, test incrementally after each change to catch layout/alignment problems early

## LSP Support

Use LSP features for code navigation and analysis:
- `documentSymbol` - List all symbols in a file
- `goToDefinition` - Jump to symbol definition
- `findReferences` - Find symbol references
- `hover` - Get symbol documentation
- `workspaceSymbol` - Search symbols across the codebase
- `incomingCalls` / `outgoingCalls` - Analyze call hierarchies

Example: When exploring unfamiliar code, use `documentSymbol` to understand the structure, then `goToDefinition` to navigate to implementations.

## Key Constraints

- **NEVER** suppress type errors with `as any`, `@ts-ignore`, or similar
- **NEVER** commit without explicit user request
- **NEVER** leave code in broken state after failures
