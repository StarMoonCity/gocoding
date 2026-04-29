package components

import (
	"github.com/charmbracelet/bubbles/list"

	"gocoding/internal/models"
)

// listItem 项目列表项
type listItem struct {
	project models.Project
}

func (i listItem) Title() string {
	return i.project.Alias
}

func (i listItem) Description() string {
	return i.project.Path
}

func (i listItem) FilterValue() string { return i.project.Alias }

// newListItems 创建列表项
func newListItems(projects []models.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = listItem{project: proj}
	}
	return items
}
