package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateViewport() {
	if len(m.list.Items()) == 0 {
		m.viewport.SetContent("")
		return
	}
	current := m.list.SelectedItem().(listItem)
	p := current.project

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
	switch m.state {
	case StateAddProject:
		return m.viewAddProject()
	case StateRenameProject:
		return m.viewRenameProject()
	case StateDeleteConfirm:
		return m.viewDeleteConfirm()
	case StateIDEMenu:
		return m.viewIDEMenu()
	case StateViewDetail:
		return m.viewViewDetail()
	case StateEditDescription:
		return m.viewEditDescription()
	case StateSearch:
		return m.viewSearch()
	case StateProviderList:
		return m.viewProviderList()
	case StateProviderAdd:
		return m.viewProviderForm(false)
	case StateProviderEdit:
		return m.viewProviderForm(true)
	case StateProviderDelete:
		return m.viewProviderDelete()
	default:
		return m.viewList()
	}
}

func (m *Model) viewList() string {
	// 标题 - 带装饰
	headerStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center)

	titleText := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(MutedText).Render("┌"),
		lipgloss.NewStyle().Foreground(PrimaryColor).Render(" Gocoding "),
		lipgloss.NewStyle().Foreground(MutedText).Render("项目管理 "),
		lipgloss.NewStyle().Foreground(MutedText).Render("┐"),
	)

	header := headerStyle.Render(titleText)

	helpNav := m.renderHelpText()
	config := m.calculateLayout(m.width, m.height)
	content := m.renderContent(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Align(lipgloss.Center).
			Padding(1, 0).
			Render("暂无项目按 [n] 添加")
	}

	// 错误消息
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = ErrorBoxStyle.Render("✗ " + m.errMsg)
	}

	// 主内容
	mainContent := lipgloss.NewStyle().
		Width(m.width - config.paddingX*2).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				header,
				"",
				content,
				emptyMsg,
				errDisplay,
				"",
				helpNav,
			),
		)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)
}

// renderHelpText 根据屏幕宽度渲染不同长度的帮助文本
func (m *Model) renderHelpText() string {
	config := m.calculateLayout(m.width, m.height)

	// 快捷键分隔符
	sep := lipgloss.NewStyle().Foreground(MutedText).Render(" │ ")

	quit := lipgloss.NewStyle().Foreground(SecondaryText).Render("退出")
	nav := lipgloss.NewStyle().Foreground(SecondaryText).Render("导航")

	switch config.helpTextMode {
	case HelpTextCompact:
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, sep,
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[↑↓]"),
						nav,
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[n]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("添加"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[/]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("搜索"),
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
	case HelpTextNormal:
		return lipgloss.NewStyle().
			Foreground(SecondaryText).
			MarginLeft(config.paddingX).
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, sep,
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[↑↓]"),
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
						HelpKeyStyle.Render("[r]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("重命名"),
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
						HelpKeyStyle.Render("[↑↓/j/k]"),
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
						HelpKeyStyle.Render("[v]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("详情"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[e]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("描述"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[r]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("重命名"),
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
						HelpKeyStyle.Render("[1/2/3]"),
						lipgloss.NewStyle().Foreground(SecondaryText).Render("快速打开"),
					),
					lipgloss.JoinHorizontal(lipgloss.Left, " ",
						HelpKeyStyle.Render("[q]"),
						quit,
					),
				),
			)
	}
}

// renderContent 渲染项目列表内容（支持单/双列）
func (m *Model) renderContent(config LayoutConfig) string {
	if m.layoutMode == LayoutDouble && len(m.list.Items()) > 0 {
		return m.renderDoubleColumn(config)
	}
	return m.renderSingleColumn(config)
}

// renderSingleColumn 单列布局
func (m *Model) renderSingleColumn(config LayoutConfig) string {
	return lipgloss.NewStyle().
		Width(m.width - config.paddingX*2).
		Align(lipgloss.Center).
		Render(m.list.View())
}

// renderDoubleColumn 双列布局
func (m *Model) renderDoubleColumn(config LayoutConfig) string {
	items := m.list.Items()
	itemCount := len(items)
	halfCount := (itemCount + 1) / 2

	var leftItems, rightItems []list.Item
	for i, item := range items {
		if i < halfCount {
			leftItems = append(leftItems, item)
		} else {
			rightItems = append(rightItems, item)
		}
	}

	// 创建左右两个列表
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)

	leftList := list.New(leftItems, delegate, config.listWidth, config.listHeight)
	rightList := list.New(rightItems, delegate, config.listWidth, config.listHeight)

	leftView := lipgloss.NewStyle().
		Width(config.listWidth).
		Render(leftList.View())

	rightView := lipgloss.NewStyle().
		Width(config.listWidth).
		Render(rightList.View())

	// 使用列间隙连接
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		leftView,
		lipgloss.NewStyle().Width(config.columnGap).Render(""),
		rightView,
	)
}

// viewSearch 搜索视图
func (m *Model) viewSearch() string {
	config := m.calculateLayout(m.width, m.height)

	// 搜索框 - 带图标和装饰
	searchIcon := lipgloss.NewStyle().Foreground(PrimaryColor).Render("🔍")
	searchBox := lipgloss.NewStyle().
		Background(BackgroundLight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 2).
		Width(m.width - config.paddingX*2 - 10).
		Render(searchIcon + " " + m.searchQuery + "_")

	// 搜索状态
	statusText := lipgloss.NewStyle().Foreground(SecondaryText).Render(
		fmt.Sprintf("%d / %d 项目", len(m.list.Items()), m.store.Len()),
	)

	// 帮助文本
	helpText := lipgloss.NewStyle().
		Foreground(SecondaryText).
		MarginLeft(config.paddingX).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Foreground(MutedText).Render("│"),
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
			Align(lipgloss.Center).
			Padding(1, 0).
			Render("没有匹配的项目")
	}

	mainContent := lipgloss.NewStyle().
		Width(m.width - config.paddingX*2).
		Align(lipgloss.Center).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
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

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)
}
