package components

import (
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
	"gocoding/internal/ui"
)

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
}

// ideColorByType 返回 IDE 对应的品牌色
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
