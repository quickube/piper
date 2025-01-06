## Environment Variables

Piper uses the following environment variables to configure its functionality.
The helm chart populates them using the [values.yaml](https://github.com/quickube/piper/tree/main/helm-chart/values.yaml) file.

### Git

- GIT_PROVIDER
  The git provider that Piper will use, possible variables: GitHub | GitLab | Bitbucket

* GIT_TOKEN
  The git token that will be used to connect to the git provider.

- GIT_URL
  The git URL that will be used, only relevant when running GitLab self-hosted.

- GIT_ORG_NAME
  The organization name.

* GIT_ORG_LEVEL_WEBHOOK
  Boolean variable, whether to configure the webhook at the organization level. Defaults to `false`.

* GIT_WEBHOOK_REPO_LIST
  List of repositories to configure webhooks for.

* GIT_WEBHOOK_URL
  URL of Piper ingress to configure webhooks.

* GIT_WEBHOOK_AUTO_CLEANUP
  Boolean variable that, if true, will cause Piper to automatically clean up all webhooks it creates when they are no longer necessary.
  Note that there is a race condition between a pod being terminated and a new one being scheduled.

* GIT_ENFORCE_ORG_BELONGING
  Boolean variable that, if true, will cause Piper to enforce the organizational belonging of the git event creator. Defaults to `false`.

* GIT_FULL_HEALTH_CHECK
  Boolean variable that, if true, enables full health checks on webhooks. A full health check involves expecting and validating a ping event from a webhook.
  This doesn't work for Bitbucket because the API call doesn't exist on that platform.

### Argo Workflows Server

* ARGO_WORKFLOWS_TOKEN
  This token is used to authenticate with the Argo Workflows server.

* ARGO_WORKFLOWS_ADDRESS
  The address of the Argo Workflows server.

* ARGO_WORKFLOWS_CREATE_CRD
  Boolean variable that determines whether to directly send Workflows instructions or create a CRD in the Cluster.

* ARGO_WORKFLOWS_NAMESPACE
  The namespace of Workflows creation for Argo Workflows.

* KUBE_CONFIG
  Used to configure the Argo Workflows client with local kube configurations.

### Rookout

* ROOKOUT_TOKEN
  The token used to configure the Rookout agent. If not provided, the agent will not start.

* ROOKOUT_LABELS
  The labels to label instances in Rookout, defaults to "service:piper".

* ROOKOUT_REMOTE_ORIGIN
  The repo URL for source code fetching, defaults to "https://github.com/quickube/piper.git".