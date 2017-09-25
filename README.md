# terraform-provider-gocd 0.1.8

[![GoDoc](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd?status.svg)](https://godoc.org/github.com/drewsonne/terraform-provider-gocd/gocd)
[![Build Status](https://travis-ci.org/drewsonne/terraform-provider-gocd.svg?branch=master)](https://travis-ci.org/drewsonne/terraform-provider-gocd)
[![codecov](https://codecov.io/gh/drewsonne/terraform-provider-gocd/branch/master/graph/badge.svg)](https://codecov.io/gh/drewsonne/terraform-provider-gocd)
[![Go Report Card](https://goreportcard.com/badge/github.com/drewsonne/terraform-provider-gocd)](https://goreportcard.com/report/github.com/drewsonne/terraform-provider-gocd)

Terraform provider for GoCD Server

## Installation

    $ brew tap drewsonne/tap
    $ brew install terraform-provider-gocd
    $ tf-install-provider gocd
    
__NOTE__: `terraform` does not currently provide a way to easily install 3rd party providers. Until this is implemented,
the `tf-install-provider` utility can be used to copy the provider binary to the correct location.

## Resources

 - `gocd_pipeline`
 - `gocd_pipeline_template`
 - `gocd_pipeline_stage`
 - `gocd_environment`
 - `gocd_environment_association`

### gocd\_pipeline

Provides support for creating pipelines in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline" "build" {
  name = "build"
  group = "terraform-provider-gocd"
  label_template = "0.0.$${COUNT}"
  materials = [
    {
      type = "git"
      attributes {
        name = "terraform-provider-gocd"
        url = "https://github.com/drewsonne/terraform-provider-gocd.git"
      }
    }
  ]
}
```

#### Argument Reference

 - `name` - (Required) The name of the pipeline.
 - `group` - (Required) The name of the pipeline group to deploy into.
 - `materials` - (Required) The list of materials to be used by pipeline. At least one material must be specified. Each `materials` block supports fields documented below.
 - `label_template` - (Optional)  The label template to customise the pipeline instance label. 
 - `enable_pipeline_locking` - (Optional)  Whether pipeline is locked to run single instance or not.
 - `template` - (Optional)  The name of the template used by pipeline. A `gocd_pipeline_stage` can not be assigned to a `gocd_pipeline` it `template` is set.
 - `parameters` - (Optional) A [map](https://www.terraform.io/docs/configuration/variables.html#maps) of parameters to be used for substitution in a pipeline or pipeline template.
 - `environment_variables` - (Optional) The list of environment variables that will be passed to all tasks (commands) that are part of this pipeline. Each `environment_variables` block supports fields documented below.
 
The `environment_variables` block supports:

 - `name` - (Required) The name of the environment variable.
 - `value` - (Optional) The value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `encrypted_value` - (Optional) The encrypted value of the environment variable. One of `value` or `encrypted_value` must be set.
 - `secure` - Whether environment variable is secure or not. When set to `true`, encrypts the value if one is specified. The default value is `false`.

Type `materials` block supports:

 - `type` (Required) The type of a material. Can be one of git, dependency.
 - `attributes` (Required) A [map](https://www.terraform.io/docs/configuration/variables.html#maps) of attributes for each material type. See the [GoCD API Documentation](https://api.gocd.org/current/#the-pipeline-material-object) for each material type attributes.
   

#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.
 
### gocd\_pipeline\_template

Provides support for creating pipeline templates in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline_template" "terraform-builder" {
  name = "terraform-build-template"
}
```

#### Argument Reference

 - `name` - (Required) The name of the pipeline template.

#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

### gocd\_pipeline\_stage

Provides support for creating stages for pipelines or pipeline templates in GoCD.

#### Example Usage

```hcl
resource "gocd_pipeline_stage" "build" {
  name = "plan"
  pipeline = "plan"
  jobs = [
  <<JOB
 {
  "name": "plan",
  "tasks": [{
    "type": "exec",
    "attributes": {
      "run_if": ["passed"],
      "command": "terraform",
      "arguments": ["plan"]
    }
  }]
 }
  JOB
  ]
}
```

### gocd\_environment

Provides support for creating environmnets in GoCD.

#### Example Usage

```hcl
resource "gocd_environment" "testing" {
  name = "testing"
}
```

#### Argument Rference

 - `name` - (Required) Name of the environment to create.
 
#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.
 
### gocd\_environment\_association

Provides support for associating pipelines and environments in GoCD.

__NOTE:__ There is an intention to support agents and environment variables in the future.

#### Example Usage

```hcl
resource "gocd_environment_association" "build-in-testing" {
  environment = "${gocd_environment.testing.name}"
  pipeline = "${gocd_pipeline.build.name}"
}

resource "gocd_environment" "testing" {
  name = "testing"
}

resource "gocd_pipeline" "build" {
  name = "build"
  # ...
}
```

#### Argument Reference

 - `environment` - (Required) The name of the environment which the resource is being associated to.
 - `pipeline` - (Required) The name of the pipeline to associate to the environment
 

#### Attributes Reference

 - `version` - The current version of the resource configuration in GoCD.

## Demo

You will need docker and terraform >= 0.10.0 installed for this demo to work.

Either build the provider with `go build` or download it from the gihub repository. If you download it, make sure the binary is in the current folder.

	$ go build

Spin up the test gocd server, with endpoint at http://127.0.0.1:8153/go/

    $ make provision-test-gocd && sh ./scripts/wait-for-test-server.sh

Then initialise and apply the configuration.

    $ terraform init
    $ terraform apply

When you're finished, run:

    $ make teardown-test-gocd