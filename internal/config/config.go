package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ProjectsPath string `mapstructure:"projects_path"`
}

var v *viper.Viper

func Init() error {
	v = viper.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("cannot determine user home directory: " + err.Error())
	}
	configDir := filepath.Join(homeDir, ".config", "gocoding")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	v.SetDefault("projects_path", filepath.Join(configDir, "projects.json"))

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	if err := v.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); !ok {
			if writeErr := v.WriteConfig(); writeErr != nil {
				return writeErr
			}
		}
	}

	return nil
}

func GetConfig() *Config {
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil
	}
	return &cfg
}

func GetProjectsPath() string {
	return v.GetString("projects_path")
}
