package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"lumine/internal/config"
)

type ComposeService struct {
	Image       string            `yaml:"image,omitempty"`
	Build       string            `yaml:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
}

type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]ComposeService `yaml:"services"`
	Networks map[string]interface{}    `yaml:"networks,omitempty"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

// GenerateForProject generates docker-compose.yml for a project
func GenerateForProject(project *config.Project, cfg *config.Config) error {
	compose := ComposeFile{
		Version:  "3.8",
		Services: make(map[string]ComposeService),
		Networks: map[string]interface{}{
			"lumine": map[string]string{"driver": "bridge"},
		},
	}

	// Add project service based on type
	switch project.Type {
	case "laravel":
		compose.Services["app"] = ComposeService{
			Image: fmt.Sprintf("php:%s", cfg.Runtimes.PHP),
			Ports: []string{fmt.Sprintf("%d:8000", project.Port)},
			Volumes: []string{
				fmt.Sprintf("%s:/var/www/html", project.Path),
			},
			Environment: project.Env,
			Networks:    []string{"lumine"},
		}
	case "nextjs", "vue", "express":
		compose.Services["app"] = ComposeService{
			Image: fmt.Sprintf("node:%s", cfg.Runtimes.Node),
			Ports: []string{fmt.Sprintf("%d:3000", project.Port)},
			Volumes: []string{
				fmt.Sprintf("%s:/app", project.Path),
			},
			Environment: project.Env,
			Networks:    []string{"lumine"},
		}
	case "django", "fastapi":
		compose.Services["app"] = ComposeService{
			Image: fmt.Sprintf("python:%s", cfg.Runtimes.Python),
			Ports: []string{fmt.Sprintf("%d:8000", project.Port)},
			Volumes: []string{
				fmt.Sprintf("%s:/app", project.Path),
			},
			Environment: project.Env,
			Networks:    []string{"lumine"},
		}
	}

	// Write docker-compose.yml
	composePath := filepath.Join(project.Path, "docker-compose.yml")
	data, err := yaml.Marshal(compose)
	if err != nil {
		return err
	}

	return os.WriteFile(composePath, data, 0644)
}
