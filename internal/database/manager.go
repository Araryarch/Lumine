package database

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

type DatabaseInfo struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	AdminURL string
}

func NewManager(dockerMgr *docker.Manager, cfg *config.Config) *Manager {
	return &Manager{
		docker: dockerMgr,
		config: cfg,
	}
}

// GetDatabaseInfo returns connection info for a database service
func (m *Manager) GetDatabaseInfo(serviceName string) (*DatabaseInfo, error) {
	for _, service := range m.config.Services {
		if service.Name == serviceName {
			info := &DatabaseInfo{
				Type: service.Type,
				Host: "localhost",
				Port: service.Port,
			}

			// Set default credentials based on type
			switch service.Type {
			case "mysql", "mariadb":
				info.Username = "root"
				info.Password = service.Env["MYSQL_ROOT_PASSWORD"]
				info.Database = service.Env["MYSQL_DATABASE"]
				info.AdminURL = "http://localhost:8080" // phpMyAdmin
			case "postgres", "postgresql":
				info.Username = "postgres"
				info.Password = service.Env["POSTGRES_PASSWORD"]
				info.Database = service.Env["POSTGRES_DB"]
				info.AdminURL = "http://localhost:8084" // pgAdmin
			case "mongodb", "mongo":
				info.Username = "root"
				info.Password = service.Env["MONGO_INITDB_ROOT_PASSWORD"]
				info.Database = service.Env["MONGO_INITDB_DATABASE"]
				info.AdminURL = "http://localhost:8082" // Mongo Express
			case "redis":
				info.AdminURL = "http://localhost:8083" // Redis Commander
			}

			return info, nil
		}
	}

	return nil, fmt.Errorf("database service not found: %s", serviceName)
}

// GetConnectionString returns a connection string for the database
func (m *Manager) GetConnectionString(serviceName string) (string, error) {
	info, err := m.GetDatabaseInfo(serviceName)
	if err != nil {
		return "", err
	}

	switch info.Type {
	case "mysql", "mariadb":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			info.Username, info.Password, info.Host, info.Port, info.Database), nil
	case "postgres", "postgresql":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			info.Host, info.Port, info.Username, info.Password, info.Database), nil
	case "mongodb", "mongo":
		return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			info.Username, info.Password, info.Host, info.Port, info.Database), nil
	case "redis":
		return fmt.Sprintf("redis://%s:%d", info.Host, info.Port), nil
	default:
		return "", fmt.Errorf("unsupported database type: %s", info.Type)
	}
}

// TestConnection tests the database connection
func (m *Manager) TestConnection(ctx context.Context, serviceName string) error {
	// Check if container is running
	containers, err := m.docker.ListContainers(ctx)
	if err != nil {
		return err
	}

	containerName := fmt.Sprintf("lumine-%s", serviceName)
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName && container.State == "running" {
				return nil // Container is running
			}
		}
	}

	return fmt.Errorf("database container is not running: %s", serviceName)
}

// GetAllDatabases returns info for all database services
func (m *Manager) GetAllDatabases() []*DatabaseInfo {
	var databases []*DatabaseInfo

	dbTypes := map[string]bool{
		"mysql":      true,
		"mariadb":    true,
		"postgres":   true,
		"postgresql": true,
		"mongodb":    true,
		"mongo":      true,
		"redis":      true,
	}

	for _, service := range m.config.Services {
		if dbTypes[service.Type] {
			info, err := m.GetDatabaseInfo(service.Name)
			if err == nil {
				databases = append(databases, info)
			}
		}
	}

	return databases
}
