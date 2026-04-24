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

	// 检查主模型（允许为空）
	if provider.Model != "" {
		if model, ok := settings["model"].(string); !ok || model != provider.Model {
			return false
		}
	}

	// 检查推理模型（允许为空）
	if provider.ThinkingModel != "" {
		if reasoningModel, ok := env["ANTHROPIC_REASONING_MODEL"].(string); !ok || reasoningModel != provider.ThinkingModel {
			return false
		}
	} else {
		// 为空时，settings 中也不应存在该字段
		if _, exists := env["ANTHROPIC_REASONING_MODEL"]; exists {
			return false
		}
	}

	// 检查 Haiku 默认模型（允许为空）
	if provider.DefaultHaikuModel != "" {
		if haikuModel, ok := env["ANTHROPIC_DEFAULT_HAIKU_MODEL"].(string); !ok || haikuModel != provider.DefaultHaikuModel {
			return false
		}
	} else {
		if _, exists := env["ANTHROPIC_DEFAULT_HAIKU_MODEL"]; exists {
			return false
		}
	}

	// 检查 Sonnet 默认模型（允许为空）
	if provider.DefaultSonnetModel != "" {
		if sonnetModel, ok := env["ANTHROPIC_DEFAULT_SONNET_MODEL"].(string); !ok || sonnetModel != provider.DefaultSonnetModel {
			return false
		}
	} else {
		if _, exists := env["ANTHROPIC_DEFAULT_SONNET_MODEL"]; exists {
			return false
		}
	}

	// 检查 Opus 默认模型（允许为空）
	if provider.DefaultOpusModel != "" {
		if opusModel, ok := env["ANTHROPIC_DEFAULT_OPUS_MODEL"].(string); !ok || opusModel != provider.DefaultOpusModel {
			return false
		}
	} else {
		if _, exists := env["ANTHROPIC_DEFAULT_OPUS_MODEL"]; exists {
			return false
		}
	}

	// 检查 SubAgent 模型（允许为空）
	if provider.SubagentModel != "" {
		if subagentModel, ok := env["CLAUDE_CODE_SUBAGENT_MODEL"].(string); !ok || subagentModel != provider.SubagentModel {
			return false
		}
	} else {
		if _, exists := env["CLAUDE_CODE_SUBAGENT_MODEL"]; exists {
			return false
		}
	}

	// 检查 DisableNonessential（允许为空）
	if provider.DisableNonessential != "" {
		if v, ok := env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"].(string); !ok || v != provider.DisableNonessential {
			return false
		}
	} else {
		if _, exists := env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"]; exists {
			return false
		}
	}

	// 检查 DisableNonstreaming（允许为空）
	if provider.DisableNonstreaming != "" {
		if v, ok := env["CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK"].(string); !ok || v != provider.DisableNonstreaming {
			return false
		}
	} else {
		if _, exists := env["CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK"]; exists {
			return false
		}
	}

	// 检查 EffortLevel（允许为空）
	if provider.EffortLevel != "" {
		if v, ok := env["CLAUDE_CODE_EFFORT_LEVEL"].(string); !ok || v != provider.EffortLevel {
			return false
		}
	} else {
		if _, exists := env["CLAUDE_CODE_EFFORT_LEVEL"]; exists {
			return false
		}
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

	// 主模型 - 对应 ANTHROPIC_MODEL
	if provider.Model != "" {
		env["ANTHROPIC_MODEL"] = provider.Model
	} else {
		delete(env, "ANTHROPIC_MODEL")
	}

	// 推理模型 - 为空时移除
	if provider.ThinkingModel != "" {
		env["ANTHROPIC_REASONING_MODEL"] = provider.ThinkingModel
	} else {
		delete(env, "ANTHROPIC_REASONING_MODEL")
	}

	// Haiku 默认模型 - 为空时移除
	if provider.DefaultHaikuModel != "" {
		env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = provider.DefaultHaikuModel
	} else {
		delete(env, "ANTHROPIC_DEFAULT_HAIKU_MODEL")
	}

	// Sonnet 默认模型 - 为空时移除
	if provider.DefaultSonnetModel != "" {
		env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = provider.DefaultSonnetModel
	} else {
		delete(env, "ANTHROPIC_DEFAULT_SONNET_MODEL")
	}

	// Opus 默认模型 - 为空时移除
	if provider.DefaultOpusModel != "" {
		env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = provider.DefaultOpusModel
	} else {
		delete(env, "ANTHROPIC_DEFAULT_OPUS_MODEL")
	}

	// SubAgent 模型 - 为空时移除
	if provider.SubagentModel != "" {
		env["CLAUDE_CODE_SUBAGENT_MODEL"] = provider.SubagentModel
	} else {
		delete(env, "CLAUDE_CODE_SUBAGENT_MODEL")
	}

	// 禁用非必要流量 - 为空时移除
	if provider.DisableNonessential != "" {
		env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = provider.DisableNonessential
	} else {
		delete(env, "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC")
	}

	// 禁用非流式回退 - 为空时移除
	if provider.DisableNonstreaming != "" {
		env["CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK"] = provider.DisableNonstreaming
	} else {
		delete(env, "CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK")
	}

	// 推理力度 - 为空时移除
	if provider.EffortLevel != "" {
		env["CLAUDE_CODE_EFFORT_LEVEL"] = provider.EffortLevel
	} else {
		delete(env, "CLAUDE_CODE_EFFORT_LEVEL")
	}

	// 清理空值键
	for k := range env {
		if strings.HasPrefix(k, "ANTHROPIC_") && strings.HasSuffix(k, "_MODEL") && env[k] == "" {
			delete(env, k)
		}
		if k == "CLAUDE_CODE_SUBAGENT_MODEL" && env[k] == "" {
			delete(env, k)
		}
		if k == "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC" && env[k] == "" {
			delete(env, k)
		}
		if k == "CLAUDE_CODE_DISABLE_NONSTREAMING_FALLBACK" && env[k] == "" {
			delete(env, k)
		}
		if k == "CLAUDE_CODE_EFFORT_LEVEL" && env[k] == "" {
			delete(env, k)
		}
	}

	settings["env"] = env

	// 更新 model 字段（主模型为空时使用默认）
	if provider.Model != "" {
		settings["model"] = provider.Model
	}

	// 写回文件
	output, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return false, err
	}

	return false, os.WriteFile(settingsPath, output, 0644)
}
