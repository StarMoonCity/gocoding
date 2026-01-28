# Gocoding 功能说明文档

## 项目简介

Gocoding 是一个基于 Go 和 Bubble Tea 开发的 TUI（终端用户界面）项目管理工具。它提供了在终端中管理项目、打开 IDE 以及记录项目使用统计的功能，数据存储于 `~/.config/gocoding/projects.json`。

---

## 核心功能

### 项目管理

Gocoding 提供了完整的项目生命周期管理能力，支持以下操作：

- **添加项目**：输入项目路径后，系统会自动从路径中提取项目名称，用户也可以手动指定别名
- **重命名项目**：允许用户修改项目的显示名称
- **删除项目**：删除操作需要二次确认，防止误删
- **编辑描述**：为项目添加详细描述信息
- **查看详情**：展示项目的完整信息，包括创建时间、最后打开时间、打开次数等

### 搜索功能

内置模糊搜索功能，支持实时过滤项目列表。搜索时同时匹配项目别名和路径，快速定位目标项目。

### IDE 集成

支持一键从终端打开项目到多种 IDE：

| IDE | 快捷键 | 打开方式 |
|-----|--------|----------|
| Claude Code | `1` | 通过 iTerm2 创建新窗口启动 |
| Visual Studio Code | `2` | `code` 命令打开 |
| OpenCode | `3` | 系统应用启动 |

系统会自动检测 IDE 是否已安装，未安装的选项会显示为不可用状态。

---

## 界面状态说明

Gocoding 采用状态机模式管理界面，完整的状态包括：

| 状态 | 说明 |
|------|------|
| `StateList` | 主项目列表视图 |
| `StateAddProject` | 添加项目表单 |
| `StateRenameProject` | 重命名项目 |
| `StateDeleteConfirm` | 删除确认对话框 |
| `StateIDEMenu` | IDE 选择菜单 |
| `StateViewDetail` | 项目详情视图 |
| `StateEditDescription` | 编辑项目描述 |
| `StateSearch` | 搜索模式 |

---

## 快捷键参考

### 主列表模式（StateList）

| 按键 | 功能 |
|------|------|
| `↑` `↓` `j` `k` | 上下移动选择 |
| `n` | 添加新项目 |
| `/` | 进入搜索模式 |
| `r` | 重命名选中项目 |
| `d` | 删除选中项目 |
| `v` | 切换详情面板 |
| `e` | 编辑项目描述 |
| `Enter` | 打开 IDE 选择菜单 |
| `1` `2` `3` | 快速打开（对应 Claude/VSCode/OpenCode） |
| `q` `Ctrl+Q` `Ctrl+C` | 退出程序 |

### 添加项目模式（StateAddProject）

| 按键 | 功能 |
|------|------|
| `Enter` | 确认添加 |
| `Tab` | 切换输入焦点（路径与名称之间） |
| `Esc` | 返回列表 |

### IDE 选择菜单（StateIDEMenu）

| 按键 | 功能 |
|------|------|
| `↑` `↓` `j` `k` | 选择 IDE |
| `Enter` | 确认打开 |
| `1` `2` `3` | 快速选择 |
| `Esc` | 返回列表 |

### 搜索模式（StateSearch）

| 按键 | 功能 |
|------|------|
| 输入字符 | 搜索查询 |
| `Backspace` `Ctrl+H` | 删除字符 |
| `↑` `↓` `j` `k` | 选择项目 |
| `Enter` | 打开 IDE 选择菜单 |
| `1` `2` `3` | 快速打开 |
| `Esc` | 退出搜索 |

---

## 数据结构

### Project 结构体

| 字段 | 类型 | 说明 |
|------|------|------|
| `ID` | string | 唯一标识符 |
| `Path` | string | 项目路径 |
| `Alias` | string | 项目名称/别名 |
| `Description` | string | 项目描述 |
| `CreatedAt` | time.Time | 创建时间 |
| `LastOpened` | time.Time | 最后打开时间 |
| `OpenCount` | int | 打开次数 |

项目列表默认按最后打开时间降序排列，最近使用的项目显示在最前面。

---

## 文件结构

```
gocoding/
├── cmd/gocoding/
│   └── main.go              # 程序入口
└── internal/
    ├── models/
    │   └── project.go       # Project 数据结构
    ├── commands/
    │   └── executor.go      # IDE 执行器
    ├── store/
    │   └── json_store.go    # JSON 持久化
    ├── config/
    │   └── config.go        # Viper 配置管理
    └── ui/
        ├── main.go          # 主 Model 和状态定义
        ├── model.go         # UI Model
        ├── update.go        # 消息更新逻辑
        ├── view.go          # 视图渲染
        ├── styles.go        # 样式定义
        ├── layout.go        # 响应式布局
        ├── dialogs.go       # 对话框视图
        └── utils.go         # 工具函数
```

---

## 快速开始

### 构建项目

```bash
go build ./cmd/gocoding/
```

### 运行应用

```bash
go run ./cmd/gocoding/
```

### 添加依赖

```bash
go get github.com/charmbracelet/bubbletea@v1.1.0 \
      github.com/charmbracelet/lipgloss@v0.13.0 \
      github.com/charmbracelet/bubbles@v0.20.0 \
      github.com/spf13/viper@v1.18.2
```

---

## 技术栈

- **bubbletea**：终端 UI 框架
- **lipgloss**：终端样式库
- **bubbles**：UI 组件库（list、textinput 等）
- **viper**：配置管理
- **go-humanize**：人类可读的时间格式
