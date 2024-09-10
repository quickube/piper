## Installation

You must deploy Piper to a cluster with a pre-existing Argo Workflows deployment.
Piper will create a CRD that Argo Workflows will pick, so install or configure Piper to create those CRDs in the right namespace.

Please check out [values.yaml](https://github.com/quickube/piper/tree/main/helm-chart/values.yaml) file of the helm chart configurations.

To add piper helm repo run:

```bash
helm repo add piper https://piper.quickube.com
```

After configuring Piper [values.yaml](https://github.com/quickube/piper/tree/main/helm-chart/values.yaml), run the following command for installation:

```bash
helm upgrade --install piper piper/piper \
-f YOUR_VALUES_FILE.yaml
```

---

## Required Configuration

### Ingress

Piper works best when it is able to listen to webhooks from your git provider.
Expose Piper using an ingress or service, then provide the address to `piper.webhook.url` as follows:
`https://PIPER_EXPOSED_URL/webhook`

Refer to [values.yaml](https://github.com/quickube/piper/tree/main/helm-chart/values.yaml) for more information.

### Git

Piper will use git to fetch the `.workflows` folder and receive events using webhooks.

To pick which git provider you are using provide `gitProvider.name` configuration in helm chart (Currently we only support GitHub and Bitbucket).

You must also configure your organization (GitHub) or workspace (Bitbucket) name using `gitProvider.organization.name` in the helm chart.

#### Git Token Permissions

The token should have access for creating webhooks and read repositories content.
For GitHub, configure `admin:org` and `write:org` permissions in Classic Token.
For Bitbucket, configure `Repositories:read`, `Webhooks:read and write` and `Pull requests:read` permissions (for multiple repos use workspace token).

#### Token

The git token should be passed as secret in the helm chart at `gitProvider.token`.
The token can be passed as parameter via helm install command using `--set piper.gitProvider.token=YOUR_GIT_TOKEN`

Alternatively, you can consume an already existing secret by configuring `piper.gipProvider.existingSecret`.
The key should have the name `token`. You can be create a Secret using this command:

```bash
kubectl create secret generic piper-git-token --from-literal=token=YOUR_GIT_OKEN
```

#### Webhook creation

Piper will create a webhook configuration for you, for the whole organization or for each repo you configure.

Configure `piper.webhook.url` with the address of Piper that you exposed using an Ingress or Service with `/webhook` postfix.

For organization level configuration: `gitProvider.webhook.orgLevel` to `true`.

For granular repo webhook provide list of repos at: `gitProvider.webhook.repoList`.

Piper implements graceful shutdown, it will delete all the webhooks when terminated.

#### Status check

Piper will handle status checks for you.
It will notify the GitProvider for the status of Workflow for specific commit that triggered Piper.
For linking provide valid URL of your Argo Workflows server address at: `argoWorkflows.server.address`

---

### Argo Workflow Server (On development)

Piper will use REST API to communicate with Argo Workflows server for linting or for creation of workflows (ARGO_WORKFLOWS_CREATE_CRD). Please follow this [configuration](https://argoproj.github.io/argo-workflows/rest-api/).

To lint the workflow before submitting it, please configure the internal address of Argo Workflows server (for example, `argo-server.workflows.svc.cluster.local`) in the field: `argoWorkflows.server.address`. Argo will need a [token](https://argoproj.github.io/argo-workflows/access-token/) to authenticate. please provide the secret in `argoWorkflows.server.token`, Better to pass as a references to a secret in the field `argoWorkflows.server.token`.

#### Skip CRD Creation (On development)

Piper can communicate directly to Argo Workflow using ARGO_WORKFLOWS_CREATE_CRD environment variable, if you want to skip the creation of CRD change `argoWorkflows.crdCreation` to `false`.
