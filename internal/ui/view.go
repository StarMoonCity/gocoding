package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateViewport() {
	if len(m.list.Items()) == 0 {
		m.viewport.SetContent("")
		return
	}
	current := m.safeGetSelectedProject()
	if current == nil {
		m.viewport.SetContent("")
		return
	}
	p := *current

	// 彩色分隔线
	sectionDivider := SeparatorHighlightStyle.Render("─── 详情 ───")

	infoLine := func(label, value string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(ForegroundDim).Width(10).Render(label),
			lipgloss.NewStyle().Foreground(Foreground).Render(value),
		)
	}

	// 打开次数徽章 - 使用 BadgeStyle
	var openCountBadge string
	if p.OpenCount >= 10 {
		openCountBadge = FeaturedBadgeStyle.Render("打开 " + fmt.Sprintf("%d", p.OpenCount) + " 次")
	} else {
		openCountBadge = BadgeStyle.
			Foreground(SuccessColor).
			Render("打开 " + fmt.Sprintf("%d", p.OpenCount) + " 次")
	}

	// 时间戳颜色区分
	lastOpenedColor := SecondaryText
	createdColor := ForegroundDim
	if time.Since(p.LastOpened) < 24*time.Hour {
		lastOpenedColor = AccentCyan
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(p.Alias),
		"",
		sectionDivider,
		"",
		infoLine("路径", p.Path),
		infoLine("打开次数", openCountBadge),
		infoLine("创建时间", lipgloss.NewStyle().Foreground(createdColor).Render(p.CreatedAt.Format("2006-01-02 15:04:05"))),
		infoLine("最后打开", lipgloss.NewStyle().Foreground(lastOpenedColor).Render(p.LastOpened.Format("2006-01-02 15:04:05"))),
		"",
		lipgloss.NewStyle().Foreground(SecondaryText).Bold(true).MarginTop(1).Render("描述"),
		"",
		lipgloss.NewStyle().Foreground(Foreground).Render(p.Description),
	)

	m.viewport.SetContent(content)
}

func (m *Model) View() string {
	// 渲染主内容
	var content string
	switch m.state {
	case StateAddProject:
		content = m.viewAddProject()
	case StateRenameProject:
		content = m.viewRenameProject()
	case StateDeleteConfirm:
		content = m.viewDeleteConfirm()
	case StateIDEMenu:
		content = lipgloss.Place(
			m.width,
			max(0, m.height-m.debugPanelHeight()),
			lipgloss.Center,
			lipgloss.Center,
			m.viewIDEMenu(),
		)
	case StateViewDetail:
		content = m.viewViewDetail()
	case StateEditDescription:
		content = m.viewEditDescription()
	case StateSearch:
		content = m.viewSearch()
	case StateProviderList:
		content = m.viewProviderList()
	case StateProviderAdd:
		content = m.viewProviderForm(false)
	case StateProviderEdit:
		content = m.viewProviderForm(true)
	case StateProviderDelete:
		content = m.viewProviderDelete()
	case StateBatchAddProject:
		content = m.viewBatchAddProject()
	default:
		content = m.viewList()
	}

	// 调试面板（位于底部）
	if m.debug {
		content = lipgloss.JoinVertical(lipgloss.Left, content, m.renderDebugPanel())
	}

	return content
}

func (m *Model) viewList() string {
	// 渐变色标题栏
	gradientBar := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(PrimaryColor).Render("█"),
		lipgloss.NewStyle().Foreground(PrimaryColorAlt).Render("▓"),
		lipgloss.NewStyle().Foreground(PrimaryDim).Render("▒"),
		lipgloss.NewStyle().Foreground(PrimaryDark).Render("░"),
	)
	titleText := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Render("  Gocoding · 项目管理")
	headerBlock := lipgloss.JoinVertical(lipgloss.Left, gradientBar, titleText)

	config := m.calculateLayout(m.width, m.height-m.debugPanelHeight())
	helpNav := m.renderHelpText(config)
	content := m.renderContent(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = SurfaceStyle.Render("┃  暂无项目按 [n] 添加  ┃")
	}

	// 错误消息
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.Render("✗ " + m.errMsg)
	}

	// 主内容
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBlock,
		"",
		content,
		emptyMsg,
		errDisplay,
		"",
		helpNav,
	)

	return mainContent
}

// renderHelpText 根据屏幕宽度渲染不同长度的帮助文本
func (m *Model) renderHelpText(config LayoutConfig) string {
	// 使用竖线分隔符
	sep := lipgloss.NewStyle().Foreground(PrimaryDim).Render("│")

	quit := lipgloss.NewStyle().Foreground(SecondaryText).Render("退出")

	switch config.helpTextMode {
	case HelpTextCompact:
		// 超精简版: 只显示最常用的快捷键，垂直排列以节省空间
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyActionStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeySearchStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyQuitStyle.Render("[q]"),
						quit,
					),
				),
			)
	case HelpTextNormal:
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, sep,
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyNavStyle.Render("↑↓"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("导航"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeySearchStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyActionStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyActionStyle.Render("[b]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("批量"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyActionStyle.Render("[e]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("编辑"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyDangerStyle.Render("[d]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("删除"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyActionStyle.Render("[p]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("配置"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyQuitStyle.Render("[q]"),
						quit,
					),
				),
			)
	default:
		// 完整版：分成两行，避免超出80列
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					lipgloss.JoinHorizontal(lipgloss.Left, sep,
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyNavStyle.Render("↑↓/j/k"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("导航"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeySearchStyle.Render("[/]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[n]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[b]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("批量"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[e]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("编辑"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyDangerStyle.Render("[d]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("删除"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[p]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("配置"),
						),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, sep,
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[Enter]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("IDE"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyActionStyle.Render("[1-4]"),
							lipgloss.NewStyle().Foreground(SecondaryText).Render("快速"),
						),
						lipgloss.JoinHorizontal(lipgloss.Left, " ",
							HelpKeyQuitStyle.Render("[q]"),
							quit,
						),
					),
				),
			)
	}
}

// renderContent 渲染项目列表内容
func (m *Model) renderContent(config LayoutConfig) string {
	return m.renderSingleColumn(config)
}

// renderSingleColumn 单列布局
func (m *Model) renderSingleColumn(config LayoutConfig) string {
	return m.list.View()
}

// viewSearch 搜索视图
func (m *Model) viewSearch() string {
	config := m.calculateLayout(m.width, m.height-m.debugPanelHeight())

	// 搜索框 - 底部边框高亮
	searchValue := m.searchQuery
	if searchValue == "" {
		searchValue = lipgloss.NewStyle().Foreground(MutedText).Render("输入关键词")
	}
	borderColor := PrimaryDim
	if m.searchQuery != "" {
		borderColor = AccentCyan
	}
	searchBox := lipgloss.NewStyle().
		Foreground(Foreground).
		Background(BackgroundSurface).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render("Search: " + searchValue + "_")

	// 搜索状态 - 颜色编码
	matchCount := len(m.list.Items())
	totalCount := m.store.Len()
	var statusText string
	if m.searchQuery != "" {
		if matchCount > 0 {
			statusText = lipgloss.NewStyle().Foreground(SuccessColor).Render(
				fmt.Sprintf("匹配 %d/%d 项目", matchCount, totalCount),
			)
		} else {
			statusText = lipgloss.NewStyle().Foreground(WarningColor).Render(
				fmt.Sprintf("无匹配 (%d 项目总计)", totalCount),
			)
		}
	} else {
		statusText = lipgloss.NewStyle().Foreground(ForegroundDim).Render(
			fmt.Sprintf("%d 个项目", totalCount),
		)
	}

	// 帮助文本 - 分类颜色
	helpText := lipgloss.NewStyle().
		Foreground(SecondaryText).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				" ",
				HelpKeyNavStyle.Render("[↑↓]"),
				" 选择 ",
				HelpKeyActionStyle.Render("[Enter]"),
				" 打开 ",
				HelpKeyQuitStyle.Render("[Esc]"),
				" 返回",
			),
		)

	listView := m.renderSingleColumn(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 && m.searchQuery != "" {
		emptyMsg = lipgloss.NewStyle().
			Foreground(WarningColor).
			Padding(1, 0).
			Render("没有匹配的项目")
	}

	mainContent := lipgloss.NewStyle().
		Width(m.width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				searchBox,
				"",
				listView,
				emptyMsg,
				"",
				statusText,
				"",
				helpText,
			),
		)

	return mainContent
}
