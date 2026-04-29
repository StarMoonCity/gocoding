package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/commands"
	"gocoding/internal/models"
)

// AppModel 路由中枢 - 精简版只负责组合和路由
type AppModel struct {
	store         *models.ProjectStore
	providerStore *models.ModelProviderStore
	ideExec       *commands.IDEExecutor

	// 页面实例
	projectPage  *ProjectListPage
	providerPage *ProviderListPage
	searchPage   *SearchPage
	batchAddPage *BatchAddPage

	// 当前活动页面
	currentPage PageModel

	// 模态框管理
	modalManager ModalManager

	// Toast 管理
	toastManager ToastManager

	// 状态栏
	statusBar StatusBar

	// 尺寸
	width  int
	height int

	// 调试模式
	debug   bool
	lastKey string
}

// NewAppModel 创建 AppModel
func NewAppModel(store *models.ProjectStore, providerStore *models.ModelProviderStore) *AppModel {
	ideExec := commands.NewIDEExecutor()

	m := &AppModel{
		store:         store,
		providerStore: providerStore,
		ideExec:       ideExec,
	}

	// 初始化页面
	m.projectPage = NewProjectListPage(store)
	m.providerPage = NewProviderListPage(providerStore)
	m.searchPage = NewSearchPage(store)
	m.batchAddPage = NewBatchAddPage(store)

	// 设置页面引用
	m.projectPage.SetApp(m)
	m.providerPage.SetApp(m)
	m.searchPage.SetApp(m)
	m.batchAddPage.SetApp(m)

	// 默认显示项目列表
	m.currentPage = m.projectPage

	return m
}

// SetSize 设置终端尺寸
func (m *AppModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	// 同步到所有页面
	m.projectPage.SetSize(width, height)
	m.providerPage.SetSize(width, height)
	m.searchPage.SetSize(width, height)
	m.batchAddPage.SetSize(width, height)
}

// SetDebugMode 设置调试模式
func (m *AppModel) SetDebugMode(enabled bool) {
	m.debug = enabled
}

// Init 初始化
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update 更新逻辑 - 路由到当前页面或模态框
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		m.lastKey = msg.String()

		// 先尝试模态框
		if m.modalManager.HasModal() {
			modal := m.modalManager.Top()
			if modal != nil {
				if cmd, consumed := modal.Update(msg); consumed {
					return m, cmd
				}
			}
		}

		// 全局快捷键
		switch msg.String() {
		case "ctrl+c", "ctrl+q", "q":
			return m, tea.Quit
		}

		// 路由到当前页面
		if m.currentPage != nil {
			if cmd, consumed := m.currentPage.Update(msg); consumed {
				return m, cmd
			}
		}

	case tea.MouseMsg:
		if m.currentPage != nil {
			m.currentPage.HandleMouse(msg)
		}
	}

	// 更新 Toast 管理器
	if cmd := m.toastManager.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View 渲染视图
func (m *AppModel) View() string {
	var content string

	// 渲染当前页面
	if m.currentPage != nil {
		content = m.currentPage.View(m.width, m.height)
	}

	// 叠加模态框
	if m.modalManager.HasModal() {
		content = m.modalManager.Overlay(content, m.width, m.height)
	}

	// 叠加 Toast
	toastView := m.toastManager.View(m.width)
	if toastView != "" {
		content += "\n" + toastView
	}

	// 渲染状态栏
	statusView := m.statusBar.View(m.width)
	if statusView != "" {
		content += "\n" + statusView
	}

	// 调试面板
	if m.debug {
		content += "\n" + m.renderDebugPanel()
	}

	return content
}

// SwitchPage 切换页面
func (m *AppModel) SwitchPage(page PageModel) {
	m.currentPage = page
}

// PushModal 推送模态框
func (m *AppModel) PushModal(modal Modal) {
	m.modalManager.Push(modal)
}

// PopModal 弹出模态框
func (m *AppModel) PopModal() {
	m.modalManager.Pop()
}

// CloseModal 关闭所有模态框
func (m *AppModel) CloseModal() {
	m.modalManager.Close()
}

// ShowToast 显示 Toast 消息
func (m *AppModel) ShowToast(message string, toastType string) {
	m.toastManager.Show(message, toastType, 3*time.Second)
}

// SetStatusBar 设置状态栏内容
func (m *AppModel) SetStatusBar(left, center, right string) {
	m.statusBar.Set(left, center, right)
}

// Store 访问器
func (m *AppModel) Store() *models.ProjectStore {
	return m.store
}

func (m *AppModel) ProviderStore() *models.ModelProviderStore {
	return m.providerStore
}

func (m *AppModel) IDEExec() *commands.IDEExecutor {
	return m.ideExec
}

// ProviderPage 返回 providerPage
func (m *AppModel) ProviderPage() *ProviderListPage {
	return m.providerPage
}

// SearchPage 返回 searchPage
func (m *AppModel) SearchPage() *SearchPage {
	return m.searchPage
}

// ProjectPage 返回 projectPage
func (m *AppModel) ProjectPage() *ProjectListPage {
	return m.projectPage
}

// BatchAddPage 返回 batchAddPage
func (m *AppModel) BatchAddPage() *BatchAddPage {
	return m.batchAddPage
}

// renderDebugPanel 渲染调试面板
func (m *AppModel) renderDebugPanel() string {
	pageName := "nil"
	if m.currentPage != nil {
		pageName = pageTypeName(m.currentPage.PageType())
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666")).
		Render(fmt.Sprintf("Page: %s | Size: %dx%d | Key: %s", pageName, m.width, m.height, m.lastKey))
}

func pageTypeName(t PageType) string {
	switch t {
	case PageProject:
		return "Project"
	case PageProvider:
		return "Provider"
	case PageSearch:
		return "Search"
	case PageBatchAdd:
		return "BatchAdd"
	default:
		return "Unknown"
	}
}
