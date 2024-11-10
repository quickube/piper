helm repo add gitlab https://charts.gitlab.io/ 
helm upgrade --install gitlab gitlab/gitlab -n gitlab --create-namespace -f gitlab.values.yaml

# set -o errexit

# if [ -z "$(helm list -n gitlab | grep gitlab)" ]; then
#   # 7. Install gitlab self hosted
#     helm repo add gitlab https://charts.gitlab.io/ 
#     helm upgrade --install gitlab gitlab/gitlab -n gitlab --create-namespace -f gitlab.values.yaml
# else
#     echo "gitlab release exists, skipping installation"
# fi

# kubectl get secret gitlab-gitlab-initial-root-password -n=gitlab -ojsonpath='{.data.password}' | base64 --decode; echo
# glab auth login --hostname gitlab.local --token=9Rkj6d1GMnBuMLXNjAhfoYoBISVYeYCok7dpMX17ctWy1taZyzk4tW4Meoi0CXVF

# http://gitlab-webservice-default.gitlab.svc.cluster.local:8080


apt install glab
GITLAB_ROOT_PASSWORD="$(kubectl get secret gitlab-gitlab-initial-root-password --namespace=gitlab -ojsonpath='{.data.password}' | base64 --decode)"
sed -i -e "s/GITLABTOKEN/${GITLAB_ROOT_PASSWORD}/g" $HOME/.config/glab-cli/config.yml
cp ./gitlab-stuff/glab-config.yaml $HOME/.config/glab-cli/config.yml