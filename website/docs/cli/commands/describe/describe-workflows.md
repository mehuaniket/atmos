---
title: atmos describe workflows
sidebar_label: workflows
sidebar_class_name: command
id: workflows
description: Use this command to show all configured Atmos workflows.
---

:::note Purpose
Use this command to show all configured Atmos workflows.
:::

## Usage

Execute the `describe workflows` command like this:

```shell
atmos describe workflows [options]
```

<br/>

:::tip
Run `atmos describe workflows --help` to see all the available options
:::

## Examples

```shell
atmos describe workflows
atmos describe workflows --output map
atmos describe workflows -o list
atmos describe workflows -o all
atmos describe workflows -o list --format json
atmos describe workflows -o all -f yaml
atmos describe workflows -f json
```

## Flags

| Flag       | Description                                                         | Alias | Required |
|:-----------|:--------------------------------------------------------------------|:------|:---------|
| `--format` | Specify the output format: `yaml` or `json` (`yaml` is default)     | `-f`  | no       |
| `--output` | Specify the output type: `list`, `map` or `all` (`list` is default) | `-o`  | no       |

<br/>

When the `--output list` flag is passed (default), the output of the command is a list of objects. Each object has the following schema:

- `file` - the workflow manifest file name
- `workflow` - the name of the workflow defined in the workflow manifest file

For example:

```shell
atmos describe workflows
atmos describe workflows -o list
```

```yaml
- file: compliance.yaml
  workflow: deploy/aws-config/global-collector
- file: compliance.yaml
  workflow: deploy/aws-config/superadmin
- file: compliance.yaml
  workflow: destroy/aws-config/global-collector
- file: compliance.yaml
  workflow: destroy/aws-config/superadmin
- file: datadog.yaml
  workflow: deploy/datadog-integration
- file: helpers.yaml
  workflow: save/docker-config-json
- file: networking.yaml
  workflow: apply-all-components
- file: networking.yaml
  workflow: plan-all-vpc-components
- file: networking.yaml
  workflow: plan-all-vpc-flow-logs-bucket-components
- file: vpc.yaml
  workflow: vpc-tgw-eks
```

<br/>

When the `--output map` flag is passed, the output of the command is a map of workflow manifests to the lists of workflows defined in each manifest.
For example:

```shell
atmos describe workflows -o map
```

```yaml
compliance.yaml:
  - deploy/aws-config/global-collector
  - deploy/aws-config/superadmin
  - destroy/aws-config/global-collector
  - destroy/aws-config/superadmin
datadog.yaml:
  - deploy/datadog-integration
helpers.yaml:
  - save/docker-config-json
networking.yaml:
  - apply-all-components
  - plan-all-vpc-components
  - plan-all-vpc-flow-logs-bucket-components
vpc.yaml:
  - vpc-tgw-eks
```

<br/>

When the `--output all` flag is passed, the output of the command is a map of workflow manifests to the maps of all workflow definitions. For example:

```shell
atmos describe workflows -o all
```

```yaml
networking.yaml:
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
```

<br/>

:::tip
Use the [atmos workflow](/cli/commands/workflow) CLI command to execute an Atmos workflow
:::
