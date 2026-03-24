package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"gocoding/internal/models"
)

const defaultProviderFileName = "model_providers.json"

type ModelProviderConfig struct {
	filePath string
	store    *models.ModelProviderStore
}

func NewModelProviderConfig(store *models.ModelProviderStore) *ModelProviderConfig {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "gocoding")
	return &ModelProviderConfig{
		filePath: filepath.Join(configDir, defaultProviderFileName),
		store:    store,
	}
}

func (c *ModelProviderConfig) SetFilePath(path string) {
	c.filePath = path
}

func (c *ModelProviderConfig) Load() error {
	return c.store.Load(c.filePath)
}

func (c *ModelProviderConfig) Save() error {
	return c.store.Save(c.filePath)
}

func (c *ModelProviderConfig) GetStore() *models.ModelProviderStore {
	return c.store
}

// ClaudeSettingsPath 返回 Claude Code settings.json 路径
func ClaudeSettingsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".claude", "settings.json"), nil
}

// IsProviderConfigMatch 检查当前 settings.json 是否与配置一致
func IsProviderConfigMatch(provider *models.ModelProvider, settings map[string]interface{}) bool {
	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		return false
	}

	// 检查 ANTHROPIC_BASE_URL
	if baseURL, ok := env["ANTHROPIC_BASE_URL"].(string); !ok || baseURL != provider.BaseURL {
		return false
	}

	// 检查 ANTHROPIC_AUTH_TOKEN
	if apiKey, ok := env["ANTHROPIC_AUTH_TOKEN"].(string); !ok || apiKey != provider.APIKey {
		return false
	}

	// 检查 model 字段
	if model, ok := settings["model"].(string); !ok || model != provider.Model {
		return false
	}

	return true
}

// WriteToClaudeSettings 将指定配置写入 Claude settings.json
// 返回 (isAlreadySet bool, err error)
// isAlreadySet=true 表示配置已生效，无需写入
func WriteToClaudeSettings(provider *models.ModelProvider) (bool, error) {
	settingsPath, err := ClaudeSettingsPath()
	if err != nil {
		return false, err
	}

	// 读取现有配置
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return false, err
	}

	// 解析 JSON
	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false, err
	}

	// 检查是否与当前配置一致
	if IsProviderConfigMatch(provider, settings) {
		return true, nil // 配置已一致，无需写入
	}

	// 更新 env 字段
	env := make(map[string]interface{})
	if existingEnv, ok := settings["env"].(map[string]interface{}); ok {
		for k, v := range existingEnv {
			env[k] = v
		}
	}
	env["ANTHROPIC_BASE_URL"] = provider.BaseURL
	env["ANTHROPIC_AUTH_TOKEN"] = provider.APIKey

	// 更新所有 ANTHROPIC_*MODEL 结尾的键
	for k := range env {
		if strings.HasPrefix(k, "ANTHROPIC_") && strings.HasSuffix(k, "MODEL") {
			env[k] = provider.Model
		}
	}
	// 也更新 ANTHROPIC_MODEL
	env["ANTHROPIC_MODEL"] = provider.Model

	settings["env"] = env

	// 更新 model 字段
	settings["model"] = provider.Model

	// 写回文件
	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return false, err
	}

	return false, os.WriteFile(settingsPath, output, 0644)
}
