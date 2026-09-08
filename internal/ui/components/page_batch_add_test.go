package components

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"gocoding/internal/models"
)

// TestBatchAddPage_QConsumed 验证 q 在批量添加页被消费并返回退出命令
func TestBatchAddPage_QConsumed(t *testing.T) {
	p := NewBatchAddPage(models.NewProjectStore())

	cmd, consumed := p.Update(keyRunesMsg("q"))
	if !consumed {
		t.Error("q should be consumed by the batch add page")
	}
	if cmd == nil {
		t.Error("q should return a quit command")
	}
}

// TestBatchAddPage_EnterWithoutSelectionShowsToast 验证未选择任何项目时给出提示
func TestBatchAddPage_EnterWithoutSelectionShowsToast(t *testing.T) {
	app := NewAppModel(models.NewProjectStore(), models.NewModelProviderStore())
	p := NewBatchAddPage(app.Store())
	p.SetApp(app)
	app.SwitchPage(p)
	p.OnActivate()

	_, consumed := p.Update(keyMsg(tea.KeyEnter))
	if !consumed {
		t.Error("enter should be consumed by the batch add page")
	}
	if !app.toastManager.visible || len(app.toastManager.toasts) == 0 {
		t.Fatal("expected a toast when nothing is selected")
	}
	if got := app.toastManager.toasts[0].Message; got != "未选择任何项目" {
		t.Errorf("expected '未选择任何项目' toast, got %q", got)
	}
	// 没有选中任何项目时不应切换页面
	if app.currentPage != p {
		t.Error("page should stay on batch add when nothing is selected")
	}
	if app.Store().Len() != 0 {
		t.Errorf("no projects should be added, got %d", app.Store().Len())
	}
}
