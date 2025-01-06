## .workflows Folder

Piper will look in each of the target branches for a `.workflows` folder. [Example](https://github.com/quickube/piper/tree/main/examples/.workflows).
We will explain each of the files that should be included in the `.workflows` folder:

### triggers.yaml (convention name)

This file holds a list of triggers that will be executed `onStart` by `events` from specific `branches`.
Piper will execute each of the matching triggers, so configure it wisely.

```yaml
- events:
    - push
    - pull_request.synchronize
  branches: ["main"]
  onStart: ["main.yaml"]
  onExit: ["exit.yaml"]
  templates: ["templates.yaml"]
  config: "default"
```

This example can be found [here](https://github.com/quickube/piper/tree/main/examples/.workflows/triggers.yaml).

In this example, `main.yaml` will be executed as a DAG when `push` or `pull_request.synchronize` events are applied in the `main` branch.
`onExit` will execute `exit.yaml` when the workflow finishes as an exit handler.

`onExit` can overwrite the default `onExit` configuration by referencing existing DAG tasks as in the [example](https://github.com/quickube/piper/tree/main/examples/.workflows/exit.yaml).

The `config` field is used for workflow configuration selection. The default value is the `default` configuration.

#### events

The `events` field is used to determine when the trigger will be executed. The name of the event depends on the git provider.

For instance, the GitHub `pull_request` event has a few actions, one of which is `synchronize`.

#### branches

The branch for which the trigger will be executed.

#### onStart

This [file](https://github.com/quickube/piper/tree/main/examples/.workflows/main.yaml) can be named as you wish and will be referenced in the `triggers.yaml` file. It will define an entrypoint DAG that the Workflow will execute.

As a best practice, this file should contain the dependency logic and parameterization of each referenced template. It should not implement new templates; for this, use the `template.yaml` file.

#### onExit

This field is used to pass a verbose exit handler to the triggered workflow.
It will override the default `onExit` from the provided `config` or the default `config`.

The provided `exit.yaml` describes a DAG that will overwrite the default `onExit` configuration.
[Example](https://github.com/quickube/piper/tree/main/examples/.workflows/exit.yaml)

#### templates

This field will have additional templates that will be injected into the workflows.
The purpose of this field is to create repository-scope templates that can be referenced from the DAG templates at `onStart` or `onExit`.
[Example](https://github.com/quickube/piper/tree/main/examples/.workflows/templates.yaml)

As a best practice, use this field for template implementation and reference them from the executed DAGs.
[Example](https://github.com/quickube/piper/tree/main/examples/.workflows/main.yaml).

### config

Configured by the `piper-workflows-config` [ConfigMap](workflows_config.md).
It can be passed explicitly, or it will use the `default` configuration.

### parameters.yaml (convention name)

It will hold a list of global parameters for the Workflow.
These can be referenced from any template with `{{ workflow.parameters.___ }}`.

[Example](https://github.com/quickube/piper/tree/main/examples/.workflows/parameters.yaml)