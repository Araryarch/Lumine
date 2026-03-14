package repository

import (
	"fmt"

	"github.com/Araryarch/lumine/internal/domain/service"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
	"github.com/Araryarch/lumine/internal/infrastructure/docker"
)

type ServiceRepository struct {
	docker *docker.Client
	config *config.Config
}

func NewServiceRepository(docker *docker.Client, cfg *config.Config) *ServiceRepository {
	return &ServiceRepository{
		docker: docker,
		config: cfg,
	}
}

func (r *ServiceRepository) GetAll() ([]service.Service, error) {
	var services []service.Service
	for name, svc := range r.config.Services {
		if !svc.Enabled {
			continue
		}
		services = append(services, service.Service{
			Name:    name,
			Image:   svc.Image,
			Port:    svc.Port,
			Enabled: svc.Enabled,
			Env:     svc.Env,
		})
	}
	return services, nil
}

func (r *ServiceRepository) GetStatus(name string) (service.Status, error) {
	containerName := "lumine-" + name
	state, err := r.docker.GetContainerStatus(containerName)
	if err != nil {
		return service.Status{}, err
	}

	svc, exists := r.config.Services[name]
	if !exists {
		return service.Status{}, fmt.Errorf("service not found: %s", name)
	}

	return service.Status{
		Name:    name,
		Running: state == "running",
		State:   state,
		Port:    svc.Port,
		Image:   svc.Image,
	}, nil
}

func (r *ServiceRepository) Start(name string) error {
	return r.docker.StartContainer("lumine-" + name)
}

func (r *ServiceRepository) Stop(name string) error {
	return r.docker.StopContainer("lumine-" + name)
}

func (r *ServiceRepository) Restart(name string) error {
	return r.docker.RestartContainer("lumine-" + name)
}

func (r *ServiceRepository) StartAll() error {
	services, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := r.Start(svc.Name); err != nil {
			return fmt.Errorf("failed to start %s: %w", svc.Name, err)
		}
	}
	return nil
}

func (r *ServiceRepository) StopAll() error {
	services, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := r.Stop(svc.Name); err != nil {
			return fmt.Errorf("failed to stop %s: %w", svc.Name, err)
		}
	}
	return nil
}

func (r *ServiceRepository) RestartAll() error {
	services, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, svc := range services {
		if err := r.Restart(svc.Name); err != nil {
			return fmt.Errorf("failed to restart %s: %w", svc.Name, err)
		}
	}
	return nil
}

func (r *ServiceRepository) GetLogs(name string, tail string) (string, error) {
	containerName := "lumine-" + name
	return r.docker.GetLogs(containerName, tail)
}
