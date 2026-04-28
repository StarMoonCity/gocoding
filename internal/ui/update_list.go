package ui

import (
	"errors"
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
	case "4":
		if m.safeGetSelectedProject() != nil {
			return m.openWithIDE(models.IDECodexCLI)
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

// openWithIDE 使用指定 IDE 打开项目（异步）
func (m *Model) openWithIDE(ideType models.IDEType) (tea.Model, tea.Cmd) {
	return m, m.openProjectCmd(ideType)
}

// openProjectCmd 创建异步打开项目的 Cmd（带超时保护）
func (m *Model) openProjectCmd(ideType models.IDEType) tea.Cmd {
	return func() tea.Msg {
		idx := m.list.Index()
		current := m.store.GetByIndex(idx)
		if current == nil {
			return openProjectResultMsg{err: errors.New("未选择项目")}
		}

		// 使用超时包装 IDE 启动
		done := make(chan error, 1)
		go func() {
			done <- m.ideExec.OpenProject(current, ideType)
		}()

		select {
		case err := <-done:
			if err != nil {
				return openProjectResultMsg{err: err}
			}
			// 更新项目并排序
			current.UpdateLastOpened()
			m.store.SortByLastOpened()
			m.syncListItems()
			m.list.Select(0)
			return openProjectResultMsg{}
		case <-time.After(30 * time.Second):
			return openProjectResultMsg{err: errors.New("打开超时（30秒）")}
		}
	}
}
