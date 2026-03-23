package ui

import (
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"gocoding/internal/config"
	"gocoding/internal/models"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch m.state {
		case StateList:
			return m.handleListKeyMsg(msg)
		case StateAddProject:
			return m.handleAddProjectKeyMsg(msg)
		case StateRenameProject:
			return m.handleRenameKeyMsg(msg)
		case StateDeleteConfirm:
			return m.handleDeleteConfirmKeyMsg(msg)
		case StateIDEMenu:
			return m.handleIDEMenuKeyMsg(msg)
		case StateViewDetail:
			return m.handleViewDetailKeyMsg(msg)
		case StateEditDescription:
			return m.handleEditDescriptionKeyMsg(msg)
		case StateSearch:
			return m.handleSearchKeyMsg(msg)
		case StateProviderList:
			return m.handleProviderListKeyMsg(msg)
		case StateProviderAdd:
			return m.handleProviderAddKeyMsg(msg)
		case StateProviderEdit:
			return m.handleProviderEditKeyMsg(msg)
		case StateProviderDelete:
			return m.handleProviderDeleteKeyMsg(msg)
		}
	}

	// 通用组件更新
	switch m.state {
	case StateEditDescription:
		m.ta, cmd = m.ta.Update(msg)
		cmds = append(cmds, cmd)
	case StateViewDetail:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// 更新详情视图
	if m.showDetails && m.state == StateList {
		m.updateViewport()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.list.CursorDown()
	case "k", "up":
		m.list.CursorUp()
	case "/":
		// 进入搜索状态
		m.state = StateSearch
		m.searchQuery = ""
		return m, nil
	case "n":
		m.state = StateAddProject
		m.input.Reset()
		m.input.Placeholder = "输入项目路径"
		m.input.SetValue("")
		m.input.Focus()
		m.inputFocus = FocusPath
		m.secondaryInput.Reset()
		m.secondaryInput.Placeholder = "输入项目名称"
		m.secondaryInput.SetValue("")
		m.secondaryInput.Blur()
		m.tempPath = ""
		return m, textinput.Blink
	case "e":
		// 编辑选中项目的描述
		if len(m.list.Items()) > 0 {
			m.state = StateEditDescription
			current := m.list.SelectedItem().(listItem)
			m.ta.SetValue(current.project.Description)
		}
		return m, nil
	case "r":
		if len(m.list.Items()) > 0 {
			m.state = StateRenameProject
			current := m.list.SelectedItem().(listItem)
			m.input.SetValue(current.project.Alias)
			m.input.Focus()
		}
		return m, textinput.Blink
	case "d":
		if len(m.list.Items()) > 0 {
			m.state = StateDeleteConfirm
		}
	case "v":
		// 切换详情视图
		if len(m.list.Items()) > 0 {
			m.showDetails = !m.showDetails
			m.updateViewport()
		}
	case "enter":
		if len(m.list.Items()) > 0 {
			m.state = StateIDEMenu
			for _, opt := range m.ideMenu.options {
				m.ideMenu.available[opt.Type] = m.ideExec.IsIDEAvailable(opt.Type)
			}
		}
	case "1":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEOpenCode)
		}
	case "p":
		// 进入模型配置列表
		m.state = StateProviderList
		return m, nil
	case "ctrl+c", "ctrl+q", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleAddProjectKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// 如果路径和名称都有效，添加项目
		path := m.input.Value()
		name := m.secondaryInput.Value()

		// 验证路径
		if err := m.store.ValidatePath(path); err != nil {
			m.errMsg = err.Error()
			return m, nil
		}

		if path != "" && name != "" {
			project := models.Project{
				ID:        generateID(),
				Path:      path,
				Alias:     name,
				CreatedAt: time.Now(),
			}
			m.store.Add(project)
			m.state = StateList
			m.tempPath = ""
			m.tempName = ""
			m.input.Reset()
			m.secondaryInput.Reset()
			items := make([]list.Item, len(m.store.Projects))
			for i, p := range m.store.Projects {
				items[i] = listItem{project: p}
			}
			m.list.SetItems(items)
		}
	case "tab":
		// 切换焦点
		if m.inputFocus == FocusPath {
			m.inputFocus = FocusName
			m.input.Blur()
			m.secondaryInput.Focus()
		} else {
			m.inputFocus = FocusPath
			m.secondaryInput.Blur()
			m.input.Focus()
		}
	case "esc":
		m.state = StateList
		m.input.Reset()
		m.secondaryInput.Reset()
		m.tempPath = ""
		m.tempName = ""
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	var cmd tea.Cmd
	// 只更新当前焦点所在输入框
	if m.inputFocus == FocusPath {
		m.input, cmd = m.input.Update(msg)
		// 更新后再检查路径变化
		path := m.input.Value()
		if path != m.tempPath {
			m.tempPath = path
			defaultAlias := filepath.Base(path)
			if defaultAlias == "" || defaultAlias == "/" || defaultAlias == "\\" {
				defaultAlias = ""
			}
			m.secondaryInput.SetValue(defaultAlias)
		}
	} else {
		m.secondaryInput, cmd = m.secondaryInput.Update(msg)
	}
	return m, cmd
}

func (m *Model) handleRenameKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		alias := m.input.Value()
		if alias != "" {
			current := m.list.SelectedItem().(listItem)
			current.project.Alias = alias
			m.store.Update(alias, current.project.ID)
			items := make([]list.Item, len(m.store.Projects))
			for i, p := range m.store.Projects {
				items[i] = listItem{project: p}
			}
			m.list.SetItems(items)
		}
		m.state = StateList
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) handleDeleteConfirmKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		current := m.list.SelectedItem().(listItem)
		m.store.Remove(current.project.ID)
		items := make([]list.Item, len(m.store.Projects))
		for i, p := range m.store.Projects {
			items[i] = listItem{project: p}
		}
		m.list.SetItems(items)
		m.state = StateList
	case "n", "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleIDEMenuKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.ideMenu.selected = min(m.ideMenu.selected+1, len(m.ideMenu.options)-1)
	case "k", "up":
		m.ideMenu.selected = max(m.ideMenu.selected-1, 0)
	case "enter":
		selectedIDE := m.ideMenu.options[m.ideMenu.selected].Type
		return m.openWithIDE(selectedIDE)
	case "1":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEOpenCode)
		}
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) openWithIDE(ideType models.IDEType) (tea.Model, tea.Cmd) {
	current := m.list.SelectedItem().(listItem)
	current.project.UpdateLastOpened()
	m.store.Update(current.project.Alias, current.project.ID)

	if err := m.ideExec.OpenProject(&current.project, ideType); err != nil {
		m.errMsg = err.Error()
		// 3秒后清除错误消息
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			m.errMsg = ""
			return nil
		})
	}

	// 按最后打开时间排序
	m.store.SortByLastOpened()

	// 重新渲染列表
	items := make([]list.Item, len(m.store.Projects))
	for i, p := range m.store.Projects {
		items[i] = listItem{project: p}
	}
	m.list.SetItems(items)
	// 定位到最近打开的项目（第一个）
	m.list.Select(0)

	m.state = StateList
	return m, nil
}

func (m *Model) handleViewDetailKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleEditDescriptionKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "ctrl+s":
		if len(m.list.Items()) > 0 {
			current := m.list.SelectedItem().(listItem)
			current.project.Description = m.ta.Value()
			m.store.UpdateDescription(current.project.ID, m.ta.Value())
		}
		m.state = StateList
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleSearchKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// 退出搜索状态
		m.state = StateList
		m.searchQuery = ""
		m.updateListItems()
		return m, nil
	case "enter":
		// 打开选中的项目
		if len(m.list.Items()) > 0 {
			m.state = StateIDEMenu
			for _, opt := range m.ideMenu.options {
				m.ideMenu.available[opt.Type] = m.ideExec.IsIDEAvailable(opt.Type)
			}
		}
		return m, nil
	case "j", "down":
		m.list.CursorDown()
		return m, nil
	case "k", "up":
		m.list.CursorUp()
		return m, nil
	case "backspace":
		// 删除最后一个字符
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		m.updateListItems()
		return m, nil
	case "ctrl+h":
		// Ctrl+H 也支持删除
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
		m.updateListItems()
		return m, nil
	case "1":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDEOpenCode)
		}
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	default:
		// 处理搜索输入
		if len(msg.Runes) > 0 {
			m.searchQuery += string(msg.Runes[0])
		}
		// 更新列表显示
		m.updateListItems()
	}
	return m, nil
}

// updateListItems 根据搜索查询更新列表项目
func (m *Model) updateListItems() {
	results := m.store.Search(m.searchQuery)
	items := make([]list.Item, len(results))
	for i, p := range results {
		items[i] = listItem{project: p}
	}
	m.list.SetItems(items)
}

// handleProviderListKeyMsg 处理配置列表按键
func (m *Model) handleProviderListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		m.providerList.CursorDown()
	case "k", "up":
		m.providerList.CursorUp()
	case "n":
		// 新增配置
		m.state = StateProviderAdd
		m.providerInputFocus = FocusProviderName
		m.providerNameInput.Reset()
		m.providerNameInput.Placeholder = "配置名称 (如 MiniMax)"
		m.providerNameInput.SetValue("")
		m.providerNameInput.Focus()
		m.providerBaseURLInput.Reset()
		m.providerBaseURLInput.Placeholder = "Base URL (如 https://api.minimax.chat)"
		m.providerBaseURLInput.SetValue("")
		m.providerBaseURLInput.Blur()
		m.providerAPIKeyInput.Reset()
		m.providerAPIKeyInput.Placeholder = "API Key"
		m.providerAPIKeyInput.SetValue("")
		m.providerAPIKeyInput.Blur()
		m.providerModelInput.Reset()
		m.providerModelInput.Placeholder = "模型名称 (如 MiniMax-M2.7-highspeed)"
		m.providerModelInput.SetValue("")
		m.providerModelInput.Blur()
		m.editingProviderID = ""
		m.errMsg = ""
		return m, textinput.Blink
	case "e":
		// 编辑选中配置
		if len(m.providerList.Items()) > 0 {
			m.state = StateProviderEdit
			current := m.providerList.SelectedItem().(providerListItem)
			m.editingProviderID = current.provider.ID
			m.providerInputFocus = FocusProviderName
			m.providerNameInput.SetValue(current.provider.Name)
			m.providerNameInput.Focus()
			m.providerBaseURLInput.SetValue(current.provider.BaseURL)
			m.providerBaseURLInput.Blur()
			m.providerAPIKeyInput.SetValue(current.provider.APIKey)
			m.providerAPIKeyInput.Blur()
			m.providerModelInput.SetValue(current.provider.Model)
			m.providerModelInput.Blur()
			m.errMsg = ""
			return m, textinput.Blink
		}
	case "d":
		// 删除确认
		if len(m.providerList.Items()) > 0 {
			m.state = StateProviderDelete
		}
	case "a":
		// 激活选中配置
		if len(m.providerList.Items()) > 0 {
			current := m.providerList.SelectedItem().(providerListItem)
			m.providerStore.SetActive(current.provider.ID)
			// 写入 Claude settings.json 并备份
			alreadySet, err := config.WriteToClaudeSettings(&current.provider)
			if err != nil {
				m.errMsg = "激活失败: " + err.Error()
			} else if alreadySet {
				m.errMsg = "配置已生效，无需更新"
			} else {
				m.errMsg = "激活成功"
			}
			m.updateProviderListItems()
		}
	case "enter":
		// 编辑选中配置
		if len(m.providerList.Items()) > 0 {
			m.state = StateProviderEdit
			current := m.providerList.SelectedItem().(providerListItem)
			m.editingProviderID = current.provider.ID
			m.providerInputFocus = FocusProviderName
			m.providerNameInput.SetValue(current.provider.Name)
			m.providerNameInput.Focus()
			m.providerBaseURLInput.SetValue(current.provider.BaseURL)
			m.providerBaseURLInput.Blur()
			m.providerAPIKeyInput.SetValue(current.provider.APIKey)
			m.providerAPIKeyInput.Blur()
			m.providerModelInput.SetValue(current.provider.Model)
			m.providerModelInput.Blur()
			m.errMsg = ""
			return m, textinput.Blink
		}
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	// 更新列表
	m.providerList, _ = m.providerList.Update(msg)
	return m, nil
}

// handleProviderAddKeyMsg 处理新增配置按键
func (m *Model) handleProviderAddKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := m.providerNameInput.Value()
		baseURL := m.providerBaseURLInput.Value()
		apiKey := m.providerAPIKeyInput.Value()
		model := m.providerModelInput.Value()

		// 验证
		if name == "" {
			m.errMsg = "配置名称不能为空"
			return m, nil
		}
		if baseURL == "" {
			m.errMsg = "Base URL不能为空"
			return m, nil
		}
		if apiKey == "" {
			m.errMsg = "API Key不能为空"
			return m, nil
		}
		if model == "" {
			m.errMsg = "模型名称不能为空"
			return m, nil
		}

		provider := models.ModelProvider{
			ID:        models.GenerateProviderID(),
			Name:      name,
			BaseURL:   baseURL,
			APIKey:    apiKey,
			Model:     model,
			CreatedAt: time.Now(),
		}
		m.providerStore.Add(provider)
		m.updateProviderListItems()
		m.state = StateProviderList
		m.errMsg = ""
		return m, nil
	case "tab":
		// 切换焦点
		m.providerInputFocus = (m.providerInputFocus + 1) % FocusProviderCount
		m.updateProviderFocus()
	case "esc":
		m.state = StateProviderList
		m.errMsg = ""
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	// 更新当前焦点输入框
	return m, m.updateProviderInput(msg)
}

// handleProviderEditKeyMsg 处理编辑配置按键
func (m *Model) handleProviderEditKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := m.providerNameInput.Value()
		baseURL := m.providerBaseURLInput.Value()
		apiKey := m.providerAPIKeyInput.Value()
		model := m.providerModelInput.Value()

		// 验证
		if name == "" {
			m.errMsg = "配置名称不能为空"
			return m, nil
		}
		if baseURL == "" {
			m.errMsg = "Base URL不能为空"
			return m, nil
		}
		if apiKey == "" {
			m.errMsg = "API Key不能为空"
			return m, nil
		}
		if model == "" {
			m.errMsg = "模型名称不能为空"
			return m, nil
		}

		m.providerStore.Update(m.editingProviderID, name, baseURL, apiKey, model)
		m.updateProviderListItems()
		m.state = StateProviderList
		m.errMsg = ""
		return m, nil
	case "tab":
		// 切换焦点
		m.providerInputFocus = (m.providerInputFocus + 1) % FocusProviderCount
		m.updateProviderFocus()
	case "esc":
		m.state = StateProviderList
		m.errMsg = ""
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	// 更新当前焦点输入框
	return m, m.updateProviderInput(msg)
}

// handleProviderDeleteKeyMsg 处理删除确认按键
func (m *Model) handleProviderDeleteKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		current := m.providerList.SelectedItem().(providerListItem)
		m.providerStore.Remove(current.provider.ID)
		m.updateProviderListItems()
		m.state = StateProviderList
	case "n", "esc":
		m.state = StateProviderList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// updateProviderFocus 更新焦点状态
func (m *Model) updateProviderFocus() {
	switch m.providerInputFocus {
	case FocusProviderName:
		m.providerNameInput.Focus()
		m.providerBaseURLInput.Blur()
		m.providerAPIKeyInput.Blur()
		m.providerModelInput.Blur()
	case FocusProviderBaseURL:
		m.providerNameInput.Blur()
		m.providerBaseURLInput.Focus()
		m.providerAPIKeyInput.Blur()
		m.providerModelInput.Blur()
	case FocusProviderAPIKey:
		m.providerNameInput.Blur()
		m.providerBaseURLInput.Blur()
		m.providerAPIKeyInput.Focus()
		m.providerModelInput.Blur()
	case FocusProviderModel:
		m.providerNameInput.Blur()
		m.providerBaseURLInput.Blur()
		m.providerAPIKeyInput.Blur()
		m.providerModelInput.Focus()
	}
}

// updateProviderInput 更新焦点的输入框
func (m *Model) updateProviderInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.providerInputFocus {
	case FocusProviderName:
		m.providerNameInput, cmd = m.providerNameInput.Update(msg)
	case FocusProviderBaseURL:
		m.providerBaseURLInput, cmd = m.providerBaseURLInput.Update(msg)
	case FocusProviderAPIKey:
		m.providerAPIKeyInput, cmd = m.providerAPIKeyInput.Update(msg)
	case FocusProviderModel:
		m.providerModelInput, cmd = m.providerModelInput.Update(msg)
	}
	return cmd
}
