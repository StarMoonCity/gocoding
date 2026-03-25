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
