---
title: Component Catalog Template Atmos Design Pattern
sidebar_position: 7
sidebar_label: Component Catalog Template
description: Component Catalog Template Atmos Design Pattern
---

# Component Catalog Template

The **Component Catalog Template** Design Pattern is used when you have an unbounded number of a component's instances provisioned in one environment
(the same organization, OU/tenant, account and region). New instances with different settings can be configured and provisioned anytime. The old
instances must be kept unchanged and never destroyed.

This is achieved by using [`Go` Templates in Imports](/core-concepts/stacks/imports#go-templates-in-imports) and
[Hierarchical Imports with Context](/core-concepts/stacks/imports#hierarchical-imports-with-context).

The **Component Catalog Template** pattern recommends the following:

- In the component's catalog folder, create a [`Go` template](https://pkg.go.dev/text/template) manifest with all the configurations for the
  component (refer to [`Go` Templates in Imports](/core-concepts/stacks/imports#go-templates-in-imports) for more details)

- Import the `Go` template manifest into a top-level stack many times to configure the component's instances
  using [Hierarchical Imports with Context](/core-concepts/stacks/imports#hierarchical-imports-with-context) and providing different template values
  for each import

## Applicability

Use the **Component Catalog Template** pattern when:

- You have an unbounded number of a component's instances provisioned in one environment (the same organization, OU/tenant, account and region)

- New instances of the component with different settings can be configured and provisioned anytime

- The old instances of the component must be kept unchanged and never destroyed

- You want to keep the configurations [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

## Example

Suppose that we have an EKS cluster provisioned in one of the accounts and regions.
The cluster is running many different applications, each one requires
an [IAM role for service accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) (IRSA) with permissions
to access various AWS resources.

The Development team can create a new application anytime, and we need to provision a new IRSA in the EKS cluster.
We'll use the **Component Catalog Template** Design Pattern to configure the IAM roles with different settings for each application.

```console
   │   # Centralized stacks configuration (stack manifests)
   ├── stacks
   │   └── catalog  # component-specific defaults
   │       └── eks
   │           └── iam-role
   │               └── defaults.tmpl
   │   # Centralized components configuration
   └── components
       └── terraform  # Terraform components (a.k.a Terraform "root" modules)
           └── eks
               └── iam-role
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

Add the following configuration to the `stacks/catalog/eks/iam-role/defaults.tmpl` `Go` template manifest:

```yaml title="stacks/catalog/eks/iam-role/defaults.tmpl"
components:
  terraform:
    eks/iam-role/{{ .app_name }}:
      metadata:
        # Point to the Terraform component
        component: eks/iam-role
      vars:
        enabled: true
        tags:
          Service: { { .app_name } }
        service_account_name: { { .service_account_name } }
        service_account_namespace: { { .service_account_namespace } }
        # Example of using the Sprig functions in `Go` templates.
        # Refer to https://masterminds.github.io/sprig for more details.
        { { if hasKey . "iam_managed_policy_arns" } }
        iam_managed_policy_arns:
          { { range $i, $iam_managed_policy_arn := .iam_managed_policy_arns } }
          - '{{ $iam_managed_policy_arn }}'
          { { end } }
        { { - end } }
```

Import the `stacks/catalog/eks/iam-role/defaults.tmpl` manifest template into a top-level stack,
for example `stacks/orgs/acme/plat/prod/us-east-2.yaml`, and provide the configuration for each application in the `context` object:

```yaml title="stacks/orgs/acme/plat/prod/us-east-2.yaml"
import:
  - orgs/acme/plat/prod/_defaults
  - mixins/region/us-east-2

  # This import will dynamically generate a new Atmos component `eks/iam-role/admin-ui`
  - path: catalog/eks/iam-role/defaults.tmpl
    context:
      app_name: "admin-ui"
      service_account_name: "admin-ui"
      service_account_namespace: "admin"
      iam_managed_policy_arns: [ "<arn1>", "<arn2>" ]

  # This import will dynamically generate a new Atmos component `eks/iam-role/auth`
  - path: catalog/eks/iam-role/defaults.tmpl
    context:
      app_name: "auth"
      service_account_name: "auth"
      service_account_namespace: "auth"
      iam_managed_policy_arns: [ "<arn3>" ]

  # This import will dynamically generate a new Atmos component `eks/iam-role/payment-processing`
  - path: catalog/eks/iam-role/defaults.tmpl
    context:
      app_name: "payment-processing"
      service_account_name: "payment-processing"
      service_account_namespace: "payments"
      iam_managed_policy_arns: [ "<arn4>", "<arn5>" ]

  # Add new application configurations here
```

To provision the Atmos components in the stack, execute the following commands:

```shell
atmos terraform apply eks/iam-role/admin-ui --stack plat-ue2-prod
atmos terraform apply eks/iam-role/auth --stack plat-ue2-prod
atmos terraform apply eks/iam-role/payment-processing --stack plat-ue2-prod
```

## Benefits

The **Component Catalog Template** pattern provides the following benefits:

- All settings for a component are defined in just one place (in the component's template) making the entire
  configuration [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)

- Many instances of the component can be provisioned without repeating all the configuration values

- New Atmos components are generated dynamically

## Limitations

The **Component Catalog Template** pattern has the following limitations and drawbacks:

- Since new Atmos components are generated dynamically, sometimes it's not easy to know the names of the Atmos components that need to be provisioned
  without looking at the `Go` template and figuring out all the Atmos component names

- With great power comes great responsibility. The Go templating engine leads to infinite possibilities, which makes the stack configs more
  challenging to maintain and reduces consistency. Try to leverage the inheritance model as much as possible, over templated stack configs.

:::note

To address the limitations of the **Component Catalog Template** Design Pattern, consider using the following patterns:

- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)

:::

## Related Patterns

- [Component Catalog](/design-patterns/component-catalog)
- [Component Catalog with Mixins](/design-patterns/component-catalog-with-mixins)
- [Component Inheritance](/design-patterns/component-inheritance)
- [Inline Component Configuration](/design-patterns/inline-component-configuration)
- [Inline Component Customization](/design-patterns/inline-component-customization)
- [Organizational Structure Configuration](/design-patterns/organizational-structure-configuration)

## References

- [Catalogs](/core-concepts/stacks/catalogs)
