#!/usr/bin/env bash
# An overlay action whose target matches nothing is ignored by Speakeasy without an error.
set -euo pipefail
work=$(mktemp -d); trap 'rm -rf "$work"' EXIT
speakeasy merge -s automation-api-oas.yaml -s automation-api-aim-oas.yaml -o "$work/0.yaml" >/dev/null
i=0
for overlay in $(yq '.sources.automation-api.overlays[].location' .speakeasy/workflow.yaml); do
  strict=--strict
  # merge.yaml removes a root security that the merge may or may not leave behind
  [ "$overlay" = .speakeasy/overlays/common/merge.yaml ] && strict=
  speakeasy overlay apply $strict -s "$work/$i.yaml" -o "$overlay" --out "$work/$((i+1)).yaml" >/dev/null \
    || { echo "$overlay: an action targets nothing"; exit 1; }
  i=$((i+1))
done
