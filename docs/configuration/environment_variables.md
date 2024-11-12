## Environment Variables

Piper uses the following environment variables to configure its functionality.
The helm chart populates them using [values.yaml](https://github.com/quickube/piper/tree/main/helm-chart/values.yaml) file

### Git

- GIT_PROVIDER
  The git provider that Piper will use, possible variables: gitHub, gitlab or bitbucket

- GIT_TOKEN
  The git token that will be used.

- GIT_URL
  the git url that will be used, only relevant when running gitlab self hosted

- GIT_ORG_NAME
  The organization name.

* GIT_ORG_LEVEL_WEBHOOK
  Boolean variable, whether to config webhook at the organization level. Defaults to `false`.

- GIT_WEBHOOK_REPO_LIST
  List of repositories to configure webhooks to.

* GIT_WEBHOOK_URL
  URL of Piper ingress, to configure webhooks.

* GIT_WEBHOOK_AUTO_CLEANUP
  Boolean variable that, if true, will cause Piper to automatically cleanup all webhooks that it creates when they are no longer necessary.
  Notice that there is a race condition between a pod that is being terminated and the new one being scheduled.

* GIT_ENFORCE_ORG_BELONGING
  Boolean variable that, if true, will cause Piper to enforce organizational belonging of git event creator. Defaults to `false`.

* GIT_FULL_HEALTH_CHECK
  Boolean variable that, if true, enables full health checks on webhooks. A full health check contains expecting and validating ping event from a webhook.
  Doesn't work for Bitbucket, because the API call doesn't exist on that platform.

### Argo Workflows Server

* ARGO_WORKFLOWS_TOKEN
  This token is used to authenticate with the Argo Workflows server.

- ARGO_WORKFLOWS_ADDRESS
  The address of Argo Workflows Server.

* ARGO_WORKFLOWS_CREATE_CRD
  Boolean variable that deterines whether to directly send Workflows instructions or create a CRD in the Cluster.

- ARGO_WORKFLOWS_NAMESPACE
  The namespace of Workflows creation for Argo Workflows.

- KUBE_CONFIG
  Used to configure Argo Workflows client with local kube configurations.

### Rookout

* ROOKOUT_TOKEN
  The token used to configure Rookout agent. If not provided, will not start the agent.
* ROOKOUT_LABELS
  The labels to label instances at Rookout, defaults to "service:piper"
* ROOKOUT_REMOTE_ORIGIN
  The repo URL for source code fetching, default:"https://github.com/quickube/piper.git".
