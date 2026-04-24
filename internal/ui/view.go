package ui

import (
	"fmt"
	"strings"

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

	// 使用更丰富的详情展示
	infoLine := func(label, value string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(SecondaryText).Width(10).Render(label),
			lipgloss.NewStyle().Foreground(Foreground).Render(value),
		)
	}

	openCountBadge := lipgloss.NewStyle().
		Foreground(SuccessColor).
		Background(lipgloss.Color("#0F2618")).
		Padding(0, 1).
		Render("打开 " + string(rune(p.OpenCount+'0')) + " 次")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(p.Alias),
		"",
		infoLine("路径", p.Path),
		infoLine("打开次数", openCountBadge),
		infoLine("创建时间", p.CreatedAt.Format("2006-01-02 15:04:05")),
		infoLine("最后打开", p.LastOpened.Format("2006-01-02 15:04:05")),
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
	// 计算分隔线宽度 (使用终端宽度减去边距)
	sepWidth := m.width - 4
	if sepWidth < 0 {
		sepWidth = 0
	}

	// 使用线条构建标题栏
	titleBar := lipgloss.NewStyle().
		Foreground(MutedText).
		Render(
			"┏" + lipgloss.NewStyle().Foreground(PrimaryDim).Render(strings.Repeat("━", sepWidth)) + "┓",
		)

	titleContent := "⚙ Gocoding " + lipgloss.NewStyle().Foreground(MutedText).Render("·") + " 项目管理"
	titleText := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(MutedText).Render("┃ "),
		lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true).Render(titleContent),
		lipgloss.NewStyle().Foreground(MutedText).Render(" ┃"),
	)

	bottomBar := lipgloss.NewStyle().
		Foreground(MutedText).
		Render(
			"┗" + lipgloss.NewStyle().Foreground(PrimaryDim).Render(strings.Repeat("━", sepWidth)) + "┛",
		)

	config := m.calculateLayout(m.width, m.height-m.debugPanelHeight())
	helpNav := m.renderHelpText(config)
	content := m.renderContent(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Padding(1, 0).
			Render("┃  暂无项目按 [n] 添加  ┃")
	}

	// 错误消息
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.Render("✗ " + m.errMsg)
	}

	// 主内容
	mainContent := lipgloss.NewStyle().
		Width(m.width).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				titleBar,
				titleText,
				bottomBar,
				"",
				content,
				emptyMsg,
				errDisplay,
				"",
				helpNav,
			),
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
						HelpKeyStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[q]"),
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
						HelpKeyStyle.Render("↑↓"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("导航"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[b]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("批量"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[e]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("编辑"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[d]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("删除"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[p]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("配置"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[q]"),
						quit,
					),
				),
			)
	default:
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, sep,
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("↑↓/j/k"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("导航"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[b]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("批量"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[e]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("编辑"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[d]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("删除"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[p]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("配置"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[Enter]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("IDE"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[1-4]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("快速"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[q]"),
						quit,
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

	// 搜索框 - 明确前缀，避免某些终端下输入文本不明显
	searchValue := m.searchQuery
	if searchValue == "" {
		searchValue = lipgloss.NewStyle().Foreground(MutedText).Render("输入关键词")
	}
	searchBox := lipgloss.NewStyle().
		Foreground(Foreground).
		Background(Background).
		Padding(0, 1).
		Render("Search: " + searchValue + "_")

	// 搜索状态
	statusText := lipgloss.NewStyle().Foreground(SecondaryText).Render(
		fmt.Sprintf("%d / %d 项目", len(m.list.Items()), m.store.Len()),
	)

	// 帮助文本
	helpText := lipgloss.NewStyle().
		Foreground(SecondaryText).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				" ",
				HelpKeyStyle.Render("[↑↓]"),
				" 选择 ",
				HelpKeyStyle.Render("[Enter]"),
				" 打开 ",
				HelpKeyStyle.Render("[Esc]"),
				" 返回",
			),
		)

	listView := m.renderSingleColumn(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = lipgloss.NewStyle().
			Foreground(SecondaryText).
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
