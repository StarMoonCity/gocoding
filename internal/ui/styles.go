package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// 暗色主题色彩系统
var (
	// 背景层次
	BackgroundDeep   = lipgloss.Color("#0F1419") // 最深背景
	Background       = lipgloss.Color("#1A2332") // 主背景
	BackgroundLight  = lipgloss.Color("#242D3D") // 表面/卡片背景

	// 主色调
	PrimaryColor   = lipgloss.Color("#00D4FF")   // 青色主色
	PrimaryDim     = lipgloss.Color("#0099CC")   // 次级主色
	AccentColor    = lipgloss.Color("#FFD700")   // 金色强调（选中/高亮）

	// 状态色
	SuccessColor   = lipgloss.Color("#22C55E")
	WarningColor   = lipgloss.Color("#F59E0B")
	ErrorColor     = lipgloss.Color("#EF4444")

	// 文字色
	Foreground     = lipgloss.Color("#F8FAFC")   // 主文字
	SecondaryText  = lipgloss.Color("#94A3B8")   // 次级文字
	MutedText      = lipgloss.Color("#64748B")   // 淡化文字

	// 选中状态
	SelectedBg     = lipgloss.Color("#2D3A4F")
	SelectedBorder = lipgloss.Color("#00D4FF")
)

// 通用样式定义
var (
	// 标题样式
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(Background).
			Bold(true).
			Padding(0, 1)

	// 面板/对话框样式
	SurfaceStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Foreground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	// 选中项样式
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(SelectedBg).
				BorderLeft(true).
				BorderLeftForeground(SelectedBorder).
				Padding(0, 1)

	// 普通项样式
	NormalItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(Background).
				Padding(0, 1)

	// 帮助文本样式
	HelpStyle = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Background(Background)

	// 信息文本样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Background(Background)

	// IDE选中样式
	IDESelectedStyle = lipgloss.NewStyle().
			Foreground(Background).
			Background(SuccessColor).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	// 状态栏样式
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(BackgroundLight).
			Padding(0, 1)

	// 列表项别名样式
	ItemAliasStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Bold(true)

	// 列表项路径样式
	ItemPathStyle = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Italic(true)

	// 徽章/标签样式
	BadgeStyle = lipgloss.NewStyle().
			Foreground(SecondaryText).
			Background(BackgroundLight).
			Padding(0, 1).
			MarginLeft(1)

	// 焦点输入框样式（左侧高亮边框）
	FocusedInputBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, false, false, true).
				BorderForeground(PrimaryColor)

	// 普通输入框样式
	InputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, true).
			BorderForeground(SecondaryText)
)

// GetStatusStyle 根据可用性返回状态样式
func GetStatusStyle(available bool) lipgloss.Style {
	if available {
		return lipgloss.NewStyle().
			Foreground(SuccessColor).
			Background(Background)
	}
	return lipgloss.NewStyle().
		Foreground(WarningColor).
		Background(Background)
}

// ErrorBoxStyle 错误消息框样式
var ErrorBoxStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(lipgloss.Color("#2C1810")).
			Padding(1, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ErrorColor).
			Width(40)

// TipBoxStyle 提示消息框样式
var TipBoxStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(lipgloss.Color("#0F1926")).
			Padding(1, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryDim).
			Width(40)

// SuccessBoxStyle 成功消息框样式
var SuccessBoxStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Background(lipgloss.Color("#0F2618")).
			Padding(1, 2).
			MarginBottom(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SuccessColor).
			Width(40)

// HelpKeyStyle 快捷键样式
var HelpKeyStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Background(Background).
			Bold(true)
