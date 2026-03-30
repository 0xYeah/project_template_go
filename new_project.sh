#!/usr/bin/env bash
# Usage:
#   bash new_project.sh <project_name> [module_path]
#   wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- <project_name> [module_path]
#
# Examples:
#   bash new_project.sh my_service
#   bash new_project.sh my_service github.com/myorg/my_service

set -euo pipefail

TEMPLATE_MODULE="github.com/0xYeah/project_template_go"
TEMPLATE_NAME="project_template_go"

# Workspace root — defaults to current directory
# Override: wget -qO- ... | PROJECT_WORKSPACE=/path/to/ws bash -s -- my_service
WORKSPACE="${PROJECT_WORKSPACE:-$(pwd)}"

usage() {
    echo "Usage: new <project_name> [module_path]"
    echo ""
    echo "  project_name   directory name for the new project"
    echo "  module_path    Go module path (default: <project_name>)"
    echo ""
    echo "Environment:"
    echo "  PROJECT_WORKSPACE   root workspace dir (default: ${WORKSPACE})"
    exit 1
}

[[ $# -lt 1 ]] && usage

PROJECT_NAME="$1"
NEW_MODULE="${2:-${PROJECT_NAME}}"
TARGET_DIR="${WORKSPACE}/${PROJECT_NAME}"

echo "Template : ${TEMPLATE_MODULE}"
echo "New      : ${NEW_MODULE}"
echo "Target   : ${TARGET_DIR}"
echo ""

# ── 0. Guard: target must not exist ──────────────────────────────────────────
if [[ -e "${TARGET_DIR}" ]]; then
    echo "Error: target already exists: ${TARGET_DIR}"
    echo "Remove it first or choose a different project name."
    exit 1
fi

# ── 1. Install gonew if missing ──────────────────────────────────────────────
if ! command -v gonew &>/dev/null; then
    echo "[1/4] Installing gonew..."
    go install golang.org/x/tools/cmd/gonew@latest
else
    echo "[1/4] gonew: $(which gonew)"
fi

# ── 2. Clone + rename all Go files via gonew ─────────────────────────────────
echo "[2/4] Running gonew..."
gonew "${TEMPLATE_MODULE}" "${NEW_MODULE}" "${TARGET_DIR}"

cd "${TARGET_DIR}"

# ── 3. Reset constants in config/config.go ───────────────────────────────────
echo "[3/4] Patching config/config.go..."

CONFIG_FILE="config/config.go"
NEW_BUNDLE_ID="com.${PROJECT_NAME}.${PROJECT_NAME}"

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' \
        -e "s|ProjectName     = \".*\"|ProjectName     = \"${PROJECT_NAME}\"|" \
        -e "s|ProjectVersion  = \".*\"|ProjectVersion  = \"v0.0.1\"|" \
        -e "s|ProjectBundleID = \".*\"|ProjectBundleID = \"${NEW_BUNDLE_ID}\"|" \
        "${CONFIG_FILE}"
else
    sed -i \
        -e "s|ProjectName     = \".*\"|ProjectName     = \"${PROJECT_NAME}\"|" \
        -e "s|ProjectVersion  = \".*\"|ProjectVersion  = \"v0.0.1\"|" \
        -e "s|ProjectBundleID = \".*\"|ProjectBundleID = \"${NEW_BUNDLE_ID}\"|" \
        "${CONFIG_FILE}"
fi

# ── 4. Replace template references in non-Go text files ──────────────────────
echo "[4/4] Replacing references in non-Go files..."

while IFS= read -r -d '' file; do
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' \
            -e "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g" \
            -e "s|${TEMPLATE_NAME}|${PROJECT_NAME}|g" \
            "${file}"
    else
        sed -i \
            -e "s|${TEMPLATE_MODULE}|${NEW_MODULE}|g" \
            -e "s|${TEMPLATE_NAME}|${PROJECT_NAME}|g" \
            "${file}"
    fi
done < <(find . -type f \
    \( -name "*.md" -o -name "*.yml" -o -name "*.yaml" \
    -o -name "*.xml" -o -name "*.iml" -o -name "*.sh"  \
    -o -name "*.json" -o -name "*.txt" \) \
    ! -path "./.git/*" \
    -print0)

echo ""
echo "Done!"
echo "  cd ${TARGET_DIR}"
echo "  module: ${NEW_MODULE}"
echo ""
echo "Next steps:"
echo "  git init && git add . && git commit -m 'chore: init from project_template_go'"

rm -f -- "$0"
