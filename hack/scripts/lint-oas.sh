#!/usr/bin/env bash
# Lints one OpenAPI document and fails when Speakeasy counts an error or a warning.
# speakeasy lint exits 0 on warnings, --max-validation-errors and
# --max-validation-warnings included: they cap the output and gate nothing. The
# gate reads the summary line instead.
set -euo pipefail

schema=${1:?usage: lint-oas.sh <openapi document>}
out=$(mktemp); trap 'rm -f "$out"' EXIT

status=0
speakeasy lint openapi --schema "$schema" --non-interactive 2>&1 | tee "$out" || status=$?

summary=$(sed -e 's/\x1b\[[0-9;]*m//g' "$out" | grep -E 'linting complete\. [0-9]+ errors?, [0-9]+ warnings?' | tail -1 || true)
if [ -z "$summary" ]; then
  echo "$schema: no lint summary (speakeasy exited $status)"
  exit 1
fi

errors=$(sed -E 's/.* ([0-9]+) errors?,.*/\1/' <<<"$summary")
warnings=$(sed -E 's/.* ([0-9]+) warnings?.*/\1/' <<<"$summary")
if [ "$errors" -ne 0 ] || [ "$warnings" -ne 0 ] || [ "$status" -ne 0 ]; then
  echo "$schema: $errors error(s), $warnings warning(s); none is accepted"
  exit 1
fi
