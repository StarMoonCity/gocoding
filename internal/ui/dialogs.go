package ui

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewAddProject() string {
	// 动态计算对话框宽度：使用屏幕宽度的 70%，最小 40，最大 50
	dialogWidth := min(50, max(40, int(float64(m.width)*0.7)))

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

	return dialog
}

func (m *Model) viewRenameProject() string {
	// 动态计算对话框宽度
	dialogWidth := min(50, max(40, int(float64(m.width)*0.7)))

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryDim).
		Width(dialogWidth-6).
		Padding(0, 1)

	// 当前焦点的输入框高亮
	pathStyle := inputStyle
	nameStyle := inputStyle
	if m.inputFocus == FocusPath {
		pathStyle = pathStyle.BorderForeground(PrimaryColor)
	} else if m.inputFocus == FocusName {
		nameStyle = nameStyle.BorderForeground(PrimaryColor)
	}

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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("✎ 编辑项目"),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("项目路径"),
				pathStyle.Render(m.input.View()),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("项目名称"),
				nameStyle.Render(m.secondaryInput.View()),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 确认  ·  [Tab] 切换  ·  [Esc] 取消"),
			),
		)

	return dialog
}

func (m *Model) viewDeleteConfirm() string {
	current := m.safeGetSelectedProject()
	if current == nil {
		return ""
	}
	message := fmt.Sprintf("确认删除项目 '%s' ？", current.Alias)

	// 动态计算对话框宽度
	dialogWidth := min(50, max(35, int(float64(m.width)*0.6)))

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

	return dialog
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

	dialogWidth := min(40, max(30, int(float64(m.width)*0.5)))

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

	return dialog
}

func (m *Model) viewViewDetail() string {
	m.updateViewport()

	header := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render("▤ 项目详情")

	// 动态计算详情对话框宽度
	dialogWidth := max(50, int(float64(m.width)*0.8))

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
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

	return dialog
}

func (m *Model) viewEditDescription() string {
	// 动态计算对话框宽度
	dialogWidth := min(60, max(45, int(float64(m.width)*0.75)))

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

	return dialog
}

// viewProviderList 显示模型配置列表
func (m *Model) viewProviderList() string {
	header := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render("⚙ 模型配置")

	dialogWidth := m.providerDialogWidth()
	contentWidth := m.providerListWidth()
	listView := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(m.providerList.View())

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
		Width(contentWidth).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				m.renderProviderHelpItem("[k↑/j↓]", "选择"),
				"  ",
				m.renderProviderHelpItem("[N]", "新增"),
				"  ",
				m.renderProviderHelpItem("[E]", "编辑"),
				"  ",
				m.renderProviderHelpItem("[D]", "删除"),
				"  ",
				m.renderProviderHelpItem("[A]", "激活"),
			),
			m.renderProviderHelpItem("[Esc]", "退出"),
		))

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

	return dialog
}

// viewProviderForm 模型配置表单（新增/编辑共用）
func (m *Model) viewProviderForm(isEdit bool) string {
	title := "＋ 添加配置"
	if isEdit {
		title = "✎ 编辑配置"
	}

	dialogWidth := m.providerDialogWidth()

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(PrimaryDim).
		Width(dialogWidth-10).
		Padding(0, 1)

	// 当前焦点的输入框高亮 - 聚焦时使用主色
	nameStyle := inputStyle
	baseURLStyle := inputStyle
	apiKeyStyle := inputStyle
	modelStyle := inputStyle
	thinkingModelStyle := inputStyle
	defaultHaikuStyle := inputStyle
	defaultSonnetStyle := inputStyle
	defaultOpusStyle := inputStyle

	if m.providerInputFocus == FocusProviderName {
		nameStyle = nameStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderBaseURL {
		baseURLStyle = baseURLStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderAPIKey {
		apiKeyStyle = apiKeyStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderModel {
		modelStyle = modelStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderThinkingModel {
		thinkingModelStyle = thinkingModelStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderDefaultHaiku {
		defaultHaikuStyle = defaultHaikuStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderDefaultSonnet {
		defaultSonnetStyle = defaultSonnetStyle.BorderForeground(PrimaryColor)
	} else if m.providerInputFocus == FocusProviderDefaultOpus {
		defaultOpusStyle = defaultOpusStyle.BorderForeground(PrimaryColor)
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
					lipgloss.NewStyle().Foreground(SecondaryText).Render("主模型"),
					modelStyle.Render(m.providerModelInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("推理模型"),
					thinkingModelStyle.Render(m.providerThinkingModelInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("Haiku 默认模型"),
					defaultHaikuStyle.Render(m.providerDefaultHaikuInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("Sonnet 默认模型"),
					defaultSonnetStyle.Render(m.providerDefaultSonnetInput.View()),
					"",
					lipgloss.NewStyle().Foreground(SecondaryText).Render("Opus 默认模型"),
					defaultOpusStyle.Render(m.providerDefaultOpusInput.View()),
				),
				tipDisplay,
				errDisplay,
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 保存  ·  [Tab] 切换  ·  [Esc] 取消"),
			),
		)

	return dialog
}

// viewProviderDelete 确认删除配置
func (m *Model) viewProviderDelete() string {
	current := m.safeGetSelectedProvider()
	if current == nil {
		return ""
	}
	message := fmt.Sprintf("确认删除配置 '%s' ？", current.Name)

	// 动态计算对话框宽度
	dialogWidth := min(50, max(35, int(float64(m.width)*0.6)))

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

	return dialog
}

func (m *Model) renderProviderHelpItem(key, label string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		HelpKeyStyle.Render(key),
		" ",
		lipgloss.NewStyle().Foreground(SecondaryText).Render(label),
	)
}

// viewBatchAddProject 批量添加项目视图
func (m *Model) viewBatchAddProject() string {
	// 动态计算对话框宽度
	dialogWidth := min(60, max(45, int(float64(m.width)*0.75)))

	var items []string
	for i, path := range m.batchProjects {
		selected := m.batchSelected[i]
		cursor := "  "
		if i == m.batchCursor {
			cursor = "▸ "
		}
		checkbox := "[ ]"
		if selected {
			checkbox = "[×]"
		}
		items = append(items, cursor+checkbox+"  "+path)
	}

	if len(items) == 0 {
		items = append(items, "  没有找到可添加的项目")
	}

	selectedCount := len(m.filterBatchSelected())
	statusText := fmt.Sprintf("已选择: %d/%d 个项目", selectedCount, len(m.batchProjects))

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
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("＋ 批量添加项目"),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render("~/.claude/projects"),
				"",
				lipgloss.JoinVertical(lipgloss.Left, items...),
				"",
				lipgloss.NewStyle().Foreground(SecondaryText).Render(statusText),
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Space] 选择  ·  [↑↓] 移动  ·  [Enter] 确认  ·  [Esc] 取消"),
			),
		)

	return dialog
}

// filterBatchSelected 返回选中的项目路径列表
func (m *Model) filterBatchSelected() []string {
	var result []string
	for i, isSelected := range m.batchSelected {
		if isSelected {
			result = append(result, m.batchProjects[i])
		}
	}
	return result
}

// loadBatchProjects 加载 ~/.claude/projects 目录下的项目
func (m *Model) loadBatchProjects() {
	claudeProjectsDir := os.ExpandEnv("$HOME/.claude/projects")

	entries, err := os.ReadDir(claudeProjectsDir)
	if err != nil {
		return
	}

	m.batchProjects = nil
	m.batchSelected = make(map[int]bool)
	m.batchCursor = 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 读取 sessions-index.json 获取 originalPath
		sessionIndexPath := claudeProjectsDir + "/" + entry.Name() + "/sessions-index.json"
		data, err := os.ReadFile(sessionIndexPath)
		if err != nil {
			continue
		}

		var sessionIndex struct {
			OriginalPath string `json:"originalPath"`
		}
		if err := json.Unmarshal(data, &sessionIndex); err != nil {
			continue
		}

		if sessionIndex.OriginalPath == "" {
			continue
		}

		// 检查路径是否有效
		if _, err := os.Stat(sessionIndex.OriginalPath); err != nil {
			continue
		}

		// 检查是否已经添加过
		if m.isProjectAlreadyAdded(sessionIndex.OriginalPath) {
			continue
		}

		m.batchProjects = append(m.batchProjects, sessionIndex.OriginalPath)
	}
}

// isProjectAlreadyAdded 检查项目是否已经添加过
func (m *Model) isProjectAlreadyAdded(path string) bool {
	for _, p := range m.store.Projects {
		if p.Path == path {
			return true
		}
	}
	return false
}
