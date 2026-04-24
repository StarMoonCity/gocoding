package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"gocoding/internal/config"
	"gocoding/internal/models"
)

// handleProviderListKeyMsg 处理配置列表按键
func (m *Model) handleProviderListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down", "ctrl+n":
		m.providerList.CursorDown()
	case "k", "up", "ctrl+p":
		m.providerList.CursorUp()
	case "n":
		// 新增配置
		m.state = StateProviderAdd
		m.providerInputFocus = FocusProviderName
		m.providerNameInput.Reset()
		m.providerNameInput.Placeholder = "Name (e.g. MiniMax)"
		m.providerNameInput.SetValue("")
		m.providerNameInput.Focus()
		m.providerBaseURLInput.Reset()
		m.providerBaseURLInput.Placeholder = "Base URL (e.g. https://api.minimax.chat)"
		m.providerBaseURLInput.SetValue("")
		m.providerBaseURLInput.Blur()
		m.providerAPIKeyInput.Reset()
		m.providerAPIKeyInput.Placeholder = "API Key"
		m.providerAPIKeyInput.SetValue("")
		m.providerAPIKeyInput.Blur()
		m.providerModelInput.Reset()
		m.providerModelInput.Placeholder = "Model (e.g. MiniMax-M2.7-highspeed)"
		m.providerModelInput.SetValue("")
		m.providerModelInput.Blur()
		m.providerThinkingModelInput.Reset()
		m.providerThinkingModelInput.Placeholder = "Reasoning Model"
		m.providerThinkingModelInput.SetValue("")
		m.providerThinkingModelInput.Blur()
		m.providerDefaultHaikuInput.Reset()
		m.providerDefaultHaikuInput.Placeholder = "Default Haiku Model"
		m.providerDefaultHaikuInput.SetValue("")
		m.providerDefaultHaikuInput.Blur()
		m.providerDefaultSonnetInput.Reset()
		m.providerDefaultSonnetInput.Placeholder = "Default Sonnet Model"
		m.providerDefaultSonnetInput.SetValue("")
		m.providerDefaultSonnetInput.Blur()
		m.providerDefaultOpusInput.Reset()
		m.providerDefaultOpusInput.Placeholder = "Default Opus Model"
		m.providerDefaultOpusInput.SetValue("")
		m.providerDefaultOpusInput.Blur()
		m.providerSubagentInput.Reset()
		m.providerSubagentInput.Placeholder = "SubAgent Model"
		m.providerSubagentInput.SetValue("")
		m.providerSubagentInput.Blur()
		m.providerNonessentialInput.Reset()
		m.providerNonessentialInput.Placeholder = "1=禁用非必要流量"
		m.providerNonessentialInput.SetValue("")
		m.providerNonessentialInput.Blur()
		m.providerNonstreamingInput.Reset()
		m.providerNonstreamingInput.Placeholder = "1=禁用非流式回退"
		m.providerNonstreamingInput.SetValue("")
		m.providerNonstreamingInput.Blur()
		m.providerEffortInput.Reset()
		m.providerEffortInput.Placeholder = "max/high/medium/low"
		m.providerEffortInput.SetValue("")
		m.providerEffortInput.Blur()
		m.editingProviderID = ""
		m.errMsg = ""
		return m, textinput.Blink
	case "e":
		// 编辑选中配置
		if current := m.safeGetSelectedProvider(); current != nil {
			m.state = StateProviderEdit
			m.editingProviderID = current.ID
			m.providerInputFocus = FocusProviderName
			m.providerNameInput.SetValue(current.Name)
			m.providerNameInput.Focus()
			m.providerBaseURLInput.SetValue(current.BaseURL)
			m.providerBaseURLInput.Blur()
			m.providerAPIKeyInput.SetValue(current.APIKey)
			m.providerAPIKeyInput.Blur()
			m.providerModelInput.SetValue(current.Model)
			m.providerModelInput.Blur()
			m.providerThinkingModelInput.SetValue(current.ThinkingModel)
			m.providerThinkingModelInput.Blur()
			m.providerDefaultHaikuInput.SetValue(current.DefaultHaikuModel)
			m.providerDefaultHaikuInput.Blur()
			m.providerDefaultSonnetInput.SetValue(current.DefaultSonnetModel)
			m.providerDefaultSonnetInput.Blur()
			m.providerDefaultOpusInput.SetValue(current.DefaultOpusModel)
			m.providerDefaultOpusInput.Blur()
			m.providerSubagentInput.SetValue(current.SubagentModel)
			m.providerSubagentInput.Blur()
			m.providerNonessentialInput.SetValue(current.DisableNonessential)
			m.providerNonessentialInput.Blur()
			m.providerNonstreamingInput.SetValue(current.DisableNonstreaming)
			m.providerNonstreamingInput.Blur()
			m.providerEffortInput.SetValue(current.EffortLevel)
			m.providerEffortInput.Blur()
			m.errMsg = ""
			return m, textinput.Blink
		}
	case "d":
		// 删除确认
		if m.safeGetSelectedProvider() != nil {
			m.state = StateProviderDelete
		}
	case "a":
		// 激活选中配置
		if current := m.safeGetSelectedProvider(); current != nil {
			m.providerStore.SetActive(current.ID)
			// 写入 Claude settings.json
			provider := m.providerStore.Get(current.ID)
			alreadySet, err := config.WriteToClaudeSettings(provider)
			m.updateProviderListItems()
			if err != nil {
				m.errMsg = "激活失败: " + err.Error()
				// 5秒后清除错误消息
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return clearProviderMessageMsg{msgType: msgTypeErr}
				})
			} else if alreadySet {
				m.tipMsg = current.Name + " 配置已生效，无需更新"
				// 5秒后清除提示消息
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return clearProviderMessageMsg{msgType: msgTypeTip}
				})
			} else {
				m.tipMsg = "已激活 " + current.Name
				// 5秒后清除提示消息
				return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
					return clearProviderMessageMsg{msgType: msgTypeTip}
				})
			}
		}
	case "enter":
		// 编辑选中配置
		if current := m.safeGetSelectedProvider(); current != nil {
			m.state = StateProviderEdit
			m.editingProviderID = current.ID
			m.providerInputFocus = FocusProviderName
			m.providerNameInput.SetValue(current.Name)
			m.providerNameInput.Focus()
			m.providerBaseURLInput.SetValue(current.BaseURL)
			m.providerBaseURLInput.Blur()
			m.providerAPIKeyInput.SetValue(current.APIKey)
			m.providerAPIKeyInput.Blur()
			m.providerModelInput.SetValue(current.Model)
			m.providerModelInput.Blur()
			m.providerThinkingModelInput.SetValue(current.ThinkingModel)
			m.providerThinkingModelInput.Blur()
			m.providerDefaultHaikuInput.SetValue(current.DefaultHaikuModel)
			m.providerDefaultHaikuInput.Blur()
			m.providerDefaultSonnetInput.SetValue(current.DefaultSonnetModel)
			m.providerDefaultSonnetInput.Blur()
			m.providerDefaultOpusInput.SetValue(current.DefaultOpusModel)
			m.providerDefaultOpusInput.Blur()
			m.providerSubagentInput.SetValue(current.SubagentModel)
			m.providerSubagentInput.Blur()
			m.providerNonessentialInput.SetValue(current.DisableNonessential)
			m.providerNonessentialInput.Blur()
			m.providerNonstreamingInput.SetValue(current.DisableNonstreaming)
			m.providerNonstreamingInput.Blur()
			m.providerEffortInput.SetValue(current.EffortLevel)
			m.providerEffortInput.Blur()
			m.errMsg = ""
			return m, textinput.Blink
		}
	case "esc":
		m.state = StateList
	case "ctrl+c", "ctrl+q", "q":
		return m, tea.Quit
	}
	return m, nil
}

// handleProviderFormKeyMsg 处理配置表单按键（新增/编辑共用）
func (m *Model) handleProviderFormKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	isEdit := m.editingProviderID != ""

	switch msg.String() {
	case "enter":
		name := m.providerNameInput.Value()
		baseURL := m.providerBaseURLInput.Value()
		apiKey := m.providerAPIKeyInput.Value()
		model := m.providerModelInput.Value()
		thinkingModel := m.providerThinkingModelInput.Value()
		defaultHaikuModel := m.providerDefaultHaikuInput.Value()
		defaultSonnetModel := m.providerDefaultSonnetInput.Value()
		defaultOpusModel := m.providerDefaultOpusInput.Value()
		subagentModel := m.providerSubagentInput.Value()
		disableNonessential := m.providerNonessentialInput.Value()
		disableNonstreaming := m.providerNonstreamingInput.Value()
		effortLevel := m.providerEffortInput.Value()

		// 验证
		if name == "" {
			return m.setProviderErrMsg("配置名称不能为空")
		}
		if baseURL == "" {
			return m.setProviderErrMsg("Base URL不能为空")
		}
		if apiKey == "" {
			return m.setProviderErrMsg("API Key不能为空")
		}
		// 主模型允许为空

		if isEdit {
			m.providerStore.Update(m.editingProviderID, name, baseURL, apiKey, model, thinkingModel, defaultHaikuModel, defaultSonnetModel, defaultOpusModel, subagentModel, disableNonessential, disableNonstreaming, effortLevel)
		} else {
			provider := models.ModelProvider{
				ID:                    models.GenerateProviderID(),
				Name:                  name,
				BaseURL:               baseURL,
				APIKey:                apiKey,
				Model:                 model,
				ThinkingModel:         thinkingModel,
				DefaultHaikuModel:     defaultHaikuModel,
				DefaultSonnetModel:    defaultSonnetModel,
				DefaultOpusModel:      defaultOpusModel,
				SubagentModel:         subagentModel,
				DisableNonessential:   disableNonessential,
				DisableNonstreaming:   disableNonstreaming,
				EffortLevel:           effortLevel,
				CreatedAt:             time.Now(),
			}
			m.providerStore.Add(provider)
		}
		m.updateProviderListItems()
		m.state = StateProviderList
		m.errMsg = ""
		return m, nil
	case "tab":
		// 切换焦点
		m.providerInputFocus = (m.providerInputFocus + 1) % FocusProviderCount
		m.updateProviderFocus()
	case "esc":
		m.state = StateProviderList
		m.errMsg = ""
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}

	// 更新当前焦点输入框
	return m, m.updateProviderInput(msg)
}

// handleProviderDeleteKeyMsg 处理删除确认按键
func (m *Model) handleProviderDeleteKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		if current := m.safeGetSelectedProvider(); current != nil {
			m.providerStore.Remove(current.ID)
			m.updateProviderListItems()
		}
		m.state = StateProviderList
	case "n", "esc":
		m.state = StateProviderList
	case "ctrl+c", "ctrl+q":
		return m, tea.Quit
	}
	return m, nil
}

// setProviderErrMsg 设置错误消息并在5秒后自动清除
func (m *Model) setProviderErrMsg(msg string) (tea.Model, tea.Cmd) {
	m.errMsg = msg
	return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return clearProviderMessageMsg{msgType: msgTypeErr}
	})
}

// setProviderTipMsg 设置提示消息并在5秒后自动清除
func (m *Model) setProviderTipMsg(msg string) (tea.Model, tea.Cmd) {
	m.tipMsg = msg
	return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return clearProviderMessageMsg{msgType: msgTypeTip}
	})
}

// updateProviderFocus 更新焦点状态
func (m *Model) updateProviderFocus() {
	all := []*textinput.Model{
		&m.providerNameInput,
		&m.providerBaseURLInput,
		&m.providerAPIKeyInput,
		&m.providerModelInput,
		&m.providerThinkingModelInput,
		&m.providerDefaultHaikuInput,
		&m.providerDefaultSonnetInput,
		&m.providerDefaultOpusInput,
		&m.providerSubagentInput,
		&m.providerNonessentialInput,
		&m.providerNonstreamingInput,
		&m.providerEffortInput,
	}
	for i := range all {
		if ProviderInputFocus(i) == m.providerInputFocus {
			all[i].Focus()
		} else {
			all[i].Blur()
		}
	}
}

// updateProviderInput 更新焦点的输入框
func (m *Model) updateProviderInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.providerInputFocus {
	case FocusProviderName:
		m.providerNameInput, cmd = m.providerNameInput.Update(msg)
	case FocusProviderBaseURL:
		m.providerBaseURLInput, cmd = m.providerBaseURLInput.Update(msg)
	case FocusProviderAPIKey:
		m.providerAPIKeyInput, cmd = m.providerAPIKeyInput.Update(msg)
	case FocusProviderModel:
		m.providerModelInput, cmd = m.providerModelInput.Update(msg)
	case FocusProviderThinkingModel:
		m.providerThinkingModelInput, cmd = m.providerThinkingModelInput.Update(msg)
	case FocusProviderDefaultHaiku:
		m.providerDefaultHaikuInput, cmd = m.providerDefaultHaikuInput.Update(msg)
	case FocusProviderDefaultSonnet:
		m.providerDefaultSonnetInput, cmd = m.providerDefaultSonnetInput.Update(msg)
	case FocusProviderDefaultOpus:
		m.providerDefaultOpusInput, cmd = m.providerDefaultOpusInput.Update(msg)
	case FocusProviderSubagent:
		m.providerSubagentInput, cmd = m.providerSubagentInput.Update(msg)
	case FocusProviderNonessential:
		m.providerNonessentialInput, cmd = m.providerNonessentialInput.Update(msg)
	case FocusProviderNonstreaming:
		m.providerNonstreamingInput, cmd = m.providerNonstreamingInput.Update(msg)
	case FocusProviderEffort:
		m.providerEffortInput, cmd = m.providerEffortInput.Update(msg)
	}
	return cmd
}
