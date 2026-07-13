import type { RawServiceInfo, ServiceStatus, RawProjectInfo, Project } from "../types";

export function parseStatus(info: { status: any }): ServiceStatus {
  if (info.status === "Running") return "running";
  if (info.status === "Stopped") return "stopped";
  if (typeof info.status === "object" && "Error" in info.status) return "error";
  return "stopped";
}

export function toService(info: RawServiceInfo) {
  return {
    id: info.id,
    name: info.name,
    description: info.description,
    status: parseStatus(info),
    config: {
      port: info.config.port,
      executablePath: info.config.executable_path,
      arguments: info.config.arguments,
      serviceType: info.config.service_type,
      runner: (info.config.runner || "binary") as "binary" | "docker",
      containerPort: info.config.container_port,
      volumePath: info.config.volume_path,
      env: info.config.env,
      auto_start: info.config.auto_start,
    },
    log: info.log,
  };
}

export function toProject(info: RawProjectInfo): Project {
  return {
    id: info.id,
    name: info.name,
    framework: info.framework,
    path: info.path,
    startCommand: info.start_command,
    port: info.port,
    status: parseStatus(info),
    log: info.log,
  };
}
