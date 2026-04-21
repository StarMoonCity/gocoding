# Gocoding

一款使用 Go 和 Bubble Tea 构建的 TUI 项目管理工具。

## 功能特性

- 项目管理，支持别名和描述
- 支持在 Claude Code、VSCode、OpenCode、Codex CLI 中打开项目
- AI 模型提供商配置
- 调试模式（开发时使用）

## 使用方法

```bash
# 运行应用
gocoding

# 开启调试模式
gocoding -d

# 直接打开提供商配置
gocoding -p
```

## 快捷键

- `↑/↓` 或 `j/k` - 导航
- `n` - 添加新项目
- `r` - 重命名项目
- `d` - 删除项目
- `v` - 查看详情
- `e` - 编辑描述
- `/` - 搜索
- `p` - 打开提供商配置
- `q` - 退出

## 配置文件

项目数据：`~/.config/gocoding/projects.json`
提供商配置：`~/.config/gocoding/providers.json`
