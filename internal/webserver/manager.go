package webserver

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

type WebServerType string

const (
	Nginx  WebServerType = "nginx"
	Apache WebServerType = "apache"
	Caddy  WebServerType = "caddy"
)

func NewManager(dockerMgr *docker.Manager, cfg *config.Config) *Manager {
	return &Manager{
		docker: dockerMgr,
		config: cfg,
	}
}

// StartWebServer starts a web server with the specified type
func (m *Manager) StartWebServer(ctx context.Context, serverType WebServerType, port int) error {
	service := config.Service{
		Name:    string(serverType),
		Type:    string(serverType),
		Version: m.getDefaultVersion(serverType),
		Port:    port,
	}

	return m.docker.StartService(ctx, service)
}

// getDefaultVersion returns the default version for a web server type
func (m *Manager) getDefaultVersion(serverType WebServerType) string {
	defaults := map[WebServerType]string{
		Nginx:  "latest",
		Apache: "latest",
		Caddy:  "latest",
	}

	if version, ok := defaults[serverType]; ok {
		return version
	}

	return "latest"
}

// GetWebServerConfig returns configuration for a web server type
func (m *Manager) GetWebServerConfig(serverType WebServerType, projectPath string, domain string) string {
	switch serverType {
	case Nginx:
		return m.generateNginxConfig(projectPath, domain)
	case Apache:
		return m.generateApacheConfig(projectPath, domain)
	case Caddy:
		return m.generateCaddyConfig(projectPath, domain)
	default:
		return ""
	}
}

func (m *Manager) generateNginxConfig(projectPath, domain string) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s;
    root %s;

    index index.html index.htm index.php;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass php:9000;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.ht {
        deny all;
    }
}`, domain, projectPath)
}

func (m *Manager) generateApacheConfig(projectPath, domain string) string {
	return fmt.Sprintf(`<VirtualHost *:80>
    ServerName %s
    DocumentRoot %s

    <Directory %s>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>

    <FilesMatch \.php$>
        SetHandler "proxy:fcgi://php:9000"
    </FilesMatch>

    ErrorLog ${APACHE_LOG_DIR}/%s-error.log
    CustomLog ${APACHE_LOG_DIR}/%s-access.log combined
</VirtualHost>`, domain, projectPath, projectPath, domain, domain)
}

func (m *Manager) generateCaddyConfig(projectPath, domain string) string {
	return fmt.Sprintf(`%s {
    root * %s
    encode gzip
    php_fastcgi php:9000
    file_server
}`, domain, projectPath)
}

// GetAvailableWebServers returns list of available web servers
func GetAvailableWebServers() []struct {
	Name        string
	Type        WebServerType
	Description string
	DefaultPort int
} {
	return []struct {
		Name        string
		Type        WebServerType
		Description string
		DefaultPort int
	}{
		{
			Name:        "Nginx",
			Type:        Nginx,
			Description: "High-performance web server and reverse proxy",
			DefaultPort: 80,
		},
		{
			Name:        "Apache",
			Type:        Apache,
			Description: "Most popular web server with .htaccess support",
			DefaultPort: 80,
		},
		{
			Name:        "Caddy",
			Type:        Caddy,
			Description: "Modern web server with automatic HTTPS",
			DefaultPort: 80,
		},
	}
}
