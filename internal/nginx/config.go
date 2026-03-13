package nginx

import (
	"fmt"
	"os"
	"path/filepath"

	"lumine/internal/config"
)

// GenerateConfig generates nginx config for a project
func GenerateConfig(project *config.Project) error {
	var template string

	switch project.Type {
	case "laravel":
		template = laravelNginxConfig(project)
	case "nextjs", "vue", "express":
		template = nodeNginxConfig(project)
	case "django", "fastapi":
		template = pythonNginxConfig(project)
	default:
		template = defaultNginxConfig(project)
	}

	configDir := filepath.Join(config.ConfigDir, "nginx")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, fmt.Sprintf("%s.conf", project.Name))
	return os.WriteFile(configPath, []byte(template), 0644)
}

func laravelNginxConfig(project *config.Project) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    root %s/public;

    index index.php index.html;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass app:%d;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}`, project.Domain, project.Path, project.Port)
}

func nodeNginxConfig(project *config.Project) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;

    location / {
        proxy_pass http://localhost:%d;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}`, project.Domain, project.Port)
}

func pythonNginxConfig(project *config.Project) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;

    location / {
        proxy_pass http://localhost:%d;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}`, project.Domain, project.Port)
}

func defaultNginxConfig(project *config.Project) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    root %s;

    index index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }
}`, project.Domain, project.Path)
}
