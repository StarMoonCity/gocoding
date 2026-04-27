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

## Snapshot Testing (Regression Testing)

The project uses Golden File snapshot testing for UI regression prevention.

### Test Files
- `internal/ui/snapshot_test.go` - Main snapshot test file
- `internal/ui/testdata/*.golden` - Golden file snapshots

### Run Tests
```bash
# Run all UI tests
go test -v ./internal/ui/...

# Run snapshot tests only
go test -v ./internal/ui/... -run TestViewSnapshot

# Update snapshots (after intentional UI changes)
go test -v ./internal/ui/... -run TestViewSnapshot -update
```

### Test Coverage
- **Views**: List views (empty, with projects), add project form, provider list with tips
- **Tip/Error messages**: Activation feedback, error display
- **Layout**: Help text alignment, dialog borders

### Important Rules
1. **Do NOT use `-update` flag casually** - it overwrites golden files, hiding real regressions
2. **Test does NOT modify config files** - activation tests simulate state without calling `WriteToClaudeSettings`
3. **Fixed terminal size** - tests use 80x24 to ensure consistent rendering
4. **Dynamic content normalization** - timestamps and counts are replaced with placeholders

### Adding New Snapshot Tests
```go
{
    name: "new_view",
    setup: func(m *Model) {
        m.state = StateXXX
        // setup state...
    },
    golden: "new_view",
},
```

### Notes on lipgloss.Place
- `lipgloss.Place` alignment can behave differently in tests vs real terminal
- Prefer standard layout composition over Place for reliability
- Always run tests after changes to catch regressions

## Problem Analysis Framework

遇到 Bug/显示问题时，按此顺序排查：

### 显示问题（裁剪 vs 截断）

1. **先确认症状类型**：
   - 裁剪（内容被切掉）→ 尺寸/布局问题
   - 截断（内容被…替代）→ 渲染/样式问题

2. **TUI 问题排查顺序**：
   - 检查 `SetSize()` → `calculateLayout()` → `SetWidth/SetHeight` 尺寸链
   - 检查布局配置是否考虑了所有元素（header、helpText、padding 等）
   - 最后才检查 Render 函数的样式逻辑

3. **布局计算要点**：
   - 对话框宽度 = 内容宽度 + padding + border
   - 对话框高度 = 所有元素高度之和 + padding + border
   - 列表高度 = 对话框高度 - 其他元素高度

### 代码问题（逻辑错误）

1. **先重现问题**：确认复现步骤
2. **隔离范围**：缩小到具体函数/文件
3. **假设验证**：列出可能原因，逐一排除
4. **改动最小化**：只改必要的地方

### 分析习惯

- 不要急于修改代码，先说出分析计划
- 列出 2-3 种可能的原因
- 从最可能的开始排查
- 修改后验证是否解决，确认没有引入新问题

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
