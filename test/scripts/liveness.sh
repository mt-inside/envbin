source common.sh

set -x
curl -w "%{http_code}\n" ${base_url}/health
curl ${base_url}/handlers/liveness?value=false
curl -w "%{http_code}\n" ${base_url}/health
curl ${base_url}/handlers/liveness?value=true
curl -w "%{http_code}\n" ${base_url}/health
