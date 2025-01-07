## Health Check

Currently not supported for GitLab / Bitbucket

The following example shows a health check being executed every 1 minute as configured in the helm chart under `livenessProbe`, and triggered by the `/healthz` endpoint:

```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
    scheme: HTTP
  initialDelaySeconds: 10
  timeoutSeconds: 10
  periodSeconds: 60
  successThreshold: 1
  failureThreshold: 4
```

The mechanism for checking the health of Piper is:

1. Piper sets the health status of all webhooks to not healthy

2. Piper requests a ping from all the configured webhooks.

3. The Git provider sends a ping to the `/webhook` endpoint, which will set the health status to `healthy` with a timeout of 5 seconds.

4. Piper checks the status of all configured webhooks.

Therefore, the criteria for health checking are:

1. The registered webhook exists.
2. The webhook sends a ping within 5 seconds.
