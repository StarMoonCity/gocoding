package components

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
	"gocoding/internal/ui"
)

// ConfirmDialog 确认对话框
type ConfirmDialog struct {
	title     string
	message   string
	onConfirm func()
	onCancel  func()
	selected  int
	hoverIdx  int
}

// NewConfirmDialog 创建确认对话框
func NewConfirmDialog(title, message string, onConfirm, onCancel func()) *ConfirmDialog {
	return &ConfirmDialog{
		title:     title,
		message:   message,
		onConfirm: onConfirm,
		onCancel:  onCancel,
		selected:  1, // 默认选择"否"
	}
}

// Update 处理消息
func (d *ConfirmDialog) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			d.selected = 0
		case "right", "l":
			d.selected = 1
		case "enter", "y":
			if d.selected == 1 && d.onConfirm != nil {
				d.onConfirm()
			} else if d.selected == 0 && d.onCancel != nil {
				d.onCancel()
			}
			return nil, true
		case "n", "esc":
			if d.onCancel != nil {
				d.onCancel()
			}
			return nil, true
		}
		return nil, true
	}
	return nil, false
}

// View 渲染对话框
func (d *ConfirmDialog) View(width, height int) string {
	dialogWidth := min(50, max(35, int(float64(width)*0.6)))

	buttonWidth := 10

	// 确定按钮（右侧）- 霓虹红
	confirmStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(ui.ErrorColor).
		Background(lipgloss.Color("#1A0D10")).
		Padding(0, 2)
	if d.hoverIdx == 1 {
		confirmStyle = confirmStyle.Background(ui.ErrorColor).Foreground(ui.Background).Bold(true)
	}

	// 取消按钮（左侧）- 霓虹灰
	cancelStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(ui.SecondaryText).
		Background(ui.BackgroundLight).
		Padding(0, 2)
	if d.hoverIdx == 0 {
		cancelStyle = cancelStyle.Background(ui.BackgroundHover).Foreground(ui.Foreground)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.ErrorColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().
					Foreground(ui.ErrorColor).
					Bold(true).
					Render(d.title),
				"",
				lipgloss.NewStyle().
					Foreground(ui.Foreground).
					Render(d.message),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					cancelStyle.Render("[ N ] 否"),
					"  ",
					confirmStyle.Render("[ Y ] 是"),
				),
				"",
				lipgloss.NewStyle().
					Foreground(ui.MutedText).
					Render("此操作不可恢复"),
			),
		)

	return dialog
}

// Overlay 将对话框叠加到内容上
func (d *ConfirmDialog) Overlay(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, d.View(width, height))
}

// IDEOption IDE 选项
type IDEOption struct {
	Type        models.IDEType
	Name        string
	Description string
}

// IDEMenu IDE 菜单
type IDEMenu struct {
	title     string
	options   []IDEOption
	selected  int
	available map[models.IDEType]bool
	onSelect  func(ideType models.IDEType)
	onClose   func()
}

// NewIDEMenu 创建 IDE 菜单
func NewIDEMenu(onSelect func(ideType models.IDEType), onClose func()) *IDEMenu {
	return &IDEMenu{
		title: "选择 IDE",
		options: []IDEOption{
			{Type: models.IDEClaudeCode, Name: "Claude", Description: "Claude Code IDE"},
			{Type: models.IDEVSCode, Name: "VSCode", Description: "Visual Studio Code"},
			{Type: models.IDEOpenCode, Name: "OpenCode", Description: "OpenCode IDE"},
			{Type: models.IDECodexCLI, Name: "Codex", Description: "Codex CLI"},
		},
		available: make(map[models.IDEType]bool),
		onSelect:  onSelect,
		onClose:   onClose,
	}
}

// Update 处理消息
func (m *IDEMenu) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.selected = min(m.selected+1, len(m.options)-1)
		case "k", "up":
			m.selected = max(m.selected-1, 0)
		case "enter":
			if m.onSelect != nil {
				m.onSelect(m.options[m.selected].Type)
			}
			return nil, true
		case "1", "2", "3", "4":
			idx := int(msg.Runes[0] - '1')
			if idx < len(m.options) && m.onSelect != nil {
				m.onSelect(m.options[idx].Type)
			}
			return nil, true
		case "esc":
			if m.onClose != nil {
				m.onClose()
			}
			return nil, true
		}
		return nil, true
	}
	return nil, false
}

// ideColor 返回 IDE 对应的品牌色
func ideColorByType(ideType models.IDEType) lipgloss.Color {
	switch ideType {
	case models.IDEClaudeCode:
		return ui.IDEClaudeColor
	case models.IDEVSCode:
		return ui.IDEVSCodeColor
	case models.IDEOpenCode:
		return ui.IDEOpenCodeColor
	case models.IDECodexCLI:
		return ui.IDECodexColor
	default:
		return ui.PrimaryColor
	}
}

// View 渲染 IDE 菜单
func (m *IDEMenu) View(width, height int) string {
	var options []string
	for i, opt := range m.options {
		available := m.available[opt.Type]
		isSelected := i == m.selected

		ideClr := ideColorByType(opt.Type)

		var prefix string
		var nameStyle lipgloss.Style
		if isSelected {
			prefix = lipgloss.NewStyle().Foreground(ideClr).Render("▸ ")
			nameStyle = lipgloss.NewStyle().Foreground(ideClr).Bold(true)
		} else {
			prefix = "  "
			nameStyle = lipgloss.NewStyle().Foreground(ui.Foreground).Bold(true)
		}

		var statusIcon string
		if available {
			statusIcon = lipgloss.NewStyle().Foreground(ideClr).Render("●")
		} else {
			statusIcon = lipgloss.NewStyle().Foreground(ui.MutedText).Render("○")
		}

		colorBar := lipgloss.NewStyle().Foreground(ideClr).Render("▌")
		options = append(options,
			prefix+statusIcon+"  "+
				nameStyle.Render(opt.Name)+
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("  "+opt.Description)+
				" "+colorBar)
	}

	dialogWidth := min(45, max(35, int(float64(width)*0.5)))

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.AccentCyan).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.AccentCyan).Bold(true).Render("▣ "+m.title),
				"",
				lipgloss.JoinVertical(lipgloss.Left, options...),
				"",
				lipgloss.NewStyle().
					Foreground(ui.SecondaryText).
					Render("[↑↓] 选择  ·  [Enter] 打开  ·  [Esc] 返回"),
			),
		)

	return dialog
}

// Overlay 将菜单叠加到内容上
func (m *IDEMenu) Overlay(content string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, m.View(width, height))
}
