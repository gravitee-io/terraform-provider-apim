#!/bin/bash
set -euo pipefail

# Ensure HOME is set
if [ -z "${HOME:-}" ]; then
    echo "Error: HOME environment variable is not set" >&2
    exit 1
fi

# Get working directory with proper quoting
WORKDIR="$(pwd)"



# Ensure target directory exists
CONFIG_DIR="${HOME}"
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Error: Home directory $CONFIG_DIR does not exist" >&2
    exit 1
fi

# Create template with proper variable expansion
TF_TEMPLATE=$(cat <<EOF
provider_installation {
  dev_overrides {
    "registry.terraform.io/gravitee-io/apim" = "$WORKDIR"
    "registry.terraform.io/hashicorp/apim" = "$WORKDIR"
  }
  direct {}
}
EOF
)
# Write configuration with error checking
if ! echo "$TF_TEMPLATE"  > "${CONFIG_DIR}/.terraformrc"; then
    echo "Error: Failed to write configuration file" >&2
    exit 1
fi
echo "Configuration file successfully created at ${CONFIG_DIR}/.terraformrc"

# Create template with proper variable expansion
TF_TEMPLATE=$(cat <<EOF
provider_installation {
  dev_overrides {
    "registry.opentofu.org/gravitee-io/apim" = "$WORKDIR"
  }
  direct {}
}
EOF
)
if ! echo "$TF_TEMPLATE" > "${CONFIG_DIR}/.tofurc"; then
    echo "Error: Failed to write configuration file" >&2
    exit 1
fi
echo "Configuration file successfully created at ${CONFIG_DIR}/.tofurc"