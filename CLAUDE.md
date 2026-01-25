# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gocoding is a TUI (Terminal User Interface) project management tool built with Go and Bubble Tea. It manages projects, opens IDEs, and records recently opened projects.

**Data Storage**: `~/.config/gocoding/projects.json`

## Commands

```bash
# Build the project
go build ./cmd/gocoding/

# Run the application
go run ./cmd/gocoding/

# Add dependencies
go get github.com/charmbracelet/bubbletea@v1.1.0 github.com/charmbracelet/lipgloss@v0.13.0 github.com/charmbracelet/bubbles@v0.20.0 github.com/spf13/viper@v1.18.2
```

## Architecture

### Dependency Flow
```
cmd/gocoding/main.go
    └── internal/ui/main.go (Bubble Tea Model)
        ├── internal/models/project.go (Project data structures)
        ├── internal/commands/executor.go (IDE execution)
        ├── internal/config/config.go (Viper config)
        └── internal/store/json_store.go (JSON persistence)
```

### UI State Machine
The main `Model` in `internal/ui/main.go` uses a state machine pattern:
- `StateList` - Main project list view
- `StateAddProject` - Enter project path and name (auto-fills name from path)
- `StateRenameProject` - Rename existing project
- `StateDeleteConfirm` - Delete confirmation dialog
- `StateIDEMenu` - IDE selection menu (Claude/VSCode/OpenCode)

### Key Types
- `Project` struct in `internal/models/project.go` contains: ID, Path, Alias, CreatedAt, LastOpened, OpenCount
- `IDEExecutor` in `internal/commands/executor.go` handles opening projects in different IDEs
- `ProjectStore` manages CRUD operations for projects

### TUI Framework
- Uses `bubbletea` for the main application loop
- Uses `lipgloss` for terminal styling
- Uses `bubbles/list` for project list and `bubbles/textinput` for input fields

# 开发
- 写代码时，检查语法时，优先使用gopls-lsp
