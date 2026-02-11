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

import {APIM} from "./lib/apim.mjs";
import {LOG, PROJECT_DIR, setNoQuoteEscape, setQuoteEscape, time, toggleVerbosity,} from "./lib/index.mjs";

const KIND_CONFIG = path.join(PROJECT_DIR, "hack", "kind");
const PKI = path.join(
    KIND_CONFIG,
    "pki",
);

const APIM_IMAGE_REGISTRY = await getAPIMImageRegistry();
const APIM_IMAGE_TAG = await getAPIMImageTag();
const APIM_CHART_REGISTRY = await getAPIMChartRegistry();
const APIM_CHART_VERSION = await getAPIMChartVersion();

const APIM_MINIMAL = $.env.APIM_MINIMAL === "true";
const APIM_VALUES = APIM_MINIMAL ? "values-minimal.yaml" : `${$.env.APIM_VALUES || "values.yaml"}`;

const IMAGES = new Map([
    [
        `${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}`,
        `gravitee-apim-gateway:dev`,
    ],
    [
        `${APIM_IMAGE_REGISTRY}/apim-management-api:${APIM_IMAGE_TAG}`,
        `gravitee-apim-management-api:dev`,
    ],
    [
        `${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}`,
        `gravitee-apim-management-ui:dev`,
    ],
    [`mccutchen/go-httpbin:latest`, `go-httpbin:dev`],
    [`mongo:7.0.23-jammy`, `mongo:7.0.23-jammy`],
]);

if (APIM_MINIMAL) {
    IMAGES.delete(`${APIM_IMAGE_REGISTRY}/apim-gateway:${APIM_IMAGE_TAG}`)
    IMAGES.delete(`${APIM_IMAGE_REGISTRY}/apim-management-ui:${APIM_IMAGE_TAG}`)
    IMAGES.delete(`mccutchen/go-httpbin:latest`)
}

if (!argv.verbose) {
    $.quiet = true;
}

toggleVerbosity(argv.verbose);

async function getAPIMImageRegistry() {
    if ($.env.APIM_IMAGE_REGISTRY) {
        return $.env.APIM_IMAGE_REGISTRY;
    }
    return await APIM.getImageRegistry();
}

async function getAPIMImageTag() {
    if ($.env.APIM_IMAGE_TAG) {
        return $.env.APIM_IMAGE_TAG;
    }
    return await APIM.getImageTag();
}

async function getAPIMChartRegistry() {
    if ($.env.APIM_CHART_REGISTRY) {
        return $.env.APIM_CHART_REGISTRY;
    }
    return await APIM.getChartRegistry();
}

async function getAPIMChartVersion() {
    if ($.env.APIM_CHART_VERSION) {
        return $.env.APIM_CHART_VERSION;
    }
    return await APIM.getChartVersion();
}

async function createKindCluster() {
    setNoQuoteEscape();
    await $`kind create cluster --config ${KIND_CONFIG}/kind.yaml`;
    setQuoteEscape();
}

async function imageExists(image) {
    try {
        await $`docker image inspect ${image}`;
        return true;
    } catch {
        return false;
    }
}

async function loadImages() {
    setNoQuoteEscape();

    for (const [image, tag] of IMAGES.entries()) {
        if (await imageExists(image)) {
            LOG.blue(`image ${image} found locally, skipping pull`);
        } else {
            LOG.blue(`pulling image ${image}`);
            await $`docker pull ${image}`;
        }
        LOG.blue(`tagging image ${image} with ${tag}`);
        await $`docker tag ${image} ${tag}`;
        LOG.blue(`loading image tag ${tag}`);
        await $`kind load docker-image ${tag} --name gravitee`;
    }

    setQuoteEscape();
}

async function createGraviteeNamespace() {
    await $`kubectl create ns gravitee`;
}

async function createTLSSecret() {
    await $`kubectl create secret tls tls-server --cert=${PKI}/server.crt --key=${PKI}/server.key`;
}

async function helmInstallAPIM() {
    await $`helm repo add graviteeio https://helm.gravitee.io`;
    await $`helm repo update graviteeio`;
    if ($.env.APIM_GRAVITEE_LICENSE) {
        let lic = $.env.APIM_GRAVITEE_LICENSE
        await $`helm install apim ${APIM_CHART_REGISTRY} -f ${KIND_CONFIG}/apim/${APIM_VALUES} --set license.key=${lic} --version ${APIM_CHART_VERSION}`;
    } else {
        await $`helm install apim ${APIM_CHART_REGISTRY} -f ${KIND_CONFIG}/apim/${APIM_VALUES} --version ${APIM_CHART_VERSION}`;
    }

}

async function deployHTTPBin() {
    await $`kubectl apply -f ${KIND_CONFIG}/httpbin`;
}

async function waitForApim() {
    await $`kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=apim3 --timeout=360s`;
}

async function fetchWithRetry(url, options, {retries = 10, baseDelay = 1000, maxDelay = 30000} = {}) {
    for (let attempt = 1; attempt <= retries; attempt++) {
        try {
            const resp = await fetch(url, options);
            return resp;
        } catch (err) {
            if (attempt === retries) {
                throw err;
            }
            const delay = Math.min(baseDelay * 2 ** (attempt - 1), maxDelay);
            LOG.blue(`  fetch ${url} failed (${err.cause?.code || err.message}), retrying in ${delay / 1000}s (${attempt}/${retries})`);
            await sleep(delay);
        }
    }
}

async function configureAPIM() {
    const apimConfig = path.join(PROJECT_DIR, "hack", "apim");
    const settings = await fs.readFile(path.join(apimConfig, "settings.json"), "utf-8");
    const dcr = await fs.readFile(path.join(apimConfig, "dcr.json"), "utf-8");

    const adminAuth = "Basic " + Buffer.from("admin:admin").toString("base64");
    const api1Auth = "Basic " + Buffer.from("api1:api1").toString("base64");
    const baseUrl = "http://localhost:30083/management/organizations/DEFAULT/environments/DEFAULT";

    let resp;

    resp = await fetchWithRetry(`${baseUrl}/settings`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": adminAuth,
        },
        body: settings,
    });
    if (resp.status >= 400) {
        throw new Error(`POST settings failed (${resp.status}): ${await resp.text()}`);
    }

    resp = await fetchWithRetry(`${baseUrl}/configuration/applications/registration/providers`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Authorization": adminAuth,
        },
        body: dcr,
    });
    if (resp.status >= 400) {
        throw new Error(`POST DCR provider failed (${resp.status}): ${await resp.text()}`);
    }

    resp = await fetchWithRetry("http://localhost:30083/management/organizations/DEFAULT/environments", {
        headers: {
            "Authorization": api1Auth,
        },
    });
    if (resp.status >= 400) {
        throw new Error(`GET environments failed (${resp.status}): ${await resp.text()}`);
    }
}

let steps = 8
if (APIM_MINIMAL) {
    steps = 6
}
let step = 1;

LOG.blue(`
  ☸ [${step++}/${steps}] Initializing kind cluster
`);

await time(createKindCluster);

LOG.blue(`
  🐳 [${step++}/${steps}] Loading docker images
`);

await time(loadImages);

LOG.blue(`
  ☸ [${step++}/${steps}] Creating gravitee namespace
`);

await time(createGraviteeNamespace);

if (!APIM_MINIMAL) {
    LOG.blue(`
  ☸ [${step++}/${steps}] Creating APIM gateway TLS secret
    `);

    await time(createTLSSecret);

    LOG.blue(`
  ☸ [${step++}/${steps}] Deploying httpbin
`);
    await time(deployHTTPBin);
}

LOG.blue(`
  ☸ [${step++}/${steps}] Installing APIM
`);

await time(helmInstallAPIM);

if (APIM_MINIMAL) {
    LOG.magenta(`
    APIM Minimal containers are starting ...

    Version: ${APIM_IMAGE_TAG}

    Management API      http://localhost:30083/management
    Automation API      http://localhost:30083/automation
`);
} else {
    LOG.magenta(`
    APIM containers are starting ...

    Version: ${APIM_IMAGE_TAG}

    Available endpoints are:
        Gateway             http://localhost:30082
        Gateway with mTLS   https://localhost:30084
        Management API      http://localhost:30083/management
        Automation API      http://localhost:30083/automation
        Console             http://localhost:30080
`);
}
LOG.blue(`
  ⚙ [${step++}/${steps}] Waiting for services to be ready ...
`);

await time(waitForApim);

LOG.blue(`
  ⚙ [${step++}/${steps}] Configuring APIM
`);

await time(configureAPIM);
