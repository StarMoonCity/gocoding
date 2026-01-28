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

	content := fmt.Sprintf(`项目详情

名称: %s
路径: %s
打开次数: %d
创建时间: %s
最后打开: %s

描述:
%s`, p.Alias, p.Path, p.OpenCount,
		p.CreatedAt.Format("2006-01-02 15:04:05"),
		p.LastOpened.Format("2006-01-02 15:04:05"),
		p.Description)

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
	default:
		return m.viewList()
	}
}

func (m *Model) viewList() string {
	// 响应式标题大小
	header := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Render("Gocoding · 项目管理")

	// 响应式帮助文本
	helpNav := m.renderHelpText()

	config := m.calculateLayout(m.width, m.height)

	content := m.renderContent(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Align(lipgloss.Center).
			Padding(1, 0).
			Render("没有项目，按 [n] 添加")
	}

	// 错误消息显示
	var errDisplay string
	if m.errMsg != "" {
		errDisplay = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(lipgloss.Color("#2C1810")).
			Padding(1, 2).
			MarginBottom(1).
			Align(lipgloss.Center).
			Width(m.width - config.paddingX*2 - 4).
			Render("⚠ " + m.errMsg)
	}

	// 主内容区域（设置宽度以居中显示）
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

	switch config.helpTextMode {
	case HelpTextCompact:
		return lipgloss.NewStyle().
			Foreground(SecondaryColor).
			MarginLeft(config.paddingX).
			Render("↑/↓ 移动 | n 添加 | / 搜索 | q 退出")
	case HelpTextNormal:
		return lipgloss.NewStyle().
			Foreground(SecondaryColor).
			MarginLeft(config.paddingX).
			Render("↑/↓ j/k 移动 | n 添加 | / 搜索 | r 重命名 | d 删除 | enter 选择 | q 退出")
	default:
		return lipgloss.NewStyle().
			Foreground(SecondaryColor).
			MarginLeft(config.paddingX).
			Render("↑/↓ j/k 导航 | / 搜索 | n 添加 | v 详情 | e 描述 | r 重命名 | d 删除 | enter IDE | 1/2/3 快速打开 | q 退出")
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

	// 搜索框
	searchBox := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Background(lipgloss.Color("#1a1a2e")).
		Padding(0, 1).
		Render("搜索: " + m.searchQuery + "_")

	// 搜索状态提示
	statusText := fmt.Sprintf("%d/%d 项目", len(m.list.Items()), m.store.Len())

	status := lipgloss.NewStyle().
		Foreground(SecondaryColor).
		MarginLeft(config.paddingX).
		Render(statusText)

	// 帮助文本
	helpText := lipgloss.NewStyle().
		Foreground(SecondaryColor).
		MarginLeft(config.paddingX).
		Render("↑/↓ j/k 选择 | enter IDE | 1/2/3 快速打开 | backspace 删除 | esc 返回")

	// 渲染列表
	listView := m.renderSingleColumn(config)

	emptyMsg := ""
	if len(m.list.Items()) == 0 {
		emptyMsg = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Align(lipgloss.Center).
			Padding(1, 0).
			Render("没有匹配的项目")
	}

	// 主内容区域
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
				status,
				"",
				helpText,
			),
		)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		mainContent,
	)
}
