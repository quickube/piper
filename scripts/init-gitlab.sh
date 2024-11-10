#!/bin/sh
set -o errexit

if [ -z "$(helm list -n gitlab | grep gitlab)" ]; then
  # 8. Install gitlab
  helm repo add gitlab https://charts.gitlab.io/
  helm upgrade --install gitlab --create-namespace -n gitlab gitlab/gitlab -f gitlab.values.yaml
else
  echo "Gitlab release exists, skipping installation"
fi