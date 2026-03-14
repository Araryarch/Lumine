package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version     string              `yaml:"version"`
	DefaultPHP  string              `yaml:"default_php"`
	PHPVersions []string            `yaml:"php_versions"`
	Services    map[string]*Service `yaml:"services"`
	Projects    []string            `yaml:"projects"`
}

type Service struct {
	Image   string            `yaml:"image"`
	Port    int               `yaml:"port"`
	Enabled bool              `yaml:"enabled"`
	Env     map[string]string `yaml:"env,omitempty"`
}

func Load() (*Config, error) {
	configDir := GetConfigDir()
	configFile := filepath.Join(configDir, "config.yaml")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefault(configFile); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func createDefault(path string) error {
	cfg := &Config{
		Version:     "2.0.0",
		DefaultPHP:  "8.2",
		PHPVersions: []string{"7.4", "8.0", "8.1", "8.2", "8.3"},
		Services: map[string]*Service{
			"nginx": {
				Image:   "nginx:alpine",
				Port:    80,
				Enabled: true,
			},
			"mysql": {
				Image:   "mysql:8.0",
				Port:    3306,
				Enabled: true,
				Env: map[string]string{
					"MYSQL_ROOT_PASSWORD": "root",
				},
			},
			"redis": {
				Image:   "redis:7-alpine",
				Port:    6379,
				Enabled: true,
			},
			"mailhog": {
				Image:   "mailhog/mailhog",
				Port:    8025,
				Enabled: true,
			},
			"phpmyadmin": {
				Image:   "phpmyadmin:latest",
				Port:    8080,
				Enabled: true,
			},
		},
		Projects: []string{},
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func GetConfigDir() string {
	return filepath.Join(os.Getenv("HOME"), ".lumine")
}

func GetProjectsDir() string {
	return filepath.Join(os.Getenv("HOME"), "lumine-projects")
}
