---
title: Create Workflows
sidebar_position: 8
sidebar_label: Create Workflows
---

Atmos workflows are a way of combining multiple commands into executable units of work.

<br/>

:::tip
Refer to [Atmos Workflows](/core-concepts/workflows) for more information about configuring workflows
:::

<br/>

:::note
You can use [Atmos Custom Commands](/core-concepts/custom-commands) in [Atmos Workflows](/core-concepts/workflows),
and [Atmos Workflows](/core-concepts/workflows)
in [Atmos Custom Commands](/core-concepts/custom-commands)
:::

<br/>

To define workflows, add the following configurations:

- In `atmos.yaml`, add the `workflows` section and configure the base path to the workflows:

```yaml
workflows:
  # Can also be set using 'ATMOS_WORKFLOWS_BASE_PATH' ENV var, or '--workflows-dir' command-line arguments
  # Supports both absolute and relative paths
  base_path: "stacks/workflows"
```

- Add workflow manifests in the `stacks/workflows` folder. In this Quick Start example, we will define Atmos workflows in the `networking.yaml`
  workflow manifest:

```console
   │   # Centralized stacks configuration
   ├── stacks
   │   └── workflows
   │       ├── networking.yaml
```

- Add the following Atmos workflows to the `stacks/workflows/networking.yaml` file:

```yaml
workflows:

  plan-all-vpc-flow-logs-bucket-components:
    description: |
      Run 'terraform plan' on all 'vpc-flow-logs-bucket' components in all stacks
    steps:
      - command: terraform plan vpc-flow-logs-bucket -s plat-ue2-dev
      - command: terraform plan vpc-flow-logs-bucket -s plat-uw2-dev
      - command: terraform plan vpc-flow-logs-bucket -s plat-ue2-staging
      - command: terraform plan vpc-flow-logs-bucket -s plat-uw2-staging
      - command: terraform plan vpc-flow-logs-bucket -s plat-ue2-prod
      - command: terraform plan vpc-flow-logs-bucket -s plat-uw2-prod

  plan-all-vpc-components:
    description: |
      Run 'terraform plan' on all 'vpc' components in all stacks
    steps:
      - command: terraform plan vpc -s plat-ue2-dev
      - command: terraform plan vpc -s plat-uw2-dev
      - command: terraform plan vpc -s plat-ue2-staging
      - command: terraform plan vpc -s plat-uw2-staging
      - command: terraform plan vpc -s plat-ue2-prod
      - command: terraform plan vpc -s plat-uw2-prod

  apply-all-components:
    description: |
      Run 'terraform apply' on all components in all stacks
    steps:
      - command: terraform apply vpc-flow-logs-bucket -s plat-ue2-dev -auto-approve
      - command: terraform apply vpc -s plat-ue2-dev -auto-approve
      - command: terraform apply vpc-flow-logs-bucket -s plat-uw2-dev -auto-approve
      - command: terraform apply vpc -s plat-uw2-dev -auto-approve
      - command: terraform apply vpc-flow-logs-bucket -s plat-ue2-staging -auto-approve
      - command: terraform apply vpc -s plat-ue2-staging -auto-approve
      - command: terraform apply vpc-flow-logs-bucket -s plat-uw2-staging -auto-approve
      - command: terraform apply vpc -s plat-uw2-staging -auto-approve
      - command: terraform apply vpc-flow-logs-bucket -s plat-ue2-prod -auto-approve
      - command: terraform apply vpc -s plat-ue2-prod -auto-approve
      - command: terraform apply vpc-flow-logs-bucket -s plat-uw2-prod -auto-approve
      - command: terraform apply vpc -s plat-uw2-prod -auto-approve
```

- Run the following Atmos commands to execute the workflows:

```shell
# Execute the workflow `plan-all-vpc-flow-logs-bucket-components` from the workflow manifest `networking`
atmos workflow plan-all-vpc-flow-logs-bucket-components -f networking

# Execute the workflow `plan-all-vpc-components` from the workflow manifest `networking`
atmos workflow plan-all-vpc-components -f networking

# Execute the workflow `apply-all-components` from the workflow manifest `networking`
atmos workflow apply-all-components -f networking
```

<br/>

:::tip
Refer to [atmos workflow](/cli/commands/workflow) for more information on the `atmos workflow` CLI command
:::

<br/>

The `atmos workflow` CLI command supports the `--dry-run` flag. If passed, the command will just print information about the executed workflow steps
without executing them. For example:

<br/>

```shell
atmos workflow plan-all-vpc-components -f networking --dry-run
```

```console
Executing the workflow 'plan-all-vpc-components' from 'stacks/workflows/networking.yaml'

Executing workflow step: terraform plan vpc -s plat-ue2-dev
Executing workflow step: terraform plan vpc -s plat-uw2-dev
Executing workflow step: terraform plan vpc -s plat-ue2-staging
Executing workflow step: terraform plan vpc -s plat-uw2-staging
Executing workflow step: terraform plan vpc -s plat-ue2-prod
Executing workflow step: terraform plan vpc -s plat-uw2-prod
```
