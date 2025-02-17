---
title: Partial Stack Configuration Atmos Design Pattern
sidebar_position: 12
sidebar_label: Partial Stack Configuration
description: Partial Stack Configuration Atmos Design Pattern
---

# Partial Stack Configuration

The **Partial Stack Configuration** Design Pattern describes the mechanism of splitting an Atmos top-level stack's configuration across many Atmos
stack manifests to manage and modify them separately and independently.

Each partial top-level stack manifest imports or configures a set of related Atmos components. Each Atmos component belongs to just one of the partial
top-level stack manifests. The pattern helps to group all components by category or function and to make each partial stack manifest smaller and
easier to manage.

## Applicability

Use the **Partial Stack Configuration** pattern when:

- You have top-level stacks with complex configurations. Some parts of the configurations must be managed and modified independent of the other
  parts, possibly by different people or teams

- You need to group the components in a top-level stack by category or function

- You want to keep the configuration easy to manage and [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

In the following structure, we have many Terraform components (Terraform root modules) in the `components/terraform` folder.

In the `stacks/catalog` folder, we define the defaults for each component using the [Component Catalog](/design-patterns/component-catalog) Atmos
Design Pattern.

In the `stacks/orgs/acme/plat/dev/us-east-2` folder, we split the top-level stack manifest into the following parts by category:

- `load-balancers.yaml`
- `data.yaml`
- `dns.yaml`
- `logs.yaml`
- `notifications.yaml`
- `firewalls.yaml`
- `networking.yaml`
- `eks.yaml`

<br/>

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   ├── catalog  # component-specific defaults
   │   │   ├── alb
   │   │   │   └── defaults.yaml
   │   │   ├── aurora-postgres
   │   │   │   └── defaults.yaml
   │   │   ├── dns
   │   │   │   └── defaults.yaml
   │   │   ├── eks
   │   │   │   └── defaults.yaml
   │   │   ├── efs
   │   │   │   └── defaults.yaml
   │   │   ├── msk
   │   │   │   └── defaults.yaml
   │   │   ├── ses
   │   │   │   └── defaults.yaml
   │   │   ├── sns-topic
   │   │   │   └── defaults.yaml
   │   │   ├── network-firewall
   │   │   │   └── defaults.yaml
   │   │   ├── network-firewall-logs-bucket
   │   │   │   └── defaults.yaml
   │   │   ├── waf
   │   │   │   └── defaults.yaml
   │   │   ├── vpc
   │   │   │   └── defaults.yaml
   │   │   └── vpc-flow-logs-bucket
   │   │       └── defaults.yaml
   │   ├── mixins
   │   │   ├── tenant  # tenant-specific defaults
   │   │   │   └── plat.yaml
   │   │   ├── region  # region-specific defaults
   │   │   │   ├── us-east-2.yaml
   │   │   │   └── us-west-2.yaml
   │   │   └── stage  # stage-specific defaults
   │   │       ├── dev.yaml
   │   │       ├── prod.yaml
   │   │       └── staging.yaml
   │   └── orgs  # Organizations
   │       └── acme
   │           ├── _defaults.yaml
   │           └── plat  # 'plat' represents the "Platform" OU (a.k.a tenant)
   │               ├── _defaults.yaml
   │               └── dev  # 'dev' account
   │                  ├── _defaults.yaml
   │                  ├── # Split the top-level stack 'plat-ue2-dev' by category of components
   │                  └── us-east-2
   │                      ├── load-balancers.yaml
   │                      ├── data.yaml
   │                      ├── dns.yaml
   │                      ├── logs.yaml
   │                      ├── notifications.yaml
   │                      ├── firewalls.yaml
   │                      ├── networking.yaml
   │                      └── eks.yaml
   │   # Centralized components configuration
   └── components
       └── terraform  # Terraform components (a.k.a Terraform "root" modules)
           ├── alb
           ├── aurora-postgres
           ├── dns
           ├── eks
           ├── efs
           ├── msk
           ├── ses
           ├── sns-topic
           ├── network-firewall
           ├── network-firewall-logs-bucket
           ├── waf
           ├── vpc
           └── vpc-flow-logs-bucket
```

Note that the partial stack manifests are parts of the same top-level Atmos stack `plat-ue2-dev` since they all import the same context variables
`namespace`, `tenant`, `environment` and `stage`. A top-level Atmos stack is defined by the context variables, not by the file names or locations
in the filesystem (file names can be anything, they are for people to better organize the entire configuration).

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

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/load-balancers.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/load-balancers.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/alb/defaults
  # Import other Load Balancer components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/data.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/data.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/aurora-postgres/defaults
  - catalog/msk/defaults
  - catalog/efs/defaults
  # Import other Data components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/dns.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/dns.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/dns/defaults
  # Import other DNS components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/logs.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/logs.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/network-firewall-logs-bucket/defaults
  - catalog/vpc-flow-logs-bucket/defaults
  # Import other Logs components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/notifications.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/notifications.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/ses/defaults
  - catalog/sns-topic/defaults
  # Import other Notification components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/firewalls.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/firewalls.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/network-firewall/defaults
  - catalog/waf/defaults
  # Import other Firewall components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/networking.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/networking.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/vpc/defaults
  # Import other Networking components
```

Add the following configuration to the `stacks/orgs/acme/plat/dev/us-east-2/eks.yaml` partial stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2/eks.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the related component manifests into this partial stack manifest
  - catalog/eks/defaults
```

## Benefits

The **Partial Stack Configuration** pattern provides the following benefits:

- Allows defining Atmos stacks with complex configurations by splitting the configurations into smaller manifests and by grouping the components by
  category or function

- Makes the configurations easier to understand

- Allows creating and modifying the partial stack manifests independently, possibly by different teams

## Limitations

The **Partial Stack Configuration** pattern has the following limitations and drawbacks:

- The structure described by the pattern can become big and complex in a production-ready enterprise-grade infrastructure

- In the example above, we showed how to split just one top-level stack manifest (one organization, one OU/tenant, one account, one region) into
  smaller parts and import the related components. This is useful and not complicated with one or a few top-level stacks, but the configuration
  can become too complex when we need to do the same for all organizations, OUs/tenants, accounts and regions. We'll have to repeat the same
  filesystem layout many times and import the same components into many partial stack manifests

:::note

To address the limitations of the **Partial Stack Configuration** Design Pattern, consider
the [Layered Stack Configuration](/design-patterns/layered-stack-configuration) Design Pattern

:::

## Related Patterns

- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)
- [Layered Stack Configuration](/design-patterns/layered-stack-configuration)
- [Component Overrides](/design-patterns/component-overrides)
- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)

## References

- [Catalogs](/core-concepts/stacks/catalogs)
- [Mixins](/core-concepts/stacks/mixins)
