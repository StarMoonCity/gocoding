package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewProjectStore(t *testing.T) {
	store := NewProjectStore()
	if store == nil {
		t.Fatal("NewProjectStore() returned nil")
	}
	if store.Projects == nil {
		t.Error("Projects slice should not be nil")
	}
	if store.index == nil {
		t.Error("index map should not be nil")
	}
	if store.Len() != 0 {
		t.Errorf("expected length 0, got %d", store.Len())
	}
}

func TestProjectStore_Add(t *testing.T) {
	store := NewProjectStore()
	project := Project{
		ID:        "test-id-1",
		Path:      "/tmp/test-project",
		Alias:     "Test Project",
		CreatedAt: time.Now(),
	}

	store.Add(project)

	if store.Len() != 1 {
		t.Errorf("expected length 1, got %d", store.Len())
	}

	// 测试添加后能正确获取
	retrieved := store.Get("test-id-1")
	if retrieved == nil {
		t.Fatal("Get returned nil for existing ID")
	}
	if retrieved.ID != project.ID {
		t.Errorf("expected ID %s, got %s", project.ID, retrieved.ID)
	}
	if retrieved.Alias != project.Alias {
		t.Errorf("expected Alias %s, got %s", project.Alias, retrieved.Alias)
	}
}

func TestProjectStore_AddMultiple(t *testing.T) {
	store := NewProjectStore()

	for i := 0; i < 5; i++ {
		store.Add(Project{
			ID:    string(rune('a' + i)),
			Path:  "/tmp/test" + string(rune('0'+i)),
			Alias: "Project " + string(rune('A'+i)),
		})
	}

	if store.Len() != 5 {
		t.Errorf("expected length 5, got %d", store.Len())
	}
}

func TestProjectStore_Remove(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
	store.Add(Project{ID: "id3", Path: "/tmp/p3", Alias: "P3"})

	store.Remove("id2")

	if store.Len() != 2 {
		t.Errorf("expected length 2 after remove, got %d", store.Len())
	}

	if store.Get("id2") != nil {
		t.Error("removed project should not be found")
	}

	// 确认其他项目仍在
	if store.Get("id1") == nil {
		t.Error("id1 should still exist")
	}
	if store.Get("id3") == nil {
		t.Error("id3 should still exist")
	}
}

func TestProjectStore_RemoveBoundary(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*ProjectStore)
		removeID    string
		wantLen     int
		wantGet     map[string]bool // id -> 是否应存在
	}{
		{
			name: "删除不存在的项目",
			setup: func(s *ProjectStore) {
				s.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
			},
			removeID: "non-existent",
			wantLen:  1,
			wantGet:  map[string]bool{"id1": true},
		},
		{
			name: "删除唯一项目",
			setup: func(s *ProjectStore) {
				s.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
			},
			removeID: "id1",
			wantLen:  0,
			wantGet:  map[string]bool{"id1": false},
		},
		{
			name: "删除中间项目",
			setup: func(s *ProjectStore) {
				s.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
				s.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
				s.Add(Project{ID: "id3", Path: "/tmp/p3", Alias: "P3"})
			},
			removeID: "id2",
			wantLen:  2,
			wantGet:  map[string]bool{"id1": true, "id2": false, "id3": true},
		},
		{
			name: "删除首项目",
			setup: func(s *ProjectStore) {
				s.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
				s.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
			},
			removeID: "id1",
			wantLen:  1,
			wantGet:  map[string]bool{"id1": false, "id2": true},
		},
		{
			name: "删除尾项目",
			setup: func(s *ProjectStore) {
				s.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
				s.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
			},
			removeID: "id2",
			wantLen:  1,
			wantGet:  map[string]bool{"id1": true, "id2": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewProjectStore()
			tt.setup(store)
			store.Remove(tt.removeID)

			if store.Len() != tt.wantLen {
				t.Errorf("after Remove(%q): Len() = %d, want %d", tt.removeID, store.Len(), tt.wantLen)
			}
			for id, shouldExist := range tt.wantGet {
				got := store.Get(id)
				if shouldExist && got == nil {
					t.Errorf("after Remove(%q): Get(%q) = nil, want non-nil", tt.removeID, id)
				}
				if !shouldExist && got != nil {
					t.Errorf("after Remove(%q): Get(%q) should be nil", tt.removeID, id)
				}
			}
		})
	}
}

func TestProjectStore_Get(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})

	// 测试获取存在的项目
	result := store.Get("id1")
	if result == nil {
		t.Fatal("Get returned nil for existing ID")
	}
	if result.ID != "id1" {
		t.Errorf("expected ID id1, got %s", result.ID)
	}

	// 测试获取不存在的项目
	result = store.Get("non-existent")
	if result != nil {
		t.Error("Get should return nil for non-existent ID")
	}
}

func TestProjectStore_Update(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "Old Alias"})

	store.Update("id1", "New Alias", "/tmp/new-path")

	p := store.Get("id1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Alias != "New Alias" {
		t.Errorf("expected Alias 'New Alias', got '%s'", p.Alias)
	}
	if p.Path != "/tmp/new-path" {
		t.Errorf("expected Path '/tmp/new-path', got '%s'", p.Path)
	}
}

func TestProjectStore_UpdatePartial(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "Original"})

	// 只更新别名
	store.Update("id1", "New Alias", "")

	p := store.Get("id1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Alias != "New Alias" {
		t.Errorf("expected Alias 'New Alias', got '%s'", p.Alias)
	}
	if p.Path != "/tmp/p1" {
		t.Errorf("Path should remain unchanged, got '%s'", p.Path)
	}

	// 只更新路径
	store.Update("id1", "", "/tmp/updated-path")
	p = store.Get("id1")
	if p.Path != "/tmp/updated-path" {
		t.Errorf("expected Path '/tmp/updated-path', got '%s'", p.Path)
	}
	if p.Alias != "New Alias" {
		t.Errorf("Alias should remain 'New Alias', got '%s'", p.Alias)
	}
}

func TestProjectStore_UpdateNonExistent(t *testing.T) {
	store := NewProjectStore()
	// 更新不存在的项目不应 panic
	store.Update("non-existent", "alias", "/tmp/path")
}

func TestProjectStore_UpdateDescription(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})

	store.UpdateDescription("id1", "A great project")

	p := store.Get("id1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Description != "A great project" {
		t.Errorf("expected Description 'A great project', got '%s'", p.Description)
	}
}

func TestProjectStore_UpdateDescriptionNonExistent(t *testing.T) {
	store := NewProjectStore()
	// 更新不存在项目的描述不应 panic
	store.UpdateDescription("non-existent", "desc")
}

func TestProjectStore_GetByIndex(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
	store.Add(Project{ID: "id3", Path: "/tmp/p3", Alias: "P3"})

	// 有效索引
	p := store.GetByIndex(1)
	if p == nil {
		t.Fatal("GetByIndex(1) returned nil")
	}
	if p.ID != "id2" {
		t.Errorf("expected ID id2, got %s", p.ID)
	}

	// 无效索引
	if store.GetByIndex(-1) != nil {
		t.Error("GetByIndex(-1) should return nil")
	}
	if store.GetByIndex(100) != nil {
		t.Error("GetByIndex(100) should return nil")
	}
}

func TestProjectStore_Len(t *testing.T) {
	store := NewProjectStore()
	if store.Len() != 0 {
		t.Errorf("expected 0, got %d", store.Len())
	}

	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	if store.Len() != 1 {
		t.Errorf("expected 1, got %d", store.Len())
	}

	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
	if store.Len() != 2 {
		t.Errorf("expected 2, got %d", store.Len())
	}

	store.Remove("id1")
	if store.Len() != 1 {
		t.Errorf("expected 1 after remove, got %d", store.Len())
	}
}

func TestProjectStore_SortByLastOpened(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		projects  []Project
		wantOrder []string // 期望的 ID 顺序（从新到旧）
	}{
		{
			name: "三项目按时间倒序排列",
			projects: []Project{
				{ID: "id1", Path: "/tmp/p1", Alias: "P1", LastOpened: now.Add(-2 * time.Hour)},
				{ID: "id2", Path: "/tmp/p2", Alias: "P2", LastOpened: now},
				{ID: "id3", Path: "/tmp/p3", Alias: "P3", LastOpened: now.Add(-1 * time.Hour)},
			},
			wantOrder: []string{"id2", "id3", "id1"},
		},
		{
			name: "相同时间",
			projects: []Project{
				{ID: "id1", Path: "/tmp/p1", Alias: "P1", LastOpened: now},
				{ID: "id2", Path: "/tmp/p2", Alias: "P2", LastOpened: now},
			},
			wantOrder: []string{"id1", "id2"},
		},
		{
			name: "单项目",
			projects: []Project{
				{ID: "id1", Path: "/tmp/p1", Alias: "P1", LastOpened: now},
			},
			wantOrder: []string{"id1"},
		},
		{
			name:      "空存储",
			projects:  []Project{},
			wantOrder: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewProjectStore()
			for _, p := range tt.projects {
				store.Add(p)
			}
			store.SortByLastOpened()

			if len(store.Projects) != len(tt.wantOrder) {
				t.Errorf("expected %d projects, got %d", len(tt.wantOrder), len(store.Projects))
				return
			}
			for i, wantID := range tt.wantOrder {
				if store.Projects[i].ID != wantID {
					t.Errorf("position %d: expected ID %s, got %s", i, wantID, store.Projects[i].ID)
				}
			}
		})
	}
}

func TestProjectStore_Search(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/gocoding", Alias: "Go Project", Description: "A Go project"})
	store.Add(Project{ID: "id2", Path: "/tmp/webapp", Alias: "Web App", Description: "A web application"})
	store.Add(Project{ID: "id3", Path: "/tmp/golang-tools", Alias: "Go Tools", Description: "Go developer tools"})

	tests := []struct {
		name     string
		query    string
		wantLen  int
		wantIDs  []string
	}{
		{
			name:    "按别名搜索小写",
			query:   "go",
			wantLen: 2,
			wantIDs: []string{"id1", "id3"},
		},
		{
			name:    "按别名搜索大写",
			query:   "GO",
			wantLen: 2,
			wantIDs: []string{"id1", "id3"},
		},
		{
			name:    "按路径搜索",
			query:   "webapp",
			wantLen: 1,
			wantIDs: []string{"id2"},
		},
		{
			name:    "按描述搜索",
			query:   "developer",
			wantLen: 1,
			wantIDs: []string{"id3"},
		},
		{
			name:    "空搜索返回所有",
			query:   "",
			wantLen: 3,
			wantIDs: []string{"id1", "id2", "id3"},
		},
		{
			name:    "无匹配",
			query:   "python",
			wantLen: 0,
			wantIDs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := store.Search(tt.query)
			if len(results) != tt.wantLen {
				t.Errorf("Search(%q) returned %d results, want %d", tt.query, len(results), tt.wantLen)
			}
			for i, wantID := range tt.wantIDs {
				if i < len(results) && results[i].ID != wantID {
					t.Errorf("Search(%q) result[%d] ID = %s, want %s", tt.query, i, results[i].ID, wantID)
				}
			}
		})
	}
}

func TestProjectStore_GetMostRecentlyOpened(t *testing.T) {
	store := NewProjectStore()

	// 空存储
	if store.GetMostRecentlyOpened() != nil {
		t.Error("GetMostRecentlyOpened on empty store should return nil")
	}

	now := time.Now()
	store.Add(Project{ID: "id1", Path: "/tmp/old", Alias: "Old", LastOpened: now.Add(-1 * time.Hour)})
	store.Add(Project{ID: "id2", Path: "/tmp/new", Alias: "New", LastOpened: now})

	mostRecent := store.GetMostRecentlyOpened()
	if mostRecent == nil {
		t.Fatal("GetMostRecentlyOpened returned nil")
	}
	if mostRecent.ID != "id2" {
		t.Errorf("expected most recent ID id2, got %s", mostRecent.ID)
	}
}

func TestProjectStore_GetIndexByProject(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
	store.Add(Project{ID: "id3", Path: "/tmp/p3", Alias: "P3"})

	if store.GetIndexByProject("id2") != 1 {
		t.Errorf("expected index 1 for id2, got %d", store.GetIndexByProject("id2"))
	}

	if store.GetIndexByProject("id1") != 0 {
		t.Errorf("expected index 0 for id1, got %d", store.GetIndexByProject("id1"))
	}

	if store.GetIndexByProject("non-existent") != -1 {
		t.Error("expected -1 for non-existent ID")
	}
}

func TestProject_UpdateLastOpened(t *testing.T) {
	before := time.Now().Add(-1 * time.Hour)
	p := Project{
		ID:         "id1",
		Path:       "/tmp/p1",
		Alias:      "P1",
		OpenCount:  5,
		LastOpened: before,
	}

	p.UpdateLastOpened()

	if p.LastOpened.Before(before) {
		t.Error("LastOpened should be updated to current time")
	}
	if p.OpenCount != 6 {
		t.Errorf("expected OpenCount 6, got %d", p.OpenCount)
	}
}

func TestValidatePath(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "有效目录路径",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "空路径",
			path:    "",
			wantErr: true,
		},
		{
			name:    "不存在的路径",
			path:    "/tmp/non-existent-path-12345",
			wantErr: true,
		},
		{
			name:    "文件路径（非目录）",
			path:    tmpFile,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if tt.wantErr && err == nil {
				t.Errorf("ValidatePath(%q) expected error, got nil", tt.path)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ValidatePath(%q) expected no error, got %v", tt.path, err)
			}
		})
	}
}

func TestProjectStore_ValidatePath(t *testing.T) {
	store := NewProjectStore()
	tmpDir := t.TempDir()

	// Store 的 ValidatePath 应该委托给全局 ValidatePath
	if err := store.ValidatePath(tmpDir); err != nil {
		t.Errorf("Store.ValidatePath should not return error for existing directory: %v", err)
	}

	if err := store.ValidatePath(""); err == nil {
		t.Error("Store.ValidatePath should return error for empty path")
	}
}

func TestProjectStore_Load_Save(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "projects.json")

	// 创建并保存存储
	store1 := NewProjectStore()
	store1.Add(Project{
		ID:          "id1",
		Path:        "/tmp/p1",
		Alias:       "Project 1",
		Description: "A test project",
		OpenCount:   3,
		CreatedAt:   time.Now(),
		LastOpened:  time.Now(),
	})

	if err := store1.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 加载到新存储
	store2 := NewProjectStore()
	if err := store2.Load(path); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if store2.Len() != 1 {
		t.Errorf("expected loaded store to have 1 project, got %d", store2.Len())
	}

	p := store2.Get("id1")
	if p == nil {
		t.Fatal("loaded project not found")
	}
	if p.Alias != "Project 1" {
		t.Errorf("expected Alias 'Project 1', got '%s'", p.Alias)
	}
	if p.Description != "A test project" {
		t.Errorf("expected Description 'A test project', got '%s'", p.Description)
	}
	if p.OpenCount != 3 {
		t.Errorf("expected OpenCount 3, got %d", p.OpenCount)
	}
}

func TestProjectStore_LoadNonExistent(t *testing.T) {
	store := NewProjectStore()
	// 加载不存在的文件应该返回 nil（按代码逻辑处理）
	if err := store.Load("/tmp/non-existent-file-12345.json"); err != nil {
		t.Errorf("Load should not return error for non-existent file: %v", err)
	}
	if store.Len() != 0 {
		t.Errorf("store should remain empty after loading non-existent file, got %d", store.Len())
	}
}

func TestProjectStore_SaveInvalidDir(t *testing.T) {
	store := NewProjectStore()
	// 保存到无效目录应该返回错误
	err := store.Save("/root/impossible-dir/projects.json")
	if err == nil {
		t.Error("Save should return error for invalid directory")
	}
}

func TestProject_JSONSerialization(t *testing.T) {
	now := time.Now()
	p := Project{
		ID:          "test-id",
		Path:        "/tmp/test",
		Alias:       "Test",
		Description: "A test project",
		CreatedAt:   now,
		LastOpened:  now,
		OpenCount:   5,
	}

	// 序列化
	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化
	var p2 Project
	if err := json.Unmarshal(data, &p2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if p2.ID != p.ID {
		t.Errorf("ID mismatch: expected %s, got %s", p.ID, p2.ID)
	}
	if p2.Alias != p.Alias {
		t.Errorf("Alias mismatch: expected %s, got %s", p.Alias, p2.Alias)
	}
	if p2.OpenCount != p.OpenCount {
		t.Errorf("OpenCount mismatch: expected %d, got %d", p.OpenCount, p2.OpenCount)
	}
}

func TestProjectStore_IndexAfterAddRemove(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})
	store.Add(Project{ID: "id3", Path: "/tmp/p3", Alias: "P3"})

	// 删除中间元素后，索引应该正确重建
	store.Remove("id2")

	// 添加新项目，索引应该正确
	store.Add(Project{ID: "id4", Path: "/tmp/p4", Alias: "P4"})

	// 验证所有项目仍可正确获取
	if store.Get("id1") == nil {
		t.Error("id1 should be accessible")
	}
	if store.Get("id3") == nil {
		t.Error("id3 should be accessible")
	}
	if store.Get("id4") == nil {
		t.Error("id4 should be accessible")
	}
	if store.Get("id2") != nil {
		t.Error("id2 should not exist")
	}
}

func TestProjectStore_RebuildIndex(t *testing.T) {
	store := NewProjectStore()
	store.Add(Project{ID: "id1", Path: "/tmp/p1", Alias: "P1"})
	store.Add(Project{ID: "id2", Path: "/tmp/p2", Alias: "P2"})

	// 手动触发索引重建
	store.rebuildIndex()

	// 验证索引仍然正确
	if store.Get("id1") == nil {
		t.Error("id1 should be accessible after rebuildIndex")
	}
	if store.Get("id2") == nil {
		t.Error("id2 should be accessible after rebuildIndex")
	}
}
