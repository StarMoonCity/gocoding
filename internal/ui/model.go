package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"gocoding/internal/commands"
	"gocoding/internal/models"
)

type IDEExecutor = commands.IDEExecutor

type Model struct {
	store            *models.ProjectStore
	providerStore    *models.ModelProviderStore
	ideExec          *IDEExecutor
	list             list.Model
	providerList     list.Model // 用于显示配置列表
	width            int
	height           int
	state            AppState
	dialog           *DialogModel
	input            textinput.Model
	secondaryInput   textinput.Model // 用于项目名称输入
	ideMenu          *IDEMenuModel
	ta               textarea.Model
	viewport         viewport.Model
	tempPath         string
	tempName         string
	inputFocus       InputFocus // 当前输入焦点
	errMsg           string
	tipMsg           string // 提示信息
	editingProjectID string // 编辑中的项目ID
	layoutMode       LayoutMode
	showDetails      bool
	searchQuery      string // 搜索查询字符串
	// Provider 配置输入
	providerInputFocus         ProviderInputFocus
	providerNameInput          textinput.Model
	providerBaseURLInput       textinput.Model
	providerAPIKeyInput        textinput.Model
	providerModelInput         textinput.Model
	providerThinkingModelInput textinput.Model // 推理模型
	providerDefaultHaikuInput  textinput.Model // Haiku 默认模型
	providerDefaultSonnetInput textinput.Model // Sonnet 默认模型
	providerDefaultOpusInput   textinput.Model // Opus 默认模型
	editingProviderID          string          // 编辑中的配置ID，为空表示新增
	itemIndexCache             map[string]int  // listItem.FilterValue -> index 缓存
	// 调试模式
	debug   bool
	lastKey string
	// 批量添加项目
	batchProjects []string     // ~/.claude/projects 下的项目路径列表
	batchSelected map[int]bool // 选中的项目索引
	batchCursor   int          // 批量选择列表的当前光标位置
}

type LayoutMode int

const (
	LayoutSingle LayoutMode = iota // 单列布局
)

// InputFocus 输入框焦点状态
type InputFocus int

const (
	FocusPath InputFocus = iota
	FocusName
)

// ProviderInputFocus 模型提供商配置输入焦点
type ProviderInputFocus int

const (
	FocusProviderName ProviderInputFocus = iota
	FocusProviderBaseURL
	FocusProviderAPIKey
	FocusProviderModel
	FocusProviderThinkingModel
	FocusProviderDefaultHaiku
	FocusProviderDefaultSonnet
	FocusProviderDefaultOpus
	FocusProviderCount // 焦点数量
)

type AppState int

const (
	StateList AppState = iota
	StateAddProject
	StateRenameProject
	StateDeleteConfirm
	StateIDEMenu
	StateViewDetail
	StateEditDescription
	StateSearch
	StateProviderList    // 模型配置列表
	StateProviderAdd     // 添加新配置
	StateProviderEdit    // 编辑配置
	StateProviderDelete  // 删除确认
	StateBatchAddProject // 批量添加项目
)

type DialogModel struct {
	title    string
	message  string
	buttons  []string
	selected int
}

type IDEMenuModel struct {
	title     string
	options   []IDEOption
	selected  int
	available map[models.IDEType]bool
}

type IDEOption struct {
	Type        models.IDEType
	Name        string
	Description string
}

type listItem struct {
	project models.Project
}

type providerListItem struct {
	provider models.ModelProvider
}

func (i providerListItem) Title() string {
	return i.provider.Name
}

func (i providerListItem) Description() string {
	if i.provider.Active {
		return i.provider.BaseURL + " • " + i.provider.Model + " [激活]"
	}
	return i.provider.BaseURL + " • " + i.provider.Model
}

func (i providerListItem) FilterValue() string { return i.provider.Name }

func (i listItem) Title() string {
	// 只显示别名
	return i.project.Alias
}

func (i listItem) Description() string {
	// 在描述中显示路径（会换行）
	return i.project.Path
}

func (i listItem) FilterValue() string { return i.project.Alias }

// newListItems converts a slice of projects to list items
func newListItems(projects []models.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = listItem{project: p}
	}
	return items
}

// newProviderListItems converts a slice of providers to list items
func newProviderListItems(providers []models.ModelProvider) []list.Item {
	items := make([]list.Item, len(providers))
	for i, p := range providers {
		items[i] = providerListItem{provider: p}
	}
	return items
}

func NewModel(store *models.ProjectStore) *Model {
	ideExec := commands.NewIDEExecutor()
	m := &Model{
		store:       store,
		ideExec:     ideExec,
		state:       StateList,
		searchQuery: "",
	}

	items := newListItems(store.Projects)

	delegate := projectListDelegate{}

	// 初始大小，后续 SetSize 会重新设置
	m.list = list.New(items, delegate, 60, 14)
	m.list.Title = ""
	m.list.Styles.Title = lipgloss.NewStyle().Foreground(SecondaryText).Bold(true).MarginBottom(1)
	m.list.SetShowTitle(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(false) // 禁用列表内置帮助，使用自定义帮助文本
	m.list.SetFilteringEnabled(true)

	m.dialog = &DialogModel{}
	m.ideMenu = &IDEMenuModel{
		options: []IDEOption{
			{Type: models.IDEClaudeCode, Name: "Claude", Description: "Claude Code IDE"},
			{Type: models.IDEVSCode, Name: "VSCode", Description: "Visual Studio Code"},
			{Type: models.IDEOpenCode, Name: "OpenCode", Description: "OpenCode IDE"},
			{Type: models.IDECodexCLI, Name: "Codex", Description: "Codex CLI"},
		},
		available: make(map[models.IDEType]bool),
	}

	// 初始化单行输入框
	m.input = textinput.New()
	m.input.Placeholder = "输入项目路径"
	m.secondaryInput = textinput.New()
	m.secondaryInput.Placeholder = "输入项目名称"

	// 初始化文本区域
	m.ta = textarea.New()
	m.ta.Focus()

	// 初始化Viewport
	m.viewport = viewport.New(0, 0)

	// 自动定位到最近打开的项目
	if store.Len() > 0 {
		recent := store.GetMostRecentlyOpened()
		if recent != nil {
			index := store.GetIndexByProject(recent.ID)
			if index >= 0 {
				m.list.Select(index)
			}
		}
	}

	// 初始化模型配置输入框
	m.initProviderInputs()

	// 初始化模型配置列表
	m.initProviderList()

	return m
}

// initProviderInputs 初始化模型配置输入框
func (m *Model) initProviderInputs() {
	m.providerNameInput = textinput.New()
	m.providerNameInput.Placeholder = "配置名称 (如 MiniMax)"

	m.providerBaseURLInput = textinput.New()
	m.providerBaseURLInput.Placeholder = "Base URL (如 https://api.minimax.chat)"

	m.providerAPIKeyInput = textinput.New()
	m.providerAPIKeyInput.Placeholder = "API Key"

	m.providerModelInput = textinput.New()
	m.providerModelInput.Placeholder = "主模型 (如 MiniMax-M2.7-highspeed)"

	m.providerThinkingModelInput = textinput.New()
	m.providerThinkingModelInput.Placeholder = "推理模型 (如 MiniMax-M2.7-highspeed)"

	m.providerDefaultHaikuInput = textinput.New()
	m.providerDefaultHaikuInput.Placeholder = "Haiku 默认模型"

	m.providerDefaultSonnetInput = textinput.New()
	m.providerDefaultSonnetInput.Placeholder = "Sonnet 默认模型"

	m.providerDefaultOpusInput = textinput.New()
	m.providerDefaultOpusInput.Placeholder = "Opus 默认模型"
}

// initProviderList 初始化模型配置列表
func (m *Model) initProviderList() {
	delegate := providerListDelegate{}
	m.providerList = list.New(nil, delegate, 0, 0)
	m.providerList.SetShowTitle(false)
	m.providerList.SetShowStatusBar(false)
	m.providerList.SetShowHelp(false)
	m.providerList.SetFilteringEnabled(true)
}

// projectListDelegate 项目列表自定义渲染
type projectListDelegate struct{}

func (d projectListDelegate) Height() int                               { return 1 }
func (d projectListDelegate) Spacing() int                              { return 1 }
func (d projectListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d projectListDelegate) Render(w io.Writer, m list.Model, width int, item list.Item) {
	proj, ok := item.(listItem)
	if !ok {
		return
	}

	isSelected := item == m.SelectedItem()

	// 选中标记
	selector := "  "
	selectorStyle := lipgloss.NewStyle().Foreground(MutedText)
	if isSelected {
		selector = "▸ "
		selectorStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	}

	// 名称
	nameStyle := lipgloss.NewStyle().Bold(true)
	if isSelected {
		nameStyle = nameStyle.Foreground(PrimaryColor)
	} else {
		nameStyle = nameStyle.Foreground(Foreground)
	}

	// 打开次数徽章
	countBadge := lipgloss.NewStyle().
		Foreground(SecondaryText).
		Background(BackgroundLight).
		Padding(0, 1).
		MarginLeft(1).
		Render(fmt.Sprintf("×%d", proj.project.OpenCount))

	// 组合内容
	content := lipgloss.JoinHorizontal(
		lipgloss.Left,
		selectorStyle.Render(selector),
		nameStyle.Render(proj.project.Alias),
		countBadge,
	)

	fmt.Fprintf(w, "%s", content)
}

// providerListDelegate 自定义列表项渲染
type providerListDelegate struct{}

func (d providerListDelegate) Height() int                               { return 2 }
func (d providerListDelegate) Spacing() int                              { return 1 }
func (d providerListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d providerListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	p, ok := item.(providerListItem)
	if !ok {
		return
	}

	provider := p.provider
	isSelected := index == m.Index()
	rowWidth := max(20, m.Width())

	// 激活标签样式
	activeTag := ""
	if provider.Active {
		activeTag = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Render("● 激活")
	}

	// 名称样式 - 选中时使用主色高亮
	nameStyle := lipgloss.NewStyle().Bold(true)
	if isSelected {
		nameStyle = nameStyle.Foreground(PrimaryColor)
	} else {
		nameStyle = nameStyle.Foreground(Foreground)
	}

	// URL 和模型样式
	infoStyle := lipgloss.NewStyle().Foreground(SecondaryText)

	// 选中标记
	selector := "  "
	selectorStyle := lipgloss.NewStyle().Foreground(MutedText)
	if isSelected {
		selector = "▸ "
		selectorStyle = lipgloss.NewStyle().Foreground(PrimaryColor)
	}

	selectorWidth := lipgloss.Width(selector)
	activeWidth := lipgloss.Width(activeTag)
	nameWidth := rowWidth - selectorWidth
	if activeWidth > 0 {
		nameWidth -= activeWidth + 1
	}
	nameWidth = max(8, nameWidth)

	name := ansi.Truncate(provider.Name, nameWidth, "…")
	line := selectorStyle.Render(selector) + nameStyle.Render(name)
	if activeWidth > 0 {
		gapWidth := rowWidth - selectorWidth - lipgloss.Width(name) - activeWidth
		if gapWidth < 1 {
			name = ansi.Truncate(provider.Name, max(6, nameWidth-1+gapWidth), "…")
			line = selectorStyle.Render(selector) + nameStyle.Render(name)
			gapWidth = rowWidth - selectorWidth - lipgloss.Width(name) - activeWidth
		}
		line = lipgloss.JoinHorizontal(
			lipgloss.Left,
			line,
			strings.Repeat(" ", max(1, gapWidth)),
			activeTag,
		)
	}
	line = lipgloss.NewStyle().Width(rowWidth).Render(line)

	// 第二行：URL 和模型
	infoText := provider.BaseURL
	if provider.Model != "" {
		infoText += " • " + provider.Model
	}
	infoPrefix := strings.Repeat(" ", selectorWidth)
	infoWidth := max(8, rowWidth-selectorWidth)
	secondLine := infoPrefix + infoStyle.Render(ansi.Truncate(infoText, infoWidth, "…"))

	// 使用换行符连接
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		line,
		secondLine,
	)

	fmt.Fprintf(w, "%s", content)
}

// SetProviderStore 设置模型配置存储
func (m *Model) SetProviderStore(store *models.ModelProviderStore) {
	m.providerStore = store
	// 更新配置列表
	m.updateProviderListItems()
}

// updateProviderListItems 更新配置列表
func (m *Model) updateProviderListItems() {
	if m.providerStore == nil {
		return
	}
	m.providerList.SetItems(newProviderListItems(m.providerStore.Providers))
}

// syncListItems 同步项目列表项
func (m *Model) syncListItems() {
	m.list.SetItems(newListItems(m.store.Projects))
}

// safeGetSelectedProject 安全获取选中的项目
func (m *Model) safeGetSelectedProject() *models.Project {
	if len(m.list.Items()) == 0 {
		return nil
	}
	item := m.list.SelectedItem()
	if item == nil {
		return nil
	}
	p, ok := item.(listItem)
	if !ok {
		return nil
	}
	return &p.project
}

// safeGetSelectedProvider 安全获取选中的提供商
func (m *Model) safeGetSelectedProvider() *models.ModelProvider {
	if len(m.providerList.Items()) == 0 {
		return nil
	}
	item := m.providerList.SelectedItem()
	if item == nil {
		return nil
	}
	p, ok := item.(providerListItem)
	if !ok {
		return nil
	}
	return &p.provider
}

// SetState 设置应用状态
func (m *Model) SetState(state AppState) {
	m.state = state
}

// SetDebugMode 设置调试模式
func (m *Model) SetDebugMode(enabled bool) {
	m.debug = enabled
}

// IsDebugMode 检查是否开启调试模式
func (m *Model) IsDebugMode() bool {
	return m.debug
}

func (m *Model) Init() tea.Cmd {
	return nil
}
