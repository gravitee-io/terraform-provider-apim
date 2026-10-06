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

import { isEmptyString, LOG } from "./index.mjs";

const API_BASE = "https://circleci.com/api/v2";
const APP_BASE = "https://app.circleci.com";
const ORG = "gravitee-io";
const SCM = "github";
const PROJECT = "terraform-provider-apim";
const CIRCLECI_TOKEN = process.env.CIRCLECI_TOKEN;

if (isEmptyString(CIRCLECI_TOKEN)) {
  LOG.red("CIRCLECI_TOKEN cannot be found");
  process.exit(1);
}

export async function triggerPipeline(parameters, branch = "master") {
  const { url } = await runPipeline(parameters, branch);
  return url;
}

export async function runPipeline(parameters, branch = "master") {
  const response = await fetch(
    `${API_BASE}/project/${SCM}/${ORG}/${PROJECT}/pipeline`,
    {
      method: "POST",
      headers: {
        "Circle-Token": CIRCLECI_TOKEN,
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify({ parameters, branch }),
    },
  );

  if (response.status === 201) {
    const json = await response.json();
    return {
      id: json.id,
      number: json.number,
      url: `${APP_BASE}/pipelines/${SCM}/${ORG}/${PROJECT}/${json.number}`,
    };
  }

  throw new Error(`Unable to run pipeline (HTTP status ${response.status})`);
}

const PENDING_WORKFLOW_STATUSES = ["running", "failing", "on_hold"];

// Resolves with the pipeline's workflows once none is pending. Rejects when the
// pipeline's config fails to compile or the timeout elapses.
export async function waitForPipeline(id, { timeout = 15 * 60_000, interval = 10_000 } = {}) {
  const deadline = Date.now() + timeout;
  while (Date.now() < deadline) {
    const pipeline = await get(`${API_BASE}/pipeline/${id}`);
    if (pipeline.state === "errored") {
      throw new Error(`Pipeline errored: ${JSON.stringify(pipeline.errors)}`);
    }
    const { items: workflows } = await get(`${API_BASE}/pipeline/${id}/workflow`);
    const settled =
      workflows.length > 0 &&
      workflows.every((w) => !PENDING_WORKFLOW_STATUSES.includes(w.status));
    if (settled) {
      return workflows;
    }
    await new Promise((resolve) => setTimeout(resolve, interval));
  }
  throw new Error(`Pipeline ${id} still running after ${timeout / 1000}s`);
}

async function get(url) {
  const response = await fetch(url, {
    headers: { "Circle-Token": CIRCLECI_TOKEN, Accept: "application/json" },
  });
  if (!response.ok) {
    throw new Error(`GET ${url} answered HTTP ${response.status}`);
  }
  return response.json();
}
