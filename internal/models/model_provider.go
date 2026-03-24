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
	ID              string    `json:"id"`
	Name            string    `json:"name"`       // 配置名称，如 "MiniMax", "Claude"
	BaseURL         string    `json:"base_url"`   // API Base URL
	APIKey          string    `json:"api_key"`    // API Key
	Model           string    `json:"model"`      // 默认模型名
	ThinkingModel   string    `json:"thinking_model"`    // 推理模型
	DefaultHaikuModel  string `json:"default_haiku_model"`  // Haiku 默认模型
	DefaultSonnetModel string `json:"default_sonnet_model"` // Sonnet 默认模型
	DefaultOpusModel   string `json:"default_opus_model"`  // Opus 默认模型
	Active          bool      `json:"active"`     // 是否激活
	CreatedAt       time.Time `json:"created_at"`
}

type ModelProviderStore struct {
	Providers []ModelProvider  `json:"providers"`
	ActiveID  string           `json:"active_id"` // 当前激活的配置ID
	index     map[string]int    `json:"-"`         // id -> index 映射，不持久化
}

func NewModelProviderStore() *ModelProviderStore {
	return &ModelProviderStore{
		Providers: make([]ModelProvider, 0),
		index:     make(map[string]int),
	}
}

// rebuildIndex 重建索引
func (s *ModelProviderStore) rebuildIndex() {
	s.index = make(map[string]int)
	for i, p := range s.Providers {
		s.index[p.ID] = i
	}
}

// Add 添加新配置
func (s *ModelProviderStore) Add(provider ModelProvider) {
	s.Providers = append(s.Providers, provider)
	s.index[provider.ID] = len(s.Providers) - 1
}

// Remove 删除配置
func (s *ModelProviderStore) Remove(id string) {
	idx, ok := s.index[id]
	if !ok {
		return
	}
	s.Providers = append(s.Providers[:idx], s.Providers[idx+1:]...)
	s.rebuildIndex()
}

// Get 获取配置
func (s *ModelProviderStore) Get(id string) *ModelProvider {
	idx, ok := s.index[id]
	if !ok {
		return nil
	}
	return &s.Providers[idx]
}

// Update 更新配置
func (s *ModelProviderStore) Update(id, name, baseURL, apiKey, model, thinkingModel, defaultHaikuModel, defaultSonnetModel, defaultOpusModel string) {
	idx, ok := s.index[id]
	if !ok {
		return
	}
	s.Providers[idx].Name = name
	s.Providers[idx].BaseURL = baseURL
	s.Providers[idx].APIKey = apiKey
	s.Providers[idx].Model = model
	s.Providers[idx].ThinkingModel = thinkingModel
	s.Providers[idx].DefaultHaikuModel = defaultHaikuModel
	s.Providers[idx].DefaultSonnetModel = defaultSonnetModel
	s.Providers[idx].DefaultOpusModel = defaultOpusModel
}

// SetActive 设置激活配置
func (s *ModelProviderStore) SetActive(id string) {
	// 先取消所有激活
	for i := range s.Providers {
		s.Providers[i].Active = false
	}
	// 激活指定配置
	idx, ok := s.index[id]
	if !ok {
		return
	}
	s.Providers[idx].Active = true
	s.ActiveID = id
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
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	s.rebuildIndex()
	return nil
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
