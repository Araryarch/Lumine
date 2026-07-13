export type ServiceStatus = "running" | "stopped" | "error" | "starting" | "stopping"

export interface ServiceConfig {
  port: number
  executablePath: string
  arguments: string
  serviceType: string
  runner?: "binary" | "docker"
  containerPort?: number
  volumePath?: string
  env?: Record<string, string>
  auto_start?: boolean
}

export interface Service {
  id: string
  name: string
  description: string
  status: ServiceStatus
  config: ServiceConfig
  log: string[]
}

// Raw shape from Rust backend (snake_case)
// Serde externally-tagged enum: "Running" | "Stopped" | { "Error": "message" }
export interface RawServiceInfo {
  id: string
  name: string
  description: string
  status: "Running" | "Stopped" | { Error: string }
  config: { port: number; executable_path: string; arguments: string; service_type: string; runner?: string; container_port?: number; volume_path?: string; env?: Record<string, string>; auto_start?: boolean }
  log: string[]
}

export interface HostEntry {
  name: string
  php: string | null
  comment: string
  is_enable: boolean
  is_ipv6: boolean
}

export interface CertEntry {
  domain: string
  cert_path: string
  key_path: string
  created_at: string
}

export type ProjectStatus = "running" | "stopped" | "error" | "starting" | "stopping"

export interface Project {
  id: string
  name: string
  framework: string
  path: string
  startCommand: string
  port: number
  status: ProjectStatus
  log: string[]
  env?: Record<string, string>
}

export interface RawProjectInfo {
  id: string
  name: string
  framework: string
  path: string
  start_command: string
  port: number
  status: "Running" | "Stopped" | { Error: string }
  log: string[]
  env?: Record<string, string>
}

export interface AppSettings {
  startOnBoot: boolean;
  autoStartServices: boolean;
  minimizeToTray: boolean;
  checkUpdates: boolean;
  defaultPhp: string;
  defaultNode: string;
  documentRoot: string;
  terminalEmulator: string;
  codeEditor: string;
  fileExplorer: string;
}

export interface Stack {
  id: string
  name: string
  description: string
  color: string
  services: string[]
}
