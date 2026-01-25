package main

import (
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

	m := ui.NewModel(projectStore)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	projectsPath := config.GetProjectsPath()
	if err := projectStore.Save(projectsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save projects: %v\n", err)
		os.Exit(1)
	}
}
