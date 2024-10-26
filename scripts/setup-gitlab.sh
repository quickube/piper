helm repo add gitlab https://charts.gitlab.io/ 
helm repo update
helm upgrade --install gitlab gitlab/gitlab \
--namespace gitlab --create-namespace \
--timeout 600s \
--set global.hosts.domain=pipelab.com \
--set global.hosts.externalIP=0.0.0.0 \
--set certmanager-issuer.email=omriassa@gmail.com \
--set gitlab-runner.install=false \
--set prometheus.install=false