package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// 暗色主题色彩系统 - 统一单色调风格
var (
	// 背景层次
	BackgroundDeep  = lipgloss.Color("#0F1419") // 最深背景
	Background      = lipgloss.Color("#1A2332") // 主背景
	BackgroundLight = lipgloss.Color("#243246") // 表面/卡片背景
	BackgroundHover = lipgloss.Color("#363F52") // 悬停状态

	// 主色调 - 统一使用青色系
	PrimaryColor = lipgloss.Color("#00D4FF") // 青色主色
	PrimaryDim   = lipgloss.Color("#0099CC") // 次级主色
	PrimaryDark  = lipgloss.Color("#006B8F") // 深色变体

	// 状态色 - 使用同一色系
	SuccessColor = lipgloss.Color("#10B981") // 绿色（同一色调）
	WarningColor = lipgloss.Color("#F59E0B") // 保留警告色
	ErrorColor   = lipgloss.Color("#EF4444") // 保留错误色

	// 文字色
	Foreground    = lipgloss.Color("#F8FAFC") // 主文字
	SecondaryText = lipgloss.Color("#94A3B8") // 次级文字
	MutedText     = lipgloss.Color("#64748B") // 淡化文字

	// 选中状态 - 使用主色调作为唯一强调
	SelectedBg     = lipgloss.Color("#1E3A5F") // 选中背景（深青色）
	SelectedBorder = lipgloss.Color("#00D4FF") // 选中边框

	// 边框线条样式
	BorderVertical  = "│"
	BorderHorizontal = "─"
	BorderCornerTL  = "┌"
	BorderCornerTR  = "┐"
	BorderCornerBL  = "└"
	BorderCornerBR  = "┘"
	BorderDVertical = "║"
	BorderDHorizontal = "═"
	BorderDCornerTL = "╔"
	BorderDCornerTR = "╗"
	BorderDCornerBL = "╚"
	BorderDCornerBR = "╝"
	BorderTeeLeft   = "├"
	BorderTeeRight  = "┤"
	BorderTeeTop    = "┬"
	BorderTeeBottom = "┴"
	BorderCross     = "┼"
)

// 通用样式定义
var (
	// 标题样式
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	// 面板/对话框样式
	SurfaceStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Foreground).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryDim).
			Padding(1, 2)

	// 选中项样式 - 统一使用主色调
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(SelectedBg).
				BorderLeft(true).
				BorderLeftForeground(PrimaryColor).
				Padding(0, 1)

	// 普通项样式
	NormalItemStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(Background).
			Padding(0, 1)

	// 帮助文本样式
	HelpStyle = lipgloss.NewStyle().
			Foreground(SecondaryText)

	// 信息文本样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryText)

	// IDE选中样式 - 使用成功色
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

	// 焦点输入框样式（左侧高亮边框）- 使用主色调
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
			Foreground(SuccessColor)
	}
	return lipgloss.NewStyle().
		Foreground(WarningColor)
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
	Foreground(PrimaryColor).
	Bold(true)

// TitleBarStyle 标题栏样式
var TitleBarStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	Background(BackgroundDeep)

// TitleTextStyle 标题文字样式
var TitleTextStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// SeparatorStyle 分隔线样式
var SeparatorStyle = lipgloss.NewStyle().
	Foreground(PrimaryDim)

// ListItemStyle 列表项样式
var ListItemStyle = lipgloss.NewStyle().
	Foreground(Foreground)

// SelectedListItemStyle 选中列表项样式
var SelectedListItemStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// ButtonStyle 按钮样式
var ButtonStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundLight).
	Padding(0, 2).
	Margin(0, 1)

// ButtonHoverStyle 按钮悬停样式
var ButtonHoverStyle = lipgloss.NewStyle().
	Foreground(Background).
	Background(PrimaryColor).
	Padding(0, 2).
	Margin(0, 1).
	Bold(true)

// DangerButtonStyle 危险按钮样式
var DangerButtonStyle = lipgloss.NewStyle().
	Foreground(ErrorColor).
	Background(lipgloss.Color("#2C1810")).
	Padding(0, 2).
	Margin(0, 1)

// DangerButtonHoverStyle 危险按钮悬停样式
var DangerButtonHoverStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(ErrorColor).
	Padding(0, 2).
	Margin(0, 1).
	Bold(true)

// SuccessButtonStyle 成功按钮样式
var SuccessButtonStyle = lipgloss.NewStyle().
	Foreground(SuccessColor).
	Background(lipgloss.Color("#0F2618")).
	Padding(0, 2).
	Margin(0, 1)
