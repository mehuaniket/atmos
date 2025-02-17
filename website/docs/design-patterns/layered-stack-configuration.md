---
title: Layered Stack Configuration Atmos Design Pattern
sidebar_position: 13
sidebar_label: Layered Stack Configuration
description: Layered Stack Configuration Atmos Design Pattern
---

# Layered Stack Configuration

The **Layered Stack Configuration** Design Pattern describes the mechanism of grouping Atmos components by category or function,
adding the groups of components to layers, and importing the layers into top-level Atmos stacks.

Each layer imports or configures a set of related Atmos components. Each Atmos component belongs to just one layer.
Each layer can be managed separately, possibly by different teams.

<br/>

:::note
The **Layered Stack Configuration** Design Pattern works around the limitations of
the [Partial Stack Configuration](/design-patterns/partial-stack-configuration) pattern. Instead of splitting the top-level Atmos stacks into parts,
the **Layered Stack Configuration** pattern adds separate layers to group the related Atmos components by category, and then import the layer
manifests into the top-level Atmos stacks.
:::

## Applicability

Use the **Layered Stack Configuration** pattern when:

- You have many Atmos components, and you need to group the components by category or function

- You want to split the components into layers. Each layer should be managed and modified independent of the other layers, possibly by different
  people or teams

- You want to keep the configuration easy to manage and [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

In the following structure, we have various Terraform components (Terraform root modules) in the `components/terraform` folder.

In the `stacks/catalog` folder, we define the defaults for each component using the [Component Catalog](/design-patterns/component-catalog) Atmos
Design Pattern.

In the `stacks/layers` folder, we define layers (groups of components), and import the related components into the layer manifests:

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
   │   ├── layers  # grouping of components by category/function
   │   │   ├── load-balancers.yaml
   │   │   ├── data.yaml
   │   │   ├── dns.yaml
   │   │   ├── logs.yaml
   │   │   ├── notifications.yaml
   │   │   ├── firewalls.yaml
   │   │   ├── networking.yaml
   │   │   └── eks.yaml
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
   │               │   ├── global-region.yaml
   │               │   ├── us-east-2.yaml
   │               │   └── us-west-2.yaml
   │               ├── staging
   │               │   ├── _defaults.yaml
   │               │   ├── global-region.yaml
   │               │   ├── us-east-2.yaml
   │               │   └── us-west-2.yaml
   │               └── prod
   │                   ├── _defaults.yaml
   │                   ├── global-region.yaml
   │                   ├── us-east-2.yaml
   │                   └── us-west-2.yaml
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

Add the following configuration to the `stacks/layers/load-balancers.yaml` layer manifest:

```yaml title="stacks/layers/load-balancers.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/alb/defaults
  # Import other Load Balancer components
```

Add the following configuration to the `stacks/layers/data.yaml` layer manifest:

```yaml title="stacks/layers/data.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/aurora-postgres/defaults
  - catalog/msk/defaults
  - catalog/efs/defaults
  # Import other Data components
```

Add the following configuration to the `stacks/layers/dns.yaml` layer manifest:

```yaml title="stacks/layers/dns.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/dns/defaults
  # Import other DNS components
```

Add the following configuration to the `stacks/layers/logs.yaml` layer manifest:

```yaml title="stacks/layers/logs.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/network-firewall-logs-bucket/defaults
  - catalog/vpc-flow-logs-bucket/defaults
  # Import other Logs components
```

Add the following configuration to the `stacks/layers/notifications.yaml` layer manifest:

```yaml title="stacks/layers/notifications.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/ses/defaults
  - catalog/sns-topic/defaults
  # Import other Notification components
```

Add the following configuration to the `stacks/layers/firewalls.yaml` layer manifest:

```yaml title="stacks/layers/firewalls.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/network-firewall/defaults
  - catalog/waf/defaults
  # Import other Firewall components
```

Add the following configuration to the `stacks/layers/networking.yaml` layer manifest:

```yaml title="stacks/layers/networking.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/vpc/defaults
  # Import other Networking components
```

Add the following configuration to the `stacks/layers/eks.yaml` layer manifest:

```yaml title="stacks/layers/eks.yaml"
import:
  # Import the related component manifests into this layer manifest
  - catalog/eks/defaults
```

Import the required layers into the `stacks/orgs/acme/plat/dev/us-east-2.yaml` top-level stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-east-2.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-east-2
  # Import the layers (groups of components)
  - layers/load-balancers
  - layers/data
  - layers/dns
  - layers/logs
  - layers/notifications
  - layers/firewalls
  - layers/networking
  - layers/eks
```

Import the required layers into the `stacks/orgs/acme/plat/dev/us-west-2.yaml` top-level stack manifest:

```yaml title="stacks/orgs/acme/plat/dev/us-west-2.yaml"
import:
  # The `orgs/acme/plat/dev/_defaults` and `mixins/region/us-west-2` manifests 
  # define the top-level Atmos stack `plat-uw2-dev`
  - orgs/acme/plat/dev/_defaults
  - mixins/region/us-west-2
  # Import the layers (groups of components)
  - layers/load-balancers
  - layers/data
  - layers/dns
  - layers/logs
  - layers/notifications
  - layers/firewalls
  - layers/networking
  - layers/eks
```

Import the required layers into the `stacks/orgs/acme/plat/prod/us-east-2.yaml` top-level stack manifest:

```yaml title="stacks/orgs/acme/plat/prod/us-east-2.yaml"
import:
  # The `orgs/acme/plat/prod/_defaults` and `mixins/region/us-east-2` manifests 
  # define the top-level Atmos stack `plat-ue2-prod`
  - orgs/acme/plat/prod/_defaults
  - mixins/region/us-east-2
  # Import the layers (groups of components)
  - layers/load-balancers
  - layers/data
  - layers/dns
  - layers/logs
  - layers/notifications
  - layers/firewalls
  - layers/networking
  - layers/eks
```

<br/>

Similarly, import the required layers into the other top-level stacks for the other organizations, OUs/tenants, accounts and regions.
Make sure to import only the layers that define the component that need to be provisioned in the stacks.

## Benefits

The **Layered Stack Configuration** pattern provides the following benefits:

- Allows to group Atmos components by category or function

  people or teams. Furthermore, controls like
  GitHub's [`CODEOWNERS`](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
  can be leveraged so that specific teams or individuals must review changes to these files.

- Allows importing only the required layers into the top-level stacks (only the groups of components that need to be provisioned in the stacks)

- Makes the configurations easier to understand

## Related Patterns

- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)
- [Partial Stack Configuration](/design-patterns/partial-stack-configuration)
- [Component Overrides](/design-patterns/component-overrides)

## References

- [Catalogs](/core-concepts/stacks/catalogs)
- [Mixins](/core-concepts/stacks/mixins)
