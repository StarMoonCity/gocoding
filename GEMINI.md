# Gocoding

## Project Overview
Gocoding is a TUI (Terminal User Interface) project management tool written in Go. It is designed to help developers manage their local projects, offering features like:
- **Project Tracking**: Maintains a list of projects with aliases and paths.
- **Quick Access**: Open projects in various IDEs (Claude Code, VSCode, OpenCode).
- **Statistics**: Tracks "last opened" times and open counts.
- **Visual Interface**: A rich terminal interface built with Bubble Tea, supporting adaptive layouts (single/double column) and various dialogs.

## Architecture

### Directory Structure
- **`cmd/gocoding/`**: Application entry point (`main.go`). Initializes config, store, and starts the TUI program.
- **`internal/ui/`**: Core UI logic using the Bubble Tea framework.
    - `main.go`: Defines the `Model`, `Update`, and `View` methods. Implements the state machine.
    - `styles.go`: (Inferred) Likely contains Lipgloss style definitions.
- **`internal/models/`**: Data structures.
    - `project.go`: Defines the `Project` struct (ID, Path, Alias, Stats).
- **`internal/store/`**: Data persistence.
    - `json_store.go`: Handles loading/saving projects to `~/.config/gocoding/projects.json`.
- **`internal/commands/`**: System interaction.
    - `executor.go`: Handles executing external IDE commands.
- **`internal/config/`**: Configuration management using Viper.

### Key Concepts
- **State Machine**: The UI transitions between states defined in `AppState` (e.g., `StateList`, `StateAddProject`, `StateIDEMenu`).
- **Data Storage**: Projects are persisted in a JSON file at `~/.config/gocoding/projects.json`.
- **IDE Integration**: The tool checks for available IDEs and attempts to launch them with the selected project path.

## Building and Running

### Prerequisites
- Go 1.20 or higher.

### Commands
- **Run**:
  ```bash
  go run ./cmd/gocoding/
  ```
- **Build**:
  ```bash
  go build ./cmd/gocoding/
  ```

### Dependencies
Key libraries used:
- `github.com/charmbracelet/bubbletea` (TUI runtime)
- `github.com/charmbracelet/bubbles` (UI components)
- `github.com/charmbracelet/lipgloss` (Styling)
- `github.com/spf13/viper` (Configuration)

## Development Conventions

- **UI Framework**: This project relies heavily on the [Bubble Tea](https://github.com/charmbracelet/bubbletea) Elm architecture (`Model`, `Update`, `View`). Changes to the UI should respect this pattern.
- **Styling**: Use `lipgloss` for all terminal styling.
- **Path Handling**: Be mindful of cross-platform path handling, though the current environment is Darwin.
- **Error Handling**: The application generally catches errors and displays them in the TUI or logs to stderr on fatal startup errors.
