# project_template_go

A Go project template. Run one command inside your project directory to scaffold the full structure — all module paths, project name, version, and bundle ID are rewritten automatically.

## Quick Start

**已有 `go.mod`**（module path 自动读取，无需传参）：

```bash
cd my_project
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash
```

**没有 `go.mod`**（手动传 module path）：

```bash
cd my_project
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- my_project
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- github.com/myorg/my_service
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- mycompany.com/backend
```

## What it does

1. Clones this template into a temp directory
2. Rewrites all module paths and project name references across every file
3. Patches `config/config.go` constants:
   - `ProjectName` → last segment of module path
   - `ProjectVersion` → reset to `v0.0.1`
   - `ProjectBundleID` → `com.<project_name>.<project_name>`
4. Copies everything into the current directory
5. Deletes itself

No extra tools required — only `git` and `bash`.
