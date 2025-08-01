#!/bin/bash
set -euo pipefail

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Check for terraform installation
if ! command -v terraform >/dev/null 2>&1; then
    echo -e "${RED}ERROR: terraform is not installed${NC}" >&2
    exit 1
fi

# Store the initial directory
readonly INITIAL_DIR=$(pwd) || {
    echo -e "${RED}ERROR: Failed to get current directory${NC}" >&2
    exit 1
}

# Ensure we return to initial directory on exit
trap 'cd "${INITIAL_DIR}"' EXIT

handle_error() {
    local step=$1
    echo -e "${RED}================================================================" >&2
    echo -e "ERROR: Failed during step: $step" >&2
    echo -e "================================================================${NC}" >&2
    exit 1
}

# Check if directory argument is provided
if [ -z "${1:-}" ]; then
    echo -e "${RED}ERROR: Directory argument is required${NC}" >&2
    exit 1
fi

# Check if directory exists
if [ ! -d "${1}" ]; then
    echo -e "${RED}ERROR: Directory ${1} does not exist${NC}" >&2
    exit 1
fi

cd "${1}" || handle_error "Changing to directory ${1}"
echo -e "${BLUE}Testing example: ${1}${NC}"

# Clean up with error handling
rm -rf .terraform* || handle_error "Cleaning .terraform files"
rm -f terraform.tfstate* || handle_error "Cleaning terraform state files"

echo -e "${BLUE}Running terraform apply...${NC}"
terraform apply -auto-approve || handle_error "terraform apply"

echo -e "${BLUE}Running terraform destroy...${NC}"
terraform destroy -auto-approve || handle_error "terraform destroy"

echo -e "${GREEN}All steps completed successfully${NC}"