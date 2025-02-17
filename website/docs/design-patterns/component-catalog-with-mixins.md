---
title: Component Catalog with Mixins Atmos Design Pattern
sidebar_position: 6
sidebar_label: Component Catalog with Mixins
description: Component Catalog with Mixins Atmos Design Pattern
---

# Component Catalog with Mixins

The **Component Catalog with Mixins** Design Pattern is a variation of the [Component Catalog](/design-patterns/component-catalog) pattern, with the
difference being that we first create parts of a component's configuration related to different environments (e.g. in `mixins` folder), then
assemble environment-specific manifests from the parts, and then import the environment-specific manifests themselves into the top-level stacks.

It's similar to how [Helm](https://helm.sh/) and [helmfile](https://helmfile.readthedocs.io/en/latest/#environment) handle environments.

The **Component Catalog with Mixins** Design Pattern prescribes the following:

- For a Terraform component, create a folder with the same name in `stacks/catalog` to make it symmetrical and easy to find.
  For example, the `stacks/catalog/vpc` folder should mirror the `components/terraform/vpc` folder.

- In the component's catalog folder, in the `mixins` sub-folder, add manifests with component configurations for specific environments (organizations,
  tenants, regions, accounts). For example:

  - `stacks/catalog/vpc/mixins/defaults.yaml` - component manifest with all the default values for the component (the defaults that can be reused
    across multiple environments)
  - `stacks/catalog/vpc/mixins/dev.yaml` - component manifest with the settings related to the `dev` account
  - `stacks/catalog/vpc/mixins/prod.yaml` - component manifest with the settings related to the `prod` account
  - `stacks/catalog/vpc/mixins/staging.yaml` - component manifest with the settings related to the `staging` account
  - `stacks/catalog/vpc/mixins/ue2.yaml` - component manifest with the settings for `us-east-2` region
  - `stacks/catalog/vpc/mixins/uw2.yaml` - component manifest with the settings for `us-west-2` region
  - `stacks/catalog/vpc/mixins/core.yaml` - component manifest with the settings related to the `core` tenant
  - `stacks/catalog/vpc/mixins/plat.yaml` - component manifest with the settings related to the `plat` tenant
  - `stacks/catalog/vpc/mixins/org1.yaml` - component manifest with the settings related to the `org1` organization
  - `stacks/catalog/vpc/mixins/org2.yaml` - component manifest with the settings related to the `org2` organization

- In the component's catalog folder, add manifests for specific environments by assembling the corresponding mixins together (using imports). For
  example:

  - `stacks/catalog/vpc/org1-plat-ue2-dev.yaml` - manifest for the `org1` organization, `plat` tenant, `ue2` region, `dev` account
  - `stacks/catalog/vpc/org1-plat-ue2-prod.yaml` - manifest for the `org1` organization, `plat` tenant, `ue2` region, `prod` account
  - `stacks/catalog/vpc/org1-plat-ue2-staging.yaml` - manifest for the `org1` organization, `plat` tenant, `ue2` region, `staging` account
  - `stacks/catalog/vpc/org1-plat-uw2-dev.yaml` - manifest for the `org1` organization, `plat` tenant, `uw2` region, `dev` account
  - `stacks/catalog/vpc/org1-plat-uw2-prod.yaml` - manifest for the `org1` organization, `plat` tenant, `uw2` region, `prod` account
  - `stacks/catalog/vpc/org1-plat-uw2-staging.yaml` - manifest for the `org1` organization, `plat` tenant, `uw2` region, `staging` account
  - `stacks/catalog/vpc/org2-plat-ue2-dev.yaml` - manifest for the `org2` organization, `plat` tenant, `ue2` region, `dev` account
  - `stacks/catalog/vpc/org2-plat-ue2-prod.yaml` - manifest for the `org2` organization, `plat` tenant, `ue2` region, `prod` account
  - `stacks/catalog/vpc/org2-plat-ue2-staging.yaml` - manifest for the `org2` organization, `plat` tenant, `ue2` region, `staging` account
  - `stacks/catalog/vpc/org2-plat-uw2-dev.yaml` - manifest for the `org2` organization, `plat` tenant, `uw2` region, `dev` account
  - `stacks/catalog/vpc/org2-plat-uw2-prod.yaml` - manifest for the `org2` organization, `plat` tenant, `uw2` region, `prod` account
  - `stacks/catalog/vpc/org2-plat-uw2-staging.yaml` - manifest for the `org2` organization, `plat` tenant, `uw2` region, `staging` account

- Import the environment manifests into the top-level stacks. For example:

  - import the `stacks/catalog/vpc/org1-plat-ue2-dev.yaml` manifest into the `stacks/orgs/org1/plat/dev/us-east-2.yaml` top-level stack
  - import the `stacks/catalog/vpc/org1-plat-ue2-prod.yaml` manifest into the `stacks/orgs/org1/plat/prod/us-east-2.yaml` top-level stack
  - import the `stacks/catalog/vpc/org1-plat-uw2-staging.yaml` manifest into the `stacks/orgs/org1/plat/staging/us-west-2.yaml` top-level stack
  - import the `stacks/catalog/vpc/org2-plat-ue2-dev.yaml` manifest into the `stacks/orgs/org2/plat/dev/us-east-2.yaml` top-level stack
  - etc.

## Applicability

Use the **Component Catalog** pattern when:

- You have many components that are provisioned in multiple stacks (many OUs, accounts, regions) with different configurations for each stack

- You need to make the component configurations reusable across different environments

- You want to keep the configurations [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

:::note
Having the environment-specific manifests in the component's catalog makes the most sense for multi-Org, multi-OU and/or
multi-region architectures, such that there will be multiple dev/staging/prod or region configurations, which get imported into multiple Org/OU
top-level stack manifests.
:::

## Example

The following example shows the Atmos stack and component configurations to provision the `vpc` component into
a multi-org, multi-tenant, multi-account, multi-region environment. The component's configuration for each organization, tenant, account and region
are defined as mixins in the component's catalog. The mixins then combined into the environment manifests, and the environment manifests are imported
into the top-level Atmos stacks.

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   └── catalog  # component-specific defaults
   │       └── vpc
   │           ├── mixins
   │           │   ├── defaults.yaml
   │           │   ├── dev.yaml
   │           │   ├── prod.yaml
   │           │   ├── staging.yaml
   │           │   ├── ue2.yaml
   │           │   ├── uw2.yaml
   │           │   ├── core.yaml
   │           │   ├── plat.yaml
   │           │   ├── org1.yaml
   │           │   └── org2.yaml
   │           ├── org1-plat-ue2-dev.yaml
   │           ├── org1-plat-ue2-prod.yaml
   │           ├── org1-plat-ue2-staging.yaml
   │           ├── org1-plat-uw2-dev.yaml
   │           ├── org1-plat-uw2-prod.yaml
   │           ├── org1-plat-uw2-staging.yaml
   │           ├── org2-plat-ue2-dev.yaml
   │           ├── org2-plat-ue2-prod.yaml
   │           ├── org2-plat-ue2-staging.yaml
   │           ├── org2-plat-uw2-dev.yaml
   │           ├── org2-plat-uw2-prod.yaml
   │           └── org2-plat-uw2-staging.yaml
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
  included_paths:
    # Tell Atmos to search for the top-level stack manifests in the `orgs` folder and its sub-folders
    - "orgs/**/*"
  excluded_paths:
    # Tell Atmos that all `_defaults.yaml` files are not top-level stack manifests
    - "**/_defaults.yaml"
  # If you are using multiple organizations (namespaces), use the following `name_pattern`:
  name_pattern: "{namespace}-{tenant}-{environment}-{stage}"
  # If you are using a single organization (namespace), use the following `name_pattern`:
  # name_pattern: "{tenant}-{environment}-{stage}"

schemas:
  jsonschema:
    base_path: "stacks/schemas/jsonschema"
  opa:
    base_path: "stacks/schemas/opa"
  atmos:
    manifest: "stacks/schemas/atmos/atmos-manifest/1.0/atmos-manifest.json"
```

Add the following default configuration to the `stacks/catalog/vpc/mixins/defaults.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/defaults.yaml"
components:
  terraform:
    vpc:
      metadata:
        # Point to the Terraform component in `components/terraform/vpc`
        component: vpc
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

Add the following configuration to the `stacks/catalog/vpc/mixins/ue2.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/ue2.yaml"
components:
  terraform:
    vpc:
      vars:
        availability_zones:
          - us-east-2a
          - us-east-2b
          - us-east-2c
```

Add the following configuration to the `stacks/catalog/vpc/mixins/uw2.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/uw2.yaml"
components:
  terraform:
    vpc:
      vars:
        availability_zones:
          - us-west-2a
          - us-west-2b
          - us-west-2c
```

Add the following configuration to the `stacks/catalog/vpc/mixins/dev.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/dev.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `ipv4_primary_cidr_block`, `max_subnet_count` and `vpc_flow_logs_enabled` from the defaults
        ipv4_primary_cidr_block: 10.7.0.0/18
        # In `dev`, use only 2 subnets
        max_subnet_count: 2
        # In `dev`, disable the VPC flow logs
        vpc_flow_logs_enabled: false
```

Add the following configuration to the `stacks/catalog/vpc/mixins/prod.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/prod.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `ipv4_primary_cidr_block`, `map_public_ip_on_launch` and `assign_generated_ipv6_cidr_block` from the defaults
        ipv4_primary_cidr_block: 10.8.0.0/18
        # In `prod`, don't map public IPs on launch
        map_public_ip_on_launch: false
        # In `prod`, use IPv6
        assign_generated_ipv6_cidr_block: true
```

Add the following configuration to the `stacks/catalog/vpc/mixins/staging.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/staging.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `ipv4_primary_cidr_block`, `max_subnet_count` and `map_public_ip_on_launch` from the defaults
        ipv4_primary_cidr_block: 10.9.0.0/18
        # In `staging`, use only 2 subnets
        max_subnet_count: 2
        # In `staging`, don't map public IPs on launch
        map_public_ip_on_launch: false
```

Add the following configuration to the `stacks/catalog/vpc/mixins/core.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/core.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `vpc_flow_logs_traffic_type` from the defaults
        # In `core` tenant, set VPC Flow Logs traffic type to `REJECT`
        vpc_flow_logs_traffic_type: "REJECT"
```

Add the following configuration to the `stacks/catalog/vpc/mixins/plat.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/plat.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `nat_eip_aws_shield_protection_enabled` from the defaults
        # In `plat` tenant, enable NAT EIP shield protection
        nat_eip_aws_shield_protection_enabled: true
```

Add the following configuration to the `stacks/catalog/vpc/mixins/org1.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/org1.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `subnet_type_tag_key` from the defaults
        subnet_type_tag_key: "org1/subnet/type"
```

Add the following configuration to the `stacks/catalog/vpc/mixins/org2.yaml` manifest:

```yaml title="stacks/catalog/vpc/mixins/org2.yaml"
components:
  terraform:
    vpc:
      vars:
        # Override `subnet_type_tag_key` from the defaults
        subnet_type_tag_key: "org2/subnet/type"
```

Assemble the `stacks/catalog/vpc/org1-plat-ue2-dev.yaml` environment manifest from the corresponding mixins:

```yaml title="stacks/catalog/vpc/org1-plat-ue2-dev.yaml"
import:
  # The imports are processed in the order they are defined.
  # The next imported manifest will override the configurations from the previously imported manifests
  - catalog/vpc/mixins/defaults
  - catalog/vpc/mixins/org1
  - catalog/vpc/mixins/plat
  - catalog/vpc/mixins/ue2
  - catalog/vpc/mixins/dev
```

Assemble the `stacks/catalog/vpc/org1-plat-ue2-prod.yaml` environment manifest from the corresponding mixins:

```yaml title="stacks/catalog/vpc/org1-plat-ue2-prod.yaml"
import:
  # The imports are processed in the order they are defined.
  # The next imported manifest will override the configurations from the previously imported manifests
  - catalog/vpc/mixins/defaults
  - catalog/vpc/mixins/org1
  - catalog/vpc/mixins/plat
  - catalog/vpc/mixins/ue2
  - catalog/vpc/mixins/prod
```

Assemble the `stacks/catalog/vpc/org1-plat-uw2-staging.yaml` environment manifest from the corresponding mixins:

```yaml title="stacks/catalog/vpc/org1-plat-uw2-staging.yaml"
import:
  # The imports are processed in the order they are defined.
  # The next imported manifest will override the configurations from the previously imported manifests
  - catalog/vpc/mixins/defaults
  - catalog/vpc/mixins/org1
  - catalog/vpc/mixins/plat
  - catalog/vpc/mixins/uw2
  - catalog/vpc/mixins/staging
```

Assemble the `stacks/catalog/vpc/org2-core-ue2-dev.yaml` environment manifest from the corresponding mixins:

```yaml title="stacks/catalog/vpc/org2-core-ue2-dev.yaml"
import:
  # The imports are processed in the order they are defined.
  # The next imported manifest will override the configurations from the previously imported manifests
  - catalog/vpc/mixins/defaults
  - catalog/vpc/mixins/org2
  - catalog/vpc/mixins/core
  - catalog/vpc/mixins/ue2
  - catalog/vpc/mixins/dev
```

Similarly, assemble the mixins for the other environments.

Import the `stacks/catalog/vpc/org1-plat-ue2-dev.yaml` environment manifest into the `stacks/orgs/org1/plat/dev/us-east-2.yaml` top-level stack:

```yaml title="stacks/orgs/org1/plat/dev/us-east-2.yaml"
import:
  - orgs/org1/plat/dev/_defaults
  - mixins/region/us-east-2
  - catalog/vpc/org1-plat-ue2-dev
```

Import the `stacks/catalog/vpc/org1-plat-ue2-prod.yaml` environment manifest into the `stacks/orgs/org1/plat/prod/us-east-2.yaml` top-level stack:

```yaml title="stacks/orgs/org1/plat/prod/us-east-2.yaml"
import:
  - orgs/org1/plat/prod/_defaults
  - mixins/region/us-east-2
  - catalog/vpc/org1-plat-ue2-prod
```

Import the `stacks/catalog/vpc/org1-plat-uw2-staging.yaml` environment manifest into the `stacks/orgs/org1/plat/staging/us-west-2.yaml` top-level
stack:

```yaml title="stacks/orgs/org1/plat/staging/us-west-2.yaml"
import:
  - orgs/org1/plat/staging/_defaults
  - mixins/region/us-west-2
  - catalog/vpc/org1-plat-uw2-staging
```

Import the `stacks/catalog/vpc/org2-core-ue2-dev.yaml` environment manifest into the `stacks/orgs/org2/core/dev/us-east-2.yaml` top-level
stack:

```yaml title="stacks/orgs/org2/core/dev/us-east-2.yaml"
import:
  - orgs/org2/core/dev/_defaults
  - mixins/region/us-east-2
  - catalog/vpc/org2-core-ue2-dev
```

Similarly, import the other environment mixins into the corresponding top-level stacks.

## Benefits

The **Component Catalog with Mixins** pattern provides the following benefits:

- Easy to see where the configuration for each environment is defined

- Easy to manage different variations of the configurations

- The defaults for the components are defined in just one place (in the catalog) making the entire
  configuration [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

- The defaults for the components are reusable across many environments by using hierarchical [imports](/core-concepts/stacks/imports)

## Limitations

The **Component Catalog with Mixins** pattern has the following limitations and drawbacks:

- The structure described by the pattern can be complex for basic infrastructures, e.g. for a very simple organizational structure (one organization
  and OU), and just a few components deployed into a few accounts and regions

:::note

To address the limitations of the **Component Catalog with Mixins** Design Pattern when you are provisioning a very basic infrastructure, consider
using the following patterns:

- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)
- [Component Catalog](/design-patterns/component-catalog)

:::

## Related Patterns

- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog Template](/design-patterns/component-catalog-template)
- [Component Inheritance](/design-patterns/component-inheritance)
- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)
- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)

## References

- [Catalogs](/core-concepts/stacks/catalogs)
- [Mixins](/core-concepts/stacks/mixins)
