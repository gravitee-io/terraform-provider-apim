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

// Runs the verify-aim-oas CircleCI workflow on a branch and exits with its
// result. CircleCI holds the GitHub token that reads the private aim repository.

const VERBOSE = argv.verbose;
const BRANCH = argv["branch"];

toggleVerbosity(VERBOSE);

if (!BRANCH) {
  LOG.red("--branch is required");
  process.exit(1);
}

const pipeline = await runPipeline({ trigger: "verify_aim_oas" }, BRANCH);

LOG.blue(`Pipeline is running at ${pipeline.url}`);

const workflows = await waitForPipeline(pipeline.id);
const failed = workflows.filter((w) => w.status !== "success");

for (const w of workflows) {
  LOG.log(`  ${w.name}: ${w.status}`);
}

if (failed.length > 0) {
  LOG.red(`The AI Management OAS fragment check failed, see ${pipeline.url}`);
  process.exit(1);
}

LOG.green("automation-api-aim-oas.yaml matches gravitee-gamma-module-aim at AIM_OAS_REF");
