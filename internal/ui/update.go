package ui

import (
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

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
