package models

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ModelProvider struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`      // 配置名称，如 "MiniMax", "Claude"
	BaseURL  string    `json:"base_url"`  // API Base URL
	APIKey   string    `json:"api_key"`   // API Key
	Model    string    `json:"model"`     // 默认模型名
	Active   bool      `json:"active"`    // 是否激活
	CreatedAt time.Time `json:"created_at"`
}

type ModelProviderStore struct {
	Providers []ModelProvider `json:"providers"`
	ActiveID  string         `json:"active_id"` // 当前激活的配置ID
}

func NewModelProviderStore() *ModelProviderStore {
	return &ModelProviderStore{
		Providers: make([]ModelProvider, 0),
	}
}

// Add 添加新配置
func (s *ModelProviderStore) Add(provider ModelProvider) {
	s.Providers = append(s.Providers, provider)
}

// Remove 删除配置
func (s *ModelProviderStore) Remove(id string) {
	for i, p := range s.Providers {
		if p.ID == id {
			s.Providers = append(s.Providers[:i], s.Providers[i+1:]...)
			return
		}
	}
}

// Get 获取配置
func (s *ModelProviderStore) Get(id string) *ModelProvider {
	for i := range s.Providers {
		if s.Providers[i].ID == id {
			return &s.Providers[i]
		}
	}
	return nil
}

// Update 更新配置
func (s *ModelProviderStore) Update(id string, name, baseURL, apiKey, model string) {
	for i := range s.Providers {
		if s.Providers[i].ID == id {
			s.Providers[i].Name = name
			s.Providers[i].BaseURL = baseURL
			s.Providers[i].APIKey = apiKey
			s.Providers[i].Model = model
			return
		}
	}
}

// SetActive 设置激活配置
func (s *ModelProviderStore) SetActive(id string) {
	// 先取消所有激活
	for i := range s.Providers {
		s.Providers[i].Active = false
	}
	// 激活指定配置
	for i := range s.Providers {
		if s.Providers[i].ID == id {
			s.Providers[i].Active = true
			s.ActiveID = id
			return
		}
	}
}

// GetActive 获取当前激活的配置
func (s *ModelProviderStore) GetActive() *ModelProvider {
	for i := range s.Providers {
		if s.Providers[i].Active {
			return &s.Providers[i]
		}
	}
	// 如果没有激活的，尝试用 ActiveID 找
	if s.ActiveID != "" {
		return s.Get(s.ActiveID)
	}
	return nil
}

// Len 返回配置数量
func (s *ModelProviderStore) Len() int {
	return len(s.Providers)
}

// GenerateID 生成唯一ID
func GenerateProviderID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *ModelProviderStore) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, s)
}

func (s *ModelProviderStore) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
