package components

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gocoding/internal/models"
)

func keyRunesMsg(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func keyMsg(kt tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: kt}
}

func newTestStore(projects ...models.Project) *models.ProjectStore {
	store := models.NewProjectStore()
	for _, p := range projects {
		store.Add(p)
	}
	return store
}

// TestProjectListPage_KeyBindingsFollowDocumentation 验证 r=重命名、e=编辑描述与文档一致
func TestProjectListPage_KeyBindingsFollowDocumentation(t *testing.T) {
	store := newTestStore(models.Project{
		ID: "id1", Path: "/tmp/demo", Alias: "demo", CreatedAt: time.Now(),
	})
	p := NewProjectListPage(store)

	// r 打开重命名对话框并预填当前路径/名称
	_, consumed := p.Update(keyRunesMsg("r"))
	if !consumed {
		t.Fatal("r should be consumed")
	}
	if p.state != ProjectStateRename {
		t.Fatalf("r should open rename dialog, got state %v", p.state)
	}
	if p.editingID != "id1" {
		t.Errorf("rename should target selected project, got %s", p.editingID)
	}
	if p.input.Value() != "/tmp/demo" || p.secondaryInput.Value() != "demo" {
		t.Errorf("rename dialog should prefill path/alias, got %q / %q", p.input.Value(), p.secondaryInput.Value())
	}

	// 返回列表后 e 打开描述编辑器
	_, _ = p.Update(keyMsg(tea.KeyEsc))
	if p.state != ProjectStateList {
		t.Fatalf("esc should return to list, got %v", p.state)
	}
	_, consumed = p.Update(keyRunesMsg("e"))
	if !consumed {
		t.Fatal("e should be consumed")
	}
	if p.state != ProjectStateEditDescription {
		t.Fatalf("e should open description editor, got state %v", p.state)
	}
}

// TestProjectListPage_AddRejectsDuplicatePath 验证重复路径不允许添加
func TestProjectListPage_AddRejectsDuplicatePath(t *testing.T) {
	existing := t.TempDir()
	store := newTestStore(models.Project{
		ID: "id1", Path: existing, Alias: "dup", CreatedAt: time.Now(),
	})
	p := NewProjectListPage(store)

	_, _ = p.Update(keyRunesMsg("n"))
	if p.state != ProjectStateAdd {
		t.Fatalf("n should open add dialog, got %v", p.state)
	}

	p.input.SetValue(existing)
	p.secondaryInput.SetValue("new-alias")
	_, _ = p.Update(keyMsg(tea.KeyEnter))

	if p.state != ProjectStateAdd {
		t.Errorf("add dialog should stay open on duplicate, got %v", p.state)
	}
	if !strings.Contains(p.errMsg, "已存在") {
		t.Errorf("expected duplicate-path error, got %q", p.errMsg)
	}
	if store.Len() != 1 {
		t.Errorf("duplicate path must not be added, store len = %d", store.Len())
	}

	// 错误必须在对话框中可见
	view := p.View(80, 30)
	if !strings.Contains(view, p.errMsg) {
		t.Error("validation error should be rendered inside the add dialog")
	}
}

// TestProjectListPage_AddFallsBackToPathBasename 验证名称为空时自动使用路径目录名
func TestProjectListPage_AddFallsBackToPathBasename(t *testing.T) {
	store := models.NewProjectStore()
	p := NewProjectListPage(store)

	_, _ = p.Update(keyRunesMsg("n"))
	tmpDir := t.TempDir()
	p.input.SetValue(tmpDir)
	p.secondaryInput.SetValue("")
	_, _ = p.Update(keyMsg(tea.KeyEnter))

	if p.state != ProjectStateList {
		t.Fatalf("add should return to list, got %v", p.state)
	}
	if store.Len() != 1 {
		t.Fatalf("expected 1 project after add, got %d", store.Len())
	}
	want := filepath.Base(tmpDir)
	if got := store.Projects[0].Alias; got != want {
		t.Errorf("alias should fall back to basename %q, got %q", want, got)
	}
}

// TestProjectListPage_AddShowsValidationErrorInDialog 验证路径校验错误在对话框中可见
func TestProjectListPage_AddShowsValidationErrorInDialog(t *testing.T) {
	p := NewProjectListPage(models.NewProjectStore())

	_, _ = p.Update(keyRunesMsg("n"))
	p.input.SetValue("/definitely/not/exist/xyz-12345")
	p.secondaryInput.SetValue("any")
	_, _ = p.Update(keyMsg(tea.KeyEnter))

	if p.errMsg == "" {
		t.Fatal("expected validation error for nonexistent path")
	}
	view := p.View(80, 30)
	if !strings.Contains(view, p.errMsg) {
		t.Error("validation error should be rendered inside the add dialog")
	}
}

// TestProjectListPage_AddErrorClearsOnTyping 验证输入变化后错误提示被清除
func TestProjectListPage_AddErrorClearsOnTyping(t *testing.T) {
	p := NewProjectListPage(models.NewProjectStore())

	_, _ = p.Update(keyRunesMsg("n"))
	p.input.SetValue("/definitely/not/exist/xyz-12345")
	p.secondaryInput.SetValue("any")
	_, _ = p.Update(keyMsg(tea.KeyEnter))
	if p.errMsg == "" {
		t.Fatal("precondition: error should be set")
	}

	// 继续输入一个字符，错误应被清除
	_, _ = p.Update(keyRunesMsg("a"))
	if p.errMsg != "" {
		t.Errorf("typing should clear the stale error, got %q", p.errMsg)
	}
}

// TestProjectListPage_DetailShowsNeverOpened 验证从未打开的项目显示“从未打开”
func TestProjectListPage_DetailShowsNeverOpened(t *testing.T) {
	store := newTestStore(models.Project{
		ID: "id1", Path: "/tmp/demo", Alias: "demo", CreatedAt: time.Now(),
	})
	p := NewProjectListPage(store)
	p.SetSize(80, 30)

	_, _ = p.Update(keyRunesMsg("v"))
	if p.state != ProjectStateViewDetail {
		t.Fatalf("v should open detail view, got %v", p.state)
	}

	view := p.View(80, 30)
	if !strings.Contains(view, "从未打开") {
		t.Error("detail view should show 从未打开 for never-opened project")
	}
	if strings.Contains(view, "0001-01-01") {
		t.Error("detail view should not render zero time as 0001-01-01")
	}
}

// TestProjectListPage_DelegateTruncatesLongAlias 验证超长别名被截断，不撑破列表布局
func TestProjectListPage_DelegateTruncatesLongAlias(t *testing.T) {
	longAlias := strings.Repeat("很长的项目名称", 10) // 140 个字符宽度

	m := list.New(nil, projectListDelegate{searchQuery: new(string)}, 30, 10)
	m.SetSize(30, 10)
	m.SetItems(newListItems([]models.Project{
		{ID: "id1", Path: "/tmp/long", Alias: longAlias, CreatedAt: time.Now()},
	}))

	var b bytes.Buffer
	d := projectListDelegate{searchQuery: new(string)}
	d.Render(&b, m, 0, m.Items()[0])

	if width := lipgloss.Width(b.String()); width > 30 {
		t.Errorf("rendered row width %d exceeds list width 30", width)
	}
	if !strings.Contains(b.String(), "…") {
		t.Error("truncated alias should end with an ellipsis")
	}
}

// TestProjectListPage_DeleteClampsSelection 验证删除最后一项后选中项被钳制
func TestProjectListPage_DeleteClampsSelection(t *testing.T) {
	store := newTestStore(
		models.Project{ID: "id1", Path: "/tmp/a", Alias: "a", CreatedAt: time.Now()},
		models.Project{ID: "id2", Path: "/tmp/b", Alias: "b", CreatedAt: time.Now()},
	)
	p := NewProjectListPage(store)
	p.SetSize(80, 30)

	// 选中最后一项并删除
	p.list.Select(1)
	_, _ = p.Update(keyRunesMsg("d"))
	if p.state != ProjectStateDeleteConfirm {
		t.Fatalf("d should open delete confirm, got %v", p.state)
	}
	p.hoverButton = 1 // 切到“是”后回车确认删除
	_, _ = p.Update(keyMsg(tea.KeyEnter))

	if p.state != ProjectStateList {
		t.Fatalf("enter should close delete confirm, got %v", p.state)
	}
	if store.Len() != 1 {
		t.Fatalf("expected 1 project after delete, got %d", store.Len())
	}
	// 删除后列表变短，选中项应被钳制到最后一个有效项
	if idx := p.list.Index(); idx >= len(p.list.Items()) {
		t.Errorf("selection index %d out of range (items %d)", idx, len(p.list.Items()))
	}
	if p.safeGetSelectedProject() == nil {
		t.Error("selection should still resolve to a project")
	}
}
