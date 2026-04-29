package components

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/config"
	"gocoding/internal/models"
	"gocoding/internal/ui"
)

// ProviderPageState 提供商页面状态
type ProviderPageState int

const (
	ProviderStateList ProviderPageState = iota
	ProviderStateAdd
	ProviderStateEdit
	ProviderStateDeleteConfirm
)

// ProviderFocus 输入框焦点
type ProviderFocus int

const (
	ProviderFocusName ProviderFocus = iota
	ProviderFocusBaseURL
	ProviderFocusAPIKey
	ProviderFocusModel
	ProviderFocusThinkingModel
	ProviderFocusDefaultHaiku
	ProviderFocusDefaultSonnet
	ProviderFocusDefaultOpus
	ProviderFocusSubagent
	ProviderFocusNonessential
	ProviderFocusNonstreaming
	ProviderFocusEffort
	ProviderFocusCount
)

// ProviderListPage 模型配置列表页面
type ProviderListPage struct {
	store     *models.ModelProviderStore
	app       *AppModel
	list      list.Model
	width     int
	height    int
	editingID string

	// 状态机
	state ProviderPageState

	// 输入框
	nameInput              textinput.Model
	baseURLInput           textinput.Model
	apiKeyInput            textinput.Model
	modelInput             textinput.Model
	thinkingModelInput     textinput.Model
	defaultHaikuInput      textinput.Model
	defaultSonnetInput     textinput.Model
	defaultOpusInput       textinput.Model
	subagentInput          textinput.Model
	nonessentialInput      textinput.Model
	nonstreamingInput      textinput.Model
	effortInput            textinput.Model
	inputFocus             ProviderFocus

	// 删除确认
	hoverButton int

	// 错误/提示消息
	errMsg string
	tipMsg string
}

// NewProviderListPage 创建配置列表页面
func NewProviderListPage(store *models.ModelProviderStore) *ProviderListPage {
	p := &ProviderListPage{
		store: store,
		state: ProviderStateList,
	}

	delegate := providerListDelegate{}
	p.list = list.New(nil, delegate, 60, 14)
	p.list.SetShowTitle(false)
	p.list.SetShowStatusBar(false)
	p.list.SetShowHelp(false)
	p.list.SetFilteringEnabled(true)

	p.initInputs()
	p.updateListItems()

	return p
}

func (p *ProviderListPage) initInputs() {
	p.nameInput = textinput.New()
	p.nameInput.Placeholder = "Name (e.g. MiniMax)"
	p.nameInput.Focus()

	p.baseURLInput = textinput.New()
	p.baseURLInput.Placeholder = "Base URL (e.g. https://api.minimax.chat)"

	p.apiKeyInput = textinput.New()
	p.apiKeyInput.Placeholder = "API Key"

	p.modelInput = textinput.New()
	p.modelInput.Placeholder = "Model (e.g. MiniMax-M2.7-highspeed)"

	p.thinkingModelInput = textinput.New()
	p.thinkingModelInput.Placeholder = "Reasoning Model"

	p.defaultHaikuInput = textinput.New()
	p.defaultHaikuInput.Placeholder = "Default Haiku Model"

	p.defaultSonnetInput = textinput.New()
	p.defaultSonnetInput.Placeholder = "Default Sonnet Model"

	p.defaultOpusInput = textinput.New()
	p.defaultOpusInput.Placeholder = "Default Opus Model"

	p.subagentInput = textinput.New()
	p.subagentInput.Placeholder = "SubAgent Model"

	p.nonessentialInput = textinput.New()
	p.nonessentialInput.Placeholder = "1=禁用非必要流量"

	p.nonstreamingInput = textinput.New()
	p.nonstreamingInput.Placeholder = "1=禁用非流式回退"

	p.effortInput = textinput.New()
	p.effortInput.Placeholder = "max/high/medium/low"
}

// SetApp 设置 AppModel 引用
func (p *ProviderListPage) SetApp(app *AppModel) {
	p.app = app
}

// PageType 返回页面类型
func (p *ProviderListPage) PageType() PageType {
	return PageProvider
}

// OnActivate 页面激活时调用
func (p *ProviderListPage) OnActivate() {
	p.state = ProviderStateList
	p.syncListItems()
}

// OnDeactivate 页面停用时调用
func (p *ProviderListPage) OnDeactivate() {
}

// SetSize 设置尺寸
func (p *ProviderListPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// Update 处理消息
func (p *ProviderListPage) Update(msg tea.Msg) (tea.Cmd, bool) {
	// 先更新组件（包括 textinput）
	switch p.state {
	case ProviderStateAdd, ProviderStateEdit:
		var cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8, cmd9, cmd10, cmd11, cmd12 tea.Cmd
		p.nameInput, cmd1 = p.nameInput.Update(msg)
		p.baseURLInput, cmd2 = p.baseURLInput.Update(msg)
		p.apiKeyInput, cmd3 = p.apiKeyInput.Update(msg)
		p.modelInput, cmd4 = p.modelInput.Update(msg)
		p.thinkingModelInput, cmd5 = p.thinkingModelInput.Update(msg)
		p.defaultHaikuInput, cmd6 = p.defaultHaikuInput.Update(msg)
		p.defaultSonnetInput, cmd7 = p.defaultSonnetInput.Update(msg)
		p.defaultOpusInput, cmd8 = p.defaultOpusInput.Update(msg)
		p.subagentInput, cmd9 = p.subagentInput.Update(msg)
		p.nonessentialInput, cmd10 = p.nonessentialInput.Update(msg)
		p.nonstreamingInput, cmd11 = p.nonstreamingInput.Update(msg)
		p.effortInput, cmd12 = p.effortInput.Update(msg)

		// 处理按键
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch p.state {
			case ProviderStateAdd:
				if cmd := p.handleAddKeyMsg(msg); cmd != nil {
					return tea.Batch(append([]tea.Cmd{cmd}, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8, cmd9, cmd10, cmd11, cmd12)...), true
				}
			case ProviderStateEdit:
				if cmd := p.handleEditKeyMsg(msg); cmd != nil {
					return tea.Batch(append([]tea.Cmd{cmd}, cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8, cmd9, cmd10, cmd11, cmd12)...), true
				}
			}
		}
		return tea.Batch(cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8, cmd9, cmd10, cmd11, cmd12), true
	}

	// 处理按键
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch p.state {
		case ProviderStateList:
			return p.handleListKeyMsg(msg), true
		case ProviderStateDeleteConfirm:
			return p.handleDeleteConfirmKeyMsg(msg), true
		}
	}

	// 更新列表
	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return cmd, false
}

// View 渲染页面
func (p *ProviderListPage) View(width, height int) string {
	p.width = width
	p.height = height

	switch p.state {
	case ProviderStateAdd:
		return p.viewAdd()
	case ProviderStateEdit:
		return p.viewEdit()
	case ProviderStateDeleteConfirm:
		return p.viewDeleteConfirm()
	default:
		return p.viewList()
	}
}

// HandleMouse 处理鼠标消息
func (p *ProviderListPage) HandleMouse(msg tea.MouseMsg) {
	switch p.state {
	case ProviderStateList:
		p.handleListMouse(msg)
	case ProviderStateDeleteConfirm:
		p.handleDeleteConfirmMouse(msg)
	}
}

func (p *ProviderListPage) handleListMouse(msg tea.MouseMsg) {
	listStartY := 4
	listEndY := p.height - 8
	itemHeight := 3

	switch {
	case msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress:
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := (msg.Y - listStartY) / itemHeight
			items := p.list.Items()
			if clickedIndex >= 0 && clickedIndex < len(items) {
				p.list.Select(clickedIndex)
			}
		}
	case msg.Button == tea.MouseButtonWheelUp:
		p.list.CursorUp()
	case msg.Button == tea.MouseButtonWheelDown:
		p.list.CursorDown()
	}
}

func (p *ProviderListPage) handleDeleteConfirmMouse(msg tea.MouseMsg) {
	if msg.Button != tea.MouseButtonLeft || msg.Action != tea.MouseActionPress {
		return
	}

	dialogWidth := min(50, max(35, int(float64(p.width)*0.6)))
	dialogHeight := 10
	dialogX := (p.width - dialogWidth) / 2
	dialogY := (p.height - dialogHeight) / 2

	buttonY := dialogY + 6
	buttonWidth := 10

	if msg.Y == buttonY && msg.X >= dialogX+10 && msg.X < dialogX+10+buttonWidth {
		p.hoverButton = 1
	} else if msg.Y == buttonY && msg.X >= dialogX+25 && msg.X < dialogX+25+buttonWidth {
		p.hoverButton = 0
	}
}

// ============== 视图方法 ==============

func (p *ProviderListPage) viewList() string {
	dialogWidth := p.providerDialogWidth()
	contentWidth := p.providerListWidth()

	listView := lipgloss.NewStyle().
		Width(contentWidth).
		Align(lipgloss.Center).
		Render(p.list.View())

	helpText := p.renderHelpText()

	var errDisplay string
	if p.errMsg != "" {
		errDisplay = ui.ErrorBoxStyle.Render("✗ " + p.errMsg)
		p.errMsg = ""
	}

	var tipDisplay string
	if p.tipMsg != "" {
		tipDisplay = ui.SuccessBoxStyle.Render("✓ " + p.tipMsg)
		p.tipMsg = ""
	}

	dialog := lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.AccentMagenta).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.AccentMagenta).Bold(true).Render("⚙ 模型配置"),
				"",
				listView,
				errDisplay,
				tipDisplay,
				helpText,
			),
		)

	return dialog
}

func (p *ProviderListPage) viewAdd() string {
	dialogWidth := min(p.width-4, max(50, int(float64(p.width)*0.8)))

	inactiveInput := ui.InputBorder.Width(dialogWidth - 6).Padding(0, 1)
	focusedInput := ui.FocusedInputBorder.Width(dialogWidth - 6).Padding(0, 1)

	inputs := []struct {
		focus    ProviderFocus
		label    string
		input    textinput.Model
		required bool
	}{
		{ProviderFocusName, "配置名称", p.nameInput, true},
		{ProviderFocusBaseURL, "API Base URL", p.baseURLInput, true},
		{ProviderFocusAPIKey, "API Key", p.apiKeyInput, true},
		{ProviderFocusModel, "模型名称", p.modelInput, false},
		{ProviderFocusThinkingModel, "推理模型", p.thinkingModelInput, false},
		{ProviderFocusDefaultHaiku, "Haiku 默认模型", p.defaultHaikuInput, false},
		{ProviderFocusDefaultSonnet, "Sonnet 默认模型", p.defaultSonnetInput, false},
		{ProviderFocusDefaultOpus, "Opus 默认模型", p.defaultOpusInput, false},
		{ProviderFocusSubagent, "SubAgent 模型", p.subagentInput, false},
		{ProviderFocusNonessential, "禁用非必要流量 (1/空)", p.nonessentialInput, false},
		{ProviderFocusNonstreaming, "禁用非流式回退 (1/空)", p.nonstreamingInput, false},
		{ProviderFocusEffort, "推理力度 (max/high/medium/low)", p.effortInput, false},
	}

	var items []string
	items = append(items, lipgloss.NewStyle().Foreground(ui.AccentMagenta).Bold(true).Render("＋ 新增配置"))

	for _, inp := range inputs {
		style := inactiveInput
		if p.inputFocus == inp.focus {
			style = focusedInput
		}
		requiredMark := ""
		if inp.required {
			requiredMark = " *"
		}
		items = append(items, "")
		items = append(items, lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render(inp.label+requiredMark))
		items = append(items, style.Render(inp.input.View()))
	}

	items = append(items, "")
	items = append(items, lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Enter] 确认  ·  [Tab/Shift+Tab] 切换  ·  [Esc] 取消"))

	if p.errMsg != "" {
		items = append(items, "")
		items = append(items, ui.ErrorBoxStyle.Render("✗ "+p.errMsg))
	}

	return lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.AccentMagenta).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Center, items...))
}

func (p *ProviderListPage) viewEdit() string {
	dialogWidth := min(p.width-4, max(50, int(float64(p.width)*0.8)))

	inactiveInput := ui.InputBorder.Width(dialogWidth - 6).Padding(0, 1)
	focusedInput := ui.FocusedInputBorder.Width(dialogWidth - 6).Padding(0, 1)

	inputs := []struct {
		focus    ProviderFocus
		label    string
		input    textinput.Model
		required bool
	}{
		{ProviderFocusName, "配置名称", p.nameInput, true},
		{ProviderFocusBaseURL, "API Base URL", p.baseURLInput, true},
		{ProviderFocusAPIKey, "API Key（留空则不修改）", p.apiKeyInput, false},
		{ProviderFocusModel, "模型名称", p.modelInput, false},
		{ProviderFocusThinkingModel, "推理模型", p.thinkingModelInput, false},
		{ProviderFocusDefaultHaiku, "Haiku 默认模型", p.defaultHaikuInput, false},
		{ProviderFocusDefaultSonnet, "Sonnet 默认模型", p.defaultSonnetInput, false},
		{ProviderFocusDefaultOpus, "Opus 默认模型", p.defaultOpusInput, false},
		{ProviderFocusSubagent, "SubAgent 模型", p.subagentInput, false},
		{ProviderFocusNonessential, "禁用非必要流量 (1/空)", p.nonessentialInput, false},
		{ProviderFocusNonstreaming, "禁用非流式回退 (1/空)", p.nonstreamingInput, false},
		{ProviderFocusEffort, "推理力度 (max/high/medium/low)", p.effortInput, false},
	}

	var items []string
	items = append(items, lipgloss.NewStyle().Foreground(ui.AccentMagenta).Bold(true).Render("✎ 编辑配置"))

	for _, inp := range inputs {
		style := inactiveInput
		if p.inputFocus == inp.focus {
			style = focusedInput
		}
		requiredMark := ""
		if inp.required {
			requiredMark = " *"
		}
		items = append(items, "")
		items = append(items, lipgloss.NewStyle().Foreground(ui.ForegroundDim).Render(inp.label+requiredMark))
		items = append(items, style.Render(inp.input.View()))
	}

	items = append(items, "")
	items = append(items, lipgloss.NewStyle().Foreground(ui.SecondaryText).Render("[Enter] 确认  ·  [Tab/Shift+Tab] 切换  ·  [Esc] 取消"))

	if p.errMsg != "" {
		items = append(items, "")
		items = append(items, ui.ErrorBoxStyle.Render("✗ "+p.errMsg))
	}

	return lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.AccentMagenta).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Center, items...))
}

func (p *ProviderListPage) viewDeleteConfirm() string {
	items := p.list.Items()
	if len(items) == 0 || p.list.Index() >= len(items) {
		p.state = ProviderStateList
		return p.viewList()
	}

	selectedItem := items[p.list.Index()]
	item, ok := selectedItem.(providerListItem)
	if !ok {
		p.state = ProviderStateList
		return p.viewList()
	}

	dialogWidth := min(50, max(35, int(float64(p.width)*0.6)))
	buttonWidth := 10

	confirmStyle := lipgloss.NewStyle().Width(buttonWidth).Foreground(ui.ErrorColor).Background(lipgloss.Color("#1A0D10")).Padding(0, 2)
	if p.hoverButton == 1 {
		confirmStyle = confirmStyle.Background(ui.ErrorColor).Foreground(ui.Background).Bold(true)
	}

	cancelStyle := lipgloss.NewStyle().Width(buttonWidth).Foreground(ui.SecondaryText).Background(ui.BackgroundLight).Padding(0, 2)
	if p.hoverButton == 0 {
		cancelStyle = cancelStyle.Background(ui.BackgroundHover).Foreground(ui.Foreground)
	}

	return lipgloss.NewStyle().
		Width(dialogWidth).
		Border(ui.NeonBorder).
		BorderForeground(ui.ErrorColor).
		Background(ui.BackgroundSurface).
		Foreground(ui.Foreground).
		Padding(1, 2).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.NewStyle().Foreground(ui.ErrorColor).Bold(true).Render("✗ 删除配置"),
				"",
				lipgloss.NewStyle().Foreground(ui.Foreground).Render(fmt.Sprintf("确认删除配置 '%s'？", item.provider.Name)),
				"",
				lipgloss.JoinHorizontal(lipgloss.Center, cancelStyle.Render("[ N ] 否"), "  ", confirmStyle.Render("[ Y ] 是")),
				"",
				lipgloss.NewStyle().Foreground(ui.MutedText).Render("此操作不可恢复"),
			),
		)
}

// ============== 处理器 ==============

func (p *ProviderListPage) handleListKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down", "ctrl+n":
		p.list.CursorDown()
	case "k", "up", "ctrl+p":
		p.list.CursorUp()
	case "n":
		p.state = ProviderStateAdd
		p.resetInputs()
		p.inputFocus = ProviderFocusName
		p.nameInput.Focus()
		p.errMsg = ""
		return textinput.Blink
	case "e":
		if p.getSelectedProvider() != nil {
			p.loadProviderToInputs(p.getSelectedProvider())
			p.state = ProviderStateEdit
			p.editingID = p.getSelectedProvider().ID
			return textinput.Blink
		}
	case "d":
		if p.getSelectedProvider() != nil {
			p.state = ProviderStateDeleteConfirm
			p.hoverButton = 0
		}
	case "a":
		if p.getSelectedProvider() != nil {
			provider := p.getSelectedProvider()
			p.store.SetActive(provider.ID)
			// 写入 Claude settings.json
			activeProvider := p.store.Get(provider.ID)
			alreadySet, err := config.WriteToClaudeSettings(activeProvider)
			p.updateListItems()
			if err != nil {
				p.errMsg = "激活失败: " + err.Error()
			} else if alreadySet {
				p.tipMsg = provider.Name + " 配置已生效，无需更新"
			} else {
				p.tipMsg = "已激活 " + provider.Name
			}
		}
	case "enter":
		if p.getSelectedProvider() != nil {
			p.loadProviderToInputs(p.getSelectedProvider())
			p.state = ProviderStateEdit
			p.editingID = p.getSelectedProvider().ID
			return textinput.Blink
		}
	case "esc":
		if p.app != nil {
			p.app.SwitchPage(p.app.projectPage)
			p.app.projectPage.OnActivate()
		}
	case "ctrl+c", "ctrl+q", "q":
		return tea.Quit
	}
	return nil
}

func (p *ProviderListPage) handleAddKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		name := p.nameInput.Value()
		baseURL := p.baseURLInput.Value()
		apiKey := p.apiKeyInput.Value()
		model := p.modelInput.Value()
		thinkingModel := p.thinkingModelInput.Value()
		defaultHaiku := p.defaultHaikuInput.Value()
		defaultSonnet := p.defaultSonnetInput.Value()
		defaultOpus := p.defaultOpusInput.Value()
		subagent := p.subagentInput.Value()
		nonessential := p.nonessentialInput.Value()
		nonstreaming := p.nonstreamingInput.Value()
		effort := p.effortInput.Value()

		if name == "" {
			p.errMsg = "配置名称不能为空"
			return nil
		}
		if baseURL == "" {
			p.errMsg = "API Base URL 不能为空"
			return nil
		}
		if apiKey == "" {
			p.errMsg = "API Key 不能为空"
			return nil
		}

		provider := models.ModelProvider{
			ID:                   models.GenerateProviderID(),
			Name:                 name,
			BaseURL:              baseURL,
			APIKey:               apiKey,
			Model:                model,
			ThinkingModel:        thinkingModel,
			DefaultHaikuModel:    defaultHaiku,
			DefaultSonnetModel:   defaultSonnet,
			DefaultOpusModel:     defaultOpus,
			SubagentModel:        subagent,
			DisableNonessential:  nonessential,
			DisableNonstreaming:  nonstreaming,
			EffortLevel:          effort,
			CreatedAt:            time.Now(),
		}

		p.store.Add(provider)
		p.updateListItems()
		p.state = ProviderStateList

		if p.app != nil {
			p.app.ShowToast("已添加: "+name, "success")
		}

	case "tab":
		p.inputFocus = (p.inputFocus + 1) % ProviderFocusCount
		p.updateFocus()
	case "shift+tab":
		if p.inputFocus == 0 {
			p.inputFocus = ProviderFocusCount - 1
		} else {
			p.inputFocus--
		}
		p.updateFocus()
	case "esc":
		p.state = ProviderStateList
	}
	return nil
}

func (p *ProviderListPage) handleEditKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		name := p.nameInput.Value()
		baseURL := p.baseURLInput.Value()
		apiKey := p.apiKeyInput.Value()
		model := p.modelInput.Value()
		thinkingModel := p.thinkingModelInput.Value()
		defaultHaiku := p.defaultHaikuInput.Value()
		defaultSonnet := p.defaultSonnetInput.Value()
		defaultOpus := p.defaultOpusInput.Value()
		subagent := p.subagentInput.Value()
		nonessential := p.nonessentialInput.Value()
		nonstreaming := p.nonstreamingInput.Value()
		effort := p.effortInput.Value()

		if name == "" {
			p.errMsg = "配置名称不能为空"
			return nil
		}
		if baseURL == "" {
			p.errMsg = "API Base URL 不能为空"
			return nil
		}

		// 如果 API Key 为空，则不更新
		if apiKey == "" {
			provider := p.store.Get(p.editingID)
			if provider != nil {
				apiKey = provider.APIKey
			}
		}

		p.store.Update(p.editingID, name, baseURL, apiKey, model, thinkingModel, defaultHaiku, defaultSonnet, defaultOpus, subagent, nonessential, nonstreaming, effort)
		p.updateListItems()
		p.state = ProviderStateList

		if p.app != nil {
			p.app.ShowToast("已更新: "+name, "success")
		}

	case "tab":
		p.inputFocus = (p.inputFocus + 1) % ProviderFocusCount
		p.updateFocus()
	case "shift+tab":
		if p.inputFocus == 0 {
			p.inputFocus = ProviderFocusCount - 1
		} else {
			p.inputFocus--
		}
		p.updateFocus()
	case "esc":
		p.state = ProviderStateList
	}
	return nil
}

func (p *ProviderListPage) handleDeleteConfirmKeyMsg(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y", "enter":
		if p.hoverButton == 1 {
			provider := p.getSelectedProvider()
			if provider != nil {
				name := provider.Name
				p.store.Remove(provider.ID)
				p.updateListItems()
				if p.app != nil {
					p.app.ShowToast("已删除: "+name, "success")
				}
			}
		}
		p.state = ProviderStateList
	case "n", "esc":
		p.state = ProviderStateList
	}
	return nil
}

// ============== 辅助方法 ==============

func (p *ProviderListPage) updateFocus() {
	p.nameInput.Blur()
	p.baseURLInput.Blur()
	p.apiKeyInput.Blur()
	p.modelInput.Blur()
	p.thinkingModelInput.Blur()
	p.defaultHaikuInput.Blur()
	p.defaultSonnetInput.Blur()
	p.defaultOpusInput.Blur()
	p.subagentInput.Blur()
	p.nonessentialInput.Blur()
	p.nonstreamingInput.Blur()
	p.effortInput.Blur()

	switch p.inputFocus {
	case ProviderFocusName:
		p.nameInput.Focus()
	case ProviderFocusBaseURL:
		p.baseURLInput.Focus()
	case ProviderFocusAPIKey:
		p.apiKeyInput.Focus()
	case ProviderFocusModel:
		p.modelInput.Focus()
	case ProviderFocusThinkingModel:
		p.thinkingModelInput.Focus()
	case ProviderFocusDefaultHaiku:
		p.defaultHaikuInput.Focus()
	case ProviderFocusDefaultSonnet:
		p.defaultSonnetInput.Focus()
	case ProviderFocusDefaultOpus:
		p.defaultOpusInput.Focus()
	case ProviderFocusSubagent:
		p.subagentInput.Focus()
	case ProviderFocusNonessential:
		p.nonessentialInput.Focus()
	case ProviderFocusNonstreaming:
		p.nonstreamingInput.Focus()
	case ProviderFocusEffort:
		p.effortInput.Focus()
	}
}

func (p *ProviderListPage) resetInputs() {
	p.nameInput.Reset()
	p.nameInput.Placeholder = "Name (e.g. MiniMax)"
	p.baseURLInput.Reset()
	p.baseURLInput.Placeholder = "Base URL (e.g. https://api.minimax.chat)"
	p.apiKeyInput.Reset()
	p.apiKeyInput.Placeholder = "API Key"
	p.modelInput.Reset()
	p.modelInput.Placeholder = "Model (e.g. MiniMax-M2.7-highspeed)"
	p.thinkingModelInput.Reset()
	p.thinkingModelInput.Placeholder = "Reasoning Model"
	p.defaultHaikuInput.Reset()
	p.defaultHaikuInput.Placeholder = "Default Haiku Model"
	p.defaultSonnetInput.Reset()
	p.defaultSonnetInput.Placeholder = "Default Sonnet Model"
	p.defaultOpusInput.Reset()
	p.defaultOpusInput.Placeholder = "Default Opus Model"
	p.subagentInput.Reset()
	p.subagentInput.Placeholder = "SubAgent Model"
	p.nonessentialInput.Reset()
	p.nonessentialInput.Placeholder = "1=禁用非必要流量"
	p.nonstreamingInput.Reset()
	p.nonstreamingInput.Placeholder = "1=禁用非流式回退"
	p.effortInput.Reset()
	p.effortInput.Placeholder = "max/high/medium/low"
}

func (p *ProviderListPage) loadProviderToInputs(provider *models.ModelProvider) {
	p.nameInput.SetValue(provider.Name)
	p.nameInput.Focus()
	p.baseURLInput.SetValue(provider.BaseURL)
	p.apiKeyInput.SetValue("") // 不显示 API Key
	p.apiKeyInput.Placeholder = "API Key（留空则不修改）"
	p.modelInput.SetValue(provider.Model)
	p.thinkingModelInput.SetValue(provider.ThinkingModel)
	p.defaultHaikuInput.SetValue(provider.DefaultHaikuModel)
	p.defaultSonnetInput.SetValue(provider.DefaultSonnetModel)
	p.defaultOpusInput.SetValue(provider.DefaultOpusModel)
	p.subagentInput.SetValue(provider.SubagentModel)
	p.nonessentialInput.SetValue(provider.DisableNonessential)
	p.nonstreamingInput.SetValue(provider.DisableNonstreaming)
	p.effortInput.SetValue(provider.EffortLevel)
	p.inputFocus = ProviderFocusName
	p.errMsg = ""
}

func (p *ProviderListPage) getSelectedProvider() *models.ModelProvider {
	items := p.list.Items()
	if len(items) == 0 || p.list.Index() >= len(items) {
		return nil
	}

	item, ok := items[p.list.Index()].(providerListItem)
	if !ok {
		return nil
	}
	return &item.provider
}

func (p *ProviderListPage) syncListItems() {
	p.updateListItems()
}

// renderHelpText 渲染帮助文本
func (p *ProviderListPage) renderHelpText() string {
	return lipgloss.NewStyle().
		Foreground(ui.SecondaryText).
		Width(p.providerListWidth()).
		Render(
			lipgloss.JoinVertical(
				lipgloss.Center,
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					p.renderHelpItem("[k↑/j↓]", "选择", ui.HelpKeyNavStyle),
					"  ",
					p.renderHelpItem("[N]", "新增", ui.HelpKeyActionStyle),
					"  ",
					p.renderHelpItem("[E]", "编辑", ui.HelpKeyActionStyle),
					"  ",
					p.renderHelpItem("[D]", "删除", ui.HelpKeyDangerStyle),
					"  ",
					p.renderHelpItem("[A]", "激活", ui.HelpKeyActionStyle),
				),
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					p.renderHelpItem("[Esc]", "退出", ui.HelpKeyQuitStyle),
				),
			),
		)
}

func (p *ProviderListPage) renderHelpItem(key, label string, keyStyle lipgloss.Style) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		keyStyle.Render(key),
		" ",
		lipgloss.NewStyle().Foreground(ui.SecondaryText).Render(label),
	)
}

// providerDialogWidth 计算对话框宽度
func (p *ProviderListPage) providerDialogWidth() int {
	return min(p.width-4, max(50, int(float64(p.width)*0.8)))
}

// providerListWidth 计算列表宽度
func (p *ProviderListPage) providerListWidth() int {
	return p.providerDialogWidth() - 4
}

// updateListItems 更新列表项
func (p *ProviderListPage) updateListItems() {
	if p.store == nil {
		return
	}
	p.list.SetItems(newProviderListItems(p.store.Providers))
}

// newProviderListItems 创建配置列表项
func newProviderListItems(providers []models.ModelProvider) []list.Item {
	items := make([]list.Item, len(providers))
	for i, prov := range providers {
		items[i] = providerListItem{provider: prov}
	}
	return items
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

	accentClr := providerAccentColor(provider.BaseURL)

	activeTag := ""
	if provider.Active {
		activeTag = ui.ActiveBadgeStyle.Render("● 激活")
	}

	selector := "  "
	if isSelected {
		selector = lipgloss.NewStyle().Foreground(accentClr).Render("▸ ")
	}

	selectorWidth := lipgloss.Width(selector)
	activeWidth := lipgloss.Width(activeTag)
	nameWidth := max(20, m.Width()) - selectorWidth
	if activeWidth > 0 {
		nameWidth -= activeWidth + 1
	}
	nameWidth = max(8, nameWidth)

	name := provider.Name

	var firstLine string
	nameStyle := ui.ProviderItemStyle
	if isSelected && provider.Active {
		nameStyle = ui.ProviderActiveItemStyle
	} else if isSelected {
		nameStyle = ui.ProviderSelectedItemStyle
	} else if provider.Active {
		nameStyle = ui.ProviderActiveItemStyle
	}

	if provider.Active {
		firstLine = lipgloss.JoinHorizontal(
			lipgloss.Left,
			selector,
			nameStyle.Bold(true).Width(nameWidth).Render(name),
			ui.ActiveBadgeStyle.Render("● 激活"),
		)
	} else if isSelected {
		firstLine = lipgloss.JoinHorizontal(
			lipgloss.Left,
			selector,
			nameStyle.Bold(true).Width(nameWidth).Render(name),
		)
	} else {
		firstLine = lipgloss.JoinHorizontal(
			lipgloss.Left,
			selector,
			nameStyle.Width(nameWidth).Render(name),
		)
	}

	infoText := provider.BaseURL
	if provider.Model != "" {
		infoText += " • " + provider.Model
	}
	infoColor := ui.SecondaryText
	if !provider.Active {
		infoColor = ui.ForegroundDim
	}
	infoPrefix := lipgloss.NewStyle().Foreground(ui.MutedText).Width(selectorWidth).Render("")
	secondLine := infoPrefix + lipgloss.NewStyle().Foreground(infoColor).Width(max(20, m.Width())-selectorWidth).Render(infoText)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		firstLine,
		secondLine,
	)

	fmt.Fprintf(w, "%s", content)
}

// providerAccentColor 根据 URL 返回 provider 类型色彩提示
func providerAccentColor(url string) lipgloss.Color {
	lower := strings.ToLower(url)
	switch {
	case strings.Contains(lower, "anthropic"):
		return ui.AccentGold
	case strings.Contains(lower, "minimax"):
		return ui.AccentMagenta
	case strings.Contains(lower, "openai"):
		return ui.SuccessColor
	case strings.Contains(lower, "deepseek"):
		return ui.AccentCyan
	default:
		return ui.PrimaryColor
	}
}
