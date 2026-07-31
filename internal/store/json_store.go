package store

import (
	"gocoding/internal/config"
	"gocoding/internal/models"
)

type JSONStore struct {
	filePath string
	store    *models.ProjectStore
}

func NewJSONStore(store *models.ProjectStore) *JSONStore {
	return &JSONStore{
		filePath: config.GetProjectsPath(),
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
