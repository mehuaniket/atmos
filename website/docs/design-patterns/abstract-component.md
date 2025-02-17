---
title: Abstract Component Atmos Design Pattern
sidebar_position: 9
sidebar_label: Abstract Component
description: Abstract Component Atmos Design Pattern
---

# Abstract Component

The **Abstract Component** Design Pattern describes the mechanism of deriving Atmos components from one or more **abstract** base components (
blueprints), allowing reusing the base components' configurations and prohibiting the abstract base component from being provisioned.

Atmos supports two types of components:

- `real` - is a ["concrete"](https://en.wikipedia.org/wiki/Concrete_class) component instance that can be provisioned

- `abstract` - a component configuration, which cannot be instantiated directly. The concept is borrowed
  from ["abstract base classes"](https://en.wikipedia.org/wiki/Abstract_type) of Object-Oriented Programming

The type of component is expressed in the `metadata.type` parameter of a given component configuration.

:::note
For more details, refer to:

- [Atmos Components](/core-concepts/components)
- [Glossary](/reference/glossary)

:::

## Applicability

Use the **Abstract Component** pattern when:

- You need to have reusable base components that serve as blueprints for the derived Atmos components

- You need to prevent the abstract base components from being provisioned

- You need to keep the configuration of all components [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

The following example shows the Atmos stack and component configurations to provision the `vpc` component into
a multi-account, multi-region environment. In the `catalog/vpc` folder, we have the `defaults.yaml` manifest that configures the base **abstract**
component `vpc/defaults` to be inherited by all the derived VPC components in all stacks. By being **abstract**, the base component `vpc/defaults`
is prohibited from being provisioned.

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
   │   # Centralized components configuration
   └── components
       └── terraform  # Terraform components (a.k.a Terraform "root" modules)
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

Add the following configuration for the abstract base component `vpc/defaults` to the `stacks/catalog/vpc/defaults.yaml` manifest:

```yaml title="stacks/catalog/vpc/defaults.yaml"
components:
  terraform:
    vpc/defaults:
      metadata:
        # Abstract components can't be provisioned, they just serve as base components (blueprints) for real components
        # `metadata.type` can be `abstract` or `real`
        # `real` is the default and can be omitted
        type: abstract
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
the `vpc/defaults` abstract base component:

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
        # `metadata.type` can be `abstract` or `real`
        # `real` is the default and can be omitted
        # Real components can be provisioned
        type: real
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

To provision the `vpc` component into the `plat-ue2-prod` top-level stack, execute the following command:

```shell
atmos terraform apply vpc -s plat-ue2-prod
```

If you try to execute the following commands to provision the `vpc/defaults` abstract base component:

```shell
atmos terraform plan vpc/defaults -s plat-ue2-prod
atmos terraform apply vpc/defaults -s plat-ue2-prod
```

the following error will be thrown:

```console
abstract component 'vpc/defaults' cannot be provisioned since it's explicitly prohibited from 
being deployed by 'metadata.type: abstract' attribute
```

## Benefits

The **Abstract Component** pattern provides the following benefits:

- Allows creating very DRY and reusable configurations that are built upon existing abstract base components (blueprints)

- Prevents the abstract base components from being provisioned

- The `metadata.type: abstract` attribute serves as a guard against accidentally deploying the components that are not meant to be deployed

## Related Patterns

- [Component Inheritance](/design-patterns/component-inheritance)
- [Multiple Component Instances](/design-patterns/multiple-component-instances)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)

## References

- [Component-Oriented Programming](/core-concepts/components/component-oriented-programming)
- [Component Inheritance](/core-concepts/components/inheritance)
