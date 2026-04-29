package ui

import (
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
)

// 霓虹风格颜色系统 - 深色背景 + 青色主色调
var (
	// 背景层次（从深到浅）
	BackgroundDeep   = lipgloss.Color("#0A0E14") // 最深背景（近黑）
	Background       = lipgloss.Color("#11151C") // 主背景
	BackgroundSurface = lipgloss.Color("#1A2332") // 面板/卡片
	BackgroundLight  = lipgloss.Color("#243246") // 表面层
	BackgroundHover  = lipgloss.Color("#2D3B4F") // 悬停状态

	// 主色调 - 柔和青色系（护眼）
	PrimaryColor    = lipgloss.Color("#4DB6AC") // 柔和青色
	PrimaryColorAlt = lipgloss.Color("#26A69A") // 备用青色
	PrimaryDim      = lipgloss.Color("#00796B") // 深青色
	PrimaryDark     = lipgloss.Color("#004D40") // 最深青
	PrimaryGlow     = lipgloss.Color("#4DB6AC") // 发光色（同主色）

	// 状态色 - 柔和色系
	SuccessColor = lipgloss.Color("#4CAF50") // 柔和绿色（非霓虹）
	WarningColor = lipgloss.Color("#FFB300") // 琥珀色
	ErrorColor   = lipgloss.Color("#E57373") // 柔和红色
	ErrorDim     = lipgloss.Color("#B06060") // 深红

	// 文字色（层次分明）
	Foreground     = lipgloss.Color("#FFFFFF") // 主文字（纯白）
	ForegroundDim  = lipgloss.Color("#E0E6ED") // 次级文字
	SecondaryText  = lipgloss.Color("#8892A0") // 中性文字
	MutedText      = lipgloss.Color("#4A5568") // 淡化文字

	// 选中/激活状态（霓虹发光效果）
	SelectedBg      = lipgloss.Color("#0D3B4D") // 选中背景
	SelectedBgAlt   = lipgloss.Color("#1A4A5E") // 选中背景备用
	SelectedBorder  = lipgloss.Color("#00E5FF") // 选中边框（发光）

	// 悬停状态
	HoverBg         = lipgloss.Color("#1E3A5F") // 悬停背景
	HoverBorder     = lipgloss.Color("#00D4FF") // 悬停边框

	// 输入框焦点
	FocusBorder     = lipgloss.Color("#00E5FF") // 焦点边框
	FocusBg         = lipgloss.Color("#0D2A36") // 焦点背景

	// 特殊效果色 - 柔和色系（舒适护眼）
	AccentCyan      = lipgloss.Color("#5BBFBA") // 柔和青色（非霓虹）
	AccentMagenta   = lipgloss.Color("#B060B0") // 柔和洋红（非霓虹）
	AccentGold      = lipgloss.Color("#D4A574") // 柔和金色（非霓虹）

	// IDE 品牌色 - 柔和色系
	IDEClaudeColor   = lipgloss.Color("#E6A370") // 柔和橙色 - Claude
	IDEVSCodeColor   = lipgloss.Color("#4A90A4") // 柔和蓝色 - VSCode
	IDEOpenCodeColor = lipgloss.Color("#4CAF50") // 柔和绿色 - OpenCode
	IDECodexColor    = lipgloss.Color("#9B6B9B") // 柔和洋红 - Codex
)

// 边框定义 - 霓虹风格
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

	// 双线边框（对话框用）
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

	// 强调边框（顶部亮色）
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
	// 标题样式 - 霓虹发光效果
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

	// 选中项样式 - 霓虹发光
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

	// 焦点输入框样式 - 霓虹发光边框
	FocusedInputBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), false, false, false, true).
				BorderForeground(FocusBorder).
				Background(FocusBg)

	// 普通输入框样式
	InputBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), false, false, false, true).
			BorderForeground(SecondaryText).
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

// 消息框样式 - 霓虹风格
var (
	// ErrorBoxStyle 错误消息框 - 霓虹红边框
	ErrorBoxStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Background(lipgloss.Color("#1A0D10")).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(ErrorColor).
			Width(40)

	// TipBoxStyle 提示消息框 - 霓虹青边框
	TipBoxStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(lipgloss.Color("#0A1A1F")).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(PrimaryDim).
			Width(40)

	// SuccessBoxStyle 成功消息框 - 霓虹绿边框
	SuccessBoxStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Background(lipgloss.Color("#0A1F12")).
			Padding(1, 2).
			MarginBottom(1).
			Border(NeonBorder).
			BorderForeground(SuccessColor).
			Width(40)
)

// HelpKeyStyle 快捷键样式 - 高亮显示 (默认青色)
var HelpKeyStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// HelpKey 分类样式 - 按操作类型颜色编码
var (
	HelpKeyNavStyle    = lipgloss.NewStyle().Foreground(PrimaryColor).Bold(true)   // 导航: 青色
	HelpKeyActionStyle = lipgloss.NewStyle().Foreground(AccentGold).Bold(true)     // 操作: 金色
	HelpKeyDangerStyle = lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)     // 危险: 红色
	HelpKeyQuitStyle   = lipgloss.NewStyle().Foreground(SecondaryText).Bold(true)  // 退出: 灰色
	HelpKeySearchStyle = lipgloss.NewStyle().Foreground(AccentCyan).Bold(true)     // 搜索: 亮青
)

// TitleBarStyle 标题栏样式
var TitleBarStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true).
	Background(BackgroundDeep)

// TitleTextStyle 标题文字样式 - 发光效果
var TitleTextStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// SeparatorStyle 分隔线样式
var SeparatorStyle = lipgloss.NewStyle().
	Foreground(PrimaryDark)

// SeparatorHighlightStyle 高亮分隔线（用于区段分隔）
var SeparatorHighlightStyle = lipgloss.NewStyle().
	Foreground(PrimaryDim)

// FeaturedBadgeStyle 突出显示徽章（金色，用于重要/活跃标记）
var FeaturedBadgeStyle = lipgloss.NewStyle().
	Foreground(BackgroundDeep).
	Background(AccentGold).
	Padding(0, 1).
	MarginLeft(1).
	Bold(true)

// ActiveBadgeStyle 激活状态徽章（柔和绿色背景）
var ActiveBadgeStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(lipgloss.Color("#2D5A4A")).
	Padding(0, 1).
	MarginLeft(1).
	Bold(true)

// ListItemStyle 列表项样式
var ListItemStyle = lipgloss.NewStyle().
	Foreground(Foreground)

// SelectedListItemStyle 选中列表项样式 - 霓虹高亮
var SelectedListItemStyle = lipgloss.NewStyle().
	Foreground(PrimaryColor).
	Bold(true)

// ButtonStyle 按钮样式
var ButtonStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(BackgroundLight).
	Padding(0, 2).
	Margin(0, 1)

// ButtonHoverStyle 按钮悬停样式 - 霓虹效果
var ButtonHoverStyle = lipgloss.NewStyle().
	Foreground(Background).
	Background(PrimaryColor).
	Padding(0, 2).
	Margin(0, 1).
	Bold(true)

// DangerButtonStyle 危险按钮样式
var DangerButtonStyle = lipgloss.NewStyle().
	Foreground(ErrorColor).
	Background(lipgloss.Color("#1A0D10")).
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
	Background(lipgloss.Color("#0A1F12")).
	Padding(0, 2).
	Margin(0, 1)

// 对话框样式 - 双线霓虹边框
var DialogStyle = lipgloss.NewStyle().
	Background(BackgroundSurface).
	Foreground(Foreground).
	Border(DoubleNeonBorder).
	BorderForeground(PrimaryColor).
	Padding(1, 3)

// Provider 表单样式 - 强调边框
var ProviderFormStyle = lipgloss.NewStyle().
	Background(BackgroundSurface).
	Foreground(Foreground).
	Border(TopHeavyBorder).
	BorderForeground(PrimaryColor).
	Padding(1, 2)

// Provider 列表项样式
var ProviderItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(Background).
	Padding(0, 1)

// Provider 选中项样式 - 霓虹发光
var ProviderSelectedItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(SelectedBg).
	Foreground(PrimaryColor).
	BorderLeft(true).
	BorderLeftForeground(PrimaryColor).
	Padding(0, 1)

// Provider 激活项样式 - 柔和发光
var ProviderActiveItemStyle = lipgloss.NewStyle().
	Foreground(Foreground).
	Background(lipgloss.Color("#2D5A4A")). // 柔和绿色背景
	Bold(true).
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
	Background(BackgroundDeep).
	Bold(true).
	Padding(0, 1)

// FooterStyle 底部状态栏样式
var FooterStyle = lipgloss.NewStyle().
	Foreground(SecondaryText).
	Background(BackgroundLight).
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

// ideColor 返回 IDE 对应的品牌色
func ideColor(ideType models.IDEType) lipgloss.Color {
	switch ideType {
	case models.IDEClaudeCode:
		return IDEClaudeColor
	case models.IDEVSCode:
		return IDEVSCodeColor
	case models.IDEOpenCode:
		return IDEOpenCodeColor
	case models.IDECodexCLI:
		return IDECodexColor
	default:
		return PrimaryColor
	}
}
