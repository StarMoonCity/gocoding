# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gocoding is a TUI project management tool built with Go and Bubble Tea. It manages projects, opens IDEs, and configures AI model providers.

## Commands

```bash
# Build the project
go build ./cmd/gocoding/

# Run the application
go run ./cmd/gocoding/

# Install to ~/go/bin (requires ~/go/bin in PATH)
./install.sh

# Run with provider config (skip to model provider list)
go run ./cmd/gocoding/ -p

# Add dependencies
go get github.com/charmbracelet/bubbletea@v1.1.0 github.com/charmbracelet/lipgloss@v0.13.0 github.com/charmbracelet/bubbles@v0.20.0 github.com/spf13/viper@v1.18.2

# Tidy dependencies
go mod tidy
```

## Data Storage

- Projects: `~/.config/gocoding/projects.json`
- Model Providers: `~/.config/gocoding/providers.json`

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
The `Model` in `internal/ui/model.go` uses a state machine pattern. Key states:

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

- `Project` struct: ID, Path, Alias, Description, CreatedAt, LastOpened, OpenCount
- `ProjectStore`: CRUD operations, Search, SortByLastOpened
- `ModelProvider` struct: ID, Name, BaseURL, APIKey, Model, Active
- `ModelProviderStore`: Provider management with active provider tracking
- `IDEExecutor` in `internal/commands/executor.go`: Opens projects in Claude/VSCode/OpenCode

### TUI Framework

- `bubbletea` - Main application loop
- `lipgloss` - Terminal styling
- `bubbles/list` - Project and provider lists
- `bubbles/textinput` - Form inputs
- `bubbles/textarea` - Multi-line text (descriptions)
- `bubbles/viewport` - Scrollable detail views

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
