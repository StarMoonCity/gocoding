package components

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestToastManager_AutoExpireWithoutInput 验证 toast 无需按键也能自动消失
func TestToastManager_AutoExpireWithoutInput(t *testing.T) {
	m := NewToastManager()

	m.Show("hello", string(ToastTip), 50*time.Millisecond)

	// 可见时 Update 应返回过期计时命令
	var tick tea.Cmd
	if tick = m.Update(nil); tick == nil {
		t.Fatal("expected a tick cmd when toast is visible")
	}

	// 计时触发后 toast 被移除
	if cmd := m.Update(tick()); cmd != nil {
		t.Fatalf("expected nil cmd after last toast expired, got %v", cmd)
	}
	if m.visible {
		t.Error("toast should be invisible after expiry")
	}
	if len(m.toasts) != 0 {
		t.Errorf("toasts should be empty, got %d", len(m.toasts))
	}
}

// TestToastManager_MultipleToastsExpireSequentially 验证多条 toast 依次过期
func TestToastManager_MultipleToastsExpireSequentially(t *testing.T) {
	m := NewToastManager()

	m.Show("first", string(ToastTip), time.Millisecond)
	m.Show("second", string(ToastTip), time.Millisecond)

	// 第一条 toast 的过期计时
	cmd := m.Update(nil)
	if cmd == nil {
		t.Fatal("expected tick for the first toast")
	}
	// 第一条过期后，应自动为第二条继续计时
	if cmd = m.Update(cmd()); cmd == nil {
		t.Fatal("expected tick for the second toast")
	}
	// 第二条过期后全部清空
	if cmd = m.Update(cmd()); cmd != nil {
		t.Fatalf("expected no more ticks, got %v", cmd)
	}
	if m.visible || len(m.toasts) != 0 {
		t.Errorf("expected empty toasts, visible=%v len=%d", m.visible, len(m.toasts))
	}
}
