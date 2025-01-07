## Workflow Configuration

Piper can inject configuration for Workflows that Piper creates.

`default` config is used as a convention for all Workflows that Piper will create, even if not explicitly mentioned in the `triggers.yaml` file.

### ConfigMap

Piper will mount a ConfigMap when Helm is used.
The `piper.workflowsConfig` variable in the Helm chart will create a ConfigMap that holds a set of configurations for Piper.
Here is an [example](https://github.com/quickube/piper/tree/main/examples/config.yaml) of such a configuration.

### Spec

This will be injected into the Workflow spec field and can hold all configurations of the Workflow.
> :warning: Please note that the fields `entrypoint` and `onExit` should not exist in the spec; both of them are managed fields.

### onExit

This is the exit handler for each of the Workflows created by Piper.
It configures a DAG that will be executed when the workflow ends.
You can provide the templates to it as shown in the following [Examples](https://github.com/quickube/piper/tree/main/examples/config.yaml).