#!/usr/bin/env bash
# An overlay action whose target matches nothing is ignored by Speakeasy without an error.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"
work=$(mktemp -d); trap 'rm -rf "$work"' EXIT
overlays=$(yq -e '.sources.automation-api.overlays[].location' .speakeasy/workflow.yaml)
speakeasy merge -s automation-api-oas.yaml -s automation-api-aim-oas.yaml -o "$work/0.yaml" >/dev/null
i=0
for overlay in $overlays; do
  case "$overlay" in
    # merge.yaml removes a root security that the merge may or may not leave behind
    .speakeasy/overlays/common/merge.yaml) strict= ;;
    # Speakeasy keeps method names for the paths held-back.yaml removes
    .speakeasy/speakeasy-modifications-overlay.yaml) strict= ;;
    *) strict=--strict ;;
  esac
  if ! speakeasy overlay apply $strict -s "$work/$i.yaml" -o "$overlay" --out "$work/$((i+1)).yaml" >"$work/apply.log" 2>&1; then
    cat "$work/apply.log"
    echo "$overlay: failed to apply, see above ('did not match any targets' names a dead target)"
    exit 1
  fi
  i=$((i+1))
done
