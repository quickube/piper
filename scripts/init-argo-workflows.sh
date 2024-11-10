#!/bin/sh
set -o errexit

if [ -z "$(helm list -n workflows | grep argo-workflow)" ]; then
  # 7. Install argo workflows
  helm repo add argo https://argoproj.github.io/argo-helm
  helm upgrade --install argo-workflow argo/argo-workflows -n workflows --create-namespace -f workflows.values.yaml
else
  echo "Workflows release exists, skipping installation"
fi


cat <<EOF | kubectl apply -n cas -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argo-workflows-sa
automountServiceAccountToken: true
---
apiVersion: v1
kind: Secret
metadata:
  name: argo-workflows-sa.service-account-token
  annotations:
    kubernetes.io/service-account.name: argo-workflows-sa
type: kubernetes.io/service-account-token
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: argo-workflows-sa
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: argo-workflow-argo-workflows-workflow
subjects:
  - kind: ServiceAccount
    name: argo-workflows-sa
EOF