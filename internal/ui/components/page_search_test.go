package components

import (
	"testing"
	"time"

	"gocoding/internal/models"
)

// TestSearchPage_ClampsSelectionAfterFilter 验证过滤后选中项被钳制到有效范围
func TestSearchPage_ClampsSelectionAfterFilter(t *testing.T) {
	store := models.NewProjectStore()
	store.Add(models.Project{ID: "id1", Path: "/p/alpha", Alias: "alpha", CreatedAt: time.Now()})
	store.Add(models.Project{ID: "id2", Path: "/p/beta", Alias: "beta", CreatedAt: time.Now()})
	store.Add(models.Project{ID: "id3", Path: "/p/gamma", Alias: "gamma", CreatedAt: time.Now()})

	p := NewSearchPage(store)
	p.SetSize(80, 30)
	// 光标放到最后一个项目
	p.list.Select(2)

	// 输入只匹配第一个项目的查询：al
	_, _ = p.Update(keyRunesMsg("a"))
	_, consumed := p.Update(keyRunesMsg("l"))
	if !consumed {
		t.Fatal("search input should be consumed")
	}

	if len(p.list.Items()) != 1 {
		t.Fatalf("expected 1 match, got %d", len(p.list.Items()))
	}
	if idx := p.list.Index(); idx != 0 {
		t.Errorf("selection should clamp to 0, got %d", idx)
	}
	if p.getSelectedProject() == nil {
		t.Error("clamped selection should still resolve to a project")
	}
	if p.getSelectedProject().ID != "id1" {
		t.Errorf("expected selected project id1, got %q", p.getSelectedProject().ID)
	}
}
