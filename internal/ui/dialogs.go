package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewAddProject() string {
	dialogWidth := min(50, m.width-10)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryDim).
		Width(dialogWidth-6).
		Padding(0, 1)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryDim).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("＋ 添加项目"),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("项目路径"),
				inputStyle.Render(m.input.View()),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("项目名称"),
				inputStyle.Render(m.secondaryInput.View()),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 确认  ·  [Tab] 切换  ·  [Esc] 取消"),
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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("✎ 重命名项目"),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("新名称"),
				m.input.View(),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 确认  ·  [Esc] 取消"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}

func (m *Model) viewDeleteConfirm() string {
	current := m.safeGetSelectedProject()
	if current == nil {
		return ""
	}
	message := fmt.Sprintf("确认删除项目 '%s' ？", current.Alias)

	dialogWidth := min(50, m.width-10)

	// 使用 DoubleBorder 更严肃
	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ErrorColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().
					Foreground(ErrorColor).
					Bold(true).
					Render("✗ 删除项目"),
				"",
				lipgloss.NewStyle().
					Foreground(Foreground).
					Render(message),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					lipgloss.NewStyle().
						Foreground(ErrorColor).
						Background(lipgloss.Color("#2C1810")).
						Padding(0, 2).
						Render("[ Y ] 是"),
					"  ",
					lipgloss.NewStyle().
						Foreground(SecondaryText).
						Background(BackgroundLight).
						Padding(0, 2).
						Render("[ N ] 否"),
				),
				"",
				lipgloss.NewStyle().
					Foreground(MutedText).
					Render("此操作不可恢复"),
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
		var nameStyle lipgloss.Style
		if i == m.ideMenu.selected {
			prefix = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Render("▸ ")
			nameStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true)
		} else {
			prefix = "  "
			nameStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Bold(true)
		}
		var statusIcon string
		if available {
			statusIcon = lipgloss.NewStyle().Foreground(SuccessColor).Render("●")
		} else {
			statusIcon = lipgloss.NewStyle().Foreground(MutedText).Render("○")
		}
		options = append(options,
			prefix+statusIcon+"  "+
				nameStyle.Render(opt.Name)+
				lipgloss.NewStyle().Foreground(SecondaryText).Render("  "+opt.Description))
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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("▣ 选择 IDE"),
				"",
				lipgloss.JoinVertical(lipgloss.Left, options...),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[↑↓] 选择  ·  [Enter] 打开  ·  [Esc] 返回"),
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
		Render("▤ 项目详情")

	dialog := lipgloss.NewStyle().
		Width(m.width-10).
		Border(lipgloss.NormalBorder()).
		BorderForeground(SecondaryText).
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
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Esc] 返回  ·  [Q] 退出"),
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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("✎ 编辑描述"),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("描述（支持多行）"),
				m.ta.View(),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter/Ctrl+S] 保存  ·  [Esc] 取消"),
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
		Render("⚙ 模型配置")

	listView := m.providerList.View()

	dialogWidth := min(60, m.width-10)

	// 错误消息
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.Width(dialogWidth - 4).Render("✗ " + m.errMsg)
	}

	// 提示消息
	var tipDisplay string
	if m.tipMsg != "" {
		tipDisplay = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(lipgloss.Color("#0F1926")).
			Padding(0, 1).
			Render("ℹ " + m.tipMsg)
	}

	helpText := lipgloss.NewStyle().
		Foreground(SecondaryText).
		Align(lipgloss.Center).
		Width(dialogWidth - 4).
		Render(
			HelpKeyStyle.Render("[↑↓]"), " 选择 ",
			HelpKeyStyle.Render("[N]"), " 新增 ",
			HelpKeyStyle.Render("[E]"), " 编辑 ",
			HelpKeyStyle.Render("[D]"), " 删除 ",
			HelpKeyStyle.Render("[A]"), " 激活 ",
			HelpKeyStyle.Render("[Esc]"), " 退出",
		)

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
				tipDisplay,
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
	title := "＋ 添加配置"
	if isEdit {
		title = "✎ 编辑配置"
	}

	dialogWidth := min(60, m.width-10)

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryDim).
		Width(dialogWidth - 10).
		Padding(0, 1)

	// 当前焦点的输入框高亮 - 聚焦时使用主色
	nameStyle := inputStyle
	baseURLStyle := inputStyle
	apiKeyStyle := inputStyle
	modelStyle := inputStyle

	if m.providerInputFocus == FocusProviderName {
		nameStyle = nameStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderBaseURL {
		baseURLStyle = baseURLStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderAPIKey {
		apiKeyStyle = apiKeyStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderModel {
		modelStyle = modelStyle.BorderForeground(PrimaryColor)
	}

	// 错误消息
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.Width(dialogWidth - 4).Render("✗ " + m.errMsg)
	}

	// 提示消息
	var tipDisplay string
	if m.tipMsg != "" {
		tipDisplay = TipBoxStyle.Width(dialogWidth - 4).Render("ℹ " + m.tipMsg)
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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(title),
				"",
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.NewStyle().Foreground(SecondaryText).Render("配置名称"),
					nameStyle.Render(m.providerNameInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("Base URL"),
					baseURLStyle.Render(m.providerBaseURLInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("API Key"),
					apiKeyStyle.Render(m.providerAPIKeyInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("模型名称"),
					modelStyle.Render(m.providerModelInput.View()),
				),
				tipDisplay,
				errDisplay,
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 保存  ·  [Tab] 切换  ·  [Esc] 取消"),
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
	current := m.safeGetSelectedProvider()
	if current == nil {
		return ""
	}
	message := fmt.Sprintf("确认删除配置 '%s' ？", current.Name)

	dialogWidth := min(50, m.width-10)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ErrorColor).
		Background(Background).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Render("✗ 删除配置"),
				"",
				lipgloss.NewStyle().Foreground(Foreground).Render(message),
				"",
				lipgloss.JoinHorizontal(
					lipgloss.Center,
					lipgloss.NewStyle().
						Foreground(ErrorColor).
						Background(lipgloss.Color("#2C1810")).
						Padding(0, 2).
						Render("[ Y ] 是"),
					"  ",
					lipgloss.NewStyle().
						Foreground(SecondaryText).
						Background(BackgroundLight).
						Padding(0, 2).
						Render("[ N ] 否"),
				),
				"",
				lipgloss.NewStyle().Foreground(MutedText).Render("此操作不可恢复"),
			),
		)

	overlay := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialog,
	)

	return overlay
}
