package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/commands"
	"gocoding/internal/models"
	"gocoding/internal/ui"
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

	// Toast 管理
	toastManager ToastManager

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

	// Toast 管理器先处理过期消息，避免被页面消费导致 toast 不消失
	if cmd := m.toastManager.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, tea.Batch(cmds...)

	case tea.KeyMsg:
		m.lastKey = msg.String()

		// 全局快捷键
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		}

		// 路由到当前页面
		if m.currentPage != nil {
			if cmd, consumed := m.currentPage.Update(msg); consumed {
				cmds = append(cmds, cmd)
				// 页面处理可能触发 ShowToast，这里补一次计时
				if tick := m.toastManager.Update(msg); tick != nil {
					cmds = append(cmds, tick)
				}
				return m, tea.Batch(cmds...)
			}
		}

	case openProjectSuccessMsg:
		// 打开 IDE 成功后由项目页刷新列表与状态
		if cmd, consumed := m.projectPage.Update(msg); consumed {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case searchOpenedMsg:
		// 打开 IDE 成功后由搜索页刷新列表与状态
		if cmd, consumed := m.searchPage.Update(msg); consumed {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case errMsg:
		// 展示 IDE 打开等异步操作产生的错误
		m.ShowToast(msg.err.Error(), string(ToastError))

	case tea.MouseMsg:
		if m.currentPage != nil {
			m.currentPage.HandleMouse(msg)
		}
	}

	// 更新 Toast 管理器（页面处理可能触发 ShowToast，这里补一次计时）
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

	// 叠加 Toast
	toastView := m.toastManager.View(m.width)
	if toastView != "" {
		content += "\n" + toastView
	}

	// 调试面板
	if m.debug {
		content += "\n" + m.renderDebugPanel()
	}

	// 全屏背景：逐行填充确保无透明区域
	if m.width > 0 && m.height > 0 {
		content = ui.FillBackground(content, m.width, m.height)
	}

	return content
}

// SwitchPage 切换页面
func (m *AppModel) SwitchPage(page PageModel) {
	m.currentPage = page
}

// ShowToast 显示 Toast 消息
func (m *AppModel) ShowToast(message string, toastType string) {
	m.toastManager.Show(message, toastType, 3*time.Second)
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
		Foreground(ui.MutedText).
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
