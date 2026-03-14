package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Araryarch/lumine/internal/domain/project"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
)

type ProjectRepository struct {
	projectsDir string
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		projectsDir: config.GetProjectsDir(),
	}
}

func (r *ProjectRepository) List() ([]project.Project, error) {
	if _, err := os.Stat(r.projectsDir); os.IsNotExist(err) {
		return []project.Project{}, nil
	}

	entries, err := os.ReadDir(r.projectsDir)
	if err != nil {
		return nil, err
	}

	var projects []project.Project
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			projectPath := filepath.Join(r.projectsDir, name)
			projectType := r.GetType(projectPath)

			projects = append(projects, project.Project{
				Name: name,
				Type: projectType,
				Path: projectPath,
				URL:  fmt.Sprintf("http://%s.test", name),
			})
		}
	}

	return projects, nil
}

func (r *ProjectRepository) GetType(projectPath string) project.Type {
	if fileExists(filepath.Join(projectPath, "artisan")) {
		return project.TypeLaravel
	}
	if fileExists(filepath.Join(projectPath, "wp-config.php")) {
		return project.TypeWordPress
	}
	if fileExists(filepath.Join(projectPath, "package.json")) {
		return project.TypeNodeJS
	}
	if fileExists(filepath.Join(projectPath, "composer.json")) {
		return project.TypePHP
	}
	if fileExists(filepath.Join(projectPath, "index.html")) {
		return project.TypeStatic
	}
	return project.TypeUnknown
}

func (r *ProjectRepository) Delete(name string) error {
	projectPath := filepath.Join(r.projectsDir, name)
	nginxConf := filepath.Join(config.GetConfigDir(), "nginx", name+".conf")

	if err := os.RemoveAll(projectPath); err != nil {
		return err
	}

	if err := os.Remove(nginxConf); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (r *ProjectRepository) Create(name string, projectType project.Type, phpVersion string) error {
	return fmt.Errorf("not implemented")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
