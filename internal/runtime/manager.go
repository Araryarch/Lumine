package runtime

import (
	"context"
	"fmt"

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

// SwitchRuntime changes the runtime version
func (m *Manager) SwitchRuntime(runtime, version string) error {
	switch runtime {
	case "php":
		m.config.Runtimes.PHP = version
	case "node":
		m.config.Runtimes.Node = version
	case "python":
		m.config.Runtimes.Python = version
	case "bun":
		m.config.Runtimes.Bun = version
	case "deno":
		m.config.Runtimes.Deno = version
	case "go":
		m.config.Runtimes.Go = version
	default:
		return fmt.Errorf("unsupported runtime: %s", runtime)
	}

	return config.SaveConfig(m.config)
}

// GetRuntimeVersion returns the current version of a runtime
func (m *Manager) GetRuntimeVersion(runtime string) string {
	switch runtime {
	case "php":
		return m.config.Runtimes.PHP
	case "node":
		return m.config.Runtimes.Node
	case "python":
		return m.config.Runtimes.Python
	case "bun":
		return m.config.Runtimes.Bun
	case "deno":
		return m.config.Runtimes.Deno
	case "go":
		return m.config.Runtimes.Go
	default:
		return "unknown"
	}
}
