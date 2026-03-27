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

	m.list.SetSize(config.listWidth, config.listHeight)

	// 设置模型配置列表大小
	m.providerList.SetSize(m.providerListWidth(), config.listHeight)

	providerInputWidth := max(20, m.providerDialogWidth()-14)
	m.providerNameInput.Width = providerInputWidth
	m.providerBaseURLInput.Width = providerInputWidth
	m.providerAPIKeyInput.Width = providerInputWidth
	m.providerModelInput.Width = providerInputWidth
	m.providerThinkingModelInput.Width = providerInputWidth
	m.providerDefaultHaikuInput.Width = providerInputWidth
	m.providerDefaultSonnetInput.Width = providerInputWidth
	m.providerDefaultOpusInput.Width = providerInputWidth

	// 设置文本区域大小
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
	return min(72, max(46, m.width-10))
}

func (m *Model) providerListWidth() int {
	return max(30, m.providerDialogWidth()-4)
}

func (m *Model) calculateLayout(width, height int) LayoutConfig {
	config := LayoutConfig{
		columnGap:   2,
		paddingX:    0,
		columnCount: 1, // 始终使用单列布局
	}

	// 列表宽度使用整个屏幕宽度
	config.listWidth = width

	// 确保最小宽度
	if config.listWidth < 30 {
		config.listWidth = 30
	}

	// 根据屏幕宽度确定帮助文本模式
	switch {
	case width < 60:
		config.helpTextMode = HelpTextCompact
	case width < 100:
		config.helpTextMode = HelpTextNormal
	default:
		config.helpTextMode = HelpTextFull
	}

	// 计算列表高度（预留头部、帮助文本和错误消息的空间）
	minHeight := 8
	maxHeight := height - 6
	if maxHeight > 24 {
		maxHeight = 24
	}
	config.listHeight = max(minHeight, maxHeight)

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
