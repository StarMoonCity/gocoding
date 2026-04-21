package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"gocoding/internal/models"
)

type IDEExecutor struct{}

func NewIDEExecutor() *IDEExecutor {
	return &IDEExecutor{}
}

func (e *IDEExecutor) OpenProject(project *models.Project, ideType models.IDEType) error {
	var cmd *exec.Cmd

	switch ideType {
	case models.IDEClaudeCode:
		if runtime.GOOS == "darwin" {
			// 检查项目路径是否存在
			if _, err := os.Stat(project.Path); os.IsNotExist(err) {
				return fmt.Errorf("项目路径不存在: %s", project.Path)
			}
			// 通过 iTerm2 打开新窗口，进入项目目录并运行 claude
			script := fmt.Sprintf(`
				tell application "iTerm2"
					create window with default profile
					tell current session of current window
						write text "cd %s && claude"
					end tell
				end tell
			`, project.Path)
			cmd = exec.Command("osascript", "-e", script)
		} else {
			cmd = exec.Command("claude", project.Path)
		}
	case models.IDEVSCode:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "Visual Studio Code", project.Path)
		} else {
			cmd = exec.Command("code", project.Path)
		}
	case models.IDEOpenCode:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", "-a", "OpenCode", project.Path)
		} else {
			cmd = exec.Command("opencode", project.Path)
		}
	case models.IDECodexCLI:
		if runtime.GOOS == "darwin" {
			// 检查项目路径是否存在
			if _, err := os.Stat(project.Path); os.IsNotExist(err) {
				return fmt.Errorf("项目路径不存在: %s", project.Path)
			}
			// 通过 iTerm2 打开新窗口，进入项目目录并运行 codex
			script := fmt.Sprintf(`
				tell application "iTerm2"
					create window with default profile
					tell current session of current window
						write text "cd %s && codex"
					end tell
				end tell
			`, project.Path)
			cmd = exec.Command("osascript", "-e", script)
		} else {
			cmd = exec.Command("codex", project.Path)
		}
	default:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("open", project.Path)
		} else {
			cmd = exec.Command("xdg-open", project.Path)
		}
	}

	if cmd != nil {
		return cmd.Start()
	}
	return nil
}

func (e *IDEExecutor) IsIDEAvailable(ideType models.IDEType) bool {
	var cmd *exec.Cmd
	switch ideType {
	case models.IDEClaudeCode:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("sh", "-c", "osascript -e 'id of app \"Claude\"' 2>/dev/null")
		} else {
			cmd = exec.Command("which", "claude")
		}
	case models.IDEVSCode:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("sh", "-c", "osascript -e 'id of app \"Visual Studio Code\"' 2>/dev/null")
		} else {
			cmd = exec.Command("which", "code")
		}
	case models.IDEOpenCode:
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("sh", "-c", "osascript -e 'id of app \"OpenCode\"' 2>/dev/null")
		} else {
			cmd = exec.Command("which", "opencode")
		}
	case models.IDECodexCLI:
		cmd = exec.Command("which", "codex")
	default:
		return false
	}
	return cmd.Run() == nil
}
