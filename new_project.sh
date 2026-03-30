#!/usr/bin/env bash
# Usage:
#   bash new_project.sh <module_path>
#   wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- <module_path>
#
# Run from inside the cloned (empty) project directory:
#   git clone git@github.com:myorg/my_service.git && cd my_service
#   wget -qO- https://raw.githubusercontent.com/0xYeah/project_template_go/main/new_project.sh | bash -s -- github.com/myorg/my_service

set -euo pipefail

TEMPLATE_MODULE="github.com/0xYeah/project_template_go"
TEMPLATE_NAME="project_template_go"

usage() {
    echo "Usage: bash new_project.sh <module_path>"
    echo ""
    echo "  module_path   Go module path, must contain a dot"
    echo "                e.g. github.com/myorg/my_service"
    echo ""
    echo "Run from inside your cloned (empty) project directory."
    exit 1
}

[[ $# -lt 1 ]] && usage

NEW_MODULE="$1"

if [[ "${NEW_MODULE}" != *.* ]]; then
    echo "Error: invalid module path \"${NEW_MODULE}\": missing dot in first path element."
    echo "  Example: github.com/myorg/my_service"
    exit 1
fi

# Derive project name from last path segment
PROJECT_NAME="${NEW_MODULE##*/}"
TARGET_DIR="$(pwd)"
TMP_DIR="$(mktemp -d)"

echo "Template : ${TEMPLATE_MODULE}"
echo "New      : ${NEW_MODULE}"
echo "Target   : ${TARGET_DIR}"
echo ""

# ── 1. Install gonew if missing ──────────────────────────────────────────────
if ! command -v gonew &>/dev/null; then
    echo "[1/4] Installing gonew..."
    go install golang.org/x/tools/cmd/gonew@latest
else
    echo "[1/4] gonew: $(which gonew)"
fi

# ── 2. Clone + rename all Go files via gonew into tmp dir ────────────────────
echo "[2/4] Running gonew..."
gonew "${TEMPLATE_MODULE}" "${NEW_MODULE}" "${TMP_DIR}/scaffold"
mv "${TMP_DIR}/scaffold"/* "${TMP_DIR}/scaffold"/.[!.]* "${TARGET_DIR}"/ 2>/dev/null || true
rm -rf "${TMP_DIR}"

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
echo "Done! module: ${NEW_MODULE}"
echo ""
echo "Next steps:"
echo "  git add . && git commit -m 'chore: init from project_template_go'"

rm -f -- "$0"
