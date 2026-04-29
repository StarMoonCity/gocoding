package components

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Modal 模态框接口
type Modal interface {
	// Update 处理消息，返回是否已处理
	Update(msg tea.Msg) (tea.Cmd, bool)
	// View 渲染模态框
	View(width, height int) string
	// Overlay 将模态框叠加到内容上
	Overlay(content string, width, height int) string
}

// ModalManager 模态框管理器 - 支持模态框栈
type ModalManager struct {
	stack     []Modal
	focused   int
	width     int
	height    int
	hoverButton int
	mouseEnabled bool
}

// NewModalManager 创建新的模态框管理器
func NewModalManager() *ModalManager {
	return &ModalManager{
		stack:   make([]Modal, 0),
		focused: -1,
	}
}

// Push 推送模态框到栈
func (m *ModalManager) Push(modal Modal) {
	m.stack = append(m.stack, modal)
	m.focused = len(m.stack) - 1
}

// Pop 弹出栈顶模态框
func (m *ModalManager) Pop() {
	if len(m.stack) == 0 {
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
	m.focused = len(m.stack) - 1
}

// Close 关闭所有模态框
func (m *ModalManager) Close() {
	m.stack = nil
	m.focused = -1
}

// HasModal 是否有模态框
func (m *ModalManager) HasModal() bool {
	return len(m.stack) > 0
}

// Top 获取栈顶模态框
func (m *ModalManager) Top() Modal {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

// Update 更新模态框
func (m *ModalManager) Update(msg tea.Msg) (tea.Cmd, bool) {
	if len(m.stack) == 0 {
		return nil, false
	}

	modal := m.stack[len(m.stack)-1]
	return modal.Update(msg)
}

// View 渲染栈顶模态框
func (m *ModalManager) View(width, height int) string {
	if len(m.stack) == 0 {
		return ""
	}
	return m.stack[len(m.stack)-1].View(width, height)
}

// Overlay 将模态框叠加到内容上
func (m *ModalManager) Overlay(content string, width, height int) string {
	if len(m.stack) == 0 {
		return content
	}
	modalView := m.stack[len(m.stack)-1].View(width, height)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modalView)
}

// SetHoverButton 设置悬停按钮索引
func (m *ModalManager) SetHoverButton(idx int) {
	m.hoverButton = idx
}

// HoverButton 获取悬停按钮索引
func (m *ModalManager) HoverButton() int {
	return m.hoverButton
}

// SetMouseEnabled 设置鼠标启用状态
func (m *ModalManager) SetMouseEnabled(enabled bool) {
	m.mouseEnabled = enabled
}

// IsMouseEnabled 获取鼠标启用状态
func (m *ModalManager) IsMouseEnabled() bool {
	return m.mouseEnabled
}
