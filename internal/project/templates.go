package project

import (
	"context"
	"fmt"
	"os/exec"
)

func (m *Manager) createLaravelProject(ctx context.Context, name, path string) error {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		fmt.Sprintf("php:%s", m.config.Runtimes.PHP),
		"composer", "create-project", "laravel/laravel", name)
	return cmd.Run()
}

func (m *Manager) createNextJSProject(ctx context.Context, name, path string) error {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		fmt.Sprintf("node:%s", m.config.Runtimes.Node),
		"npx", "create-next-app@latest", name, "--typescript", "--tailwind", "--app")
	return cmd.Run()
}

func (m *Manager) createVueProject(ctx context.Context, name, path string) error {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		fmt.Sprintf("node:%s", m.config.Runtimes.Node),
		"npm", "create", "vue@latest", name)
	return cmd.Run()
}

func (m *Manager) createDjangoProject(ctx context.Context, name, path string) error {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		fmt.Sprintf("python:%s", m.config.Runtimes.Python),
		"django-admin", "startproject", name)
	return cmd.Run()
}

func (m *Manager) createExpressProject(ctx context.Context, name, path string) error {
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		fmt.Sprintf("node:%s", m.config.Runtimes.Node),
		"npx", "express-generator", name)
	return cmd.Run()
}

func (m *Manager) createFastAPIProject(ctx context.Context, name, path string) error {
	// Create basic FastAPI structure
	return fmt.Errorf("FastAPI project creation not yet implemented")
}

func (m *Manager) createAxumProject(ctx context.Context, name, path string) error {
	// Create Axum project using cargo
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		"rust:latest",
		"sh", "-c",
		fmt.Sprintf("cargo new %s && cd %s && cargo add axum tokio --features tokio/full", name, name))
	return cmd.Run()
}

func (m *Manager) createActixProject(ctx context.Context, name, path string) error {
	// Create Actix-web project
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		"rust:latest",
		"sh", "-c",
		fmt.Sprintf("cargo new %s && cd %s && cargo add actix-web actix-rt", name, name))
	return cmd.Run()
}

func (m *Manager) createRocketProject(ctx context.Context, name, path string) error {
	// Create Rocket project
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", path),
		"-w", "/app",
		"rust:latest",
		"sh", "-c",
		fmt.Sprintf("cargo new %s && cd %s && cargo add rocket", name, name))
	return cmd.Run()
}
