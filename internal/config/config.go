package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Services []Service `yaml:"services"`
	Projects []Project `yaml:"projects"`
	Runtimes Runtimes  `yaml:"runtimes"`
}

type Service struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"` // php, mysql, nginx, etc
	Version string            `yaml:"version"`
	Port    int               `yaml:"port"`
	Env     map[string]string `yaml:"env,omitempty"`
}

type Project struct {
	Name      string            `yaml:"name"`
	Type      string            `yaml:"type"` // laravel, nextjs, vue, django, etc
	Path      string            `yaml:"path"`
	Domain    string            `yaml:"domain"` // e.g., myapp.test
	Runtime   string            `yaml:"runtime"` // php, node, python, bun, deno
	Version   string            `yaml:"version"`
	Port      int               `yaml:"port"`
	Env       map[string]string `yaml:"env,omitempty"`
	Status    string            `yaml:"status"` // running, stopped
}

type Runtimes struct {
	PHP    string `yaml:"php"`
	Node   string `yaml:"node"`
	Python string `yaml:"python"`
	Bun    string `yaml:"bun"`
	Deno   string `yaml:"deno"`
	Go     string `yaml:"go"`
	Rust   string `yaml:"rust"`
}

var (
	ConfigDir  string
	ConfigFile string
)

func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	ConfigDir = filepath.Join(home, ".lumine")
	ConfigFile = filepath.Join(ConfigDir, "config.yaml")

	// Create config directory if not exists
	if err := os.MkdirAll(ConfigDir, 0755); err != nil {
		return err
	}

	// Create default config if not exists
	if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
		defaultConfig := Config{
			Services: []Service{
				{Name: "nginx", Type: "nginx", Version: "latest", Port: 80},
				{Name: "mysql", Type: "mysql", Version: "8.0", Port: 3306, Env: map[string]string{"MYSQL_ROOT_PASSWORD": "root", "MYSQL_DATABASE": "lumine"}},
				{Name: "redis", Type: "redis", Version: "7.2", Port: 6379},
			},
			Projects: []Project{},
			Runtimes: Runtimes{
				PHP:    "8.2-fpm",
				Node:   "20-alpine",
				Python: "3.11-slim",
				Bun:    "latest",
				Deno:   "latest",
				Go:     "1.21-alpine",
				Rust:   "latest",
			},
		}
		return SaveConfig(&defaultConfig)
	}

	return nil
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(ConfigFile, data, 0644)
}
