# project_template_go

A Go project template. Run one command to scaffold a new project — all module paths, project name, version, and bundle ID are rewritten automatically.

## Quick Start

```bash
# 先 cd 到你想放新项目的父目录，再运行
cd ~/your/workspace
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- <project_name> <module_path>
```

Examples:

```bash
cd ~/workspace
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- my_service github.com/myorg/my_service

# custom workspace root (pass env to bash, not wget)
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | PROJECT_WORKSPACE=/path/to/ws bash -s -- my_service github.com/myorg/my_service
```

New project lands at `./<project_name>/` relative to where you run the command.

> `module_path` is required and must be a valid Go module path (must contain a dot, e.g. `github.com/myorg/my_service`).

## What it does

1. Checks `gonew` is installed (auto-installs if missing)
2. Clones this template via `gonew`, rewrites all Go import paths to the new module path
3. Patches `config/config.go` constants:
   - `ProjectName` → new project name
   - `ProjectVersion` → reset to `v0.0.1`
   - `ProjectBundleID` → `com.<project_name>.<project_name>`
4. Replaces template references in `.md / .yml / .yaml / .xml / .sh / .json / .txt` files
5. Deletes itself

## After scaffolding

```bash
cd <project_name>
git init && git add . && git commit -m "chore: init from project_template_go"
```
