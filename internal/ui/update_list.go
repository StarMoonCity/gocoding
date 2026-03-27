package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"gocoding/internal/models"
)

// handleListKeyMsg 处理列表状态按键
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
	case "b":
		// 批量添加项目
		m.state = StateBatchAddProject
		m.loadBatchProjects()
		return m, nil
	case "e":
		// 重命名选中项目
		if current := m.safeGetSelectedProject(); current != nil {
			m.state = StateRenameProject
			m.editingProjectID = current.ID
			m.input.SetValue(current.Path)
			m.input.Placeholder = "输入项目路径"
			m.input.Focus()
			m.inputFocus = FocusPath
			m.secondaryInput.SetValue(current.Alias)
			m.secondaryInput.Placeholder = "输入项目名称"
			m.secondaryInput.Blur()
		}
		return m, textinput.Blink
	case "d":
		if m.safeGetSelectedProject() != nil {
			m.state = StateDeleteConfirm
		}
	case "v":
		// 切换详情视图
		if len(m.list.Items()) > 0 {
			m.showDetails = !m.showDetails
			m.updateViewport()
		}
	case "enter":
		if m.safeGetSelectedProject() != nil {
			m.state = StateIDEMenu
			for _, opt := range m.ideMenu.options {
				m.ideMenu.available[opt.Type] = m.ideExec.IsIDEAvailable(opt.Type)
			}
		}
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
	case "p":
		// 进入模型配置列表
		m.state = StateProviderList
		return m, nil
	case "ctrl+c", "ctrl+q", "q":
		return m, tea.Quit
	}
	return m, nil
}

// openWithIDE 使用指定 IDE 打开项目
func (m *Model) openWithIDE(ideType models.IDEType) (tea.Model, tea.Cmd) {
	idx := m.list.Index()
	current := m.store.GetByIndex(idx)
	if current == nil {
		m.state = StateList
		return m, nil
	}
	current.UpdateLastOpened()

	if err := m.ideExec.OpenProject(current, ideType); err != nil {
		m.errMsg = err.Error()
		// 3秒后清除错误消息
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return clearProviderMessageMsg{msgType: msgTypeErr}
		})
	}

	// 按最后打开时间排序
	m.store.SortByLastOpened()

	// 重新渲染列表
	m.syncListItems()
	// 定位到最近打开的项目（第一个）
	m.list.Select(0)

	m.state = StateList
	return m, nil
}
