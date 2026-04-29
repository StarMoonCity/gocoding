package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
	"gocoding/internal/ui"
)

// SearchPageState 搜索页面状态
type SearchPageState int

const (
	SearchStateList SearchPageState = iota
	SearchStateIDEMenu
)

// SearchPage 搜索页面
type SearchPage struct {
	store       *models.ProjectStore
	app         *AppModel
	list        list.Model
	width       int
	height      int
	searchQuery string
	state       SearchPageState

	// IDE 菜单
	ideMenu      *IDEMenu
	ideAvailable map[models.IDEType]bool
}

// NewSearchPage 创建搜索页面
func NewSearchPage(store *models.ProjectStore) *SearchPage {
	p := &SearchPage{
		store: store,
		state: SearchStateList,
	}

	items := newListItems(store.Projects)
	delegate := projectListDelegate{searchQuery: &p.searchQuery}

	p.list = list.New(items, delegate, 60, 14)
	p.list.SetShowTitle(false)
	p.list.SetShowStatusBar(false)
	p.list.SetShowHelp(false)
	p.list.SetFilteringEnabled(true)

	p.initIDEMenu()

	return p
}

func (p *SearchPage) initIDEMenu() {
	p.ideMenu = &IDEMenu{
		title: "选择 IDE",
		options: []IDEOption{
			{Type: models.IDEClaudeCode, Name: "Claude", Description: "Claude Code IDE"},
			{Type: models.IDEVSCode, Name: "VSCode", Description: "Visual Studio Code"},
			{Type: models.IDEOpenCode, Name: "OpenCode", Description: "OpenCode IDE"},
			{Type: models.IDECodexCLI, Name: "Codex", Description: "Codex CLI"},
		},
		available: make(map[models.IDEType]bool),
	}
	p.ideAvailable = make(map[models.IDEType]bool)
}

// SetApp 设置 AppModel 引用
func (p *SearchPage) SetApp(app *AppModel) {
	p.app = app
}

// PageType 返回页面类型
func (p *SearchPage) PageType() PageType {
	return PageSearch
}

// OnActivate 页面激活时调用
func (p *SearchPage) OnActivate() {
	p.state = SearchStateList
	p.searchQuery = ""
	p.updateListItems()
}

// OnDeactivate 页面停用时调用
func (p *SearchPage) OnDeactivate() {
}

// SetSize 设置尺寸
func (p *SearchPage) SetSize(width, height int) {
	p.width = width
	p.height = height

	// 动态计算列表尺寸
	listWidth := min(80, max(50, width-4))
	listHeight := max(8, height-10)
	p.list.SetSize(listWidth, listHeight)
}

// Update 处理消息
func (p *SearchPage) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch p.state {
		case SearchStateList:
			return p.handleListKeyMsg(msg), true
		case SearchStateIDEMenu:
			return p.handleIDEMenuKeyMsg(msg), true
		}
	}

	// 更新列表
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return cmd, false
}

// View 渲染页面
func (p *SearchPage) View(width, height int) string {
	p.width = width
	p.height = height

	switch p.state {
	case SearchStateIDEMenu:
		return p.viewIDEMenu()
	default:
		return p.viewList()
	}
}

// HandleMouse 处理鼠标消息
func (p *SearchPage) HandleMouse(msg tea.MouseMsg) {
	switch p.state {
	case SearchStateList:
		p.handleListMouse(msg)
	}
}

// ============== 视图方法 ==============

func (p *SearchPage) viewList() string {
	// 搜索框
	searchValue := p.searchQuery
	if searchValue == "" {
		searchValue = lipgloss.NewStyle().Foreground(ui.MutedText).Render("输入关键词")
	}
	borderColor := ui.PrimaryDim
	if p.searchQuery != "" {
		borderColor = ui.AccentCyan
	}
	searchBoxWidth := max(40, p.width-10)
	searchBox := lipgloss.NewStyle().
		Width(searchBoxWidth).
		Foreground(ui.Foreground).
		Background(ui.BackgroundSurface).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(borderColor).
		Padding(0, 1).
		Render("Search: " + searchValue + "_")

	// 搜索状态
	matchCount := len(p.list.Items())
	totalCount := p.store.Len()
	var statusText string
	if p.searchQuery != "" {
		if matchCount > 0 {
			statusText = lipgloss.NewStyle().Foreground(ui.SuccessColor).Render(
				fmt.Sprintf("匹配 %d/%d 项目", matchCount, totalCount),
			)
		} else {
			statusText = lipgloss.NewStyle().Foreground(ui.WarningColor).Render(
				fmt.Sprintf("无匹配 (%d 项目总计)", totalCount),
			)
		}
	} else {
		statusText = lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render(
			fmt.Sprintf("%d 个项目", totalCount),
		)
	}

	// 帮助文本
	helpTextWidth := max(40, p.width-10)
	helpText := lipgloss.NewStyle().
		Width(helpTextWidth).
		Foreground(ui.SecondaryText).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				ui.HelpKeyNavStyle.Render("[↑↓]"),
				" 选择 ",
				ui.HelpKeyActionStyle.Render("[Enter]"),
				" 打开 ",
				ui.HelpKeyNavStyle.Render("[1-4]"),
				" 快速打开 ",
				ui.HelpKeyQuitStyle.Render("[Esc]"),
				" 返回",
			),
		)

	listView := p.list.View()

	emptyMsg := ""
	if len(p.list.Items()) == 0 && p.searchQuery != "" {
		emptyMsg = lipgloss.NewStyle().
			Foreground(ui.WarningColor).
			Padding(1, 0).
			Render("没有匹配的项目")
	}

	mainContent := lipgloss.JoinVertical(
				lipgloss.Left,
				searchBox,
				"",
				listView,
				emptyMsg,
				"",
				statusText,
				"",
				helpText,
			)

	// 上下居中，左右左对齐带间距
	return lipgloss.Place(p.width, p.height, lipgloss.Left, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, "  ", mainContent),
	)
}

func (p *SearchPage) viewIDEMenu() string {
	return lipgloss.Place(
		p.width,
		p.height,
		lipgloss.Center,
		lipgloss.Center,
		p.ideMenu.View(p.width, p.height),
	)
}

// ============== 处理器 ==============

func (p *SearchPage) handleListKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// 返回项目列表
		if p.app != nil {
			p.app.SwitchPage(p.app.projectPage)
			p.app.projectPage.OnActivate()
		}
	case "enter":
		if p.getSelectedProject() != nil {
			p.state = SearchStateIDEMenu
			for _, opt := range p.ideMenu.options {
				p.ideAvailable[opt.Type] = p.app.IDEExec().IsIDEAvailable(opt.Type)
				p.ideMenu.available[opt.Type] = p.ideAvailable[opt.Type]
			}
		}
	case "j", "down":
		p.list.CursorDown()
	case "k", "up":
		p.list.CursorUp()
	case "backspace", "ctrl+h":
		// 删除最后一个字符
		searchRunes := []rune(p.searchQuery)
		if len(searchRunes) > 0 {
			p.searchQuery = string(searchRunes[:len(searchRunes)-1])
		}
		p.updateListItems()
	case "1":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEOpenCode)
		}
	case "4":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDECodexCLI)
		}
	case "ctrl+c", "ctrl+q":
		return tea.Quit
	default:
		// 处理搜索输入
		if len(msg.Runes) > 0 {
			p.searchQuery += string(msg.Runes[0])
			p.updateListItems()
		}
	}
	return nil
}

func (p *SearchPage) handleIDEMenuKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down":
		p.ideMenu.selected = min(p.ideMenu.selected+1, len(p.ideMenu.options)-1)
	case "k", "up":
		p.ideMenu.selected = max(p.ideMenu.selected-1, 0)
	case "enter":
		selectedIDE := p.ideMenu.options[p.ideMenu.selected].Type
		return p.openWithIDE(selectedIDE)
	case "1":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDEOpenCode)
		}
	case "4":
		if p.getSelectedProject() != nil {
			return p.openWithIDE(models.IDECodexCLI)
		}
	case "esc":
		p.state = SearchStateList
	}
	return nil
}

// ============== 辅助方法 ==============

func (p *SearchPage) getSelectedProject() *models.Project {
	items := p.list.Items()
	if len(items) == 0 || p.list.Index() >= len(items) {
		return nil
	}

	item, ok := items[p.list.Index()].(listItem)
	if !ok {
		return nil
	}
	return &item.project
}

func (p *SearchPage) openWithIDE(ideType models.IDEType) tea.Cmd {
	project := p.getSelectedProject()
	if project == nil {
		return nil
	}

	// 检查 IDE 是否可用
	if !p.app.IDEExec().IsIDEAvailable(ideType) {
		p.app.ShowToast("IDE 不可用", "error")
		p.state = SearchStateList
		return nil
	}

	// 执行打开
	if err := p.app.IDEExec().OpenProject(project, ideType); err != nil {
		p.app.ShowToast("打开失败: "+err.Error(), "error")
	} else {
		p.app.ShowToast("已打开: "+project.Alias, "success")
	}

	p.state = SearchStateList
	return tea.Quit
}

func (p *SearchPage) handleListMouse(msg tea.MouseMsg) {
	listStartY := 4
	listEndY := p.height - 7

	switch {
	case msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress:
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := msg.Y - listStartY + p.list.Index()
			items := p.list.Items()
			if clickedIndex >= 0 && clickedIndex < len(items) {
				p.list.Select(clickedIndex)
			}
		}
	case msg.Button == tea.MouseButtonWheelUp:
		p.list.CursorUp()
	case msg.Button == tea.MouseButtonWheelDown:
		p.list.CursorDown()
	}
}

// updateListItems 根据搜索查询更新列表项
func (p *SearchPage) updateListItems() {
	results := p.store.Search(p.searchQuery)
	p.list.SetItems(newListItems(results))
}
