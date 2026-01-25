package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewAddProject() string {
	dialogWidth := min(50, m.width-10)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryColor).
		Width(dialogWidth-6).
		Padding(0, 1)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				TitleStyle.Render("添加项目"),
				"",
				InfoStyle.Render("项目路径"),
				inputStyle.Render(m.input.View()),
				"",
				InfoStyle.Render("项目名称"),
				inputStyle.Render(m.secondaryInput.View()),
				"",
				HelpStyle.Render("[enter] 确认 · [tab] 切换 · [esc] 取消"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewRenameProject() string {
	dialogWidth := min(50, m.width-10)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				TitleStyle.Render("重命名项目"),
				"",
				InfoStyle.Render("新名称"),
				m.input.View(),
				"",
				HelpStyle.Render("[enter] 确认 · [esc] 取消"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewDeleteConfirm() string {
	current := m.list.SelectedItem().(listItem)
	message := fmt.Sprintf("确认删除项目 '%s' ?", current.project.Alias)

	dialogWidth := min(50, m.width-10)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ErrorColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				TitleStyle.Render("删除项目"),
				"",
				lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Render(message),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					GetStatusStyle(true).Render("[y] 是"),
					"  ",
					GetStatusStyle(true).Render("[n] 否"),
				),
				"",
				HelpStyle.Render("此操作不可恢复"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewIDEMenu() string {
	var options []string
	for i, opt := range m.ideMenu.options {
		available := m.ideMenu.available[opt.Type]
		var prefix string
		if i == m.ideMenu.selected {
			prefix = IDESelectedStyle.Render(" ● ")
		} else {
			prefix = "   "
		}
		availability := GetStatusStyle(available).Render("✓ ")
		if !available {
			availability = GetStatusStyle(false).Render("✗ ")
		}
		options = append(options, prefix+availability+lipgloss.NewStyle().Bold(true).Render(opt.Name)+"  "+opt.Description)
	}

	dialogWidth := min(40, m.width-10)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				TitleStyle.Render("选择 IDE"),
				"",
				lipgloss.JoinVertical(lipgloss.Left, options...),
				"",
				HelpStyle.Render("[enter] 打开 · [esc] 返回"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewViewDetail() string {
	m.updateViewport()

	header := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render("项目详情")

	dialog := lipgloss.NewStyle().
		Width(m.width-10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SecondaryColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				header,
				"",
				m.viewport.View(),
				"",
				HelpStyle.Render("[esc] 返回 · [q] 退出"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewEditDescription() string {
	dialogWidth := min(60, m.width-10)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				TitleStyle.Render("编辑项目描述"),
				"",
				InfoStyle.Render("描述 (支持多行)"),
				m.ta.View(),
				"",
				HelpStyle.Render("[enter/ctrl+s] 保存 · [esc] 取消"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}
