#!/bin/bash
set -euo pipefail

# Check for required commands
if ! command -v envsubst >/dev/null 2>&1; then
    echo "Error: envsubst command not found. Please install gettext package." >&2
    exit 1
fi

# Ensure HOME is set
if [ -z "${HOME:-}" ]; then
    echo "Error: HOME environment variable is not set" >&2
    exit 1
fi

# Get working directory with proper quoting
WORKDIR="$(pwd)"

# Create template with proper variable expansion
TEMPLATE=$(cat <<EOF
provider_installation {
  dev_overrides {
    "registry.terraform.io/gravitee-io/apim" = "$WORKDIR"
    "registry.terraform.io/hashicorp/apim" = "$WORKDIR"
  }
  direct {}
}
EOF
)

# Ensure target directory exists
CONFIG_DIR="${HOME}"
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Error: Home directory $CONFIG_DIR does not exist" >&2
    exit 1
fi

# Write configuration with error checking
if ! echo "$TEMPLATE" | envsubst > "${CONFIG_DIR}/.terraformrc"; then
    echo "Error: Failed to write configuration file" >&2
    exit 1
fi

echo "Configuration file successfully created at ${CONFIG_DIR}/.terraformrc"