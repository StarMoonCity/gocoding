package components

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/ui"
)

// ToastType Toast 消息类型
type ToastType string

const (
	ToastTip     ToastType = "tip"
	ToastSuccess ToastType = "success"
	ToastError   ToastType = "error"
)

// Toast Toast 消息结构
type Toast struct {
	Message  string
	Type     ToastType
	Duration time.Duration
	CreatedAt time.Time
}

// ToastManager Toast 通知管理器
type ToastManager struct {
	toasts    []Toast
	timer     *time.Timer
	width     int
	visible   bool
}

// NewToastManager 创建新的 Toast 管理器
func NewToastManager() *ToastManager {
	return &ToastManager{
		toasts:  make([]Toast, 0),
		width:   40,
		visible: false,
	}
}

// Show 显示 Toast 消息
func (m *ToastManager) Show(message string, toastType string, duration time.Duration) {
	toast := Toast{
		Message:   message,
		Type:      ToastType(toastType),
		Duration:  duration,
		CreatedAt: time.Now(),
	}
	m.toasts = append(m.toasts, toast)
	m.visible = true

	// 启动定时器
	if m.timer != nil {
		m.timer.Stop()
	}
	m.timer = time.AfterFunc(duration, func() {
		m.toasts = m.toasts[1:]
		if len(m.toasts) == 0 {
			m.visible = false
		}
	})
}

// Update 更新消息（用于定时器）
func (m *ToastManager) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View 渲染 Toast
func (m *ToastManager) View(width int) string {
	if !m.visible || len(m.toasts) == 0 {
		return ""
	}

	toast := m.toasts[0]

	var style lipgloss.Style
	switch toast.Type {
	case ToastSuccess:
		style = ui.SuccessBoxStyle
	case ToastError:
		style = ui.ErrorBoxStyle
	default:
		style = ui.TipBoxStyle
	}

	// 计算位置（底部居中）
	toastWidth := min(40, width-10)
	marginLeft := (width - toastWidth) / 2

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		style.Width(toastWidth).MarginLeft(marginLeft).Render(toast.Message),
	)

	return content
}

// SetWidth 设置渲染宽度
func (m *ToastManager) SetWidth(width int) {
	m.width = width
}

// Clear 清除所有 Toast
func (m *ToastManager) Clear() {
	m.toasts = nil
	m.visible = false
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
}
