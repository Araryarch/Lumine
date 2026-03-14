package service

import (
	"github.com/Araryarch/lumine/internal/domain/service"
)

type Service struct {
	repo service.Repository
}

func NewService(repo service.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllStatuses() ([]service.Status, error) {
	services, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var statuses []service.Status
	for _, svc := range services {
		status, err := s.repo.GetStatus(svc.Name)
		if err != nil {
			continue
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (s *Service) StartAll() error {
	return s.repo.StartAll()
}

func (s *Service) StopAll() error {
	return s.repo.StopAll()
}

func (s *Service) RestartAll() error {
	return s.repo.RestartAll()
}

func (s *Service) Start(name string) error {
	return s.repo.Start(name)
}

func (s *Service) Stop(name string) error {
	return s.repo.Stop(name)
}

func (s *Service) Restart(name string) error {
	return s.repo.Restart(name)
}

func (s *Service) GetLogs(name string, tail string) (string, error) {
	return s.repo.GetLogs(name, tail)
}
