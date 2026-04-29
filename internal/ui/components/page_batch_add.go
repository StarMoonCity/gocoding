package components

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
	"gocoding/internal/ui"
)

// BatchAddPage 批量添加项目页面
type BatchAddPage struct {
	store     *models.ProjectStore
	app       *AppModel
	width     int
	height    int
	projects  []string
	selected  map[int]bool
	cursor    int
}

// NewBatchAddPage 创建批量添加页面
func NewBatchAddPage(store *models.ProjectStore) *BatchAddPage {
	p := &BatchAddPage{
		store:    store,
		selected: make(map[int]bool),
		cursor:   0,
	}

	p.loadProjects()

	return p
}

// PageType 返回页面类型
func (p *BatchAddPage) PageType() PageType {
	return PageBatchAdd
}

// OnActivate 页面激活时调用
func (p *BatchAddPage) OnActivate() {
	p.loadProjects()
}

// OnDeactivate 页面停用时调用
func (p *BatchAddPage) OnDeactivate() {
}

// SetApp 设置 AppModel 引用
func (p *BatchAddPage) SetApp(app *AppModel) {
	p.app = app
}

// SetSize 设置尺寸
func (p *BatchAddPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// Update 处理消息
func (p *BatchAddPage) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			// 切换选中状态
			if p.cursor >= 0 && p.cursor < len(p.projects) {
				p.selected[p.cursor] = !p.selected[p.cursor]
			}
			return nil, true
		case "up", "k":
			if p.cursor > 0 {
				p.cursor--
			}
			return nil, true
		case "down", "j":
			if p.cursor < len(p.projects)-1 {
				p.cursor++
			}
			return nil, true
		case "enter":
			// 确认添加选中的项目
			p.addSelectedProjects()
			return nil, true
		case "esc":
			// 取消，返回项目列表
			if p.app != nil {
				p.app.SwitchPage(p.app.projectPage)
				p.app.projectPage.OnActivate()
			}
			return nil, true
		case "ctrl+c", "ctrl+q":
			return tea.Quit, false
		}
	}
	return nil, false
}

// View 渲染页面
func (p *BatchAddPage) View(width, height int) string {
	p.width = width
	p.height = height

	dialogWidth := min(60, max(45, int(float64(p.width)*0.75)))

	var items []string
	for i, path := range p.projects {
		isSelected := p.selected[i]
		cursor := "  "
		isHovered := i == p.hoverIndex()

		if i == p.cursor {
			cursor = lipgloss.NewStyle().Foreground(ui.PrimaryColor).Render("▸ ")
		} else if isHovered {
			cursor = "  "
		}
		checkbox := "[ ]"
		if isSelected {
			checkbox = lipgloss.NewStyle().Foreground(ui.SuccessColor).Render("[×]")
		}

		itemText := cursor + checkbox + "  " + path
		if isHovered && i != p.cursor {
			itemStyle := lipgloss.NewStyle().
				Foreground(ui.PrimaryDim).
				Background(ui.BackgroundHover)
			itemText = itemStyle.Render(itemText)
		}
		items = append(items, itemText)
	}

	if len(items) == 0 {
		items = append(items, lipgloss.NewStyle().Foreground(ui.WarningColor).Render("  没有找到可添加的项目"))
	}

	selectedCount := p.filterSelectedCount()
	var statusText string
	if selectedCount > 0 {
		statusText = ui.FeaturedBadgeStyle.Render(fmt.Sprintf("已选择: %d/%d", selectedCount, len(p.projects)))
	} else {
		statusText = lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render(fmt.Sprintf("已选择: %d/%d", selectedCount, len(p.projects)))
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.AccentGold).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.AccentGold).Bold(true).Render("＋ 批量添加项目"),
				"",
				lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render("~/.claude/projects"),
				"",
				lipgloss.JoinVertical(lipgloss.Left, items...),
				"",
				statusText,
				"",
				lipgloss.NewStyle().
					Foreground(ui.SecondaryText).
					Render("[Space/点击] 选择  ·  [↑↓] 移动  ·  [Enter] 确认  ·  [Esc] 取消"),
			),
		)

	// 上下左右居中
	return lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, dialog)
}

// HandleMouse 处理鼠标消息
func (p *BatchAddPage) HandleMouse(msg tea.MouseMsg) {
	listStartY := 4
	listEndY := p.height - 7

	switch {
	case msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress:
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := msg.Y - listStartY
			if clickedIndex >= 0 && clickedIndex < len(p.projects) {
				p.cursor = clickedIndex
				p.selected[clickedIndex] = !p.selected[clickedIndex]
			}
		}
	case msg.Button == tea.MouseButtonWheelUp:
		if p.cursor > 0 {
			p.cursor--
		}
	case msg.Button == tea.MouseButtonWheelDown:
		if p.cursor < len(p.projects)-1 {
			p.cursor++
		}
	}
}

// hoverIndex 获取悬停索引
func (p *BatchAddPage) hoverIndex() int {
	listStartY := 4
	listEndY := p.height - 7
	if p.cursor >= listStartY && p.cursor < listEndY {
		return p.cursor - listStartY
	}
	return -1
}

// filterSelectedCount 返回选中的项目数量
func (p *BatchAddPage) filterSelectedCount() int {
	count := 0
	for _, isSelected := range p.selected {
		if isSelected {
			count++
		}
	}
	return count
}

// loadProjects 加载项目列表
func (p *BatchAddPage) loadProjects() {
	claudeProjectsDir := os.ExpandEnv("$HOME/.claude/projects")

	entries, err := os.ReadDir(claudeProjectsDir)
	if err != nil {
		p.projects = nil
		return
	}

	p.projects = nil
	p.selected = make(map[int]bool)
	p.cursor = 0

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
		if p.isAlreadyAdded(sessionIndex.OriginalPath) {
			continue
		}

		p.projects = append(p.projects, sessionIndex.OriginalPath)
	}
}

// isAlreadyAdded 检查项目是否已经添加过
func (p *BatchAddPage) isAlreadyAdded(path string) bool {
	for _, proj := range p.store.Projects {
		if proj.Path == path {
			return true
		}
	}
	return false
}

// addSelectedProjects 添加选中的项目
func (p *BatchAddPage) addSelectedProjects() {
	count := 0
	for i, isSelected := range p.selected {
		if isSelected && i < len(p.projects) {
			path := p.projects[i]
			alias := filepath.Base(path)
			if alias == "" || alias == "/" || alias == "\\" {
				alias = path
			}
			project := models.Project{
				ID:        generateID(),
				Path:      path,
				Alias:     alias,
			}
			p.store.Add(project)
			count++
		}
	}

	if count > 0 && p.app != nil {
		p.app.ShowToast(fmt.Sprintf("已添加 %d 个项目", count), "success")
		p.app.SwitchPage(p.app.projectPage)
		p.app.projectPage.OnActivate()
	}
}
