---
title: Multiple Component Instances Atmos Design Pattern
sidebar_position: 10
sidebar_label: Multiple Component Instances
description: Multiple Component Instances Atmos Design Pattern
---

# Multiple Component Instances

The **Multiple Component Instances** Design Pattern describes how to provision many instances of a Terraform component in the same environment by
configuring multiple Atmos component instances.

An Atmos component can have any name that can be different from the Terraform component name. More than one instance of the same Terraform component
(with the same or different settings) can be provisioned into the same environment by defining multiple Atmos components that provide configuration
for the Terraform component.

For example, two different Atmos components `vpc/1` and `vpc/2` can provide configuration for the same Terraform component `vpc`.

## Applicability

Use the **Multiple Component Instances** pattern when:

- You need to provision multiple instances of a Terraform component in the same environment (organization, OU, account, region)

## Example

The following example shows the Atmos stack and component configurations to provision two Atmos components (`vpc/1` and `vpc/2`) that use the
same Terraform component `vpc`. The `vpc/1` and `vpc/2` components are provisioned in the same stack. In the `catalog/vpc` folder, we have
the `defaults.yaml` manifest that configures the base abstract component `vpc/defaults` to be inherited by the derived VPC components `vpc/1`
and `vpc/2`.

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

Add the following default configuration to the `stacks/catalog/vpc/defaults.yaml` manifest:

```yaml title="stacks/catalog/vpc/defaults.yaml"
components:
  terraform:
    vpc/defaults:
      metadata:
        # Abstract components can't be provisioned and serve as base components (blueprints) for real components
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

Configure multiple `vpc` component instances in the `stacks/orgs/acme/plat/prod/us-east-2.yaml` top-level stack:

```yaml title="stacks/orgs/acme/plat/prod/us-east-2.yaml"
import:
  import:
    - orgs/acme/plat/prod/_defaults
    - mixins/region/us-east-2
    # Import the defaults for all VPC components
    - catalog/vpc/defaults

components:
  terraform:
    # Atmos component `vpc/1`
    vpc/1:
      metadata:
        # Point to the Terraform component in `components/terraform/vpc`
        component: vpc
        # Inherit the defaults for all VPC components
        inherits:
          - vpc/defaults
      # Define/override variables specific to this `vpc/1` component
      vars:
        name: vpc-1
        ipv4_primary_cidr_block: 10.9.0.0/18

    # Atmos component `vpc/2`
    vpc/2:
      metadata:
        # Point to the Terraform component in `components/terraform/vpc`
        component: vpc
        # Inherit the defaults for all VPC components
        inherits:
          - vpc/defaults
      # Define/override variables specific to this `vpc/2` component
      vars:
        name: vpc-2
        ipv4_primary_cidr_block: 10.10.0.0/18
        map_public_ip_on_launch: false
```

<br/>

To provision the components in the stack, execute the following commands:

```shell
atmos terraform apply vpc/1 -s plat-ue2-prod
atmos terraform apply vpc/2 -s plat-ue2-prod
```

The provisioned VPCs will have the following names:

```console
acme-plat-ue2-prod-vpc-1
acme-plat-ue2-prod-vpc-2
```

The names are generated from the context using the following template:

```console
{namespace}-{tenant}-{environment}-{stage}-{name}
```

## Benefits

The **Multiple Component Instances** pattern provides the following benefits:

- Separation of the code (Terraform component) from the configuration (Atmos components)

- The same Terraform code is reused by multiple Atmos component instances with different configurations

- The defaults for the components are defined in just one place making the entire
  configuration [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Related Patterns

- [Component Inheritance](/design-patterns/component-inheritance)
- [Abstract Component](/design-patterns/abstract-component)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)
