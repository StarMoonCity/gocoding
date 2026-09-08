package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Graphite Teal 暗色主题 - 语义化颜色系统
var (
	// 背景层次（从深到浅）
	BackgroundDeep    = lipgloss.Color("#090C10") // 最深层，强调层次区域
	Background        = lipgloss.Color("#0D1117") // 应用全屏画布背景
	BackgroundSurface = lipgloss.Color("#161B22") // 页头、页脚、弹窗、输入框、面板
	BackgroundLight   = lipgloss.Color("#21262D") // 按钮、徽章、次级表面
	BackgroundHover   = lipgloss.Color("#292F36") // 鼠标悬停状态

	// 主色调 - 品牌青色系
	PrimaryColor    = lipgloss.Color("#56B6C2") // 页面标题、焦点边框、主要操作
	PrimaryColorAlt = lipgloss.Color("#3FA7B3") // 品牌层次
	PrimaryDim      = lipgloss.Color("#397C84") // 普通强调边框和分隔线
	PrimaryDark     = lipgloss.Color("#294F55") // 弱分隔线

	// 状态色
	SuccessColor = lipgloss.Color("#3FB950") // 成功、可用、激活
	WarningColor = lipgloss.Color("#D29922") // 警告、最近打开提示
	ErrorColor   = lipgloss.Color("#F85149") // 错误和危险操作
	ErrorDim     = lipgloss.Color("#A83A3A") // 错误背景或弱错误边框

	// 文字色（层次分明）
	Foreground    = lipgloss.Color("#E6EDF3") // 主要文字
	ForegroundDim = lipgloss.Color("#B1BAC4") // 次级文字、路径、说明
	SecondaryText = lipgloss.Color("#8B949E") // 帮助文本、非关键元数据
	MutedText     = lipgloss.Color("#6E7681") // 占位符、不可用内容

	// 选中/激活状态
	SelectedBg     = lipgloss.Color("#1B3438") // 键盘当前项整行背景
	SelectedBgAlt  = lipgloss.Color("#23454A") // 选中且激活的组合状态
	SelectedBorder = lipgloss.Color("#56B6C2") // 当前项左侧指示条

	// 悬停状态
	HoverBg     = lipgloss.Color("#292F36") // 鼠标悬停背景
	HoverBorder = lipgloss.Color("#397C84") // 鼠标悬停指示条

	// 输入框焦点
	FocusBorder = lipgloss.Color("#56B6C2") // 焦点边框
	FocusBg     = lipgloss.Color("#14272A") // 焦点背景

	// 功能色
	AccentCyan    = lipgloss.Color("#58A6FF") // 搜索、信息状态
	AccentMagenta = lipgloss.Color("#B060B0") // 保留变量，不再作为页面主色
	AccentGold    = lipgloss.Color("#D29922") // 少量重要标记

	// IDE 品牌色 - 仅用于 IDE 标识和状态点
	IDEClaudeColor   = lipgloss.Color("#E6A370") // Claude
	IDEVSCodeColor   = lipgloss.Color("#4A90A4") // VSCode
	IDEOpenCodeColor = lipgloss.Color("#4CAF50") // OpenCode
	IDECodexColor    = lipgloss.Color("#9B6B9B") // Codex

	// 语义化背景（消息框、按钮等使用）
	ErrorBgDim   = lipgloss.Color("#1A0D10") // 错误弱背景
	SuccessBgDim = lipgloss.Color("#0A1F12") // 成功弱背景
	TipBgDim     = lipgloss.Color("#0A1A1F") // 提示弱背景
	ActiveBgDim  = lipgloss.Color("#1B3438") // 激活弱背景（同 SelectedBg）
)

// 边框定义 - 圆角单线
var (
	// 标准圆角边框
	NeonBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// 双线边框（保留兼容）
	DoubleNeonBorder = lipgloss.Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	// 强调边框（保留兼容）
	TopHeavyBorder = lipgloss.Border{
		Top:         "═",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
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
			Background(BackgroundSurface).
			Foreground(Foreground).
			Border(NeonBorder).
			BorderForeground(PrimaryDim).
			Padding(1, 2)

	// 选中项样式 - 整行高亮
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(SelectedBg).
				Bold(true).
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
			Foreground(SecondaryText)

	// 信息文本样式
	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryText)

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

	// 焦点输入框样式
	FocusedInputBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, false, false, true).
				BorderForeground(FocusBorder).
				Background(FocusBg)

	// 普通输入框样式
	InputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, true).
			BorderForeground(PrimaryDark).
			Background(BackgroundSurface)
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

// 消息框样式
var (
	// ErrorBoxStyle 错误消息框
	ErrorBoxStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(ErrorBgDim).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(ErrorDim).
			Width(40)

	// TipBoxStyle 提示消息框
	TipBoxStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(TipBgDim).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(PrimaryDim).
			Width(40)

	// SuccessBoxStyle 成功消息框
	SuccessBoxStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Background(SuccessBgDim).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(SuccessColor).
			Width(40)
)

// HelpKeyStyle 快捷键样式
var HelpKeyStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// HelpKey 分类样式
var (
	HelpKeyNavStyle    = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)  // 导航
	HelpKeyActionStyle = lipgloss.NewStyle().Foreground(AccentGold).Bold(true)    // 操作
	HelpKeyDangerStyle = lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)    // 危险
	HelpKeyQuitStyle   = lipgloss.NewStyle().Foreground(SecondaryText).Bold(true) // 退出
	HelpKeySearchStyle = lipgloss.NewStyle().Foreground(AccentCyan).Bold(true)    // 搜索
)

// TitleBarStyle 标题栏样式
var TitleBarStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	Background(BackgroundSurface)

// TitleTextStyle 标题文字样式
var TitleTextStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// SeparatorStyle 分隔线样式
var SeparatorStyle = lipgloss.NewStyle().
	Foreground(PrimaryDark)

// SeparatorHighlightStyle 高亮分隔线
var SeparatorHighlightStyle = lipgloss.NewStyle().
	Foreground(PrimaryDim)

// FeaturedBadgeStyle 突出显示徽章
var FeaturedBadgeStyle = lipgloss.NewStyle().
	Foreground(BackgroundDeep).
	Background(AccentGold).
	Padding(0, 1).
	MarginLeft(1).
	Bold(true)

// ActiveBadgeStyle 激活状态徽章
var ActiveBadgeStyle = lipgloss.NewStyle().
	Foreground(SuccessColor).
	Background(ActiveBgDim).
	Padding(0, 1).
	MarginLeft(1).
	Bold(true)

// ListItemStyle 列表项样式
var ListItemStyle = lipgloss.NewStyle().
	Foreground(Foreground)

// SelectedListItemStyle 选中列表项样式
var SelectedListItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Bold(true)

// ButtonStyle 按钮样式
var ButtonStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundLight).
	Padding(0, 2).
	Margin(0, 1)

// ButtonHoverStyle 按钮悬停样式
var ButtonHoverStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundHover).
	Padding(0, 2).
	Margin(0, 1).
	Bold(true)

// DangerButtonStyle 危险按钮样式
var DangerButtonStyle = lipgloss.NewStyle().
	Foreground(ErrorColor).
	Background(ErrorBgDim).
	Padding(0, 2).
	Margin(0, 1)

// DangerButtonHoverStyle 危险按钮悬停/聚焦样式
var DangerButtonHoverStyle = lipgloss.NewStyle().
	Foreground(Background).
	Background(ErrorColor).
	Padding(0, 2).
	Margin(0, 1).
	Bold(true)

// SuccessButtonStyle 成功按钮样式
var SuccessButtonStyle = lipgloss.NewStyle().
	Foreground(SuccessColor).
	Background(SuccessBgDim).
	Padding(0, 2).
	Margin(0, 1)

// 对话框样式 - 圆角单线边框
var DialogStyle = lipgloss.NewStyle().
	Background(BackgroundSurface).
	Foreground(Foreground).
	Border(NeonBorder).
	BorderForeground(PrimaryDim).
	Padding(1, 3)

// Provider 表单样式
var ProviderFormStyle = lipgloss.NewStyle().
	Background(BackgroundSurface).
	Foreground(Foreground).
	Border(NeonBorder).
	BorderForeground(PrimaryDim).
	Padding(1, 2)

// Provider 列表项样式 - 普通状态
var ProviderItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundSurface).
	Padding(0, 1)

// Provider 选中项样式
var ProviderSelectedItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(SelectedBg).
	Bold(true).
	BorderLeft(true).
	BorderLeftForeground(SelectedBorder).
	Padding(0, 1)

// Provider 激活项样式 - 不使用大面积绿色背景
var ProviderActiveItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundSurface).
	Padding(0, 1)

// Provider 输入框样式
var ProviderInputStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundSurface).
	Border(lipgloss.RoundedBorder(), false, false, false, true).
	BorderForeground(PrimaryDark).
	Padding(0, 1)

// Provider 焦点输入框样式
var ProviderFocusedInputStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(FocusBg).
	Border(lipgloss.RoundedBorder(), false, false, false, true).
	BorderForeground(FocusBorder).
	Padding(0, 1)

// HeaderStyle 顶部标题栏样式
var HeaderStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Background(BackgroundSurface).
	Bold(true).
	Padding(0, 1)

// FooterStyle 底部状态栏样式
var FooterStyle = lipgloss.NewStyle().
	Foreground(SecondaryText).
	Background(BackgroundSurface).
	Padding(0, 1)

// 悬停样式（动态使用）
func HoverStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(Foreground).
		Background(HoverBg).
		BorderLeft(true).
		BorderLeftForeground(HoverBorder).
		Padding(0, 1)
}

// ActiveStyle 激活状态样式
func ActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(Background).
		Background(PrimaryColor).
		Bold(true).
		Padding(0, 1)
}

// FillBackground 强制将每一行填充到指定宽度，确保无透明区域
// 对已渲染的带 ANSI 样式的内容，逐行补齐背景色空格
func FillBackground(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return content
	}

	bgStyle := lipgloss.NewStyle().Background(Background)
	lines := strings.Split(content, "\n")

	// 补齐行数到 height
	for len(lines) < height {
		lines = append(lines, "")
	}
	// 截断多余行
	if len(lines) > height {
		lines = lines[:height]
	}

	for i, line := range lines {
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			padding := strings.Repeat(" ", width-lineWidth)
			lines[i] = line + bgStyle.Render(padding)
		}
	}

	return strings.Join(lines, "\n")
}
