#!/usr/bin/env bash
# Installs Terraform <version> (or "latest") on linux/amd64, for CircleCI machine executors.
set -euo pipefail

version="${1:?usage: install-terraform.sh <version|latest>}"
if [[ "$version" == "latest" ]]; then
  version=$(curl -fsSL https://checkpoint-api.hashicorp.com/v1/check/terraform | jq -r .current_version)
fi

work=$(mktemp -d); trap 'rm -rf "$work"' EXIT
curl -fsSLo "$work/terraform.zip" "https://releases.hashicorp.com/terraform/${version}/terraform_${version}_linux_amd64.zip"
unzip -q "$work/terraform.zip" -d "$work"
sudo install -m 0755 "$work/terraform" /usr/local/bin/terraform
terraform version
