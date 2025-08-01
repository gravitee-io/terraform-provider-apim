#!/bin/bash
set -euo pipefail

# Color definitions
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly NC='\033[0m' # No Color

# Initialize counters
total_success=0
total_failures=0


# Check if test script exists and is executable
test_script="./tests/examples/tf_apply_destroy.sh"
if [[ ! -x "$test_script" ]]; then
    echo -e "${RED}Error: Test script '$test_script' not found or not executable${NC}" >&2
    exit 1
fi

# Store find results in array
directories=$(find examples/use-cases -name "provider.tf" -exec dirname {} \; || { echo "Error: Find command failed" >&2; exit 1; })

if [[ -z "$directories" ]]; then
    echo -e "${YELLOW}Warning: No directories found containing provider.tf${NC}" >&2
    exit 0
fi
failed_dirs=()
while read -r directory; do
    if  "$test_script" "$directory"; then
        ((total_success++))
    else
        ((total_failures++))
        failed_dirs+=("$directory")
    fi
done <<< "$directories"

total=$((total_success + total_failures))

if [[ "$total" -eq 0 ]]; then
    echo -e "${YELLOW}Warning: No tests were run${NC}" >&2
    exit 1
elif [[ "$total_failures" -gt 0 ]]; then
    echo -e "\nTest Results: ${RED}$total_failures/$total failed${NC}" >&2
    echo -e "\n${RED}Failure Details:${NC}"
    for dir in "${failed_dirs[@]}"; do
        echo -e "${RED}$dir${NC}" >&2
    done
    
    exit 1
else
    echo -e "Test Results: ${GREEN}$total passed${NC}"
    exit 0
fi