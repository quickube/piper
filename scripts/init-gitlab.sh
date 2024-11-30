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
  helm upgrade --install gitlab -n gitlab gitlab/gitlab -f gitlab.values.yaml

  timeout 5m bash -c '
  while true; do
    # Get the pod status using kubectl
    POD_STATUS=$(kubectl get pods -n gitlab -l app=webservice -o jsonpath='{.items[0].status.phase}')

    # Check if the pod is in 'Running' state
    if [ "${POD_STATUS}" = "Running" ]; then
      echo "Pod is ready, continuing with next steps..."
      sleep 10
      break  # Exit the loop
    else
      echo "Pod is not ready yet. Current status: $POD_STATUS"
      sleep 5  # Wait for 5 seconds before checking again
    fi
  done
  '

  GITLAB_TOOLBOX_POD=$(kubectl get pods --namespace gitlab -l app=toolbox -o name | sed 's|pod/||')
  TOKENS_OUTPUT=$(kubectl exec -it -c toolbox ${GITLAB_TOOLBOX_POD} -n gitlab -- gitlab-rails runner /tmp/scripts/piper-setup.rb)
  echo $TOKENS_OUTPUT
else
  echo "Gitlab release exists, skipping installation"
fi