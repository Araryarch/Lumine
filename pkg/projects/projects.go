package projects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Araryarch/lumine/pkg/config"
)

type ProjectType string

const (
	TypePHP       ProjectType = "PHP"
	TypeLaravel   ProjectType = "Laravel"
	TypeWordPress ProjectType = "WordPress"
	TypeNodeJS    ProjectType = "Node.js"
	TypeStatic    ProjectType = "Static"
	TypeUnknown   ProjectType = "Unknown"
)

type Project struct {
	Name string
	Type ProjectType
	Path string
	URL  string
}

func ListProjects() ([]Project, error) {
	projectsDir := config.GetProjectsDir()
	
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return []Project{}, nil
	}

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, err
	}

	var projects []Project
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			projectPath := filepath.Join(projectsDir, name)
			projectType := DetectProjectType(projectPath)
			
			projects = append(projects, Project{
				Name: name,
				Type: projectType,
				Path: projectPath,
				URL:  fmt.Sprintf("http://%s.test", name),
			})
		}
	}

	return projects, nil
}

func DetectProjectType(projectPath string) ProjectType {
	// Check for Laravel
	if fileExists(filepath.Join(projectPath, "artisan")) {
		return TypeLaravel
	}

	// Check for WordPress
	if fileExists(filepath.Join(projectPath, "wp-config.php")) {
		return TypeWordPress
	}

	// Check for Node.js
	if fileExists(filepath.Join(projectPath, "package.json")) {
		return TypeNodeJS
	}

	// Check for PHP
	if fileExists(filepath.Join(projectPath, "composer.json")) {
		return TypePHP
	}

	// Check for static HTML
	if fileExists(filepath.Join(projectPath, "index.html")) {
		return TypeStatic
	}

	return TypeUnknown
}

func CreatePHPProject(name string, phpVersion string) error {
	projectPath := filepath.Join(config.GetProjectsDir(), name)
	
	if err := os.MkdirAll(filepath.Join(projectPath, "public"), 0755); err != nil {
		return err
	}

	indexPHP := `<?php
phpinfo();
`
	if err := os.WriteFile(filepath.Join(projectPath, "public", "index.php"), []byte(indexPHP), 0644); err != nil {
		return err
	}

	return CreateNginxVHost(name, phpVersion, "php")
}

func CreateLaravelProject(name string, phpVersion string) error {
	projectsDir := config.GetProjectsDir()
	
	// Use composer docker image to create Laravel project
	cmd := fmt.Sprintf("docker run --rm -v %s:/app -w /app composer:latest create-project laravel/laravel %s --prefer-dist", projectsDir, name)
	
	// This would need to be executed via exec
	// For now, return instruction
	return fmt.Errorf("run: %s", cmd)
}

func CreateWordPressProject(name string, phpVersion string) error {
	projectPath := filepath.Join(config.GetProjectsDir(), name)
	
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	// Use WordPress CLI docker image
	cmd := fmt.Sprintf("docker run --rm -v %s:/var/www/html wordpress:cli core download", projectPath)
	
	return fmt.Errorf("run: %s", cmd)
}

func CreateNodeJSProject(name string) error {
	projectPath := filepath.Join(config.GetProjectsDir(), name)
	
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

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

func CreateStaticProject(name string) error {
	projectPath := filepath.Join(config.GetProjectsDir(), name)
	
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

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

	if err := os.WriteFile(filepath.Join(projectPath, "index.html"), []byte(indexHTML), 0644); err != nil {
		return err
	}

	return CreateNginxVHost(name, "", "static")
}

func CreateNginxVHost(projectName string, phpVersion string, projectType string) error {
	nginxDir := filepath.Join(config.GetConfigDir(), "nginx")
	if err := os.MkdirAll(nginxDir, 0755); err != nil {
		return err
	}

	vhostFile := filepath.Join(nginxDir, projectName+".conf")
	
	var vhostContent string
	phpContainer := strings.ReplaceAll(phpVersion, ".", "")

	switch projectType {
	case "laravel":
		vhostContent = fmt.Sprintf(`server {
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

	case "static":
		vhostContent = fmt.Sprintf(`server {
    listen 80;
    server_name %s.test;
    root /var/www/html/%s;
    
    index index.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
}`, projectName, projectName)

	default:
		vhostContent = fmt.Sprintf(`server {
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
	}

	return os.WriteFile(vhostFile, []byte(vhostContent), 0644)
}

func DeleteProject(name string) error {
	projectPath := filepath.Join(config.GetProjectsDir(), name)
	nginxConf := filepath.Join(config.GetConfigDir(), "nginx", name+".conf")
	
	// Remove project directory
	if err := os.RemoveAll(projectPath); err != nil {
		return err
	}

	// Remove nginx config
	if err := os.Remove(nginxConf); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
