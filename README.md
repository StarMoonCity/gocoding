# Gocoding

A TUI project management tool built with Go and Bubble Tea.

## Features

- Project management with aliases and descriptions
- Open projects in Claude Code, VSCode, OpenCode, or Codex CLI
- AI model provider configuration
- Debug mode for development

## Usage

```bash
# Run the application
gocoding

# Run with debug mode
gocoding -d

# Open provider config directly
gocoding -p
```

## Keybindings

- `↑/↓` or `j/k` - Navigate
- `n` - Add new project
- `r` - Rename project
- `d` - Delete project
- `v` - View details
- `e` - Edit description
- `/` - Search
- `p` - Open provider config
- `q` - Quit

## Configuration

Projects are stored in `~/.config/gocoding/projects.json`
Providers are stored in `~/.config/gocoding/providers.json`
