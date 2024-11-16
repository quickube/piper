#!/bin/bash

GITLAB_WEBSERVICE_POD=$(kubectl get pods --namespace gitlab -l app=webservice -o name | sed 's|pod/||')

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

GITLAB_TOOLBOX_POD=$(kubectl get pods --namespace gitlab -l app=toolbox -o name | sed 's|pod/||')
kubectl cp ./scripts/gitlab-setup.rb gitlab/${GITLAB_TOOLBOX_POD}:/tmp 

GITLAB_ROOT_TOKEN=$(kubectl exec -it -c toolbox ${GITLAB_TOOLBOX_POD} -n gitlab -- gitlab-rails runner /tmp/gitlab-setup.rb | grep -oP '(?<=TOKEN: )\S+')
echo ${GITLAB_ROOT_TOKEN}

CONTENT_JSON_HEADER="Content-Type: application/json"
GITLAB_TOKEN_HEADER="PRIVATE-TOKEN: ${GITLAB_ROOT_TOKEN}" 
GITLAB_BASE_URL="http://localhost:8080/api/v4"


# create a new user
user_id=$(curl --location "${GITLAB_BASE_URL}/users" --header "${CONTENT_JSON_HEADER}" --header "${GITLAB_TOKEN_HEADER}" \
--data '{
    "email": "piper@example.com",
    "password": "Aa123456",
    "password_confirmation": "Aa123456",
    "username":"piper-user",
    "name":"piper-user"
}' | jq -r ".id")
sleep 3

IMPERSONATE_USER_HEADER="Sudo: ${user_id}"

# create a new group as the user created
group_id=$(curl --location "${GITLAB_BASE_URL}/groups" --header "${CONTENT_JSON_HEADER}" --header "${GITLAB_TOKEN_HEADER}" \
--header "${IMPERSONATE_USER_HEADER}"  --data '{"name": "Pied Pipers", "path": "pied-pipers"}' | jq -r ".id")
sleep 3

# create a project on group namespace
reponse=$(curl --location "${GITLAB_BASE_URL}/projects" --header "${CONTENT_JSON_HEADER}" --header "${GITLAB_TOKEN_HEADER}" \
--header "${IMPERSONATE_USER_HEADER}" --data "{\"name\":\"piper-e2e-test\", \"namespace_id\": \"${group_id}\"}")
sleep 3

# create a user personal access token
USER_TOKEN=$(curl --location "${GITLAB_BASE_URL}/users/${user_id}/personal_access_tokens" --header "${CONTENT_JSON_HEADER}" --header "${GITLAB_TOKEN_HEADER}" \
--data '{"name":"p-token", "scopes": ["api", "read_api", "write_repository"]}' | jq -r ".token")
sleep 3


EXPIRES_AT_ONE_YEAR=$(date -d "+1 year" +"%Y-%m-%d")
# create group access token
GROUP_TOKEN=$(curl --location "${GITLAB_BASE_URL}/groups/${group_id}/access_tokens" --header "${CONTENT_JSON_HEADER}" --header "${GITLAB_TOKEN_HEADER}" \
--header "${IMPERSONATE_USER_HEADER}" --data "{\"name\":\"g-token\", \"scopes\": [\"api\", \"read_api\", \"write_repository\"], \"expires_at\":\"${EXPIRES_AT_ONE_YEAR}\", \"access_level\": 40 }" | jq -r ".token")
sleep 3

echo $GROUP_TOKEN, $GITLAB_ROOT_TOKEN



















