# Releasing `terraform-provider-apim`

## Versioning policy

- **0.x versions were beta**: they followed Terraform pre-1.0 semantics (breaking changes may occur without an upgrader).
- **1.x.y releases**: treat these as the stable release line and follow semver expectations.

## When to release

- **After merging to `main` or a maintenance branch**.
- You must release **after APIM is released**, so every resource endpoint needed by the provider is available.

## Automated release flow (default)

### 1. Prepare the release

Go to **GitHub Actions** → workflow **Prepare Release** → **Run workflow**.

- Select the branch:
  - **`main`** → bumps the **minor** version (e.g. `1.4.0` → `1.5.0`)
  - **`1.<minor>.x`** (maintenance branch) → bumps the **patch** version (e.g. `1.4.1` → `1.4.2`)

The workflow:
1. Bumps the version via `speakeasy bump`
2. Validates the new version (no tag collision, no skipped version)
3. Runs `make speakeasy` to regenerate code
4. Creates a PR titled `release: X.Y.Z`

### 2. Review and merge

- Review the PR as usual (check generated diffs, version in `gen.lock`)
- **Squash-merge the PR** — this is required for the automated release to trigger
- Keep the **squash commit title** exactly `release: X.Y.Z` (the default from the prepare-release PR). Do not change it during squash-merge.
- **Do not** use rebase-merge or a regular merge commit on the release PR — the Release workflow will not run (use the manual Release workflow as fallback).

#### Branch protection (recommended)

On `main` and each `1.<minor>.x` maintenance branch, configure branch protection so release PRs can only land via **squash merge** (disable merge commits and rebase merge for those branches, or restrict them to admins).

### 3. Automatic release

Once the squash-merge lands on `main` or a `1.*.x` branch (with `.speakeasy/gen.lock` updated), the **Release** workflow detects a commit whose subject is `release: X.Y.Z` and whose version matches `.speakeasy/gen.lock`, then:
1. Verifies the tag does not already exist
2. Creates and pushes an annotated tag `vX.Y.Z`
3. Runs GoReleaser (multi-platform binaries, checksums, GPG signing)
4. Sends a Slack notification
5. Publishes the changelog to the docs repo

The Terraform Registry picks up the GitHub release and makes the new version publicly available.

### 4. Post-release checks

1. **Check Terraform Registry** — the version should appear within minutes.
2. **Smoke test locally** using the released provider:
   - Remove any `~/.terraformrc` dev override if present
   - Pick an example (e.g. `examples/use-cases/v4api-proxy`) and run: `terraform init` → `terraform apply` → `terraform destroy`
3. **Notify** the team in `#gravite-release-alerts`.
4. If changelog publishing failed, check CircleCI workflow **publish-changelog** and the CI Slack channel.

## Manual release (fallback)

Use the manual trigger when:
- The release PR was merged with a regular merge (not squash) and the auto-trigger didn't fire
- You need a **dry run** to test the release pipeline
- You need to re-release or debug a failed release

Go to **GitHub Actions** → workflow **Release** → **Run workflow**:

- **Branch**: select the branch to release from (`main` or `1.<minor>.x`)
- **Dry run** (`dry_run`): checked by default — runs GoReleaser in snapshot mode without creating a tag or publishing artifacts. Uncheck for a real release.

## Maintenance releases (patch versions)

Maintenance releases are **patch-only** releases on a maintenance branch, for example: `1.4.1` → `1.4.2`.

### When the maintenance branch exists

After releasing a new **minor** from `main` (e.g. `1.4.0`), create a maintenance branch:
- Branch name: `1.4.x`
- This branch receives bugfixes and compatible changes for the **1.4.\*** line.
- Add a matching rule in `.mergify.yml` (label `apply-on-1-4-x` → branch `1.4.x`).

### How to release a maintenance patch

1. Merge or cherry-pick the desired fixes into `1.<minor>.x`.
2. Run the **Prepare Release** workflow **from the maintenance branch** (`1.<minor>.x`).
3. Review and squash-merge the PR it creates (title `release: X.Y.Z`).
4. The release triggers automatically.

### Backporting with Mergify

Use dash-style labels: `apply-on-1-4-x` to backport to branch `1.4.x` (dots in the branch name, dashes in the label). Each maintenance branch needs its own rule in `.mergify.yml` — copy the `1.4.x` example and adjust the label and branch name.

## Changelog path

The changelog file path is configured in `hack/apim.yaml` under `apim.doc.changelog`. This path is passed directly to the changelog publishing script — no manual APIM docs version input is needed. Update this field when the APIM docs version changes.

GoReleaser release notes include only commits matching `feat`, `fix`, `test`, and `docs` (see `.goreleaser.yaml`).
