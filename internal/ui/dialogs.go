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

// viewProviderList 显示模型配置列表
func (m *Model) viewProviderList() string {
	header := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render("模型配置")

	listView := m.providerList.View()

	dialogWidth := min(60, m.width-10)

	// 错误消息显示
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(lipgloss.Color("#2C1810")).
			Padding(1, 2).
			MarginBottom(1).
			Align(lipgloss.Center).
			Width(dialogWidth - 4).
			Render("⚠ " + m.errMsg)
	}

	helpText := lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Align(lipgloss.Center).
		Render("↑/↓ j/k 选择 | n 新增 | e 编辑 | d 删除 | a 激活 | esc 返回")

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
				header,
				"",
				listView,
				errDisplay,
				"",
				helpText,
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

// viewProviderForm 模型配置表单（新增/编辑共用）
func (m *Model) viewProviderForm(isEdit bool) string {
	title := "添加配置"
	if isEdit {
		title = "编辑配置"
	}

	dialogWidth := min(60, m.width-10)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryColor).
		Width(dialogWidth - 10).
		Padding(0, 1)

	// 当前焦点的输入框高亮显示
	nameStyle := inputStyle
	baseURLStyle := inputStyle
	apiKeyStyle := inputStyle
	modelStyle := inputStyle

	if m.providerInputFocus == FocusProviderName {
		nameStyle = nameStyle.BorderForeground(SecondaryColor)
	} else if m.providerInputFocus == FocusProviderBaseURL {
		baseURLStyle = baseURLStyle.BorderForeground(SecondaryColor)
	} else if m.providerInputFocus == FocusProviderAPIKey {
		apiKeyStyle = apiKeyStyle.BorderForeground(SecondaryColor)
	} else if m.providerInputFocus == FocusProviderModel {
		modelStyle = modelStyle.BorderForeground(SecondaryColor)
	}

	// 错误消息显示
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(lipgloss.Color("#2C1810")).
			Padding(1, 2).
			MarginBottom(1).
			Width(dialogWidth - 4).
			Render("⚠ " + m.errMsg)
	}

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
				TitleStyle.Render(title),
				"",
				InfoStyle.Render("配置名称"),
				nameStyle.Render(m.providerNameInput.View()),
				"",
				InfoStyle.Render("Base URL"),
				baseURLStyle.Render(m.providerBaseURLInput.View()),
				"",
				InfoStyle.Render("API Key"),
				apiKeyStyle.Render(m.providerAPIKeyInput.View()),
				"",
				InfoStyle.Render("模型名称"),
				modelStyle.Render(m.providerModelInput.View()),
				errDisplay,
				"",
				HelpStyle.Render("[enter] 保存 · [tab] 切换 · [esc] 取消"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

// viewProviderDelete 确认删除配置
func (m *Model) viewProviderDelete() string {
	current := m.providerList.SelectedItem().(providerListItem)
	message := fmt.Sprintf("确认删除配置 '%s' ?", current.provider.Name)

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
				TitleStyle.Render("删除配置"),
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
