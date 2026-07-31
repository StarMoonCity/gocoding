package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestInitReadsCustomProjectsPath 验证 Init 会读取已有配置中的 projects_path
func TestInitReadsCustomProjectsPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".config", "gocoding")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	customPath := filepath.Join(home, "data", "projects.json")
	content := "projects_path: " + customPath + "\n"
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if got := GetProjectsPath(); got != customPath {
		t.Errorf("GetProjectsPath() = %q, want %q", got, customPath)
	}
}

// TestInitDefaultsWithoutConfig 验证首次运行（无配置文件）使用默认路径
func TestInitDefaultsWithoutConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	want := filepath.Join(home, ".config", "gocoding", "projects.json")
	if got := GetProjectsPath(); got != want {
		t.Errorf("GetProjectsPath() = %q, want %q", got, want)
	}
}
