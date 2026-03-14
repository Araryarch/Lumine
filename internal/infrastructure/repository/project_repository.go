package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	projectPath := filepath.Join(r.projectsDir, name)

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	switch projectType {
	case project.TypeStatic:
		return r.createStaticProject(name, projectPath)
	case project.TypeNodeJS:
		return r.createNodeJSProject(name, projectPath)
	case project.TypePHP:
		return r.createPHPProject(name, projectPath, phpVersion)
	case project.TypeLaravel:
		return r.createLaravelProject(name, projectPath, phpVersion)
	case project.TypeWordPress:
		return r.createWordPressProject(name, projectPath)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

func (r *ProjectRepository) createStaticProject(name, projectPath string) error {
	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>Welcome to %s</h1>
    <p>Your static site is ready!</p>
</body>
</html>`, name, name)

	return os.WriteFile(filepath.Join(projectPath, "index.html"), []byte(indexHTML), 0644)
}

func (r *ProjectRepository) createNodeJSProject(name, projectPath string) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "main": "index.js",
  "scripts": {
    "start": "node index.js"
  }
}`, name)

	if err := os.WriteFile(filepath.Join(projectPath, "package.json"), []byte(packageJSON), 0644); err != nil {
		return err
	}

	indexJS := `const http = require('http');

const server = http.createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'text/html' });
  res.end('<h1>Node.js Server Running!</h1>');
});

const PORT = process.env.PORT || 3000;
server.listen(PORT, () => {
  console.log('Server running on port ' + PORT);
});
`

	return os.WriteFile(filepath.Join(projectPath, "index.js"), []byte(indexJS), 0644)
}

func (r *ProjectRepository) createPHPProject(name, projectPath, phpVersion string) error {
	if err := os.MkdirAll(filepath.Join(projectPath, "public"), 0755); err != nil {
		return err
	}

	indexPHP := `<?php
phpinfo();
`
	if err := os.WriteFile(filepath.Join(projectPath, "public", "index.php"), []byte(indexPHP), 0644); err != nil {
		return err
	}

	return r.createNginxConfig(name, phpVersion)
}

func (r *ProjectRepository) createLaravelProject(name, projectPath, phpVersion string) error {
	laravelReadme := `Laravel Project
================

To complete setup, run:

docker run --rm -v %s:/app -w /app composer create-project laravel/laravel .

Then configure your nginx to point to /public
`
	return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(fmt.Sprintf(laravelReadme, projectPath)), 0644)
}

func (r *ProjectRepository) createWordPressProject(name, projectPath string) error {
	wordpressReadme := `WordPress Project
================

To complete setup, run:

docker run --rm -v %s:/var/www/html wordpress:cli core download
`
	return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(fmt.Sprintf(wordpressReadme, projectPath)), 0644)
}

func (r *ProjectRepository) createNginxConfig(projectName, phpVersion string) error {
	nginxDir := filepath.Join(config.GetConfigDir(), "nginx")
	if err := os.MkdirAll(nginxDir, 0755); err != nil {
		return err
	}

	phpContainer := strings.ReplaceAll(phpVersion, ".", "")
	if phpContainer == "" {
		phpContainer = "82"
	}

	vhostContent := fmt.Sprintf(`server {
    listen 80;
    server_name %s.test;
    root /var/www/html/%s/public;
    
    index index.php index.html;
    
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
    
    location ~ \.php$ {
        fastcgi_pass php%s:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}`, projectName, projectName, phpContainer)

	return os.WriteFile(filepath.Join(nginxDir, projectName+".conf"), []byte(vhostContent), 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
