# project_template_go

A Go project template. One command scaffolds a new project — all module paths, project name, version, and bundle ID are rewritten automatically.

## Quick Start

```bash
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- <module_path> [target_dir]
```

### Mode 1 — already inside cloned repo

```bash
git clone git@github.com:myorg/my_service.git && cd my_service
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- github.com/myorg/my_service
git add . && git commit -m "chore: init from project_template_go"
```

### Mode 2 — create a new directory anywhere

```bash
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- github.com/myorg/my_service ./my_service
# or any other module path format:
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- mycompany.com/backend ./backend
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- gitlab.com/team/api ./api
```

> `module_path` must be a valid Go module path (must contain a dot).
> `target_dir` defaults to the current directory.

## What it does

1. Checks `gonew` is installed (auto-installs if missing)
2. Clones this template via `gonew`, rewrites all Go import paths to the new module path
3. Patches `config/config.go` constants:
   - `ProjectName` → last segment of module path
   - `ProjectVersion` → reset to `v0.0.1`
   - `ProjectBundleID` → `com.<project_name>.<project_name>`
4. Replaces template references in `.md / .yml / .yaml / .xml / .sh / .json / .txt` files
5. Deletes itself
