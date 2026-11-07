---
page_title: "Authorization as code"
subcategory: "Authorization"
---

# PDP gateways, schemas, entities and policies

This example declares what the Authorization Management module enforces for one engine:

1. A PDP gateway, the engine of a `target_pdp_id` run on the gateways of a sharding `tag` (`apim_authz_pdp`)
2. The GAPL schema that engine checks against (`apim_authz_schema`)
3. The principals and resources that policies name (`apim_authz_entity`)
4. A policy deployed to that engine (`apim_authz_policy`)

It requires APIM 4.14 with the Authorization Management module and a license that includes it.

## Routing scope

Schemas, entities and policies name a PDP gateway by its routing scope: `target_pdp_id@tag` for a tagged gateway,
`target_pdp_id` for an untagged one, or `default`. `apim_authz_pdp` reports it in `routing_scope`: reference that
attribute rather than writing the scope by hand, so Terraform creates the PDP gateway first.

## PDP gateway

`target_pdp_id` and `tag` cannot be changed on an existing PDP gateway: another value is refused when planning.
`name`, `strict` and `endpoint` converge on apply. `strict` is shared by every PDP gateway of one `target_pdp_id`.

```terraform
# The engine of target_pdp_id "travel", run on the gateways tagged "eu", with its
# AuthZEN endpoint. target_pdp_id and tag are fixed at creation.
resource "apim_authz_pdp" "travel_eu" {
  hrid          = "travel-eu"
  name          = "Travel PDP (EU)"
  target_pdp_id = "travel"
  tag           = "eu"
  endpoint = {
    path = "/authz/travel-eu"
  }
}

```

~> APIM refuses with a `409` to delete a PDP gateway that is the only target of an entity or a policy, and deletes
nothing: retarget them, or destroy them in the same run.

## Schema

One schema per PDP gateway: a `target` another schema holds is refused. Without a `target` the schema is stored and
enforced nowhere.

```terraform
# One schema per PDP gateway, named by its routing scope: travel@eu.
resource "apim_authz_schema" "travel" {
  hrid   = "travel-schema"
  name   = "Travel schema"
  target = apim_authz_pdp.travel_eu.routing_scope
  schema = <<-EOT
    entity User in [Group];
    entity Group { department?: String };
    entity Hotel;
    action read appliesTo { principal: [User], resource: [Hotel] };
  EOT
}

```

## Entities

`entity_id`, `kind` and `entity_type` cannot be changed on an existing entity. `attributes`, `parents`, `targets` and
`active` are replaced whole on every apply: an attribute removed from the configuration is removed from the platform.
Write `attributes` with `jsonencode`: APIM reports the object in its own key order and spacing, which plans no change.

```terraform
resource "apim_authz_entity" "agents" {
  hrid        = "travel-agents"
  entity_id   = "agents"
  kind        = "PRINCIPAL"
  entity_type = "Group"
  attributes = jsonencode({
    department = "travel"
  })
  targets = [apim_authz_pdp.travel_eu.routing_scope]
}

resource "apim_authz_entity" "alice" {
  hrid        = "travel-alice"
  entity_id   = "alice"
  kind        = "PRINCIPAL"
  entity_type = "User"
  parents = [
    {
      entity_type = apim_authz_entity.agents.entity_type
      entity_id   = apim_authz_entity.agents.entity_id
    }
  ]
  targets = [apim_authz_pdp.travel_eu.routing_scope]
}

resource "apim_authz_entity" "grand_budapest" {
  hrid        = "travel-grand-budapest"
  entity_id   = "grand-budapest"
  kind        = "RESOURCE"
  entity_type = "Hotel"
  targets     = [apim_authz_pdp.travel_eu.routing_scope]
}

```

The platform reports `parents` and `targets` sorted. Declaring them in another order plans no change.

~> Deleting an entity deletes the RESOURCE entities that have it as a parent, and detaches the PRINCIPAL ones. APIM
refuses with a `413` a delete that would remove or detach more entities than its cascade limit.

## Policy

`entity_id` and `policy_type` cannot be changed on an existing policy. An apply that changes the policy deploys it
again, a policy disabled in the Console included. `policy_status` reports where it stands.

```terraform
# Deployed on every apply to the PDP gateways its targets name.
resource "apim_authz_policy" "agents_read_hotels" {
  hrid        = "travel-agents-read-hotels"
  entity_id   = "policy.travel.agents-read-hotels"
  name        = "Travel agents read hotels"
  policy_type = "CUSTOM"
  policy_text = <<-EOT
    permit(
      principal in Group::"agents",
      action == Action::"read",
      resource == Hotel::"grand-budapest"
    );
  EOT
  targets     = [apim_authz_pdp.travel_eu.routing_scope]
}

```

An empty `targets` deploys the policy to the default engine on every gateway. An MCP proxy asks the default engine
only: an `MCP` policy targets `default`.
