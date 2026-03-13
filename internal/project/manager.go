package project

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"lumine/internal/config"
	"lumine/internal/docker"
)

type Manager struct {
	docker *docker.Manager
	config *config.Config
}

func NewManager(dockerMgr *docker.Manager, cfg *config.Config) *Manager {
	return &Manager{
		docker: dockerMgr,
		config: cfg,
	}
}

// CreateProject creates a new project with the specified framework
func (m *Manager) CreateProject(ctx context.Context, name, projectType, path string) error {
	switch projectType {
	case "laravel":
		return m.createLaravelProject(ctx, name, path)
	case "nextjs":
		return m.createNextJSProject(ctx, name, path)
	case "vue":
		return m.createVueProject(ctx, name, path)
	case "django":
		return m.createDjangoProject(ctx, name, path)
	case "express":
		return m.createExpressProject(ctx, name, path)
	case "fastapi":
		return m.createFastAPIProject(ctx, name, path)
	case "axum":
		return m.createAxumProject(ctx, name, path)
	case "actix":
		return m.createActixProject(ctx, name, path)
	case "rocket":
		return m.createRocketProject(ctx, name, path)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}
