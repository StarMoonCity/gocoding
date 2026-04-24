package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"gocoding/internal/config"
	"gocoding/internal/models"
	"gocoding/internal/store"
	"gocoding/internal/ui"
)

func main() {
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to init config: %v\n", err)
		os.Exit(1)
	}

	projectStore := models.NewProjectStore()
	jsonStore := store.NewJSONStore(projectStore)

	if err := jsonStore.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load projects: %v\n", err)
		os.Exit(1)
	}

	projectStore.SortByLastOpened()

	// 初始化模型提供商配置
	providerStore := models.NewModelProviderStore()
	providerConfig := config.NewModelProviderConfig(providerStore)

	if err := providerConfig.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load provider config: %v\n", err)
		os.Exit(1)
	}

	m := ui.NewModel(projectStore)
	m.SetProviderStore(providerStore)

	// 解析命令行标志
	providerMode := flag.Bool("p", false, "直接打开模型配置界面")
	debugMode := flag.Bool("d", false, "开启调试模式")
	flag.Parse()
	if *providerMode {
		m.SetState(ui.StateProviderList)
	}
	if *debugMode {
		m.SetDebugMode(true)
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	projectsPath := config.GetProjectsPath()
	if err := projectStore.Save(projectsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save projects: %v\n", err)
		os.Exit(1)
	}

	// 保存模型提供商配置
	if err := providerConfig.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save provider config: %v\n", err)
		os.Exit(1)
	}
}
