package components

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
	"gocoding/internal/ui"
)

// ProjectPageState 项目页面状态
type ProjectPageState int

const (
	ProjectStateList ProjectPageState = iota
	ProjectStateAdd
	ProjectStateRename
	ProjectStateDeleteConfirm
	ProjectStateIDEMenu
	ProjectStateViewDetail
	ProjectStateEditDescription
)

// projectListDelegate 项目列表自定义渲染
type projectListDelegate struct {
	searchQuery *string
}

func (d projectListDelegate) Height() int                               { return 1 }
func (d projectListDelegate) Spacing() int                              { return 0 }
func (d projectListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d projectListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	proj, ok := item.(listItem)
	if !ok {
		return
	}

	isSelected := item == m.SelectedItem()
	listWidth := m.Width()

	separator := ui.SeparatorStyle.Render(" │ ")

	selector := "  "
	if isSelected {
		selector = lipgloss.NewStyle().Foreground(ui.SelectedBorder).Render("▸ ")
	}

	var recentDot string
	if time.Since(proj.project.LastOpened) < 1*time.Hour {
		recentDot = lipgloss.NewStyle().Foreground(ui.WarningColor).Render("•")
	}

	var countBadge string
	if proj.project.OpenCount >= 10 {
		countBadge = ui.FeaturedBadgeStyle.Render(fmt.Sprintf("×%d", proj.project.OpenCount))
	} else {
		countBadge = ui.BadgeStyle.Foreground(ui.SuccessColor).Render(fmt.Sprintf("×%d", proj.project.OpenCount))
	}

	// 截断过长的别名，避免撑破列表布局
	alias := proj.project.Alias
	nameBudget := m.Width() - lipgloss.Width(selector) - lipgloss.Width(recentDot) - lipgloss.Width(separator) - lipgloss.Width(countBadge)
	if lipgloss.Width(alias) > nameBudget {
		alias = truncateToWidth(alias, nameBudget)
	}

	var namePart string
	if d.searchQuery != nil && *d.searchQuery != "" {
		query := *d.searchQuery
		aliasLower := strings.ToLower(alias)
		queryLower := strings.ToLower(query)
		if idx := strings.Index(aliasLower, queryLower); idx >= 0 {
			before := alias[:idx]
			match := alias[idx : idx+len(query)]
			after := alias[idx+len(query):]
			if isSelected {
				namePart = lipgloss.JoinHorizontal(lipgloss.Left,
					ui.SelectedListItemStyle.Render(before),
					lipgloss.NewStyle().Foreground(ui.AccentGold).Bold(true).Render(match),
					ui.SelectedListItemStyle.Render(after),
				)
			} else {
				namePart = lipgloss.JoinHorizontal(lipgloss.Left,
					ui.ListItemStyle.Render(before),
					lipgloss.NewStyle().Foreground(ui.AccentGold).Bold(true).Render(match),
					ui.ListItemStyle.Render(after),
				)
			}
		} else {
			if isSelected {
				namePart = ui.SelectedListItemStyle.Render(alias)
			} else {
				namePart = ui.ListItemStyle.Render(alias)
			}
		}
	} else {
		if isSelected {
			namePart = ui.SelectedListItemStyle.Render(alias)
		} else {
			namePart = ui.ListItemStyle.Render(alias)
		}
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		selector,
		recentDot,
		namePart,
		separator,
		countBadge,
	)

	// 整行背景 - 逐行填充确保无透明区域
	var rowBg lipgloss.Color
	if isSelected {
		rowBg = ui.SelectedBg
	} else {
		rowBg = ui.Background
	}
	bgStyle := lipgloss.NewStyle().Background(rowBg)
	contentWidth := lipgloss.Width(content)
	if contentWidth < listWidth {
		content += bgStyle.Render(strings.Repeat(" ", listWidth-contentWidth))
	}
	row := content

	fmt.Fprintf(w, "%s", row)
}

// truncateToWidth 按显示宽度截断字符串（rune 安全），超宽时以省略号结尾
func truncateToWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	if maxWidth == 1 {
		return "…"
	}
	var sb strings.Builder
	w := 0
	for _, r := range s {
		rw := lipgloss.Width(string(r))
		if w+rw > maxWidth-1 {
			break
		}
		sb.WriteRune(r)
		w += rw
	}
	return sb.String() + "…"
}

// ProjectListPage 项目列表页面 - 完全自治
type ProjectListPage struct {
	store *models.ProjectStore
	app   *AppModel
	list  list.Model

	// 页面状态
	state ProjectPageState

	// 添加/编辑表单
	input          textinput.Model
	secondaryInput textinput.Model
	inputFocus     InputFocus
	tempPath       string
	editingID      string

	// 详情视图
	viewport viewport.Model

	// 描述编辑
	ta textarea.Model

	// IDE 菜单
	ideMenu      *IDEMenu
	ideAvailable map[models.IDEType]bool

	// 错误消息
	errMsg string

	// 鼠标状态
	hoverIndex   int
	hoverButton  int
	mouseEnabled bool

	// 尺寸
	width  int
	height int
}

// NewProjectListPage 创建项目列表页面
func NewProjectListPage(store *models.ProjectStore) *ProjectListPage {
	p := &ProjectListPage{
		store: store,
		state: ProjectStateList,
	}

	items := newListItems(store.Projects)
	delegate := projectListDelegate{searchQuery: new(string)}

	// 初始化列表，使用占位尺寸，实际尺寸在 SetSize 时计算
	p.list = list.New(items, delegate, 60, 14)
	p.list.Title = ""
	p.list.Styles.Title = lipgloss.NewStyle().Foreground(ui.SecondaryText).Bold(true).MarginBottom(1)
	p.list.SetShowTitle(false)
	p.list.SetShowStatusBar(false)
	p.list.SetShowHelp(false)
	p.list.SetFilteringEnabled(true)

	p.input = textinput.New()
	p.input.Placeholder = "输入项目路径"
	p.secondaryInput = textinput.New()
	p.secondaryInput.Placeholder = "输入项目名称"

	p.ta = textarea.New()
	p.ta.Focus()

	p.viewport = viewport.New(0, 0)

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

	// 自动定位到最近打开的项目
	if store.Len() > 0 {
		recent := store.GetMostRecentlyOpened()
		if recent != nil {
			index := store.GetIndexByProject(recent.ID)
			if index >= 0 {
				p.list.Select(index)
			}
		}
	}

	return p
}

// SetApp 设置 App 引用
func (p *ProjectListPage) SetApp(app *AppModel) {
	p.app = app
}

// PageType 返回页面类型
func (p *ProjectListPage) PageType() PageType {
	return PageProject
}

// OnActivate 页面激活时
func (p *ProjectListPage) OnActivate() {
	p.state = ProjectStateList
	p.syncListItems()
}

// OnDeactivate 页面停用时
func (p *ProjectListPage) OnDeactivate() {
}

// SetSize 设置尺寸
func (p *ProjectListPage) SetSize(width, height int) {
	p.width = width
	p.height = height

	// 动态计算列表尺寸 - 与 viewList 的 contentWidth 保持一致
	listWidth := max(60, width-10)
	listHeight := max(8, height-10)
	p.list.SetSize(listWidth, listHeight)
}

// Update 处理消息
func (p *ProjectListPage) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case openProjectSuccessMsg:
		if project := p.store.Get(msg.projectID); project != nil {
			project.UpdateLastOpened()
			p.store.SortByLastOpened()
			p.syncListItems()
			p.list.Select(0)
		}
		p.state = ProjectStateList
		return nil, true
	case tea.KeyMsg:
		switch p.state {
		case ProjectStateList:
			return p.handleListKeyMsg(msg), true
		case ProjectStateAdd:
			return p.handleAddKeyMsg(msg), true
		case ProjectStateRename:
			return p.handleRenameKeyMsg(msg), true
		case ProjectStateDeleteConfirm:
			return p.handleDeleteConfirmKeyMsg(msg), true
		case ProjectStateIDEMenu:
			return p.handleIDEMenuKeyMsg(msg), true
		case ProjectStateViewDetail:
			return p.handleViewDetailKeyMsg(msg), true
		case ProjectStateEditDescription:
			return p.handleEditDescKeyMsg(msg), true
		}
	}

	// 更新组件
	switch p.state {
	case ProjectStateAdd, ProjectStateRename:
		var cmd1, cmd2 tea.Cmd
		p.input, cmd1 = p.input.Update(msg)
		p.secondaryInput, cmd2 = p.secondaryInput.Update(msg)
		return tea.Batch(cmd1, cmd2), true
	case ProjectStateViewDetail:
		var cmd tea.Cmd
		p.viewport, cmd = p.viewport.Update(msg)
		return cmd, true
	case ProjectStateEditDescription:
		var cmd tea.Cmd
		p.ta, cmd = p.ta.Update(msg)
		return cmd, true
	}

	// 更新列表
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return cmd, false
}

// View 渲染页面
func (p *ProjectListPage) View(width, height int) string {
	p.width = width
	p.height = height

	switch p.state {
	case ProjectStateAdd:
		return p.viewAdd()
	case ProjectStateRename:
		return p.viewRename()
	case ProjectStateDeleteConfirm:
		return p.viewDeleteConfirm()
	case ProjectStateIDEMenu:
		return p.viewIDEMenu()
	case ProjectStateViewDetail:
		return p.viewDetail()
	case ProjectStateEditDescription:
		return p.viewEditDesc()
	default:
		return p.viewList()
	}
}

// HandleMouse 处理鼠标
func (p *ProjectListPage) HandleMouse(msg tea.MouseMsg) {
	p.mouseEnabled = true

	// 列表状态下支持滚轮导航，与其他页面保持一致
	if p.state == ProjectStateList {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			p.list.CursorUp()
		case tea.MouseButtonWheelDown:
			p.list.CursorDown()
		}
	}
}

// ============== 视图方法 ==============

func (p *ProjectListPage) viewList() string {
	titleText := lipgloss.NewStyle().
		Foreground(ui.PrimaryColor).
		Bold(true).
		Render("  Gocoding · 项目管理")

	helpNav := p.renderHelpText()
	content := p.list.View()

	emptyMsg := ""
	if len(p.list.Items()) == 0 {
		emptyMsg = ui.SurfaceStyle.Render("┃  暂无项目按 [n] 添加  ┃")
	}

	var errDisplay string
	if p.errMsg != "" {
		errDisplay = ui.ErrorBoxStyle.Render("✗ " + p.errMsg)
		p.errMsg = ""
	}

	// 计算内容区域
	listContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		content,
		emptyMsg,
		errDisplay,
		"",
		helpNav,
	)

	// 整体内容区域 - 带背景确保无透明区域
	contentWidth := max(60, p.width-10)
	fullContent := lipgloss.NewStyle().
		Width(contentWidth).
		Background(ui.Background).
		Render(
			lipgloss.JoinVertical(lipgloss.Left, titleText, listContent),
		)

	// 全屏背景，左对齐垂直居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Left, lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, "  ", fullContent))
}

func (p *ProjectListPage) viewAdd() string {
	dialogWidth := min(50, max(40, int(float64(p.width)*0.7)))

	inactiveInput := ui.InputBorder.Width(dialogWidth-6).Padding(0, 1)
	focusedInput := ui.FocusedInputBorder.Width(dialogWidth-6).Padding(0, 1)

	pathStyle := inactiveInput
	nameStyle := inactiveInput
	if p.inputFocus == FocusPath {
		pathStyle = focusedInput
	} else {
		nameStyle = focusedInput
	}

	var errDisplay string
	if p.errMsg != "" {
		errDisplay = ui.ErrorBoxStyle.Width(dialogWidth - 6).Render("✗ " + p.errMsg)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.PrimaryColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("＋ 添加项目"),
				"",
				lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render("项目路径"),
				pathStyle.Render(p.input.View()),
				"",
				lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render("项目名称"),
				nameStyle.Render(p.secondaryInput.View()),
				"",
				errDisplay,
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Enter] 确认  ·  [Tab] 切换  ·  [Esc] 取消"),
			),
		)

	// 全屏背景居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(dialog)
}

func (p *ProjectListPage) viewRename() string {
	dialogWidth := min(50, max(40, int(float64(p.width)*0.7)))

	inactiveInput := ui.InputBorder.Width(dialogWidth-6).Padding(0, 1)
	focusedInput := ui.FocusedInputBorder.Width(dialogWidth-6).Padding(0, 1)

	pathStyle := inactiveInput
	nameStyle := inactiveInput
	if p.inputFocus == FocusPath {
		pathStyle = focusedInput
	} else {
		nameStyle = focusedInput
	}

	var errDisplay string
	if p.errMsg != "" {
		errDisplay = ui.ErrorBoxStyle.Width(dialogWidth - 6).Render("✗ " + p.errMsg)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.PrimaryColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("✎ 编辑项目"),
				"",
				lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render("项目路径"),
				pathStyle.Render(p.input.View()),
				"",
				lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render("项目名称"),
				nameStyle.Render(p.secondaryInput.View()),
				"",
				errDisplay,
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Enter] 确认  ·  [Tab] 切换  ·  [Esc] 取消"),
			),
		)

	// 全屏背景居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(dialog)
}

func (p *ProjectListPage) viewDeleteConfirm() string {
	current := p.safeGetSelectedProject()
	if current == nil {
		return ""
	}

	dialogWidth := min(50, max(35, int(float64(p.width)*0.6)))
	buttonWidth := 10

	confirmStyle := ui.DangerButtonStyle.Width(buttonWidth)
	if p.hoverButton == 1 {
		confirmStyle = ui.DangerButtonHoverStyle.Width(buttonWidth)
	}

	cancelStyle := ui.ButtonStyle.Width(buttonWidth)
	if p.hoverButton == 0 {
		cancelStyle = ui.ButtonHoverStyle.Width(buttonWidth)
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.PrimaryDim).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.ErrorColor).Bold(true).Render("✗ 删除项目"),
				"",
				lipgloss.NewStyle().Foreground(ui.Foreground).Render(fmt.Sprintf("确认删除项目 '%s' ？", current.Alias)),
				"",
				lipgloss.JoinHorizontal(lipgloss.Center, cancelStyle.Render("[ N ] 否"), "  ", confirmStyle.Render("[ Y ] 是")),
				"",
				lipgloss.NewStyle().Foreground(ui.MutedText).Render("此操作不可恢复"),
			),
		)

	// 全屏背景居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(dialog)
}

func (p *ProjectListPage) viewIDEMenu() string {
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(p.ideMenu.View(p.width, p.height))
}

func (p *ProjectListPage) viewDetail() string {
	p.updateViewport()

	dialogWidth := max(50, int(float64(p.width)*0.8))

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.PrimaryColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Align(lipgloss.Center).Render("▤ 项目详情"),
				"",
				p.viewport.View(),
				"",
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Esc] 返回"),
			),
		)

	// 全屏背景居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(dialog)
}

func (p *ProjectListPage) viewEditDesc() string {
	dialogWidth := min(60, max(45, int(float64(p.width)*0.75)))
	p.ta.SetWidth(dialogWidth - 6)
	p.ta.SetHeight(max(3, p.height-16))

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.PrimaryColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render("✎ 编辑描述"),
				"",
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("描述（支持多行）"),
				p.ta.View(),
				"",
				lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Ctrl+S] 保存  ·  [Esc] 取消"),
			),
		)

	// 全屏背景居中
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Background(ui.Background).
		Align(lipgloss.Center, lipgloss.Center).
		Render(dialog)
}

// ============== 处理器 ==============

func (p *ProjectListPage) handleListKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down":
		p.list.CursorDown()
	case "k", "up":
		p.list.CursorUp()
	case "pgup", "pgdown", "home", "end":
		var cmd tea.Cmd
		p.list, cmd = p.list.Update(msg)
		return cmd
	case "/":
		p.app.SwitchPage(p.app.searchPage)
		p.app.searchPage.OnActivate()
	case "n":
		p.state = ProjectStateAdd
		p.input.Reset()
		p.input.Placeholder = "输入项目路径"
		p.input.SetValue("")
		p.input.Focus()
		p.inputFocus = FocusPath
		p.secondaryInput.Reset()
		p.secondaryInput.Placeholder = "输入项目名称"
		p.secondaryInput.SetValue("")
		p.secondaryInput.Blur()
		p.tempPath = ""
		p.errMsg = ""
		return textinput.Blink
	case "r":
		if current := p.safeGetSelectedProject(); current != nil {
			p.state = ProjectStateRename
			p.editingID = current.ID
			p.input.SetValue(current.Path)
			p.input.Placeholder = "输入项目路径"
			p.input.Focus()
			p.inputFocus = FocusPath
			p.secondaryInput.SetValue(current.Alias)
			p.secondaryInput.Placeholder = "输入项目名称"
			p.secondaryInput.Blur()
			p.errMsg = ""
		}
		return textinput.Blink
	case "e":
		if current := p.safeGetSelectedProject(); current != nil {
			p.ta.SetValue(current.Description)
			p.ta.Focus()
			p.state = ProjectStateEditDescription
			p.errMsg = ""
		}
	case "d":
		if p.safeGetSelectedProject() != nil {
			p.state = ProjectStateDeleteConfirm
		}
	case "v":
		if len(p.list.Items()) > 0 {
			p.state = ProjectStateViewDetail
			p.updateViewport()
		}
	case "enter":
		if p.safeGetSelectedProject() != nil {
			p.state = ProjectStateIDEMenu
			for _, opt := range p.ideMenu.options {
				p.ideAvailable[opt.Type] = p.app.IDEExec().IsIDEAvailable(opt.Type)
				p.ideMenu.available[opt.Type] = p.ideAvailable[opt.Type]
			}
		}
	case "1":
		if p.safeGetSelectedProject() != nil {
			return p.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if p.safeGetSelectedProject() != nil {
			return p.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if p.safeGetSelectedProject() != nil {
			return p.openWithIDE(models.IDEOpenCode)
		}
	case "4":
		if p.safeGetSelectedProject() != nil {
			return p.openWithIDE(models.IDECodexCLI)
		}
	case "p":
		p.app.SwitchPage(p.app.providerPage)
		p.app.providerPage.OnActivate()
	case "b":
		p.app.SwitchPage(p.app.batchAddPage)
		p.app.batchAddPage.OnActivate()
	case "q":
		return tea.Quit
	}
	return nil
}

func (p *ProjectListPage) handleAddKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		path := p.input.Value()
		name := p.secondaryInput.Value()

		if err := p.store.ValidatePath(path); err != nil {
			p.errMsg = err.Error()
			return nil
		}

		if p.store.PathExists(path) {
			p.errMsg = "该项目已存在"
			return nil
		}

		// 名称为空时自动从路径提取
		if name == "" {
			name = defaultAliasForPath(path)
			if name == "" {
				p.errMsg = "项目名称不能为空"
				return nil
			}
			p.secondaryInput.SetValue(name)
		}

		project := models.Project{
			ID:        generateID(),
			Path:      path,
			Alias:     name,
			CreatedAt: time.Now(),
		}
		p.store.Add(project)
		p.syncListItems()
		p.state = ProjectStateList
	case "tab":
		if p.inputFocus == FocusPath {
			p.inputFocus = FocusName
			p.input.Blur()
			p.secondaryInput.Focus()
		} else {
			p.inputFocus = FocusPath
			p.secondaryInput.Blur()
			p.input.Focus()
		}
	case "esc":
		p.state = ProjectStateList
	}

	var cmd tea.Cmd
	if p.inputFocus == FocusPath {
		p.input, cmd = p.input.Update(msg)
		path := p.input.Value()
		if path != p.tempPath {
			p.tempPath = path
			p.secondaryInput.SetValue(defaultAliasForPath(path))
		}
	} else {
		p.secondaryInput, cmd = p.secondaryInput.Update(msg)
	}
	// 输入变化时清除上一次的校验错误
	p.errMsg = ""
	return cmd
}

func (p *ProjectListPage) handleRenameKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		path := p.input.Value()
		name := p.secondaryInput.Value()
		if p.editingID != "" {
			current := p.store.Get(p.editingID)
			if current != nil && path != "" && path != current.Path {
				if err := p.store.ValidatePath(path); err != nil {
					p.errMsg = err.Error()
					return nil
				}
				if p.store.PathExists(path) {
					p.errMsg = "该项目已存在"
					return nil
				}
			}
			if name != "" || path != "" {
				p.store.Update(p.editingID, name, path)
				p.syncListItems()
			}
		}
		p.state = ProjectStateList
		p.editingID = ""
	case "tab":
		if p.inputFocus == FocusPath {
			p.inputFocus = FocusName
			p.input.Blur()
			p.secondaryInput.Focus()
		} else {
			p.inputFocus = FocusPath
			p.secondaryInput.Blur()
			p.input.Focus()
		}
	case "esc":
		p.state = ProjectStateList
		p.editingID = ""
	}

	var cmd tea.Cmd
	if p.inputFocus == FocusPath {
		p.input, cmd = p.input.Update(msg)
	} else {
		p.secondaryInput, cmd = p.secondaryInput.Update(msg)
	}
	// 输入变化时清除上一次的校验错误
	p.errMsg = ""
	return cmd
}

func (p *ProjectListPage) handleDeleteConfirmKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "left", "h":
		p.hoverButton = 0
	case "right", "l":
		p.hoverButton = 1
	case "y", "enter":
		if p.hoverButton == 1 {
			current := p.list.SelectedItem().(listItem)
			p.store.Remove(current.project.ID)
			p.syncListItems()
		}
		p.state = ProjectStateList
	case "n", "esc":
		p.state = ProjectStateList
	}
	return nil
}

func (p *ProjectListPage) handleIDEMenuKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down":
		p.ideMenu.selected = min(p.ideMenu.selected+1, len(p.ideMenu.options)-1)
	case "k", "up":
		p.ideMenu.selected = max(p.ideMenu.selected-1, 0)
	case "enter":
		selectedIDE := p.ideMenu.options[p.ideMenu.selected].Type
		return p.openWithIDE(models.IDEType(selectedIDE))
	case "esc":
		p.state = ProjectStateList
	}
	return nil
}

func (p *ProjectListPage) handleViewDetailKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		p.state = ProjectStateList
	case "up", "down", "pgup", "pgdown", "home", "end":
		var cmd tea.Cmd
		p.viewport, cmd = p.viewport.Update(msg)
		return cmd
	}
	return nil
}

func (p *ProjectListPage) handleEditDescKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+s":
		if current := p.safeGetSelectedProject(); current != nil {
			p.store.UpdateDescription(current.ID, p.ta.Value())
			p.syncListItems()
		}
		p.state = ProjectStateList
	case "esc":
		p.state = ProjectStateList
	default:
		var cmd tea.Cmd
		p.ta, cmd = p.ta.Update(msg)
		return cmd
	}
	return nil
}

func (p *ProjectListPage) openWithIDE(ideType models.IDEType) tea.Cmd {
	idx := p.list.Index()
	current := p.store.GetByIndex(idx)
	if current == nil {
		return func() tea.Msg {
			return errMsg{err: fmt.Errorf("未选择项目")}
		}
	}
	if !p.app.IDEExec().IsIDEAvailable(ideType) {
		p.app.ShowToast("IDE 不可用", string(ToastError))
		return nil
	}
	projectID := current.ID
	proj := *current
	return func() tea.Msg {
		done := make(chan error, 1)
		go func() {
			done <- p.app.IDEExec().OpenProject(&proj, ideType)
		}()

		select {
		case err := <-done:
			if err != nil {
				return errMsg{err: err}
			}
			return openProjectSuccessMsg{projectID: projectID}
		case <-time.After(30 * time.Second):
			return errMsg{err: fmt.Errorf("打开超时（30秒）")}
		}
	}
}

// ============== 辅助方法 ==============

func (p *ProjectListPage) renderHelpText() string {
	sep := lipgloss.NewStyle().Foreground(ui.PrimaryDim).Render("│")
	quit := lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("退出")

	return lipgloss.NewStyle().
		Foreground(ui.SecondaryText).
		MarginLeft(2).
		Render(
			lipgloss.JoinHorizontal(lipgloss.Left, sep,
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyNavStyle.Render("↑↓"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("导航"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeySearchStyle.Render("[/]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("搜索"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyActionStyle.Render("[n]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("添加"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyActionStyle.Render("[r]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("改名"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyActionStyle.Render("[e]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("描述"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyActionStyle.Render("[b]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("批量"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyActionStyle.Render("[p]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("配置"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyDangerStyle.Render("[d]"),
					lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("删除"),
				),
				lipgloss.JoinHorizontal(lipgloss.Left, " ",
					ui.HelpKeyQuitStyle.Render("[q]"),
					quit,
				),
			),
		)
}

func (p *ProjectListPage) syncListItems() {
	projects := p.store.Projects
	p.list.SetItems(newListItems(projects))
	// 列表变短后修正选中项，避免光标落在窗口外
	if len(projects) > 0 && p.list.Index() >= len(projects) {
		p.list.Select(len(projects) - 1)
	}
}

func (p *ProjectListPage) safeGetSelectedProject() *models.Project {
	if len(p.list.Items()) == 0 {
		return nil
	}
	item := p.list.SelectedItem()
	if item == nil {
		return nil
	}
	prj, ok := item.(listItem)
	if !ok {
		return nil
	}
	return &prj.project
}

func (p *ProjectListPage) updateViewport() {
	dialogWidth := max(50, int(float64(p.width)*0.8))
	p.viewport.Width = dialogWidth - 6
	p.viewport.Height = max(5, p.height-16)

	if len(p.list.Items()) == 0 {
		p.viewport.SetContent("")
		return
	}
	current := p.safeGetSelectedProject()
	if current == nil {
		p.viewport.SetContent("")
		return
	}
	proj := *current

	sectionDivider := ui.SeparatorHighlightStyle.Render("─── 详情 ───")

	infoLine := func(label, value string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(ui.ForegroundDim).Width(10).Render(label),
			lipgloss.NewStyle().Foreground(ui.Foreground).Render(value),
		)
	}

	var openCountBadge string
	if proj.OpenCount >= 10 {
		openCountBadge = ui.FeaturedBadgeStyle.Render(fmt.Sprintf("打开 %d 次", proj.OpenCount))
	} else {
		openCountBadge = ui.BadgeStyle.Foreground(ui.SuccessColor).Render(fmt.Sprintf("打开 %d 次", proj.OpenCount))
	}

	lastOpenedColor := ui.SecondaryText
	createdColor := ui.ForegroundDim
	if time.Since(proj.LastOpened) < 24*time.Hour {
		lastOpenedColor = ui.AccentCyan
	}

	lastOpenedText := proj.LastOpened.Format("2006-01-02 15:04:05")
	if proj.LastOpened.IsZero() {
		lastOpenedText = "从未打开"
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Foreground(ui.PrimaryColor).Bold(true).Render(proj.Alias),
		"",
		sectionDivider,
		"",
		infoLine("路径", proj.Path),
		infoLine("打开次数", openCountBadge),
		infoLine("创建时间", lipgloss.NewStyle().Foreground(createdColor).Render(proj.CreatedAt.Format("2006-01-02 15:04:05"))),
		infoLine("最后打开", lipgloss.NewStyle().Foreground(lastOpenedColor).Render(lastOpenedText)),
		"",
		lipgloss.NewStyle().Foreground(ui.SecondaryText).Bold(true).MarginTop(1).Render("描述"),
		"",
		lipgloss.NewStyle().Foreground(ui.Foreground).Render(proj.Description),
	)

	p.viewport.SetContent(content)
}

// errMsg 错误消息
type errMsg struct {
	err error
}

// openProjectSuccessMsg 项目打开成功消息
type openProjectSuccessMsg struct {
	projectID string
}

// generateID 生成唯一 ID
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// defaultAliasForPath 从路径中提取默认项目名称
func defaultAliasForPath(path string) string {
	alias := filepath.Base(path)
	if alias == "" || alias == "/" || alias == "\\" {
		return ""
	}
	return alias
}
