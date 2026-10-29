#!/usr/bin/env bash
# Fail if Speakeasy removes a hand-maintained request region during generation.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

check_region_count() {
  local file="$1"
  local marker="$2"
  local expected="$3"
  local actual

  actual=$(grep -F -c "$marker" "$file" || true)
  if [[ "$actual" -ne "$expected" ]]; then
    echo "$file: expected $expected occurrence(s) of '$marker', found $actual" >&2
    return 1
  fi
}

check_region_count internal/provider/provider.go '// BEGIN GRAVITEE CLOUD INIT' 1
check_region_count internal/provider/provider.go '// END GRAVITEE CLOUD INIT' 1

check_region_count internal/provider/apilink_resource.go '// BEGIN GRAVITEE API LINK REQUEST' 2
check_region_count internal/provider/apilink_resource.go '// END GRAVITEE API LINK REQUEST' 2
check_region_count internal/provider/documentationapi_resource.go '// BEGIN GRAVITEE API DOCUMENTATION REQUEST' 2
check_region_count internal/provider/documentationapi_resource.go '// END GRAVITEE API DOCUMENTATION REQUEST' 2
check_region_count internal/provider/documentationportal_resource.go '// BEGIN GRAVITEE PORTAL DOCUMENTATION REQUEST' 2
check_region_count internal/provider/documentationportal_resource.go '// END GRAVITEE PORTAL DOCUMENTATION REQUEST' 2
check_region_count internal/provider/portallink_resource.go '// BEGIN GRAVITEE PORTAL LINK REQUEST' 2
check_region_count internal/provider/portallink_resource.go '// END GRAVITEE PORTAL LINK REQUEST' 2
check_region_count internal/provider/portal_resource.go '// BEGIN GRAVITEE PORTAL REQUEST' 2
check_region_count internal/provider/portal_resource.go '// END GRAVITEE PORTAL REQUEST' 2
