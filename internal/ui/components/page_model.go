package components

import "github.com/charmbracelet/bubbletea"

// PageType 页面类型枚举
type PageType int

const (
	PageProject PageType = iota
	PageProvider
	PageSearch
	PageBatchAdd
)

// InputFocus 输入框焦点状态
type InputFocus int

const (
	FocusPath InputFocus = iota
	FocusName
)

// PageModel 页面模型接口 - 企业级架构基础
type PageModel interface {
	// Update 处理消息，返回是否已处理
	Update(msg tea.Msg) (tea.Cmd, bool)
	// View 渲染页面
	View(width, height int) string
	// OnActivate 页面激活时调用
	OnActivate()
	// OnDeactivate 页面停用时调用
	OnDeactivate()
	// HandleMouse 处理鼠标消息
	HandleMouse(msg tea.MouseMsg)
	// PageType 返回页面类型
	PageType() PageType
}
