package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// LayoutConfig 布局配置
type LayoutConfig struct {
	listWidth    int
	listHeight   int
	columnCount  int
	columnGap    int
	paddingX     int
	helpTextMode HelpTextMode
}

// HelpTextMode 帮助文本模式
type HelpTextMode int

const (
	HelpTextCompact HelpTextMode = iota // 精简版: j/k ↑/↓
	HelpTextNormal                      // 标准版: ↑/↓ j/k 导航
	HelpTextFull                        // 完整版: 所有快捷键
)

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height

	// 计算布局配置
	config := m.calculateLayout(width, height-m.debugPanelHeight())
	m.layoutMode = LayoutSingle

	// 动态设置列表大小，让列表占满可用高度
	m.list.SetSize(config.listWidth, config.listHeight)

	// 设置模型配置列表大小（预留更多行给header、colTitle、helpText等）
	providerWidth := m.providerListWidth()
	m.providerList.SetSize(providerWidth, config.listHeight-6)

	// 动态计算 provider 输入框宽度
	providerInputWidth := max(20, m.providerDialogWidth()-14)
	m.providerNameInput.Width = providerInputWidth
	m.providerBaseURLInput.Width = providerInputWidth
	m.providerAPIKeyInput.Width = providerInputWidth
	m.providerModelInput.Width = providerInputWidth
	m.providerThinkingModelInput.Width = providerInputWidth
	m.providerDefaultHaikuInput.Width = providerInputWidth
	m.providerDefaultSonnetInput.Width = providerInputWidth
	m.providerDefaultOpusInput.Width = providerInputWidth
	m.providerSubagentInput.Width = providerInputWidth
	m.providerNonessentialInput.Width = providerInputWidth
	m.providerNonstreamingInput.Width = providerInputWidth
	m.providerEffortInput.Width = providerInputWidth

	// 设置文本区域大小 - 使用可用宽度
	m.ta.SetWidth(width - 10)
	m.ta.SetHeight(5)

	// 设置Viewport大小
	detailWidth := width - config.paddingX*2 - config.listWidth - config.columnGap - 4
	detailHeight := config.listHeight
	if detailWidth > 0 {
		m.viewport.Width = detailWidth
		m.viewport.Height = detailHeight
	}
}

func (m *Model) debugPanelHeight() int {
	if m.debug {
		return 1
	}
	return 0
}

func (m *Model) providerDialogWidth() int {
	// 使用屏幕宽度的 80%，但限制在 46-72 之间
	suggestedWidth := int(float64(m.width) * 0.8)
	return min(72, max(46, suggestedWidth))
}

func (m *Model) providerListWidth() int {
	// 使用对话框宽度的 90%
	return max(30, m.providerDialogWidth()*9/10)
}

func (m *Model) calculateLayout(width, height int) LayoutConfig {
	config := LayoutConfig{
		columnGap:   2,
		paddingX:    0,
		columnCount: 1,
	}

	// 列表宽度使用整个屏幕宽度
	config.listWidth = width

	// 确保最小宽度
	if config.listWidth < 30 {
		config.listWidth = 30
	}

	// 根据屏幕宽度确定帮助文本模式
	switch {
	case width < 50:
		config.helpTextMode = HelpTextCompact
	case width < 80:
		config.helpTextMode = HelpTextNormal
	default:
		config.helpTextMode = HelpTextFull
	}

	// 计算列表高度（预留头部、帮助文本和错误消息的空间）
	// 动态计算，留出空间给头部(3行)、帮助文本(2行)、错误消息(1行)、底部边距(1行)
	headerHeight := 3
	helpHeight := 2
	errorHeight := 1
	marginHeight := 1
	 reservedHeight := headerHeight + helpHeight + errorHeight + marginHeight

	minHeight := 8
	listHeight := height - reservedHeight
	if listHeight < minHeight {
		listHeight = minHeight
	}
	config.listHeight = listHeight

	return config
}

// getSelectedItemName 获取当前选中的项目或配置名称
func (m *Model) getSelectedItemName() string {
	switch m.state {
	case StateList, StateAddProject, StateRenameProject, StateDeleteConfirm,
		StateIDEMenu, StateViewDetail, StateEditDescription, StateSearch, StateBatchAddProject:
		if proj := m.safeGetSelectedProject(); proj != nil {
			return proj.Alias
		}
	case StateProviderList, StateProviderAdd, StateProviderEdit, StateProviderDelete:
		if provider := m.safeGetSelectedProvider(); provider != nil {
			return provider.Name
		}
	}
	return ""
}

// renderDebugPanel 渲染调试面板
func (m *Model) renderDebugPanel() string {
	debugStyle := lipgloss.NewStyle().
		Foreground(MutedText).
		Background(BackgroundDeep).
		Width(max(0, m.width-4)).
		Padding(0, 2)

	stateName := m.stateName()
	itemName := m.getSelectedItemName()

	info := fmt.Sprintf("State: %s | Item: %s | Key: %s", stateName, itemName, m.lastKey)
	return debugStyle.Render(info)
}

// stateName 获取状态的字符串名称
func (m *Model) stateName() string {
	switch m.state {
	case StateList:
		return "StateList"
	case StateAddProject:
		return "StateAddProject"
	case StateRenameProject:
		return "StateRenameProject"
	case StateDeleteConfirm:
		return "StateDeleteConfirm"
	case StateIDEMenu:
		return "StateIDEMenu"
	case StateViewDetail:
		return "StateViewDetail"
	case StateEditDescription:
		return "StateEditDescription"
	case StateSearch:
		return "StateSearch"
	case StateProviderList:
		return "StateProviderList"
	case StateProviderAdd:
		return "StateProviderAdd"
	case StateProviderEdit:
		return "StateProviderEdit"
	case StateProviderDelete:
		return "StateProviderDelete"
	case StateBatchAddProject:
		return "StateBatchAddProject"
	default:
		return "Unknown"
	}
}
