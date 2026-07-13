use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;

#[cfg(target_os = "windows")]
const CREATE_NO_WINDOW: u32 = 0x08000000;

const MKCERT_VERSION: &str = "v1.4.4";

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct CertEntry {
    pub domain: String,
    pub cert_path: String,
    pub key_path: String,
    pub created_at: String,
}

pub struct MkCertManager {
    data_file: PathBuf,
    certs_dir: PathBuf,
    bin_dir: PathBuf,
}

impl MkCertManager {
    pub fn new(app_data_dir: PathBuf) -> Self {
        let certs_dir = app_data_dir.join("certs");
        if !certs_dir.exists() {
            let _ = fs::create_dir_all(&certs_dir);
        }
        let bin_dir = app_data_dir.join("bin");
        if !bin_dir.exists() {
            let _ = fs::create_dir_all(&bin_dir);
        }
        let data_file = app_data_dir.join("certs.json");
        if !data_file.exists() {
            let _ = fs::write(&data_file, "[]");
        }
        Self { data_file, certs_dir, bin_dir }
    }

    /// Get the path to the mkcert binary. Checks:
    /// 1. Our bundled copy in AppData/bin/
    /// 2. System PATH
    fn get_mkcert_path(&self) -> Option<String> {
        // Check our bundled binary first
        let bundled = self.get_bundled_path();
        if bundled.exists() {
            return Some(bundled.to_string_lossy().to_string());
        }

        // Fallback: check system PATH
        let mut cmd = crate::utils::create_command("mkcert");
        cmd.arg("--version");
        #[cfg(target_os = "windows")]
        cmd.creation_flags(CREATE_NO_WINDOW);
        if cmd.output().is_ok() {
            return Some("mkcert".to_string());
        }

        None
    }

    fn get_bundled_path(&self) -> PathBuf {
        #[cfg(target_os = "windows")]
        { self.bin_dir.join("mkcert.exe") }
        #[cfg(not(target_os = "windows"))]
        { self.bin_dir.join("mkcert") }
    }

    pub fn check_installed(&self) -> bool {
        self.get_mkcert_path().is_some()
    }

    /// Download mkcert binary from GitHub releases
    pub async fn download_mkcert(&self) -> Result<String, String> {
        let bundled = self.get_bundled_path();
        if bundled.exists() {
            return Ok("mkcert is already downloaded.".to_string());
        }

        #[cfg(target_os = "windows")]
        let asset_name = "mkcert-v1.4.4-windows-amd64.exe";
        #[cfg(all(target_os = "linux", target_arch = "x86_64"))]
        let asset_name = "mkcert-v1.4.4-linux-amd64";
        #[cfg(all(target_os = "macos", target_arch = "x86_64"))]
        let asset_name = "mkcert-v1.4.4-darwin-amd64";
        #[cfg(all(target_os = "macos", target_arch = "aarch64"))]
        let asset_name = "mkcert-v1.4.4-darwin-arm64";

        let url = format!(
            "https://github.com/FiloSottile/mkcert/releases/download/{}/{}",
            MKCERT_VERSION, asset_name
        );

        let response = reqwest::Client::new()
            .get(&url)
            .header("User-Agent", "Lumine")
            .send()
            .await
            .map_err(|e| format!("Failed to download mkcert: {}", e))?;

        if !response.status().is_success() {
            return Err(format!("Download failed with status: {}", response.status()));
        }

        let bytes = response.bytes()
            .await
            .map_err(|e| format!("Failed to read download: {}", e))?;

        fs::write(&bundled, &bytes)
            .map_err(|e| format!("Failed to save mkcert binary: {}", e))?;

        // Make executable on Unix
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let _ = fs::set_permissions(&bundled, fs::Permissions::from_mode(0o755));
        }

        Ok(format!("mkcert {} downloaded successfully!", MKCERT_VERSION))
    }

    pub fn install_root_ca(&self) -> Result<String, String> {
        let mkcert = self.get_mkcert_path()
            .ok_or("mkcert is not installed. Please download it first.")?;
        let mut cmd = crate::utils::create_command(&mkcert);
        cmd.arg("-install");
        #[cfg(target_os = "windows")]
        cmd.creation_flags(CREATE_NO_WINDOW);
        let output = cmd.output().map_err(|e| format!("Failed to run mkcert -install: {}", e))?;
        if output.status.success() {
            let msg = String::from_utf8_lossy(&output.stderr).to_string();
            Ok(if msg.is_empty() { "Root CA installed successfully".to_string() } else { msg })
        } else {
            Err(String::from_utf8_lossy(&output.stderr).to_string())
        }
    }

    pub fn generate_cert(&self, domain: &str) -> Result<CertEntry, String> {
        let mkcert = self.get_mkcert_path()
            .ok_or("mkcert is not installed. Please download it first.")?;

        // Validate domain
        if domain.is_empty() || domain.len() > 255 {
            return Err("Invalid domain length".to_string());
        }
        if !domain.chars().all(|c| c.is_ascii_alphanumeric() || c == '.' || c == '-' || c == '*') {
            return Err("Invalid domain characters".to_string());
        }

        // Check if cert already exists
        let existing = self.get_certs()?;
        if existing.iter().any(|c| c.domain == domain) {
            return Err(format!("Certificate for '{}' already exists", domain));
        }

        let safe_name = domain.replace('*', "_wildcard_").replace('.', "_");
        let cert_file = self.certs_dir.join(format!("{}.pem", safe_name));
        let key_file = self.certs_dir.join(format!("{}-key.pem", safe_name));

        let mut cmd = crate::utils::create_command(&mkcert);
        cmd.arg("-cert-file").arg(&cert_file)
           .arg("-key-file").arg(&key_file)
           .arg(domain);
        #[cfg(target_os = "windows")]
        cmd.creation_flags(CREATE_NO_WINDOW);

        let output = cmd.output().map_err(|e| format!("Failed to run mkcert: {}", e))?;
        if !output.status.success() {
            return Err(format!("mkcert failed: {}", String::from_utf8_lossy(&output.stderr)));
        }

        let entry = CertEntry {
            domain: domain.to_string(),
            cert_path: cert_file.to_string_lossy().to_string(),
            key_path: key_file.to_string_lossy().to_string(),
            created_at: chrono::Local::now().format("%Y-%m-%d %H:%M:%S").to_string(),
        };

        let mut certs = self.get_certs()?;
        certs.push(entry.clone());
        self.save_certs(&certs)?;

        Ok(entry)
    }

    pub fn get_certs(&self) -> Result<Vec<CertEntry>, String> {
        let content = fs::read_to_string(&self.data_file)
            .map_err(|e| format!("Failed to read certs data: {}", e))?;
        let certs: Vec<CertEntry> = serde_json::from_str(&content).unwrap_or_default();
        Ok(certs)
    }

    fn save_certs(&self, certs: &[CertEntry]) -> Result<(), String> {
        let json = serde_json::to_string_pretty(certs)
            .map_err(|e| format!("Failed to serialize certs: {}", e))?;
        fs::write(&self.data_file, json)
            .map_err(|e| format!("Failed to save certs data: {}", e))?;
        Ok(())
    }

    pub fn delete_cert(&self, domain: &str) -> Result<(), String> {
        let mut certs = self.get_certs()?;
        if let Some(pos) = certs.iter().position(|c| c.domain == domain) {
            let cert = certs.remove(pos);
            // Delete cert files from disk
            let _ = fs::remove_file(&cert.cert_path);
            let _ = fs::remove_file(&cert.key_path);
            self.save_certs(&certs)?;
            Ok(())
        } else {
            Err(format!("Certificate for '{}' not found", domain))
        }
    }
}
