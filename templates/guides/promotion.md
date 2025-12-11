---
page_title: "Promote resources"
subcategory: "Tutorials"
---

## Purpose and must-follow rules
Terraform provider allows easy promotion of APIs, but any resources it can manage, which is more that you can do with the console.

There is only a simple rule that one must follow to promote: **one environment (and organization) per state**.

~> We recommend a **remote backend** with **separate state keys per env**, `.tfstate` should not be committed.
The examples above do not show state management.

This usually means you use either:
* a separate directory per environment
* a separate branch per environment

## Step-by-step example with branches

For the sake of simplicity, the following example is a **Shared Policy Group** as it is quite a small resource,
**but note it is applicable to any resources**.

~> If you are testing this setup with local states, make sure state files are removed/restored when you switch between branches.

### Step 1: resource creation

Consider a Git-managed directory containing the following file:

```terraform
terraform {
  required_providers {
    apim = {
      source = "gravitee-io/apim"
    }
  }
}

variable "environment_id" {
  type = string
}

provider "apim" {
  server_url = "https://gravitee-apim-ctrl-plane.apis.acme.com/automation"
  organization_id = "<ORGANIZATION_UUID>"
  environment_id = var.environment_id
}

resource "apim_shared_policy_group" "simple" {
  hrid        = "simple"
  name        = "[Terraform] Http callout shared policy"
  api_type    = "PROXY"
  phase       = "REQUEST"
  steps = [
    {
      name = "HTTP Callout",
      enabled = true,
      policy = "policy-http-callout",
      configuration = jsonencode(
        {
          method = "GET",
          // no URL yet
          exitOnError = true
        })
    },
  ]
}
```

You can commit (and push) this file.

You may have noted that `url` lacks in the `policy-http-callout` configuration. 
This is meant to illustrate that the `main` branch should contain what is stable regardless of environments. 

### Step 2: apply to the first environment 
Then, you can use a new branch: 

```shell
git checkout -b gravitee/dev
```

Edit the file and set what is missing: the URL for the dev environment. 

```terraform
resource "apim_shared_policy_group" "simple" {
  // ...
  steps = [
    {
      // ...
      configuration = jsonencode(
        {
          method = "GET",
          url = "http://dev.callout-tracing.acme.com"     // +
         // ...
        })
    },
  ]
}
```

Resource is ready to be committed and pushed.

Resources can be applied with appropriate variables set (environment and credentials):

Here is what should be done using the CLI:
```shell
TF_VAR_env_id='<ENV_UUID>' APIM_SA_TOKEN='<SERVICE_ACCOUNT_TOKEN>' terraform apply
```

Now that the state is created, from now on, the branch is bound to `dev` environment.

### Step 3: Promote

Let's create a new branch. 

```shell
git checkout -b gravitee/prod
```

Let's make some changes for the production environment:
* Edit the URL
* Add some credentials

```terraform
// ... 

variable "callout_key" {
  type = "string"
}

resource "apim_shared_policy_group" "simple" {
  // ...
  steps = [
    {
      // ...
      configuration = jsonencode(
        {
          method = "GET",
          headers = [                                    // +
            {                                            // +
              name = "Authorization",                    // +
              value = "ApiKey ${var.callout_key}"        // +
            }                                            // +
          ],                                             // +
          // prod URL
          url = "http://prod.callout-tracing.acme.com"   // -/+
          // ..
        })
    },
  ]
}

// ...
```

You can see here we need to add a new configuration section (headers) not just a value. 
In that case it is convenient to have a git branch to implement those kinds of env-specific configurations.

Resource is ready versioned and deployed.

Here is what should be done using the CLI:
```shell
TF_VAR_env_id='<ENV_UUID>' \
  APIM_SA_TOKEN='<SERVICE_ACCOUNT_TOKEN>' \
  TF_VAR_callout_key='<API_KEY>' \
  terraform apply
```

Promotion is complete! The resource is applied to the `prod` environment.

### Step 4: make some changes

Let's go back to the main branch.

```shell
git checkout main
```

Consider the following change 

```terraform
// ...
resource "apim_shared_policy_group" "simple" {
  // ...
  steps = [
    {
      //...
      configuration = jsonencode(
        {
          method = "GET",
          // no URL yet
          exitOnError    = false    // -/+
          fireAndForget  = true     // + 
          useSystemProxy = true     // +
        })
    },
  ]
}
```
With this change the call is asynchronous, system proxy is used, and error tolerance is applied.  Here we assume this is core and needs to apply on all environments.

Apply it to `dev` environment (we assume it will be applied, CLI command is the same as above)

```shell
git rebase main dev
git push -f
```

Apply it to `prod` environment (same as above)

```shell
git rebase dev prod
git push -f
```

Et voilà! Changes are applied to all environments by changing only what matters most, avoid side effects.

### Discussion

It might be a surprising approach in the Terraform world, but some aspects rather specific to Gravitee should be considered.

#### One branch / env.

Having a branch per environment approach is better suited if environment changes are mostly owned by API publishers.

In this scenario:
- _Branch-per-env:_ URLs/hosts, quotas/sampling, toggles, per-env routing
- _Variables/secret manager:_ credentials, API keys, tokens, client secrets
- _Main branch:_ shared defaults and structural resource changes

Sensitive values, such as credentials, should be managed in appropriate secured storage and only set using Terraform variables 
as environment vars: `TF_VAR_...`. Avoid using `.tfvars` for credentials as they would be visible in source control.

#### One branch for all env.

This means that all changes (including the one shown above: add a section to a resource that is specific to an environment)
must be done differently (templating, externalize policy configuration, conditions...).  

Your resource will be more complex. 
For instance, API is a complex resource, you may want to consider using the simplest possible approach to ease maintenance of your resources.

This approach also means that DevOps teams should manage all env-specific changes (URLs/hosts, quotas/sampling...) 
as well as sensitive values management; it might be a source of error in some cases.

## Directory-based configuration

The above tutorial can be done with plain old directories instead of branches.
* `dev/` to contains dev files
* `prod/` to contains prod files

Each could optionally contain `.tfvar` files to set urls/credentials/environment information.

With this approach, changes are be harder to carry out, 
especially if more complex resources like APIs contain environment-specific configurations (endpoint targets, env-specific API resources...). 
As simple copy and paste may not suffice, although Terraform variables can reduce the burden as discussed above.

## Gravitee cloud

The above examples can be conducted with Gravitee Cloud.

```terraform
provider "apim" {
  cloud_auth = "<CLOUD_TOKEN>"
}
```

Only `cloud_auth` or (`APIM_CLOUD_TOKEN` environment var) needs to be changed for the `provider` section for a given environment (i.e branch or environment) to another
as organization/environment/url are inferred from the token. Other env-specific changes can be done as mentioned above. 
