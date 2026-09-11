package ui

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var sgrPattern = regexp.MustCompile(`\x1b\[([0-9;]*)m`)

// countUncoloredCells 统计没有背景色生效的可见单元格数量
func countUncoloredCells(content string) int {
	total := 0
	for _, line := range strings.Split(content, "\n") {
		total += countUncoloredLineCells(line)
	}
	return total
}

func countUncoloredLineCells(line string) int {
	count := 0
	hasBg := false
	for i := 0; i < len(line); {
		if loc := sgrPattern.FindStringIndex(line[i:]); loc != nil && loc[0] == 0 {
			seq := line[i : i+loc[1]]
			hasBg = sgrSetsBackground(seq, hasBg)
			i += loc[1]
			continue
		}
		if !hasBg {
			count++
		}
		i++
	}
	return count
}

func TestFillGapsColorsEveryCell(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	// 典型场景：内层只设置前景色，拼接时补齐的空格不带背景色，浅色终端下会透出白底
	title := lipgloss.NewStyle().Foreground(PrimaryColor).Render("模型配置")
	row := lipgloss.NewStyle().Background(BackgroundSurface).Render(strings.Repeat("-", 40))
	content := lipgloss.JoinVertical(lipgloss.Center, title, row)

	if countUncoloredCells(content) == 0 {
		t.Fatal("测试前提不成立：原始内容应存在未着色单元格")
	}

	filled := FillGaps(content, BackgroundSurface)
	if got := countUncoloredCells(filled); got != 0 {
		t.Errorf("FillGaps 后仍有 %d 个未着色单元格", got)
	}
	if stripANSI(filled) != stripANSI(content) {
		t.Error("FillGaps 不应改变可见文本")
	}
}

func TestFillGapsKeepsExistingBackgrounds(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	content := lipgloss.NewStyle().Background(SelectedBg).Render("  已选中  ")
	filled := FillGaps(content, Background)
	if filled != content {
		t.Errorf("已有背景色的内容不应被修改\n原始: %q\n结果: %q", content, filled)
	}
}

func TestSgrSetsBackground(t *testing.T) {
	tests := []struct {
		seq     string
		current bool
		want    bool
	}{
		{"\x1b[0m", true, false},
		{"\x1b[49m", true, false},
		{"\x1b[48;2;13;17;23m", false, true},
		{"\x1b[48;5;236m", false, true},
		{"\x1b[41m", false, true},
		{"\x1b[1;38;2;100;44;20m", false, false}, // 前景色参数不应误判为背景色
		{"\x1b[38;5;44m", false, false},
		{"\x1b[1m", true, true},
	}
	for _, tt := range tests {
		if got := sgrSetsBackground(tt.seq, tt.current); got != tt.want {
			t.Errorf("sgrSetsBackground(%q, %v) = %v, want %v",
				strconv.Quote(tt.seq), tt.current, got, tt.want)
		}
	}
}

func stripANSI(s string) string {
	return strings.ReplaceAll(sgrPattern.ReplaceAllString(s, ""), "\x00", "")
}
