package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/commands"
	"gocoding/internal/models"
)

type IDEExecutor = commands.IDEExecutor

type Model struct {
	store          *models.ProjectStore
	providerStore  *models.ModelProviderStore
	ideExec        *IDEExecutor
	list           list.Model
	providerList   list.Model // 用于显示配置列表
	width          int
	height         int
	state          AppState
	dialog         *DialogModel
	input          textinput.Model
	secondaryInput textinput.Model // 用于项目名称输入
	ideMenu        *IDEMenuModel
	ta             textarea.Model
	viewport       viewport.Model
	tempPath       string
	tempName       string
	inputFocus     InputFocus // 当前输入焦点
	errMsg         string
	layoutMode     LayoutMode
	showDetails    bool
	searchQuery    string // 搜索查询字符串
	// Provider 配置输入
	providerInputFocus      ProviderInputFocus
	providerNameInput      textinput.Model
	providerBaseURLInput   textinput.Model
	providerAPIKeyInput    textinput.Model
	providerModelInput     textinput.Model
	editingProviderID      string // 编辑中的配置ID，为空表示新增
}

type LayoutMode int

const (
	LayoutSingle LayoutMode = iota // 单列布局
	LayoutDouble                   // 双列布局
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
	StateProviderList  // 模型配置列表
	StateProviderAdd   // 添加新配置
	StateProviderEdit  // 编辑配置
	StateProviderDelete // 删除确认
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

func NewModel(store *models.ProjectStore) *Model {
	ideExec := commands.NewIDEExecutor()
	m := &Model{
		store:        store,
		ideExec:      ideExec,
		state:        StateList,
		searchQuery:  "",
	}

	items := make([]list.Item, len(store.Projects))
	for i, p := range store.Projects {
		items[i] = listItem{project: p}
	}

	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0) // 紧凑间距

	// 初始大小，后续 SetSize 会重新设置
	m.list = list.New(items, delegate, 60, 14)
	m.list.Title = ""
	m.list.Styles.Title = lipgloss.NewStyle().Foreground(SecondaryColor).Bold(true).MarginBottom(1)
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
	m.providerModelInput.Placeholder = "模型名称 (如 MiniMax-M2.7-highspeed)"
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
	items := make([]list.Item, len(m.providerStore.Providers))
	for i, p := range m.providerStore.Providers {
		items[i] = providerListItem{provider: p}
	}
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(0)
	m.providerList = list.New(items, delegate, 60, 14)
	m.providerList.SetShowTitle(false)
	m.providerList.SetShowStatusBar(false)
	m.providerList.SetShowHelp(false)
	m.providerList.SetFilteringEnabled(true)
}

// SetState 设置应用状态
func (m *Model) SetState(state AppState) {
	m.state = state
}

func (m *Model) Init() tea.Cmd {
	return nil
}
