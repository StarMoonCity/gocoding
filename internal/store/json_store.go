package store

import (
	"os"
	"path/filepath"

	"gocoding/internal/models"
)

const defaultFileName = "projects.json"

type JSONStore struct {
	filePath string
	store    *models.ProjectStore
}

func NewJSONStore(store *models.ProjectStore) *JSONStore {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "gocoding")
	return &JSONStore{
		filePath: filepath.Join(configDir, defaultFileName),
		store:    store,
	}
}

func (s *JSONStore) SetFilePath(path string) {
	s.filePath = path
}

func (s *JSONStore) Load() error {
	return s.store.Load(s.filePath)
}

func (s *JSONStore) Save() error {
	return s.store.Save(s.filePath)
}

func (s *JSONStore) AddProject(project models.Project) {
	s.store.Add(project)
}

func (s *JSONStore) RemoveProject(id string) {
	s.store.Remove(id)
}

func (s *JSONStore) UpdateProject(id, alias, path string) {
	s.store.Update(id, alias, path)
}

func (s *JSONStore) GetProject(id string) *models.Project {
	return s.store.Get(id)
}

func (s *JSONStore) GetProjectByIndex(index int) *models.Project {
	return s.store.GetByIndex(index)
}

func (s *JSONStore) GetAllProjects() []models.Project {
	return s.store.Projects
}

func (s *JSONStore) Len() int {
	return s.store.Len()
}
