use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::PathBuf;
use std::sync::{Arc, Mutex};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StackInfo {
    pub id: String,
    pub name: String,
    pub description: String,
    pub color: String,
    pub services: Vec<String>,
}

#[derive(Clone)]
pub struct StackManager {
    stacks: Arc<Mutex<HashMap<String, StackInfo>>>,
    storage_path: PathBuf,
}

impl StackManager {
    pub fn new(app_data_dir: PathBuf) -> Self {
        let storage_path = app_data_dir.join("stacks.json");
        let stacks = Self::load_stacks(&storage_path);

        Self {
            stacks: Arc::new(Mutex::new(stacks)),
            storage_path,
        }
    }

    fn load_stacks(path: &PathBuf) -> HashMap<String, StackInfo> {
        if path.exists() {
            if let Ok(data) = fs::read_to_string(path) {
                if let Ok(parsed) = serde_json::from_str(&data) {
                    return parsed;
                }
            }
        }
        let mut default_stacks = HashMap::new();
        
        // 1. XAMPP Migration
        default_stacks.insert(
            "xampp_migration".to_string(),
            StackInfo {
                id: "xampp_migration".to_string(),
                name: "XAMPP Migration".to_string(),
                description: "if u migrate from xampp".to_string(),
                color: "#3b82f6".to_string(), // Blue
                services: vec!["apache-default".to_string(), "mysql-default".to_string(), "pma-default".to_string()],
            }
        );

        // 2. MERN Stack
        default_stacks.insert(
            "mern_stack".to_string(),
            StackInfo {
                id: "mern_stack".to_string(),
                name: "MERN Stack".to_string(),
                description: "MongoDB, Express, React, Node".to_string(),
                color: "#10b981".to_string(), // Emerald
                services: vec!["mongo-default".to_string(), "node-default".to_string()],
            }
        );

        // 3. Laravel Local
        default_stacks.insert(
            "laravel_dev".to_string(),
            StackInfo {
                id: "laravel_dev".to_string(),
                name: "Laravel Local".to_string(),
                description: "Standard Laravel Dev Environment".to_string(),
                color: "#ef4444".to_string(), // Red
                services: vec!["mysql-default".to_string(), "redis-default".to_string(), "mailpit-default".to_string()],
            }
        );
        
        // 4. LEMP Stack
        default_stacks.insert(
            "lemp_stack".to_string(),
            StackInfo {
                id: "lemp_stack".to_string(),
                name: "LEMP Stack".to_string(),
                description: "Nginx, MySQL".to_string(),
                color: "#8b5cf6".to_string(), // Violet
                services: vec!["nginx-default".to_string(), "mysql-default".to_string()],
            }
        );

        // 5. LAMP Stack
        default_stacks.insert(
            "lamp_stack".to_string(),
            StackInfo {
                id: "lamp_stack".to_string(),
                name: "LAMP Stack".to_string(),
                description: "Apache, MySQL".to_string(),
                color: "#f59e0b".to_string(), // Amber
                services: vec!["apache-default".to_string(), "mysql-default".to_string()],
            }
        );

        // 6. Python Data Science
        default_stacks.insert(
            "python_data".to_string(),
            StackInfo {
                id: "python_data".to_string(),
                name: "Python Data".to_string(),
                description: "PostgreSQL & Redis for Python apps".to_string(),
                color: "#0ea5e9".to_string(), // Sky Blue
                services: vec!["postgres-default".to_string(), "redis-default".to_string(), "python-default".to_string()],
            }
        );

        // 7. Go Microservices
        default_stacks.insert(
            "go_microservices".to_string(),
            StackInfo {
                id: "go_microservices".to_string(),
                name: "Go Microservices".to_string(),
                description: "High performance Go stack".to_string(),
                color: "#06b6d4".to_string(), // Cyan
                services: vec!["postgres-default".to_string(), "redis-default".to_string(), "go-default".to_string()],
            }
        );

        // 8. Spring Boot Dev
        default_stacks.insert(
            "spring_boot".to_string(),
            StackInfo {
                id: "spring_boot".to_string(),
                name: "Spring Boot".to_string(),
                description: "Java Spring Boot environment".to_string(),
                color: "#84cc16".to_string(), // Lime
                services: vec!["mysql-default".to_string(), "redis-default".to_string(), "java-default".to_string()],
            }
        );

        // 9. Ruby on Rails
        default_stacks.insert(
            "rails_dev".to_string(),
            StackInfo {
                id: "rails_dev".to_string(),
                name: "Ruby on Rails".to_string(),
                description: "Standard Rails stack".to_string(),
                color: "#be123c".to_string(), // Rose
                services: vec!["postgres-default".to_string(), "redis-default".to_string(), "ruby-default".to_string()],
            }
        );

        // 10. Rust Web Backend
        default_stacks.insert(
            "rust_backend".to_string(),
            StackInfo {
                id: "rust_backend".to_string(),
                name: "Rust Backend".to_string(),
                description: "Blazing fast Rust stack".to_string(),
                color: "#f97316".to_string(), // Orange
                services: vec!["postgres-default".to_string(), "redis-default".to_string(), "rust-default".to_string()],
            }
        );

        default_stacks
    }

    fn save_stacks(&self, stacks: &HashMap<String, StackInfo>) {
        if let Ok(data) = serde_json::to_string_pretty(stacks) {
            let _ = fs::write(&self.storage_path, data);
        }
    }

    pub fn get_all(&self) -> Vec<StackInfo> {
        let guard = self.stacks.lock().unwrap();
        let mut list: Vec<StackInfo> = guard.values().cloned().collect();
        // Sort by name
        list.sort_by(|a, b| a.name.cmp(&b.name));
        list
    }

    pub fn add_stack(&self, info: StackInfo) -> Result<(), String> {
        let mut guard = self.stacks.lock().unwrap();
        if guard.contains_key(&info.id) {
            return Err("Stack with this ID already exists".to_string());
        }
        guard.insert(info.id.clone(), info);
        self.save_stacks(&guard);
        Ok(())
    }

    pub fn edit_stack(&self, id: &str, info: StackInfo) -> Result<(), String> {
        let mut guard = self.stacks.lock().unwrap();
        if !guard.contains_key(id) {
            return Err("Stack not found".to_string());
        }
        guard.insert(id.to_string(), info);
        self.save_stacks(&guard);
        Ok(())
    }

    pub fn delete_stack(&self, id: &str) -> Result<(), String> {
        let mut guard = self.stacks.lock().unwrap();
        if guard.remove(id).is_some() {
            self.save_stacks(&guard);
            Ok(())
        } else {
            Err("Stack not found".to_string())
        }
    }
}
