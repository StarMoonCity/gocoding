# AGENTS.md

This file provides guidance for agentic coding agents operating in this repository.

## Project Overview

Gocoding is a TUI (Terminal User Interface) project management tool built with Go and Bubble Tea. It manages projects, opens IDEs, and records recently opened projects.

**Data Storage**: `~/.config/gocoding/projects.json`

## Build Commands

```bash
# Build the application
go build ./cmd/gocoding/

# Run the application
go run ./cmd/gocoding/

# Add dependencies
go get github.com/charmbracelet/bubbletea@v1.1.0 github.com/charmbracelet/lipgloss@v0.13.0 github.com/charmbracelet/bubbles@v0.20.0 github.com/spf13/viper@v1.18.2

# Run tests (if any exist)
go test ./...

# Run single test file
go test -v ./internal/models/...

# Check for linting issues (if golangci-lint is installed)
golangci-lint run ./...
```

## Code Style Guidelines

### Imports

Group imports in this order:
1. Standard library
2. External dependencies
3. Internal packages (relative imports)

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"

    "github.com/charmbracelet/bubbletea"
    "github.com/spf13/viper"

    "gocoding/internal/config"
    "gocoding/internal/models"
)
```

### Naming Conventions

- **Packages**: lowercase, short, descriptive (e.g., `ui`, `models`, `commands`)
- **Variables/Functions**: camelCase for unexported, PascalCase for exported
- **Constants**: PascalCase for exported, camelCase for unexported; use descriptive names
- **Types**: PascalCase, singular nouns (e.g., `Project`, `ProjectStore`)
- **Interfaces**: PascalCase, often ending with `-er` (e.g., `IDEExecutor`)
- **Type aliases**: PascalCase (e.g., `type IDEExecutor = commands.IDEExecutor`)

### Error Handling

- Always check errors explicitly: `if err != nil { ... }`
- Never use `_` to ignore errors unless explicitly intentional
- Return meaningful error messages with context (Chinese for user-facing errors)
- Handle `os.IsNotExist(err)` specifically for file-not-found cases
- Use `fmt.Errorf("context: %w", err)` for wrapping errors

```go
func (s *ProjectStore) Load(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // Silent handling for non-existent files
        }
        return err
    }
    return json.Unmarshal(data, s)
}
```

### Structs and JSON

- Use JSON tags for all exported struct fields
- Use snake_case for JSON keys (e.g., `json:"last_opened"`)

```go
type Project struct {
    ID         string    `json:"id"`
    Path       string    `json:"path"`
    Alias      string    `json:"alias"`
    Description string   `json:"description"`
    CreatedAt  time.Time `json:"created_at"`
    LastOpened time.Time `json:"last_opened"`
    OpenCount  int       `json:"open_count"`
}
```

### Comments

- Comment exported types, functions, and constants
- Use Chinese comments for internal implementation details
- Use English comments for exported/public APIs
- Place comments on the line above the code they describe

### File Organization

Follow the established architecture:
```
cmd/gocoding/main.go           # Entry point
internal/ui/main.go            # Bubble Tea Model (main UI logic)
internal/models/project.go     # Data structures
internal/commands/executor.go  # IDE execution
internal/config/config.go      # Viper configuration
internal/store/json_store.go   # JSON persistence
internal/ui/{layout,view,styles,utils,dialogs,update}.go # UI components
```

### UI/Styling Code

For Bubble Tea UI components:
- Use `lipgloss` for terminal styling
- Define colors as package-level variables in `styles.go`
- Define style constants as package-level variables
- Delegate all visual/styling changes to `frontend-ui-ux-engineer` agent

### Testing

- Write tests for model and store operations
- Use table-driven tests where appropriate
- Test file naming: `<file>_test.go`
- Run specific tests: `go test -v ./internal/models/...`

### Development Workflow

1. **Syntax checking**: Use `gopls-lsp` (`lsp_diagnostics`) for Go syntax verification
2. **Build before commit**: Always run `go build ./cmd/gocoding/` before committing
3. **Incremental changes**: Make small, focused changes
4. **Verify changes**: Run `lsp_diagnostics` on modified files

### Key Constraints

- **NEVER** suppress type errors with `as any`, `@ts-ignore`, `@ts-expect-error`
- **NEVER** commit without explicit user request
- **NEVER** leave code in broken state after failures
- **DELEGATE** frontend visual/styling changes to `frontend-ui-ux-engineer`
- **CONSULT** `oracle` for architecture decisions or after 2+ failed fix attempts

### Dependency Injection Pattern

The codebase uses simple constructor injection:

```go
func NewJSONStore(store *models.ProjectStore) *JSONStore { ... }
func NewModel(store *models.ProjectStore) *Model { ... }
```

Follow this pattern for new components that depend on stores or services.
