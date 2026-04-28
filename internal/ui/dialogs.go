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

	inactiveInput := InputBorder.
		Width(dialogWidth - 6).
		Padding(0, 1)
	focusedInput := FocusedInputBorder.
		Width(dialogWidth - 6).
		Padding(0, 1)

	pathInputStyle := inactiveInput
	nameInputStyle := inactiveInput
	if m.inputFocus == FocusPath {
		pathInputStyle = focusedInput
	} else if m.inputFocus == FocusName {
		nameInputStyle = focusedInput
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(PrimaryColor).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("＋ 添加项目"),
				"",
				lipgloss.NewStyle().Foreground(ForegroundDim).Render("项目路径"),
				pathInputStyle.Render(m.input.View()),
				"",
				lipgloss.NewStyle().Foreground(ForegroundDim).Render("项目名称"),
				nameInputStyle.Render(m.secondaryInput.View()),
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

	inactiveInput := InputBorder.
		Width(dialogWidth - 6).
		Padding(0, 1)
	focusedInput := FocusedInputBorder.
		Width(dialogWidth - 6).
		Padding(0, 1)

	// 当前焦点的输入框高亮
	pathStyle := inactiveInput
	nameStyle := inactiveInput
	if m.inputFocus == FocusPath {
		pathStyle = focusedInput
	} else if m.inputFocus == FocusName {
		nameStyle = focusedInput
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(PrimaryColor).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render("✎ 编辑项目"),
				"",
				lipgloss.NewStyle().Foreground(ForegroundDim).Render("项目路径"),
				pathStyle.Render(m.input.View()),
				"",
				lipgloss.NewStyle().Foreground(ForegroundDim).Render("项目名称"),
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

	// 按钮样式 - 根据悬停状态
	buttonWidth := 10

	// 确定按钮（右侧）- 霓虹红
	confirmStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(ErrorColor).
		Background(lipgloss.Color("#1A0D10")).
		Padding(0, 2)
	if m.hoverButton == 1 {
		confirmStyle = confirmStyle.Background(ErrorColor).Foreground(Background).Bold(true)
	}

	// 取消按钮（左侧）- 霓虹灰
	cancelStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(SecondaryText).
		Background(BackgroundLight).
		Padding(0, 2)
	if m.hoverButton == 0 {
		cancelStyle = cancelStyle.Background(BackgroundHover).Foreground(Foreground)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(ErrorColor).
		Background(BackgroundSurface).
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
					cancelStyle.Render("[ N ] 否"),
					"  ",
					confirmStyle.Render("[ Y ] 是"),
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
		isSelected := i == m.ideMenu.selected
		isHovered := m.mouseEnabled && i == m.hoverIndex && !isSelected
		ideClr := ideColor(opt.Type)

		var prefix string
		var nameStyle lipgloss.Style
		if isSelected {
			prefix = lipgloss.NewStyle().Foreground(ideClr).Render("▸ ")
			nameStyle = lipgloss.NewStyle().Foreground(ideClr).Bold(true)
		} else if isHovered {
			prefix = "  "
			nameStyle = lipgloss.NewStyle().Foreground(ideClr).Bold(true)
		} else {
			prefix = "  "
			nameStyle = lipgloss.NewStyle().Foreground(Foreground).Bold(true)
		}
		var statusIcon string
		if available {
			statusIcon = lipgloss.NewStyle().Foreground(ideClr).Render("●")
		} else {
			statusIcon = lipgloss.NewStyle().Foreground(MutedText).Render("○")
		}
		// 每个 IDE 行末添加色彩指示条
		colorBar := lipgloss.NewStyle().Foreground(ideClr).Render("▌")
		options = append(options,
			prefix+statusIcon+"  "+
				nameStyle.Render(opt.Name)+
				lipgloss.NewStyle().Foreground(SecondaryText).Render("  "+opt.Description)+
				" "+colorBar)
	}

	dialogWidth := min(45, max(35, int(float64(m.width)*0.5)))

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(AccentCyan).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(AccentCyan).Bold(true).Render("▣ 选择 IDE"),
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
		Border(NeonBorder).
		BorderForeground(PrimaryColor).
		Background(BackgroundSurface).
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
		Border(NeonBorder).
		BorderForeground(PrimaryColor).
		Background(BackgroundSurface).
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
	dialogWidth := m.providerDialogWidth()
	contentWidth := m.providerListWidth()

	listView := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(m.providerList.View())

	// 错误消息 - 使用霓虹风格
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.
			Width(dialogWidth - 6).
			Render("✗ " + m.errMsg)
	}

	// 提示消息 - 使用霓虹风格
	var tipDisplay string
	if m.tipMsg != "" {
		tipDisplay = TipBoxStyle.
			Width(dialogWidth - 6).
			Render("ℹ " + m.tipMsg)
	}

	helpText := lipgloss.NewStyle().
		Foreground(SecondaryText).
		Width(contentWidth).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.renderProviderHelpItem("[k↑/j↓]", "选择", HelpKeyNavStyle),
					"  ",
					m.renderProviderHelpItem("[N]", "新增", HelpKeyActionStyle),
					"  ",
					m.renderProviderHelpItem("[E]", "编辑", HelpKeyActionStyle),
					"  ",
					m.renderProviderHelpItem("[D]", "删除", HelpKeyDangerStyle),
					"  ",
					m.renderProviderHelpItem("[A]", "激活", HelpKeyActionStyle),
				),
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					m.renderProviderHelpItem("[Esc]", "退出", HelpKeyQuitStyle),
				),
			),
		)

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(AccentMagenta).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(AccentMagenta).Bold(true).Render("⚙ 模型配置"),
				tipDisplay,
				errDisplay,
				"",
				listView,
				"",
				helpText,
			),
		)

	return dialog
}

// providerFormContent 生成表单内容（用于滚动区域）
func (m *Model) providerFormContent() string {
	dialogWidth := m.providerDialogWidth()

	// 基础输入框样式
	inactiveInput := ProviderInputStyle.Width(dialogWidth - 10).Padding(0, 1)
	focusedInput := ProviderFocusedInputStyle.Width(dialogWidth - 10).Padding(0, 1)

	// 辅助函数：根据焦点返回样式
	inputFor := func(focus ProviderInputFocus) lipgloss.Style {
		if m.providerInputFocus == focus {
			return focusedInput
		}
		return inactiveInput
	}

	// 区块标题辅助函数
	sectionHeader := func(title string, color lipgloss.Color) string {
		bar := lipgloss.NewStyle().Foreground(color).Render("▌")
		text := lipgloss.NewStyle().Foreground(color).Bold(true).Render(title)
		return lipgloss.JoinHorizontal(lipgloss.Left, bar, " ", text)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		// === Core 区块 ===
		sectionHeader("Core", PrimaryColor),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("配置名称"),
		inputFor(FocusProviderName).Render(m.providerNameInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("Base URL"),
		inputFor(FocusProviderBaseURL).Render(m.providerBaseURLInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("API Key"),
		inputFor(FocusProviderAPIKey).Render(m.providerAPIKeyInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("主模型"),
		inputFor(FocusProviderModel).Render(m.providerModelInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("推理模型"),
		inputFor(FocusProviderThinkingModel).Render(m.providerThinkingModelInput.View()),
		"",
		// === Default Models 区块 ===
		sectionHeader("Default Models", AccentCyan),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("Haiku 默认模型"),
		inputFor(FocusProviderDefaultHaiku).Render(m.providerDefaultHaikuInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("Sonnet 默认模型"),
		inputFor(FocusProviderDefaultSonnet).Render(m.providerDefaultSonnetInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("Opus 默认模型"),
		inputFor(FocusProviderDefaultOpus).Render(m.providerDefaultOpusInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("SubAgent 模型"),
		inputFor(FocusProviderSubagent).Render(m.providerSubagentInput.View()),
		"",
		// === Advanced 区块 ===
		sectionHeader("Advanced", AccentGold),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("禁用非必要流量"),
		inputFor(FocusProviderNonessential).Render(m.providerNonessentialInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("禁用非流式回退"),
		inputFor(FocusProviderNonstreaming).Render(m.providerNonstreamingInput.View()),
		"",
		lipgloss.NewStyle().Foreground(ForegroundDim).Render("推理力度"),
		inputFor(FocusProviderEffort).Render(m.providerEffortInput.View()),
	)
}

// viewProviderForm
func (m *Model) viewProviderForm(isEdit bool) string {
	title := "＋ 添加配置"
	borderColor := SuccessColor
	titleColor := SuccessColor
	if isEdit {
		title = "✎ 编辑配置"
		borderColor = PrimaryColor
		titleColor = PrimaryColor
	}

	// 更新表单内容的滚动区域
	m.providerFormViewport.SetContent(m.providerFormContent())

	dialogWidth := m.providerDialogWidth()

	// 错误消息 - 霓虹风格
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.
			Width(dialogWidth - 6).
			Render("✗ " + m.errMsg)
	}

	// 提示消息 - 霓虹风格
	var tipDisplay string
	if m.tipMsg != "" {
		tipDisplay = TipBoxStyle.
			Width(dialogWidth - 6).
			Render("ℹ " + m.tipMsg)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(borderColor).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(titleColor).Bold(true).Render(title),
				"",
				m.providerFormViewport.View(),
				tipDisplay,
				errDisplay,
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Enter] 保存  ·  [Tab/⇧Tab] 切换  ·  [↑↓] 滚动  ·  [Esc] 取消"),
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

	// 按钮样式 - 根据悬停状态
	buttonWidth := 10

	// 确定按钮（右侧）- 霓虹红
	confirmStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(ErrorColor).
		Background(lipgloss.Color("#1A0D10")).
		Padding(0, 2)
	if m.hoverButton == 1 {
		confirmStyle = confirmStyle.Background(ErrorColor).Foreground(Background).Bold(true)
	}

	// 取消按钮（左侧）- 霓虹灰
	cancelStyle := lipgloss.NewStyle().
		Width(buttonWidth).
		Foreground(SecondaryText).
		Background(BackgroundLight).
		Padding(0, 2)
	if m.hoverButton == 0 {
		cancelStyle = cancelStyle.Background(BackgroundHover).Foreground(Foreground)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(ErrorColor).
		Background(BackgroundSurface).
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
					cancelStyle.Render("[ N ] 否"),
					"  ",
					confirmStyle.Render("[ Y ] 是"),
				),
				"",
				lipgloss.NewStyle().Foreground(MutedText).Render("此操作不可恢复"),
			),
		)

	return dialog
}

func (m *Model) renderProviderHelpItem(key, label string, keyStyle lipgloss.Style) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		keyStyle.Render(key),
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
		isHovered := m.mouseEnabled && i == m.hoverIndex

		if i == m.batchCursor {
			cursor = lipgloss.NewStyle().Foreground(PrimaryColor).Render("▸ ")
		} else if isHovered {
			cursor = "  "
		}
		checkbox := "[ ]"
		if selected {
			checkbox = lipgloss.NewStyle().Foreground(SuccessColor).Render("[×]")
		}

		itemText := cursor + checkbox + "  " + path
		if isHovered && i != m.batchCursor {
			itemStyle := lipgloss.NewStyle().
				Foreground(PrimaryDim).
				Background(BackgroundHover)
			itemText = itemStyle.Render(itemText)
		}
		items = append(items, itemText)
	}

	if len(items) == 0 {
		items = append(items, lipgloss.NewStyle().Foreground(WarningColor).Render("  没有找到可添加的项目"))
	}

	selectedCount := len(m.filterBatchSelected())
	var statusText string
	if selectedCount > 0 {
		statusText = FeaturedBadgeStyle.Render(fmt.Sprintf("已选择: %d/%d", selectedCount, len(m.batchProjects)))
	} else {
		statusText = lipgloss.NewStyle().Foreground(ForegroundDim).Render(fmt.Sprintf("已选择: %d/%d", selectedCount, len(m.batchProjects)))
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(NeonBorder).
		BorderForeground(AccentGold).
		Background(BackgroundSurface).
		Foreground(Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(AccentGold).Bold(true).Render("＋ 批量添加项目"),
				"",
				lipgloss.NewStyle().Foreground(ForegroundDim).Render("~/.claude/projects"),
				"",
				lipgloss.JoinVertical(lipgloss.Left, items...),
				"",
				statusText,
				"",
				lipgloss.NewStyle().
					Foreground(SecondaryText).
					Render("[Space/点击] 选择  ·  [↑↓] 移动  ·  [Enter] 确认  ·  [Esc] 取消"),
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
