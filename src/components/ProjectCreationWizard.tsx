import { useState, useEffect, useRef } from "react";
import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import { X, Check, Loader2, Terminal, Code2, Copy, Import } from "lucide-react";
import { open } from "@tauri-apps/plugin-dialog";
import { CustomSelect } from "./CustomSelect";
import { 
  SiNodedotjs, SiPython, SiPhp, SiGo, SiBun,
  SiNextdotjs, SiReact, SiVuedotjs, SiSvelte, SiNuxt, SiNestjs, SiExpress,
  SiFastapi, SiDjango, SiFlask,
  SiLaravel, SiSymfony, SiHtml5
} from "react-icons/si";

interface WizardProps {
  onClose: () => void;
  onCreated: () => void;
}

type Framework = {
  id: string;
  name: string;
  icon: any;
  createImage: string;
  createCmd: string;
  devImage: string;
  devCmd: string;
  defaultPort: number;
};

type Ecosystem = {
  id: string;
  name: string;
  icon: any;
  frameworks: Framework[];
};

const ECOSYSTEMS: Ecosystem[] = [
  {
    id: "node",
    name: "Node.js / TS",
    icon: SiNodedotjs,
    frameworks: [
      { id: "nextjs", name: "Next.js", icon: SiNextdotjs, createImage: "node:lts-alpine", createCmd: "npx -y create-next-app@latest {name} --ts --tailwind --eslint --app --src-dir --import-alias '@/*'", devImage: "node:lts-alpine", devCmd: "npm install && npm run dev -- -p {port} -H 0.0.0.0", defaultPort: 3000 },
      { id: "react", name: "React (Vite)", icon: SiReact, createImage: "node:lts-alpine", createCmd: "npx -y create-vite {name} --template react", devImage: "node:lts-alpine", devCmd: "npm install && npm run dev -- --host 0.0.0.0 --port {port}", defaultPort: 5173 },
      { id: "vue", name: "Vue (Vite)", icon: SiVuedotjs, createImage: "node:lts-alpine", createCmd: "npx -y create-vite {name} --template vue", devImage: "node:lts-alpine", devCmd: "npm install && npm run dev -- --host 0.0.0.0 --port {port}", defaultPort: 5173 },
      { id: "svelte", name: "SvelteKit", icon: SiSvelte, createImage: "node:lts-alpine", createCmd: "npx -y sv create {name} --template minimal --types ts --no-add-ons", devImage: "node:lts-alpine", devCmd: "npm install && npm run dev -- --host 0.0.0.0 --port {port}", defaultPort: 5173 },
      { id: "nuxt", name: "Nuxt", icon: SiNuxt, createImage: "node:lts-alpine", createCmd: "npx -y nuxi@latest init {name}", devImage: "node:lts-alpine", devCmd: "npm install && npm run dev -- -o -H 0.0.0.0 -p {port}", defaultPort: 3000 },
      { id: "nestjs", name: "NestJS", icon: SiNestjs, createImage: "node:lts-alpine", createCmd: "npx -y @nestjs/cli new {name} --package-manager npm", devImage: "node:lts-alpine", devCmd: "npm run start:dev", defaultPort: 3000 },
      { id: "native-node", name: "Native Node.js", icon: SiNodedotjs, createImage: "node:lts-alpine", createCmd: "mkdir -p {name} && cd {name} && npm init -y && echo 'const http = require(\"http\");\nhttp.createServer((req, res) => res.end(\"Hello from Native Node!\")).listen(parseInt(process.env.PORT || 3000));\nconsole.log(\"Server running\");' > index.js", devImage: "node:lts-alpine", devCmd: "PORT={port} node index.js", defaultPort: 3000 },
      { id: "express", name: "Express", icon: SiExpress, createImage: "node:lts-alpine", createCmd: "npx -y express-generator {name}", devImage: "node:lts-alpine", devCmd: "npm install && npm run start", defaultPort: 3000 },
    ]
  },
  {
    id: "python",
    name: "Python",
    icon: SiPython,
    frameworks: [
      { id: "fastapi", name: "FastAPI", icon: SiFastapi, createImage: "python:3.12-alpine", createCmd: "mkdir -p {name} && cd {name} && echo 'from fastapi import FastAPI\n\napp = FastAPI()\n\n@app.get(\"/\")\ndef read_root():\n    return {\"Hello\": \"World\"}' > main.py && echo 'fastapi[standard]' > requirements.txt", devImage: "python:3.12-alpine", devCmd: "pip install -r requirements.txt && fastapi dev main.py --host 0.0.0.0 --port {port}", defaultPort: 8000 },
      { id: "django", name: "Django", icon: SiDjango, createImage: "python:3.12-alpine", createCmd: "pip install django && django-admin startproject {name}", devImage: "python:3.12-alpine", devCmd: "pip install django && python manage.py runserver 0.0.0.0:{port}", defaultPort: 8000 },
      { id: "native-python", name: "Native Python", icon: SiPython, createImage: "python:3.12-alpine", createCmd: "mkdir -p {name} && cd {name} && echo 'print(\"Starting Native Python server...\")' > main.py", devImage: "python:3.12-alpine", devCmd: "python -m http.server {port}", defaultPort: 8000 },
      { id: "flask", name: "Flask", icon: SiFlask, createImage: "python:3.12-alpine", createCmd: "mkdir -p {name} && cd {name} && echo 'from flask import Flask\napp = Flask(__name__)\n\n@app.route(\"/\")\ndef hello():\n    return \"Hello World!\"' > app.py && echo 'flask' > requirements.txt", devImage: "python:3.12-alpine", devCmd: "pip install -r requirements.txt && flask run --host=0.0.0.0 --port={port}", defaultPort: 5000 },
    ]
  },
  {
    id: "php",
    name: "PHP",
    icon: SiPhp,
    frameworks: [
      { id: "laravel", name: "Laravel", icon: SiLaravel, createImage: "composer:latest", createCmd: "composer create-project laravel/laravel {name}", devImage: "composer:latest", devCmd: "php artisan serve --host=0.0.0.0 --port={port}", defaultPort: 8000 },
      { id: "symfony", name: "Symfony", icon: SiSymfony, createImage: "composer:latest", createCmd: "composer create-project symfony/skeleton {name}", devImage: "composer:latest", devCmd: "php -S 0.0.0.0:{port} -t public", defaultPort: 8000 },
      { id: "native-php", name: "Native PHP", icon: SiPhp, createImage: "alpine:latest", createCmd: "mkdir -p {name} && cd {name} && echo '<?php\\n\\necho \"Hello from Native PHP!\";\\n' > index.php", devImage: "alpine:latest", devCmd: "echo 'Native PHP project is served automatically by Lumine Nginx Proxy.' && sleep infinity", defaultPort: 80 },
    ]
  },
  {
    id: "go",
    name: "Go",
    icon: SiGo,
    frameworks: [
      { id: "gin", name: "Gin", icon: SiGo, createImage: "golang:1.22-alpine", createCmd: "mkdir -p {name} && cd {name} && go mod init {name} && go get -u github.com/gin-gonic/gin && echo 'package main\n\nimport \"github.com/gin-gonic/gin\"\n\nfunc main() {\n\tr := gin.Default()\n\tr.GET(\"/\", func(c *gin.Context) {\n\t\tc.JSON(200, gin.H{\"message\": \"hello world\"})\n\t})\n\tr.Run(\":8080\")\n}' > main.go", devImage: "golang:1.22-alpine", devCmd: "go run main.go", defaultPort: 8080 },
      { id: "native-go", name: "Native Go", icon: SiGo, createImage: "golang:1.22-alpine", createCmd: "mkdir -p {name} && cd {name} && go mod init {name} && echo 'package main\nimport (\"fmt\"; \"net/http\")\nfunc main() {\n\thttp.HandleFunc(\"/\", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, \"Hello Native Go\") })\n\tfmt.Println(\"Running...\")\n\thttp.ListenAndServe(\":8080\", nil)\n}' > main.go", devImage: "golang:1.22-alpine", devCmd: "go run main.go", defaultPort: 8080 },
      { id: "fiber", name: "Fiber", icon: SiGo, createImage: "golang:1.22-alpine", createCmd: "mkdir -p {name} && cd {name} && go mod init {name} && go get -u github.com/gofiber/fiber/v2 && echo 'package main\n\nimport \"github.com/gofiber/fiber/v2\"\n\nfunc main() {\n\tapp := fiber.New()\n\tapp.Get(\"/\", func(c *fiber.Ctx) error {\n\t\treturn c.SendString(\"Hello, World 👋!\")\n\t})\n\tapp.Listen(\":3000\")\n}' > main.go", devImage: "golang:1.22-alpine", devCmd: "go run main.go", defaultPort: 3000 },
    ]
  },
  {
    id: "bun",
    name: "Bun",
    icon: SiBun,
    frameworks: [
      { id: "native-bun", name: "Native Bun", icon: SiBun, createImage: "oven/bun:latest", createCmd: "mkdir -p {name} && cd {name} && bun init -y && echo 'Bun.serve({ port: process.env.PORT || 3000, fetch(req) { return new Response(\"Hello Native Bun\"); } }); console.log(\"Running\");' > index.ts", devImage: "oven/bun:latest", devCmd: "PORT={port} bun run --watch index.ts", defaultPort: 3000 },
      { id: "elysia", name: "Elysia", icon: SiBun, createImage: "oven/bun:latest", createCmd: "bun create elysia {name}", devImage: "oven/bun:latest", devCmd: "bun install && bun run --watch src/index.ts", defaultPort: 3000 },
    ]
  },
  {
    id: "html",
    name: "Native HTML/JS",
    icon: SiHtml5,
    frameworks: [
      { id: "native-html", name: "Static HTML", icon: SiHtml5, createImage: "alpine:latest", createCmd: "mkdir -p {name} && cd {name} && echo '<!DOCTYPE html>\\n<html>\\n<head>\\n<title>Native App</title>\\n</head>\\n<body>\\n<h1>Hello World!</h1>\\n</body>\\n</html>' > index.html", devImage: "alpine:latest", devCmd: "echo 'Static HTML is served automatically by Lumine Nginx Proxy.' && sleep infinity", defaultPort: 80 },
    ]
  }
];

export function ProjectCreationWizard({ onClose, onCreated }: WizardProps) {
  const [mode, setMode] = useState<"scaffold" | "import">("scaffold");
  const [step, setStep] = useState<"ecosystem" | "framework" | "details" | "install">("ecosystem");
  const [selectedEcosystem, setSelectedEcosystem] = useState(ECOSYSTEMS[0]);
  const [selectedFramework, setSelectedFramework] = useState(ECOSYSTEMS[0].frameworks[0]);
  
  const [projectName, setProjectName] = useState("");
  const [port, setPort] = useState(8000);
  
  // Import mode specific state
  const [importPath, setImportPath] = useState("");
  const [importFramework, setImportFramework] = useState("Custom");
  const [importCmd, setImportCmd] = useState("npm run dev");
  const [importDockerImage, setImportDockerImage] = useState("node:lts-alpine");

  const [logs, setLogs] = useState<string[]>([]);
  const [installStatus, setInstallStatus] = useState<"installing" | "done" | "error">("installing");
  const logEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logEndRef.current) {
      logEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [logs]);

  const handleImport = async () => {
    try {
      const workspacePosix = importPath.replace(/\\/g, '/');
      const devCmdReplaced = importCmd.replace("{port}", port.toString());
      const dockerStartArgs = ["docker", "run", "--rm", "--name", `lumine_proj_imported_${Date.now()}`, "-v", `${workspacePosix}:/app`, "-w", "/app", "-p", `${port}:${port}`, importDockerImage, "sh", "-c", devCmdReplaced];
      const dockerStart = JSON.stringify(dockerStartArgs);

      await invoke("add_project", {
        project: {
          id: `imported-${projectName}-${Date.now()}`,
          name: projectName,
          framework: importFramework,
          path: importPath,
          start_command: dockerStart,
          port: port,
          status: "Stopped",
          log: []
        }
      });

      try {
        const domainName = `${projectName}.test`;
        await invoke("add_host", { name: domainName, php: null, comment: projectName, isEnable: true, isIpv6: false });
        await invoke("generate_cert", { domain: domainName });
      } catch (e) { console.warn(e); }

      onCreated();
      onClose();
    } catch (e) {
      alert("Failed to import project: " + e);
    }
  };

  const startInstallation = async () => {
    setStep("install");
    setLogs([]);
    setInstallStatus("installing");

    const unlistenLog = await listen<string>("project-creation-log", (e) => {
      const cleanLine = e.payload.replace(/\x1B\[[0-9;]*[a-zA-Z]/g, "");
      setLogs((prev) => [...prev, cleanLine]);
    });
    const unlistenDone = await listen<string>("project-creation-done", () => {
      setInstallStatus("done");
    });
    const unlistenError = await listen<string>("project-creation-error", (e) => {
      setLogs((prev) => [...prev, `ERROR: ${e.payload}`]);
      setInstallStatus("error");
    });

    try {
      const workspace = await invoke<string>("get_workspace_dir");
      const workspacePosix = workspace.replace(/\\/g, '/');

      const createCommandRaw = selectedFramework.createCmd.replace(/\{name\}/g, `/tmp/${projectName}`);
      const createCommandArgs = ["sh", "-c", `${createCommandRaw} && mkdir -p /app/${projectName} && cp -a /tmp/${projectName}/. /app/${projectName}/`];
      const dockerCreateArgs = ["docker", "run", "--rm", "-v", `${workspacePosix}:/app`, "-w", "/app", selectedFramework.createImage, ...createCommandArgs];
      const dockerCreate = JSON.stringify(dockerCreateArgs);

      await invoke("create_new_project", {
        name: projectName,
        framework: selectedFramework.name,
        command: dockerCreate,
      });

      const devCmdReplaced = selectedFramework.devCmd.replace(/\{port\}/g, port.toString());
      const dockerStartArgs = ["docker", "run", "--rm", "--name", `lumine_proj_${selectedFramework.id}_${Date.now()}`, "-v", `${workspacePosix}/${projectName}:/app`, "-w", "/app", "-p", `${port}:${port}`, selectedFramework.devImage, "sh", "-c", devCmdReplaced];
      const dockerStart = JSON.stringify(dockerStartArgs);

      await invoke("add_project", {
        project: {
          id: `${selectedFramework.id}-${projectName}-${Date.now()}`,
          name: projectName,
          framework: selectedFramework.name,
          path: `${workspace}\\${projectName}`,
          start_command: dockerStart,
          port: port,
          status: "Stopped",
          log: []
        }
      });

      try {
        const domainName = `${projectName}.test`;
        await invoke("add_host", { name: domainName, php: null, comment: projectName, isEnable: true, isIpv6: false });
        await invoke("generate_cert", { domain: domainName });
      } catch (e) { console.warn(e); }

      setInstallStatus("done");
      onCreated();
    } catch (e) {
      setInstallStatus("error");
      setLogs((prev) => [...prev, `INVOKE ERROR: ${e}`]);
    }

    unlistenLog();
    unlistenDone();
    unlistenError();
  };

  return (
    <div 
      className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4 animate-fade-in"
      onMouseDown={step !== "install" ? onClose : undefined}
    >
      <div 
        className="bg-[#1e1e2e] border border-[#2a2a35] rounded-2xl shadow-2xl max-w-2xl w-full flex flex-col overflow-hidden max-h-[90vh] animate-modal"
        onMouseDown={e => e.stopPropagation()}
      >
        
        <div className="flex flex-col border-b border-[#2a2a35] shrink-0">
          <div className="flex items-center justify-between p-6 pb-4">
            <h2 className="text-xl font-bold text-white flex items-center gap-2">
              <Code2 className="w-5 h-5 text-blue-400" />
              {mode === "scaffold" ? "Create New Project" : "Import Project"}
            </h2>
            {step !== "install" && (
              <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
                <X className="w-5 h-5" />
              </button>
            )}
          </div>
          
          {step !== "install" && (
            <div className="flex gap-6 px-6">
              <button
                onClick={() => { setMode("scaffold"); setStep("ecosystem"); }}
                className={`pb-3 text-sm font-medium border-b-2 transition-colors ${mode === "scaffold" ? "border-blue-500 text-blue-400" : "border-transparent text-gray-400 hover:text-gray-300"}`}
              >
                Scaffold New
              </button>
              <button
                onClick={() => setMode("import")}
                className={`pb-3 text-sm font-medium border-b-2 transition-colors ${mode === "import" ? "border-blue-500 text-blue-400" : "border-transparent text-gray-400 hover:text-gray-300"}`}
              >
                Import Existing
              </button>
            </div>
          )}
        </div>

        <div className="p-6 overflow-y-auto custom-scrollbar flex-1">
          {mode === "scaffold" ? (
            <>
              {step === "ecosystem" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-300 mb-4">1. Select Ecosystem</h3>
                  <div className="grid grid-cols-2 gap-4">
                    {ECOSYSTEMS.map((eco) => {
                      const Icon = eco.icon;
                      return (
                        <div
                          key={eco.id}
                          onClick={() => {
                            setSelectedEcosystem(eco);
                            setStep("framework");
                          }}
                          className="p-4 rounded-xl border bg-[#242533] border-[#2a2a35] text-gray-400 hover:border-blue-500/50 hover:bg-blue-500/5 cursor-pointer transition-all flex items-center gap-3"
                        >
                          <div className="p-2 bg-black/20 rounded-lg">
                            <Icon className="w-5 h-5" />
                          </div>
                          <div>
                            <div className="font-bold text-gray-200">{eco.name}</div>
                            <div className="text-xs opacity-70">{eco.frameworks.length} Frameworks</div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {step === "framework" && (
                <div>
                  <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center justify-between">
                    <span>2. Select Framework for {selectedEcosystem.name}</span>
                    <button onClick={() => setStep("ecosystem")} className="text-xs text-gray-400 hover:text-white">Change</button>
                  </h3>
                  <div className="grid grid-cols-2 gap-4">
                    {selectedEcosystem.frameworks.map((fw) => {
                      const FwIcon = fw.icon;
                      return (
                        <div
                          key={fw.id}
                          onClick={() => {
                            setSelectedFramework(fw);
                            setPort(fw.defaultPort);
                          }}
                          className={`p-4 rounded-xl border cursor-pointer transition-all flex items-center gap-3 ${
                            selectedFramework.id === fw.id
                              ? "bg-blue-500/10 border-blue-500/50 text-white"
                              : "bg-[#242533] border-[#2a2a35] text-gray-400 hover:border-gray-500 hover:text-gray-300"
                          }`}
                        >
                          <div className={`p-2 rounded-lg ${selectedFramework.id === fw.id ? "bg-blue-500/20 text-blue-400" : "bg-black/20"}`}>
                            <FwIcon className="w-6 h-6" />
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="font-bold mb-0.5">{fw.name}</div>
                            <div className="text-[10px] opacity-70 font-mono truncate" title={fw.createCmd}>{fw.createCmd}</div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                  <div className="mt-8 flex justify-between">
                    <button
                      onClick={() => setStep("ecosystem")}
                      className="text-gray-400 hover:text-white px-4 py-2 rounded-lg font-medium transition-colors"
                    >
                      Back
                    </button>
                    <button
                      onClick={() => setStep("details")}
                      className="bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg font-medium transition-colors"
                    >
                      Next Step
                    </button>
                  </div>
                </div>
              )}

              {step === "details" && (
                <div className="space-y-6">
                  <h3 className="text-sm font-semibold text-gray-300">3. Project Details</h3>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1.5">Project Name (Folder Name)</label>
                    <input
                      type="text"
                      value={projectName}
                      onChange={(e) => setProjectName(e.target.value.replace(/[^a-zA-Z0-9-_]/g, "").toLowerCase())}
                      placeholder="my-awesome-app"
                      className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-blue-500 transition-colors"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1.5">Dev Server Port</label>
                    <input
                      type="number"
                      value={port}
                      onChange={(e) => setPort(parseInt(e.target.value))}
                      className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-4 py-2.5 text-white focus:outline-none focus:border-blue-500 transition-colors"
                    />
                  </div>

                  <div className="bg-blue-500/10 border border-blue-500/20 p-4 rounded-lg flex items-start gap-3">
                    <Terminal className="w-5 h-5 text-blue-500 shrink-0 mt-0.5" />
                    <div className="text-sm text-blue-200/80">
                      This will run <code className="bg-black/30 px-1.5 py-0.5 rounded text-blue-100">{selectedFramework.createCmd.replace(/\{name\}/g, projectName || "project-name")}</code>
                      <br />Command will be safely executed inside an isolated Docker container.
                    </div>
                  </div>

                  <div className="mt-8 flex justify-between">
                    <button
                      onClick={() => setStep("framework")}
                      className="text-gray-400 hover:text-white px-4 py-2 rounded-lg font-medium transition-colors"
                    >
                      Back
                    </button>
                    <button
                      onClick={startInstallation}
                      disabled={!projectName.trim()}
                      className="bg-blue-500 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed text-white px-6 py-2 rounded-lg font-medium transition-colors"
                    >
                      Create Project
                    </button>
                  </div>
                </div>
              )}

              {step === "install" && (
                <div className="h-full flex flex-col">
                  <div className="flex items-center justify-between mb-4">
                    <h3 className="text-sm font-semibold text-gray-300 flex items-center gap-2">
                      {installStatus === "installing" && <Loader2 className="w-4 h-4 animate-spin text-blue-400" />}
                      {installStatus === "done" && <Check className="w-4 h-4 text-green-400" />}
                      {installStatus === "error" && <X className="w-4 h-4 text-red-400" />}
                      Installing {selectedFramework.name}
                    </h3>
                    <button
                      onClick={() => navigator.clipboard.writeText(logs.join('\n'))}
                      className="flex items-center gap-1.5 text-xs bg-[#181825] hover:bg-[#2a2a35] text-gray-400 hover:text-white px-2.5 py-1.5 rounded transition-colors"
                    >
                      <Copy className="w-3.5 h-3.5" />
                      Copy Logs
                    </button>
                  </div>
                  
                  <div className="h-[300px] bg-[#0a0a0f] rounded-lg border border-[#2a2a35] p-4 text-sm font-mono overflow-y-auto custom-scrollbar select-text">
                    {logs.length === 0 ? (
                      <span className="text-gray-400 animate-pulse">Starting installation...</span>
                    ) : (
                      <>
                        {logs.map((line, i) => (
                          <div key={i} className="text-gray-200 leading-relaxed break-all">
                            {line}
                          </div>
                        ))}
                        {installStatus === "done" && (
                          <div className="mt-4 text-green-400 font-bold border-t border-[#2a2a35] pt-4">
                            ✅ Installation completed successfully! You may now close this window.
                          </div>
                        )}
                      </>
                    )}
                    <div ref={logEndRef} />
                  </div>

                  {installStatus !== "installing" && (
                    <div className="mt-6 flex justify-end">
                      <button
                        onClick={onClose}
                        className={`${installStatus === "done" ? "bg-green-600 hover:bg-green-500" : "bg-gray-700 hover:bg-gray-600"} text-white px-6 py-2 rounded-lg font-medium transition-colors`}
                      >
                        {installStatus === "done" ? "Finish & Close" : "Close"}
                      </button>
                    </div>
                  )}
                </div>
              )}
            </>
          ) : (
            <div className="space-y-5">
              <div className="bg-yellow-500/10 border border-yellow-500/20 p-4 rounded-lg flex items-start gap-3">
                <Import className="w-5 h-5 text-yellow-500 shrink-0 mt-0.5" />
                <div className="text-sm text-yellow-200/80 leading-relaxed">
                  Import an existing project folder into Lumine. You must provide the absolute path and the Docker execution command to run your development server.
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1">Project Name</label>
                <input
                  type="text"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value.replace(/[^a-zA-Z0-9-_]/g, "").toLowerCase())}
                  placeholder="my-imported-app"
                  className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-gray-400 mb-1">Absolute Directory Path</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={importPath}
                    onChange={(e) => setImportPath(e.target.value)}
                    placeholder="C:\Lumine\www\my-app"
                    className="flex-1 w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500 font-mono"
                  />
                  <button
                    onClick={async () => {
                      const selected = await open({ directory: true, multiple: false });
                      if (selected && typeof selected === 'string') {
                        setImportPath(selected);
                      }
                    }}
                    className="px-4 py-2 bg-[#2a2a35] hover:bg-[#32323d] text-sm font-medium rounded-lg text-white transition-colors"
                  >
                    Browse
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-400 mb-1">Framework / Type</label>
                  <CustomSelect
                    options={[
                      { label: "Custom", value: "Custom" },
                      { label: "Next.js", value: "Next.js" },
                      { label: "React", value: "React" },
                      { label: "Vue", value: "Vue" },
                      { label: "Nuxt", value: "Nuxt" },
                      { label: "SvelteKit", value: "SvelteKit" },
                      { label: "Laravel", value: "Laravel" },
                      { label: "NestJS", value: "NestJS" },
                      { label: "Express", value: "Express" },
                    ]}
                    value={importFramework}
                    onChange={(v) => setImportFramework(v)}
                    allowCustom={true}
                    searchable={true}
                    placeholder="Select or type..."
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-400 mb-1">Local Port</label>
                  <input
                    type="number"
                    value={port}
                    onChange={(e) => setPort(parseInt(e.target.value))}
                    className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-gray-400 mb-1">Docker Image</label>
                  <CustomSelect
                    options={[
                      { label: "node:lts-alpine", value: "node:lts-alpine" },
                      { label: "oven/bun:latest", value: "oven/bun:latest" },
                      { label: "denoland/deno:latest", value: "denoland/deno:latest" },
                      { label: "php:8.2-cli", value: "php:8.2-cli" },
                      { label: "python:3.11-alpine", value: "python:3.11-alpine" },
                      { label: "golang:alpine", value: "golang:alpine" },
                    ]}
                    value={importDockerImage}
                    onChange={(v) => setImportDockerImage(v)}
                    allowCustom={true}
                    searchable={true}
                    placeholder="e.g. node:lts-alpine"
                  />
                </div>
                <div>
                  <label className="block text-xs font-medium text-gray-400 mb-1">Dev Command (Use {"{port}"} for port)</label>
                  <input
                    type="text"
                    value={importCmd}
                    onChange={(e) => setImportCmd(e.target.value)}
                    placeholder="npm run dev -- --port {port}"
                    className="w-full bg-[#181825] border border-[#2a2a35] rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500 font-mono"
                  />
                </div>
              </div>

              <div className="mt-6 flex justify-end pt-4 border-t border-[#2a2a35]">
                <button
                  onClick={handleImport}
                  disabled={!projectName.trim() || !importPath.trim()}
                  className="bg-blue-500 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed text-white px-6 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
                >
                  <Import className="w-4 h-4" />
                  Import Project
                </button>
              </div>
            </div>
          )}
        </div>

      </div>
    </div>
  );
}
