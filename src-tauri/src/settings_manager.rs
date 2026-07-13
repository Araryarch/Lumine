use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct CustomTerminalImage {
    pub name: String,
    pub image: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct AppSettings {
    pub start_on_boot: bool,
    pub auto_start_services: bool,
    pub minimize_to_tray: bool,
    pub check_updates: bool,
    pub default_php: String,
    pub default_node: String,
    pub document_root: String,
    pub terminal_emulator: String,
    pub code_editor: String,
    pub file_explorer: String,
    #[serde(default)]
    pub custom_terminal_images: Vec<CustomTerminalImage>,
}

impl Default for AppSettings {
    fn default() -> Self {
        Self {
            start_on_boot: false,
            auto_start_services: false,
            minimize_to_tray: true,
            check_updates: true,
            default_php: "8.3".to_string(),
            default_node: "20.11".to_string(),
            document_root: "C:\\Lumine\\www".to_string(),
            terminal_emulator: "cmd".to_string(),
            code_editor: "code".to_string(),
            file_explorer: "explorer".to_string(),
            custom_terminal_images: vec![],
        }
    }
}

pub struct SettingsManager {
    data_file: PathBuf,
}

impl SettingsManager {
    pub fn new(app_data_dir: PathBuf) -> Self {
        if !app_data_dir.exists() {
            let _ = fs::create_dir_all(&app_data_dir);
        }
        let data_file = app_data_dir.join("settings.json");
        
        if !data_file.exists() {
            let default_settings = AppSettings::default();
            let _ = fs::write(&data_file, serde_json::to_string_pretty(&default_settings).unwrap());
        }

        Self { data_file }
    }

    pub fn get_settings(&self) -> Result<AppSettings, String> {
        let content = fs::read_to_string(&self.data_file)
            .map_err(|e| format!("Failed to read settings: {}", e))?;
        let settings: AppSettings = serde_json::from_str(&content)
            .unwrap_or_default();
        Ok(settings)
    }

    pub fn save_settings(&self, settings: AppSettings) -> Result<(), String> {
        let json = serde_json::to_string_pretty(&settings)
            .map_err(|e| format!("Failed to serialize settings: {}", e))?;
        fs::write(&self.data_file, json)
            .map_err(|e| format!("Failed to save settings: {}", e))?;
        Ok(())
    }
}
