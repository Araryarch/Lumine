package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"lumine/internal/docker"
)

// Common versions for popular services
var commonVersions = map[string][]string{
	"php": {
		"8.3-fpm", "8.3-apache", "8.3-cli",
		"8.2-fpm", "8.2-apache", "8.2-cli",
		"8.1-fpm", "8.1-apache", "8.1-cli",
		"8.0-fpm", "8.0-apache", "8.0-cli",
		"7.4-fpm", "7.4-apache", "7.4-cli",
	},
	"mysql": {
		"8.3", "8.2", "8.1", "8.0", "5.7", "5.6",
	},
	"mariadb": {
		"11.2", "11.1", "11.0", "10.11", "10.10", "10.6", "10.5",
	},
	"postgres": {
		"16", "15", "14", "13", "12", "11",
	},
	"postgresql": {
		"16", "15", "14", "13", "12", "11",
	},
	"nginx": {
		"latest", "1.25", "1.24", "1.23", "1.22", "alpine",
	},
	"apache": {
		"latest", "2.4", "2.4-alpine",
	},
	"httpd": {
		"latest", "2.4", "2.4-alpine",
	},
	"redis": {
		"7.2", "7.0", "6.2", "6.0", "alpine",
	},
	"mongodb": {
		"7.0", "6.0", "5.0", "4.4",
	},
	"phpmyadmin": {
		"latest", "5.2", "5.1", "5.0",
	},
	"adminer": {
		"latest", "4.8", "4.7",
	},
	"elasticsearch": {
		"8.11.0", "8.10.0", "7.17.0", "7.16.0",
	},
	"rabbitmq": {
		"3.12", "3.11", "3.10", "3.12-management", "3.11-management",
	},
	"memcached": {
		"latest", "1.6", "1.5", "alpine",
	},
}

func (m model) fetchVersions(serviceType string) tea.Cmd {
	return func() tea.Msg {
		// First check if we have common versions
		if versions, ok := commonVersions[serviceType]; ok {
			return versionListMsg{
				serviceType: serviceType,
				versions:    versions,
			}
		}

		// Otherwise try to fetch from Docker Hub
		versions, err := docker.FetchDockerHubTags(serviceType)
		if err != nil || len(versions) == 0 {
			// Fallback to generic versions
			versions = []string{"latest", "stable", "alpine"}
		}

		return versionListMsg{
			serviceType: serviceType,
			versions:    versions,
		}
	}
}
