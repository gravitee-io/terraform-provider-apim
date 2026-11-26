/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { runPipeline, waitForPipeline } from "./lib/circleci.mjs";

import { LOG, toggleVerbosity } from "./lib/index.mjs";

// Runs an OAS check in CircleCI on a branch and exits with its result. CircleCI holds the GitHub token that reads the
// source repositories, some of them private.

const CHECKS = {
  apim: {
    trigger: "verify_oas",
    name: "APIM",
    done: "automation-api-oas.yaml matches gravitee-api-management at APIM_OAS_BRANCH",
  },
  aim: {
    trigger: "verify_aim_oas",
    name: "AI Management",
    done: "automation-api-aim-oas.yaml matches gravitee-gamma-module-aim at AIM_OAS_REF",
  },
  authz: {
    trigger: "verify_authz_oas",
    name: "AuthZ",
    done: "automation-api-authz-oas.yaml matches gravitee-gamma-module-authz at AUTHZ_OAS_REF",
  },
};

const VERBOSE = argv.verbose;
const BRANCH = argv["branch"];
const CHECK_ID = argv["check"] ?? "";
// Own keys only: "constructor" or "toString" would otherwise pass the check below.
const CHECK = Object.hasOwn(CHECKS, CHECK_ID) ? CHECKS[CHECK_ID] : undefined;

toggleVerbosity(VERBOSE);

if (!BRANCH) {
  LOG.red("--branch is required");
  process.exit(1);
}

if (!CHECK) {
  LOG.red(`--check is required, one of: ${Object.keys(CHECKS).join(", ")}`);
  process.exit(1);
}

const pipeline = await runPipeline({ trigger: CHECK.trigger }, BRANCH);

LOG.blue(`Pipeline is running at ${pipeline.url}`);

const workflows = await waitForPipeline(pipeline.id);
const failed = workflows.filter((w) => w.status !== "success");

for (const w of workflows) {
  LOG.log(`  ${w.name}: ${w.status}`);
}

if (failed.length > 0) {
  LOG.red(`The ${CHECK.name} OAS check failed, see ${pipeline.url}`);
  process.exit(1);
}

LOG.green(CHECK.done);
