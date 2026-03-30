# project_template_go

A Go project template. Run one command inside your cloned repo to scaffold the full project structure — all module paths, project name, version, and bundle ID are rewritten automatically.

## Quick Start

```bash
# 1. Create repo on GitHub, then clone and enter it
git clone git@github.com:myorg/my_service.git && cd my_service

# 2. Scaffold
wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- github.com/myorg/my_service

# 3. Commit
git add . && git commit -m "chore: init from project_template_go"
```

The script writes directly into the current directory. No subdirectory is created.

## What it does

1. Checks `gonew` is installed (auto-installs if missing)
2. Clones this template via `gonew`, rewrites all Go import paths to the new module path
3. Patches `config/config.go` constants:
   - `ProjectName` → last segment of module path
   - `ProjectVersion` → reset to `v0.0.1`
   - `ProjectBundleID` → `com.<project_name>.<project_name>`
4. Replaces template references in `.md / .yml / .yaml / .xml / .sh / .json / .txt` files
5. Deletes itself
