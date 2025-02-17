---
title: Component Validation
sidebar_position: 2
sidebar_label: Validation
description: Use JSON Schema and OPA policies to validate Components.
id: validation
---

Validation is critical to maintaining hygienic configurations in distributed team environments.

Atmos component validation allows:

* Validate component config (`vars`, `settings`, `backend`, `env`, `overrides` and other sections) using JSON Schema

* Check if the component config (including relations between different component variables) is correct to allow or deny component provisioning using
  OPA/Rego policies

:::tip

Refer to [atmos validate component](/cli/commands/validate/component) CLI command for more information

:::

## JSON Schema

Atmos has native support for [JSON Schema](https://json-schema.org/), which can validate the schema of configurations. JSON Schema is an industry
standard and provides a vocabulary to annotate and validate JSON documents for correctness.

This is powerful stuff: because you can define many schemas, it's possible to validate components differently for different environments or teams.

## Open Policy Agent (OPA)

The [Open Policy Agent](https://www.openpolicyagent.org/docs/latest/) (OPA, pronounced “oh-pa”) is another open-source industry standard that provides
a general-purpose policy engine to unify policy enforcement across your stacks. The OPA language (Rego) is a high-level declarative language for
specifying policy as code. Atmos has native support for the OPA decision-making engine to enforce policies across all the components in your stacks
(e.g. for microservice configurations).

## Usage

Atmos `validate component` command supports `--schema-path`, `--schema-type` and `--module-paths` command line arguments.
If the arguments are not provided, Atmos will try to find and use the `settings.validation` section defined in the component's YAML config.

```bash
atmos validate component vpc -s plat-ue2-prod --schema-path vpc/validate-vpc-component.json --schema-type jsonschema

atmos validate component vpc -s plat-ue2-prod --schema-path vpc/validate-vpc-component.rego --schema-type opa

atmos validate component vpc -s plat-ue2-dev --schema-path vpc/validate-vpc-component.rego --schema-type opa --module-paths catalog/constants

atmos validate component vpc -s plat-ue2-dev --schema-path vpc/validate-vpc-component.rego --schema-type opa --module-paths catalog

atmos validate component vpc -s plat-ue2-prod

atmos validate component vpc -s plat-ue2-dev

atmos validate component vpc -s plat-ue2-dev --timeout 15
```

### Configure Component Validation

In [`atmos.yaml`](https://github.com/cloudposse/atmos/blob/master/examples/quick-start/rootfs/usr/local/etc/atmos/atmos.yaml), add the `schemas`
section:

```yaml title="atmos.yaml"
# Validation schemas (for validating atmos stacks and components)
schemas:
  # https://json-schema.org
  jsonschema:
    # Can also be set using `ATMOS_SCHEMAS_JSONSCHEMA_BASE_PATH` ENV var, or `--schemas-jsonschema-dir` command-line arguments
    # Supports both absolute and relative paths
    base_path: "stacks/schemas/jsonschema"
  # https://www.openpolicyagent.org
  opa:
    # Can also be set using `ATMOS_SCHEMAS_OPA_BASE_PATH` ENV var, or `--schemas-opa-dir` command-line arguments
    # Supports both absolute and relative paths
    base_path: "stacks/schemas/opa"
```

In the component [manifest](https://github.com/cloudposse/atmos/blob/master/examples/quick-start/stacks/catalog/vpc/defaults.yaml), add
the `settings.validation` section:

```yaml title="examples/quick-start/stacks/catalog/vpc/defaults.yaml"
components:
  terraform:
    vpc:
      metadata:
        # Point to the Terraform component
        component: vpc
      settings:
        # Validation
        # Supports JSON Schema and OPA policies
        # All validation steps must succeed to allow the component to be provisioned
        validation:
          validate-vpc-component-with-jsonschema:
            schema_type: jsonschema
            # 'schema_path' can be an absolute path or a path relative to 'schemas.jsonschema.base_path' defined in `atmos.yaml`
            schema_path: "vpc/validate-vpc-component.json"
            description: Validate 'vpc' component variables using JSON Schema
          check-vpc-component-config-with-opa-policy:
            schema_type: opa
            # 'schema_path' can be an absolute path or a path relative to 'schemas.opa.base_path' defined in `atmos.yaml`
            schema_path: "vpc/validate-vpc-component.rego"
            # An array of filesystem paths (folders or individual files) to the additional modules for schema validation
            # Each path can be an absolute path or a path relative to `schemas.opa.base_path` defined in `atmos.yaml`
            # In this example, we have the additional Rego modules in `stacks/schemas/opa/catalog/constants`
            module_paths:
              - "catalog/constants"
            description: Check 'vpc' component configuration using OPA policy
            # Set `disabled` to `true` to skip the validation step
            # `disabled` is set to `false` by default, the step is allowed if `disabled` is not declared
            disabled: false
            # Validation timeout in seconds
            timeout: 10
```

Add the following JSON Schema in the
file [`stacks/schemas/jsonschema/vpc/validate-vpc-component.json`](https://github.com/cloudposse/atmos/blob/master/examples/quick-start/stacks/schemas/jsonschema/vpc/validate-vpc-component.json):

```json title="examples/quick-start/stacks/schemas/jsonschema/vpc/validate-vpc-component.json"
{
  "$id": "vpc-component",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "vpc component validation",
  "description": "JSON Schema for the 'vpc' Atmos component.",
  "type": "object",
  "properties": {
    "vars": {
      "type": "object",
      "properties": {
        "region": {
          "type": "string"
        },
        "cidr_block": {
          "type": "string",
          "pattern": "^([0-9]{1,3}\\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$"
        },
        "map_public_ip_on_launch": {
          "type": "boolean"
        }
      },
      "additionalProperties": true,
      "required": [
        "region",
        "cidr_block",
        "map_public_ip_on_launch"
      ]
    }
  }
}
```

Add the following Rego package in the file [`stacks/schemas/opa/catalog/constants/constants.rego`](https://github.com/cloudposse/atmos/blob/master/examples/quick-start/stacks/schemas/opa/catalog/constants/constants.rego):

```rego title="examples/quick-start/stacks/schemas/opa/catalog/constants/constants.rego"
package atmos.constants

vpc_dev_max_availability_zones_error_message := "In 'dev', only 2 Availability Zones are allowed"

vpc_prod_map_public_ip_on_launch_error_message := "Mapping public IPs on launch is not allowed in 'prod'. Set 'map_public_ip_on_launch' variable to 'false'"

vpc_name_regex := "^[a-zA-Z0-9]{2,20}$"

vpc_name_regex_error_message := "VPC name must be a valid string from 2 to 20 alphanumeric chars"
```

Add the following OPA policy in the file [`stacks/schemas/opa/vpc/validate-vpc-component.rego`](https://github.com/cloudposse/atmos/blob/master/examples/quick-start/stacks/schemas/opa/vpc/validate-vpc-component.rego):

```rego title="examples/quick-start/stacks/schemas/opa/vpc/validate-vpc-component.rego"
# Atmos looks for the 'errors' (array of strings) output from all OPA policies
# If the 'errors' output contains one or more error messages, Atmos considers the policy failed

# 'package atmos' is required in all `atmos` OPA policies
package atmos

import future.keywords.in

# Import the constants from the file `stacks/schemas/opa/catalog/constants/constants.rego`
import data.atmos.constants.vpc_dev_max_availability_zones_error_message
import data.atmos.constants.vpc_prod_map_public_ip_on_launch_error_message
import data.atmos.constants.vpc_name_regex
import data.atmos.constants.vpc_name_regex_error_message

# In production, don't allow mapping public IPs on launch
errors[vpc_prod_map_public_ip_on_launch_error_message] {
    input.vars.stage == "prod"
    input.vars.map_public_ip_on_launch == true
}

# In 'dev', only 2 Availability Zones are allowed
errors[vpc_dev_max_availability_zones_error_message] {
    input.vars.stage == "dev"
    count(input.vars.availability_zones) != 2
}

# Check VPC name
errors[vpc_name_regex_error_message] {
    not re_match(vpc_name_regex, input.vars.name)
}
```

<br/>

:::note

Atmos supports OPA policies for components validation in a single Rego file and in multiple Rego files.

As shown in the example above, you can define some Rego constants, modules and helper functions in a separate
file `stacks/schemas/opa/catalog/constants/constants.rego`, and then import them into the main policy
file `stacks/schemas/opa/vpc/validate-vpc-component.rego`.

You also need to specify the `module_paths` attribute in the component's `settings.validation` section.
The `module_paths` attribute is an array of filesystem paths (folders or individual files) to the additional modules for schema validation.
Each path can be an absolute path or a path relative to `schemas.opa.base_path` defined in `atmos.yaml`.
If a folder is specified in `module_paths`, Atmos will recursively process the folder and all its sub-folders and load all Rego files into the OPA
engine.

This allows you to separate the common OPA modules, constants and helper functions into a catalog of reusable Rego modules,
and to structure your OPA policies to make them DRY.

:::

<br/>

Run the following commands to validate the component in the stacks:

```bash
> atmos validate component vpc -s plat-ue2-prod

Mapping public IPs on launch is not allowed in 'prod'. Set 'map_public_ip_on_launch' variable to 'false'

exit status 1
```

```bash
> atmos validate component vpc -s plat-ue2-dev

In 'dev', only 2 Availability Zones are allowed
VPC name must be a valid string from 2 to 20 alphanumeric chars

exit status 1
```

```bash
> atmos validate component vpc -s plat-ue2-staging

Validate 'vpc' component variables using JSON Schema
{
  "valid": false,
  "errors": [
    {
      "keywordLocation": "",
      "absoluteKeywordLocation": "file:///examples/quick-start/stacks/schemas/jsonschema/vpc-component#",
      "instanceLocation": "",
      "error": "doesn't validate with file:///examples/quick-start/stacks/schemas/jsonschema/vpc-component#"
    },
    {
      "keywordLocation": "/properties/vars/properties/cidr_block/pattern",
      "absoluteKeywordLocation": "file:///examples/quick-start/stacks/schemas/jsonschema/vpc-component#/properties/vars/properties/cidr_block/pattern",
      "instanceLocation": "/vars/cidr_block",
      "error": "does not match pattern '^([0-9]{1,3}\\\\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$'"
    }
  ]
}

exit status 1
```

Try to run the following commands to provision the component in the stacks:

```bash
atmos terraform apply vpc -s plat-ue2-prod
atmos terraform apply vpc -s plat-ue2-dev
```

Since the OPA validation policies don't pass, Atmos does not allow provisioning the component in the stacks:

![atmos-validate-vpc-in-plat-ue2-prod](/img/atmos-validate-infra-vpc-in-tenant1-ue2-dev.png)
![atmos-validate-vpc-in-plat-ue2-dev](/img/atmos-validate-infra-vpc-in-tenant1-ue2-dev.png)

## OPA Policy Examples

```rego
# 'atmos' looks for the 'errors' (array of strings) output from all OPA policies
# If the 'errors' output contains one or more error messages, 'atmos' considers the policy failed

# 'package atmos' is required in all 'atmos' OPA policies
package atmos

import future.keywords.in

# Import the constants
import data.atmos.constants.vpc_dev_max_availability_zones_error_message
import data.atmos.constants.vpc_prod_map_public_ip_on_launch_error_message
import data.atmos.constants.vpc_name_regex
import data.atmos.constants.vpc_name_regex_error_message

# Function `object_has_key` checks if an object has the specified key with a string value
# https://www.openpolicyagent.org/docs/latest/policy-reference/#types
object_has_key(o, k) {
    some item
    item = o[k]
    type_name(item) == "string"
}

# In production, don't allow mapping public IPs on launch
errors[vpc_prod_map_public_ip_on_launch_error_message] {
    input.vars.stage == "prod"
    input.vars.map_public_ip_on_launch == true
}

# In 'dev', only 2 Availability Zones are allowed
errors[vpc_dev_max_availability_zones_error_message] {
    input.vars.stage == "dev"
    count(input.vars.availability_zones) != 2
}

# Check VPC name
errors[vpc_name_regex_error_message] {
    not re_match(vpc_name_regex, input.vars.name)
}

# Check the app hostname usign Regex
errors[message] {
    not re_match("^([a-z0-9]+([\\-a-z0-9]*[a-z0-9]+)?\\.){1,}([a-z0-9]+([\\-a-z0-9]*[a-z0-9]+)?){1,63}(\\.[a-z0-9]{2,7})+$", input.vars.app_config.hostname)
    message = "'app_config.hostname' must contain at least a subdomain and a top level domain. Example: subDomain1.topLevelDomain.com"
}

# Check the email address usign Regex
errors[message] {
    not re_match("^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$", input.vars.app_config.contact.email)
    message = "'app_config.contact.email' must be a valid email address"
}

# Check the phone number usign Regex
errors[message] {
    not re_match("^[\\+]?[(]?[0-9]{3}[)]?[-\\s\\.]?[0-9]{3}[-\\s\\.]?[0-9]{4,6}", input.vars.app_config.contact.phone)
    message = "'app_config.contact.phone' must be a valid phone number"
}

# Check if the component has a `Team` tag
errors[message] {
    not object_has_key(input.vars.tags, "Team")
    message = "All components must have 'Team' tag defined to specify which team is responsible for managing and provisioning them"
}

# Check if the Team has permissions to provision components in an OU (tenant)
errors[message] {
    input.vars.tags.Team == "devs"
    input.vars.tenant == "corp"
    message = "'devs' team cannot provision components into 'corp' OU"
}

# Check the message of the day from the manager
# If `settings.notes.allowed` is set to `false`, output the message from the manager
errors[message] {
    input.settings.notes.allowed == false
    message = concat("", [input.settings.notes.manager, " says: ", input.settings.notes.message])
}

# Check `notes2` config in the free-form Atmos section `settings`
errors[message] {
    input.settings.notes2.message == ""
    message = "'notes2.message' should not be empty"
}

# Check that the `app_config.hostname` variable is defined only once for the stack accross all stack manifests
# Refer to https://atmos.tools/cli/commands/describe/component#sources-of-component-variables for details on how 
# 'atmos' detects sources for all variables
# https://www.openpolicyagent.org/docs/latest/policy-language/#universal-quantification-for-all
errors[message] {
    hostnames := {app_config | some app_config in input.sources.vars.app_config; app_config.hostname}
    count(hostnames) > 0
    message = "'app_config.hostname' variable must be defined only once for the stack accross all stack manifests"
}

# This policy checks that the 'bar' variable is not defined in any of the '_defaults.yaml' Atmos stack manifests
# Refer to https://atmos.tools/cli/commands/describe/component#sources-of-component-variables for details on how 
# 'atmos' detects sources for all variables
# https://www.openpolicyagent.org/docs/latest/policy-language/#universal-quantification-for-all
errors[message] {
    # Get all 'stack_dependencies' of the 'bar' variable
    stack_dependencies := input.sources.vars.bar.stack_dependencies
    # Get all stack dependencies of the 'bar' variable where 'stack_file' ends with '_defaults'
    defaults_stack_dependencies := {stack_dependency | some stack_dependency in stack_dependencies; endswith(stack_dependency.stack_file, "_defaults")}
    # Check the count of the stack dependencies of the 'bar' variable where 'stack_file' ends with '_defaults'
    count(defaults_stack_dependencies) > 0
    # Generate the error message
    message = "The 'bar' variable must not be defined in any of '_defaults.yaml' stack manifests"
}

# This policy checks that if the 'foo' variable is defined in the 'stack1.yaml' stack manifest, it cannot be overriden in 'stack2.yaml'
# Refer to https://atmos.tools/cli/commands/describe/component#sources-of-component-variables for details 
# on how 'atmos' detects sources for all variables
# https://www.openpolicyagent.org/docs/latest/policy-language/#universal-quantification-for-all
errors[message] {
    # Get all 'stack_dependencies' of the 'foo' variable
    stack_dependencies := input.sources.vars.foo.stack_dependencies
    # Check if the 'foo' variable is defined in the 'stack1.yaml' stack manifest
    stack1_dependency := endswith(stack_dependencies[0].stack_file, "stack1")
    stack1_dependency == true
    # Get all stack dependencies of the 'foo' variable where 'stack_file' ends with 'stack2' (this means that the variable 
    # is redefined in one of the files 'stack2')
    stack2_dependencies := {stack_dependency | some stack_dependency in stack_dependencies; endswith(stack_dependency.stack_file, "stack2")}
    # Check the count of the stack dependencies of the 'foo' variable where 'stack_file' ends with 'stack2'
    count(stack2_dependencies) > 0
    # Generate the error message
    message = "If the 'foo' variable is defined in 'stack1.yaml', it cannot be overriden in 'stack2.yaml'"
}

# This policy shows an example on how to check the imported files in the stacks
# All stack files (root stacks and imported) that the current component depends on are in the `deps` section
# For example:
# deps:
# - catalog/xxx
# - catalog/yyy
# - orgs/zzz/_defaults
errors[message] {
    input.vars.tags.Team == "devs"
    input.vars.tenant == "corp"
    input.deps[_] == "catalog/xxx"
    message = "'devs' team cannot import the 'catalog/xxx' file when provisioning components into 'corp' OU"
}

errors["'service_1_name' variable length must be greater than 10 chars"] {
    count(input.vars.service_1_name) <= 10
}
```

<br/>

:::note

- If a regex pattern in the 're_match' function contains a backslash to escape special chars (e.g. '\.' or '\-'),
  it must be escaped with another backslash when represented as a regular Go string ('\\.', '\\-').

- The reason is that backslash is also used to escape special characters in Go strings like newline (\n).

- If you want to match the backslash character itself, you'll need four slashes.

:::
