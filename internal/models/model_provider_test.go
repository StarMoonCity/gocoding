package models

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

func TestNewModelProviderStore(t *testing.T) {
	store := NewModelProviderStore()
	if store == nil {
		t.Fatal("NewModelProviderStore() returned nil")
	}
	if store.Providers == nil {
		t.Error("Providers slice should not be nil")
	}
	if store.index == nil {
		t.Error("index map should not be nil")
	}
	if store.Len() != 0 {
		t.Errorf("expected length 0, got %d", store.Len())
	}
}

func TestModelProviderStore_Add(t *testing.T) {
	store := NewModelProviderStore()
	provider := ModelProvider{
		ID:        "provider-1",
		Name:      "MiniMax",
		BaseURL:   "https://api.minimax.chat",
		APIKey:    "test-key",
		Model:     "MiniMax-M2.7",
		Active:    false,
		CreatedAt: time.Now(),
	}

	store.Add(provider)

	if store.Len() != 1 {
		t.Errorf("expected length 1, got %d", store.Len())
	}

	retrieved := store.Get("provider-1")
	if retrieved == nil {
		t.Fatal("Get returned nil for existing ID")
	}
	if retrieved.Name != "MiniMax" {
		t.Errorf("expected Name 'MiniMax', got '%s'", retrieved.Name)
	}
	if retrieved.BaseURL != "https://api.minimax.chat" {
		t.Errorf("expected BaseURL 'https://api.minimax.chat', got '%s'", retrieved.BaseURL)
	}
}

func TestModelProviderStore_AddMultiple(t *testing.T) {
	store := NewModelProviderStore()

	providers := []ModelProvider{
		{ID: "id1", Name: "Provider 1", BaseURL: "https://p1.example.com", APIKey: "key1", Model: "model1"},
		{ID: "id2", Name: "Provider 2", BaseURL: "https://p2.example.com", APIKey: "key2", Model: "model2"},
		{ID: "id3", Name: "Provider 3", BaseURL: "https://p3.example.com", APIKey: "key3", Model: "model3"},
	}

	for _, p := range providers {
		store.Add(p)
	}

	if store.Len() != 3 {
		t.Errorf("expected length 3, got %d", store.Len())
	}
}

func TestModelProviderStore_Remove(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
	store.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
	store.Add(ModelProvider{ID: "id3", Name: "P3", BaseURL: "https://p3.com", APIKey: "key", Model: "m3"})

	store.Remove("id2")

	if store.Len() != 2 {
		t.Errorf("expected length 2 after remove, got %d", store.Len())
	}

	if store.Get("id2") != nil {
		t.Error("removed provider should not be found")
	}

	if store.Get("id1") == nil {
		t.Error("id1 should still exist")
	}
	if store.Get("id3") == nil {
		t.Error("id3 should still exist")
	}
}

func TestModelProviderStore_RemoveBoundary(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*ModelProviderStore)
		removeID string
		wantLen  int
		wantGet  map[string]bool
	}{
		{
			name: "删除不存在的配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
			},
			removeID: "non-existent",
			wantLen: 1,
			wantGet: map[string]bool{"id1": true},
		},
		{
			name: "删除唯一配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
			},
			removeID: "id1",
			wantLen:  0,
			wantGet:  map[string]bool{"id1": false},
		},
		{
			name: "删除中间配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
				s.Add(ModelProvider{ID: "id3", Name: "P3", BaseURL: "https://p3.com", APIKey: "key", Model: "m3"})
			},
			removeID: "id2",
			wantLen:  2,
			wantGet:  map[string]bool{"id1": true, "id2": false, "id3": true},
		},
		{
			name: "删除首配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
			},
			removeID: "id1",
			wantLen:  1,
			wantGet:  map[string]bool{"id1": false, "id2": true},
		},
		{
			name: "删除尾配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
			},
			removeID: "id2",
			wantLen:  1,
			wantGet:  map[string]bool{"id1": true, "id2": false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewModelProviderStore()
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

func TestModelProviderStore_Get(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})

	result := store.Get("id1")
	if result == nil {
		t.Fatal("Get returned nil for existing ID")
	}
	if result.ID != "id1" {
		t.Errorf("expected ID id1, got %s", result.ID)
	}

	result = store.Get("non-existent")
	if result != nil {
		t.Error("Get should return nil for non-existent ID")
	}
}

func TestModelProviderStore_Update(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{
		ID:      "id1",
		Name:    "Old Name",
		BaseURL: "https://old.com",
		APIKey:  "old-key",
		Model:   "old-model",
	})

	store.Update(
		"id1",
		"New Name",
		"https://new.com",
		"new-key",
		"new-model",
		"new-thinking",
		"new-haiku",
		"new-sonnet",
		"new-opus",
			"new-subagent",
			"1",
			"1",
			"max",
	)

	p := store.Get("id1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Name != "New Name" {
		t.Errorf("expected Name 'New Name', got '%s'", p.Name)
	}
	if p.BaseURL != "https://new.com" {
		t.Errorf("expected BaseURL 'https://new.com', got '%s'", p.BaseURL)
	}
	if p.APIKey != "new-key" {
		t.Errorf("expected APIKey 'new-key', got '%s'", p.APIKey)
	}
	if p.Model != "new-model" {
		t.Errorf("expected Model 'new-model', got '%s'", p.Model)
	}
	if p.ThinkingModel != "new-thinking" {
		t.Errorf("expected ThinkingModel 'new-thinking', got '%s'", p.ThinkingModel)
	}
	if p.DefaultHaikuModel != "new-haiku" {
		t.Errorf("expected DefaultHaikuModel 'new-haiku', got '%s'", p.DefaultHaikuModel)
	}
	if p.DefaultSonnetModel != "new-sonnet" {
		t.Errorf("expected DefaultSonnetModel 'new-sonnet', got '%s'", p.DefaultSonnetModel)
	}
	if p.DefaultOpusModel != "new-opus" {
		t.Errorf("expected DefaultOpusModel 'new-opus', got '%s'", p.DefaultOpusModel)
	}
}

func TestModelProviderStore_UpdateNonExistent(t *testing.T) {
	store := NewModelProviderStore()
	// 更新不存在的配置不应 panic
	store.Update("non-existent", "name", "url", "key", "model", "thinking", "haiku", "sonnet", "opus", "", "", "", "")
}

func TestModelProviderStore_SetActive(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(*ModelProviderStore)
		setActiveID   string
		wantActiveID  string
		wantProviders map[string]bool // id -> 期望的 Active 状态
	}{
		{
			name: "设置激活配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
				s.Add(ModelProvider{ID: "id3", Name: "P3", BaseURL: "https://p3.com", APIKey: "key", Model: "m3"})
			},
			setActiveID:  "id2",
			wantActiveID: "id2",
			wantProviders: map[string]bool{
				"id1": false,
				"id2": true,
				"id3": false,
			},
		},
		{
			name: "切换激活配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
				s.SetActive("id2")
			},
			setActiveID:  "id1",
			wantActiveID: "id1",
			wantProviders: map[string]bool{
				"id1": true,
				"id2": false,
			},
		},
		{
			name: "设置不存在的配置为激活",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
			},
			setActiveID:  "non-existent",
			wantActiveID: "",
			wantProviders: map[string]bool{
				"id1": false,
			},
		},
		{
			name: "空存储设置激活",
			setup: func(s *ModelProviderStore) {},
			setActiveID:  "id1",
			wantActiveID: "",
			wantProviders: map[string]bool{},
		},
		{
			name: "SetActive 清除其他激活状态",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1", Active: true})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2", Active: true})
			},
			setActiveID:  "id2",
			wantActiveID: "id2",
			wantProviders: map[string]bool{
				"id1": false,
				"id2": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewModelProviderStore()
			tt.setup(store)
			store.SetActive(tt.setActiveID)

			if store.ActiveID != tt.wantActiveID {
				t.Errorf("ActiveID = %q, want %q", store.ActiveID, tt.wantActiveID)
			}

			for id, wantActive := range tt.wantProviders {
				p := store.Get(id)
				if p == nil {
					t.Errorf("Get(%q) = nil", id)
					continue
				}
				if p.Active != wantActive {
					t.Errorf("Get(%q).Active = %v, want %v", id, p.Active, wantActive)
				}
			}

			// 验证只有一个激活的配置
			activeCount := 0
			for _, p := range store.Providers {
				if p.Active {
					activeCount++
				}
			}
			if activeCount > 1 {
				t.Errorf("expected at most 1 active provider, got %d", activeCount)
			}
		})
	}
}

func TestModelProviderStore_GetActive(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*ModelProviderStore)
		wantActive string // 期望激活的 ID，为空表示期望 nil
	}{
		{
			name: "空存储",
			setup: func(s *ModelProviderStore) {},
			wantActive: "",
		},
		{
			name: "无激活配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
			},
			wantActive: "",
		},
		{
			name: "有 Active=true 的配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2", Active: true})
			},
			wantActive: "id2",
		},
		{
			name: "通过 SetActive 设置激活",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
				s.SetActive("id2")
			},
			wantActive: "id2",
		},
		{
			name: "切换激活配置",
			setup: func(s *ModelProviderStore) {
				s.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
				s.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
				s.SetActive("id2")
				s.SetActive("id1")
			},
			wantActive: "id1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewModelProviderStore()
			tt.setup(store)
			active := store.GetActive()

			if tt.wantActive == "" {
				if active != nil {
					t.Errorf("GetActive() = %v, want nil", active.ID)
				}
			} else {
				if active == nil {
					t.Fatalf("GetActive() = nil, want %s", tt.wantActive)
				}
				if active.ID != tt.wantActive {
					t.Errorf("GetActive().ID = %s, want %s", active.ID, tt.wantActive)
				}
			}
		})
	}
}

func TestModelProviderStore_Len(t *testing.T) {
	store := NewModelProviderStore()
	if store.Len() != 0 {
		t.Errorf("expected 0, got %d", store.Len())
	}

	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
	if store.Len() != 1 {
		t.Errorf("expected 1, got %d", store.Len())
	}

	store.Remove("id1")
	if store.Len() != 0 {
		t.Errorf("expected 0 after remove, got %d", store.Len())
	}
}

func TestGenerateProviderID(t *testing.T) {
	id1 := GenerateProviderID()
	id2 := GenerateProviderID()

	if id1 == "" {
		t.Error("GenerateProviderID should not return empty string")
	}

	if id1 == id2 {
		t.Error("GenerateProviderID should generate unique IDs")
	}

	// ID 应该是 8 字符的十六进制字符串
	if len(id1) != 8 {
		t.Errorf("expected ID length 8, got %d", len(id1))
	}
}

func TestModelProviderStore_Load_Save(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "providers.json")

	// 创建并保存存储
	store1 := NewModelProviderStore()
	store1.Add(ModelProvider{
		ID:              "id1",
		Name:            "Test Provider",
		BaseURL:         "https://api.test.com",
		APIKey:          "secret-key",
		Model:           "test-model",
		ThinkingModel:   "test-reasoning",
		DefaultHaikuModel:  "haiku-3",
		DefaultSonnetModel: "sonnet-4",
		DefaultOpusModel:   "opus-3",
		Active:          true,
		CreatedAt:       time.Now(),
	})
	store1.SetActive("id1")

	if err := store1.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// 加载到新存储
	store2 := NewModelProviderStore()
	if err := store2.Load(path); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if store2.Len() != 1 {
		t.Errorf("expected loaded store to have 1 provider, got %d", store2.Len())
	}

	p := store2.Get("id1")
	if p == nil {
		t.Fatal("loaded provider not found")
	}
	if p.Name != "Test Provider" {
		t.Errorf("expected Name 'Test Provider', got '%s'", p.Name)
	}
	if p.BaseURL != "https://api.test.com" {
		t.Errorf("expected BaseURL 'https://api.test.com', got '%s'", p.BaseURL)
	}
	if p.APIKey != "secret-key" {
		t.Errorf("expected APIKey 'secret-key', got '%s'", p.APIKey)
	}
	if p.Model != "test-model" {
		t.Errorf("expected Model 'test-model', got '%s'", p.Model)
	}
	if p.ThinkingModel != "test-reasoning" {
		t.Errorf("expected ThinkingModel 'test-reasoning', got '%s'", p.ThinkingModel)
	}
}

func TestModelProviderStore_LoadNonExistent(t *testing.T) {
	store := NewModelProviderStore()
	// 加载不存在的文件应该返回 nil
	if err := store.Load("/tmp/non-existent-provider-12345.json"); err != nil {
		t.Errorf("Load should not return error for non-existent file: %v", err)
	}
	if store.Len() != 0 {
		t.Errorf("store should remain empty, got %d", store.Len())
	}
}

func TestModelProviderStore_SaveInvalidDir(t *testing.T) {
	store := NewModelProviderStore()
	err := store.Save("/root/impossible-dir/providers.json")
	if err == nil {
		t.Error("Save should return error for invalid directory")
	}
}

func TestModelProvider_JSONSerialization(t *testing.T) {
	now := time.Now()
	p := ModelProvider{
		ID:                "test-id",
		Name:              "Test Provider",
		BaseURL:           "https://api.test.com",
		APIKey:            "secret",
		Model:             "test-model",
		ThinkingModel:      "reasoning",
		DefaultHaikuModel:  "haiku",
		DefaultSonnetModel: "sonnet",
		DefaultOpusModel:   "opus",
		Active:            true,
		CreatedAt:         now,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var p2 ModelProvider
	if err := json.Unmarshal(data, &p2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if p2.ID != p.ID {
		t.Errorf("ID mismatch: expected %s, got %s", p.ID, p2.ID)
	}
	if p2.Name != p.Name {
		t.Errorf("Name mismatch: expected %s, got %s", p.Name, p2.Name)
	}
	if p2.BaseURL != p.BaseURL {
		t.Errorf("BaseURL mismatch: expected %s, got %s", p.BaseURL, p2.BaseURL)
	}
	if p2.Active != p.Active {
		t.Errorf("Active mismatch: expected %v, got %v", p.Active, p2.Active)
	}
}

func TestModelProviderStore_JSONSerialization(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
	store.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
	store.ActiveID = "id2"

	data, err := json.Marshal(store)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var store2 ModelProviderStore
	if err := json.Unmarshal(data, &store2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if store2.Len() != 2 {
		t.Errorf("expected 2 providers, got %d", store2.Len())
	}
	if store2.ActiveID != "id2" {
		t.Errorf("expected ActiveID 'id2', got '%s'", store2.ActiveID)
	}
}

func TestModelProviderStore_IndexAfterRemove(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
	store.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})
	store.Add(ModelProvider{ID: "id3", Name: "P3", BaseURL: "https://p3.com", APIKey: "key", Model: "m3"})

	store.Remove("id2")

	// 添加新项目验证索引正确
	store.Add(ModelProvider{ID: "id4", Name: "P4", BaseURL: "https://p4.com", APIKey: "key", Model: "m4"})

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

func TestModelProviderStore_RebuildIndex(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "P1", BaseURL: "https://p1.com", APIKey: "key", Model: "m1"})
	store.Add(ModelProvider{ID: "id2", Name: "P2", BaseURL: "https://p2.com", APIKey: "key", Model: "m2"})

	// 手动触发索引重建
	store.rebuildIndex()

	if store.Get("id1") == nil {
		t.Error("id1 should be accessible after rebuildIndex")
	}
	if store.Get("id2") == nil {
		t.Error("id2 should be accessible after rebuildIndex")
	}
}

func TestModelProviderStore_GetAfterUpdate(t *testing.T) {
	store := NewModelProviderStore()
	store.Add(ModelProvider{ID: "id1", Name: "Original", BaseURL: "https://orig.com", APIKey: "key", Model: "m1"})

	store.Update("id1", "Updated", "https://updated.com", "new-key", "new-model", "", "", "", "", "", "", "", "")

	p := store.Get("id1")
	if p == nil {
		t.Fatal("Get returned nil")
	}
	if p.Name != "Updated" {
		t.Errorf("expected Name 'Updated', got '%s'", p.Name)
	}
	if p.BaseURL != "https://updated.com" {
		t.Errorf("expected BaseURL 'https://updated.com', got '%s'", p.BaseURL)
	}
}

func TestModelProvider_Fields(t *testing.T) {
	now := time.Now()
	p := ModelProvider{
		ID:                "id",
		Name:              "Name",
		BaseURL:           "BaseURL",
		APIKey:            "APIKey",
		Model:             "Model",
		ThinkingModel:      "ThinkingModel",
		DefaultHaikuModel:  "DefaultHaikuModel",
		DefaultSonnetModel: "DefaultSonnetModel",
		DefaultOpusModel:   "DefaultOpusModel",
		Active:            true,
		CreatedAt:         now,
	}

	if p.ID != "id" {
		t.Error("ID mismatch")
	}
	if p.Name != "Name" {
		t.Error("Name mismatch")
	}
	if p.BaseURL != "BaseURL" {
		t.Error("BaseURL mismatch")
	}
	if p.APIKey != "APIKey" {
		t.Error("APIKey mismatch")
	}
	if p.Model != "Model" {
		t.Error("Model mismatch")
	}
	if p.ThinkingModel != "ThinkingModel" {
		t.Error("ThinkingModel mismatch")
	}
	if p.DefaultHaikuModel != "DefaultHaikuModel" {
		t.Error("DefaultHaikuModel mismatch")
	}
	if p.DefaultSonnetModel != "DefaultSonnetModel" {
		t.Error("DefaultSonnetModel mismatch")
	}
	if p.DefaultOpusModel != "DefaultOpusModel" {
		t.Error("DefaultOpusModel mismatch")
	}
	if !p.Active {
		t.Error("Active should be true")
	}
	if !p.CreatedAt.Equal(now) {
		t.Error("CreatedAt mismatch")
	}
}

func TestIDETypeConstants(t *testing.T) {
	if IDEClaudeCode != "claude" {
		t.Errorf("expected IDEClaudeCode 'claude', got '%s'", IDEClaudeCode)
	}
	if IDEVSCode != "code" {
		t.Errorf("expected IDEVSCode 'code', got '%s'", IDEVSCode)
	}
	if IDEOpenCode != "opencode" {
		t.Errorf("expected IDEOpenCode 'opencode', got '%s'", IDEOpenCode)
	}
}
