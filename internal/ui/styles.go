package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor   = lipgloss.Color("#00D4FF")   // 淡蓝色标题
	SecondaryColor = lipgloss.Color("#808890")
	SuccessColor   = lipgloss.Color("#27AE60")
	WarningColor   = lipgloss.Color("#F39C12")
	ErrorColor     = lipgloss.Color("#E74C3C")
	Background     = lipgloss.Color("#1E272E")   // 深灰/炭黑背景
	Foreground     = lipgloss.Color("#F5F6FA")
	SelectedBg     = lipgloss.Color("#2C3E50")
	SelectedBorder = lipgloss.Color("#00D4FF")   // 选中项左侧高亮
)

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Background(Background).
			Bold(true).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(Background).
			Bold(true).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(SelectedBg).
				BorderLeft(true).
				BorderLeftForeground(SelectedBorder).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
				Foreground(Foreground).
				Background(Background).
				Padding(0, 1)

	// 路径样式 - 灰色斜体，更紧凑
	PathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			MarginLeft(2)

	HelpStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(Background)

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(Background).
			Bold(true)

	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Background(Background).
			Foreground(Foreground).
			Padding(1, 2)

	IDESelectedStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(SuccessColor).
			Padding(0, 2).
			MarginRight(1).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Background(Background)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(SelectedBg).
			Padding(0, 1)

	// 列表项别名样式
	ItemAliasStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Bold(true)

	// 列表项路径样式
	ItemPathStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Italic(true)
)

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
