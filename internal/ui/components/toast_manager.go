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
	Message   string
	Type      ToastType
	Duration  time.Duration
	CreatedAt time.Time
	ExpiresAt time.Time
}

// ToastManager Toast 通知管理器
type ToastManager struct {
	toasts      []Toast
	width       int
	visible     bool
	tickPending bool
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
		ExpiresAt: time.Now().Add(duration),
	}
	m.toasts = append(m.toasts, toast)
	m.visible = true
}

// toastExpiredMsg toast 过期消息
type toastExpiredMsg struct{}

// Update 驱动 toast 过期：通过 tea.Tick 在主循环内移除，避免 goroutine 直接修改状态
func (m *ToastManager) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case toastExpiredMsg:
		if len(m.toasts) > 0 {
			m.toasts = m.toasts[1:]
		}
		m.visible = len(m.toasts) > 0
		m.tickPending = false
	}

	// 有可见 toast 且没有挂起的计时器时，安排过期
	if m.tickPending || !m.visible || len(m.toasts) == 0 {
		return nil
	}
	remaining := time.Until(m.toasts[0].ExpiresAt)
	if remaining < 0 {
		remaining = 0
	}
	m.tickPending = true
	return tea.Tick(remaining, func(time.Time) tea.Msg {
		return toastExpiredMsg{}
	})
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
	m.tickPending = false
}
