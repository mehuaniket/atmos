---
title: Describe Stacks
sidebar_position: 7
sidebar_label: Describing
id: describing
description: Describe stacks to view the fully deep-merged configuration
---

Describing stacks is helpful to understand what the final, fully computed and deep-merged configuration of a stack will look like. Use this to slice
and dice the Stack configuration to show different information about stacks and component.

For example, if we wanted to understand what the final configuration looks like for the "production" stack, we could do that by calling
the [`atmos describe stacks`](/cli/commands/describe/stacks) command to view the YAML output.

The output can be written to a file by passing the `--file` command-line flag to `atmos` or even formatted as YAML or JSON by using `--format`
command-line flag.

Since the output of a Stack might be overwhelming, and we're only interested in some particular section of the configuration, the output can be
filtered using flags to narrow the output by `stack`, `component-types`, `components`, and `sections`. The component sections can be further filtered
by `atmos_component`, `atmos_stack`, `atmos_stack_file`, `backend`, `backend_type`, `command`, `component`, `env`, `inheritance`, `metadata`,
`overrides`, `remote_state_backend`, `remote_state_backend_type`, `settings`, `vars`, `workspace`.

For example:

```shell
atmos descrive stacks
```

```yaml
plat-ue2-dev:
  components:
    terraform:
      vpc:
        backend: {}
        backend_type: s3
        command: terraform
        component: vpc
        env: {}
        inheritance: []
        metadata:
          component: vpc
        overrides: {}
        remote_state_backend: {}
        remote_state_backend_type: ""
        settings:
          validation:
            check-vpc-component-config-with-opa-policy:
              description: Check 'vpc' component configuration using OPA policy
              disabled: false
              module_paths:
                - catalog/constants
              schema_path: vpc/validate-vpc-component.rego
              schema_type: opa
              timeout: 10
            validate-vpc-component-with-jsonschema:
              description: Validate 'vpc' component variables using JSON Schema
              schema_path: vpc/validate-vpc-component.json
              schema_type: jsonschema
        vars:
          availability_zones:
            - us-east-2a
            - us-east-2b
            - us-east-2c
          dns_hostnames_enabled: true
          enabled: true
          environment: ue2
          map_public_ip_on_launch: true
          max_subnet_count: 3
          name: common
          namespace: acme
          nat_gateway_enabled: true
          nat_instance_enabled: false
          region: us-east-2
          stage: dev
          tenant: plat
          vpc_flow_logs_enabled: true
          vpc_flow_logs_log_destination_type: s3
          vpc_flow_logs_traffic_type: ALL
        workspace: plat-ue2-dev
      vpc-flow-logs-bucket:
        backend: {}
        backend_type: s3
        command: terraform
        component: vpc-flow-logs-bucket
        env: {}
        inheritance: []
        metadata:
          component: vpc-flow-logs-bucket
        overrides: {}
        remote_state_backend: {}
        remote_state_backend_type: ""
        settings: {}
        vars:
          enabled: true
          environment: ue2
          force_destroy: true
          lifecycle_rule_enabled: false
          name: vpc-flow-logs
          namespace: acme
          region: us-east-2
          stage: dev
          tenant: plat
          traffic_type: ALL
        workspace: plat-ue2-dev

# Other stacks here
```

```shell
atmos descrive stacks --components vpc --sections metadata --format json
```

```json
{
  "plat-ue2-dev": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  },
  "plat-ue2-prod": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  },
  "plat-ue2-staging": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  },
  "plat-uw2-dev": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  },
  "plat-uw2-prod": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  },
  "plat-uw2-staging": {
    "components": {
      "terraform": {
        "vpc": {
          "metadata": {
            "component": "vpc"
          }
        }
      }
    }
  }
}
```

<br/>

:::tip PRO TIP

If the filtering options built-in to Atmos are not sufficient, redirect the output to [`jq`](https://stedolan.github.io/jq/) for very powerful
filtering options.

:::
