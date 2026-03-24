package ui

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
	config := m.calculateLayout(width, height)
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

func (m *Model) providerDialogWidth() int {
	return min(72, max(46, m.width-10))
}

func (m *Model) providerListWidth() int {
	return max(30, m.providerDialogWidth()-4)
}

func (m *Model) calculateLayout(width, height int) LayoutConfig {
	config := LayoutConfig{
		columnGap:   2,
		paddingX:    2,
		columnCount: 1, // 始终使用单列布局
	}

	// 根据屏幕宽度确定列表宽度和帮助文本模式
	switch {
	case width < 60:
		// 超窄屏
		config.listWidth = width - config.paddingX*2 - 2
		config.helpTextMode = HelpTextCompact
	case width < 100:
		// 窄屏到中屏
		config.listWidth = width - 20
		config.helpTextMode = HelpTextNormal
	default:
		// 宽屏
		config.listWidth = min(60, width-20)
		config.helpTextMode = HelpTextNormal
	}

	// 确保最小宽度
	if config.listWidth < 30 {
		config.listWidth = 30
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
