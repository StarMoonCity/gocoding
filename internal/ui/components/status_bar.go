package components

import (
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/ui"
)

// StatusBar 状态栏组件
type StatusBar struct {
	left    string
	center  string
	right   string
	style   lipgloss.Style
}

// NewStatusBar 创建新的状态栏
func NewStatusBar() *StatusBar {
	return &StatusBar{
		style: ui.StatusBarStyle,
	}
}

// Set 设置状态栏内容
func (m *StatusBar) Set(left, center, right string) {
	m.left = left
	m.center = center
	m.right = right
}

// View 渲染状态栏
func (m *StatusBar) View(width int) string {
	if m.left == "" && m.center == "" && m.right == "" {
		return ""
	}

	// 计算各部分宽度
	padding := 1
	separator := m.style.Render("│")

	leftWidth := lipgloss.Width(m.left)
	centerWidth := lipgloss.Width(m.center)
	rightWidth := lipgloss.Width(m.right)

	// 动态调整
	available := width - leftWidth - centerWidth - rightWidth - padding*6
	if available < 0 {
		available = 0
	}

	// 构建状态栏
	leftPart := m.style.PaddingLeft(padding).Render(m.left)
	centerPart := m.style.Render(m.center)
	rightPart := m.style.PaddingRight(padding).Render(m.right)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftPart,
		separator,
		lipgloss.Place(available, 1, lipgloss.Center, lipgloss.Center, centerPart),
		separator,
		rightPart,
	)
}

// SetStyle 设置状态栏样式
func (m *StatusBar) SetStyle(style lipgloss.Style) {
	m.style = style
}
