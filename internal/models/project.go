package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type IDEType string

const (
	IDEClaudeCode IDEType = "claude"
	IDEVSCode     IDEType = "code"
	IDEOpenCode   IDEType = "opencode"
)

type Project struct {
	ID         string    `json:"id"`
	Path       string    `json:"path"`
	Alias      string    `json:"alias"`
	Description string   `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
	LastOpened time.Time `json:"last_opened"`
	OpenCount  int       `json:"open_count"`
}

func (p *Project) UpdateLastOpened() {
	p.LastOpened = time.Now()
	p.OpenCount++
}

type ProjectStore struct {
	Projects []Project `json:"projects"`
}

func NewProjectStore() *ProjectStore {
	return &ProjectStore{
		Projects: make([]Project, 0),
	}
}

func (s *ProjectStore) Add(project Project) {
	s.Projects = append(s.Projects, project)
}

func (s *ProjectStore) Remove(id string) {
	for i, p := range s.Projects {
		if p.ID == id {
			s.Projects = append(s.Projects[:i], s.Projects[i+1:]...)
			return
		}
	}
}

func (s *ProjectStore) Get(id string) *Project {
	for i := range s.Projects {
		if s.Projects[i].ID == id {
			return &s.Projects[i]
		}
	}
	return nil
}

func (s *ProjectStore) Update(alias string, id string) {
	for i := range s.Projects {
		if s.Projects[i].ID == id {
			s.Projects[i].Alias = alias
			return
		}
	}
}

func (s *ProjectStore) UpdateDescription(id string, description string) {
	for i := range s.Projects {
		if s.Projects[i].ID == id {
			s.Projects[i].Description = description
			return
		}
	}
}

func (s *ProjectStore) GetByIndex(index int) *Project {
	if index < 0 || index >= len(s.Projects) {
		return nil
	}
	return &s.Projects[index]
}

func (s *ProjectStore) Len() int {
	return len(s.Projects)
}

func (s *ProjectStore) SortByLastOpened() {
	s.Projects = sortProjectsByLastOpened(s.Projects)
}

// Search 搜索项目（模糊匹配别名、路径、描述）
func (s *ProjectStore) Search(query string) []Project {
	if query == "" {
		return s.Projects
	}
	var results []Project
	lowerQuery := toLower(query)
	for _, p := range s.Projects {
		if containsLower(p.Alias, lowerQuery) ||
			containsLower(p.Path, lowerQuery) ||
			containsLower(p.Description, lowerQuery) {
			results = append(results, p)
		}
	}
	return results
}

// GetMostRecentlyOpened 返回最近打开的项目
func (s *ProjectStore) GetMostRecentlyOpened() *Project {
	if len(s.Projects) == 0 {
		return nil
	}
	s.SortByLastOpened()
	return &s.Projects[0]
}

// GetIndexByProject 返回项目在列表中的索引
func (s *ProjectStore) GetIndexByProject(id string) int {
	for i, p := range s.Projects {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

func containsLower(s, substr string) bool {
	return toLower(s) == toLower(substr) ||
		len(s) > len(substr) && toLower(s[:len(substr)]) == substr ||
		len(s) > len(substr) && toLower(s[len(s)-len(substr):]) == substr ||
		len(s) > len(substr) && containsSubstring(toLower(s), substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func sortProjectsByLastOpened(projects []Project) []Project {
	for i := 0; i < len(projects)-1; i++ {
		for j := i + 1; j < len(projects); j++ {
			if projects[j].LastOpened.After(projects[i].LastOpened) {
				projects[i], projects[j] = projects[j], projects[i]
			}
		}
	}
	return projects
}

func (s *ProjectStore) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, s)
}

func (s *ProjectStore) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
