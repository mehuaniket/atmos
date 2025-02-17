---
title: Inline Component Customization Atmos Design Pattern
sidebar_position: 3
sidebar_label: Inline Component Customization
description: Inline Component Customization Atmos Design Pattern
---

# Inline Component Customization

The **Inline Component Customization** pattern is used when the defaults for the [components](/core-concepts/components) in
a [stack](/core-concepts/stacks)
are configured in default manifests, the manifests are [imported](/core-concepts/stacks/imports) into the top-level stacks, and the components
are customized inline in each top-level stack overriding the configuration for each environment (OU, account, region).

## Applicability

Use the **Inline Component Customization** pattern when:

- You have components that are provisioned in multiple stacks (e.g. `dev`, `staging`, `prod` accounts) with different configurations for each stack

- You need to make the components' default/baseline configurations reusable across different stacks

- You want to keep the configurations [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

Suppose you need a simple setup with only `dev`, `staging` and `prod` stages (accounts). Here's how you might organize the stacks and
component configurations for the `vpc` and `vpc-flow-logs-bucket` components, and then customize the components in the stacks.

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   ├── defaults  # component-specific defaults
   │   │   ├── vpc-flow-logs-bucket.yaml
   │   │   └── vpc.yaml
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
  excluded_paths:
    # Tell Atmos that the `defaults` folder and all sub-folders don't contain top-level stack manifests
    - "defaults/**/*"

schemas:
  jsonschema:
    base_path: "stacks/schemas/jsonschema"
  opa:
    base_path: "stacks/schemas/opa"
  atmos:
    manifest: "stacks/schemas/atmos/atmos-manifest/1.0/atmos-manifest.json"
```

Add the following default configuration to the `stacks/defaults/vpc-flow-logs-bucket.yaml` manifest:

```yaml title="stacks/defaults/vpc-flow-logs-bucket.yaml"
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
```

Add the following default configuration to the `stacks/defaults/vpc.yaml` manifest:

```yaml title="stacks/defaults/vpc.yaml"
components:
  terraform:
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

Configure the `stacks/dev.yaml` top-level stack manifest:

```yaml title="stacks/dev.yaml"
vars:
  stage: dev

# Import the component default configurations
import:
  - defaults/vpc

components:
  terraform:
    # Customize the `vpc` component for the `dev` account
    # You can define variables or override the imported defaults
    vpc:
      vars:
        max_subnet_count: 2
        vpc_flow_logs_enabled: false
```

Configure the `stacks/staging.yaml` top-level stack manifest:

```yaml title="stacks/staging.yaml"
vars:
  stage: staging

# Import the component default configurations
import:
  - defaults/vpc-flow-logs-bucket
  - defaults/vpc

components:
  terraform:
    # Customize the `vpc` component for the `staging` account
    # You can define variables or override the imported defaults
    vpc:
      vars:
        map_public_ip_on_launch: false
        vpc_flow_logs_traffic_type: "REJECT"
```

Configure the `stacks/prod.yaml` top-level stack manifest:

```yaml title="stacks/prod.yaml"
vars:
  stage: prod

# Import the component default configurations
import:
  - defaults/vpc-flow-logs-bucket
  - defaults/vpc

components:
  terraform:
    # Customize the `vpc` component for the `prod` account
    # You can define variables or override the imported defaults
    vpc:
      vars:
        map_public_ip_on_launch: false
```

To provision the components, execute the following commands:

```shell
# `dev` stack
atmos terraform apply vpc -s dev

# `staging` stack
atmos terraform apply vpc-flow-logs-bucket -s staging
atmos terraform apply vpc -s staging

# `prod` stack
atmos terraform apply vpc-flow-logs-bucket -s prod
atmos terraform apply vpc -s prod
```

## Benefits

The **Inline Component Customization** pattern provides the following benefits:

- The defaults for the components are defined in just one place making the entire
  configuration [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

- The defaults for the components are reusable across many stacks

- Simple stack and component configurations

## Limitations

The **Inline Component Customization** pattern has the following limitations and drawbacks:

- The pattern is useful to customize components per account or region, but if you have more than one Organization, Organizational Unit (OU) or region,
  then the inline customizations would be repeated in the stack manifests, making the entire stack configuration
  not [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

- Should be used only for specific use-cases, e.g. when you use just one region, Organization or Organizational Unit (OU)

:::note

To address the limitations of the **Inline Component Customization** Design Pattern, consider using the following patterns:

- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)

:::

## Related Patterns

- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Component Inheritance](/design-patterns/component-inheritance)
- [Partial Component Configuration](/design-patterns/partial-component-configuration)
