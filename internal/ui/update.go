package ui

import (
	"github.com/charmbracelet/bubbletea"
)

// 消息类型常量
const (
	msgTypeErr = "err"
	msgTypeTip = "tip"
)

// clearProviderMessageMsg 自定义消息，用于清除提供商表单的消息
type clearProviderMessageMsg struct {
	msgType string
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case clearProviderMessageMsg:
		// 清除消息
		if msg.msgType == msgTypeErr {
			m.errMsg = ""
		} else if msg.msgType == msgTypeTip {
			m.tipMsg = ""
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		// 更新最后按键（用于调试模式）
		m.lastKey = msg.String()
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
		case StateProviderAdd, StateProviderEdit:
			return m.handleProviderFormKeyMsg(msg)
		case StateProviderDelete:
			return m.handleProviderDeleteKeyMsg(msg)
		case StateBatchAddProject:
			return m.handleBatchAddProjectKeyMsg(msg)
		}
	case tea.MouseMsg:
		return m.handleMouseMsg(msg)
	}

	// 通用组件更新
	switch m.state {
	case StateEditDescription:
		m.ta, cmd = m.ta.Update(msg)
		cmds = append(cmds, cmd)
	case StateViewDetail:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	case StateProviderAdd, StateProviderEdit:
		m.providerFormViewport, cmd = m.providerFormViewport.Update(msg)
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

// handleMouseMsg 处理鼠标消息
func (m *Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	m.mouseEnabled = true

	switch m.state {
	case StateList:
		return m.handleListMouseMsg(msg)
	case StateProviderList:
		return m.handleProviderListMouseMsg(msg)
	case StateDeleteConfirm, StateProviderDelete:
		return m.handleDialogMouseMsg(msg)
	case StateIDEMenu:
		return m.handleIDEMenuMouseMsg(msg)
	case StateViewDetail:
		return m.handleViewDetailMouseMsg(msg)
	case StateBatchAddProject:
		return m.handleBatchAddProjectMouseMsg(msg)
	}

	return m, nil
}

// handleListMouseMsg 处理列表鼠标消息
func (m *Model) handleListMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// 计算列表区域 (排除 header 和 footer)
	listStartY := 3 // headerHeight
	listEndY := m.height - 4 - 1 // footer + margin

	switch msg.Type {
	case tea.MouseLeft:
		// 单击选择列表项
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := msg.Y - listStartY + m.list.Index()
			items := m.list.Items()
			if clickedIndex >= 0 && clickedIndex < len(items) {
				m.list.Select(clickedIndex)
				m.hoverIndex = clickedIndex
			}
		}
	case tea.MouseLeft + 256: // double click
		// 双击打开 IDE 菜单
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := msg.Y - listStartY + m.list.Index()
			items := m.list.Items()
			if clickedIndex >= 0 && clickedIndex < len(items) {
				m.list.Select(clickedIndex)
				m.state = StateIDEMenu
			}
		}
	case tea.MouseMotion:
		// 悬停更新 hoverIndex
		if msg.Y >= listStartY && msg.Y < listEndY {
			hoverIdx := msg.Y - listStartY + m.list.Index()
			items := m.list.Items()
			if hoverIdx >= 0 && hoverIdx < len(items) {
				m.hoverIndex = hoverIdx
			}
		}
	case tea.MouseWheelUp:
		m.list.CursorUp()
	case tea.MouseWheelDown:
		m.list.CursorDown()
	}
	return m, nil
}

// handleProviderListMouseMsg 处理提供商列表鼠标消息
func (m *Model) handleProviderListMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	listStartY := 3
	listEndY := m.height - 4 - 1
	itemHeight := 3 // 2 行内容 + 1 spacing

	switch msg.Type {
	case tea.MouseLeft:
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := (msg.Y - listStartY) / itemHeight
			items := m.providerList.Items()
			if clickedIndex >= 0 && clickedIndex < len(items) {
				m.providerList.Select(clickedIndex)
			}
		}
	case tea.MouseLeft + 256:
		if msg.Y >= listStartY && msg.Y < listEndY {
			m.state = StateProviderEdit
		}
	case tea.MouseMotion:
		if msg.Y >= listStartY && msg.Y < listEndY {
			m.hoverIndex = (msg.Y - listStartY) / itemHeight
		}
	case tea.MouseWheelUp:
		m.providerList.CursorUp()
	case tea.MouseWheelDown:
		m.providerList.CursorDown()
	}
	return m, nil
}

// handleDialogMouseMsg 处理对话框鼠标消息（确认/取消按钮）
func (m *Model) handleDialogMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// 对话框按钮区域（水平居中）
	buttonWidth := 10
	buttonSpacing := 4
	buttonsCount := 2 // 取消、确定

	// 计算按钮区域
	totalWidth := buttonWidth*buttonsCount + buttonSpacing*(buttonsCount-1)
	buttonStartX := (m.width - totalWidth) / 2
	buttonY := m.height - 7

	if msg.Y == buttonY && msg.Type == tea.MouseLeft {
		// 判断点击了哪个按钮
		relX := msg.X - buttonStartX
		if relX >= 0 && relX < totalWidth {
			buttonIndex := relX / (buttonWidth + buttonSpacing)
			if buttonIndex == 0 {
				// 取消
				if m.state == StateDeleteConfirm {
					m.state = StateList
				} else if m.state == StateProviderDelete {
					m.state = StateProviderList
				}
			} else if buttonIndex == 1 {
				// 确定
				return m.handleDialogConfirm()
			}
		}
	}

	// 悬停高亮按钮
	if msg.Type == tea.MouseMotion && msg.Y == buttonY {
		relX := msg.X - buttonStartX
		if relX >= 0 && relX < totalWidth {
			m.hoverButton = relX / (buttonWidth + buttonSpacing)
		} else {
			m.hoverButton = -1
		}
	}

	return m, nil
}

// handleDialogConfirm 处理对话框确认
func (m *Model) handleDialogConfirm() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateDeleteConfirm:
		return m.handleDeleteConfirmKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})
	case StateProviderDelete:
		return m.handleProviderDeleteKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})
	}
	return m, nil
}

// handleIDEMenuMouseMsg 处理 IDE 菜单鼠标消息
func (m *Model) handleIDEMenuMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	menuStartY := 4
	itemHeight := 2

	if msg.Type == tea.MouseLeft && msg.Y >= menuStartY && msg.Y < menuStartY+len(m.ideMenu.options)*itemHeight {
		clickedIndex := (msg.Y - menuStartY) / itemHeight
		if clickedIndex >= 0 && clickedIndex < len(m.ideMenu.options) {
			m.ideMenu.selected = clickedIndex
			selectedIDE := m.ideMenu.options[m.ideMenu.selected].Type
			return m.openWithIDE(selectedIDE)
		}
	}

	if msg.Type == tea.MouseMotion && msg.Y >= menuStartY && msg.Y < menuStartY+len(m.ideMenu.options)*itemHeight {
		m.hoverIndex = (msg.Y - menuStartY) / itemHeight
	}

	return m, nil
}

// handleViewDetailMouseMsg 处理详情视图鼠标消息
func (m *Model) handleViewDetailMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseWheelUp:
		m.viewport.LineUp(3)
	case tea.MouseWheelDown:
		m.viewport.LineDown(3)
	}
	return m, nil
}

// handleBatchAddProjectMouseMsg 处理批量添加项目鼠标消息
func (m *Model) handleBatchAddProjectMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	listStartY := 4
	listEndY := m.height - 7

	switch msg.Type {
	case tea.MouseLeft:
		if msg.Y >= listStartY && msg.Y < listEndY {
			clickedIndex := msg.Y - listStartY
			if clickedIndex >= 0 && clickedIndex < len(m.batchProjects) {
				m.batchCursor = clickedIndex
				// 切换选中状态
				if m.batchSelected[clickedIndex] {
					delete(m.batchSelected, clickedIndex)
				} else {
					m.batchSelected[clickedIndex] = true
				}
			}
		}
	case tea.MouseMotion:
		if msg.Y >= listStartY && msg.Y < listEndY {
			m.hoverIndex = msg.Y - listStartY
		}
	case tea.MouseWheelUp:
		if m.batchCursor > 0 {
			m.batchCursor--
		}
	case tea.MouseWheelDown:
		if m.batchCursor < len(m.batchProjects)-1 {
			m.batchCursor++
		}
	}
	return m, nil
}
