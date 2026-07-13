use std::fs;
use std::path::Path;
use tauri::{App, Manager};
use serde_json::json;

pub fn initialize_default_packages(app: &mut App) -> Result<(), String> {
    let app_data_dir = app.path().app_data_dir().map_err(|e| e.to_string())?;
    
    // We only need to check if services.json exists to avoid overwriting it
    let services_file = app_data_dir.join("services.json");
    if services_file.exists() {
        return Ok(());
    }

    // Ensure app_data_dir exists
    fs::create_dir_all(&app_data_dir).map_err(|e| e.to_string())?;

    register_default_services(&app_data_dir);

    Ok(())
}

fn register_default_services(app_data_dir: &Path) {
    let services_file = app_data_dir.join("services.json");
    
    // IMPORTANT: This JSON must match the exact structure that service.rs
    // deserializes via serde. ServiceInfo has nested ServiceConfig.
    // Serde uses snake_case by default for field names.
    let services = json!([
        // ===== Web Servers =====
        {
            "id": "nginx-default",
            "name": "Nginx",
            "description": "High-performance web server & reverse proxy",
            "status": "Stopped",
            "config": {
                "port": 8082,
                "container_port": 80,
                "executable_path": "nginx:alpine",
                "arguments": "",
                "service_type": "Web Server",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "apache-default",
            "name": "Apache",
            "description": "Classic & reliable web server with .htaccess support",
            "status": "Stopped",
            "config": {
                "port": 8090,
                "container_port": 80,
                "executable_path": "httpd:alpine",
                "arguments": "",
                "service_type": "Web Server",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "caddy-default",
            "name": "Caddy",
            "description": "Modern web server with automatic HTTPS",
            "status": "Stopped",
            "config": {
                "port": 8443,
                "container_port": 80,
                "executable_path": "caddy:alpine",
                "arguments": "",
                "service_type": "Web Server",
                "runner": "docker"
            },
            "log": []
        },

        // ===== Languages =====
        {
            "id": "php-default",
            "name": "PHP",
            "description": "PHP scripting language runtime",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "php:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "node-default",
            "name": "Node.js",
            "description": "Node.js JavaScript runtime for Next.js & Express apps",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "node:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "bun-default",
            "name": "Bun",
            "description": "Blazing-fast JavaScript runtime & bundler",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "oven/bun:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "python-default",
            "name": "Python",
            "description": "Python runtime for Django, Flask & FastAPI apps",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "python:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "go-default",
            "name": "Go",
            "description": "Go runtime for Gin, Fiber & Echo APIs",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "golang:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "deno-default",
            "name": "Deno",
            "description": "Secure TypeScript runtime with built-in tooling",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "denoland/deno:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "rust-default",
            "name": "Rust",
            "description": "Blazing fast and memory-safe language",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "rust:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "ruby-default",
            "name": "Ruby",
            "description": "Ruby language runtime for Rails apps",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "ruby:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "java-default",
            "name": "Java",
            "description": "Eclipse Temurin Java JDK",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "eclipse-temurin:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "gcc-default",
            "name": "C/C++",
            "description": "GCC compiler collection for C and C++",
            "status": "Stopped",
            "config": {
                "port": 0,
                "executable_path": "gcc:latest",
                "arguments": "",
                "service_type": "Language",
                "runner": "docker"
            },
            "log": []
        },

        // ===== Databases =====
        {
            "id": "mysql-default",
            "name": "MySQL",
            "description": "MySQL 8.0 open-source relational database",
            "status": "Stopped",
            "config": {
                "port": 3306,
                "executable_path": "mysql:8.0",
                "arguments": "-e MYSQL_ROOT_PASSWORD=root",
                "service_type": "Database",
                "runner": "docker",
                "volume_path": "/var/lib/mysql"
            },
            "log": []
        },
        {
            "id": "mariadb-default",
            "name": "MariaDB",
            "description": "MariaDB 11 MySQL-compatible database with extra features",
            "status": "Stopped",
            "config": {
                "port": 3307,
                "container_port": 3306,
                "executable_path": "mariadb:11",
                "arguments": "-e MARIADB_ROOT_PASSWORD=root",
                "service_type": "Database",
                "runner": "docker",
                "volume_path": "/var/lib/mysql"
            },
            "log": []
        },
        {
            "id": "postgres-default",
            "name": "PostgreSQL",
            "description": "PostgreSQL 16 advanced open-source relational database",
            "status": "Stopped",
            "config": {
                "port": 5432,
                "executable_path": "postgres:16-alpine",
                "arguments": "-e POSTGRES_PASSWORD=root",
                "service_type": "Database",
                "runner": "docker",
                "volume_path": "/var/lib/postgresql/data"
            },
            "log": []
        },
        {
            "id": "mongo-default",
            "name": "MongoDB",
            "description": "MongoDB 7 NoSQL document database for modern apps",
            "status": "Stopped",
            "config": {
                "port": 27017,
                "executable_path": "mongo:7",
                "arguments": "",
                "service_type": "Database",
                "runner": "docker",
                "volume_path": "/data/db"
            },
            "log": []
        },
        {
            "id": "redis-default",
            "name": "Redis",
            "description": "Redis 7 in-memory data store for caching & queues",
            "status": "Stopped",
            "config": {
                "port": 6379,
                "executable_path": "redis:7-alpine",
                "arguments": "",
                "service_type": "Database",
                "runner": "docker",
                "volume_path": "/data"
            },
            "log": []
        },

        // ===== Admin Panel =====
        {
            "id": "pma-default",
            "name": "phpMyAdmin",
            "description": "Web UI for MySQL & MariaDB management",
            "status": "Stopped",
            "config": {
                "port": 8080,
                "container_port": 80,
                "executable_path": "phpmyadmin:latest",
                "arguments": "-e PMA_HOST=host.docker.internal",
                "service_type": "Admin Panel",
                "runner": "docker"
            },
            "log": []
        },
        {
            "id": "adminer-default",
            "name": "Adminer",
            "description": "Lightweight DB admin for MySQL, PostgreSQL & more",
            "status": "Stopped",
            "config": {
                "port": 8081,
                "container_port": 8080,
                "executable_path": "adminer:latest",
                "arguments": "",
                "service_type": "Admin Panel",
                "runner": "docker"
            },
            "log": []
        },
        
        // ===== Mailer =====
        {
            "id": "mailpit-default",
            "name": "Mailpit",
            "description": "Modern MailHog replacement for SMTP testing",
            "status": "Stopped",
            "config": {
                "port": 8025,
                "container_port": 8025,
                "executable_path": "axllent/mailpit",
                "arguments": "-p 1025:1025",
                "service_type": "Mailer",
                "runner": "docker",
                "volume_path": "/data"
            },
            "log": []
        },

        // ===== Storage =====
        {
            "id": "minio-default",
            "name": "MinIO",
            "description": "High Performance Object Storage (S3 Compatible)",
            "status": "Stopped",
            "config": {
                "port": 9001,
                "container_port": 9001,
                "executable_path": "minio/minio",
                "arguments": "server /data --console-address \":9001\" -p 9000:9000",
                "service_type": "Storage",
                "runner": "docker",
                "volume_path": "/data"
            },
            "log": []
        },
        {
            "id": "pureftpd-default",
            "name": "Pure-FTPd",
            "description": "Secure and lightweight FTP server",
            "status": "Stopped",
            "config": {
                "port": 21,
                "container_port": 21,
                "executable_path": "stilliard/pure-ftpd",
                "arguments": "",
                "service_type": "Storage",
                "runner": "docker",
                "volume_path": "/home/ftpusers"
            },
            "log": []
        }
    ]);

    let _ = fs::write(&services_file, serde_json::to_string_pretty(&services).unwrap());
}
