source common.sh

set -x
curl -w "%{http_code}\n" ${base_url}/ready
curl ${base_url}/handlers/readiness?value=false
curl -w "%{http_code}\n" ${base_url}/ready
curl ${base_url}/handlers/readiness?value=true
curl -w "%{http_code}\n" ${base_url}/ready
