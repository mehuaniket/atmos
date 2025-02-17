---
title: Inline Component Configuration Atmos Design Pattern
sidebar_position: 2
sidebar_label: Inline Component Configuration
description: Inline Component Configuration Atmos Design Pattern
---

# Inline Component Configuration

The **Inline Component Configuration** pattern is used when the [components](/core-concepts/components) in a [stack](/core-concepts/stacks) 
are configured inline in the stack manifest without [importing](/core-concepts/stacks/imports) and reusing default/base configurations. If you're starting with Atmos, this will be the easiest to understand until you need the more advanced features.

## Applicability

Use the **Inline Component Configuration** pattern when:

- You have a very simple organizational structure, e.g. one OU, one or a few accounts, one region

- You have a component that is provisioned only in one stack (e.g. only in the `dev` account). In this case, the component is configured inline in the
  stack manifest and is not used in other stacks

- For testing or development purposes

## Example

Suppose you need a simple setup with only `dev`, `staging` and `prod` stages (accounts). Here's how you might organize the stacks and
inline-component configuration for the `vpc` and `vpc-flow-logs-bucket` components.

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   ├── dev.yaml
   │   ├── staging.yaml
   │   └── prod.yaml
   │  
   │   # Centralized components configuration
   └── components
       └── terraform  # Terraform components (a.k.a Terraform "root" modules)
           ├── vpc
           ├── vpc-flow-logs-bucket
           ├── < other components >
```

Add the following minimal configuration to `atmos.yaml` [CLI config file](/cli/configuration) :

```yaml title="atmos.yaml"
components:
  terraform:
    base_path: "components/terraform"

stacks:
  base_path: "stacks"
  name_pattern: "{stage}"

schemas:
  jsonschema:
    base_path: "stacks/schemas/jsonschema"
  opa:
    base_path: "stacks/schemas/opa"
  atmos:
    manifest: "stacks/schemas/atmos/atmos-manifest/1.0/atmos-manifest.json"
```

Add the following component configurations to the `stacks/dev.yaml` stack manifest:

```yaml title="stacks/dev.yaml"
vars:
  stage: dev

components:
  terraform:
    vpc-flow-logs-bucket:
      metadata:
        # Point to the Terraform component
        component: vpc-flow-logs-bucket
      vars:
        enabled: true
        name: "vpc-flow-logs"
        traffic_type: "ALL"
        force_destroy: true
        lifecycle_rule_enabled: false

    vpc:
      metadata:
        # Point to the Terraform component
        component: vpc
      settings:
        # All validation steps must succeed to allow the component to be provisioned
        validation:
          validate-vpc-component-with-jsonschema:
            schema_type: jsonschema
            schema_path: "vpc/validate-vpc-component.json"
            description: Validate 'vpc' component variables using JSON Schema
          check-vpc-component-config-with-opa-policy:
            schema_type: opa
            schema_path: "vpc/validate-vpc-component.rego"
            # An array of filesystem paths (folders or individual files) to the additional modules for schema validation
            # Each path can be an absolute path or a path relative to `schemas.opa.base_path` defined in `atmos.yaml`
            # In this example, we have the additional Rego modules in `stacks/schemas/opa/catalog/constants`
            module_paths:
              - "catalog/constants"
            description: Check 'vpc' component configuration using OPA policy
      vars:
        enabled: true
        name: "common"
        max_subnet_count: 3
        map_public_ip_on_launch: true
        dns_hostnames_enabled: true
        assign_generated_ipv6_cidr_block: false
        nat_gateway_enabled: true
        nat_instance_enabled: false
        vpc_flow_logs_enabled: true
        vpc_flow_logs_traffic_type: "ALL"
        vpc_flow_logs_log_destination_type: "s3"
        nat_eip_aws_shield_protection_enabled: false
        subnet_type_tag_key: "acme/subnet/type"
        ipv4_primary_cidr_block: 10.9.0.0/18
```

To provision the components, execute the following commands:

```shell
# `dev` stack
atmos terraform apply vpc-flow-logs-bucket -s dev
atmos terraform apply vpc -s dev
```

## Benefits

The **Inline Component Configuration** pattern provides the following benefits:

- Very simple stack and component configurations

- All components are defined in just one place (in one stack manifest) - easier to see what is provisioned and where

## Limitations

The **Inline Component Configuration** pattern has the following limitations and drawbacks:

- If you have more than one stack (e.g. `dev`, `staging`, `prod`), then the component definitions would be repeated in the stack manifests,
  which makes them not reusable and the entire stack configuration not [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

- Should be used only for specific use-cases (e.g. you have just one stack, or you are designing and testing the components)

:::note

Use the [Inline Component Customization](/design-patterns/inline-component-customization) Design Pattern to address the limitations of the
**Inline Component Configuration** pattern.

:::

## Related Patterns

- [Inline Component Customization](/design-patterns/inline-component-customization)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Component Inheritance](/design-patterns/component-inheritance)
- [Partial Component Configuration](/design-patterns/partial-component-configuration)
