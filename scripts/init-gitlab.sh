#!/bin/sh
set -o errexit


LICENSE="${1:-$GITLAB_LICENSE}"

if [ -z "$LICENSE" ]; then
  echo "no gitlab license was entered for init-gitlab.sh as argument or env: GITLAB_LICENSE"
  exit 2
fi

if [ -z "$(helm list -n gitlab | grep gitlab)" ]; then
  #start gitlab namespace
  kubectl create namespace gitlab
  # add gitlab secret
  kubectl create secret generic gitlab-license -n gitlab --from-literal=license_key=$1
  kubectl apply -f ./scripts/gitlab-setup.yaml -n gitlab
  # 8. Install gitlab
  helm repo add gitlab https://charts.gitlab.io/
  helm upgrade --install gitlab -n gitlab gitlab/gitlab --version 8.6.1 -f gitlab.values.yaml

  echo "waiting for gitlab toolbox pod to ready"
  kubectl wait --namespace gitlab --for=condition=ready pod -l app=toolbox --timeout=600s
  echo "waiting for gitlab webservice pod to ready"
  kubectl wait --namespace gitlab --for=condition=ready pod -l app=webservice --timeout=600s

  echo "setup gitlab configs"
  GITLAB_TOOLBOX_POD=$(kubectl get pods --namespace gitlab -l app=toolbox -o name)

  sleep 10
  TOKENS_OUTPUT=$(kubectl exec -it -c toolbox ${GITLAB_TOOLBOX_POD} -n gitlab -- gitlab-rails runner /tmp/scripts/piper-setup.rb)
  echo $TOKENS_OUTPUT
else
  echo "Gitlab release exists, skipping installation"
fi