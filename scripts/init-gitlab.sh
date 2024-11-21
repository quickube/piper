#!/bin/sh
set -o errexit


LICENSE="${1:-$GITLAB_LICENSE}"
echo "$LICENSE"

if [ -z "$LICENSE" ]; then
  echo "no gitlab license was entered for init-gitlab.sh as argument or env: GITLAB_LICENSE"
  exit 2
fi

if [ -z "$(helm list -n gitlab | grep gitlab)" ]; then
  #start gitlab namespace
  kubectl create namespace gitlab
  # add gitlab secret
  kubectl create secret generic gitlab-license -n gitlab --from-literal=license_key=$1
  # 8. Install gitlab
  helm repo add gitlab https://charts.gitlab.io/
  helm upgrade --install gitlab -n gitlab gitlab/gitlab -f gitlab.values.yaml
else
  echo "Gitlab release exists, skipping installation"
fi