package ui

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"gocoding/internal/models"
)

var updateSnapshots = flag.Bool("update", false, "update golden files")

func normalizeView(view string) string {
	// 替换时间戳格式: 2006-01-02 15:04:05 -> [TIMESTAMP]
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
	view = re.ReplaceAllString(view, "[TIMESTAMP]")

	// 替换 OpenCount: "打开 X 次" -> "打开 N 次"
	re = regexp.MustCompile(`打开 \d 次`)
	view = re.ReplaceAllString(view, "打开 N 次")

	return view
}

func snapshotFilePath(name string) string {
	return filepath.Join("testdata", name+".golden")
}

func TestViewSnapshot(t *testing.T) {
	// 固定测试时间，避免时间相关的差异
	fixedTime := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)

	project1 := models.Project{
		ID:          "test-1",
		Path:        "/tmp/test-project-1",
		Alias:       "Test Project 1",
		Description: "A test project",
		OpenCount:   5,
		CreatedAt:   fixedTime,
		LastOpened:  fixedTime.Add(1 * time.Hour),
	}

	project2 := models.Project{
		ID:          "test-2",
		Path:        "/tmp/test-project-2",
		Alias:       "Test Project 2",
		Description: "Another test project",
		OpenCount:   3,
		CreatedAt:   fixedTime,
		LastOpened:  fixedTime,
	}

	tests := []struct {
		name   string
		setup  func(*Model)
		golden string
	}{
		{
			name:   "empty_list",
			setup:  func(m *Model) {},
			golden: "list_empty",
		},
		{
			name: "list_with_one_project",
			setup: func(m *Model) {
				m.store.Add(project1)
				m.list.SetItems(newListItems(m.store.Projects))
			},
			golden: "list_one_project",
		},
		{
			name: "list_with_two_projects",
			setup: func(m *Model) {
				m.store.Add(project1)
				m.store.Add(project2)
				m.list.SetItems(newListItems(m.store.Projects))
			},
			golden: "list_two_projects",
		},
		{
			name: "add_project",
			setup: func(m *Model) {
				m.state = StateAddProject
				m.input.SetValue("/tmp/test-project")
				m.secondaryInput.SetValue("Test Project")
			},
			golden: "add_project",
		},
		{
			name: "provider_list_with_tip",
			setup: func(m *Model) {
				provider := models.ModelProvider{
					ID:                  "test-provider-1",
					Name:                "i 公司",
					BaseURL:             "https://api.test.com",
					APIKey:              "test-key",
					Model:               "test-model",
					Active:              false,
					DisableNonessential: "false",
					DisableNonstreaming: "false",
					EffortLevel:         "low",
				}
				providerStore := models.NewModelProviderStore()
				providerStore.Add(provider)
				m.SetProviderStore(providerStore)
				m.state = StateProviderList
				m.updateProviderListItems()
				// 模拟激活后的提示
				m.tipMsg = "i 公司 配置已生效，无需更新"
			},
			golden: "provider_list_with_tip",
		},
	}

	flag.Parse()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := models.NewProjectStore()
			m := NewModel(store)
			m.SetSize(80, 24)
			tt.setup(m)

			got := normalizeView(m.View())
			goldenPath := snapshotFilePath(tt.golden)

			if *updateSnapshots || !exists(goldenPath) {
				if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
					t.Fatalf("failed to write golden file: %v", err)
				}
				t.Logf("updated golden file: %s", goldenPath)
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			if string(want) != got {
				t.Errorf("snapshot mismatch for %s\nExpected:\n%s\n\nGot:\n%s", tt.name, want, got)
			}
		})
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// TestActivateProviderTip 测试激活供应商配置后的提示功能
func TestActivateProviderTip(t *testing.T) {
	provider := models.ModelProvider{
		ID:                  "test-provider-1",
		Name:                "Test Provider",
		BaseURL:             "https://api.test.com",
		APIKey:             "test-key",
		Model:              "test-model",
		Active:             false,
		DisableNonessential: "false",
		DisableNonstreaming: "false",
		EffortLevel:        "low",
	}

	providerStore := models.NewModelProviderStore()
	providerStore.Add(provider)

	store := models.NewProjectStore()
	m := NewModel(store)
	m.SetProviderStore(providerStore)
	m.SetSize(80, 24)

	// 切换到提供商列表状态
	m.state = StateProviderList
	m.updateProviderListItems()

	// 直接设置 tipMsg 来模拟激活后的提示（不实际调用激活逻辑，避免写入配置文件）
	m.tipMsg = "已激活 Test Provider"

	// 验证视图渲染
	view := m.View()
	if view == "" {
		t.Error("视图不应为空")
	}

	// 验证 tipMsg 包含在视图中
	if !strings.Contains(view, "已激活 Test Provider") {
		t.Errorf("视图应包含提示信息，实际视图: %s", view)
	}

	t.Log("激活提示视图测试通过")
}

// TestActivateProviderError 测试激活失败时的错误提示
func TestActivateProviderError(t *testing.T) {
	// 创建一个空的 provider store（没有 provider）
	store := models.NewProjectStore()
	m := NewModel(store)
	m.SetProviderStore(models.NewModelProviderStore())
	m.SetSize(80, 24)

	m.state = StateProviderList

	// 直接设置 errMsg 来模拟错误场景（不实际执行激活逻辑）
	m.errMsg = "激活失败: 配置无效"

	// 验证视图渲染
	view := m.View()
	if view == "" {
		t.Error("视图不应为空")
	}

	// 验证 errMsg 包含在视图中
	if !strings.Contains(view, "激活失败") {
		t.Errorf("视图应包含错误信息，实际视图: %s", view)
	}

	t.Log("激活错误视图测试通过")
}
