package project

import (
	"github.com/Araryarch/lumine/internal/domain/project"
)

type Service struct {
	repo project.Repository
}

func NewService(repo project.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]project.Project, error) {
	return s.repo.List()
}

func (s *Service) Delete(name string) error {
	return s.repo.Delete(name)
}

func (s *Service) Create(name string, projectType project.Type, phpVersion string) error {
	return s.repo.Create(name, projectType, phpVersion)
}
