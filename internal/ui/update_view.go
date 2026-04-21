package ui

import (
	"github.com/charmbracelet/bubbletea"

	"gocoding/internal/models"
)

// handleIDEMenuKeyMsg 处理 IDE 菜单状态按键
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
	case "4":
		if len(m.list.Items()) > 0 {
			return m.openWithIDE(models.IDECodexCLI)
		}
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// handleViewDetailKeyMsg 处理查看详情状态按键
func (m *Model) handleViewDetailKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// handleEditDescriptionKeyMsg 处理编辑描述状态按键
func (m *Model) handleEditDescriptionKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "ctrl+s":
		if current := m.safeGetSelectedProject(); current != nil {
			m.store.UpdateDescription(current.ID, m.ta.Value())
		}
		m.state = StateList
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// handleSearchKeyMsg 处理搜索状态按键
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
		if m.safeGetSelectedProject() != nil {
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
		searchRunes := []rune(m.searchQuery)
		if len(searchRunes) > 0 {
			m.searchQuery = string(searchRunes[:len(searchRunes)-1])
		}
		m.updateListItems()
		return m, nil
	case "ctrl+h":
		// Ctrl+H 也支持删除
		searchRunes := []rune(m.searchQuery)
		if len(searchRunes) > 0 {
			m.searchQuery = string(searchRunes[:len(searchRunes)-1])
		}
		m.updateListItems()
		return m, nil
	case "1":
		if m.safeGetSelectedProject() != nil {
			return m.openWithIDE(models.IDEClaudeCode)
		}
	case "2":
		if m.safeGetSelectedProject() != nil {
			return m.openWithIDE(models.IDEVSCode)
		}
	case "3":
		if m.safeGetSelectedProject() != nil {
			return m.openWithIDE(models.IDEOpenCode)
		}
	case "4":
		if m.safeGetSelectedProject() != nil {
			return m.openWithIDE(models.IDECodexCLI)
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
