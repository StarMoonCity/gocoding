package ui

import (
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbletea"

	"gocoding/internal/models"
)

// handleAddProjectKeyMsg 处理添加项目状态按键
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
			m.list.SetItems(newListItems(m.store.Projects))
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

// handleRenameKeyMsg 处理重命名状态按键
func (m *Model) handleRenameKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		path := m.input.Value()
		name := m.secondaryInput.Value()
		if m.editingProjectID != "" {
			// 验证路径（如果改变了）
			if path != "" {
				current := m.store.Get(m.editingProjectID)
				if current != nil && current.Path != path {
					if err := m.store.ValidatePath(path); err != nil {
						m.errMsg = err.Error()
						return m, nil
					}
				}
			}
			if name != "" || path != "" {
				m.store.Update(m.editingProjectID, name, path)
				m.syncListItems()
			}
		}
		m.state = StateList
		m.editingProjectID = ""
		m.input.Reset()
		m.secondaryInput.Reset()
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
		m.editingProjectID = ""
		m.input.Reset()
		m.secondaryInput.Reset()
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	var cmd tea.Cmd
	// 只更新当前焦点所在输入框
	if m.inputFocus == FocusPath {
		m.input, cmd = m.input.Update(msg)
	} else {
		m.secondaryInput, cmd = m.secondaryInput.Update(msg)
	}
	return m, cmd
}

// handleDeleteConfirmKeyMsg 处理删除确认状态按键
func (m *Model) handleDeleteConfirmKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		current := m.list.SelectedItem().(listItem)
		m.store.Remove(current.project.ID)
		m.syncListItems()
		m.state = StateList
	case "n", "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// updateListItems 根据搜索查询更新列表项目
func (m *Model) updateListItems() {
	results := m.store.Search(m.searchQuery)
	m.list.SetItems(newListItems(results))
}

// handleBatchAddProjectKeyMsg 处理批量添加项目状态按键
func (m *Model) handleBatchAddProjectKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case " ":
		// 切换选中状态
		if m.batchCursor >= 0 && m.batchCursor < len(m.batchProjects) {
			m.batchSelected[m.batchCursor] = !m.batchSelected[m.batchCursor]
		}
	case "up", "k":
		if m.batchCursor > 0 {
			m.batchCursor--
		}
	case "down", "j":
		if m.batchCursor < len(m.batchProjects)-1 {
			m.batchCursor++
		}
	case "enter":
		// 确认添加选中的项目
		selected := m.filterBatchSelected()
		if len(selected) > 0 {
			for _, path := range selected {
				alias := filepath.Base(path)
				if alias == "" || alias == "/" || alias == "\\" {
					alias = path
				}
				project := models.Project{
					ID:        generateID(),
					Path:      path,
					Alias:     alias,
					CreatedAt: time.Now(),
				}
				m.store.Add(project)
			}
			m.syncListItems()
		}
		m.state = StateList
		m.batchProjects = nil
		m.batchSelected = nil
		m.batchCursor = 0
	case "esc":
		m.state = StateList
		m.batchProjects = nil
		m.batchSelected = nil
		m.batchCursor = 0
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}
