---
title: atmos workflow
sidebar_label: workflow
sidebar_class_name: command
description: Use this command to perform sequential execution of `atmos` and `shell` commands defined as workflow steps.
---

:::note Purpose
Use this command to perform sequential execution of `atmos` and `shell` commands defined as workflow steps.
:::

<br/>

An Atmos workflow is a series of steps that are run in order to achieve some outcome. Every workflow has a name and is easily executed from the
command line by calling `atmos workflow`. Use workflows to orchestrate any number of commands. Workflows can call any `atmos` subcommand (including
[Atmos Custom Commands](/core-concepts/custom-commands)), shell commands, and have access to the stack configurations.

:::note
You can use [Atmos Custom Commands](/core-concepts/custom-commands) in [Atmos Workflows](/core-concepts/workflows),
and [Atmos Workflows](/core-concepts/workflows)
in [Atmos Custom Commands](/core-concepts/custom-commands)
:::

## Usage

Execute the `atmos workflow` command like this:

```shell
atmos workflow <workflow_name> --file <workflow_file> [options]
```

### Examples

```shell
atmos workflow
atmos workflow plan-all-vpc-components --file networking
atmos workflow apply-all-components -f networking --dry-run
atmos workflow test-1 -f workflow1 --from-step step2
```

<br/>

:::tip
Run `atmos workflow --help` to see all the available options
:::

<br/>

Alternatively, just run `atmos workflow` to start an interactive UI to view, search and execute the configured Atmos workflows:

```shell
atmos workflow
```

<br/>

![`atmos workflow` CLI command 1](/img/cli/workflow/atmos-workflow-command-1.png)

- Use the `right/left` arrow keys to navigate between the "Workflow Manifests", "Workflows" and the selected workflow views

- Use the `up/down` arrow keys (or the mouse wheel) to select a workflow manifest and a workflow to execute

- Use the `/` key to filter/search for the workflow manifests and workflows in the corresponding views

- Press `Enter` to execute the selected workflow from the selected workflow manifest

<br/>

![`atmos workflow` CLI command 2](/img/cli/workflow/atmos-workflow-command-2.png)

## Arguments

| Argument         | Description   | Required |
|:-----------------|:--------------|:---------|
| `workflow_name ` | Workflow name | yes      |

## Flags

| Flag          | Description                                                                                   | Alias | Required |
|:--------------|:----------------------------------------------------------------------------------------------|:------|:---------|
| `--file`      | File name where the workflow is defined                                                       | `-f`  | yes      |
| `--stack`     | Atmos stack<br/>(if provided, will override stacks defined in the workflow or workflow steps) | `-s`  | no       |
| `--from-step` | Start the workflow from the named step                                                        |       | no       |
| `--dry-run`   | Dry run. Print information about the executed workflow steps without executing them           |       | no       |
