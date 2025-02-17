---
title: Component Inheritance Atmos Design Pattern
sidebar_position: 8
sidebar_label: Component Inheritance
description: Component Inheritance Atmos Design Pattern
---

# Component Inheritance

The **Component Inheritance** Design Pattern describes the mechanism of deriving Atmos components from one or more base components, allowing reusing
the base components' configurations.

In Atmos, **Component Inheritance** is the mechanism of deriving a component from one or more base components, inheriting all the
properties of the base component(s) and overriding only some fields specific to the derived component. The derived component acquires all the
properties of the base component(s), allowing creating very DRY configurations that are built upon existing components.

:::note
Atmos supports many different types on component inheritance. Refer to [Component Inheritance](/core-concepts/components/inheritance) for
more details.
:::

## Applicability

Use the **Component Inheritance** pattern when:

- You need to have reusable base components that serve as blueprints for the derived Atmos components

- You need to keep the configuration of all components [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

The following example shows the Atmos stack and component configurations to provision the `vpc` component into
a multi-account, multi-region environment. In the `catalog/vpc` folder, we have the `defaults.yaml` manifest that configures the base
component `vpc/defaults` to be inherited by all the derived VPC components in all stacks.

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   ├── catalog  # component-specific defaults
   │   │   └── vpc
   │   │       └── defaults.yaml
   │   ├── mixins
   │   │   ├── tenant  # tenant-specific defaults
   │   │   │   └── plat.yaml
   │   │   ├── region  # region-specific defaults
   │   │   │   ├── us-east-2.yaml
   │   │   │   └── us-west-2.yaml
   │   │   └── stage  # stage-specific defaults
   │   │       ├── dev.yaml
   │   │       ├── staging.yaml
   │   │       └── prod.yaml
   │   └── orgs  # Organizations
   │       └── acme
   │           ├── _defaults.yaml
   │           └── plat  # 'plat' represents the "Platform" OU (a.k.a tenant)
   │               ├── _defaults.yaml
   │               ├── dev
   │               │   ├── _defaults.yaml
   │               │   ├── us-east-2.yaml
   │               │   └── us-west-2.yaml
   │               ├── staging
   │               │   ├── _defaults.yaml
   │               │   ├── us-east-2.yaml
   │               │   └── us-west-2.yaml
   │               └── prod
   │                   ├── _defaults.yaml
   │                   ├── us-east-2.yaml
   │                   └── us-west-2.yaml
   │   # Centralized library of reusable components
   └── components
       └── terraform  # Terraform components (a.k.a. Terraform "root" modules)
           └── vpc
```

Add the following minimal configuration to `atmos.yaml` [CLI config file](/cli/configuration) :

```yaml title="atmos.yaml"
components:
  terraform:
    base_path: "components/terraform"

stacks:
  base_path: "stacks"
  name_pattern: "{tenant}-{environment}-{stage}"
  included_paths:
    # Tell Atmos to search for the top-level stack manifests in the `orgs` folder and its sub-folders
    - "orgs/**/*"
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

Add the following configuration for the base component `vpc/defaults` to the `stacks/catalog/vpc/defaults.yaml` manifest:

```yaml title="stacks/catalog/vpc/defaults.yaml"
components:
  terraform:
    vpc/defaults:
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
        ipv4_primary_cidr_block: 10.0.0.0/18
```

Configure the `vpc` Atmos component in the `stacks/orgs/acme/plat/prod/us-east-2.yaml` top-level stack. The `vpc` component inherits from
the `vpc/defaults` base component:

```yaml title="stacks/orgs/acme/plat/prod/us-east-2.yaml"
import:
  import:
    - orgs/acme/plat/prod/_defaults
    - mixins/region/us-east-2
    # Import the `vpc/defaults` component from the `catalog/vpc/defaults.yaml` manifest
    - catalog/vpc/defaults

components:
  terraform:
    # Atmos component `vpc`
    vpc:
      metadata:
        # Point to the Terraform component in `components/terraform/vpc`
        component: vpc
        # Inherit from the `vpc/defaults` Atmos base component
        # This is Single Inheritance: the Atmos component inherits from one base Atmos component
        inherits:
          - vpc/defaults
      # Define/override variables specific to this `vpc` component
      vars:
        name: my-vpc
        vpc_flow_logs_enabled: false
        ipv4_primary_cidr_block: 10.9.0.0/18
```

<br/>

## Benefits

The **Component Inheritance** pattern provides the following benefits:

- Allows creating very DRY, consistent, and reusable configurations that are built upon existing components

- Any Atmos component can serve as a building block for other Atmos components

## Related Patterns

- [Abstract Component](/design-patterns/abstract-component)
- [Multiple Component Instances](/design-patterns/multiple-component-instances)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)

## References

- [Component-Oriented Programming](/core-concepts/components/component-oriented-programming)
- [Component Inheritance](/core-concepts/components/inheritance)
